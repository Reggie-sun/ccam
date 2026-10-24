package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"sync"
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

// servers is the account's user-scope servers: name → command.
func servers(t *testing.T, dir string) map[string]string {
	t.Helper()
	out := map[string]string{}
	s, _ := readDoc(t, dir)["mcpServers"].(map[string]any)
	for n, v := range s {
		def, _ := v.(map[string]any)
		cmd, _ := def["command"].(string)
		out[n] = cmd
	}
	return out
}

// saveGroup writes the group's mcp.json with a new modification time, as
// another ccam command would.
func saveGroup(t *testing.T, group string, defs map[string]string) {
	t.Helper()
	s := &Store{Servers: map[string]json.RawMessage{}}
	for n, cmd := range defs {
		s.Servers[n] = json.RawMessage(`{"command":"` + cmd + `"}`)
	}
	if err := s.Save(Path(group)); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(Path(group)); err == nil {
		later := info.ModTime().Add(time.Second)
		os.Chtimes(Path(group), later, later)
	}
}

func backups(state string) int {
	b, _ := filepath.Glob(filepath.Join(state, "backups", ".claude.json.*"))
	return len(b)
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
	if info, _ := os.Stat(path); runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatal("mcp.json must be 0600")
	}
	got, _ := Load(path)
	if !reflect.DeepEqual(got.Names(), []string{"ctx"}) {
		t.Fatalf("names %v", got.Names())
	}
}

// Commands that change one group's servers at the same moment all keep
// their change.
func TestUpdateKeepsConcurrentChanges(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g", "mcp.json")
	var wg sync.WaitGroup
	for i := range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := Update(path, func(s *Store) error {
				time.Sleep(time.Millisecond) // a window for the others
				s.Servers[fmt.Sprint("s", i)] = json.RawMessage(`{"command":"x"}`)
				return nil
			}); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if s, _ := Load(path); len(s.Servers) != 16 {
		t.Fatalf("kept %d of 16: %v", len(s.Servers), s.Names())
	}
	// A change that fails saves nothing.
	before, _ := os.ReadFile(path)
	if err := Update(path, func(s *Store) error { s.Servers = nil; return fmt.Errorf("no") }); err == nil {
		t.Fatal("the error is lost")
	}
	if after, _ := os.ReadFile(path); string(after) != string(before) {
		t.Fatal("saved after a failed change")
	}
}

