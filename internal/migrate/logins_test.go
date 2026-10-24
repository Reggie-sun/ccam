package migrate

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/share"
)

// signInSource is ~/.claude with one user server (gh), its sign-in, a plugin
// server's sign-in, a sign-in of a server the source no longer has, and an
// entry without a token (a server that still needs a sign-in there).
func signInSource(t *testing.T) fixture {
	f := newFixture(t, map[string]string{
		".claude/CLAUDE.md":                      "x",
		".claude/plugins/installed_plugins.json": `{"plugins":{}}`,
		".claude.json":                           `{"mcpServers":{"gh":{"type":"http","url":"https://gh.example/mcp"},"lin":{"type":"http","url":"https://lin.example/mcp"}}}`,
		".claude/.credentials.json": `{"claudeAiOauth":{"accessToken":"main-login"},"mcpOAuth":{` +
			`"gh|aaaa":{"serverName":"gh","accessToken":"gh-token","refreshToken":"gh-refresh"},` +
			`"plugin:p:srv|bbbb":{"serverName":"plugin:p:srv","accessToken":"plugin-token"},` +
			`"gone|cccc":{"serverName":"gone","accessToken":"gone-token"},` +
			`"lin|dddd":{"serverName":"lin","accessToken":"","discoveryState":{"authorizationServerUrl":"https://lin.example"}}}}`,
	})
	return f
}

func account(t *testing.T, root, name, mcp, ext, creds string) Account {
	dir := filepath.Join(root, "profiles", name)
	os.MkdirAll(dir, 0o700)
	if creds != "" {
		write(t, dir, map[string]string{".credentials.json": creds})
	}
	return Account{Name: name, Dir: dir, MCP: mcp, Extensions: ext}
}

// A refresh token works for one holder: copying a sign-in into every
// account would let the first renewal log out all the other copies (the
// source's too). migrate writes no login file; it says which accounts must
// sign in to which servers.
func TestMCPSignInsAreNeverCopied(t *testing.T) {
	f := signInSource(t)
	root := t.TempDir()
	login := `{"claudeAiOauth":{"accessToken":"own"}}`
	servers := account(t, root, "servers", "default", "local", login)
	plugins := account(t, root, "plugins", "local", "default", login)
	signedIn := account(t, root, "signedin", "default", "local", `{"claudeAiOauth":{"accessToken":"own"},"mcpOAuth":{"gh|aaaa":{"accessToken":"own-token"}}}`)
	stub := account(t, root, "stub", "default", "local", `{"claudeAiOauth":{"accessToken":"own"},"mcpOAuth":{"gh|aaaa":{"accessToken":"","discoveryState":{}}}}`)
	loggedOut := account(t, root, "out", "default", "local", "")
	other := account(t, root, "other", "team", "team", login)
	accts := []Account{servers, plugins, signedIn, stub, loggedOut, other}
	before := map[string]string{}
	for _, a := range accts {
		before[a.Name] = readIfThere(filepath.Join(a.Dir, ".credentials.json"))
	}

	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: DefaultItems, Accounts: accts})
	if err != nil {
		t.Fatal(err)
	}
	var plan strings.Builder
	pl.Describe(&plan, false)
	want := "sign in with /mcp: gh in servers, stub, out; plugin:p:srv in plugins"
	if !strings.Contains(plan.String(), want) || strings.Contains(plan.String(), "lin,") {
		t.Fatalf("plan:\n%s", plan.String())
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	for _, a := range accts {
		if got := readIfThere(filepath.Join(a.Dir, ".credentials.json")); got != before[a.Name] {
			t.Fatalf("%s's login file changed:\n%s\n→\n%s", a.Name, before[a.Name], got)
		}
	}
	if s := pl.SignIns(); s != "gh in servers, stub, out; plugin:p:srv in plugins" {
		t.Fatalf("summary %q", s)
	}
	if strings.Contains(strings.Join(pl.Copied(), ","), "sign-in") {
		t.Fatalf("copied: %v", pl.Copied())
	}
	// No token may reach the shared group.
	filepath.WalkDir(f.group, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.Contains(read(t, path), "-token") {
			t.Fatalf("token in the group: %s", path)
		}
		return nil
	})
}

