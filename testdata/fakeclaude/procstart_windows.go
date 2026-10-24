package main

import (
	"strconv"
	"syscall"
)

// procStart is this process's start as Claude Code records it on Windows:
// the creation time as a FILETIME, 100 ns since 1601.
func procStart() string {
	h, err := syscall.GetCurrentProcess()
	if err != nil {
		return ""
	}
	var created, exited, kernel, user syscall.Filetime
	if syscall.GetProcessTimes(h, &created, &exited, &kernel, &user) != nil {
		return ""
	}
	return strconv.FormatUint(uint64(created.HighDateTime)<<32|uint64(created.LowDateTime), 10)
}
