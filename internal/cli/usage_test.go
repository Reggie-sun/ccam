package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// testCLI is an isolated CCAM_HOME with profiles whose tokens tell the fake
// usage server which percentage to report (-1: the lookup fails).
type testCLI struct {
	t *testing.T
	p paths.Paths
}

func newTestCLI(t *testing.T, fiveHour map[string]int) *testCLI {
	t.Setenv("CCAM_HOME", t.TempDir())
	c := &testCLI{t: t, p: paths.FromEnv()}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer unknown" {
			w.WriteHeader(500)
			return
		}
		var pct int
		fmt.Sscanf(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer pct-"), "%d", &pct)
		fmt.Fprintf(w, `{"limits":[{"kind":"session","percent":%d,"resets_at":"2030-01-01T00:00:00Z"},{"kind":"weekly_all","percent":10,"resets_at":"2030-01-05T00:00:00Z"}]}`, pct)
	}))
	t.Cleanup(srv.Close)
	orig := newFetcher
	newFetcher = func(dir string) *usage.Fetcher {
		f := usage.NewFetcher(dir)
		f.URL = srv.URL
		return f
	}
	t.Cleanup(func() { newFetcher = orig })

	cfg := &config.Config{}
	slot := 1
	for _, name := range []string{"personal", "work", "side"} {
		pct, ok := fiveHour[name]
		dir := c.p.ProfileDir(name)
		os.MkdirAll(dir, 0o700)
		if ok {
			token := fmt.Sprintf("pct-%d", pct)
			if pct < 0 {
				token = "unknown"
			}
			cred := fmt.Sprintf(`{"claudeAiOauth":{"accessToken":%q,"expiresAt":%d}}`, token, time.Now().Add(time.Hour).UnixMilli())
			os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(cred), 0o600)
		}
		cfg.Profiles = append(cfg.Profiles, config.Profile{Name: name, Slot: slot})
		slot++
	}
	if err := cfg.AddGroup("team", []string{"history"}, []string{"personal", "work", "side"}); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Save(c.p.ConfigFile()); err != nil {
		t.Fatal(err)
	}
	return c
}

func (c *testCLI) run(stdin string, args ...string) (string, string, int) {
	var out, errb bytes.Buffer
	code := ExecuteWith(args, c.p, strings.NewReader(stdin), &out, &errb)
	return out.String(), errb.String(), code
}

func TestUsageTableAndJSON(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 72, "work": 8, "side": 100})
	out, errOut, code := c.run("", "usage")
	if code != 0 {
		t.Fatalf("%d %s", code, errOut)
	}
	for _, s := range []string{"72%", "8%", "100%", "Most headroom: work"} {
		if !strings.Contains(out, s) {
			t.Errorf("missing %q:\n%s", s, out)
		}
	}
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("escape codes in piped output: %q", out)
	}
	out, _, _ = c.run("", "u", "--json")
	var reports []usage.Report
	if err := json.Unmarshal([]byte(out), &reports); err != nil || len(reports) != 3 || reports[2].Limits.FiveHour.Percent != 100 {
		t.Fatalf("json: %v %s", err, out)
	}
}

// The advice line follows next --auto: an account whose weekly window is
// full is not recommended however low its 5-hour usage is.
func TestUsageAdviceSkipsWeeklyFull(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 5, "work": 40})
	cache := filepath.Join(c.p.State, "usage", "personal.json")
	os.MkdirAll(filepath.Dir(cache), 0o700)
	reset := time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339)
	os.WriteFile(cache, []byte(`{"at":"`+time.Now().UTC().Format(time.RFC3339)+`","limits":{"five_hour":{"percent":5,"known":true,"resets_at":"`+reset+`"},"weekly":{"percent":100,"known":true,"resets_at":"`+reset+`"}}}`), 0o600)
	out, _, _ := c.run("", "usage")
	if !strings.Contains(out, "Most headroom: work") {
		t.Fatalf("%s", out)
	}
}

func TestUsageCacheModeSaysNotLoggedIn(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10})
	out, _, _ := c.run("", "usage", "--source", "cache")
	if !strings.Contains(out, "not logged in → ccam login work") || strings.Contains(out, "lookup failed: no saved usage") && strings.Contains(out, " work ") && strings.Count(out, "lookup failed") > 1 {
		t.Fatalf("%s", out)
	}
}

func TestUsageReportsLoggedOut(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10})
	out, _, _ := c.run("", "usage")
	if !strings.Contains(out, "not logged in → ccam login work") {
		t.Fatalf("%s", out)
	}
}

