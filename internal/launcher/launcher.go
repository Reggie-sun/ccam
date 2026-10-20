// Package launcher turns "profile + claude args" into an exec of claude.
//
// This is the hot path behind `ccam <name>` and shortcuts: no cobra, no TUI,
// no HTTP, one config read, lstat-only checks.
package launcher

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
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

	sessions []string // session folders this process is registered in
}

// Prepare resolves the profile (sel == "" means the one set for the
// current folder, else the default one) and builds the exec plan. User args
// are passed through untouched, except that one leading "--" separator is
// dropped.
func Prepare(p paths.Paths, cfg *config.Config, sel string, args []string) (*Plan, error) {
	var prof *config.Profile
	var err error
	if sel == "" {
		prof, err = HereProfile(cfg)
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
		return nil, fmt.Errorf("the folder of account %s is missing (%s). Check with `ccam doctor`", prof.Name, dir)
	}
	argv := make([]string, 0, len(args)+1)
	argv = append(argv, "claude")
	argv = append(argv, args...)
	return &Plan{Profile: prof, Dir: dir, Bin: bin, Argv: argv, Env: ProfileEnv(p, os.Environ(), dir, prof.Name)}, nil
}

// HereProfile is the account a launch without a name uses: the one set for
// the current folder with `ccam folder`, else the default account.
func HereProfile(cfg *config.Config) (*config.Profile, error) {
	if len(cfg.Folders) > 0 {
		if f, ok := FolderHere(cfg); ok {
			if prof, ok := cfg.Find(f.Account); ok {
				return prof, nil
			}
			return nil, fmt.Errorf("%s is set to account %q, which does not exist. Change it with `ccam folder`", f.Path, f.Account)
		}
	}
	return cfg.DefaultProfile()
}

// FolderHere finds the `ccam folder` entry for the working directory. Entries
// are stored with symlinks resolved; the working directory is resolved only
// when it does not match as it is, which keeps the common case to one getcwd.
func FolderHere(cfg *config.Config) (*config.Folder, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, false
	}
	if f, ok := cfg.FolderFor(cwd); ok {
		return f, true
	}
	if real, err := filepath.EvalSymlinks(cwd); err == nil && real != cwd {
		return cfg.FolderFor(real)
	}
	return nil, false
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

// DryRun prints what would be executed, as a line for sh, or for
// PowerShell on Windows.
func (pl *Plan) DryRun(w io.Writer) {
	if runtime.GOOS == "windows" {
		// The value of an assignment must be a string, so it is always quoted.
		parts := []string{"$env:CLAUDE_CONFIG_DIR='" + strings.ReplaceAll(pl.Dir, "'", "''") + "';", "&", PSQuote(pl.Bin)}
		for _, a := range pl.Argv[1:] {
			parts = append(parts, PSQuote(a))
		}
		fmt.Fprintln(w, strings.Join(parts, " "))
		return
	}
	parts := []string{"CLAUDE_CONFIG_DIR=" + ShellQuote(pl.Dir), ShellQuote(pl.Bin)}
	for _, a := range pl.Argv[1:] {
		parts = append(parts, ShellQuote(a))
	}
	fmt.Fprintln(w, strings.Join(parts, " "))
}

// PSQuote quotes s for PowerShell when needed.
func PSQuote(s string) string {
	if s != "" && strings.IndexFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-_./:\\", r))
	}) < 0 {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// BeforeExec runs right before exec, e.g. to flush a CPU profile.
var BeforeExec func()

// ExitCode is returned when claude ran and ended with a non-zero code that
// ccam must exit with, without a message of its own.
type ExitCode int

func (e ExitCode) Error() string { return fmt.Sprintf("claude exited with code %d", int(e)) }

// Run starts claude as a child process with the given standard streams,
// waits for it and returns its exit code (128+n when signal n ended it).
// Ctrl+C is not passed on: the terminal sends it to claude as well.
func (pl *Plan) Run(p paths.Paths, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	setTitle(pl.Profile.Name)
	MarkUsed(p, pl.Profile.Name)
	cmd := &exec.Cmd{Path: pl.Bin, Args: pl.Argv, Env: pl.Env, Stdin: stdin, Stdout: stdout, Stderr: stderr}
	sig := make(chan os.Signal, 4)
	signal.Notify(sig, passOn...)
	defer signal.Stop(sig)
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	done := make(chan struct{})
	go func() {
		for {
			select {
			case s := <-sig:
				if s != os.Interrupt {
					_ = cmd.Process.Signal(s)
				}
			case <-done:
				return
			}
		}
	}()
	err := cmd.Wait()
	close(done)
	var ee *exec.ExitError
	switch {
	case err == nil:
		return 0, nil
	case errors.As(err, &ee):
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
			return 128 + int(ws.Signal()), nil
		}
		return ee.ExitCode(), nil
	}
	return 0, err
}

