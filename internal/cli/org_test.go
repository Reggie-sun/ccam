package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
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
	cfg.ClaudeBin, _ = exec.LookPath("true")
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
