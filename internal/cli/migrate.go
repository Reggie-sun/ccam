package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/migrate"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
	"github.com/jhyunwoo/ccam/internal/tui"
)

func migrateCmd(e *env) *cobra.Command {
	var group string
	var yes bool
	picked := map[string]*bool{}
	cmd := &cobra.Command{
		Use:   "migrate [folder]",
		Short: "Copy an existing Claude Code setup (~/.claude) into the shared group",
		Long: "Copies settings.json, CLAUDE.md, extensions (agents, commands, skills, plugins,\n" +
			"output-styles, rules) and user MCP servers from a Claude Code folder into a share\n" +
			"group, so every account that shares those items gets them. Name items to copy only\n" +
			"those. Without item flags, a checklist on a terminal lets you choose (everything\n" +
			"but conversation history is checked); elsewhere, or with --yes, everything except\n" +
			"conversation history is copied.\n\n" +
			"The folder (default ~/.claude, with MCP servers from ~/.claude.json) is only read:\n" +
			"nothing in it is moved, linked or changed. Logins (.credentials.json, .claude.json)\n" +
			"are never copied into the group. Conversation history is copied only with --history.\n\n" +
			"MCP servers' saved sign-ins (OAuth) go along: each account that gets a copied server\n" +
			"through the group gets its sign-in in its own login file, unless it already has one.\n" +
			"Accounts that are running, not logged in, imported or using the macOS Keychain are\n" +
			"skipped; sign in to the servers there with /mcp.\n\n" +
			"Prints the plan and asks first. When the group already has a different copy of a\n" +
			"file, you choose: keep group, use mine, or merge; with --yes the group's copy is kept.\n" +
			"Replaced group files are saved in <group>/" + share.BackupDir + "/<time>/.",
		Example: "  ccam migrate                     # ~/.claude into group default\n" +
			"  ccam --dry-run migrate           # show the plan, copy nothing\n" +
			"  ccam migrate ~/.claude-work --group work --settings --mcp\n" +
			"  ccam migrate --history --yes     # conversation history only",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			src := filepath.Join(e.paths.Home, ".claude")
			if len(args) == 1 {
				src = args[0]
			}
			g := "default"
			if cmd.Flags().Changed("group") {
				var s config.Share
				if err := s.Set("settings", group); err != nil {
					return err
				}
				if g = s.Get("settings"); g == config.Local {
					return fmt.Errorf("--group takes a group name, not %q", group)
				}
			}
			var items []string
			for _, item := range config.Items {
				if *picked[item] {
					items = append(items, item)
				}
			}
			defaulted := len(items) == 0
			if defaulted && !yes && interactive(e) && fsutil.Exists(src) {
				picked, ok, err := tui.RunMigratePicker(src, g, migrate.Found(src, e.paths.Home), migrate.DefaultItems)
				if err != nil {
					return err
				}
				if !ok {
					e.printf("Cancelled.\n")
					return nil
				}
				if len(picked) == 0 {
					e.printf("Nothing selected.\n")
					return nil
				}
				items, defaulted = picked, false
			}
			if defaulted {
				items = migrate.DefaultItems
			}
			pl, err := migrate.Build(migrate.Options{Source: src, Home: e.paths.Home, GroupsRoot: e.paths.GroupsRoot(), Group: g, Items: items, Accounts: migrateAccounts(e, cfg)})
			if err != nil {
				return err
			}
			if yes {
				for _, c := range pl.Conflicts {
					c.Choice = share.KeepGroup
				}
			}
			pl.Describe(e.stdout, !yes)
			if defaulted && hasHistory(pl.Source) {
				e.printf("  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.\n")
			}
			if pl.Empty() {
				e.printf("Nothing to copy.\n")
				return nil
			}
			if e.dryRun {
				return nil
			}
			if !yes {
				for _, c := range pl.Conflicts {
					c.Choice = askConflict(e, c)
				}
				if !pl.Changes() {
					e.printf("Nothing to copy: the group keeps its own copies.\n")
					return nil
				}
				if !e.confirm(fmt.Sprintf("Copy these into group %s?", g)) {
					e.printf("Cancelled.\n")
					return nil
				}
			}
			if err := pl.Apply(); err != nil {
				return fmt.Errorf("failed while copying: %w\nThe source folder was not changed. Group files replaced so far are saved in %s", err, pl.BackupDir)
			}
			migrateSummary(e, cfg, pl, yes, rerunArgs(args, cmd.Flags().Changed("group"), g))
			return nil
		},
	}
	cmd.Flags().StringVarP(&group, "group", "g", "default", "share group to copy into")
	for _, item := range config.Items {
		picked[item] = new(bool)
		help := "copy " + itemLabels[item]
		switch item {
		case "extensions":
			help = "copy agents, commands, skills, plugins, output-styles, rules"
		case "mcp":
			help = "copy user MCP servers (mcpServers in .claude.json)"
		case "history":
			help = "copy conversation history (projects/, history.jsonl); off unless given"
		}
		cmd.Flags().BoolVar(picked[item], item, false, help)
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "copy without asking; on a conflict the group's copy is kept")
	return cmd
}

