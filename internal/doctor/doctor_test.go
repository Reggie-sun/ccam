package doctor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/paths"
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
	if len(r.out) != 1 || r.out[0].Level != Warn || strings.Contains(r.out[0].String(), "secret") {
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
	all := fmt.Sprint(r.out)
	if !strings.Contains(all, "b: logged in to the same Claude account as a") || !strings.Contains(all, "c: now logged in as other@x, but ccam recorded c@x") || len(r.out) != 2 {
		t.Fatalf("%v", r.out)
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
