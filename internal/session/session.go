// Package session tracks Claude Code sessions started by ccam.
//
// Before exec, ccam writes a marker named after its PID holding the
// process start time. exec keeps both, so the marker identifies the claude
// process itself: children it spawns do not inherit it (unlike a lock fd),
// and a reused PID has a different start time.
package session

import (
	"os"
	"path/filepath"
	"strconv"
)

// Register marks the current process as a session in dir.
func Register(dir string) error {
	pid := os.Getpid()
	start, ok := startTime(pid)
	if !ok {
		return nil // unsupported platform: skip tracking rather than fail the launch
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	// Write under a name Active skips, then rename: a concurrent Active
	// must never read a half-written marker and delete it as stale.
	name := strconv.Itoa(pid)
	tmp := filepath.Join(dir, ".tmp-"+name)
	if err := os.WriteFile(tmp, []byte(start), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, name))
}

// Active reports whether any session registered in dir is still running.
// Markers of finished sessions are removed on the way.
func Active(dir string) bool { return Count(dir) > 0 }

// Count is the number of sessions registered in dir that still run.
func Count(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		path := filepath.Join(dir, e.Name())
		want, _ := os.ReadFile(path)
		if got, ok := startTime(pid); ok && got == string(want) {
			n++
			continue
		}
		os.Remove(path)
	}
	return n
}

// Unregister removes the marker of the current process from dir, for a
// session that ended while ccam kept running.
func Unregister(dir string) {
	os.Remove(filepath.Join(dir, strconv.Itoa(os.Getpid())))
}
