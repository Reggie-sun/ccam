package tui

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/jhyunwoo/ccam/internal/config"
)

// groupsTab is the editor of the Groups tab: what a group shares and who
// is in it, changed on a copy until the plan is applied.
type groupsTab struct {
	editing bool
	isNew   bool
	name    string
	items   checklist
	members checklist
	pane    int // 0: items, 1: accounts
}

// sync closes the editor when its group is gone.
func (g *groupsTab) sync(cfg *config.Config) {
	if g.editing && !g.isNew {
		if _, ok := cfg.FindGroup(g.name); !ok {
			g.editing = false
		}
	}
}

func (d *Dashboard) edit(name string, isNew bool) {
	var items, members []string
	if gr, ok := d.cfg.FindGroup(name); ok {
		items, members = gr.Items, gr.Members
	}
	var names []string
	for _, p := range d.cfg.Profiles {
		names = append(names, p.Name)
	}
	d.groups = groupsTab{editing: true, isNew: isNew, name: name,
		items:   newChecklist(config.Items, func(i string) string { return itemLabels[i] }, items),
		members: newChecklist(names, func(n string) string { return n }, members)}
}

// itemBlocked says why item cannot be shared by the group being edited:
// one of its accounts gets it from another group.
func (d *Dashboard) itemBlocked(item string) string {
	g := &d.groups
	for _, m := range g.members.items() {
		if from := d.cfg.GroupFor(m, item); from != nil && from.Name != g.name {
			return m + ": from " + from.Name
		}
	}
	return ""
}

// memberBlocked says why account cannot join the group being edited.
func (d *Dashboard) memberBlocked(account string) string {
	g := &d.groups
	for _, item := range g.items.items() {
		if from := d.cfg.GroupFor(account, item); from != nil && from.Name != g.name {
			return item + " from " + from.Name
		}
	}
	return ""
}

func (d *Dashboard) groupsKey(k tea.KeyPressMsg) tea.Cmd {
	if d.groups.editing {
		return d.editorKey(k)
	}
	key := k.String()
	var gr *config.Group
	if len(d.cfg.Groups) > 0 {
		gr = &d.cfg.Groups[d.cursor[tabGroups]]
	}
	switch key {
	case "n":
		d.dlg = newInputDialog("New group", "a-z, 0-9, - and _; then choose what it shares and who is in it", "", "team", func(d *Dashboard, v string) (tea.Cmd, error) {
			if err := config.ValidateGroupName(v); err != nil {
				return nil, err
			}
			if _, ok := d.cfg.FindGroup(v); ok {
				return nil, fmt.Errorf("a group named %q already exists", v)
			}
			d.edit(v, true)
			return nil, nil
		})
	case "enter":
		if gr != nil {
			d.edit(gr.Name, false)
		}
	case "e":
		if gr == nil {
			return nil
		}
		name := gr.Name
		d.dlg = newInputDialog("Rename group "+name, "its folder keeps its name ("+gr.Folder()+"), so nothing moves", name, "", func(d *Dashboard, v string) (tea.Cmd, error) {
			if v == name || v == "" {
				return nil, nil
			}
			if err := config.ValidateGroupName(v); err != nil {
				return nil, err
			}
			return d.planChange("Rename "+name, func(c *config.Config) error { return c.RenameGroup(name, v) }, "group "+name+" is now "+v), nil
		})
	case "x":
		if gr == nil {
			return nil
		}
		name := gr.Name
		ch, err := d.be.PlanGroups(func(c *config.Config) error { return c.RemoveGroup(name) })
		if err != nil {
			d.notice = notice{text: err.Error(), err: true}
			return nil
		}
		body := "No account is in it."
		if names := ch.Names(); len(names) > 0 {
			body = joinAnd(names) + " leave it and keep a copy of " + cmp.Or(strings.Join(gr.Items, ", "), "its items") + "."
		}
		if len(d.cfg.FolderUsers(gr.Folder())) == 1 {
			body += " Its folder groups/" + gr.Folder() + " goes to ccam's trash. Nothing is deleted."
		}
		d.dlg = newDestructive("Remove group "+name+"?", body, "Yes, remove "+name, func(d *Dashboard) tea.Cmd {
			return d.run("Removing group "+name+"…", func() doneMsg {
				cfg, msg, err := d.be.RemoveGroup(name, ch)
				return doneMsg{cfg: cfg, notice: msg, err: err}
			})
		})
	default:
		return d.moveCursor(k)
	}
	return nil
}

