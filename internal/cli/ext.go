package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/ext"
	"github.com/jhyunwoo/ccam/internal/profile"
)

type extReport struct {
	Accounts []ext.Set  `json:"accounts"`
	Project  ext.Set    `json:"project"`
	Problems []extIssue `json:"problems"`
	Folder   string     `json:"folder"`
}

type extIssue struct {
	Account string `json:"account"`
	Problem string `json:"problem"`
}

func extCmd(e *env) *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "ext [name|number] [--json]",
		Short: "Show the skills, commands, agents, plugins, MCP servers and hooks each account loads",
		Long: "Lists, for every account (or the one named), the extensions Claude Code loads and\n" +
			"where each comes from: a share group, the account's own folder, or this project's\n" +
			".claude folder and .mcp.json. Below the table it lists problems: broken links, skill\n" +
			"folders without SKILL.md, and MCP servers defined differently in two places.\n" +
			"It only reads; nothing is changed.",
		Example: "  ccam ext\n  ccam ext work\n  ccam ext --json",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			profs := make([]*config.Profile, 0, len(cfg.Profiles))
			if len(args) == 1 {
				prof, err := cfg.MustFind(args[0])
				if err != nil {
					return err
				}
				profs = append(profs, prof)
			} else {
				for i := range cfg.Profiles {
					profs = append(profs, &cfg.Profiles[i])
				}
			}
			if len(profs) == 0 {
				e.printf("No accounts yet. Add one with `ccam add <name>`.\n")
				return nil
			}
			cwd, _ := os.Getwd()
			rep := extReport{Project: ext.ScanProject(cwd), Folder: cwd, Problems: []extIssue{}}
			for _, p := range profs {
				s := ext.ScanAccount(p.Name, profile.Dir(e.paths, p), e.paths.GroupsRoot(), p.Share)
				s.Problems = append(s.Problems, ext.Overlaps(s, rep.Project)...)
				rep.Accounts = append(rep.Accounts, s)
				for _, pr := range s.Problems {
					rep.Problems = append(rep.Problems, extIssue{p.Name, pr})
				}
			}
			if asJSON {
				data, _ := json.MarshalIndent(rep, "", "  ")
				e.printf("%s\n", data)
				return nil
			}
			printExt(e, rep)
			return nil
		},
	}
	cmd.ValidArgsFunction = profileNames(e)
	cmd.Flags().BoolVar(&asJSON, "json", false, "print everything as JSON")
	return cmd
}

func printExt(e *env, rep extReport) {
	kinds := []struct {
		title string
		pick  func(ext.Set) []ext.Item
	}{
		{"Skills", func(s ext.Set) []ext.Item { return s.Skills }},
		{"Commands", func(s ext.Set) []ext.Item { return s.Commands }},
		{"Agents", func(s ext.Set) []ext.Item { return s.Agents }},
		{"Plugins", func(s ext.Set) []ext.Item { return s.Plugins }},
		{"MCP servers", func(s ext.Set) []ext.Item { return s.MCP }},
	}
	empty := true
	for _, k := range kinds {
		names := ext.Names(rep.Accounts, k.pick)
		if len(names) == 0 {
			continue
		}
		empty = false
		e.printf("%s\n", k.title)
		tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
		header := []string{" "}
		for _, s := range rep.Accounts {
			header = append(header, s.Name)
		}
		fmt.Fprintln(tw, strings.Join(header, "\t"))
		for _, n := range names {
			row := []string{"  " + n}
			for _, s := range rep.Accounts {
				row = append(row, cell(k.pick(s), n))
			}
			fmt.Fprintln(tw, strings.Join(row, "\t"))
		}
		_ = tw.Flush()
		e.printf("\n")
	}
	if hookRows := hookTable(rep.Accounts); len(hookRows) > 0 {
		empty = false
		e.printf("Hooks (commands per event)\n")
		tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
		header := []string{" "}
		for _, s := range rep.Accounts {
			header = append(header, s.Name)
		}
		fmt.Fprintln(tw, strings.Join(header, "\t"))
		for _, r := range hookRows {
			fmt.Fprintln(tw, strings.Join(r, "\t"))
		}
		_ = tw.Flush()
		e.printf("\n")
	}
	if empty {
		e.printf("No skills, commands, agents, plugins, MCP servers or hooks in the accounts.\n\n")
	}
	if p := rep.Project; len(p.Skills)+len(p.Commands)+len(p.Agents)+len(p.MCP)+len(p.Hooks) > 0 {
		e.printf("This project (%s), for every account:\n", rep.Folder)
		for _, part := range []struct {
			label string
			items []ext.Item
		}{{"skills", p.Skills}, {"commands", p.Commands}, {"agents", p.Agents}, {"MCP servers", p.MCP}} {
			if len(part.items) > 0 {
				e.printf("  %s: %s\n", part.label, joinNames(part.items))
			}
		}
		for _, h := range p.Hooks {
			e.printf("  hook %s: %d\n", h.Event, h.Count)
		}
		e.printf("\n")
	}
	e.printf("Sources: a group name (shared), account (its own folder), - (not loaded).\n")
	if len(rep.Problems) == 0 {
		e.printf("✓ no problems found\n")
		return
	}
	for _, pr := range rep.Problems {
		e.printf("! %s: %s\n", pr.Account, pr.Problem)
	}
}

// cell is an account's source for one item, with ! when it has a problem.
func cell(items []ext.Item, name string) string {
	for _, it := range items {
		if it.Name == name {
			if it.Problem != "" {
				return it.Source + " !"
			}
			return it.Source
		}
	}
	return "-"
}

func hookTable(sets []ext.Set) [][]string {
	var events []string
	for _, s := range sets {
		for _, h := range s.Hooks {
			if !contains(events, h.Event) {
				events = append(events, h.Event)
			}
		}
	}
	var rows [][]string
	for _, ev := range events {
		row := []string{"  " + ev}
		for _, s := range sets {
			c := "-"
			for _, h := range s.Hooks {
				if h.Event == ev {
					c = fmt.Sprintf("%d (%s)", h.Count, h.Source)
				}
			}
			row = append(row, c)
		}
		rows = append(rows, row)
	}
	return rows
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func joinNames(items []ext.Item) string {
	names := make([]string, len(items))
	for i, it := range items {
		names[i] = it.Name
	}
	return strings.Join(names, ", ")
}
