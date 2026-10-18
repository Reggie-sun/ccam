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
	"strings"
)

// Register marks the current process as a session in dir.
func Register(dir string) error {
	pid := os.Getpid()
	start, ok := startTime(pid)
	if !ok {
		return nil // no /proc: skip tracking rather than fail the launch
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, strconv.Itoa(pid)), []byte(start), 0o600)
}

// Active reports whether any session registered in dir is still running.
// Markers of finished sessions are removed on the way.
func Active(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	active := false
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		path := filepath.Join(dir, e.Name())
		want, _ := os.ReadFile(path)
		if got, ok := startTime(pid); ok && got == string(want) {
			active = true
			continue
		}
		os.Remove(path)
	}
	return active
}

// startTime reads field 22 of /proc/<pid>/stat (clock ticks since boot).
// The command name in field 2 may contain spaces, so parse after its ')'.
func startTime(pid int) (string, bool) {
	data, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return "", false
	}
	s := string(data)
	i := strings.LastIndexByte(s, ')')
	if i < 0 {
		return "", false
	}
	fields := strings.Fields(s[i+1:])
	if len(fields) < 20 {
		return "", false
	}
	return fields[19], true
}
