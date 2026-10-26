package group

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/lock"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/share"
)

// perms is whether folder permissions can make a move fail (not as root,
// not on Windows).
var perms = runtime.GOOS != "windows" && os.Geteuid() != 0

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
	if err := ch.Check(func(string, string) bool { return false }); err != nil {
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
// with --fresh removes the copies ccam made in .claude.json.
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

	// A copy Sync made for the account goes when MCP sharing stops (--fresh).
	dir := p.ProfileDir("a")
	if _, err := mcp.Sync(dir, p.GroupDir("h"), launcher.MCPStateDir(p, "a")); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(read(t, filepath.Join(dir, ".claude.json")), `"fs"`) {
		t.Fatal("sync did not copy")
	}
	ch, _ = apply(t, p, cfg, func(c *config.Config) error { return c.LeaveGroup("a", "h") }, Options{Fresh: true})
	var buf bytes.Buffer
	ch.Describe(&buf)
	if !strings.Contains(buf.String(), "MCP servers ccam copied into .claude.json: fs") || !strings.Contains(buf.String(), "MCP servers fs stop loading for a (--fresh)") {
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
	err = ch.Check(func(string, string) bool { return true })
	if err == nil || !strings.Contains(err.Error(), "a, b are running") {
		t.Fatal(err)
	}
}

// A plan whose questions waited is checked again on the config about to be
// saved: an account that started meanwhile, or groups another ccam changed,
// stop it; unrelated changes do not.
func TestRecheckRefusesWhatChangedMeanwhile(t *testing.T) {
	p, cfg := setup(t, "a", "b")
	add := func(c *config.Config) error { return c.AddGroup("g", []string{"memory"}, []string{"a"}) }
	next := cfg.Clone()
	if err := add(next); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, cfg, next, Options{})
	if err != nil {
		t.Fatal(err)
	}
	idle := func(string, string) bool { return false }
	cur := cfg.Clone()
	cur.Default = "b"
	if err := ch.Recheck(cur, add, idle); err != nil {
		t.Fatalf("another default account: %v", err)
	}
	cur = cfg.Clone()
	if err := cur.AddGroup("h", []string{"settings"}, []string{"b"}); err != nil {
		t.Fatal(err)
	}
	if err := ch.Recheck(cur, add, idle); err == nil || !strings.Contains(err.Error(), "the groups changed while this plan waited") {
		t.Fatalf("groups changed meanwhile: %v", err)
	}
	var asked []string
	err = ch.Recheck(cfg.Clone(), add, func(name, dir string) bool { asked = append(asked, name+" "+dir); return true })
	if err == nil || !strings.Contains(err.Error(), "a is running") || len(asked) != 1 || asked[0] != "a "+p.ProfileDir("a") {
		t.Fatalf("started meanwhile: %v, asked %q", err, asked)
	}
}

