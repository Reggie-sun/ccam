package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func writeClaudeJSON(t *testing.T, dir, content string) {
	t.Helper()
	os.MkdirAll(dir, 0o700)
	if err := os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readDoc(t *testing.T, dir string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, ".claude.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestStoreRoundTripAndPerm(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g", "mcp.json")
	s, err := Load(path)
	if err != nil || len(s.Servers) != 0 {
		t.Fatalf("empty load: %v %v", s, err)
	}
	s.Servers["ctx"] = json.RawMessage(`{"type":"http","url":"https://x"}`)
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o600 {
		t.Fatal("mcp.json must be 0600")
	}
	got, _ := Load(path)
	if !reflect.DeepEqual(got.Names(), []string{"ctx"}) {
		t.Fatalf("names %v", got.Names())
	}
}

func TestValidate(t *testing.T) {
	if Validate(json.RawMessage(`{"command":"npx"}`)) != nil || Validate(json.RawMessage(`{"url":"u"}`)) != nil {
		t.Fatal("valid rejected")
	}
	if Validate(json.RawMessage(`{"type":"stdio"}`)) == nil || Validate(json.RawMessage(`[1]`)) == nil {
		t.Fatal("invalid accepted")
	}
}

func TestFlagArgs(t *testing.T) {
	g := t.TempDir()
	if FlagArgs(g) != nil {
		t.Fatal("no mcp.json → no args")
	}
	os.WriteFile(Path(g), []byte(`{"mcpServers":{}}`), 0o600)
	if got := FlagArgs(g); !reflect.DeepEqual(got, []string{"--mcp-config=" + Path(g)}) {
		t.Fatalf("%v", got)
	}
}

func TestNewLocalReportsOnceAndSkipsShared(t *testing.T) {
	dir := t.TempDir()
	seen := filepath.Join(t.TempDir(), "seen")
	group := filepath.Join(t.TempDir(), "mcp.json")
	(&Store{Servers: map[string]json.RawMessage{"shared": json.RawMessage(`{}`)}}).Save(group)
	writeClaudeJSON(t, dir, `{"mcpServers":{"shared":{},"mine":{"command":"x"}}}`)
	got, err := NewLocal(dir, group, seen)
	if err != nil || !reflect.DeepEqual(got, []string{"mine"}) {
		t.Fatalf("first: %v %v", got, err)
	}
	if got, _ := NewLocal(dir, group, seen); got != nil {
		t.Fatalf("unchanged file reported again: %v", got)
	}
	// Claude rewrites the file with one more server.
	writeClaudeJSON(t, dir, `{"mcpServers":{"shared":{},"mine":{"command":"x"},"new":{"command":"y"}}}`)
	future := time.Now().Add(time.Second)
	os.Chtimes(filepath.Join(dir, ".claude.json"), future, future)
	if got, _ := NewLocal(dir, group, seen); !reflect.DeepEqual(got, []string{"new"}) {
		t.Fatalf("after change: %v", got)
	}
}

func TestSyncMergesKeepsOtherKeysAndBacksUp(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	writeClaudeJSON(t, prof, `{"oauthAccount":{"emailAddress":"a@b"},"projects":{"/x":{"k":1}},"mcpServers":{"mine":{"command":"m"}}}`)
	s := &Store{Servers: map[string]json.RawMessage{"g1": json.RawMessage(`{"command":"one"}`), "g2": json.RawMessage(`{"url":"u"}`)}}
	s.Save(Path(group))

	if err := Sync(prof, group, state); err != nil {
		t.Fatal(err)
	}
	doc := readDoc(t, prof)
	servers := doc["mcpServers"].(map[string]any)
	if len(servers) != 3 || servers["mine"] == nil || servers["g1"] == nil {
		t.Fatalf("servers %v", servers)
	}
	if doc["oauthAccount"].(map[string]any)["emailAddress"] != "a@b" || doc["projects"] == nil {
		t.Fatal("other keys lost")
	}
	if b, _ := filepath.Glob(filepath.Join(state, "backups", ".claude.json.*")); len(b) != 1 {
		t.Fatalf("backups %v", b)
	}
	if _, err := os.Stat(filepath.Join(prof, ".claude.json.lock")); err == nil {
		t.Fatal("lock dir left behind")
	}

	// Removing g2 from the group removes it from the profile, but never "mine".
	delete(s.Servers, "g2")
	s.Save(Path(group))
	future := time.Now().Add(time.Second)
	os.Chtimes(Path(group), future, future)
	if err := Sync(prof, group, state); err != nil {
		t.Fatal(err)
	}
	servers = readDoc(t, prof)["mcpServers"].(map[string]any)
	if servers["g2"] != nil || servers["mine"] == nil || servers["g1"] == nil {
		t.Fatalf("after removal %v", servers)
	}
}

