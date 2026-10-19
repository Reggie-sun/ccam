package tui

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func pressAll(m tea.Model, keys ...string) tea.Model {
	for _, k := range keys {
		m, _ = m.Update(key(k))
	}
	return m
}

func newMigrate() MigratePicker {
	return NewMigratePicker("/h/.claude", "default", map[string]bool{"settings": true, "mcp": true}, DefaultItems())
}

func TestMigratePickerDefaultsAndToggle(t *testing.T) {
	m := pressAll(newMigrate(), "enter").(MigratePicker)
	if !m.Done() || !reflect.DeepEqual(m.Items(), []string{"mcp", "settings", "memory", "extensions"}) {
		t.Fatalf("defaults: done=%v %v", m.Done(), m.Items())
	}
	// Turn mcp off, move to history (5th) and turn it on.
	m = pressAll(newMigrate(), "space", "down", "down", "down", "down", "space", "enter").(MigratePicker)
	if !reflect.DeepEqual(m.Items(), []string{"settings", "memory", "extensions", "history"}) {
		t.Fatalf("toggled: %v", m.Items())
	}
}

func TestMigratePickerCancelAndNothingSelected(t *testing.T) {
	for _, k := range []string{"esc", "q", "ctrl+c"} {
		if m := pressAll(newMigrate(), k).(MigratePicker); m.Done() {
			t.Errorf("%s did not cancel", k)
		}
	}
	m := pressAll(newMigrate(), "space", "down", "space", "down", "space", "down", "space", "enter").(MigratePicker)
	if !m.Done() || len(m.Items()) != 0 {
		t.Fatalf("nothing selected: done=%v %v", m.Done(), m.Items())
	}
}

func TestMigratePickerMarksMissingItems(t *testing.T) {
	view := newMigrate().View().Content
	for _, line := range strings.Split(view, "\n") {
		missing := strings.Contains(line, "not found")
		switch {
		case strings.Contains(line, itemLabels["settings"]) && missing,
			strings.Contains(line, itemLabels["memory"]) && !missing:
			t.Errorf("line %q", line)
		}
	}
	if !strings.Contains(view, "Copy into group default") || !strings.Contains(view, "/h/.claude") {
		t.Fatalf("view:\n%s", view)
	}
}
