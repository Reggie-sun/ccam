package group

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/paths"
)

func setup(t *testing.T, names ...string) (paths.Paths, *config.Config) {
	t.Helper()
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{}
	for i, n := range names {
		cfg.Profiles = append(cfg.Profiles, config.Profile{Name: n, Slot: i + 1})
		os.MkdirAll(p.ProfileDir(n), 0o700)
		os.WriteFile(filepath.Join(p.ProfileDir(n), "CLAUDE.md"), []byte(n), 0o600)
	}
	return p, cfg
}

func apply(t *testing.T, p paths.Paths, old *config.Config, edit func(*config.Config) error, opt Options) (*Change, *config.Config) {
	t.Helper()
	next := old.Clone()
	if err := edit(next); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, old, next, opt)
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Check(func(string) bool { return false }); err != nil {
		t.Fatal(err)
	}
	if err := ch.Apply(); err != nil {
		t.Fatal(err)
	}
	return ch, next
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestSeedOrder(t *testing.T) {
	p, cfg := setup(t, "a", "b", "c")
	cfg.Default = "b"
	add := func(c *config.Config) error { return c.AddGroup("g", []string{"memory"}, []string{"a", "b", "c"}) }
	// The default account goes first without --from.
	ch, _ := apply(t, p, cfg, add, Options{})
	if got := strings.Join(ch.Names(), ","); got != "b,a,c" {
		t.Fatalf("order %s", got)
	}
	if read(t, filepath.Join(p.ProfileDir("c"), "CLAUDE.md")) != "b" {
		t.Fatal("b's copy should win")
	}

	p, cfg = setup(t, "a", "b", "c")
	cfg.Default = "b"
	ch, _ = apply(t, p, cfg, add, Options{Seed: "c"})
	if got := strings.Join(ch.Names(), ","); got != "c,b,a" || len(ch.Conflicts()) != 2 {
		t.Fatalf("order %s, %d conflicts", got, len(ch.Conflicts()))
	}
	if read(t, filepath.Join(p.ProfileDir("a"), "CLAUDE.md")) != "c" {
		t.Fatal("--from wins")
	}
}

