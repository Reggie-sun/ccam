package cli

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/migrate"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/shim"
	"github.com/jhyunwoo/ccam/internal/tui"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func addCmd(e *env) *cobra.Command {
	var noLogin bool
	var shares []string
	var shortcut string
	var api config.API
	var useAPI bool
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add an account: create profile → log in → optional shortcut",
		Long: "Run on a terminal without options, a wizard asks for the name, the shared items\n" +
			"and whether to create a shortcut command (and its name).\n" +
			"Otherwise no shortcut is created unless --shortcut names one.\n\n" +
			"--api adds an account that uses an API key instead of a Claude login: Anthropic's\n" +
			"API, or a compatible endpoint (--base-url) such as LiteLLM or OpenRouter. The key\n" +
			"is not stored; --key-from says where ccam reads it at each launch:\n" +
			"  env:NAME          an environment variable\n" +
			"  command:COMMAND   the first line a command prints (a password manager, a keychain)",
		Example: "  ccam add\n  ccam add work --shortcut ccwork\n  ccam add side --no-login --share history=team\n" +
			"  ccam add api --api --key-from env:ANTHROPIC_KEY_WORK\n" +
			"  ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from 'command:pass show llm/key'",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, f := range []string{"base-url", "key-from", "auth", "model"} {
				if cmd.Flags().Changed(f) && !useAPI {
					return fmt.Errorf("--%s works with --api", f)
				}
			}
			if useAPI {
				if len(args) == 0 {
					return fmt.Errorf("an account name is required: ccam add <name> --api --key-from …")
				}
				if api.KeyFrom == "" {
					return fmt.Errorf("--api needs --key-from env:NAME or command:COMMAND")
				}
				if err := api.Check(); err != nil {
					return err
				}
			}
			if interactive(e) && !useAPI && !cmd.Flags().Changed("share") && !cmd.Flags().Changed("shortcut") && !noLogin {
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
				item, group, ok := strings.Cut(s, "=")
				if !ok {
					return fmt.Errorf("--share takes item=group (e.g. history=team): %q", s)
				}
				if err := share.Set(item, group); err != nil {
					return err
				}
			}
			if useAPI {
				return addAPIProfile(e, args[0], share, shortcut, &api)
			}
			return addProfile(e, args[0], share, !noLogin, shortcut, nil)
		},
	}
	cmd.Flags().BoolVar(&useAPI, "api", false, "an API key account instead of a Claude login")
	cmd.Flags().StringVar(&api.BaseURL, "base-url", "", "with --api: the endpoint (default: Anthropic's API)")
	cmd.Flags().StringVar(&api.KeyFrom, "key-from", "", "with --api: env:NAME or command:COMMAND")
	cmd.Flags().StringVar(&api.Auth, "auth", "", "with --api: api-key (x-api-key, default) or bearer")
	cmd.Flags().StringVar(&api.Model, "model", "", "with --api: the model to use (ANTHROPIC_MODEL)")
	cmd.Flags().BoolVar(&noLogin, "no-login", false, "log in later with ccam login <name>")
	cmd.Flags().StringArrayVar(&shares, "share", nil, "share item=group|local (repeatable)")
	cmd.Flags().StringVar(&shortcut, "shortcut", "", "also create this shortcut command (e.g. ccwork); none without it")
	return cmd
}

// checkShortcut validates a shortcut name for account prof.
func checkShortcut(e *env, cfg *config.Config, prof, name string) error {
	self, _ := ccamPath()
	return shim.Check(name, e.paths.Bin, self, cfg, prof)
}

// addAPIProfile adds an API key account: no login, the key read at launch.
func addAPIProfile(e *env, name string, share config.Share, shortcut string, api *config.API) error {
	if _, err := ReadKeyForCheck(api); err != nil {
		fmt.Fprintf(e.stderr, "Note: %v. Set it before starting the account.\n", err)
	}
	return addProfile(e, name, share, false, shortcut, api)
}

// ReadKeyForCheck reads an env: key to tell early that it is missing;
// command: sources are not run outside a launch.
func ReadKeyForCheck(api *config.API) (string, error) {
	if !strings.HasPrefix(api.KeyFrom, "env:") {
		return "", nil
	}
	return launcher.ReadAPIKey(api.KeyFrom)
}

