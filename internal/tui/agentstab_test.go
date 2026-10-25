package tui

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
)

// sampleAgents: work's open session (newest), personal's background one,
// and side, which could not be asked.
func sampleAgents(now time.Time) AgentList {
	return AgentList{Sessions: []Agent{
		{Account: "work", State: "open (busy)", Folder: "~/code/web", StartedAt: now.Add(-30 * time.Second)},
		{Account: "personal", ID: "49ec0028", Name: "probe", State: "blocked, idle", Folder: "~/code/app", StartedAt: now.Add(-12 * time.Minute)},
	}, Failed: []string{"side: could not list its sessions: exit status 1"}}
}

// longAgents are wider than their columns: a background session first,
// then an open one.
func longAgents() []Agent {
	return []Agent{
		{Account: "a-very-long-profile-name-32chars", ID: "49ec0028", Name: "refactor the billing and invoicing service end to end",
			State: "blocked, waiting: a permission prompt", Folder: "~/code/clients/a-long-client-name/services/billing-and-invoicing/api"},
		{Account: "work", State: "open (waiting): a permission prompt for Bash", Folder: "/tmp/web"},
	}
}

// asked counts the times the dashboard asked for the list.
func asked(f *fake) int {
	n := 0
	for _, c := range f.calls {
		if c == "agents" {
			n++
		}
	}
	return n
}

// A narrower screen drops Started, then Name, then Folder; Account, ID and
// State always stay. Folder is no wider than the widest folder, so Started
// follows it on a wide screen, and Name takes the room Folder leaves.
func TestAgentColumns(t *testing.T) {
	now := time.Now()
	sample := sampleAgents(now).Sessions
	longName := []Agent{{Account: "personal", ID: "49ec0028", Name: longAgents()[0].Name, State: "blocked, idle", Folder: "~/code/app"}}
	// Double-width characters: a cut may fall inside one.
	wide := []Agent{{Account: "personal", ID: "49ec0028", Name: "修复计费服务的测试", State: "blocked, idle", Folder: "~/代码/项目/网站/前端/组件"}}
	cases := []struct {
		name string
		list []Agent
		w    int
		want agentCols
	}{
		{"sample", sample, 200, agentCols{account: 8, name: 5, state: 13, folder: 10, started: true}},
		{"sample", sample, 80, agentCols{account: 8, name: 5, state: 13, folder: 10, started: true}},
		{"sample", sample, 60, agentCols{account: 8, name: 5, state: 13, folder: 10}},
		{"sample", sample, 50, agentCols{account: 8, state: 13, folder: 10}},
		{"sample", sample, 40, agentCols{account: 8, state: 13}},
		{"long", longAgents(), 200, agentCols{account: 16, name: 53, state: 24, folder: 68, started: true}},
		{"long", longAgents(), 120, agentCols{account: 16, name: 25, state: 24, folder: 26, started: true}},
		{"long", longAgents(), 83, agentCols{account: 16, name: 4, state: 24, folder: 10, started: true}},
		{"long", longAgents(), 82, agentCols{account: 16, name: 11, state: 24, folder: 11}},
		{"long", longAgents(), 80, agentCols{account: 16, name: 10, state: 24, folder: 10}},
		{"long", longAgents(), 60, agentCols{account: 16, state: 24}},
		{"long name", longName, 120, agentCols{account: 8, name: 53, state: 13, folder: 10, started: true}},
		{"long name", longName, 100, agentCols{account: 8, name: 40, state: 13, folder: 10, started: true}},
		{"long name", longName, 63, agentCols{account: 8, name: 12, state: 13, folder: 10}},
	}
	for _, c := range cases {
		if got := agentColumns(c.list, c.w); got != c.want {
			t.Errorf("%s at %d: %+v, want %+v", c.name, c.w, got, c.want)
		}
	}
	// Every line fits once Account, ID and State do, and every row lines up
	// under the header; Name and Folder, when shown, get their least and at
	// most their widest value, and are cut only when the row has no column
	// to spare.
	d := dash(newFake(t))
	d.now = func() time.Time { return now }
	for _, c := range []struct {
		list               []Agent
		from, name, folder int // name, folder: the widest ones
	}{{sample, minWidth, 5, 10}, {longAgents(), 56, 53, 68}, {longName, minWidth, 53, 10}, {wide, minWidth, 18, 26}} {
		d.agents = &AgentList{Sessions: c.list}
		for w := c.from; w <= 200; w++ {
			cols := agentColumns(c.list, w)
			if cols.name > 0 && cols.name < len("Name") ||
				cols.folder > 0 && (cols.folder < min(agentFolderMin, c.folder) || cols.folder > c.folder) {
				t.Errorf("%d columns: %+v", w, cols)
			}
			lines, _ := d.agentsView(w)
			for _, l := range lines {
				if lw := ansi.StringWidth(l); lw > w {
					t.Errorf("%d columns: a line of %d: %q", w, lw, ansi.Strip(l))
				}
			}
			for i, a := range c.list {
				row, head := ansi.StringWidth(lines[i+1]), ansi.StringWidth(lines[0])
				if cols.started {
					row, head = row-ansi.StringWidth(shortAgo(a.StartedAt, now)), head-len("Started")
				}
				if row != head {
					t.Errorf("%d columns: a row of %d under a header of %d: %q", w, row, head, ansi.Strip(lines[i+1]))
				}
			}
			if hw := ansi.StringWidth(lines[0]); hw < w && (cols.name > 0 && cols.name < c.name || cols.folder > 0 && cols.folder < c.folder) {
				t.Errorf("%d columns: %+v cut with %d columns to spare", w, cols, w-hw)
			}
		}
	}
}

