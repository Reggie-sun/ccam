package session

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// claudeStart is the start time Claude Code records for pid on macOS: what
// `ps -o lstart=` prints in the C locale and UTC.
func claudeStart(t *testing.T, pid int) string {
	t.Helper()
	c := exec.Command("ps", "-o", "lstart=", "-p", strconv.Itoa(pid))
	c.Env = append(os.Environ(), "LC_ALL=C", "TZ=UTC")
	out, err := c.Output()
	if err != nil {
		t.Fatalf("ps for %d: %v", pid, err)
	}
	return strings.TrimSpace(string(out))
}

// otherStart is a start time in Claude Code's form that no test process has.
const otherStart = "Thu Jan  1 00:00:00 1970"

func TestSameStartReadsPsLstart(t *testing.T) {
	at := func(y int, mo time.Month, d, h, mi, s int) string {
		return strconv.FormatInt(time.Date(y, mo, d, h, mi, s, 0, time.UTC).Unix(), 10) + ".250000"
	}
	for _, c := range []struct {
		recorded, now string
		same, known   bool
	}{
		{"Thu Oct 22 03:28:37 2026", at(2026, 10, 22, 3, 28, 37), true, true},
		{"Fri Oct  2 03:28:37 2026", at(2026, 10, 2, 3, 28, 37), true, true},
		{"Fri Oct  2 03:28:37 2026", at(2026, 10, 2, 3, 28, 38), true, true},
		{"Fri Oct  2 03:28:37 2026", at(2026, 10, 2, 3, 28, 40), false, true},
		{"Fri Oct  2 03:28:37 2026", at(2026, 10, 22, 3, 28, 37), false, true},
		{"36646188", at(2026, 10, 2, 3, 28, 37), false, false},
	} {
		if same, known := sameStart(c.recorded, c.now); same != c.same || known != c.known {
			t.Errorf("sameStart(%q, %q) = %v, %v", c.recorded, c.now, same, known)
		}
	}
	// What `ps` prints here is what this compares.
	s, _ := StartTime(os.Getpid())
	if same, known := sameStart(claudeStart(t, os.Getpid()), s); !same || !known {
		t.Fatalf("ps says %q for this process, StartTime %s", claudeStart(t, os.Getpid()), s)
	}
}
