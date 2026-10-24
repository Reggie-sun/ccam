package cli

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/migrate"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
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
	return shim.Check(name, e.paths, self, cfg, prof)
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

// canLinkFiles is fsutil.CanLinkFiles; tests replace it.
var canLinkFiles = fsutil.CanLinkFiles

// fileItems lists the items of items that link files (settings, memory,
// history), which Windows does only in Developer Mode.
func fileItems(items []string) []string {
	return slices.DeleteFunc(slices.Clone(items), func(item string) bool {
		return !slices.ContainsFunc(share.Layout[item], func(e share.Entry) bool { return !e.Dir })
	})
}

// checkFileLinks runs before a new account's folder exists, as the group
// engine's Check does for a change of groups. Where files cannot be linked,
// the group default that `ccam add` creates leaves out the items that need
// file links, and a group that already shares one is refused.
func checkFileLinks(e *env, cfg *config.Config, a *newAccount) error {
	items := a.items
	if items == nil {
		items = config.DefaultItems
	}
	var need []string
	for _, g := range a.groups {
		have := items
		if grp, ok := cfg.FindGroup(g); ok {
			have = grp.Items
		}
		if len(fileItems(have)) > 0 {
			need = append(need, g)
		}
	}
	if len(need) == 0 {
		return nil
	}
	root := e.paths.GroupsRoot()
	if err := os.MkdirAll(root, 0o700); err != nil {
		return err
	}
	err := canLinkFiles(root)
	if err == nil {
		return nil
	}
	for _, g := range need {
		if grp, ok := cfg.FindGroup(g); ok {
			return fmt.Errorf("%v, and group %s shares %s. Turn it on, or add the account to a group that shares only mcp and extensions", err, g, strings.Join(fileItems(grp.Items), ", "))
		}
	}
	dropped := fileItems(items)
	a.items = slices.DeleteFunc(slices.Clone(items), func(item string) bool { return slices.Contains(dropped, item) })
	shares, rest := "only "+strings.Join(a.items, ", "), "the rest"
	if len(a.items) == 0 {
		shares, rest = "nothing yet", "them"
	}
	e.printf("Note: %v, so group default shares %s. Once it is on, `ccam group set default --add %s` shares %s.\n",
		err, shares, strings.Join(dropped, ","), rest)
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
	if err := checkFileLinks(e, cfg, &a); err != nil {
		return err
	}
	if dir := e.paths.ProfileDir(name); !fsutil.Exists(dir) && keychainLeftover(dir, e.paths.Home) {
		svc := creds.KeychainService(dir, e.paths.Home)
		return fmt.Errorf("the macOS Keychain still holds the login of an earlier account named %s, which the new account would be logged in as. "+
			"Remove it with `security delete-generic-password -a %s -s '%s'`, or choose another name", name, creds.KeychainAccount(), svc)
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
		os.RemoveAll(e.paths.ProfileDir(name)) // as above
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
		fmt.Fprintln(e.stderr, "Could not create shortcuts:", err)
	}
	prof, err = cfg.MustFind(name)
	if err != nil {
		return err
	}
	e.printf("✓ %s (%s) added\n", name, accountLabel(prof))
	printShortcut(e, prof)
	if claudeDir := filepath.Join(e.paths.Home, ".claude"); len(cfg.Profiles) == 1 {
		// Claude Code reads an account's own .claude.json, so the MCP
		// servers plain claude has in ~/.claude.json are not there.
		from, servers := migrate.HomeOnly(claudeDir, e.paths.Home, nil, nil)
		switch {
		case migrate.HasContent(claudeDir) && len(servers) > 0:
			e.printf("Tip: run `ccam migrate` to copy your existing Claude Code settings (%s) and MCP servers (%s) into the shared group.\n", claudeDir, strings.Join(servers, ", "))
		case migrate.HasContent(claudeDir):
			e.printf("Tip: run `ccam migrate` to copy your existing Claude Code settings (%s) into the shared settings.\n", claudeDir)
		case len(servers) > 0:
			e.printf("Tip: the MCP servers in %s (%s) are not in ccam's accounts; `ccam migrate --mcp` copies them into the shared group.\n", from, strings.Join(servers, ", "))
		}
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
	// Before the browser: the login's record would not be saved.
	if err := cfg.Writable(e.paths.ConfigFile()); err != nil {
		return nil, err
	}
	name, dir := prof.Name, profile.Dir(e.paths, prof)
	before := prof.Login
	if before == nil {
		before, _ = profile.ReadIdentity(dir)
	}
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
	was := ""
	switch {
	case prof.Login != nil && id != nil && prof.Login.AccountUUID != id.AccountUUID:
		was = cmp.Or(prof.Login.Email, "another account")
	case prof.Login == nil && prof.Email != "" && email != "" && !strings.EqualFold(prof.Email, email):
		was = prof.Email // an older ccam, or an import, recorded no login: only the email
	}
	if was != "" {
		fmt.Fprintf(e.stderr, "ccam: this login is %s, but %s was %s before.\n", cmp.Or(email, "another account"), name, was)
		// A login replaces the old one and keeps the MCP sign-ins; a logout
		// would delete those.
		fmt.Fprintf(e.stderr, "To undo it, run `ccam login %s` again and sign in as %s.\n", name, was)
	}
	// A usage limit Claude Code reported was the account's logged in
	// before; it holds only when the login is that account again.
	if before == nil || id == nil || before.AccountUUID != id.AccountUUID {
		usage.DropLimit(usageStateDir(e), name)
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
	c, err := launcher.ClaudeCommand(context.Background(), bin, args...)
	if err != nil {
		return err
	}
	c.Env = launcher.ProfileEnv(e.paths, os.Environ(), profile.Dir(e.paths, prof), prof.Name)
	c.Stdin, c.Stdout, c.Stderr = e.stdin, e.stdout, e.stderr
	return c.Run()
}

// printShortcut prints the account's shortcut line under a ✓ line.
func printShortcut(e *env, prof *config.Profile) {
	self, _ := ccamPath()
	switch {
	case prof.Shortcut == "":
		e.printf("  no shortcut (add one with `ccam shortcut %s %s`)\n", prof.Name, shim.Suggest(prof.Name))
	case shim.Linked(e.paths.Bin, self, prof.Shortcut):
		e.printf("  shortcut: %s\n", prof.Shortcut)
	default:
		if bin, ok := shim.Other(e.paths, self, prof.Shortcut); ok {
			e.printf("  shortcut: %s (runs another ccam, %s; `ccam doctor --fix` points it at this one)\n", prof.Shortcut, bin)
			return
		}
		e.printf("  shortcut: %s (not created; see above)\n", prof.Shortcut)
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
	// This also renews hard-linked shortcuts (Windows) that still run the
	// binary a new install or an update replaced.
	created, err := shim.Install(e.paths, self, cfg)
	if len(created) > 0 && !shim.InPath(e.paths.Bin) {
		fmt.Fprint(e.stderr, pathHint(e.paths.Bin))
	}
	return created, err
}

// pathHint tells how to put dir on PATH for good.
func pathHint(dir string) string {
	if runtime.GOOS == "windows" {
		// The user Path is REG_EXPAND_SZ: read and written as one, its
		// %VARIABLE% entries stay; the last call tells Explorer it changed.
		return fmt.Sprintf("Note: %s is not in PATH. Add it for your user in PowerShell, then open a new terminal:\n"+
			"  $k = 'HKCU:\\Environment'; Set-ItemProperty $k Path (%s + (Get-Item $k).GetValue('Path', '', 'DoNotExpandEnvironmentNames')) -Type ExpandString; "+
			"[Environment]::SetEnvironmentVariable('CCAM_PATH_CHANGED', $null, 'User')\n", dir, launcher.PSString(dir+";"))
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
	was      string // the email ccam recorded, when the folder is logged in as another now
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
			// Whom the folder is logged in as now: a /login inside a
			// session, or Claude Code's own first login, changes it.
			if org.live[p.Name] && id.Email != "" {
				if p.Email != "" && !strings.EqualFold(p.Email, id.Email) {
					r.was = p.Email
				}
				r.Email = id.Email
			}
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
		switch {
		case r.API != "":
			email = "(" + r.API + ")"
		case email == "":
			email = "(not logged in)"
		case r.was != "":
			email += " (was " + r.was + ")"
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
// stays in place) and forgets the account, unless it is running.
func removeAccount(e *env, name string) (string, error) {
	// Reload and check again: a prompt may have waited while another ccam
	// saved, or while the account started.
	cfg, p, err := e.loadProfile(name)
	if err != nil {
		return "", err
	}
	if err := group.Unfinished(e.paths); err != nil {
		// The note names the accounts: one of a new account of this name
		// would be finished on it.
		return "", err
	}
	if err := refuseRunning(e, p); err != nil {
		return "", err
	}
	// A logout cannot be undone: the removal must be able to finish.
	if err := canSave(e, cfg); err != nil {
		return "", err
	}
	unlock, err := holdRefresh(e, p)
	if err != nil {
		return "", err
	}
	defer unlock()
	left := ""
	if p.Dir != "" && !p.Share.Equal(config.Share{}) {
		// A folder registered in place stays, and plain claude goes on
		// using it: it must not stay linked into the groups, whose folders
		// a later `ccam group remove` moves to the trash. It leaves them
		// as `ccam group leave` does, with a copy of what they share. If
		// that stops halfway, it stays registered for doctor --fix.
		if err := leaveGroups(e, p.Name); err != nil {
			return "", fmt.Errorf("%s could not leave its groups, so it is still registered: %w", p.Name, err)
		}
		if cfg, p, err = e.loadProfile(p.Name); err != nil {
			return "", err
		}
		left = "; it left its groups with a copy of what they share"
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
	// A managed folder goes to the trash: its login is revoked first, so no
	// working token waits there (or, on macOS, in a Keychain item named
	// after the folder, which a later account of this name would inherit).
	// A copy of the login in another account's folder would be revoked as
	// well, so then the login stays.
	loggedOut, revoked := "", false
	if prof.Dir == "" && prof.API == nil && storedLogin(e, &prof) {
		if q := sameLogin(e, cfg, &prof); q != "" {
			loggedOut = fmt.Sprintf("; it was not logged out, since %s holds the same login", q)
		} else if err := logOut(e, cfg, &prof); err != nil {
			loggedOut = fmt.Sprintf("; it could not be logged out first (%v), so its login was not revoked", err)
		} else {
			loggedOut, revoked = "; logged out first", true
		}
	}
	trashed, err := profile.Trash(e.paths, &prof)
	if err != nil {
		if revoked {
			err = fmt.Errorf("%w\n%s is still registered, but logged out now; `ccam login %s` logs it in again", err, prof.Name, prof.Name)
		}
		return "", err
	}
	cfg.Remove(prof.Name)
	if err := e.save(cfg); err != nil {
		// config.toml still lists the account: its folder comes back.
		if trashed != "" {
			if uerr := fsutil.Move(trashed, profile.Dir(e.paths, &prof)); uerr != nil {
				err = errors.Join(err, uerr)
			}
		}
		if revoked {
			err = fmt.Errorf("%w\n%s is logged out now; `ccam login %s` logs it in again", err, prof.Name, prof.Name)
		}
		return "", err
	}
	stuck := "" // backups that could not be moved, and stay
	if prof.Dir != "" {
		// The folder stays in use: so do the backups of its .claude.json,
		// which leaving the groups may just have made.
		backups := filepath.Join(launcher.MCPStateDir(e.paths, prof.Name), "backups")
		if fsutil.Exists(backups) {
			kept := fsutil.UniqueDir(filepath.Join(profile.Dir(e.paths, &prof), share.BackupDir, "mcp-"+fsutil.Stamp(time.Now())))
			if err := fsutil.Move(backups, kept); err != nil {
				stuck = backups
				left += fmt.Sprintf("; the backups of its .claude.json could not be moved into its folder (%v), so they stay in %s (a later account named %s would find them there)", err, backups, prof.Name)
			} else {
				left += "; the backups of its .claude.json are in " + kept
			}
		}
	}
	for _, path := range nameState(e, prof.Name) {
		if stuck != "" && path == filepath.Dir(stuck) {
			removeAllBut(path, filepath.Base(stuck))
			continue
		}
		os.RemoveAll(path) // a later account of this name must not inherit it
	}
	if self, err := ccamPath(); err == nil {
		shim.Remove(e.paths, self, prof.Shortcut)
	}
	if trashed != "" {
		return fmt.Sprintf("%s removed (folder moved to %s%s)", prof.Name, trashed, loggedOut), nil
	}
	return fmt.Sprintf("%s unregistered (folder %s left in place%s)", prof.Name, profile.Dir(e.paths, &prof), left), nil
}

// removeAllBut removes what is in dir except the entry named keep.
func removeAllBut(dir, keep string) {
	entries, _ := os.ReadDir(dir)
	for _, en := range entries {
		if en.Name() != keep {
			os.RemoveAll(filepath.Join(dir, en.Name()))
		}
	}
}

// leaveGroups takes an account out of every group, as `ccam group leave`
// does without --fresh: what it shared becomes a copy of the group's in its
// folder, and the group's MCP servers go into its .claude.json.
func leaveGroups(e *env, name string) error {
	_, _, err := changeGroups(e, func(c *config.Config) error {
		c.DropMember(name)
		return c.Resolve()
	}, changeFlags{yes: true})
	return err
}

func removeCmd(e *env) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:     "remove <name|number>",
		Aliases: []string{"rm"},
		Short:   "Remove an account (it is logged out; its folder goes to the trash)",
		Long: "The account is logged out first (not when another account's folder holds the\n" +
			"same login), then its folder goes to ccam's trash.\n\n" +
			"An account registered in place (ccam import) keeps its folder and its login, for\n" +
			"plain claude: it leaves its groups first and keeps a copy of what they share, as\n" +
			"`ccam group leave` does, then ccam forgets it.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, p, err := e.loadProfile(args[0])
			if err != nil {
				return err
			}
			if err := refuseRunning(e, p); err != nil {
				return err
			}
			q := fmt.Sprintf("Remove account %s (%s)?", p.Name, profile.Dir(e.paths, p))
			switch copied := sameLogin(e, cfg, p); {
			case p.Dir == "" && p.API == nil && copied == "":
				q += " It is logged out first; its folder goes to ccam's trash."
			case p.Dir == "":
				q += " Its folder goes to ccam's trash."
			case !p.Share.Equal(config.Share{}):
				q += " It leaves its groups first, keeping a copy of what they share; its folder stays where it is."
			default:
				q += " Its folder stays where it is."
			}
			if !yes && !e.confirm(q) {
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
// its name, and returns the old name and a note when the folder kept its
// name: where Claude Code names the login after the folder's path (the
// macOS Keychain), a moved folder would leave the login behind. Refused
// while a change of groups is unfinished: its note names the accounts.
func renameAccount(e *env, sel, newName string) (string, string, error) {
	cfg, old, err := e.loadProfile(sel)
	if err != nil {
		return "", "", err
	}
	if err := group.Unfinished(e.paths); err != nil {
		return "", "", err
	}
	oldName, oldDir := old.Name, profile.Dir(e.paths, old)
	if err := refuseRunning(e, old); err != nil {
		return "", "", err
	}
	unlock, err := holdRefresh(e, old)
	if err != nil {
		return "", "", err
	}
	defer unlock()
	note := ""
	if store := namedByFolder(oldDir); old.Dir == "" && store != "" {
		note = fmt.Sprintf("its folder stays %s: %s keeps the login under the folder's path", oldDir, store)
	}
	moved := old.Dir == "" && note == "" && fsutil.Exists(oldDir)
	if err := profile.Rename(e.paths, cfg, oldName, newName, note != ""); err != nil {
		return "", "", err
	}
	if err := e.save(cfg); err != nil {
		// config.toml still names the old account: its folder goes back.
		if moved {
			if uerr := profile.MoveFolder(e.paths.ProfileDir(newName), oldDir); uerr != nil {
				return "", "", errors.Join(err, uerr)
			}
		}
		return "", "", err
	}
	newState := nameState(e, newName)
	for i, path := range nameState(e, oldName) {
		_ = os.Rename(path, newState[i]) // best effort: last use and cached usage
	}
	// Notes name the account and run through its folder.
	if p, ok := cfg.Find(newName); ok {
		handoff.Rename(e.paths.State, oldName, newName, oldDir, profile.Dir(e.paths, p))
	}
	if _, err := installShims(e, cfg); err != nil {
		fmt.Fprintln(e.stderr, "Could not update shortcuts:", err)
	}
	return oldName, note, nil
}

func renameCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "rename <name> <new-name>",
		Short: "Rename an account (its shortcut command stays the same)",
		Long: "Renames the account, its folder and what ccam keeps under its name. On macOS the folder\n" +
			"keeps its name: the Keychain keeps the login (and the MCP sign-ins) under the folder's path,\n" +
			"so a moved folder would be logged out. The same holds on Windows where Claude Code keeps\n" +
			"logins in Credential Manager.",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName, note, err := renameAccount(e, args[0], args[1])
			if err != nil {
				return err
			}
			e.printf("✓ %s → %s\n", oldName, args[1])
			if note != "" {
				e.printf("  %s\n", note)
			}
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
	if err := cfg.Writable(e.paths.ConfigFile()); err != nil {
		return err
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
	var yes bool
	cmd := &cobra.Command{
		Use:   "logout <name|number>",
		Short: "Log out: remove the account's login and its MCP server sign-ins",
		Long: "Runs `claude auth logout` in the account: Claude Code revokes the login and deletes what it\n" +
			"stores for the account, its MCP server sign-ins too (sign in to those servers again with\n" +
			"/mcp after the next login). ccam asks first when that loses sign-ins, or when the account\n" +
			"is running, since its sessions are signed out as well.\n\n" +
			"To log in as another Claude account, `ccam login <name>` alone is enough; it keeps the\n" +
			"MCP sign-ins.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, p, err := e.loadProfile(args[0])
			if err != nil {
				return err
			}
			if p.API != nil {
				return fmt.Errorf("%s is an API key account: it has no login to remove (its key comes from %s)", p.Name, p.API.KeyFrom)
			}
			if !yes && !confirmLogout(e, cfg, p) {
				e.printf("Cancelled.\n")
				return nil
			}
			unlock, err := holdRefresh(e, p)
			if err != nil {
				return err
			}
			defer unlock()
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
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "log out without asking")
	return cmd
}

// confirmLogout says what a logout of p loses besides the login (its MCP
// sign-ins, its running sessions' login, a copy of the login in another
// account) and asks; true when there is nothing to lose or the answer is
// yes.
func confirmLogout(e *env, cfg *config.Config, p *config.Profile) bool {
	dir := profile.Dir(e.paths, p)
	var lose []string
	if launcher.Running(e.paths, p.Name, dir) {
		lose = append(lose, fmt.Sprintf("%s is running: its Claude Code sessions are signed out too.", p.Name))
	}
	if q := sameLogin(e, cfg, p); q != "" {
		lose = append(lose, fmt.Sprintf("%s holds the same login (a copied folder), which the logout revokes too; `ccam login %s` gives %s a login of its own instead.", q, p.Name, p.Name))
	}
	if s := mcpSignIns(dir, e.paths.Home); len(s) > 0 {
		lose = append(lose, fmt.Sprintf("Logging out also removes %s's MCP server sign-ins (%s); sign in to them again with /mcp.", p.Name, strings.Join(s, ", ")))
	}
	if len(lose) == 0 {
		return true
	}
	for _, l := range lose {
		e.printf("%s\n", l)
	}
	return e.confirm(fmt.Sprintf("Log out %s?", p.Name))
}

// mcpSignIns names the MCP servers dir's login holds a sign-in for.
func mcpSignIns(dir, home string) []string {
	data, _, err := creds.Read(dir, home)
	if err != nil {
		return nil
	}
	var names []string
	for k := range creds.MCPSignIns(data) {
		name, _, _ := strings.Cut(k, "|")
		if !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	return names
}

// storedLogin reports whether anything of a login is stored for prof (a
// dead one too, or one kept where ccam cannot read it).
func storedLogin(e *env, prof *config.Profile) bool {
	_, _, err := creds.Read(profile.Dir(e.paths, prof), e.paths.Home)
	return !errors.Is(err, creds.ErrNone)
}

// sameLogin names another account whose folder holds prof's login (a
// copied folder: the same refresh token), or plain claude when ~/.claude,
// registered as no account, holds it; or "".
func sameLogin(e *env, cfg *config.Config, prof *config.Profile) string {
	key := creds.LoginKey(profile.Dir(e.paths, prof))
	if key == "" {
		return ""
	}
	for i := range cfg.Profiles {
		if q := &cfg.Profiles[i]; q.Name != prof.Name && q.API == nil && creds.LoginKey(profile.Dir(e.paths, q)) == key {
			return q.Name
		}
	}
	if dir := profile.PlainClaude(e.paths, cfg); dir != "" && creds.LoginKey(dir) == key {
		return profile.PlainClaudeName
	}
	return ""
}

// logOut runs `claude auth logout` in prof's folder, its output kept back
// unless it fails: Claude Code revokes the login and deletes what it stores
// for that folder (the login file, or its macOS Keychain or Windows
// Credential Manager entry), the MCP sign-ins with it.
func logOut(e *env, cfg *config.Config, prof *config.Profile) error {
	bin, err := launcher.ResolveClaude(e.paths, cfg.ClaudeBin)
	if err != nil {
		return err
	}
	c, err := launcher.ClaudeCommand(context.Background(), bin, "auth", "logout")
	if err != nil {
		return err
	}
	c.Env = launcher.ProfileEnv(e.paths, os.Environ(), profile.Dir(e.paths, prof), prof.Name)
	if out, err := c.CombinedOutput(); err != nil {
		if msg := strings.TrimSpace(string(out)); msg != "" {
			return fmt.Errorf("%v: %s", err, msg)
		}
		return err
	}
	return nil
}

// namedByFolder and keychainLeftover are replaced in tests, which cannot
// reach a Keychain.
var (
	namedByFolder = creds.NamedByFolder
	// keychainLeftover reports whether the macOS Keychain holds a login (or
	// MCP sign-ins) for a folder that does not exist yet: one an earlier
	// account of that name left behind, as older ccam versions did on
	// remove and rename. An item Claude Code's logout leaves with its device
	// key only is no login.
	keychainLeftover = func(dir, home string) bool {
		if runtime.GOOS != "darwin" {
			return false
		}
		data, src, err := creds.Read(dir, home)
		return err == nil && src == creds.Keychain && holdsLogin(data)
	}
)

// holdsLogin reports whether stored login data holds what a new account of
// the folder would inherit: a login, dead or not, or MCP sign-ins. The
// device key Claude Code's logout writes back is neither.
func holdsLogin(data []byte) bool {
	return creds.Judge(data, time.Now()) != "not logged in" || len(creds.MCPSignIns(data)) > 0
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
			if err == nil || len(created) > 0 {
				e.printf("✓ %d new shortcuts created (%s)\n", len(created), e.paths.Bin)
			}
			return err
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
