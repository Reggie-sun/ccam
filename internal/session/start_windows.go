package session

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// startTime asks Windows when the process started (GetProcessTimes). A
// process that has exited but is still referenced has an exit code other
// than STILL_ACTIVE and does not count.
func startTime(pid int) (string, bool) {
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
