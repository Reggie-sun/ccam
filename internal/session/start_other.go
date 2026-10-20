//go:build !linux && !darwin && !windows

package session

// startTime is unsupported here; sessions are then not tracked.
func startTime(pid int) (string, bool) { return "", false }
