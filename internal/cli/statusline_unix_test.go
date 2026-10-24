//go:build !windows

package cli

import (
	"os"
	"path/filepath"
	"runtime"
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
	if groupRuns(pid) {
		syscall.Kill(-pid, syscall.SIGKILL)
		t.Fatal("the command's process group outlived ccam")
	}
}

// groupRuns reports whether a process of process group pgid still runs.
// One that has ended but that nobody has reaped yet (a zombie: in a
// container whose PID 1 reaps nothing, an orphan stays one) does not: on
// Linux /proc/<pid>/stat tells them apart, elsewhere any member counts.
func groupRuns(pgid int) bool {
	if runtime.GOOS != "linux" {
		return syscall.Kill(-pgid, 0) == nil
	}
	entries, _ := os.ReadDir("/proc")
	for _, e := range entries {
		if _, err := strconv.Atoi(e.Name()); err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join("/proc", e.Name(), "stat"))
		s := string(data)
		i := strings.LastIndexByte(s, ')')
		if err != nil || i < 0 {
			continue // ended meanwhile
		}
		// State, parent and process group follow the command's name.
		if f := strings.Fields(s[i+1:]); len(f) > 2 && f[2] == strconv.Itoa(pgid) && f[0] != "Z" && f[0] != "X" {
			return true
		}
	}
	return false
}
