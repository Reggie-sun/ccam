//go:build !linux && !darwin && !windows

package session

import "testing"

// claudeStart is unknown here; the tests that need it skip.
func claudeStart(t *testing.T, pid int) string { return "" }

const otherStart = "1"
