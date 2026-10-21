// Package tui is the `ccam` dashboard, the `ccam add` wizard and the
// `ccam migrate` picker.
package tui

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/doctor"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// Backend does the work the dashboard asks for. A method that changes
// something returns the config as it is afterwards.
type Backend interface {
	Usage(ctx context.Context, name string, fresh bool) usage.Report
	Running(name string) bool
	// Sessions is how many ccam sessions of the account run now.
	Sessions(name string) int
	// Takeover lists the accounts that may take over from from by the rules
	// of next --auto (the same organization, not the same Claude account).
	Takeover(from string) []string
	LastUsed(name string) time.Time
	SetDefault(name string) (*config.Config, error)
	SetFolder(dir, account string) (*config.Config, error)
	UnsetFolder(dir string) (*config.Config, error)
	Rename(oldName, newName string) (*config.Config, error)
	// SetShortcut gives account the shortcut command; "" removes it.
	SetShortcut(account, shortcut string) (*config.Config, error)
	SuggestShortcut(account string) string
	CheckShortcut(account, shortcut string) error
	Remove(name string) (*config.Config, string, error)
	// RemoveGroup applies ch (the group's accounts leave it), drops the
	// group and moves its folder to the trash.
	RemoveGroup(name string, ch *group.Change) (*config.Config, string, error)
	// PlanGroups builds what edit would move on disk, refusing a change
	// that cannot be applied now; ApplyGroups saves edit and moves.
	PlanGroups(edit func(*config.Config) error) (*group.Change, error)
	ApplyGroups(edit func(*config.Config) error, ch *group.Change) (*config.Config, error)
	Extensions() (ExtView, error)
	// Doctor runs the checks; with fix it repairs what is safe, only the
	// findings of account only (doctor.General: of no account) when set.
	Doctor(fix bool, only string) ([]doctor.Finding, error)
}

// ExtView is `ccam ext` prepared for the Extensions tab.
type ExtView struct {
	Kinds    []ExtKind
	Project  []string // lines about the current folder's own extensions
	Problems []string
}

type ExtKind struct {
	Title string
	Rows  []ExtRow
}

// ExtRow is one extension and where it is loaded: Groups are group names,
// Own the accounts that have it in their own folder.
type ExtRow struct {
	Name    string
	Groups  []string
	Own     []string
	Problem bool
}

// Action is what the user chose that needs the terminal; the caller does
// it after the program has exited.
type Action int

const (
	None Action = iota
	Launch
	Add
	Login
)

type Result struct {
	Action  Action
	Profile string
	Args    []string // for Launch: arguments for claude
	DryRun  bool     // for Launch: print the command only
	Add     AddAnswers
}

const (
	tabAccounts = iota
	tabGroups
	tabFolders
	tabExtensions
	tabHealth
)

var tabNames = []string{"Accounts", "Groups", "Folders", "Extensions", "Health"}

// tabAbout says in one line what each tab is for; it shows under the tabs.
var tabAbout = []string{
	"Your Claude accounts and their usage limits; enter launches the selected one.",
	"Groups share settings, memory, extensions, MCP or history between accounts.",
	"A folder rule picks the account that plain `ccam` starts in that folder.",
	"Skills, commands, agents, plugins, MCP and hooks, and which accounts see them.",
	"Checks Claude Code and each account's folder, login and links (ccam doctor).",
}

type (
	usageMsg struct {
		name   string
		report usage.Report
	}
	tickMsg struct{}
	doneMsg struct {
		cfg    *config.Config
		notice string
		err    error
	}
	extMsg struct {
		view ExtView
		err  error
	}
	healthMsg struct {
		findings []doctor.Finding
		err      error
		fixed    bool   // the result of f or F (doctor --fix)
		scope    string // what was repaired, for the notice
	}
)

// notice is the line above the keys: a result or an error, cleared by the
// next key.
type notice struct {
	text string
	err  bool
}

// accInfo is what the Accounts tab shows beside the config.
type accInfo struct {
	last    time.Time
	running bool
}

