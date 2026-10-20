package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddLinksDefaultShares(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	for _, name := range []string{"settings.json", "CLAUDE.md", "skills", "plugins"} {
		target, err := os.Readlink(filepath.Join(h.profileDir("work"), name))
		if err != nil || target != filepath.Join(h.home, "data", "groups", "default", name) {
			t.Fatalf("%s link: %q %v", name, target, err)
		}
	}
	if _, err := os.Lstat(filepath.Join(h.profileDir("work"), "projects")); err == nil {
		t.Fatal("history is local by default")
	}
	// Writing through one account is visible in the other.
	os.WriteFile(filepath.Join(h.profileDir("personal"), "CLAUDE.md"), []byte("shared rule"), 0o600)
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("work"), "CLAUDE.md")); string(b) != "shared rule" {
		t.Fatalf("not shared: %q", b)
	}
}

func TestShareRoundTripCLI(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	work := h.profileDir("work")
	os.MkdirAll(filepath.Join(work, "projects", "p"), 0o700)
	os.WriteFile(filepath.Join(work, "projects", "p", "s1.jsonl"), []byte("work session"), 0o600)
	os.WriteFile(filepath.Join(work, "history.jsonl"), []byte("{\"display\":\"w\"}\n"), 0o600)

	// Dry run changes nothing.
	out := h.ok("ccam", "--dry-run", "group", "set", "default", "--add", "history")
	if !strings.Contains(out, "projects") {
		t.Fatalf("plan output: %s", out)
	}
	if info, _ := os.Lstat(filepath.Join(work, "projects")); info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("dry run applied changes")
	}

	h.ok("ccam", "group", "set", "default", "--add", "history", "--yes")
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "projects", "p", "s1.jsonl")); string(b) != "work session" {
		t.Fatal("personal cannot see work's session")
	}
	h.ok("ccam", "group", "leave", "default", "work", "--yes")
	h.ok("ccam", "group", "join", "default", "work", "--yes")
	if b, _ := os.ReadFile(filepath.Join(work, "projects", "p", "s1.jsonl")); string(b) != "work session" {
		t.Fatal("session lost in round trip")
	}
	if out := h.ok("ccam", "share", "work"); !strings.Contains(out, "group default") {
		t.Fatalf("share status: %s", out)
	}
	h.run("y\n", nil, "ccam", "remove", "work")
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "projects", "p", "s1.jsonl")); string(b) != "work session" {
		t.Fatal("removing an account must not remove shared data")
	}
}

func TestShareConflictPrompt(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "group", "leave", "default", "work", "--yes")
	os.WriteFile(filepath.Join(h.home, "data", "groups", "default", "settings.json"), []byte(`{"model":"opus"}`), 0o600)
	os.Remove(filepath.Join(h.profileDir("work"), "settings.json"))
	os.WriteFile(filepath.Join(h.profileDir("work"), "settings.json"), []byte(`{"theme":"dark"}`), 0o600)
	// Choose 3) JSON merge, then confirm.
	r := h.run("3\ny\n", nil, "ccam", "group", "join", "default", "work")
	if r.code != 0 {
		t.Fatalf("%+v", r)
	}
	b, _ := os.ReadFile(filepath.Join(h.profileDir("work"), "settings.json"))
	if !strings.Contains(string(b), "opus") || !strings.Contains(string(b), "dark") {
		t.Fatalf("merged settings: %s", b)
	}
}

func TestLaunchRepairsMissingLinkButStopsOnRealFile(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	os.Remove(filepath.Join(h.profileDir("work"), "skills"))
	h.ok("ccam", "work")
	if _, err := os.Readlink(filepath.Join(h.profileDir("work"), "skills")); err != nil {
		t.Fatal("missing link not recreated at launch")
	}
	os.Remove(filepath.Join(h.profileDir("work"), "CLAUDE.md"))
	os.WriteFile(filepath.Join(h.profileDir("work"), "CLAUDE.md"), []byte("real"), 0o600)
	os.Remove(h.log)
	r := h.run("", nil, "ccam", "work")
	if r.code == 0 || !strings.Contains(r.stderr, "doctor --fix") {
		t.Fatalf("expected stop: %+v", r)
	}
	if _, err := os.Stat(h.log); err == nil {
		t.Fatal("claude must not start with broken links")
	}
	r = h.run("", nil, "ccam", "doctor")
	if r.code == 0 {
		t.Fatalf("doctor should fail: %s", r.stdout)
	}
	h.ok("ccam", "doctor", "--fix")
	h.ok("ccam", "work")
	// The group had no CLAUDE.md yet, so the real file became the group's.
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "CLAUDE.md")); string(b) != "real" {
		t.Fatalf("doctor --fix must keep the real file's content, got %q", b)
	}
}

// Paths added to history later (checkpoints, task lists, pastes) never stop
// a launch of an account that has its own copy; doctor --fix moves it in.
func TestLatePathsWaitForDoctorFix(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "group", "set", "default", "--add", "history", "--yes")
	fh := filepath.Join(h.profileDir("work"), "file-history")
	os.Remove(fh) // as before this ccam: the account's own folder
	os.MkdirAll(filepath.Join(fh, "s1"), 0o700)
	os.WriteFile(filepath.Join(fh, "s1", "abc@v1"), []byte("backup"), 0o600)
	h.ok("ccam", "work")
	r := h.run("", nil, "ccam", "doctor")
	if !strings.Contains(r.stdout, "file-history is not shared with the group yet") {
		t.Fatalf("%+v", r)
	}
	h.ok("ccam", "doctor", "--fix")
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "file-history", "s1", "abc@v1")); string(b) != "backup" {
		t.Fatal("not moved into the group")
	}
	// The account's own link (to a dotfiles repository) waits as well.
	dotfiles := filepath.Join(h.home, "dotfiles", "keybindings.json")
	os.MkdirAll(filepath.Dir(dotfiles), 0o700)
	os.WriteFile(dotfiles, []byte("{}"), 0o600)
	kb := filepath.Join(h.profileDir("work"), "keybindings.json")
	os.Remove(kb)
	if err := os.Symlink(dotfiles, kb); err != nil {
		t.Skip("no symlinks:", err)
	}
	h.ok("ccam", "work")
	if r := h.run("", nil, "ccam", "doctor"); !strings.Contains(r.stdout, "keybindings.json is not shared with the group yet") {
		t.Fatalf("%+v", r)
	}
}
