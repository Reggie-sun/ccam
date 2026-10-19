package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func (h *harness) link(name string) string {
	h.t.Helper()
	target, err := os.Readlink(filepath.Join(h.home, "bin", name))
	if err != nil {
		return ""
	}
	return target
}

func TestAddWithShortcutFlag(t *testing.T) {
	h := newHarness(t)
	r := h.run("", []string{"FAKE_EMAIL=dev@company.com"}, "ccam", "add", "work", "--shortcut", "ccwork")
	if r.code != 0 || r.stdout != "✓ work (dev@company.com) added\n  shortcut: ccwork\n" {
		t.Fatalf("%+v", r)
	}
	if h.link("ccwork") != ccamBin {
		t.Fatalf("link: %q", h.link("ccwork"))
	}
	out := h.ok("ccam", "add", "side", "--no-login")
	if out != "✓ side (not logged in) added\n  no shortcut (add one with `ccam shortcut side ccside`)\n" {
		t.Fatalf("no shortcut: %q", out)
	}

	// A custom shortcut runs Claude Code as its account, arguments untouched.
	h.ok("ccwork", "-p", "hi there", "--continue")
	dir, args := h.called()
	if dir != h.profileDir("work") || !reflect.DeepEqual(args, []string{"-p", "hi there", "--continue"}) {
		t.Fatalf("%s %q", dir, args)
	}
	out = h.ok("ccam", "list")
	if !strings.Contains(out, "ccwork") || !strings.Contains(out, "side  -") {
		t.Fatalf("list:\n%s", out)
	}
	if out := h.ok("ccam", "list", "--json"); !strings.Contains(out, `"shortcut": "ccwork"`) || !strings.Contains(out, `"shortcut": ""`) {
		t.Fatalf("list --json:\n%s", out)
	}
	if out := h.ok("ccam", "doctor"); !strings.Contains(out, "shortcuts OK (ccwork)") {
		t.Fatalf("doctor:\n%s", out)
	}
}

func TestAddRejectsBadShortcuts(t *testing.T) {
	h := newHarness(t)
	elsewhere := filepath.Join(h.home, "elsewhere")
	os.MkdirAll(elsewhere, 0o755)
	os.WriteFile(filepath.Join(elsewhere, "mytool"), []byte("#!/bin/sh\n"), 0o755)
	env := []string{"PATH=" + filepath.Join(h.home, "bin") + ":" + elsewhere + ":/usr/bin:/bin"}
	for name, want := range map[string]string{
		"Bad":    "1-64 characters",
		"ccam":   "taken by ccam",
		"list":   "taken by ccam",
		"mytool": "would hide " + filepath.Join(elsewhere, "mytool"),
	} {
		r := h.run("", env, "ccam", "add", "work", "--no-login", "--shortcut", name)
		if r.code == 0 || !strings.Contains(r.stderr, want) {
			t.Errorf("%s: %+v", name, r)
		}
	}
	if out := h.ok("ccam", "list"); !strings.Contains(out, "No accounts yet") {
		t.Fatalf("a rejected shortcut must not leave the account half added:\n%s", out)
	}
	h.ok("ccam", "add", "work", "--no-login", "--shortcut", "ccwork")
	r := h.run("", nil, "ccam", "add", "home", "--no-login", "--shortcut", "ccwork")
	if r.code == 0 || !strings.Contains(r.stderr, "already the shortcut of account work") {
		t.Fatalf("taken: %+v", r)
	}
}

