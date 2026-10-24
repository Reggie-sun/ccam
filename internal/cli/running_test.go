package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/session"
)

// claudeRuns makes the config folder dir look as Claude Code leaves it while
// it runs there without ccam (plain claude, an IDE, a background session):
// it lists this test process in sessions/<pid>.json as Claude Code lists its
// own processes. The returned func removes the record again.
func claudeRuns(t *testing.T, dir string) func() {
	t.Helper()
	if _, ok := session.StartTime(os.Getpid()); !ok {
		t.Skip("no process start times on this platform")
	}
	path := filepath.Join(dir, "sessions", strconv.Itoa(os.Getpid())+".json")
	os.MkdirAll(filepath.Dir(path), 0o700)
	if err := os.WriteFile(path, fmt.Appendf(nil, `{"pid":%d,"kind":"bg","jobId":"49ec0028","startedAt":%d}`, os.Getpid(), time.Now().UnixMilli()), 0o644); err != nil {
		t.Fatal(err)
	}
	return func() { os.Remove(path) }
}

// The commands that move or remove an account's folder refuse while Claude
// Code runs there, also when ccam did not start it, and go ahead once it
// has ended.
func TestGuardsSeeClaudeCodeStartedElsewhere(t *testing.T) {
	c := newTestCLI(t, nil)
	work := c.p.ProfileDir("work")
	ended := claudeRuns(t, work)
	for _, args := range [][]string{
		{"rename", "work", "job"},
		{"remove", "work", "--yes"},
		{"group", "leave", "team", "work", "--yes"},
		{"group", "set", "team", "--add", "memory", "--yes"},
	} {
		_, errOut, code := c.run("", args...)
		if code == 0 || !strings.Contains(errOut, "work is running") {
			t.Fatalf("%v: want a refusal, got %d %q", args, code, errOut)
		}
	}
	_, errOut, _ := c.run("", "rename", "work", "job")
	if !strings.Contains(errOut, "`ccam agents` lists them, `ccam work stop <id>` ends a background one") {
		t.Fatalf("the refusal says how to find the sessions: %q", errOut)
	}
	if _, err := os.Stat(work); err != nil {
		t.Fatal("the folder moved")
	}
	if p, _ := c.cfg().FindGroup("team"); !slices.Contains(p.Members, "work") || slices.Contains(p.Items, "memory") {
		t.Fatalf("groups saved: %+v", p)
	}
	ended()
	c.ok("rename", "work", "job")
}

// mcp import leaves the account's own entries in .claude.json, and changes
// to the group's servers leave its copy, while Claude Code runs there.
func TestMCPLeavesTheClaudeJSONOfARunningAccount(t *testing.T) {
	c := newTestCLI(t, nil)
	c.ok("group", "set", "team", "--add", "mcp", "--yes")
	work := c.p.ProfileDir("work")
	local := filepath.Join(work, ".claude.json")
	os.WriteFile(local, []byte(`{"mcpServers":{"fs":{"command":"fs"}}}`), 0o600)
	ended := claudeRuns(t, work)

	out := c.ok("mcp", "import", "work")
	if !strings.Contains(out, "account work is running, so its own entries stay") {
		t.Fatal(out)
	}
	if b, _ := os.ReadFile(local); !strings.Contains(string(b), `"fs"`) {
		t.Fatalf("entries removed under a running Claude Code: %s", b)
	}

	// The copy made for Remote Control (as a launch makes it).
	if _, err := mcp.Sync(work, c.p.GroupDir("team"), launcher.MCPStateDir(c.p, "work")); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(local)
	out = c.ok("mcp", "add-json", "gh", `{"type":"stdio","command":"gh"}`, "--group", "team")
	if !strings.Contains(out, "account work is running, so the copy in its .claude.json stays") {
		t.Fatal(out)
	}
	if after, _ := os.ReadFile(local); !bytes.Equal(before, after) {
		t.Fatalf("copy rewritten under a running Claude Code: %s", after)
	}
	ended()
	if out := c.ok("mcp", "import", "work"); !strings.Contains(out, "Removed account work's own entries") {
		t.Fatal(out)
	}
}

// doctor --fix only makes missing links in a folder Claude Code runs on.
func TestDoctorFixLeavesAFolderClaudeCodeRunsOn(t *testing.T) {
	c := newTestCLI(t, nil)
	work := c.p.ProfileDir("work")
	defer claudeRuns(t, work)()
	os.MkdirAll(filepath.Join(work, "projects"), 0o700) // a real folder where team's link belongs
	out, _, _ := c.run("", "doctor", "--fix")
	box, _, _ := strings.Cut(out[max(strings.Index(out, "╭─ ✘ work"), 0):], "╰")
	if !strings.Contains(box, "work") || !strings.Contains(box, "is running; quit Claude Code there") {
		t.Fatal(out)
	}
	if info, err := os.Lstat(filepath.Join(work, "projects")); err != nil || info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("moved under a running Claude Code")
	}
}

// The dashboard shows such an account as running and refuses to remove or
// rename it.
func TestDashboardSeesClaudeCodeStartedElsewhere(t *testing.T) {
	c := newTestCLI(t, nil)
	e, _ := c.env()
	b := backend{e: e}
	if b.Running("work") {
		t.Fatal("nothing runs yet")
	}
	defer claudeRuns(t, c.p.ProfileDir("work"))()
	if !b.Running("work") || b.Running("side") {
		t.Fatal("running marks")
	}
	if _, _, err := b.Remove("work"); err == nil || !strings.Contains(err.Error(), "work is running") {
		t.Fatalf("remove: %v", err)
	}
	if _, err := b.Rename("work", "job"); err == nil || !strings.Contains(err.Error(), "work is running") {
		t.Fatalf("rename: %v", err)
	}
	if _, ok := c.cfg().Find("work"); !ok {
		t.Fatal("work is gone")
	}
}

