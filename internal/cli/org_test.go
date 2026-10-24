package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/tui"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// logIn writes the oauthAccount Claude Code keeps in .claude.json.
func logIn(t *testing.T, dir, account, org, orgType string) {
	t.Helper()
	data := fmt.Sprintf(`{"oauthAccount":{"accountUuid":%q,"emailAddress":"%s@example.com","organizationUuid":%q,"organizationName":"Org %s","organizationType":%q}}`,
		account, account, org, org, orgType)
	if err := os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
}

// Automatic choices stay in the organization of the account they leave.
func TestAutoStaysInTheOrganization(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 72, "work": 20, "side": 10})
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin = stubClaude
	cfg.Save(c.p.ConfigFile())
	logIn(t, c.p.ProfileDir("personal"), "u1", "acme", "claude_team")
	logIn(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team")
	logIn(t, c.p.ProfileDir("side"), "u3", "home", "claude_max")

	// side has the most room, but it is a personal plan.
	out, errOut, code := c.run("", "--dry-run", "next", "--auto")
	if code != 0 || !strings.Contains(out, c.p.ProfileDir("work")) {
		t.Fatalf("%d out=%s err=%s", code, out, errOut)
	}
	out, _, code = c.run("", "--dry-run", "next", "--auto", "--cross-org")
	if code != 0 || !strings.Contains(out, c.p.ProfileDir("side")) {
		t.Fatalf("--cross-org: %d %s", code, out)
	}
	if out, _, _ := c.run("", "--dry-run", "pick"); strings.TrimSpace(out) != "work" {
		t.Fatalf("pick from personal's folder: %q", out)
	}
	// work is the same Claude account as personal in another folder: it
	// shares personal's limits, so it never takes over from it.
	logIn(t, c.p.ProfileDir("work"), "u1", "acme", "claude_team")
	out, errOut, code = c.run("", "--dry-run", "next", "--auto")
	if code != 0 || strings.Contains(out, c.p.ProfileDir("work")) {
		t.Fatalf("%d out=%s err=%s", code, out, errOut)
	}
	// With nobody left in the organization, the message says how to look further.
	os.Remove(filepath.Join(c.p.ProfileDir("personal"), ".credentials.json"))
	if _, errOut, code := c.run("", "--dry-run", "pick"); code != exitNoRoom || !strings.Contains(errOut, "Org acme") || !strings.Contains(errOut, "--cross-org") {
		t.Fatalf("%d %s", code, errOut)
	}
}

// --fallback tries one account after another, and each must fit every
// account tried before it, not only the last one: a second folder of a
// Claude account that reached its limit shares that limit, and an account
// whose organization is unknown (it fits any) opens no other organization.
func TestFallbackFitsEveryTriedAccount(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 100, "work": 100, "side": 10})
	e, cfg := c.env()
	tried := map[string]bool{"personal": true, "work": true}
	next := func() string {
		if p, _ := fallbackTarget(e, cfg, &cfg.Profiles[1], tried, []string{"-p", "x"}, usage.Policy{}); p != nil {
			return p.Name
		}
		return ""
	}
	logIn(t, c.p.ProfileDir("personal"), "u1", "acme", "claude_team")
	logIn(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team")
	logIn(t, c.p.ProfileDir("side"), "u1", "acme", "claude_team")
	if got := next(); got != "" {
		t.Fatalf("after personal and work: %s, personal's Claude account", got)
	}
	os.Remove(filepath.Join(c.p.ProfileDir("work"), ".claude.json"))
	logIn(t, c.p.ProfileDir("side"), "u3", "home", "claude_max")
	if got := next(); got != "" {
		t.Fatalf("after personal and work: %s, another organization than personal's", got)
	}
	logIn(t, c.p.ProfileDir("side"), "u3", "acme", "claude_team")
	if got := next(); got != "side" {
		t.Fatalf("another seat of personal's organization: %q", got)
	}
}

