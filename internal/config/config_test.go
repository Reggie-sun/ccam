package config

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestLoadMissingIsEmpty(t *testing.T) {
	c, err := Load(filepath.Join(t.TempDir(), "none.toml"))
	if err != nil || len(c.Profiles) != 0 {
		t.Fatalf("got %+v %v", c, err)
	}
}

func TestRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	c := &Config{Default: "personal", MCPStrategy: MCPFlag}
	c.Profiles = append(c.Profiles,
		Profile{Name: "personal", Slot: 1, Email: "me@example.com"},
		Profile{Name: "work", Slot: 2, Dir: "/x/.claude"},
	)
	if err := c.AddGroup("default", DefaultItems, []string{"personal"}); err != nil {
		t.Fatal(err)
	}
	if err := c.Save(path); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), "[[profile]]") || !strings.Contains(string(raw), "[[group]]") || strings.Contains(string(raw), "[profile.share]") {
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
	src := `default = "personal"
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

func TestShareEqual(t *testing.T) {
	if !(Share{History: Local}).Equal(Share{}) || (Share{MCP: "default"}).Equal(Share{}) {
		t.Fatal("Equal")
	}
}

func TestLoadRejectsNamesThatLeaveTheFolders(t *testing.T) {
	for _, body := range []string{
		"[[profile]]\nname = \"../x\"\nslot = 1\n",
		"[[profile]]\nname = \"work\"\nslot = 1\n[profile.share]\nhistory = \"../../y\"\n",
	} {
		path := filepath.Join(t.TempDir(), "config.toml")
		os.WriteFile(path, []byte(body), 0o600)
		if _, err := Load(path); err == nil {
			t.Errorf("accepted:\n%s", body)
		}
	}
}

func TestIsNameMatchesRule(t *testing.T) {
	rule := regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,31}$`)
	for _, s := range []string{"", "a", "0", "_a", "-a", "a_", "a-", "A", "aé", "a.b", "a/b", "..", "a b",
		strings.Repeat("a", 32), strings.Repeat("a", 33), "work\n", "x\x00"} {
		if got, want := isName(s), rule.MatchString(s); got != want {
			t.Errorf("isName(%q) = %v, want %v", s, got, want)
		}
	}
}

func TestFolderForPicksTheLongestContainingPath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "code")
	c := &Config{}
	c.SetFolder(root, "personal")
	c.SetFolder(filepath.Join(root, "work"), "work")
	cases := map[string]string{
		root:                         "personal",
		filepath.Join(root, "notes"): "personal",
		filepath.Join(root, "work"):  "work",
		filepath.Join(root, "work", "api", "src"):  "work",
		filepath.Join(root, "workshop"):            "personal", // not inside root/work
		filepath.Dir(root):                         "",
		filepath.Join(filepath.Dir(root), "code2"): "",
	}
	for dir, want := range cases {
		f, ok := c.FolderFor(dir)
		got := ""
		if ok {
			got = f.Account
		}
		if got != want {
			t.Errorf("FolderFor(%s) = %q, want %q", dir, got, want)
		}
	}
}

func TestSetFolderReplacesAndUnsetRemoves(t *testing.T) {
	dir := t.TempDir()
	c := &Config{}
	c.SetFolder(dir, "a")
	c.SetFolder(dir+string(filepath.Separator), "b") // same folder, written differently
	if len(c.Folders) != 1 || c.Folders[0].Account != "b" {
		t.Fatalf("%+v", c.Folders)
	}
	if !c.UnsetFolder(dir) || len(c.Folders) != 0 || c.UnsetFolder(dir) {
		t.Fatalf("%+v", c.Folders)
	}
}

func TestRemoveDropsTheAccountsFolders(t *testing.T) {
	c := &Config{Profiles: []Profile{{Name: "a", Slot: 1}, {Name: "b", Slot: 2}}}
	c.SetFolder(filepath.Join(t.TempDir(), "x"), "a")
	c.SetFolder(filepath.Join(t.TempDir(), "y"), "b")
	c.Remove("a")
	if len(c.Folders) != 1 || c.Folders[0].Account != "b" {
		t.Fatalf("%+v", c.Folders)
	}
}

func TestFoldersRoundTripAndValidate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	dir := t.TempDir()
	c := &Config{Profiles: []Profile{{Name: "work", Slot: 1}}}
	c.SetFolder(dir, "work")
	if err := c.Save(path); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), "[[folder]]") {
		t.Fatalf("unexpected toml layout:\n%s", raw)
	}
	got, err := Load(path)
	if err != nil || len(got.Folders) != 1 || got.Folders[0].Path != dir || got.Folders[0].Account != "work" {
		t.Fatalf("%+v %v", got, err)
	}
	for _, bad := range []string{
		"[[folder]]\npath = \"relative/dir\"\naccount = \"work\"\n",
		"[[folder]]\npath = " + tomlLiteral(dir) + "\naccount = \"../x\"\n",
	} {
		os.WriteFile(path, []byte(bad), 0o600)
		if _, err := Load(path); err == nil {
			t.Fatalf("accepted %q", bad)
		}
	}
}

func tomlLiteral(s string) string { return `'` + s + `'` }

func TestAPICheckNeedsAKeySource(t *testing.T) {
	for _, from := range []string{"", "env:", "command:", "command:  ", "env:KEY", "command:pass show k", "file:/k", "KEY"} {
		err := (&API{KeyFrom: from}).Check()
		ok := from == "env:KEY" || from == "command:pass show k"
		if (err == nil) != ok {
			t.Errorf("%q: %v", from, err)
		}
	}
}
