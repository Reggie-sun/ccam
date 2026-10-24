package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFallbackRunsAgainOnAnotherAccountAtTheLimit(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	runs := filepath.Join(h.home, "runs")
	stdinFile := filepath.Join(h.home, "stdin")
	env := []string{"FAKE_RUNS=" + runs, "FAKE_CLAUDE_STDIN=" + stdinFile, "FAKE_STDOUT=answer"}
	ranAs := func() string {
		data, _ := os.ReadFile(runs)
		os.Remove(runs)
		return strings.Join(strings.Fields(string(data)), ",")
	}

	r := h.run("piped input\n", append(env, "FAKE_LIMITED=personal"), "ccam", "--fallback", "personal", "-p", "summarise")
	if r.code != 0 || ranAs() != "personal,work" {
		t.Fatalf("retry: %+v", r)
	}
	if r.stdout != "answer (work)\n" {
		t.Fatalf("only the run that worked may print: %q", r.stdout)
	}
	if !strings.Contains(r.stderr, "personal reached its usage limit; running the command again as work") {
		t.Fatalf("stderr: %s", r.stderr)
	}
	if b, _ := os.ReadFile(stdinFile); string(b) != "piped input\n" {
		t.Fatalf("the retry must get the same input: %q", b)
	}
	if dir, args := h.called(); dir != h.profileDir("work") || strings.Join(args, " ") != "-p summarise" {
		t.Fatalf("retry ran %s %q", dir, args)
	}

	// Claude Code's usual message names the window.
	r = h.run("", append(env, "FAKE_LIMITED=personal", "FAKE_LIMIT_TEXT=You've hit your session limit · resets 3:45pm (Asia/Seoul)"), "ccam", "--fallback", "personal", "-p", "x")
	if r.code != 0 || ranAs() != "personal,work" || r.stdout != "answer (work)\n" {
		t.Fatalf("session limit: %+v", r)
	}

	// A failure that is not a usage limit is not retried; its output and code pass through.
	r = h.run("", append(env, "FAKE_EXIT=3"), "ccam", "--fallback", "personal", "-p", "x")
	if r.code != 3 || ranAs() != "personal" || r.stdout != "answer (personal)\n" {
		t.Fatalf("plain failure: %+v", r)
	}

	// Every account at its limit: the last attempt's output, and exit code
	// 75 (EX_TEMPFAIL) so a script can try again later.
	r = h.run("", append(env, "FAKE_LIMITED=personal,work"), "ccam", "run", "--fallback", "work", "-p", "x")
	if r.code != 75 || ranAs() != "work,personal" || !strings.Contains(r.stdout, "hit your limit") || !strings.Contains(r.stderr, "no other logged-in account has room") {
		t.Fatalf("all limited: %+v", r)
	}
	// --wait without --fallback is refused; a bad duration too.
	if r := h.run("", env, "ccam", "--wait", "work", "-p", "x"); r.code == 0 || !strings.Contains(r.stderr, "--fallback") || ranAs() != "" {
		t.Fatalf("--wait alone: %+v", r)
	}
	if r := h.run("", env, "ccam", "--fallback", "--wait=soon", "work", "-p", "x"); r.code == 0 || !strings.Contains(r.stderr, "duration") {
		t.Fatalf("bad --wait: %+v", r)
	}
	// A value after a space would be the account: the message says how to write it.
	if r := h.run("", env, "ccam", "--fallback", "--wait", "30m", "-p", "x"); r.code == 0 || !strings.Contains(r.stderr, "--wait=30m") || ranAs() != "" {
		t.Fatalf("--wait 30m: %+v", r)
	}
	// A zero or a number without a unit gets examples that --wait takes.
	for _, v := range []string{"0", "30"} {
		if r := h.run("", env, "ccam", "--fallback", "--wait", v, "-p", "x"); r.code == 0 || !strings.Contains(r.stderr, "such as --wait=30m or --wait=2h") || ranAs() != "" {
			t.Fatalf("--wait %s: %+v", v, r)
		}
	}

	// No name: the account for the folder or the default, as without --fallback.
	h.ok("ccam", "use", "work")
	r = h.run("", append(env, "FAKE_LIMITED=work"), "ccam", "--fallback", "-p", "x")
	if r.code != 0 || ranAs() != "work,personal" {
		t.Fatalf("default account: %+v", r)
	}

	// Interactive runs are refused rather than run without the safety net.
	r = h.run("", env, "ccam", "--fallback", "work")
	if r.code == 0 || !strings.Contains(r.stderr, "-p or --print") || ranAs() != "" {
		t.Fatalf("interactive: %+v", r)
	}
	out := h.ok("ccam", "--dry-run", "--fallback", "work", "-p", "x")
	if !strings.Contains(out, "CLAUDE_CONFIG_DIR=") || ranAs() != "" {
		t.Fatalf("dry run: %s", out)
	}

	// A leading "--", as for any launch.
	for _, argv := range [][]string{{"--fallback", "personal", "--", "-p", "x"}, {"--fallback", "run", "personal", "--", "-p", "x"}, {"run", "--fallback", "personal", "--", "--print", "x"}} {
		r = h.run("", append(env, "FAKE_LIMITED=personal"), append([]string{"ccam"}, argv...)...)
		if r.code != 0 || ranAs() != "personal,work" {
			t.Fatalf("%q: %+v", argv, r)
		}
		if _, args := h.called(); strings.Join(args, " ") != strings.Join(argv[len(argv)-2:], " ") {
			t.Fatalf("%q: claude got %q", argv, args)
		}
	}
	if r := h.run("", env, "ccam", "--fallback", "personal", "--", "--", "-p", "x"); r.code == 0 || !strings.Contains(r.stderr, "-p or --print") {
		t.Fatalf("after a second -- -p is a prompt: %+v", r)
	}
}

// A pipe that stays open and sends nothing (ssh without -n, Node's spawn)
// does not keep --fallback from starting claude: like claude -p it goes on
// without input after a moment.
func TestFallbackGoesOnWhenStdinStaysOpenAndSilent(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	runs := filepath.Join(h.home, "runs")
	stdinFile := filepath.Join(h.home, "stdin")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	cmd := exec.Command(ccamBin, "--fallback", "personal", "-p", "x")
	cmd.Env = append(append([]string{}, h.env...), "FAKE_RUNS="+runs, "FAKE_CLAUDE_STDIN="+stdinFile, "FAKE_STDOUT=answer", "FAKE_LIMITED=personal")
	cmd.Stdin = r
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	r.Close()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		data, _ := os.ReadFile(runs)
		if err != nil || out.String() != "answer (work)\n" || strings.Join(strings.Fields(string(data)), ",") != "personal,work" {
			t.Fatalf("%v: stdout %q, stderr %q, runs %q", err, out.String(), errb.String(), data)
		}
		if b, err := os.ReadFile(stdinFile); err != nil || len(b) != 0 {
			t.Fatalf("claude got %q (%v)", b, err)
		}
		if !strings.Contains(errb.String(), "no input on stdin") {
			t.Fatalf("stderr: %s", errb.String())
		}
	case <-time.After(15 * time.Second):
		cmd.Process.Kill()
		t.Fatalf("--fallback still waits for stdin; stderr: %q", errb.String())
	}
}
