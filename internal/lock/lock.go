// Package lock provides the two kinds of locks ccam needs: a non-blocking
// flock for ccam's own critical sections, and DirLock, which speaks Claude
// Code's `.claude.json.lock` convention (a directory created with mkdir,
// considered stale after 10 seconds).
package lock

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// Flock is an exclusive, non-blocking lock on a file: flock(2), or
// LockFileEx on Windows.
type Flock struct {
	f *os.File
}

// ErrBusy means another process holds the lock.
var ErrBusy = errors.New("already in use")

// TryFlock takes the lock without waiting.
func TryFlock(path string) (*Flock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, err
	}
	if err := tryLock(f); err != nil {
		f.Close()
		return nil, err
	}
	return &Flock{f: f}, nil
}

func (l *Flock) Unlock() {
	unlock(l.f) // closing releases it anyway
	l.f.Close()
}

// staleAfter matches proper-lockfile's default used by Claude Code.
const staleAfter = 10 * time.Second

// DirLock acquires a mkdir-style lock at path, waiting up to timeout.
// A lock directory older than 10 seconds is treated as stale and replaced.
func DirLock(path string, timeout time.Duration) (unlock func(), err error) {
	deadline := time.Now().Add(timeout)
	for {
		err := os.Mkdir(path, 0o700)
		if err == nil {
			return func() { os.Remove(path) }, nil
		}
		if !errors.Is(err, fs.ErrExist) {
			return nil, err
		}
		// Retry at once only when the stale lock is really gone; one that
		// cannot be removed (not empty, no permission) waits like a live one.
		if info, err := os.Stat(path); err == nil && time.Since(info.ModTime()) > staleAfter && os.Remove(path) == nil {
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("could not lock %s (Claude Code is writing the file)", path)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
