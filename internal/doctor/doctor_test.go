package doctor

import (
	"encoding/json"
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
