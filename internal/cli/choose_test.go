package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// fakeTime makes clock and sleep follow a fake clock.
func fakeTime(t *testing.T, start time.Time) *time.Time {
	now := start
	oc, os := clock, sleep
	clock = func() time.Time { return now }
	sleep = func(d time.Duration) { now = now.Add(d) }
	t.Cleanup(func() { clock, sleep = oc, os })
	return &now
}

func fullReport(name string, resets time.Time) *usage.Report {
	return &usage.Report{Name: name, Status: usage.StatusOK, Limits: &usage.Limits{FiveHour: usage.Window{Percent: 100, ResetsAt: resets, Known: true}}}
}

func TestUntilRoomWaitsForTheFirstReset(t *testing.T) {
	start := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	now := fakeTime(t, start)
	var errOut bytes.Buffer
	e := &env{stderr: &errOut}
	reset := start.Add(40 * time.Minute)
	tries := 0
	prof, err := untilRoom(e, usage.Policy{}, 2*time.Hour, "none", func() (*config.Profile, []*usage.Report, error) {
		tries++
		if now.Before(reset) {
			return nil, []*usage.Report{fullReport("a", reset), fullReport("b", start.Add(3*time.Hour))}, nil
		}
		return &config.Profile{Name: "a"}, nil, nil
	})
	if err != nil || prof.Name != "a" || tries != 2 {
		t.Fatalf("%v %v tries=%d", prof, err, tries)
	}
	if got := now.Sub(start); got != 40*time.Minute+afterReset {
		t.Fatalf("waited %v", got)
	}
	if !strings.Contains(errOut.String(), "waiting until") {
		t.Fatalf("stderr %q", errOut.String())
	}
}

func TestUntilRoomGivesUpWithExitCode75(t *testing.T) {
	start := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	fakeTime(t, start)
	e := &env{stderr: &bytes.Buffer{}}
	none := func() (*config.Profile, []*usage.Report, error) {
		return nil, []*usage.Report{fullReport("a", start.Add(3*time.Hour))}, nil
	}
	for _, wait := range []time.Duration{0, time.Hour} {
		_, err := untilRoom(e, usage.Policy{}, wait, "no room", none)
		var coded codedError
		if err == nil || !asCoded(err, &coded) || coded.code != exitNoRoom {
			t.Fatalf("wait %v: %v", wait, err)
		}
	}
}

func asCoded(err error, out *codedError) bool {
	c, ok := err.(codedError)
	*out = c
	return ok
}

func TestModelArgAndWaitFor(t *testing.T) {
	for in, want := range map[string]string{
		"-p hi --model opus": "opus", "--model=claude-sonnet-5-5 -p x": "claude-sonnet-5-5", "-p x": "", "-- --model opus": "",
	} {
		if got := modelArg(strings.Fields(in)); got != want {
			t.Errorf("%q: %q", in, got)
		}
	}
	if d, err := waitFor("90m"); err != nil || d != 90*time.Minute {
		t.Fatal(d, err)
	}
	for _, bad := range []string{"soon", "-1h", "0s"} {
		if _, err := waitFor(bad); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

func TestNextAutoHonoursMaxUsageAndExits75(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 70, "work": 85, "side": 100})
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin = "/bin/true"
	cfg.Save(c.p.ConfigFile())
	out, errOut, code := c.run("", "--dry-run", "next", "--auto", "--max-usage", "60")
	// 70 and 85 are both over 60, so the lower one overflows; 100 never.
	if code != 0 || !strings.Contains(errOut, "→ ") || !strings.Contains(out, "profiles/personal") && !strings.Contains(out, `profiles\personal`) {
		t.Fatalf("%d %q %q", code, out, errOut)
	}
	c2 := newTestCLI(t, map[string]int{"personal": 100, "work": 100})
	cfg2, _ := config.Load(c2.p.ConfigFile())
	cfg2.ClaudeBin = "/bin/true"
	cfg2.Save(c2.p.ConfigFile())
	_, errOut, code = c2.run("", "--dry-run", "next", "--auto")
	if code != exitNoRoom || !strings.Contains(errOut, "headroom") {
		t.Fatalf("all full: %d %q", code, errOut)
	}
	if _, errOut, code := c.run("", "next", "--max-usage", "50"); code == 0 || !strings.Contains(errOut, "--auto") {
		t.Fatalf("--max-usage without --auto: %d %q", code, errOut)
	}
	if _, errOut, code := c.run("", "--dry-run", "next", "--auto", "--strategy", "random"); code == 0 || !strings.Contains(errOut, "auto_strategy") {
		t.Fatalf("bad strategy: %d %q", code, errOut)
	}
}
