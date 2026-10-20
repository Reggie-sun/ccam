package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/session"
)

// If applying a share change fails halfway, the target is already saved so
// `doctor --fix` (or running the same command again) finishes the change
// instead of undoing the half that was done.
func TestFailedShareChangeCanBeFinished(t *testing.T) {
	if os.Geteuid() == 0 || runtime.GOOS == "windows" {
		t.Skip("root, and Windows, ignore the read-only folder this test relies on")
	}
	c := newTestCLI(t, nil)
	work := c.p.ProfileDir("work")
	os.WriteFile(filepath.Join(work, "CLAUDE.md"), []byte("mine"), 0o600)
	group := c.p.GroupDir("default")
	os.MkdirAll(c.p.GroupsRoot(), 0o700)
	os.Mkdir(group, 0o500) // moving into the group fails
	t.Cleanup(func() { os.Chmod(group, 0o700) })

	_, errOut, code := c.run("", "share", "work", "--memory", "on", "--yes")
	if code == 0 || !strings.Contains(errOut, "doctor --fix") {
		t.Fatalf("want failure pointing at doctor --fix, got %d %q", code, errOut)
	}
	cfg, _ := config.Load(c.p.ConfigFile())
	if p, _ := cfg.Find("work"); p.Share.Get("memory") != "default" {
		t.Fatal("target share must be saved before applying")
	}
	os.Chmod(group, 0o700)
	c.run("", "doctor", "--fix") // exit code also reflects unrelated checks
	if b, _ := os.ReadFile(filepath.Join(work, "CLAUDE.md")); string(b) != "mine" {
		t.Fatalf("content after finishing: %q", b)
	}
	if target, _ := os.Readlink(filepath.Join(work, "CLAUDE.md")); target != filepath.Join(group, "CLAUDE.md") {
		t.Fatalf("not linked: %q", target)
	}
}

// Moving data under a running session would lose what it writes meanwhile.
func TestShareRefusesWhileRunning(t *testing.T) {
	c := newTestCLI(t, nil)
	if err := session.Register(launcher.SessionDir(c.p, "work")); err != nil {
		t.Fatal(err)
	}
	if !launcher.Running(c.p, "work") {
		t.Skip("session tracking is not supported on this platform")
	}
	_, errOut, code := c.run("", "share", "work", "--memory", "off", "--yes")
	if code == 0 || !strings.Contains(errOut, "is running") {
		t.Fatalf("want refusal, got %d %q", code, errOut)
	}
	if _, _, code := c.run("", "--dry-run", "share", "work", "--memory", "off"); code != 0 {
		t.Fatal("a dry run only shows the plan and may run")
	}
}
