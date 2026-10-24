package session

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// StartTime asks the kernel (kern.proc.pid) when the process started.
func StartTime(pid int) (string, bool) {
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

// sameStart compares a start time Claude Code recorded (procStart) with
// StartTime's. Claude Code records what `ps -o lstart=` prints in the C
// locale and UTC ("Thu Oct 22 03:28:37 2026"), to the second. known is
// false for a form ccam cannot read.
func sameStart(recorded, now string) (same, known bool) {
	at, err := time.Parse("Mon Jan 2 15:04:05 2006", strings.Join(strings.Fields(recorded), " "))
	start, ok := wallStart(now)
	if err != nil || !ok {
		return false, false
	}
	d := at.Unix() - start.Unix()
	return -1 <= d && d <= 1, true
}

// wallStart is StartTime's seconds and microseconds as a time.
func wallStart(now string) (time.Time, bool) {
	sec, usec, _ := strings.Cut(now, ".")
	s, err := strconv.ParseInt(sec, 10, 64)
	u, err2 := strconv.ParseInt(usec, 10, 64)
	if err != nil || err2 != nil {
		return time.Time{}, false
	}
	return time.Unix(s, u*1000), true
}
