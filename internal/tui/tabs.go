package tui

import (
	"cmp"
	"slices"
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

// healthBox is the findings of one account, or of ccam as a whole when
// profile is "".
type healthBox struct {
	profile  string
	findings []doctor.Finding
}

// healthBoxes groups findings into a General box (when it has any) and one
// box per account, in the order doctor checked them.
func healthBoxes(fs []doctor.Finding) []healthBox {
	var boxes []healthBox
	at := map[string]int{}
	for _, f := range fs {
		i, ok := at[f.Profile]
		if !ok {
			i = len(boxes)
			at[f.Profile] = i
			boxes = append(boxes, healthBox{profile: f.Profile})
		}
		boxes[i].findings = append(boxes[i].findings, f)
	}
	slices.SortStableFunc(boxes, func(a, b healthBox) int {
		return cmp.Compare(min(len(a.profile), 1), min(len(b.profile), 1)) // General first
	})
	return boxes
}

// problems counts the findings that still need attention.
func (b healthBox) problems() int {
	n := 0
	for _, f := range b.findings {
		if f.Level >= doctor.Warn && !f.Fixed {
			n++
		}
	}
	return n
}

// worst is the level the box's title shows.
func (b healthBox) worst() doctor.Level {
	l := doctor.OK
	for _, f := range b.findings {
		if f.Level >= doctor.Warn && !f.Fixed {
			l = max(l, f.Level)
		}
	}
	return l
}

func (t Theme) levelIcon(l doctor.Level) string {
	return map[doctor.Level]string{
		doctor.OK: t.success().Render(symOK), doctor.Info: t.dim().Render(symInfo),
		doctor.Warn: t.warning().Render(symWarn), doctor.Error: t.error().Render(symFail),
	}[l]
}

// healthAreaW fits the longest area name, "Share links" and "Folder rules".
const healthAreaW = 12

// healthView draws a summary line, then one box per account with a row per
// checked area, so a problem shows where it is at a glance.
func (d *Dashboard) healthView(w int) ([]string, int) {
	t := d.theme
	switch {
	case d.hErr != nil:
		return []string{" " + t.error().Render(symFail+" "+d.hErr.Error())}, -1
	case d.health == nil:
		return []string{" " + t.accent().Render(spinner[d.frame%len(spinner)]) + " Checking…"}, -1
	}
	out, cur, fails := healthLines(t, d.health, w, min(d.cursor[tabHealth], len(healthBoxes(d.health))-1))
	out = append(out, "")
	if fails == 0 {
		out = append(out, " "+t.success().Render(symOK)+" everything looks fine")
	} else {
		out = append(out, " "+plural(fails, "problem")+"; f repairs the selected box, F all of them (ccam doctor --fix).")
	}
	return out, cur
}

// HealthReport is `ccam doctor`'s output: the Health tab's summary line and
// boxes, width columns wide (0 fits each box to its longest line and wraps
// nothing, for output that is not a terminal), in the terminal's own colors.
func HealthReport(fs []doctor.Finding, width int, color bool) []string {
	out, _, _ := healthLines(ThemeFor("light-ansi", false), fs, width, -1)
	if !color {
		for i, l := range out {
			out[i] = ansi.Strip(l)
		}
	}
	return out
}

// healthLines is the summary line and one box per account, w columns wide
// (0: see HealthReport). Box sel (-1: none) gets the accent border; cur is
// its first line. fails counts the findings that still need attention.
func healthLines(t Theme, fs []doctor.Finding, w, sel int) (out []string, cur, fails int) {
	boxes := healthBoxes(fs)
	var summary []string
	for _, b := range boxes {
		fails += b.problems()
		summary = append(summary, t.levelIcon(b.worst())+" "+cmp.Or(b.profile, "General"))
	}
	out = []string{" " + strings.Join(summary, "   "), ""}
	boxW := w - 2
	if w <= 0 {
		// Every box as wide as the widest line, so they line up.
		boxW = 0
		for _, b := range boxes {
			_, rows := healthRows(t, b, 0)
			for _, r := range rows {
				boxW = max(boxW, ansi.StringWidth(r)+4)
			}
		}
	}
	cur = -1
	for i, b := range boxes {
		if i == sel {
			cur = len(out)
		}
		out = append(out, healthBoxLines(t, b, boxW, i == sel)...)
	}
	return out, cur, fails
}

// healthBoxLines draws one box boxW wide (at least as wide as its title):
// the title names the account and its problem count, and each row is an
// area with its findings.
func healthBoxLines(t Theme, b healthBox, boxW int, selected bool) []string {
	border := t.subtle()
	if selected {
		border = t.accent()
	}
	title, rows := healthRows(t, b, boxW)
	boxW = max(boxW, ansi.StringWidth(title)+6)
	inner := boxW - 4
	top := border.Render("╭─ ") + title + " "
	top += border.Render(strings.Repeat("─", max(0, boxW-ansi.StringWidth(top)-1)) + "╮")
	out := []string{" " + top}
	for _, r := range rows {
		out = append(out, " "+border.Render("│")+" "+padRight(ansi.Truncate(r, inner, symMore), inner)+" "+border.Render("│"))
	}
	out = append(out, " "+border.Render("╰"+strings.Repeat("─", max(0, boxW-2))+"╯"))
	return out
}

// healthRows is a box's title and rows, wrapped to fit a box boxW wide
// (0: not wrapped). An account's folder has a finding only when something
// is wrong with it, so its row says OK otherwise.
func healthRows(t Theme, b healthBox, boxW int) (string, []string) {
	title := t.levelIcon(b.worst()) + " " + t.bold().Render(cmp.Or(b.profile, "General"))
	if n := b.problems(); n > 0 {
		title += t.dim().Render(symDot + plural(n, "problem"))
	}
	fs := b.findings
	if b.profile != "" {
		// The same rows in the same order in every account's box.
		if !slices.ContainsFunc(fs, func(f doctor.Finding) bool { return f.Area == "Folder" && f.Level != doctor.Info }) {
			fs = append([]doctor.Finding{{Level: doctor.OK, Area: "Folder", Msg: "OK"}}, fs...)
		}
		rank := func(f doctor.Finding) int {
			if i := slices.Index([]string{"Folder", "Login", "Share links"}, f.Area); i >= 0 {
				return i
			}
			return 3
		}
		fs = slices.Clone(fs)
		slices.SortStableFunc(fs, func(a, b doctor.Finding) int { return cmp.Compare(rank(a), rank(b)) })
	}
	var rows []string
	prev := ""
	for _, f := range fs {
		label := ""
		if f.Area != prev {
			label = f.Area
		}
		prev = f.Area
		text := f.Msg
		if f.Fixed {
			text += " " + t.success().Render("→ fixed")
		}
		lines := strings.Split(text, "\n")
		if boxW > 0 {
			lines = wrap(text, boxW-4-2-healthAreaW-1)
		}
		for i, l := range lines {
			if i == 0 {
				rows = append(rows, t.levelIcon(f.Level)+" "+padRight(t.dim().Render(label), healthAreaW)+" "+l)
			} else {
				rows = append(rows, strings.Repeat(" ", 2+healthAreaW+1)+l)
			}
		}
	}
	return title, rows
}

// healthKey: f repairs the selected box (one account, or General), F
// everything, like `ccam doctor --fix`.
func (d *Dashboard) healthKey(k tea.KeyPressMsg) tea.Cmd {
	only, title, scope := "", "Repair everything that can be repaired?", ""
	switch k.String() {
	case "F":
	case "f":
		boxes := healthBoxes(d.health)
		if len(boxes) == 0 {
			return nil
		}
		b := boxes[min(d.cursor[tabHealth], len(boxes)-1)]
		only, scope = cmp.Or(b.profile, doctor.General), " in "+cmp.Or(b.profile, "General")
		if b.problems() == 0 {
			d.notice = notice{text: "nothing to repair" + scope}
			return nil
		}
		title = "Repair what can be repaired" + scope + "?"
	default:
		return d.moveCursor(k) // ↑/↓ go from box to box
	}
	d.dlg = newConfirm(title, "Like `ccam doctor --fix`: missing links are made, and a real file where a link belongs is backed up first. Nothing is deleted.", "Yes, repair", func(d *Dashboard) tea.Cmd {
		d.busy = "Repairing…"
		return tea.Batch(func() tea.Msg {
			f, err := d.be.Doctor(true, only)
			return healthMsg{findings: f, err: err, fixed: true, scope: scope}
		}, d.tick())
	})
	return nil
}
