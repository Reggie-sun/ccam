package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// An account Claude Code refused at its limit is not chosen again before
// the reset, even when its saved usage is old and low: here personal's
// token has expired and its last lookup, three hours ago, said 10%.
func TestLimitRecordKeepsTheAccountOut(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 30, "side": 40})
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin = stubClaude
	cfg.Save(c.p.ConfigFile())
	now := time.Now()
	os.WriteFile(filepath.Join(c.p.ProfileDir("personal"), ".credentials.json"),
		fmt.Appendf(nil, `{"claudeAiOauth":{"accessToken":"pct-10","expiresAt":%d}}`, now.Add(-time.Hour).UnixMilli()), 0o600)
	dir := filepath.Join(c.p.State, "usage")
	os.MkdirAll(dir, 0o700)
	os.WriteFile(filepath.Join(dir, "personal.json"), fmt.Appendf(nil, `{"at":%q,"limits":{"five_hour":{"percent":10,"known":true,"resets_at":%q}}}`,
		now.Add(-3*time.Hour).UTC().Format(time.RFC3339), now.Add(time.Hour).UTC().Format(time.RFC3339)), 0o600)
	if err := usage.SaveLimit(dir, "personal", "five_hour", now.Add(2*time.Hour), now.Add(-time.Minute)); err != nil {
		t.Fatal(err)
	}
	// work's session left a note in this folder: continuing leaves work out.
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("work"), cwd), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	os.WriteFile(conv, []byte("{}\n"), 0o600)
	handoff.Save(c.p.State, handoff.Ticket{Account: "work", SessionID: "s1", Transcript: conv, Cwd: cwd, At: now, Reason: handoff.Switch})

	choose := func() (string, string) {
		_, errOut, code := c.run("", "--dry-run", "next", "--auto", "--continue")
		out, _, _ := c.run("", "--dry-run", "pick")
		if code != 0 {
			t.Fatalf("next: %d %s", code, errOut)
		}
		return errOut, strings.TrimSpace(out)
	}
	if next, pick := choose(); !strings.Contains(next, "→ work → side") || pick != "work" {
		t.Fatalf("with the record: next %q, pick %q", next, pick)
	}
	// Without it, the old 10% wins: the record is what keeps personal out.
	os.Remove(filepath.Join(dir, "personal.limit.json"))
	if next, pick := choose(); !strings.Contains(next, "→ work → personal") || pick != "personal" {
		t.Fatalf("without the record: next %q, pick %q", next, pick)
	}
}

// The hook may look before Claude Code has written the refused request;
// carrying the conversation on records the limit from it then.
func TestNextContinueRecordsTheLimitTheHookMissed(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40})
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin = stubClaude
	cfg.Save(c.p.ConfigFile())
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	now := time.Now().Truncate(time.Second)
	reset := now.Add(2 * time.Hour)
	limitConversation(t, conv, now.Add(-30*time.Second), reset)
	handoff.Save(c.p.State, handoff.Ticket{Account: "personal", SessionID: "s1", Transcript: conv, Cwd: cwd, At: now, Reason: handoff.Limit, Limit: "session"})
	record := filepath.Join(c.p.State, "usage", "personal.limit.json")
	// A dry run writes nothing.
	if _, errOut, code := c.run("", "--dry-run", "next", "--continue"); code != 0 {
		t.Fatalf("%d %s", code, errOut)
	}
	if _, err := os.Stat(record); err == nil {
		t.Fatal("a dry run recorded the limit")
	}
	// A real run (claude cannot be found, so nothing replaces the test).
	cfg.ClaudeBin = "ccam-test-no-such-claude"
	cfg.Save(c.p.ConfigFile())
	if _, errOut, code := c.run("", "next", "--continue"); code == 0 || !strings.Contains(errOut, "cannot find the claude executable") {
		t.Fatalf("%d %s", code, errOut)
	}
	out, _, _ := c.run("", "usage", "--source", "cache")
	if !strings.Contains(out, "100%") {
		t.Fatalf("after continuing:\n%s", out)
	}
	if _, err := os.Stat(record); err != nil {
		t.Fatal(err)
	}
}

// Prompt caches are kept per organization, and every personal plan is one:
// the size line says the cache starts over when the organizations differ.
func TestContinueSaysWhenCachesAreNotShared(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40})
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin = stubClaude
	cfg.Save(c.p.ConfigFile())
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	limitConversation(t, conv, time.Now(), time.Now().Add(time.Hour))
	note := "prompt caches are not shared between organizations"
	logIn(t, c.p.ProfileDir("personal"), "u1", "home-1", "claude_max")
	logIn(t, c.p.ProfileDir("work"), "u2", "home-2", "claude_pro")
	_, errOut, _ := c.run("", "--dry-run", "next", "--continue")
	if !strings.Contains(errOut, "about 100k tokens") || !strings.Contains(errOut, note) {
		t.Fatalf("two personal plans: %s", errOut)
	}
	logIn(t, c.p.ProfileDir("work"), "u2", "home-1", "claude_max")
	if _, errOut, _ := c.run("", "--dry-run", "next", "--continue"); strings.Contains(errOut, note) {
		t.Fatalf("one organization: %s", errOut)
	}
}

// rename and remove carry the record along with the rest of the account's
// saved usage.
func TestRenameMovesTheLimitRecord(t *testing.T) {
	c := newTestCLI(t, map[string]int{"work": 10})
	dir := filepath.Join(c.p.State, "usage")
	usage.SaveLimit(dir, "work", "five_hour", time.Now().Add(time.Hour), time.Now())
	if _, errOut, code := c.run("", "rename", "work", "job"); code != 0 {
		t.Fatalf("%d %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(dir, "job.limit.json")); err != nil || fileThere(filepath.Join(dir, "work.limit.json")) {
		t.Fatalf("after rename: %v", err)
	}
}

func fileThere(path string) bool { _, err := os.Stat(path); return err == nil }

// A pick mark dated after now (the clock was set back) is not recent.
func TestPickIgnoresMarksFromTheFuture(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 50})
	mark := filepath.Join(c.p.State, "picks", "personal")
	os.MkdirAll(filepath.Dir(mark), 0o700)
	os.WriteFile(mark, nil, 0o600)
	later := time.Now().Add(time.Hour)
	os.Chtimes(mark, later, later)
	if out, _, _ := c.run("", "--dry-run", "pick"); strings.TrimSpace(out) != "personal" {
		t.Fatalf("picked %q", out)
	}
}