// Dashboard is the full-screen `ccam`. Build it with NewDashboard.
type Dashboard struct {
	be        Backend
	cfg       *config.Config
	cwd       string
	version   string
	themeName string
	theme     Theme
	width     int
	height    int
	tab       int
	cursor    [5]int

	usage     map[string]*usage.Report
	info      map[string]accInfo
	frame     int
	ticking   bool
	tickEvery time.Duration // 0 stops the spinner (tests)
	digits    string
	digitAt   time.Time

	groups groupsTab
	ext    *ExtView
	extErr error
	health []doctor.Finding
	hErr   error
	loaded [5]bool

	dlg    dialog
	help   bool
	busy   string // a change is being applied
	notice notice
	result Result
	now    func() time.Time
}

// NewDashboard shows cfg; cwd is the current folder (for ccam folder).
func NewDashboard(be Backend, cfg *config.Config, cwd, version string) *Dashboard {
	d := &Dashboard{be: be, cwd: cwd, version: version, themeName: cfg.Theme, width: 80, height: 24,
		usage: map[string]*usage.Report{}, now: time.Now, tickEvery: 120 * time.Millisecond}
	d.theme = ThemeFor(cfg.Theme, true)
	d.setConfig(cfg)
	def, _ := cfg.DefaultProfile()
	for i, p := range cfg.Profiles {
		if def != nil && p.Name == def.Name {
			d.cursor[tabAccounts] = i
		}
	}
	if f, ok := cfg.FolderFor(cwd); ok {
		for i, p := range cfg.Profiles {
			if p.Name == f.Account {
				d.cursor[tabAccounts] = i
			}
		}
	}
	return d
}

// Result is valid after the program has finished.
func (d *Dashboard) Result() Result { return d.result }

func (d *Dashboard) setConfig(cfg *config.Config) {
	d.cfg = cfg
	d.info = map[string]accInfo{}
	for _, p := range cfg.Profiles {
		d.info[p.Name] = accInfo{last: d.be.LastUsed(p.Name), running: d.be.Running(p.Name)}
	}
	for name := range d.usage {
		if _, ok := cfg.Find(name); !ok {
			delete(d.usage, name)
		}
	}
	for t := range d.cursor {
		d.cursor[t] = max(0, min(d.cursor[t], d.rowCount(t)-1))
	}
	d.groups.sync(cfg)
	// Tabs that read the disk load again when opened next.
	d.loaded[tabExtensions], d.loaded[tabHealth] = false, false
}

func (d *Dashboard) rowCount(tab int) int {
	switch tab {
	case tabAccounts:
		return len(d.cfg.Profiles)
	case tabGroups:
		return len(d.cfg.Groups)
	case tabFolders:
		return len(d.cfg.Folders)
	case tabHealth:
		return len(healthBoxes(d.health))
	}
	return 0
}

func (d *Dashboard) Init() tea.Cmd {
	return tea.Batch(tea.RequestBackgroundColor, d.fetchAll(false), d.tick())
}

func (d *Dashboard) fetchAll(fresh bool) tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(d.cfg.Profiles))
	for _, p := range d.cfg.Profiles {
		name := p.Name
		if fresh {
			delete(d.usage, name)
		}
		cmds = append(cmds, func() tea.Msg {
			return usageMsg{name: name, report: d.be.Usage(context.Background(), name, fresh)}
		})
	}
	return tea.Batch(cmds...)
}

// tick runs the spinner while something is loading.
func (d *Dashboard) tick() tea.Cmd {
	if d.ticking || d.tickEvery == 0 {
		return nil
	}
	d.ticking = true
	return tea.Tick(d.tickEvery, func(time.Time) tea.Msg { return tickMsg{} })
}

