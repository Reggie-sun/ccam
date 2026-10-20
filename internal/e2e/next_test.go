package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

func TestNextRotatesAndContinues(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "personal")
	h.called()
	time.Sleep(1100 * time.Millisecond) // last-used has 1 s resolution
	h.ok("ccam", "next", "--", "-p", "hi")
	dir, args := h.called()
	if dir != h.profileDir("work") || !reflect.DeepEqual(args, []string{"-p", "hi"}) {
		t.Fatalf("next: %s %q", dir, args)
	}
	// --continue needs a conversation in the folder.
	project := filepath.Join(h.home, "app")
	os.MkdirAll(project, 0o700)
	r := h.runIn(project, "", nil, "ccam", "next", "--continue")
	if r.code == 0 || !strings.Contains(r.stderr, "no account has a conversation in") {
		t.Fatalf("expected an error: %+v", r)
	}
	// work has one; personal does not share its history, so Claude Code
	// copies it over.
	real, _ := filepath.EvalSymlinks(project)
	conv := filepath.Join(h.profileDir("work"), "projects", transcript.Slug(real), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	os.WriteFile(conv, []byte(`{"type":"ai-title","aiTitle":"Fix the tests"}`+"\n"+
		`{"type":"assistant","message":{"usage":{"input_tokens":1000,"cache_read_input_tokens":181000,"output_tokens":400}}}`+"\n"), 0o600)
	r = h.runIn(project, "", nil, "ccam", "next", "--continue", "--", "-p", "go on")
	if r.code != 0 || !strings.Contains(r.stderr, `continuing "Fix the tests" (about 182k tokens); personal does not share work's history`) {
		t.Fatalf("%+v", r)
	}
	dir, args = h.called()
	// ccam passes the file with links resolved (/private/var on macOS, long
	// names for 8.3 ones on Windows).
	resolved, _ := filepath.EvalSymlinks(conv)
	if dir != h.profileDir("personal") || !reflect.DeepEqual(args, []string{"--resume", resolved, "--fork-session", "-p", "go on"}) {
		t.Fatalf("next --continue: %s %q", dir, args)
	}
	// Sharing history: that very conversation by its id (--continue could
	// open a newer one).
	h.ok("ccam", "group", "set", "default", "--add", "history", "--yes")
	time.Sleep(1100 * time.Millisecond)
	r = h.runIn(project, "", nil, "ccam", "next", "--continue", "--", "-p", "go on")
	dir, args = h.called()
	if r.code != 0 || dir != h.profileDir("work") || !reflect.DeepEqual(args, []string{"--resume", "s1", "-p", "go on"}) {
		t.Fatalf("next --continue: %s %q %+v", dir, args, r)
	}
}

