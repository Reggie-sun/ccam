package cli

import (
	"context"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/share"
	"github.com/jhyunwoo/ccam/internal/shim"
	"github.com/jhyunwoo/ccam/internal/tui"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// interactive reports whether both ends are a terminal.
func interactive(e *env) bool { return isTerminal(e.stdin) && isTerminal(e.stdout) }

// runHome is bare `ccam`: the picker on a terminal, the list otherwise.
func runHome(e *env) error {
	if !interactive(e) {
		return listProfiles(e, false)
	}
	for {
		cfg, err := e.load()
		if err != nil {
			return err
		}
		res, err := pick(e, cfg)
		if err != nil {
			return err
		}
		switch res.Action {
		case tui.Launch:
			return launcher.LaunchWith(e.paths, cfg, res.Profile, nil, e.dryRun, e.stdout)
		case tui.Add:
			return addInteractive(e, cfg, "")
		case tui.SetDefault:
			cfg.Default = res.Profile
			if err := e.save(cfg); err != nil {
				return err
			}
		case tui.Share:
			prof, err := cfg.MustFind(res.Profile)
			if err != nil {
				return err
			}
			if err := changeShare(e, cfg, prof, res.Share, share.Options{}, false, nil); err != nil {
				return err
			}
			e.printf("Press Enter to return to the list.")
			_, _ = e.in.ReadString('\n')
		default:
			return nil
		}
	}
}

func pick(e *env, cfg *config.Config) (tui.Result, error) {
	def, _ := cfg.DefaultProfile()
	accts := map[string]usage.Account{}
	for _, a := range usageAccounts(e, cfg) {
		accts[a.Name] = a
	}
	rows := make([]tui.Row, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		last := "-"
		if t := launcher.LastUsed(e.paths, p.Name); !t.IsZero() {
			last = Ago(time.Since(t))
		}
		rows = append(rows, tui.Row{Name: p.Name, Slot: p.Slot, Email: p.Email, Share: p.Share, LastUsed: last, Default: def == p})
	}
	f := e.fetcher()
	fetch := func(ctx context.Context, name string) usage.Report {
		return f.Fetch(ctx, accts[name], usage.Options{})
	}
	final, err := tea.NewProgram(tui.NewPicker(rows, fetch), tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout)).Run()
	if err != nil {
		return tui.Result{}, err
	}
	return final.(tui.Model).Result(), nil
}

// addInteractive runs the wizard, then the same steps as `ccam add`.
func addInteractive(e *env, cfg *config.Config, name string) error {
	check := func(name, shortcut string) error { return checkShortcut(e, cfg, name, shortcut) }
	ans, err := tui.RunAddWizard(cfg, name, shim.Suggest, check)
	if err != nil {
		return err
	}
	e.printf("→ Log in in your browser. ccam continues when you are done.\n")
	return addProfile(e, ans.Name, ans.Share, true, ans.Shortcut)
}
