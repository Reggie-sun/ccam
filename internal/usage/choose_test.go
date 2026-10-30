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

func weekIn(r *Report, d time.Duration) *Report {
	r.Limits.Weekly.ResetsAt = fixedNow.Add(d)
	return r
}

func sessions(r *Report, n int) *Report {
	r.Sessions = n
	return r
}

func TestPolicyWhy(t *testing.T) {
	a, b := rep("a", 12, 42*time.Minute, 41), sessions(rep("b", 30, time.Hour, 20), 1)
	cases := []struct {
		policy Policy
		in     *Report
		want   string
	}{
		{Policy{}, a, "12% of its 5 hours used, the least"},
		{Policy{Strategy: Expiring}, a, "its 5-hour window resets first, in 42m (12% of its 5 hours used)"},
		{Policy{Strategy: Balance}, b, "80% of its week is left, resetting in 3d 0h; 1 running"},
		{Policy{MaxUsage: 10}, a, "every account is at 10% or more; it has the lowest use, 12% of its 5 hours used"},
		{Policy{}, rep("c", -1, 0, -1), ""},
	}
	for _, c := range cases {
		if got := c.policy.Why(c.in, fixedNow); got != c.want {
			t.Errorf("%+v for %s: got %q, want %q", c.policy, c.in.Name, got, c.want)
		}
	}
}

// Use says what Why says without comparing the account with the others:
// for an account taken while a better one was passed over.
func TestPolicyUse(t *testing.T) {
	a, b := rep("a", 12, 42*time.Minute, 41), sessions(rep("b", 30, time.Hour, 20), 1)
	cases := []struct {
		policy Policy
		in     *Report
		want   string
	}{
		{Policy{}, a, "12% of its 5 hours used"},
		{Policy{Strategy: Expiring}, a, "its 5-hour window resets in 42m (12% of its 5 hours used)"},
		{Policy{Strategy: Balance}, b, "80% of its week is left, resetting in 3d 0h; 1 running"},
		{Policy{MaxUsage: 10}, a, "12% of its 5 hours used"},
		{Policy{MaxUsage: 10, Strategy: Balance}, b, "30% of its 5 hours used"},
		{Policy{}, rep("c", -1, 0, -1), ""},
	}
	for _, c := range cases {
		if got := c.policy.Use(c.in, fixedNow); got != c.want {
			t.Errorf("%+v for %s: got %q, want %q", c.policy, c.in.Name, got, c.want)
		}
	}
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
		// a: 90% left over 3 days (30/day); b: 50% left over 1 day (50/day).
		{"balance takes the week that would go unused", Policy{Strategy: Balance}, []*Report{rep("a", 10, time.Hour, 10), weekIn(rep("b", 60, time.Hour, 50), 24*time.Hour)}, "b"},
		{"balance counts running sessions", Policy{Strategy: Balance}, []*Report{rep("a", 10, time.Hour, 10), sessions(weekIn(rep("b", 60, time.Hour, 50), 24*time.Hour), 3)}, "a"},
		{"balance still skips full", Policy{Strategy: Balance}, []*Report{rep("a", 10, time.Hour, 10), weekIn(rep("b", 100, time.Hour, 50), 24*time.Hour)}, "a"},
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
	if (Policy{MaxUsage: 80, Strategy: Expiring}).Check() != nil || (Policy{Strategy: Balance}).Check() != nil || (Policy{}).Check() != nil {
		t.Fatal("valid policy refused")
	}
}

// Full is the limit Choose sees: the 5-hour, weekly and --model windows.
func TestPolicyFull(t *testing.T) {
	opus := Scoped{Name: "Opus", Window: Window{Percent: 100, Known: true}}
	cases := []struct {
		policy Policy
		in     *Report
		want   bool
	}{
		{Policy{}, nil, false},
		{Policy{}, &Report{Name: "x"}, false},
		{Policy{}, rep("a", 100, time.Hour, 10), true},
		{Policy{}, rep("a", 20, time.Hour, 100), true},
		{Policy{}, rep("a", 20, time.Hour, 10, opus), false},
		{Policy{Model: "opus"}, rep("a", 20, time.Hour, 10, opus), true},
		{Policy{Model: "sonnet"}, rep("a", 20, time.Hour, 10, opus), false},
		{Policy{Model: "opus"}, rep("a", 20, time.Hour, 10, Scoped{Name: "Opus", Window: Window{Percent: 100}}), false}, // unknown
	}
	for i, tc := range cases {
		if got := tc.policy.Full(tc.in); got != tc.want {
			t.Errorf("%d: %v", i, got)
		}
	}
}

