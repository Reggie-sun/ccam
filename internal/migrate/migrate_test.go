package migrate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/share"
)

type fixture struct {
	home, src, groups, group string
}

// newFixture lays out home/.claude (the source) and an empty groups root.
// files are relative to home, so ".claude.json" is ~/.claude.json.
func newFixture(t *testing.T, files map[string]string) fixture {
	t.Helper()
	root := t.TempDir()
	f := fixture{home: filepath.Join(root, "home"), groups: filepath.Join(root, "groups")}
	f.src = filepath.Join(f.home, ".claude")
	f.group = filepath.Join(f.groups, "default")
	os.MkdirAll(f.src, 0o700)
	write(t, f.home, files)
	return f
}

func write(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		path := filepath.Join(dir, rel)
		os.MkdirAll(filepath.Dir(path), 0o700)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func (f fixture) build(t *testing.T, items ...string) *Plan {
	t.Helper()
	if len(items) == 0 {
		items = DefaultItems
	}
	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: items})
	if err != nil {
		t.Fatal(err)
	}
	return pl
}

func (f fixture) migrate(t *testing.T, items ...string) *Plan {
	t.Helper()
	pl := f.build(t, items...)
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	return pl
}

// snapshot records content hash, mode and mtime of every path under the
// roots, so a test can prove the source was never touched.
func snapshot(t *testing.T, roots ...string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, root := range roots {
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				t.Fatal(err)
			}
			info, err := os.Lstat(path)
			if err != nil {
				t.Fatal(err)
			}
			v := fmt.Sprintf("%v %d", info.Mode(), info.ModTime().UnixNano())
			switch {
			case info.Mode()&os.ModeSymlink != 0:
				link, _ := os.Readlink(path)
				v += " -> " + link
			case info.Mode().IsRegular():
				data, _ := os.ReadFile(path)
				sum := sha256.Sum256(data)
				v += " " + hex.EncodeToString(sum[:])
			}
			out[path] = v
			return nil
		})
	}
	return out
}

func noFileNamed(t *testing.T, root, name string) {
	t.Helper()
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err == nil && d.Name() == name {
			t.Errorf("%s must never be copied: %s", name, path)
		}
		return nil
	})
}

var fullSource = map[string]string{
	".claude/settings.json":                   `{"model":"opus"}`,
	".claude/CLAUDE.md":                       "# my rules",
	".claude/skills/review/SKILL.md":          "skill",
	".claude/commands/fix.md":                 "fix it",
	".claude/projects/p/s1.jsonl":             "{}\n",
	".claude/history.jsonl":                   "{\"display\":\"hi\"}\n",
	".claude/.credentials.json":               `{"claudeAiOauth":{"accessToken":"secret"}}`,
	".claude/skills/review/.credentials.json": "secret",
	".claude.json": `{"oauthAccount":{"emailAddress":"me@example.com"},` +
		`"mcpServers":{"fs":{"type":"stdio","command":"npx","args":["fs"]},"ctx":{"type":"http","url":"https://mcp.example/mcp"}},` +
		`"projects":{"/code":{"allowedTools":[]}}}`,
}

func TestCopyIntoEmptyGroup(t *testing.T) {
	f := newFixture(t, fullSource)
	before := snapshot(t, f.home)
	pl := f.migrate(t)

	if read(t, filepath.Join(f.group, "settings.json")) != `{"model":"opus"}` ||
		read(t, filepath.Join(f.group, "CLAUDE.md")) != "# my rules" ||
		read(t, filepath.Join(f.group, "skills", "review", "SKILL.md")) != "skill" ||
		read(t, filepath.Join(f.group, "commands", "fix.md")) != "fix it" {
		t.Fatal("content not copied")
	}
	if info, _ := os.Lstat(filepath.Join(f.group, "settings.json")); !info.Mode().IsRegular() {
		t.Fatal("the group gets a real file, not a link to the source")
	}
	for _, name := range []string{"projects", "history.jsonl", ".claude.json"} {
		if _, err := os.Lstat(filepath.Join(f.group, name)); err == nil {
			t.Errorf("%s copied without being asked for", name)
		}
	}
	noFileNamed(t, f.groups, ".credentials.json")
	var store struct {
		Servers map[string]map[string]any `json:"mcpServers"`
	}
	if err := json.Unmarshal([]byte(read(t, filepath.Join(f.group, "mcp.json"))), &store); err != nil {
		t.Fatal(err)
	}
	if store.Servers["fs"]["command"] != "npx" || store.Servers["ctx"]["url"] != "https://mcp.example/mcp" || len(store.Servers) != 2 {
		t.Fatalf("mcp.json: %v", store.Servers)
	}
	if got := pl.Copied(); !reflect.DeepEqual(got, []string{"2 MCP servers", "settings.json", "CLAUDE.md", "commands/", "skills/"}) {
		t.Fatalf("Copied() = %q", got)
	}
	if _, err := os.Stat(pl.BackupDir); err == nil {
		t.Fatal("nothing was replaced, so there is nothing to back up")
	}
	if after := snapshot(t, f.home); !reflect.DeepEqual(before, after) {
		t.Fatal("the source changed")
	}
}

