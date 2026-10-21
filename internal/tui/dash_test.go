package tui

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/teatest/v2"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/doctor"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// fake is a Backend over a config in memory; group changes are planned
// with the real engine in a scratch home.
type fake struct {
	p       paths.Paths
	cfg     *config.Config
	usage   map[string]usage.Report
	slow    time.Duration
	running map[string]bool
	calls   []string
	newer   string // what NewerRelease answers
	updErr  error
}

func newFake(t *testing.T) *fake {
	t.Setenv("CCAM_HOME", t.TempDir())
	cfg := &config.Config{Default: "personal", Profiles: []config.Profile{
		{Name: "personal", Slot: 1, Email: "me@example.com"},
		{Name: "work", Slot: 2, Email: "dev@company.com"},
		{Name: "side", Slot: 3},
	}}
	must(t, cfg.AddGroup("default", config.DefaultItems, []string{"personal", "work"}))
	must(t, cfg.AddGroup("team", []string{"history"}, []string{"work"}))
	return &fake{p: paths.FromEnv(), cfg: cfg, usage: map[string]usage.Report{}, running: map[string]bool{}}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func (f *fake) Usage(ctx context.Context, name string, fresh bool) usage.Report {
	time.Sleep(f.slow)
	if r, ok := f.usage[name]; ok {
		return r
	}
	return usage.Report{Name: name, Status: usage.StatusLoggedOut}
}
func (f *fake) Running(name string) bool { return f.running[name] }
func (f *fake) Sessions(name string) int {
	if f.running[name] {
		return 1
	}
	return 0
}
func (f *fake) Takeover(from string) []string {
	var out []string
	for _, p := range f.cfg.Profiles {
		if p.Name != from {
			out = append(out, p.Name)
		}
	}
	return out
}
func (f *fake) LastUsed(string) time.Time           { return time.Time{} }
func (f *fake) SuggestShortcut(name string) string  { return "cc" + name }
func (f *fake) CheckShortcut(name, sc string) error { return nil }
func (f *fake) Extensions() (ExtView, error)        { return ExtView{}, nil }
func (f *fake) Doctor(fix bool, only string) ([]doctor.Finding, error) {
	f.calls = append(f.calls, fmt.Sprint("doctor ", fix, " ", only))
	return []doctor.Finding{{Level: doctor.OK, Msg: "fine"}}, nil
}

func (f *fake) NewerRelease(context.Context) string { return f.newer }
func (f *fake) Update(context.Context) error {
	f.calls = append(f.calls, "update")
	return f.updErr
}

func (f *fake) change(call string, edit func(*config.Config) error) (*config.Config, error) {
	f.calls = append(f.calls, call)
	next := f.cfg.Clone()
	if err := edit(next); err != nil {
		return nil, err
	}
	f.cfg = next
	return next, nil
}

func (f *fake) SetDefault(name string) (*config.Config, error) {
	return f.change("default "+name, func(c *config.Config) error { c.Default = name; return nil })
}
func (f *fake) SetFolder(dir, account string) (*config.Config, error) {
	return f.change("folder "+account, func(c *config.Config) error { c.SetFolder(dir, account); return nil })
}
func (f *fake) UnsetFolder(dir string) (*config.Config, error) {
	return f.change("unfolder", func(c *config.Config) error { c.UnsetFolder(dir); return nil })
}
func (f *fake) Rename(o, n string) (*config.Config, error) {
	return f.change("rename "+o+" "+n, func(c *config.Config) error { p, _ := c.Find(o); p.Name = n; return c.Resolve() })
}
func (f *fake) SetShortcut(a, s string) (*config.Config, error) {
	return f.change("shortcut "+a+" "+s, func(c *config.Config) error { p, _ := c.Find(a); p.Shortcut = s; return nil })
}
func (f *fake) Remove(name string) (*config.Config, string, error) {
	cfg, err := f.change("remove "+name, func(c *config.Config) error { c.Remove(name); return nil })
	return cfg, name + " removed", err
}
func (f *fake) PlanGroups(edit func(*config.Config) error) (*group.Change, error) {
	next := f.cfg.Clone()
	if err := edit(next); err != nil {
		return nil, err
	}
	return group.Build(f.p, f.cfg, next, group.Options{})
}
func (f *fake) ApplyGroups(edit func(*config.Config) error, ch *group.Change) (*config.Config, error) {
	return f.change("groups", edit)
}
func (f *fake) RemoveGroup(name string, ch *group.Change) (*config.Config, string, error) {
	cfg, err := f.change("remove group "+name, func(c *config.Config) error { return c.RemoveGroup(name) })
	return cfg, "group " + name + " removed", err
}

func dash(f *fake) *Dashboard {
	d := NewDashboard(f, f.cfg, "/home/me/code", "v0.9.0")
	d.tickEvery = 0
	return d
}

func key(s string) tea.KeyPressMsg {
	switch s {
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	case "up":
		return tea.KeyPressMsg{Code: tea.KeyUp}
	case "right":
		return tea.KeyPressMsg{Code: tea.KeyRight}
	case "left":
		return tea.KeyPressMsg{Code: tea.KeyLeft}
	case "tab":
		return tea.KeyPressMsg{Code: tea.KeyTab}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	case "space":
		return tea.KeyPressMsg{Code: tea.KeySpace, Text: " "}
	case "backspace":
		return tea.KeyPressMsg{Code: tea.KeyBackspace}
	}
	return tea.KeyPressMsg{Code: rune(s[0]), Text: s}
}

// press sends keys and runs the commands they return, as the program
// would, except that it stops at Quit.
func press(t *testing.T, m tea.Model, keys ...string) tea.Model {
	t.Helper()
	for _, k := range keys {
		var cmd tea.Cmd
		m, cmd = m.Update(key(k))
		m = drain(m, cmd)
	}
	return m
}

func drain(m tea.Model, cmd tea.Cmd) tea.Model {
	if cmd == nil {
		return m
	}
	switch msg := cmd().(type) {
	case nil, tea.QuitMsg:
	case tea.BatchMsg:
		for _, c := range msg {
			m = drain(m, c)
		}
	default:
		var next tea.Cmd
		m, next = m.Update(msg)
		m = drain(m, next)
	}
	return m
}

func typeText(t *testing.T, m tea.Model, s string) tea.Model {
	t.Helper()
	for _, r := range s {
		m = press(t, m, string(r))
	}
	return m
}

func screen(m tea.Model) string { return ansi.Strip(m.(*Dashboard).render()) }

func TestStartsOnDefaultAndLaunches(t *testing.T) {
	f := newFake(t)
	m := press(t, dash(f), "down", "enter")
	if r := m.(*Dashboard).Result(); r.Action != Launch || r.Profile != "work" {
		t.Fatalf("%+v", r)
	}
	f.cfg.Folders = []config.Folder{{Path: "/home/me", Account: "side"}}
	m = press(t, dash(f), "enter")
	if r := m.(*Dashboard).Result(); r.Profile != "side" {
		t.Fatalf("the folder's account first: %+v", r)
	}
}

// Digits move the cursor; only Enter launches. Two digits make one number.
func TestDigitsMoveTheCursor(t *testing.T) {
	f := newFake(t)
	for i := 4; i <= 12; i++ {
		f.cfg.Profiles = append(f.cfg.Profiles, config.Profile{Name: fmt.Sprintf("a%d", i), Slot: i})
	}
	d := dash(f)
	m := press(t, d, "3")
	if d.Result().Action != None || d.current().Name != "side" {
		t.Fatalf("on %s, %+v", d.current().Name, d.Result())
	}
	press(t, m, "1", "2")
	if d.current().Name != "a12" {
		t.Fatalf("1 then 2: on %s", d.current().Name)
	}
	press(t, m, "9", "9") // no 99: the second 9 starts again
	if d.current().Name != "a9" {
		t.Fatalf("on %s", d.current().Name)
	}
	press(t, m, "0")
	if !strings.Contains(screen(m), "no account with number 0") {
		t.Fatal(screen(m))
	}
	press(t, m, "enter")
	if r := d.Result(); r.Action != Launch {
		t.Fatalf("%+v", r)
	}
}

func TestDefaultLoginAndQuit(t *testing.T) {
	f := newFake(t)
	m := press(t, dash(f), "down", "down", "d")
	if f.cfg.Default != "side" || !strings.Contains(screen(m), "side is the default account") {
		t.Fatalf("%v\n%s", f.calls, screen(m))
	}
	m = press(t, dash(f), "l")
	if r := m.(*Dashboard).Result(); r.Action != Login || r.Profile != "side" {
		t.Fatalf("%+v", r)
	}
	f.cfg.Profiles[0].API = &config.API{KeyFrom: "env:K"}
	m = press(t, dash(f), "up", "up", "l")
	if r := m.(*Dashboard).Result(); r.Action != None || !strings.Contains(screen(m), "API key account") {
		t.Fatalf("%+v", r)
	}
	m = press(t, dash(f), "q")
	if r := m.(*Dashboard).Result(); r.Action != None {
		t.Fatalf("%+v", r)
	}
}

func TestMostRoomAndLaunchOptions(t *testing.T) {
	f := newFake(t)
	w := func(p float64) usage.Window { return usage.Window{Percent: p, Known: true} }
	f.usage["personal"] = usage.Report{Name: "personal", Status: usage.StatusOK, Limits: &usage.Limits{FiveHour: w(70), Weekly: w(10)}}
	f.usage["work"] = usage.Report{Name: "work", Status: usage.StatusOK, Limits: &usage.Limits{FiveHour: w(20), Weekly: w(100)}}
	f.usage["side"] = usage.Report{Name: "side", Status: usage.StatusOK, Limits: &usage.Limits{FiveHour: w(40), Weekly: w(30)}}
	d := dash(f)
	m := drain(d, d.Init())
	press(t, m, "n")
	if r := d.Result(); r.Action != Launch || r.Profile != "side" {
		t.Fatalf("work's week is full, so side: %+v", r)
	}
	d = dash(f)
	press(t, d, "o", "2")
	if r := d.Result(); r.Profile != "personal" || !slices.Equal(r.Args, []string{"--continue"}) {
		t.Fatalf("%+v", r)
	}
	d = dash(f)
	m = press(t, d, "o", "4")
	typeText(t, m, "-p hi")
	press(t, m, "enter")
	if r := d.Result(); !slices.Equal(r.Args, []string{"-p", "hi"}) {
		t.Fatalf("%+v", r)
	}
}

func TestAccountGroupsGoThroughThePlan(t *testing.T) {
	f := newFake(t)
	m := press(t, dash(f), "down", "down", "s") // side, in no group
	if !strings.Contains(screen(m), "Groups of side") {
		t.Fatal(screen(m))
	}
	m = press(t, m, "down", "space", "enter") // join team
	if !strings.Contains(screen(m), "Do you want to apply this plan?") || !strings.Contains(screen(m), "+history (group team)") {
		t.Fatal(screen(m))
	}
	m = press(t, m, "enter")
	if g, _ := f.cfg.FindGroup("team"); !g.HasMember("side") || !strings.Contains(screen(m), "side is in team") {
		t.Fatalf("%v\n%s", f.calls, screen(m))
	}
	// Two groups that share the same item are refused before any plan.
	must(t, f.cfg.AddGroup("pair", []string{"settings"}, nil))
	m = press(t, dash(f), "s", "down", "down", "space", "enter")
	if !strings.Contains(screen(m), "both share settings") {
		t.Fatal(screen(m))
	}
}

func TestGroupsTabCreatesAndEdits(t *testing.T) {
	f := newFake(t)
	m := press(t, dash(f), "right")
	if s := screen(m); !strings.Contains(s, "default") || !strings.Contains(s, "no group") || !strings.Contains(s, "side") {
		t.Fatal(s)
	}
	m = press(t, m, "n")
	m = typeText(t, m, "docs")
	m = press(t, m, "enter")
	if !strings.Contains(screen(m), "new group docs") {
		t.Fatal(screen(m))
	}
	// memory: work gets it from default, so it cannot be added for work;
	// side can have it.
	m = press(t, m, "tab", "down", "down", "space", "tab", "down", "down", "space", "enter")
	if !strings.Contains(screen(m), "Do you want to apply this plan?") {
		t.Fatal(screen(m))
	}
	press(t, m, "enter")
	g, ok := f.cfg.FindGroup("docs")
	if !ok || !slices.Equal(g.Items, []string{"memory"}) || !slices.Equal(g.Members, []string{"side"}) {
		t.Fatalf("%+v", g)
	}
	// Editing team: adding settings is blocked while work gets it from default.
	m = press(t, dash(f), "right", "down", "enter", "down", "space")
	if !strings.Contains(screen(m), "settings (work: from default)") {
		t.Fatal(screen(m))
	}
	if g, _ := f.cfg.FindGroup("team"); g.Has("settings") {
		t.Fatal("blocked item toggled")
	}
}

func TestRemoveStartsOnNo(t *testing.T) {
	f := newFake(t)
	m := press(t, dash(f), "down", "down", "x")
	if !strings.Contains(screen(m), "Remove account side?") {
		t.Fatal(screen(m))
	}
	press(t, m, "enter")
	if len(f.calls) != 0 {
		t.Fatalf("enter on the first option must keep it: %v", f.calls)
	}
	m = press(t, dash(f), "down", "down", "x", "down", "enter")
	if _, ok := f.cfg.Find("side"); ok || !strings.Contains(screen(m), "side removed") {
		t.Fatalf("%v", f.calls)
	}
	// A group with accounts: they leave it first.
	m = press(t, dash(f), "right", "down", "x")
	if s := screen(m); !strings.Contains(s, "Remove group team?") || !strings.Contains(s, "work leave") {
		t.Fatal(s)
	}
	press(t, m, "down", "enter")
	if _, ok := f.cfg.FindGroup("team"); ok {
		t.Fatal("not removed")
	}
}

func TestRenameAndShortcut(t *testing.T) {
	f := newFake(t)
	m := press(t, dash(f), "e")
	for range "personal" {
		m = press(t, m, "backspace")
	}
	m = typeText(t, m, "home")
	m = press(t, m, "enter")
	if _, ok := f.cfg.Find("home"); !ok || !strings.Contains(screen(m), "personal is now home") {
		t.Fatalf("%v", f.calls)
	}
	press(t, m, "c", "enter")
	if p, _ := f.cfg.Find("home"); p.Shortcut != "cchome" {
		t.Fatalf("suggested shortcut: %v", f.calls)
	}
}

func TestTabsHelpAndHealth(t *testing.T) {
	f := newFake(t)
	m := press(t, dash(f), "?")
	if !strings.Contains(screen(m), "Shortcuts") {
		t.Fatal(screen(m))
	}
	m = press(t, m, "x", "left") // any key closes the help; left wraps to Health
	if s := screen(m); !strings.Contains(s, "fine") || !strings.Contains(s, "everything looks fine") {
		t.Fatal(s)
	}
	m = press(t, m, "f") // the only box has nothing to repair
	if !strings.Contains(screen(m), "nothing to repair in General") || slices.Contains(f.calls, "doctor true ") {
		t.Fatalf("%v\n%s", f.calls, screen(m))
	}
	m = press(t, m, "F", "enter")
	if !slices.Contains(f.calls, "doctor true ") || !strings.Contains(screen(m), "repaired") {
		t.Fatalf("%v", f.calls)
	}
}

func TestAddFromTheDashboard(t *testing.T) {
	f := newFake(t)
	d := dash(f)
	m := press(t, d, "a")
	m = typeText(t, m, "laptop")
	press(t, m, "enter", "3", "enter", "n", "enter") // log in later, default group, no shortcut
	r := d.Result()
	if r.Action != Add || r.Add.Name != "laptop" || r.Add.Login || !slices.Equal(r.Add.Groups, []string{"default"}) {
		t.Fatalf("%+v", r)
	}
}

func TestFirstFrameIsFastAndUsageArrives(t *testing.T) {
	f := newFake(t)
	f.slow = 300 * time.Millisecond
	f.usage["personal"] = usage.Report{Name: "personal", Status: usage.StatusOK, Limits: &usage.Limits{FiveHour: usage.Window{Percent: 72, Known: true}}}
	start := time.Now()
	tm := teatest.NewTestModel(t, dash(f), teatest.WithInitialTermSize(80, 24))
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("side")) },
		teatest.WithDuration(time.Second), teatest.WithCheckInterval(time.Millisecond))
	if took := time.Since(start); took > 100*time.Millisecond {
		t.Fatalf("first frame took %v", took)
	}
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("72%")) },
		teatest.WithDuration(2*time.Second))
	tm.Send(key("q"))
	tm.WaitFinished(t, teatest.WithFinalTimeout(2*time.Second))
}

