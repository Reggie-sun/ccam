package config

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func load(t *testing.T, src string) (*Config, string, error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	c, err := Load(path)
	return c, path, err
}

// The conversion must give every account exactly the folders it had, so
// nothing has to move on disk.
func TestConvertSharesKeepsEveryAccountsFolders(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	dirs := []string{Local, Local, "default", "default", "team", "none", strings.Repeat("x", 32)}
	for round := range 2000 {
		c := &Config{}
		n := 1 + r.IntN(6)
		for i := range n {
			p := Profile{Name: fmt.Sprintf("a%d", i), Slot: i + 1}
			for _, item := range Items {
				*p.Share.field(item) = dirs[r.IntN(len(dirs))]
			}
			c.Profiles = append(c.Profiles, p)
		}
		before := map[string]Share{}
		for _, p := range c.Profiles {
			before[p.Name] = p.Share
		}
		c.ConvertShares()
		got, err := c.resolved()
		if err != nil {
			t.Fatalf("round %d: %v\n%+v", round, err, c.Groups)
		}
		for name, s := range before {
			if !got[name].Equal(s) {
				t.Fatalf("round %d: %s had %+v, groups give %+v\n%+v", round, name, s, got[name], c.Groups)
			}
		}
		seen := map[string]bool{}
		for _, g := range c.Groups {
			if err := ValidateGroupName(g.Name); err != nil || seen[g.Name] {
				t.Fatalf("round %d: bad or repeated name %q: %v", round, g.Name, err)
			}
			seen[g.Name] = true
		}
	}
}

func TestConvertSharesNames(t *testing.T) {
	all := Share{MCP: "default", Settings: "default", Memory: "default", Extensions: "default", History: "default"}
	noHistory := all
	noHistory.History = Local
	c := &Config{Profiles: []Profile{{Name: "personal", Share: noHistory}, {Name: "work", Share: all}, {Name: "side", Share: all}}}
	c.ConvertShares()
	want := []Group{
		{Name: "default", Items: DefaultItems, Members: []string{"personal", "work", "side"}},
		{Name: "default-history", Dir: "default", Items: []string{"history"}, Members: []string{"work", "side"}},
	}
	if fmt.Sprint(c.Groups) != fmt.Sprint(want) {
		t.Fatalf("got %+v", c.Groups)
	}
}

func TestLoadConvertsOldFormatAndKeepsBackup(t *testing.T) {
	src := `default = "personal"
[[profile]]
name = "personal"
slot = 1
[profile.share]
mcp = "default"
settings = "default"
history = "team"

[[profile]]
name = "work"
slot = 2
[profile.share]
mcp = "default"
settings = "default"
history = "team"

[[profile]]
name = "solo"
slot = 3
[profile.share]
`
	c, path, err := load(t, src)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Groups) != 2 || c.Groups[0].Name != "default" || c.Groups[1].Name != "team" || c.Profiles[1].Share.Get("history") != "team" {
		t.Fatalf("got %+v", c.Groups)
	}
	if !c.Profiles[2].Share.Equal(Share{}) {
		t.Fatal("solo shares nothing")
	}
	if _, err := os.Stat(path + ".bak"); err == nil {
		t.Fatal("loading must not write")
	}
	if err := c.Save(path); err != nil {
		t.Fatal(err)
	}
	bak, _ := os.ReadFile(path + ".bak")
	if string(bak) != src {
		t.Fatalf("backup:\n%s", bak)
	}
	raw, _ := os.ReadFile(path)
	if strings.Contains(string(raw), "[profile.share]") || !strings.Contains(string(raw), "[[group]]") {
		t.Fatalf("saved:\n%s", raw)
	}
	again, err := Load(path)
	if err != nil || again.legacy || again.Profiles[0].Share.Get("history") != "team" {
		t.Fatalf("%v %+v", err, again)
	}
	// A second save leaves the backup alone.
	again.Default = "work"
	if err := again.Save(path); err != nil {
		t.Fatal(err)
	}
	if bak2, _ := os.ReadFile(path + ".bak"); string(bak2) != src {
		t.Fatal("backup rewritten")
	}
}

