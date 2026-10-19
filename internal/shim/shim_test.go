package shim

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
)

func TestLookup(t *testing.T) {
	cfg := &config.Config{Version: 2, Profiles: []config.Profile{
		{Name: "work", Slot: 1, Shortcut: "ccwork"},
		{Name: "side", Slot: 2, Shortcut: "ccside"},
		{Name: "home", Slot: 5},
	}}
	cases := map[string]string{"/home/u/.local/bin/ccwork": "work", "ccside": "side"}
	for in, want := range cases {
		if got, ok := Lookup(cfg, in); !ok || got != want {
			t.Errorf("%s -> %q %v", in, got, ok)
		}
	}
	// No guessing from claude<N> once shortcuts are stored.
	for _, in := range []string{"ccam", "ccs", "claude", "claude2", "claude-home", "work"} {
		if _, ok := Lookup(cfg, in); ok {
			t.Errorf("%s should not be a shortcut", in)
		}
	}
	// A version 1 config (not yet upgraded by the CLI) keeps the old names working.
	old := &config.Config{Version: 1, Profiles: []config.Profile{{Name: "work", Slot: 2}}}
	for in, want := range map[string]string{"claude2": "2", "/x/claude12": "12", "claude-work": "work"} {
		if got, ok := Lookup(old, in); !ok || got != want {
			t.Errorf("v1 %s -> %q %v", in, got, ok)
		}
	}
	for _, in := range []string{"ccam", "ccs", "claude", "claude-", "claudex"} {
		if _, ok := Lookup(old, in); ok {
			t.Errorf("v1 %s should not be a shortcut", in)
		}
	}
}

type binFixture struct{ dir, bin, ccam string }

func newBin(t *testing.T) binFixture {
	dir := t.TempDir()
	f := binFixture{dir: dir, bin: filepath.Join(dir, "bin"), ccam: filepath.Join(dir, "ccam")}
	os.WriteFile(f.ccam, []byte("x"), 0o755)
	os.MkdirAll(f.bin, 0o755)
	return f
}

func TestInstallCreatesOnlyConfiguredShortcuts(t *testing.T) {
	f := newBin(t)
	os.WriteFile(filepath.Join(f.bin, "claude9"), []byte("user file"), 0o755) // not ours
	os.Symlink(f.ccam, filepath.Join(f.bin, "claude7"))                       // stale, ours
	os.Symlink(f.ccam, filepath.Join(f.bin, "claude-old"))                    // stale, ours
	os.Symlink(f.ccam, filepath.Join(f.bin, "ccs"))                           // removed alias, ours
	os.Symlink(f.ccam, filepath.Join(f.bin, "mytool"))                        // ours, but a name we never made: left alone

	cfg := &config.Config{Version: 2, Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}, {Name: "b", Slot: 2}}}
	created, err := Install(f.bin, f.ccam, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 || created[0] != "cca" {
		t.Fatalf("created %v", created)
	}
	for _, gone := range []string{"claude7", "claude-old", "ccs", "claude1", "claude2"} {
		if _, err := os.Lstat(filepath.Join(f.bin, gone)); err == nil {
			t.Errorf("%s should not exist", gone)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(f.bin, "claude9")); string(b) != "user file" {
		t.Fatal("user file touched")
	}
	if !Linked(f.bin, f.ccam, "mytool") {
		t.Fatal("a link we did not name must be left alone")
	}
	// Idempotent.
	if created, err := Install(f.bin, f.ccam, cfg); err != nil || len(created) != 0 {
		t.Fatalf("second install %v %v", created, err)
	}
	Remove(f.bin, f.ccam, "cca")
	Remove(f.bin, f.ccam, "claude9") // not ours
	if Linked(f.bin, f.ccam, "cca") || !fileExists(filepath.Join(f.bin, "claude9")) {
		t.Fatal("Remove must delete only our link")
	}
}

func TestInstallWithoutShortcutsCreatesNothing(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "bin")
	cfg := &config.Config{Version: 2, Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	if created, err := Install(bin, filepath.Join(dir, "ccam"), cfg); err != nil || len(created) != 0 {
		t.Fatalf("%v %v", created, err)
	}
	if _, err := os.Stat(bin); err == nil {
		t.Fatal("no shortcuts, no bin folder")
	}
}

func fileExists(p string) bool { _, err := os.Lstat(p); return err == nil }

func TestInstallRefusesForeignFile(t *testing.T) {
	f := newBin(t)
	os.WriteFile(filepath.Join(f.bin, "claude1"), []byte("mine"), 0o755)
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "claude1"}}}
	if _, err := Install(f.bin, f.ccam, cfg); err == nil {
		t.Fatal("want error for foreign claude1")
	}
}

