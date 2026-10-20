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
	Project = "project" // the .claude folder or .mcp.json of the current project
)

// Item is one extension.
type Item struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Problem string `json:"problem,omitempty"`
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
}

// ScanAccount reads the extensions in an account folder. groupsRoot is the
// folder of the share groups, so links into it are named by their group.
func ScanAccount(name, dir, groupsRoot string, share config.Share) Set {
	s := Set{Name: name}
	for _, folder := range []string{"skills", "commands", "agents", "plugins"} {
		path := filepath.Join(dir, folder)
		if target, ok := fsutil.LinkTarget(path); ok && !fsutil.Exists(target) {
			s.Problems = append(s.Problems, folder+": broken link to "+target+" (ccam doctor --fix)")
		}
	}
	s.Skills = skills(filepath.Join(dir, "skills"), sourceOf(dir, "skills", groupsRoot))
	s.Commands = markdown(filepath.Join(dir, "commands"), sourceOf(dir, "commands", groupsRoot))
	s.Agents = markdown(filepath.Join(dir, "agents"), sourceOf(dir, "agents", groupsRoot))
	s.Plugins = plugins(filepath.Join(dir, "plugins"), sourceOf(dir, "plugins", groupsRoot))
	s.Hooks = hooks(filepath.Join(dir, "settings.json"), sourceOf(dir, "settings.json", groupsRoot))

	defs := map[string][]def{}
	if g := share.Get("mcp"); g != config.Local {
		if store, err := mcp.Load(mcp.Path(filepath.Join(groupsRoot, g))); err == nil {
			for n, raw := range store.Servers {
				defs[n] = append(defs[n], def{g, "group " + g + "'s mcp.json", raw})
			}
		}
	}
	if local, err := mcp.LocalServers(dir); err == nil {
		for n, raw := range local {
			defs[n] = append(defs[n], def{Account, ".claude.json", raw})
		}
	}
	var mcpProblems []string
	s.MCP, mcpProblems = mcpItems(defs)
	s.Problems = append(s.Problems, mcpProblems...)
	for _, list := range [][]Item{s.Skills, s.Commands, s.Agents, s.Plugins} {
		for _, it := range list {
			if it.Problem != "" {
				s.Problems = append(s.Problems, it.Name+": "+it.Problem)
			}
		}
	}
	return s
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
		it := Item{Name: n, Source: ds[0].source}
		if len(ds) > 1 && !sameJSON(ds[0].raw, ds[1].raw) {
			it.Problem = fmt.Sprintf("defined differently in %s and %s", ds[0].where, ds[1].where)
			problems = append(problems, "MCP server "+n+" is "+it.Problem)
		}
		items = append(items, it)
	}
	return items, problems
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
	if source == Account {
		return "the account's .claude.json"
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
// one) with a SKILL.md.
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

// plugins lists the installed plugins (installed_plugins.json).
func plugins(dir, source string) []Item {
	data, err := os.ReadFile(filepath.Join(dir, "installed_plugins.json"))
	if err != nil {
		return nil
	}
	var reg struct {
		Plugins map[string]json.RawMessage `json:"plugins"`
	}
	if json.Unmarshal(data, &reg) != nil {
		return []Item{{Name: "installed_plugins.json", Source: source, Problem: "cannot be read as JSON"}}
	}
	var out []Item
	for _, n := range sortedKeys(reg.Plugins) {
		out = append(out, Item{Name: n, Source: source})
	}
	return out
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