func TestLoadRefusesBothFormats(t *testing.T) {
	_, _, err := load(t, `[[group]]
name = "default"
items = ["mcp"]
members = ["a"]

[[profile]]
name = "a"
slot = 1
[profile.share]
mcp = "default"
`)
	if err == nil || !strings.Contains(err.Error(), "[profile.share]") {
		t.Fatal(err)
	}
}

func TestLoadRefusesBadGroups(t *testing.T) {
	cases := map[string]string{
		"both group a and group b": `[[group]]
name = "a"
items = ["mcp", "history"]
members = ["x"]
[[group]]
name = "b"
items = ["history"]
members = ["x"]
[[profile]]
name = "x"
slot = 1
`,
		"unknown item": `[[group]]
name = "a"
items = ["tokens"]
members = []
`,
		"defined twice": `[[group]]
name = "a"
items = []
members = []
[[group]]
name = "a"
items = []
members = []
`,
		"cannot be a group name": `[[group]]
name = "local"
items = []
members = []
`,
		"invalid folder": `[[group]]
name = "a"
dir = "../x"
items = []
members = []
`,
		`invalid folder "local"`: `[[group]]
name = "a"
dir = "local"
items = ["mcp"]
members = []
`,
		"both share history through the folder groups/team": `[[group]]
name = "team"
items = ["history"]
members = ["x"]
[[group]]
name = "pair"
dir = "team"
items = ["history"]
members = ["y"]
[[profile]]
name = "x"
slot = 1
[[profile]]
name = "y"
slot = 2
`,
	}
	for want, src := range cases {
		if _, _, err := load(t, src); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v", want, err)
		}
	}
}

func TestUnknownMembersAreIgnored(t *testing.T) {
	c, _, err := load(t, `[[group]]
name = "default"
items = ["mcp"]
members = ["gone", "a"]
[[profile]]
name = "a"
slot = 1
`)
	if err != nil || c.Profiles[0].Share.Get("mcp") != "default" {
		t.Fatalf("%v %+v", err, c)
	}
}

