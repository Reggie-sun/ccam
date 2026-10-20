package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// A dialog takes the keys while it is open. On a choice it sets d.dlg to
// nil before running what was chosen, which may open the next dialog.
type dialog interface {
	update(d *Dashboard, msg tea.Msg) tea.Cmd
	view(d *Dashboard, w int) []string
	keys() string
}

func (d *Dashboard) updateDialog(msg tea.Msg) tea.Cmd { return d.dlg.update(d, msg) }

// dialogTop is a rule in the accent color and the title under it.
func dialogTop(t Theme, w int, title string) []string {
	return []string{" " + t.accent().Render(strings.Repeat("─", w-2)), "  " + t.accent().Bold(true).Render(title), ""}
}

// wrap breaks text into lines of at most w columns.
func wrap(text string, w int) []string {
	var out []string
	for _, para := range strings.Split(text, "\n") {
		out = append(out, strings.Split(lipgloss.NewStyle().Width(max(10, w)).Render(para), "\n")...)
	}
	return out
}

// options draws a numbered choice list with the cursor.
func options(t Theme, labels []string, cursor int) []string {
	out := make([]string, len(labels))
	for i, l := range labels {
		if i == cursor {
			out[i] = "  " + t.accent().Render(symCursor+" "+strconv.Itoa(i+1)+". "+l)
		} else {
			out[i] = "    " + t.dim().Render(strconv.Itoa(i+1)+".") + " " + l
		}
	}
	return out
}

// pickKey moves cursor over n options or picks one by number; it reports
// whether an option was picked.
func pickKey(k string, cursor *int, n int) bool {
	switch k {
	case "up", "k":
		*cursor = max(0, *cursor-1)
	case "down", "j":
		*cursor = min(n-1, *cursor+1)
	case "enter":
		return true
	default:
		if i, err := strconv.Atoi(k); err == nil && i >= 1 && i <= n {
			*cursor = i - 1
			return true
		}
	}
	return false
}

// confirmDialog asks a yes/no question. A destructive one is drawn in a
// box in the error color and starts on No.
type confirmDialog struct {
	title       string
	body        string
	yes, no     string
	destructive bool
	cursor      int // 0 is the first option shown
	onYes       func(d *Dashboard) tea.Cmd
}

func newConfirm(title, body, yes string, onYes func(d *Dashboard) tea.Cmd) *confirmDialog {
	return &confirmDialog{title: title, body: body, yes: yes, no: "No", onYes: onYes}
}

func newDestructive(title, body, yes string, onYes func(d *Dashboard) tea.Cmd) *confirmDialog {
	return &confirmDialog{title: title, body: body, yes: yes, no: "No, keep it", destructive: true, onYes: onYes}
}

func (c *confirmDialog) labels() []string {
	if c.destructive {
		return []string{c.no, c.yes}
	}
	return []string{c.yes, c.no}
}

func (c *confirmDialog) update(d *Dashboard, msg tea.Msg) tea.Cmd {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return nil
	}
	if k.String() == "esc" || k.String() == "q" {
		d.dlg = nil
		return nil
	}
	if !pickKey(k.String(), &c.cursor, 2) {
		return nil
	}
	d.dlg = nil
	if c.labels()[c.cursor] == c.yes {
		return c.onYes(d)
	}
	return nil
}

func (c *confirmDialog) view(d *Dashboard, w int) []string {
	t := d.theme
	if !c.destructive {
		out := dialogTop(t, w, c.title)
		for _, l := range wrap(c.body, w-4) {
			out = append(out, "  "+l)
		}
		if c.body != "" {
			out = append(out, "")
		}
		return append(out, options(t, c.labels(), c.cursor)...)
	}
	var b strings.Builder
	b.WriteString(t.bold().Render(c.title) + "\n\n")
	for _, l := range wrap(c.body, w-8) {
		b.WriteString(l + "\n")
	}
	b.WriteString("\n" + strings.Join(options(t, c.labels(), c.cursor), "\n"))
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(t.Error).Padding(0, 1).Width(w - 2)
	out := []string{""}
	for _, l := range strings.Split(box.Render(b.String()), "\n") {
		out = append(out, " "+l)
	}
	return out
}

func (c *confirmDialog) keys() string {
	return "↑/↓ to select" + symDot + "enter to confirm" + symDot + "esc to cancel"
}

// inputDialog asks for one line of text. onOK returns an error to show,
// keeping the dialog open, or nil and what to run.
type inputDialog struct {
	title, hint string
	in          input
	err         string
	onOK        func(d *Dashboard, value string) (tea.Cmd, error)
}

func newInputDialog(title, hint, value, placeholder string, onOK func(d *Dashboard, value string) (tea.Cmd, error)) *inputDialog {
	return &inputDialog{title: title, hint: hint, in: newInput(value, placeholder), onOK: onOK}
}

