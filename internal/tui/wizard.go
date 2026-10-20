package tui

import (
	"cmp"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/jhyunwoo/ccam/internal/config"
)

// The wizard is plain Bubble Tea on purpose: huh pulls in a clipboard
// package whose init() searches PATH, which cost ~110 ms on every launch.

// AddAnswers is what the add wizard collects.
type AddAnswers struct {
	Name   string
	Groups []string // groups to join
	// Items is what a new group default shares, when there were no groups.
	Items    []string
	Shortcut string      // "" = no shortcut
	Login    bool        // log in now (a Claude login account)
	API      *config.API // an API key account
}

// ValidateNewName checks an account name against the rules and existing names.
func ValidateNewName(cfg *config.Config) func(string) error {
	return func(name string) error {
		if err := config.ValidateName(name); err != nil {
			return err
		}
		if _, ok := cfg.Find(name); ok {
			return fmt.Errorf("an account with this name already exists")
		}
		return nil
	}
}

// DefaultItems are pre-selected: everything except conversation history.
func DefaultItems() []string { return config.DefaultItems }

type step int

const (
	stepName step = iota
	stepKind
	stepKey
	stepURL
	stepAuth
	stepModel
	stepGroups
	stepAskShortcut
	stepShortcut
	stepSummary
)

const (
	kindLogin = iota
	kindAPI
	kindLater
)

var kindLabels = []string{
	"Claude subscription: log in now in the browser",
	"API key: Anthropic's API or a compatible gateway",
	"Claude subscription: log in later (ccam login <name>)",
}

// Wizard asks for a name (unless given), the kind of account, the groups
// to join (or, before there are groups, what group default shares) and,
// when WithShortcut set it up, a shortcut command; then shows a summary.
// Esc goes back one step.
type Wizard struct {
	validate  func(string) error
	theme     Theme
	themeName string // the theme setting; "auto" follows the background
	step      step
	back      []step
	askedName bool
	name      input
	kind      int
	cursor    int // in choice steps
	keyFrom   input
	baseURL   input
	auth      int
	model     input
	list      checklist
	groups    []GroupInfo // nil: no groups yet, the list holds items
	shortcut  input
	err       string
	done      bool
	cancelled bool

	checkShortcut func(name, shortcut string) error // nil: no shortcut questions
	suggestFor    func(name string) string          // default shortcut for a name, e.g. ccwork
}

func NewWizard(cfg *config.Config, name string) Wizard {
	w := Wizard{validate: ValidateNewName(cfg), themeName: cfg.Theme, theme: ThemeFor(cfg.Theme, true), name: newInput(name, ""),
		keyFrom: newInput("", "env:ANTHROPIC_API_KEY or command:pass show anthropic"),
		baseURL: newInput("", "https://api.anthropic.com"), model: newInput("", "Claude Code's default"),
		list: itemChecklist(DefaultItems())}
	if name == "" {
		w.askedName = true
	} else {
		w.step = stepKind
	}
	if len(cfg.Groups) > 0 {
		w.groups = GroupInfos(cfg)
		w.list = groupChecklist(w.groups, cfg.DefaultGroups())
	}
	return w
}

// WithShortcut adds the shortcut questions after the groups. check
// validates a shortcut for the account name; suggest gives the default for
// a name, offered only when it passes check.
func (w Wizard) WithShortcut(suggest func(name string) string, check func(name, shortcut string) error) Wizard {
	w.suggestFor, w.checkShortcut = suggest, check
	return w
}

// Answers is valid when Done reports true.
func (w Wizard) Answers() AddAnswers {
	a := AddAnswers{Name: w.name.Value(), Shortcut: w.shortcut.Value(), Login: w.kind == kindLogin}
	switch picked := w.list.items(); {
	case w.groups != nil:
		a.Groups = picked
	case len(picked) > 0:
		a.Groups, a.Items = []string{"default"}, picked
	}
	if w.kind == kindAPI {
		a.API = &config.API{KeyFrom: w.keyFrom.Value(), BaseURL: w.baseURL.Value(), Model: w.model.Value()}
		if w.auth == 1 {
			a.API.Auth = config.Bearer
		}
	}
	return a
}

func (w Wizard) Done() bool { return w.done }

func (w Wizard) Init() tea.Cmd { return tea.RequestBackgroundColor }

func (w *Wizard) next(s step) {
	w.back = append(w.back, w.step)
	w.step, w.err, w.cursor = s, "", 0
}

func (w Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if bg, ok := msg.(tea.BackgroundColorMsg); ok {
		w.theme = ThemeFor(w.themeName, bg.IsDark())
		return w, nil
	}
	quit := w.update(msg)
	if quit {
		return w, tea.Quit
	}
	return w, nil
}

