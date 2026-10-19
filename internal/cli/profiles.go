package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/shim"
	"github.com/jhyunwoo/ccam/internal/tui"
)

func addCmd(e *env) *cobra.Command {
	var noLogin bool
	var shares []string
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add an account: create profile → log in → shortcut",
		Long:  "Run on a terminal without options, a wizard asks for the name and shared items.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactive(e) && !cmd.Flags().Changed("share") && !noLogin {
				cfg, err := e.load()
				if err != nil {
					return err
				}
				name := ""
				if len(args) == 1 {
					if err := tui.ValidateNewName(cfg)(args[0]); err != nil {
						return err
					}
					name = args[0]
				}
				return addInteractive(e, cfg, name)
			}
			if len(args) == 0 {
				return fmt.Errorf("an account name is required: ccam add <name>")
			}
			share := config.DefaultShare()
			for _, s := range shares {
				item, group, ok := cutEq(s)
				if !ok {
					return fmt.Errorf("--share takes item=group (e.g. history=team): %q", s)
				}
				if err := share.Set(item, group); err != nil {
					return err
				}
			}
			return addProfile(e, args[0], share, !noLogin)
		},
	}
	cmd.Flags().BoolVar(&noLogin, "no-login", false, "log in later (`ccam login <name>`)")
	cmd.Flags().StringArrayVar(&shares, "share", nil, "share item=group|local (repeatable)")
	return cmd
}

func addProfile(e *env, name string, share config.Share, login bool) error {
	cfg, err := e.load()
	if err != nil {
		return err
	}
	prof, err := profile.Add(e.paths, cfg, name, share)
	if err != nil {
		return err
	}
	if err := applyShare(e, cfg, prof); err != nil {
		return err
	}
	if err := e.save(cfg); err != nil {
		return err
	}
	if login {
		if err := loginProfile(e, cfg, name); err != nil {
			fmt.Fprintf(e.stderr, "Login failed: %v\nTry again later with `ccam login %s`.\n", err, name)
		}
	}
	if _, err := installShims(e, cfg); err != nil {
		fmt.Fprintln(e.stderr, "Could not create shortcuts:", err)
	}
	prof, _ = cfg.Find(name)
	email := prof.Email
	if email == "" {
		email = "not logged in"
	}
	e.printf("✓ %s (%s) added — shortcut: claude%d\n", name, email, prof.Slot)
	return nil
}

// loginProfile runs `claude auth login` inside the profile and caches the email.
func loginProfile(e *env, cfg *config.Config, sel string) error {
	prof, err := cfg.MustFind(sel)
	if err != nil {
		return err
	}
	if err := runClaude(e, cfg, prof, "auth", "login"); err != nil {
		return err
	}
	email, err := profile.ReadEmail(profile.Dir(e.paths, prof))
	if err != nil {
		return fmt.Errorf("could not read the email after login: %w", err)
	}
	prof.Email = email
	return e.save(cfg)
}

// runClaude runs a claude subcommand (not exec) inside a profile.
func runClaude(e *env, cfg *config.Config, prof *config.Profile, args ...string) error {
	bin, err := launcher.ResolveClaude(e.paths, cfg.ClaudeBin)
	if err != nil {
		return err
	}
	c := exec.Command(bin, args...)
	c.Env = launcher.ProfileEnv(e.paths, os.Environ(), profile.Dir(e.paths, prof), prof.Name)
	c.Stdin, c.Stdout, c.Stderr = e.stdin, e.stdout, e.stderr
	return c.Run()
}

func installShims(e *env, cfg *config.Config) ([]string, error) {
	self, err := ccamPath()
	if err != nil {
		return nil, err
	}
	created, err := shim.Install(e.paths.Bin, self, cfg)
	if err == nil && !shim.InPath(e.paths.Bin) {
		fmt.Fprintf(e.stderr, "Note: %s is not in PATH. Add it to ~/.bashrc or ~/.zshrc:\n  export PATH=\"%s:$PATH\"\n", e.paths.Bin, e.paths.Bin)
	}
	return created, err
}

func cutEq(s string) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}

type listRow struct {
	Slot     int          `json:"slot"`
	Name     string       `json:"name"`
	Email    string       `json:"email"`
	Default  bool         `json:"default"`
	Dir      string       `json:"dir"`
	LastUsed *time.Time   `json:"last_used,omitempty"`
	Share    config.Share `json:"share"`
}

func listCmd(e *env) *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List accounts",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listProfiles(e, asJSON)
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "print as JSON")
	return cmd
}