func groupsConfig(t *testing.T) *Config {
	t.Helper()
	c := &Config{Profiles: []Profile{{Name: "personal", Slot: 1}, {Name: "work", Slot: 2}, {Name: "side", Slot: 3}}}
	if err := c.AddGroup("default", DefaultItems, []string{"personal", "work", "3"}); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestJoinRefusesOverlapUnlessReplace(t *testing.T) {
	c := groupsConfig(t)
	if err := c.AddGroup("team", []string{"settings", "history"}, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := c.JoinGroup("work", "team", false); err == nil || !strings.Contains(err.Error(), "already gets settings from group default") {
		t.Fatal(err)
	}
	left, err := c.JoinGroup("work", "team", true)
	if err != nil || !slices.Equal(left, []string{"default"}) {
		t.Fatalf("%v %v", left, err)
	}
	w, _ := c.Find("work")
	if w.Share.Get("history") != "team" || w.Share.Get("mcp") != Local {
		t.Fatalf("%+v", w.Share)
	}
	if err := c.LeaveGroup("work", "team"); err != nil || !w.Share.Equal(Share{}) {
		t.Fatalf("%v %+v", err, w.Share)
	}
	if err := c.LeaveGroup("work", "team"); err == nil {
		t.Fatal("not a member")
	}
}

func TestSetItemsRefusesOverlap(t *testing.T) {
	c := groupsConfig(t)
	if err := c.AddGroup("team", []string{"history"}, []string{"work"}); err != nil {
		t.Fatal(err)
	}
	if err := c.SetGroupItems("team", []string{"history", "mcp"}); err == nil {
		t.Fatal("work would get mcp twice")
	}
}

func TestRenameKeepsFolder(t *testing.T) {
	c := groupsConfig(t)
	if err := c.RenameGroup("default", "main"); err != nil {
		t.Fatal(err)
	}
	g, _ := c.FindGroup("main")
	if g.Folder() != "default" || c.Profiles[0].Share.Get("mcp") != "default" {
		t.Fatalf("%+v", g)
	}
	if err := c.AddGroup("default", nil, nil); err == nil {
		t.Fatal("folder default is in use")
	}
	if err := c.RenameGroup("main", "default"); err != nil {
		t.Fatal(err)
	}
	if g, _ := c.FindGroup("default"); g.Dir != "" {
		t.Fatalf("%+v", g)
	}
}

func TestRemoveAccountLeavesGroups(t *testing.T) {
	c := groupsConfig(t)
	c.Remove("work")
	g, _ := c.FindGroup("default")
	if !slices.Equal(g.Members, []string{"personal", "side"}) {
		t.Fatal(g.Members)
	}
}

func TestSaveRefusesShareSetBehindTheGroups(t *testing.T) {
	c := groupsConfig(t)
	c.Profiles[0].Share.History = "default"
	if err := c.Save(filepath.Join(t.TempDir(), "c.toml")); err == nil {
		t.Fatal("expected an error")
	}
}

func TestCloneIsDeep(t *testing.T) {
	c := groupsConfig(t)
	n := c.Clone()
	if err := n.LeaveGroup("work", "default"); err != nil {
		t.Fatal(err)
	}
	if g, _ := c.FindGroup("default"); len(g.Members) != 3 || c.Profiles[1].Share.Get("mcp") != "default" {
		t.Fatal("clone shares memory with the original")
	}
}

// Load runs on every launch: the groups format, and the old one that is
// converted until the first save.
func BenchmarkLoad(b *testing.B) {
	var groups, old strings.Builder
	groups.WriteString("[[group]]\nname = \"default\"\nitems = [\"mcp\", \"settings\", \"memory\", \"extensions\"]\nmembers = [\"a1\", \"a2\", \"a3\", \"a4\", \"a5\"]\n")
	for i := 1; i <= 5; i++ {
		fmt.Fprintf(&groups, "[[profile]]\nname = \"a%d\"\nslot = %d\n", i, i)
		fmt.Fprintf(&old, "[[profile]]\nname = \"a%d\"\nslot = %d\n[profile.share]\nmcp = \"default\"\nsettings = \"default\"\nmemory = \"default\"\nextensions = \"default\"\n", i, i)
	}
	for name, src := range map[string]string{"groups": groups.String(), "old": old.String()} {
		path := filepath.Join(b.TempDir(), "config.toml")
		os.WriteFile(path, []byte(src), 0o600)
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				if _, err := Load(path); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// A member left over from an account that is gone must not pass its groups
// to a new account of that name, nor break the save.
func TestStaleMemberDoesNotPassToANewAccount(t *testing.T) {
	c := groupsConfig(t)
	c.Groups = append(c.Groups, Group{Name: "team", Items: []string{"history"}, Members: []string{"job"}})
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(c.Resolve())
	c.DropMember("job")
	c.Profiles = append(c.Profiles, Profile{Name: "job", Slot: 9})
	must(c.Resolve())
	if p, _ := c.Find("job"); !p.Share.Equal(Share{}) {
		t.Fatalf("%+v", p.Share)
	}
	must(c.Save(filepath.Join(t.TempDir(), "c.toml")))
	if !c.IsAccount("job") || c.IsAccount("9") {
		t.Fatal("IsAccount takes names only")
	}
}

func TestDefaultGroupsAfterARename(t *testing.T) {
	c := &Config{Profiles: []Profile{{Name: "a", Slot: 1}, {Name: "b", Slot: 2}}}
	c.Groups = []Group{
		{Name: "default-history", Dir: "default", Items: []string{"history"}, Members: []string{"a"}},
		{Name: "main", Dir: "default", Items: DefaultItems, Members: []string{"a", "b"}},
	}
	if got := c.DefaultGroups(); len(got) != 1 || got[0] != "main" {
		t.Fatalf("%v", got)
	}
}
