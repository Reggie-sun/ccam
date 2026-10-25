package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/transcript"
)

// briefWait is how ccam's first prompt of a brief conversation ends when
// the person gave none.
const briefWait = " Then reply with one or two sentences on where the work stands, and wait for my next message."

// briefConversation writes work's conversation in the folder app of h:
// a title, a prompt, an edit and a reply of about 182k tokens. It returns
// the folder (links resolved) and the conversation's file.
func briefConversation(t *testing.T, h *harness) (string, string) {
	t.Helper()
	project := filepath.Join(h.home, "app")
	os.MkdirAll(project, 0o700)
	real, _ := filepath.EvalSymlinks(project)
	conv := filepath.Join(h.profileDir("work"), "projects", transcript.Slug(real), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	rows := []string{
		`{"type":"ai-title","aiTitle":"Fix the tests"}`,
		`{"parentUuid":null,"isSidechain":false,"type":"user","message":{"role":"user","content":"(fix the failing tests)"},"origin":{"kind":"human"}}`,
		fmt.Sprintf(`{"isSidechain":false,"type":"assistant","message":{"id":"m1","model":"claude-opus-5-5","role":"assistant","content":[{"type":"tool_use","id":"t1","name":"Edit","input":{"file_path":%q}}]}}`, filepath.Join(real, "a_test.go")),
		`{"isSidechain":false,"type":"assistant","message":{"id":"m2","model":"claude-opus-5-5","role":"assistant","content":[{"type":"text","text":"(two of three pass)"}],"usage":{"input_tokens":1000,"cache_read_input_tokens":181000,"output_tokens":400}}}`,
	}
	if err := os.WriteFile(conv, []byte(strings.Join(rows, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return real, conv
}

// `ccam next --continue --brief` starts a new conversation on the next
// account whose first prompt mentions ccam's summary of the saved one.
func TestNextContinueBrief(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	real, conv := briefConversation(t, h)
	brief := filepath.Join(h.home, "state", "brief", "s1.md")
	ask := `Read @"` + brief + `", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here.`

	if r := h.runIn(real, "", nil, "ccam", "next", "--brief"); r.code == 0 || !strings.Contains(r.stderr, "--brief works with --continue") {
		t.Fatalf("%+v", r)
	}
	// A dry run prints the summary and saves none.
	r := h.runIn(real, "", nil, "ccam", "--dry-run", "next", "--continue", "--brief")
	if r.code != 0 || !strings.Contains(r.stdout, "## Claude's last reply\n\n(two of three pass)") || !strings.Contains(r.stdout, "the summary ccam wrote") {
		t.Fatalf("%+v", r)
	}
	if fileExists(brief) {
		t.Fatal("a dry run saved the summary")
	}

	// The person's prompt ends ccam's; their options follow it.
	r = h.runIn(real, "", nil, "ccam", "next", "--continue", "--brief", "--", "--model", "sonnet", "(run them again)")
	if r.code != 0 || !strings.Contains(r.stderr, "→ work → personal") || !strings.Contains(r.stderr, "summary (") ||
		!strings.Contains(r.stderr, `starting a new conversation with ccam's summary of "Fix the tests" instead of resuming it (about 182k tokens)`) {
		t.Fatalf("%+v", r)
	}
	dir, args := h.called()
	if dir != h.profileDir("personal") || !reflect.DeepEqual(args, []string{ask + " Then: (run them again)", "--model", "sonnet"}) {
		t.Fatalf("%s %q", dir, args)
	}
	data, err := os.ReadFile(brief)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, s := range []string{"## Last requests, oldest first\n\n1. (fix the failing tests)", "## Files Claude changed\n\n- a_test.go", filepath.Base(conv)} {
		if !strings.Contains(text, s) {
			t.Errorf("no %q in\n%s", s, text)
		}
	}
	if !strings.HasPrefix(text, "# Summary of an earlier Claude Code conversation: \"Fix the tests\"\n") {
		t.Fatalf("%s", text)
	}
	if info, _ := os.Stat(brief); perms && info.Mode().Perm() != 0o600 {
		t.Fatalf("mode %v", info.Mode())
	}

	// ccam's --mcp-config is one argument, so the prompt after it stays a
	// prompt.
	h.ok("ccam", "mcp", "add-json", "fs", `{"type":"stdio","command":"fs-mcp"}`)
	if r := h.runIn(real, "", nil, "ccam", "next", "--continue", "--brief"); r.code != 0 {
		t.Fatalf("%+v", r)
	}
	want := []string{"--mcp-config=" + filepath.Join(h.home, "data", "groups", "default", "mcp.json"), ask + briefWait}
	if _, args := h.called(); !reflect.DeepEqual(args, want) {
		t.Fatalf("%q", args)
	}

	// A note /ccam:switch left is used up, as by --continue.
	writeNote(t, h, "work", conv, real, "personal")
	if r := h.runIn(real, "", nil, "ccam", "next", "--continue", "--brief"); r.code != 0 {
		t.Fatalf("%+v", r)
	}
	h.called()
	if fileExists(filepath.Join(h.home, "state", "handoff", "work", "s1.json")) {
		t.Fatal("the note stays")
	}
}

// Through npm's claude.cmd the prompt, with its quotes and commas, reaches
// Claude Code as ccam wrote it.
func TestBatchClaudeTakesABrief(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("batch files are Windows'")
	}
	h := newHarness(t)
	shim := filepath.Join(h.home, "npm", "claude.cmd")
	os.MkdirAll(filepath.Dir(shim), 0o700)
	os.WriteFile(shim, []byte("@ECHO off\r\nSETLOCAL\r\nSET \"_prog="+fakeBin+"\"\r\n"+
		"endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & \"%_prog%\" %*\r\n"), 0o600)
	os.WriteFile(filepath.Join(h.home, "config", "config.toml"), fmt.Appendf(nil, "claude_bin = %q\n", shim), 0o600)
	h.addAccounts()
	real, _ := briefConversation(t, h)
	r := h.runIn(real, "", []string{"FAKE_ARGV=ucrt"}, "ccam", "next", "--continue", "--brief")
	if r.code != 0 {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	brief := filepath.Join(h.home, "state", "brief", "s1.md")
	want := []string{`Read @"` + brief + `", the summary ccam wrote of my earlier Claude Code conversation, which is not loaded here.` + briefWait}
	if _, args := h.called(); !reflect.DeepEqual(args, want) {
		t.Fatalf("got\n%q\nwant\n%q", args, want)
	}
}
