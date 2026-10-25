package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"time"
)

// The shape Claude Code 2.1.308 writes for `claude mcp add` without -s
// user: folders in projects, a git repository under its top folder.
func TestFolderServersReadsWhatClaudeCodeWrites(t *testing.T) {
	data, err := os.ReadFile("testdata/local-scope-2.1.308.json")
	if err != nil {
		t.Fatal(err)
	}
	prof := t.TempDir()
	os.WriteFile(filepath.Join(prof, ".claude.json"), data, 0o600)
	got, err := FolderServers(prof)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string][]string{}
	for folder, s := range got {
		names[folder] = sortedKeys(s)
	}
	want := map[string][]string{"/w/plain": {"linksrv", "plainsrv"}, "/w/plain/sub": {"plainsubsrv"}, "/w/repo": {"reposubsrv", "wtsrv"}}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("got %v", names)
	}
	if user, _ := LocalServers(prof); len(user) != 1 || user["usersrv"] == nil {
		t.Fatalf("user scope %v", user)
	}
	if got, err := FolderServers(t.TempDir()); err != nil || got != nil {
		t.Fatalf("no .claude.json: %v %v", got, err)
	}
}

// Claude Code files a folder's servers under the folder with links
// resolved, and a git repository's (its folders, its worktrees) under the
// top of the main worktree (checked with 2.1.308).
func TestProjectKey(t *testing.T) {
	root, _ := filepath.EvalSymlinks(t.TempDir())
	plain := filepath.Join(root, "plain")
	repo := filepath.Join(root, "repo")
	wt := filepath.Join(root, "wt")
	bare, bwt := filepath.Join(root, "proj.git"), filepath.Join(root, "bwt")
	for _, d := range []string{filepath.Join(plain, "sub"), filepath.Join(repo, ".git", "worktrees", "wt"), filepath.Join(repo, "sub", "deeper"), wt, filepath.Join(bare, "worktrees", "bwt"), bwt} {
		os.MkdirAll(d, 0o700)
	}
	os.WriteFile(filepath.Join(wt, ".git"), []byte("gitdir: "+filepath.Join(repo, ".git", "worktrees", "wt")+"\n"), 0o600)
	os.WriteFile(filepath.Join(repo, ".git", "worktrees", "wt", "commondir"), []byte("../..\n"), 0o600)
	// A worktree of a bare repository, as `git worktree add` makes it.
	os.WriteFile(filepath.Join(bwt, ".git"), []byte("gitdir: "+filepath.Join(bare, "worktrees", "bwt")+"\n"), 0o600)
	os.WriteFile(filepath.Join(bare, "worktrees", "bwt", "commondir"), []byte("../..\n"), 0o600)
	link := filepath.Join(root, "link")
	linked := os.Symlink(plain, link) == nil
	for _, c := range []struct{ dir, want string }{
		{plain, plain},
		{filepath.Join(plain, "sub"), filepath.Join(plain, "sub")}, // outside git a folder is its own
		{repo, repo},
		{filepath.Join(repo, "sub", "deeper"), repo},
		{wt, repo},
		{bwt, bare}, // checked with 2.1.308: the bare repository's folder
	} {
		if got := ProjectKey(c.dir); got != c.want {
			t.Errorf("%s: got %s, want %s", c.dir, got, c.want)
		}
	}
	if linked {
		if got := ProjectKey(link); got != plain {
			t.Errorf("through a link: %s", got)
		}
	}
}

