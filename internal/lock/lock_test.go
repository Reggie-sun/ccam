package lock

import (
	"os"
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
	a.Unlock()
	b, err := TryFlock(path)
	if err != nil {
		t.Fatalf("after unlock: %v", err)
	}
	b.Unlock()
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

// A stale lock that cannot be removed (not empty) must not spin past the
// timeout.
func TestDirLockGivesUpOnUnremovableStaleLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.lock")
	os.Mkdir(path, 0o700)
	os.WriteFile(filepath.Join(path, "keep"), nil, 0o600)
	old := time.Now().Add(-time.Minute)
	os.Chtimes(path, old, old)
	done := make(chan error, 1)
	go func() {
		_, err := DirLock(path, 100*time.Millisecond)
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("locked although the lock folder is still there")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("DirLock ignored its timeout")
	}
}
