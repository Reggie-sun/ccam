package shim

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/paths"
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

type binFixture struct {
	dir, bin, ccam string
	p              paths.Paths
}

func newBin(t *testing.T) binFixture {
	dir := t.TempDir()
	f := binFixture{dir: dir, bin: filepath.Join(dir, "bin"), ccam: filepath.Join(dir, "ccam")}
	f.p = paths.Paths{Bin: f.bin, State: filepath.Join(dir, "state")}
	os.WriteFile(f.ccam, []byte("x"), 0o755)
	os.MkdirAll(f.bin, 0o755)
	return f
}

func TestInstallCreatesOnlyConfiguredShortcuts(t *testing.T) {
	f := newBin(t)
	os.WriteFile(filepath.Join(f.bin, "userfile"), []byte("user file"), 0o755) // not ours
	os.Symlink(f.ccam, filepath.Join(f.bin, "mytool"+exe))                     // ours, but not configured: left alone

	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}, {Name: "b", Slot: 2}}}
	created, err := Install(f.p, f.ccam, cfg)
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
	if created, err := Install(f.p, f.ccam, cfg); err != nil || len(created) != 0 {
		t.Fatalf("second install %v %v", created, err)
	}
	Remove(f.p, f.ccam, "cca")
	Remove(f.p, f.ccam, "userfile") // not ours
	if Linked(f.bin, f.ccam, "cca") || !fileExists(filepath.Join(f.bin, "userfile")) {
		t.Fatal("Remove must delete only our link")
	}
}

func TestInstallWithoutShortcutsCreatesNothing(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "bin")
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1}}}
	if created, err := Install(paths.Paths{Bin: bin, State: filepath.Join(dir, "state")}, filepath.Join(dir, "ccam"), cfg); err != nil || len(created) != 0 {
		t.Fatalf("%v %v", created, err)
	}
	if _, err := os.Stat(bin); err == nil {
		t.Fatal("no shortcuts, no bin folder")
	}
}

func fileExists(p string) bool { _, err := os.Lstat(p); return err == nil }

func TestInstallRefusesForeignFile(t *testing.T) {
	f := newBin(t)
	os.WriteFile(filepath.Join(f.bin, "cca"+exe), []byte("mine"), 0o755)
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}}}
	if _, err := Install(f.p, f.ccam, cfg); err == nil {
		t.Fatal("want error for a foreign cca")
	}
}

func TestCheck(t *testing.T) {
	f := newBin(t)
	elsewhere := filepath.Join(f.dir, "usr-bin")
	os.MkdirAll(elsewhere, 0o755)
	os.WriteFile(filepath.Join(elsewhere, "taken"+exe), []byte("#!/bin/sh\n"), 0o755)
	os.WriteFile(filepath.Join(elsewhere, "data"), []byte("not executable"), 0o644)
	os.WriteFile(filepath.Join(f.bin, "foreign"+exe), []byte("x"), 0o755)
	os.Symlink(f.ccam, filepath.Join(f.bin, "ccwork"+exe))
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
		"taken":    "would hide " + filepath.Join(elsewhere, "taken"+exe),
		"foreign":  "is not a ccam shortcut",
	} {
		err := Check(name, f.p, f.ccam, cfg, "home")
		if err == nil || !strings.Contains(err.Error(), wantErr) {
			t.Errorf("%s: %v, want %q", name, err, wantErr)
		}
	}
	for _, ok := range []string{"cchome", "home", "c.h", "data", "h_2"} {
		if err := Check(ok, f.p, f.ccam, cfg, "home"); err != nil {
			t.Errorf("%s: %v", ok, err)
		}
	}
	// The account's own shortcut is fine (renaming to the same name).
	if err := Check("ccwork", f.p, f.ccam, cfg, "work"); err != nil {
		t.Errorf("own shortcut: %v", err)
	}
}

func TestSuggest(t *testing.T) {
	if Suggest("work") != "ccwork" || Suggest("my-side_2") != "ccmy-side_2" {
		t.Fatal(Suggest("work"))
	}
	for _, name := range []string{"work", "a", "my-side_2", "x1234567890123456789012345678901"} {
		if !config.IsShortcut(Suggest(name)) {
			t.Errorf("suggestion for %s is not a valid shortcut name", name)
		}
	}
}

