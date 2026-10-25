package doctor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/ext"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/group"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/share"
	"github.com/jhyunwoo/ccam/internal/shim"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func TestPluginPathsRewrittenToGroup(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	pluginsDir := filepath.Join(p.GroupDir("default"), "plugins")
	os.MkdirAll(pluginsDir, 0o700)
	viaProfile := filepath.Join(p.ProfileDir("a"), "plugins", "cache", "mk", "demo", "1.0")
	file := filepath.Join(pluginsDir, "installed_plugins.json")
	install, _ := json.Marshal(viaProfile)
	os.WriteFile(file, []byte(`{"plugins":{"demo@mk":[{"installPath":`+string(install)+`}]}}`), 0o600)

	r := &report{}
	checkPlugins(r, p, cfg, false)
	if len(r.out) != 1 || r.out[0].Level != Warn || r.out[0].Fixed {
		t.Fatalf("findings %v", r.out)
	}
	r = &report{}
	checkPlugins(r, p, cfg, true)
	if !r.out[0].Fixed {
		t.Fatal("not fixed")
	}
	data, _ := os.ReadFile(file)
	want, _ := json.Marshal(filepath.Join(pluginsDir, "cache", "mk", "demo", "1.0"))
	if !strings.Contains(string(data), string(want)) {
		t.Fatalf("rewritten: %s", data)
	}
}

func TestCredentialsChecks(t *testing.T) {
	dir := t.TempDir()
	prof := &config.Profile{Name: "a"}
	cred := filepath.Join(dir, ".credentials.json")
	os.WriteFile(cred, []byte(`{"claudeAiOauth":{"expiresAt":1,"refreshTokenExpiresAt":2}}`), 0o644)
	r := &report{}
	checkCredentials(r, prof, dir, "/nonexistent-home", true)
	if runtime.GOOS == "windows" { // no permission bits to check
		if len(r.out) != 1 || r.out[0].Level != Error {
			t.Fatalf("%v", r.out)
		}
		return
	}
	if len(r.out) != 2 || !r.out[0].Fixed || r.out[1].Level != Error {
		t.Fatalf("%v", r.out)
	}
	if info, _ := os.Stat(cred); info.Mode().Perm() != 0o600 {
		t.Fatal("perm not fixed")
	}
}

func TestAuthEnvWarnsWithoutValue(t *testing.T) {
	t.Setenv("CLAUDE_CODE_OAUTH_TOKEN", "secret-token")
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "")
	r := &report{}
	checkAuthEnv(r)
	if len(r.out) != 1 || r.out[0].Level != Warn || strings.Contains(r.out[0].Msg, "secret") {
		t.Fatalf("findings %v", r.out)
	}
}

func TestFolderChecks(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	cfg.SetFolder(dir, "a")
	cfg.SetFolder(filepath.Join(dir, "gone"), "a")
	cfg.SetFolder(filepath.Join(dir, "x"), "nobody")
	r := &report{}
	checkFolders(r, cfg)
	if len(r.out) != 2 || r.out[0].Level != Warn || !strings.Contains(r.out[0].Msg, "gone") ||
		r.out[1].Level != Error || !strings.Contains(r.out[1].Msg, "nobody") {
		t.Fatalf("%v", r.out)
	}
}

func TestGroupChecks(t *testing.T) {
	cfg := &config.Config{Theme: "solarized", Profiles: []config.Profile{{Name: "a", Slot: 1}, {Name: "group", Slot: 2}}}
	cfg.Groups = []config.Group{{Name: "default", Items: []string{"mcp"}, Members: []string{"a", "gone"}}}
	r := &report{}
	checkGroups(r, cfg)
	all := fmt.Sprint(r.out)
	for _, want := range []string{"group default lists gone", `theme "solarized"`, "`ccam group` is now a ccam command"} {
		if !strings.Contains(all, want) {
			t.Errorf("missing %q in %v", want, r.out)
		}
	}
	r = &report{}
	checkGroups(r, &config.Config{Theme: "dark-ansi"})
	if len(r.out) != 0 {
		t.Fatalf("%v", r.out)
	}
}

// A setting written below a table is read as the table's, so ccam does
// not see it, and refuses to save over it: doctor says where it belongs.
func TestConfigLinesCcamDoesNotRead(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	os.WriteFile(path, []byte("[[profile]]\nname = \"a\"\nslot = 1\n\nrelay = true\n"), 0o600)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	r := &report{}
	checkGroups(r, cfg)
	if len(r.out) != 1 || r.out[0].Level != Warn || r.out[0].Area != "Config" ||
		!strings.Contains(r.out[0].Msg, "config.toml line 5: relay is read as part of the [[profile]] table above it") {
		t.Fatalf("%v", r.out)
	}
}