func TestStatuslineSuggestsTheNextAccountNearTheLimit(t *testing.T) {
	c := newTestCLI(t, nil)
	t.Setenv("CCAM_PROFILE", "work")
	input := `{"model":{"display_name":"Opus"},"rate_limits":{"five_hour":{"used_percentage":%v,"resets_at":4102444800},"seven_day":{"used_percentage":%v,"resets_at":4102444800}}}`
	for _, tc := range []struct {
		five, week float64
		hint       bool
	}{{89, 10, false}, {90, 10, true}, {20, 95, true}} {
		out, _, _ := c.run(fmt.Sprintf(input, tc.five, tc.week), "statusline")
		if got := strings.Contains(out, "→ ccam next --auto --continue"); got != tc.hint {
			t.Errorf("5h %v 7d %v: %q", tc.five, tc.week, out)
		}
	}
	// With one account there is nothing to switch to.
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.Profiles = cfg.Profiles[:1]
	cfg.Save(c.p.ConfigFile())
	if out, _, _ := c.run(fmt.Sprintf(input, 99, 99), "statusline"); strings.Contains(out, "next") {
		t.Errorf("single account: %q", out)
	}
}

func TestStatuslineCollectsAndPrints(t *testing.T) {
	c := newTestCLI(t, nil)
	raw, _ := os.ReadFile("../usage/testdata/statusline-input.json")
	// The fixture's resets_at are fixed dates; move them ahead of the clock,
	// or usage --source cache counts the saved windows as reset.
	input := strings.NewReplacer(
		"1792346400", fmt.Sprint(time.Now().Add(2*time.Hour).Unix()),
		"1792634400", fmt.Sprint(time.Now().Add(72*time.Hour).Unix()),
	).Replace(string(raw))
	if input == string(raw) {
		t.Fatal("the fixture's resets_at changed; update this test")
	}
	t.Setenv("CCAM_PROFILE", "work")
	out, _, code := c.run(input, "statusline")
	if code != 0 || strings.TrimSpace(out) != "[work] · Opus · 5h 24% · 7d 41%" && strings.TrimSpace(out) != "[work] · Opus · 5h 23% · 7d 41%" {
		t.Fatalf("%q", out)
	}
	out, _, _ = c.run("", "usage", "--source", "cache")
	if !strings.Contains(out, "41%") || !strings.Contains(out, "(cached)") {
		t.Fatalf("cached usage:\n%s", out)
	}
	out, _, _ = c.run(input, "statusline", "--wrap", "echo mine")
	if strings.TrimSpace(out) != "[work] mine" {
		t.Fatalf("wrap: %q", out)
	}
}

func TestStatuslineInstallKeepsSymlinkAndWraps(t *testing.T) {
	c := newTestCLI(t, nil)
	group := filepath.Join(c.p.GroupDir("default"), "settings.json")
	os.MkdirAll(filepath.Dir(group), 0o700)
	os.WriteFile(group, []byte(`{"model":"opus","statusLine":{"type":"command","command":"my-line"}}`), 0o600)
	for _, n := range []string{"personal", "work"} {
		os.Symlink(group, filepath.Join(c.p.ProfileDir(n), "settings.json"))
	}
	out, errOut, code := c.run("", "statusline", "install")
	if code != 0 {
		t.Fatalf("%s %s", out, errOut)
	}
	if strings.Count(out, group) != 1 {
		t.Fatalf("shared file should be written once:\n%s", out)
	}
	if info, _ := os.Lstat(filepath.Join(c.p.ProfileDir("work"), "settings.json")); info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("symlink replaced")
	}
	var doc map[string]any
	data, _ := os.ReadFile(group)
	json.Unmarshal(data, &doc)
	cmd := doc["statusLine"].(map[string]any)["command"].(string)
	wrapped := "statusline --wrap my-line"
	if runtime.GOOS == "windows" {
		wrapped = "statusline --wrap-base64 bXktbGluZQ=="
	}
	if !strings.Contains(cmd, wrapped) || doc["model"] != "opus" {
		t.Fatalf("settings: %s", data)
	}
	out, _, _ = c.run("", "statusline", "install")
	if !strings.Contains(out, "already set") {
		t.Fatalf("second install: %s", out)
	}
}

