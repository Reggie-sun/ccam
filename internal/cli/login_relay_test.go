//go:build !windows

package cli

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// lineClaude makes every account of c log in with a claude whose `auth
// login` says it started (a file "started" in the account's folder),
// reads one line into "got" (with skip: nothing), and logs in.
func lineClaude(t *testing.T, c *testCLI, skip bool) {
	t.Helper()
	read := `IFS= read -r line; printf '%s' "$line" > "$CLAUDE_CONFIG_DIR/got"`
	if skip {
		read = ":"
	}
	fake := filepath.Join(t.TempDir(), "claude")
	script := "#!/bin/sh\n" +
		`: > "$CLAUDE_CONFIG_DIR/started"` + "\n" +
		read + "\n" +
		`printf '{"oauthAccount":{"emailAddress":"%s@example.com","accountUuid":"u-%s"}}' "$CCAM_PROFILE" "$CCAM_PROFILE" > "$CLAUDE_CONFIG_DIR/.claude.json"` + "\n"
	if err := os.WriteFile(fake, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := c.cfg()
	cfg.ClaudeBin = fake
	if err := cfg.Save(c.p.ConfigFile()); err != nil {
		t.Fatal(err)
	}
}

// appears waits up to d for path to exist.
func appears(path string, d time.Duration) bool {
	for deadline := time.Now().Add(d); time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
		if fileThere(path) {
			return true
		}
	}
	return fileThere(path)
}

// `ccam login --expired` logs in one account after another: the line typed
// for a later login goes to that login, not to the relay of an earlier one,
// which used to stay blocked on stdin after its claude ended.
func TestLoginExpiredGivesEachLoginItsLine(t *testing.T) {
	c := newTestCLI(t, nil)
	lineClaude(t, c, false)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	done := make(chan int, 1)
	go func() { _, _, code := c.runWith(r, "login", "--expired"); done <- code }()
	for _, name := range []string{"personal", "work", "side"} {
		dir := c.p.ProfileDir(name)
		if !appears(filepath.Join(dir, "started"), 10*time.Second) {
			t.Fatalf("the login of %s did not start", name)
		}
		if _, err := io.WriteString(w, "code-for-"+name+"\n"); err != nil {
			t.Fatal(err)
		}
		if !appears(filepath.Join(dir, "got"), 5*time.Second) {
			w.Close() // ends the login that waits for a line
			<-done
			t.Fatalf("the line typed for %s did not reach its login", name)
		}
		if line := readFile(t, filepath.Join(dir, "got")); line != "code-for-"+name {
			t.Fatalf("%s's login got %q", name, line)
		}
	}
	select {
	case code := <-done:
		if code != 0 {
			t.Fatalf("exit %d", code)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("login --expired did not end")
	}
}

// After a login, the next prompt (the dashboard's "Press Enter to return")
// gets the next line typed, and the one after it the line after that.
func TestLinesAfterALoginReachTheNextPrompts(t *testing.T) {
	c := newTestCLI(t, nil)
	lineClaude(t, c, true)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	e := &env{paths: c.p, stdin: r, in: bufio.NewReader(r), stdout: io.Discard, stderr: io.Discard}
	if _, err := loginProfile(e, "work"); err != nil {
		t.Fatal(err)
	}
	got := make(chan string, 1)
	go func() { line, _ := e.lines().readLine(); got <- line }()
	io.WriteString(w, "\ny\n")
	select {
	case line := <-got:
		if line != "\n" {
			t.Fatalf("the prompt got %q", line)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("the Enter typed after the login never reached the next prompt")
	}
	if !e.confirm("Proceed?") {
		t.Fatal("the answer after it was lost")
	}
}

// A read that its reader stopped waiting for brings its line to the next
// reader.
func TestLineSourceHandsOnAReadUnderWay(t *testing.T) {
	src := &lineSource{in: bufio.NewReader(strings.NewReader("a\nb\n"))}
	src.next(nil)
	for _, want := range []string{"a\n", "b\n", ""} {
		if line, _ := src.readLine(); line != want {
			t.Fatalf("got %q, want %q", line, want)
		}
	}
}
