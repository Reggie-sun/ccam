package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// seat is an account of seatsHarness: its Claude account, organization and
// saved 5-hour usage.
type seat struct {
	name, account, org string
	fiveHour           int
}

// seatsHarness is the owner's setup: five seats (c1 c2 c3 c4 cj) of one
// Claude Team organization in a config.toml from before groups
// ([profile.share] tables) without [profile.login] records, so ccam reads
// who each folder is from its .claude.json. The usage was saved just now;
// c1's numbers show the most room, but Claude Code reported its 5-hour
// limit (the limit record). extra are more accounts after them.
func seatsHarness(t *testing.T, extra ...seat) *harness {
	t.Helper()
	h := newHarness(t)
	seats := append([]seat{{"c1", "u1", "team", 5}, {"c2", "u2", "team", 40}, {"c3", "u3", "team", 20}, {"c4", "u4", "team", 60}, {"cj", "u5", "team", 35}}, extra...)
	usageDir := filepath.Join(h.home, "state", "usage")
	os.MkdirAll(usageDir, 0o700)
	now := time.Now().UTC()
	at, in2h, in3d := now.Format(time.RFC3339), now.Add(2*time.Hour).Format(time.RFC3339), now.Add(72*time.Hour).Format(time.RFC3339)
	var cfg strings.Builder
	fmt.Fprintf(&cfg, "claude_bin = %q\ndefault = \"c1\"\n", fakeBin)
	for i, s := range seats {
		fmt.Fprintf(&cfg, "\n[[profile]]\nname = %q\nslot = %d\nemail = \"%s@example.com\"\nshortcut = \"cc%s\"\n", s.name, i+1, s.account, s.name)
		cfg.WriteString("\n[profile.share]\nmcp = \"default\"\nsettings = \"default\"\nmemory = \"default\"\nextensions = \"default\"\nhistory = \"default\"\n")
		dir := h.profileDir(s.name)
		os.MkdirAll(dir, 0o700)
		os.WriteFile(filepath.Join(dir, ".claude.json"), fmt.Appendf(nil,
			`{"oauthAccount":{"accountUuid":%q,"emailAddress":"%s@example.com","organizationUuid":%q,"organizationName":"Org %s","organizationType":"claude_team"}}`,
			s.account, s.account, s.org, s.org), 0o600)
		os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"fake","expiresAt":4102444800000}}`), 0o600)
		os.WriteFile(filepath.Join(usageDir, s.name+".json"), fmt.Appendf(nil,
			`{"at":%q,"limits":{"five_hour":{"percent":%d,"resets_at":%q,"known":true},"weekly":{"percent":10,"resets_at":%q,"known":true}}}`,
			at, s.fiveHour, in2h, in3d), 0o600)
	}
	if err := os.WriteFile(filepath.Join(h.home, "config", "config.toml"), []byte(cfg.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(usageDir, "c1.limit.json"), fmt.Appendf(nil, `{"at":%q,"windows":{"five_hour":%q}}`, at, in2h), 0o600)
	return h
}

// Seats of one Team organization stand in for each other in automatic
// choices, as personal plans do: next --auto, pick and --fallback leave the
// seat at its limit and take the seat with the most room.
func TestSeatsOfOneTeamStandInForEachOther(t *testing.T) {
	h := seatsHarness(t)
	cfgFile := filepath.Join(h.home, "config", "config.toml")
	before, _ := os.ReadFile(cfgFile)

	r := h.run("", nil, "ccam", "--dry-run", "next", "--auto")
	if r.code != 0 || !strings.Contains(r.stderr, "→ c1 → c3 (20% of its 5 hours used, the least)") || !strings.Contains(r.stdout, h.profileDir("c3")) {
		t.Fatalf("next --auto: %+v", r)
	}
	h.ok("ccam", "next", "--auto")
	if dir, _ := h.called(); dir != h.profileDir("c3") {
		t.Fatalf("next --auto ran %s", dir)
	}
	var picked struct{ Name, Reason string }
	if err := json.Unmarshal([]byte(h.ok("ccam", "pick", "--json")), &picked); err != nil || picked.Name != "c3" || picked.Reason != "20% of its 5 hours used, the least" {
		t.Fatalf("pick --json: %v %+v", err, picked)
	}
	// c3 was picked a moment ago: the next worker gets the seat after it.
	if out := h.ok("ccam", "pick"); out != "cj\n" {
		t.Fatalf("pick: %q", out)
	}
	runs := filepath.Join(h.home, "runs")
	r = h.run("", []string{"FAKE_RUNS=" + runs, "FAKE_LIMITED=c1", "FAKE_STDOUT=answer"}, "ccam", "--fallback", "c1", "-p", "x")
	if got := strings.Join(runsIn(runs), ","); r.code != 0 || got != "c1,c3" || r.stdout != "answer (c3)\n" {
		t.Fatalf("--fallback ran %s: %+v", got, r)
	}
	// The limit record is what keeps c1 out: without it, its saved 5% is the least.
	os.Remove(filepath.Join(h.home, "state", "usage", "c1.limit.json"))
	if r := h.run("", nil, "ccam", "--dry-run", "next", "--auto"); r.code != 0 || !strings.Contains(r.stdout, h.profileDir("c1")) {
		t.Fatalf("without the record: %+v", r)
	}
	// Choosing left the config as it was: no login records, the old format.
	if after, _ := os.ReadFile(cfgFile); string(after) != string(before) {
		t.Fatalf("config.toml changed:\n%s", after)
	}
}

// The choice never goes to another organization (unless --cross-org) nor
// to a second folder of a Claude account, which shares that account's
// limits; with --fallback, also after other seats were tried first.
func TestSeatsStayInTheirOrganizationAndAccount(t *testing.T) {
	// other has the most room but another organization; dup is c1's Claude
	// account in a second folder.
	h := seatsHarness(t, seat{"other", "u6", "home", 1}, seat{"dup", "u1", "team", 1})
	choose := func(args ...string) string {
		t.Helper()
		r := h.run("", nil, append([]string{"ccam", "--dry-run"}, args...)...)
		if r.code != 0 {
			t.Fatalf("%q: %+v", args, r)
		}
		if args[0] == "pick" {
			return strings.TrimSpace(r.stdout)
		}
		_, to, _ := strings.Cut(r.stderr, "→ c1 → ") // next: → c1 → <account> (why)
		name, _, _ := strings.Cut(to, " ")
		return name
	}
	for args, want := range map[string]string{"next --auto": "c3", "next --auto --cross-org": "other", "pick": "c3", "pick --cross-org": "other"} {
		if got := choose(strings.Fields(args)...); got != want {
			t.Errorf("%s chose %q, want %s", args, got, want)
		}
	}
	runs := filepath.Join(h.home, "runs")
	r := h.run("", []string{"FAKE_RUNS=" + runs, "FAKE_LIMITED=c1,c2,c3,c4,cj", "FAKE_STDOUT=answer"}, "ccam", "--fallback", "c1", "-p", "x")
	if got := strings.Join(runsIn(runs), ","); r.code != 75 || got != "c1,c3,cj,c2,c4" || !strings.Contains(r.stderr, "no other logged-in account has room left") {
		t.Fatalf("--fallback ran %s: %+v", got, r)
	}
}
