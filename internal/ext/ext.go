// Package ext lists the extensions each account loads: skills, commands,
// agents, plugins, MCP servers and hooks, with where each one comes from
// (a share group, the account's own folder, or the project) and the
// problems that make one not load as expected. It only reads.
package ext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/mcp"
)

// Sources other than a group name.
const (
	Account = "account" // the account's own folder
	Local   = "local"   // the account's .claude.json, for the current folder only (local scope)
	Project = "project" // the .claude folder or .mcp.json of the current project
)

// Item is one extension.
type Item struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Problem string `json:"problem,omitempty"`
	// Off marks a plugin that is installed but not turned on in
	// enabledPlugins, so Claude Code loads nothing from it.
	Off bool `json:"off,omitempty"`
	// Remote marks an MCP server reached by URL (http, sse): one that may
	// need a sign-in, which holds for its account only.
	Remote bool `json:"remote,omitempty"`
}

// Hook counts the hook commands for one event.
type Hook struct {
	Event  string `json:"event"`
	Count  int    `json:"count"`
	Source string `json:"source"`
}

// Set is what one account, or the project, loads.
type Set struct {
	Name     string   `json:"name"`
	Skills   []Item   `json:"skills"`
	Commands []Item   `json:"commands"`
	Agents   []Item   `json:"agents"`
	Plugins  []Item   `json:"plugins"`
	MCP      []Item   `json:"mcp_servers"`
	Hooks    []Hook   `json:"hooks"`
	Problems []string `json:"problems,omitempty"`
	// GroupProblems are problems of items that come from a group, by
	// group: every account in it would report them.
	GroupProblems map[string][]string `json:"-"`
}

// ScanAccount reads the extensions in an account folder. groupsRoot is the
// folder of the share groups, so links into it are named by their group;
// names maps a group folder to the group's name where they differ. project
// is the folder Claude Code would start in ("" for none): its settings can
// turn plugins on or off, plugins installed for a project apply there only,
// and so do the account's local-scope MCP servers of that folder.
func ScanAccount(name, dir, groupsRoot string, share config.Share, names map[string]string, project string) Set {
	s := Set{Name: name}
	label := func(source string) string {
		if n, ok := names[source]; ok {
			return n
		}
		return source
	}
	src := func(entry string) string { return label(sourceOf(dir, entry, groupsRoot)) }
	for _, folder := range []string{"skills", "commands", "agents", "plugins"} {
		path := filepath.Join(dir, folder)
		if target, ok := fsutil.LinkTarget(path); ok && !fsutil.Exists(target) {
			s.Problems = append(s.Problems, folder+": broken link to "+target+" (ccam doctor --fix)")
		}
	}
	s.Skills = skills(filepath.Join(dir, "skills"), src("skills"))
	s.Commands = markdown(filepath.Join(dir, "commands"), src("commands"))
	s.Agents = markdown(filepath.Join(dir, "agents"), src("agents"))
	var defs map[string][]def // MCP servers by name, where each is defined
	s.Plugins, defs = plugins(filepath.Join(dir, "plugins"), src("plugins"), enabledPlugins(dir, project), project)
	s.Hooks = hooks(filepath.Join(dir, "settings.json"), src("settings.json"))

	if g := share.Get("mcp"); g != config.Local {
		if store, err := mcp.Load(mcp.Path(filepath.Join(groupsRoot, g))); err == nil {
			for n, raw := range store.Servers {
				defs[n] = append(defs[n], def{label(g), "group " + label(g) + "'s mcp.json", raw})
			}
		}
	}
	if local, err := mcp.LocalServers(dir); err == nil {
		for n, raw := range local {
			defs[n] = append(defs[n], def{Account, ".claude.json", raw})
		}
	}
	var here []string // local-scope servers of this folder
	if project != "" {
		key := mcp.ProjectKey(project)
		if folders, err := mcp.FolderServers(dir); err == nil {
			if servers, ok := mcp.FolderLookup(folders, key); ok {
				for _, n := range sortedKeys(servers) {
					// Claude Code takes the local scope first, so it goes first.
					defs[n] = append([]def{{Local, ".claude.json's local scope for " + key, servers[n]}}, defs[n]...)
					here = append(here, n)
				}
			}
		}
	}
	var mcpProblems []string
	s.MCP, mcpProblems = mcpItems(defs)
	s.Problems = append(s.Problems, mcpProblems...)
	for _, n := range here {
		s.Problems = append(s.Problems, fmt.Sprintf("MCP server %s is in this account's local scope for this folder (what `claude mcp add` does unless given -s user), so other accounts do not load it: `ccam mcp import %s %s` shares it with the accounts of its group, in every folder", n, name, n))
	}
	for _, list := range [][]Item{s.Skills, s.Commands, s.Agents, s.Plugins} {
		for _, it := range list {
			source := it.Source
			if it.Off {
				// Whether a plugin is on is up to settings.json, which may be
				// the account's own.
				source = src("settings.json")
			}
			switch {
			case it.Problem == "":
			case source == Account:
				s.Problems = append(s.Problems, it.Name+": "+it.Problem)
			default:
				if s.GroupProblems == nil {
					s.GroupProblems = map[string][]string{}
				}
				s.GroupProblems[source] = append(s.GroupProblems[source], it.Name+": "+it.Problem)
			}
		}
	}
	return s
}

