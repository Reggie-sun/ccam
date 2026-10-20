package e2e

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestAPIKeyAccount(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	envFile := filepath.Join(h.home, "env")
	names := "FAKE_ENV_NAMES=ANTHROPIC_API_KEY,ANTHROPIC_AUTH_TOKEN,ANTHROPIC_BASE_URL,ANTHROPIC_MODEL,CLAUDE_CODE_OAUTH_TOKEN,CLAUDE_CODE_USE_BEDROCK"
	seen := func() string {
		data, _ := os.ReadFile(envFile)
		os.Remove(envFile)
		return string(data)
	}

	r := h.run("", nil, "ccam", "add", "gw", "--api", "--base-url", "https://llm.example.com/v1", "--auth", "bearer",
		"--model", "claude-sonnet-5-5", "--key-from", "env:GW_KEY")
	if r.code != 0 || !strings.Contains(r.stdout, "gw (API: llm.example.com) added") || !strings.Contains(r.stderr, "$GW_KEY") {
		t.Fatalf("add: %+v", r)
	}
	if out := h.ok("ccam", "list"); !strings.Contains(out, "(API: llm.example.com)") {
		t.Fatalf("list:\n%s", out)
	}

	// The key and endpoint reach claude; inherited credentials do not.
	inherited := []string{"FAKE_ENV_FILE=" + envFile, names, "GW_KEY=sk-gw-123", "ANTHROPIC_API_KEY=sk-from-shell",
		"CLAUDE_CODE_OAUTH_TOKEN=oauth-from-shell", "CLAUDE_CODE_USE_BEDROCK=1"}
	if r := h.run("", inherited, "ccam", "gw", "-p", "hi"); r.code != 0 {
		t.Fatalf("launch: %+v", r)
	}
	want := "ANTHROPIC_AUTH_TOKEN=sk-gw-123\nANTHROPIC_BASE_URL=https://llm.example.com/v1\nANTHROPIC_MODEL=claude-sonnet-5-5\n"
	if got := seen(); got != want {
		t.Fatalf("env:\n%s\nwant:\n%s", got, want)
	}
	// A plain account keeps whatever the shell has.
	h.run("", inherited, "ccam", "work", "-p", "hi")
	if got := seen(); !strings.Contains(got, "ANTHROPIC_API_KEY=sk-from-shell") {
		t.Fatalf("plain account env:\n%s", got)
	}

	// An empty key stops the launch before claude starts.
	os.Remove(h.log)
	if r := h.run("", []string{"GW_KEY="}, "ccam", "gw", "-p", "hi"); r.code == 0 || !strings.Contains(r.stderr, "$GW_KEY") {
		t.Fatalf("empty key: %+v", r)
	}
	if _, err := os.Stat(h.log); err == nil {
		t.Fatal("claude ran without a key")
	}

	// command: sources are run at launch only, never by a dry run.
	marker := filepath.Join(h.home, "command-ran")
	cmd := "echo sk-cmd-456 && echo ran > " + marker
	if runtime.GOOS == "windows" {
		cmd = "echo sk-cmd-456&& echo ran > " + marker
	}
	h.ok("ccam", "add", "direct", "--api", "--key-from", "command:"+cmd)
	h.ok("ccam", "--dry-run", "direct", "-p", "hi")
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("a dry run ran the key command")
	}
	h.run("", []string{"FAKE_ENV_FILE=" + envFile, names}, "ccam", "direct", "-p", "hi")
	if got := seen(); got != "ANTHROPIC_API_KEY=sk-cmd-456\n" {
		t.Fatalf("command key env: %q", got)
	}

	// No login to run or remove, and usage says API.
	if r := h.run("", nil, "ccam", "login", "gw"); r.code == 0 || !strings.Contains(r.stderr, "API key account") {
		t.Fatalf("login: %+v", r)
	}
	if out := h.run("", nil, "ccam", "usage").stdout; !strings.Contains(out, "gw") || !strings.Contains(out, "API") {
		t.Fatalf("usage:\n%s", out)
	}
	if r := h.run("", nil, "ccam", "add", "bad", "--api", "--key-from", "file:/x"); r.code == 0 {
		t.Fatal("a bad key source was accepted")
	}
	if r := h.run("", nil, "ccam", "add", "bad", "--base-url", "https://x"); r.code == 0 || !strings.Contains(r.stderr, "--api") {
		t.Fatalf("--base-url without --api: %+v", r)
	}
}

// A key or token in the shell would quietly replace a login account's own.
func TestLaunchWarnsAboutAuthVariables(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	r := h.run("", []string{"ANTHROPIC_API_KEY=sk-x", "CLAUDE_CODE_OAUTH_TOKEN=", "ANTHROPIC_AUTH_TOKEN="}, "ccam", "work", "-p", "hi")
	if r.code != 0 || !strings.Contains(r.stderr, "$ANTHROPIC_API_KEY is set") || !strings.Contains(r.stderr, "work's login") {
		t.Fatalf("%+v", r)
	}
	r = h.run("", []string{"ANTHROPIC_API_KEY=", "CLAUDE_CODE_OAUTH_TOKEN=", "ANTHROPIC_AUTH_TOKEN="}, "ccam", "work", "-p", "hi")
	if strings.Contains(r.stderr, "is set") {
		t.Fatalf("no warning without a value: %+v", r)
	}
}
