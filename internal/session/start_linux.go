package session

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// StartTime reads field 22 of /proc/<pid>/stat (clock ticks since boot).
func StartTime(pid int) (string, bool) {
	fields, ok := stat(pid)
	if !ok || len(fields) < 20 {
		return "", false
	}
	return fields[19], true
}

// stat is /proc/<pid>/stat from field 3 on. The command name in field 2 may
// contain spaces, so it is split after its ')'.
func stat(pid int) ([]string, bool) {
	data, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return nil, false
	}
	s := string(data)
	i := strings.LastIndexByte(s, ')')
	if i < 0 {
		return nil, false
	}
	return strings.Fields(s[i+1:]), true
}

// sameStart compares a start time Claude Code recorded (procStart) with
// StartTime's: Claude Code reads the same field 22 of /proc/<pid>/stat.
// known is false for a value that is no such number.
func sameStart(recorded, now string) (same, known bool) {
	if _, err := strconv.ParseUint(recorded, 10, 64); err != nil {
		return false, false
	}
	return recorded == now, true
}

// wallStart is when the process of StartTime's ticks started by the clock:
// the boot time (btime in /proc/stat) plus 100 ticks a second.
func wallStart(now string) (time.Time, bool) {
	ticks, err := strconv.ParseInt(now, 10, 64)
	data, err2 := os.ReadFile("/proc/stat")
	if err != nil || err2 != nil {
		return time.Time{}, false
	}
	for line := range strings.Lines(string(data)) {
		if v, ok := strings.CutPrefix(line, "btime "); ok {
			boot, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			return time.Unix(boot, 0).Add(time.Duration(ticks) * 10 * time.Millisecond), err == nil
		}
	}
	return time.Time{}, false
}
