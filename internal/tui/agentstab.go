package tui

import (
	"cmp"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

const (
	agentIDW        = 8 // ids are 8 hex digits
	agentAccountMax = 16
	agentStateMax   = 24
	agentStartedW   = len("Started")
	agentFolderMin  = 10
)

// agentCols are the Agents tab's column widths; 0 leaves Name or Folder out.
type agentCols struct {
	account, name, state, folder int
	started                      bool
}

// agentColumns fits the columns of `ccam agents` into w: Account, ID and
// State always; Started while Name and Folder still get their least; Name
// at most half of the rest, Folder the remainder, then Name what Folder
// leaves, neither wider than its widest value. A narrower screen drops
// Started, then Name, then Folder.
func agentColumns(list []Agent, w int) agentCols {
	c := agentCols{account: len("Account"), state: len("State")}
	name, folder := len("Name"), len("Folder")
	for _, a := range list {
		c.account = max(c.account, ansi.StringWidth(a.Account))
		c.state = max(c.state, ansi.StringWidth(a.State))
		name = max(name, ansi.StringWidth(a.Name))
		folder = max(folder, ansi.StringWidth(a.Folder))
	}
	c.account, c.state = min(c.account, agentAccountMax), min(c.state, agentStateMax)
	base := 4 + c.account + 2 + agentIDW + 2 + c.state
	const least = 2 + len("Name") + 2 + agentFolderMin
	c.started = base+2+agentStartedW+least <= w
	rest := w - base
	if c.started {
		rest -= 2 + agentStartedW
	}
	switch {
	case rest >= least:
		c.name = min(name, max(len("Name"), (rest-4)/2), rest-4-agentFolderMin)
		c.folder = min(folder, rest-4-c.name)
		c.name = min(name, rest-4-c.folder)
	case rest >= 2+agentFolderMin:
		c.folder = min(folder, rest-2)
	}
	return c
}

// agentsView is the sessions of `ccam agents`, or its empty sentence, then
// a line for each account that could not be asked.
func (d *Dashboard) agentsView(w int) ([]string, int) {
	t := d.theme
	if d.agents == nil {
		return []string{" " + t.accent().Render(spinner[d.frame%len(spinner)]) + " Asking every account's Claude Code for its sessions…"}, -1
	}
	var out []string
	list := d.agents.Sessions
	if len(list) > 0 {
		c := agentColumns(list, w)
		row := func(account, id, name, state, folder, started string) string {
			line := padRight(account, c.account) + "  " + padRight(id, agentIDW)
			if c.name > 0 {
				line += "  " + padRight(name, c.name)
			}
			line += "  " + padRight(state, c.state)
			if c.folder > 0 {
				line += "  " + padRight(folder, c.folder)
			}
			if c.started {
				line += "  " + started
			}
			return line
		}
		out = append(out, t.dim().Render("    "+row("Account", "ID", "Name", "State", "Folder", "Started")))
		none, now := t.dim().Render("-"), d.now()
		for i, a := range list {
			cur := "  "
			if i == d.cursor[tabAgents] {
				cur = t.accent().Render(symCursor) + " "
			}
			out = append(out, " "+cur+" "+row(ansi.Truncate(a.Account, c.account, symMore),
				cmp.Or(ansi.Truncate(a.ID, agentIDW, symMore), none), cmp.Or(ansi.Truncate(a.Name, c.name, symMore), none),
				ansi.Truncate(a.State, c.state, symMore), cutLeft(a.Folder, c.folder), t.dim().Render(shortAgo(a.StartedAt, now))))
		}
	} else if len(d.agents.Failed) == 0 {
		for _, l := range wrap("No sessions. Start one in the background with `ccam <account> --bg \"…\"`.", w-2) {
			out = append(out, " "+l)
		}
	}
	if len(out) > 0 && len(d.agents.Failed) > 0 {
		out = append(out, "")
	}
	for _, f := range d.agents.Failed {
		for i, l := range wrap(f, w-5) {
			mark := "   "
			if i == 0 {
				mark = " " + t.error().Render(symFail) + " "
			}
			out = append(out, mark+l)
		}
	}
	if len(list) == 0 {
		return out, -1
	}
	return out, d.cursor[tabAgents] + 1
}

// agentsKey: enter attaches to the selected background session (the
// dashboard ends, as a launch does), x stops it after asking.
func (d *Dashboard) agentsKey(k tea.KeyPressMsg) tea.Cmd {
	var a *Agent
	if d.agents != nil && len(d.agents.Sessions) > 0 {
		a = &d.agents.Sessions[d.cursor[tabAgents]]
	}
	switch k.String() {
	case "enter":
		if a == nil {
			return nil
		}
		if a.ID == "" {
			d.notice = notice{text: "this session has no id; only a background session can be attached", err: true}
			return nil
		}
		return d.quit(Result{Action: Launch, Profile: a.Account, Args: []string{"attach", a.ID}})
	case "x":
		if a == nil {
			return nil
		}
		if a.ID == "" {
			d.notice = notice{text: "this session has no id; only a background session can be stopped", err: true}
			return nil
		}
		account, id := a.Account, a.ID
		d.dlg = newDestructive("Stop background session "+id+"?",
			"It runs as "+account+". Like `ccam "+account+" stop "+id+"`, Claude Code stops it and keeps its conversation, which `ccam attach "+id+"` opens again.",
			"Yes, stop "+id, func(d *Dashboard) tea.Cmd {
				d.loaded[tabAgents] = false // read the list again once it is done (doneMsg opens the tab)
				return d.run("Stopping "+id+"…", func() doneMsg {
					return doneMsg{notice: id + " (" + account + ") stopped; `ccam attach " + id + "` opens it again", err: d.be.StopSession(account, id)}
				})
			})
	default:
		return d.moveCursor(k)
	}
	return nil
}