func TestAgentsTabListsAndAttaches(t *testing.T) {
	f := newFake(t)
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	f.agents = sampleAgents(now)
	d := dash(f)
	d.width, d.height = 80, 30
	d.now = func() time.Time { return now }
	m := press(t, d, "right")
	lines := map[string]bool{}
	for _, l := range strings.Split(screen(m), "\n") {
		lines[strings.TrimRight(l, " ")] = true
	}
	for _, want := range []string{
		" Background agents and open sessions of every account; enter attaches to one.",
		"    Account   ID        Name   State          Folder      Started",
		" ❯  work      -         -      open (busy)    ~/code/web  now",
		"    personal  49ec0028  probe  blocked, idle  ~/code/app  12m",
		" ✘ side: could not list its sessions: exit status 1",
		" enter to attach · x to stop · r to reload · ? for shortcuts",
	} {
		if !lines[want] {
			t.Errorf("missing %q in\n%s", want, screen(m))
		}
	}
	// An open session has no id to attach to.
	m = press(t, m, "enter")
	if r := d.Result(); r.Action != None || !strings.Contains(screen(m), "this session has no id; only a background session can be attached") {
		t.Fatalf("%+v\n%s", r, screen(m))
	}
	press(t, m, "down", "enter")
	if r := d.Result(); r.Action != Launch || r.Profile != "personal" || !slices.Equal(r.Args, []string{"attach", "49ec0028"}) {
		t.Fatalf("%+v", r)
	}

	// The list is read when the tab first opens and again on r.
	f.calls = nil
	m = press(t, dash(f), "right", "left", "right")
	if n := asked(f); n != 1 {
		t.Fatalf("asked %d times: %v", n, f.calls)
	}
	m = press(t, m, "r")
	if n := asked(f); n != 2 {
		t.Fatalf("r: asked %d times: %v", n, f.calls)
	}
	s := screen(press(t, m, "?"))
	for _, want := range []string{"attach the background session", "stop the background session"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in the help:\n%s", want, s)
		}
	}
}

