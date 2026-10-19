package usage

import (
	"strings"
	"testing"
	"time"
)

func TestTableMatchesPlanExample(t *testing.T) {
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	w := func(p float64, d time.Duration) Window { return Window{Percent: p, ResetsAt: now.Add(d), Known: true} }
	reports := []*Report{
		{Name: "personal", Slot: 1, Status: StatusOK, Source: "live", At: now, Limits: &Limits{FiveHour: w(72, 100*time.Minute), Weekly: w(41, 78*time.Hour)}},
		{Name: "work", Slot: 2, Status: StatusOK, Source: "live", At: now, Limits: &Limits{FiveHour: w(8, 252*time.Minute), Weekly: w(19, 122*time.Hour)}},
		{Name: "side", Slot: 3, Status: StatusError, Source: "cache", At: now.Add(-12 * time.Minute), Limits: &Limits{FiveHour: w(100, 22*time.Minute), Weekly: w(68, 33*time.Hour)}},
		{Name: "api", Slot: 4, Status: StatusAPI},
		{Name: "old", Slot: 5, Status: StatusExpired},
	}
	lines := Table(reports, RenderOptions{Now: now})
	want := []string{
		"██████░░  72% 1시간 40분",
		"█░░░░░░░   8% 4시간 12분",
		"████████ 100% 22분",
		"3일 6시간",
		"12분 전 (캐시)",
		"한도 없음 (API)",
		"로그인 만료, `ccam old`",
	}
	all := strings.Join(lines, "\n")
	for _, s := range want {
		if !strings.Contains(all, s) {
			t.Errorf("missing %q in\n%s", s, all)
		}
	}
	for _, l := range lines {
		if width.StringWidth(l) > 80 {
			t.Errorf("wider than 80 columns (%d): %s", width.StringWidth(l), l)
		}
	}
	if b := Best(reports); b == nil || b.Name != "work" {
		t.Fatalf("best %v", b)
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

func TestVerboseShowsScoped(t *testing.T) {
	r := &Report{Name: "a", Slot: 1, Status: StatusOK, Source: "live", Limits: &Limits{
		FiveHour: Window{Percent: 1, Known: true},
		Scoped:   []Scoped{{Name: "Fable", Window: Window{Percent: 88, Known: true}}},
	}}
	out := strings.Join(Table([]*Report{r}, RenderOptions{Verbose: true, Now: time.Now()}), "\n")
	if !strings.Contains(out, "Fable 주간") || !strings.Contains(out, "88%") || !strings.Contains(out, "출처: live") {
		t.Fatal(out)
	}
}

func TestDuration(t *testing.T) {
	cases := map[time.Duration]string{
		100 * time.Minute: "1시간 40분",
		78 * time.Hour:    "3일 6시간",
		22 * time.Minute:  "22분",
		10 * time.Second:  "1분",
		-time.Minute:      "곧",
	}
	for d, want := range cases {
		if got := Duration(d); got != want {
			t.Errorf("%v: %q want %q", d, got, want)
		}
	}
}