// update handles one message and reports whether the wizard is finished
// (done or cancelled).
func (w *Wizard) update(msg tea.Msg) bool {
	k, isKey := msg.(tea.KeyPressMsg)
	if isKey {
		switch k.String() {
		case "ctrl+c":
			w.cancelled = true
			return true
		case "esc":
			if len(w.back) == 0 {
				w.cancelled = true
				return true
			}
			w.step, w.back, w.err = w.back[len(w.back)-1], w.back[:len(w.back)-1], ""
			return false
		}
	}
	field := map[step]*input{stepName: &w.name, stepKey: &w.keyFrom, stepURL: &w.baseURL, stepModel: &w.model, stepShortcut: &w.shortcut}[w.step]
	if field != nil && !(isKey && k.String() == "enter") {
		if field.update(msg) {
			w.err = ""
		}
		return false
	}
	if !isKey {
		return false
	}
	key := k.String()
	switch w.step {
	case stepName:
		if err := w.validate(w.name.Value()); err != nil {
			w.err = err.Error()
			return false
		}
		w.next(stepKind)
	case stepKind:
		if pickKey(key, &w.cursor, len(kindLabels)) {
			w.kind = w.cursor
			if w.kind == kindAPI {
				w.next(stepKey)
			} else {
				w.next(stepGroups)
			}
		}
	case stepKey:
		api := config.API{KeyFrom: w.keyFrom.Value()}
		if err := api.Check(); err != nil {
			w.err = err.Error()
			return false
		}
		w.next(stepURL)
	case stepURL:
		api := config.API{KeyFrom: "env:X", BaseURL: w.baseURL.Value()}
		if err := api.Check(); err != nil {
			w.err = err.Error()
			return false
		}
		w.next(stepAuth)
	case stepAuth:
		if pickKey(key, &w.cursor, 2) {
			w.auth = w.cursor
			w.next(stepModel)
		}
	case stepModel:
		w.next(stepGroups)
	case stepGroups:
		if w.list.key(key) {
			w.err = ""
			return false
		}
		if key != "enter" {
			return false
		}
		if w.groups != nil {
			if w.err = overlap(w.groups, w.list.items()); w.err != "" {
				return false
			}
		}
		if w.checkShortcut != nil {
			w.next(stepAskShortcut)
		} else {
			w.next(stepSummary)
		}
	case stepAskShortcut:
		switch key {
		case "y", "Y":
			s := w.suggestFor(w.name.Value())
			if w.checkShortcut(w.name.Value(), s) != nil {
				s = ""
			}
			w.shortcut = newInput(s, "cc"+w.name.Value())
			w.next(stepShortcut)
		case "n", "N", "enter":
			w.shortcut = newInput("", "")
			w.next(stepSummary)
		}
	case stepShortcut:
		if w.shortcut.Value() == "" {
			w.err = "type a name for the shortcut, or Esc to go back"
			return false
		}
		if err := w.checkShortcut(w.name.Value(), w.shortcut.Value()); err != nil {
			w.err = err.Error()
			return false
		}
		w.next(stepSummary)
	case stepSummary:
		if pickKey(key, &w.cursor, 2) {
			if w.cursor == 1 {
				w.cancelled = true
			} else {
				w.done = true
			}
			return true
		}
	}
	return false
}