func TestIdenticalContentIsNoOp(t *testing.T) {
	f := newFixture(t, fullSource)
	f.migrate(t)
	group := snapshot(t, f.groups)
	pl := f.build(t)
	if !pl.Empty() {
		var b bytes.Buffer
		pl.Describe(&b, true)
		t.Fatalf("second run should have nothing to do:\n%s", b.String())
	}
	for _, en := range pl.Entries {
		if en.State != Same && en.State != Missing {
			t.Errorf("%s: state %d", en.Entry.Name, en.State)
		}
	}
	if !reflect.DeepEqual(pl.MCP.Same, []string{"ctx", "fs"}) {
		t.Fatalf("MCP same: %v", pl.MCP.Same)
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(group, snapshot(t, f.groups)) {
		t.Fatal("a no-op run changed the group")
	}
}

// conflictFixture has a group whose settings, CLAUDE.md and fs server
// differ from the source.
func conflictFixture(t *testing.T) fixture {
	f := newFixture(t, fullSource)
	write(t, f.group, map[string]string{
		"settings.json": `{"theme":"dark"}`,
		"CLAUDE.md":     "# group rules",
		"mcp.json":      `{"mcpServers":{"fs":{"type":"stdio","command":"uvx","args":["fs"]}}}`,
	})
	return f
}

func TestConflictsDefaultToKeepGroup(t *testing.T) {
	f := conflictFixture(t)
	before := snapshot(t, f.home)
	pl := f.build(t)
	if len(pl.Conflicts) != 3 {
		t.Fatalf("conflicts: %d", len(pl.Conflicts))
	}
	// --yes forces this; it is also every non-history default.
	for _, c := range pl.Conflicts {
		c.Choice = share.KeepGroup
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	if read(t, filepath.Join(f.group, "settings.json")) != `{"theme":"dark"}` || read(t, filepath.Join(f.group, "CLAUDE.md")) != "# group rules" {
		t.Fatal("keep group must leave the group's files alone")
	}
	mcpJSON := read(t, filepath.Join(f.group, "mcp.json"))
	if !strings.Contains(mcpJSON, "uvx") || !strings.Contains(mcpJSON, "ctx") {
		t.Fatalf("fs kept, ctx added: %s", mcpJSON)
	}
	if got := pl.Kept(); !reflect.DeepEqual(got, []string{"MCP server fs", "settings.json", "CLAUDE.md"}) {
		t.Fatalf("Kept() = %q", got)
	}
	if !reflect.DeepEqual(before, snapshot(t, f.home)) {
		t.Fatal("the source changed")
	}
}

func TestConflictUseMineAndMergeBackUpTheGroup(t *testing.T) {
	f := conflictFixture(t)
	before := snapshot(t, f.home)
	pl := f.build(t)
	for _, c := range pl.Conflicts {
		switch c.Entry.Name {
		case "settings.json":
			c.Choice = share.MergeJSON
		default:
			c.Choice = share.UseMine
		}
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	var settings map[string]any
	json.Unmarshal([]byte(read(t, filepath.Join(f.group, "settings.json"))), &settings)
	if settings["theme"] != "dark" || settings["model"] != "opus" {
		t.Fatalf("merged settings: %v", settings)
	}
	if read(t, filepath.Join(f.group, "CLAUDE.md")) != "# my rules" {
		t.Fatal("use mine must write the source's CLAUDE.md")
	}
	if !strings.Contains(read(t, filepath.Join(f.group, "mcp.json")), `"npx"`) {
		t.Fatal("use mine must replace the fs server")
	}
	for name, want := range map[string]string{
		"settings.json": `{"theme":"dark"}`,
		"CLAUDE.md":     "# group rules",
	} {
		if got := read(t, filepath.Join(pl.BackupDir, name)); got != want {
			t.Errorf("backup %s = %q", name, got)
		}
	}
	if !strings.Contains(read(t, filepath.Join(pl.BackupDir, "mcp.json")), "uvx") {
		t.Error("mcp.json backup missing")
	}
	if !strings.HasPrefix(pl.BackupDir, filepath.Join(f.group, share.BackupDir)+string(filepath.Separator)) {
		t.Errorf("backup dir %s", pl.BackupDir)
	}
	if !reflect.DeepEqual(before, snapshot(t, f.home)) {
		t.Fatal("the source changed")
	}
}

func TestBackupsNeverOverwriteEachOther(t *testing.T) {
	f := conflictFixture(t)
	var dirs []string
	for _, mine := range []string{`{"v":1}`, `{"v":2}`} {
		write(t, f.home, map[string]string{".claude/settings.json": mine})
		pl := f.build(t, "settings")
		pl.Conflicts[0].Choice = share.UseMine
		if err := pl.Apply(); err != nil {
			t.Fatal(err)
		}
		dirs = append(dirs, pl.BackupDir)
	}
	if dirs[0] == dirs[1] {
		t.Fatalf("two runs share %s", dirs[0])
	}
	if read(t, filepath.Join(dirs[0], "settings.json")) != `{"theme":"dark"}` || read(t, filepath.Join(dirs[1], "settings.json")) != `{"v":1}` {
		t.Fatal("each run keeps the group copy it replaced")
	}
}

func TestMissingItemsAreSkipped(t *testing.T) {
	f := newFixture(t, map[string]string{".claude/CLAUDE.md": "only memory"})
	pl := f.migrate(t)
	for _, en := range pl.Entries {
		if en.Entry.Name != "CLAUDE.md" && en.State != Missing {
			t.Errorf("%s: state %d, want Missing", en.Entry.Name, en.State)
		}
	}
	if pl.MCP == nil || pl.MCP.From != "" {
		t.Fatalf("no .claude.json anywhere: %+v", pl.MCP)
	}
	entries, _ := os.ReadDir(f.group)
	if len(entries) != 1 || entries[0].Name() != "CLAUDE.md" {
		t.Fatalf("group: %v", entries)
	}
	var b bytes.Buffer
	pl.Describe(&b, true)
	if !strings.Contains(b.String(), "not found, skipped: settings.json") ||
		!strings.Contains(b.String(), "not found, skipped: agents/, commands/, skills/, plugins/, output-styles/, rules/") {
		t.Fatalf("plan:\n%s", b.String())
	}
}

func TestHistoryIsOptIn(t *testing.T) {
	f := newFixture(t, fullSource)
	pl := f.migrate(t, "history")
	if read(t, filepath.Join(f.group, "projects", "p", "s1.jsonl")) != "{}\n" || read(t, filepath.Join(f.group, "history.jsonl")) != "{\"display\":\"hi\"}\n" {
		t.Fatal("history not copied")
	}
	if _, err := os.Lstat(filepath.Join(f.group, "settings.json")); err == nil {
		t.Fatal("only the named item is copied")
	}
	if pl.MCP != nil {
		t.Fatal("mcp was not asked for")
	}

	// A group with its own history: new sessions are added, prompt history
	// defaults to merging by line.
	g := newFixture(t, fullSource)
	write(t, g.group, map[string]string{"projects/q/s2.jsonl": "group\n", "history.jsonl": "{\"display\":\"group\"}\n"})
	pl = g.build(t, "history")
	if len(pl.Conflicts) != 1 || pl.Conflicts[0].Choice != share.MergeLines {
		t.Fatalf("conflicts: %+v", pl.Conflicts)
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	if read(t, filepath.Join(g.group, "history.jsonl")) != "{\"display\":\"group\"}\n{\"display\":\"hi\"}\n" {
		t.Fatalf("history.jsonl: %q", read(t, filepath.Join(g.group, "history.jsonl")))
	}
	if read(t, filepath.Join(g.group, "projects", "q", "s2.jsonl")) != "group\n" || read(t, filepath.Join(g.group, "projects", "p", "s1.jsonl")) != "{}\n" {
		t.Fatal("projects not merged")
	}
}

func TestDirMergeKeepsGroupFiles(t *testing.T) {
	f := newFixture(t, map[string]string{
		".claude/skills/a/SKILL.md": "mine a",
		".claude/skills/b/SKILL.md": "mine b",
		".claude/skills/c.md":       "same",
	})
	write(t, f.group, map[string]string{"skills/a/SKILL.md": "group a", "skills/c.md": "same"})
	pl := f.build(t, "extensions")
	var skills *Entry
	for _, en := range pl.Entries {
		if en.Entry.Name == "skills" {
			skills = en
		}
	}
	if skills.State != Merge || skills.Files != 1 || skills.Differ != 1 {
		t.Fatalf("skills: %+v", skills)
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	if read(t, filepath.Join(f.group, "skills", "a", "SKILL.md")) != "group a" || read(t, filepath.Join(f.group, "skills", "b", "SKILL.md")) != "mine b" {
		t.Fatal("merge must add missing files and keep the group's")
	}
}

func TestCredentialsNeverCopied(t *testing.T) {
	f := newFixture(t, fullSource)
	write(t, f.home, map[string]string{".claude/plugins/.credentials.json": "x", ".claude/projects/.credentials.json": "x"})
	f.migrate(t, "settings", "memory", "extensions", "mcp", "history")
	noFileNamed(t, f.groups, ".credentials.json")
	if _, err := os.Lstat(filepath.Join(f.group, ".claude.json")); err == nil {
		t.Fatal(".claude.json copied")
	}
}

func TestLinkedSourceIsSkippedAndDotfileLinksAreCopied(t *testing.T) {
	f := newFixture(t, map[string]string{"dotfiles/CLAUDE.md": "from dotfiles"})
	write(t, f.group, map[string]string{"settings.json": `{"a":1}`, "skills/x.md": "x"})
	// settings.json and skills/ already link into the group (an imported
	// ~/.claude that shares them); CLAUDE.md links to a dotfiles repo.
	os.Symlink(filepath.Join(f.group, "settings.json"), filepath.Join(f.src, "settings.json"))
	os.Symlink(filepath.Join(f.group, "skills"), filepath.Join(f.src, "skills"))
	os.Symlink(filepath.Join(f.home, "dotfiles", "CLAUDE.md"), filepath.Join(f.src, "CLAUDE.md"))
	before := snapshot(t, f.home)
	pl := f.migrate(t, "settings", "memory", "extensions")
	states := map[string]State{}
	for _, en := range pl.Entries {
		states[en.Entry.Name] = en.State
	}
	if states["settings.json"] != Linked || states["skills"] != Linked || states["CLAUDE.md"] != New {
		t.Fatalf("states: %v", states)
	}
	info, err := os.Lstat(filepath.Join(f.group, "CLAUDE.md"))
	if err != nil || !info.Mode().IsRegular() || read(t, filepath.Join(f.group, "CLAUDE.md")) != "from dotfiles" {
		t.Fatal("a dotfiles link's content must be copied as a real file")
	}
	if !reflect.DeepEqual(before, snapshot(t, f.home)) {
		t.Fatal("the source changed")
	}
}

func TestPluginPathsPointAtTheGroup(t *testing.T) {
	f := newFixture(t, nil)
	cache := filepath.Join(f.src, "plugins", "cache", "mkt", "p", "1.0")
	reg := func(path string) string {
		return fmt.Sprintf(`{"version":2,"plugins":{"p@mkt":[{"installPath":%q}]}}`, path)
	}
	write(t, f.src, map[string]string{
		"plugins/installed_plugins.json": reg(cache),
		"plugins/cache/mkt/p/1.0/x.md":   "plugin",
	})
	before := snapshot(t, f.home)
	f.migrate(t, "extensions")
	got := read(t, filepath.Join(f.group, "plugins", "installed_plugins.json"))
	want := filepath.Join(f.group, "plugins", "cache", "mkt", "p", "1.0")
	if !strings.Contains(got, want) || strings.Contains(got, f.src) {
		t.Fatalf("registry: %s", got)
	}
	if !reflect.DeepEqual(before, snapshot(t, f.home)) {
		t.Fatal("the source registry changed")
	}

	// The group already has a registry: the lists are merged, the old one
	// backed up.
	g := newFixture(t, nil)
	other := filepath.Join(g.group, "plugins", "cache", "mkt", "q", "1.0")
	write(t, g.src, map[string]string{
		"plugins/installed_plugins.json": reg(filepath.Join(g.src, "plugins", "cache", "mkt", "p", "1.0")),
		"plugins/cache/mkt/p/1.0/x.md":   "plugin",
	})
	write(t, g.group, map[string]string{"plugins/installed_plugins.json": `{"version":2,"plugins":{"q@mkt":[{"installPath":"` + other + `"}]}}`})
	pl := g.migrate(t, "extensions")
	merged := read(t, filepath.Join(g.group, "plugins", "installed_plugins.json"))
	if !strings.Contains(merged, "q@mkt") || !strings.Contains(merged, filepath.Join(g.group, "plugins", "cache", "mkt", "p", "1.0")) || strings.Contains(merged, g.src) {
		t.Fatalf("merged registry: %s", merged)
	}
	if !strings.Contains(read(t, filepath.Join(pl.BackupDir, "plugins", "installed_plugins.json")), "q@mkt") {
		t.Fatal("old group registry not backed up")
	}
}

func TestMCPFromClaudeJSON(t *testing.T) {
	// ~/.claude keeps its servers in ~/.claude.json.
	f := newFixture(t, map[string]string{
		".claude/CLAUDE.md": "x",
		".claude.json":      `{"projects":{"/a":{"mcpServers":{"proj":{"command":"p"}}}},"mcpServers":{"gh":{"type":"http","url":"https://gh.example/mcp"}}}`,
	})
	if got := ClaudeJSON(f.src, f.home); got != filepath.Join(f.home, ".claude.json") {
		t.Fatalf("ClaudeJSON = %s", got)
	}
	pl := f.build(t, "mcp")
	if !reflect.DeepEqual(pl.MCP.New, []string{"gh"}) {
		t.Fatalf("new servers: %v (project-scope servers must not be read)", pl.MCP.New)
	}

	// Any other folder (or an imported ~/.claude with the link) has it inside.
	other := filepath.Join(f.home, ".claude-work")
	write(t, other, map[string]string{
		"settings.json": "{}",
		".claude.json":  `{"mcpServers":{"fs":{"command":"npx"}}}`,
	})
	if got := ClaudeJSON(other, f.home); got != filepath.Join(other, ".claude.json") {
		t.Fatalf("ClaudeJSON = %s", got)
	}
	pl, err := Build(Options{Source: other, Home: f.home, GroupsRoot: f.groups, Group: "work", Items: []string{"mcp"}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(pl.MCP.New, []string{"fs"}) || pl.GroupDir != filepath.Join(f.groups, "work") {
		t.Fatalf("plan: %+v %+v", pl, pl.MCP)
	}

	// Same server with other key order and spacing: not a conflict.
	write(t, f.group, map[string]string{"mcp.json": "{\"mcpServers\": {\"gh\": {\"url\": \"https://gh.example/mcp\",\n \"type\": \"http\"}}}"})
	if pl := f.build(t, "mcp"); !reflect.DeepEqual(pl.MCP.Same, []string{"gh"}) || !pl.Empty() {
		t.Fatalf("same server: %+v", pl.MCP)
	}
	write(t, f.group, map[string]string{"mcp.json": `{"mcpServers":{"gh":{"type":"http","url":"https://other.example/mcp"}}}`})
	if pl := f.build(t, "mcp"); len(pl.MCP.Conflicts) != 1 || pl.MCP.Conflicts[0].Choice != share.KeepGroup {
		t.Fatalf("different server: %+v", pl.MCP)
	}
}

func TestNotAClaudeFolder(t *testing.T) {
	root := t.TempDir()
	for name, files := range map[string]map[string]string{
		"random":     {"notes.txt": "x"},
		"login-only": {".credentials.json": "x", ".claude.json": "{}"},
	} {
		dir := filepath.Join(root, name)
		write(t, dir, files)
		_, err := Build(Options{Source: dir, Home: root, GroupsRoot: filepath.Join(root, "groups"), Group: "default", Items: DefaultItems})
		if err == nil || !strings.Contains(err.Error(), "does not look like a Claude Code folder") {
			t.Errorf("%s: %v", name, err)
		}
	}
	_, err := Build(Options{Source: filepath.Join(root, "nope"), Home: root, GroupsRoot: root, Group: "default", Items: DefaultItems})
	if err == nil || !strings.Contains(err.Error(), "no Claude Code folder") {
		t.Errorf("missing folder: %v", err)
	}
}

func TestHasContent(t *testing.T) {
	f := newFixture(t, map[string]string{".claude/.credentials.json": "x", ".claude/projects/p/s.jsonl": "{}"})
	if HasContent(f.src) {
		t.Fatal("logins and history alone are not worth a migrate hint")
	}
	write(t, f.src, map[string]string{"skills/x.md": "x"})
	if !HasContent(f.src) {
		t.Fatal("skills/ is")
	}
	if HasContent(filepath.Join(f.home, "missing")) {
		t.Fatal("missing folder")
	}
}

func TestDescribe(t *testing.T) {
	f := conflictFixture(t)
	pl := f.build(t)
	var b bytes.Buffer
	pl.Describe(&b, true)
	out := b.String()
	for _, want := range []string{
		"Copy from " + f.src + " into group default (" + f.group + "):\n",
		"  mcp         ~/.claude.json  new: ctx; differ: fs\n",
		"              fs              differs from the group's copy; you choose: keep group / use mine\n",
		"  settings    settings.json   differs from the group's copy; you choose: keep group / use mine / merge JSON\n",
		"  extensions  commands/       new: 1 file, 6 B\n",
		"              skills/         new: 1 file, 5 B\n",
		"              not found, skipped: agents/, plugins/, output-styles/, rules/\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	for _, c := range pl.Conflicts {
		c.Choice = share.KeepGroup
	}
	b.Reset()
	pl.Describe(&b, false)
	if !strings.Contains(b.String(), "settings.json   differs from the group's copy: keep group\n") {
		t.Errorf("decided plan:\n%s", b.String())
	}
}

func TestSize(t *testing.T) {
	for n, want := range map[int64]string{0: "0 B", 1023: "1023 B", 1536: "1.5 KB", 5 << 20: "5.0 MB", 3 << 30: "3.0 GB"} {
		if got := Size(n); got != want {
			t.Errorf("Size(%d) = %q, want %q", n, got, want)
		}
	}
}

// A stow-style relative link out of ~/.claude still resolves in the group,
// and folders migrate creates are private even when the source's are not.
func TestCopiedLinksResolveAndGroupIsPrivate(t *testing.T) {
	f := newFixture(t, map[string]string{"dotfiles/skill/SKILL.md": "stowed"})
	os.MkdirAll(filepath.Join(f.src, "skills"), 0o755)
	os.Symlink("../../dotfiles/skill", filepath.Join(f.src, "skills", "stow"))
	os.MkdirAll(filepath.Join(f.src, "agents"), 0o755)
	os.WriteFile(filepath.Join(f.src, "agents", "a.md"), []byte("a"), 0o644)
	f.migrate(t, "extensions")
	if got := read(t, filepath.Join(f.group, "skills", "stow", "SKILL.md")); got != "stowed" {
		t.Fatalf("link: %q", got)
	}
	for _, d := range []string{f.groups, f.group} {
		if info, _ := os.Stat(d); info.Mode().Perm() != 0o700 {
			t.Errorf("%s is %v", d, info.Mode().Perm())
		}
	}
}

// The group already lists the plugin the source has: one record, not two.
func TestMergedRegistryHasNoDuplicates(t *testing.T) {
	f := newFixture(t, nil)
	rec := func(dir string) string {
		return fmt.Sprintf(`{"plugins":{"p@mkt":[{"installPath":%q,"version":"1.0"}]}}`, filepath.Join(dir, "plugins", "cache", "mkt", "p"))
	}
	write(t, f.src, map[string]string{"plugins/installed_plugins.json": rec(f.src), "plugins/cache/mkt/p/x.md": "x"})
	write(t, f.group, map[string]string{"plugins/installed_plugins.json": rec(f.group), "plugins/cache/mkt/p/x.md": "x"})
	f.migrate(t, "extensions")
	var reg struct {
		Plugins map[string][]any `json:"plugins"`
	}
	json.Unmarshal([]byte(read(t, filepath.Join(f.group, "plugins", "installed_plugins.json"))), &reg)
	if n := len(reg.Plugins["p@mkt"]); n != 1 {
		t.Fatalf("%d records: %v", n, reg.Plugins)
	}
}

func TestFoundReportsItemsTheSourceHas(t *testing.T) {
	f := newFixture(t, map[string]string{
		".claude/settings.json":        "{}",
		".claude/skills/demo/SKILL.md": "x",
		".claude.json":                 `{"mcpServers":{"fs":{"command":"npx"}}}`,
	})
	got := Found(f.src, f.home)
	want := map[string]bool{"settings": true, "extensions": true, "mcp": true}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
	f = newFixture(t, map[string]string{".claude.json": `{"mcpServers":{}}`})
	if got := Found(f.src, f.home); got["mcp"] {
		t.Fatalf("empty mcpServers counted: %v", got)
	}
}
