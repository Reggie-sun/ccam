package doctor

import (
	"os"
	"path/filepath"
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
	os.WriteFile(file, []byte(`{"plugins":{"demo@mk":[{"installPath":"`+viaProfile+`"}]}}`), 0o600)

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
	if !strings.Contains(string(data), filepath.Join(pluginsDir, "cache", "mk", "demo", "1.0")) {
		t.Fatalf("rewritten: %s", data)
	}
}

func TestCredentialsChecks(t *testing.T) {
	dir := t.TempDir()
	prof := &config.Profile{Name: "a"}
	cred := filepath.Join(dir, ".credentials.json")
	os.WriteFile(cred, []byte(`{"claudeAiOauth":{"expiresAt":1,"refreshTokenExpiresAt":2}}`), 0o644)
	r := &report{}
	checkCredentials(r, prof, dir, true)
	if len(r.out) != 2 || !r.out[0].Fixed || r.out[1].Level != Error {
		t.Fatalf("%v", r.out)
	}
	if info, _ := os.Stat(cred); info.Mode().Perm() != 0o600 {
		t.Fatal("perm not fixed")
	}
}
