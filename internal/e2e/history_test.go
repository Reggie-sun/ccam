package e2e

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// holdFolder starts `ccam <account>` in dir as a session that lasts a few
// seconds, and returns once its claude runs, with a func that ends it.
func (h *harness) holdFolder(dir, account string) (stop func()) {
	h.t.Helper()
	log := filepath.Join(h.t.TempDir(), "holder.log")
	cmd := exec.Command(ccamBin, account)
	cmd.Dir = dir
	cmd.Env = append(append([]string{}, h.env...), "PWD="+dir, "FAKE_SLEEP=6", "FAKE_CLAUDE_LOG="+log)
	if err := cmd.Start(); err != nil {
		h.t.Fatal(err)
	}
	stop = sync.OnceFunc(func() {
		if runtime.GOOS == "windows" {
			// ccam waits for claude there instead of becoming it: end both.
			exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
		}
		cmd.Process.Kill()
		cmd.Wait()
	})
	h.t.Cleanup(stop)
	deadline := time.Now().Add(3 * time.Second)
	for !fileExists(log) && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if !fileExists(log) {
		h.t.Fatal("the session that holds the folder did not start")
	}
	return stop
}

// onTerminal runs the sh command line on a terminal of its own (ptyCommand)
// in dir, with input typed at the start, and returns what the terminal
// showed. It fails the test when the line has not finished after 20 s.
func (h *harness) onTerminal(dir, line, input string, extraEnv ...string) string {
	h.t.Helper()
	term := h.startTerminal(dir, line, extraEnv...)
	term.in.Write([]byte(input))
	out, err := term.wait(20 * time.Second)
	if err != nil {
		h.t.Fatal(err)
	}
	return out
}

// terminal is a command line running on a terminal of its own.
type terminal struct {
	cmd  *exec.Cmd
	in   io.WriteCloser // typed on the terminal
	out  *syncBuffer    // what the terminal shows
	done chan error
}

// startTerminal starts the sh command line on a terminal of its own in dir.
func (h *harness) startTerminal(dir, line string, extraEnv ...string) *terminal {
	h.t.Helper()
	cmd := ptyCommand(h.t, line)
	cmd.Dir = dir
	cmd.Env = append(append(append([]string{}, h.env...), "PWD="+dir), extraEnv...)
	in, err := cmd.StdinPipe()
	if err != nil {
		h.t.Fatal(err)
	}
	term := &terminal{cmd: cmd, in: in, out: &syncBuffer{}, done: make(chan error, 1)}
	cmd.Stdout, cmd.Stderr = term.out, term.out
	if err := cmd.Start(); err != nil {
		h.t.Fatal(err)
	}
	go func() { term.done <- cmd.Wait() }()
	h.t.Cleanup(func() { cmd.Process.Kill(); <-term.done; term.done <- nil })
	return term
}

// wait waits for the line to finish, at most d, and returns what the
// terminal showed.
func (term *terminal) wait(d time.Duration) (string, error) {
	select {
	case <-term.done:
		term.done <- nil
		return term.out.String(), nil
	case <-time.After(d):
		term.cmd.Process.Kill()
		return "", fmt.Errorf("%s did not finish:\n%s", term.cmd.Args, term.out.String())
	}
}

// syncBuffer is a bytes.Buffer for a writer and readers in other goroutines.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func sharedHistory(t *testing.T) (*harness, string) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "group", "set", "default", "--add", "history", "--yes")
	dir := filepath.Join(h.home, "app")
	os.MkdirAll(dir, 0o700)
	dir, _ = filepath.EvalSymlinks(dir)
	return h, dir
}

