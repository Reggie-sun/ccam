package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func window(p float64, reset time.Time) Window {
	return Window{Percent: p, ResetsAt: reset, Known: true}
}

// A window Claude Code refused a request in counts as full until its reset,
// whatever the saved numbers say and however they were found.
func TestLimitRecordCountsAsFull(t *testing.T) {
	ctx := context.Background()
	now := time.Now().Truncate(time.Second)
	reset := now.Add(3 * time.Hour)
	full := func(t *testing.T, r Report) {
		t.Helper()
		if l := r.Limits; l == nil || !l.FiveHour.Known || l.FiveHour.Percent != 100 || !l.FiveHour.ResetsAt.Equal(reset) {
			t.Fatalf("not full until the reset: %+v", r)
		}
		if (Policy{}).Choose([]*Report{&r}, now) != nil || !(Policy{}).FreeAt([]*Report{&r}, now).Equal(reset) {
			t.Fatalf("choice: %+v", r.Limits)
		}
	}
	t.Run("live cache 20 s old at 40%", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		a := e.account("a", "ok-new-format.json", future)
		e.f.writeCache("a", "live", cached{At: now.Add(-20 * time.Second), Limits: &Limits{FiveHour: window(40, now.Add(time.Hour))}})
		if err := SaveLimit(e.f.StateDir, "a", "five_hour", reset, now); err != nil {
			t.Fatal(err)
		}
		r := e.f.Fetch(ctx, a, Options{})
		if e.hits.Load() != 0 || r.Source != "live" {
			t.Fatalf("the cache within 60 s: %+v hits=%d", r, e.hits.Load())
		}
		full(t, r)
		r = e.f.Fetch(ctx, a, Options{NoCache: true}) // a live answer, too
		if e.hits.Load() != 1 {
			t.Fatal("--no-cache must ask")
		}
		full(t, r)
	})
	t.Run("expired token and a 3-hour-old cache at 10%", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		a := e.account("a", "ok-new-format.json", now.Add(-time.Hour).UnixMilli())
		e.f.writeCache("a", "live", cached{At: now.Add(-3 * time.Hour), Limits: &Limits{FiveHour: window(10, now.Add(time.Hour))}})
		SaveLimit(e.f.StateDir, "a", "five_hour", reset, now)
		r := e.f.Fetch(ctx, a, Options{})
		if r.Status != StatusExpired || r.Source != "cache" || e.hits.Load() != 0 {
			t.Fatalf("%+v", r)
		}
		full(t, r)
	})
	t.Run("saved values only, status line", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		a := e.account("a", "boom", future)
		SaveStatusLine(e.f.StateDir, "a", Limits{FiveHour: window(23, now.Add(time.Hour))}, now)
		SaveLimit(e.f.StateDir, "a", "five_hour", reset, now)
		for _, opt := range []Options{{Source: "cache"}, {}} {
			full(t, e.f.Fetch(ctx, a, opt))
		}
	})
	t.Run("Claude Code's lookup", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		e.f.Source = "claude"
		e.f.Probe = func(context.Context, Account) ([]byte, error) {
			return []byte(`{"type":"control_response","response":{"subtype":"success","request_id":"ccam","response":{"rate_limits_available":true,` +
				`"rate_limits":{"five_hour":{"utilization":30,"resets_at":"` + now.Add(time.Hour).UTC().Format(time.RFC3339) + `"}}}}}`), nil
		}
		a := e.account("a", "ok-new-format.json", future)
		SaveLimit(e.f.StateDir, "a", "five_hour", reset, now)
		full(t, e.f.Fetch(ctx, a, Options{}))
		full(t, e.f.Fetch(ctx, a, Options{})) // its 60 s cache
	})
	t.Run("nothing saved", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		a := e.account("a", "boom", future)
		SaveLimit(e.f.StateDir, "a", "five_hour", reset, now.Add(-time.Minute))
		r := e.f.Fetch(ctx, a, Options{})
		full(t, r)
		if r.Status != StatusError || r.Source != "limit" || !r.At.Equal(now.Add(-time.Minute)) || r.Limits.Weekly.Known {
			t.Fatalf("%+v", r)
		}
	})
	t.Run("no subscription", func(t *testing.T) {
		e := newEnv(t)
		e.f.Now = func() time.Time { return now }
		out := e.account("out", "", 0)
		api := e.account("api", "", 0)
		os.WriteFile(filepath.Join(api.Dir, ".claude.json"), []byte(`{"primaryApiKey":"sk-x"}`), 0o600)
		for _, a := range []Account{out, api} {
			SaveLimit(e.f.StateDir, a.Name, "five_hour", reset, now)
			if r := e.f.Fetch(ctx, a, Options{}); r.Limits != nil {
				t.Fatalf("%s: %+v", a.Name, r)
			}
		}
	})
}

