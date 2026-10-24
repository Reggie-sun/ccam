//go:build !linux && !darwin && !windows

package session

// Parent is unknown here.
func Parent(pid int) (int, bool) { return 0, false }

// Exe is unknown here.
func Exe(pid int) (string, bool) { return "", false }