func TestAgentsTabEmptyFailedAndLoading(t *testing.T) {
	f := newFake(t)
	d := dash(f)
	if s := screen(press(t, d, "right", "enter", "x")); !strings.Contains(s, "No sessions. Start one in the background with `ccam <account> --bg \"…\"`.") ||
		d.Result().Action != None || d.dlg != nil {
		t.Fatalf("%+v\n%s", d.Result(), s)
	}
	// A failure is not an empty list.
	f.agents = AgentList{Failed: []string{"claude was not found"}}
	if s := screen(press(t, dash(f), "right")); !strings.Contains(s, " ✘ claude was not found") || strings.Contains(s, "No sessions") {
		t.Fatal(s)
	}

	d = dash(f)
	drain(d, d.fetchAll(false)) // the usage is in: only the list loads
	d.tab, d.loaded[tabAgents] = tabAgents, true
	if s := screen(d); !strings.Contains(s, "Asking every account's Claude Code for its sessions…") || !d.loading() {
		t.Fatalf("loading %v:\n%s", d.loading(), s)
	}
	press(t, d, "enter", "x")
	if d.Result().Action != None || d.dlg != nil {
		t.Fatalf("while it loads: %+v", d.Result())
	}
	// Fewer sessions than before the reload: the cursor comes back to one.
	d.cursor[tabAgents] = 1
	m, _ := d.Update(agentsMsg{d.agentsGen, AgentList{Sessions: []Agent{{Account: "work", ID: "77ab0011", State: "working, busy"}}}})
	if d.cursor[tabAgents] != 0 || d.loading() {
		t.Fatalf("cursor %d, loading %v", d.cursor[tabAgents], d.loading())
	}
	press(t, m, "enter")
	if r := d.Result(); r.Action != Launch || r.Profile != "work" || !slices.Equal(r.Args, []string{"attach", "77ab0011"}) {
		t.Fatalf("%+v", r)
	}
}

// x stops a background session after a question that starts on No, then
// reads the list again, whether the stop worked or not.
func TestAgentsTabStopAsksFirst(t *testing.T) {
	f := newFake(t)
	f.agents = sampleAgents(time.Now())
	d := dash(f)
	d.height = 30
	m := press(t, d, "right", "x")
	if d.dlg != nil || !strings.Contains(screen(m), "this session has no id; only a background session can be stopped") {
		t.Fatal(screen(m))
	}
	m = press(t, m, "down", "x")
	s := screen(m)
	for _, want := range []string{"Stop background session 49ec0028?", "It runs as personal", "No, keep it", "Yes, stop 49ec0028"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in\n%s", want, s)
		}
	}
	m = press(t, m, "enter")
	if slices.Contains(f.calls, "stop personal 49ec0028") || asked(f) != 1 {
		t.Fatalf("enter on the first option must keep it: %v", f.calls)
	}
	m = press(t, m, "x", "down", "enter")
	if !slices.Contains(f.calls, "stop personal 49ec0028") || asked(f) != 2 ||
		!strings.Contains(screen(m), "49ec0028 (personal) stopped; `ccam attach 49ec0028` opens it again") {
		t.Fatalf("%v\n%s", f.calls, screen(m))
	}

	f.calls, f.stopErr = nil, errors.New("could not stop 49ec0028: exit status 1")
	m = press(t, m, "x", "down", "enter")
	if s := screen(m); !slices.Contains(f.calls, "stop personal 49ec0028") || asked(f) != 1 ||
		!strings.Contains(s, "✘ could not stop 49ec0028: exit status 1") || strings.Contains(s, "stopped;") {
		t.Fatalf("%v\n%s", f.calls, s)
	}
}

// Only the newest load's answer counts: an account may answer one late,
// after r or a stop has asked again.
func TestAgentsTabDropsAnOlderAnswer(t *testing.T) {
	f := newFake(t)
	newer := AgentList{Sessions: []Agent{{Account: "work", ID: "77ab0011", State: "working, busy"}}}
	d := dash(f)
	d.height = 30
	_, first := d.Update(key("right")) // the tab opens and asks
	_, second := d.Update(key("r"))    // r asks again before the answer
	f.agents = newer
	m := drain(d, second) // the second answer comes first
	f.agents = sampleAgents(time.Now())
	older, ok := first().(agentsMsg) // the first one last
	if !ok {
		t.Fatal("the tab did not ask for the list")
	}
	m, _ = m.Update(older)
	if s := screen(m); !strings.Contains(s, "77ab0011") || strings.Contains(s, "49ec0028") {
		t.Fatalf("the older answer replaced the newer one:\n%s", s)
	}
	// An older answer while the newest load runs: the tab still waits.
	_, newest := m.Update(key("r"))
	m, _ = m.Update(older)
	if s := screen(m); !strings.Contains(s, "Asking every account's Claude Code for its sessions…") {
		t.Fatalf("the older answer ended the wait:\n%s", s)
	}
	f.agents = newer
	if s := screen(drain(m, newest)); !strings.Contains(s, "77ab0011") || strings.Contains(s, "49ec0028") || d.loading() {
		t.Fatalf("loading %v:\n%s", d.loading(), s)
	}
}
