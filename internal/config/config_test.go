package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadMissingIsEmpty(t *testing.T) {
	c, err := Load(filepath.Join(t.TempDir(), "none.toml"))
	if err != nil || c.Version != CurrentVersion || len(c.Profiles) != 0 {
		t.Fatalf("got %+v %v", c, err)
	}
}

func TestRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	c := &Config{Version: 1, Default: "personal", MCPStrategy: MCPFlag}
	c.Profiles = append(c.Profiles,
		Profile{Name: "personal", Slot: 1, Email: "me@example.com", Share: DefaultShare()},
		Profile{Name: "work", Slot: 2, Dir: "/x/.claude"},
	)
	if err := c.Save(path); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), "[[profile]]") || !strings.Contains(string(raw), "[profile.share]") {
		t.Fatalf("unexpected toml layout:\n%s", raw)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Default != "personal" || len(got.Profiles) != 2 || got.Profiles[0].Share.MCP != "default" || got.Profiles[1].Dir != "/x/.claude" {
		t.Fatalf("round trip lost data: %+v", got)
	}
}

func TestPlanExampleParses(t *testing.T) {
	src := `version = 1
default = "personal"
claude_bin = "claude"
mcp_strategy = "flag"

[[profile]]
name = "personal"
slot = 1
email = "me@example.com"
[profile.share]
mcp = "default"
settings = "default"
memory = "default"
extensions = "default"
history = "local"
`
	path := filepath.Join(t.TempDir(), "c.toml")
	os.WriteFile(path, []byte(src), 0o600)
	c, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	p, ok := c.Find("1")
	if !ok || p.Name != "personal" || p.Share.Get("history") != Local || p.Share.Get("mcp") != "default" {
		t.Fatalf("got %+v", c)
	}
}

func TestNewerVersionRejected(t *testing.T) {
	path := filepath.Join(t.TempDir(), "c.toml")
	os.WriteFile(path, []byte("version = 99\n"), 0o600)
	if _, err := Load(path); err == nil {
		t.Fatal("want error for future version")
	}
}

func TestValidateName(t *testing.T) {
	for _, bad := range []string{"", "Work", "add", "run", "list", "12", "a b", "-x", strings.Repeat("a", 33)} {
		if ValidateName(bad) == nil {
			t.Errorf("%q should be rejected", bad)
		}
	}
	for _, good := range []string{"work", "personal-2", "a1", "my_team"} {
		if err := ValidateName(good); err != nil {
			t.Errorf("%q rejected: %v", good, err)
		}
	}
}

func TestFindAndSlots(t *testing.T) {
	c := &Config{Profiles: []Profile{{Name: "a", Slot: 1}, {Name: "b", Slot: 3}}}
	if c.NextSlot() != 2 {
		t.Fatalf("next slot %d", c.NextSlot())
	}
	if p, ok := c.Find("3"); !ok || p.Name != "b" {
		t.Fatal("find by slot")
	}
	if _, ok := c.Find("2"); ok {
		t.Fatal("slot 2 should not exist")
	}
	if p, _ := c.DefaultProfile(); p.Name != "a" {
		t.Fatal("default should be lowest slot")
	}
	c.Default = "b"
	c.Remove("b")
	if c.Default != "" || len(c.Profiles) != 1 {
		t.Fatal("remove")
	}
}

func TestShareSet(t *testing.T) {
	var s Share
	if s.Get("mcp") != Local {
		t.Fatal("empty means local")
	}
	if err := s.Set("history", "on"); err != nil || s.History != "default" {
		t.Fatal("on maps to default group")
	}
	if err := s.Set("history", "team"); err != nil || s.Get("history") != "team" {
		t.Fatal("named group")
	}
	if s.Set("bogus", "on") == nil || s.Set("mcp", "Bad Group") == nil {
		t.Fatal("bad input accepted")
	}
}