// writeNote leaves what `ccam hook switch` writes for /ccam:switch.
func writeNote(t *testing.T, h *harness, account, conv, cwd, target string) {
	t.Helper()
	dir := filepath.Join(h.home, "state", "handoff")
	os.MkdirAll(dir, 0o700)
	note := fmt.Sprintf(`{"account":%q,"session_id":"s1","transcript_path":%q,"cwd":%q,"at":%q,"reason":"switch","target":%q}`,
		account, conv, cwd, time.Now().Format(time.RFC3339Nano), target)
	if err := os.WriteFile(filepath.Join(dir, account+".json"), []byte(note), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestNoteWithAGoneTargetOrAFailedLaunch(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	project := filepath.Join(h.home, "app")
	os.MkdirAll(project, 0o700)
	real, _ := filepath.EvalSymlinks(project)
	conv := filepath.Join(h.profileDir("personal"), "projects", transcript.Slug(real), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	os.WriteFile(conv, []byte("{}\n"), 0o600)
	note := filepath.Join(h.home, "state", "handoff", "personal.json")

	// The account /ccam:switch named was renamed since: choose as usual.
	writeNote(t, h, "personal", conv, real, "side")
	r := h.runIn(project, "", nil, "ccam", "next", "--continue")
	if dir, _ := h.called(); r.code != 0 || dir != h.profileDir("work") || !strings.Contains(r.stderr, "names side, which is no longer an account") {
		t.Fatalf("%s %+v", dir, r)
	}
	if _, err := os.Stat(note); err == nil {
		t.Fatal("a note that was carried on stays")
	}

	// A launch that fails keeps the note for the next try.
	writeNote(t, h, "personal", conv, real, "work")
	os.Rename(h.profileDir("work"), h.profileDir("work")+".away")
	if r := h.runIn(project, "", nil, "ccam", "next", "--continue"); r.code == 0 || !strings.Contains(r.stderr, "the folder of account work is missing") {
		t.Fatalf("%+v", r)
	}
	if _, err := os.Stat(note); err != nil {
		t.Fatal("the failed launch used the note up")
	}
}

func TestHistoryLockWarnsOnConcurrentUse(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "group", "set", "default", "--add", "history", "--yes")
	first := exec.Command(ccamBin, "personal")
	first.Env = append(append([]string{}, h.env...), "FAKE_SLEEP=2", "FAKE_CLAUDE_LOG="+filepath.Join(h.home, "first.log"))
	if err := first.Start(); err != nil {
		t.Fatal(err)
	}
	defer first.Wait()
	deadline := time.Now().Add(2 * time.Second)
	for !fileExists(filepath.Join(h.home, "first.log")) && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	r := h.run("", nil, "ccam", "work")
	if r.code != 0 || !strings.Contains(r.stderr, "another Claude Code") {
		t.Fatalf("expected warning, got %+v", r)
	}
	// A different folder is not affected.
	cmd := exec.Command(ccamBin, "work")
	cmd.Dir = t.TempDir()
	cmd.Env = h.env
	out, _ := cmd.CombinedOutput()
	if strings.Contains(string(out), "another Claude Code") {
		t.Fatalf("other folder warned: %s", out)
	}
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// A background job started inside the session must not keep the account
// looking "running" after claude itself exits.
func TestLeftoverChildDoesNotHoldSession(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "group", "set", "default", "--add", "history", "--yes")
	r := h.run("", []string{"FAKE_BG=3"}, "ccam", "personal")
	if r.code != 0 {
		t.Fatalf("%+v", r)
	}
	r = h.run("", nil, "ccam", "work")
	if strings.Contains(r.stderr, "another Claude Code") {
		t.Fatalf("finished session still counted as running: %s", r.stderr)
	}
}

// With relay on, ccam stays beside Claude Code, and after /ccam:switch it
// carries the conversation on as the other account by itself.
func TestRelayCarriesASwitchOn(t *testing.T) {
	script, err := exec.LookPath("script")
	if err != nil || runtime.GOOS != "linux" {
		t.Skip("needs util-linux script for a terminal")
	}
	h := newHarness(t)
	h.addAccounts()
	cfgPath := filepath.Join(h.home, "config", "config.toml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Relay = true
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatal(err)
	}
	project := filepath.Join(h.home, "app")
	real, _ := filepath.EvalSymlinks(h.home)
	conv := filepath.Join(h.profileDir("personal"), "projects", transcript.Slug(filepath.Join(real, "app")), "s1.jsonl")
	os.MkdirAll(filepath.Dir(conv), 0o700)
	os.MkdirAll(project, 0o700)
	os.WriteFile(conv, []byte("{}\n"), 0o600)
	envFile := filepath.Join(h.home, "env")

	cmd := exec.Command(script, "-qec", ccamBin+" personal", "/dev/null")
	cmd.Dir = project
	cmd.Env = append(append([]string{}, h.env...), "FAKE_SWITCH_FROM=personal", "FAKE_SWITCH_TO=work",
		"FAKE_SWITCH_TRANSCRIPT="+conv, "FAKE_SWITCH_STATE="+filepath.Join(h.home, "state"),
		"FAKE_ENV_FILE="+envFile, "FAKE_ENV_NAMES=CCAM_RELAY")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	dir, args := h.called()
	resolved, _ := filepath.EvalSymlinks(conv)
	if dir != h.profileDir("work") || !reflect.DeepEqual(args, []string{"--resume", resolved, "--fork-session"}) {
		t.Fatalf("the relay started %s %q\n%s", dir, args, out)
	}
	if b, _ := os.ReadFile(envFile); !regexp.MustCompile(`CCAM_RELAY=\d+`).Match(b) {
		t.Fatalf("the session does not know it is relayed: %q", b)
	}

	// A note another session wrote (another relay, or none) is not this
	// session's /ccam:switch: the relay ends with claude.
	cmd = exec.Command(script, "-qec", ccamBin+" personal", "/dev/null")
	cmd.Dir = project
	cmd.Env = append(append([]string{}, h.env...), "FAKE_SWITCH_FROM=personal", "FAKE_SWITCH_TO=work", "FAKE_SWITCH_RELAY=1",
		"FAKE_SWITCH_TRANSCRIPT="+conv, "FAKE_SWITCH_STATE="+filepath.Join(h.home, "state"))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if dir, _ := h.called(); dir != h.profileDir("personal") {
		t.Fatalf("the relay carried on another session's note as %s", dir)
	}
}