func TestNextAutoPicksLowest(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 72, "work": 100, "side": 30})
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin = stubClaude
	cfg.Save(c.p.ConfigFile())
	out, errOut, code := c.run("", "--dry-run", "next", "--auto")
	if code != 0 || !strings.Contains(out, c.p.ProfileDir("side")) || !strings.Contains(errOut, "(30% of its 5 hours used, the least)") {
		t.Fatalf("%d out=%s err=%s", code, out, errOut)
	}
	// balance: every week is 10% used and resets at the same time, so the
	// reason is about the week.
	out, errOut, code = c.run("", "--dry-run", "next", "--auto", "--strategy", "balance")
	if code != 0 || !strings.Contains(errOut, "90% of its week is left") {
		t.Fatalf("%d out=%s err=%s", code, out, errOut)
	}
	if _, errOut, code := c.run("", "next", "--auto", "--strategy", "random"); code == 0 || !strings.Contains(errOut, `"balance"`) {
		t.Fatalf("%d %s", code, errOut)
	}
}

// Scripts get JSON even before the first account is added.
func TestJSONWithNoAccounts(t *testing.T) {
	c := newTestCLI(t, nil)
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.Profiles = nil
	cfg.Groups = nil
	cfg.Save(c.p.ConfigFile())
	for _, args := range [][]string{{"usage", "--json"}, {"u", "--json"}, {"ext", "--json"}, {"list", "--json"}} {
		out, errOut, code := c.run("", args...)
		var v any
		if code != 0 || json.Unmarshal([]byte(out), &v) != nil {
			t.Errorf("%v: %d %q %q", args, code, out, errOut)
		}
	}
	out, _, _ := c.run("", "ext", "--json")
	var rep struct{ Accounts []any }
	if json.Unmarshal([]byte(out), &rep) != nil || rep.Accounts == nil {
		t.Errorf("ext: accounts must be [], not null: %s", out)
	}
	if out, _, _ := c.run("", "usage"); !strings.Contains(out, "No accounts yet") {
		t.Errorf("table: %q", out)
	}
}

// The live view redraws in place by moving the cursor up over the last
// draw. A draw that does not fit the screen cannot be moved over, so it
// waits, and the table is drawn once more at the end.
func TestLiveTableStaysOnScreen(t *testing.T) {
	var out bytes.Buffer
	live := &liveTable{w: &out, tty: true, rows: 8}
	table := func(n int) []string {
		var lines []string
		for i := range n {
			lines = append(lines, fmt.Sprintf("line %d", i))
		}
		return lines
	}
	for _, n := range []int{5, 7, 8, 11, 14, 17, 20} {
		live.draw(table(n), false)
	}
	if !live.held {
		t.Fatal("not held")
	}
	live.draw(table(20), true)
	ups := regexp.MustCompile(`\x1b\[(\d+)A`).FindAllStringSubmatch(out.String(), -1)
	for _, m := range ups {
		if n, _ := strconv.Atoi(m[1]); n >= live.rows {
			t.Errorf("moved up %d rows on a %d-row screen", n, live.rows)
		}
	}
	last := out.String()[strings.LastIndex(out.String(), "\x1b["):]
	if len(ups) != 2 || !strings.Contains(last, "line 19") || strings.Count(out.String(), "line 19") != 1 {
		t.Fatalf("%d moves:\n%q", len(ups), out.String())
	}
	// Without a terminal, the table is printed once, as it is.
	out.Reset()
	plain := &liveTable{w: &out}
	plain.draw(table(3), true)
	if out.String() != "line 0\nline 1\nline 2\n" {
		t.Fatalf("%q", out.String())
	}
}

// The advice line names what next --auto would pick and, unless that is the
// one with most headroom, why.
func TestUsageAdviceFollowsTheStrategy(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 5, "work": 40})
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.AutoStrategy = usage.Expiring
	cfg.Save(c.p.ConfigFile())
	// work's 5-hour window resets first.
	cache := filepath.Join(c.p.State, "usage", "work.json")
	os.MkdirAll(filepath.Dir(cache), 0o700)
	now := time.Now().UTC()
	os.WriteFile(cache, fmt.Appendf(nil, `{"at":%q,"limits":{"five_hour":{"percent":40,"known":true,"resets_at":%q}}}`,
		now.Format(time.RFC3339), now.Add(20*time.Minute).Format(time.RFC3339)), 0o600)
	out, _, _ := c.run("", "usage")
	if strings.Contains(out, "Most headroom") || !strings.Contains(out, "→ next --auto picks work: its 5-hour window resets first") {
		t.Fatalf("%s", out)
	}
}
