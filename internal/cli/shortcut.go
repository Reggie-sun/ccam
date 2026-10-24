package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/shim"
)

// saveShortcut saves cfg, in which an account's shortcut changed from old
// to name, and updates the commands in the bin folder.
func saveShortcut(e *env, cfg *config.Config, old, name string) error {
	if err := e.save(cfg); err != nil {
		return err
	}
	if self, err := ccamPath(); err == nil && old != "" {
		shim.Remove(e.paths, self, old)
	}
	return makeShortcuts(e, cfg, name)
}

// makeShortcuts creates the shortcut links cfg names. It fails only when
// name's cannot be made; a file in the way of another account's is said.
func makeShortcuts(e *env, cfg *config.Config, name string) error {
	_, err := installShims(e, cfg)
	if err == nil {
		return nil
	}
	if self, _ := ccamPath(); name == "" || shim.Linked(e.paths.Bin, self, name) {
		fmt.Fprintln(e.stderr, "Note: some shortcuts could not be created:", err)
		return nil
	}
	return fmt.Errorf("could not create the shortcut: %w", err)
}

func shortcutCmd(e *env) *cobra.Command {
	var remove bool
	cmd := &cobra.Command{
		Use:   "shortcut <name|number> [command]",
		Short: "Show, create, rename or remove an account's shortcut command",
		Long: "A shortcut is a command in ~/.local/bin that runs Claude Code as one account,\n" +
			"e.g. `ccwork --resume` or `ccwork -p \"hi\"`. Arguments go to claude unchanged.\n\n" +
			"  ccam shortcut <account>            print the current shortcut\n" +
			"  ccam shortcut <account> <command>  create it, or rename the existing one\n" +
			"  ccam shortcut <account> --remove   delete it\n\n" +
			"Names are a-z, 0-9, '.', '-' and '_'; a name cannot be a ccam command, another\n" +
			"account's shortcut, or hide another program on PATH.",
		Example: "  ccam shortcut work ccwork\n  ccam shortcut 2 ccside\n  ccam shortcut work --remove",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			prof, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			if remove && len(args) == 2 {
				return fmt.Errorf("give a command or --remove, not both")
			}
			old := prof.Shortcut
			switch {
			case len(args) == 1 && !remove:
				if old == "" {
					e.printf("%s: none (add one with `ccam shortcut %s %s`)\n", prof.Name, prof.Name, shim.Suggest(prof.Name))
				} else {
					e.printf("%s: %s\n", prof.Name, old)
				}
				return nil
			case remove:
				if old == "" {
					e.printf("%s has no shortcut.\n", prof.Name)
					return nil
				}
				prof.Shortcut = ""
			default:
				if args[1] == old {
					e.printf("%s: %s (unchanged)\n", prof.Name, old)
					if e.dryRun {
						return nil
					}
					return makeShortcuts(e, cfg, old)
				}
				if err := checkShortcut(e, cfg, prof.Name, args[1]); err != nil {
					return err
				}
				prof.Shortcut = args[1]
			}
			if e.dryRun {
				if prof.Shortcut == "" {
					e.printf("Would remove %s, the shortcut of %s.\n", old, prof.Name)
				} else {
					e.printf("Would set the shortcut of %s to %s.\n", prof.Name, prof.Shortcut)
				}
				return nil
			}
			if err := saveShortcut(e, cfg, old, prof.Shortcut); err != nil {
				return err
			}
			switch {
			case prof.Shortcut == "":
				e.printf("✓ %s removed; %s has no shortcut now\n", old, prof.Name)
			case old == "":
				e.printf("✓ shortcut for %s: %s\n", prof.Name, prof.Shortcut)
			default:
				e.printf("✓ shortcut for %s: %s (was %s)\n", prof.Name, prof.Shortcut, old)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&remove, "remove", false, "delete the account's shortcut")
	return cmd
}
