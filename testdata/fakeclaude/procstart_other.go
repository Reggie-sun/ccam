//go:build !windows

package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// procStart is this process's start as Claude Code records it: field 22 of
// /proc/<pid>/stat on Linux, else what `ps -o lstart=` prints in the C
// locale and UTC.
func procStart() string {
	pid := strconv.Itoa(os.Getpid())
	if data, err := os.ReadFile("/proc/" + pid + "/stat"); err == nil {
		s := string(data)
		if f := strings.Fields(s[strings.LastIndexByte(s, ')')+1:]); len(f) >= 20 {
			return f[19]
		}
		return ""
	}
	c := exec.Command("ps", "-o", "lstart=", "-p", pid)
	c.Env = append(os.Environ(), "LC_ALL=C", "TZ=UTC")
	out, _ := c.Output()
	return strings.TrimSpace(string(out))
}
