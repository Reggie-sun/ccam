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
	out := h.ok("ccam", "--dry-run", "share", "work", "--history", "on")
	if !strings.Contains(out, "projects") {
		t.Fatalf("plan output: %s", out)
	}
	if info, _ := os.Lstat(filepath.Join(work, "projects")); info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("dry run applied changes")
	}

	h.ok("ccam", "share", "work", "--history", "on", "--yes")
	h.ok("ccam", "share", "personal", "--history", "on", "--yes")
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "projects", "p", "s1.jsonl")); string(b) != "work session" {
		t.Fatal("personal cannot see work's session")
	}
	h.ok("ccam", "share", "work", "--history", "off", "--yes")
	h.ok("ccam", "share", "work", "--history", "on", "--yes")
	if b, _ := os.ReadFile(filepath.Join(work, "projects", "p", "s1.jsonl")); string(b) != "work session" {
		t.Fatal("session lost in round trip")
	}
	if out := h.ok("ccam", "share", "work"); !strings.Contains(out, "공유 (default 그룹)") {
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
	h.ok("ccam", "share", "work", "--settings", "off", "--yes")
	os.WriteFile(filepath.Join(h.home, "data", "groups", "default", "settings.json"), []byte(`{"model":"opus"}`), 0o600)
	os.Remove(filepath.Join(h.profileDir("work"), "settings.json"))
	os.WriteFile(filepath.Join(h.profileDir("work"), "settings.json"), []byte(`{"theme":"dark"}`), 0o600)
	// Choose 3) JSON merge, then confirm.
	r := h.run("3\ny\n", nil, "ccam", "share", "work", "--settings", "on")
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
	h.ok("ccs", "work")
	if _, err := os.Readlink(filepath.Join(h.profileDir("work"), "skills")); err != nil {
		t.Fatal("missing link not recreated at launch")
	}
	os.Remove(filepath.Join(h.profileDir("work"), "CLAUDE.md"))
	os.WriteFile(filepath.Join(h.profileDir("work"), "CLAUDE.md"), []byte("real"), 0o600)
	os.Remove(h.log)
	r := h.run("", nil, "ccs", "work")
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
	h.ok("ccs", "work")
	// The group had no CLAUDE.md yet, so the real file became the group's.
	if b, _ := os.ReadFile(filepath.Join(h.profileDir("personal"), "CLAUDE.md")); string(b) != "real" {
		t.Fatalf("doctor --fix must keep the real file's content, got %q", b)
	}
}
