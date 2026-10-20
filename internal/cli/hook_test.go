package cli

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

func hookJSON(t *testing.T, v map[string]string) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// A usage limit leaves a note and a notification; next --continue then
// carries that very conversation on as another account, never the one
// that hit the limit.
func TestLimitHookAndNextContinue(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 40, "side": 30})
	cfg, _ := config.Load(c.p.ConfigFile())
	cfg.ClaudeBin, _ = exec.LookPath("true")
	cfg.Save(c.p.ConfigFile())
	cwd, _ := os.Getwd()
	conv := filepath.Join(transcript.ProjectDir(c.p.ProfileDir("personal"), cwd), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	os.WriteFile(conv, []byte("{}\n"), 0o600)

	t.Setenv("CCAM_PROFILE", "personal")
	t.Setenv("KITTY_WINDOW_ID", "")
	t.Setenv("TERM_PROGRAM", "iTerm.app")
	// A short throttle is not a usage limit.
	throttle := hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": "API Error: Server is temporarily limiting requests (not your usage limit)", "transcript_path": conv, "cwd": cwd})
	if out, _, code := c.run(throttle, "hook", "stop-failure"); code != 0 || out != "" {
		t.Fatalf("%d %q", code, out)
	}
	limit := hookJSON(t, map[string]string{"error": "rate_limit", "last_assistant_message": "You've hit your session limit · resets 3:45pm",
		"session_id": "s1", "transcript_path": conv, "cwd": cwd})
	out, _, code := c.run(limit, "hook", "stop-failure")
	var got map[string]string
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil || !strings.HasPrefix(got["terminalSequence"], "\x1b]9;ccam: personal hit its 5-hour limit (resets 3:45pm).") {
		t.Fatalf("%d %q", code, out)
	}
	// personal has the most room on paper, but just hit its limit.
	_, errOut, code := c.run("", "--dry-run", "next", "--auto", "--continue")
	if code != 0 || !strings.Contains(errOut, "→ personal → side") || !strings.Contains(errOut, "copies it into side") {
		t.Fatalf("%d %s", code, errOut)
	}
	// A dry run leaves the note for the real run.
	if _, ok := handoff.ForFolder(c.p.State, cwd, time.Now()); !ok {
		t.Fatal("the note is gone after a dry run")
	}
}

func TestSwitchHook(t *testing.T) {
	c := newTestCLI(t, nil)
	cwd, _ := os.Getwd()
	conv := filepath.Join(t.TempDir(), "s2.jsonl")
	os.WriteFile(conv, []byte("{}\n"), 0o600)
	in := func(args string) string {
		return hookJSON(t, map[string]string{"session_id": "s2", "transcript_path": conv, "cwd": cwd, "command_args": args})
	}
	stopped := func(out string) string {
		var got struct {
			Continue   *bool  `json:"continue"`
			StopReason string `json:"stopReason"`
		}
		if json.Unmarshal([]byte(out), &got) != nil || got.Continue == nil || *got.Continue {
			t.Fatalf("not stopped: %q", out)
		}
		return got.StopReason
	}
	t.Setenv("CCAM_PROFILE", "")
	if out, _, _ := c.run(in(""), "hook", "switch"); !strings.Contains(stopped(out), "not started by ccam") {
		t.Fatal(out)
	}
	t.Setenv("CCAM_PROFILE", "work")
	t.Setenv("CCAM_RELAY", "")
	if out, _, _ := c.run(in("nobody"), "hook", "switch"); !strings.Contains(stopped(out), `no account "nobody"`) {
		t.Fatal(out)
	}
	out, _, _ := c.run(in("3"), "hook", "switch")
	if r := stopped(out); !strings.Contains(r, "/exit") || !strings.Contains(r, "as side") {
		t.Fatal(r)
	}
	tk, ok := handoff.ForFolder(c.p.State, cwd, time.Now())
	if !ok || tk.Account != "work" || tk.Target != "side" || tk.Reason != handoff.Switch {
		t.Fatalf("%+v", tk)
	}
}
