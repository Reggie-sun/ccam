package usage

import (
	"strings"
	"testing"
	"time"
)

func TestTableShowsABarPerLimit(t *testing.T) {
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	w := func(p float64, d time.Duration) Window { return Window{Percent: p, ResetsAt: now.Add(d), Known: true} }
	reports := []*Report{
		{Name: "personal", Slot: 1, Status: StatusOK, Source: "live", At: now, Limits: &Limits{FiveHour: w(72, 100*time.Minute), Weekly: w(41, 78*time.Hour),
			Scoped: []Scoped{{Name: "Fable", Window: w(88, 30*time.Hour)}}}},
		{Name: "work", Slot: 2, Status: StatusOK, Source: "live", At: now, Limits: &Limits{FiveHour: w(8, 252*time.Minute), Weekly: w(19, 122*time.Hour)}},
		{Name: "side", Slot: 3, Status: StatusError, Source: "cache", At: now.Add(-12 * time.Minute), Limits: &Limits{FiveHour: w(100, 22*time.Minute), Weekly: w(68, 33*time.Hour)}},
		{Name: "fable", Slot: 4, Status: StatusOK, Source: "live", At: now, Limits: &Limits{FiveHour: w(2, time.Hour),
			Scoped: []Scoped{{Name: "Fable", Window: w(0, 50*time.Hour)}}}},
		{Name: "api", Slot: 5, Status: StatusAPI},
		{Name: "old", Slot: 6, Status: StatusExpired},
	}
	lines := Table(reports, RenderOptions{Now: now})
	want := []string{
		" 1 personal  just now",
		"   5-hour   ██████████████░░░░░░  72%  resets in 1h 40m",
		"   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 6h",
		"   Fable    ██████████████████░░  88%  resets in 1d 6h",
		" 2 work      just now",
		"   5-hour   ██░░░░░░░░░░░░░░░░░░   8%  resets in 4h 12m",
		"   Weekly   ████░░░░░░░░░░░░░░░░  19%  resets in 5d 2h",
		" 3 side      12m ago (cached)",
		"   5-hour   ████████████████████ 100%  resets in 22m",
		"   Weekly   ██████████████░░░░░░  68%  resets in 1d 9h",
		" 4 fable     just now",
		"   5-hour   ░░░░░░░░░░░░░░░░░░░░   2%  resets in 1h 0m",
		"   Weekly   -",
		"   Fable    ░░░░░░░░░░░░░░░░░░░░   0%  resets in 2d 2h",
		" 5 api       no limits (API)",
		" 6 old       login expired, run `ccam old` once to refresh",
	}
	if got, exp := strings.Join(lines, "\n"), strings.Join(want, "\n"); got != exp {
		t.Fatalf("got\n%s\nwant\n%s", got, exp)
	}
	for _, l := range lines {
		if width.StringWidth(l) > 80 {
			t.Errorf("wider than 80 columns (%d): %s", width.StringWidth(l), l)
		}
	}
}

func TestColorsKeepPercentText(t *testing.T) {
	g := gauge(Window{Percent: 85, Known: true}, RenderOptions{Color: true})
	if !strings.Contains(g, "\x1b[33m") || !strings.Contains(g, "85%") {
		t.Fatalf("%q", g)
	}
	g = gauge(Window{Percent: 100, Known: true}, RenderOptions{Color: true})
	if !strings.Contains(g, "\x1b[31m") || !strings.Contains(g, "100%") {
		t.Fatalf("%q", g)
	}
	if g := gauge(Window{Percent: 50, Known: true}, RenderOptions{Color: true}); strings.Contains(g, "\x1b[") {
		t.Fatal("no color under 80%")
	}
}

func TestVerboseShowsSource(t *testing.T) {
	r := &Report{Name: "a", Slot: 1, Status: StatusOK, Source: "live", Limits: &Limits{
		FiveHour: Window{Percent: 1, Known: true},
		Scoped:   []Scoped{{Name: "A long model name", Window: Window{Percent: 88, Known: true}}},
	}}
	out := strings.Join(Table([]*Report{r}, RenderOptions{Now: time.Now()}), "\n")
	if strings.Contains(out, "source:") || !strings.Contains(out, "   A long … ") || !strings.Contains(out, "88%") {
		t.Fatal(out)
	}
	out = strings.Join(Table([]*Report{r}, RenderOptions{Verbose: true, Now: time.Now()}), "\n")
	if !strings.Contains(out, "source: live") {
		t.Fatal(out)
	}
}

func TestDuration(t *testing.T) {
	cases := map[time.Duration]string{
		100 * time.Minute: "1h 40m",
		78 * time.Hour:    "3d 6h",
		22 * time.Minute:  "22m",
		10 * time.Second:  "1m",
		-time.Minute:      "soon",
	}
	for d, want := range cases {
		if got := Duration(d); got != want {
			t.Errorf("%v: %q want %q", d, got, want)
		}
	}
}
