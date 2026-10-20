package tui

import (
	"errors"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/jhyunwoo/ccam/internal/config"
)

func wizardView(m tea.Model) string { return ansi.Strip(m.(Wizard).View().Content) }

func TestWizardNameKindItems(t *testing.T) {
	cfg := &config.Config{Profiles: []config.Profile{{Name: "work", Slot: 1}}}
	m := typeText(t, NewWizard(cfg, ""), "work")
	m = press(t, m, "enter")
	if !strings.Contains(wizardView(m), "already exists") {
		t.Fatal(wizardView(m))
	}
	for range 4 {
		m = press(t, m, "backspace")
	}
	m = typeText(t, m, "side")
	// Log in now; no groups yet, so items for group default: turn mcp off
	// and history on.
	m = press(t, m, "enter", "enter")
	if !strings.Contains(wizardView(m), "group default") {
		t.Fatal(wizardView(m))
	}
	m = press(t, m, "space", "down", "down", "down", "down", "space", "enter")
	if !strings.Contains(wizardView(m), "Add this account?") {
		t.Fatal(wizardView(m))
	}
	m = press(t, m, "enter")
	w := m.(Wizard)
	a := w.Answers()
	if !w.Done() || a.Name != "side" || !a.Login || strings.Join(a.Groups, ",") != "default" ||
		strings.Join(a.Items, ",") != "settings,memory,extensions,history" {
		t.Fatalf("%+v", a)
	}
}

func TestWizardAPIKeyAccount(t *testing.T) {
	m := press(t, NewWizard(&config.Config{}, "gw"), "2")
	m = press(t, m, "enter")
	if !strings.Contains(wizardView(m), "env:NAME or command:COMMAND") {
		t.Fatalf("an empty key source is refused:\n%s", wizardView(m))
	}
	m = typeText(t, m, "env:GW_KEY")
	m = press(t, m, "enter")
	m = typeText(t, m, "https://llm.example.com")
	m = press(t, m, "enter", "2", "enter", "enter", "enter")
	w := m.(Wizard)
	a := w.Answers()
	if !w.Done() || a.Login || a.API == nil || a.API.KeyFrom != "env:GW_KEY" || a.API.Auth != config.Bearer || a.API.BaseURL != "https://llm.example.com" {
		t.Fatalf("%+v %+v", a, a.API)
	}
}

func TestWizardEscGoesBack(t *testing.T) {
	m := press(t, NewWizard(&config.Config{}, "x1"), "3", "esc")
	if w := m.(Wizard); w.step != stepKind || w.cancelled {
		t.Fatalf("back to the kind step: %v", w.step)
	}
	m = press(t, m, "esc")
	if w := m.(Wizard); !w.cancelled || w.Done() {
		t.Fatal("esc on the first step cancels")
	}
	// No on the summary cancels too.
	m = press(t, NewWizard(&config.Config{}, "x1"), "3", "enter", "2")
	if w := m.(Wizard); !w.cancelled || w.Done() {
		t.Fatal("no on the summary")
	}
}

func TestWizardOffersGroups(t *testing.T) {
	cfg := &config.Config{Profiles: []config.Profile{{Name: "work", Slot: 1}}}
	must(t, cfg.AddGroup("default", config.DefaultItems, []string{"work"}))
	must(t, cfg.AddGroup("team", []string{"history"}, nil))
	must(t, cfg.AddGroup("pair", []string{"settings"}, nil))
	m := press(t, NewWizard(cfg, "side"), "1", "down", "space", "enter", "enter")
	if a := m.(Wizard).Answers(); !m.(Wizard).Done() || strings.Join(a.Groups, ",") != "default,team" || a.Items != nil {
		t.Fatalf("%+v", a)
	}
	m = press(t, NewWizard(cfg, "side"), "1", "down", "down", "space", "enter")
	if !strings.Contains(wizardView(m), "both share settings") {
		t.Fatal(wizardView(m))
	}
	m = press(t, NewWizard(cfg, "side"), "1", "space", "enter", "enter")
	if a := m.(Wizard).Answers(); len(a.Groups) != 0 {
		t.Fatalf("no group: %+v", a)
	}
}

func shortcutWizard() tea.Model {
	check := func(name, shortcut string) error {
		if shortcut == "list" || shortcut == "ccbusy" {
			return errors.New(`"` + shortcut + `" is taken`)
		}
		return nil
	}
	return NewWizard(&config.Config{}, "work").WithShortcut(func(name string) string { return "cc" + name }, check)
}

func TestWizardShortcut(t *testing.T) {
	m := press(t, shortcutWizard(), "3", "enter")
	if !strings.Contains(wizardView(m), "Create a shortcut command for work? [y/N]") || !strings.Contains(wizardView(m), "ccwork") {
		t.Fatal(wizardView(m))
	}
	// No by default.
	w := press(t, m, "enter", "enter").(Wizard)
	if a := w.Answers(); !w.Done() || a.Shortcut != "" {
		t.Fatalf("%+v", a)
	}
	// Yes takes the suggestion, which can be edited.
	w = press(t, m, "y", "enter", "enter").(Wizard)
	if a := w.Answers(); a.Shortcut != "ccwork" {
		t.Fatalf("%+v", a)
	}
	m2 := press(t, m, "y")
	for range "ccwork" {
		m2 = press(t, m2, "backspace")
	}
	m2 = typeText(t, m2, "list")
	m2 = press(t, m2, "enter")
	if !strings.Contains(wizardView(m2), "is taken") {
		t.Fatal(wizardView(m2))
	}
	for range "list" {
		m2 = press(t, m2, "backspace")
	}
	m2 = typeText(t, m2, "ccjob")
	w = press(t, m2, "enter", "enter").(Wizard)
	if a := w.Answers(); a.Shortcut != "ccjob" {
		t.Fatalf("%+v", a)
	}
}

func TestWizardAcceptsPastedName(t *testing.T) {
	m, _ := NewWizard(&config.Config{}, "").Update(tea.PasteMsg{Content: "side\n"})
	if w := m.(Wizard); w.Answers().Name != "side" {
		t.Fatalf("pasted name: %q", w.Answers().Name)
	}
}

func TestValidateNewName(t *testing.T) {
	v := ValidateNewName(&config.Config{Profiles: []config.Profile{{Name: "work", Slot: 1}}})
	if v("work") == nil || v("add") == nil || v("Bad") == nil || v("group") == nil || v("new-one") != nil {
		t.Fatal("validation")
	}
	if !slices.Equal(DefaultItems(), config.DefaultItems) {
		t.Fatal("default items")
	}
}