// Waiting lists the servers of the project's .mcp.json, of names, that
// wait for approval in the account of dir (Claude Code 2.1.308): the
// account has not trusted the project folder (trust is per account; an
// approval in the project's .claude/settings.local.json counts only where
// it is trusted), and its settings.json neither approves them
// (enableAllProjectMcpServers, enabledMcpjsonServers) nor turns them off
// (disabledMcpjsonServers). Claude Code asks once; -p runs take them.
func Waiting(dir, project string, names []string) []string {
	if len(names) == 0 || mcp.Trusts(dir, project) {
		return nil
	}
	var settings struct {
		All      bool     `json:"enableAllProjectMcpServers"`
		Enabled  []string `json:"enabledMcpjsonServers"`
		Disabled []string `json:"disabledMcpjsonServers"`
	}
	if data, err := os.ReadFile(filepath.Join(dir, "settings.json")); err == nil {
		_ = json.Unmarshal(data, &settings)
	}
	if settings.All {
		return nil
	}
	var out []string
	for _, n := range names {
		if !slices.Contains(settings.Enabled, n) && !slices.Contains(settings.Disabled, n) {
			out = append(out, n)
		}
	}
	return out
}

// Gap is an MCP server that some accounts loading it are signed in to and
// others not: a sign-in (OAuth, with /mcp) holds for its account only.
type Gap struct {
	Server   string
	Missing  []string // accounts that load it without a sign-in
	SignedIn []string // accounts that have one
}

// SignInAccounts lists the accounts of sets whose sign-ins SignInGaps
// needs: those that load a remote server another account loads too.
// Reading sign-ins reads the login (on macOS from the Keychain), so the
// other accounts' are left alone.
func SignInAccounts(sets []Set) []string {
	loaders := map[string][]string{} // remote server → accounts that load it
	for _, s := range sets {
		for _, it := range s.MCP {
			if it.Remote {
				loaders[it.Name] = append(loaders[it.Name], s.Name)
			}
		}
	}
	var out []string
	for _, n := range sortedKeys(loaders) {
		for _, a := range loaders[n] {
			if len(loaders[n]) > 1 && !slices.Contains(out, a) {
				out = append(out, a)
			}
		}
	}
	return out
}

