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
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/migrate"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/shim"
	"github.com/jhyunwoo/ccam/internal/tui"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func addCmd(e *env) *cobra.Command {
	var noLogin, noGroup bool
	var shares, groups []string
	var shortcut string
	var api config.API
	var useAPI bool
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add an account: create profile → log in → optional shortcut",
		Long: "Run on a terminal without options, a wizard asks for the name, the groups\n" +
			"and whether to create a shortcut command (and its name).\n" +
			"Otherwise no shortcut is created unless --shortcut names one.\n\n" +
			"The account joins group default (created to share MCP servers, settings,\n" +
			"CLAUDE.md and extensions if there is none), the groups given with --group,\n" +
			"or none with --no-group. See `ccam group`.\n\n" +
			"--api adds an account that uses an API key instead of a Claude login: Anthropic's\n" +
			"API, or a compatible endpoint (--base-url) such as LiteLLM or OpenRouter. The key\n" +
			"is not stored; --key-from says where ccam reads it at each launch:\n" +
			"  env:NAME          an environment variable\n" +
			"  command:COMMAND   the first line a command prints (a password manager, a keychain)",
		Example: "  ccam add\n  ccam add work --shortcut ccwork\n  ccam add side --no-login --group default --group team\n" +
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
			if len(shares) > 0 {
				return fmt.Errorf("--share was replaced by groups: use --group <group> (repeatable) or --no-group; see `ccam group`")
			}
			if noGroup && len(groups) > 0 {
				return fmt.Errorf("--no-group and --group cannot be used together")
			}
			if interactive(e) && !useAPI && !cmd.Flags().Changed("group") && !noGroup && !cmd.Flags().Changed("shortcut") && !noLogin {
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
			acct := newAccount{name: args[0], groups: parseList(groups), login: !noLogin && !useAPI, shortcut: shortcut}
			acct.defaults = !noGroup && len(acct.groups) == 0
			if useAPI {
				acct.api = &api
				if _, err := ReadKeyForCheck(&api); err != nil {
					fmt.Fprintf(e.stderr, "Note: %v. Set it before starting the account.\n", err)
				}
			}
			return addProfile(e, acct)
		},
	}
	cmd.Flags().BoolVar(&useAPI, "api", false, "an API key account instead of a Claude login")
	cmd.Flags().StringVar(&api.BaseURL, "base-url", "", "with --api: the endpoint (default: Anthropic's API)")
	cmd.Flags().StringVar(&api.KeyFrom, "key-from", "", "with --api: env:NAME or command:COMMAND")
	cmd.Flags().StringVar(&api.Auth, "auth", "", "with --api: api-key (x-api-key, default) or bearer")
	cmd.Flags().StringVar(&api.Model, "model", "", "with --api: the model to use (ANTHROPIC_MODEL)")
	cmd.Flags().BoolVar(&noLogin, "no-login", false, "log in later with ccam login <name>")
	cmd.Flags().StringArrayVar(&groups, "group", nil, "join this group (repeatable; default: group default)")
	cmd.Flags().BoolVar(&noGroup, "no-group", false, "join no group: share nothing")
	cmd.Flags().StringArrayVar(&shares, "share", nil, "replaced by --group")
	_ = cmd.Flags().MarkHidden("share")
	cmd.Flags().StringVar(&shortcut, "shortcut", "", "also create this shortcut command (e.g. ccwork); none without it")
	return cmd
}

// checkShortcut validates a shortcut name for account prof.
func checkShortcut(e *env, cfg *config.Config, prof, name string) error {
	self, _ := ccamPath()
	return shim.Check(name, e.paths.Bin, self, cfg, prof)
}

// ReadKeyForCheck reads an env: key to tell early that it is missing;
// command: sources are not run outside a launch.
func ReadKeyForCheck(api *config.API) (string, error) {
	if !strings.HasPrefix(api.KeyFrom, "env:") {
		return "", nil
	}
	return launcher.ReadAPIKey(api.KeyFrom)
}

// newAccount is what `ccam add` and its wizard create.
type newAccount struct {
	name     string
	groups   []string // groups to join; default is created when missing
	defaults bool     // join cfg.DefaultGroups instead
	items    []string // what a new group default shares (config.DefaultItems if nil)
	login    bool
	shortcut string
	api      *config.API
}

// joinNew puts a new account in groups, creating a missing one with items.
func joinNew(cfg *config.Config, name string, groups, items []string) error {
	if items == nil {
		items = config.DefaultItems
	}
	for _, g := range groups {
		if _, ok := cfg.FindGroup(g); !ok {
			if err := cfg.AddGroup(g, items, nil); err != nil {
				return err
			}
		}
		if _, err := cfg.JoinGroup(name, g, false); err != nil {
			return err
		}
	}
	return nil
}