func (d *Dashboard) loading() bool {
	if d.busy != "" {
		return true
	}
	for _, p := range d.cfg.Profiles {
		if _, ok := d.usage[p.Name]; !ok {
			return true
		}
	}
	return d.tab == tabExtensions && d.ext == nil && d.extErr == nil ||
		d.tab == tabHealth && d.health == nil && d.hErr == nil
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width, d.height = msg.Width, msg.Height
		return d, nil
	case tea.BackgroundColorMsg:
		d.theme = ThemeFor(d.themeName, msg.IsDark())
		return d, nil
	case tickMsg:
		d.ticking = false
		d.frame++
		if d.loading() {
			return d, d.tick()
		}
		return d, nil
	case usageMsg:
		r := msg.report
		d.usage[msg.name] = &r
		return d, nil
	case doneMsg:
		d.busy = ""
		if msg.cfg != nil {
			d.setConfig(msg.cfg)
		}
		if msg.err != nil {
			d.notice = notice{text: msg.err.Error(), err: true}
		} else if msg.notice != "" {
			d.notice = notice{text: msg.notice}
		}
		return d, d.openTab()
	case extMsg:
		d.ext, d.extErr = &msg.view, msg.err
		return d, nil
	case healthMsg:
		d.health, d.hErr = msg.findings, msg.err
		if d.health == nil && d.hErr == nil {
			d.health = []doctor.Finding{}
		}
		if msg.fixed {
			d.busy, d.loaded[tabHealth] = "", true
			if msg.err == nil {
				d.notice = notice{text: "repaired what could be repaired" + msg.scope}
			}
		}
		return d, nil
	case tea.PasteMsg:
		if d.dlg != nil {
			return d, d.updateDialog(msg)
		}
		return d, nil
	case tea.KeyPressMsg:
		return d, d.key(msg)
	}
	return d, nil
}

// quit ends the program with r.
func (d *Dashboard) quit(r Result) tea.Cmd {
	d.result = r
	return tea.Quit
}

func (d *Dashboard) key(k tea.KeyPressMsg) tea.Cmd {
	if k.String() == "ctrl+c" {
		return tea.Quit
	}
	if d.busy != "" {
		return nil // a change is being applied; wait for it
	}
	if d.dlg != nil {
		return d.updateDialog(k)
	}
	d.notice = notice{}
	if d.help {
		d.help = false
		return nil
	}
	switch k.String() {
	case "q":
		return tea.Quit
	case "esc":
		if d.tab == tabGroups && d.groups.editing {
			d.groups.editing = false
			return nil
		}
		return tea.Quit
	case "?":
		d.help = true
		return nil
	case "tab", "right":
		if d.tab == tabGroups && d.groups.editing && k.String() == "tab" {
			break
		}
		d.tab = (d.tab + 1) % len(tabNames)
		return d.openTab()
	case "shift+tab", "left":
		if d.tab == tabGroups && d.groups.editing && k.String() == "shift+tab" {
			break
		}
		d.tab = (d.tab + len(tabNames) - 1) % len(tabNames)
		return d.openTab()
	case "r":
		d.loaded[tabExtensions], d.loaded[tabHealth] = false, false
		d.info = map[string]accInfo{}
		for _, p := range d.cfg.Profiles {
			d.info[p.Name] = accInfo{last: d.be.LastUsed(p.Name), running: d.be.Running(p.Name)}
		}
		return tea.Batch(d.fetchAll(true), d.openTab(), d.tick())
	}
	switch d.tab {
	case tabAccounts:
		return d.accountsKey(k)
	case tabGroups:
		return d.groupsKey(k)
	case tabFolders:
		return d.foldersKey(k)
	case tabHealth:
		return d.healthKey(k)
	}
	return d.moveCursor(k)
}

// moveCursor handles ↑/↓ for the current tab's rows.
func (d *Dashboard) moveCursor(k tea.KeyPressMsg) tea.Cmd {
	n := d.rowCount(d.tab)
	switch k.String() {
	case "up", "k":
		d.cursor[d.tab] = max(0, d.cursor[d.tab]-1)
	case "down", "j":
		d.cursor[d.tab] = min(n-1, d.cursor[d.tab]+1)
	case "home", "g":
		d.cursor[d.tab] = 0
	case "end", "G":
		d.cursor[d.tab] = max(0, n-1)
	}
	return nil
}

