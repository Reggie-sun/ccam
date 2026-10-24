package session

import (
	"path/filepath"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Parent is the process that started process pid, from a process snapshot
// (CreateToolhelp32Snapshot). The snapshot still names a parent that has
// ended, whose pid Windows may give to a newer process: a parent that has
// ended, or that started after pid did, is none.
func Parent(pid int) (int, bool) {
	if pid <= 0 {
		return 0, false
	}
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, false
	}
	defer func() { _ = windows.CloseHandle(snap) }()
	e := windows.ProcessEntry32{Size: uint32(unsafe.Sizeof(windows.ProcessEntry32{}))}
	for err := windows.Process32First(snap, &e); err == nil; err = windows.Process32Next(snap, &e) {
		if e.ProcessID != uint32(pid) {
			continue
		}
		ppid := int(e.ParentProcessID)
		child, ok1 := started(pid)
		parent, ok2 := started(ppid)
		return ppid, ok1 && ok2 && parent <= child
	}
	return 0, false
}

// started is StartTime as nanoseconds since 1970.
func started(pid int) (int64, bool) {
	s, ok := StartTime(pid)
	ns, err := strconv.ParseInt(s, 10, 64)
	return ns, ok && err == nil
}

// Exe is the program process pid runs (QueryFullProcessImageName) while it
// runs, links resolved.
func Exe(pid int) (string, bool) {
	if pid <= 0 {
		return "", false
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return "", false
	}
	defer func() { _ = windows.CloseHandle(h) }()
	var code uint32
	buf := make([]uint16, windows.MAX_LONG_PATH)
	n := uint32(len(buf))
	if windows.GetExitCodeProcess(h, &code) != nil || code != 259 || // STILL_ACTIVE
		windows.QueryFullProcessImageName(h, 0, &buf[0], &n) != nil {
		return "", false
	}
	path := windows.UTF16ToString(buf[:n])
	if real, err := filepath.EvalSymlinks(path); err == nil {
		path = real
	}
	return path, true
}
