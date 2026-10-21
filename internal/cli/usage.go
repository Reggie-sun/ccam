package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func usageCmd(e *env) *cobra.Command {
	var opt usage.Options
	var asJSON, verbose bool
	var watch int
	cmd := &cobra.Command{
		Use:     "usage",
		Aliases: []string{"u"},
		Short:   "5-hour, weekly and per-model usage of every account, with reset times",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opt.Source != "" && opt.Source != "live" && opt.Source != "cache" && opt.Source != "claude" {
				return fmt.Errorf("--source must be live, claude or cache")
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
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show where the numbers came from")
	cmd.Flags().StringVar(&opt.Source, "source", "", "live (ccam asks the usage endpoint), claude (Claude Code does, slower) or cache (saved values) (config: usage_source)")
	cmd.Flags().BoolVar(&opt.NoCache, "no-cache", false, "ignore the 60-second cache")
	cmd.Flags().BoolVar(&opt.Refresh, "refresh", false, "refresh expired tokens (only for accounts not running)")
	return cmd
}

// newFetcher is replaced in tests to point at a fake server.
var newFetcher = usage.NewFetcher

func (e *env) fetcher() *usage.Fetcher {
	f := newFetcher(usageStateDir(e))
	f.Home = e.paths.Home
	if cfg, err := e.load(); err == nil {
		f.Source = cfg.UsageSource
		f.Probe = claudeProbe(e, cfg)
		if e.dryRun {
			// Claude Code writes to the account's folder and may refresh its
			// login, so a dry run goes by the saved usage.
			f.Probe = func(context.Context, usage.Account) ([]byte, error) {
				return nil, errors.New("a dry run does not start Claude Code; the usage shown is the saved one")
			}
		}
	}
	return f
}

// claudeProbe runs Claude Code in an account's folder and asks it for the
// usage (usage.ProbeRequest): with --safe-mode, so the account's hooks, MCP
// servers and plugins do not start, and without saving a session.
func claudeProbe(e *env, cfg *config.Config) func(context.Context, usage.Account) ([]byte, error) {
	return func(ctx context.Context, a usage.Account) ([]byte, error) {
		bin, err := launcher.ResolveClaude(e.paths, cfg.ClaudeBin)
		if err != nil {
			return nil, err
		}
		c := exec.CommandContext(ctx, bin, "-p", "--input-format", "stream-json", "--output-format", "stream-json",
			"--verbose", "--safe-mode", "--no-session-persistence")
		// Its own login only, from a folder that is not a project.
		c.Env = launcher.ProfileEnv(e.paths, launcher.WithoutAuth(os.Environ()), a.Dir, a.Name)
		c.Dir = os.TempDir()
		c.Stdin = strings.NewReader(usage.ProbeRequest)
		c.WaitDelay = 2 * time.Second // a child that keeps the pipe open must not hold the lookup
		return c.Output()
	}
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
	// Lines are cut to the terminal width: a wrapped line would take two
	// rows, and moving up by the number of lines would then land too low.
	cols := 0
	if f, ok := e.stdout.(*os.File); ok && tty {
		if w, _, err := term.GetSize(f.Fd()); err == nil && w > 1 {
			cols = w - 1
		}
	}
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
				if cols > 0 {
					l = ansi.Truncate(l, cols, "…")
				}
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
	// The same rule next --auto and pick use, so the advice matches them.
	p, err := choiceFlags{}.policy(cfg, nil)
	if err != nil {
		p = usage.Policy{}
	}
	// next --auto stays in the organization of the last used account and
	// counts running sessions.
	from, _, _ := next.Rotate(candidates(e, cfg), "", nil)
	org := newOrgs(e, cfg, false)
	var pool []*usage.Report
	for _, r := range view {
		if prof, ok := cfg.Find(r.Name); ok && org.fits(from, prof) {
			r.Sessions = session.Count(launcher.SessionDir(e.paths, r.Name))
			pool = append(pool, r)
		}
	}
	if best := p.Choose(pool, time.Now()); best != nil && len(view) > 1 {
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