// openTab loads what the current tab needs from disk the first time.
func (d *Dashboard) openTab() tea.Cmd {
	if d.loaded[d.tab] {
		return nil
	}
	d.loaded[d.tab] = true
	switch d.tab {
	case tabExtensions:
		d.ext, d.extErr = nil, nil
		return tea.Batch(func() tea.Msg { v, err := d.be.Extensions(); return extMsg{v, err} }, d.tick())
	case tabHealth:
		d.health, d.hErr = nil, nil
		return tea.Batch(func() tea.Msg { f, err := d.be.Doctor(false, ""); return healthMsg{findings: f, err: err} }, d.tick())
	}
	return nil
}

// run applies a change in the background; the screen shows label meanwhile.
func (d *Dashboard) run(label string, f func() doneMsg) tea.Cmd {
	d.busy = label
	return tea.Batch(func() tea.Msg { return f() }, d.tick())
}

func (d *Dashboard) View() tea.View {
	v := tea.NewView(d.render())
	v.AltScreen = true
	v.WindowTitle = "ccam"
	return v
}

func (d *Dashboard) render() string {
	t := d.theme
	w := max(40, d.width)
	var top []string
	head := " " + t.bold().Render("ccam") + t.dim().Render(" "+d.version+symDot+plural(len(d.cfg.Profiles), "account")+symDot+plural(len(d.cfg.Groups), "group"))
	if d.tab == tabAccounts {
		legend := t.dim().Render(t.success().Render(symOK)+" default"+symDot+symHere+" this folder"+symDot+symRunning+" running") + " "
		if gap := w - ansi.StringWidth(head) - ansi.StringWidth(legend); gap > 1 {
			head += strings.Repeat(" ", gap) + legend
		}
	}
	top = append(top, head, "")
	var tabs []string
	for i, n := range tabNames {
		if i == d.tab {
			tabs = append(tabs, t.badge().Render(n))
		} else {
			tabs = append(tabs, t.dim().Render(" "+n+" "))
		}
	}
	top = append(top, " "+strings.Join(tabs, " "), " "+t.dim().Render(tabAbout[d.tab]), "")

	var dialog []string
	if d.dlg != nil {
		dialog = d.dlg.view(d, w)
	}
	tail := []string{""}
	switch {
	case d.busy != "":
		tail = append(tail, " "+t.accent().Render(spinner[d.frame%len(spinner)])+" "+d.busy)
	case d.notice.err:
		// An error says what to do next; keep it whole (up to 4 lines).
		lines := wrap(symFail+" "+d.notice.text, w-2)
		for _, l := range lines[:min(len(lines), 4)] {
			tail = append(tail, " "+t.error().Render(l))
		}
	case d.notice.text != "":
		tail = append(tail, " "+t.success().Render(symOK)+" "+d.notice.text)
	}
	tail = append(tail, " "+t.dim().Render(d.keysLine()))

	// A dialog taller than the screen leaves no room for the tab, and
	// loses its first lines rather than its choices.
	room := d.height - len(top) - len(dialog) - len(tail)
	if room < 0 {
		dialog = dialog[min(len(dialog), -room):]
		room = 0
	}
	if room < 3 && d.dlg != nil {
		room = 0 // a line or two of the tab says nothing
	}
	lines, cur := d.body(w)
	if d.help {
		lines, cur = d.helpView(w), -1
	}
	lines = window(t, lines, cur, room)
	all := append(append(append(top, lines...), dialog...), tail...)
	for i, l := range all {
		if ansi.StringWidth(l) > w {
			all[i] = ansi.Truncate(l, w-1, symMore)
		}
	}
	return strings.Join(all, "\n")
}

// body is the current tab's lines and the line the cursor is on (-1: none).
func (d *Dashboard) body(w int) ([]string, int) {
	switch d.tab {
	case tabAccounts:
		return d.accountsView(w)
	case tabGroups:
		return d.groupsView(w)
	case tabFolders:
		return d.foldersView(w)
	case tabExtensions:
		return d.extView(w), -1
	}
	return d.healthView(w)
}

