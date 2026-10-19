package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/jhyunwoo/ccam/internal/config"
)

var itemLabels = map[string]string{
	"mcp":        "MCP servers",
	"settings":   "settings.json",
	"memory":     "global CLAUDE.md",
	"extensions": "agents, commands, skills",
	"history":    "conversation history",
}

// shareEdit toggles share items on a copy; nothing moves until the caller
// shows the migration plan and the user confirms.
type shareEdit struct {
	orig   config.Share
	next   config.Share
	cursor int
}

func newShareEdit(s config.Share) shareEdit { return shareEdit{orig: s, next: s} }

func (e *shareEdit) toggle() {
	item := config.Items[e.cursor]
	if e.next.Get(item) != config.Local {
		_ = e.next.Set(item, config.Local)
		return
	}
	// Turning back on restores the original group, else "default".
	if g := e.orig.Get(item); g != config.Local {
		_ = e.next.Set(item, g)
	} else {
		_ = e.next.Set(item, "default")
	}
}

func (m Model) updateShare(k tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.mode = modeList
	case "up", "k":
		if m.share.cursor > 0 {
			m.share.cursor--
		}
	case "down", "j":
		if m.share.cursor < len(config.Items)-1 {
			m.share.cursor++
		}
	case "space", " ", "x":
		m.share.toggle()
	case "enter":
		if m.share.next.Equal(m.share.orig) {
			m.mode = modeList
			return m, nil
		}
		m.result = Result{Action: Share, Profile: m.rows[m.cursor].Name, Share: m.share.next}
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) viewShare(b *strings.Builder) {
	fmt.Fprintf(b, "%s\n\n", titleStyle.Render(" Sharing for "+m.rows[m.cursor].Name))
	for i, item := range config.Items {
		marker := "  "
		if i == m.share.cursor {
			marker = cursorStyle.Render("› ")
		}
		g := m.share.next.Get(item)
		box, state := "[x]", "group "+g
		if g == config.Local {
			box, state = "[ ]", "local"
		}
		changed := ""
		if g != m.share.orig.Get(item) {
			changed = warnStyle.Render(" (changed)")
		}
		fmt.Fprintf(b, "%s%s %s %s%s\n", marker, box, padRight(itemLabels[item], 24), dimStyle.Render(state), changed)
	}
	b.WriteString("\n" + dimStyle.Render(" Space on/off · Enter apply (shows the migration plan, then asks) · Esc cancel") + "\n")
}
