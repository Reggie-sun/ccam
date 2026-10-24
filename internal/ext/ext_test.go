package ext

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	os.MkdirAll(filepath.Dir(path), 0o700)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func names(items []Item) []string {
	var out []string
	for _, it := range items {
		out = append(out, it.Name+"@"+it.Source)
	}
	return out
}

func TestScanAccount(t *testing.T) {
	root := t.TempDir()
	groups := filepath.Join(root, "groups")
	group := filepath.Join(groups, "default")
	dir := filepath.Join(root, "profiles", "work")

	// Skills shared through the group: one fine, one without SKILL.md.
	write(t, filepath.Join(group, "skills", "review", "SKILL.md"), "---\nname: review\n---\n")
	os.MkdirAll(filepath.Join(group, "skills", "empty"), 0o700)
	// Claude Code's own store of skills synced from claude.ai: not a skill
	// and not a problem.
	write(t, filepath.Join(group, "skills", "synced", "org_account", "pdf", "SKILL.md"), "---\nname: pdf\n---\n")
	write(t, filepath.Join(group, "skills", "synced", "org_account", "manifest.json"), "{}")
	os.MkdirAll(dir, 0o700)
	if err := os.Symlink(filepath.Join(group, "skills"), filepath.Join(dir, "skills")); err != nil {
		t.Skip("no symlinks:", err)
	}
	// The account's own commands, one a broken link, one in a subfolder.
	write(t, filepath.Join(dir, "commands", "deploy.md"), "deploy")
	write(t, filepath.Join(dir, "commands", "git", "pr.md"), "pr")
	os.Symlink(filepath.Join(root, "gone.md"), filepath.Join(dir, "commands", "old.md"))
	// Agents shared through the group: the folder itself is a link.
	write(t, filepath.Join(group, "agents", "reviewer.md"), "agent")
	os.Symlink(filepath.Join(group, "agents"), filepath.Join(dir, "agents"))
	write(t, filepath.Join(dir, "plugins", "installed_plugins.json"), `{"version":2,"plugins":{"fmt@mkt":[{}]}}`)
	write(t, filepath.Join(dir, "settings.json"), `{"enabledPlugins":{"fmt@mkt":true},"hooks":{"PreToolUse":[{"matcher":"Bash","hooks":[{"type":"command","command":"a"},{"type":"command","command":"b"}]}],"Stop":[{"hooks":[{"type":"command","command":"c"}]}]}}`)
	// MCP: the group's servers and the account's own; fs differs, ctx is the same in both.
	write(t, filepath.Join(group, "mcp.json"), `{"mcpServers":{"ctx":{"type":"http","url":"https://ctx"},"fs":{"command":"fs-server"}}}`)
	write(t, filepath.Join(dir, ".claude.json"), `{"mcpServers":{"ctx":{"url":"https://ctx","type":"http"},"fs":{"command":"other-fs"},"own":{"command":"x"}}}`)

	s := ScanAccount("work", dir, groups, config.Share{MCP: "default", Extensions: "default"}, nil, "")
	if got := names(s.Skills); !slices.Equal(got, []string{"empty@default", "review@default"}) {
		t.Errorf("skills %v", got)
	}
	if got := names(s.Commands); !slices.Equal(got, []string{"deploy@account", "git/pr@account", "old@account"}) {
		t.Errorf("commands %v", got)
	}
	if got := names(s.Agents); !slices.Equal(got, []string{"reviewer@default"}) {
		t.Errorf("agents through a linked folder %v", got)
	}
	if got := names(s.Plugins); !slices.Equal(got, []string{"fmt@mkt@account"}) {
		t.Errorf("plugins %v", got)
	}
	if got := names(s.MCP); !slices.Equal(got, []string{"ctx@default", "fs@default", "own@account"}) {
		t.Errorf("mcp %v", got)
	}
	if len(s.Hooks) != 2 || s.Hooks[0].Event != "PreToolUse" || s.Hooks[0].Count != 2 || s.Hooks[1].Count != 1 {
		t.Errorf("hooks %+v", s.Hooks)
	}
	// A problem of a group's skill is the group's, reported once for all.
	if got := s.GroupProblems["default"]; len(got) != 1 || !strings.HasPrefix(got[0], "empty: no SKILL.md") {
		t.Errorf("group problems %v", s.GroupProblems)
	}
	problems := strings.Join(s.Problems, "\n")
	for _, want := range []string{"old: broken link to", "MCP server fs is defined differently in group default's mcp.json and .claude.json"} {
		if !strings.Contains(problems, want) {
			t.Errorf("missing problem %q in:\n%s", want, problems)
		}
	}
	if strings.Contains(problems, "ctx") {
		t.Errorf("the same server twice is no problem:\n%s", problems)
	}
}

