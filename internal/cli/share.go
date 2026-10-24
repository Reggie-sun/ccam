package cli

import (
	"fmt"
	"slices"
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
	cmd := &cobra.Command{
		Use:   "share <name|number>",
		Short: "Show what an account shares, and through which group",
		Long:  "Sharing is set by groups; change it with `ccam group` (join, leave, set).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, f := range append(slices.Clone(config.Items), "fresh", "yes", "prefer") {
				if cmd.Flags().Changed(f) {
					return fmt.Errorf("sharing is set by groups now: `ccam group join <group> %s`, `ccam group leave <group> %s`, or `ccam group set <group> --add <item>` (see `ccam group --help`)", args[0], args[0])
				}
			}
			cfg, prof, err := e.loadProfile(args[0])
			if err != nil {
				return err
			}
			return showShare(e, cfg, prof)
		},
	}
	// The old flags stay, hidden, to point to ccam group.
	for _, item := range config.Items {
		cmd.Flags().String(item, "", "replaced by ccam group")
		_ = cmd.Flags().MarkHidden(item)
	}
	cmd.Flags().Bool("fresh", false, "replaced by ccam group")
	cmd.Flags().BoolP("yes", "y", false, "replaced by ccam group")
	cmd.Flags().String("prefer", "", "replaced by ccam group")
	for _, f := range []string{"fresh", "yes", "prefer"} {
		_ = cmd.Flags().MarkHidden(f)
	}
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

func showShare(e *env, cfg *config.Config, prof *config.Profile) error {
	tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Sharing for %s\n", prof.Name)
	for _, item := range config.Items {
		state := "its own"
		if g := cfg.GroupFor(prof.Name, item); g != nil {
			state = "group " + g.Name
		}
		fmt.Fprintf(tw, "  %s\t%s\t%s\n", item, itemLabels[item], state)
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	var names []string
	for _, g := range cfg.GroupsOf(prof.Name) {
		names = append(names, g.Name)
	}
	if len(names) == 0 {
		e.printf("%s is in no group. `ccam group join <group> %s` adds it to one.\n", prof.Name, prof.Name)
		return nil
	}
	e.printf("Groups: %s. Change them with `ccam group`.\n", strings.Join(names, ", "))
	return nil
}

// askConflict asks how to settle c; who names the account when there are
// several in the change.
func askConflict(e *env, c *share.Conflict, who string) share.Resolution {
	grp := c.Group
	if c.Seed != "" {
		grp = c.Seed + " (moves into the group first)"
	}
	if who != "" {
		who += ": "
	}
	e.printf("\n%s%s differs from the group's copy.\n  mine:  %s\n  group: %s\n", who, c.Entry.Name, c.Mine, grp)
	for i, r := range c.Candidates {
		e.printf("  %d) %s\n", i+1, r)
	}
	for {
		e.printf("Choice [1]: ")
		line, err := e.lines().readLine()
		line = strings.TrimSpace(line)
		if line == "" || err != nil {
			return c.Candidates[0]
		}
		if n, err := strconv.Atoi(line); err == nil && n >= 1 && n <= len(c.Candidates) {
			return c.Candidates[n-1]
		}
	}
}