// startsOnRead is stdin that finds the account started when its question is
// first answered: started in another terminal while the question waited.
type startsOnRead struct {
	start func()
	r     io.Reader
}

func (s *startsOnRead) Read(b []byte) (int, error) {
	if s.start != nil {
		s.start()
		s.start = nil
	}
	return s.r.Read(b)
}

// startWork registers this process as a session of work, as a launch does.
func startWork(t *testing.T, c *testCLI) func() {
	return func() {
		dir := launcher.SessionDir(c.p, "work")
		if err := session.Register(dir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { session.Unregister(dir) })
	}
}

func (c *testCLI) runWith(stdin io.Reader, args ...string) (string, string, int) {
	var out, errb bytes.Buffer
	code := ExecuteWith(args, c.p, stdin, &out, &errb)
	return out.String(), errb.String(), code
}

// remove checks again after its question: the account may have started
// while it waited.
func TestRemoveChecksAgainAfterTheQuestion(t *testing.T) {
	c := newTestCLI(t, nil)
	if _, ok := session.StartTime(os.Getpid()); !ok {
		t.Skip("no process start times on this platform")
	}
	out, errOut, code := c.runWith(&startsOnRead{start: startWork(t, c), r: strings.NewReader("y\n")}, "remove", "work")
	if code == 0 || !strings.Contains(errOut, "work is running") {
		t.Fatalf("want a refusal, got %d %q %q", code, out, errOut)
	}
	if _, err := os.Stat(c.p.ProfileDir("work")); err != nil {
		t.Fatal("the folder went to the trash")
	}
	if _, ok := c.cfg().Find("work"); !ok {
		t.Fatal("the account was forgotten")
	}
}

// A group change checks again after its questions, before it saves or moves
// anything.
func TestGroupChangeChecksAgainAfterTheQuestion(t *testing.T) {
	c := newTestCLI(t, nil)
	if _, ok := session.StartTime(os.Getpid()); !ok {
		t.Skip("no process start times on this platform")
	}
	work := c.p.ProfileDir("work")
	os.WriteFile(filepath.Join(work, "CLAUDE.md"), []byte("mine"), 0o600)
	out, errOut, code := c.runWith(&startsOnRead{start: startWork(t, c), r: strings.NewReader("y\n")}, "group", "set", "team", "--add", "memory")
	if code == 0 || !strings.Contains(out, "Proceed?") || !strings.Contains(errOut, "work is running") {
		t.Fatalf("want a refusal after the question, got %d %q %q", code, out, errOut)
	}
	if g, _ := c.cfg().FindGroup("team"); slices.Contains(g.Items, "memory") {
		t.Fatal("the groups were saved")
	}
	if info, err := os.Lstat(filepath.Join(work, "CLAUDE.md")); err != nil || info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("CLAUDE.md moved")
	}
}

// The dashboard's plan dialog may stay open: applying checks again that no
// account started and that the groups are still those the plan was made
// from.
func TestDashboardGroupPlanIsCheckedAgain(t *testing.T) {
	c := newTestCLI(t, nil)
	if _, ok := session.StartTime(os.Getpid()); !ok {
		t.Skip("no process start times on this platform")
	}
	e, _ := c.env()
	b := backend{e: e}
	addMemory := func(cfg *config.Config) error {
		g, err := cfg.MustFindGroup("team")
		if err != nil {
			return err
		}
		return cfg.SetGroupItems(g.Name, append(slices.Clone(g.Items), "memory"))
	}
	ch, err := b.PlanGroups(addMemory)
	if err != nil {
		t.Fatal(err)
	}
	startWork(t, c)()
	if _, err := b.ApplyGroups(addMemory, ch); err == nil || !strings.Contains(err.Error(), "work is running") {
		t.Fatalf("started while the dialog was open: %v", err)
	}
	session.Unregister(launcher.SessionDir(c.p, "work"))

	ch, err = b.PlanGroups(addMemory)
	if err != nil {
		t.Fatal(err)
	}
	c.ok("group", "leave", "team", "side", "--yes") // another ccam, meanwhile
	if _, err := b.ApplyGroups(addMemory, ch); err == nil || !strings.Contains(err.Error(), "the groups changed while this plan waited") {
		t.Fatalf("groups changed while the dialog was open: %v", err)
	}
	if g, _ := c.cfg().FindGroup("team"); slices.Contains(g.Items, "memory") {
		t.Fatal("the stale plan was saved")
	}
}

// A group without members that another ccam changes while a join waits at
// its question: the join stops instead of saving shares it never moved.
func TestGroupJoinStopsWhenTheGroupChangedMeanwhile(t *testing.T) {
	c := newTestCLI(t, nil)
	c.ok("group", "add", "docs", "--items", "settings,memory")
	work := c.p.ProfileDir("work")
	os.WriteFile(filepath.Join(work, "CLAUDE.md"), []byte("private"), 0o600)
	other := func() { c.ok("group", "set", "docs", "--remove", "memory", "--yes") }
	out, errOut, code := c.runWith(&startsOnRead{start: other, r: strings.NewReader("y\n")}, "group", "join", "docs", "work")
	if code == 0 || !strings.Contains(errOut, "the groups changed while this plan waited") {
		t.Fatalf("want a refusal, got %d %q %q", code, out, errOut)
	}
	if g, _ := c.cfg().FindGroup("docs"); slices.Contains(g.Members, "work") {
		t.Fatal("work joined docs")
	}
	if info, err := os.Lstat(filepath.Join(work, "CLAUDE.md")); err != nil || info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("CLAUDE.md moved into the group")
	}
}