func TestScanProjectAndOverlaps(t *testing.T) {
	proj := t.TempDir()
	write(t, filepath.Join(proj, ".claude", "skills", "lint", "SKILL.md"), "x")
	write(t, filepath.Join(proj, ".claude", "settings.local.json"), `{"hooks":{"PostToolUse":[{"hooks":[{"type":"command","command":"fmt"}]}]}}`)
	write(t, filepath.Join(proj, ".mcp.json"), `{"mcpServers":{"fs":{"command":"proj-fs"}}}`)
	p := ScanProject(proj)
	if got := names(p.Skills); !slices.Equal(got, []string{"lint@project"}) {
		t.Errorf("skills %v", got)
	}
	if len(p.Hooks) != 1 || p.Hooks[0].Event != "PostToolUse" || names(p.MCP)[0] != "fs@project" {
		t.Errorf("%+v", p)
	}
	acct := Set{MCP: []Item{{Name: "fs", Source: "default"}, {Name: "own", Source: Account}}}
	if got := Overlaps(acct, p); len(got) != 1 || !strings.Contains(got[0], "group default's mcp.json") {
		t.Errorf("overlaps %v", got)
	}
}

func TestBrokenFolderLink(t *testing.T) {
	dir := t.TempDir()
	if err := os.Symlink(filepath.Join(dir, "gone"), filepath.Join(dir, "skills")); err != nil {
		t.Skip(err)
	}
	s := ScanAccount("a", dir, t.TempDir(), config.Share{}, nil, "")
	if len(s.Problems) != 1 || !strings.Contains(s.Problems[0], "skills: broken link") {
		t.Fatalf("%v", s.Problems)
	}
}

func TestEmptyAccount(t *testing.T) {
	s := ScanAccount("new", t.TempDir(), t.TempDir(), config.Share{}, nil, "")
	if len(s.Skills)+len(s.Commands)+len(s.Agents)+len(s.Plugins)+len(s.MCP)+len(s.Hooks)+len(s.Problems) != 0 {
		t.Fatalf("%+v", s)
	}
}

