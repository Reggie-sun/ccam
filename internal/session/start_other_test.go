//go:build !linux && !darwin && !windows

package session

import (
	"os/exec"
	"testing"
)

// claudeStart is unknown here; the tests that need it skip.
func claudeStart(t *testing.T, pid int) string { return "" }

const otherStart = "1"

// waitEnded is never reached here: the tests that need it skip.
func waitEnded(t *testing.T, c *exec.Cmd) { t.Skip("no process start times on this platform") }
