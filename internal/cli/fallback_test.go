package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func (c *testCLI) env() (*env, *config.Config) {
	e := &env{paths: c.p, stdin: strings.NewReader(""), in: bufio.NewReader(strings.NewReader("")), stdout: os.Stdout, stderr: os.Stderr}
	cfg, err := e.load()
	if err != nil {
		c.t.Fatal(err)
	}
	return e, cfg
}

func TestFallbackTargetPrefersRoomThenUnknown(t *testing.T) {
	cases := []struct {
		name  string
		usage map[string]int
		tried []string
		want  string
	}{
		{"lowest 5-hour usage", map[string]int{"personal": 100, "work": 60, "side": 20}, []string{"personal"}, "side"},
		{"full accounts are skipped", map[string]int{"personal": 100, "work": 100, "side": 30}, []string{"personal"}, "side"},
		{"unknown usage comes after room", map[string]int{"personal": 100, "work": -1, "side": 90}, []string{"personal"}, "side"},
		{"unknown usage is still tried", map[string]int{"personal": 100, "work": -1, "side": 100}, []string{"personal"}, "work"},
		{"tried accounts are not", map[string]int{"personal": 100, "work": 10, "side": 100}, []string{"personal", "work"}, ""},
		{"logged-out accounts are not", map[string]int{"personal": 100}, []string{"personal"}, ""},
	}
	for _, tc := range cases {
		c := newTestCLI(t, tc.usage)
		e, cfg := c.env()
		tried := map[string]bool{}
		for _, n := range tc.tried {
			tried[n] = true
		}
		from, _ := cfg.Find(tc.tried[0])
		got := ""
		if p, _ := fallbackTarget(e, cfg, from, tried, []string{"-p", "hi"}, usage.Policy{}); p != nil {
			got = p.Name
		}
		if got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestFallbackTargetKeepsHistoryWhenResuming(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 100, "work": 50, "side": 10})
	e, cfg := c.env()
	cfg.Profiles[2].Share.History = config.Local // side does not share history
	from := &cfg.Profiles[0]
	tried := map[string]bool{"personal": true}
	if p, _ := fallbackTarget(e, cfg, from, tried, []string{"-p", "go on", "--continue"}, usage.Policy{}); p == nil || p.Name != "work" {
		t.Fatalf("--continue must stay in the history group: %v", p)
	}
	if p, _ := fallbackTarget(e, cfg, from, tried, []string{"-p", "new"}, usage.Policy{}); p == nil || p.Name != "side" {
		t.Fatalf("a new run may use any account: %v", p)
	}
}

func TestLimitedByMessageOrUsage(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 100, "work": 40})
	e, cfg := c.env()
	personal, work := &cfg.Profiles[0], &cfg.Profiles[1]
	p := usage.Policy{}
	for _, out := range []string{
		"You've hit your limit · resets 3pm",
		"You've hit your session limit · resets 3:45pm (Asia/Seoul)",
		"Some output first\nYou’ve hit your weekly limit · resets Oct 24, 9am (Asia/Seoul)\n",
		"You've hit your Opus limit · resets Oct 24, 9am (Asia/Seoul)",
		// --output-format json and stream-json
		`{"type":"result","subtype":"success","is_error":true,"result":"You've hit your session limit · resets 3:45pm"}`,
		`{"type":"assistant","message":{"model":"<synthetic>","content":[{"type":"text","text":"You're out of usage credits"}]}}`,
	} {
		if !limited(e, cfg, work, p, []byte(out), nil) {
			t.Errorf("message on stdout: %q", out)
		}
	}
	if !limited(e, cfg, work, p, nil, []byte("Claude AI usage limit reached|1760000000")) {
		t.Error("older message on stderr")
	}
	if !limited(e, cfg, personal, p, []byte("API Error: 500"), nil) {
		t.Error("full by usage")
	}
	for _, out := range []string{"API Error: 500", "API Error: Server is temporarily limiting requests (not your usage limit)",
		`{"type":"result","is_error":true,"result":"API Error: Request rejected (429) · Number of request tokens has exceeded your rate limit"}`,
		// What Claude replies or a tool prints is no message of Claude Code's.
		`{"type":"assistant","message":{"model":"claude-opus-5-5","content":[{"type":"text","text":"You've reached your goal: the tests pass."}]}}`,
		`{"type":"user","message":{"content":[{"type":"tool_result","content":"You've hit your quota of retries"}]}}`} {
		if limited(e, cfg, work, p, []byte(out), nil) {
			t.Errorf("another failure with room left is not a limit: %q", out)
		}
	}
}

