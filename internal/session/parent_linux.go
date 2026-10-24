package session

import (
	"os"
	"strconv"
)

// Parent is the parent of process pid: field 4 of /proc/<pid>/stat.
func Parent(pid int) (int, bool) {
	fields, ok := stat(pid)
	if !ok || len(fields) < 2 {
		return 0, false
	}
	ppid, err := strconv.Atoi(fields[1])
	return ppid, err == nil
}

// Exe is the program process pid runs, links resolved (/proc/<pid>/exe).
func Exe(pid int) (string, bool) {
	path, err := os.Readlink("/proc/" + strconv.Itoa(pid) + "/exe")
	return path, err == nil
}
