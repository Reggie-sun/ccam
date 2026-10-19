// Package launcher turns "profile + claude args" into an exec of claude.
//
// This is the hot path behind `ccam <name>` and `claude2`: no cobra, no TUI,
// no HTTP, one config read, lstat-only checks.
package launcher

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
)

type Plan struct {
	Profile *config.Profile
	Dir     string   // CLAUDE_CONFIG_DIR
	Bin     string   // absolute path of claude
	Argv    []string // argv including argv[0]
	Env     []string
}

// Prepare resolves the profile (sel == "" means the default one) and builds
// the exec plan. User args are passed through untouched, except that one
// leading "--" separator is dropped.
func Prepare(p paths.Paths, cfg *config.Config, sel string, args []string) (*Plan, error) {
	var prof *config.Profile
	var err error
	if sel == "" {
		prof, err = cfg.DefaultProfile()
	} else {
		prof, err = cfg.MustFind(sel)
	}
	if err != nil {
		return nil, err
	}
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	bin, err := ResolveClaude(p, cfg.ClaudeBin)
	if err != nil {
		return nil, err
	}
	dir := profile.Dir(p, prof)
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("%s 계정 폴더가 없습니다 (%s). `ccam doctor` 로 확인하세요", prof.Name, dir)
	}
	argv := make([]string, 0, len(args)+1)
	argv = append(argv, "claude")
	argv = append(argv, args...)
	return &Plan{Profile: prof, Dir: dir, Bin: bin, Argv: argv, Env: ProfileEnv(p, os.Environ(), dir, prof.Name)}, nil
}

// InsertArgs puts ccam-generated args (e.g. --mcp-config) before the user's,
// so user flags of the same kind still apply on top.
func (pl *Plan) InsertArgs(extra ...string) {
	argv := make([]string, 0, len(pl.Argv)+len(extra))
	argv = append(argv, pl.Argv[0])
	argv = append(argv, extra...)
	pl.Argv = append(argv, pl.Argv[1:]...)
}

// ProfileEnv is the environment claude runs with for a profile.
func ProfileEnv(p paths.Paths, environ []string, dir, name string) []string {
	return buildEnv(environ, dir, name, runtime.GOOS == "darwin" && creds.DefaultStore(dir, p.Home))
}

// buildEnv sets CLAUDE_CONFIG_DIR. defaultStore keeps an imported ~/.claude
// on the unsuffixed macOS Keychain item that plain `claude` uses.
func buildEnv(environ []string, dir, name string, defaultStore bool) []string {
	env := make([]string, 0, len(environ)+3)
	for _, kv := range environ {
		if strings.HasPrefix(kv, "CLAUDE_CONFIG_DIR=") || strings.HasPrefix(kv, "CCAM_PROFILE=") ||
			strings.HasPrefix(kv, "CLAUDE_SECURESTORAGE_CONFIG_DIR=") {
			continue
		}
		env = append(env, kv)
	}
	env = append(env, "CLAUDE_CONFIG_DIR="+dir, "CCAM_PROFILE="+name)
	if defaultStore {
		env = append(env, "CLAUDE_SECURESTORAGE_CONFIG_DIR=")
	}
	return env
}

// DryRun prints what would be executed.
func (pl *Plan) DryRun(w io.Writer) {
	parts := []string{"CLAUDE_CONFIG_DIR=" + ShellQuote(pl.Dir), ShellQuote(pl.Bin)}
	for _, a := range pl.Argv[1:] {
		parts = append(parts, ShellQuote(a))
	}
	fmt.Fprintln(w, strings.Join(parts, " "))
}

// BeforeExec runs right before exec, e.g. to flush a CPU profile.
var BeforeExec func()

// Exec replaces the current process with claude. It only returns on error.
func (pl *Plan) Exec(p paths.Paths) error {
	setTitle(pl.Profile.Name)
	MarkUsed(p, pl.Profile.Name)
	if BeforeExec != nil {
		BeforeExec()
	}
	return syscall.Exec(pl.Bin, pl.Argv, pl.Env)
}

func setTitle(name string) {
	if fi, err := os.Stderr.Stat(); err == nil && fi.Mode()&os.ModeCharDevice != 0 {
		fmt.Fprintf(os.Stderr, "\x1b]0;claude · %s\a", name)
	}
}

func lastUsedPath(p paths.Paths, name string) string {
	return filepath.Join(p.State, "last-used", name)
}

// MarkUsed records the launch time; losing it is harmless, so no fsync.
func MarkUsed(p paths.Paths, name string) {
	_ = fsutil.WriteFast(lastUsedPath(p, name), []byte(strconv.FormatInt(time.Now().Unix(), 10)))
}

// LastUsed returns the last launch time, or zero.
func LastUsed(p paths.Paths, name string) time.Time {
	data, err := os.ReadFile(lastUsedPath(p, name))
	if err != nil {
		return time.Time{}
	}
	sec, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(sec, 0)
}

// ResolveClaude finds the claude executable. Lookups through PATH are cached
// in the state dir together with the file's mtime; a changed or missing file
// triggers a fresh lookup.
func ResolveClaude(p paths.Paths, bin string) (string, error) {
	if bin == "" {
		bin = "claude"
	}
	if strings.Contains(bin, "/") {
		return bin, nil
	}
	cache := filepath.Join(p.State, "claude-path")
	if data, err := os.ReadFile(cache); err == nil {
		path, mtime, ok := strings.Cut(strings.TrimSpace(string(data)), "\t")
		if ok {
			if fi, err := os.Stat(path); err == nil && strconv.FormatInt(fi.ModTime().UnixNano(), 10) == mtime {
				return path, nil
			}
		}
	}
	path, err := exec.LookPath(bin)
	if err != nil {
		return "", fmt.Errorf("claude 실행 파일을 찾을 수 없습니다 (%s). Claude Code를 설치했는지 확인하세요", bin)
	}
	path, _ = filepath.Abs(path)
	if isSelf(path) {
		return "", fmt.Errorf("%s 가 ccam 자신을 가리킵니다. config.toml 의 claude_bin 을 확인하세요", path)
	}
	if fi, err := os.Stat(path); err == nil {
		_ = fsutil.WriteFast(cache, []byte(path+"\t"+strconv.FormatInt(fi.ModTime().UnixNano(), 10)))
	}
	return path, nil
}

func isSelf(path string) bool {
	self, err := os.Executable()
	if err != nil {
		return false
	}
	a, err1 := filepath.EvalSymlinks(self)
	b, err2 := filepath.EvalSymlinks(path)
	return err1 == nil && err2 == nil && a == b
}

// ShellQuote quotes s for sh when needed.
func ShellQuote(s string) string {
	if s == "" {
		return "''"
	}
	if strings.IndexFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-_./=:@,+%", r))
	}) < 0 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// checkLinks verifies share links with lstat only. Missing links are created
// on the spot; anything that could lose data is left for `ccam doctor --fix`.
func checkLinks(p paths.Paths, pl *Plan) error {
	problems := share.Check(pl.Dir, pl.Profile.Share, p.GroupsRoot())
	if len(problems) == 0 {
		return nil
	}
	rest, err := share.QuickFix(problems)
	if err != nil {
		return err
	}
	if len(rest) == 0 {
		return nil
	}
	msg := make([]string, 0, len(rest))
	for _, pr := range rest {
		msg = append(msg, "  "+pr.String())
	}
	return fmt.Errorf("공유 링크에 문제가 있어 실행을 멈췄습니다:\n%s\n`ccam doctor --fix` 로 고친 뒤 다시 실행하세요", strings.Join(msg, "\n"))
}