// Opus and Sonnet records fill the model's own window: the one whose name
// has the family in it, else a new one.
func TestLimitRecordFillsModelWindows(t *testing.T) {
	e := newEnv(t)
	now := time.Now().Truncate(time.Second)
	e.f.Now = func() time.Time { return now }
	a := e.account("a", "boom", future)
	e.f.writeCache("a", "live", cached{At: now.Add(-time.Hour), Limits: &Limits{FiveHour: window(10, now.Add(time.Hour)),
		Scoped: []Scoped{{Name: "Fable", Window: window(30, now.Add(48*time.Hour))}, {Name: "Claude Opus", Window: window(20, now.Add(48*time.Hour))}}}})
	SaveLimit(e.f.StateDir, "a", "seven_day_opus", now.Add(50*time.Hour), now)
	SaveLimit(e.f.StateDir, "a", "seven_day_sonnet", now.Add(60*time.Hour), now)
	r := e.f.Fetch(context.Background(), a, Options{Source: "cache"})
	s := r.Limits.Scoped
	if len(s) != 3 || s[0].Name != "Fable" || s[0].Window.Percent != 30 || s[1].Name != "Claude Opus" || s[1].Window.Percent != 100 ||
		!s[1].Window.ResetsAt.Equal(now.Add(50*time.Hour)) || s[2].Name != "Sonnet" || s[2].Window.Percent != 100 {
		t.Fatalf("%+v", s)
	}
	if !(Policy{Model: "opus"}).Full(&r) || !(Policy{Model: "claude-sonnet-5-5"}).Full(&r) || (Policy{Model: "fable"}).Full(&r) || (Policy{}).Full(&r) {
		t.Fatal("full by model")
	}
	// The saved cache itself is left as it was.
	if c, _ := e.f.readCache("a", "live"); c.Limits.Scoped[1].Window.Percent != 20 || len(c.Limits.Scoped) != 2 {
		t.Fatalf("cache changed: %+v", c.Limits)
	}
}

// A record whose window has reset no longer counts, and the next record
// drops it.
func TestLimitRecordExpires(t *testing.T) {
	e := newEnv(t)
	now := time.Now().Truncate(time.Second)
	a := e.account("a", "boom", future)
	e.f.writeCache("a", "live", cached{At: now, Limits: &Limits{FiveHour: window(40, time.Time{})}})
	SaveLimit(e.f.StateDir, "a", "five_hour", now.Add(time.Hour), now)
	later := now.Add(2 * time.Hour)
	e.f.Now = func() time.Time { return later }
	if r := e.f.Fetch(context.Background(), a, Options{Source: "cache"}); r.Limits.FiveHour.Percent != 40 {
		t.Fatalf("an expired record: %+v", r.Limits)
	}
	if err := SaveLimit(e.f.StateDir, "a", "seven_day", later.Add(48*time.Hour), later); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(e.f.StateDir, "a.limit.json")
	data, _ := os.ReadFile(path)
	var rec struct{ Windows map[string]time.Time }
	if json.Unmarshal(data, &rec) != nil || len(rec.Windows) != 1 || !rec.Windows["seven_day"].Equal(later.Add(48*time.Hour)) {
		t.Fatalf("%s", data)
	}
	if info, _ := os.Stat(path); runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatal(info.Mode())
	}
	// Windows ccam does not choose by, and limits that are over, are not kept.
	if err := SaveLimit(e.f.StateDir, "a", "overage", later.Add(time.Hour), later); err == nil {
		t.Fatal("overage recorded")
	}
	if err := SaveLimit(e.f.StateDir, "b", "five_hour", later.Add(-time.Minute), later); err != nil || fileExists(filepath.Join(e.f.StateDir, "b.limit.json")) {
		t.Fatal("a limit that is over was recorded")
	}
	found := false
	for _, p := range CachePaths(e.f.StateDir, "a") {
		found = found || p == path
	}
	if !found {
		t.Fatal("CachePaths must list the record, for rename and remove")
	}
}

func fileExists(path string) bool { _, err := os.Stat(path); return err == nil }

// BenchmarkFetchCache is `ccam usage --source cache` for five accounts:
// saved values only, without and with a limit record.
func BenchmarkFetchCache(b *testing.B) {
	for _, limit := range []bool{false, true} {
		b.Run(map[bool]string{false: "no-limit", true: "limit"}[limit], func(b *testing.B) {
			root := b.TempDir()
			f := NewFetcher(filepath.Join(root, "state"))
			now := time.Now()
			var accts []Account
			for i := range 5 {
				name := fmt.Sprintf("p%d", i)
				accts = append(accts, Account{Name: name, Slot: i + 1, Dir: filepath.Join(root, name)})
				f.writeCache(name, "live", cached{At: now.Add(-time.Minute), Limits: &Limits{FiveHour: window(20, now.Add(time.Hour)),
					Weekly: window(30, now.Add(48*time.Hour)), Scoped: []Scoped{{Name: "Fable", Window: window(10, now.Add(48*time.Hour))}}}})
				if limit {
					if err := SaveLimit(f.StateDir, name, "five_hour", now.Add(time.Hour), now); err != nil {
						b.Fatal(err)
					}
				}
			}
			for b.Loop() {
				f.FetchAll(context.Background(), accts, Options{Source: "cache"}, nil)
			}
		})
	}
}

// A limit record laid over a lookup says so: in --json (limited) and on
// the -v line, which otherwise names the lookup's source for a 100% it
// never reported.
func TestLimitOverALookupSaysSo(t *testing.T) {
	e := newEnv(t)
	now := time.Now().Truncate(time.Second)
	e.f.Now = func() time.Time { return now }
	a := e.account("a", "ok-new-format.json", future)
	if err := SaveLimit(e.f.StateDir, "a", "five_hour", now.Add(2*time.Hour), now.Add(-5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	r := e.f.Fetch(context.Background(), a, Options{NoCache: true})
	if r.Source != "live" || !r.Limited.Equal(now.Add(-5*time.Minute)) {
		t.Fatalf("%+v", r)
	}
	out := fmt.Sprint(Table([]*Report{&r}, RenderOptions{Verbose: true, Now: now}))
	if !strings.Contains(out, "5-hour full: Claude Code reported the usage limit 5m ago") {
		t.Fatalf("%s", out)
	}
	if data, _ := json.Marshal(r); !strings.Contains(string(data), `"limited":`) {
		t.Fatalf("json: %s", data)
	}
}
