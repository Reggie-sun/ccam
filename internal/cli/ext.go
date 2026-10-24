package cli

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/ext"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/migrate"
	"github.com/jhyunwoo/ccam/internal/profile"
)

type extReport struct {
	Accounts []ext.Set  `json:"accounts"`
	Project  ext.Set    `json:"project"`
	Problems []extIssue `json:"problems"`
	Folder   string     `json:"folder"`
}

type extIssue struct {
	Account string `json:"account,omitempty"`
	Group   string `json:"group,omitempty"` // a problem every account in the group has
	Problem string `json:"problem"`
}

// scanExt reads the extensions of profs and of the project in cwd.
func scanExt(e *env, cfg *config.Config, profs []*config.Profile, cwd string) extReport {
	rep := extReport{Accounts: []ext.Set{}, Project: ext.ScanProject(cwd), Folder: cwd, Problems: []extIssue{}}
	names := map[string]string{}
	for _, g := range cfg.Groups {
		if g.Dir != "" {
			names[g.Dir] = g.Name
		}
	}
	var project []string // the project's .mcp.json servers
	for _, it := range rep.Project.MCP {
		project = append(project, it.Name)
	}
	seen := map[extIssue]bool{}
	dirs := map[string]string{}
	connectors := map[string][]string{}
	for _, p := range profs {
		dir := profile.Dir(e.paths, p)
		dirs[p.Name] = dir
		s := ext.ScanAccount(p.Name, dir, e.paths.GroupsRoot(), p.Share, names, cwd)
		s.Problems = append(s.Problems, ext.Overlaps(s, rep.Project)...)
		if g := cfg.GroupFor(p.Name, "mcp"); g != nil && len(accounts(cfg, g)) > 1 {
			for _, it := range s.MCP {
				if it.Source == ext.Account && !ext.PluginServer(it.Name) {
					s.Problems = append(s.Problems, fmt.Sprintf("MCP server %s is only in this account's .claude.json, so the other accounts of group %s do not load it: `ccam mcp import %s %s` shares it", it.Name, g.Name, p.Name, it.Name))
				}
			}
		}
		if w := ext.Waiting(dir, cwd, project); len(w) > 0 {
			whose := "this account"
			if g := cfg.GroupFor(p.Name, "settings"); g != nil {
				whose = "the accounts of group " + g.Name // its settings.json is the group's
			}
			s.Problems = append(s.Problems, fmt.Sprintf("has not trusted this folder yet: Claude Code asks once, and until then the project's MCP servers %s wait for approval there. "+
				"To approve them for %s: \"enabledMcpjsonServers\": [%s] in its settings.json", strings.Join(w, ", "), whose, quoted(w)))
		}
		var used []string
		if ok, _ := profile.DecodeFileKey(filepath.Join(dir, ".claude.json"), "claudeAiMcpEverConnected", &used); ok {
			connectors[p.Name] = used
		}
		rep.Accounts = append(rep.Accounts, s)
		for _, pr := range s.Problems {
			rep.Problems = append(rep.Problems, extIssue{Account: p.Name, Problem: pr})
		}
		for _, g := range slices.Sorted(maps.Keys(s.GroupProblems)) {
			for _, pr := range s.GroupProblems[g] {
				if is := (extIssue{Group: g, Problem: pr}); !seen[is] {
					seen[is] = true
					rep.Problems = append(rep.Problems, is)
				}
			}
		}
	}
	signedIn := map[string]map[string]bool{}
	for _, a := range ext.SignInAccounts(rep.Accounts) {
		signedIn[a] = mcp.SignedIn(dirs[a], e.paths.Home)
	}
	for _, g := range ext.SignInGaps(rep.Accounts, signedIn) {
		for _, a := range g.Missing {
			rep.Problems = append(rep.Problems, extIssue{Account: a, Problem: fmt.Sprintf("MCP server %s needs a sign-in in this account: %s signed in, and a sign-in holds for one account only. Sign in with /mcp in a session of %s", g.Server, strings.Join(g.SignedIn, ", "), a)})
		}
	}
	rep.Problems = append(rep.Problems, connectorGaps(connectors)...)
	if len(profs) == len(cfg.Profiles) {
		if from, servers := homeOnlyMCP(e, cfg); len(servers) > 0 {
			rep.Problems = append(rep.Problems, extIssue{Problem: homeOnlyText(from, servers, cfg)})
		}
	}
	return rep
}

