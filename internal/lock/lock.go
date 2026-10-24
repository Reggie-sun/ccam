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

// WaitFlock takes the lock, waiting up to timeout for another process to
// let it go; after that it gives up with an error that is ErrBusy.
func WaitFlock(path string, timeout time.Duration) (*Flock, error) {
	deadline := time.Now().Add(timeout)
	for {
		l, err := TryFlock(path)
		if !errors.Is(err, ErrBusy) {
			return l, err
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("another ccam holds %s: %w", path, ErrBusy)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// staleAfter matches proper-lockfile's default used by Claude Code.
const staleAfter = 10 * time.Second

// heldError is DirLock's answer when another process still holds the lock
// after the wait; it is ErrBusy.
type heldError struct{ path string }

func (e heldError) Error() string {
	return fmt.Sprintf("could not lock %s (Claude Code is writing the file)", e.path)
}

func (heldError) Is(target error) bool { return target == ErrBusy }

// DirLock acquires a mkdir-style lock at path, waiting up to timeout.
// A lock directory older than 10 seconds is treated as stale and replaced.
// Making sure a lock is stale and removing it are two steps: two processes
// could both remove it, the second one the lock the first just made, and
// both would hold it. So ccam's processes take turns on a flock at
// path+".ccam" first. Claude Code's own lockfile has the same race, which
// ccam cannot rule out; the callers leave an account that runs alone.
func DirLock(path string, timeout time.Duration) (unlock func(), err error) {
	deadline := time.Now().Add(timeout)
	turn, err := waitTurn(path, deadline)
	if err != nil {
		return nil, err
	}
	release, err := DirLockStale(path, time.Until(deadline), staleAfter)
	if err != nil {
		turn.Unlock()
		return nil, err
	}
	return func() { release(); turn.Unlock() }, nil
}

// waitTurn takes the flock at lock+".ccam", waiting until deadline. Unlike
// TryFlock it does not make the folder: a lock in a folder that is gone
// fails, as the mkdir would.
func waitTurn(lock string, deadline time.Time) (*Flock, error) {
	f, err := os.OpenFile(lock+".ccam", os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, err
	}
	for {
		err := tryLock(f)
		if err == nil {
			return &Flock{f: f}, nil
		}
		if !errors.Is(err, ErrBusy) || time.Now().After(deadline) {
			f.Close()
			if errors.Is(err, ErrBusy) {
				return nil, fmt.Errorf("could not lock %s (another ccam is writing the file): %w", lock, ErrBusy)
			}
			return nil, err
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// DirLockStale is DirLock for a lock its holders count as stale after
// stale (proper-lockfile's option, which Claude Code sets per lock),
// without DirLock's turns among ccam's processes: its callers take a lock
// of their own around it.
func DirLockStale(path string, timeout, stale time.Duration) (unlock func(), err error) {
	deadline := time.Now().Add(timeout)
	for {
		err := os.Mkdir(path, 0o700)
		if err == nil {
			return own(path), nil
		}
		if !errors.Is(err, fs.ErrExist) {
			return nil, err
		}
		// Retry at once only when the stale lock is really gone; one that
		// cannot be removed (not empty, no permission) waits like a live one.
		if info, err := os.Stat(path); err == nil && time.Since(info.ModTime()) > stale && os.Remove(path) == nil {
			continue
		}
		if time.Now().After(deadline) {
			return nil, heldError{path}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// own is the unlock of the lock folder at path this process just made. A
// holder slower than the stale age may have lost the lock to another,
// whose folder then stays: only the very folder made here is removed.
func own(path string) func() {
	made, err := os.Lstat(path)
	return func() {
		if err == nil {
			now, err := os.Lstat(path)
			if err != nil || !os.SameFile(made, now) || !now.ModTime().Equal(made.ModTime()) {
				return
			}
		}
		os.Remove(path)
	}
}