// Every screen fits the terminal: at 80 columns, and at 60 with columns left out.
func TestScreensFit(t *testing.T) {
	f := newFake(t)
	f.cfg.Profiles[0].Name = "a-very-long-profile-name-32chars"
	f.cfg.Profiles[0].Email = "someone.with.a.long.address@example-company.com"
	must(t, f.cfg.Resolve())
	f.cfg.Default = ""
	f.cfg.SetFolder("/home/me/code/clients/a-long-client-name/services/billing-and-invoicing/api", "work")
	f.usage["work"] = usage.Report{Status: usage.StatusOK, Limits: &usage.Limits{FiveHour: usage.Window{Percent: 100, Known: true}}}
	screens := [][]string{
		{}, {"?"}, {"u"}, {"o"}, {"s"}, {"x"}, {"a"}, {"e"},
		{"right"}, {"right", "enter"}, {"right", "n"}, {"right", "down", "x"},
		{"right", "right"}, {"right", "right", "n"}, {"left", "left"}, {"left"}, {"left", "f"},
	}
	for _, width := range []int{80, 60} {
		for _, keys := range screens {
			d := dash(f)
			d.Update(tea.WindowSizeMsg{Width: width, Height: 24})
			m := drain(d, d.Init())
			m = press(t, m, keys...)
			for _, line := range strings.Split(m.(*Dashboard).render(), "\n") {
				if w := ansi.StringWidth(line); w > width {
					t.Errorf("%d columns after %v: line is %d: %q", width, keys, w, ansi.Strip(line))
				}
			}
			if n := strings.Count(m.(*Dashboard).render(), "\n") + 1; n > 24 {
				t.Errorf("%v: %d lines for 24 rows", keys, n)
			}
		}
	}
}

