package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/doctor"
)

func (d *Dashboard) foldersView(w int) ([]string, int) {
	t := d.theme
	here, _ := d.cfg.FolderFor(d.cwd)
	var out []string
	if len(d.cfg.Folders) == 0 {
		out = append(out, " No folder rules. In a folder with a rule, `ccam` without a name starts that",
			" folder's account (ccam folder).")
	} else {
		pathW := 0
		for _, f := range d.cfg.Folders {
			pathW = max(pathW, ansi.StringWidth(f.Path))
		}
		pathW = max(min(pathW, w-30), 8)
		out = append(out, t.dim().Render("    "+padRight("Folder", pathW)+"  Account"))
		for i, f := range d.cfg.Folders {
			cur := "  "
			if i == d.cursor[tabFolders] {
				cur = t.accent().Render(symCursor) + " "
			}
			path := f.Path
			if cut := ansi.StringWidth(path) - pathW; cut > 0 {
				path = ansi.TruncateLeft(path, cut+ansi.StringWidth(symMore), symMore)
			}
			path = padRight(path, pathW)
			line := " " + cur + " " + path + "  " + f.Account
			if here != nil && here.Path == f.Path {
				line += "  " + t.accent().Render(symHere+" this folder")
			}
			if _, ok := d.cfg.Find(f.Account); !ok {
				line += "  " + t.error().Render(symFail+" no such account")
			}
			out = append(out, line)
		}
	}
	out = append(out, "", " "+t.dim().Render("This folder: "+d.cwd))
	if here != nil {
		out = append(out, " "+t.dim().Render("uses "+here.Account+" (rule for "+here.Path+")"))
	}
	if len(d.cfg.Folders) == 0 {
		return out, -1
	}
	return out, d.cursor[tabFolders] + 1
}

func (d *Dashboard) foldersKey(k tea.KeyPressMsg) tea.Cmd {
	var f *config.Folder
	if len(d.cfg.Folders) > 0 {
		f = &d.cfg.Folders[d.cursor[tabFolders]]
	}
	switch k.String() {
	case "n":
		d.dlg = d.chooseAccount("Use which account in "+d.cwd+"?", d.cwd)
	case "enter":
		if f != nil {
			d.dlg = d.chooseAccount("Use which account in "+f.Path+"?", f.Path)
		}
	case "x":
		if f == nil {
			return nil
		}
		path := f.Path
		d.dlg = newConfirm("Remove the rule for "+path+"?", "`ccam` without a name then starts the default account there.", "Yes, remove it", func(d *Dashboard) tea.Cmd {
			cfg, err := d.be.UnsetFolder(path)
			d.after(cfg, err, "no rule for "+path+" now")
			return nil
		})
	default:
		return d.moveCursor(k)
	}
	return nil
}

func (d *Dashboard) chooseAccount(title, path string) dialog {
	if len(d.cfg.Profiles) == 0 {
		return &infoDialog{title: title, lines: []string{"No accounts yet."}}
	}
	labels := make([]string, len(d.cfg.Profiles))
	for i, p := range d.cfg.Profiles {
		labels[i] = p.Name + "  " + d.theme.dim().Render(emailOf(&p))
	}
	return &choiceDialog{title: title, labels: labels, cursor: d.cursor[tabAccounts], onPick: func(d *Dashboard, i int) tea.Cmd {
		name := d.cfg.Profiles[i].Name
		cfg, err := d.be.SetFolder(path, name)
		d.after(cfg, err, name+" is used in "+path)
		return nil
	}}
}

func (d *Dashboard) extView(w int) []string {
	t := d.theme
	switch {
	case d.extErr != nil:
		return []string{" " + t.error().Render(symFail+" "+d.extErr.Error())}
	case d.ext == nil:
		return []string{" " + t.accent().Render(spinner[d.frame%len(spinner)]) + " Reading the accounts' extensions…"}
	}
	var out []string
	for _, k := range d.ext.Kinds {
		if len(k.Rows) == 0 {
			continue
		}
		out = append(out, " "+t.bold().Render(k.Title))
		nameW := 0
		for _, r := range k.Rows {
			nameW = max(nameW, ansi.StringWidth(r.Name))
		}
		nameW = min(nameW, 28)
		for _, r := range k.Rows {
			var where []string
			for _, g := range r.Groups {
				where = append(where, t.Group(g))
			}
			if len(r.Own) > 0 {
				where = append(where, t.dim().Render("own: ")+strings.Join(r.Own, ", "))
			}
			mark := "  "
			if r.Problem {
				mark = t.warning().Render(symWarn) + " "
			}
			out = append(out, "  "+mark+padRight(ansi.Truncate(r.Name, nameW, symMore), nameW)+"  "+strings.Join(where, symDot))
		}
		out = append(out, "")
	}
	if len(out) == 0 {
		out = append(out, " No skills, commands, agents, plugins, MCP servers or hooks in the accounts.", "")
	}
	if len(d.ext.Project) > 0 {
		out = append(out, " "+t.bold().Render("This folder's project, for every account"))
		for _, l := range d.ext.Project {
			out = append(out, "   "+l)
		}
		out = append(out, "")
	}
	if len(d.ext.Problems) == 0 {
		out = append(out, " "+t.success().Render(symOK)+" no problems found")
	}
	for _, p := range d.ext.Problems {
		for i, l := range wrap(p, w-5) {
			mark := "   "
			if i == 0 {
				mark = " " + t.warning().Render(symWarn) + " "
			}
			out = append(out, mark+l)
		}
	}
	return out
}

func (d *Dashboard) healthView(w int) ([]string, int) {
	t := d.theme
	switch {
	case d.hErr != nil:
		return []string{" " + t.error().Render(symFail+" "+d.hErr.Error())}, -1
	case d.health == nil:
		return []string{" " + t.accent().Render(spinner[d.frame%len(spinner)]) + " Checking…"}, -1
	}
	var out []string
	fails := 0
	for _, f := range d.health {
		icon := map[doctor.Level]string{
			doctor.OK: t.success().Render(symOK), doctor.Info: t.dim().Render(symInfo),
			doctor.Warn: t.warning().Render(symWarn), doctor.Error: t.error().Render(symFail),
		}[f.Level]
		if f.Level >= doctor.Warn && !f.Fixed {
			fails++
		}
		text := f.Msg
		if f.Profile != "" {
			text = f.Profile + ": " + text
		}
		if f.Fixed {
			text += " " + t.success().Render("→ fixed")
		}
		for i, l := range wrap(text, w-5) {
			if i == 0 {
				out = append(out, " "+icon+" "+l)
			} else {
				out = append(out, "   "+l)
			}
		}
	}
	out = append(out, "")
	if fails == 0 {
		out = append(out, " "+t.success().Render(symOK)+" everything looks fine")
	} else {
		out = append(out, " "+plural(fails, "problem")+"; f repairs what can be repaired safely (ccam doctor --fix).")
	}
	return out, -1
}

func (d *Dashboard) healthKey(k tea.KeyPressMsg) tea.Cmd {
	if k.String() != "f" {
		return nil
	}
	d.dlg = newConfirm("Repair what can be repaired?", "Like `ccam doctor --fix`: missing links are made, and a real file where a link belongs is backed up first. Nothing is deleted.", "Yes, repair", func(d *Dashboard) tea.Cmd {
		d.busy = "Repairing…"
		return tea.Batch(func() tea.Msg {
			f, err := d.be.Doctor(true)
			return healthMsg{findings: f, err: err, fixed: true}
		}, d.tick())
	})
	return nil
}
