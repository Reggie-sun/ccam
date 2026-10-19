package tui

import (
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// MigratePicker asks which items `ccam migrate` copies.
type MigratePicker struct {
	source, group string
	found         map[string]bool // items the source has something for
	list          checklist
	done          bool // false after Esc: cancelled
}

// NewMigratePicker starts with checked pre-selected.
func NewMigratePicker(source, group string, found map[string]bool, checked []string) MigratePicker {
	return MigratePicker{source: source, group: group, found: found, list: newChecklist(checked)}
}

// Items is valid when Done reports true; it may be empty.
func (m MigratePicker) Items() []string { return m.list.items() }

func (m MigratePicker) Done() bool { return m.done }

func (m MigratePicker) Init() tea.Cmd { return nil }

func (m MigratePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok || m.list.key(k.String()) {
		return m, nil
	}
	switch k.String() {
	case "ctrl+c", "esc", "q":
		return m, tea.Quit
	case "enter":
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m MigratePicker) View() tea.View {
	var b strings.Builder
	b.WriteString(titleStyle.Render(" Copy into group "+m.group) + "\n\n")
	b.WriteString(" From " + m.source + "\n ? Choose what to copy (Space to toggle)\n")
	m.list.view(&b, func(item string) string {
		if !m.found[item] {
			return "not found"
		}
		return ""
	})
	b.WriteString("\n" + dimStyle.Render(" Enter show the plan · Esc cancel") + "\n")
	return tea.NewView(b.String())
}

// RunMigratePicker runs the picker on the terminal. ok is false when the
// user cancelled.
func RunMigratePicker(source, group string, found map[string]bool, checked []string) (items []string, ok bool, err error) {
	final, err := tea.NewProgram(NewMigratePicker(source, group, found, checked), tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout)).Run()
	if err != nil {
		return nil, false, err
	}
	m := final.(MigratePicker)
	return m.Items(), m.Done(), nil
}
