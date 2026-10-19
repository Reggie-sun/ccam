package cli

import (
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
)

// applyShare links a freshly added (empty) profile into its share groups.
func applyShare(e *env, cfg *config.Config, prof *config.Profile) error {
	pl, err := share.Build(profile.Dir(e.paths, prof), e.paths.GroupsRoot(), config.Share{}, prof.Share, share.Options{})
	if err != nil {
		return err
	}
	return pl.Apply()
}

var itemLabels = map[string]string{
	"mcp":        "MCP servers",
	"settings":   "settings.json",
	"memory":     "global CLAUDE.md",
	"extensions": "agents, commands, skills, plugins",
	"history":    "conversation history",
}

func shareCmd(e *env) *cobra.Command {
	var fresh, yes bool
	var prefer string
	values := map[string]*string{}
	cmd := &cobra.Command{
		Use:   "share <name|number> [--mcp on|off|<group>] [--history …]",
		Short: "Change shared items (migrates data)",
		Long: "Give each item on (= the default group), off (= local), or a group name.\n" +
			"Items: --mcp --settings --memory --extensions --history\n" +
			"Without options, prints the current state.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			prof, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			target := prof.Share
			changed := false
			for _, item := range config.Items {
				if cmd.Flags().Changed(item) {
					if err := target.Set(item, *values[item]); err != nil {
						return err
					}
					changed = true
				}
			}
			if !changed {
				return showShare(e, prof)
			}
			res, err := parsePrefer(prefer)
			if err != nil {
				return err
			}
			return changeShare(e, cfg, prof, target, share.Options{Fresh: fresh}, yes, res)
		},
	}
	for _, item := range config.Items {
		values[item] = new(string)
		cmd.Flags().StringVar(values[item], item, "", itemLabels[item]+": on|off|<group>")
	}
	cmd.Flags().BoolVar(&fresh, "fresh", false, "when turning sharing off, start empty instead of copying the group's content")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "run without asking (conflicts follow --prefer)")
	cmd.Flags().StringVar(&prefer, "prefer", "", "conflict handling: group|mine|merge (default: ask; with --yes, each item's default)")
	return cmd
}

func parsePrefer(s string) (*share.Resolution, error) {
	var r share.Resolution
	switch s {
	case "":
		return nil, nil
	case "group":
		r = share.KeepGroup
	case "mine":
		r = share.UseMine
	case "merge":
		r = share.MergeJSON
	default:
		return nil, fmt.Errorf("--prefer must be group, mine or merge")
	}
	return &r, nil
}

func showShare(e *env, prof *config.Profile) error {
	tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Sharing for %s\n", prof.Name)
	for _, item := range config.Items {
		g := prof.Share.Get(item)
		if g == config.Local {
			g = "local"
		} else {
			g = "shared (group " + g + ")"
		}
		fmt.Fprintf(tw, "  %s\t%s\t--%s\n", itemLabels[item], g, item)
	}
	return tw.Flush()
}

// changeShare prints the plan, asks, resolves conflicts and applies.
func changeShare(e *env, cfg *config.Config, prof *config.Profile, target config.Share, opt share.Options, yes bool, prefer *share.Resolution) error {
	dir := profile.Dir(e.paths, prof)
	pl, err := share.Build(dir, e.paths.GroupsRoot(), prof.Share, target, opt)
	if err != nil {
		return err
	}
	for _, c := range pl.Conflicts {
		if prefer != nil {
			for _, cand := range c.Candidates {
				if cand == *prefer || (*prefer == share.MergeJSON && cand == share.MergeLines) {
					c.Choice = cand
				}
			}
			continue
		}
		if yes {
			continue
		}
		c.Choice = askConflict(e, c)
	}
	e.printf("Plan (%s):\n", prof.Name)
	pl.Describe(e.stdout)
	if e.dryRun {
		return nil
	}
	if !pl.Empty() && !yes {
		if !e.confirm("Proceed?") {
			e.printf("Cancelled.\n")
			return nil
		}
	}
	// Save the target first: if Apply stops halfway, doctor --fix and a
	// rerun build their plan toward it and finish the change.
	prof.Share = target
	if err := e.save(cfg); err != nil {
		return err
	}
	if err := pl.Apply(); err != nil {
		return fmt.Errorf("failed while applying: %w\nThe originals are in %s/%s. Fix the cause, then finish with `ccam doctor --fix`", err, dir, share.BackupDir)
	}
	e.printf("✓ Sharing changed for %s\n", prof.Name)
	return nil
}

func askConflict(e *env, c *share.Conflict) share.Resolution {
	e.printf("\n%s differs from the group's copy.\n  mine:  %s\n  group: %s\n", c.Entry.Name, c.Mine, c.Group)
	for i, r := range c.Candidates {
		e.printf("  %d) %s\n", i+1, r)
	}
	for {
		e.printf("Choice [1]: ")
		line, err := e.in.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" || err != nil {
			return c.Candidates[0]
		}
		if n, err := strconv.Atoi(line); err == nil && n >= 1 && n <= len(c.Candidates) {
			return c.Candidates[n-1]
		}
	}
}
