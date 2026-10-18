package shim

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
)

func TestSelector(t *testing.T) {
	cases := map[string]string{"claude2": "2", "/home/u/.local/bin/claude12": "12", "claude-work": "work"}
	for in, want := range cases {
		if got, ok := Selector(in); !ok || got != want {
			t.Errorf("%s -> %q %v", in, got, ok)
		}
	}
	for _, in := range []string{"ccam", "ccs", "claude", "claude-", "claudex"} {
		if _, ok := Selector(in); ok {
			t.Errorf("%s should not be a shim", in)
		}
	}
}

func TestInstallSyncsLinks(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "bin")
	ccam := filepath.Join(dir, "ccam")
	os.WriteFile(ccam, []byte("x"), 0o755)
	os.MkdirAll(bin, 0o755)
	os.WriteFile(filepath.Join(bin, "claude9"), []byte("user file"), 0o755) // not ours
	os.Symlink(ccam, filepath.Join(bin, "claude7"))                      // stale, ours

	cfg := &config.Config{NameShims: true, Profiles: []config.Profile{{Name: "a", Slot: 1}, {Name: "b", Slot: 2}}}
	created, err := Install(bin, ccam, cfg)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(created)
	want := []string{"ccs", "claude-a", "claude-b", "claude1", "claude2"}
	if len(created) != len(want) {
		t.Fatalf("created %v", created)
	}
	if _, err := os.Lstat(filepath.Join(bin, "claude7")); err == nil {
		t.Fatal("stale link kept")
	}
	if b, _ := os.ReadFile(filepath.Join(bin, "claude9")); string(b) != "user file" {
		t.Fatal("user file touched")
	}
	// Idempotent.
	if created, err := Install(bin, ccam, cfg); err != nil || len(created) != 0 {
		t.Fatalf("second install %v %v", created, err)
	}
	RemoveFor(bin, ccam, &cfg.Profiles[1])
	if _, err := os.Lstat(filepath.Join(bin, "claude2")); err == nil {
		t.Fatal("RemoveFor left claude2")
	}
}

func TestInstallRefusesForeignFile(t *testing.T) {
	dir := t.TempDir()
	ccam := filepath.Join(dir, "ccam")
	os.WriteFile(ccam, []byte("x"), 0o755)
	os.WriteFile(filepath.Join(dir, "claude1"), []byte("mine"), 0o755)
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	if _, err := Install(dir, ccam, cfg); err == nil {
		t.Fatal("want error for foreign claude1")
	}
}
