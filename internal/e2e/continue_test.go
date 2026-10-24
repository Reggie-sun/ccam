package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/transcript"
)

// ranIt writes into account's .claude.json what Claude Code (2.1.308) keeps
// there when a session in folder ends: projects[folder].lastSessionId.
func ranIt(t *testing.T, h *harness, account, folder, session string) {
	t.Helper()
	path := filepath.Join(h.profileDir(account), ".claude.json")
	data, _ := os.ReadFile(path)
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	v["projects"] = map[string]any{folder: map[string]any{"lastSessionId": session, "lastCost": 0.4}}
	data, _ = json.Marshal(v)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// usedAgo sets when ccam last started each account, seconds ago.
func usedAgo(t *testing.T, h *harness, ago map[string]int) {
	t.Helper()
	dir := filepath.Join(h.home, "state", "last-used")
	os.MkdirAll(dir, 0o700)
	for name, s := range ago {
		os.WriteFile(filepath.Join(dir, name), []byte(strconv.FormatInt(time.Now().Unix()-int64(s), 10)), 0o600)
	}
}

// With shared history every account sees the conversation; the one it
// belongs to is the one that ran it, which Claude Code notes in that
// account's .claude.json, not the account used last in another folder.
func TestContinueFromTheAccountThatRanIt(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "add", "side")
	h.ok("ccam", "group", "set", "default", "--add", "history", "--yes")
	project := filepath.Join(h.home, "app")
	os.MkdirAll(project, 0o700)
	real, _ := filepath.EvalSymlinks(project)
	conv := filepath.Join(h.profileDir("personal"), "projects", transcript.Slug(real), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	at := time.Now().Add(-9 * time.Minute).UTC().Format(time.RFC3339Nano)
	os.WriteFile(conv, fmt.Appendf(nil, `{"type":"user","message":{"role":"user","content":"Fix the tests"},"timestamp":%q}`+"\n"+
		`{"type":"assistant","message":{"model":"claude-opus-5-5","content":[{"type":"text","text":"Done."}]},"timestamp":%q}`+"\n", at, at), 0o600)
	// personal ran it 10 minutes ago; side was started since, elsewhere.
	usedAgo(t, h, map[string]int{"personal": 600, "work": 3000, "side": 60})
	cont := func() (string, string) {
		t.Helper()
		r := h.runIn(project, "", nil, "ccam", "--dry-run", "next", "--continue")
		if r.code != 0 {
			t.Fatalf("%+v", r)
		}
		return r.stderr, r.stdout
	}
	// Nobody says who ran it: the account used last counts, as before.
	if errOut, _ := cont(); !strings.Contains(errOut, "→ side → personal") {
		t.Fatalf("without a word from Claude Code: %s", errOut)
	}
	ranIt(t, h, "personal", real, "s1")
	if errOut, out := cont(); !strings.Contains(errOut, "→ personal → work") || !strings.Contains(out, h.profileDir("work")) {
		t.Fatalf("personal ran it: %s\n%s", errOut, out)
	}
	// It moved on to work with --resume; work ran it last of the two.
	ranIt(t, h, "work", real, "s1")
	usedAgo(t, h, map[string]int{"personal": 900, "work": 300, "side": 60})
	if errOut, _ := cont(); !strings.Contains(errOut, "→ work → side") {
		t.Fatalf("personal and work ran it: %s", errOut)
	}
}
