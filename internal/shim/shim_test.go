package shim

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
)

func TestLookup(t *testing.T) {
	cfg := &config.Config{Profiles: []config.Profile{
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
	// Only configured shortcuts resolve: not ccam, account names or numbers.
	for _, in := range []string{"ccam", "claude", "cchome", "work", "5"} {
		if _, ok := Lookup(cfg, in); ok {
			t.Errorf("%s should not be a shortcut", in)
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
	os.WriteFile(filepath.Join(f.bin, "userfile"), []byte("user file"), 0o755) // not ours
	os.Symlink(f.ccam, filepath.Join(f.bin, "mytool"))                         // ours, but not configured: left alone

	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}, {Name: "b", Slot: 2}}}
	created, err := Install(f.bin, f.ccam, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 || created[0] != "cca" {
		t.Fatalf("created %v", created)
	}
	if entries, _ := os.ReadDir(f.bin); len(entries) != 3 {
		t.Errorf("only cca is created: %v", entries)
	}
	if b, _ := os.ReadFile(filepath.Join(f.bin, "userfile")); string(b) != "user file" {
		t.Fatal("user file touched")
	}
	if !Linked(f.bin, f.ccam, "mytool") {
		t.Fatal("Install never removes links")
	}
	// Idempotent.
	if created, err := Install(f.bin, f.ccam, cfg); err != nil || len(created) != 0 {
		t.Fatalf("second install %v %v", created, err)
	}
	Remove(f.bin, f.ccam, "cca")
	Remove(f.bin, f.ccam, "userfile") // not ours
	if Linked(f.bin, f.ccam, "cca") || !fileExists(filepath.Join(f.bin, "userfile")) {
		t.Fatal("Remove must delete only our link")
	}
}

func TestInstallWithoutShortcutsCreatesNothing(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "bin")
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}}}
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
	os.WriteFile(filepath.Join(f.bin, "cca"), []byte("mine"), 0o755)
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}}}
	if _, err := Install(f.bin, f.ccam, cfg); err == nil {
		t.Fatal("want error for a foreign cca")
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