func TestUnchangedAccountsAreLeftOut(t *testing.T) {
	p, cfg := setup(t, "a", "b")
	if err := cfg.AddGroup("g", []string{"memory"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	ch, _ := apply(t, p, cfg, func(c *config.Config) error { _, err := c.JoinGroup("b", "g", false); return err }, Options{})
	if got := strings.Join(ch.Names(), ","); got != "b" {
		t.Fatalf("got %s", got)
	}
	var buf bytes.Buffer
	ch.Describe(&buf)
	if !strings.Contains(buf.String(), "b: +memory (group g)") {
		t.Fatal(buf.String())
	}
}

// Moving MCP to a new group brings the servers along; leaving MCP sharing
// removes the copies ccam made in .claude.json.
func TestMCPMovesAndLeaves(t *testing.T) {
	p, cfg := setup(t, "a")
	if err := cfg.AddGroup("g", []string{"mcp"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(p.GroupDir("g"), 0o700)
	os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"}}}`), 0o600)
	if err := cfg.AddGroup("h", nil, nil); err != nil {
		t.Fatal(err)
	}
	ch, cfg := apply(t, p, cfg, func(c *config.Config) error {
		if err := c.SetGroupItems("g", nil); err != nil {
			return err
		}
		if err := c.SetGroupItems("h", []string{"mcp"}); err != nil {
			return err
		}
		_, err := c.JoinGroup("a", "h", false)
		return err
	}, Options{})
	if ch.Accounts[0].SeedMCP == nil || !strings.Contains(read(t, mcp.Path(p.GroupDir("h"))), `"fs"`) {
		t.Fatal("servers not carried to group h")
	}

	// A copy Sync made for the account goes when MCP sharing stops.
	dir := p.ProfileDir("a")
	if err := mcp.Sync(dir, p.GroupDir("h"), launcher.MCPStateDir(p, "a")); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(read(t, filepath.Join(dir, ".claude.json")), `"fs"`) {
		t.Fatal("sync did not copy")
	}
	ch, _ = apply(t, p, cfg, func(c *config.Config) error { return c.LeaveGroup("a", "h") }, Options{})
	var buf bytes.Buffer
	ch.Describe(&buf)
	if !strings.Contains(buf.String(), "MCP servers ccam copied into .claude.json: fs") {
		t.Fatal(buf.String())
	}
	if strings.Contains(read(t, filepath.Join(dir, ".claude.json")), `"fs"`) || mcp.HasCopy(launcher.MCPStateDir(p, "a")) {
		t.Fatal("the copied server must go")
	}
}

func TestCheckRefusesRunningAccounts(t *testing.T) {
	p, cfg := setup(t, "a", "b")
	next := cfg.Clone()
	if err := next.AddGroup("g", []string{"memory"}, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, cfg, next, Options{})
	if err != nil {
		t.Fatal(err)
	}
	err = ch.Check(func(n string) bool { return true })
	if err == nil || !strings.Contains(err.Error(), "a, b are running") {
		t.Fatal(err)
	}
}

// Three accounts that bring their own plugin lists into a group whose
// plugins folder has none: every list ends up in the group's.
func TestEveryPluginListIsMerged(t *testing.T) {
	p, cfg := setup(t, "a", "b", "c")
	os.MkdirAll(filepath.Join(p.GroupDir("g"), "plugins"), 0o700)
	for _, n := range []string{"a", "b", "c"} {
		dir := filepath.Join(p.ProfileDir(n), "plugins")
		os.MkdirAll(dir, 0o700)
		os.WriteFile(filepath.Join(dir, "installed_plugins.json"), []byte(`{"version":2,"plugins":{"`+n+`@m":[{"scope":"user"}]}}`), 0o600)
	}
	apply(t, p, cfg, func(c *config.Config) error { return c.AddGroup("g", []string{"extensions"}, []string{"a", "b", "c"}) }, Options{})
	got := read(t, filepath.Join(p.GroupDir("g"), "plugins", "installed_plugins.json"))
	for _, n := range []string{"a@m", "b@m", "c@m"} {
		if !strings.Contains(got, n) {
			t.Errorf("%s missing from %s", n, got)
		}
	}
}

// Moving to a group whose folder is there but empty brings the old group's
// content along, as moving to a missing one does.
func TestEmptyGroupFolderIsSeeded(t *testing.T) {
	p, cfg := setup(t, "a")
	skill := filepath.Join(p.ProfileDir("a"), "skills", "s", "SKILL.md")
	os.MkdirAll(filepath.Dir(skill), 0o700)
	os.WriteFile(skill, []byte("skill"), 0o600)
	_, cfg = apply(t, p, cfg, func(c *config.Config) error {
		if err := c.AddGroup("x", []string{"extensions"}, []string{"a"}); err != nil {
			return err
		}
		return c.AddGroup("y", []string{"extensions"}, nil)
	}, Options{})
	os.MkdirAll(filepath.Join(p.GroupDir("y"), "skills"), 0o700) // there, but empty
	apply(t, p, cfg, func(c *config.Config) error {
		if err := c.LeaveGroup("a", "x"); err != nil {
			return err
		}
		_, err := c.JoinGroup("a", "y", false)
		return err
	}, Options{})
	if got := read(t, skill); got != "skill" {
		t.Fatalf("a sees %q", got)
	}
}

// Moving MCP to another group drops the old group's servers that ccam had
// copied into .claude.json; under the flag strategy they would win.
func TestMCPMoveDropsCopiedServers(t *testing.T) {
	p, cfg := setup(t, "a")
	if err := cfg.AddGroup("g", []string{"mcp"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	if err := cfg.AddGroup("h", nil, nil); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(p.GroupDir("g"), 0o700)
	os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"}}}`), 0o600)
	dir := p.ProfileDir("a")
	if err := mcp.Sync(dir, p.GroupDir("g"), launcher.MCPStateDir(p, "a")); err != nil {
		t.Fatal(err)
	}
	apply(t, p, cfg, func(c *config.Config) error {
		if err := c.SetGroupItems("g", nil); err != nil {
			return err
		}
		if err := c.SetGroupItems("h", []string{"mcp"}); err != nil {
			return err
		}
		_, err := c.JoinGroup("a", "h", false)
		return err
	}, Options{})
	if strings.Contains(read(t, filepath.Join(dir, ".claude.json")), `"fs"`) {
		t.Fatal("group g's copied server is still in .claude.json")
	}
}