func TestIdentityChecks(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{Profiles: []config.Profile{
		{Name: "a", Slot: 1, Login: &config.Identity{AccountUUID: "u1", Email: "a@x"}},
		{Name: "b", Slot: 2},
		{Name: "c", Slot: 3, Login: &config.Identity{AccountUUID: "u3", Email: "c@x"}},
	}}
	write := func(name, uuid, email string) {
		dir := p.ProfileDir(name)
		os.MkdirAll(dir, 0o700)
		os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"oauthAccount":{"accountUuid":"`+uuid+`","emailAddress":"`+email+`","organizationUuid":"o"}}`), 0o600)
	}
	write("a", "u1", "a@x")
	write("b", "u1", "a@x")
	write("c", "u9", "other@x")
	r := &report{}
	checkIdentities(r, p, cfg)
	if len(r.out) != 2 || r.out[0].Profile != "b" || !strings.HasPrefix(r.out[0].Msg, "logged in to the same Claude account as a") ||
		r.out[1].Profile != "c" || !strings.HasPrefix(r.out[1].Msg, "now logged in as other@x, but ccam recorded c@x") {
		t.Fatalf("%v", r.out)
	}
}

// An account ccam has no login record of (an older ccam's, an imported
// one) is checked against the email it recorded, and the login line names
// whom the folder is logged in as now.
func TestIdentityChecksWithoutALoginRecord(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{Profiles: []config.Profile{{Name: "d", Slot: 1, Email: "d@x"}, {Name: "e", Slot: 2, Email: "e@x"}}}
	for name, email := range map[string]string{"d": "other@x", "e": "E@x"} {
		dir := p.ProfileDir(name)
		os.MkdirAll(dir, 0o700)
		os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"oauthAccount":{"accountUuid":"u-`+name+`","emailAddress":"`+email+`","organizationUuid":"o"}}`), 0o600)
		os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"t","refreshToken":"r-`+name+`","expiresAt":4102444800000}}`), 0o600)
	}
	r := &report{}
	checkIdentities(r, p, cfg)
	if len(r.out) != 1 || r.out[0].Profile != "d" || !strings.HasPrefix(r.out[0].Msg, "now logged in as other@x, but ccam recorded d@x") {
		t.Fatalf("%v", r.out)
	}
	r = &report{}
	checkCredentials(r, &cfg.Profiles[0], p.ProfileDir("d"), p.Home, false)
	if len(r.out) != 1 || !strings.HasPrefix(r.out[0].Msg, "login valid (other@x") {
		t.Fatalf("%v", r.out)
	}
}

// A login Claude Code has cleared (both tokens emptied after the server
// refused the refresh token) cannot be renewed on the next launch, and a
// logout leaves a login file without a login.
func TestClearedLoginIsExpired(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("macOS asks the Keychain first")
	}
	prof := &config.Profile{Name: "a"}
	later := fmt.Sprint(time.Now().Add(20 * 24 * time.Hour).UnixMilli())
	for data, want := range map[string]Level{
		`{"claudeAiOauth":{"accessToken":"","refreshToken":"","expiresAt":0,"refreshTokenExpiresAt":` + later + `}}`: Error,
		`{"coworkRemoteDevice":{"k":"v"}}`: Warn,
		`{"claudeAiOauth":{"accessToken":"at","refreshToken":"rt","expiresAt":1,"refreshTokenExpiresAt":` + later + `}}`: Info,
	} {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(data), 0o600)
		r := &report{}
		checkCredentials(r, prof, dir, "/nonexistent-home", false)
		if len(r.out) != 1 || r.out[0].Level != want || want != Info && !strings.Contains(r.out[0].Msg, "ccam login a") {
			t.Errorf("%s: %v", data, r.out)
		}
	}
}

// Folders that are copies of one another hold one refresh token: whichever
// renews it first logs the other out. Two logins to one account do not.
func TestCopiedLoginIsAnError(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}, {Name: "b", Slot: 2}, {Name: "c", Slot: 3}}}
	for name, rt := range map[string]string{"a": "FAKE-RT-1", "b": "FAKE-RT-1", "c": "FAKE-RT-2"} {
		dir := p.ProfileDir(name)
		os.MkdirAll(dir, 0o700)
		os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"oauthAccount":{"accountUuid":"u1","emailAddress":"a@x","organizationUuid":"o"}}`), 0o600)
		os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"at","refreshToken":"`+rt+`"}}`), 0o600)
	}
	r := &report{}
	checkIdentities(r, p, cfg)
	if len(r.out) != 2 || r.out[0].Profile != "b" || r.out[0].Level != Error || !strings.Contains(r.out[0].Msg, "same login token as a") ||
		!strings.Contains(r.out[0].Msg, "`ccam login b`") || r.out[1].Profile != "c" || r.out[1].Level != Warn {
		t.Fatalf("%v", r.out)
	}
	for _, f := range r.out {
		if strings.Contains(f.Msg, "FAKE-RT") {
			t.Fatal("a token in the report")
		}
	}
}

// A missing share link says how it gets made, unless --fix makes it now.
func TestMissingLinkSaysHowItIsMade(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	if err := cfg.AddGroup("default", []string{"memory"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(p.ProfileDir("a"), 0o700)
	os.MkdirAll(p.GroupDir("default"), 0o700)
	os.WriteFile(filepath.Join(p.GroupDir("default"), "CLAUDE.md"), []byte("shared"), 0o600)
	line := func(fix bool) string {
		for _, f := range Run(p, cfg, Options{Fix: fix}) {
			if f.Profile == "a" && strings.Contains(f.Msg, "CLAUDE.md") {
				return f.Msg
			}
		}
		return ""
	}
	if got := line(false); !strings.Contains(got, "share link is missing") || !strings.HasSuffix(got, "(the next launch or ccam doctor --fix makes it)") {
		t.Fatalf("%q", got)
	}
	if got := line(true); got == "" || strings.Contains(got, "makes it") {
		t.Fatalf("with --fix: %q", got)
	}
}

// The dashboard's Health tab groups findings by account and area.
func TestFindingsSayTheirArea(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	os.MkdirAll(p.ProfileDir("a"), 0o700)
	areas := map[string]bool{}
	for _, f := range Run(p, cfg, Options{}) {
		if f.Profile == "a" {
			areas[f.Area] = true
		}
		if f.Area == "" {
			t.Errorf("no area: %v", f)
		}
	}
	if !areas["Login"] || !areas["Share links"] {
		t.Fatalf("areas of a: %v", areas)
	}
}

// The dashboard repairs one box at a time: one account's findings, or with
// General those of no account.
func TestFixOnlyOneAccount(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}, {Name: "b", Slot: 2}}}
	for _, n := range []string{"a", "b"} {
		os.MkdirAll(p.ProfileDir(n), 0o755)
	}
	if runtime.GOOS == "windows" {
		t.Skip("no permission bits to repair")
	}
	perm := func(n string) os.FileMode { info, _ := os.Stat(p.ProfileDir(n)); return info.Mode().Perm() }
	Run(p, cfg, Options{Fix: true, FixOnly: General})
	if perm("a") != 0o755 || perm("b") != 0o755 {
		t.Fatal("General repaired an account")
	}
	Run(p, cfg, Options{Fix: true, FixOnly: "b"})
	if perm("a") != 0o755 || perm("b") != 0o700 {
		t.Fatalf("a %o, b %o", perm("a"), perm("b"))
	}
}

// Servers plain Claude Code has in ~/.claude.json that no account loads
// are reported with the command that copies them; those a group or an
// account has are not.
func TestHomeMCPServersNoAccountLoads(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}, {Name: "b", Slot: 2}}}
	if err := cfg.AddGroup("default", []string{"mcp"}, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	cfg.Resolve()
	os.MkdirAll(filepath.Join(p.Home, ".claude"), 0o700)
	os.WriteFile(filepath.Join(p.Home, ".claude.json"), []byte(`{"mcpServers":{"gone":{"command":"g"},"shared":{"command":"s"},"own":{"command":"o"}}}`), 0o600)
	os.MkdirAll(p.GroupDir("default"), 0o700)
	os.WriteFile(filepath.Join(p.GroupDir("default"), "mcp.json"), []byte(`{"mcpServers":{"shared":{"command":"s"}}}`), 0o600)
	os.MkdirAll(p.ProfileDir("b"), 0o700)
	os.WriteFile(filepath.Join(p.ProfileDir("b"), ".claude.json"), []byte(`{"mcpServers":{"own":{"command":"o"}}}`), 0o600)
	r := &report{}
	checkHomeMCP(r, p, cfg)
	if len(r.out) != 1 || r.out[0].Level != Warn || !strings.Contains(r.out[0].Msg, "no account loads: gone →") || !strings.Contains(r.out[0].Msg, "`ccam migrate --mcp`") {
		t.Fatalf("%v", r.out)
	}
}

// Servers only one account of a group has, and servers another account is
// signed in to, are said for the account that lacks them.
func TestAccountMCPGaps(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}, {Name: "b", Slot: 2}}}
	if err := cfg.AddGroup("default", []string{"mcp"}, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	cfg.Resolve()
	os.MkdirAll(p.GroupDir("default"), 0o700)
	os.WriteFile(filepath.Join(p.GroupDir("default"), "mcp.json"), []byte(`{"mcpServers":{"remote":{"type":"http","url":"https://r"}}}`), 0o600)
	for _, n := range []string{"a", "b"} {
		os.MkdirAll(p.ProfileDir(n), 0o700)
	}
	os.WriteFile(filepath.Join(p.ProfileDir("a"), ".claude.json"), []byte(`{"mcpServers":{"x":{"command":"x"}},"projects":{"/proj":{"mcpServers":{"y":{"command":"y"}}}}}`), 0o600)
	os.WriteFile(filepath.Join(p.ProfileDir("a"), ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"t"},"mcpOAuth":{"remote|0123456789abcdef":{"accessToken":"t","refreshToken":"r"}}}`), 0o600)
	os.WriteFile(filepath.Join(p.ProfileDir("b"), ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"t"},"mcpOAuth":{"remote|0123456789abcdef":{"accessToken":"","discoveryState":{}}}}`), 0o600)
	// A plugin of a's own with an MCP server: it comes with the plugin, which
	// `ccam mcp import` cannot share.
	installPlugin(t, p.ProfileDir("a"), "demo@mk", `{"srv":{"command":"s"}}`)
	r := &report{}
	checkAccountMCP(r, p, cfg)
	var got []string
	for _, f := range r.out {
		got = append(got, fmt.Sprintf("%d %s: %s", f.Level, f.Profile, f.Msg))
	}
	all := strings.Join(got, "\n")
	for _, want := range []string{
		fmt.Sprintf("%d a: MCP servers the other accounts of group default do not load: x; for one folder only: y (/proj) → `ccam mcp import a x y`", Info),
		fmt.Sprintf("%d b: not signed in to MCP servers other accounts are signed in to: remote (a)", Warn),
	} {
		if !strings.Contains(all, want) {
			t.Errorf("missing %q in\n%s", want, all)
		}
	}
	if len(got) != 2 {
		t.Errorf("findings:\n%s", all)
	}
}

// installPlugin installs a plugin, turned on, in an account folder's own
// plugins folder, with the MCP servers of mcpJSON (its .mcp.json).
func installPlugin(t *testing.T, dir, name, mcpJSON string) {
	t.Helper()
	install := filepath.Join(dir, "plugins", "cache", "mk", "demo", "1.0")
	if err := os.MkdirAll(install, 0o700); err != nil {
		t.Fatal(err)
	}
	ip, _ := json.Marshal(install)
	for path, data := range map[string]string{
		filepath.Join(install, ".mcp.json"):                     mcpJSON,
		filepath.Join(dir, "plugins", "installed_plugins.json"): `{"version":2,"plugins":{"` + name + `":[{"scope":"user","installPath":` + string(ip) + `}]}}`,
		filepath.Join(dir, "settings.json"):                     `{"enabledPlugins":{"` + name + `":true}}`,
	} {
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// A share link whose group folder is gone is reported, as ccam ext tells
// to run doctor --fix, and --fix makes the folder again, also while the
// account runs: that loses nothing.
func TestMissingGroupFolderIsMadeAgain(t *testing.T) {
	for _, running := range []bool{false, true} {
		t.Setenv("CCAM_HOME", t.TempDir())
		p := paths.FromEnv()
		// A claude that is not there: doctor runs no real one.
		cfg := &config.Config{ClaudeBin: filepath.Join(p.Data, "no-claude"), Profiles: []config.Profile{{Name: "a", Slot: 1}}}
		if err := cfg.AddGroup("default", []string{"extensions", "history"}, []string{"a"}); err != nil {
			t.Fatal(err)
		}
		dir := p.ProfileDir("a")
		os.MkdirAll(dir, 0o700)
		pl, err := share.Build(dir, p.GroupsRoot(), config.Share{}, cfg.Profiles[0].Share, share.Options{})
		if err == nil {
			err = pl.Apply()
		}
		if err != nil {
			t.Fatal(err)
		}
		if running {
			if err := session.Register(launcher.SessionDir(p, "a")); err != nil {
				t.Fatal(err)
			}
		}
		for _, name := range []string{"skills", "projects"} {
			os.Remove(filepath.Join(p.GroupDir("default"), name))
		}
		links := func(fix bool) []Finding {
			var out []Finding
			for _, f := range Run(p, cfg, Options{Fix: fix}) {
				if f.Area == "Share links" && f.Profile == "a" {
					out = append(out, f)
				}
			}
			return out
		}
		if got := links(false); len(got) != 2 || got[0].Level != Warn || !strings.Contains(got[0].Msg, "missing") || !strings.HasSuffix(got[0].Msg, "(ccam doctor --fix makes it)") {
			t.Fatalf("running %v: %v", running, got)
		}
		got := links(true)
		if len(got) < 2 || !got[0].Fixed || !got[1].Fixed {
			t.Fatalf("running %v, --fix: %v", running, got)
		}
		for _, name := range []string{"skills", "projects"} {
			if info, err := os.Stat(filepath.Join(dir, name)); err != nil || !info.IsDir() {
				t.Fatalf("running %v: %s: %v", running, name, err)
			}
		}
		if s := ext.ScanAccount("a", dir, p.GroupsRoot(), cfg.Profiles[0].Share, nil, ""); len(s.Problems) != 0 {
			t.Fatalf("running %v: ext still reports %v", running, s.Problems)
		}
		session.Unregister(launcher.SessionDir(p, "a"))
	}
}

// doctor --fix links a shared path whose own copy differs from the
// group's: the account then sees the group's copy. It says so, and where
// the account's copy went: a real settings.json, a later account's copy
// of a path shared since an upgrade (keybindings.json, which the first
// account's copy fills the group with), a skill in a real skills folder.
// A dry run says what --fix would do and changes nothing.
func TestFixSaysWhatItSetAside(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{ClaudeBin: filepath.Join(p.Data, "no-claude"), Profiles: []config.Profile{{Name: "a", Slot: 1}, {Name: "b", Slot: 2}}}
	if err := cfg.AddGroup("default", []string{"settings", "extensions"}, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	group := p.GroupDir("default")
	os.MkdirAll(filepath.Join(group, "skills", "x"), 0o700)
	os.WriteFile(filepath.Join(group, "settings.json"), []byte(`{"group":true}`), 0o600)
	os.WriteFile(filepath.Join(group, "skills", "x", "SKILL.md"), []byte("group"), 0o600)
	for _, n := range []string{"a", "b"} {
		dir := p.ProfileDir(n)
		os.MkdirAll(dir, 0o700)
		os.WriteFile(filepath.Join(dir, "keybindings.json"), []byte(`{"`+n+`":1}`), 0o600)
	}
	a, b := p.ProfileDir("a"), p.ProfileDir("b")
	os.WriteFile(filepath.Join(a, "settings.json"), []byte(`{"mine":true}`), 0o600)
	os.MkdirAll(filepath.Join(a, "skills", "x"), 0o700)
	os.WriteFile(filepath.Join(a, "skills", "x", "SKILL.md"), []byte("mine"), 0o600)
	os.Symlink(filepath.Join(group, "settings.json"), filepath.Join(b, "settings.json"))

	// said reports the finding of account n that contains text, and the
	// path that ends it.
	said := func(fs []Finding, n, text string) (string, string) {
		t.Helper()
		for _, f := range fs {
			if f.Profile == n && strings.Contains(f.Msg, text) {
				return f.Msg, f.Msg[strings.LastIndex(f.Msg, " ")+1:]
			}
		}
		t.Fatalf("no finding of %s says %q: %v", n, text, fs)
		return "", ""
	}
	dry := Run(p, cfg, Options{DryRun: true})
	said(dry, "a", "settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to")
	said(dry, "b", "keybindings.json differs from the group's copy: --fix would use the group's")
	said(dry, "a", "skills: 1 file differs from the group's copy ("+filepath.Join("x", "SKILL.md")+")")
	if got := readFile(t, filepath.Join(a, "settings.json")); got != `{"mine":true}` {
		t.Fatalf("a dry run changed settings.json: %s", got)
	}

	fixed := Run(p, cfg, Options{Fix: true})
	_, saved := said(fixed, "a", "settings.json differed from the group's copy, which it now uses; this account's copy is in")
	if got := readFile(t, saved); got != `{"mine":true}` {
		t.Fatalf("a's copy in %s: %s", saved, got)
	}
	_, saved = said(fixed, "b", "keybindings.json differed from the group's copy, which it now uses; this account's copy is in")
	if got := readFile(t, saved); got != `{"b":1}` {
		t.Fatalf("b's copy in %s: %s", saved, got)
	}
	said(fixed, "a", "skills: 1 file differs from the group's copy")
	if got := readFile(t, filepath.Join(a, "settings.json")); got != `{"group":true}` {
		t.Fatalf("a's settings.json: %s", got)
	}
	if got := readFile(t, filepath.Join(b, "keybindings.json")); got != `{"a":1}` {
		t.Fatalf("b's keybindings.json: %s", got)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// The usage format check has one time limit for all accounts: on a
// network that never answers it took 3 s per account. When no account
// answers it says so.
func TestUsageFormatCheckIsBounded(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{}
	for i, name := range []string{"a", "b", "c", "d", "e"} {
		cfg.Profiles = append(cfg.Profiles, config.Profile{Name: name, Slot: i + 1})
		dir := p.ProfileDir(name)
		os.MkdirAll(dir, 0o700)
		os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(fmt.Sprintf(`{"claudeAiOauth":{"accessToken":"t","refreshToken":"r","expiresAt":%d}}`, time.Now().Add(time.Hour).UnixMilli())), 0o600)
	}
	for name, rt := range map[string]http.RoundTripper{
		"no answer": roundTrip(func(r *http.Request) (*http.Response, error) { <-r.Context().Done(); return nil, r.Context().Err() }),
		"refused":   roundTrip(func(*http.Request) (*http.Response, error) { return nil, fmt.Errorf("connection refused") }),
	} {
		f := usage.NewFetcher(filepath.Join(p.State, "usage", name))
		f.Client = &http.Client{Transport: rt}
		r := &report{}
		start := time.Now()
		checkUsageFormat(r, p, cfg, f)
		if d := time.Since(start); d > usage.LookupTimeout+2*time.Second {
			t.Errorf("%s: took %v", name, d)
		}
		if len(r.out) != 1 || r.out[0].Level != Info || !strings.Contains(r.out[0].Msg, "could not check the usage response format") {
			t.Errorf("%s: %+v", name, r.out)
		}
	}
}

type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// The entries Claude Code 2.1.309 keeps in a config folder are known; a
// name it does not is reported with the version ccam knows.
func TestUnknownEntriesOfClaudeCode2_1_309(t *testing.T) {
	dir := t.TempDir()
	for _, n := range []string{"policy-limits.json.signature-iat.json", "remote-settings.json.signature-iat.json", "stats-cache.json",
		"usage-data", "teams", "logs", ".claude.json.backup", "remote-settings.json.next-suffix"} {
		os.WriteFile(filepath.Join(dir, n), nil, 0o600)
	}
	r := &report{}
	checkUnknownEntries(r, &config.Profile{Name: "a"}, dir)
	if len(r.out) != 0 {
		t.Fatalf("%+v", r.out)
	}
	os.WriteFile(filepath.Join(dir, "x.json"), nil, 0o600)
	checkUnknownEntries(r, &config.Profile{Name: "a"}, dir)
	if len(r.out) != 1 || !strings.Contains(r.out[0].Msg, "x.json") || !strings.Contains(r.out[0].Msg, LayoutVersion) {
		t.Fatalf("%+v", r.out)
	}
}

// A stopped change of groups that cannot be finished yet (one of its
// accounts runs) leaves its accounts' links alone: repairing them in config
// order would decide the group's first content against the change (--from).
func TestFixWaitsForAStoppedChangeOfGroups(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{ClaudeBin: stubClaude}
	for i, n := range []string{"a", "b", "c"} {
		cfg.Profiles = append(cfg.Profiles, config.Profile{Name: n, Slot: i + 1})
		os.MkdirAll(p.ProfileDir(n), 0o700)
		os.WriteFile(filepath.Join(p.ProfileDir(n), "CLAUDE.md"), []byte(n), 0o600)
	}
	if err := cfg.AddGroup("h", []string{"memory"}, nil); err != nil {
		t.Fatal(err)
	}
	next := cfg.Clone()
	for _, n := range []string{"a", "b", "c"} {
		if _, err := next.JoinGroup(n, "h", true); err != nil {
			t.Fatal(err)
		}
	}
	ch, err := group.Build(p, cfg, next, group.Options{Seed: "b"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	// ccam stopped right after saving the config: the note stays.
	note := mustRead(t, group.NotePath(p))
	ch.Abort()
	os.WriteFile(group.NotePath(p), note, 0o600)
	if err := session.Register(launcher.SessionDir(p, "c")); err != nil {
		t.Fatal(err)
	}
	waits := false
	for _, f := range Run(p, next, Options{Fix: true}) {
		if f.Area == "Share links" && f.Fixed && (f.Profile == "a" || f.Profile == "b") {
			t.Errorf("%s repaired before the change of groups: %s", f.Profile, f.Msg)
		}
		waits = waits || f.Profile == "a" && strings.Contains(f.Msg, "change of groups is finished")
	}
	if !waits {
		t.Error("no finding says a waits for the change of groups")
	}
	session.Unregister(launcher.SessionDir(p, "c"))
	Run(p, next, Options{Fix: true})
	if data, _ := os.ReadFile(filepath.Join(p.GroupDir("h"), "CLAUDE.md")); string(data) != "b" {
		t.Errorf("group memory = %q, want b's (--from b)", data)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// Finishing a stopped change keeps the answers given to its conflicts, and
// says where each account's copy went.
func TestFinishingAChangeKeepsTheConflictAnswers(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{ClaudeBin: stubClaude}
	for i, n := range []string{"a", "b"} {
		cfg.Profiles = append(cfg.Profiles, config.Profile{Name: n, Slot: i + 1})
		os.MkdirAll(p.ProfileDir(n), 0o700)
		os.WriteFile(filepath.Join(p.ProfileDir(n), "settings.json"), []byte(`{"who":"`+n+`"}`), 0o600)
	}
	if err := cfg.AddGroup("h", []string{"settings"}, nil); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(p.GroupDir("h"), 0o700)
	os.WriteFile(filepath.Join(p.GroupDir("h"), "settings.json"), []byte(`{"who":"group"}`), 0o600)
	next := cfg.Clone()
	for _, n := range []string{"a", "b"} {
		if _, err := next.JoinGroup(n, "h", true); err != nil {
			t.Fatal(err)
		}
	}
	ch, err := group.Build(p, cfg, next, group.Options{})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range ch.Conflicts() {
		if c.Account == "b" {
			c.Choice = share.UseMine
		}
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	note := mustRead(t, group.NotePath(p))
	ch.Abort()
	os.WriteFile(group.NotePath(p), note, 0o600)
	if err := ch.Accounts[0].Plan.Apply(); err != nil { // a's plan ran, then ccam stopped
		t.Fatal(err)
	}
	said := false
	for _, f := range Run(p, next, Options{Fix: true}) {
		said = said || f.Profile == "b" && strings.Contains(f.Msg, "settings.json") && strings.Contains(f.Msg, "as chosen")
	}
	if data, _ := os.ReadFile(filepath.Join(p.GroupDir("h"), "settings.json")); string(data) != `{"who":"b"}` {
		t.Errorf("group settings = %s, want b's (its answer was use mine)", data)
	}
	if !said {
		t.Error("no finding says what was done with b's settings.json")
	}
}

// A leave that stopped inside a copy back, after the link was removed,
// is finished with the copy put in place.
func TestFinishingALeaveStoppedInACopyBack(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{ClaudeBin: stubClaude, Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	dir := p.ProfileDir("a")
	os.MkdirAll(filepath.Join(dir, "skills", "s"), 0o700)
	os.WriteFile(filepath.Join(dir, "skills", "s", "SKILL.md"), []byte("mine"), 0o600)
	if err := cfg.AddGroup("g", []string{"extensions"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	pl, _ := share.Build(dir, p.GroupsRoot(), config.Share{}, cfg.Profiles[0].Share, share.Options{})
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	next := cfg.Clone()
	next.LeaveGroup("a", "g")
	ch, err := group.Build(p, cfg, next, group.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	note := mustRead(t, group.NotePath(p))
	ch.Abort()
	os.WriteFile(group.NotePath(p), note, 0o600)
	stage := filepath.Join(dir, share.BackupDir, "incomplete-skills")
	if err := fsutil.CopyTree(filepath.Join(p.GroupDir("g"), "skills"), stage); err != nil {
		t.Fatal(err)
	}
	os.Remove(filepath.Join(dir, "skills")) // the link; stopped here
	Run(p, next, Options{Fix: true})
	if data, err := os.ReadFile(filepath.Join(dir, "skills", "s", "SKILL.md")); string(data) != "mine" || fsutil.Exists(stage) {
		t.Fatalf("skills/s/SKILL.md = %q (%v), stage left: %v", data, err, fsutil.Exists(stage))
	}
}

// When the account started between the stop and doctor, Claude Code made
// the entry anew: doctor keeps naming the complete copy, and --fix does
// not claim to have repaired it.
func TestStageKeptAfterALaunch(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{ClaudeBin: stubClaude, Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	dir := p.ProfileDir("a")
	os.MkdirAll(filepath.Join(dir, "skills", "s"), 0o700)
	os.WriteFile(filepath.Join(dir, "skills", "s", "SKILL.md"), []byte("mine"), 0o600)
	if err := cfg.AddGroup("g", []string{"extensions"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	pl, _ := share.Build(dir, p.GroupsRoot(), config.Share{}, cfg.Profiles[0].Share, share.Options{})
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	next := cfg.Clone()
	next.LeaveGroup("a", "g")
	stage := filepath.Join(dir, share.BackupDir, "incomplete-skills")
	if err := fsutil.CopyTree(filepath.Join(p.GroupDir("g"), "skills"), stage); err != nil {
		t.Fatal(err)
	}
	os.Remove(filepath.Join(dir, "skills"))                 // the link; stopped here
	os.MkdirAll(filepath.Join(dir, "skills", "new"), 0o700) // the next launch
	for _, fix := range []bool{false, true} {
		found := false
		for _, f := range Run(p, next, Options{Fix: fix}) {
			if f.Area == "Share links" && strings.Contains(f.Msg, stage) {
				found = true
				if f.Level != Warn || f.Fixed {
					t.Errorf("fix %v: %+v", fix, f)
				}
			}
		}
		if !found {
			t.Errorf("fix %v: the stage is not named", fix)
		}
	}
	if !fsutil.Exists(filepath.Join(stage, "s", "SKILL.md")) || !fsutil.Exists(filepath.Join(dir, "skills", "new")) {
		t.Fatal("something moved")
	}
}

// stoppedJoin makes a, b and c, each with its own CLAUDE.md, join group h
// with --from b, and stops the change right after the config was saved.
func stoppedJoin(t *testing.T) (paths.Paths, *config.Config) {
	t.Helper()
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{ClaudeBin: stubClaude}
	for i, n := range []string{"a", "b", "c"} {
		cfg.Profiles = append(cfg.Profiles, config.Profile{Name: n, Slot: i + 1})
		os.MkdirAll(p.ProfileDir(n), 0o700)
		os.WriteFile(filepath.Join(p.ProfileDir(n), "CLAUDE.md"), []byte(n), 0o600)
	}
	if err := cfg.AddGroup("h", []string{"memory"}, nil); err != nil {
		t.Fatal(err)
	}
	next := cfg.Clone()
	for _, n := range []string{"a", "b", "c"} {
		if _, err := next.JoinGroup(n, "h", true); err != nil {
			t.Fatal(err)
		}
	}
	ch, err := group.Build(p, cfg, next, group.Options{Seed: "b"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	note := mustRead(t, group.NotePath(p))
	ch.Abort()
	os.WriteFile(group.NotePath(p), note, 0o600)
	return p, next
}

// A dry run of doctor --fix previews what --fix does with a stopped change
// of groups: finish it in its order (--from b), so a's and c's copies are
// set aside and b's becomes the group's.
func TestDryRunOfAStoppedChangeFollowsIt(t *testing.T) {
	p, next := stoppedJoin(t)
	var lines []string
	for _, f := range Run(p, next, Options{DryRun: true}) {
		if f.Area == "Share links" || f.Area == "Groups" {
			lines = append(lines, f.Profile+": "+f.Msg)
		}
	}
	all := strings.Join(lines, "\n")
	for _, n := range []string{"a", "c"} {
		if !strings.Contains(all, n+": CLAUDE.md differs from the group's copy: --fix would use the group's, and move this account's copy to "+filepath.Join(p.ProfileDir(n), share.BackupDir)) {
			t.Errorf("%s's copy is not set aside in:\n%s", n, all)
		}
	}
	if strings.Contains(all, "b: CLAUDE.md differs") || !strings.Contains(all, "--fix finishes it in its order (b, a, c)") {
		t.Errorf("dry run:\n%s", all)
	}
	if fsutil.Exists(p.GroupDir("h")) || !group.Pending(p) {
		t.Fatal("the dry run changed something")
	}
	Run(p, next, Options{Fix: true})
	if data, _ := os.ReadFile(filepath.Join(p.GroupDir("h"), "CLAUDE.md")); string(data) != "b" {
		t.Errorf("group memory = %q, want b's (--from b)", data)
	}
}

// doctor --fix while another ccam carries out a change of groups repairs
// none of the share links: a link made between the other ccam's steps
// would make its next step fail.
func TestFixLeavesAChangeInProgressAlone(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{ClaudeBin: stubClaude}
	cfg.Profiles = append(cfg.Profiles, config.Profile{Name: "a", Slot: 1})
	os.MkdirAll(p.ProfileDir("a"), 0o700)
	os.WriteFile(filepath.Join(p.ProfileDir("a"), "CLAUDE.md"), []byte("a"), 0o600)
	next := cfg.Clone()
	if err := next.AddGroup("h", []string{"memory"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	next.Resolve()
	ch, err := group.Build(p, cfg, next, group.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Begin(); err != nil {
		t.Fatal(err)
	}
	defer ch.Abort()
	pl := ch.Accounts[0].Plan
	k := slicesIndex(pl.Steps, share.OpLink)
	first, rest := *pl, *pl
	first.Steps, rest.Steps = pl.Steps[:k], pl.Steps[k:]
	if err := first.Apply(); err != nil { // the other ccam is between its steps
		t.Fatal(err)
	}
	said := false
	for _, f := range Run(p, next, Options{Fix: true}) {
		if (f.Area == "Share links" || f.Area == "Groups") && f.Fixed {
			t.Errorf("repaired during the change: %s", f.Msg)
		}
		said = said || strings.Contains(f.Msg, "being carried out right now")
	}
	if !said {
		t.Error("no finding says a change is being carried out")
	}
	if err := rest.Apply(); err != nil {
		t.Fatalf("the other ccam's next step: %v", err)
	}
}

func slicesIndex(steps []share.Step, op share.Op) int {
	for i, s := range steps {
		if s.Op == op {
			return i
		}
	}
	return len(steps)
}

// doctor --fix repairs toward the config saved when it took the lock, not
// one it loaded before: a leave finished in between stays a leave.
func TestFixRepairsTowardTheSavedConfig(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{ClaudeBin: stubClaude}
	cfg.Profiles = append(cfg.Profiles, config.Profile{Name: "a", Slot: 1})
	os.MkdirAll(p.ProfileDir("a"), 0o700)
	os.WriteFile(filepath.Join(p.ProfileDir("a"), "CLAUDE.md"), []byte("a"), 0o600)
	joined := cfg.Clone()
	if err := joined.AddGroup("h", []string{"memory"}, []string{"a"}); err != nil {
		t.Fatal(err)
	}
	joined.Resolve()
	left := joined.Clone()
	if err := left.LeaveGroup("a", "h"); err != nil {
		t.Fatal(err)
	}
	for _, step := range [][2]*config.Config{{cfg, joined}, {joined, left}} {
		ch, err := group.Build(p, step[0], step[1], group.Options{})
		if err == nil {
			err = ch.Begin()
		}
		if err == nil {
			err = ch.Apply()
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	// doctor loaded joined; the leave was saved since.
	Run(p, joined, Options{Fix: true, Reload: func() (*config.Config, error) { return left, nil }})
	if _, linked := fsutil.LinkTarget(filepath.Join(p.ProfileDir("a"), "CLAUDE.md")); linked {
		t.Fatal("CLAUDE.md was linked into the group it left")
	}
}

// A shortcut ccam made for another ccam that is still installed works: it
// is reported as running that one, not as missing, and --fix points it at
// the ccam that runs doctor, saying so only when it did.
func TestShortcutOfAnotherCcam(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("needs a symlink")
	}
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	self, other := filepath.Join(p.Home, "new", "ccam"), filepath.Join(p.Home, "old", "ccam")
	for _, bin := range []string{self, other} {
		os.MkdirAll(filepath.Dir(bin), 0o755)
		os.WriteFile(bin, []byte("ccam"), 0o755)
	}
	cfg := &config.Config{ClaudeBin: stubClaude, Profiles: []config.Profile{{Name: "work", Slot: 1, Shortcut: "ccwork"}}}
	os.MkdirAll(p.ProfileDir("work"), 0o700)
	if _, err := shim.Install(p, other, cfg); err != nil { // made by the old ccam
		t.Fatal(err)
	}
	shortcut := func(fix bool) Finding {
		var got []Finding
		for _, f := range Run(p, cfg, Options{Fix: fix, CcamPath: self}) {
			if f.Area == "Shortcuts" && !strings.Contains(f.Msg, "not in PATH") {
				got = append(got, f)
			}
		}
		if len(got) != 1 {
			t.Fatalf("shortcut findings: %+v", got)
		}
		return got[0]
	}
	if f := shortcut(false); strings.Contains(f.Msg, "missing") || !strings.Contains(f.Msg, "run another ccam") || !strings.Contains(f.Msg, other) || f.Level == Warn || f.Level == Error {
		t.Fatalf("doctor: %+v", f)
	}
	if f := shortcut(true); !f.Fixed {
		t.Fatalf("doctor --fix: %+v", f)
	}
	if !shim.Linked(p.Bin, self, "ccwork") {
		t.Fatal("--fix said fixed, but ccwork still runs the other ccam")
	}
	if f := shortcut(false); f.Level != OK {
		t.Fatalf("after --fix: %+v", f)
	}
}

// With auto_switch on, doctor says what it does, quotes Anthropic's seat
// rule for each plan that has seats among the accounts (read as the
// automatic choice reads it: .claude.json, else the recorded login), and
// warns about what keeps a session where it is: an account without the ccam
// plugin turned on, whose hook would not run, and a Claude Code that gives
// hooks no CLAUDE_PID. Windows ignores auto_switch and says only that; with
// auto_switch off, none of these lines are there.
func TestAutoSwitchChecks(t *testing.T) {
	t.Setenv("CCAM_HOME", t.TempDir())
	p := paths.FromEnv()
	cfg := &config.Config{AutoSwitch: true, Profiles: []config.Profile{
		{Name: "personal", Slot: 1}, {Name: "work", Slot: 2}, {Name: "lab", Slot: 3},
		{Name: "corp", Slot: 4, Login: &config.Identity{AccountUUID: "u4", OrgUUID: "big", OrgType: "claude_enterprise"}},
		{Name: "api", Slot: 5, API: &config.API{KeyFrom: "env:API_KEY"}},
	}}
	for name, org := range map[string]string{"personal": `"organizationUuid":"p","organizationType":"claude_max"`,
		"work": `"organizationUuid":"acme","organizationName":"Acme","organizationType":"claude_team"`,
		"lab":  `"organizationUuid":"acme","organizationName":"Acme","organizationType":"claude_team"`} {
		os.MkdirAll(p.ProfileDir(name), 0o700)
		os.WriteFile(filepath.Join(p.ProfileDir(name), ".claude.json"), []byte(`{"oauthAccount":{"accountUuid":"u-`+name+`",`+org+`}}`), 0o600)
	}
	// personal and corp have the plugin on, work has it turned off, lab has
	// none, and an API account never moves.
	for _, name := range []string{"personal", "work", "corp"} {
		installPlugin(t, p.ProfileDir(name), "ccam@ccam", "{}")
	}
	os.WriteFile(filepath.Join(p.ProfileDir("work"), "settings.json"), []byte(`{"enabledPlugins":{"ccam@ccam":false}}`), 0o600)

	check := func(claude, goos string) []string {
		r := &report{}
		checkAutoSwitch(r, p, cfg, claude, goos)
		var out []string
		for _, f := range r.out {
			if f.Area != "auto_switch" {
				t.Errorf("area %q: %+v", f.Area, f)
			}
			out = append(out, fmt.Sprintf("%d %s: %s", f.Level, f.Profile, f.Msg))
		}
		return out
	}
	on := fmt.Sprintf("%d : auto_switch is on: at a usage limit, ccam moves an interactive session it started to another account like the one it leaves (a seat of the same organization, or another personal plan) with room. Whether that fits Anthropic's terms is your call: https://ccam.dev/en/features/next/#auto-switch", Info)
	plugins := []string{
		fmt.Sprintf("%d work: auto_switch is on, but work has the ccam plugin turned off, so its sessions are not moved. In Claude Code as work: /plugin enable ccam@ccam", Warn),
		fmt.Sprintf("%d lab: auto_switch is on, but lab does not have the ccam plugin, so its sessions are not moved. In Claude Code as lab: /plugin install ccam --marketplace jhyunwoo/ccam", Warn),
	}
	seats := []string{
		fmt.Sprintf(`%d : Claude Team seats (work, lab): auto_switch moves a session between seats of one organization. Anthropic's help center says "Each seat on a Team plan is assigned to one named individual, and only that person may use it." and "Users can't route their work through another user's account to use that user's remaining usage." Whether you use it there is your call: https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129`, Info),
		fmt.Sprintf(`%d : Claude Enterprise seats (corp): auto_switch moves a session between seats of one organization. Anthropic's help center says "Each Enterprise seat is assigned to one named individual, and only that person may use it." and "Users can't use another user's account to get around their own spend limits or usage limits." Whether you use it there is your call: https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0`, Info),
	}
	old := fmt.Sprintf("%d : auto_switch needs Claude Code 2.1.214 or later, which tells hooks its process (CLAUDE_PID); this one is 2.1.200, so no session is moved. Update it with `claude update`", Warn)
	windows := []string{fmt.Sprintf("%d : auto_switch is ignored on Windows; with relay = true, /ccam:switch and then /exit carry a conversation on.", Info)}
	for _, c := range []struct {
		claude, goos string
		want         []string
	}{
		{"2.1.310 (Claude Code)", "linux", slices.Concat([]string{on}, plugins, seats)},
		{"2.1.214 (Claude Code)", "darwin", slices.Concat([]string{on}, plugins, seats)},
		{"", "linux", slices.Concat([]string{on}, plugins, seats)}, // claude --version failed
		{"2.1.200 (Claude Code)", "linux", slices.Concat([]string{on, old}, plugins, seats)},
		{"2.1.200 (Claude Code)", "windows", windows},
	} {
		if got := check(c.claude, c.goos); !slices.Equal(got, c.want) {
			t.Errorf("%s on %s:\n got %q\nwant %q", c.claude, c.goos, got, c.want)
		}
	}
	cfg.AutoSwitch = false
	if got := check("2.1.200 (Claude Code)", "linux"); len(got) != 0 {
		t.Errorf("auto_switch off: %q", got)
	}

	// Run passes on what `claude --version` printed, and the platform.
	cfg.AutoSwitch, cfg.ClaudeBin = true, stubClaude
	want := slices.Concat([]string{on, old}, plugins, seats)
	if runtime.GOOS == "windows" {
		want = windows // the stub claude does not start there
	} else {
		cfg.ClaudeBin = filepath.Join(t.TempDir(), "claude")
		os.WriteFile(cfg.ClaudeBin, []byte("#!/bin/sh\necho '2.1.200 (Claude Code)'\n"), 0o755)
	}
	var got []string
	for _, f := range Run(p, cfg, Options{}) {
		if f.Area == "auto_switch" {
			got = append(got, fmt.Sprintf("%d %s: %s", f.Level, f.Profile, f.Msg))
		}
	}
	if !slices.Equal(got, want) {
		t.Errorf("Run:\n got %q\nwant %q", got, want)
	}
}
