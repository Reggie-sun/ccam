package handoff

import (
	"github.com/jhyunwoo/ccam/internal/transcript"

	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLimitKind(t *testing.T) {
	cases := map[string]string{
		"You've hit your session limit · resets 3:45pm":                                             "session",
		"You’ve hit your weekly limit · resets Oct 23, 9am":                                         "weekly",
		"You've hit your Opus limit · resets Oct 23":                                                "model",
		"You've hit your Fable limit":                                                               "model",
		"You've hit your monthly spend limit":                                                       "spend",
		"API Error: Server is temporarily limiting requests (not your usage limit)":                 "",
		"API Error: Request rejected (429) · Number of request tokens has exceeded your rate limit": "",
	}
	for in, want := range cases {
		if got := LimitKind(in); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
	if got := Reset("You've hit your session limit · resets 3:45pm"); got != "3:45pm" {
		t.Fatal(got)
	}
	if Reset("no reset here") != "" {
		t.Fatal("reset")
	}
}

func TestForFolder(t *testing.T) {
	state, cwd := t.TempDir(), t.TempDir()
	conv := filepath.Join(t.TempDir(), "s.jsonl")
	os.WriteFile(conv, nil, 0o600)
	now := time.Now()
	for _, tk := range []Ticket{
		{Account: "a", Cwd: cwd, Transcript: conv, At: now.Add(-time.Hour), Reason: Limit},
		{Account: "b", Cwd: cwd, Transcript: conv, At: now.Add(-time.Minute), Reason: Switch, Target: "c"},
		{Account: "c", Cwd: t.TempDir(), Transcript: conv, At: now, Reason: Limit},              // another folder
		{Account: "d", Cwd: cwd, Transcript: conv + "x", At: now, Reason: Limit},                // gone
		{Account: "e", Cwd: cwd, Transcript: conv, At: now.Add(-13 * time.Hour), Reason: Limit}, // too old
	} {
		if err := Save(state, tk); err != nil {
			t.Fatal(err)
		}
	}
	got, ok := ForFolder(state, cwd, now)
	if !ok || got.Account != "b" || got.Target != "c" {
		t.Fatalf("%+v", got)
	}
	Remove(state, "b")
	if got, _ := ForFolder(state, cwd, now); got.Account != "a" {
		t.Fatalf("%+v", got)
	}
}

// The hook's cwd follows a cd in Claude Code's shell; the conversation
// stays under the folder the session started in.
func TestNoteFromASubfolderBelongsToItsStartFolder(t *testing.T) {
	cwd := t.TempDir()
	real, _ := filepath.EvalSymlinks(cwd)
	conv := filepath.Join(t.TempDir(), "projects", transcript.Slug(real), "s.jsonl")
	tk := Ticket{Cwd: filepath.Join(cwd, "sub"), Transcript: conv}
	if !tk.In(cwd) || tk.In(filepath.Join(cwd, "sub", "other")) {
		t.Fatal("matched by the conversation's folder")
	}
}
