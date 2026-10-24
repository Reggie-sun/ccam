package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/lock"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// On macOS Claude Code keeps the login in a Keychain item named after the
// folder's path, so a rename that moved the folder would log the account
// out (and leave the login for a later account of the old name): there the
// folder keeps its name. Tested here with the Keychain stood in for.
func TestRenameKeepsAFolderTheLoginIsNamedAfter(t *testing.T) {
	c := newTestCLI(t, map[string]int{"work": 10})
	orig := namedByFolder
	namedByFolder = func(string) string { return "the macOS Keychain" }
	t.Cleanup(func() { namedByFolder = orig })
	work := c.p.ProfileDir("work")
	out := c.ok("rename", "work", "job")
	if !strings.Contains(out, "its folder stays "+work+": the macOS Keychain keeps the login under the folder's path") {
		t.Fatalf("%s", out)
	}
	job, _ := c.cfg().Find("job")
	if profile.Dir(c.p, job) != work || fileThere(c.p.ProfileDir("job")) {
		t.Fatalf("folder %s", profile.Dir(c.p, job))
	}
	var rows []struct{ Name, Dir string }
	if err := json.Unmarshal([]byte(c.ok("list", "--json")), &rows); err != nil || !slices.Contains(rows, struct{ Name, Dir string }{"job", work}) {
		t.Fatalf("job runs elsewhere: %v %v", rows, err)
	}
	if _, errOut, code := c.run("", "add", "work", "--no-login"); code == 0 || !strings.Contains(errOut, "account job's") {
		t.Fatalf("a new work took job's folder: %d %s", code, errOut)
	}
	// Elsewhere the folder moves with the name, as before.
	namedByFolder = func(string) string { return "" }
	c.ok("rename", "job", "side2")
	if s, _ := c.cfg().Find("side2"); profile.Dir(c.p, s) != c.p.ProfileDir("side2") || fileThere(work) {
		t.Fatalf("folder %s", profile.Dir(c.p, s))
	}
}

// A Keychain item that an older ccam left behind when it removed or
// renamed an account would log a new account of that name in as the old
// one: ccam add stops and says how to remove it.
func TestAddStopsAtAKeychainLeftover(t *testing.T) {
	c := newTestCLI(t, nil)
	orig := keychainLeftover
	keychainLeftover = func(dir, home string) bool { return dir == c.p.ProfileDir("old") }
	t.Cleanup(func() { keychainLeftover = orig })
	_, errOut, code := c.run("", "add", "old", "--no-login")
	svc := creds.KeychainService(c.p.ProfileDir("old"), c.p.Home)
	if code == 0 || !strings.Contains(errOut, "security delete-generic-password -a "+creds.KeychainAccount()+" -s '"+svc+"'") {
		t.Fatalf("%d %s", code, errOut)
	}
	if fileThere(c.p.ProfileDir("old")) {
		t.Fatal("folder created")
	}
	c.ok("add", "new", "--no-login")
	if _, err := os.Stat(filepath.Join(c.p.ProfileDir("new"))); err != nil {
		t.Fatal(err)
	}
}

// Importing a folder copied from a registered one (cp -r ~/.claude
// ~/.claude-2) registers two accounts holding one refresh token; the note
// says so, and how to give the copy its own login. A separate login to the
// same Claude account is no copy.
func TestImportNotesACopiedLogin(t *testing.T) {
	c := newTestCLI(t, nil)
	login := func(name, rt string) string {
		dir := filepath.Join(c.p.Home, name)
		os.MkdirAll(dir, 0o700)
		os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"oauthAccount":{"accountUuid":"u1","emailAddress":"me@example.com"}}`), 0o600)
		os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"at","refreshToken":"`+rt+`"}}`), 0o600)
		return dir
	}
	if out := c.ok("import", login(".claude-one", "FAKE-RT-SAME")); strings.Contains(out, "same login") {
		t.Fatal(out)
	}
	out := c.ok("import", login(".claude-two", "FAKE-RT-SAME"))
	if !strings.Contains(out, "holds the same login as one (a copied folder)") || !strings.Contains(out, "`ccam login two`") || strings.Contains(out, "FAKE-RT") {
		t.Fatal(out)
	}
	if out := c.ok("import", login(".claude-three", "FAKE-RT-OTHER")); strings.Contains(out, "same login") {
		t.Fatalf("a login of its own: %s", out)
	}
}

// A Keychain item left for a folder counts when a new account of it would
// inherit something: a login (a dead one too) or MCP sign-ins, not the
// device key Claude Code's logout writes back.
func TestHoldsLogin(t *testing.T) {
	for data, want := range map[string]bool{
		`{"coworkRemoteDevice":{"k":"v"}}`:                                       false,
		`{"mcpOAuth":{"lin|aa":{"accessToken":"","discoveryState":{}}}}`:         false,
		`{"claudeAiOauth":{"accessToken":"at","refreshToken":"rt"}}`:             true,
		`{"claudeAiOauth":{"accessToken":"","refreshToken":"","expiresAt":0}}`:   true,
		`{"coworkRemoteDevice":{},"mcpOAuth":{"gh|aa":{"accessToken":"token"}}}`: true,
	} {
		if got := holdsLogin([]byte(data)); got != want {
			t.Errorf("%s: %v", data, got)
		}
	}
}

// While `ccam usage --refresh` renews a login, rename, remove and logout
// wait for it: the renewed login is saved into the folder the request
// was sent for.
func TestFolderChangesWaitForARefresh(t *testing.T) {
	c := newTestCLI(t, nil)
	l, err := lock.TryFlock(usage.RefreshLock(launcher.SessionDir(c.p, "work")))
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"rename", "work", "job"}, {"remove", "work", "--yes"}, {"logout", "work", "--yes"}} {
		if _, errOut, code := c.run("", args...); code == 0 || !strings.Contains(errOut, "another ccam is using the login of work right now") {
			t.Errorf("%v: %d %s", args, code, errOut)
		}
	}
	if !fileThere(c.p.ProfileDir("work")) {
		t.Fatal("the folder moved")
	}
	l.Unlock()
	c.ok("rename", "work", "job")
}
