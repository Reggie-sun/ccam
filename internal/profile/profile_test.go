package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/paths"
)

func testPaths(t *testing.T) paths.Paths {
	t.Setenv("CCAM_HOME", t.TempDir())
	return paths.FromEnv()
}

func TestAddCreatesPrivateDirAndSlots(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	a, err := Add(p, cfg, "personal")
	if err != nil {
		t.Fatal(err)
	}
	if a.Slot != 1 || cfg.Default != "personal" {
		t.Fatalf("first profile: %+v default=%q", a, cfg.Default)
	}
	b, err := Add(p, cfg, "work")
	if err != nil || b.Slot != 2 {
		t.Fatalf("second: %+v %v", b, err)
	}
	info, err := os.Stat(p.ProfileDir("work"))
	if err != nil || runtime.GOOS != "windows" && info.Mode().Perm() != 0o700 {
		t.Fatalf("dir perm %v %v", info, err)
	}
	if _, err := Add(p, cfg, "work"); err == nil {
		t.Fatal("duplicate accepted")
	}
	if _, err := Add(p, cfg, "run"); err == nil {
		t.Fatal("reserved name accepted")
	}
}

func TestTrashMovesManagedOnly(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	prof, _ := Add(p, cfg, "work")
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
	Add(p, cfg, "work")
	cfg.SetFolder(t.TempDir(), "work")
	if err := cfg.AddGroup("default", config.DefaultItems, []string{"work"}); err != nil {
		t.Fatal(err)
	}
	if err := Rename(p, cfg, "work", "job", false); err != nil {
		t.Fatal(err)
	}
	if _, ok := cfg.Find("job"); !ok || cfg.Default != "job" || cfg.Folders[0].Account != "job" || cfg.Groups[0].Members[0] != "job" {
		t.Fatal("config not renamed")
	}
	if _, err := os.Stat(p.ProfileDir("job")); err != nil {
		t.Fatal("folder not renamed")
	}
}

// The macOS Keychain names a login after its folder's path, so a rename
// there keeps the folder (Profile.Folder): the login, and the MCP sign-ins
// stored with it, stay the account's. No new account takes the folder over.
func TestRenameKeepsTheFolder(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	Add(p, cfg, "work")
	dir := p.ProfileDir("work")
	if err := Rename(p, cfg, "work", "job", true); err != nil {
		t.Fatal(err)
	}
	job, _ := cfg.Find("job")
	if Dir(p, job) != dir || job.Folder != "work" || fsutil.Exists(p.ProfileDir("job")) {
		t.Fatalf("%+v", job)
	}
	if _, err := Add(p, cfg, "work"); err == nil || !strings.Contains(err.Error(), "account job's") {
		t.Fatalf("a new work would take job's folder: %v", err)
	}
	cfgFile := filepath.Join(t.TempDir(), "config.toml")
	if err := cfg.Save(cfgFile); err != nil {
		t.Fatal(err)
	}
	if loaded, err := config.Load(cfgFile); err != nil || Dir(p, &loaded.Profiles[0]) != dir {
		t.Fatalf("saved and loaded: %v %+v", err, loaded)
	}
	// Renamed on, and back to the folder's name.
	if err := Rename(p, cfg, "job", "side", true); err != nil || cfg.Profiles[0].Folder != "work" {
		t.Fatalf("%v %+v", err, cfg.Profiles[0])
	}
	if err := Rename(p, cfg, "side", "work", true); err != nil || cfg.Profiles[0].Folder != "" || Dir(p, &cfg.Profiles[0]) != dir {
		t.Fatalf("%v %+v", err, cfg.Profiles[0])
	}
	// A rename that may move the folder moves the kept one, and Trash takes
	// the kept one too.
	Rename(p, cfg, "work", "a", true)
	if err := Rename(p, cfg, "a", "b", false); err != nil || cfg.Profiles[0].Folder != "" || !fsutil.Exists(p.ProfileDir("b")) || fsutil.Exists(dir) {
		t.Fatalf("%v %+v", err, cfg.Profiles[0])
	}
	Rename(p, cfg, "b", "c", true)
	if dst, err := Trash(p, &cfg.Profiles[0]); err != nil || dst == "" || fsutil.Exists(p.ProfileDir("b")) {
		t.Fatalf("trash %q %v", dst, err)
	}
}