func joinAnd(names []string) string {
	if len(names) <= 1 {
		return strings.Join(names, "")
	}
	return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

func (d *Dashboard) editorKey(k tea.KeyPressMsg) tea.Cmd {
	g := &d.groups
	list := &g.items
	if g.pane == 1 {
		list = &g.members
	}
	switch k.String() {
	case "tab", "shift+tab":
		g.pane = 1 - g.pane
		return nil
	case "space", " ", "x":
		if len(list.keys) == 0 {
			return nil
		}
		key := list.keys[list.cursor]
		if !list.selected[key] {
			why := d.itemBlocked(key)
			if g.pane == 1 {
				why = d.memberBlocked(key)
			}
			if why != "" {
				d.notice = notice{text: key + " (" + why + "): an item can come from one group only", err: true}
				return nil
			}
		}
		list.selected[key] = !list.selected[key]
		return nil
	case "enter":
		return d.applyEditor()
	}
	list.key(k.String())
	return nil
}

func (d *Dashboard) applyEditor() tea.Cmd {
	g := d.groups
	name, items, members := g.name, g.items.items(), g.members.items()
	if g.isNew {
		if len(items) == 0 {
			d.notice = notice{text: "choose at least one item for the group to share", err: true}
			return nil
		}
		d.groups.editing = false
		return d.planChange("New group "+name, func(c *config.Config) error { return c.AddGroup(name, items, members) },
			"group "+name+" shares "+strings.Join(items, ", "))
	}
	edit := func(c *config.Config) error {
		gr, err := c.MustFindGroup(name)
		if err != nil {
			return err
		}
		for _, m := range slices.Clone(gr.Members) {
			if c.IsAccount(m) && !slices.Contains(members, m) {
				if err := c.LeaveGroup(m, name); err != nil {
					return err
				}
			}
		}
		if err := c.SetGroupItems(name, items); err != nil {
			return err
		}
		for _, m := range members {
			if _, err := c.JoinGroup(m, name, false); err != nil {
				return err
			}
		}
		return nil
	}
	d.groups.editing = false
	return d.planChange("Change group "+name, edit, "group "+name+" saved")
}

func (d *Dashboard) groupsView(w int) ([]string, int) {
	t := d.theme
	if len(d.cfg.Groups) == 0 && !d.groups.editing {
		return []string{
			" No groups yet. A group shares items (settings, CLAUDE.md, MCP servers, extensions,",
			" history) between the accounts in it. Press " + t.accent().Render("n") + " to create one.",
		}, -1
	}
	nameW, itemsW := len("no group"), len("Shares")
	for _, g := range d.cfg.Groups {
		nameW = max(nameW, ansi.StringWidth(g.Name))
		itemsW = max(itemsW, ansi.StringWidth(strings.Join(g.Items, symDot)))
	}
	nameW = min(nameW, 24)
	out := []string{t.dim().Render("    " + padRight("Group", nameW) + "  " + padRight("Shares", itemsW) + "  Accounts")}
	for i, g := range d.cfg.Groups {
		cur := "  "
		if i == d.cursor[tabGroups] {
			cur = t.accent().Render(symCursor) + " "
			if d.groups.editing {
				cur = t.dim().Render(symCursor) + " "
			}
		}
		name := padRight(t.Group(ansi.Truncate(g.Name, nameW, symMore)), nameW)
		items := padRight(cmp.Or(strings.Join(g.Items, symDot), t.dim().Render("nothing")), itemsW)
		var accts []string
		for _, m := range g.Members {
			if d.cfg.IsAccount(m) {
				accts = append(accts, m)
			}
		}
		out = append(out, " "+cur+" "+name+"  "+items+"  "+cmp.Or(strings.Join(accts, ", "), t.dim().Render("-")))
	}
	var none []string
	for _, p := range d.cfg.Profiles {
		if len(d.cfg.GroupsOf(p.Name)) == 0 {
			none = append(none, p.Name)
		}
	}
	if len(none) > 0 {
		out = append(out, t.dim().Render("    "+padRight("no group", nameW)+"  "+strings.Repeat(" ", itemsW)+"  "+strings.Join(none, ", ")))
	}
	cur := d.cursor[tabGroups] + 1
	if !d.groups.editing {
		return out, cur
	}
	return append(out, d.editorView(w)...), len(out) + 3 + max(d.groups.items.cursor, d.groups.members.cursor)
}

func (d *Dashboard) editorView(w int) []string {
	t := d.theme
	g := &d.groups
	title := g.name
	if g.isNew {
		title = "new group " + g.name
	} else if gr, ok := d.cfg.FindGroup(g.name); ok {
		title += symDot + "folder groups/" + gr.Folder()
	}
	out := []string{"", " " + t.subtle().Render(strings.Repeat("─", w-2)), "  " + t.bold().Render(title)}
	// Items need the room: their labels are longer than account names.
	colW := (w - 4) * 3 / 5
	left := []string{t.dim().Render("Shares")}
	for i, item := range g.items.keys {
		left = append(left, d.editorRow(g.pane == 0 && i == g.items.cursor, g.items.selected[item], editorLabels[item], d.itemBlocked(item), colW))
	}
	left = append(left, "    "+t.dim().Render(".credentials.json  never shared"))
	right := []string{t.dim().Render("Accounts")}
	for i, name := range g.members.keys {
		label := name
		if d.info[name].running {
			label += " " + t.accent().Render(symRunning+" running")
		}
		right = append(right, d.editorRow(g.pane == 1 && i == g.members.cursor, g.members.selected[name], label, d.memberBlocked(name), colW))
	}
	for i := 0; i < max(len(left), len(right)); i++ {
		l, r := "", ""
		if i < len(left) {
			l = left[i]
		}
		if i < len(right) {
			r = right[i]
		}
		out = append(out, "  "+padRight(ansi.Truncate(l, colW, symMore), colW)+"  "+ansi.Truncate(r, w-colW-5, symMore))
	}
	return out
}

var editorLabels = map[string]string{
	"mcp":        "MCP servers",
	"settings":   "settings.json",
	"memory":     "CLAUDE.md",
	"extensions": "skills, agents, plugins…",
	"history":    "conversation history",
}

func (d *Dashboard) editorRow(focus, on bool, label, blocked string, w int) string {
	t := d.theme
	box := "[ ]"
	if on {
		box = "[" + t.success().Render(symOK) + "]"
	}
	cur := "  "
	if focus {
		cur = t.accent().Render(symCursor) + " "
		label = t.accent().Render(label)
	}
	if blocked != "" && !on {
		return cur + t.inactive().Render("[ ] "+ansi.Strip(label)) + " " + t.dim().Render(blocked)
	}
	return cur + box + " " + label
}
