package profile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/paths"
)

func testPaths(t *testing.T) paths.Paths {
	t.Setenv("CCAM_HOME", t.TempDir())
	return paths.FromEnv()
}

func TestAddCreatesPrivateDirAndSlots(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	a, err := Add(p, cfg, "personal", config.DefaultShare())
	if err != nil {
		t.Fatal(err)
	}
	if a.Slot != 1 || cfg.Default != "personal" {
		t.Fatalf("first profile: %+v default=%q", a, cfg.Default)
	}
	b, err := Add(p, cfg, "work", config.Share{})
	if err != nil || b.Slot != 2 {
		t.Fatalf("second: %+v %v", b, err)
	}
	info, err := os.Stat(p.ProfileDir("work"))
	if err != nil || info.Mode().Perm() != 0o700 {
		t.Fatalf("dir perm %v %v", info, err)
	}
	if _, err := Add(p, cfg, "work", config.Share{}); err == nil {
		t.Fatal("duplicate accepted")
	}
	if _, err := Add(p, cfg, "run", config.Share{}); err == nil {
		t.Fatal("reserved name accepted")
	}
}

func TestTrashMovesManagedOnly(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	prof, _ := Add(p, cfg, "work", config.Share{})
	os.WriteFile(filepath.Join(p.ProfileDir("work"), ".credentials.json"), []byte("{}"), 0o600)
	dst, err := Trash(p, prof)
	if err != nil || !strings.HasPrefix(dst, p.TrashDir()) {
		t.Fatalf("%q %v", dst, err)
	}
	if _, err := os.Stat(filepath.Join(dst, ".credentials.json")); err != nil {
		t.Fatal("trash lost the file")
	}

	ext := t.TempDir()
	imported := &config.Profile{Name: "old", Dir: ext}
	if _, err := Trash(p, imported); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(ext); err != nil {
		t.Fatal("imported dir must stay")
	}
}

func TestRename(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	Add(p, cfg, "work", config.Share{})
	if err := Rename(p, cfg, "work", "job"); err != nil {
		t.Fatal(err)
	}
	if _, ok := cfg.Find("job"); !ok || cfg.Default != "job" {
		t.Fatal("config not renamed")
	}
	if _, err := os.Stat(p.ProfileDir("job")); err != nil {
		t.Fatal("folder not renamed")
	}
}

func TestReadEmailStreams(t *testing.T) {
	dir := t.TempDir()
	big := `{"projects":{"/a":{"x":[1,2,3]}},"numStartups":3,"oauthAccount":{"emailAddress":"dev@company.com","accountUuid":"u"},"mcpServers":{}}`
	os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(big), 0o600)
	email, err := ReadEmail(dir)
	if err != nil || email != "dev@company.com" {
		t.Fatalf("%q %v", email, err)
	}

	os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{"numStartups":1}`), 0o600)
	if email, err := ReadEmail(dir); err != nil || email != "" {
		t.Fatalf("logged-out file: %q %v", email, err)
	}
}