func TestNameOfTheRunningCommand(t *testing.T) {
	cases := map[string]string{"/home/u/.local/bin/ccwork": "ccwork", "ccam": "ccam"}
	if runtime.GOOS == "windows" {
		cases = map[string]string{`C:\Users\u\.local\bin\ccwork.exe`: "ccwork", "CCWORK.EXE": "ccwork", "ccam": "ccam", `C:\x\ccam.exe`: "ccam"}
	}
	for in, want := range cases {
		if got := Name(in); got != want {
			t.Errorf("%s: %s", in, got)
		}
	}
}

// Windows without the right to symlink: shortcuts are hard links, and
// Relink renews them after an update replaced the binary.
func TestHardLinkShortcutsAndRelink(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows only")
	}
	orig := symlink
	symlink = func(string, string) error { return errors.New("A required privilege is not held by the client.") }
	t.Cleanup(func() { symlink = orig })
	f := newBin(t)
	os.WriteFile(f.ccam, []byte("old "+ccamBuild), 0o755)
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}}}
	if created, err := Install(f.p, f.ccam, cfg); err != nil || len(created) != 1 {
		t.Fatalf("%v %v", created, err)
	}
	if !Linked(f.bin, f.ccam, "cca") {
		t.Fatal("a hard link is ours")
	}
	// ccam update: the old binary is renamed, a new one takes its place.
	os.Rename(f.ccam, f.ccam+".old")
	os.WriteFile(f.ccam, []byte("new "+ccamBuild), 0o755)
	if Linked(f.bin, f.ccam, "cca") {
		t.Fatal("still the old binary")
	}
	if renewed := Relink(f.p, f.ccam, cfg); len(renewed) != 1 {
		t.Fatalf("renewed %v", renewed)
	}
	if b, _ := os.ReadFile(filepath.Join(f.bin, "cca.exe")); string(b) != "new "+ccamBuild || !Linked(f.bin, f.ccam, "cca") {
		t.Fatalf("after relink: %q", b)
	}
}

// A hard link ccam made (Windows) keeps running the binary an update or a
// reinstall replaced: Install links it to this ccam (doctor --fix), Check
// takes its name and Remove deletes it. A copy of ccam ccam did not make is
// the user's (on Windows, where hard links were not recorded before, it
// counts as one ccam made), and another program's file stays.
func TestShortcutsThatRunAnOlderCcamAreRenewed(t *testing.T) {
	f := newBin(t)
	os.WriteFile(f.ccam, []byte("old "+ccamBuild), 0o755)
	cca, ccb, ccc := filepath.Join(f.bin, "cca"+exe), filepath.Join(f.bin, "ccb"+exe), filepath.Join(f.bin, "ccc"+exe)
	if err := os.Link(f.ccam, cca); err != nil {
		t.Skip("no hard links here:", err)
	}
	rec := readRecord(f.p)
	rec.add(cca, f.ccam)
	if err := rec.save(f.p); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(ccb, []byte("a copy of ccam: "+ccamBuild), 0o755)
	os.WriteFile(ccc, []byte("another program"), 0o755)
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}, {Name: "b", Slot: 2, Shortcut: "ccb"}, {Name: "c", Slot: 3, Shortcut: "ccc"}}}
	if Stale(f.p, "cca") {
		t.Fatal("a hard link of the binary it was made for is stale")
	}
	// An update: the old binary is renamed, a new one takes its place.
	os.Rename(f.ccam, f.ccam+".old")
	os.WriteFile(f.ccam, []byte("new "+ccamBuild), 0o755)
	windows := runtime.GOOS == "windows"
	if !Stale(f.p, "cca") || Stale(f.p, "ccb") != windows || Stale(f.p, "ccc") {
		t.Fatal("Stale")
	}
	if err := Check("cca", f.p, f.ccam, cfg, "a"); err != nil {
		t.Errorf("Check: %v", err)
	}
	if err := Check("ccc", f.p, f.ccam, cfg, "c"); err == nil {
		t.Error("Check takes another program's file")
	}
	created, err := Install(f.p, f.ccam, cfg)
	if want := map[bool]string{true: "cca ccb", false: "cca"}[windows]; strings.Join(created, " ") != want || err == nil || !strings.Contains(err.Error(), "ccc") {
		t.Fatalf("created %v, %v", created, err)
	}
	if !Linked(f.bin, f.ccam, "cca") {
		t.Fatal("cca still runs the older ccam")
	}
	if b, _ := os.ReadFile(ccc); string(b) != "another program" {
		t.Fatal("another program was replaced")
	}
	if b, _ := os.ReadFile(ccb); !windows && string(b) != "a copy of ccam: "+ccamBuild {
		t.Fatal("the user's copy of ccam was replaced")
	}
	os.Remove(cca)
	os.WriteFile(cca, []byte("an older ccam: "+ccamBuild), 0o755)
	if renewed := Relink(f.p, f.ccam, cfg); strings.Join(renewed, " ") != "cca" {
		t.Fatalf("renewed %v", renewed)
	}
	os.Remove(cca)
	os.WriteFile(cca, []byte("an older ccam: "+ccamBuild), 0o755)
	Remove(f.p, f.ccam, "cca")
	Remove(f.p, f.ccam, "ccc")
	if fileExists(cca) || !fileExists(ccc) {
		t.Fatal("Remove")
	}
}