// Trust is per account; a folder above dir counts only up to the top of
// its git repository, and the main worktree's for its other worktrees
// (checked with 2.1.308).
func TestTrusts(t *testing.T) {
	root, _ := filepath.EvalSymlinks(t.TempDir())
	repo := filepath.Join(root, "repo")
	os.MkdirAll(filepath.Join(repo, ".git", "worktrees", "wt"), 0o700)
	os.MkdirAll(filepath.Join(repo, "sub"), 0o700)
	wt := filepath.Join(root, "wt")
	os.MkdirAll(wt, 0o700)
	os.WriteFile(filepath.Join(wt, ".git"), []byte("gitdir: "+filepath.Join(repo, ".git", "worktrees", "wt")+"\n"), 0o600)
	os.WriteFile(filepath.Join(repo, ".git", "worktrees", "wt", "commondir"), []byte("../..\n"), 0o600)
	plain := filepath.Join(root, "plain")
	os.MkdirAll(plain, 0o700)
	prof := t.TempDir()
	trust := func(folders ...string) {
		p := map[string]any{}
		for _, f := range folders {
			p[f] = map[string]any{"hasTrustDialogAccepted": true}
		}
		data, _ := json.Marshal(map[string]any{"projects": p})
		os.WriteFile(filepath.Join(prof, ".claude.json"), data, 0o600)
	}
	for _, c := range []struct {
		trusted []string
		dir     string
		want    bool
	}{
		{nil, plain, false},
		{[]string{plain}, plain, true},
		{[]string{root}, plain, false}, // outside git a parent does not count
		{[]string{repo}, filepath.Join(repo, "sub"), true},
		{[]string{filepath.Join(repo, "sub")}, filepath.Join(repo, "sub"), true},
		{[]string{filepath.Join(repo, "sub")}, repo, false},
		{[]string{root}, repo, false}, // above the repository
		{[]string{repo}, wt, true},    // the main worktree's trust
		{[]string{wt}, wt, true},
		{[]string{wt}, repo, false},
	} {
		trust(c.trusted...)
		if got := Trusts(prof, c.dir); got != c.want {
			t.Errorf("trusted %v, in %s: got %v", c.trusted, c.dir, got)
		}
	}
	if Trusts(t.TempDir(), plain) {
		t.Error("an account without .claude.json trusts nothing")
	}
}

// A session starts without the trust dialog in a folder the account
// trusted, or below one: inside git up to the top of the repository,
// outside git any folder above (Claude Code 2.1.310, read from its code).
func TestSessionTrusts(t *testing.T) {
	root, _ := filepath.EvalSymlinks(t.TempDir())
	repo := filepath.Join(root, "repo")
	os.MkdirAll(filepath.Join(repo, ".git", "worktrees", "wt"), 0o700)
	os.MkdirAll(filepath.Join(repo, "sub", "deeper"), 0o700)
	wt := filepath.Join(root, "wt")
	os.MkdirAll(wt, 0o700)
	os.WriteFile(filepath.Join(wt, ".git"), []byte("gitdir: "+filepath.Join(repo, ".git", "worktrees", "wt")+"\n"), 0o600)
	os.WriteFile(filepath.Join(repo, ".git", "worktrees", "wt", "commondir"), []byte("../..\n"), 0o600)
	plain := filepath.Join(root, "plain", "sub")
	os.MkdirAll(plain, 0o700)
	prof := t.TempDir()
	trust := func(folders ...string) {
		p := map[string]any{}
		for _, f := range folders {
			p[f] = map[string]any{"hasTrustDialogAccepted": true}
		}
		data, _ := json.Marshal(map[string]any{"projects": p})
		os.WriteFile(filepath.Join(prof, ".claude.json"), data, 0o600)
	}
	for _, c := range []struct {
		trusted []string
		dir     string
		want    bool
	}{
		{nil, plain, false},
		{[]string{plain}, plain, true},
		{[]string{filepath.Dir(plain)}, plain, true}, // outside git any folder above counts
		{[]string{root}, plain, true},
		{[]string{filepath.Join(plain, "x")}, plain, false},
		{[]string{repo}, filepath.Join(repo, "sub", "deeper"), true},
		{[]string{filepath.Join(repo, "sub")}, filepath.Join(repo, "sub", "deeper"), true},
		{[]string{root}, filepath.Join(repo, "sub"), false}, // above the repository
		{[]string{repo}, wt, true},                          // the main worktree's trust
		{[]string{wt}, repo, false},
	} {
		trust(c.trusted...)
		if got := SessionTrusts(prof, c.dir); got != c.want {
			t.Errorf("trusted %v, in %s: got %v", c.trusted, c.dir, got)
		}
	}
	// Trusts, the rule for a project's .mcp.json, does not count a parent
	// outside git.
	trust(filepath.Dir(plain))
	if Trusts(prof, plain) {
		t.Error("Trusts counts a parent outside git")
	}
	if SessionTrusts(t.TempDir(), plain) {
		t.Error("an account without .claude.json trusts nothing")
	}
}