func TestEmptyStateOffersAdd(t *testing.T) {
	f := newFake(t)
	f.cfg = &config.Config{}
	m := press(t, dash(f))
	if !strings.Contains(screen(m), "No accounts yet") {
		t.Fatal(screen(m))
	}
	m = press(t, m, "a")
	if !strings.Contains(screen(m), "Account name") {
		t.Fatal(screen(m))
	}
}

// A shortened folder path keeps the account column under its heading.
func TestFolderColumnsLineUp(t *testing.T) {
	f := newFake(t)
	f.cfg.SetFolder("/home/me/code/clients/a-long-client-name/services/billing-and-invoicing/api", "work")
	f.cfg.SetFolder("/home/me/side", "side")
	d := dash(f)
	d.Update(tea.WindowSizeMsg{Width: 60, Height: 24})
	m := press(t, drain(d, d.Init()), "right", "right")
	col := -1
	for _, line := range strings.Split(screen(m), "\n") {
		switch i := strings.Index(line, "Account"); {
		case strings.Contains(line, "Folder") && i >= 0:
			col = ansi.StringWidth(line[:i])
		case strings.HasSuffix(strings.TrimSpace(line), "work") || strings.HasSuffix(strings.TrimSpace(line), "side"):
			if got := ansi.StringWidth(line[:strings.LastIndex(line, " ")+1]); got != col {
				t.Errorf("account at column %d, heading at %d: %q", got, col, line)
			}
		}
	}
	if col < 0 {
		t.Fatal(screen(m))
	}
}

