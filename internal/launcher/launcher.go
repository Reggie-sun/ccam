// Package launcher turns "profile + claude args" into an exec of claude.
//
// This is the hot path behind `ccs <name>` and `claude2`: no cobra, no TUI,
// no HTTP, one config read, lstat-only checks.
package launcher

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/trace"
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
	argv := make([]string, 0, len(args)+1)
	argv = append(argv, "claude")
	argv = append(argv, args...)
	return &Plan{Profile: prof, Dir: dir, Bin: bin, Argv: argv, Env: buildEnv(os.Environ(), dir, prof.Name)}, nil
}

// InsertArgs puts ccam-generated args (e.g. --mcp-config) before the user's,
// so user flags of the same kind still apply on top.
func (pl *Plan) InsertArgs(extra ...string) {
	argv := make([]string, 0, len(pl.Argv)+len(extra))
	argv = append(argv, pl.Argv[0])
	argv = append(argv, extra...)
	pl.Argv = append(argv, pl.Argv[1:]...)
}

func buildEnv(environ []string, dir, name string) []string {
	env := make([]string, 0, len(environ)+2)
	for _, kv := range environ {
		if strings.HasPrefix(kv, "CLAUDE_CONFIG_DIR=") || strings.HasPrefix(kv, "CCAM_PROFILE=") {
			continue
		}
		env = append(env, kv)
	}
	return append(env, "CLAUDE_CONFIG_DIR="+dir, "CCAM_PROFILE="+name)
}

// DryRun prints what would be executed.
func (pl *Plan) DryRun(w io.Writer) {
	parts := []string{"CLAUDE_CONFIG_DIR=" + shellQuote(pl.Dir), shellQuote(pl.Bin)}
	for _, a := range pl.Argv[1:] {
		parts = append(parts, shellQuote(a))
	}
	fmt.Fprintln(w, strings.Join(parts, " "))
}

// Exec replaces the current process with claude. It only returns on error.
func (pl *Plan) Exec(p paths.Paths) error {
	setTitle(pl.Profile.Name)
	MarkUsed(p, pl.Profile.Name)
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

func shellQuote(s string) string {
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

// Launch is the whole `ccs <name>` sequence: load config, prepare, then exec
// (or print the plan when dry).
func Launch(p paths.Paths, sel string, args []string, dry bool, stdout io.Writer) error {
	trace.Step("start")
	cfg, err := config.Load(p.ConfigFile())
	if err != nil {
		return err
	}
	trace.Step("load config")
	pl, err := Prepare(p, cfg, sel, args)
	if err != nil {
		return err
	}
	trace.Step("prepare")
	if dry {
		pl.DryRun(stdout)
		return nil
	}
	return pl.Exec(p)
}