// lines is the current step, without a title.
func (w Wizard) lines() []string {
	t := w.theme
	var out []string
	add := func(s ...string) { out = append(out, s...) }
	field := func(question string, in input, hint string) {
		add(question, "  "+t.subtle().Render("› ")+in.view(t))
		if hint != "" {
			add("  " + t.dim().Render(hint))
		}
	}
	name := w.name.Value()
	if w.step != stepName {
		add(t.dim().Render("Account: ") + name)
		add("")
	}
	switch w.step {
	case stepName:
		field("Account name", w.name, "a-z, 0-9, - and _ (e.g. work, personal)")
	case stepKind:
		add("How does "+name+" sign in?", "")
		add(options(t, kindLabels, w.cursor)...)
	case stepKey:
		field("Where does ccam read the API key at each launch?", w.keyFrom,
			"env:NAME (a variable) or command:COMMAND (the first line it prints); the key is never stored")
	case stepURL:
		field("Endpoint (empty: Anthropic's API)", w.baseURL, "for a gateway such as LiteLLM or OpenRouter")
	case stepAuth:
		add("How is the key sent?", "")
		add(options(t, []string{"x-api-key header (ANTHROPIC_API_KEY), Anthropic's API", "Authorization: Bearer (ANTHROPIC_AUTH_TOKEN), most gateways"}, w.cursor)...)
	case stepModel:
		field("Model (empty: Claude Code's default)", w.model, "e.g. claude-sonnet-5-5")
	case stepGroups:
		var b strings.Builder
		if w.groups != nil {
			add("Choose the groups to join (Space to toggle)")
			members := map[string]string{}
			for _, g := range w.groups {
				members[g.Name] = strings.Join(g.Members, ", ")
			}
			w.list.styled(&b, t, func(n string) string { return members[n] })
		} else {
			add("Choose what to share with the accounts you add (group default)")
			w.list.styled(&b, t, nil)
		}
		for _, l := range strings.Split(strings.TrimRight(b.String(), "\n"), "\n") {
			add(l)
		}
		add("  " + t.dim().Render("Login tokens are never shared."))
	case stepAskShortcut:
		hint := "A command that opens Claude Code as this account."
		if s := w.suggestFor(name); w.checkShortcut(name, s) == nil {
			hint = "A command such as " + s + " that opens Claude Code as this account."
		}
		add(fmt.Sprintf("Create a shortcut command for %s? [y/N]", name), "  "+t.dim().Render(hint))
	case stepShortcut:
		field("Shortcut name", w.shortcut, "a-z, 0-9, '.', '-' and '_'")
	case stepSummary:
		a := w.Answers()
		kind := []string{"Claude login, now", "API key", "Claude login, later"}[w.kind]
		add(t.dim().Render("Sign-in   ") + kind)
		if a.API != nil {
			add(t.dim().Render("Key from  ")+a.API.KeyFrom, t.dim().Render("Endpoint  ")+a.API.Label())
			if a.API.Model != "" {
				add(t.dim().Render("Model     ") + a.API.Model)
			}
		}
		groups := "none"
		if len(a.Groups) > 0 {
			groups = strings.Join(a.Groups, ", ")
		}
		if a.Items != nil {
			groups += " (shares " + strings.Join(a.Items, ", ") + ")"
		}
		add(t.dim().Render("Groups    ")+groups, t.dim().Render("Shortcut  ")+cmp.Or(a.Shortcut, "none"), "", "Add this account?")
		add(options(t, []string{"Yes", "No"}, w.cursor)...)
	}
	if w.err != "" {
		add("  " + t.error().Render(symFail+" "+w.err))
	}
	return out
}

func (w Wizard) keys() string {
	back := "esc to go back"
	if len(w.back) == 0 {
		back = "esc to cancel"
	}
	switch w.step {
	case stepKind, stepAuth, stepSummary:
		return "↑/↓ to select" + symDot + "enter to choose" + symDot + back
	case stepGroups:
		return "space to toggle" + symDot + "enter for the next step" + symDot + back
	case stepAskShortcut:
		return "y for yes" + symDot + "n or enter for no" + symDot + back
	}
	return "enter for the next step" + symDot + back
}

func (w Wizard) View() tea.View {
	var b strings.Builder
	b.WriteString(" " + w.theme.bold().Render("Add account") + "\n\n")
	for _, l := range w.lines() {
		b.WriteString(" " + l + "\n")
	}
	b.WriteString("\n " + w.theme.dim().Render(w.keys()) + "\n")
	return tea.NewView(b.String())
}

// RunAddWizard runs the wizard on the terminal, with the shortcut questions.
func RunAddWizard(cfg *config.Config, name string, suggest func(name string) string, check func(name, shortcut string) error) (AddAnswers, error) {
	wiz := NewWizard(cfg, name).WithShortcut(suggest, check)
	final, err := tea.NewProgram(wiz, tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout)).Run()
	if err != nil {
		return AddAnswers{}, err
	}
	w := final.(Wizard)
	if !w.Done() {
		return AddAnswers{}, fmt.Errorf("cancelled")
	}
	return w.Answers(), nil
}

// wizardDialog is the wizard inside the dashboard.
type wizardDialog struct{ w Wizard }

func newWizardDialog(d *Dashboard) dialog {
	w := NewWizard(d.cfg, "").WithShortcut(d.be.SuggestShortcut, d.be.CheckShortcut)
	w.theme = d.theme
	return &wizardDialog{w: w}
}

func (wd *wizardDialog) update(d *Dashboard, msg tea.Msg) tea.Cmd {
	wd.w.theme = d.theme
	if !wd.w.update(msg) {
		return nil
	}
	d.dlg = nil
	if wd.w.done {
		return d.quit(Result{Action: Add, Add: wd.w.Answers()})
	}
	return nil
}

func (wd *wizardDialog) view(d *Dashboard, w int) []string {
	out := dialogTop(d.theme, w, "Add account")
	for _, l := range wd.w.lines() {
		out = append(out, "  "+l)
	}
	return out
}

func (wd *wizardDialog) keys() string { return wd.w.keys() }
