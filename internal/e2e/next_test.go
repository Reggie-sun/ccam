package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNextRotatesAndContinues(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccs", "personal")
	h.called()
	time.Sleep(1100 * time.Millisecond) // last-used has 1 s resolution
	h.ok("ccam", "next", "--", "-p", "hi")
	dir, args := h.called()
	if dir != h.profileDir("work") || !reflect.DeepEqual(args, []string{"-p", "hi"}) {
		t.Fatalf("next: %s %q", dir, args)
	}
	// --continue needs a shared history group.
	r := h.run("", nil, "ccam", "next", "--continue")
	if r.code == 0 || !strings.Contains(r.stderr, "--history on") {
		t.Fatalf("expected history error: %+v", r)
	}
	h.ok("ccam", "share", "personal", "--history", "on", "--yes")
	h.ok("ccam", "share", "work", "--history", "on", "--yes")
	h.ok("ccam", "next", "--continue", "--", "-p", "go on")
	dir, args = h.called()
	if dir != h.profileDir("personal") || !reflect.DeepEqual(args, []string{"--continue", "-p", "go on"}) {
		t.Fatalf("next --continue: %s %q", dir, args)
	}
}

func TestHistoryLockWarnsOnConcurrentUse(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "share", "personal", "--history", "on", "--yes")
	h.ok("ccam", "share", "work", "--history", "on", "--yes")
	first := exec.Command(filepath.Join(h.home, "bin", "ccs"), "personal")
	first.Env = append(append([]string{}, h.env...), "FAKE_SLEEP=2", "FAKE_CLAUDE_LOG="+filepath.Join(h.home, "first.log"))
	if err := first.Start(); err != nil {
		t.Fatal(err)
	}
	defer first.Wait()
	deadline := time.Now().Add(2 * time.Second)
	for !fileExists(filepath.Join(h.home, "first.log")) && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	r := h.run("", nil, "ccs", "work")
	if r.code != 0 || !strings.Contains(r.stderr, "다른 Claude Code") {
		t.Fatalf("expected warning, got %+v", r)
	}
	// A different folder is not affected.
	cmd := exec.Command(filepath.Join(h.home, "bin", "ccs"), "work")
	cmd.Dir = t.TempDir()
	cmd.Env = h.env
	out, _ := cmd.CombinedOutput()
	if strings.Contains(string(out), "다른 Claude Code") {
		t.Fatalf("other folder warned: %s", out)
	}
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// A background job started inside the session must not keep the account
// looking "running" after claude itself exits.
func TestLeftoverChildDoesNotHoldSession(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "share", "personal", "--history", "on", "--yes")
	h.ok("ccam", "share", "work", "--history", "on", "--yes")
	r := h.run("", []string{"FAKE_BG=3"}, "ccs", "personal")
	if r.code != 0 {
		t.Fatalf("%+v", r)
	}
	r = h.run("", nil, "ccs", "work")
	if strings.Contains(r.stderr, "다른 Claude Code") {
		t.Fatalf("finished session still counted as running: %s", r.stderr)
	}
}