// Plugins load only when enabledPlugins turns them on (the account's
// settings, then the project's, then its settings.local.json), and bring
// their MCP servers as plugin:<plugin>:<server>: from .mcp.json, flat or
// wrapped, and from plugin.json, inline or as a file (Claude Code 2.1.306).
func TestPluginsFollowEnabledPlugins(t *testing.T) {
	root := t.TempDir()
	groups := filepath.Join(root, "groups")
	plugins := filepath.Join(groups, "default", "plugins")
	dir := filepath.Join(root, "profiles", "work")
	project := filepath.Join(root, "proj")
	cache := filepath.Join(plugins, "cache", "mk")
	write(t, filepath.Join(cache, "on", ".mcp.json"), `{"flat":{"command":"a"}}`)
	write(t, filepath.Join(cache, "on", ".claude-plugin", "plugin.json"), `{"name":"on","mcpServers":{"inline":{"command":"b"}}}`)
	write(t, filepath.Join(cache, "path", ".claude-plugin", "plugin.json"), `{"name":"path","mcpServers":"./servers.json"}`)
	write(t, filepath.Join(cache, "path", "servers.json"), `{"mcpServers":{"nested":{"type":"http","url":"https://x"}}}`)
	write(t, filepath.Join(cache, "off", ".mcp.json"), `{"offsrv":{"command":"c"}}`)
	reg := `{"version":2,"plugins":{` +
		`"on@mk":[{"scope":"user","installPath":"` + filepath.ToSlash(filepath.Join(cache, "on")) + `"}],` +
		`"path@mk":[{"scope":"user","installPath":"` + filepath.ToSlash(filepath.Join(cache, "path")) + `"}],` +
		`"off@mk":[{"scope":"user","installPath":"` + filepath.ToSlash(filepath.Join(cache, "off")) + `"}],` +
		`"unnamed@mk":[{"scope":"user","installPath":"` + filepath.ToSlash(filepath.Join(cache, "on")) + `"}],` +
		`"gone@mk":[{"scope":"user","installPath":"` + filepath.ToSlash(filepath.Join(cache, "gone")) + `"}],` +
		`"elsewhere@mk":[{"scope":"project","projectPath":"` + filepath.ToSlash(filepath.Join(root, "other")) + `","installPath":"` + filepath.ToSlash(filepath.Join(cache, "on")) + `"}],` +
		`"here@mk":[{"scope":"local","projectPath":"` + filepath.ToSlash(project) + `","installPath":"` + filepath.ToSlash(filepath.Join(cache, "off")) + `"}]}}`
	write(t, filepath.Join(plugins, "installed_plugins.json"), reg)
	os.MkdirAll(dir, 0o700)
	if err := os.Symlink(plugins, filepath.Join(dir, "plugins")); err != nil {
		t.Skip("no symlinks:", err)
	}
	write(t, filepath.Join(dir, "settings.json"), `{"enabledPlugins":{"on@mk":true,"path@mk":true,"off@mk":false,"gone@mk":true,"here@mk":false,"elsewhere@mk":true}}`)
	write(t, filepath.Join(project, ".claude", "settings.json"), `{"enabledPlugins":{"here@mk":true,"path@mk":false}}`)
	write(t, filepath.Join(project, ".claude", "settings.local.json"), `{"enabledPlugins":{"path@mk":true}}`)

	s := ScanAccount("work", dir, groups, config.Share{Extensions: "default"}, nil, project)
	var got []string
	for _, it := range s.Plugins {
		got = append(got, fmt.Sprintf("%s off=%v problem=%v", it.Name, it.Off, it.Problem != ""))
	}
	want := []string{"gone@mk off=false problem=true", "here@mk off=false problem=false", "off@mk off=true problem=false",
		"on@mk off=false problem=false", "path@mk off=false problem=false", "unnamed@mk off=true problem=true"}
	if !slices.Equal(got, want) {
		t.Errorf("plugins\n got %q\nwant %q", got, want)
	}
	if got := names(s.MCP); !slices.Equal(got, []string{"plugin:here:offsrv@default", "plugin:on:flat@default", "plugin:on:inline@default", "plugin:path:nested@default"}) {
		t.Errorf("mcp %v", got)
	}
	// Whether a plugin is on is up to settings.json, here the account's own;
	// a missing folder is a problem of the group the plugins come from.
	if p := strings.Join(s.Problems, "\n"); !strings.Contains(p, "unnamed@mk: installed but not in enabledPlugins") || strings.Contains(p, "gone@mk") {
		t.Errorf("account problems %v", s.Problems)
	}
	if g := s.GroupProblems["default"]; len(g) != 1 || !strings.HasPrefix(g[0], "gone@mk: its folder") {
		t.Errorf("group problems %v", s.GroupProblems)
	}
	// Outside the project, the project's settings and plugins do not count.
	s = ScanAccount("work", dir, groups, config.Share{Extensions: "default"}, nil, "")
	if got := names(s.MCP); !slices.Equal(got, []string{"plugin:on:flat@default", "plugin:on:inline@default", "plugin:path:nested@default"}) {
		t.Errorf("mcp without a project %v", got)
	}
}