// The Health tab draws a box per account with a row per area, so where a
// problem is shows at a glance; ↓ moves to the next box.
func TestHealthBoxPerAccount(t *testing.T) {
	f := newFake(t)
	d := dash(f)
	d.width, d.height = 80, 40
	d.tab, d.loaded[tabHealth] = tabHealth, true
	d.health = []doctor.Finding{
		{Level: doctor.OK, Area: "Claude Code", Msg: "Claude Code 2.1.304"},
		{Level: doctor.OK, Profile: "personal", Area: "Login", Msg: "login valid"},
		{Level: doctor.OK, Profile: "personal", Area: "Share links", Msg: "share links OK"},
		{Level: doctor.Info, Profile: "personal", Area: "Folder", Msg: "unknown entries: x.json"},
		{Level: doctor.OK, Profile: "work", Area: "Login", Msg: "login valid"},
		{Level: doctor.Error, Profile: "work", Area: "Share links", Msg: "settings.json is a real file"},
	}
	s := screen(d)
	for _, want := range []string{
		"Checks Claude Code and each account's folder",
		"✔ General   ✔ personal   ✘ work",
		"╭─ ✘ work · 1 problem ─",
		"│ ✔ Folder       OK                                                          │\n │ ℹ              unknown entries: x.json",
		"│ ✘ Share links  settings.json is a real file",
		"1 problem; f repairs",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in\n%s", want, s)
		}
	}
	for _, l := range strings.Split(s, "\n") {
		if w := ansi.StringWidth(l); w > 80 {
			t.Errorf("wider than 80 (%d): %s", w, l)
		}
	}
	m := press(t, d, "down", "down")
	if c := m.(*Dashboard).cursor[tabHealth]; c != 2 {
		t.Fatalf("cursor %d", c)
	}
	// f repairs only the selected account's box.
	m = press(t, m, "f")
	if s := screen(m); !strings.Contains(s, "Repair what can be repaired in work?") {
		t.Fatal(s)
	}
	m = press(t, m, "enter")
	if !slices.Contains(f.calls, "doctor true work") || !strings.Contains(screen(m), "repaired what could be repaired in work") {
		t.Fatalf("%v\n%s", f.calls, screen(m))
	}
}