// SignInGaps finds the gaps among sets, given the servers each account is
// signed in to (signedIn, by account; SignInAccounts says whose are needed).
func SignInGaps(sets []Set, signedIn map[string]map[string]bool) []Gap {
	loaders := map[string][]string{} // server → accounts that load it
	for _, s := range sets {
		for _, it := range s.MCP {
			loaders[it.Name] = append(loaders[it.Name], s.Name)
		}
	}
	var out []Gap
	for _, n := range sortedKeys(loaders) {
		g := Gap{Server: n}
		for _, a := range loaders[n] {
			if signedIn[a][n] {
				g.SignedIn = append(g.SignedIn, a)
			} else {
				g.Missing = append(g.Missing, a)
			}
		}
		if len(g.SignedIn) > 0 && len(g.Missing) > 0 {
			out = append(out, g)
		}
	}
	return out
}

// ScanProject reads the project-level extensions of the folder dir.
func ScanProject(dir string) Set {
	s := Set{Name: Project}
	claude := filepath.Join(dir, ".claude")
	s.Skills = skills(filepath.Join(claude, "skills"), Project)
	s.Commands = markdown(filepath.Join(claude, "commands"), Project)
	s.Agents = markdown(filepath.Join(claude, "agents"), Project)
	for _, f := range []string{"settings.json", "settings.local.json"} {
		s.Hooks = append(s.Hooks, hooks(filepath.Join(claude, f), Project)...)
	}
	if store, err := mcp.Load(filepath.Join(dir, ".mcp.json")); err == nil {
		for _, n := range store.Names() {
			s.MCP = append(s.MCP, Item{Name: n, Source: Project})
		}
	}
	return s
}

// Overlaps reports MCP servers an account defines that the project's
// .mcp.json defines too: the two definitions compete for one name.
func Overlaps(account, project Set) []string {
	var out []string
	for _, a := range account.MCP {
		for _, p := range project.MCP {
			if a.Name == p.Name {
				out = append(out, fmt.Sprintf("MCP server %s is defined in the project's .mcp.json and in %s", a.Name, describe(a.Source)))
			}
		}
	}
	return out
}

type def struct {
	source, where string
	raw           json.RawMessage
}

// mcpItems turns server definitions into items. A name defined in more
// than one place with different settings is a problem; the same settings
// twice (mcp_strategy = "sync" copies the group's servers) is not.
func mcpItems(defs map[string][]def) ([]Item, []string) {
	var items []Item
	var problems []string
	for _, n := range sortedKeys(defs) {
		ds := defs[n]
		it := Item{Name: n, Source: ds[0].source, Remote: remote(ds[0].raw)}
		if len(ds) > 1 && !sameJSON(ds[0].raw, ds[1].raw) {
			it.Problem = fmt.Sprintf("defined differently in %s and %s", ds[0].where, ds[1].where)
			problems = append(problems, "MCP server "+n+" is "+it.Problem)
		}
		items = append(items, it)
	}
	return items, problems
}

// remote reports whether a server definition reaches its server by URL.
func remote(def json.RawMessage) bool {
	var d struct {
		URL string `json:"url"`
	}
	return json.Unmarshal(def, &d) == nil && d.URL != ""
}

func sameJSON(a, b json.RawMessage) bool {
	var x, y any
	if json.Unmarshal(a, &x) != nil || json.Unmarshal(b, &y) != nil {
		return bytes.Equal(a, b)
	}
	xa, _ := json.Marshal(x)
	ya, _ := json.Marshal(y)
	return bytes.Equal(xa, ya)
}

func describe(source string) string {
	switch source {
	case Account:
		return "the account's .claude.json"
	case Local:
		return "the account's local scope for this folder"
	}
	return "group " + source + "'s mcp.json"
}

// sourceOf names where dir/name comes from: the group it links into, or
// the account itself.
func sourceOf(dir, name, groupsRoot string) string {
	target, ok := fsutil.LinkTarget(filepath.Join(dir, name))
	if !ok {
		return Account
	}
	if rel, err := filepath.Rel(groupsRoot, target); err == nil && !strings.HasPrefix(rel, "..") {
		if g, _, _ := strings.Cut(filepath.ToSlash(rel), "/"); g != "" && g != "." {
			return g
		}
	}
	return Account
}

