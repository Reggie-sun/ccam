package config

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
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
	c := &Config{Default: "personal", MCPStrategy: MCPFlag, AutoSwitch: true}
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
	if !strings.Contains(string(raw), "[[profile]]") || !strings.Contains(string(raw), "[[group]]") || strings.Contains(string(raw), "[profile.share]") ||
		!strings.Contains(string(raw), "auto_switch = true") {
		t.Fatalf("unexpected toml layout:\n%s", raw)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Default != "personal" || len(got.Profiles) != 2 || got.Profiles[0].Share.MCP != "default" || got.Profiles[1].Dir != "/x/.claude" || !got.AutoSwitch {
		t.Fatalf("round trip lost data: %+v", got)
	}
}

// auto_switch is a setting for all accounts: at the top it is read, below
// a table it is named as the table's, like relay.
func TestAutoSwitchIsASetting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	os.WriteFile(path, []byte("auto_switch = true\n\n[[profile]]\nname = \"work\"\nslot = 1\n"), 0o600)
	c, err := Load(path)
	if err != nil || len(c.Unknown()) != 0 || !c.AutoSwitch {
		t.Fatalf("%v %q %+v", err, c.Unknown(), c)
	}
	os.WriteFile(path, []byte("[[profile]]\nname = \"work\"\nslot = 1\nauto_switch = true\n"), 0o600)
	c, err = Load(path)
	if err != nil || c.AutoSwitch || len(c.Unknown()) != 1 || !strings.Contains(c.Unknown()[0], "line 4: auto_switch is read as part of the [[profile]] table above it") {
		t.Fatalf("%v %q", err, c.Unknown())
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
		// A shortcut is a file name in the bin folder.
		"[[profile]]\nname = \"work\"\nslot = 1\nshortcut = \"../escaped\"\n",
		"[[profile]]\nname = \"work\"\nslot = 1\nshortcut = \"a/b\"\n",
		"[[profile]]\nname = \"work\"\nslot = 1\nshortcut = \"Upper\"\n",
		// One that would stand for Claude Code or ccam itself.
		"[[profile]]\nname = \"work\"\nslot = 1\nshortcut = \"claude\"\n",
		"[[profile]]\nname = \"work\"\nslot = 1\nshortcut = \"ccam\"\n",
		// Two accounts cannot share one.
		"[[profile]]\nname = \"work\"\nslot = 1\nshortcut = \"cc\"\n[[profile]]\nname = \"side\"\nslot = 2\nshortcut = \"cc\"\n",
	} {
		path := filepath.Join(t.TempDir(), "config.toml")
		os.WriteFile(path, []byte(body), 0o600)
		if _, err := Load(path); err == nil {
			t.Errorf("accepted:\n%s", body)
		}
	}
}

