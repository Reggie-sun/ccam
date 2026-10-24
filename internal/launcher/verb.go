package launcher

import "strings"

// Verbs of claudeVerb besides the ones claude takes.
const (
	verbRemoteControl = "remote-control" // also rc, and remote, sync, bridge first
)

// claudeVerb is the subcommand claude args run, as far as ccam needs it to
// place --mcp-config (Claude Code 2.1.308): verbRemoteControl; a
// background-session verb given first (attach, logs, stop, kill, respawn,
// rm), which Claude Code reads by name only there; "mcp <verb>"; or "" for
// anything else, a session among them. Options are skipped with their
// values; an option claude has that ccam does not know counts as one
// without a value, so its value ends the search as a prompt and ccam
// passes --mcp-config: where it does not belong, claude refuses it loudly
// rather than start a session without the group's servers.
func claudeVerb(args []string) string {
	if len(args) == 0 {
		return ""
	}
	switch args[0] {
	case "remote-control", "rc", "remote", "sync", "bridge":
		return verbRemoteControl // read before any option
	case "attach", "logs", "stop", "kill", "respawn", "rm":
		return args[0]
	}
	for _, a := range args {
		if a == "--" {
			break
		}
		if a == "-p" || a == "--print" {
			return "" // claude has no subcommands with --print
		}
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--":
			if i+1 < len(args) {
				return verbAt(args, i+1)
			}
			return ""
		case isFlag(a):
			name, _, inline := strings.Cut(a, "=")
			if !inline {
				i += takes(args, i, arity(name))
			}
		default:
			return verbAt(args, i)
		}
	}
	return ""
}

// verbAt is claudeVerb for the first argument that is no option, at i.
func verbAt(args []string, i int) string {
	switch args[i] {
	case "remote-control", "rc":
		return verbRemoteControl
	case "mcp":
		if i+1 < len(args) && !isFlag(args[i+1]) {
			return "mcp " + args[i+1]
		}
		return "mcp"
	}
	return ""
}

// takes is how many of the arguments after args[i], a flag with n values
// (arity), are its values.
func takes(args []string, i, n int) int {
	switch {
	case (n == oneValue || n == values) && i+1 < len(args):
		take := 1
		for n == values && i+1+take < len(args) && !isFlag(args[i+1+take]) {
			take++
		}
		return take
	case n == optValue && i+1 < len(args) && !isFlag(args[i+1]):
		return 1
	}
	return 0
}

// arity is how many values a claude option takes: carried's count, and
// for options a hop does not need, from Claude Code 2.1.308's --help, its
// Remote Control allow-list and some of its hidden options.
func arity(flag string) int {
	switch flag {
	case "--workload", "--plugin-dir-no-mcp", "--thinking", "--thinking-display", "--max-thinking-tokens",
		"--remote-control-session-name-prefix", "--output-format", "--input-format", "--permission-prompts",
		"--permission-prompt-tool", "--json-schema", "--max-budget-usd", "--environment", "--max-turns",
		"--teammate-mode":
		return oneValue
	case "--file", "--channels", "--dangerously-load-development-channels":
		return values
	case "--cloud", "--prompt-suggestions", "--remote":
		return optValue
	}
	_, _, n := carried(flag)
	return n
}
