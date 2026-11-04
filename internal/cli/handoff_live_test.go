package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func saveStatus(t *testing.T, c *testCLI, account string, fiveHour float64) {
	t.Helper()
	l := usage.Limits{FiveHour: usage.Window{Percent: fiveHour, ResetsAt: time.Now().Add(time.Hour), Known: true}}
	if err := usage.SaveStatusLine(filepath.Join(c.p.State, "usage"), account, l, time.Now()); err != nil {
		t.Fatal(err)
	}
}

// Under handoff_at the hook says nothing; past it the first call asks for
// the whole handoff, later tool calls say nothing until the sixth, each
// prompt asks again, and writing the handoff itself is no step.
func TestHookHandoffAsksPastTheThreshold(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10})
	t.Setenv("CCAM_PROFILE", "personal")
	t.Setenv("CLAUDE_CODE_SESSION_ATTENDED", "")
	t.Setenv("CLAUDE_CODE_ENTRYPOINT", "cli")
	call := func(event, extra string) string {
		out, errOut, code := c.run(`{"session_id":"s1","hook_event_name":"`+event+`"`+extra+`}`, "hook", "handoff")
		if code != 0 {
			t.Fatalf("%d %s", code, errOut)
		}
		return out
	}
	saveStatus(t, c, "personal", 60)
	if out := call("PostToolUse", ""); out != "" {
		t.Fatalf("below the threshold: %q", out)
	}
	saveStatus(t, c, "personal", 86)
	path := livePath(&env{paths: c.p}, "s1")
	first := call("PostToolUse", "")
	if !strings.Contains(first, "86%") || !strings.Contains(first, path) || !strings.Contains(first, "Write it now") {
		t.Fatalf("first: %q", first)
	}
	for i := 2; i < remindEvery; i++ {
		if out := call("PostToolUse", ""); out != "" {
			t.Fatalf("call %d: %q", i, out)
		}
	}
	if out := call("PostToolUse", `,"tool_input":{"file_path":"`+path+`"}`); out != "" {
		t.Fatalf("writing the handoff: %q", out)
	}
	if out := call("UserPromptSubmit", ""); !strings.Contains(out, "rewrite it now") {
		t.Fatalf("a prompt: %q", out)
	}
	t.Setenv("CLAUDE_CODE_SESSION_ATTENDED", "0")
	if out := call("UserPromptSubmit", ""); out != "" {
		t.Fatalf("unattended: %q", out)
	}
}

// A brief names the handoff the earlier session kept.
func TestBriefNamesTheLiveHandoff(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "side": 30})
	trueClaude(t, c)
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	writeConversation(t, conv, time.Now())
	live := livePath(&env{paths: c.p}, "s1")
	os.MkdirAll(filepath.Dir(live), 0o700)
	os.WriteFile(live, []byte("# handoff"), 0o600)
	out, errOut, code := c.run("", "--dry-run", "next", "--continue", "--brief")
	if code != 0 || !strings.Contains(errOut, "handoff kept by that session") || !strings.Contains(out, "the detailed handoff") || !strings.Contains(out, "s1.md") {
		t.Fatalf("%d %s\n%s", code, errOut, out)
	}
}