func TestEveryTabSaysWhatItIsFor(t *testing.T) {
	if len(tabAbout) != len(tabNames) {
		t.Fatal("a tab has no description")
	}
	for _, a := range tabAbout {
		if len(a) > 79 {
			t.Errorf("wider than the screen: %s", a)
		}
	}
}

func TestAccountsShowFiveHourWeeklyAndFable(t *testing.T) {
	f := newFake(t)
	w := func(p float64) usage.Window { return usage.Window{Percent: p, Known: true} }
	f.usage["personal"] = usage.Report{Name: "personal", Status: usage.StatusOK, Limits: &usage.Limits{
		FiveHour: w(50), Weekly: w(25), Scoped: []usage.Scoped{{Name: "Fable", Window: w(100)}}}}
	f.usage["work"] = usage.Report{Name: "work", Status: usage.StatusOK, Limits: &usage.Limits{FiveHour: w(2)}}
	d := dash(f)
	d.width, d.height = 100, 30
	m := drain(d, d.fetchAll(false))
	s := screen(m)
	for _, want := range []string{"5-hour       Weekly       Fable", "███░░░  50%  █▌░░░░  25%  ██████ 100%", "  2%  -            -"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in\n%s", want, s)
		}
	}
}

// `ccam doctor` into a pipe: plain text, nothing wrapped or cut, and every
// box as wide as the widest.
func TestHealthReportForAPipe(t *testing.T) {
	long := strings.Repeat("x", 120)
	lines := HealthReport([]doctor.Finding{
		{Level: doctor.OK, Area: "Claude Code", Msg: "2.1.304"},
		{Level: doctor.Warn, Profile: "a", Area: "Share links", Msg: long},
	}, 0, false)
	out := strings.Join(lines, "\n")
	if strings.Contains(out, "\x1b[") || !strings.Contains(out, long) || !strings.Contains(out, "✔ General   ⚠ a") {
		t.Fatal(out)
	}
	w := ansi.StringWidth(lines[2])
	for _, l := range lines[2:] {
		if ansi.StringWidth(l) != w {
			t.Fatalf("boxes of different widths:\n%s", out)
		}
	}
}

