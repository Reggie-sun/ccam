// Command fakeclaude stands in for claude in the end-to-end tests. It
// records its environment and arguments (NUL separated) in FAKE_CLAUDE_LOG
// and fakes `auth login` by writing an oauthAccount into .claude.json.
// The other FAKE_* variables make it misbehave on purpose.
package main

import (
	"cmp"
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
	if path := os.Getenv("FAKE_ENV_FILE"); path != "" {
		// NAME=value lines for the variables named in FAKE_ENV_NAMES.
		var b strings.Builder
		for _, name := range strings.Split(os.Getenv("FAKE_ENV_NAMES"), ",") {
			if v, ok := os.LookupEnv(name); ok {
				fmt.Fprintf(&b, "%s=%s\n", name, v)
			}
		}
		must(os.WriteFile(path, []byte(b.String()), 0o600))
	}
	if runs := os.Getenv("FAKE_RUNS"); runs != "" {
		f, err := os.OpenFile(runs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		must(err)
		fmt.Fprintln(f, profile)
		f.Close()
	}
	if from := os.Getenv("FAKE_SWITCH_FROM"); from != "" && from == profile {
		// What `ccam hook switch` leaves for /ccam:switch.
		cwd, _ := os.Getwd()
		ticket := fmt.Sprintf(`{"account":%q,"session_id":"s1","transcript_path":%q,"cwd":%q,"at":%q,"reason":"switch","target":%q,"relay":%q}`,
			from, os.Getenv("FAKE_SWITCH_TRANSCRIPT"), cwd, time.Now().Format(time.RFC3339Nano), os.Getenv("FAKE_SWITCH_TO"),
			cmp.Or(os.Getenv("FAKE_SWITCH_RELAY"), os.Getenv("CCAM_RELAY")))
		dir := filepath.Join(os.Getenv("FAKE_SWITCH_STATE"), "handoff")
		must(os.MkdirAll(dir, 0o700))
		must(os.WriteFile(filepath.Join(dir, from+".json"), []byte(ticket), 0o600))
	}
	if slices.Contains(args, "--input-format") && slices.Contains(args, "stream-json") {
		// The get_usage control request: answer from FAKE_USAGE_DIR/<account>.json.
		in, _ := io.ReadAll(os.Stdin)
		if strings.Contains(string(in), "get_usage") {
			data, err := os.ReadFile(filepath.Join(os.Getenv("FAKE_USAGE_DIR"), profile+".json"))
			if err != nil {
				data = []byte(`{"type":"control_response","response":{"subtype":"success","request_id":"ccam","response":{"rate_limits_available":false,"rate_limits":null}}}`)
			}
			fmt.Println(`{"type":"system","subtype":"init"}`)
			os.Stdout.Write(append(data, '\n'))
		}
		return
	}
	if len(args) >= 2 && args[0] == "agents" && args[1] == "--json" {
		// FAKE_AGENTS_DIR/<account>.json holds what the account lists.
		data, err := os.ReadFile(filepath.Join(os.Getenv("FAKE_AGENTS_DIR"), profile+".json"))
		if err != nil {
			data = []byte("[]")
		}
		os.Stdout.Write(data)
		return
	}
	if len(args) >= 2 && args[0] == "auth" && args[1] == "login" {
		email := os.Getenv("FAKE_EMAIL")
		if email == "" {
			email = "user@example.com"
		}
		// The account is named after the email unless FAKE_ACCOUNT says
		// otherwise; FAKE_LOGIN_ARGS records the arguments.
		account := os.Getenv("FAKE_ACCOUNT")
		if account == "" {
			account = "uuid-" + email
		}
		if path := os.Getenv("FAKE_LOGIN_ARGS"); path != "" {
			must(os.WriteFile(path, []byte(strings.Join(args, " ")), 0o600))
		}
		must(os.WriteFile(filepath.Join(dir, ".claude.json"), fmt.Appendf(nil,
			`{"numStartups":1,"oauthAccount":{"emailAddress":%q,"accountUuid":%q,"organizationUuid":"org-1","organizationName":"Acme","organizationType":"claude_team"}}`,
			email, account), 0o600))
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