// migrateAccounts are the accounts that may receive the source's MCP sign-ins.
func migrateAccounts(e *env, cfg *config.Config) []migrate.Account {
	out := make([]migrate.Account, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		out = append(out, migrate.Account{
			Name: p.Name, Dir: profile.Dir(e.paths, p),
			MCP: p.Share.Get("mcp"), Extensions: p.Share.Get("extensions"),
			SessionDir: launcher.SessionDir(e.paths, p.Name), Imported: p.Dir != "",
		})
	}
	return out
}

func hasHistory(src string) bool {
	for _, en := range share.Layout["history"] {
		if fsutil.Exists(filepath.Join(src, en.Name)) {
			return true
		}
	}
	return false
}

func rerunArgs(args []string, groupSet bool, g string) string {
	out := "ccam migrate"
	if len(args) == 1 {
		out += " " + args[0]
	}
	if groupSet {
		out += " --group " + g
	}
	return out
}

// migrateSummary says what was copied and which accounts see it, since the
// group only reaches accounts that share each item through it.
func migrateSummary(e *env, cfg *config.Config, pl *migrate.Plan, yes bool, rerun string) {
	if copied := pl.Copied(); len(copied) > 0 {
		e.printf("✓ Copied into group %s: %s\n", pl.Group, strings.Join(copied, ", "))
	}
	if kept := pl.Kept(); len(kept) > 0 {
		e.printf("Kept the group's copy of: %s\n", strings.Join(kept, ", "))
		if yes {
			e.printf("To choose instead, run `%s` without --yes.\n", rerun)
		}
	}
	if fsutil.Exists(pl.BackupDir) {
		e.printf("The group's previous copies are saved in %s\n", pl.BackupDir)
	}
	if l := pl.Logins; l != nil && len(l.Skipped) > 0 {
		e.printf("MCP sign-ins were not copied to: %s. Sign in to those servers there with /mcp.\n", strings.Join(l.Skipped, ", "))
	}
	items := pl.Items()
	if len(items) == 0 {
		return
	}
	if len(cfg.Profiles) == 0 {
		if pl.Group == "default" {
			e.printf("No accounts yet. `ccam add <name>` shares settings, CLAUDE.md, extensions and MCP through group default.\n")
		} else {
			e.printf("No accounts yet. Add one with `ccam add <name> --share settings=%s` (one --share per item).\n", pl.Group)
		}
		return
	}
	value := pl.Group
	if value == "default" {
		value = "on"
	}
	e.printf("Accounts that share these through group %s:\n", pl.Group)
	missing := map[string][]string{} // profile → item flags to turn on
	for _, item := range items {
		var names []string
		for _, p := range cfg.Profiles {
			if p.Share.Get(item) == pl.Group {
				names = append(names, p.Name)
			} else {
				missing[p.Name] = append(missing[p.Name], "--"+item+" "+value)
			}
		}
		who := strings.Join(names, ", ")
		if who == "" {
			who = "none yet"
		}
		e.printf("  %-11s %s\n", item, who)
	}
	if len(missing) == 0 {
		return
	}
	e.printf("To share them with the other accounts:\n")
	for _, p := range cfg.Profiles {
		if flags := missing[p.Name]; len(flags) > 0 {
			e.printf("  ccam share %s %s\n", p.Name, strings.Join(flags, " "))
		}
	}
	if pl.Logins != nil {
		e.printf("Then run `%s` again to give them the MCP sign-ins too.\n", rerun)
	}
}