// The dashboard's n launches the account with the most room among those
// that may take over from the selected one (Backend.Takeover), by the rules
// of next --auto: the seats of its Team organization, not another
// organization's account, and not the seat Claude Code reported full.
func TestDashboardMostRoomAmongSeats(t *testing.T) {
	c := newTestCLI(t, nil)
	cfg := &config.Config{Default: "c1"}
	for i, s := range []struct {
		name, account, org string
		fiveHour           int
	}{{"c1", "u1", "acme", 5}, {"c2", "u2", "acme", 40}, {"c3", "u3", "acme", 20}, {"c4", "u4", "acme", 60}, {"cj", "u5", "acme", 35}, {"other", "u6", "home", 1}} {
		dir := c.p.ProfileDir(s.name)
		os.MkdirAll(dir, 0o700)
		logIn(t, dir, s.account, s.org, "claude_team")
		os.WriteFile(filepath.Join(dir, ".credentials.json"), fmt.Appendf(nil, `{"claudeAiOauth":{"accessToken":"pct-%d","expiresAt":%d}}`,
			s.fiveHour, time.Now().Add(time.Hour).UnixMilli()), 0o600)
		cfg.Profiles = append(cfg.Profiles, config.Profile{Name: s.name, Slot: i + 1})
	}
	if err := cfg.Save(c.p.ConfigFile()); err != nil {
		t.Fatal(err)
	}
	if err := usage.SaveLimit(filepath.Join(c.p.State, "usage"), "c1", "five_hour", time.Now().Add(2*time.Hour), time.Now()); err != nil {
		t.Fatal(err)
	}
	e, cfg := c.env()
	d := tui.NewDashboard(backend{e: e}, cfg, "", "dev")
	m := runCmds(d, d.Init())
	m, cmd := m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	runCmds(m, cmd)
	if r := d.Result(); r.Action != tui.Launch || r.Profile != "c3" {
		t.Fatalf("n on c1: %+v", r)
	}
}

// runCmds runs cmd and the commands its messages lead to, as the program
// would, until Quit.
func runCmds(m tea.Model, cmd tea.Cmd) tea.Model {
	if cmd == nil {
		return m
	}
	switch msg := cmd().(type) {
	case nil, tea.QuitMsg:
	case tea.BatchMsg:
		for _, c := range msg {
			m = runCmds(m, c)
		}
	default:
		var next tea.Cmd
		m, next = m.Update(msg)
		m = runCmds(m, next)
	}
	return m
}

func TestIdentityPools(t *testing.T) {
	team := &config.Identity{AccountUUID: "a", OrgUUID: "o1", OrgType: "claude_team"}
	team2 := &config.Identity{AccountUUID: "b", OrgUUID: "o2", OrgType: "claude_enterprise"}
	max := &config.Identity{AccountUUID: "c", OrgUUID: "o3", OrgType: "claude_max"}
	pro := &config.Identity{AccountUUID: "d", OrgUUID: "o4", OrgType: "claude_pro"}
	if config.SamePool(team, team2) || config.SamePool(team, max) || !config.SamePool(max, pro) || !config.SamePool(nil, team) {
		t.Fatal("pools")
	}
	if !config.SameAccount(team, &config.Identity{AccountUUID: "a", OrgUUID: "o1"}) || config.SameAccount(team, &config.Identity{AccountUUID: "a", OrgUUID: "o9"}) || config.SameAccount(nil, nil) {
		t.Fatal("same account")
	}
}

func TestListShowsOrganizationsWhenThereAreSeveral(t *testing.T) {
	c := newTestCLI(t, nil)
	logIn(t, c.p.ProfileDir("personal"), "u1", "acme", "claude_team")
	logIn(t, c.p.ProfileDir("work"), "u2", "acme", "claude_team")
	if out, _, _ := c.run("", "list"); strings.Contains(out, "Organization") {
		t.Fatal(out)
	}
	logIn(t, c.p.ProfileDir("side"), "u3", "home", "claude_max")
	out, _, _ := c.run("", "list")
	if !strings.Contains(out, "Organization") || !strings.Contains(out, "Org acme") {
		t.Fatal(out)
	}
	for _, l := range strings.Split(out, "\n") {
		if strings.Contains(l, " side ") && !strings.Contains(l, "personal") {
			t.Fatalf("side has a personal plan: %q", l)
		}
	}
	if out, _, _ := c.run("", "list", "--json"); !strings.Contains(out, `"plan": "Claude Max"`) {
		t.Fatal(out)
	}
}
