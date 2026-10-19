package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/profile"
)

func mcpCmd(e *env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Manage shared MCP servers (the group's mcp.json)",
	}
	var group string
	cmd.PersistentFlags().StringVarP(&group, "group", "g", "", "target group (default: the only group sharing MCP, otherwise default)")
	cmd.AddCommand(mcpListCmd(e), mcpAddCmd(e, &group), mcpAddJSONCmd(e, &group), mcpRemoveCmd(e, &group), mcpImportCmd(e))
	return cmd
}

// mcpGroups lists groups that some profile shares MCP through.
func mcpGroups(cfg *config.Config) []string {
	seen := map[string]bool{}
	var out []string
	for _, p := range cfg.Profiles {
		if g := p.Share.Get("mcp"); g != config.Local && !seen[g] {
			seen[g] = true
			out = append(out, g)
		}
	}
	sort.Strings(out)
	return out
}

func pickGroup(cfg *config.Config, flag string) string {
	if flag != "" {
		return flag
	}
	if gs := mcpGroups(cfg); len(gs) == 1 {
		return gs[0]
	}
	return "default"
}

func mcpListCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Shared MCP servers per group, and servers only one account has",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			groups := mcpGroups(cfg)
			if len(groups) == 0 {
				e.printf("No account shares MCP. Turn it on with `ccam share <name> --mcp on`.\n")
			}
			shared := map[string]map[string]json.RawMessage{}
			for _, g := range groups {
				store, err := mcp.Load(mcp.Path(e.paths.GroupDir(g)))
				if err != nil {
					return err
				}
				shared[g] = store.Servers
				var members []string
				for _, p := range cfg.Profiles {
					if p.Share.Get("mcp") == g {
						members = append(members, p.Name)
					}
				}
				e.printf("group %s (%s)\n", g, strings.Join(members, ", "))
				if len(store.Servers) == 0 {
					e.printf("  (none) add one with `ccam mcp add`\n")
				}
				for _, n := range store.Names() {
					e.printf("  %-20s %s\n", n, describeServer(store.Servers[n]))
				}
			}
			for i := range cfg.Profiles {
				p := &cfg.Profiles[i]
				local, err := mcp.LocalServers(profile.Dir(e.paths, p))
				if err != nil || len(local) == 0 {
					continue
				}
				var only []string
				for n := range local {
					if _, ok := shared[p.Share.Get("mcp")][n]; !ok {
						only = append(only, n)
					}
				}
				if len(only) > 0 {
					sort.Strings(only)
					e.printf("only in account %s: %s\n", p.Name, strings.Join(only, ", "))
				}
			}
			return nil
		},
	}
}

