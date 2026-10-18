package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
)

func TestSuggestName(t *testing.T) {
	cases := map[string]string{
		"/h/.claude":      "main",
		"/h/.claude-work": "work",
		"/h/.claude-2":    "acct2",
		"/h/.claude-Team": "team",
		"/h/.claude-run":  "", // reserved
		"/h/.claude-":     "",
	}
	for in, want := range cases {
		if got := SuggestName(in); got != want {
			t.Errorf("%s: %q want %q", in, got, want)
		}
	}
}

func TestImportDefaultFolderLinksHomeClaudeJSON(t *testing.T) {
	p := testPaths(t)
	claudeDir := filepath.Join(p.Home, ".claude")
	os.MkdirAll(claudeDir, 0o700)
	os.WriteFile(filepath.Join(claudeDir, ".credentials.json"), []byte(`{}`), 0o600)
	os.WriteFile(filepath.Join(p.Home, ".claude.json"), []byte(`{"oauthAccount":{"emailAddress":"me@example.com"}}`), 0o600)
	os.MkdirAll(filepath.Join(p.Home, ".claude-work"), 0o700)
	os.WriteFile(filepath.Join(p.Home, ".claude-work", ".claude.json"), []byte(`{"oauthAccount":{"emailAddress":"dev@co.com"}}`), 0o600)

	cfg := &config.Config{}
	cands := FindCandidates(p, cfg)
	if len(cands) != 2 || cands[0].Email != "me@example.com" || cands[1].Name != "work" {
		t.Fatalf("candidates %+v", cands)
	}
	prof, notes, err := Import(p, cfg, claudeDir, "")
	if err != nil || prof.Name != "main" || prof.Dir != claudeDir || len(notes) != 1 || prof.Email != "me@example.com" {
		t.Fatalf("%+v %v %v", prof, notes, err)
	}
	if target, _ := os.Readlink(filepath.Join(claudeDir, ".claude.json")); target != filepath.Join(p.Home, ".claude.json") {
		t.Fatalf("link %q", target)
	}
	if !prof.Share.Equal(config.Share{}) || cfg.Default != "main" {
		t.Fatal("imported profiles start unshared and the first becomes default")
	}
	prof, notes, err = Import(p, cfg, filepath.Join(p.Home, ".claude-work"), "")
	if err != nil || prof.Name != "work" || prof.Slot != 2 || len(notes) != 0 || prof.Email != "dev@co.com" {
		t.Fatalf("%+v %v %v", prof, notes, err)
	}
	if len(FindCandidates(p, cfg)) != 0 {
		t.Fatal("registered folders must not be offered again")
	}
	if _, _, err := Import(p, cfg, claudeDir, "other"); err == nil {
		t.Fatal("double import accepted")
	}
}

func TestImportNeedsNameWhenNoneDerivable(t *testing.T) {
	p := testPaths(t)
	dir := filepath.Join(p.Home, ".claude-run")
	os.MkdirAll(dir, 0o700)
	if _, _, err := Import(p, &config.Config{}, dir, ""); err == nil {
		t.Fatal("want error")
	}
	if prof, _, err := Import(p, &config.Config{}, dir, "runner"); err != nil || prof.Name != "runner" {
		t.Fatalf("%v %v", prof, err)
	}
}
