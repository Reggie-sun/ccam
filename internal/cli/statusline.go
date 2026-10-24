package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
)

func readAll(r *bufio.Reader) ([]byte, error) { return io.ReadAll(r) }

// runShell runs a wrapped status line command with the same input, in
// the shell Claude Code would use. It has no time limit of its own: Claude
// Code gives a status line 10 minutes and stops a run when a newer update
// starts, so a slow command shows as it would without ccam. When Claude
// Code ends ccam, the command (on Unix its whole process group) ends too.
func runShell(command string, input []byte) (string, error) {
	ctx, stop := signal.NotifyContext(context.Background(), endSignals...)
	defer stop()
	sh, flag := shell()
	c := exec.CommandContext(ctx, sh, flag, command)
	ownGroup(c)
	c.Stdin = bytes.NewReader(input)
	// A background child the command leaves holds the pipe open; without
	// WaitDelay, Output would wait for that child however long it runs.
	c.WaitDelay = 200 * time.Millisecond
	out, err := c.Output()
	if errors.Is(err, exec.ErrWaitDelay) {
		err = nil // the command itself ended well
	}
	return string(out), err
}

// shell is what Claude Code runs status line commands with: sh, or on
// Windows Git Bash (CLAUDE_CODE_GIT_BASH_PATH or the usual install
// folders), else PowerShell.
func shell() (string, string) {
	if runtime.GOOS != "windows" {
		return "sh", "-c"
	}
	candidates := []string{os.Getenv("CLAUDE_CODE_GIT_BASH_PATH"), `C:\Program Files\Git\bin\bash.exe`, `C:\Program Files (x86)\Git\bin\bash.exe`}
	if git, err := exec.LookPath("git"); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(filepath.Dir(git)), "bin", "bash.exe"))
	}
	for _, c := range candidates {
		if c != "" && fsutil.Exists(c) {
			return c, "-c"
		}
	}
	for _, c := range []string{`C:\Program Files\PowerShell\7\pwsh.exe`, `C:\Program Files\PowerShell\7-preview\pwsh.exe`} {
		if fsutil.Exists(c) {
			return c, "-Command"
		}
	}
	return "powershell", "-Command"
}

// statusLineCommand is how settings.json runs this ccam. On Windows the
// command goes through Git Bash or PowerShell, so the path uses forward
// slashes and ~ for the home folder, which both understand. A path that
// needs quotes even so is quoted for the shell Claude Code will use: in
// single quotes for Git Bash, after & for PowerShell, and without ~, which
// neither expands inside quotes. sh is that shell (shell).
func statusLineCommand(self, home, sh string) string {
	if runtime.GOOS != "windows" {
		return launcher.ShellQuote(self) + " statusline"
	}
	path, rest := filepath.ToSlash(self), filepath.ToSlash(self)
	if rel, err := filepath.Rel(home, self); err == nil && !strings.HasPrefix(rel, "..") {
		path, rest = "~/"+filepath.ToSlash(rel), filepath.ToSlash(rel)
	}
	if launcher.ShellQuote(rest) == rest && !strings.ContainsAny(rest, ",@%") {
		return path + " statusline"
	}
	abs := filepath.ToSlash(self)
	if strings.HasSuffix(strings.ToLower(sh), "bash.exe") {
		return launcher.ShellQuote(abs) + " statusline"
	}
	return "& " + launcher.PSString(abs) + " statusline"
}

func statuslineInstallCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Set the ccam status line in every account's settings.json (an existing one is wrapped and kept)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			self, err := ccamPath()
			if err != nil {
				return err
			}
			sh, _ := shell()
			done := map[string]bool{}
			for i := range cfg.Profiles {
				target := writeTarget(filepath.Join(profile.Dir(e.paths, &cfg.Profiles[i]), "settings.json"))
				if done[target] {
					continue
				}
				done[target] = true
				changed, err := installStatusLine(target, statusLineCommand(self, e.paths.Home, sh))
				if err != nil {
					return err
				}
				if changed {
					e.printf("✓ %s\n", target)
				} else {
					e.printf("· %s (already set)\n", target)
				}
			}
			return nil
		},
	}
}

