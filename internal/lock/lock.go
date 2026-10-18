// Package lock provides the two kinds of locks ccam needs.
//
// Flock locks can be handed over to the exec'd claude process, so they stay
// held exactly as long as that Claude Code session runs. DirLock speaks
// Claude Code's own `.claude.json.lock` convention (a directory created
// with mkdir, considered stale after 10 seconds).
package lock

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// Flock is an exclusive, non-blocking flock on a file.
type Flock struct {
	f *os.File
}

// ErrBusy means another process holds the lock.
var ErrBusy = errors.New("이미 사용 중입니다")

// TryFlock takes the lock without waiting.
func TryFlock(path string) (*Flock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, ErrBusy
		}
		return nil, err
	}
	return &Flock{f: f}, nil
}

// Held reports whether someone currently holds the lock at path.
func Held(path string) bool {
	l, err := TryFlock(path)
	if err != nil {
		return errors.Is(err, ErrBusy)
	}
	l.Unlock()
	return false
}

// KeepAcrossExec clears close-on-exec so the lock survives syscall.Exec and
// is released only when the new program (claude) exits.
func (l *Flock) KeepAcrossExec() error {
	_, _, errno := syscall.Syscall(syscall.SYS_FCNTL, l.f.Fd(), syscall.F_SETFD, 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func (l *Flock) Unlock() {
	_ = syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN) // closing releases it anyway
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
		if info, err := os.Stat(path); err == nil && time.Since(info.ModTime()) > staleAfter {
			os.Remove(path)
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("%s 잠금을 얻지 못했습니다 (Claude Code가 파일을 쓰는 중)", path)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
