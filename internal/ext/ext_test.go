package ext

import (
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
	write(t, filepath.Join(dir, "settings.json"), `{"hooks":{"PreToolUse":[{"matcher":"Bash","hooks":[{"type":"command","command":"a"},{"type":"command","command":"b"}]}],"Stop":[{"hooks":[{"type":"command","command":"c"}]}]}}`)
	// MCP: the group's servers and the account's own; fs differs, ctx is the same in both.
	write(t, filepath.Join(group, "mcp.json"), `{"mcpServers":{"ctx":{"type":"http","url":"https://ctx"},"fs":{"command":"fs-server"}}}`)
	write(t, filepath.Join(dir, ".claude.json"), `{"mcpServers":{"ctx":{"url":"https://ctx","type":"http"},"fs":{"command":"other-fs"},"own":{"command":"x"}}}`)

	s := ScanAccount("work", dir, groups, config.Share{MCP: "default", Extensions: "default"}, nil)
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
	s := ScanAccount("a", dir, t.TempDir(), config.Share{}, nil)
	if len(s.Problems) != 1 || !strings.Contains(s.Problems[0], "skills: broken link") {
		t.Fatalf("%v", s.Problems)
	}
}

func TestEmptyAccount(t *testing.T) {
	s := ScanAccount("new", t.TempDir(), t.TempDir(), config.Share{}, nil)
	if len(s.Skills)+len(s.Commands)+len(s.Agents)+len(s.Plugins)+len(s.MCP)+len(s.Hooks)+len(s.Problems) != 0 {
		t.Fatalf("%+v", s)
	}
}
