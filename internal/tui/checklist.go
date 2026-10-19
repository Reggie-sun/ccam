package tui

import (
	"fmt"
	"strings"

	"github.com/jhyunwoo/ccam/internal/config"
)

// checklist is the list of share items with checkboxes, used by the add
// wizard and the migrate picker.
type checklist struct {
	selected map[string]bool
	cursor   int
}

func newChecklist(checked []string) checklist {
	c := checklist{selected: map[string]bool{}}
	for _, item := range checked {
		c.selected[item] = true
	}
	return c
}

// key moves the cursor or toggles the item under it; false for other keys.
func (c *checklist) key(k string) bool {
	switch k {
	case "up", "k":
		if c.cursor > 0 {
			c.cursor--
		}
	case "down", "j":
		if c.cursor < len(config.Items)-1 {
			c.cursor++
		}
	case "space", " ", "x":
		item := config.Items[c.cursor]
		c.selected[item] = !c.selected[item]
	default:
		return false
	}
	return true
}

// items returns the checked items in display order.
func (c checklist) items() []string {
	var out []string
	for _, item := range config.Items {
		if c.selected[item] {
			out = append(out, item)
		}
	}
	return out
}

// view writes one line per item; note, if not nil, adds dim text after it.
func (c checklist) view(b *strings.Builder, note func(item string) string) {
	for i, item := range config.Items {
		marker := "  "
		if i == c.cursor {
			marker = cursorStyle.Render("› ")
		}
		box := "[ ]"
		if c.selected[item] {
			box = "[x]"
		}
		line := itemLabels[item]
		if note != nil {
			if n := note(item); n != "" {
				line = padRight(line, 24) + " " + dimStyle.Render(n)
			}
		}
		fmt.Fprintf(b, " %s%s %s\n", marker, box, line)
	}
}
