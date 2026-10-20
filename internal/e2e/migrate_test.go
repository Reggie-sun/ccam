package e2e

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// seedClaude writes an existing Claude Code setup into the harness's
// ~/.claude and ~/.claude.json and returns the user home.
func (h *harness) seedClaude() string {
	h.t.Helper()
	home := filepath.Join(h.home, "home")
	for rel, content := range map[string]string{
		".claude/settings.json":          `{"model":"opus"}`,
		".claude/CLAUDE.md":              "# my rules",
		".claude/skills/review/SKILL.md": "skill",
		".claude/projects/p/s1.jsonl":    "{}\n",
		".claude/.credentials.json":      `{"claudeAiOauth":{"accessToken":"secret"}}`,
		".claude.json":                   `{"oauthAccount":{"emailAddress":"old@example.com"},"mcpServers":{"fs":{"type":"stdio","command":"npx","args":["fs"]}}}`,
	} {
		path := filepath.Join(home, rel)
		os.MkdirAll(filepath.Dir(path), 0o700)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			h.t.Fatal(err)
		}
	}
	return home
}

func treeState(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		info, _ := os.Lstat(path)
		v := fmt.Sprintf("%v %d", info.Mode(), info.ModTime().UnixNano())
		if info.Mode().IsRegular() {
			data, _ := os.ReadFile(path)
			sum := sha256.Sum256(data)
			v += " " + hex.EncodeToString(sum[:])
		}
		out[path] = v
		return nil
	})
	return out
}

func TestMigrateDryRunChangesNothing(t *testing.T) {
	h := newHarness(t)
	home := h.seedClaude()
	before := treeState(t, home)
	out := h.ok("ccam", "--dry-run", "migrate")
	for _, want := range []string{
		"Copy from " + filepath.Join(home, ".claude") + " into group default",
		"  mcp         ~/.claude.json  new: fs\n",
		"  settings    settings.json   new: 16 B\n",
		"  extensions  skills/         new: 1 file, 5 B\n",
		"add --history to copy it",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "[y/N]") {
		t.Error("a dry run does not ask")
	}
	if _, err := os.Stat(filepath.Join(h.home, "data", "groups")); err == nil {
		t.Error("dry run created the groups folder")
	}
	if !reflect.DeepEqual(before, treeState(t, home)) {
		t.Error("dry run changed ~/.claude")
	}

	// Without --yes and without a "y", nothing is copied either.
	r := h.run("\n", nil, "ccam", "migrate")
	if r.code != 0 || !strings.Contains(r.stdout, "Copy these into group default? [y/N] Cancelled.") {
		t.Fatalf("%+v", r)
	}
	if _, err := os.Stat(filepath.Join(h.home, "data", "groups", "default", "settings.json")); err == nil {
		t.Error("answering no copied anyway")
	}
}

func TestMigrateCopiesIntoTheSharedGroup(t *testing.T) {
	h := newHarness(t)
	home := h.seedClaude()
	// The first account gets a hint, the second does not.
	r := h.run("", []string{"FAKE_EMAIL=me@example.com"}, "ccam", "add", "personal")
	if !strings.Contains(r.stdout, "Tip: run `ccam migrate`") {
		t.Fatalf("first add: %+v", r)
	}
	r = h.run("", []string{"FAKE_EMAIL=dev@company.com"}, "ccam", "add", "work", "--no-group")
	if r.code != 0 || strings.Contains(r.stdout, "ccam migrate") {
		t.Fatalf("second add: %+v", r)
	}
	before := treeState(t, home)

	out := h.ok("ccam", "migrate", "--yes")
	if !strings.Contains(out, "✓ Copied into group default: 1 MCP server, settings.json, CLAUDE.md, skills/\n") {
		t.Fatalf("summary:\n%s", out)
	}
	for _, want := range []string{
		"Group default: personal\n",
		"  ccam group join default work\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	// The group's accounts see them through their links.
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "settings.json")); string(b) != `{"model":"opus"}` {
		t.Errorf("settings.json: %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "skills", "review", "SKILL.md")); string(b) != "skill" {
		t.Errorf("skills: %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "CLAUDE.md")); string(b) != "# my rules" {
		t.Errorf("personal CLAUDE.md: %q", b)
	}
	if _, err := os.Lstat(filepath.Join(h.profileDir("work"), "CLAUDE.md")); err == nil {
		t.Error("work is in no group")
	}
	// The group's MCP server reaches the launch.
	h.ok("ccam", "personal")
	if _, args := h.called(); len(args) != 1 || args[0] != "--mcp-config="+filepath.Join(h.home, "data", "groups", "default", "mcp.json") {
		t.Fatalf("launch args: %q", args)
	}
	group := filepath.Join(h.home, "data", "groups", "default")
	for _, name := range []string{"projects", "history.jsonl", ".credentials.json", ".claude.json"} {
		if _, err := os.Lstat(filepath.Join(group, name)); err == nil {
			t.Errorf("%s must not be copied", name)
		}
	}
	if !reflect.DeepEqual(before, treeState(t, home)) {
		t.Error("migrate changed ~/.claude")
	}

	// A changed source with --yes keeps the group's copy and says how to choose.
	os.WriteFile(filepath.Join(home, ".claude", "settings.json"), []byte(`{"theme":"dark"}`), 0o600)
	out = h.ok("ccam", "migrate", "--yes", "--settings")
	if !strings.Contains(out, "Kept the group's copy of: settings.json\nTo choose instead, run `ccam migrate` without --yes.\n") {
		t.Fatalf("--yes conflict:\n%s", out)
	}
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "settings.json")); string(b) != `{"model":"opus"}` {
		t.Fatalf("--yes overwrote the group: %q", b)
	}
	// Choosing "merge JSON" interactively writes the merge; the old group copy is saved.
	r = h.run("3\ny\n", nil, "ccam", "migrate", "--settings")
	if r.code != 0 || !strings.Contains(r.stdout, "3) merge JSON") || !strings.Contains(r.stdout, "The group's previous copies are saved in "+filepath.Join(group, ".ccam-backup")) {
		t.Fatalf("%+v", r)
	}
	b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "settings.json"))
	if !strings.Contains(string(b), "opus") || !strings.Contains(string(b), "dark") {
		t.Fatalf("merged settings: %s", b)
	}
	if out := h.ok("ccam", "doctor"); strings.Contains(out, "✗") {
		t.Fatalf("doctor after migrate:\n%s", out)
	}
}

func TestMigrateErrors(t *testing.T) {
	h := newHarness(t)
	r := h.run("", nil, "ccam", "migrate")
	if r.code == 0 || !strings.Contains(r.stderr, "no Claude Code folder at "+filepath.Join(h.home, "home", ".claude")) {
		t.Fatalf("missing ~/.claude: %+v", r)
	}
	notClaude := filepath.Join(h.home, "notes")
	os.MkdirAll(notClaude, 0o700)
	os.WriteFile(filepath.Join(notClaude, "todo.txt"), []byte("x"), 0o600)
	r = h.run("", nil, "ccam", "migrate", notClaude)
	if r.code == 0 || !strings.Contains(r.stderr, "does not look like a Claude Code folder") {
		t.Fatalf("not a Claude folder: %+v", r)
	}
	h.seedClaude()
	r = h.run("", nil, "ccam", "migrate", "--group", "bad/group")
	if r.code == 0 || !strings.Contains(r.stderr, "invalid group name") {
		t.Fatalf("bad group: %+v", r)
	}
	r = h.run("", nil, "ccam", "add", "migrate")
	if r.code == 0 {
		t.Fatal("migrate must be a reserved name")
	}
}
