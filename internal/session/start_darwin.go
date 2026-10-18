package session

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// startTime asks the kernel (kern.proc.pid) when the process started.
func startTime(pid int) (string, bool) {
	if pid <= 0 {
		return "", false
	}
	kp, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
	if err != nil || kp.Proc.P_pid != int32(pid) {
		return "", false
	}
	t := kp.Proc.P_starttime
	return fmt.Sprintf("%d.%06d", t.Sec, t.Usec), true
}
