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
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
)

func readAll(r *bufio.Reader) ([]byte, error) { return io.ReadAll(r) }

// runShell runs a wrapped status line command with the same input.
func runShell(command string, input []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, "sh", "-c", command)
	c.Stdin = bytes.NewReader(input)
	out, err := c.Output()
	return string(out), err
}

func statuslineInstallCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "모든 계정의 settings.json 에 ccam 상태줄을 설정 (기존 상태줄은 감싸서 유지)",
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
				changed, err := installStatusLine(target, self)
				if err != nil {
					return err
				}
				if changed {
					e.printf("✓ %s\n", target)
				} else {
					e.printf("· %s (이미 설정됨)\n", target)
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

func installStatusLine(path, self string) (bool, error) {
	doc := map[string]any{}
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, &doc); err != nil {
			return false, err
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return false, err
	}
	ours := launcher.ShellQuote(self) + " statusline"
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
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return false, err
	}
	return true, fsutil.AtomicWrite(path, append(out, '\n'), 0o600)
}
