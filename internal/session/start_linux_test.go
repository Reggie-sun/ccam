package session

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"golang.org/x/sys/unix"
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

// waitEnded waits until c's process has ended, without reaping it: until
// c.Wait it is a zombie.
func waitEnded(t *testing.T, c *exec.Cmd) {
	t.Helper()
	var info unix.Siginfo
	if err := unix.Waitid(unix.P_PID, c.Process.Pid, &info, unix.WEXITED|unix.WNOWAIT, nil); err != nil {
		t.Fatal(err)
	}
}

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

func TestEndedStates(t *testing.T) {
	// Fields from field 3 on: state first, num_threads at index 17.
	fields := func(state, threads string) []string {
		f := make([]string, 20)
		f[0], f[17] = state, threads
		return f
	}
	for _, c := range []struct {
		state, threads string
		ended          bool
	}{
		{"R", "1", false},
		{"S", "8", false},
		{"Z", "1", true},
		{"Z", "0", true},
		{"Z", "2", false}, // main thread gone, others still running
		{"X", "0", true},
	} {
		if got := ended(fields(c.state, c.threads)); got != c.ended {
			t.Errorf("ended(state %s, %s threads) = %v", c.state, c.threads, got)
		}
	}
}
