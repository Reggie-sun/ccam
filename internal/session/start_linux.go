package session

import (
	"os"
	"strconv"
	"strings"
)

// startTime reads field 22 of /proc/<pid>/stat (clock ticks since boot).
// The command name in field 2 may contain spaces, so parse after its ')'.
func startTime(pid int) (string, bool) {
	data, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return "", false
	}
	s := string(data)
	i := strings.LastIndexByte(s, ')')
	if i < 0 {
		return "", false
	}
	fields := strings.Fields(s[i+1:])
	if len(fields) < 20 {
		return "", false
	}
	return fields[19], true
}