// A shortcut ccam made for a ccam that has moved since (another install, a
// copy that is gone) is ccam's all the same: Install points it at this
// ccam, Remove deletes it, and Check takes its name. One that runs another
// ccam that is still there stays as it is (a second copy, such as `go run`,
// takes no shortcuts). Links ccam did not make stay, also to a ccam.
func TestStaleShortcutsArePointedAtThisCcam(t *testing.T) {
	f := newBin(t)
	other := filepath.Join(f.dir, "other", "ccam"+exe)
	os.MkdirAll(filepath.Dir(other), 0o755)
	os.WriteFile(other, []byte("another copy"), 0o755)
	tool := filepath.Join(f.dir, "tool")
	os.WriteFile(tool, []byte("#!/bin/sh\n"), 0o755)
	moved := filepath.Join(f.dir, "moved", "ccam"+exe)
	links := map[string]string{"cca": moved, "ccb": other, "ccc": tool, "cce": moved}
	rec := readRecord(f.p)
	for name, target := range links {
		if err := os.Symlink(target, filepath.Join(f.bin, name+exe)); err != nil {
			t.Skip("no symbolic links here:", err)
		}
		if name != "cce" { // a link the user made
			rec.add(filepath.Join(f.bin, name+exe), target)
		}
	}
	rec.save(f.p)
	os.Remove(filepath.Join(f.bin, "ccc"+exe))
	os.Symlink(f.ccam+"-mine", filepath.Join(f.bin, "ccc"+exe)) // changed by the user since
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}, {Name: "b", Slot: 2, Shortcut: "ccb"},
		{Name: "c", Slot: 3, Shortcut: "ccc"}, {Name: "d", Slot: 4, Shortcut: "ccd"}, {Name: "e", Slot: 5, Shortcut: "cce"}}}
	created, err := Install(f.p, f.ccam, cfg)
	in := func(name string) bool {
		return err != nil && strings.Contains(err.Error(), filepath.Join(f.bin, name+exe)+" ")
	}
	if !in("ccc") || !in("cce") || in("cca") || in("ccb") {
		t.Fatalf("want only ccc and cce refused, got %v", err)
	}
	if strings.Join(created, " ") != "cca ccd" {
		t.Fatalf("created %v", created)
	}
	for _, name := range []string{"cca", "ccd"} {
		if !Linked(f.bin, f.ccam, name) {
			t.Errorf("%s does not run this ccam", name)
		}
	}
	if target, _ := os.Readlink(filepath.Join(f.bin, "ccb"+exe)); target != other {
		t.Fatal("a link to another ccam that is still there was replaced")
	}
	if b, _ := os.ReadFile(other); string(b) != "another copy" {
		t.Fatal("the other ccam was touched")
	}
	// Asked to (doctor --fix), it runs this one; links ccam did not make,
	// or made for this one, are no other ccam's.
	for name, want := range map[string]bool{"ccb": true, "cca": false, "ccc": false, "cce": false} {
		if bin, ok := Other(f.p, f.ccam, name); ok != want || ok && bin != other {
			t.Errorf("Other %s: %q %v", name, bin, ok)
		}
	}
	if done, err := Repoint(f.p, f.ccam, []string{"ccb", "ccc", "cce"}); err != nil || strings.Join(done, " ") != "ccb" || !Linked(f.bin, f.ccam, "ccb") {
		t.Fatalf("Repoint: %v %v", done, err)
	}
	if b, _ := os.ReadFile(other); string(b) != "another copy" {
		t.Fatal("the other ccam was touched")
	}
	// Remove and Check see such a link as ccam's too.
	os.Remove(filepath.Join(f.bin, "cca"+exe))
	os.Symlink(moved, filepath.Join(f.bin, "cca"+exe))
	rec = readRecord(f.p)
	rec.add(filepath.Join(f.bin, "cca"+exe), moved)
	rec.save(f.p)
	if err := Check("cca", f.p, f.ccam, cfg, "a"); err != nil {
		t.Errorf("Check: %v", err)
	}
	if err := Check("cce", f.p, f.ccam, cfg, "e"); err == nil {
		t.Error("Check takes a link ccam did not make")
	}
	Remove(f.p, f.ccam, "cca")
	if fileExists(filepath.Join(f.bin, "cca"+exe)) {
		t.Error("Remove left the old link")
	}
	Remove(f.p, f.ccam, "ccc")
	Remove(f.p, f.ccam, "cce")
	if !fileExists(filepath.Join(f.bin, "ccc"+exe)) || !fileExists(filepath.Join(f.bin, "cce"+exe)) {
		t.Error("Remove deleted a link ccam did not make")
	}
}

