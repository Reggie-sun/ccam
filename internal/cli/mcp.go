package cli

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/ext"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/migrate"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
)

func mcpCmd(e *env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Manage shared MCP servers (the group's mcp.json)",
	}
	var group string
	cmd.PersistentFlags().StringVarP(&group, "group", "g", "", "target group (default: the only group that shares MCP, otherwise default)")
	cmd.AddCommand(mcpListCmd(e), mcpAddCmd(e, &group), mcpAddJSONCmd(e, &group), mcpRemoveCmd(e, &group), mcpImportCmd(e))
	return cmd
}

// mcpGroups lists the groups that share MCP.
func mcpGroups(cfg *config.Config) []*config.Group {
	var out []*config.Group
	for i := range cfg.Groups {
		if cfg.Groups[i].Has("mcp") {
			out = append(out, &cfg.Groups[i])
		}
	}
	return out
}

// pickGroup returns the name and folder of the group whose mcp.json a
// command changes. A name that is no group yet is a folder that a group
// created later with that name starts from.
func pickGroup(cfg *config.Config, flag string) (name, folder string, err error) {
	if flag == "" {
		if gs := mcpGroups(cfg); len(gs) == 1 {
			return gs[0].Name, gs[0].Folder(), nil
		}
		flag = "default"
	}
	if g, ok := cfg.FindGroup(flag); ok {
		return g.Name, g.Folder(), nil
	}
	if err := folderInUse(cfg, flag); err != nil {
		return "", "", err
	}
	return flag, flag, config.ValidateGroupName(flag)
}

// folderInUse refuses a name that is no group but the folder of one (a
// group renamed since): writing there would change that group.
func folderInUse(cfg *config.Config, name string) error {
	if users := cfg.FolderUsers(name); len(users) > 0 {
		return fmt.Errorf("there is no group %s; group %s uses the folder groups/%s (use --group %s)", name, users[0].Name, name, users[0].Name)
	}
	return nil
}

func mcpListCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Shared MCP servers per group, and servers only one account has",
		Long: "Lists the servers each group shares, then the servers only one account has:\n" +
			"in its .claude.json for every folder (user scope), or for one folder only (the\n" +
			"local scope `claude mcp add` uses unless given -s user). `ccam mcp import` shares\n" +
			"them. `claude mcp list` does not show a group's servers: Claude Code reads them\n" +
			"only when a session starts.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			groups := mcpGroups(cfg)
			if len(groups) == 0 {
				e.printf("No group shares MCP. Turn it on with `ccam group set <group> --add mcp`.\n")
			}
			shared := map[string]map[string]json.RawMessage{} // folder → servers
			for _, g := range groups {
				store, err := mcp.Load(mcp.Path(e.paths.GroupDir(g.Folder())))
				if err != nil {
					return err
				}
				shared[g.Folder()] = store.Servers
				e.printf("group %s (%s)\n", g.Name, strings.Join(accounts(cfg, g), ", "))
				if len(store.Servers) == 0 {
					e.printf("  (none) add one with `ccam mcp add`\n")
				}
				for _, n := range store.Names() {
					e.printf("  %-20s %s\n", n, describeServer(store.Servers[n]))
				}
			}
			for i := range cfg.Profiles {
				p := &cfg.Profiles[i]
				dir := profile.Dir(e.paths, p)
				group, sharing := shared[p.Share.Get("mcp")]
				local, _ := mcp.LocalServers(dir)
				var only, own []string
				for n, def := range local {
					switch have, ok := group[n]; {
					case !ok:
						only = append(only, n)
					case !mcp.SameDef(have, def):
						own = append(own, n)
					}
				}
				if len(only) > 0 {
					sort.Strings(only)
					e.printf("only in account %s: %s\n", p.Name, strings.Join(only, ", "))
				}
				if len(own) > 0 && sharing {
					sort.Strings(own)
					e.printf("account %s has its own %s, used in its sessions in place of the group's\n", p.Name, strings.Join(own, ", "))
				}
				folders, _ := mcp.FolderServers(dir)
				for _, folder := range slices.Sorted(maps.Keys(folders)) {
					e.printf("only in account %s, for %s only: %s\n", p.Name, folder, strings.Join(slices.Sorted(maps.Keys(folders[folder])), ", "))
				}
			}
			// The project's servers load in every account that trusts the folder.
			if cwd, err := os.Getwd(); err == nil {
				if store, err := mcp.Load(filepath.Join(cwd, ".mcp.json")); err == nil && len(store.Servers) > 0 {
					e.printf("this folder's .mcp.json: %s\n", strings.Join(store.Names(), ", "))
					for i := range cfg.Profiles {
						p := &cfg.Profiles[i]
						if w := ext.Waiting(profile.Dir(e.paths, p), cwd, store.Names()); len(w) > 0 {
							e.printf("  waiting for approval in account %s, which has not trusted the folder: %s\n", p.Name, strings.Join(w, ", "))
						}
					}
				}
			}
			if from, names := homeOnlyMCP(e, cfg); len(names) > 0 {
				e.printf("in no account, only in %s (plain claude's): %s. `%s` copies them into the shared group\n", from, strings.Join(names, ", "), migrate.MCPCommand(cfg))
			}
			return nil
		},
	}
}