// skills lists the skill folders in dir. A skill is a folder (or a link to
// one) with a SKILL.md. skills/synced without one is where Claude Code keeps
// the skills it syncs from claude.ai (a folder per account, each with a
// manifest.json), so it is neither a skill nor a problem.
func skills(dir, source string) []Item {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []Item
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		it := Item{Name: e.Name(), Source: source}
		info, err := os.Stat(path)
		switch {
		case err != nil:
			target, _ := os.Readlink(path)
			it.Problem = "broken link to " + target
		case !info.IsDir():
			continue // a stray file, not a skill
		case !fsutil.Exists(filepath.Join(path, "SKILL.md")) && e.Name() == "synced":
			continue // Claude Code's synced skills, see above
		case !fsutil.Exists(filepath.Join(path, "SKILL.md")):
			it.Problem = "no SKILL.md, so Claude Code does not load it"
		}
		out = append(out, it)
	}
	return out
}

// markdown lists the .md files under dir (commands, agents) by their path
// without the extension; subfolders become "folder/name".
func markdown(dir, source string) []Item {
	// WalkDir does not follow a link at its root, and a shared folder is one.
	if real, err := filepath.EvalSymlinks(dir); err == nil {
		dir = real
	}
	var out []Item
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if path != dir && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		it := Item{Name: strings.TrimSuffix(filepath.ToSlash(rel), ".md"), Source: source}
		if _, err := os.Stat(path); err != nil {
			target, _ := os.Readlink(path)
			it.Problem = "broken link to " + target
		}
		out = append(out, it)
		return nil
	})
	return out
}

// install is one installation of a plugin in installed_plugins.json.
type install struct {
	Scope       string `json:"scope"`
	InstallPath string `json:"installPath"`
	ProjectPath string `json:"projectPath"`
}

// plugins lists the installed plugins (installed_plugins.json) that apply
// in project, and the MCP servers of those turned on, by the names Claude
// Code gives them (plugin:<plugin>:<server>).
func plugins(dir, source string, enabled map[string]bool, project string) ([]Item, map[string][]def) {
	defs := map[string][]def{}
	data, err := os.ReadFile(filepath.Join(dir, "installed_plugins.json"))
	if err != nil {
		return nil, defs
	}
	var reg struct {
		Plugins map[string]json.RawMessage `json:"plugins"`
	}
	if json.Unmarshal(data, &reg) != nil {
		return []Item{{Name: "installed_plugins.json", Source: source, Problem: "cannot be read as JSON"}}, defs
	}
	var out []Item
	for _, n := range sortedKeys(reg.Plugins) {
		var list []install
		if json.Unmarshal(reg.Plugins[n], &list) != nil {
			var one install // version 1: one object per plugin
			_ = json.Unmarshal(reg.Plugins[n], &one)
			list = []install{one}
		}
		in, ok := applies(list, project)
		if !ok {
			continue // installed for another project only
		}
		it := Item{Name: n, Source: source}
		on, named := enabled[n]
		switch {
		case !named:
			// Installing names it; a settings.json kept from elsewhere may not.
			it.Off = true
			it.Problem = "installed but not in enabledPlugins, so Claude Code does not load it (claude plugin enable " + n + ")"
		case !on:
			it.Off = true // turned off on purpose
		case in.InstallPath != "" && !fsutil.Exists(in.InstallPath):
			it.Problem = "its folder " + in.InstallPath + " is missing, so Claude Code does not load it (reinstall it with /plugin)"
		case in.InstallPath != "":
			short, _, _ := strings.Cut(n, "@")
			for srv, raw := range pluginServers(in.InstallPath) {
				key := "plugin:" + short + ":" + srv
				defs[key] = append(defs[key], def{source, "plugin " + n, raw})
			}
		}
		out = append(out, it)
	}
	return out, defs
}

// PluginServer reports whether name is the name Claude Code gives a
// plugin's MCP server (plugin:<plugin>:<server>). Such a server comes with
// its plugin, not from .claude.json, so `ccam mcp import` cannot share it.
func PluginServer(name string) bool { return strings.HasPrefix(name, "plugin:") }

