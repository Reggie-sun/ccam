package e2e

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func (h *harness) link(name string) string {
	h.t.Helper()
	target, err := os.Readlink(filepath.Join(h.home, "bin", name+exe))
	if err != nil {
		return ""
	}
	// macOS temp dirs live under /var, a link to /private/var: compare resolved paths.
	a, err1 := filepath.EvalSymlinks(target)
	b, err2 := filepath.EvalSymlinks(ccamBin)
	if err1 == nil && err2 == nil && a == b {
		return ccamBin
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
	os.WriteFile(filepath.Join(elsewhere, "mytool"+exe), []byte("#!/bin/sh\n"), 0o755)
	sep := string(os.PathListSeparator)
	env := []string{"PATH=" + filepath.Join(h.home, "bin") + sep + elsewhere + sep + systemPath}
	for name, want := range map[string]string{
		"Bad":    "1-64 characters",
		"ccam":   "taken by ccam",
		"list":   "taken by ccam",
		"mytool": "would hide " + filepath.Join(elsewhere, "mytool"+exe),
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

// A link to ccam in the bin folder that no account claims (its shortcut was
// changed in config.toml by hand) fails instead of opening another account.
func TestUnclaimedShortcutLinkFails(t *testing.T) {
	h := newHarness(t)
	h.ok("ccam", "add", "work", "--no-login", "--shortcut", "ccwork")
	cfgPath := filepath.Join(h.home, "config", "config.toml")
	data, _ := os.ReadFile(cfgPath)
	os.WriteFile(cfgPath, []byte(strings.Replace(string(data), "shortcut = 'ccwork'", "", 1)), 0o600)
	if h.link("ccwork") != ccamBin {
		t.Fatal("the link should still be there")
	}
	r := h.run("", nil, "ccwork", "-c")
	if r.code == 0 || r.stderr != "ccam: ccwork is not the shortcut of any account (see `ccam list`, `ccam shortcut`)\n" {
		t.Fatalf("an unclaimed shortcut must not launch anything: %+v", r)
	}
	if _, err := os.Stat(h.log); err == nil {
		t.Fatal("claude was called")
	}
}