func readIfThere(path string) string {
	data, _ := os.ReadFile(path)
	return string(data)
}

// An entry without a token is no sign-in: there is nothing to sign in to
// again for it, and nothing to say.
func TestPlaceholderSignInsAreNoSignIns(t *testing.T) {
	f := newFixture(t, map[string]string{
		".claude/CLAUDE.md":         "x",
		".claude.json":              `{"mcpServers":{"lin":{"type":"http","url":"https://lin.example/mcp"}}}`,
		".claude/.credentials.json": `{"mcpOAuth":{"lin|dddd":{"accessToken":"","discoveryState":{}}}}`,
	})
	acct := account(t, t.TempDir(), "a", "default", "local", `{"claudeAiOauth":{"accessToken":"own"}}`)
	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: DefaultItems, Accounts: []Account{acct}})
	if err != nil {
		t.Fatal(err)
	}
	var plan strings.Builder
	pl.Describe(&plan, false)
	if pl.Logins != nil || strings.Contains(plan.String(), "sign-ins") {
		t.Fatalf("plan:\n%s", plan.String())
	}
}

func TestPluginSignInsNeedThePlugins(t *testing.T) {
	f := signInSource(t)
	os.RemoveAll(filepath.Join(f.src, "plugins"))
	acct := account(t, t.TempDir(), "a", "local", "default", `{"claudeAiOauth":{"accessToken":"own"}}`)
	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: DefaultItems, Accounts: []Account{acct}})
	if err != nil {
		t.Fatal(err)
	}
	if s := pl.SignIns(); s != "" {
		t.Fatalf("plugin sign-ins without plugins: %q", s)
	}
}

func TestMCPSignInNotNeededWhenTheGroupKeepsItsDefinition(t *testing.T) {
	f := signInSource(t)
	write(t, f.group, map[string]string{"mcp.json": `{"mcpServers":{"gh":{"type":"http","url":"https://other.example/mcp"}}}`})
	acct := account(t, t.TempDir(), "a", "default", "local", `{"claudeAiOauth":{"accessToken":"own"}}`)
	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: []string{"mcp"}, Accounts: []Account{acct}})
	if err != nil {
		t.Fatal(err)
	}
	var plan strings.Builder
	pl.Describe(&plan, true)
	if !strings.Contains(plan.String(), "sign in with /mcp: gh in a") {
		t.Fatalf("before the choice, the plan says what \"use mine\" would need:\n%s", plan.String())
	}
	for _, c := range pl.Conflicts {
		c.Choice = share.UseMine
	}
	if s := pl.SignIns(); s != "gh in a" {
		t.Fatalf("use mine: %q", s)
	}
	for _, c := range pl.Conflicts {
		c.Choice = share.KeepGroup
	}
	// Claude Code keys the sign-in by the definition: the group's server is
	// another one, which the source never signed in to. With the choice
	// made (--yes keeps the group's), the plan does not list it either.
	plan.Reset()
	pl.Describe(&plan, false)
	if s := pl.SignIns(); s != "" || strings.Contains(plan.String(), "sign-ins") {
		t.Fatalf("got %q:\n%s", s, plan.String())
	}
}

// Sign-ins alone are nothing to copy; the plan still says which account
// must sign in.
func TestMCPSignInsAloneAreNothingToCopy(t *testing.T) {
	f := signInSource(t)
	write(t, f.group, map[string]string{"mcp.json": `{"mcpServers":{"gh":{"type":"http","url":"https://gh.example/mcp"},"lin":{"type":"http","url":"https://lin.example/mcp"}}}`})
	acct := account(t, t.TempDir(), "a", "default", "local", `{"claudeAiOauth":{"accessToken":"own"}}`)
	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: []string{"mcp"}, Accounts: []Account{acct}})
	if err != nil {
		t.Fatal(err)
	}
	if !pl.Empty() || pl.SignIns() != "gh in a" {
		t.Fatalf("empty %v, sign-ins %q", pl.Empty(), pl.SignIns())
	}
}
