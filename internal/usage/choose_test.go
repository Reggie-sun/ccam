package usage

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)

// rep is an account whose 5-hour window resets in fiveIn and is at five%,
// with the weekly window at week% (negative: unknown).
func rep(name string, five float64, fiveIn time.Duration, week float64, scoped ...Scoped) *Report {
	l := &Limits{FiveHour: Window{Percent: five, ResetsAt: fixedNow.Add(fiveIn), Known: five >= 0}, Scoped: scoped}
	if week >= 0 {
		l.Weekly = Window{Percent: week, ResetsAt: fixedNow.Add(72 * time.Hour), Known: true}
	}
	if five < 0 {
		l.FiveHour = Window{}
	}
	return &Report{Name: name, Status: StatusOK, Limits: l}
}

func repName(r *Report) string {
	if r == nil {
		return ""
	}
	return r.Name
}

func TestPolicyChoose(t *testing.T) {
	cases := []struct {
		name   string
		policy Policy
		in     []*Report
		want   string
	}{
		{"lowest 5-hour usage", Policy{}, []*Report{rep("a", 60, time.Hour, 10), rep("b", 20, 3*time.Hour, 10)}, "b"},
		{"full is skipped", Policy{}, []*Report{rep("a", 100, time.Hour, 10), rep("b", 90, time.Hour, 10)}, "b"},
		{"weekly full is skipped", Policy{}, []*Report{rep("a", 5, time.Hour, 100), rep("b", 90, time.Hour, 10)}, "b"},
		{"unknown is skipped", Policy{}, []*Report{rep("a", -1, 0, 10), rep("b", 90, time.Hour, 10)}, "b"},
		{"nothing left", Policy{}, []*Report{rep("a", 100, time.Hour, 10)}, ""},
		{"below the threshold first", Policy{MaxUsage: 80}, []*Report{rep("a", 85, time.Hour, 10), rep("b", 70, time.Hour, 79)}, "b"},
		{"weekly over the threshold", Policy{MaxUsage: 80}, []*Report{rep("a", 10, time.Hour, 85), rep("b", 70, time.Hour, 20)}, "b"},
		{"overflow when all are over", Policy{MaxUsage: 80}, []*Report{rep("a", 95, time.Hour, 10), rep("b", 85, time.Hour, 10)}, "b"},
		{"expiring takes the soonest reset", Policy{Strategy: Expiring}, []*Report{rep("a", 10, 4*time.Hour, 10), rep("b", 50, 30*time.Minute, 10)}, "b"},
		{"expiring still skips full", Policy{Strategy: Expiring}, []*Report{rep("a", 10, 4*time.Hour, 10), rep("b", 100, 30*time.Minute, 10)}, "a"},
		{"model window counts", Policy{Model: "opus"}, []*Report{
			rep("a", 10, time.Hour, 10, Scoped{Name: "Opus", Window: Window{Percent: 100, Known: true}}),
			rep("b", 50, time.Hour, 10, Scoped{Name: "Opus", Window: Window{Percent: 20, Known: true}}),
		}, "b"},
		{"other model windows do not", Policy{Model: "claude-sonnet-5-5"}, []*Report{
			rep("a", 10, time.Hour, 10, Scoped{Name: "Opus", Window: Window{Percent: 100, Known: true}}),
			rep("b", 50, time.Hour, 10),
		}, "a"},
		{"logged out is skipped", Policy{}, []*Report{{Name: "a", Status: StatusLoggedOut, Limits: rep("a", 1, time.Hour, 1).Limits}, rep("b", 50, time.Hour, 10)}, "b"},
	}
	for _, c := range cases {
		if got := repName(c.policy.Choose(c.in, fixedNow)); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

func TestPolicyFreeAt(t *testing.T) {
	// Only full accounts block a choice (over the threshold still overflows),
	// so FreeAt is when the first of them has a window reset.
	p := Policy{MaxUsage: 80}
	in := []*Report{
		rep("a", 100, 2*time.Hour, 10),     // 5-hour full: free in 2h
		rep("b", 50, 4*time.Hour, 100),     // weekly full: free in 72h
		rep("c", 100, 30*time.Minute, 100), // both full: free when the later one resets, 72h
		rep("d", -1, 0, 10),                // unknown: no time
	}
	if got := p.FreeAt(in, fixedNow); !got.Equal(fixedNow.Add(2 * time.Hour)) {
		t.Fatalf("FreeAt %v", got.Sub(fixedNow))
	}
	if got := p.FreeAt([]*Report{rep("d", -1, 0, -1)}, fixedNow); !got.IsZero() {
		t.Fatalf("unknown only: %v", got)
	}
}

func TestPolicyCheck(t *testing.T) {
	for _, bad := range []Policy{{MaxUsage: -1}, {MaxUsage: 101}, {Strategy: "random"}} {
		if bad.Check() == nil {
			t.Errorf("accepted %+v", bad)
		}
	}
	if (Policy{MaxUsage: 80, Strategy: Expiring}).Check() != nil || (Policy{}).Check() != nil {
		t.Fatal("valid policy refused")
	}
}