// window shows at most room lines of lines, keeping cur in view, with
// "… N more" where lines are cut.
func window(t Theme, lines []string, cur, room int) []string {
	if len(lines) <= room {
		return lines
	}
	if room <= 0 {
		return nil
	}
	start := 0
	if cur >= 0 {
		start = max(0, min(cur-room/2, len(lines)-room))
	}
	end := min(len(lines), start+room)
	out := slices.Clone(lines[start:end])
	if start > 0 {
		out[0] = t.dim().Render(fmt.Sprintf("   ↑ %d more", start+1))
	}
	if end < len(lines) {
		out[len(out)-1] = t.dim().Render(fmt.Sprintf("   ↓ %d more", len(lines)-end+1))
	}
	return out
}

func plural(n int, word string) string {
	if n == 1 {
		return "1 " + word
	}
	return fmt.Sprintf("%d %ss", n, word)
}

// keysLine is the short key list for the current screen.
func (d *Dashboard) keysLine() string {
	if d.dlg != nil {
		return d.dlg.keys()
	}
	if d.help {
		return "any key to close"
	}
	var keys []string
	switch d.tab {
	case tabAccounts:
		keys = []string{"enter to launch", "n for most room", "a to add", "s for groups"}
	case tabGroups:
		if d.groups.editing {
			keys = []string{"space to toggle", "tab to switch lists", "enter for the plan", "esc to close"}
		} else {
			keys = []string{"enter to edit", "n for a new group", "e to rename", "x to remove"}
		}
	case tabFolders:
		keys = []string{"n to use an account here", "enter to change", "x to remove"}
	case tabExtensions:
		keys = []string{"r to reload"}
	case tabHealth:
		keys = []string{"f to fix the selected box", "F to fix all", "r to check again"}
	}
	if !(d.tab == tabGroups && d.groups.editing) {
		keys = append(keys, "? for shortcuts")
	}
	return strings.Join(keys, symDot)
}

func (d *Dashboard) helpView(w int) []string {
	t := d.theme
	rows := [][2]string{
		{"←/→ tab", "switch tabs"}, {"↑/↓ j/k", "move"}, {"r", "refresh"}, {"q esc", "quit"},
	}
	switch d.tab {
	case tabAccounts:
		rows = append(rows, [][2]string{
			{"enter", "launch the account"}, {"1-9", "go to that number (12: 1 2)"},
			{"n", "launch the most room"}, {"o", "launch with options"},
			{"u", "usage details"}, {"a", "add an account"}, {"l", "log in again"}, {"s", "the account's groups"},
			{"d", "make it the default"}, {"f", "use it in this folder"},
			{"e", "rename"}, {"c", "shortcut command"}, {"x", "remove"},
		}...)
	case tabGroups:
		rows = append(rows, [][2]string{
			{"enter", "edit items and accounts"}, {"space", "toggle while editing"},
			{"tab", "items or accounts"}, {"n", "new group"}, {"e", "rename"}, {"x", "remove"},
		}...)
	case tabFolders:
		rows = append(rows, [][2]string{{"n", "use an account here"}, {"enter", "change the account"}, {"x", "remove the rule"}}...)
	case tabHealth:
		rows = append(rows, [][2]string{{"↑/↓", "select a box"}, {"f", "repair the selected box"}, {"F", "repair all (doctor --fix)"}}...)
	}
	cell := func(r [2]string) string { return t.accent().Render(padRight(r[0], 8)) + " " + r[1] }
	out := []string{" " + t.bold().Render("Shortcuts"), ""}
	if w < 72 {
		for _, r := range rows {
			out = append(out, "  "+cell(r))
		}
		return out
	}
	half := (len(rows) + 1) / 2
	colW := (w - 6) / 2
	for i := 0; i < half; i++ {
		line := "  " + padRight(ansi.Truncate(cell(rows[i]), colW, symMore), colW)
		if i+half < len(rows) {
			line += "  " + cell(rows[i+half])
		}
		out = append(out, line)
	}
	return out
}
