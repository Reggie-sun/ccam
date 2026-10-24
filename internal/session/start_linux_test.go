package session

import (
	"os"
	"testing"
	"time"
)

// claudeStart is the start time Claude Code records for pid: on Linux the
// same field of /proc/<pid>/stat that StartTime reads.
func claudeStart(t *testing.T, pid int) string {
	t.Helper()
	s, ok := StartTime(pid)
	if !ok {
		t.Fatalf("no start time for %d", pid)
	}
	return s
}

// otherStart is a start time in Claude Code's form that no test process has.
const otherStart = "1"

func TestSameStartComparesTicks(t *testing.T) {
	for _, c := range []struct {
		recorded, now string
		same, known   bool
	}{
		{"36646188", "36646188", true, true},
		{"36646187", "36646188", false, true},
		{"Thu Oct 22 03:28:37 2026", "36646188", false, false},
		{"", "36646188", false, false},
	} {
		if same, known := sameStart(c.recorded, c.now); same != c.same || known != c.known {
			t.Errorf("sameStart(%q, %q) = %v, %v", c.recorded, c.now, same, known)
		}
	}
	s, _ := StartTime(os.Getpid())
	at, ok := wallStart(s)
	if !ok || at.After(time.Now().Add(2*time.Second)) || at.Before(time.Now().Add(-time.Hour)) {
		t.Fatalf("this test began at %v (%v)", at, ok)
	}
}