// connectorGaps names the claude.ai connectors some accounts have used and
// others not: they belong to the claude.ai account and come with its login,
// so ccam cannot share them.
func connectorGaps(used map[string][]string) []extIssue {
	all := map[string][]string{} // connector → accounts that used it
	for _, a := range slices.Sorted(maps.Keys(used)) {
		for _, c := range used[a] {
			all[c] = append(all[c], a)
		}
	}
	var out []extIssue
	for _, a := range slices.Sorted(maps.Keys(used)) {
		var lack []string
		for _, c := range slices.Sorted(maps.Keys(all)) {
			if !slices.Contains(used[a], c) {
				lack = append(lack, c+" ("+strings.Join(all[c], ", ")+")")
			}
		}
		if len(lack) > 0 {
			out = append(out, extIssue{Account: a, Problem: "has not used claude.ai connectors other accounts have: " + strings.Join(lack, ", ") +
				". Connectors belong to the claude.ai account; connect them there for this one (claude.ai, Settings, Connectors)"})
		}
	}
	return out
}

// homeOnlyText is the problem of servers plain Claude Code has that no
// account loads.
func homeOnlyText(from string, servers []string, cfg *config.Config) string {
	return fmt.Sprintf("%s (plain claude's) has MCP servers no account loads: %s. `%s` copies them into the shared group", from, strings.Join(servers, ", "), migrate.MCPCommand(cfg))
}

// quoted is names as JSON strings, comma-separated.
func quoted(names []string) string {
	q := make([]string, len(names))
	for i, n := range names {
		q[i] = strconv.Quote(n)
	}
	return strings.Join(q, ", ")
}

func extCmd(e *env) *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "ext [name|number] [--json]",
		Short: "Show the skills, commands, agents, plugins, MCP servers and hooks each account loads",
		Long: "Lists, for every account (or the one named), the extensions Claude Code loads and\n" +
			"where each comes from: a share group, the account's own folder (local: its own,\n" +
			"for this folder only), or this project's .claude folder and .mcp.json. MCP servers\n" +
			"include the plugins' (plugin:<plugin>:<server>); plugins turned off show (off).\n" +
			"Below the table it lists problems: broken links, skill folders without SKILL.md,\n" +
			"MCP servers defined differently in two places or only in one account, servers\n" +
			"another account is signed in to (sign in with /mcp), project servers waiting for\n" +
			"an account's approval, and servers of plain claude's ~/.claude.json no account\n" +
			"loads. It only reads; nothing is changed.",
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
			if len(profs) == 0 && !asJSON {
				e.printf("No accounts yet. Add one with `ccam add <name>`.\n")
				return nil
			}
			cwd, _ := os.Getwd()
			rep := scanExt(e, cfg, profs, cwd)
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
		e.printf("This project (%s):\n", rep.Folder)
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
	e.printf("Sources: a group name (shared), account (its own folder), local (its own, for this folder only), - (not loaded); (off): installed, not turned on.\n")
	if len(rep.Problems) == 0 {
		e.printf("✓ no problems found\n")
		return
	}
	for _, pr := range rep.Problems {
		switch {
		case pr.Group != "":
			e.printf("! group %s: %s\n", pr.Group, pr.Problem)
		case pr.Account != "":
			e.printf("! %s: %s\n", pr.Account, pr.Problem)
		default:
			e.printf("! %s\n", pr.Problem)
		}
	}
}

// cell is an account's source for one item, (off) for a plugin turned
// off, with ! when it has a problem.
func cell(items []ext.Item, name string) string {
	for _, it := range items {
		if it.Name == name {
			c := it.Source
			if it.Off {
				c += " (off)"
			}
			if it.Problem != "" {
				c += " !"
			}
			return c
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