func TestNewLocalReportsOnceAndSkipsShared(t *testing.T) {
	dir := t.TempDir()
	seen := filepath.Join(t.TempDir(), "seen")
	group := filepath.Join(t.TempDir(), "mcp.json")
	(&Store{Servers: map[string]json.RawMessage{"shared": json.RawMessage(`{}`)}}).Save(group)
	writeClaudeJSON(t, dir, `{"mcpServers":{"shared":{},"mine":{"command":"x"}}}`)
	got, err := NewLocal(dir, group, seen)
	if err != nil || !slices.Equal(got.User, []string{"mine"}) || len(got.Folder)+len(got.Shadow) != 0 {
		t.Fatalf("first: %+v %v", got, err)
	}
	if got, _ := NewLocal(dir, group, seen); len(got.User)+len(got.Folder)+len(got.Shadow) != 0 {
		t.Fatalf("unchanged file reported again: %+v", got)
	}
	// Claude rewrites the file with one more server.
	writeClaudeJSON(t, dir, `{"mcpServers":{"shared":{},"mine":{"command":"x"},"new":{"command":"y"}}}`)
	future := time.Now().Add(time.Second)
	os.Chtimes(filepath.Join(dir, ".claude.json"), future, future)
	if got, _ := NewLocal(dir, group, seen); !slices.Equal(got.User, []string{"new"}) {
		t.Fatalf("after change: %+v", got)
	}
}

// Servers `claude mcp add` puts in one folder's local scope, and the
// account's own servers that hide a group server, are reported once each.
func TestNewLocalReportsFolderServersAndShadows(t *testing.T) {
	dir := t.TempDir()
	seen := filepath.Join(t.TempDir(), "seen")
	group := filepath.Join(t.TempDir(), "mcp.json")
	(&Store{Servers: map[string]json.RawMessage{"fs": json.RawMessage(`{"command":"group"}`), "same": json.RawMessage(`{"command":"s"}`)}}).Save(group)
	writeClaudeJSON(t, dir, `{"mcpServers":{"fs":{"command":"mine"},"same":{"command":"s"}},"projects":{"/p":{"mcpServers":{"loc":{"command":"l"}}},"/q":{"allowedTools":[]}}}`)
	got, err := NewLocal(dir, group, seen)
	if err != nil || len(got.User) != 0 || !slices.Equal(got.Shadow, []string{"fs"}) || !reflect.DeepEqual(got.Folder, map[string][]string{"/p": {"loc"}}) {
		t.Fatalf("first: %+v %v", got, err)
	}
	writeClaudeJSON(t, dir, `{"mcpServers":{"fs":{"command":"mine"},"same":{"command":"s"}},"projects":{"/p":{"mcpServers":{"loc":{"command":"l"}}},"/q":{"mcpServers":{"loc":{"command":"l"}}}}}`)
	future := time.Now().Add(time.Second)
	os.Chtimes(filepath.Join(dir, ".claude.json"), future, future)
	got, _ = NewLocal(dir, group, seen)
	if len(got.User)+len(got.Shadow) != 0 || !reflect.DeepEqual(got.Folder, map[string][]string{"/q": {"loc"}}) {
		t.Fatalf("second: %+v", got)
	}
}

// Only a sign-in that holds a token counts; Claude Code writes empty
// entries for servers that wait for one.
func TestSignedIn(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"a"},"mcpOAuth":{`+
		`"plugin:context7:context7|0123456789abcdef":{"accessToken":"t","refreshToken":"r"},`+
		`"remote|fedcba9876543210":{"accessToken":"","discoveryState":{}},`+
		`"a|b|0011223344556677":{"refreshToken":"r"}}}`), 0o600)
	got := SignedIn(dir, t.TempDir())
	if !reflect.DeepEqual(got, map[string]bool{"plugin:context7:context7": true, "a|b": true}) {
		t.Fatalf("%v", got)
	}
	if got := SignedIn(t.TempDir(), t.TempDir()); len(got) != 0 {
		t.Fatalf("no login: %v", got)
	}
}
