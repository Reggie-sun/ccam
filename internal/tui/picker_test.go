package tui

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/teatest/v2"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func rows() []Row {
	return []Row{
		{Name: "personal", Slot: 1, Email: "me@example.com", Share: config.DefaultShare(), LastUsed: "just now", Default: true},
		{Name: "work", Slot: 2, Email: "dev@company.com", Share: config.DefaultShare(), LastUsed: "2h ago"},
		{Name: "side", Slot: 3, Share: config.Share{}, LastUsed: "-"},
	}
}

func fakeUsage(pct map[string]float64) UsageFunc {
	return func(ctx context.Context, name string) usage.Report {
		p, ok := pct[name]
		if !ok {
			return usage.Report{Name: name, Status: usage.StatusLoggedOut}
		}
		return usage.Report{Name: name, Status: usage.StatusOK, Limits: &usage.Limits{FiveHour: usage.Window{Percent: p, Known: true}}}
	}
}

func run(t *testing.T, m Model, keys ...tea.KeyPressMsg) (Model, string) {
	t.Helper()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("personal")) },
		teatest.WithDuration(2*time.Second))
	for _, k := range keys {
		tm.Send(k)
	}
	final := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second)).(Model)
	out, _ := io.ReadAll(tm.FinalOutput(t))
	return final, string(out)
}

func key(s string) tea.KeyPressMsg {
	switch s {
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	case "up":
		return tea.KeyPressMsg{Code: tea.KeyUp}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	case "space":
		return tea.KeyPressMsg{Code: tea.KeySpace, Text: " "}
	}
	return tea.KeyPressMsg{Code: rune(s[0]), Text: s}
}

func TestPickerStartsOnDefaultAndLaunches(t *testing.T) {
	m, _ := run(t, NewPicker(rows(), nil), key("down"), key("enter"))
	if r := m.Result(); r.Action != Launch || r.Profile != "work" {
		t.Fatalf("%+v", r)
	}
}

func TestPickerNumberKeyLaunchesSlot(t *testing.T) {
	m, _ := run(t, NewPicker(rows(), nil), key("3"))
	if r := m.Result(); r.Action != Launch || r.Profile != "side" {
		t.Fatalf("%+v", r)
	}
}

func TestPickerQuitAndAdd(t *testing.T) {
	m, _ := run(t, NewPicker(rows(), nil), key("q"))
	if m.Result().Action != None {
		t.Fatal("q should do nothing")
	}
	m, _ = run(t, NewPicker(rows(), nil), key("a"))
	if m.Result().Action != Add {
		t.Fatal("a should add")
	}
	m, _ = run(t, NewPicker(rows(), nil), key("down"), key("d"))
	if r := m.Result(); r.Action != SetDefault || r.Profile != "work" {
		t.Fatalf("%+v", r)
	}
}

func TestPickerShareToggle(t *testing.T) {
	// personal: history is local; move to it (5th item) and turn it on.
	m, out := run(t, NewPicker(rows(), nil),
		key("s"), key("down"), key("down"), key("down"), key("down"), key("space"), key("enter"))
	r := m.Result()
	if r.Action != Share || r.Profile != "personal" || r.Share.Get("history") != "default" || r.Share.Get("mcp") != "default" {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(out, "conversation history") {
		t.Fatal("share screen not shown")
	}
	// Toggling off and on again is no change: Enter just goes back.
	m, _ = run(t, NewPicker(rows(), nil), key("s"), key("space"), key("space"), key("enter"), key("q"))
	if m.Result().Action != None {
		t.Fatalf("no-op share change: %+v", m.Result())
	}
	// Esc leaves the share screen without a result.
	m, _ = run(t, NewPicker(rows(), nil), key("s"), key("space"), key("esc"), key("q"))
	if m.Result().Action != None {
		t.Fatalf("%+v", m.Result())
	}
}

func TestPickerShowsUsageWhenItArrives(t *testing.T) {
	m := NewPicker(rows(), fakeUsage(map[string]float64{"personal": 72, "work": 100}))
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool {
		return bytes.Contains(b, []byte("72%")) && bytes.Contains(b, []byte("100%"))
	}, teatest.WithDuration(2*time.Second))
	tm.Send(key("q"))
	tm.WaitFinished(t, teatest.WithFinalTimeout(2*time.Second))
}

func TestFirstFrameIsFast(t *testing.T) {
	slow := func(ctx context.Context, name string) usage.Report {
		time.Sleep(2 * time.Second) // usage must not block the first frame
		return usage.Report{Name: name}
	}
	start := time.Now()
	tm := teatest.NewTestModel(t, NewPicker(rows(), slow), teatest.WithInitialTermSize(80, 24))
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("side")) },
		teatest.WithDuration(time.Second), teatest.WithCheckInterval(time.Millisecond))
	if took := time.Since(start); took > 100*time.Millisecond {
		t.Fatalf("first frame took %v", took)
	}
	tm.Send(key("q"))
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func TestViewsFit80Columns(t *testing.T) {
	long := rows()
	long[0].Name = "a-very-long-profile-name-32chars"
	long[0].Email = "someone.with.a.long.address@example-company.com"
	long[0].Share = config.Share{MCP: "team-alpha", Settings: "team-alpha", Memory: "team-alpha", Extensions: "team-alpha", History: "team-alpha"}
	m := NewPicker(long, fakeUsage(nil))
	m.usage["work"] = &usage.Report{Status: usage.StatusOK, Limits: &usage.Limits{FiveHour: usage.Window{Percent: 100, Known: true}}}
	for _, view := range []string{m.View().Content, func() string {
		mm, _ := m.Update(key("s"))
		return mm.(Model).View().Content
	}()} {
		for _, line := range strings.Split(view, "\n") {
			if w := ansi.StringWidth(line); w > 80 {
				t.Errorf("line is %d columns: %q", w, ansi.Strip(line))
			}
		}
	}
}