func TestIsShortcutMatchesRule(t *testing.T) {
	rule := regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)
	for _, s := range []string{"", "a", "0", ".a", "_a", "-a", "a.b", "a_", "a-", "cc.work", "A", "aé", "a/b", "..", "a b",
		strings.Repeat("a", 64), strings.Repeat("a", 65), "work\n", "x\x00"} {
		if got, want := IsShortcut(s), rule.MatchString(s); got != want {
			t.Errorf("IsShortcut(%q) = %v, want %v", s, got, want)
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

// On macOS, as on Windows, a folder reached with other letter case is the
// same folder (the file system ignores case by default), and its rule
// applies there too.
func TestFolderForIgnoresCaseWhereTheFileSystemDoes(t *testing.T) {
	if want := runtime.GOOS == "windows" || runtime.GOOS == "darwin"; foldCase != want {
		t.Fatalf("foldCase = %v on %s", foldCase, runtime.GOOS)
	}
	defer func(v bool) { foldCase = v }(foldCase)
	foldCase = true
	root := filepath.Join(t.TempDir(), "Code", "App")
	c := &Config{}
	c.SetFolder(root, "work")
	lower := strings.ToLower(root)
	if f, ok := c.FolderFor(filepath.Join(lower, "src")); !ok || f.Account != "work" {
		t.Fatalf("FolderFor(%s) misses the rule for %s", lower, root)
	}
	c.SetFolder(lower, "side")
	if len(c.Folders) != 1 || c.Folders[0].Account != "side" {
		t.Fatalf("%+v", c.Folders)
	}
	foldCase = false
	if _, ok := c.FolderFor(strings.ToUpper(root)); ok {
		t.Fatal("case counts where the file system tells it apart")
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

// Settings added at the end of config.toml, as the docs' snippets invite,
// belong to the last table in TOML. Load reads the rest and names them;
// Save refuses to write the file, which would delete them.
func TestUnreadLinesAreNamedAndKept(t *testing.T) {
	for _, tc := range []struct{ name, body string }{
		{"group format", "default = \"work\"\n\n[[group]]\nname = \"default\"\nitems = [\"mcp\"]\nmembers = [\"work\"]\n\n[[profile]]\nname = \"work\"\nslot = 1\n"},
		{"old format", "[[profile]]\nname = \"work\"\nslot = 1\n[profile.share]\nmcp = \"default\"\n"},
	} {
		path := filepath.Join(t.TempDir(), "config.toml")
		body := tc.body + "\nrelay = true\nauto_max_useage = 80\n"
		os.WriteFile(path, []byte(body), 0o600)
		c, err := Load(path)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if p, ok := c.Find("work"); !ok || p.Share.Get("mcp") != "default" || c.Relay {
			t.Fatalf("%s: read %+v", tc.name, c)
		}
		lines := strings.Count(body, "\n")
		got := strings.Join(c.Unknown(), "\n")
		for _, want := range []string{
			"line " + strconv.Itoa(lines-1) + ": relay is read as part of the",
			"go above the first [[group]], [[profile]] or [[folder]]",
			"line " + strconv.Itoa(lines) + ": auto_max_useage is not a setting of",
		} {
			if !strings.Contains(got, want) {
				t.Errorf("%s: missing %q in\n%s", tc.name, want, got)
			}
		}
		c.Default = "work"
		if err := c.Save(path); err == nil || !strings.Contains(err.Error(), "not saved") || !strings.Contains(err.Error(), "relay") {
			t.Errorf("%s: Save: %v", tc.name, err)
		}
		if after, _ := os.ReadFile(path); string(after) != body {
			t.Errorf("%s: the file changed:\n%s", tc.name, after)
		}
		if _, err := os.Stat(path + ".bak"); err == nil {
			t.Errorf("%s: a backup was written", tc.name)
		}
	}
	// At the top, a name ccam does not know (a typo, or a newer ccam's).
	path := filepath.Join(t.TempDir(), "config.toml")
	os.WriteFile(path, []byte("relay_mode = true\n\n[[profile]]\nname = \"work\"\nslot = 1\n"), 0o600)
	c, err := Load(path)
	if err != nil || len(c.Unknown()) != 1 || !strings.Contains(c.Unknown()[0], "line 1: relay_mode is not a ccam setting") {
		t.Fatalf("%v %q", err, c.Unknown())
	}
}

// ccam v0.1 and v0.2 wrote keys later releases dropped (`version`, and
// v0.1's `name_shims`): such a config is read and saved again without them,
// as promised when they were dropped.
func TestKeysOfEarlierReleasesAreDropped(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	os.WriteFile(path, []byte("version = 2\nname_shims = true\n\n[[profile]]\nname = \"work\"\nslot = 1\n"), 0o600)
	c, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Unknown()) != 0 {
		t.Fatalf("unknown: %v", c.Unknown())
	}
	if err := c.Save(path); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(path); strings.Contains(string(data), "version") || strings.Contains(string(data), "name_shims") {
		t.Fatalf("saved: %s", data)
	}
}

// Two accounts that resolve to one folder would share a login, and
// removing one would trash the other's folder.
func TestTwoAccountsCannotHaveOneFolder(t *testing.T) {
	for _, body := range []string{
		"[[profile]]\nname = \"a\"\nslot = 1\n\n[[profile]]\nname = \"b\"\nslot = 2\nfolder = \"a\"\n",
		"[[profile]]\nname = \"a\"\nslot = 1\nfolder = \"x\"\n\n[[profile]]\nname = \"b\"\nslot = 2\nfolder = \"x\"\n",
	} {
		path := filepath.Join(t.TempDir(), "config.toml")
		os.WriteFile(path, []byte(body), 0o600)
		if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "accounts a and b use the same folder") {
			t.Errorf("%s: %v", body, err)
		}
	}
}