// applies picks the installation of a plugin that counts in project: one
// for the user, or one made for project or a folder above it.
func applies(list []install, project string) (install, bool) {
	for _, in := range list {
		switch in.Scope {
		case "", "user":
			return in, true
		default:
			if project == "" || in.ProjectPath == "" {
				continue
			}
			if rel, err := filepath.Rel(in.ProjectPath, project); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return in, true
			}
		}
	}
	return install{}, false
}

// enabledPlugins is enabledPlugins of the account's settings.json, then of
// the project's .claude/settings.json and settings.local.json: a later file
// wins (Claude Code 2.1.306). A value that is not false counts as on.
func enabledPlugins(dir, project string) map[string]bool {
	files := []string{filepath.Join(dir, "settings.json")}
	if project != "" {
		files = append(files, filepath.Join(project, ".claude", "settings.json"), filepath.Join(project, ".claude", "settings.local.json"))
	}
	out := map[string]bool{}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var doc struct {
			Enabled map[string]json.RawMessage `json:"enabledPlugins"`
		}
		if json.Unmarshal(data, &doc) != nil {
			continue
		}
		for k, v := range doc.Enabled {
			out[k] = strings.TrimSpace(string(v)) != "false"
		}
	}
	return out
}

// pluginServers reads the MCP servers of a plugin's folder: its .mcp.json
// and plugin.json's mcpServers (an object, a file, or a list of them),
// each with or without an mcpServers wrapper. Claude Code loads them all.
func pluginServers(root string) map[string]json.RawMessage {
	out := map[string]json.RawMessage{}
	addFile := func(path string) {
		if data, err := os.ReadFile(path); err == nil {
			addServers(out, data)
		}
	}
	addFile(filepath.Join(root, ".mcp.json"))
	data, err := os.ReadFile(filepath.Join(root, ".claude-plugin", "plugin.json"))
	if err != nil {
		return out
	}
	var manifest struct {
		MCP json.RawMessage `json:"mcpServers"`
	}
	if json.Unmarshal(data, &manifest) != nil || len(manifest.MCP) == 0 {
		return out
	}
	var refs []json.RawMessage
	if json.Unmarshal(manifest.MCP, &refs) != nil {
		refs = []json.RawMessage{manifest.MCP}
	}
	for _, ref := range refs {
		var path string
		if json.Unmarshal(ref, &path) == nil {
			addFile(filepath.Join(root, filepath.FromSlash(path)))
			continue
		}
		addServers(out, ref)
	}
	return out
}

func addServers(out map[string]json.RawMessage, data []byte) {
	var doc map[string]json.RawMessage
	if json.Unmarshal(data, &doc) != nil {
		return
	}
	if wrapped, ok := doc["mcpServers"]; ok {
		doc = nil
		if json.Unmarshal(wrapped, &doc) != nil {
			return
		}
	}
	for n, raw := range doc {
		out[n] = raw
	}
}

// hooks counts the hook commands per event in a settings file.
func hooks(path, source string) []Hook {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var doc struct {
		Hooks map[string][]struct {
			Hooks []json.RawMessage `json:"hooks"`
		} `json:"hooks"`
	}
	if json.Unmarshal(data, &doc) != nil {
		return nil
	}
	var out []Hook
	for _, event := range sortedKeys(doc.Hooks) {
		n := 0
		for _, m := range doc.Hooks[event] {
			n += len(m.Hooks)
		}
		if n > 0 {
			out = append(out, Hook{Event: event, Count: n, Source: source})
		}
	}
	return out
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Names lists the distinct item names across sets, sorted.
func Names(sets []Set, pick func(Set) []Item) []string {
	var out []string
	for _, s := range sets {
		for _, it := range pick(s) {
			if !slices.Contains(out, it.Name) {
				out = append(out, it.Name)
			}
		}
	}
	sort.Strings(out)
	return out
}
