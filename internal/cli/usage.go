package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func usageCmd(e *env) *cobra.Command {
	var opt usage.Options
	var asJSON, verbose bool
	var watch int
	cmd := &cobra.Command{
		Use:     "usage",
		Aliases: []string{"u"},
		Short:   "5-hour and weekly usage of every account, with reset times",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opt.Source != "" && opt.Source != "live" && opt.Source != "cache" {
				return fmt.Errorf("--source must be live or cache")
			}
			if e.dryRun && opt.Refresh {
				return fmt.Errorf("--dry-run cannot be combined with --refresh, which rewrites login tokens")
			}
			cfg, err := e.load()
			if err != nil {
				return err
			}
			if len(cfg.Profiles) == 0 {
				e.printf("No accounts yet. Add one with `ccam add <name>`.\n")
				return nil
			}
			f := e.fetcher()
			if !cmd.Flags().Changed("watch") {
				return showUsage(context.Background(), e, f, cfg, opt, asJSON, verbose, false)
			}
			if watch <= 0 {
				watch = 60
			}
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()
			for {
				if err := showUsage(ctx, e, f, cfg, opt, asJSON, verbose, true); err != nil {
					return err
				}
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(time.Duration(watch) * time.Second):
				}
			}
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "print as JSON")
	cmd.Flags().IntVar(&watch, "watch", 60, "refresh every N seconds")
	cmd.Flags().Lookup("watch").NoOptDefVal = "60"
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show per-model limits and the data source")
	cmd.Flags().StringVar(&opt.Source, "source", "", "live|cache")
	cmd.Flags().BoolVar(&opt.NoCache, "no-cache", false, "ignore the 60-second cache")
	cmd.Flags().BoolVar(&opt.Refresh, "refresh", false, "refresh expired tokens (only for accounts not running)")
	return cmd
}

// newFetcher is replaced in tests to point at a fake server.
var newFetcher = usage.NewFetcher

func (e *env) fetcher() *usage.Fetcher {
	f := newFetcher(usageStateDir(e))
	f.Home = e.paths.Home
	return f
}

func usageStateDir(e *env) string { return filepath.Join(e.paths.State, "usage") }

func usageAccounts(e *env, cfg *config.Config) []usage.Account {
	out := make([]usage.Account, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		out = append(out, usage.Account{
			Name: p.Name, Slot: p.Slot, Dir: profile.Dir(e.paths, p),
			SessionDir: launcher.SessionDir(e.paths, p.Name),
			Imported:   p.Dir != "",
			API:        p.API != nil,
		})
	}
	return out
}

func showUsage(ctx context.Context, e *env, f *usage.Fetcher, cfg *config.Config, opt usage.Options, asJSON, verbose, clear bool) error {
	accts := usageAccounts(e, cfg)
	if asJSON {
		reports := f.FetchAll(ctx, accts, opt, nil)
		enc := json.NewEncoder(e.stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(reports)
	}
	tty := isTerminal(e.stdout)
	ro := usage.RenderOptions{Color: tty && os.Getenv("NO_COLOR") == "", Verbose: verbose}
	view := make([]*usage.Report, len(accts))
	for i, a := range accts {
		view[i] = &usage.Report{Name: a.Name, Slot: a.Slot, Status: usage.StatusPending}
	}
	if clear && tty {
		fmt.Fprint(e.stdout, "\x1b[H\x1b[2J")
	}
	// On a terminal, draw placeholders and redraw as each account answers.
	drawn := 0
	draw := func() {
		ro.Now = time.Now()
		lines := usage.Table(view, ro)
		if drawn > 0 {
			fmt.Fprintf(e.stdout, "\x1b[%dA", drawn)
		}
		for _, l := range lines {
			if tty {
				fmt.Fprint(e.stdout, "\x1b[2K")
			}
			fmt.Fprintln(e.stdout, l)
		}
		drawn = len(lines)
	}
	if tty {
		draw()
	}
	f.FetchAll(ctx, accts, opt, func(i int, r usage.Report) {
		view[i] = &r
		if tty {
			draw()
		}
	})
	if !tty {
		draw()
	}
	if best := usage.Best(view); best != nil && len(view) > 1 {
		e.printf("→ Most headroom: %s   (open it with ccam next --auto)\n", best.Name)
	}
	return nil
}

// isTerminal reports whether w is a terminal. A character-device check is
// not enough: /dev/null is one, and `ccam add work >/dev/null` must not
// start the wizard.
func isTerminal(w any) bool {
	f, ok := w.(*os.File)
	return ok && term.IsTerminal(f.Fd())
}

func statuslineCmd(e *env) *cobra.Command {
	var wrap string
	cmd := &cobra.Command{
		Use:   "statusline",
		Short: "Claude Code status line: shows the account name and collects usage",
		Long:  "Claude Code runs this as its status line command. Install with `ccam statusline install`.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusline(e, wrap)
		},
	}
	cmd.Flags().StringVar(&wrap, "wrap", "", "existing status line command (its output gets the account name in front)")
	cmd.AddCommand(statuslineInstallCmd(e))
	return cmd
}

func runStatusline(e *env, wrap string) error {
	input, err := readAll(e.in)
	if err != nil {
		return err
	}
	name := os.Getenv("CCAM_PROFILE")
	if name == "" {
		name = profileByConfigDir(e, os.Getenv("CLAUDE_CONFIG_DIR"))
	}
	limits, model, ok, _ := usage.ParseStatusLine(input)
	if ok && name != "" {
		_ = usage.SaveStatusLine(usageStateDir(e), name, limits, time.Now())
	}
	label := name
	if label == "" {
		label = "claude"
	}
	if wrap != "" {
		out, _ := runShell(wrap, input)
		first, _, _ := strings.Cut(strings.TrimRight(out, "\n"), "\n")
		e.printf("[%s] %s\n", label, first)
		return nil
	}
	parts := []string{"[" + label + "]"}
	if model != "" {
		parts = append(parts, model)
	}
	if limits.FiveHour.Known {
		parts = append(parts, fmt.Sprintf("5h %.0f%%", limits.FiveHour.Percent))
	}
	if limits.Weekly.Known {
		parts = append(parts, fmt.Sprintf("7d %.0f%%", limits.Weekly.Percent))
	}
	if nearLimit(limits) && hasOtherAccount(e) {
		parts = append(parts, "→ ccam next --auto --continue")
	}
	e.printf("%s\n", strings.Join(parts, " · "))
	return nil
}

// hintAt is the usage at which the status line suggests the next account.
const hintAt = 90

func nearLimit(l usage.Limits) bool {
	return l.FiveHour.Known && l.FiveHour.Percent >= hintAt || l.Weekly.Known && l.Weekly.Percent >= hintAt
}

func hasOtherAccount(e *env) bool {
	cfg, err := e.load()
	return err == nil && len(cfg.Profiles) > 1
}

func profileByConfigDir(e *env, dir string) string {
	if dir == "" {
		return ""
	}
	cfg, err := e.load()
	if err != nil {
		return ""
	}
	for i := range cfg.Profiles {
		if profile.Dir(e.paths, &cfg.Profiles[i]) == dir {
			return cfg.Profiles[i].Name
		}
	}
	return ""
}