func addProfile(e *env, a newAccount) error {
	name, shortcut := a.name, a.shortcut
	cfg, err := e.load()
	if err != nil {
		return err
	}
	if a.defaults {
		a.groups = cfg.DefaultGroups()
		if a.groups == nil {
			e.printf("Note: there is no group default, so %s shares nothing; `ccam group join <group> %s` adds it to one.\n", name, name)
		}
	}
	if shortcut != "" {
		if err := checkShortcut(e, cfg, name, shortcut); err != nil {
			return err
		}
	}
	// Check the groups before creating the folder.
	for _, g := range a.groups {
		if _, ok := cfg.FindGroup(g); !ok && g != "default" {
			return fmt.Errorf("no group %q; create it with `ccam group add %s --items …`, or see `ccam group`", g, g)
		}
	}
	prof, err := profile.Add(e.paths, cfg, name)
	if err != nil {
		return err
	}
	prof.Shortcut = shortcut
	prof.API = a.api
	if err := joinNew(cfg, name, a.groups, a.items); err != nil {
		os.RemoveAll(e.paths.ProfileDir(name))
		return err
	}
	if prof, err = cfg.MustFind(name); err != nil {
		return err
	}
	if err := applyShare(e, cfg, prof); err != nil {
		// The folder is new and unregistered; left behind, it would block
		// every retry with "folder already exists".
		os.RemoveAll(e.paths.ProfileDir(name))
		return err
	}
	if err := e.save(cfg); err != nil {
		return err
	}
	if a.login {
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
	name, dir := prof.Name, profile.Dir(e.paths, prof)
	if err := runLogin(e, cfg, prof); err != nil {
		return nil, err
	}
	email, err := profile.ReadEmail(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read the email after login: %w", err)
	}
	id, _ := profile.ReadIdentity(dir)
	// Claude Code takes the email only as a hint for the login page, so a
	// login can end up in another account than the one this was.
	if was := prof.Login; was != nil && id != nil && was.AccountUUID != id.AccountUUID {
		before := cmp.Or(was.Email, "another account")
		fmt.Fprintf(e.stderr, "ccam: this login is %s, but %s was %s before.\n", cmp.Or(email, "another account"), name, before)
		fmt.Fprintf(e.stderr, "To undo it, run `ccam logout %s`, then `ccam login %s` and sign in as %s.\n", name, name, before)
	}
	if id != nil {
		for _, q := range cfg.Profiles {
			if q.Name != name && config.SameAccount(q.Login, id) {
				fmt.Fprintf(e.stderr, "ccam: %s is logged in to the same account (%s); the two share their usage limits.\n", q.Name, email)
			}
		}
	}
	var saved *config.Profile
	_, err = e.update(func(cfg *config.Config) error {
		if saved, err = cfg.MustFind(name); err == nil {
			saved.Email, saved.Login = email, id
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
	out := append(launcher.StatePaths(e.paths, name), usage.CachePaths(usageStateDir(e), name)...)
	return append(out, handoff.Paths(e.paths.State, name)...)
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
	Groups   []string          `json:"groups"`
	Share    map[string]string `json:"share"`          // item → group, or "local"
	Plan     string            `json:"plan,omitempty"` // e.g. Claude Team
	Org      string            `json:"org,omitempty"`  // the organization's name
	pool     string
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
	org := newOrgs(e, cfg, false)
	rows := make([]listRow, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		r := listRow{Slot: p.Slot, Shortcut: p.Shortcut, Name: p.Name, Email: p.Email, Default: def == p, Dir: profile.Dir(e.paths, p), Groups: []string{}, Share: map[string]string{}}
		if p.API != nil {
			r.API = p.API.Label()
		}
		for _, g := range cfg.GroupsOf(p.Name) {
			r.Groups = append(r.Groups, g.Name)
		}
		for _, item := range config.Items {
			r.Share[item] = config.Local
			if g := cfg.GroupFor(p.Name, item); g != nil {
				r.Share[item] = g.Name
			}
		}
		if id := org.of(p); id != nil && p.API == nil {
			r.Plan, r.Org, r.pool = id.PlanName(), id.OrgName, id.Pool()
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
	// The organization matters only when the accounts are in more than one.
	pools := map[string]bool{}
	for _, r := range rows {
		if r.pool != "" {
			pools[r.pool] = true
		}
	}
	showOrg := len(pools) > 1
	tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
	if showOrg {
		fmt.Fprintln(tw, " \t#\tName\tShortcut\tEmail\tOrganization\tGroups\tLast used")
	} else {
		fmt.Fprintln(tw, " \t#\tName\tShortcut\tEmail\tGroups\tLast used")
	}
	for _, r := range rows {
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
		if showOrg {
			o := r.Org
			if r.pool == config.Personal {
				o = "personal"
			}
			email += "\t" + cmp.Or(o, "-")
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\t%s\t%s\n", mark, r.Slot, r.Name, shortcut, email, cmp.Or(strings.Join(r.Groups, " "), "-"), last)
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

// removeAccount moves an account's folder to the trash (an imported one
// stays in place) and forgets the account. The caller has checked that it
// is not running.
func removeAccount(e *env, name string) (string, error) {
	// Reload: a prompt may have waited while another ccam saved.
	cfg, p, err := e.loadProfile(name)
	if err != nil {
		return "", err
	}
	prof := *p
	if g := prof.Share.Get("extensions"); g != config.Local {
		// Other accounts' plugin paths may run through this folder.
		sep := string(filepath.Separator)
		groupPlugins := filepath.Join(e.paths.GroupDir(g), "plugins")
		if _, err := profile.RewritePluginPaths(groupPlugins,
			filepath.Join(profile.Dir(e.paths, &prof), "plugins")+sep, groupPlugins+sep); err != nil {
			return "", err
		}
	}
	trashed, err := profile.Trash(e.paths, &prof)
	if err != nil {
		return "", err
	}
	cfg.Remove(prof.Name)
	if err := e.save(cfg); err != nil {
		return "", err
	}
	for _, path := range nameState(e, prof.Name) {
		os.RemoveAll(path) // a later account of this name must not inherit it
	}
	if self, err := ccamPath(); err == nil {
		shim.Remove(e.paths.Bin, self, prof.Shortcut)
	}
	if trashed != "" {
		return fmt.Sprintf("%s removed (folder moved to %s)", prof.Name, trashed), nil
	}
	return fmt.Sprintf("%s unregistered (folder %s left in place)", prof.Name, profile.Dir(e.paths, &prof)), nil
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
			msg, err := removeAccount(e, p.Name)
			if err != nil {
				return err
			}
			e.printf("✓ %s\n", msg)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "remove without asking")
	return cmd
}

// renameAccount renames an account, its folder and what ccam keeps under
// its name, and returns the old name.
func renameAccount(e *env, sel, newName string) (string, error) {
	cfg, old, err := e.loadProfile(sel)
	if err != nil {
		return "", err
	}
	oldName := old.Name
	if err := refuseRunning(e, oldName); err != nil {
		return "", err
	}
	if err := profile.Rename(e.paths, cfg, oldName, newName); err != nil {
		return "", err
	}
	if err := e.save(cfg); err != nil {
		return "", err
	}
	newState := nameState(e, newName)
	for i, path := range nameState(e, oldName) {
		_ = os.Rename(path, newState[i]) // best effort: last use and cached usage
	}
	if _, err := installShims(e, cfg); err != nil {
		fmt.Fprintln(e.stderr, "Could not update shortcuts:", err)
	}
	return oldName, nil
}

func renameCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "rename <name> <new-name>",
		Short: "Rename an account (its shortcut command stays the same)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName, err := renameAccount(e, args[0], args[1])
			if err != nil {
				return err
			}
			e.printf("✓ %s → %s\n", oldName, args[1])
			return nil
		},
	}
}

func loginCmd(e *env) *cobra.Command {
	var stale bool
	cmd := &cobra.Command{
		Use:   "login <name|number> | login --expired",
		Short: "Log in again (--expired: every account whose login no longer works)",
		Example: "  ccam login work\n" +
			"  ccam login --expired    # one after another: never logged in, or the login expired",
		Args: func(cmd *cobra.Command, args []string) error {
			if stale {
				return cobra.NoArgs(cmd, args)
			}
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if stale {
				return loginStale(e)
			}
			p, err := loginProfile(e, args[0])
			if err != nil {
				return err
			}
			e.printf("✓ %s logged in: %s\n", p.Name, p.Email)
			return nil
		},
	}
	cmd.Flags().BoolVar(&stale, "expired", false, "log in every account that was never logged in or whose login expired")
	return cmd
}

// loginStale logs in, one after another, every account whose login no
// longer works. API key accounts have no login and are left out.
func loginStale(e *env) error {
	cfg, err := e.load()
	if err != nil {
		return err
	}
	type todo struct{ name, why string }
	var list []todo
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		if p.API != nil {
			continue
		}
		if why := creds.Stale(profile.Dir(e.paths, p), e.paths.Home, time.Now()); why != "" {
			list = append(list, todo{p.Name, why})
		}
	}
	if len(list) == 0 {
		e.printf("Every account's login works.\n")
		return nil
	}
	failed := 0
	for i, t := range list {
		e.printf("[%d/%d] %s (%s): log in in your browser.\n", i+1, len(list), t.name, t.why)
		p, err := loginProfile(e, t.name)
		if err != nil {
			fmt.Fprintf(e.stderr, "ccam: %s: %v\n", t.name, err)
			failed++
			continue
		}
		e.printf("✓ %s logged in: %s\n", p.Name, p.Email)
	}
	if failed > 0 {
		return exitError(1)
	}
	return nil
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
					p.Email = "" // the login record stays, to check the next login against
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
