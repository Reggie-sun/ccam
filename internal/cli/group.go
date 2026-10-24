package cli

import (
	"cmp"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/share"
)

// changeFlags are the flags of every command that moves data between
// accounts and groups.
type changeFlags struct {
	yes    bool
	prefer string
	fresh  bool
	from   string
}

func (f *changeFlags) add(cmd *cobra.Command, seed bool) {
	cmd.Flags().BoolVarP(&f.yes, "yes", "y", false, "run without asking (conflicts follow --prefer)")
	cmd.Flags().StringVar(&f.prefer, "prefer", "", "conflict handling: group|mine|merge (default: ask; with --yes, each item's default)")
	cmd.Flags().BoolVar(&f.fresh, "fresh", false, "accounts that stop sharing an item start without it instead of with a copy (an imported folder still keeps one)")
	if seed {
		cmd.Flags().StringVar(&f.from, "from", "", "the account whose copy goes into the group first")
	}
}

// changeGroups applies edit to a copy of the config, shows what that moves
// on disk, asks, saves the config, then moves the data.
func changeGroups(e *env, edit func(*config.Config) error, f changeFlags) (*group.Change, bool, error) {
	prefer, err := parsePrefer(f.prefer)
	if err != nil {
		return nil, false, err
	}
	old, err := e.load()
	if err != nil {
		return nil, false, err
	}
	if err := group.Unfinished(e.paths); err != nil {
		return nil, false, err
	}
	next := old.Clone()
	if err := edit(next); err != nil {
		return nil, false, err
	}
	if f.from != "" {
		p, err := next.MustFind(f.from)
		if err != nil {
			return nil, false, fmt.Errorf("--from: %w", err)
		}
		f.from = p.Name
	}
	ch, err := group.Build(e.paths, old, next, group.Options{Seed: f.from, Fresh: f.fresh})
	if err != nil {
		return nil, false, err
	}
	if !e.dryRun {
		if err := ch.Check(group.Running(e.paths)); err != nil {
			return nil, false, err
		}
	}
	for _, c := range ch.Conflicts() {
		switch {
		case prefer != nil:
			for _, cand := range c.Candidates {
				if cand == *prefer || (*prefer == share.MergeJSON && cand == share.MergeLines) {
					c.Choice = cand
				}
			}
		case !f.yes:
			c.Choice = askConflict(e, c.Conflict, c.Account)
		}
	}
	if !ch.Empty() {
		e.printf("Plan (paths in %s):\n", e.paths.Data)
		ch.Describe(e.stdout)
	}
	if e.dryRun {
		return ch, false, nil
	}
	if !ch.Empty() && !f.yes && !e.confirm("Proceed?") {
		e.printf("Cancelled.\n")
		return ch, false, nil
	}
	// Save first, with a note of the change: if a move stops halfway,
	// doctor --fix finishes it. The prompts may have waited, so the edit
	// is made again on the current config, after checking again that
	// nothing started or changed the groups meanwhile.
	if _, err := e.update(func(c *config.Config) error { return recheck(ch, c, edit, group.Running(e.paths)) }); err != nil {
		ch.Abort()
		return nil, false, err
	}
	if err := ch.Apply(); err != nil {
		return nil, false, fmt.Errorf("failed while applying: %w", err)
	}
	return ch, true, nil
}

// recheck is Change.Recheck, then Begin, inside the config update.
func recheck(ch *group.Change, cur *config.Config, edit func(*config.Config) error, running func(name, dir string) bool) error {
	if err := ch.Recheck(cur, edit, running); err != nil {
		return err
	}
	return ch.Begin()
}

func groupCmd(e *env) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "group",
		Aliases: []string{"groups"},
		Short:   "Groups: which items a set of accounts share",
		Long: "A group shares items between its accounts: every account in the group\n" +
			"links those items to one folder, groups/<group>/. An account can be in\n" +
			"several groups, but each item comes from one group only.\n\n" +
			"Items: " + strings.Join(config.Items, ", ") + "\n" +
			"Login tokens are never shared.\n\n" +
			"Without a subcommand, lists the groups.",
		Example: "  ccam group\n" +
			"  ccam group add team --items history --members work,side\n" +
			"  ccam group set default --remove history\n" +
			"  ccam group join team personal\n" +
			"  ccam --dry-run group leave default side",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error { return listGroups(e) },
	}
	cmd.AddCommand(groupListCmd(e), groupShowCmd(e), groupAddCmd(e), groupSetCmd(e), groupJoinCmd(e),
		groupLeaveCmd(e), groupRenameCmd(e), groupRemoveCmd(e))
	return cmd
}

func groupListCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the groups, what they share and their accounts",
		Args:    cobra.NoArgs,
		RunE:    func(cmd *cobra.Command, args []string) error { return listGroups(e) },
	}
}

func listGroups(e *env) error {
	cfg, err := e.load()
	if err != nil {
		return err
	}
	if len(cfg.Groups) == 0 {
		e.printf("No groups. `ccam add` creates group default; or `ccam group add <name> --items …`.\n")
		return nil
	}
	tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Group\tShares\tAccounts")
	for _, g := range cfg.Groups {
		name := g.Name
		if g.Dir != "" {
			name += " (folder " + g.Dir + ")"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", name, cmp.Or(strings.Join(g.Items, " "), "nothing"), cmp.Or(strings.Join(accounts(cfg, &g), ", "), "-"))
	}
	return tw.Flush()
}

// accounts lists the members of g that are accounts.
func accounts(cfg *config.Config, g *config.Group) []string {
	var out []string
	for _, m := range g.Members {
		if cfg.IsAccount(m) {
			out = append(out, m)
		}
	}
	return out
}

func groupShowCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:               "show <group>",
		Short:             "Show a group: items, accounts, folder and MCP servers",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: groupNames(e),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			g, err := cfg.MustFindGroup(args[0])
			if err != nil {
				return err
			}
			e.printf("Group %s\n", g.Name)
			if len(g.Items) == 0 {
				e.printf("  Shares:    nothing yet (`ccam group set %s --add …`)\n", g.Name)
			}
			for i, item := range g.Items {
				label := "             "
				if i == 0 {
					label = "  Shares:    "
				}
				e.printf("%s%-11s %s\n", label, item, itemLabels[item])
			}
			e.printf("  Accounts:  %s\n", cmp.Or(strings.Join(accounts(cfg, g), ", "), "none (`ccam group join "+g.Name+" <account>`)"))
			e.printf("  Folder:    %s\n", e.paths.GroupDir(g.Folder()))
			if g.Has("mcp") {
				names := "none"
				if st, err := mcp.Load(mcp.Path(e.paths.GroupDir(g.Folder()))); err == nil && len(st.Servers) > 0 {
					names = strings.Join(st.Names(), ", ")
				}
				e.printf("  MCP:       %s\n", names)
			}
			for _, m := range g.Members {
				if !cfg.IsAccount(m) {
					e.printf("  Note: %s is listed as a member but is not an account; `ccam group leave %s %s` removes it\n", m, g.Name, m)
				}
			}
			return nil
		},
	}
}

// groupNames completes group names.
func groupNames(e *env) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return profileNames(e)(cmd, nil, "")
		}
		cfg, err := e.load()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var out []string
		for _, g := range cfg.Groups {
			out = append(out, g.Name+"\t"+strings.Join(g.Items, " "))
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	}
}

func parseList(values []string) []string {
	var out []string
	for _, v := range values {
		for _, f := range strings.Split(v, ",") {
			if f = strings.TrimSpace(f); f != "" {
				out = append(out, f)
			}
		}
	}
	return out
}

func groupAddCmd(e *env) *cobra.Command {
	var items, members []string
	var f changeFlags
	cmd := &cobra.Command{
		Use:   "add <group> --items <item,…> [--members <account,…>]",
		Short: "Create a group (its accounts' items move into it)",
		Long: "Creates a group that shares the items given with --items between the accounts\n" +
			"given with --members. The first account to bring an item fills the group's\n" +
			"folder; the others' copies are merged in or, for a file that differs, saved\n" +
			"in their " + share.BackupDir + " folder. You choose which copy of a file such as\n" +
			"settings.json wins; inside folders the group's copy stays, and the plan names\n" +
			"the files that differ.\n" +
			"Items: " + strings.Join(config.Items, ", "),
		Example: "  ccam group add team --items history --members work,side\n" +
			"  ccam group add docs --items memory,settings --members 1,3 --from work",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := config.ParseItems(strings.Join(items, ","))
			if err != nil {
				return err
			}
			if len(list) == 0 {
				return fmt.Errorf("--items is required, e.g. --items history,memory (items: %s)", strings.Join(config.Items, ", "))
			}
			names := parseList(members)
			_, done, err := changeGroups(e, func(c *config.Config) error {
				return c.AddGroup(args[0], list, names)
			}, f)
			if err != nil || !done {
				return err
			}
			e.printf("✓ Group %s shares %s\n", args[0], strings.Join(list, ", "))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&items, "items", nil, "what the group shares: "+strings.Join(config.Items, ","))
	cmd.Flags().StringSliceVar(&members, "members", nil, "accounts (names or numbers) in the group")
	f.add(cmd, true)
	return cmd
}

