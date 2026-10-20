// Command fakeclaude stands in for claude in the end-to-end tests. It
// records its environment and arguments (NUL separated) in FAKE_CLAUDE_LOG
// and fakes `auth login` by writing an oauthAccount into .claude.json.
// The other FAKE_* variables make it misbehave on purpose.
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	if s := os.Getenv("FAKE_ONLY_SLEEP"); s != "" {
		sleep(s)
		return
	}
	args := os.Args[1:]
	dir := os.Getenv("CLAUDE_CONFIG_DIR")
	profile := os.Getenv("CCAM_PROFILE")
	var log strings.Builder
	fmt.Fprintf(&log, "CLAUDE_CONFIG_DIR=%s\x00CCAM_PROFILE=%s\x00", dir, profile)
	for _, a := range args {
		log.WriteString(a + "\x00")
	}
	path := os.Getenv("FAKE_CLAUDE_LOG")
	if path == "" {
		fmt.Fprintln(os.Stderr, "FAKE_CLAUDE_LOG not set")
		os.Exit(2)
	}
	must(os.WriteFile(path, []byte(log.String()), 0o600))
	if runs := os.Getenv("FAKE_RUNS"); runs != "" {
		f, err := os.OpenFile(runs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		must(err)
		fmt.Fprintln(f, profile)
		f.Close()
	}
	if len(args) >= 2 && args[0] == "auth" && args[1] == "login" {
		email := os.Getenv("FAKE_EMAIL")
		if email == "" {
			email = "user@example.com"
		}
		must(os.WriteFile(filepath.Join(dir, ".claude.json"), fmt.Appendf(nil, `{"numStartups":1,"oauthAccount":{"emailAddress":%q}}`, email), 0o600))
		must(os.WriteFile(filepath.Join(dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"fake","expiresAt":4102444800000}}`), 0o600))
	}
	if len(args) >= 2 && args[0] == "auth" && args[1] == "logout" {
		must(os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(`{}`), 0o600))
	}
	if probe := os.Getenv("FAKE_SESSION_PROBE"); probe != "" {
		// With exec the marker carries this process's PID; where ccam waits
		// for claude instead (Windows), it carries ccam's, the parent's.
		held := "free"
		for _, pid := range []int{os.Getpid(), os.Getppid()} {
			if _, err := os.Stat(filepath.Join(probe, strconv.Itoa(pid))); err == nil {
				held = "held"
			}
		}
		must(os.WriteFile(probe+".result", []byte(held+"\n"), 0o600))
	}
	if path := os.Getenv("FAKE_CLAUDE_STDIN"); path != "" {
		data, err := io.ReadAll(os.Stdin)
		must(err)
		must(os.WriteFile(path, data, 0o600))
	}
	if s := os.Getenv("FAKE_SLEEP"); s != "" {
		sleep(s)
	}
	if s := os.Getenv("FAKE_BG"); s != "" {
		// A child that outlives claude, like a Bash-tool background job.
		self, err := os.Executable()
		must(err)
		bg := exec.Command(self)
		bg.Env = append(os.Environ(), "FAKE_ONLY_SLEEP="+s)
		must(bg.Start())
	}
	if slices.Contains(strings.Split(os.Getenv("FAKE_LIMITED"), ","), profile) && profile != "" {
		fmt.Println("You've hit your limit · resets 3pm")
		os.Exit(1)
	}
	if out := os.Getenv("FAKE_STDOUT"); out != "" {
		fmt.Println(out + " (" + profile + ")")
	}
	code, _ := strconv.Atoi(os.Getenv("FAKE_EXIT"))
	os.Exit(code)
}

func sleep(s string) {
	n, _ := strconv.ParseFloat(s, 64)
	time.Sleep(time.Duration(n * float64(time.Second)))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "fakeclaude:", err)
		os.Exit(2)
	}
}
