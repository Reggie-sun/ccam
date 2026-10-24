package lock

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
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

// WaitFlock gets the lock once the holder lets it go, and gives up with
// ErrBusy when it does not in time.
func TestWaitFlock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.lock")
	a, err := TryFlock(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := WaitFlock(path, 100*time.Millisecond); !errors.Is(err, ErrBusy) {
		t.Fatalf("held: %v", err)
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		a.Unlock()
	}()
	b, err := WaitFlock(path, 5*time.Second)
	if err != nil {
		t.Fatalf("after the holder let go: %v", err)
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

// Claude Code counts its token-refresh lock stale only after a minute: one
// 20 seconds old is still held for that lock, though not for a 10 s one.
func TestDirLockStaleAfterItsOwnAge(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".oauth_refresh.lock")
	os.Mkdir(path, 0o700)
	old := time.Now().Add(-20 * time.Second)
	os.Chtimes(path, old, old)
	if _, err := DirLockStale(path, 0, time.Minute); err == nil {
		t.Fatal("a 20 s old refresh lock was taken over")
	}
	unlock, err := DirLockStale(path, 0, 15*time.Second)
	if err != nil {
		t.Fatalf("stale for a 15 s lock: %v", err)
	}
	unlock()
	if _, err := os.Stat(path); err == nil {
		t.Fatal("unlock leaves the folder")
	}
}

// Several ccams that find the same stale lock at once take it one after
// another: making sure it is stale and removing it are two steps, and
// without taking turns a second taker removed the fresh lock of the first
// and both held it (two writers of one .claude.json).
func TestDirLockStaleTakeoverIsExclusive(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".claude.json.lock")
	overlaps := 0
	for range 40 {
		os.Mkdir(path, 0o700)
		old := time.Now().Add(-time.Minute)
		os.Chtimes(path, old, old)
		var holders, peak atomic.Int32
		start := make(chan struct{})
		var wg sync.WaitGroup
		for range 8 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				unlock, err := DirLock(path, 500*time.Millisecond)
				if err != nil {
					return
				}
				n := holders.Add(1)
				for p := peak.Load(); n > p && !peak.CompareAndSwap(p, n); p = peak.Load() {
				}
				time.Sleep(2 * time.Millisecond)
				holders.Add(-1)
				unlock()
			}()
		}
		close(start)
		wg.Wait()
		if peak.Load() > 1 {
			overlaps++
		}
		os.Remove(path)
	}
	if overlaps > 0 {
		t.Fatalf("%d of 40 times two held the lock at once", overlaps)
	}
}

// A holder slower than the stale age may lose its lock to Claude Code,
// which takes it over: its unlock then leaves the other's lock alone.
func TestDirLockUnlockLeavesAnotherHoldersLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".claude.json.lock")
	unlock, err := DirLock(path, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(path)
	os.Mkdir(path, 0o700)
	later := time.Now().Add(time.Hour)
	os.Chtimes(path, later, later) // even if the folder got the same inode
	unlock()
	if _, err := os.Stat(path); err != nil {
		t.Fatal("unlock removed the lock of another holder")
	}
}