func TestShortcutCommand(t *testing.T) {
	h := newHarness(t)
	h.ok("ccam", "add", "work", "--no-login")
	if out := h.ok("ccam", "shortcut", "work"); out != "work: none (add one with `ccam shortcut work ccwork`)\n" {
		t.Fatalf("show none: %q", out)
	}
	if out := h.ok("ccam", "--dry-run", "shortcut", "work", "ccw"); !strings.Contains(out, "Would set") || h.link("ccw") != "" {
		t.Fatalf("dry run: %q", out)
	}
	if out := h.ok("ccam", "shortcut", "work", "ccw"); out != "✓ shortcut for work: ccw\n" {
		t.Fatalf("create: %q", out)
	}
	if out := h.ok("ccam", "shortcut", "1"); out != "work: ccw\n" {
		t.Fatalf("show: %q", out)
	}
	if out := h.ok("ccam", "shortcut", "work", "ccjob"); out != "✓ shortcut for work: ccjob (was ccw)\n" {
		t.Fatalf("rename: %q", out)
	}
	if h.link("ccw") != "" || h.link("ccjob") != ccamBin {
		t.Fatal("rename must move the link")
	}
	h.ok("ccjob")
	if dir, _ := h.called(); dir != h.profileDir("work") {
		t.Fatalf("renamed shortcut runs %s", dir)
	}
	r := h.run("", nil, "ccam", "shortcut", "work", "x", "--remove")
	if r.code == 0 {
		t.Fatal("a command and --remove together")
	}
	if out := h.ok("ccam", "shortcut", "work", "--remove"); out != "✓ ccjob removed; work has no shortcut now\n" {
		t.Fatalf("remove: %q", out)
	}
	if h.link("ccjob") != "" {
		t.Fatal("link left behind")
	}
	if out := h.ok("ccam", "shortcut", "work", "--remove"); out != "work has no shortcut.\n" {
		t.Fatalf("remove again: %q", out)
	}
	// remove deletes the account's shortcut too.
	h.ok("ccam", "shortcut", "work", "ccw")
	h.run("y\n", nil, "ccam", "remove", "work")
	if h.link("ccw") != "" {
		t.Fatal("remove left the shortcut")
	}
}

// Configs from before shortcuts were stored had claude<slot> for every
// account. Those links keep working, and the CLI records them.
func TestOldConfigKeepsItsShortcuts(t *testing.T) {
	h := newHarness(t)
	fake := filepath.Join(repoRoot(), "testdata", "fake-claude")
	cfg := fmt.Sprintf("version = 1\nclaude_bin = %q\n\n[[profile]]\nname = \"work\"\nslot = 1\n\n[[profile]]\nname = \"home\"\nslot = 2\n", fake)
	os.WriteFile(filepath.Join(h.home, "config", "config.toml"), []byte(cfg), 0o600)
	for _, name := range []string{"work", "home"} {
		os.MkdirAll(h.profileDir(name), 0o700)
	}
	bin := filepath.Join(h.home, "bin")
	os.MkdirAll(bin, 0o755)
	os.Symlink(ccamBin, filepath.Join(bin, "claude1"))
	os.Symlink(ccamBin, filepath.Join(bin, "claude-home")) // an old --names link
	os.Symlink(ccamBin, filepath.Join(bin, "ccs"))

	// Before any CLI command: the launch path still maps claude1.
	h.ok("claude1", "-c")
	if dir, _ := h.called(); dir != h.profileDir("work") {
		t.Fatalf("claude1 runs %s", dir)
	}
	out := h.ok("ccam", "list")
	if !strings.Contains(out, "claude1") || !strings.Contains(out, "claude-home") {
		t.Fatalf("list:\n%s", out)
	}
	h.ok("ccam", "shims", "install")
	data, _ := os.ReadFile(filepath.Join(h.home, "config", "config.toml"))
	if !strings.Contains(string(data), "version = 2") || !strings.Contains(string(data), "shortcut = 'claude1'") || !strings.Contains(string(data), "shortcut = 'claude-home'") {
		t.Fatalf("saved config:\n%s", data)
	}
	if h.link("ccs") != "" || h.link("claude1") != ccamBin || h.link("claude-home") != ccamBin {
		t.Fatal("shims install keeps configured shortcuts and drops the old ccs link")
	}
	h.ok("claude-home")
	if dir, _ := h.called(); dir != h.profileDir("home") {
		t.Fatalf("claude-home runs %s", dir)
	}
	// Removed shortcuts no longer resolve, even with claude<N> names.
	h.ok("ccam", "shortcut", "work", "--remove")
	os.Symlink(ccamBin, filepath.Join(bin, "claude1"))
	r := h.run("", nil, "claude1", "-c")
	if r.code == 0 || !strings.Contains(r.stderr, "claude1 is not the shortcut of any account") {
		t.Fatalf("an unconfigured claude1 must not launch anything: %+v", r)
	}
}
