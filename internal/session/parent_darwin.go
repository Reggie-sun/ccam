package session

import (
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// Parent asks the kernel (kern.proc.pid) for the parent of process pid.
func Parent(pid int) (int, bool) {
	if pid <= 0 {
		return 0, false
	}
	kp, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
	if err != nil || kp.Proc.P_pid != int32(pid) {
		return 0, false
	}
	return int(kp.Eproc.Ppid), true
}

// Exe is the program process pid runs, links resolved: kern.procargs2 is
// argc (4 bytes), then the path it was started from.
func Exe(pid int) (string, bool) {
	if pid <= 0 {
		return "", false
	}
	buf, err := unix.SysctlRaw("kern.procargs2", pid)
	if err != nil || len(buf) < 5 {
		return "", false
	}
	path, _, _ := strings.Cut(string(buf[4:]), "\x00")
	if path == "" {
		return "", false
	}
	if real, err := filepath.EvalSymlinks(path); err == nil {
		path = real
	}
	return path, true
}