// An account Claude Code reported full is not tried as one of unknown
// usage, even when nothing else about it is known.
func TestFallbackTargetSkipsARecordedLimit(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 100, "work": -1, "side": 100})
	e, cfg := c.env()
	usage.SaveLimit(filepath.Join(c.p.State, "usage"), "work", "seven_day", time.Now().Add(48*time.Hour), time.Now())
	if p, _ := fallbackTarget(e, cfg, &cfg.Profiles[0], map[string]bool{"personal": true}, []string{"-p", "hi"}, usage.Policy{}); p != nil {
		t.Fatalf("tried %s", p.Name)
	}
}

// With --model, that model's weekly window counts as a limit too, as it
// does for the choice of the next account.
func TestLimitedByTheModelWindow(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"limits":[{"kind":"session","percent":10},{"kind":"weekly_all","percent":20},`+
			`{"kind":"weekly_scoped","percent":100,"resets_at":"2030-01-01T00:00:00Z","scope":{"model":{"id":null,"display_name":"Opus"}}}]}`)
	}))
	t.Cleanup(srv.Close)
	orig := newFetcher
	newFetcher = func(dir string) *usage.Fetcher {
		f := usage.NewFetcher(dir)
		f.URL = srv.URL
		return f
	}
	t.Cleanup(func() { newFetcher = orig })
	e, cfg := c.env()
	personal := &cfg.Profiles[0]
	fail := []byte("API Error: 500")
	if !limited(e, cfg, personal, usage.Policy{Model: "opus"}, fail, nil) {
		t.Error("the Opus week is full")
	}
	if limited(e, cfg, personal, usage.Policy{}, fail, nil) || limited(e, cfg, personal, usage.Policy{Model: "sonnet"}, fail, nil) {
		t.Error("without --model opus, the Opus week does not count")
	}
}

func TestPrintMode(t *testing.T) {
	for args, want := range map[string]bool{
		"-p hi": true, "--print hi": true, "hi -p": true, "--resume": false, "-- -p": false, "": false,
	} {
		if got := printMode(strings.Fields(args)); got != want {
			t.Errorf("%q: %v", args, got)
		}
	}
}

// scriptClaude makes claude_bin a shell script that appends the account it
// runs as to the returned file and then runs body.
func scriptClaude(t *testing.T, c *testCLI, body string) (runs string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("a shell script stands in for claude")
	}
	dir := t.TempDir()
	runs = filepath.Join(dir, "runs")
	bin := filepath.Join(dir, "claude")
	script := "#!/bin/sh\necho \"${CLAUDE_CONFIG_DIR##*/}\" >> '" + runs + "'\n" + body + "\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := c.cfg()
	cfg.ClaudeBin = bin
	if err := cfg.Save(c.p.ConfigFile()); err != nil {
		t.Fatal(err)
	}
	return runs
}

func ranAs(runs string) string {
	data, _ := os.ReadFile(runs)
	return strings.Join(strings.Fields(string(data)), ",")
}

// --wait bounds the waiting of the whole run: when claude keeps saying
// "limit" while the usage shows room (or cannot be read), each wait used
// to start a new --wait and the command ran again and again for hours.
func TestFallbackWaitIsBounded(t *testing.T) {
	for name, pct := range map[string]int{"usage shows room": 50, "usage unknown": -1} {
		t.Run(name, func(t *testing.T) {
			c := newTestCLI(t, map[string]int{"personal": pct, "work": pct})
			runs := scriptClaude(t, c, `echo "You've hit your limit · resets 3pm"; exit 1`)
			start := time.Date(2026, 10, 20, 9, 0, 0, 0, time.UTC)
			now := fakeTime(t, start)
			slept := sleep
			sleep = func(d time.Duration) {
				if now.Sub(start) > 2*time.Hour {
					t.Fatalf("still waiting after %v; claude ran as %s", now.Sub(start), ranAs(runs))
				}
				slept(d)
			}
			e, _ := c.env()
			var errOut bytes.Buffer
			e.stderr = &errOut
			err := fallback(e, "personal", []string{"-p", "hi"}, "30m")
			var coded codedError
			if !errors.As(err, &coded) || coded.code != exitNoRoom {
				t.Fatalf("got %v\n%s", err, errOut.String())
			}
			if now.Sub(start) > 30*time.Minute {
				t.Fatalf("waited %v with --wait=30m", now.Sub(start))
			}
			if strings.Contains(errOut.String(), "work reached its usage limit; running the command again as work") {
				t.Fatalf("after a wait the limit is not said again:\n%s", errOut.String())
			}
		})
	}
}

// A run that a signal ended (timeout, a CI cancel, systemd stopping the
// job) is passed through, not run again as another account.
func TestFallbackDoesNotRerunASignalledRun(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 100, "work": 10})
	runs := scriptClaude(t, c, "exit 143")
	e, _ := c.env()
	err := fallback(e, "personal", []string{"-p", "hi"}, "")
	var code launcher.ExitCode
	if !errors.As(err, &code) || code != 143 || ranAs(runs) != "personal" {
		t.Fatalf("got %v, ran as %s", err, ranAs(runs))
	}
}

// The leading "--" any launch accepts: what follows it is what decides -p
// and the model.
func TestFallbackTakesTheSeparator(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 10, "work": 10})
	runs := scriptClaude(t, c, `echo "$@" > "${0%/*}/args"`)
	e, _ := c.env()
	if err := fallback(e, "personal", []string{"--", "-p", "x"}, ""); err != nil || ranAs(runs) != "personal" {
		t.Fatalf("got %v, ran as %s", err, ranAs(runs))
	}
	if b, _ := os.ReadFile(filepath.Join(filepath.Dir(runs), "args")); string(b) != "-p x\n" {
		t.Fatalf("claude got %q", b)
	}
	if err := fallback(e, "personal", []string{"--", "--", "-p", "x"}, ""); err == nil || !strings.Contains(err.Error(), "-p or --print") {
		t.Fatalf("a second -- makes -p a prompt: %v", err)
	}
	cfg := c.cfg()
	if p, err := (choiceFlags{}).policy(cfg, launcher.ClaudeArgs([]string{"--", "-p", "x", "--model", "opus"})); err != nil || p.Model != "opus" {
		t.Fatalf("model after the separator: %+v %v", p, err)
	}
}

// Piped input is read once for every attempt; like Claude Code, ccam waits
// at most a moment for it to start.
func TestReadInputWaitsOnlyForTheFirstByte(t *testing.T) {
	wait := 100 * time.Millisecond
	var errOut bytes.Buffer
	// A pipe that stays open and sends nothing.
	r, w := io.Pipe()
	defer w.Close()
	start := time.Now()
	in, err := readInput(r, &errOut, wait)
	if err != nil || in == nil || len(in) != 0 || time.Since(start) > 2*time.Second || !strings.Contains(errOut.String(), "no input on stdin") {
		t.Fatalf("idle pipe: %q %v %q", in, err, errOut.String())
	}
	// Input that starts in time is read to its end, however slowly it comes.
	errOut.Reset()
	r, w = io.Pipe()
	go func() {
		w.Write([]byte("first "))
		time.Sleep(3 * wait)
		w.Write([]byte("rest"))
		w.Close()
	}()
	if in, err := readInput(r, &errOut, wait); err != nil || string(in) != "first rest" || errOut.Len() != 0 {
		t.Fatalf("slow input: %q %v %q", in, err, errOut.String())
	}
	if in, err := readInput(strings.NewReader(""), &errOut, wait); err != nil || in == nil || len(in) != 0 || errOut.Len() != 0 {
		t.Fatalf("empty input: %q %v", in, err)
	}
	// Without a deadline (stream-json input) a late start is waited for.
	errOut.Reset()
	r, w = io.Pipe()
	go func() {
		time.Sleep(3 * wait)
		w.Write([]byte("late"))
		w.Close()
	}()
	if in, err := readInput(r, &errOut, 0); err != nil || string(in) != "late" || errOut.Len() != 0 {
		t.Fatalf("no deadline: %q %v %q", in, err, errOut.String())
	}
	for args, want := range map[string]bool{"-p --input-format stream-json": true, "-p --input-format=stream-json": true, "-p --input-format text": false, "-p -- --input-format=stream-json": false} {
		if streamInput(strings.Fields(args)) != want {
			t.Errorf("streamInput(%s) != %v", args, want)
		}
	}
	// Endless input stops at Claude Code's limit.
	if _, err := readInput(endless{}, &errOut, wait); err == nil || !strings.Contains(err.Error(), "10 MiB") {
		t.Fatalf("endless input: %v", err)
	}
}

type endless struct{}

func (endless) Read(b []byte) (int, error) { return len(b), nil }
