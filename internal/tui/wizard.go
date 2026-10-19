package tui

import (
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
	Name     string
	Share    config.Share
	Shortcut string // "" = no shortcut
}

// ValidateNewName checks a profile name against the rules and existing names.
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

// ShareFromItems turns the selected item keys into a Share (default group).
func ShareFromItems(items []string) config.Share {
	var s config.Share
	for _, item := range items {
		_ = s.Set(item, "default")
	}
	return s
}

// DefaultItems are pre-selected: everything except conversation history.
func DefaultItems() []string { return []string{"mcp", "settings", "memory", "extensions"} }

// Wizard asks for a name (unless given), the shared items and, when
// WithShortcut set it up, whether to create a shortcut and its name.
type Wizard struct {
	validate  func(string) error
	askName   bool
	name      string
	err       string
	selected  map[string]bool
	cursor    int
	done      bool
	cancelled bool

	checkShortcut func(name, shortcut string) error // nil: no shortcut questions
	suggestFor    func(name string) string          // default shortcut for a name, e.g. ccwork
	suggest       string                            // suggestFor(name) when it passes the check, else ""
	askShortcut   bool                              // "Create a shortcut …? [y/N]"
	typeShortcut  bool                              // "Shortcut name [ccwork]:"
	shortcut      string
}

func NewWizard(cfg *config.Config, name string) Wizard {
	w := Wizard{validate: ValidateNewName(cfg), askName: name == "", name: name, selected: map[string]bool{}}
	for _, item := range DefaultItems() {
		w.selected[item] = true
	}
	return w
}

// WithShortcut adds the shortcut questions after the shared items. check
// validates a shortcut for the account name; suggest gives the default for
// a name, offered only when it passes check.
func (w Wizard) WithShortcut(suggest func(name string) string, check func(name, shortcut string) error) Wizard {
	w.suggestFor, w.checkShortcut = suggest, check
	return w
}

// Answers is valid when Done reports true.
func (w Wizard) Answers() AddAnswers {
	var items []string
	for _, item := range config.Items {
		if w.selected[item] {
			items = append(items, item)
		}
	}
	return AddAnswers{Name: w.name, Share: ShareFromItems(items), Shortcut: w.shortcut}
}

func (w Wizard) Done() bool { return w.done }

func (w Wizard) Init() tea.Cmd { return nil }

func (w Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Terminals deliver pasted text as one message, not as key presses.
	if p, ok := msg.(tea.PasteMsg); ok {
		switch text := strings.TrimSpace(p.Content); {
		case w.askName:
			w.name += text
		case w.typeShortcut:
			w.shortcut += text
		}
		return w, nil
	}
	k, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return w, nil
	}
	switch k.String() {
	case "ctrl+c", "esc":
		w.cancelled = true
		return w, tea.Quit
	}
	if w.askName {
		switch k.String() {
		case "enter":
			if err := w.validate(w.name); err != nil {
				w.err = err.Error()
				return w, nil
			}
			w.askName, w.err = false, ""
		case "backspace":
			if r := []rune(w.name); len(r) > 0 {
				w.name = string(r[:len(r)-1])
			}
		default:
			w.name += k.Text
		}
		return w, nil
	}
	if w.askShortcut {
		switch k.String() {
		case "y", "Y":
			w.askShortcut, w.typeShortcut = false, true
		case "n", "N", "enter":
			w.done = true
			return w, tea.Quit
		}
		return w, nil
	}
	if w.typeShortcut {
		switch k.String() {
		case "enter":
			if w.shortcut == "" {
				w.shortcut = w.suggest
			}
			if w.shortcut == "" {
				w.err = "type a name for the shortcut"
				return w, nil
			}
			if err := w.checkShortcut(w.name, w.shortcut); err != nil {
				w.err, w.shortcut = err.Error(), ""
				return w, nil
			}
			w.done = true
			return w, tea.Quit
		case "backspace":
			if r := []rune(w.shortcut); len(r) > 0 {
				w.shortcut = string(r[:len(r)-1])
			}
		default:
			w.shortcut += k.Text
		}
		return w, nil
	}
	switch k.String() {
	case "up", "k":
		if w.cursor > 0 {
			w.cursor--
		}
	case "down", "j":
		if w.cursor < len(config.Items)-1 {
			w.cursor++
		}
	case "space", " ", "x":
		item := config.Items[w.cursor]
		w.selected[item] = !w.selected[item]
	case "enter":
		if w.checkShortcut != nil {
			w.askShortcut, w.suggest = true, ""
			if s := w.suggestFor(w.name); w.checkShortcut(w.name, s) == nil {
				w.suggest = s
			}
			return w, nil
		}
		w.done = true
		return w, tea.Quit
	}
	return w, nil
}

func (w Wizard) View() tea.View {
	var b strings.Builder
	b.WriteString(titleStyle.Render(" Add account") + "\n\n")
	if w.askName {
		fmt.Fprintf(&b, " Account name: %s█\n", w.name)
		b.WriteString(dimStyle.Render(" a-z, 0-9, - and _ (e.g. work, personal)") + "\n")
		if w.err != "" {
			b.WriteString(" " + warnStyle.Render(w.err) + "\n")
		}
		b.WriteString("\n" + dimStyle.Render(" Enter next · Esc cancel") + "\n")
		return tea.NewView(b.String())
	}
	if w.askShortcut || w.typeShortcut {
		fmt.Fprintf(&b, " Name: %s\n\n", w.name)
		if w.askShortcut {
			fmt.Fprintf(&b, " Create a shortcut command for %s? [y/N]\n", w.name)
			hint := " A command that opens Claude Code as this account."
			if w.suggest != "" {
				hint = " A command such as " + w.suggest + " that opens Claude Code as this account."
			}
			b.WriteString(dimStyle.Render(hint) + "\n")
			b.WriteString("\n" + dimStyle.Render(" y yes · n or Enter no · Esc cancel") + "\n")
			return tea.NewView(b.String())
		}
		rule := " a-z, 0-9, '.', '-' and '_'"
		if w.suggest != "" {
			fmt.Fprintf(&b, " Shortcut name [%s]: %s█\n", w.suggest, w.shortcut)
			rule += "; Enter alone uses " + w.suggest
		} else {
			fmt.Fprintf(&b, " Shortcut name: %s█\n", w.shortcut)
		}
		b.WriteString(dimStyle.Render(rule) + "\n")
		if w.err != "" {
			b.WriteString(" " + warnStyle.Render(w.err) + "\n")
		}
		b.WriteString("\n" + dimStyle.Render(" Enter go to login · Esc cancel") + "\n")
		return tea.NewView(b.String())
	}
	fmt.Fprintf(&b, " Name: %s\n\n ? Choose what to share (Space to toggle)\n", w.name)
	for i, item := range config.Items {
		marker := "  "
		if i == w.cursor {
			marker = cursorStyle.Render("› ")
		}
		box := "[ ]"
		if w.selected[item] {
			box = "[x]"
		}
		fmt.Fprintf(&b, " %s%s %s\n", marker, box, itemLabels[item])
	}
	next := " Enter go to login · Esc cancel"
	if w.checkShortcut != nil {
		next = " Enter next · Esc cancel"
	}
	b.WriteString("\n" + dimStyle.Render(next) + "\n")
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
