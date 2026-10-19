package migrate

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/share"
)

// signInSource is ~/.claude with one user server (gh), its sign-in, a plugin
// server's sign-in, and a sign-in of a server the source no longer has.
func signInSource(t *testing.T) fixture {
	f := newFixture(t, map[string]string{
		".claude/CLAUDE.md":                      "x",
		".claude/plugins/installed_plugins.json": `{"plugins":{}}`,
		".claude.json":                           `{"mcpServers":{"gh":{"type":"http","url":"https://gh.example/mcp"}}}`,
		".claude/.credentials.json": `{"claudeAiOauth":{"accessToken":"main-login"},"mcpOAuth":{` +
			`"gh|aaaa":{"serverName":"gh","accessToken":"gh-token"},` +
			`"plugin:p:srv|bbbb":{"serverName":"plugin:p:srv","accessToken":"plugin-token"},` +
			`"gone|cccc":{"serverName":"gone","accessToken":"gone-token"}}}`,
	})
	return f
}

func account(t *testing.T, root, name, mcp, ext, creds string) Account {
	dir := filepath.Join(root, "profiles", name)
	os.MkdirAll(dir, 0o700)
	if creds != "" {
		write(t, dir, map[string]string{".credentials.json": creds})
	}
	return Account{Name: name, Dir: dir, MCP: mcp, Extensions: ext, SessionDir: filepath.Join(root, "run", name)}
}

func signInsOf(t *testing.T, dir string) map[string]map[string]string {
	t.Helper()
	var doc struct {
		ClaudeAiOauth map[string]string            `json:"claudeAiOauth"`
		MCPOAuth      map[string]map[string]string `json:"mcpOAuth"`
	}
	if err := json.Unmarshal([]byte(read(t, filepath.Join(dir, ".credentials.json"))), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.ClaudeAiOauth["accessToken"] == "" {
		t.Fatalf("%s lost its own login", dir)
	}
	return doc.MCPOAuth
}

func TestMCPSignInsGoToEachAccountNotTheGroup(t *testing.T) {
	f := signInSource(t)
	root := t.TempDir()
	login := `{"claudeAiOauth":{"accessToken":"own"}}`
	servers := account(t, root, "servers", "default", "local", login)
	plugins := account(t, root, "plugins", "local", "default", login)
	newer := account(t, root, "newer", "default", "local", `{"claudeAiOauth":{"accessToken":"own"},"mcpOAuth":{"gh|aaaa":{"accessToken":"newer-token"}}}`)
	loggedOut := account(t, root, "out", "default", "default", "")
	imported := account(t, root, "imported", "default", "default", login)
	imported.Imported = true
	other := account(t, root, "other", "team", "team", login)

	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default",
		Items: DefaultItems, Accounts: []Account{servers, plugins, newer, loggedOut, imported, other}})
	if err != nil {
		t.Fatal(err)
	}
	var plan strings.Builder
	pl.Describe(&plan, false)
	if !strings.Contains(plan.String(), "sign-ins") || !strings.Contains(plan.String(), "→ servers, plugins") ||
		!strings.Contains(plan.String(), "out (not logged in yet)") || !strings.Contains(plan.String(), "imported (registered in place") {
		t.Fatalf("plan:\n%s", plan.String())
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}

	if got := signInsOf(t, servers.Dir); len(got) != 1 || got["gh|aaaa"]["accessToken"] != "gh-token" {
		t.Fatalf("servers got %v (only the copied server's sign-in)", got)
	}
	if got := signInsOf(t, plugins.Dir); len(got) != 1 || got["plugin:p:srv|bbbb"]["accessToken"] != "plugin-token" {
		t.Fatalf("plugins got %v", got)
	}
	if got := signInsOf(t, newer.Dir); got["gh|aaaa"]["accessToken"] != "newer-token" {
		t.Fatalf("an account's own sign-in was replaced: %v", got)
	}
	for _, a := range []Account{imported, other} {
		if strings.Contains(read(t, filepath.Join(a.Dir, ".credentials.json")), "mcpOAuth") {
			t.Fatalf("%s got sign-ins", a.Name)
		}
	}
	if info, _ := os.Stat(filepath.Join(servers.Dir, ".credentials.json")); info.Mode().Perm() != 0o600 {
		t.Fatalf("login file mode %v", info.Mode())
	}
	if strings.Join(pl.Copied(), " ") == "" || !strings.Contains(strings.Join(pl.Copied(), ","), "MCP sign-ins for servers, plugins") {
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

func TestPluginSignInsNeedThePlugins(t *testing.T) {
	f := signInSource(t)
	os.RemoveAll(filepath.Join(f.src, "plugins"))
	acct := account(t, t.TempDir(), "a", "local", "default", `{"claudeAiOauth":{"accessToken":"own"}}`)
	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: DefaultItems, Accounts: []Account{acct}})
	if err != nil {
		t.Fatal(err)
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	if got := signInsOf(t, acct.Dir); len(got) != 0 {
		t.Fatalf("plugin sign-ins copied without plugins: %v", got)
	}
}

func TestMCPSignInSkippedWhenTheGroupKeepsItsDefinition(t *testing.T) {
	f := signInSource(t)
	write(t, f.group, map[string]string{"mcp.json": `{"mcpServers":{"gh":{"type":"http","url":"https://other.example/mcp"}}}`})
	acct := account(t, t.TempDir(), "a", "default", "local", `{"claudeAiOauth":{"accessToken":"own"}}`)
	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: []string{"mcp"}, Accounts: []Account{acct}})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range pl.Conflicts {
		c.Choice = share.KeepGroup
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	// Claude Code keys the sign-in by the definition, so the source's would
	// not match the group's server.
	if got := signInsOf(t, acct.Dir); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestMCPSignInsAloneAreSomethingToCopy(t *testing.T) {
	f := signInSource(t)
	write(t, f.group, map[string]string{"mcp.json": `{"mcpServers":{"gh":{"type":"http","url":"https://gh.example/mcp"}}}`})
	acct := account(t, t.TempDir(), "a", "default", "local", `{"claudeAiOauth":{"accessToken":"own"}}`)
	pl, err := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: []string{"mcp"}, Accounts: []Account{acct}})
	if err != nil {
		t.Fatal(err)
	}
	if pl.Empty() {
		t.Fatal("the server is already in the group, but the account still lacks its sign-in")
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	if got := signInsOf(t, acct.Dir); len(got) != 1 {
		t.Fatalf("got %v", got)
	}
	if pl, _ := Build(Options{Source: f.src, Home: f.home, GroupsRoot: f.groups, Group: "default", Items: []string{"mcp"}, Accounts: []Account{acct}}); !pl.Empty() {
		t.Fatal("second run should have nothing to copy")
	}
}