// ccamStatusLine reports whether command runs `statusline` of a program
// named ccam, its path bare or quoted for sh, Git Bash or PowerShell (after
// &), and returns what follows "statusline": the options it wraps with.
func ccamStatusLine(command string) (string, bool) {
	s, ps := strings.CutPrefix(command, "& ")
	var prog string
	var found bool
	switch {
	case ps && strings.HasPrefix(s, "'"):
		prog, s, found = psWord(s)
	case s != "" && (s[0] == '\'' || s[0] == '"'):
		prog, s, found = shWord(s)
	default:
		prog, s, found = strings.Cut(s, " ")
		s = " " + s
	}
	if !found {
		return "", false
	}
	base := strings.ToLower(path.Base(filepath.ToSlash(prog)))
	rest, ok := strings.CutPrefix(s, " statusline")
	if strings.TrimSuffix(base, ".exe") != "ccam" || !ok || rest != "" && rest[0] != ' ' {
		return "", false
	}
	return rest, true
}

// shWord reads the sh word at the start of s as launcher.ShellQuote writes
// it, with a quote inside as an escaped one between two quoted parts:
// quoted parts and backslash escapes up to a space. It returns the word
// and the rest from that space.
func shWord(s string) (string, string, bool) {
	var w strings.Builder
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case ' ':
			return w.String(), s[i:], true
		case '\'', '"':
			end := strings.IndexByte(s[i+1:], c)
			if end < 0 {
				return "", "", false
			}
			w.WriteString(s[i+1 : i+1+end])
			i += 1 + end
		case '\\':
			if i+1 < len(s) {
				i++
				w.WriteByte(s[i])
			}
		default:
			w.WriteByte(c)
		}
	}
	return w.String(), "", true
}

// psWord reads the PowerShell single-quoted string at the start of s as
// launcher.PSString writes it: a quote written twice (also a typographic
// one) is one quote of the string. It returns the string and the rest
// after it.
func psWord(s string) (string, string, bool) {
	rs := []rune(s)
	quote := func(r rune) bool { return r == '\'' || r >= '\u2018' && r <= '\u201b' }
	var w strings.Builder
	for i := 1; i < len(rs); i++ {
		if !quote(rs[i]) {
			w.WriteRune(rs[i])
			continue
		}
		if i+1 < len(rs) && quote(rs[i+1]) {
			w.WriteRune(rs[i])
			i++
			continue
		}
		return w.String(), string(rs[i+1:]), true
	}
	return "", "", false
}

// writeTarget follows a settings.json symlink so a write updates the shared
// file instead of replacing the link with a regular file.
func writeTarget(path string) string {
	if target, err := os.Readlink(path); err == nil {
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(path), target)
		}
		return target
	}
	return path
}

// installStatusLine sets ours (the command running ccam statusline) in the
// settings file at path, wrapping a status line that is already there. One
// that runs a ccam's status line already (another path, or quoted as an
// older ccam wrote it) is replaced, keeping what it wraps.
func installStatusLine(path, ours string) (bool, error) {
	doc := map[string]any{}
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, &doc); err != nil {
			return false, err
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return false, err
	}
	command := ours
	// The existing object is edited, so its other options (padding,
	// refreshInterval, hideVimModeIndicator) stay.
	sl, ok := doc["statusLine"].(map[string]any)
	if ok {
		old, _ := sl["command"].(string)
		if strings.HasPrefix(old, ours) {
			return false, nil
		}
		rest, isOurs := ccamStatusLine(old)
		switch {
		case isOurs:
			command = ours + rest
		case old == "":
		case runtime.GOOS == "windows":
			command = ours + " --wrap-base64 " + base64.StdEncoding.EncodeToString([]byte(old))
		default:
			command = ours + " --wrap " + launcher.ShellQuote(old)
		}
	} else {
		sl = map[string]any{}
	}
	sl["type"], sl["command"] = "command", command
	doc["statusLine"] = sl
	out, err := fsutil.MarshalJSON(doc)
	if err != nil {
		return false, err
	}
	return true, fsutil.AtomicWrite(path, out, 0o600)
}
