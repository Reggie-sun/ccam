package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/session"
)

func (c *testCLI) ok(args ...string) string {
	c.t.Helper()
	out, errOut, code := c.run("", args...)
	if code != 0 {
		c.t.Fatalf("%v exited %d\n%s%s", args, code, out, errOut)
	}
	return out
}

func (c *testCLI) cfg() *config.Config {
	c.t.Helper()
	cfg, err := config.Load(c.p.ConfigFile())
	if err != nil {
		c.t.Fatal(err)
	}
	return cfg
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// If applying a group change fails halfway, the groups are already saved so
// `doctor --fix` (or running the same command again) finishes the change
// instead of undoing the half that was done.
func TestFailedGroupChangeCanBeFinished(t *testing.T) {
	if os.Geteuid() == 0 || runtime.GOOS == "windows" {
		t.Skip("root, and Windows, ignore the read-only folder this test relies on")
	}
	c := newTestCLI(t, nil)
	work := c.p.ProfileDir("work")
	os.WriteFile(filepath.Join(work, "CLAUDE.md"), []byte("mine"), 0o600)
	group := c.p.GroupDir("docs")
	os.MkdirAll(c.p.GroupsRoot(), 0o700)
	os.Mkdir(group, 0o500) // moving into the group fails
	t.Cleanup(func() { os.Chmod(group, 0o700) })

	_, errOut, code := c.run("", "group", "add", "docs", "--items", "memory", "--members", "work", "--yes")
	if code == 0 || !strings.Contains(errOut, "doctor --fix") {
		t.Fatalf("want failure pointing at doctor --fix, got %d %q", code, errOut)
	}
	if p, _ := c.cfg().Find("work"); p.Share.Get("memory") != "docs" {
		t.Fatal("the groups must be saved before applying")
	}
	os.Chmod(group, 0o700)
	c.run("", "doctor", "--fix") // exit code also reflects unrelated checks
	if b, _ := os.ReadFile(filepath.Join(work, "CLAUDE.md")); string(b) != "mine" {
		t.Fatalf("content after finishing: %q", b)
	}
	if target, _ := os.Readlink(filepath.Join(work, "CLAUDE.md")); target != filepath.Join(group, "CLAUDE.md") {
		t.Fatalf("not linked: %q", target)
	}
}

// Moving data under a running session would lose what it writes meanwhile.
func TestGroupChangeRefusesWhileRunning(t *testing.T) {
	c := newTestCLI(t, nil)
	if err := session.Register(launcher.SessionDir(c.p, "work")); err != nil {
		t.Fatal(err)
	}
	if !launcher.Running(c.p, "work") {
		t.Skip("session tracking is not supported on this platform")
	}
	_, errOut, code := c.run("", "group", "set", "team", "--add", "memory", "--yes")
	if code == 0 || !strings.Contains(errOut, "work is running") {
		t.Fatalf("want refusal, got %d %q", code, errOut)
	}
	if p, _ := c.cfg().FindGroup("team"); slices.Contains(p.Items, "memory") {
		t.Fatal("nothing may be saved")
	}
	if _, _, code := c.run("", "--dry-run", "group", "set", "team", "--add", "memory"); code != 0 {
		t.Fatal("a dry run only shows the plan and may run")
	}
	// Leaving concerns only the account that leaves.
	c.ok("group", "leave", "team", "side", "--yes")
}

// Every member of a group moves together; the first content comes from
// --from, then the default account, then the others by number.
func TestGroupSetMovesEveryMember(t *testing.T) {
	c := newTestCLI(t, nil)
	for _, n := range []string{"personal", "work", "side"} {
		os.WriteFile(filepath.Join(c.p.ProfileDir(n), "CLAUDE.md"), []byte(n), 0o600)
		os.MkdirAll(filepath.Join(c.p.ProfileDir(n), "skills", n), 0o700)
		os.WriteFile(filepath.Join(c.p.ProfileDir(n), "skills", n, "SKILL.md"), []byte(n), 0o600)
	}
	out := c.ok("group", "set", "team", "--add", "memory,extensions", "--from", "work", "--yes")
	for _, s := range []string{"personal: +memory (group team), +extensions (group team)", "work:", "side:", "✓ Group team shares extensions, history, memory"} {
		if !strings.Contains(out, s) && !strings.Contains(out, strings.Replace(s, "extensions, history, memory", "memory, extensions, history", 1)) {
			t.Errorf("missing %q in\n%s", s, out)
		}
	}
	for _, n := range []string{"personal", "work", "side"} {
		dir := c.p.ProfileDir(n)
		if got := readFile(t, filepath.Join(dir, "CLAUDE.md")); got != "work" {
			t.Errorf("%s sees %q, want work's copy", n, got)
		}
		for _, s := range []string{"personal", "work", "side"} {
			if _, err := os.Stat(filepath.Join(dir, "skills", s, "SKILL.md")); err != nil {
				t.Errorf("%s misses skill %s", n, s)
			}
		}
	}
	backups, _ := filepath.Glob(filepath.Join(c.p.ProfileDir("side"), ".ccam-backup", "*", "CLAUDE.md"))
	if len(backups) != 1 || readFile(t, backups[0]) != "side" {
		t.Fatalf("side's own copy must be backed up: %v", backups)
	}

	// Taking an item out leaves every member a copy.
	c.ok("group", "set", "team", "--remove", "memory", "--yes")
	for _, n := range []string{"personal", "work", "side"} {
		path := filepath.Join(c.p.ProfileDir(n), "CLAUDE.md")
		if _, link := os.Readlink(path); link == nil || readFile(t, path) != "work" {
			t.Errorf("%s: want its own copy of the group's CLAUDE.md", n)
		}
	}
}

func TestGroupJoinLeaveAndOverlap(t *testing.T) {
	c := newTestCLI(t, nil)
	c.ok("group", "add", "default", "--items", "settings,history", "--yes")
	_, errOut, code := c.run("", "group", "join", "default", "work", "--yes")
	if code == 0 || !strings.Contains(errOut, "already gets history from group team") || !strings.Contains(errOut, "--replace") {
		t.Fatalf("%d %s", code, errOut)
	}
	out := c.ok("group", "join", "default", "work", "--replace", "--yes")
	if !strings.Contains(out, "work left group team and joined group default") {
		t.Fatal(out)
	}
	cfg := c.cfg()
	if p, _ := cfg.Find("work"); p.Share.Get("history") != "default" || p.Share.Get("settings") != "default" {
		t.Fatalf("%+v", p.Share)
	}
	c.ok("group", "leave", "default", "work", "--yes")
	if p, _ := c.cfg().Find("work"); !p.Share.Equal(config.Share{}) {
		t.Fatalf("%+v", p.Share)
	}
	if out := c.ok("share", "work"); !strings.Contains(out, "in no group") {
		t.Fatal(out)
	}
	if out := c.ok("share", "personal"); !strings.Contains(out, "group team") || !strings.Contains(out, "Groups: team") {
		t.Fatal(out)
	}
	if _, errOut, _ := c.run("", "share", "personal", "--history", "off"); !strings.Contains(errOut, "set by groups") {
		t.Fatal(errOut)
	}
}

func TestGroupListShowRenameRemove(t *testing.T) {
	c := newTestCLI(t, nil)
	os.MkdirAll(filepath.Join(c.p.GroupDir("team"), "projects", "p"), 0o700)
	os.WriteFile(filepath.Join(c.p.GroupDir("team"), "projects", "p", "s.jsonl"), []byte("x"), 0o600)
	c.run("", "doctor", "--fix") // link the accounts to team
	out := c.ok("group")
	if !strings.Contains(out, "team") || !strings.Contains(out, "history") || !strings.Contains(out, "personal, work, side") {
		t.Fatal(out)
	}
	c.ok("group", "rename", "team", "pair")
	out = c.ok("group", "show", "pair")
	if !strings.Contains(out, "Group pair") || !strings.Contains(out, filepath.Join("groups", "team")) {
		t.Fatal(out)
	}
	if p, _ := c.cfg().Find("work"); p.Share.Get("history") != "team" {
		t.Fatal("a rename moves nothing")
	}
	_, errOut, code := c.run("", "group", "remove", "pair")
	if code == 0 || !strings.Contains(errOut, "--leave") {
		t.Fatalf("%d %s", code, errOut)
	}
	out = c.ok("group", "remove", "pair", "--leave", "--yes")
	if !strings.Contains(out, "Group pair removed") || c.cfg().Groups != nil && len(c.cfg().Groups) != 0 {
		t.Fatal(out)
	}
	if _, err := os.Stat(c.p.GroupDir("team")); err == nil {
		t.Fatal("the group folder should be in the trash")
	}
	if trashed, _ := filepath.Glob(filepath.Join(c.p.TrashDir(), "group-team-*")); len(trashed) != 1 {
		t.Fatal("not in the trash")
	}
	if readFile(t, filepath.Join(c.p.ProfileDir("work"), "projects", "p", "s.jsonl")) != "x" {
		t.Fatal("each account keeps a copy")
	}
}

func TestAddJoinsDefaultGroup(t *testing.T) {
	c := newTestCLI(t, nil)
	// There are groups, but none is default: join none rather than make one.
	if out := c.ok("add", "laptop", "--no-login"); !strings.Contains(out, "no group default") {
		t.Fatal(out)
	}
	c.ok("add", "phone", "--no-login", "--group", "default")
	c.ok("add", "desk", "--no-login")
	g, ok := c.cfg().FindGroup("default")
	if !ok || !slices.Equal(g.Items, config.DefaultItems) || !slices.Equal(g.Members, []string{"phone", "desk"}) {
		t.Fatalf("%+v", g)
	}
	c.ok("add", "tablet", "--no-login", "--group", "team")
	if p, _ := c.cfg().Find("tablet"); p.Share.Get("history") != "team" || p.Share.Get("mcp") != config.Local {
		t.Fatalf("%+v", p.Share)
	}
	c.ok("add", "solo", "--no-login", "--no-group")
	if p, _ := c.cfg().Find("solo"); !p.Share.Equal(config.Share{}) {
		t.Fatalf("%+v", p.Share)
	}
	_, errOut, code := c.run("", "add", "x1", "--no-login", "--group", "nope")
	if code == 0 || !strings.Contains(errOut, `no group "nope"`) {
		t.Fatal(errOut)
	}
	if _, err := os.Stat(c.p.ProfileDir("x1")); err == nil {
		t.Fatal("a refused add must leave no folder")
	}
	if _, errOut, _ := c.run("", "add", "x2", "--no-login", "--share", "history=team"); !strings.Contains(errOut, "--group") {
		t.Fatal(errOut)
	}
}

// An old config is read as groups and rewritten as groups at the first
// save, with the old file kept.
func TestOldShareFormatIsConverted(t *testing.T) {
	c := newTestCLI(t, nil)
	old := `default = "personal"
[[profile]]
name = "personal"
slot = 1
[profile.share]
mcp = "default"
settings = "default"

[[profile]]
name = "work"
slot = 2
[profile.share]
mcp = "default"
settings = "default"
history = "default"
`
	os.WriteFile(c.p.ConfigFile(), []byte(old), 0o600)
	out := c.ok("group")
	if !strings.Contains(out, "default") || !strings.Contains(out, "default-history (folder default)") {
		t.Fatal(out)
	}
	if _, err := os.Stat(c.p.ConfigFile() + ".bak"); err == nil {
		t.Fatal("reading must not write")
	}
	c.ok("use", "work")
	if readFile(t, c.p.ConfigFile()+".bak") != old {
		t.Fatal("the old file must be kept")
	}
	if s := readFile(t, c.p.ConfigFile()); strings.Contains(s, "[profile.share]") || !strings.Contains(s, "[[group]]") {
		t.Fatal(s)
	}
}

// After a rename the old name is the folder of another group: writing MCP
// servers or migrating "into it" by that name is refused.
func TestOldGroupNameIsNotAFolder(t *testing.T) {
	c := newTestCLI(t, nil)
	c.ok("group", "rename", "team", "pair")
	_, errOut, code := c.run("", "mcp", "add", "x", "--group", "team", "--", "echo")
	if code == 0 || !strings.Contains(errOut, "group pair uses the folder groups/team") {
		t.Fatalf("%d %s", code, errOut)
	}
}

// doctor --fix leaves what a running account would still write into.
func TestDoctorFixWaitsForARunningAccount(t *testing.T) {
	c := newTestCLI(t, nil)
	if err := session.Register(launcher.SessionDir(c.p, "work")); err != nil {
		t.Fatal(err)
	}
	if !launcher.Running(c.p, "work") {
		t.Skip("session tracking is not supported on this platform")
	}
	os.MkdirAll(filepath.Join(c.p.ProfileDir("work"), "projects"), 0o700) // a real folder where team's link belongs
	out, _, _ := c.run("", "doctor", "--fix")
	if !strings.Contains(out, "work: is running; quit Claude Code there") {
		t.Fatal(out)
	}
	if info, err := os.Lstat(filepath.Join(c.p.ProfileDir("work"), "projects")); err != nil || info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("moved under a running session")
	}
}