func TestSyncSkipsWhenUnchangedAndRefusesBadJSON(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	(&Store{Servers: map[string]json.RawMessage{"g": json.RawMessage(`{"command":"x"}`)}}).Save(Path(group))
	writeClaudeJSON(t, prof, `{not json`)
	if err := Sync(prof, group, state); err == nil || !strings.Contains(err.Error(), "손대지 않음") {
		t.Fatalf("bad json: %v", err)
	}
	if b, _ := os.ReadFile(filepath.Join(prof, ".claude.json")); string(b) != "{not json" {
		t.Fatal("bad file was modified")
	}
	writeClaudeJSON(t, prof, `{}`)
	Sync(prof, group, state)
	writeClaudeJSON(t, prof, `{"untouched":true}`)
	Sync(prof, group, state) // same mcp.json mtime → no write
	if _, ok := readDoc(t, prof)["mcpServers"]; ok {
		t.Fatal("sync should have been skipped")
	}
}

func TestRemoveLocal(t *testing.T) {
	prof, state := t.TempDir(), t.TempDir()
	writeClaudeJSON(t, prof, `{"a":1,"mcpServers":{"x":{},"y":{}}}`)
	if err := RemoveLocal(prof, state, []string{"x"}); err != nil {
		t.Fatal(err)
	}
	doc := readDoc(t, prof)
	if s := doc["mcpServers"].(map[string]any); len(s) != 1 || s["y"] == nil || doc["a"] != 1.0 {
		t.Fatalf("%v", doc)
	}
}

func TestBackupRotation(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < keepBackups+3; i++ {
		if err := backup("/x/.claude.json", []byte("x"), dir); err != nil {
			t.Fatal(err)
		}
	}
	if b, _ := filepath.Glob(filepath.Join(dir, ".claude.json.*")); len(b) != keepBackups {
		t.Fatalf("kept %d", len(b))
	}
}

// An imported ~/.claude links .claude.json to ~/.claude.json; editing it
// must keep the link so plain `claude` and ccam stay one account.
func TestRemoveLocalAndSyncKeepSymlinkedClaudeJSON(t *testing.T) {
	home, prof, state, group := t.TempDir(), t.TempDir(), t.TempDir(), t.TempDir()
	real := filepath.Join(home, ".claude.json")
	os.WriteFile(real, []byte(`{"mcpServers":{"x":{"command":"x"}}}`), 0o600)
	os.Symlink(real, filepath.Join(prof, ".claude.json"))
	if err := RemoveLocal(prof, state, []string{"x"}); err != nil {
		t.Fatal(err)
	}
	(&Store{Servers: map[string]json.RawMessage{"g": json.RawMessage(`{"command":"g"}`)}}).Save(Path(group))
	if err := Sync(prof, group, state); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Lstat(filepath.Join(prof, ".claude.json")); info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("link replaced")
	}
	data, _ := os.ReadFile(real)
	if strings.Contains(string(data), `"x"`) || !strings.Contains(string(data), `"g"`) {
		t.Fatalf("real file: %s", data)
	}
}
