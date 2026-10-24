package session

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/sys/windows"
)

// StartTime asks Windows when the process started (GetProcessTimes). A
// process that has exited but is still referenced has an exit code other
// than STILL_ACTIVE and does not count.
func StartTime(pid int) (string, bool) {
	if pid <= 0 {
		return "", false
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return "", false
	}
	defer func() { _ = windows.CloseHandle(h) }()
	var code uint32
	if windows.GetExitCodeProcess(h, &code) != nil || code != 259 { // STILL_ACTIVE
		return "", false
	}
	var created, exited, kernel, user windows.Filetime
	if windows.GetProcessTimes(h, &created, &exited, &kernel, &user) != nil {
		return "", false
	}
	return fmt.Sprintf("%d", created.Nanoseconds()), true
}

// sameStart compares a start time Claude Code recorded (procStartFt, or
// procStart in sessions/) with StartTime's nanoseconds since 1970. Claude
// Code records the creation time as a FILETIME (100 ns since 1601); without
// its native helper it asks PowerShell for .NET ticks in local time
// instead, which ccam cannot compare: known is false for those and for any
// other form.
func sameStart(recorded, now string) (same, known bool) {
	ft, err := strconv.ParseInt(recorded, 10, 64)
	ns, err2 := strconv.ParseInt(now, 10, 64)
	if err != nil || err2 != nil || ft > 3e17 { // .NET ticks are above 6e17
		return false, false
	}
	d := ft - (ns/100 + 116444736000000000)
	return -1e7 <= d && d <= 1e7, true // within a second
}

// wallStart is StartTime's nanoseconds since 1970 as a time.
func wallStart(now string) (time.Time, bool) {
	ns, err := strconv.ParseInt(now, 10, 64)
	return time.Unix(0, ns), err == nil
}