// homeOnlyMCP lists the user-scope servers of plain Claude Code's
// .claude.json (~/.claude.json) that no account loads.
func homeOnlyMCP(e *env, cfg *config.Config) (from string, names []string) {
	var accts, groups []string
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		accts = append(accts, profile.Dir(e.paths, p))
		if f := p.Share.Get("mcp"); f != config.Local {
			groups = append(groups, e.paths.GroupDir(f))
		}
	}
	return migrate.HomeOnly(filepath.Join(e.paths.Home, ".claude"), e.paths.Home, accts, groups)
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
	g, folder, err := pickGroup(cfg, groupFlag)
	if err != nil {
		return err
	}
	existed := false
	noteCopies(e, cfg, folder)
	if err := mcp.Update(mcp.Path(e.paths.GroupDir(folder)), func(store *mcp.Store) error {
		_, existed = store.Servers[name]
		store.Servers[name] = def
		return nil
	}); err != nil {
		return err
	}
	verb := "added"
	if existed {
		verb = "updated"
	}
	e.printf("✓ %s %s in group %s — applies from the next launch\n", name, verb, g)
	refreshCopies(e, cfg, folder, "")
	if gr, ok := cfg.FindGroup(g); !ok {
		e.printf("Note: there is no group %s yet. `ccam group add %s --items mcp --members <account,…>` shares it.\n", g, g)
	} else if !gr.Has("mcp") {
		e.printf("Note: group %s does not share MCP yet. `ccam group set %s --add mcp`\n", g, g)
	}
	return nil
}

// noteCopies has the copies an older ccam made of the servers of the group
// in folder g noted with their hashes (mcp.NoteCopies), right before g's
// mcp.json changes: refreshCopies and Sync then update them, where they
// would take them for the accounts' own.
func noteCopies(e *env, cfg *config.Config, g string) {
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		if prof.Share.Get("mcp") == g {
			_ = mcp.NoteCopies(profile.Dir(e.paths, prof), e.paths.GroupDir(g), launcher.MCPStateDir(e.paths, prof.Name))
		}
	}
}

// refreshCopies updates the accounts whose .claude.json holds a copy of the
// servers of the group in folder g, made when they started Remote Control.
// Under the flag strategy an entry in .claude.json wins over --mcp-config,
// so a stale copy would hide the change.
func refreshCopies(e *env, cfg *config.Config, g, except string) {
	if cfg.Strategy() == config.MCPSync {
		return // every launch syncs
	}
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		state := filepath.Join(e.paths.State, "mcp", prof.Name)
		if prof.Name == except || prof.Share.Get("mcp") != g || !mcp.HasCopy(state) {
			continue
		}
		if launcher.Running(e.paths, prof.Name, profile.Dir(e.paths, prof)) {
			e.printf("Note: account %s is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.\n", prof.Name)
			continue
		}
		res, err := mcp.Sync(profile.Dir(e.paths, prof), e.paths.GroupDir(g), state)
		if err != nil {
			e.printf("Note: could not update the copy in account %s's .claude.json: %v\n", prof.Name, err)
			continue
		}
		if res.Wrote {
			e.printf("✓ Updated the copy in account %s's .claude.json (made for Remote Control)\n", prof.Name)
		}
		if len(res.Own) > 0 {
			e.printf("Note: account %s has its own %s, so the group's of those names are not copied there\n", prof.Name, strings.Join(res.Own, ", "))
		}
	}
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
			g, folder, err := pickGroup(cfg, *group)
			if err != nil {
				return err
			}
			noteCopies(e, cfg, folder)
			if err := mcp.Update(mcp.Path(e.paths.GroupDir(folder)), func(store *mcp.Store) error {
				if _, ok := store.Servers[args[0]]; !ok {
					return fmt.Errorf("group %s has no server %q", g, args[0])
				}
				delete(store.Servers, args[0])
				return nil
			}); err != nil {
				return err
			}
			e.printf("✓ %s removed from group %s\n", args[0], g)
			refreshCopies(e, cfg, folder, "")
			return nil
		},
	}
}

