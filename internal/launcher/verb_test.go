package launcher

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
)

// How Claude Code 2.1.308 reads the start of its args: Remote Control's
// names first before any option, background-session verbs by name only
// first, commander's root options with their values, no subcommands with
// --print (the dispatch and the allow-list read from the binary, the cases
// as the verifiers ran them against 2.1.306).
func TestClaudeVerb(t *testing.T) {
	for _, c := range []struct {
		args string
		want string
	}{
		{"", ""},
		{"-c", ""},
		{"fix the bug", ""},
		{"rc", verbRemoteControl},
		{"remote-control --name x", verbRemoteControl},
		{"remote", verbRemoteControl},
		{"sync", verbRemoteControl},
		{"bridge", verbRemoteControl},
		{"--verbose remote", ""}, // remote, sync and bridge only first
		{"--verbose rc", verbRemoteControl},
		{"--model opus remote-control", verbRemoteControl},
		{"--model=opus rc", verbRemoteControl},
		{"-n job rc", verbRemoteControl},
		{"--name job rc", verbRemoteControl},
		{"--debug-file x.log remote-control", verbRemoteControl},
		{"--plugin-dir /p rc", verbRemoteControl},
		{"--effort high remote-control", verbRemoteControl},
		{"--workload w rc", verbRemoteControl},
		{"--max-turns 3 rc", verbRemoteControl}, // hidden options of 2.1.308
		{"--teammate-mode auto rc", verbRemoteControl},
		{"--remote rc", ""}, // its optional value
		{"--channels a rc", ""},
		{"-- rc", verbRemoteControl},
		{"--verbose -- rc", verbRemoteControl},
		{"--add-dir /tmp rc", ""}, // rc is a second directory
		{"--add-dir=/tmp rc", verbRemoteControl},
		{"--betas a rc", ""},
		{"-d rc", ""}, // rc is the debug filter
		{"-d remote-control --name x", ""},
		{"-d --verbose rc", verbRemoteControl},
		{"--debug=api rc", verbRemoteControl},
		{"--resume rc", ""},
		{"-w rc", ""},
		{"-p rc", ""}, // no subcommands with --print
		{"--verbose rc -p", ""},
		{"rc -p", verbRemoteControl},
		{"--remote-control", ""}, // a session with Remote Control on
		{"logs myjob", "logs"},
		{"attach myjob", "attach"},
		{"stop 0123abcd", "stop"},
		{"kill x", "kill"},
		{"respawn --all", "respawn"},
		{"rm x", "rm"},
		{"--verbose logs myjob", ""}, // behind an option, a prompt
		{"mcp list", "mcp list"},
		{"--verbose mcp get fs", "mcp get"},
		{"mcp remove fs", "mcp remove"},
		{"mcp add fs -- npx x", "mcp add"},
		{"mcp", "mcp"},
		{"-p mcp list", ""},
		{"auth status", ""},
	} {
		if got := claudeVerb(strings.Fields(c.args)); got != c.want {
			t.Errorf("%q: got %q, want %q", c.args, got, c.want)
		}
	}
}

// --mcp-config goes in front of sessions and subcommands that take it, not
// in front of Remote Control (refused), background-session verbs (their
// names would not be read) or `mcp list|get|remove` (ignored there).
func TestMCPConfigOnlyWhereClaudeReadsIt(t *testing.T) {
	p, cfg := setup(t)
	cfg.Profiles[1].Share = config.Share{MCP: "default"}
	os.MkdirAll(p.GroupDir("default"), 0o700)
	os.WriteFile(filepath.Join(p.GroupDir("default"), "mcp.json"), []byte(`{"mcpServers":{"fs":{"command":"fs-mcp"}}}`), 0o600)
	flag := "--mcp-config=" + filepath.Join(p.GroupDir("default"), "mcp.json")
	for _, c := range []struct {
		args string
		flag bool
	}{
		{"-c", true},
		{"fix the bug", true},
		{"-d rc", true},
		{"--add-dir /tmp rc", true},
		{"mcp add fs -- x", true},
		{"auth status", true},
		{"--verbose logs myjob", true},
		{"rc", false},
		{"--model opus remote-control", false},
		{"-n job rc", false},
		{"logs myjob", false},
		{"attach myjob", false},
		{"mcp list", false},
		{"mcp get fs", false},
		{"mcp remove fs", false},
	} {
		pl, err := Ready(p, cfg, "work", strings.Fields(c.args), true)
		if err != nil {
			t.Fatal(err)
		}
		if got := slices.Contains(pl.Argv, flag); got != c.flag {
			t.Errorf("%q: flag %v, want %v (%q)", c.args, got, c.flag, pl.Argv)
		}
	}
}
