package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

	// A failure that is not a usage limit is not retried; its output and code pass through.
	r = h.run("", append(env, "FAKE_EXIT=3"), "ccam", "--fallback", "personal", "-p", "x")
	if r.code != 3 || ranAs() != "personal" || r.stdout != "answer (personal)\n" {
		t.Fatalf("plain failure: %+v", r)
	}

	// Every account at its limit: the last attempt's output and code.
	r = h.run("", append(env, "FAKE_LIMITED=personal,work"), "ccam", "run", "--fallback", "work", "-p", "x")
	if r.code != 1 || ranAs() != "work,personal" || !strings.Contains(r.stdout, "hit your limit") || !strings.Contains(r.stderr, "no other logged-in account has room") {
		t.Fatalf("all limited: %+v", r)
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
}
