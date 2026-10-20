package tui

import (
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// input is a one-line text field: ←/→, home/end, backspace, delete,
// ctrl+w (word), ctrl+u (line) and paste.
type input struct {
	text        []rune
	pos         int
	placeholder string
}

func newInput(value, placeholder string) input {
	r := []rune(value)
	return input{text: r, pos: len(r), placeholder: placeholder}
}

func (in input) Value() string { return string(in.text) }

func (in *input) insert(s string) {
	r := []rune(s)
	in.text = append(in.text[:in.pos], append(r, in.text[in.pos:]...)...)
	in.pos += len(r)
}

// update handles an editing key; false means the key was not for the field.
func (in *input) update(msg tea.Msg) bool {
	if p, ok := msg.(tea.PasteMsg); ok {
		in.insert(strings.TrimSpace(strings.ReplaceAll(p.Content, "\n", " ")))
		return true
	}
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return false
	}
	switch k.String() {
	case "left", "ctrl+b":
		in.pos = max(0, in.pos-1)
	case "right", "ctrl+f":
		in.pos = min(len(in.text), in.pos+1)
	case "home", "ctrl+a":
		in.pos = 0
	case "end", "ctrl+e":
		in.pos = len(in.text)
	case "backspace":
		if in.pos > 0 {
			in.text = append(in.text[:in.pos-1], in.text[in.pos:]...)
			in.pos--
		}
	case "delete", "ctrl+d":
		if in.pos < len(in.text) {
			in.text = append(in.text[:in.pos], in.text[in.pos+1:]...)
		}
	case "ctrl+u":
		in.text, in.pos = in.text[in.pos:], 0
	case "ctrl+w":
		start := in.pos
		for start > 0 && unicode.IsSpace(in.text[start-1]) {
			start--
		}
		for start > 0 && !unicode.IsSpace(in.text[start-1]) {
			start--
		}
		in.text = append(in.text[:start], in.text[in.pos:]...)
		in.pos = start
	default:
		if k.Text == "" {
			return false
		}
		in.insert(k.Text)
	}
	return true
}

// view draws the text with a block cursor, or the placeholder dimmed.
func (in input) view(t Theme) string {
	if len(in.text) == 0 {
		return "█" + t.dim().Render(in.placeholder)
	}
	before, after := string(in.text[:in.pos]), ""
	cur := " "
	if in.pos < len(in.text) {
		cur, after = string(in.text[in.pos]), string(in.text[in.pos+1:])
	}
	return before + lipgloss.NewStyle().Reverse(true).Render(cur) + after
}
