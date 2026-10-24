package session

import (
	"os"
	"strconv"
	"testing"
)

// claudeStart is the start time Claude Code records for pid on Windows: the
// creation time as a FILETIME, 100 ns since 1601.
func claudeStart(t *testing.T, pid int) string {
	t.Helper()
	s, ok := StartTime(pid)
	if !ok {
		t.Fatalf("no start time for %d", pid)
	}
	ns, _ := strconv.ParseInt(s, 10, 64)
	return strconv.FormatInt(ns/100+116444736000000000, 10)
}

// otherStart is a start time in Claude Code's form that no test process has
// (1970).
const otherStart = "116444736000000000"

func TestSameStartReadsFiletimes(t *testing.T) {
	const ns = "1792661750632000000" // 2026-10-22, nanoseconds since 1970
	const ft = "134371353506320000"  // the same as a FILETIME
	for _, c := range []struct {
		recorded    string
		same, known bool
	}{
		{ft, true, true},
		{"134371353511320000", true, true},  // half a second later
		{"134371353526320000", false, true}, // two seconds later
		{otherStart, false, true},
		{"639000000000000000", false, false}, // .NET ticks
		{"Thu Oct 22 03:28:37 2026", false, false},
	} {
		if same, known := sameStart(c.recorded, ns); same != c.same || known != c.known {
			t.Errorf("sameStart(%q) = %v, %v", c.recorded, same, known)
		}
	}
	s, _ := StartTime(os.Getpid())
	if same, known := sameStart(claudeStart(t, os.Getpid()), s); !same || !known {
		t.Fatal("this process's own FILETIME")
	}
}
