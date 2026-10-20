package tui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/x/ansi"

	"github.com/jhyunwoo/ccam/internal/config"
)

// checklist is a list with checkboxes: share items in the migrate picker,
// share items or groups in the add wizard, groups in the group editor.
type checklist struct {
	keys     []string
	label    func(key string) string
	selected map[string]bool
	cursor   int
}

func newChecklist(keys []string, label func(string) string, checked []string) checklist {
	c := checklist{keys: keys, label: label, selected: map[string]bool{}}
	for _, k := range checked {
		c.selected[k] = true
	}
	return c
}

var itemLabels = map[string]string{
	"mcp":        "MCP servers",
	"settings":   "settings.json",
	"memory":     "global CLAUDE.md",
	"extensions": "skills, agents, commands, plugins",
	"history":    "conversation history",
}

// itemChecklist lists the share items.
func itemChecklist(checked []string) checklist {
	return newChecklist(config.Items, func(item string) string { return itemLabels[item] }, checked)
}

// key moves the cursor or toggles the entry under it; false for other keys.
func (c *checklist) key(k string) bool {
	switch k {
	case "up", "k":
		if c.cursor > 0 {
			c.cursor--
		}
	case "down", "j":
		if c.cursor < len(c.keys)-1 {
			c.cursor++
		}
	case "space", " ", "x":
		if len(c.keys) > 0 {
			k := c.keys[c.cursor]
			c.selected[k] = !c.selected[k]
		}
	default:
		return false
	}
	return true
}

// items returns the checked keys in list order.
func (c checklist) items() []string {
	var out []string
	for _, k := range c.keys {
		if c.selected[k] {
			out = append(out, k)
		}
	}
	return out
}

// styled writes one line per entry; note, if not nil, adds dim text after it.
func (c checklist) styled(b *strings.Builder, t Theme, note func(key string) string) {
	w := 0
	for _, k := range c.keys {
		w = max(w, ansi.StringWidth(c.label(k)))
	}
	for i, k := range c.keys {
		cur, label := "  ", c.label(k)
		if i == c.cursor {
			cur, label = t.accent().Render(symCursor)+" ", t.accent().Render(label)
		}
		box := "[ ]"
		if c.selected[k] {
			box = "[" + t.success().Render(symOK) + "]"
		}
		line := label
		if note != nil {
			if n := note(k); n != "" {
				line = padRight(label, max(w, 24)) + "  " + t.dim().Render(n)
			}
		}
		fmt.Fprintf(b, " %s%s %s\n", cur, box, line)
	}
}

func padRight(s string, w int) string {
	if gap := w - ansi.StringWidth(s); gap > 0 {
		return s + strings.Repeat(" ", gap)
	}
	return s
}

// GroupInfo describes a group for the group lists.
type GroupInfo struct {
	Name    string
	Items   []string
	Members []string
}

// GroupInfos lists cfg's groups.
func GroupInfos(cfg *config.Config) []GroupInfo {
	out := make([]GroupInfo, 0, len(cfg.Groups))
	for _, g := range cfg.Groups {
		out = append(out, GroupInfo{Name: g.Name, Items: g.Items, Members: g.Members})
	}
	return out
}

// overlap explains why an account cannot be in all of selected, or "".
func overlap(groups []GroupInfo, selected []string) string {
	from := map[string]string{}
	for _, g := range groups {
		if !slices.Contains(selected, g.Name) {
			continue
		}
		for _, item := range g.Items {
			if prev, ok := from[item]; ok {
				return fmt.Sprintf("%s and %s both share %s; an item can come from one group only", prev, g.Name, item)
			}
			from[item] = g.Name
		}
	}
	return ""
}

// groupChecklist lists groups by name and what they share.
func groupChecklist(groups []GroupInfo, checked []string) checklist {
	keys := make([]string, len(groups))
	items := map[string]string{}
	for i, g := range groups {
		keys[i] = g.Name
		items[g.Name] = strings.Join(g.Items, " ")
	}
	return newChecklist(keys, func(name string) string {
		if items[name] == "" {
			return name
		}
		return name + " (" + items[name] + ")"
	}, checked)
}
