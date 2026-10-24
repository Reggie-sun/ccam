//go:build !windows

package cli

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Claude Code ends a status line run when a newer update starts: ending
// ccam ends the wrapped command too, and what it started, instead of
// leaving them running.
func TestWrappedStatusLineEndsWithCcam(t *testing.T) {
	pidFile := filepath.Join(t.TempDir(), "pid")
	go func() {
		for range 100 {
			if data, err := os.ReadFile(pidFile); err == nil && strings.HasSuffix(string(data), "\n") {
				syscall.Kill(os.Getpid(), syscall.SIGTERM)
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()
	start := time.Now()
	runShell("echo $$ > "+pidFile+"; sleep 30; echo late", nil)
	if time.Since(start) > 10*time.Second {
		t.Fatal("the command was not ended")
	}
	data, _ := os.ReadFile(pidFile)
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if err := syscall.Kill(-pid, 0); err == nil {
		syscall.Kill(-pid, syscall.SIGKILL)
		t.Fatal("the command's process group outlived ccam")
	}
}
