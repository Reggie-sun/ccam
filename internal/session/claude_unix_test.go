//go:build unix

package session

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"
)

// A record that is no regular file is not read: opening a FIFO would wait.
func TestClaudeRunsSkipsAFIFO(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "sessions"), 0o700)
	if err := syscall.Mkfifo(filepath.Join(dir, "sessions", strconv.Itoa(os.Getpid())+".json"), 0o600); err != nil {
		t.Skip(err)
	}
	done := make(chan bool, 1)
	go func() { done <- ClaudeRuns(dir) }()
	select {
	case got := <-done:
		if got {
			t.Fatal("a FIFO counted as a session")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("ClaudeRuns waits on a FIFO")
	}
}
