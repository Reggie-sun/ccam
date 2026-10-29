package cli

import (
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestLiveSessionQuiet(t *testing.T) {
	now := time.Now()
	ago := func(d time.Duration) int64 { return now.Add(-d).UnixMilli() }
	for _, c := range []struct {
		name string
		s    liveSession
		want bool
	}{
		{"idle long enough", liveSession{Kind: "interactive", Status: "idle", StatusUpdatedAt: ago(5 * time.Minute)}, true},
		{"idle a moment ago", liveSession{Kind: "interactive", Status: "idle", StatusUpdatedAt: ago(30 * time.Second)}, false},
		{"busy", liveSession{Kind: "interactive", Status: "busy", StatusUpdatedAt: ago(time.Hour)}, false},
		{"waiting for an answer", liveSession{Kind: "interactive", Status: "waiting", StatusUpdatedAt: ago(time.Hour)}, false},
		{"background", liveSession{Kind: "bg", Status: "idle", StatusUpdatedAt: ago(time.Hour)}, false},
		{"no time", liveSession{Kind: "interactive", Status: "idle"}, false},
	} {
		if got := c.s.quiet(now, 3*time.Minute); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestProcessStopped(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Skip(err)
	}
	defer cmd.Process.Kill()
	if processStopped(cmd.Process.Pid) {
		t.Fatal("a running process counts as stopped")
	}
	cmd.Process.Signal(syscall.SIGSTOP)
	for i := 0; i < 50 && !processStopped(cmd.Process.Pid); i++ {
		time.Sleep(20 * time.Millisecond)
	}
	if !processStopped(cmd.Process.Pid) {
		t.Fatal("a stopped process does not count as stopped")
	}
}