func TestWizardHelpers(t *testing.T) {
	cfg := &config.Config{Profiles: []config.Profile{{Name: "work", Slot: 1}}}
	v := ValidateNewName(cfg)
	if v("work") == nil || v("add") == nil || v("Bad") == nil || v("new-one") != nil {
		t.Fatal("validation")
	}
	s := ShareFromItems(DefaultItems())
	if !s.Equal(config.DefaultShare()) {
		t.Fatalf("default items %+v", s)
	}
}

func TestWizardNameThenItems(t *testing.T) {
	cfg := &config.Config{Profiles: []config.Profile{{Name: "work", Slot: 1}}}
	tm := teatest.NewTestModel(t, NewWizard(cfg, ""), teatest.WithInitialTermSize(80, 24))
	tm.Type("work")
	tm.Send(key("enter"))
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("already exists")) },
		teatest.WithDuration(2*time.Second))
	for i := 0; i < 4; i++ {
		tm.Send(tea.KeyPressMsg{Code: tea.KeyBackspace})
	}
	tm.Type("side")
	tm.Send(key("enter"))
	// Turn history on (5th item), turn mcp off (1st), confirm.
	tm.Send(key("space"))
	for i := 0; i < 4; i++ {
		tm.Send(key("down"))
	}
	tm.Send(key("space"))
	tm.Send(key("enter"))
	w := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second)).(Wizard)
	a := w.Answers()
	if !w.Done() || a.Name != "side" || a.Share.Get("mcp") != config.Local || a.Share.Get("history") != "default" || a.Share.Get("settings") != "default" {
		t.Fatalf("%+v", a)
	}
}

func TestWizardEscCancels(t *testing.T) {
	tm := teatest.NewTestModel(t, NewWizard(&config.Config{}, "x1"), teatest.WithInitialTermSize(80, 24))
	tm.Send(key("esc"))
	if w := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second)).(Wizard); w.Done() {
		t.Fatal("esc must cancel")
	}
}

func shortcutWizard(t *testing.T) *teatest.TestModel {
	check := func(name, shortcut string) error {
		if shortcut == "list" {
			return errors.New(`"list" is taken by ccam or Claude Code`)
		}
		return nil
	}
	w := NewWizard(&config.Config{}, "work").WithShortcut(func(name string) string { return "cc" + name }, check)
	tm := teatest.NewTestModel(t, w, teatest.WithInitialTermSize(80, 24))
	tm.Send(key("enter")) // keep the default items
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool {
		return bytes.Contains(b, []byte("Create a shortcut command for work? [y/N]"))
	}, teatest.WithDuration(2*time.Second))
	return tm
}

func TestWizardShortcutDefaultsToNo(t *testing.T) {
	tm := shortcutWizard(t)
	tm.Send(key("enter"))
	w := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second)).(Wizard)
	if a := w.Answers(); !w.Done() || a.Shortcut != "" || a.Name != "work" {
		t.Fatalf("%+v", a)
	}
}

func TestWizardShortcutNamed(t *testing.T) {
	tm := shortcutWizard(t)
	tm.Type("y")
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("Shortcut name [ccwork]:")) },
		teatest.WithDuration(2*time.Second))
	tm.Type("list")
	tm.Send(key("enter"))
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("is taken by ccam")) },
		teatest.WithDuration(2*time.Second))
	tm.Type("ccjob")
	tm.Send(key("enter"))
	w := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second)).(Wizard)
	if a := w.Answers(); !w.Done() || a.Shortcut != "ccjob" {
		t.Fatalf("%+v", a)
	}
}

func TestWizardShortcutSuggestion(t *testing.T) {
	tm := shortcutWizard(t)
	tm.Type("y")
	tm.Send(key("enter"))
	w := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second)).(Wizard)
	if a := w.Answers(); !w.Done() || a.Shortcut != "ccwork" {
		t.Fatalf("%+v", a)
	}
}

func TestWizardShortcutWithoutUsableSuggestion(t *testing.T) {
	check := func(name, shortcut string) error {
		if shortcut == "ccwork" {
			return errors.New(`"ccwork" would hide /usr/bin/ccwork`)
		}
		return nil
	}
	w := NewWizard(&config.Config{}, "work").WithShortcut(func(name string) string { return "cc" + name }, check)
	tm := teatest.NewTestModel(t, w, teatest.WithInitialTermSize(80, 24))
	tm.Send(key("enter"))
	tm.Type("y")
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("Shortcut name: ")) },
		teatest.WithDuration(2*time.Second))
	tm.Send(key("enter")) // nothing typed, no default
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool { return bytes.Contains(b, []byte("type a name for the shortcut")) },
		teatest.WithDuration(2*time.Second))
	tm.Type("ccjob")
	tm.Send(key("enter"))
	final := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second)).(Wizard)
	if a := final.Answers(); !final.Done() || a.Shortcut != "ccjob" {
		t.Fatalf("%+v", a)
	}
}
