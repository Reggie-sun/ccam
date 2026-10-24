package e2e

import (
	"os/exec"
	"runtime"
	"testing"
)

// ptyCommand runs the sh command line on a terminal of its own, through
// script(1): util-linux's on Linux, the BSD one on macOS. Elsewhere the test
// is skipped. What the command writes to the terminal is script's output.
func ptyCommand(t *testing.T, line string) *exec.Cmd {
	t.Helper()
	script, err := exec.LookPath("script")
	if err != nil {
		t.Skip("needs script(1) for a terminal")
	}
	switch runtime.GOOS {
	case "linux":
		return exec.Command(script, "-qec", line, "/dev/null")
	case "darwin":
		return exec.Command(script, "-q", "/dev/null", "sh", "-c", line)
	}
	t.Skip("no terminal helper on " + runtime.GOOS)
	return nil
}
