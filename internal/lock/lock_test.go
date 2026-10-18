package lock

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestFlockExclusive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.lock")
	a, err := TryFlock(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := TryFlock(path); err != ErrBusy {
		t.Fatalf("second lock: %v", err)
	}
	if !Held(path) {
		t.Fatal("Held should be true")
	}
	a.Unlock()
	if Held(path) {
		t.Fatal("Held after unlock")
	}
}

// TestLockFollowsOpenFile shows the flock lives as long as any process keeps
// the file open. The syscall.Exec hand-over itself is covered in e2e.
func TestLockFollowsOpenFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.lock")
	l, err := TryFlock(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.KeepAcrossExec(); err != nil {
		t.Fatal(err)
	}
	child := exec.Command("sleep", "0.5")
	child.ExtraFiles = []*os.File{l.f}
	if err := child.Start(); err != nil {
		t.Fatal(err)
	}
	l.f.Close() // parent lets go; the child still holds the open file description
	if !Held(path) {
		t.Fatal("lock should be held by the child")
	}
	child.Wait()
	if Held(path) {
		t.Fatal("lock should be released when the child exits")
	}
}

func TestDirLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".claude.json.lock")
	unlock, err := DirLock(path, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DirLock(path, 100*time.Millisecond); err == nil {
		t.Fatal("second DirLock should time out")
	}
	unlock()
	// Stale locks are taken over.
	os.Mkdir(path, 0o700)
	old := time.Now().Add(-time.Minute)
	os.Chtimes(path, old, old)
	unlock, err = DirLock(path, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("stale lock not taken: %v", err)
	}
	unlock()
}