func listProfiles(e *env, asJSON bool) error {
	cfg, err := e.load()
	if err != nil {
		return err
	}
	def, _ := cfg.DefaultProfile()
	rows := make([]listRow, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		r := listRow{Slot: p.Slot, Name: p.Name, Email: p.Email, Default: def == p, Dir: profile.Dir(e.paths, p), Share: p.Share}
		if t := launcher.LastUsed(e.paths, p.Name); !t.IsZero() {
			r.LastUsed = &t
		}
		rows = append(rows, r)
	}
	if asJSON {
		enc := json.NewEncoder(e.stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	}
	if len(rows) == 0 {
		e.printf("No accounts yet. Add one with `ccam add <name>`.\n")
		return nil
	}
	tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, " \tShortcut\tName\tEmail\tShared\tLast used")
	for _, r := range rows {
		mark := ""
		if r.Default {
			mark = "*"
		}
		email := r.Email
		if email == "" {
			email = "(not logged in)"
		}
		last := "-"
		if r.LastUsed != nil {
			last = Ago(time.Since(*r.LastUsed))
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", mark, "claude"+strconv.Itoa(r.Slot), r.Name, email, shareSummary(r.Share), last)
	}
	return tw.Flush()
}

// shareSummary renders e.g. "mcp settings memory ext" for shared items.
func shareSummary(s config.Share) string {
	short := map[string]string{"mcp": "mcp", "settings": "settings", "memory": "md", "extensions": "ext", "history": "history"}
	out := ""
	for _, item := range config.Items {
		g := s.Get(item)
		if g == config.Local {
			continue
		}
		if out != "" {
			out += " "
		}
		out += short[item]
		if g != "default" {
			out += "@" + g
		}
	}
	if out == "" {
		return "-"
	}
	return out
}

// Ago renders a short relative time.
func Ago(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func useCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "use <name|slot>",
		Short: "Set the default account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			p, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			cfg.Default = p.Name
			if err := e.save(cfg); err != nil {
				return err
			}
			e.printf("✓ Default account: %s\n", p.Name)
			return nil
		},
	}
}

func removeCmd(e *env) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:     "remove <name|slot>",
		Aliases: []string{"rm"},
		Short:   "Remove an account (its folder goes to the trash)",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			p, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			prof := *p
			if !yes && !e.confirm(fmt.Sprintf("Remove account %s (%s)?", prof.Name, profile.Dir(e.paths, &prof))) {
				e.printf("Cancelled.\n")
				return nil
			}
			if g := prof.Share.Get("extensions"); g != config.Local {
				// Other accounts' plugin paths may run through this folder.
				sep := string(filepath.Separator)
				groupPlugins := filepath.Join(e.paths.GroupDir(g), "plugins")
				if _, err := profile.RewritePluginPaths(groupPlugins,
					filepath.Join(profile.Dir(e.paths, &prof), "plugins")+sep, groupPlugins+sep); err != nil {
					return err
				}
			}
			trashed, err := profile.Trash(e.paths, &prof)
			if err != nil {
				return err
			}
			cfg.Remove(prof.Name)
			if err := e.save(cfg); err != nil {
				return err
			}
			if self, err := ccamPath(); err == nil {
				shim.RemoveFor(e.paths.Bin, self, &prof)
			}
			if trashed != "" {
				e.printf("✓ %s removed (folder moved to %s)\n", prof.Name, trashed)
			} else {
				e.printf("✓ %s unregistered (folder %s left in place)\n", prof.Name, profile.Dir(e.paths, &prof))
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "remove without asking")
	return cmd
}

func renameCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "rename <name> <new-name>",
		Short: "Rename an account (updates shortcuts)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			old, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			oldName := old.Name
			if err := profile.Rename(e.paths, cfg, oldName, args[1]); err != nil {
				return err
			}
			if err := e.save(cfg); err != nil {
				return err
			}
			if _, err := installShims(e, cfg); err != nil {
				fmt.Fprintln(e.stderr, "Could not update shortcuts:", err)
			}
			e.printf("✓ %s → %s\n", oldName, args[1])
			return nil
		},
	}
}

func loginCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "login <name|slot>",
		Short: "Log in again",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			if err := loginProfile(e, cfg, args[0]); err != nil {
				return err
			}
			p, _ := cfg.Find(args[0])
			e.printf("✓ %s logged in: %s\n", p.Name, p.Email)
			return nil
		},
	}
}

func logoutCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "logout <name|slot>",
		Short: "Remove the login token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			p, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			if err := runClaude(e, cfg, p, "auth", "logout"); err != nil {
				return err
			}
			p.Email = ""
			return e.save(cfg)
		},
	}
}

func shimsCmd(e *env) *cobra.Command {
	cmd := &cobra.Command{Use: "shims", Short: "Manage shortcuts (claude1…)"}
	var names bool
	install := &cobra.Command{
		Use:   "install",
		Short: "Recreate claude1…claudeN",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			if cmd.Flags().Changed("names") && cfg.NameShims != names {
				cfg.NameShims = names
				if err := e.save(cfg); err != nil {
					return err
				}
			}
			created, err := installShims(e, cfg)
			if err != nil {
				return err
			}
			e.printf("✓ %d new shortcuts created (%s)\n", len(created), e.paths.Bin)
			return nil
		},
	}
	install.Flags().BoolVar(&names, "names", false, "also create claude-<name> shortcuts")
	cmd.AddCommand(install)
	return cmd
}

func runCmd(e *env) *cobra.Command {
	// Only reached for `ccam help run`; real runs never get to cobra.
	return &cobra.Command{
		Use:                "run <name|slot> [-- claude args]",
		Short:              "Run claude as that account (same as `ccam <name>`)",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			sel := ""
			if len(args) > 0 && args[0] != "--" && args[0][0] != '-' {
				sel, args = args[0], args[1:]
			}
			return launcher.Launch(e.paths, sel, args, e.dryRun, e.stdout)
		},
	}
}