// A group without members is in no account's shares, so a change to it
// shows only in what the edit then gives the account: joining it after
// another ccam changed its items would save shares the plan never moved.
func TestRecheckSeesAMemberlessGroupChange(t *testing.T) {
	p, cfg := setup(t, "a")
	if err := cfg.AddGroup("docs", []string{"settings", "memory"}, nil); err != nil {
		t.Fatal(err)
	}
	join := func(c *config.Config) error { _, err := c.JoinGroup("a", "docs", false); return err }
	next := cfg.Clone()
	if err := join(next); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, cfg, next, Options{})
	if err != nil {
		t.Fatal(err)
	}
	idle := func(string, string) bool { return false }
	if err := ch.Recheck(cfg.Clone(), join, idle); err != nil {
		t.Fatalf("nothing changed: %v", err)
	}
	for _, items := range [][]string{{"settings"}, {"settings", "memory", "history"}} {
		cur := cfg.Clone()
		if err := cur.SetGroupItems("docs", items); err != nil {
			t.Fatal(err)
		}
		if err := ch.Recheck(cur, join, idle); err == nil || !strings.Contains(err.Error(), "the groups changed while this plan waited") {
			t.Fatalf("docs now shares %v: %v", items, err)
		}
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
	if _, err := mcp.Sync(dir, p.GroupDir("g"), launcher.MCPStateDir(p, "a")); err != nil {
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

// servers reads an account's user-scope MCP servers: name → command.
func servers(t *testing.T, dir string) map[string]string {
	t.Helper()
	var doc struct {
		Servers map[string]struct {
			Command string `json:"command"`
		} `json:"mcpServers"`
	}
	data, err := os.ReadFile(filepath.Join(dir, ".claude.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for n, d := range doc.Servers {
		out[n] = d.Command
	}
	return out
}

// An account that stops sharing MCP keeps the group's servers in its own
// .claude.json, as leave, set --remove and remove --leave promise: its own
// server of a group server's name stays, and with sync its copies become
// its own. Joining again takes the kept copies out, as the group gives
// them again.
func TestLeavingMCPKeepsTheGroupsServers(t *testing.T) {
	for _, c := range []struct {
		name string
		edit func(*config.Config) error
		sync bool
	}{
		{"leave", func(c *config.Config) error { return c.LeaveGroup("a", "g") }, false},
		{"set --remove mcp", func(c *config.Config) error { return c.SetGroupItems("g", []string{"memory"}) }, false},
		{"remove --leave", func(c *config.Config) error { return c.RemoveGroup("g") }, false},
		{"leave under sync", func(c *config.Config) error { return c.LeaveGroup("a", "g") }, true},
	} {
		t.Run(c.name, func(t *testing.T) {
			p, cfg := setup(t, "a", "b")
			if err := cfg.AddGroup("g", []string{"mcp", "memory"}, []string{"a", "b"}); err != nil {
				t.Fatal(err)
			}
			os.MkdirAll(p.GroupDir("g"), 0o700)
			os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"},"mine":{"command":"group-mine"}}}`), 0o600)
			dir := p.ProfileDir("a")
			os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"oauthAccount":{"emailAddress":"a@x"},"mcpServers":{"mine":{"command":"own-mine"}}}`), 0o600)
			if c.sync {
				if _, err := mcp.Sync(dir, p.GroupDir("g"), launcher.MCPStateDir(p, "a")); err != nil {
					t.Fatal(err)
				}
			}
			ch, next := apply(t, p, cfg, c.edit, Options{})
			var buf bytes.Buffer
			ch.Describe(&buf)
			copied := "copy     MCP servers fs: " + filepath.FromSlash("groups/g/mcp.json") + " → " + filepath.FromSlash("profiles/a/.claude.json")
			for _, want := range []string{copied, "keep     a's own MCP servers mine"} {
				if !strings.Contains(buf.String(), want) {
					t.Errorf("plan lacks %q:\n%s", want, buf.String())
				}
			}
			if got := servers(t, dir); got["fs"] != "fs" || got["mine"] != "own-mine" {
				t.Fatalf("a has %v", got)
			}
			if !strings.Contains(read(t, filepath.Join(dir, ".claude.json")), "a@x") || mcp.HasCopy(launcher.MCPStateDir(p, "a")) {
				t.Fatal("other keys lost, or a sync stamp is left")
			}
			if c.name != "leave" {
				return
			}
			// Back in the group, the kept copy goes; its own servers stay, and
			// the plan names the one the group lacks.
			data := strings.Replace(read(t, filepath.Join(dir, ".claude.json")), `"mcpServers": {`, `"mcpServers": {"solo": {"command": "s"},`, 1)
			os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(data), 0o600)
			ch, _ = apply(t, p, next, func(c *config.Config) error { _, err := c.JoinGroup("a", "g", false); return err }, Options{})
			buf.Reset()
			ch.Describe(&buf)
			for _, want := range []string{"remove   MCP servers fs kept in .claude.json when a left a group", "Note: a's own MCP servers solo stay its own"} {
				if !strings.Contains(buf.String(), want) {
					t.Errorf("rejoin plan lacks %q:\n%s", want, buf.String())
				}
			}
			if got := servers(t, dir); got["fs"] != "" || got["mine"] != "own-mine" || got["solo"] != "s" {
				t.Fatalf("after joining again a has %v", got)
			}
		})
	}
}

// A change that stops halfway is noted; doctor --fix finishes it with
// Resume: the new group's first MCP servers and the links, the --from
// account's copy first. A change that fails before its config is saved
// leaves no note behind.
func TestStoppedChangeIsFinished(t *testing.T) {
	if !perms {
		t.Skip("needs folder permissions")
	}
	p, cfg := setup(t, "a", "b")
	if err := cfg.AddGroup("g", []string{"mcp"}, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	if err := cfg.AddGroup("h", []string{"mcp", "memory"}, nil); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(p.GroupDir("g"), 0o700)
	os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"}}}`), 0o600)
	next := cfg.Clone()
	for _, n := range []string{"a", "b"} {
		if _, err := next.JoinGroup(n, "h", true); err != nil {
			t.Fatal(err)
		}
	}
	ch, err := Build(p, cfg, next, Options{Seed: "b"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil || !Pending(p) {
		t.Fatalf("not noted: %v", err)
	}
	os.MkdirAll(p.GroupDir("h"), 0o500) // the moves cannot write there
	err = ch.Apply()
	os.Chmod(p.GroupDir("h"), 0o700)
	if err == nil || !strings.Contains(err.Error(), "finish with `ccam doctor --fix`") || !Pending(p) {
		t.Fatalf("apply: %v, pending %v", err, Pending(p))
	}
	if strings.Contains(err.Error(), "The originals are in") && !fsutil.Exists(filepath.Join(p.ProfileDir("b"), share.BackupDir)) {
		t.Errorf("names a backup that is not there: %v", err)
	}
	done, err := resumeLocked(p, next, func(string, string) bool { return false })
	if err != nil || strings.Join(done.Names(), ",") != "b,a" || Pending(p) {
		t.Fatalf("resume: %v %v, pending %v", done, err, Pending(p))
	}
	if !strings.Contains(read(t, mcp.Path(p.GroupDir("h"))), `"fs"`) {
		t.Fatal("group h did not get its first servers")
	}
	if read(t, filepath.Join(p.ProfileDir("a"), "CLAUDE.md")) != "b" {
		t.Fatal("--from b: b's memory goes into the group first")
	}

	// Refused while an account runs; the note stays.
	ch, _ = Build(p, cfg, next, Options{})
	ch.Begin()
	ch.release() // ccam stopped
	if _, err := resumeLocked(p, next, func(name, _ string) bool { return name == "a" }); err == nil || !strings.Contains(err.Error(), "a is running") || !Pending(p) {
		t.Fatalf("running: %v", err)
	}
	ch.Abort()
	if Pending(p) {
		t.Fatal("Abort left the note")
	}
}

// Steps noted for a config that was never saved are not done: the account
// still gets MCP from the group the note's step would have left.
func TestResumeSkipsStepsTheConfigDoesNotHave(t *testing.T) {
	p, cfg := setup(t, "a")
	if err := cfg.AddGroup("g", []string{"mcp"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(p.GroupDir("g"), 0o700)
	os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"}}}`), 0o600)
	next := cfg.Clone()
	if err := next.LeaveGroup("a", "g"); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, cfg, next, Options{})
	if err != nil || ch.Accounts[0].KeepMCP == "" {
		t.Fatalf("%v %+v", err, ch)
	}
	ch.Begin() // and then the save failed, so cfg is what is on disk
	ch.release()
	if _, err := resumeLocked(p, cfg, func(string, string) bool { return false }); err != nil {
		t.Fatal(err)
	}
	if fsutil.Exists(filepath.Join(p.ProfileDir("a"), ".claude.json")) {
		t.Fatal("kept a copy of the group's servers for an account that still shares them")
	}
}

// No change is noted over another one; one that doctor --fix finished
// meanwhile is not carried out again; and while a change is carried out,
// another ccam neither notes nor finishes one, and doctor says it runs.
func TestNotesDoNotOverlap(t *testing.T) {
	p, cfg := setup(t, "a", "b")
	plan := func(edit func(*config.Config) error) (*Change, *config.Config) {
		next := cfg.Clone()
		if err := edit(next); err != nil {
			t.Fatal(err)
		}
		ch, err := Build(p, cfg, next, Options{})
		if err != nil {
			t.Fatal(err)
		}
		return ch, next
	}
	first, saved := plan(func(c *config.Config) error { return c.AddGroup("g", []string{"memory"}, []string{"a"}) })
	second, _ := plan(func(c *config.Config) error { return c.AddGroup("h", []string{"memory"}, []string{"b"}) })
	if err := first.Begin(); err != nil {
		t.Fatal(err)
	}
	first.release() // ccam stopped after saving the config
	if err := second.Begin(); !errors.Is(err, ErrPending) {
		t.Fatalf("noted over the first: %v", err)
	}
	if _, err := resumeLocked(p, saved, func(string, string) bool { return false }); err != nil || Pending(p) {
		t.Fatalf("resume: %v", err)
	}
	if err := first.Apply(); err != nil {
		t.Fatalf("carried out again: %v", err)
	}
	if _, ok := fsutil.LinkTarget(filepath.Join(p.ProfileDir("a"), "CLAUDE.md")); !ok {
		t.Fatal("not finished")
	}

	held, err := lock.TryFlock(lockPath(p))
	if err != nil {
		t.Fatal(err)
	}
	defer func(w time.Duration) { lockWait = w }(lockWait)
	lockWait = 50 * time.Millisecond
	if err := second.Begin(); err == nil || !strings.Contains(err.Error(), "another ccam process is changing the groups right now") {
		t.Fatalf("begin while another runs: %v", err)
	}
	if !Changing(p) {
		t.Fatal("Changing while the lock is held")
	}
	held.Unlock()
	if Changing(p) {
		t.Fatal("Changing after it was let go")
	}
}

// A note that cannot be read says how to get past it.
func TestResumeOfAnUnreadableNote(t *testing.T) {
	p, cfg := setup(t, "a")
	os.MkdirAll(p.State, 0o700)
	os.WriteFile(filepath.Join(p.State, "group-change.json"), []byte("{not json"), 0o600)
	if _, err := resumeLocked(p, cfg, func(string, string) bool { return false }); err == nil || !strings.Contains(err.Error(), "remove that file") {
		t.Fatalf("%v", err)
	}
}

// An account that joins a group from no group, in the same change as one
// that gives the group its first servers, sees those servers: its kept
// copies of them go, as when the group had them before.
func TestJoiningSeesTheSeedOfTheSameChange(t *testing.T) {
	p, cfg := setup(t, "a", "b")
	cfg.Default = "a" // planned first
	if err := cfg.AddGroup("g", []string{"mcp"}, []string{"b"}); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(p.GroupDir("g"), 0o700)
	os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"}}}`), 0o600)
	// a left g earlier and kept fs.
	if err := mcp.Keep(p.ProfileDir("a"), p.GroupDir("g"), launcher.MCPStateDir(p, "a")); err != nil {
		t.Fatal(err)
	}
	ch, _ := apply(t, p, cfg, func(c *config.Config) error {
		if err := c.AddGroup("n", []string{"mcp"}, nil); err != nil {
			return err
		}
		for _, m := range []string{"a", "b"} {
			if _, err := c.JoinGroup(m, "n", true); err != nil {
				return err
			}
		}
		return nil
	}, Options{})
	var buf bytes.Buffer
	ch.Describe(&buf)
	if !strings.Contains(buf.String(), "remove   MCP servers fs kept in .claude.json") || strings.Contains(buf.String(), "a's own MCP servers fs") {
		t.Errorf("plan:\n%s", buf.String())
	}
	if _, ok := servers(t, p.ProfileDir("a"))["fs"]; ok {
		t.Error("a's kept copy of fs stays and hides group n's")
	}
	if !strings.Contains(read(t, mcp.Path(p.GroupDir("n"))), `"fs"`) {
		t.Error("group n did not start with g's servers")
	}
}

// A .claude.json ccam cannot read is refused when the change is planned,
// before the config is saved, where the change would edit it: it would
// stop there, and every doctor --fix after it.
func TestUnreadableClaudeJSONIsRefusedInThePlan(t *testing.T) {
	p, cfg := setup(t, "a")
	if err := cfg.AddGroup("g", []string{"mcp"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	if err := cfg.AddGroup("h", []string{"mcp"}, nil); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(p.GroupDir("g"), 0o700)
	os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"}}}`), 0o600)
	dir := p.ProfileDir("a")
	if _, err := mcp.Sync(dir, p.GroupDir("g"), launcher.MCPStateDir(p, "a")); err != nil {
		t.Fatal(err)
	}
	// Cut short, as a full disk can leave it; the scanner still finds
	// mcpServers in front.
	data := strings.TrimRight(read(t, filepath.Join(dir, ".claude.json")), "\n")
	os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(data[:len(data)-1]), 0o600)
	if _, err := mcp.LocalServers(dir); err != nil {
		t.Fatalf("the scanner: %v", err)
	}
	for _, c := range []struct {
		name string
		edit func(*config.Config) error
		opt  Options
	}{
		{"move", func(c *config.Config) error { _, err := c.JoinGroup("a", "h", true); return err }, Options{}},
		{"leave", func(c *config.Config) error { return c.LeaveGroup("a", "g") }, Options{}},
		{"leave --fresh", func(c *config.Config) error { return c.LeaveGroup("a", "g") }, Options{Fresh: true}},
	} {
		next := cfg.Clone()
		if err := c.edit(next); err != nil {
			t.Fatal(err)
		}
		if _, err := Build(p, cfg, next, c.opt); err == nil || !strings.Contains(err.Error(), "repair the file first") {
			t.Errorf("%s: %v", c.name, err)
		}
	}
}

// A stopped `leave --fresh` finishes fresh: the account starts without the
// group's files, not with a copy of them.
func TestResumeOfAFreshLeave(t *testing.T) {
	p, cfg := setup(t, "a", "b")
	_, cfg = apply(t, p, cfg, func(c *config.Config) error {
		return c.AddGroup("g", []string{"memory"}, []string{"a", "b"})
	}, Options{})
	next := cfg.Clone()
	if err := next.LeaveGroup("a", "g"); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, cfg, next, Options{Fresh: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil { // the config is saved, then ccam stops
		t.Fatal(err)
	}
	ch.release()
	if _, err := resumeLocked(p, next, func(string, string) bool { return false }); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(p.ProfileDir("a"), "CLAUDE.md")); err == nil {
		t.Errorf("a has a CLAUDE.md after a fresh leave: %q", read(t, filepath.Join(p.ProfileDir("a"), "CLAUDE.md")))
	}
}

// A change that stops in a plan, after its MCP steps, runs them again when
// doctor --fix finishes it. Keep then still notes the copies it made the
// first time, beside those of servers the group got meanwhile, so joining a
// group again takes all of them out.
func TestResumeRunsKeepAgain(t *testing.T) {
	if !perms {
		t.Skip("needs file permissions")
	}
	p, cfg := setup(t, "a", "b")
	_, cfg = apply(t, p, cfg, func(c *config.Config) error {
		return c.AddGroup("g", []string{"mcp", "memory"}, []string{"a", "b"})
	}, Options{})
	os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"}}}`), 0o600)
	next := cfg.Clone()
	if err := next.LeaveGroup("a", "g"); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, cfg, next, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	memory := filepath.Join(p.GroupDir("g"), "CLAUDE.md")
	os.Chmod(memory, 0) // a's copy of the group's memory cannot be made
	err = ch.Apply()
	os.Chmod(memory, 0o600)
	if err == nil || !Pending(p) {
		t.Fatalf("the change did not stop in a's plan: %v", err)
	}
	if got := servers(t, p.ProfileDir("a")); got["fs"] != "fs" {
		t.Fatalf("Keep did not run before the plan: %v", got)
	}
	// The group gets a server before doctor --fix runs.
	os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"},"new":{"command":"new"}}}`), 0o600)
	if _, err := resumeLocked(p, next, func(string, string) bool { return false }); err != nil {
		t.Fatal(err)
	}
	ch, _ = apply(t, p, next, func(c *config.Config) error { _, err := c.JoinGroup("a", "g", true); return err }, Options{})
	var buf bytes.Buffer
	ch.Describe(&buf)
	if got := servers(t, p.ProfileDir("a")); len(got) != 0 || strings.Contains(buf.String(), "a's own MCP servers") {
		t.Errorf("after joining again a has %v; plan:\n%s", got, buf.String())
	}
}

// The first account brings a plugins folder without a plugin list (a
// marketplace added, nothing installed), the others each bring a list:
// every list is merged, not one left in a backup at random.
func TestEveryPluginListIsMergedAfterASeedWithoutOne(t *testing.T) {
	for i := range 20 {
		p, cfg := setup(t, "a", "b", "c")
		cfg.Default = "a"
		os.MkdirAll(filepath.Join(p.ProfileDir("a"), "plugins", "cache"), 0o700)
		for _, n := range []string{"b", "c"} {
			dir := filepath.Join(p.ProfileDir(n), "plugins")
			os.MkdirAll(dir, 0o700)
			os.WriteFile(filepath.Join(dir, "installed_plugins.json"), []byte(`{"version":2,"plugins":{"`+n+`@m":[{"scope":"user"}]}}`), 0o600)
		}
		apply(t, p, cfg, func(c *config.Config) error { return c.AddGroup("g", []string{"extensions"}, []string{"a", "b", "c"}) }, Options{})
		got := read(t, filepath.Join(p.GroupDir("g"), "plugins", "installed_plugins.json"))
		for _, n := range []string{"b@m", "c@m"} {
			if !strings.Contains(got, n) {
				t.Fatalf("run %d: %s missing from %s", i, n, got)
			}
		}
	}
}

// installThrough records plugin demo@m in the plugin list the way Claude
// Code does when account n installs it: through n's own plugins folder,
// which is its link into the group.
func installThrough(t *testing.T, p paths.Paths, n string) {
	t.Helper()
	dir := filepath.Join(p.ProfileDir(n), "plugins")
	inst := filepath.Join(dir, "cache", "m", "demo", "1.0")
	os.MkdirAll(inst, 0o700)
	os.WriteFile(filepath.Join(inst, "plugin.json"), []byte(`{}`), 0o600)
	q, _ := json.Marshal(inst)
	os.WriteFile(filepath.Join(dir, "installed_plugins.json"), []byte(`{"version":2,"plugins":{"demo@m":[{"scope":"user","installPath":`+string(q)+`}]}}`), 0o600)
}

// installPaths lists the install paths in the plugin list of dir.
func installPaths(t *testing.T, dir string) []string {
	t.Helper()
	var doc struct {
		Plugins map[string][]struct{ InstallPath string } `json:"plugins"`
	}
	if err := json.Unmarshal([]byte(read(t, filepath.Join(dir, "installed_plugins.json"))), &doc); err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, records := range doc.Plugins {
		for _, r := range records {
			out = append(out, r.InstallPath)
		}
	}
	return out
}

// Claude Code records plugin paths through the installing account's link
// into the group. An account that leaves the group (or moves its
// extensions to another) must not take those paths along: the accounts
// that stay would load the plugins from its folder. When the group stops
// sharing extensions, every account's copy points at itself.
func TestLeavingFreesThePluginPaths(t *testing.T) {
	leave := func(c *config.Config) error { return c.LeaveGroup("c", "g") }
	cases := []struct {
		name string
		edit func(*config.Config) error
		opt  Options
		// each lists the accounts that get their own copy of the list
		each []string
	}{
		{"leave", leave, Options{}, []string{"c"}},
		{"leave --fresh", leave, Options{Fresh: true}, nil},
		{"move to another group", func(c *config.Config) error { _, err := c.JoinGroup("c", "y", true); return err }, Options{}, nil},
		{"stop sharing", func(c *config.Config) error { return c.SetGroupItems("g", nil) }, Options{}, []string{"a", "b", "c"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, cfg := setup(t, "a", "b", "c", "d")
			d := filepath.Join(p.ProfileDir("d"), "skills", "s")
			os.MkdirAll(d, 0o700)
			os.WriteFile(filepath.Join(d, "SKILL.md"), []byte("d"), 0o600)
			_, cfg = apply(t, p, cfg, func(c *config.Config) error {
				if err := c.AddGroup("g", []string{"extensions"}, []string{"a", "b", "c"}); err != nil {
					return err
				}
				return c.AddGroup("y", []string{"extensions"}, []string{"d"})
			}, Options{})
			installThrough(t, p, "c")
			next := cfg.Clone()
			if err := tc.edit(next); err != nil {
				t.Fatal(err)
			}
			ch, err := Build(p, cfg, next, tc.opt)
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			ch.Describe(&buf)
			if !strings.Contains(buf.String(), "plugin paths under "+filepath.Join("profiles", "c", "plugins")+" → "+filepath.Join("groups", "g", "plugins")) {
				t.Errorf("the plan does not say so:\n%s", buf.String())
			}
			if err := ch.Apply(); err != nil {
				t.Fatal(err)
			}
			lists := map[string]string{} // whose list → where it is
			for _, n := range []string{"a", "b", "c"} {
				lists[n] = filepath.Join(p.GroupDir("g"), "plugins")
			}
			for _, n := range tc.each {
				lists[n] = filepath.Join(p.ProfileDir(n), "plugins")
			}
			if tc.opt.Fresh || tc.name == "move to another group" {
				delete(lists, "c")
			}
			for n, dir := range lists {
				for _, path := range installPaths(t, dir) {
					if !strings.HasPrefix(path, dir+string(filepath.Separator)) {
						t.Errorf("%s's plugin list records %s, not a path in %s", n, path, dir)
					}
					if _, err := os.Stat(path); err != nil {
						t.Errorf("%s's plugin list records a missing path: %v", n, err)
					}
				}
			}
		})
	}
}

// A change that stopped before its plans (killed right after the config
// was saved) frees the plugin paths when doctor --fix finishes it.
func TestResumeFreesThePluginPaths(t *testing.T) {
	p, cfg := setup(t, "a", "b", "c")
	_, cfg = apply(t, p, cfg, func(c *config.Config) error {
		return c.AddGroup("g", []string{"extensions"}, []string{"a", "b", "c"})
	}, Options{})
	installThrough(t, p, "c")
	next := cfg.Clone()
	if err := next.LeaveGroup("c", "g"); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, cfg, next, Options{Fresh: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	ch.release() // ccam stopped
	if _, err := resumeLocked(p, next, func(string, string) bool { return false }); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(p.GroupDir("g"), "plugins")
	for _, path := range installPaths(t, dir) {
		if !strings.HasPrefix(path, dir+string(filepath.Separator)) {
			t.Errorf("the group's plugin list records %s", path)
		}
	}
}

// Begin holds the lock until Apply: a doctor --fix that loaded the config
// before the save cannot finish the change against it in between, and
// Apply carries the change out.
func TestNoResumeBetweenBeginAndApply(t *testing.T) {
	p, base := setup(t, "a")
	os.MkdirAll(p.GroupDir("g"), 0o700)
	os.WriteFile(mcp.Path(p.GroupDir("g")), []byte(`{"mcpServers":{"fs":{"command":"fs"}}}`), 0o600)
	_, joined := apply(t, p, base, func(c *config.Config) error { return c.AddGroup("g", []string{"mcp", "memory"}, []string{"a"}) }, Options{})
	next := joined.Clone()
	if err := next.LeaveGroup("a", "g"); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, joined, next, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	if _, err := resumeLocked(p, joined, func(string, string) bool { return false }); !errors.Is(err, ErrChanging) || !Pending(p) {
		t.Fatalf("finished between Begin and Apply: %v", err)
	}
	if err := ch.Apply(); err != nil {
		t.Fatal(err)
	}
	if _, link := fsutil.LinkTarget(filepath.Join(p.ProfileDir("a"), "CLAUDE.md")); link || Pending(p) {
		t.Fatalf("Apply did not carry the leave out (link %v, pending %v)", link, Pending(p))
	}
	if got := servers(t, p.ProfileDir("a")); got["fs"] != "fs" {
		t.Fatalf("the group's MCP servers were not kept: %v", got)
	}
}

// Every answer is noted, also one that keeps the group's copy: for
// history.jsonl the first candidate is merging by line, so a stopped
// change finished by Resume must not merge what the user declined.
func TestResumeKeepsTheGroupsHistoryWhenAsked(t *testing.T) {
	p, cfg := setup(t, "a")
	os.WriteFile(filepath.Join(p.ProfileDir("a"), "history.jsonl"), []byte("{\"a\":1}\n"), 0o600)
	if err := cfg.AddGroup("h", []string{"history"}, nil); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(p.GroupDir("h"), 0o700)
	os.WriteFile(filepath.Join(p.GroupDir("h"), "history.jsonl"), []byte("{\"g\":1}\n"), 0o600)
	next := cfg.Clone()
	if _, err := next.JoinGroup("a", "h", true); err != nil {
		t.Fatal(err)
	}
	ch, err := Build(p, cfg, next, Options{})
	if err != nil {
		t.Fatal(err)
	}
	kept := 0
	for _, c := range ch.Conflicts() {
		if c.Entry.Name == "history.jsonl" && c.Choice != share.KeepGroup {
			c.Choice = share.KeepGroup
			kept++
		}
	}
	if kept != 1 {
		t.Fatalf("no history.jsonl conflict that merges by default: %+v", ch.Conflicts())
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	note, _ := os.ReadFile(NotePath(p))
	ch.Abort()
	os.WriteFile(NotePath(p), note, 0o600) // stopped after the save
	if _, err := resumeLocked(p, next, func(string, string) bool { return false }); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(filepath.Join(p.GroupDir("h"), "history.jsonl")); string(data) != "{\"g\":1}\n" {
		t.Fatalf("group history.jsonl %q (note %s)", data, note)
	}
}

// resumeLocked is Resume under the lock, as doctor --fix calls it.
func resumeLocked(p paths.Paths, cfg *config.Config, running func(name, dir string) bool) (*Change, error) {
	l, err := TryLock(p)
	if err != nil {
		return nil, err
	}
	defer l.Unlock()
	return Resume(l, cfg, running)
}