func (i *inputDialog) update(d *Dashboard, msg tea.Msg) tea.Cmd {
	if k, ok := msg.(tea.KeyPressMsg); ok {
		switch k.String() {
		case "esc":
			d.dlg = nil
			return nil
		case "enter":
			d.dlg = nil
			cmd, err := i.onOK(d, strings.TrimSpace(i.in.Value()))
			if err != nil {
				i.err = err.Error()
				d.dlg = i
				return nil
			}
			return cmd
		}
	}
	if i.in.update(msg) {
		i.err = ""
	}
	return nil
}

func (i *inputDialog) view(d *Dashboard, w int) []string {
	t := d.theme
	out := dialogTop(t, w, i.title)
	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(t.Subtle).Padding(0, 1).Width(w - 4)
	for _, l := range strings.Split(box.Render(i.in.view(t)), "\n") {
		out = append(out, "  "+l)
	}
	if i.hint != "" {
		out = append(out, "  "+t.dim().Render(i.hint))
	}
	if i.err != "" {
		for _, l := range wrap(symFail+" "+i.err, w-4) {
			out = append(out, "  "+t.error().Render(l))
		}
	}
	return out
}

func (i *inputDialog) keys() string { return "enter to confirm" + symDot + "esc to cancel" }

// choiceDialog picks one of a few options.
type choiceDialog struct {
	title  string
	body   string
	labels []string
	cursor int
	onPick func(d *Dashboard, i int) tea.Cmd
}

func (c *choiceDialog) update(d *Dashboard, msg tea.Msg) tea.Cmd {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return nil
	}
	if k.String() == "esc" || k.String() == "q" {
		d.dlg = nil
		return nil
	}
	if pickKey(k.String(), &c.cursor, len(c.labels)) {
		d.dlg = nil
		return c.onPick(d, c.cursor)
	}
	return nil
}

func (c *choiceDialog) view(d *Dashboard, w int) []string {
	out := dialogTop(d.theme, w, c.title)
	for _, l := range wrap(c.body, w-4) {
		if c.body != "" {
			out = append(out, "  "+l)
		}
	}
	if c.body != "" {
		out = append(out, "")
	}
	return append(out, options(d.theme, c.labels, c.cursor)...)
}

func (c *choiceDialog) keys() string {
	return "↑/↓ to select" + symDot + "enter to choose" + symDot + "esc to cancel"
}

// planDialog shows a group change and asks to apply it; the plan can be
// longer than the screen, so it scrolls.
type planDialog struct {
	title  string
	lines  []string
	scroll int
	cursor int
	onYes  func(d *Dashboard) tea.Cmd
}

// planRoom is how many plan lines fit with the rest of the screen.
func (p *planDialog) room(d *Dashboard) int { return max(4, d.height-15) }

func (p *planDialog) update(d *Dashboard, msg tea.Msg) tea.Cmd {
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return nil
	}
	switch k.String() {
	case "esc", "q":
		d.dlg = nil
		return nil
	case "pgdown", "space", "ctrl+d":
		p.scroll = min(max(0, len(p.lines)-p.room(d)), p.scroll+p.room(d))
		return nil
	case "pgup", "ctrl+u":
		p.scroll = max(0, p.scroll-p.room(d))
		return nil
	}
	if !pickKey(k.String(), &p.cursor, 2) {
		return nil
	}
	d.dlg = nil
	if p.cursor == 0 {
		return p.onYes(d)
	}
	return nil
}

func (p *planDialog) view(d *Dashboard, w int) []string {
	t := d.theme
	out := dialogTop(t, w, p.title)
	room := p.room(d)
	shown := p.lines[p.scroll:min(len(p.lines), p.scroll+room)]
	for _, l := range shown {
		out = append(out, "  "+ansi.Truncate(l, w-3, symMore))
	}
	if rest := len(p.lines) - p.scroll - len(shown); rest > 0 {
		out = append(out, "  "+t.dim().Render(fmt.Sprintf("… %d more lines (space for more)", rest)))
	}
	out = append(out, "", "  Do you want to apply this plan?")
	return append(out, options(t, []string{"Yes", "No"}, p.cursor)...)
}

func (p *planDialog) keys() string {
	return "↑/↓ to select" + symDot + "space for more" + symDot + "enter to confirm" + symDot + "esc to cancel"
}

// infoDialog shows text until a key is pressed.
type infoDialog struct {
	title string
	lines []string
}

func (i *infoDialog) update(d *Dashboard, msg tea.Msg) tea.Cmd {
	if _, ok := msg.(tea.KeyPressMsg); ok {
		d.dlg = nil
	}
	return nil
}

func (i *infoDialog) view(d *Dashboard, w int) []string {
	out := dialogTop(d.theme, w, i.title)
	for _, l := range i.lines {
		out = append(out, "  "+l)
	}
	return out
}

func (i *infoDialog) keys() string { return "any key to close" }