func addProfile(e *env, name string, share config.Share, login bool, shortcut string, api *config.API) error {
	cfg, err := e.load()
	if err != nil {
		return err
	}
	if shortcut != "" {
		if err := checkShortcut(e, cfg, name, shortcut); err != nil {
			return err
		}
	}
	prof, err := profile.Add(e.paths, cfg, name, share)
	if err != nil {
		return err
	}
	prof.Shortcut = shortcut
	prof.API = api
	if err := applyShare(e, cfg, prof); err != nil {
		// The folder is new and unregistered; left behind, it would block
		// every retry with "folder already exists".
		os.RemoveAll(e.paths.ProfileDir(name))
		return err
	}
	if err := e.save(cfg); err != nil {
		return err
	}
	if login {
		if _, err := loginProfile(e, name); err != nil {
			fmt.Fprintf(e.stderr, "Login failed: %v\nTry again later with `ccam login %s`.\n", err, name)
		}
		if cfg, err = e.load(); err != nil { // the login saved the email
			return err
		}
	}
	if _, err := installShims(e, cfg); err != nil {
		fmt.Fprintln(e.stderr, "Could not create the shortcut:", err)
	}
	prof, err = cfg.MustFind(name)
	if err != nil {
		return err
	}
	e.printf("✓ %s (%s) added\n", name, accountLabel(prof))
	printShortcut(e, prof)
	if claudeDir := filepath.Join(e.paths.Home, ".claude"); len(cfg.Profiles) == 1 && migrate.HasContent(claudeDir) {
		e.printf("Tip: run `ccam migrate` to copy your existing Claude Code settings (%s) into the shared settings.\n", claudeDir)
	}
	return nil
}

// loginProfile runs `claude auth login` inside the profile and caches the
// email. The login can take minutes, so the config is reloaded to save it.
func loginProfile(e *env, sel string) (*config.Profile, error) {
	cfg, prof, err := e.loadProfile(sel)
	if err != nil {
		return nil, err
	}
	if prof.API != nil {
		return nil, fmt.Errorf("%s is an API key account (%s): it has no login; its key comes from %s", prof.Name, prof.API.Label(), prof.API.KeyFrom)
	}
	name := prof.Name
	if err := runLogin(e, cfg, prof); err != nil {
		return nil, err
	}
	email, err := profile.ReadEmail(profile.Dir(e.paths, prof))
	if err != nil {
		return nil, fmt.Errorf("could not read the email after login: %w", err)
	}
	var saved *config.Profile
	_, err = e.update(func(cfg *config.Config) error {
		if saved, err = cfg.MustFind(name); err == nil {
			saved.Email = email
		}
		return err
	})
	return saved, err
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

// printShortcut prints the account's shortcut line under a ✓ line.
func printShortcut(e *env, prof *config.Profile) {
	if prof.Shortcut != "" {
		e.printf("  shortcut: %s\n", prof.Shortcut)
	} else {
		e.printf("  no shortcut (add one with `ccam shortcut %s %s`)\n", prof.Name, shim.Suggest(prof.Name))
	}
}

// nameState lists what ccam keeps under an account's name outside its folder.
func nameState(e *env, name string) []string {
	return append(launcher.StatePaths(e.paths, name), usage.CachePaths(usageStateDir(e), name)...)
}

func installShims(e *env, cfg *config.Config) ([]string, error) {
	self, err := ccamPath()
	if err != nil {
		return nil, err
	}
	created, err := shim.Install(e.paths.Bin, self, cfg)
	if err == nil && len(created) > 0 && !shim.InPath(e.paths.Bin) {
		fmt.Fprint(e.stderr, pathHint(e.paths.Bin))
	}
	return created, err
}

// pathHint tells how to put dir on PATH for good.
func pathHint(dir string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("Note: %s is not in PATH. Add it for your user in PowerShell, then open a new terminal:\n"+
			"  [Environment]::SetEnvironmentVariable(\"Path\", \"%s;\" + [Environment]::GetEnvironmentVariable(\"Path\", \"User\"), \"User\")\n", dir, dir)
	}
	return fmt.Sprintf("Note: %s is not in PATH. Add it to ~/.bashrc or ~/.zshrc:\n  export PATH=\"%s:$PATH\"\n", dir, dir)
}