// A rename that moves the folder must not move it onto the one another
// account kept.
func TestRenameDoesNotTakeAKeptFolder(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	Add(p, cfg, "work")
	Rename(p, cfg, "work", "job", true) // job keeps profiles/work
	Add(p, cfg, "side")
	if err := Rename(p, cfg, "side", "work", false); err == nil || !strings.Contains(err.Error(), "account job's") {
		t.Fatalf("%v", err)
	}
	if !fsutil.Exists(p.ProfileDir("side")) {
		t.Fatal("side's folder moved")
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

func TestRenameRewritesPluginPaths(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	Add(p, cfg, "work")
	old := p.ProfileDir("work")
	os.MkdirAll(filepath.Join(old, "plugins"), 0o700)
	install, _ := json.Marshal(filepath.Join(old, "plugins", "cache", "m", "x", "1"))
	os.WriteFile(filepath.Join(old, "plugins", "installed_plugins.json"),
		[]byte(`{"plugins":{"x@m":[{"installPath":`+string(install)+`}]}}`), 0o600)
	if err := Rename(p, cfg, "work", "job", false); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(p.ProfileDir("job"), "plugins", "installed_plugins.json"))
	want, _ := json.Marshal(filepath.Join(p.ProfileDir("job"), "plugins", "cache", "m", "x", "1"))
	if !strings.Contains(string(data), string(want)) {
		t.Fatalf("%s", data)
	}
}

func TestRenameRollsBackWhenPluginRewriteFails(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	Add(p, cfg, "work")
	os.MkdirAll(filepath.Join(p.ProfileDir("work"), "plugins"), 0o700)
	os.WriteFile(filepath.Join(p.ProfileDir("work"), "plugins", "installed_plugins.json"), []byte("{broken"), 0o600)
	if err := Rename(p, cfg, "work", "job", false); err == nil {
		t.Fatal("want error")
	}
	if _, err := os.Stat(p.ProfileDir("work")); err != nil {
		t.Fatal("folder must keep the name the config still has")
	}
	if _, ok := cfg.Find("work"); !ok {
		t.Fatal("config renamed despite the error")
	}
}

func TestTrashTwiceKeepsBothCopies(t *testing.T) {
	p := testPaths(t)
	var dsts []string
	for _, body := range []string{"first", "second"} {
		cfg := &config.Config{}
		prof, _ := Add(p, cfg, "tt")
		os.WriteFile(filepath.Join(p.ProfileDir("tt"), "note"), []byte(body), 0o600)
		dst, err := Trash(p, prof)
		if err != nil {
			t.Fatal(err)
		}
		dsts = append(dsts, dst)
	}
	if dsts[0] == dsts[1] {
		t.Fatalf("same trash folder %s", dsts[0])
	}
	for i, want := range []string{"first", "second"} {
		if b, _ := os.ReadFile(filepath.Join(dsts[i], "note")); string(b) != want {
			t.Fatalf("%s: %q", dsts[i], b)
		}
	}
}

// A group may still list a name from an account that is gone; renaming an
// account to that name must not hand it that group.
func TestRenameOverAStaleMember(t *testing.T) {
	p := testPaths(t)
	cfg := &config.Config{}
	Add(p, cfg, "work")
	cfg.Groups = []config.Group{{Name: "team", Items: []string{"history"}, Members: []string{"job"}}}
	if err := Rename(p, cfg, "work", "job", false); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Save(filepath.Join(t.TempDir(), "c.toml")); err != nil {
		t.Fatal(err)
	}
	if len(cfg.Groups[0].Members) != 0 {
		t.Fatalf("%v", cfg.Groups[0].Members)
	}
}