func mcpImportCmd(e *env) *cobra.Command {
	var keepLocal, force bool
	cmd := &cobra.Command{
		Use:   "import <account> [server…]",
		Short: "Move MCP servers only one account has into that account's shared group",
		Long: "Copies the account's own servers into the mcp.json of the group it shares MCP\n" +
			"through, then takes them out of its .claude.json, where they would win over the\n" +
			"group's (--keep-local leaves them). Without server names, all of its own.\n" +
			"A server `claude mcp add` gave one folder (local scope) then loads from the group\n" +
			"for every account in it, in every folder. A group server of the same name that\n" +
			"differs is left (--force replaces it, after a backup of mcp.json).",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			prof, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			gr := cfg.GroupFor(prof.Name, "mcp")
			if gr == nil {
				return fmt.Errorf("account %s does not share MCP. Put it in a group that does first: `ccam group join <group> %s`", prof.Name, prof.Name)
			}
			g, folder := gr.Name, gr.Folder()
			dir := profile.Dir(e.paths, prof)
			user, err := mcp.LocalServers(dir)
			if err != nil {
				return err
			}
			folders, err := mcp.FolderServers(dir)
			if err != nil {
				return err
			}
			names := args[1:]
			if len(names) == 0 {
				names = slices.Collect(maps.Keys(user))
				for _, s := range folders {
					for n := range s {
						if !slices.Contains(names, n) {
							names = append(names, n)
						}
					}
				}
				sort.Strings(names)
			}
			if len(names) == 0 {
				e.printf("Account %s has no MCP servers of its own.\n", prof.Name)
				return nil
			}
			path := mcp.Path(e.paths.GroupDir(folder))
			noteCopies(e, cfg, folder)
			// Held until the account's entries are gone: no other ccam can
			// drop a server from the group in between.
			unlock, err := mcp.Lock(path)
			if err != nil {
				return err
			}
			release := sync.OnceFunc(unlock)
			defer release()
			store, err := mcp.Load(path)
			if err != nil {
				return err
			}
			var copied, skipped, widened, fromUser []string
			fromFolders := map[string][]string{}
			replaced := false
			for _, n := range names {
				def, ok := user[n]
				var where []string
				if !ok {
					if def, where, err = folderServer(prof.Name, folders, n); err != nil {
						return err
					}
				}
				// The group's definition is what every other account uses.
				if have, ok := store.Servers[n]; ok && !mcp.SameDef(have, def) {
					if !force {
						skipped = append(skipped, n)
						continue
					}
					replaced = true
				}
				store.Servers[n] = def
				copied = append(copied, n)
				if len(where) == 0 {
					fromUser = append(fromUser, n)
					continue
				}
				widened = append(widened, n)
				for _, f := range where {
					fromFolders[f] = append(fromFolders[f], n)
				}
			}
			for _, n := range skipped {
				e.printf("Skipped %s: group %s already has a different %s (`ccam mcp import %s %s --force` replaces it)\n", n, g, n, prof.Name, n)
			}
			if len(copied) == 0 {
				return nil
			}
			if replaced {
				save := filepath.Join(fsutil.UniqueDir(filepath.Join(e.paths.GroupDir(folder), share.BackupDir, fsutil.Stamp(time.Now()))), "mcp.json")
				if err := fsutil.CopyFile(path, save); err != nil {
					return err
				}
				e.printf("The group's previous mcp.json is saved in %s\n", save)
			}
			if err := store.Save(path); err != nil {
				return err
			}
			e.printf("✓ Copied to group %s: %s\n", g, strings.Join(copied, ", "))
			if len(widened) > 0 {
				e.printf("  %s loaded for one folder of account %s only; from the group, every account in it gets them in every folder\n", strings.Join(widened, ", "), prof.Name)
			}
			err = dropOwn(e, prof, dir, keepLocal, fromUser, fromFolders)
			release()
			// The account's own entries are dealt with above.
			refreshCopies(e, cfg, folder, prof.Name)
			return err
		},
	}
	cmd.Flags().BoolVar(&keepLocal, "keep-local", false, "keep the entries in the account's .claude.json")
	cmd.Flags().BoolVar(&force, "force", false, "replace a group server of the same name that differs (the old mcp.json is backed up)")
	return cmd
}

// dropOwn takes servers just copied into the group out of the account's
// .claude.json, where they would win over the group's: user from its user
// scope, folders from those folders' local scope.
func dropOwn(e *env, prof *config.Profile, dir string, keep bool, user []string, folders map[string][]string) error {
	if keep {
		return nil
	}
	if launcher.Running(e.paths, prof.Name, dir) {
		e.printf("Note: account %s is running, so its own entries stay. Run this again after it exits to remove them.\n", prof.Name)
		return nil
	}
	if err := mcp.Remove(dir, filepath.Join(e.paths.State, "mcp", prof.Name), user, folders); err != nil {
		return fmt.Errorf("could not remove the account's entries: %w", err)
	}
	e.printf("✓ Removed account %s's own entries (.claude.json backed up)\n", prof.Name)
	return nil
}

// folderServer finds server n in the local scope of account's folders: its
// definition and the folders that have it. One name with different
// definitions in two folders is refused: which one is meant is not clear.
func folderServer(account string, folders map[string]map[string]json.RawMessage, n string) (json.RawMessage, []string, error) {
	var def json.RawMessage
	var where []string
	for _, f := range slices.Sorted(maps.Keys(folders)) {
		d, ok := folders[f][n]
		if !ok {
			continue
		}
		if def != nil && !mcp.SameDef(def, d) {
			return nil, nil, fmt.Errorf("account %s has different servers named %s for %s and %s; remove one of them (claude mcp remove %s, in that folder), then import the other", account, n, where[0], f, n)
		}
		def = d
		where = append(where, f)
	}
	if def == nil {
		return nil, nil, fmt.Errorf("account %s has no server %q", account, n)
	}
	return def, where, nil
}