func groupSetCmd(e *env) *cobra.Command {
	var items, add, remove []string
	var f changeFlags
	cmd := &cobra.Command{
		Use:   "set <group> [--add <item,…>] [--remove <item,…>] [--items <item,…>]",
		Short: "Change what a group shares (moves the data of every account in it)",
		Long: "--add shares more items, --remove stops sharing items, --items sets the whole\n" +
			"list. Every account in the group changes together. An account that stops\n" +
			"sharing an item keeps a copy of the group's (start empty with --fresh); for mcp,\n" +
			"the group's servers go into its own .claude.json.\n" +
			"Items: " + strings.Join(config.Items, ", "),
		Example: "  ccam group set default --remove history\n" +
			"  ccam group set team --add memory --from work\n" +
			"  ccam --dry-run group set default --items mcp,settings",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: groupNames(e),
		RunE: func(cmd *cobra.Command, args []string) error {
			plus, err := config.ParseItems(strings.Join(add, ","))
			if err != nil {
				return err
			}
			minus, err := config.ParseItems(strings.Join(remove, ","))
			if err != nil {
				return err
			}
			exact, err := config.ParseItems(strings.Join(items, ","))
			if err != nil {
				return err
			}
			if !cmd.Flags().Changed("items") && len(plus)+len(minus) == 0 {
				return fmt.Errorf("give --add, --remove or --items (items: %s)", strings.Join(config.Items, ", "))
			}
			if cmd.Flags().Changed("items") && len(plus)+len(minus) > 0 {
				return fmt.Errorf("--items sets the whole list; use it without --add and --remove")
			}
			var result []string
			_, done, err := changeGroups(e, func(c *config.Config) error {
				g, err := c.MustFindGroup(args[0])
				if err != nil {
					return err
				}
				list := exact
				if !cmd.Flags().Changed("items") {
					list = slices.DeleteFunc(append(slices.Clone(g.Items), plus...), func(i string) bool { return slices.Contains(minus, i) })
				}
				result = list
				if err := c.SetGroupItems(g.Name, list); err != nil {
					return fmt.Errorf("%w.\nTake the item out of the other group, or move the account out of one of them (`ccam group leave`)", err)
				}
				return nil
			}, f)
			if err != nil || !done {
				return err
			}
			e.printf("✓ Group %s shares %s\n", args[0], cmp.Or(strings.Join(config.SortItems(result), ", "), "nothing"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&add, "add", nil, "items to start sharing")
	cmd.Flags().StringSliceVar(&remove, "remove", nil, "items to stop sharing")
	cmd.Flags().StringSliceVar(&items, "items", nil, "the full list of items")
	f.add(cmd, true)
	return cmd
}

func groupJoinCmd(e *env) *cobra.Command {
	var replace bool
	var f changeFlags
	cmd := &cobra.Command{
		Use:   "join <group> <account>...",
		Short: "Put accounts in a group (their items move into it)",
		Long: "Each account starts sharing the group's items. Its own copy of an item is\n" +
			"merged into the group or, for a file that differs, saved in its " + share.BackupDir + "\n" +
			"folder. You choose which copy of a file such as settings.json wins; inside\n" +
			"folders the group's copy stays, and the plan names the files that differ.\n" +
			"An account that already gets one of the items from another group is refused;\n" +
			"--replace takes it out of that group first.",
		Example:           "  ccam group join team work side\n  ccam group join team personal --replace",
		Args:              cobra.MinimumNArgs(2),
		ValidArgsFunction: groupNames(e),
		RunE: func(cmd *cobra.Command, args []string) error {
			left := map[string][]string{}
			_, done, err := changeGroups(e, func(c *config.Config) error {
				for _, a := range args[1:] {
					l, err := c.JoinGroup(a, args[0], replace)
					if err != nil {
						return err
					}
					left[a] = l
				}
				return nil
			}, f)
			if err != nil || !done {
				return err
			}
			for _, a := range args[1:] {
				if l := left[a]; len(l) > 0 {
					e.printf("✓ %s left group %s and joined group %s\n", a, strings.Join(l, ", "), args[0])
				} else {
					e.printf("✓ %s is in group %s\n", a, args[0])
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&replace, "replace", false, "leave the groups that share the same items first")
	f.add(cmd, false)
	return cmd
}

func groupLeaveCmd(e *env) *cobra.Command {
	var f changeFlags
	cmd := &cobra.Command{
		Use:   "leave <group> <account>...",
		Short: "Take accounts out of a group (each keeps a copy)",
		Long: "Each account stops sharing the group's items and keeps its own copy of them\n" +
			"(--fresh: starts without them); the group's MCP servers go into its own\n" +
			".claude.json. The group and the other accounts are not changed.",
		Example:           "  ccam group leave team side\n  ccam group leave default work --fresh",
		Args:              cobra.MinimumNArgs(2),
		ValidArgsFunction: groupNames(e),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, done, err := changeGroups(e, func(c *config.Config) error {
				g, err := c.MustFindGroup(args[0])
				if err != nil {
					return err
				}
				for _, a := range args[1:] {
					// A member that is not an account can still be removed.
					if !c.IsAccount(a) && g.HasMember(a) {
						g.Members = slices.DeleteFunc(g.Members, func(m string) bool { return m == a })
						continue
					}
					if err := c.LeaveGroup(a, args[0]); err != nil {
						return err
					}
				}
				return c.Resolve()
			}, f)
			if err != nil || !done {
				return err
			}
			e.printf("✓ %s left group %s\n", strings.Join(args[1:], ", "), args[0])
			return nil
		},
	}
	f.add(cmd, false)
	return cmd
}

func groupRenameCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:               "rename <group> <new-name>",
		Short:             "Rename a group (its folder keeps its name)",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: groupNames(e),
		RunE: func(cmd *cobra.Command, args []string) error {
			var folder string
			if _, err := e.update(func(c *config.Config) error {
				if err := c.RenameGroup(args[0], args[1]); err != nil {
					return err
				}
				g, _ := c.FindGroup(args[1])
				folder = g.Folder()
				return nil
			}); err != nil {
				return err
			}
			e.printf("✓ Group %s is now %s (folder: %s)\n", args[0], args[1], e.paths.GroupDir(folder))
			return nil
		},
	}
}

func groupRemoveCmd(e *env) *cobra.Command {
	var leave bool
	var f changeFlags
	cmd := &cobra.Command{
		Use:   "remove <group>",
		Short: "Remove a group (with --leave, its accounts keep a copy)",
		Long: "A group with accounts is removed only with --leave: each account then keeps\n" +
			"a copy of the group's items (--fresh: starts without them). The group's\n" +
			"folder goes to ccam's trash unless another group uses it; nothing is deleted.",
		Aliases:           []string{"rm"},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: groupNames(e),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			g, err := cfg.MustFindGroup(args[0])
			if err != nil {
				return err
			}
			members := accounts(cfg, g)
			if len(members) > 0 && !leave {
				return fmt.Errorf("group %s has accounts (%s). Add --leave to take them out first; each keeps a copy of %s", g.Name, strings.Join(members, ", "), cmp.Or(strings.Join(g.Items, ", "), "its items"))
			}
			folder := g.Folder()
			trash := len(cfg.FolderUsers(folder)) == 1 && fsutil.Exists(e.paths.GroupDir(folder))
			if trash {
				e.printf("The folder %s goes to ccam's trash.\n", e.paths.GroupDir(folder))
			}
			if len(members) == 0 && !f.yes && !e.dryRun {
				// No plan to confirm, so ask here.
				if !e.confirm(fmt.Sprintf("Remove group %s?", g.Name)) {
					e.printf("Cancelled.\n")
					return nil
				}
				f.yes = true
			}
			_, done, err := changeGroups(e, func(c *config.Config) error { return c.RemoveGroup(g.Name) }, f)
			if err != nil || !done {
				return err
			}
			if trash {
				dst := fsutil.UniqueDir(filepath.Join(e.paths.TrashDir(), "group-"+folder+"-"+fsutil.Stamp(time.Now())))
				if err := fsutil.Move(e.paths.GroupDir(folder), dst); err != nil {
					return fmt.Errorf("group %s is removed, but its folder could not go to the trash: %w", g.Name, err)
				}
				e.printf("✓ Group %s removed; its folder is in %s\n", g.Name, dst)
				return nil
			}
			e.printf("✓ Group %s removed\n", g.Name)
			return nil
		},
	}
	cmd.Flags().BoolVar(&leave, "leave", false, "take the group's accounts out first (each keeps a copy)")
	f.add(cmd, false)
	return cmd
}