func describeServer(raw json.RawMessage) string {
	var v struct {
		Type    string   `json:"type"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
		URL     string   `json:"url"`
	}
	_ = json.Unmarshal(raw, &v) // unknown shapes render as empty
	if v.URL != "" {
		t := v.Type
		if t == "" {
			t = "http"
		}
		return t + " " + v.URL
	}
	return "stdio " + strings.TrimSpace(v.Command+" "+strings.Join(v.Args, " "))
}

func mcpAddCmd(e *env, group *string) *cobra.Command {
	var envs, headers []string
	var transport, url string
	cmd := &cobra.Command{
		Use:   "add <name> [-- command args…] | add <name> --url <url>",
		Short: "Add an MCP server to a group",
		Example: "  ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code\n" +
			"  ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			def := map[string]any{}
			switch {
			case url != "":
				if transport == "" || transport == "stdio" {
					transport = "http"
				}
				def["type"], def["url"] = transport, url
				if len(headers) > 0 {
					h := map[string]string{}
					for _, kv := range headers {
						k, v, ok := strings.Cut(kv, ":")
						if !ok {
							return fmt.Errorf("--header takes 'Name: value': %q", kv)
						}
						h[strings.TrimSpace(k)] = strings.TrimSpace(v)
					}
					def["headers"] = h
				}
			case len(args) > 1:
				def["type"], def["command"], def["args"] = "stdio", args[1], args[2:]
				if len(envs) > 0 {
					m := map[string]string{}
					for _, kv := range envs {
						k, v, ok := strings.Cut(kv, "=")
						if !ok {
							return fmt.Errorf("--env takes KEY=VALUE: %q", kv)
						}
						m[k] = v
					}
					def["env"] = m
				}
			default:
				return fmt.Errorf("a command to run (after --) or --url is required")
			}
			raw, _ := json.Marshal(def)
			return addServer(e, *group, args[0], raw)
		},
	}
	cmd.Flags().StringArrayVarP(&envs, "env", "e", nil, "environment variable KEY=VALUE (stdio)")
	cmd.Flags().StringVar(&url, "url", "", "remote server URL (http/sse)")
	cmd.Flags().StringVar(&transport, "transport", "", "stdio|http|sse")
	cmd.Flags().StringArrayVar(&headers, "header", nil, "HTTP header 'Name: value'")
	return cmd
}

func mcpAddJSONCmd(e *env, group *string) *cobra.Command {
	return &cobra.Command{
		Use:   "add-json <name> <json>",
		Short: "Add an MCP server from a JSON definition (same format as claude mcp add-json)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addServer(e, *group, args[0], json.RawMessage(args[1]))
		},
	}
}

func addServer(e *env, groupFlag, name string, def json.RawMessage) error {
	if err := mcp.Validate(def); err != nil {
		return err
	}
	cfg, err := e.load()
	if err != nil {
		return err
	}
	g := pickGroup(cfg, groupFlag)
	path := mcp.Path(e.paths.GroupDir(g))
	store, err := mcp.Load(path)
	if err != nil {
		return err
	}
	_, existed := store.Servers[name]
	store.Servers[name] = def
	if err := store.Save(path); err != nil {
		return err
	}
	verb := "added"
	if existed {
		verb = "updated"
	}
	e.printf("✓ %s %s in group %s — applies from the next launch\n", name, verb, g)
	if len(mcpGroups(cfg)) == 0 {
		e.printf("Note: no account shares MCP yet. `ccam share <name> --mcp %s`\n", g)
	}
	return nil
}

func mcpRemoveCmd(e *env, group *string) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an MCP server from a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			g := pickGroup(cfg, *group)
			path := mcp.Path(e.paths.GroupDir(g))
			store, err := mcp.Load(path)
			if err != nil {
				return err
			}
			if _, ok := store.Servers[args[0]]; !ok {
				return fmt.Errorf("group %s has no server %q", g, args[0])
			}
			delete(store.Servers, args[0])
			if err := store.Save(path); err != nil {
				return err
			}
			e.printf("✓ %s removed from group %s\n", args[0], g)
			return nil
		},
	}
}

func mcpImportCmd(e *env) *cobra.Command {
	var keepLocal bool
	cmd := &cobra.Command{
		Use:   "import <account> [server…]",
		Short: "Move MCP servers only one account has into that account's shared group",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			prof, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			g := prof.Share.Get("mcp")
			if g == config.Local {
				return fmt.Errorf("account %s does not share MCP. Run `ccam share %s --mcp on` first", prof.Name, prof.Name)
			}
			dir := profile.Dir(e.paths, prof)
			local, err := mcp.LocalServers(dir)
			if err != nil {
				return err
			}
			names := args[1:]
			if len(names) == 0 {
				for n := range local {
					names = append(names, n)
				}
				sort.Strings(names)
			}
			if len(names) == 0 {
				e.printf("Account %s has no MCP servers of its own.\n", prof.Name)
				return nil
			}
			path := mcp.Path(e.paths.GroupDir(g))
			store, err := mcp.Load(path)
			if err != nil {
				return err
			}
			for _, n := range names {
				def, ok := local[n]
				if !ok {
					return fmt.Errorf("account %s has no server %q", prof.Name, n)
				}
				store.Servers[n] = def
			}
			if err := store.Save(path); err != nil {
				return err
			}
			e.printf("✓ Copied to group %s: %s\n", g, strings.Join(names, ", "))
			if keepLocal {
				return nil
			}
			// Local entries win over --mcp-config on name clashes, so drop them.
			if launcher.Running(e.paths, prof.Name) {
				e.printf("Note: account %s is running, so its own entries stay. Run this again after it exits to remove them.\n", prof.Name)
				return nil
			}
			if err := mcp.RemoveLocal(dir, filepath.Join(e.paths.State, "mcp", prof.Name), names); err != nil {
				return fmt.Errorf("could not remove the account's entries: %w", err)
			}
			e.printf("✓ Removed account %s's own entries (.claude.json backed up)\n", prof.Name)
			return nil
		},
	}
	cmd.Flags().BoolVar(&keepLocal, "keep-local", false, "keep the entries in the account's .claude.json")
	return cmd
}
