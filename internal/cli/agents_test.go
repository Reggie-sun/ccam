package cli

import (
	"errors"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/jhyunwoo/ccam/internal/tui"
)

func TestAgentStateText(t *testing.T) {
	for _, c := range []struct {
		a    agentEntry
		want string
	}{
		{agentEntry{Kind: "background", State: "blocked", Status: "idle"}, "blocked, idle"},
		{agentEntry{Kind: "background", State: "working", Status: "busy"}, "working, busy"},
		{agentEntry{Kind: "background", State: "working"}, "working"},
		{agentEntry{Kind: "background", Status: "busy"}, "busy"}, // a record without a job id
		{agentEntry{Kind: "interactive", Status: "busy"}, "open (busy)"},
		{agentEntry{Kind: "interactive"}, "open (interactive)"},
		{agentEntry{Kind: "interactive", Status: "waiting", WaitingFor: "permission"}, "open (waiting): permission"},
		{agentEntry{Kind: "background", State: "blocked", Status: "waiting", WaitingFor: "permission"}, "blocked, waiting: permission"},
		{agentEntry{}, "-"},
	} {
		if got := c.a.stateText(); got != c.want {
			t.Errorf("%+v: %q, want %q", c.a, got, c.want)
		}
	}
}

// The dashboard's Agents tab lists what `ccam agents` lists, stops a
// background session as its account, and attaches to one the way `ccam
// attach` does: the dashboard ends with that launch.
func TestDashboardAgentsTab(t *testing.T) {
	c := newTestCLI(t, nil)
	home := t.TempDir()
	t.Setenv("HOME", home)
	stops := filepath.Join(t.TempDir(), "stops")
	t1, t2 := time.Now().Add(-12*time.Minute).UnixMilli(), time.Now().UnixMilli()
	scriptClaude(t, c, strings.NewReplacer("HOME", home, "T1", strconv.FormatInt(t1, 10), "T2", strconv.FormatInt(t2, 10), "STOPS", stops).Replace(
		`case "$1 ${CLAUDE_CONFIG_DIR##*/}" in
"agents personal") echo '[{"pid":11,"id":"49ec0028","cwd":"HOME/code/app","kind":"background","startedAt":T1,"sessionId":"s1","name":"probe","status":"idle","state":"blocked"}]' ;;
"agents work") echo '[{"pid":12,"cwd":"/tmp/web","kind":"interactive","startedAt":T2,"status":"busy"}]' ;;
"agents side") echo 'not json' ;;
"stop personal") echo "$2" >> 'STOPS'; echo "stopped $2" ;;
"stop side") sleep 3 & echo "stopped $2" ;;
*) echo "No job matching '$2'" >&2; exit 1 ;;
esac`))
	e, cfg := c.env()
	b := backend{e: e}

	list := b.Agents()
	want := []tui.Agent{
		{Account: "work", State: "open (busy)", Folder: "/tmp/web", StartedAt: time.UnixMilli(t2)},
		{Account: "personal", ID: "49ec0028", Name: "probe", State: "blocked, idle", Folder: filepath.Join("~", "code", "app"), StartedAt: time.UnixMilli(t1)},
	}
	if len(list.Sessions) != len(want) {
		t.Fatalf("%+v", list)
	}
	for i, w := range want {
		got := list.Sessions[i]
		if !got.StartedAt.Equal(w.StartedAt) {
			t.Errorf("%s started at %v, want %v", got.Account, got.StartedAt, w.StartedAt)
		}
		got.StartedAt, w.StartedAt = time.Time{}, time.Time{}
		if got != w {
			t.Errorf("%+v, want %+v", got, w)
		}
	}
	if len(list.Failed) != 1 || !strings.HasPrefix(list.Failed[0], "side: could not list its sessions: unexpected output of claude agents --json") {
		t.Fatalf("%q", list.Failed)
	}

	if err := b.StopSession("personal", "49ec0028"); err != nil || strings.TrimSpace(readFile(t, stops)) != "49ec0028" {
		t.Fatalf("%v", err)
	}
	if err := b.StopSession("work", "77ab0011"); err == nil || !strings.Contains(err.Error(), "could not stop 77ab0011:") || !strings.Contains(err.Error(), "No job matching") {
		t.Fatalf("%v", err)
	}
	e.dryRun = true
	if err := b.StopSession("personal", "49ec0028"); !errors.Is(err, errDryRun) || strings.TrimSpace(readFile(t, stops)) != "49ec0028" {
		t.Fatalf("a dry run: %v", err)
	}
	e.dryRun = false
	// A process claude leaves behind holds its output open; the dashboard,
	// which takes no key while it stops a session, does not wait for it.
	wait := agentsWaitDelay
	agentsWaitDelay = 100 * time.Millisecond
	t.Cleanup(func() { agentsWaitDelay = wait })
	start := time.Now()
	if err := b.StopSession("side", "aa11bb22"); err != nil || time.Since(start) > 2*time.Second {
		t.Fatalf("%v after %v", err, time.Since(start))
	}

	d := tui.NewDashboard(b, cfg, "", "dev")
	m := runCmds(d, d.Init())
	for _, k := range []tea.KeyPressMsg{{Code: tea.KeyRight}, {Code: tea.KeyDown}, {Code: tea.KeyEnter}} {
		var cmd tea.Cmd
		m, cmd = m.Update(k)
		m = runCmds(m, cmd)
	}
	if r := d.Result(); r.Action != tui.Launch || r.Profile != "personal" || !slices.Equal(r.Args, []string{"attach", "49ec0028"}) {
		t.Fatalf("%+v", r)
	}
}
