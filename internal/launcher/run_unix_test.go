//go:build !windows

package launcher

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
)

// Run records the claude process and what reached ccam: SIGHUP (passed on)
// and Ctrl+C (not passed on). The signal goroutine outlives Run, so this
// also checks under -race.
func TestRunRecordsSignals(t *testing.T) {
	p, _ := setup(t)
	// Keep the test process alive whatever reaches it.
	keep := make(chan os.Signal, 4)
	signal.Notify(keep, syscall.SIGHUP, os.Interrupt)
	defer signal.Stop(keep)
	sleep, _ := exec.LookPath("sleep")
	run := func(send syscall.Signal) (*Plan, int) {
		pl := &Plan{Profile: &config.Profile{Name: "work"}, Bin: sleep, Argv: []string{"sleep", "0.6"}}
		go func() {
			time.Sleep(200 * time.Millisecond)
			syscall.Kill(os.Getpid(), send)
		}()
		code, err := pl.Run(p, nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		return pl, code
	}
	pl, code := run(syscall.SIGHUP)
	if code != 129 || !pl.run.signaled.Load() || pl.run.interrupt.Load() != 0 || pl.run.child <= 0 {
		t.Fatalf("SIGHUP: code %d, signaled %v, interrupt %d, child %d", code, pl.run.signaled.Load(), pl.run.interrupt.Load(), pl.run.child)
	}
	before := time.Now()
	pl, code = run(syscall.SIGINT)
	if at := time.Unix(0, pl.run.interrupt.Load()); code != 0 || pl.run.signaled.Load() || at.Before(before) {
		t.Fatalf("Ctrl+C: code %d, signaled %v, at %v", code, pl.run.signaled.Load(), at)
	}
}

// A pipe is no terminal: nothing is dropped, and nothing said to be (it
// counts the bytes waiting in it too).
func TestDiscardInputLeavesAPipe(t *testing.T) {
	r, w, _ := os.Pipe()
	w.Write([]byte("typed\n"))
	if discardInput(r) {
		t.Fatal("said it dropped what a pipe holds")
	}
	w.Close()
	if rest, _ := io.ReadAll(r); string(rest) != "typed\n" {
		t.Fatalf("%q left", rest)
	}
}

// A claude_bin that exists but cannot be run (no shebang, not a program)
// is named in the error by both paths. Run leaves the account's last use
// alone; Exec marks it before syscall.Exec, which replaces ccam when it
// works (left as is: a failed exec is rare and only reorders ccam next).
func TestExecErrorNamesTheBinary(t *testing.T) {
	p, cfg := setup(t)
	bin := filepath.Join(t.TempDir(), "claude")
	os.WriteFile(bin, []byte{0, 1, 2, 3}, 0o755) // ENOEXEC
	cfg.ClaudeBin = bin
	pl, err := Prepare(p, cfg, "work", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pl.Run(p, nil, nil, nil); err == nil || !strings.Contains(err.Error(), bin) {
		t.Fatalf("Run: %v", err)
	}
	if !LastUsed(p, "work").IsZero() {
		t.Fatal("Run marked the account as used although claude did not start")
	}
	if err := pl.Exec(p); err == nil || !strings.Contains(err.Error(), bin) {
		t.Fatalf("Exec: %v", err)
	}
}

// The --dry-run line, pasted into sh, gives claude the environment a launch
// gives it: an API account's endpoint and key variable, the account's own
// CCAM_PROFILE, and none of what a launch drops (an outer session's
// CCAM_PROFILE and CCAM_RELAY, a key in the shell).
func TestDryRunLineGivesTheLaunchEnv(t *testing.T) {
	p, cfg := setup(t)
	os.MkdirAll(p.ProfileDir("gw"), 0o700)
	cfg.Profiles = append(cfg.Profiles, config.Profile{Name: "gw", Slot: 3,
		API: &config.API{BaseURL: "https://llm.example.com/v1", Model: "m1", Auth: config.Bearer, KeyFrom: "env:GW_KEY"}})
	t.Setenv("GW_KEY", "secret-key")
	t.Setenv("CCAM_PROFILE", "personal")
	t.Setenv("CCAM_RELAY", "2:1")
	t.Setenv("ANTHROPIC_API_KEY", "from-the-shell")
	envOut := filepath.Join(t.TempDir(), "env")
	cfg.ClaudeBin = filepath.Join(t.TempDir(), "claude")
	os.WriteFile(cfg.ClaudeBin, []byte("#!/bin/sh\nenv > '"+envOut+"'\n"), 0o755)
	watched := func(env []string) []string {
		var out []string
		for _, kv := range env {
			if k, _, _ := strings.Cut(kv, "="); strings.HasPrefix(k, "CCAM_") || strings.HasPrefix(k, "ANTHROPIC_") || k == "CLAUDE_CONFIG_DIR" {
				out = append(out, kv)
			}
		}
		slices.Sort(out)
		return out
	}
	for _, name := range []string{"gw", "work"} {
		dry, err := Prepare(p, cfg, name, []string{"-p", "hi"})
		if err != nil {
			t.Fatal(err)
		}
		real, _ := Prepare(p, cfg, name, []string{"-p", "hi"})
		if api := dry.Profile.API; api != nil {
			if err := dry.useAPI(api, true); err != nil {
				t.Fatal(err)
			}
			if err := real.useAPI(api, false); err != nil {
				t.Fatal(err)
			}
		}
		var line strings.Builder
		dry.DryRun(&line)
		if strings.Contains(line.String(), "secret-key") || strings.Contains(line.String(), "from-the-shell") {
			t.Fatalf("the line shows a key: %s", line.String())
		}
		if out, err := exec.Command("sh", "-c", line.String()).CombinedOutput(); err != nil {
			t.Fatalf("%s: %v %s", line.String(), err, out)
		}
		data, _ := os.ReadFile(envOut)
		if got, want := watched(strings.Split(strings.TrimSpace(string(data)), "\n")), watched(real.Env); !slices.Equal(got, want) {
			t.Errorf("%s: the pasted line gives\n%q\nwant\n%q\n(line: %s)", name, got, want, line.String())
		}
	}
}
