package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/migrate"
)

// While a change of groups that stopped halfway is noted, neither the CLI
// nor the dashboard starts another; doctor --fix finishes it, and then
// both do.
func TestNoGroupChangeOverAStoppedOne(t *testing.T) {
	c := newTestCLI(t, nil)
	e, cfg := c.env()
	next := cfg.Clone()
	if err := next.LeaveGroup("side", "team"); err != nil {
		t.Fatal(err)
	}
	ch, err := group.Build(c.p, cfg, next, group.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil { // then ccam stopped before the config was saved
		t.Fatal(err)
	}
	note, err := os.ReadFile(group.NotePath(c.p))
	if err != nil {
		t.Fatal(err)
	}
	ch.Abort() // lets go of the lock, as a stopped ccam does; the note stays
	os.WriteFile(group.NotePath(c.p), note, 0o600)
	b := backend{e: e}
	addMemory := func(cfg *config.Config) error { return cfg.SetGroupItems("team", []string{"history", "memory"}) }
	if _, err := b.PlanGroups(addMemory); err == nil || !strings.Contains(err.Error(), "finish it with `ccam doctor --fix` first") {
		t.Fatalf("dashboard: %v", err)
	}
	if _, errOut, code := c.run("", "group", "leave", "team", "work", "--yes"); code == 0 || !strings.Contains(errOut, "finish it with `ccam doctor --fix` first") {
		t.Fatalf("cli: %d %s", code, errOut)
	}
	// The note names the accounts: a renamed one would be skipped.
	if _, errOut, code := c.run("", "rename", "side", "side2"); code == 0 || !strings.Contains(errOut, "finish it with `ccam doctor --fix` first") {
		t.Fatalf("rename: %d %s", code, errOut)
	}
	if _, err := b.Rename("side", "side2"); err == nil {
		t.Fatal("dashboard rename went ahead")
	}
	if _, errOut, code := c.run("", "remove", "side", "--yes"); code == 0 || !strings.Contains(errOut, "finish it with `ccam doctor --fix` first") {
		t.Fatalf("remove: %d %s", code, errOut)
	}
	out, _, _ := c.run("", "doctor", "--fix")
	if !strings.Contains(out, "finished the last group change (side)") {
		t.Fatalf("doctor --fix:\n%s", out)
	}
	if g, _ := c.cfg().FindGroup("team"); !slices.Contains(g.Members, "side") {
		t.Fatal("the config that was never saved changed")
	}
	if _, err := b.PlanGroups(addMemory); err != nil {
		t.Fatalf("after doctor --fix: %v", err)
	}
	c.ok("group", "leave", "team", "work", "--yes")
}

// While another ccam carries a change of groups out (it holds the lock and
// its note is there), commands say so instead of calling it stopped.
func TestNoGroupChangeWhileOneRuns(t *testing.T) {
	c := newTestCLI(t, nil)
	e, cfg := c.env()
	next := cfg.Clone()
	if err := next.LeaveGroup("side", "team"); err != nil {
		t.Fatal(err)
	}
	ch, err := group.Build(c.p, cfg, next, group.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	defer ch.Abort()
	const want = "another ccam process is changing the groups right now; try again when it is done"
	b := backend{e: e}
	if _, err := b.PlanGroups(func(cfg *config.Config) error { return cfg.SetGroupItems("team", []string{"history", "memory"}) }); err == nil || err.Error() != want {
		t.Errorf("dashboard: %v", err)
	}
	for _, args := range [][]string{{"group", "leave", "team", "work", "--yes"}, {"rename", "side", "side2"}, {"remove", "side", "--yes"}} {
		if _, errOut, code := c.run("", args...); code == 0 || !strings.Contains(errOut, want) || strings.Contains(errOut, "stopped halfway") {
			t.Errorf("%v: %d %s", args, code, errOut)
		}
	}
	if out, _, _ := c.run("", "doctor"); !strings.Contains(out, "in progress right now") || strings.Contains(out, "stopped halfway") {
		t.Errorf("doctor:\n%s", out)
	}
}

// connectorGaps names, for each account, the claude.ai connectors other
// accounts used and it did not.
func TestConnectorGaps(t *testing.T) {
	got := connectorGaps(map[string][]string{"c1": {"Claude Docs"}, "cj": {"Claude Docs", "Gmail", "Google Drive"}, "c2": nil})
	var lines []string
	for _, is := range got {
		lines = append(lines, is.Account+": "+is.Problem)
	}
	all := strings.Join(lines, "\n")
	if len(got) != 2 || !strings.Contains(all, "c1: has not used claude.ai connectors other accounts have: Gmail (cj), Google Drive (cj)") ||
		!strings.Contains(all, "c2: has not used claude.ai connectors other accounts have: Claude Docs (c1, cj), Gmail (cj), Google Drive (cj)") {
		t.Fatalf("%s", all)
	}
}

// The migrate hint names the group the accounts share MCP through.
func TestMigrateMCPCommand(t *testing.T) {
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	if got := migrate.MCPCommand(cfg); got != "ccam migrate --mcp" {
		t.Fatalf("no group: %s", got)
	}
	if err := cfg.AddGroup("team", []string{"mcp"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	if got := migrate.MCPCommand(cfg); got != "ccam migrate --mcp --group team" {
		t.Fatalf("team: %s", got)
	}
	if err := cfg.AddGroup("side", []string{"mcp"}, nil); err != nil {
		t.Fatal(err)
	}
	if got := migrate.MCPCommand(cfg); got != "ccam migrate --mcp --group <team|side>" {
		t.Fatalf("two groups: %s", got)
	}
}

// A plugin's MCP server comes with its plugin: ext offers `ccam mcp import`
// for a server only in .claude.json, not for one of a plugin of the
// account's own.
func TestExtLeavesPluginServersOutOfImport(t *testing.T) {
	c := newTestCLI(t, nil)
	cfg := c.cfg()
	if err := cfg.AddGroup("m", []string{"mcp"}, []string{"personal", "work"}); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Save(c.p.ConfigFile()); err != nil {
		t.Fatal(err)
	}
	dir := c.p.ProfileDir("personal")
	install := filepath.Join(dir, "plugins", "cache", "mk", "demo", "1.0")
	if err := os.MkdirAll(install, 0o700); err != nil {
		t.Fatal(err)
	}
	ip, _ := json.Marshal(install)
	for path, data := range map[string]string{
		filepath.Join(dir, ".claude.json"):                      `{"mcpServers":{"x":{"command":"x"}}}`,
		filepath.Join(install, ".mcp.json"):                     `{"srv":{"command":"s"}}`,
		filepath.Join(dir, "plugins", "installed_plugins.json"): `{"version":2,"plugins":{"demo@mk":[{"scope":"user","installPath":` + string(ip) + `}]}}`,
		filepath.Join(dir, "settings.json"):                     `{"enabledPlugins":{"demo@mk":true}}`,
	} {
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	out := c.ok("ext")
	if !strings.Contains(out, "plugin:demo:srv") || !strings.Contains(out, "MCP server x is only in this account's .claude.json") {
		t.Fatalf("ext:\n%s", out)
	}
	if strings.Contains(out, "plugin:demo:srv is only") || strings.Contains(out, "import personal plugin:") {
		t.Errorf("ext offers to import a plugin's server:\n%s", out)
	}
}
