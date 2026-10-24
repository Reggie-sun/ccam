package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// refused runs each command, which must refuse because work is running.
func (h *harness) refused(env []string, cmds ...[]string) {
	h.t.Helper()
	for _, argv := range cmds {
		r := h.run("", env, argv...)
		if r.code == 0 || !strings.Contains(r.stderr, "work is running") {
			h.t.Fatalf("%v while work runs: %+v", argv, r)
		}
	}
}

// Claude Code started on an account's folder without ccam (plain claude
// with CLAUDE_CONFIG_DIR, an IDE) keeps the account running for ccam, and
// stops doing so when it ends, even without cleaning up after itself.
func TestClaudeCodeStartedElsewhereKeepsTheAccountRunning(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	claude := exec.Command(fakeBin)
	claude.Env = append(append([]string{}, h.env...), "CLAUDE_CONFIG_DIR="+h.profileDir("work"),
		"FAKE_RECORD=interactive", "FAKE_SLEEP=60", "FAKE_CLAUDE_LOG="+filepath.Join(h.home, "plain.log"))
	if err := claude.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { claude.Process.Kill(); claude.Wait() })
	waitFile(t, filepath.Join(h.profileDir("work"), "sessions", strconv.Itoa(claude.Process.Pid)+".json"))

	h.refused(nil,
		[]string{"ccam", "rename", "work", "job"},
		[]string{"ccam", "remove", "work", "--yes"},
		[]string{"ccam", "group", "leave", "default", "work", "--yes"},
	)
	if !fileExists(h.profileDir("work")) {
		t.Fatal("the folder moved")
	}
	// Killed, it leaves its record behind, which then names no process.
	claude.Process.Kill()
	claude.Wait()
	h.ok("ccam", "rename", "work", "job")
}

// `ccam <account> --bg` returns at once while the session goes on under
// Claude Code's daemon in the account's folder: the account stays running
// until that ends, here killed so that its records stay behind.
func TestBackgroundSessionKeepsTheAccountRunning(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	r := h.run("", []string{"FAKE_DETACH=60"}, "ccam", "work", "--bg", "fix the tests")
	if r.code != 0 || !strings.Contains(r.stdout, "Started background session") {
		t.Fatalf("%+v", r)
	}
	var daemon struct{ PID int }
	data, _ := os.ReadFile(filepath.Join(h.profileDir("work"), "daemon.lock"))
	if err := json.Unmarshal(data, &daemon); err != nil || daemon.PID == 0 {
		t.Fatalf("the background session is not running: %v", err)
	}
	worker, _ := os.FindProcess(daemon.PID)
	t.Cleanup(func() { worker.Kill() })
	h.refused(nil,
		[]string{"ccam", "rename", "work", "job"},
		[]string{"ccam", "remove", "work", "--yes"},
		[]string{"ccam", "group", "set", "default", "--add", "history", "--yes"},
	)
	worker.Kill()
	for deadline := time.Now().Add(10 * time.Second); ; time.Sleep(50 * time.Millisecond) {
		if r := h.run("", nil, "ccam", "rename", "work", "job"); r.code == 0 {
			break
		} else if time.Now().After(deadline) {
			t.Fatalf("rename after the background session ended: %+v", r)
		}
	}
}

// newHomeHarness is a harness without CCAM_HOME: ccam finds its folders
// from HOME and the XDG variables, as it does for a user.
func newHomeHarness(t *testing.T) *harness {
	t.Helper()
	h := newHarness(t)
	var env []string
	for _, kv := range h.env {
		switch name, _, _ := strings.Cut(kv, "="); name {
		case "CCAM_HOME", "HOME", "USERPROFILE", "XDG_CONFIG_HOME", "XDG_DATA_HOME", "XDG_STATE_HOME", "XDG_RUNTIME_DIR":
		default:
			env = append(env, kv)
		}
	}
	h.env = append(env, "HOME="+h.home, "USERPROFILE="+h.home)
	cfg := filepath.Join(h.home, ".config", "ccam", "config.toml")
	os.MkdirAll(filepath.Dir(cfg), 0o700)
	os.WriteFile(cfg, fmt.Appendf(nil, "claude_bin = %q\n", fakeBin), 0o600)
	return h
}

// A session started from a shell with one XDG_RUNTIME_DIR is seen from
// shells with another one or none (cron, ssh, su, a service's terminal).
func TestSessionsAreSeenFromEveryShell(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no XDG_RUNTIME_DIR on Windows")
	}
	h := newHomeHarness(t)
	h.addAccounts()
	first := exec.Command(ccamBin, "work")
	// No FAKE_RECORD: only ccam's own marker shows this session.
	first.Env = append(append([]string{}, h.env...), "XDG_RUNTIME_DIR="+t.TempDir(), "FAKE_SLEEP=60",
		"FAKE_CLAUDE_LOG="+filepath.Join(h.home, "first.log"))
	if err := first.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { first.Process.Kill(); first.Wait() })
	waitFile(t, filepath.Join(h.home, "first.log"))
	h.refused(nil, []string{"ccam", "rename", "work", "job"})
	h.refused([]string{"XDG_RUNTIME_DIR=" + t.TempDir()}, []string{"ccam", "rename", "work", "job"})
}

// A launch goes on when ccam cannot mark the session as running, and an
// XDG_RUNTIME_DIR that cannot be made does not matter.
func TestLaunchWithoutAMarker(t *testing.T) {
	h := newHomeHarness(t)
	h.addAccounts()
	locks := filepath.Join(h.home, ".local", "state", "ccam", "locks")
	os.MkdirAll(filepath.Dir(locks), 0o700)
	if err := os.WriteFile(locks, []byte("not a folder"), 0o600); err != nil {
		t.Fatal(err)
	}
	r := h.run("", nil, "ccam", "work", "-p", "hi")
	if r.code != 0 || !strings.Contains(r.stderr, "could not mark the session as running") {
		t.Fatalf("%+v", r)
	}
	h.called()
	os.Remove(locks)
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		return // the rest relies on Unix permissions
	}
	ro := filepath.Join(t.TempDir(), "ro")
	os.Mkdir(ro, 0o500)
	r = h.run("", []string{"XDG_RUNTIME_DIR=" + filepath.Join(ro, "run")}, "ccam", "work", "-p", "hi")
	if r.code != 0 || strings.Contains(r.stderr, "could not mark") {
		t.Fatalf("%+v", r)
	}
	h.called()
}