// The account's local-scope servers of the folder (keyed as Claude Code
// keys it) are listed as local and said to be the account's alone.
func TestLocalScopeServersOfThisFolder(t *testing.T) {
	root, _ := filepath.EvalSymlinks(t.TempDir())
	repo := filepath.Join(root, "repo")
	os.MkdirAll(filepath.Join(repo, ".git"), 0o700)
	os.MkdirAll(filepath.Join(repo, "sub"), 0o700)
	dir := filepath.Join(root, "profiles", "work")
	write(t, filepath.Join(dir, ".claude.json"), `{"mcpServers":{"user":{"command":"u"}},"projects":{"`+filepath.ToSlash(repo)+`":{"mcpServers":{"loc":{"command":"l"},"user":{"command":"other"}}},"/elsewhere":{"mcpServers":{"far":{"command":"f"}}}}}`)
	s := ScanAccount("work", dir, filepath.Join(root, "groups"), config.Share{}, nil, filepath.Join(repo, "sub"))
	if got := names(s.MCP); !slices.Equal(got, []string{"loc@local", "user@local"}) {
		t.Fatalf("mcp %v", got)
	}
	p := strings.Join(s.Problems, "\n")
	for _, want := range []string{"MCP server loc is in this account's local scope for this folder", "`ccam mcp import work loc`", "MCP server user is defined differently"} {
		if !strings.Contains(p, want) {
			t.Errorf("missing %q in\n%s", want, p)
		}
	}
}

// Project servers wait in an account that has not trusted the folder,
// unless its settings.json approves them (Claude Code 2.1.308).
func TestWaiting(t *testing.T) {
	project, _ := filepath.EvalSymlinks(t.TempDir())
	dir := t.TempDir()
	names := []string{"a", "b"}
	if got := Waiting(dir, project, names); !slices.Equal(got, names) {
		t.Fatalf("untrusted: %v", got)
	}
	write(t, filepath.Join(dir, "settings.json"), `{"enabledMcpjsonServers":["a"],"disabledMcpjsonServers":["b"]}`)
	if got := Waiting(dir, project, names); len(got) != 0 {
		t.Fatalf("approved and turned off: %v", got)
	}
	write(t, filepath.Join(dir, "settings.json"), `{"enableAllProjectMcpServers":true}`)
	if got := Waiting(dir, project, names); len(got) != 0 {
		t.Fatalf("all approved: %v", got)
	}
	os.Remove(filepath.Join(dir, "settings.json"))
	write(t, filepath.Join(dir, ".claude.json"), `{"projects":{"`+filepath.ToSlash(project)+`":{"hasTrustDialogAccepted":true}}}`)
	if got := Waiting(dir, project, names); len(got) != 0 {
		t.Fatalf("trusted: %v", got)
	}
}

func TestSignInGaps(t *testing.T) {
	sets := []Set{
		{Name: "c1", MCP: []Item{{Name: "ctx"}, {Name: "gh"}}},
		{Name: "c3", MCP: []Item{{Name: "ctx"}, {Name: "gh"}}},
		{Name: "c4", MCP: []Item{{Name: "ctx"}}},
	}
	got := SignInGaps(sets, map[string]map[string]bool{"c3": {"ctx": true}, "c4": {"ctx": true}, "c1": {"other": true}})
	if len(got) != 1 || got[0].Server != "ctx" || !slices.Equal(got[0].Missing, []string{"c1"}) || !slices.Equal(got[0].SignedIn, []string{"c3", "c4"}) {
		t.Fatalf("%+v", got)
	}
}

// Sign-ins are read only for accounts that load a remote server another
// account loads too: only such a server can have a gap.
func TestSignInAccounts(t *testing.T) {
	sets := []Set{
		{Name: "c1", MCP: []Item{{Name: "ctx", Remote: true}, {Name: "fs"}}},
		{Name: "c2", MCP: []Item{{Name: "fs"}, {Name: "solo", Remote: true}}},
		{Name: "c3", MCP: []Item{{Name: "ctx", Remote: true}}},
	}
	if got := SignInAccounts(sets); !slices.Equal(got, []string{"c1", "c3"}) {
		t.Fatalf("%v", got)
	}
	if got := SignInAccounts(sets[1:2]); len(got) != 0 {
		t.Fatalf("%v", got)
	}
	// ScanAccount marks a server reached by URL.
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"mcpServers":{"r":{"type":"http","url":"https://r"},"s":{"command":"s"}}}`), 0o600)
	s := ScanAccount("a", dir, t.TempDir(), config.Share{}, nil, "")
	for _, it := range s.MCP {
		if it.Remote != (it.Name == "r") {
			t.Errorf("%+v", it)
		}
	}
}