func TestSeedCopiesOnlyIntoAMissingStore(t *testing.T) {
	root := t.TempDir()
	src, dst := filepath.Join(root, "g", "mcp.json"), filepath.Join(root, "h", "mcp.json")
	saveGroup(t, filepath.Dir(src), map[string]string{"fs": "fs"})
	if err := Seed(src, dst); err != nil {
		t.Fatal(err)
	}
	if s, _ := Load(dst); !reflect.DeepEqual(s.Names(), []string{"fs"}) {
		t.Fatalf("seeded %v", s.Names())
	}
	saveGroup(t, filepath.Dir(dst), map[string]string{"own": "own"})
	if err := Seed(src, dst); err != nil {
		t.Fatal(err)
	}
	if s, _ := Load(dst); !reflect.DeepEqual(s.Names(), []string{"own"}) {
		t.Fatalf("an existing mcp.json was replaced: %v", s.Names())
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

func TestSyncMergesKeepsOtherKeysAndBacksUp(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	writeClaudeJSON(t, prof, `{"oauthAccount":{"emailAddress":"a@b"},"projects":{"/x":{"k":1}},"mcpServers":{"mine":{"command":"m"}}}`)
	s := &Store{Servers: map[string]json.RawMessage{"g1": json.RawMessage(`{"command":"one"}`), "g2": json.RawMessage(`{"url":"u"}`)}}
	s.Save(Path(group))

	if _, err := Sync(prof, group, state); err != nil {
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
	if backups(state) != 1 {
		t.Fatalf("backups %d", backups(state))
	}
	if _, err := os.Stat(filepath.Join(prof, ".claude.json.lock")); err == nil {
		t.Fatal("lock dir left behind")
	}

	// Removing g2 from the group removes it from the profile, but never "mine".
	delete(s.Servers, "g2")
	s.Save(Path(group))
	future := time.Now().Add(time.Second)
	os.Chtimes(Path(group), future, future)
	if _, err := Sync(prof, group, state); err != nil {
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
	if _, err := Sync(prof, group, state); err == nil || !strings.Contains(err.Error(), "left untouched") {
		t.Fatalf("bad json: %v", err)
	}
	if b, _ := os.ReadFile(filepath.Join(prof, ".claude.json")); string(b) != "{not json" {
		t.Fatal("bad file was modified")
	}
	writeClaudeJSON(t, prof, `{}`)
	Sync(prof, group, state)
	// The copy is there and the group did not change: nothing is written,
	// so no backup goes round either.
	writeClaudeJSON(t, prof, `{"untouched":true,"mcpServers":{"g":{"command":"x"}}}`)
	n := backups(state)
	if _, err := Sync(prof, group, state); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(prof, ".claude.json")); string(b) != `{"untouched":true,"mcpServers":{"g":{"command":"x"}}}` || backups(state) != n {
		t.Fatalf("written without a change: %s, %d backups", b, backups(state))
	}
}

// A server the account has of its own under a group server's name is never
// replaced or deleted: not by Sync, not when the group drops that name, not
// when the account stops sharing MCP. It is said once.
func TestSyncKeepsTheAccountsOwnServer(t *testing.T) {
	for _, own := range []string{"mine", "fs"} { // another definition, the same one
		t.Run(own, func(t *testing.T) {
			prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
			writeClaudeJSON(t, prof, `{"mcpServers":{"fs":{"command":"`+own+`"},"other":{"command":"o"}}}`)
			saveGroup(t, group, map[string]string{"fs": "fs", "g": "g"})
			res, err := Sync(prof, group, state)
			if err != nil {
				t.Fatal(err)
			}
			if got := servers(t, prof); got["fs"] != own || got["g"] != "g" || got["other"] != "o" {
				t.Fatalf("after Sync: %v", got)
			}
			if want := map[string][]string{"mine": {"fs"}, "fs": nil}[own]; !slices.Equal(res.Own, want) {
				t.Fatalf("said %v, want %v", res.Own, want)
			}
			if drop, _ := Copied(prof, state); !slices.Equal(drop, []string{"g"}) {
				t.Fatalf("copies %v", drop)
			}
			// Said once: the next change of the group does not say it again.
			saveGroup(t, group, map[string]string{"fs": "fs", "g": "g", "h": "h"})
			if res, _ := Sync(prof, group, state); len(res.Own) != 0 {
				t.Fatalf("said again: %v", res.Own)
			}
			// The group drops fs: the account's own stays.
			saveGroup(t, group, map[string]string{"g": "g"})
			Sync(prof, group, state)
			if got := servers(t, prof); got["fs"] != own || got["h"] != "" {
				t.Fatalf("after the group dropped fs: %v", got)
			}
			if err := Unsync(prof, state); err != nil {
				t.Fatal(err)
			}
			if got := servers(t, prof); got["fs"] != own || got["g"] != "" || got["other"] != "o" || HasCopy(state) {
				t.Fatalf("after Unsync: %v", got)
			}
		})
	}
}

// A copy the user changed in a session is the account's own from then on.
func TestSyncLeavesACopyTheUserChanged(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	writeClaudeJSON(t, prof, `{}`)
	saveGroup(t, group, map[string]string{"fs": "fs", "gh": "gh"})
	Sync(prof, group, state)
	writeClaudeJSON(t, prof, `{"mcpServers":{"fs":{"command":"edited"},"gh":{"command":"gh"}}}`)
	saveGroup(t, group, map[string]string{"fs": "fs-2", "gh": "gh-2"})
	res, err := Sync(prof, group, state)
	if err != nil {
		t.Fatal(err)
	}
	if got := servers(t, prof); got["fs"] != "edited" || got["gh"] != "gh-2" || !slices.Equal(res.Own, []string{"fs"}) {
		t.Fatalf("after the group changed: %v, said %v", got, res.Own)
	}
	if drop, keep := Copied(prof, state); !slices.Equal(drop, []string{"gh"}) || len(keep) != 0 {
		t.Fatalf("copies %v %v", drop, keep)
	}
	if err := Unsync(prof, state); err != nil {
		t.Fatal(err)
	}
	if got := servers(t, prof); got["fs"] != "edited" || got["gh"] != "" {
		t.Fatalf("after Unsync: %v", got)
	}
}

// Copied tells Sync's untouched copies from those changed since.
func TestCopiedSplitsChangedCopies(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	saveGroup(t, group, map[string]string{"fs": "fs", "gh": "gh"})
	Sync(prof, group, state)
	writeClaudeJSON(t, prof, `{"mcpServers":{"fs":{"command":"edited"},"gh":{"command":"gh"}}}`)
	if drop, keep := Copied(prof, state); !slices.Equal(drop, []string{"gh"}) || !slices.Equal(keep, []string{"fs"}) {
		t.Fatalf("drop %v keep %v", drop, keep)
	}
}

// A stamp of an older ccam lists names without hashes: such a name counts
// as ccam's copy only while it is the group's definition.
func TestOldStampNamesCountWhileTheyAreTheGroups(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	saveGroup(t, group, map[string]string{"same": "same", "differs": "group"})
	writeClaudeJSON(t, prof, `{"mcpServers":{"same":{"command":"same"},"differs":{"command":"mine"},"dropped":{"command":"d"}}}`)
	os.WriteFile(filepath.Join(state, "sync-stamp"), []byte("1 "+group+"\nsame\ndiffers\ndropped"), 0o600)
	res, err := Sync(prof, group, state)
	if err != nil {
		t.Fatal(err)
	}
	if got := servers(t, prof); got["differs"] != "mine" || got["dropped"] != "d" || got["same"] != "same" || !slices.Equal(res.Own, []string{"differs"}) {
		t.Fatalf("%v, said %v", got, res.Own)
	}
	if drop, _ := Copied(prof, state); !slices.Equal(drop, []string{"same"}) {
		t.Fatalf("copies %v", drop)
	}
}

// Under sync, a copy that is gone from .claude.json comes back at the next
// launch, also while the group's mcp.json stays as it is.
func TestSyncRestoresServersClaudeJSONLost(t *testing.T) {
	for name, lose := range map[string]func(prof string){
		"claude mcp remove": func(prof string) {
			writeClaudeJSON(t, prof, `{"oauthAccount":{"x":1},"mcpServers":{"mine":{"command":"m"}}}`)
		},
		"file made anew": func(prof string) { writeClaudeJSON(t, prof, `{"oauthAccount":{"x":1}}`) },
		"file deleted":   func(prof string) { os.Remove(filepath.Join(prof, ".claude.json")) },
	} {
		t.Run(name, func(t *testing.T) {
			prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
			writeClaudeJSON(t, prof, `{"oauthAccount":{"x":1},"mcpServers":{"mine":{"command":"m"}}}`)
			saveGroup(t, group, map[string]string{"g": "g"})
			Sync(prof, group, state)
			if !Synced(prof, group, state) {
				t.Fatal("not synced after Sync")
			}
			lose(prof)
			if Synced(prof, group, state) {
				t.Fatal("synced although the copy is gone")
			}
			res, err := Sync(prof, group, state)
			if err != nil || !slices.Equal(res.Restored, []string{"g"}) {
				t.Fatalf("restored %v, %v", res.Restored, err)
			}
			if got := servers(t, prof); got["g"] != "g" || !Synced(prof, group, state) {
				t.Fatalf("after Sync: %v", got)
			}
		})
	}
}

// An account that stops sharing MCP keeps the group's servers as its own;
// its own servers of the same names stay, Sync's copies take the group's
// definition, and copies of servers the group no longer has go.
func TestKeepGivesTheAccountTheGroupsServers(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	saveGroup(t, group, map[string]string{"fs": "fs", "gone": "gone", "copy": "copy"})
	writeClaudeJSON(t, prof, `{"oauthAccount":{"x":1},"projects":{"/p":{"mcpServers":{"loc":{"command":"l"}}}}}`)
	Sync(prof, group, state)
	writeClaudeJSON(t, prof, `{"oauthAccount":{"x":1},"projects":{"/p":{"mcpServers":{"loc":{"command":"l"}}}},"mcpServers":{"fs":{"command":"fs"},"gone":{"command":"gone"},"copy":{"command":"copy"},"mine":{"command":"own-mine"}}}`)
	saveGroup(t, group, map[string]string{"fs": "fs-2", "copy": "copy", "mine": "group-mine", "new": "new"})
	copied, own, err := KeepPlan(prof, group, state)
	if err != nil || !slices.Equal(copied, []string{"copy", "fs", "new"}) || !slices.Equal(own, []string{"mine"}) {
		t.Fatalf("plan %v %v %v", copied, own, err)
	}
	if err := Keep(prof, group, state); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"fs": "fs-2", "copy": "copy", "new": "new", "mine": "own-mine"}
	if got := servers(t, prof); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	doc := readDoc(t, prof)
	if doc["oauthAccount"] == nil || doc["projects"] == nil || HasCopy(state) {
		t.Fatalf("other keys or the sync stamp: %v", doc)
	}

	// Joining a group that has fs and copy again takes out the unchanged
	// copies; new was changed meanwhile and stays.
	writeClaudeJSON(t, prof, `{"mcpServers":{"fs":{"command":"fs-2"},"copy":{"command":"copy"},"new":{"command":"edited"},"mine":{"command":"own-mine"}}}`)
	saveGroup(t, group, map[string]string{"fs": "fs-3", "copy": "copy", "new": "new"})
	if got, _ := KeptPlan(prof, Path(group), state); !slices.Equal(got, []string{"copy", "fs"}) {
		t.Fatalf("kept plan %v", got)
	}
	if err := DropKept(prof, group, state); err != nil {
		t.Fatal(err)
	}
	if got := servers(t, prof); !reflect.DeepEqual(got, map[string]string{"new": "edited", "mine": "own-mine"}) {
		t.Fatalf("after DropKept: %v", got)
	}
	if got, _ := KeptPlan(prof, Path(group), state); got != nil {
		t.Fatalf("still noted: %v", got)
	}
}

// Keep makes .claude.json for an account that never started.
func TestKeepWithoutClaudeJSON(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	saveGroup(t, group, map[string]string{"fs": "fs"})
	if err := Keep(prof, group, state); err != nil {
		t.Fatal(err)
	}
	if got := servers(t, prof); got["fs"] != "fs" {
		t.Fatalf("%v", got)
	}
	empty := t.TempDir()
	if err := Keep(t.TempDir(), empty, t.TempDir()); err != nil {
		t.Fatalf("a group without servers: %v", err)
	}
}

// Remove takes servers out of the user scope and out of folders' local
// scope in one write, leaving everything else of those entries.
func TestRemoveUserAndFolderServers(t *testing.T) {
	prof, state := t.TempDir(), t.TempDir()
	writeClaudeJSON(t, prof, `{"mcpServers":{"u":{},"keep":{}},"projects":{"/a":{"hasTrustDialogAccepted":true,"mcpServers":{"l":{},"m":{}}},"/b":{"allowedTools":["x"]}}}`)
	if err := Remove(prof, state, []string{"u"}, map[string][]string{"/a": {"l"}, "/b": {"l"}}); err != nil {
		t.Fatal(err)
	}
	doc := readDoc(t, prof)
	if s := doc["mcpServers"].(map[string]any); len(s) != 1 || s["keep"] == nil {
		t.Fatalf("user scope %v", s)
	}
	projects := doc["projects"].(map[string]any)
	a := projects["/a"].(map[string]any)
	if s := a["mcpServers"].(map[string]any); len(s) != 1 || s["m"] == nil || a["hasTrustDialogAccepted"] != true {
		t.Fatalf("/a %v", a)
	}
	if b := projects["/b"].(map[string]any); b["allowedTools"] == nil || b["mcpServers"] != nil {
		t.Fatalf("/b %v", b)
	}
	if backups(state) != 1 {
		t.Fatalf("%d backups for one change", backups(state))
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
func TestRemoveAndSyncKeepSymlinkedClaudeJSON(t *testing.T) {
	home, prof, state, group := t.TempDir(), t.TempDir(), t.TempDir(), t.TempDir()
	real := filepath.Join(home, ".claude.json")
	os.WriteFile(real, []byte(`{"mcpServers":{"x":{"command":"x"}}}`), 0o600)
	os.Symlink(real, filepath.Join(prof, ".claude.json"))
	if err := Remove(prof, state, []string{"x"}, nil); err != nil {
		t.Fatal(err)
	}
	(&Store{Servers: map[string]json.RawMessage{"g": json.RawMessage(`{"command":"g"}`)}}).Save(Path(group))
	if _, err := Sync(prof, group, state); err != nil {
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

// Shell commands in server definitions stay readable in .claude.json.
func TestSyncAndRemoveDoNotEscapeShellCharacters(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	writeClaudeJSON(t, prof, `{"mcpServers":{"mine":{"command":"sh","args":["-c","a && b > c"]},"x":{}}}`)
	s := &Store{Servers: map[string]json.RawMessage{"g": json.RawMessage(`{"command":"sh","args":["-c","d && e > f"]}`)}}
	s.Save(Path(group))
	if _, err := Sync(prof, group, state); err != nil {
		t.Fatal(err)
	}
	if err := Remove(prof, state, []string{"x"}, nil); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(prof, ".claude.json"))
	if !strings.Contains(string(data), "a && b > c") || !strings.Contains(string(data), "d && e > f") {
		t.Fatalf("escaped:\n%s", data)
	}
}

// An older ccam's stamp is noted with hashes at the next Sync while the
// group's mcp.json is the one it synced: its copies then follow the
// group's later changes. Nothing is written to .claude.json for it.
func TestOldStampIsNotedWhileTheGroupIsUnchanged(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	saveGroup(t, group, map[string]string{"fs": "fs-1", "gh": "gh"})
	writeClaudeJSON(t, prof, `{"mcpServers":{"fs":{"command":"fs-1"},"gh":{"command":"gh"}}}`)
	info, _ := os.Stat(Path(group))
	os.WriteFile(filepath.Join(state, "sync-stamp"), []byte(syncStamp(group, info)+"\nfs\ngh"), 0o600)
	before, _ := os.ReadFile(filepath.Join(prof, ".claude.json"))
	if res, err := Sync(prof, group, state); err != nil || len(res.Own)+len(res.Restored) != 0 {
		t.Fatalf("%+v %v", res, err)
	}
	if after, _ := os.ReadFile(filepath.Join(prof, ".claude.json")); string(after) != string(before) || backups(state) != 0 {
		t.Fatal(".claude.json written for the stamp alone")
	}
	saveGroup(t, group, map[string]string{"fs": "fs-2", "gh": "gh"})
	if res, _ := Sync(prof, group, state); len(res.Own) != 0 || servers(t, prof)["fs"] != "fs-2" {
		t.Fatalf("the copy did not follow the group: %v, said %v", servers(t, prof), res.Own)
	}
}

// Keep run again (a change finished by doctor --fix after it stopped)
// keeps what the first run noted.
func TestKeepAgainKeepsTheFirstNote(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	saveGroup(t, group, map[string]string{"fs": "fs"})
	writeClaudeJSON(t, prof, `{}`)
	if err := Keep(prof, group, state); err != nil {
		t.Fatal(err)
	}
	saveGroup(t, group, map[string]string{"fs": "fs-2", "new": "new"})
	if err := Keep(prof, group, state); err != nil {
		t.Fatal(err)
	}
	if got := servers(t, prof); got["fs"] != "fs-2" || got["new"] != "new" {
		t.Fatalf("after the second Keep: %v", got)
	}
	if drop, _ := KeptPlan(prof, Path(group), state); !slices.Equal(drop, []string{"fs", "new"}) {
		t.Fatalf("noted %v", drop)
	}
}

// An older ccam's stamp is noted with hashes right before the group's
// mcp.json changes (NoteCopies, as `ccam mcp add` does): the copy then
// follows the change, and a copy the user changed stays the account's own.
func TestOldStampIsNotedBeforeTheGroupChanges(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	saveGroup(t, group, map[string]string{"fs": "fs-1", "gh": "gh"})
	writeClaudeJSON(t, prof, `{"mcpServers":{"fs":{"command":"fs-1"},"gh":{"command":"mine"}}}`)
	info, _ := os.Stat(Path(group))
	os.WriteFile(filepath.Join(state, "sync-stamp"), []byte(syncStamp(group, info)+"\nfs\ngh"), 0o600)
	before, _ := os.ReadFile(filepath.Join(prof, ".claude.json"))
	if err := NoteCopies(prof, group, state); err != nil {
		t.Fatal(err)
	}
	if after, _ := os.ReadFile(filepath.Join(prof, ".claude.json")); string(after) != string(before) {
		t.Fatal(".claude.json written")
	}
	saveGroup(t, group, map[string]string{"fs": "fs-2", "gh": "gh-2"})
	res, err := Sync(prof, group, state)
	if err != nil {
		t.Fatal(err)
	}
	if got := servers(t, prof); got["fs"] != "fs-2" || got["gh"] != "mine" || !slices.Equal(res.Own, []string{"gh"}) {
		t.Fatalf("after the group changed: %v, own %v", got, res.Own)
	}
	// Once the group's file changed, nothing can be noted any more.
	os.WriteFile(filepath.Join(state, "sync-stamp"), []byte("1 "+group+"\nfs"), 0o600)
	if err := NoteCopies(prof, group, state); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(filepath.Join(state, "sync-stamp")); string(data) != "1 "+group+"\nfs" {
		t.Fatalf("stamp rewritten: %q", data)
	}
}

// Claude Code's .claude.json.lock is honoured: while a live one is held,
// Sync and Remove give up and leave the file as it is (a write then
// could undo Claude Code's own); a lock left behind long ago is taken over.
func TestClaudeJSONLockIsHonoured(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	const before = `{"mcpServers":{"mine":{"command":"m"}}}`
	writeClaudeJSON(t, prof, before)
	s := &Store{Servers: map[string]json.RawMessage{"g1": json.RawMessage(`{"command":"one"}`)}}
	s.Save(Path(group))
	held := filepath.Join(prof, ".claude.json.lock")
	if err := os.Mkdir(held, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := Sync(prof, group, state); err == nil {
		t.Fatal("Sync wrote while Claude Code held the lock")
	}
	if err := Remove(prof, state, []string{"mine"}, nil); err == nil {
		t.Fatal("Remove wrote while Claude Code held the lock")
	}
	if data, _ := os.ReadFile(filepath.Join(prof, ".claude.json")); string(data) != before {
		t.Fatalf("changed under the lock: %s", data)
	}
	old := time.Now().Add(-time.Minute)
	os.Chtimes(held, old, old)
	if _, err := Sync(prof, group, state); err != nil {
		t.Fatalf("a stale lock: %v", err)
	}
	if readDoc(t, prof)["mcpServers"].(map[string]any)["g1"] == nil {
		t.Fatal("not synced after taking over the stale lock")
	}
}

// A launch under the sync strategy skips Sync, and the running check
// before it, while neither the group's mcp.json nor .claude.json changed
// since Sync last left them complete.
func TestSyncFreshWhileNothingChanged(t *testing.T) {
	prof, group, state := t.TempDir(), t.TempDir(), t.TempDir()
	writeClaudeJSON(t, prof, `{"mcpServers":{"mine":{"command":"m"}}}`)
	(&Store{Servers: map[string]json.RawMessage{"g": json.RawMessage(`{"command":"g"}`)}}).Save(Path(group))
	if SyncFresh(prof, group, state) {
		t.Fatal("fresh before any sync")
	}
	if _, err := Sync(prof, group, state); err != nil {
		t.Fatal(err)
	}
	if !SyncFresh(prof, group, state) {
		t.Fatal("not fresh right after a sync")
	}
	// A session removes the copy: .claude.json changed.
	writeClaudeJSON(t, prof, `{"mcpServers":{"mine":{"command":"m"}}}`)
	future := time.Now().Add(2 * time.Second)
	os.Chtimes(filepath.Join(prof, ".claude.json"), future, future)
	if SyncFresh(prof, group, state) {
		t.Fatal("fresh after .claude.json changed")
	}
	if res, err := Sync(prof, group, state); err != nil || !slices.Equal(res.Restored, []string{"g"}) || !SyncFresh(prof, group, state) {
		t.Fatalf("%+v %v", res, err)
	}
	(&Store{Servers: map[string]json.RawMessage{"g": json.RawMessage(`{"command":"g2"}`)}}).Save(Path(group))
	os.Chtimes(Path(group), future.Add(time.Second), future.Add(time.Second))
	if SyncFresh(prof, group, state) {
		t.Fatal("fresh after the group's mcp.json changed")
	}
	Sync(prof, group, state)
	if err := Unsync(prof, state); err != nil || SyncFresh(prof, group, state) {
		t.Fatalf("fresh after Unsync: %v", err)
	}
}

// A copy an earlier Keep made, unchanged since, is replaced by the next
// Keep: its plan lists it as copied, not as the account's own.
func TestKeepPlanCountsAnEarlierKeepsCopy(t *testing.T) {
	root := t.TempDir()
	acct, state := filepath.Join(root, "a"), filepath.Join(root, "state")
	g1, g2 := filepath.Join(root, "g1"), filepath.Join(root, "g2")
	writeClaudeJSON(t, acct, `{"mcpServers":{}}`)
	saveGroup(t, g1, map[string]string{"x": "def1"})
	if err := Keep(acct, g1, state); err != nil { // a leaves g1
		t.Fatal(err)
	}
	saveGroup(t, g2, map[string]string{"y": "y", "x": "def2"}) // joins g2, which gets an x later
	copied, own, err := KeepPlan(acct, g2, state)
	if err != nil || !slices.Equal(copied, []string{"x", "y"}) || len(own) != 0 {
		t.Fatalf("plan: copied %v own %v %v", copied, own, err)
	}
	if err := Keep(acct, g2, state); err != nil {
		t.Fatal(err)
	}
	if got := servers(t, acct); got["x"] != "def2" {
		t.Fatalf("Keep did not do what its plan said: %v", got)
	}
}