// A newer release is recommended under the title; U asks, then updates.
func TestOffersANewerRelease(t *testing.T) {
	f := newFake(t)
	m := tea.Model(dash(f))
	m = drain(m, m.Init())
	if s := screen(m); strings.Contains(s, "is available") {
		t.Fatalf("no newer release, but:\n%s", s)
	}
	m = press(t, m, "U")
	if m.(*Dashboard).dlg != nil {
		t.Fatal("U without a newer release opened a dialog")
	}

	f.newer = "v0.11.0"
	d := dash(f)
	m = drain(d, d.Init())
	if s := screen(m); !strings.Contains(s, "● ccam 0.11.0 is available (you have v0.9.0) · press U to update") {
		t.Fatal(s)
	}
	m = press(t, m, "U")
	if s := screen(m); !strings.Contains(s, "Update ccam v0.9.0 → 0.11.0?") {
		t.Fatal(s)
	}
	m = press(t, m, "enter")
	s := screen(m)
	if !slices.Contains(f.calls, "update") || !strings.Contains(s, "ccam 0.11.0 is installed; quit and start ccam again") || strings.Contains(s, "is available") {
		t.Fatalf("%v\n%s", f.calls, s)
	}

	// A failed update says why and keeps the offer.
	f.updErr = fmt.Errorf("/usr/bin/ccam was installed by a package")
	d = dash(f)
	m = press(t, drain(d, d.Init()), "U", "enter")
	if s := screen(m); !strings.Contains(s, "installed by a package") || !strings.Contains(s, "is available") {
		t.Fatal(s)
	}
}
