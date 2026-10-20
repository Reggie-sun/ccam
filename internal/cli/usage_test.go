package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
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
		cfg.Profiles = append(cfg.Profiles, config.Profile{Name: name, Slot: slot, Share: config.Share{History: "team"}})
		slot++
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
	input, _ := os.ReadFile("../usage/testdata/statusline-input.json")
	t.Setenv("CCAM_PROFILE", "work")
	out, _, code := c.run(string(input), "statusline")
	if code != 0 || strings.TrimSpace(out) != "[work] · Opus · 5h 24% · 7d 41%" && strings.TrimSpace(out) != "[work] · Opus · 5h 23% · 7d 41%" {
		t.Fatalf("%q", out)
	}
	out, _, _ = c.run("", "usage", "--source", "cache")
	if !strings.Contains(out, "41%") || !strings.Contains(out, "(cached)") {
		t.Fatalf("cached usage:\n%s", out)
	}
	out, _, _ = c.run(string(input), "statusline", "--wrap", "echo mine")
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
	if !strings.Contains(cmd, "statusline --wrap my-line") || doc["model"] != "opus" {
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
	cfg.ClaudeBin, _ = exec.LookPath("true") // /bin/true does not exist on macOS
	cfg.Save(c.p.ConfigFile())
	out, errOut, code := c.run("", "--dry-run", "next", "--auto")
	if code != 0 || !strings.Contains(out, c.p.ProfileDir("side")) {
		t.Fatalf("%d out=%s err=%s", code, out, errOut)
	}
}