type listRow struct {
	Slot     int               `json:"slot"`
	Shortcut string            `json:"shortcut"`
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Default  bool              `json:"default"`
	Dir      string            `json:"dir"`
	LastUsed *time.Time        `json:"last_used,omitempty"`
	API      string            `json:"api,omitempty"` // "API" or "API: <host>" for API key accounts
	Share    map[string]string `json:"share"`         // item → group, or "local"
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
		r := listRow{Slot: p.Slot, Shortcut: p.Shortcut, Name: p.Name, Email: p.Email, Default: def == p, Dir: profile.Dir(e.paths, p), Share: map[string]string{}}
		if p.API != nil {
			r.API = p.API.Label()
		}
		for _, item := range config.Items {
			r.Share[item] = p.Share.Get(item)
		}
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
	fmt.Fprintln(tw, " \t#\tName\tShortcut\tEmail\tShared\tLast used")
	for i, r := range rows {
		mark := ""
		if r.Default {
			mark = "*"
		}
		email := r.Email
		if r.API != "" {
			email = "(" + r.API + ")"
		} else if email == "" {
			email = "(not logged in)"
		}
		last := "-"
		if r.LastUsed != nil {
			last = Ago(time.Since(*r.LastUsed))
		}
		shortcut := r.Shortcut
		if shortcut == "" {
			shortcut = "-"
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\t%s\t%s\n", mark, r.Slot, r.Name, shortcut, email, cmp.Or(cfg.Profiles[i].Share.Summary(), "-"), last)
	}
	return tw.Flush()
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
		Use:   "use <name|number>",
		Short: "Set the default account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, p, err := e.loadProfile(args[0])
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
		Use:     "remove <name|number>",
		Aliases: []string{"rm"},
		Short:   "Remove an account (its folder goes to the trash)",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, p, err := e.loadProfile(args[0])
			if err != nil {
				return err
			}
			if err := refuseRunning(e, p.Name); err != nil {
				return err
			}
			if !yes && !e.confirm(fmt.Sprintf("Remove account %s (%s)?", p.Name, profile.Dir(e.paths, p))) {
				e.printf("Cancelled.\n")
				return nil
			}
			// Reload: the prompt may have waited while another ccam saved.
			cfg, p, err := e.loadProfile(p.Name)
			if err != nil {
				return err
			}
			prof := *p
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
			for _, path := range nameState(e, prof.Name) {
				os.RemoveAll(path) // a later account of this name must not inherit it
			}
			if self, err := ccamPath(); err == nil {
				shim.Remove(e.paths.Bin, self, prof.Shortcut)
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
		Short: "Rename an account (its shortcut command stays the same)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, old, err := e.loadProfile(args[0])
			if err != nil {
				return err
			}
			oldName := old.Name
			if err := refuseRunning(e, oldName); err != nil {
				return err
			}
			if err := profile.Rename(e.paths, cfg, oldName, args[1]); err != nil {
				return err
			}
			if err := e.save(cfg); err != nil {
				return err
			}
			newState := nameState(e, args[1])
			for i, path := range nameState(e, oldName) {
				_ = os.Rename(path, newState[i]) // best effort: last use and cached usage
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
		Use:   "login <name|number>",
		Short: "Log in again",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := loginProfile(e, args[0])
			if err != nil {
				return err
			}
			e.printf("✓ %s logged in: %s\n", p.Name, p.Email)
			return nil
		},
	}
}

func logoutCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "logout <name|number>",
		Short: "Remove the login token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, p, err := e.loadProfile(args[0])
			if err != nil {
				return err
			}
			if p.API != nil {
				return fmt.Errorf("%s is an API key account: it has no login to remove (its key comes from %s)", p.Name, p.API.KeyFrom)
			}
			if err := runClaude(e, cfg, p, "auth", "logout"); err != nil {
				return err
			}
			name := p.Name
			_, err = e.update(func(cfg *config.Config) error {
				p, err := cfg.MustFind(name)
				if err == nil {
					p.Email = ""
				}
				return err
			})
			return err
		},
	}
}

func shimsCmd(e *env) *cobra.Command {
	cmd := &cobra.Command{Use: "shims", Short: "Manage shortcut links (set them with `ccam shortcut`)"}
	install := &cobra.Command{
		Use:   "install",
		Short: "Recreate the shortcut links of all accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			created, err := installShims(e, cfg)
			if err != nil {
				return err
			}
			e.printf("✓ %d new shortcuts created (%s)\n", len(created), e.paths.Bin)
			return nil
		},
	}
	cmd.AddCommand(install)
	return cmd
}

func runCmd(e *env) *cobra.Command {
	// Only reached for `ccam help run`; real runs never get to cobra.
	return &cobra.Command{
		Use:   "run [--dry-run] [--fallback] [name|number] [claude args]",
		Short: "Run claude as that account (same as `ccam <name>`)",
		Long: "Runs Claude Code as the account. Without a name: the account set for this folder\n" +
			"(ccam folder), else the default one (ccam use). Everything after the name goes to\n" +
			"claude unchanged.\n\n" +
			"  --dry-run   print the environment and command instead of running it\n" +
			"  --fallback  for -p/--print runs: when the account reaches its usage limit, run\n" +
			"              the same command again as the logged-in account with the most room\n" +
			"              (output is held until a run finishes)",
		Example:            "  ccam run work --resume\n  ccam run --fallback work -p \"fix the tests\"\n  ccam --fallback -p \"summarise\" < notes.txt",
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