// A new shortcut name never takes over a file in the bin folder that is
// not this ccam's link: a copy of the binary, or a link to another one.
func TestCheckLeavesTheUsersCopiesOfCcam(t *testing.T) {
	f := newBin(t)
	os.WriteFile(filepath.Join(f.bin, "ccam-backup"+exe), []byte("a copy: "+ccamBuild), 0o755)
	os.Symlink(filepath.Join(f.dir, "opt", "ccam"), filepath.Join(f.bin, "ccam-stable"+exe))
	cfg := &config.Config{Profiles: []config.Profile{{Name: "work", Slot: 1}}}
	for _, name := range []string{"ccam-backup", "ccam-stable"} {
		if err := Check(name, f.p, f.ccam, cfg, "work"); err == nil || !strings.Contains(err.Error(), "is not a ccam shortcut") {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// One file in the way of one shortcut leaves the other shortcuts to be
// made, whatever order the accounts come in.
func TestInstallGoesOnPastAFileInTheWay(t *testing.T) {
	for range 20 {
		f := newBin(t)
		os.WriteFile(filepath.Join(f.bin, "cca"+exe), []byte("mine"), 0o755)
		cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}, {Name: "b", Slot: 2, Shortcut: "ccb"},
			{Name: "c", Slot: 3, Shortcut: "ccc"}, {Name: "d", Slot: 4, Shortcut: "ccd"}}}
		created, err := Install(f.p, f.ccam, cfg)
		if err == nil || !strings.Contains(err.Error(), filepath.Join(f.bin, "cca"+exe)+" ") {
			t.Fatalf("want cca refused, got %v", err)
		}
		if strings.Join(created, " ") != "ccb ccc ccd" {
			t.Fatalf("created %v", created)
		}
		if b, _ := os.ReadFile(filepath.Join(f.bin, "cca"+exe)); string(b) != "mine" {
			t.Fatal("the file in the way was touched")
		}
	}
}