func TestInstallKeepsForeignCcs(t *testing.T) {
	f := newBin(t)
	other := filepath.Join(f.dir, "other")
	os.WriteFile(other, []byte("y"), 0o755)
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "claude1"}}}

	os.WriteFile(filepath.Join(f.bin, "ccs"), []byte("user file"), 0o755)
	if _, err := Install(f.bin, f.ccam, cfg); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(f.bin, "ccs")); string(b) != "user file" {
		t.Fatal("regular ccs file touched")
	}

	os.Remove(filepath.Join(f.bin, "ccs"))
	os.Symlink(other, filepath.Join(f.bin, "ccs"))
	if _, err := Install(f.bin, f.ccam, cfg); err != nil {
		t.Fatal(err)
	}
	if target, err := os.Readlink(filepath.Join(f.bin, "ccs")); err != nil || target != other {
		t.Fatalf("foreign ccs link changed: %q %v", target, err)
	}
}

func TestAdoptKeepsOldLinksOnDisk(t *testing.T) {
	f := newBin(t)
	os.Symlink(f.ccam, filepath.Join(f.bin, "claude1"))
	os.Symlink(f.ccam, filepath.Join(f.bin, "claude-side")) // --names user without claude3
	os.WriteFile(filepath.Join(f.bin, "claude2"), []byte("not ours"), 0o755)
	cfg := &config.Config{Version: 1, Profiles: []config.Profile{
		{Name: "work", Slot: 1}, {Name: "home", Slot: 2}, {Name: "side", Slot: 3},
	}}
	Adopt(f.bin, f.ccam, cfg)
	got := []string{cfg.Profiles[0].Shortcut, cfg.Profiles[1].Shortcut, cfg.Profiles[2].Shortcut}
	if strings.Join(got, ",") != "claude1,,claude-side" || cfg.Version != 2 {
		t.Fatalf("adopted %q version %d", got, cfg.Version)
	}
	// Already upgraded configs are left as they are.
	cfg.Profiles[0].Shortcut = ""
	Adopt(f.bin, f.ccam, cfg)
	if cfg.Profiles[0].Shortcut != "" {
		t.Fatal("a version 2 config must not be adopted again")
	}
}

func TestCheck(t *testing.T) {
	f := newBin(t)
	elsewhere := filepath.Join(f.dir, "usr-bin")
	os.MkdirAll(elsewhere, 0o755)
	os.WriteFile(filepath.Join(elsewhere, "taken"), []byte("#!/bin/sh\n"), 0o755)
	os.WriteFile(filepath.Join(elsewhere, "data"), []byte("not executable"), 0o644)
	os.WriteFile(filepath.Join(f.bin, "foreign"), []byte("x"), 0o755)
	os.Symlink(f.ccam, filepath.Join(f.bin, "ccwork"))
	t.Setenv("PATH", f.bin+string(os.PathListSeparator)+elsewhere)
	cfg := &config.Config{Profiles: []config.Profile{{Name: "work", Slot: 1, Shortcut: "ccwork"}, {Name: "home", Slot: 2}}}

	for name, wantErr := range map[string]string{
		"Bad":      "1-64 characters",
		"-x":       "1-64 characters",
		"a/b":      "1-64 characters",
		"ccam":     "taken by ccam",
		"claude":   "taken by ccam",
		"list":     "taken by ccam",
		"shortcut": "taken by ccam",
		"ccwork":   "already the shortcut of account work",
		"taken":    "would hide " + filepath.Join(elsewhere, "taken"),
		"foreign":  "is not a ccam shortcut",
	} {
		err := Check(name, f.bin, f.ccam, cfg, "home")
		if err == nil || !strings.Contains(err.Error(), wantErr) {
			t.Errorf("%s: %v, want %q", name, err, wantErr)
		}
	}
	for _, ok := range []string{"cchome", "home", "c.h", "data", "h_2"} {
		if err := Check(ok, f.bin, f.ccam, cfg, "home"); err != nil {
			t.Errorf("%s: %v", ok, err)
		}
	}
	// The account's own shortcut is fine (renaming to the same name).
	if err := Check("ccwork", f.bin, f.ccam, cfg, "work"); err != nil {
		t.Errorf("own shortcut: %v", err)
	}
}

func TestSuggest(t *testing.T) {
	if Suggest("work") != "ccwork" || Suggest("my-side_2") != "ccmy-side_2" {
		t.Fatal(Suggest("work"))
	}
	for _, name := range []string{"work", "a", "my-side_2", "x1234567890123456789012345678901"} {
		if !shortcutRe.MatchString(Suggest(name)) {
			t.Errorf("suggestion for %s is not a valid shortcut name", name)
		}
	}
}
