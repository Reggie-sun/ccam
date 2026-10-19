// Package tui is the `ccam` selection screen and the `ccam add` wizard.
package tui

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// Action is what the user picked; the caller performs it after the
// program exits so the terminal is restored before claude starts.
type Action int

const (
	None Action = iota
	Launch
	Add
	Share
	SetDefault
)

type Result struct {
	Action  Action
	Profile string
	Share   config.Share // for Share
}

// Row is one account as shown in the list.
type Row struct {
	Name     string
	Slot     int
	Email    string
	Share    config.Share
	LastUsed string
	Default  bool
}

// UsageFunc fetches one account's usage; nil disables the column.
type UsageFunc func(ctx context.Context, name string) usage.Report

type mode int

const (
	modeList mode = iota
	modeShare
)

type usageMsg struct {
	name   string
	report usage.Report
}

// Model is the picker. Build it with NewPicker.
type Model struct {
	rows    []Row
	cursor  int
	mode    mode
	width   int
	usage   map[string]*usage.Report
	fetch   UsageFunc
	share   shareEdit
	result  Result
	message string
}

func NewPicker(rows []Row, fetch UsageFunc) Model {
	m := Model{rows: rows, width: 80, usage: map[string]*usage.Report{}, fetch: fetch}
	for i, r := range rows {
		if r.Default {
			m.cursor = i
		}
	}
	return m
}

// Result is valid after the program finished.
func (m Model) Result() Result { return m.result }

func (m Model) Init() tea.Cmd { return m.fetchAll() }

func (m Model) fetchAll() tea.Cmd {
	if m.fetch == nil {
		return nil
	}
	cmds := make([]tea.Cmd, 0, len(m.rows))
	for _, r := range m.rows {
		name := r.Name
		cmds = append(cmds, func() tea.Msg {
			return usageMsg{name: name, report: m.fetch(context.Background(), name)}
		})
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case usageMsg:
		r := msg.report
		m.usage[msg.name] = &r
	case tea.KeyPressMsg:
		if m.mode == modeShare {
			return m.updateShare(msg)
		}
		return m.updateList(msg)
	}
	return m, nil
}

func (m Model) updateList(k tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := k.String()
	m.message = ""
	switch key {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.rows)-1 {
			m.cursor++
		}
	case "enter":
		if len(m.rows) > 0 {
			m.result = Result{Action: Launch, Profile: m.rows[m.cursor].Name}
			return m, tea.Quit
		}
	case "a":
		m.result = Result{Action: Add}
		return m, tea.Quit
	case "d":
		if len(m.rows) > 0 {
			m.result = Result{Action: SetDefault, Profile: m.rows[m.cursor].Name}
			return m, tea.Quit
		}
	case "s":
		if len(m.rows) > 0 {
			m.mode = modeShare
			m.share = newShareEdit(m.rows[m.cursor].Share)
		}
	case "r":
		m.usage = map[string]*usage.Report{}
		return m, m.fetchAll()
	default:
		// Number keys launch that slot directly.
		if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
			for _, r := range m.rows {
				if fmt.Sprint(r.Slot) == key {
					m.result = Result{Action: Launch, Profile: r.Name}
					return m, tea.Quit
				}
			}
			m.message = "no account with number " + key
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	var b strings.Builder
	switch m.mode {
	case modeShare:
		m.viewShare(&b)
	default:
		m.viewList(&b)
	}
	v := tea.NewView(b.String())
	v.WindowTitle = "ccam"
	return v
}

var (
	titleStyle  = lipgloss.NewStyle().Bold(true)
	cursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	fullStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

func (m Model) viewList(b *strings.Builder) {
	b.WriteString(titleStyle.Render(" ccam · Claude Code accounts") + "\n\n")
	if len(m.rows) == 0 {
		b.WriteString(" No accounts yet. Press a to add one.\n")
	}
	nameW, emailW := 4, 6
	for _, r := range m.rows {
		nameW = max(nameW, ansi.StringWidth(r.Name))
		emailW = max(emailW, ansi.StringWidth(orDash(r.Email)))
	}
	emailW = min(emailW, 26)
	for i, r := range m.rows {
		marker := "  "
		if i == m.cursor {
			marker = cursorStyle.Render("› ")
		}
		def := " "
		if r.Default {
			def = "*"
		}
		line := fmt.Sprintf("%s%s%d %s  %s  %s  %s  %s", marker, def, r.Slot,
			padRight(r.Name, nameW),
			padRight(ansi.Truncate(orDash(r.Email), emailW, "…"), emailW),
			m.usageCell(r.Name),
			dimStyle.Render(padRight(r.LastUsed, 8)),
			dimStyle.Render(shareSummary(r.Share)))
		b.WriteString(ansi.Truncate(line, m.width-1, "…") + "\n")
	}
	b.WriteString("\n")
	if m.message != "" {
		b.WriteString(" " + warnStyle.Render(m.message) + "\n")
	}
	b.WriteString(dimStyle.Render(ansi.Truncate(" Enter launch · 1-9 jump · a add · s sharing · d default · r refresh · q quit", m.width-1, "…")) + "\n")
}

// usageCell is a fixed 8-column cell: "5h  72%", "5h …", "API", "-".
func (m Model) usageCell(name string) string {
	if m.fetch == nil {
		return padRight("", 8)
	}
	r, ok := m.usage[name]
	switch {
	case !ok:
		return dimStyle.Render(padRight("5h  …", 8))
	case r.Status == usage.StatusAPI:
		return padRight("API", 8)
	case r.Limits == nil || !r.Limits.FiveHour.Known:
		return dimStyle.Render(padRight("5h   -", 8))
	}
	p := r.Limits.FiveHour.Percent
	cell := padRight(fmt.Sprintf("5h %3.0f%%", p), 8)
	switch {
	case p >= 100:
		return fullStyle.Render(cell)
	case p >= 80:
		return warnStyle.Render(cell)
	}
	return cell
}

func padRight(s string, w int) string {
	if gap := w - ansi.StringWidth(s); gap > 0 {
		return s + strings.Repeat(" ", gap)
	}
	return s
}

func orDash(s string) string {
	if s == "" {
		return "(not logged in)"
	}
	return s
}

func shareSummary(s config.Share) string {
	short := map[string]string{"mcp": "mcp", "settings": "settings", "memory": "md", "extensions": "ext", "history": "history"}
	var parts []string
	for _, item := range config.Items {
		g := s.Get(item)
		if g == config.Local {
			continue
		}
		p := short[item]
		if g != "default" {
			p += "@" + g
		}
		parts = append(parts, p)
	}
	if len(parts) == 0 {
		return "local"
	}
	return strings.Join(parts, " ")
}