func setTitle(name string) {
	if runtime.GOOS == "windows" {
		return // the console shows escape sequences as text unless VT output is on
	}
	if fi, err := os.Stderr.Stat(); err == nil && fi.Mode()&os.ModeCharDevice != 0 {
		fmt.Fprintf(os.Stderr, "\x1b]0;claude · %s\a", name)
	}
}

func lastUsedPath(p paths.Paths, name string) string {
	return filepath.Join(p.State, "last-used", name)
}

func mcpStateDir(p paths.Paths, name string) string {
	return filepath.Join(p.State, "mcp", name)
}

// StatePaths lists what ccam keeps under an account's name in the state
// folder, so rename can carry it along and remove can clear it.
func StatePaths(p paths.Paths, name string) []string {
	return []string{lastUsedPath(p, name), mcpStateDir(p, name), filepath.Join(p.State, "picks", name)}
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

// ResolveClaude finds the claude executable, which must not be ccam itself
// (that would add ccam's flags on every pass). Lookups through PATH are
// cached in the state dir with the file's mtime and the claude_bin they
// resolved; a changed or missing file, or another claude_bin, triggers a
// fresh lookup.
func ResolveClaude(p paths.Paths, bin string) (string, error) {
	if bin == "" {
		bin = "claude"
	}
	if strings.ContainsRune(bin, '/') || strings.ContainsRune(bin, filepath.Separator) {
		if isSelf(bin) {
			return "", fmt.Errorf("%s points to ccam itself. Check claude_bin in config.toml", bin)
		}
		return bin, nil
	}
	cache := filepath.Join(p.State, "claude-path")
	if data, err := os.ReadFile(cache); err == nil {
		path, rest, _ := strings.Cut(string(data), "\t")
		mtime, cachedBin, _ := strings.Cut(rest, "\t")
		if cachedBin == bin {
			if fi, err := os.Stat(path); err == nil && strconv.FormatInt(fi.ModTime().UnixNano(), 10) == mtime {
				return path, nil
			}
		}
	}
	path, err := exec.LookPath(bin)
	if err != nil {
		return "", fmt.Errorf("cannot find the claude executable (%s). Is Claude Code installed?", bin)
	}
	path, _ = filepath.Abs(path)
	if isSelf(path) {
		return "", fmt.Errorf("%s points to ccam itself. Check claude_bin in config.toml", path)
	}
	if fi, err := os.Stat(path); err == nil {
		_ = fsutil.WriteFast(cache, []byte(path+"\t"+strconv.FormatInt(fi.ModTime().UnixNano(), 10)+"\t"+bin))
	}
	return path, nil
}

// isSelf reports whether path is the running ccam binary, through any
// symlink or hard link. Comparing inodes costs two stats, which matters on
// the launch path; resolving both paths component by component does not.
func isSelf(path string) bool {
	self, err := os.Executable()
	if err != nil {
		return false
	}
	a, err1 := os.Stat(self)
	b, err2 := os.Stat(path)
	return err1 == nil && err2 == nil && os.SameFile(a, b)
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
	return fmt.Errorf("stopped because of share link problems:\n%s\nRepair them with `ccam doctor --fix`, then run again", strings.Join(msg, "\n"))
}
