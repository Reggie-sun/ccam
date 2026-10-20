package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
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
// the shell Claude Code would use.
func runShell(command string, input []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	sh, flag := shell()
	c := exec.CommandContext(ctx, sh, flag, command)
	c.Stdin = bytes.NewReader(input)
	// Killing sh does not close the pipe a background child still holds;
	// without WaitDelay, Output waits for that child however long it runs.
	c.WaitDelay = 200 * time.Millisecond
	out, err := c.Output()
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
	return "powershell", "-Command"
}

// statusLineCommand is how settings.json runs this ccam. On Windows the
// command goes through Git Bash or PowerShell, so the path uses forward
// slashes and ~ for the home folder, which both understand.
func statusLineCommand(self, home string) string {
	if runtime.GOOS != "windows" {
		return launcher.ShellQuote(self) + " statusline"
	}
	path := filepath.ToSlash(self)
	if rel, err := filepath.Rel(home, self); err == nil && !strings.HasPrefix(rel, "..") {
		path = "~/" + filepath.ToSlash(rel)
	}
	if strings.ContainsAny(path, " '\"$") {
		path = `"` + path + `"`
	}
	return path + " statusline"
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
			done := map[string]bool{}
			for i := range cfg.Profiles {
				target := writeTarget(filepath.Join(profile.Dir(e.paths, &cfg.Profiles[i]), "settings.json"))
				if done[target] {
					continue
				}
				done[target] = true
				changed, err := installStatusLine(target, statusLineCommand(self, e.paths.Home))
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
// settings file at path, wrapping a status line that is already there.
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
	if sl, ok := doc["statusLine"].(map[string]any); ok {
		old, _ := sl["command"].(string)
		if strings.HasPrefix(old, ours) {
			return false, nil
		}
		if old != "" {
			command = ours + " --wrap " + launcher.ShellQuote(old)
		}
	}
	doc["statusLine"] = map[string]any{"type": "command", "command": command}
	out, err := fsutil.MarshalJSON(doc)
	if err != nil {
		return false, err
	}
	return true, fsutil.AtomicWrite(path, out, 0o600)
}