func TestPolicyChooseWeeklyReset(t *testing.T) {
	p := Policy{Strategy: WeeklyReset}
	soon := weekIn(rep("soon", 20, time.Hour, 80), 6*time.Hour)      // 80% of the five hours left
	later := weekIn(rep("later", 0, time.Hour, 0), 48*time.Hour)     // all left
	tieLow := weekIn(rep("tie-low", 20, time.Hour, 10), 6*time.Hour) // resets with soon, less used
	tight := weekIn(rep("tight", 70, time.Hour, 0), 6*time.Hour)     // 30% left: resets first, no room
	tightLater := weekIn(rep("tight-later", 80, time.Hour, 0), 96*time.Hour)
	full := weekIn(rep("full", 100, time.Hour, 10), time.Hour)
	unknown := rep("unknown", 0, time.Hour, -1)
	for _, c := range []struct {
		name string
		in   []*Report
		want string
	}{
		{"soonest weekly reset first", []*Report{later, soon}, "soon"},
		{"equal reset: most weekly left", []*Report{soon, tieLow}, "tie-low"},
		{"five-hour room before an earlier reset", []*Report{tight, later}, "later"},
		{"59% left is no room, 60% is", []*Report{weekIn(rep("p41", 41, time.Hour, 0), 6*time.Hour), weekIn(rep("p40", 40, time.Hour, 0), 48*time.Hour)}, "p40"},
		{"nobody has room: soonest weekly reset", []*Report{tightLater, tight}, "tight"},
		{"full account skipped", []*Report{full, later}, "later"},
		{"unknown week last", []*Report{unknown, later}, "later"},
		{"a running account is still taken when its week resets first", []*Report{sessions(weekIn(rep("busy", 0, time.Hour, 0), 6*time.Hour), 1), later}, "busy"},
		{"equal resets: idle first", []*Report{sessions(weekIn(rep("busy", 0, time.Hour, 0), 6*time.Hour), 1), weekIn(rep("idle", 0, time.Hour, 0), 6*time.Hour)}, "idle"},
		{"all running: soonest weekly reset", []*Report{sessions(weekIn(rep("few", 0, time.Hour, 0), 96*time.Hour), 1), sessions(weekIn(rep("many", 0, time.Hour, 0), 6*time.Hour), 3)}, "many"},
	} {
		if got := p.Choose(c.in, fixedNow); got == nil || got.Name != c.want {
			t.Errorf("%s: got %v, want %s", c.name, got, c.want)
		}
	}
	// Drain ignores the five-hour room.
	if got := (Policy{Strategy: Drain}).Choose([]*Report{tight, later}, fixedNow); got == nil || got.Name != "tight" {
		t.Errorf("drain: got %v, want tight", got)
	}
}

func TestPolicyChooseDrain(t *testing.T) {
	p := Policy{Strategy: Drain}
	busy := sessions(weekIn(rep("busy", 60, time.Hour, 50), 6*time.Hour), 2)
	idle := weekIn(rep("idle", 0, time.Hour, 0), 48*time.Hour)
	if got := p.Choose([]*Report{idle, busy}, fixedNow); got == nil || got.Name != "busy" {
		t.Errorf("drain keeps the running account whose week resets first: got %v", got)
	}
	full := sessions(weekIn(rep("full", 100, time.Hour, 50), 6*time.Hour), 1)
	if got := p.Choose([]*Report{full, idle}, fixedNow); got == nil || got.Name != "idle" {
		t.Errorf("drain moves on once the account is full: got %v", got)
	}
	if err := p.Check(); err != nil {
		t.Error(err)
	}
}

func TestPolicyChooseSkipsCrowdedAccounts(t *testing.T) {
	// "first" resets its week soonest, so both strategies prefer it.
	first := func(sess int, five float64) *Report {
		return sessions(weekIn(rep("first", five, time.Hour, 10), 6*time.Hour), sess)
	}
	// Over 40% of the five hours used for both accounts, so weekly-reset's
	// five-hour room rule does not decide either.
	other := sessions(weekIn(rep("other", 55, time.Hour, 10), 48*time.Hour), 1)
	for _, strategy := range []string{WeeklyReset, Drain} {
		p := Policy{Strategy: strategy}
		for _, c := range []struct {
			name string
			in   []*Report
			want string
		}{
			{"many sessions and over half used: skipped", []*Report{first(3, 60), other}, "other"},
			{"many sessions, half or less used: kept", []*Report{first(3, 50), other}, "first"},
			{"few sessions, over half used: kept", []*Report{first(2, 90), other}, "first"},
			{"every account crowded: the strategy decides", []*Report{first(4, 70), sessions(weekIn(rep("other", 60, time.Hour, 10), 48*time.Hour), 3)}, "first"},
		} {
			if got := p.Choose(c.in, fixedNow); got == nil || got.Name != c.want {
				t.Errorf("%s / %s: got %v, want %s", strategy, c.name, got, c.want)
			}
		}
	}
	// Other strategies are left as they were.
	if got := (Policy{Strategy: Balance}).Choose([]*Report{sessions(weekIn(rep("first", 60, time.Hour, 0), 6*time.Hour), 3), weekIn(rep("other", 5, time.Hour, 90), 48*time.Hour)}, fixedNow); got == nil || got.Name != "first" {
		t.Errorf("balance changed: %v", got)
	}
}