// A -p worker started beside a session that shares its history goes on:
// a new -p conversation cannot mix up the other one, and stdin of
// /dev/null, which Claude Code's Bash tool and background jobs give it, is
// a character device but no terminal anyone could answer on.
func TestHistoryWarningLetsWorkersGoOn(t *testing.T) {
	h, dir := sharedHistory(t)
	h.holdFolder(dir, "personal")
	run := func(argv ...string) result {
		t.Helper()
		cmd := exec.Command(ccamBin, argv...)
		cmd.Dir = dir
		cmd.Env = append(append([]string{}, h.env...), "PWD="+dir)
		var out, errb bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &errb // no Stdin: /dev/null
		err := cmd.Run()
		code := 0
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else if err != nil {
			t.Fatal(err)
		}
		return result{out.String(), errb.String(), code}
	}
	for _, argv := range [][]string{{"work", "-p", "task"}, {"--fallback", "work", "-p", "task"}} {
		r := run(argv...)
		if r.code != 0 || strings.Contains(r.stderr, "another Claude Code") {
			t.Fatalf("%v: %+v", argv, r)
		}
		if got, _ := h.called(); got != h.profileDir("work") {
			t.Fatalf("%v ran as %s", argv, got)
		}
	}
	// Continuing a conversation is warned about; nobody can answer, so it goes on.
	r := run("work", "-p", "--continue", "task")
	if r.code != 0 || !strings.Contains(r.stderr, "another Claude Code") || strings.Contains(r.stderr, "[y/N]") {
		t.Fatalf("%+v", r)
	}
	h.called()
}

// A job in the background of a shell with job control may not read the
// terminal: the terminal would stop it (SIGTTIN). It goes on after the
// warning instead, and a -p one is not even warned.
func TestHistoryQuestionSparesBackgroundJobs(t *testing.T) {
	ptyCommand(t, "true") // skips without a terminal helper
	h, dir := sharedHistory(t)
	h.holdFolder(dir, "personal")
	script := filepath.Join(h.home, "jobs.sh")
	os.WriteFile(script, []byte(`set -m
"$CCAM" work -p one >one.out 2>&1 &
a=$!
"$CCAM" work --continue >two.out 2>&1 &
b=$!
wait $a; echo $? >one.code
wait $b; echo $? >two.code
kill -9 $a $b 2>/dev/null
true
`), 0o600)
	shown := h.onTerminal(dir, "bash "+script, "", "CCAM="+ccamBin)
	for _, job := range []string{"one", "two"} {
		code, _ := os.ReadFile(filepath.Join(dir, job+".code"))
		out, _ := os.ReadFile(filepath.Join(dir, job+".out"))
		if strings.TrimSpace(string(code)) != "0" || strings.Contains(string(out), "[y/N]") {
			t.Fatalf("job %s ended with %q:\n%s\nterminal:\n%s", job, code, out, shown)
		}
	}
	if out, _ := os.ReadFile(filepath.Join(dir, "one.out")); strings.Contains(string(out), "another Claude Code") {
		t.Fatalf("a new -p conversation was warned about:\n%s", out)
	}
}

// Only a launch that opens a session can mix up a conversation: --version
// or `mcp list` beside a running session asks nothing. A session still
// asks on the terminal, and no means no.
func TestHistoryQuestionOnlyForSessions(t *testing.T) {
	ptyCommand(t, "true") // skips without a terminal helper
	h, dir := sharedHistory(t)
	h.holdFolder(dir, "personal")
	for _, line := range []string{ccamBin + " work --version", ccamBin + " work mcp list"} {
		if shown := h.onTerminal(dir, line, "n\n"); strings.Contains(shown, "Run anyway?") || strings.Contains(shown, "cancelled") {
			t.Fatalf("%s:\n%s", line, shown)
		}
	}
	if shown := h.onTerminal(dir, ccamBin+" work --continue", "n\n"); !strings.Contains(shown, "Run anyway? [y/N]") || !strings.Contains(shown, "ccam: cancelled") {
		t.Fatalf("a session on a terminal:\n%s", shown)
	}
}

// Claude Code keeps a folder's conversations under its real path, so the
// same folder reached through a symbolic link is the same folder: the
// warning comes either way round.
func TestHistoryLockSeesASymlinkedFolder(t *testing.T) {
	h, real := sharedHistory(t)
	link := filepath.Join(h.home, "link")
	if err := os.Symlink(real, link); err != nil {
		t.Skip("no symlinks:", err)
	}
	for _, c := range []struct{ held, run string }{{real, link}, {link, real}} {
		stop := h.holdFolder(c.held, "personal")
		r := h.runIn(c.run, "", []string{"PWD=" + c.run}, "ccam", "work")
		if r.code != 0 || !strings.Contains(r.stderr, "another Claude Code") {
			t.Fatalf("held in %s, run in %s: %+v", c.held, c.run, r)
		}
		stop()
	}
}
