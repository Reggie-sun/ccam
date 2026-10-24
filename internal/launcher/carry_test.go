package launcher

import (
	"reflect"
	"strings"
	"testing"
)

func TestCarryArgs(t *testing.T) {
	split := func(s string) []string { return strings.Fields(s) }
	for _, c := range []struct{ in, out, dropped string }{
		// Resuming restores the model and the prompt is in the conversation.
		{"fix it --model opus --permission-mode plan --add-dir /tmp", "--permission-mode plan --add-dir /tmp", ""},
		{"--model opus --verbose fix it", "--verbose", ""},
		// Variadic flags take values up to the next flag, as Claude Code reads
		// them (so a prompt right after one was one of its values); = takes no more.
		{"--add-dir /a /b --verbose --allowedTools Bash Edit", "--add-dir /a /b --verbose --allowedTools Bash Edit", ""},
		{"--add-dir /a fix it", "--add-dir /a fix it", ""},
		{"--add-dir=/a prompt --mcp-config=x.json", "--add-dir=/a --mcp-config=x.json", ""},
		// The first value is taken whatever it looks like, as the one of a flag with one value.
		{"--add-dir -odd --verbose", "--add-dir -odd --verbose", ""},
		// A value that starts with - is still the value of a flag that takes one.
		{"--append-system-prompt -terse- --effort high", "--append-system-prompt -terse- --effort high", ""},
		// "--" ends the flags: what follows is a prompt.
		{"--ide -- --add-dir /x", "--ide", ""},
		// What picks the conversation goes quietly; what cannot be done again does not.
		{"-c --resume abc --fork-session --session-id 1234 -n name --agent rev --from-pr 12 --teleport", "", ""},
		{"--plugin-url https://example.com/p.zip --plugin-dir ./p", "--plugin-url https://example.com/p.zip --plugin-dir ./p", ""},
		{"--resume --verbose", "--verbose", ""},
		{"-w feature --tmux --debug api,hooks --debug-file /tmp/d.log --bg --remote-control phone", "", "-w --tmux --debug --debug-file --bg --remote-control"},
		// An unknown flag goes with a word; a value after it goes as a prompt would.
		{"--shiny value --chrome", "--chrome", "--shiny"},
		{"--shiny=1 --shiny=2", "", "--shiny"},
		// Flags a hop drops take their values as Claude Code reads them (arity).
		{"--output-format -x --file a b --verbose", "--verbose", "--output-format --file"},
		{"--dangerously-skip-permissions --settings {\"a\":1} --plugin-dir p1 --plugin-dir p2", "--dangerously-skip-permissions --settings {\"a\":1} --plugin-dir p1 --plugin-dir p2", ""},
	} {
		out, dropped := CarryArgs(split(c.in))
		if strings.Join(out, " ") != c.out || strings.Join(dropped, " ") != c.dropped {
			t.Errorf("%s:\n  carried %q, want %q\n  dropped %q, want %q", c.in, out, c.out, dropped, c.dropped)
		}
	}
	// A second hop carries what the first one passed.
	first, _ := WithState(must(CarryArgs(split("--model opus --add-dir /tmp --permission-mode plan --effort low go on"))), "acceptEdits", "high")
	again, _ := WithState(must(CarryArgs(append([]string{"--resume", "s1"}, first...))), "acceptEdits", "high")
	if !reflect.DeepEqual(first, again) {
		t.Fatalf("first %q, then %q", first, again)
	}
}

func must(out, _ []string) []string { return out }

func TestWithState(t *testing.T) {
	split := func(s string) []string { return strings.Fields(s) }
	for _, c := range []struct{ in, mode, effort, out, dropped string }{
		// Nothing known: the args as they were.
		{"--permission-mode plan --effort high --dangerously-skip-permissions", "", "", "--permission-mode plan --effort high --dangerously-skip-permissions", ""},
		// The conversation's mode and the session's effort win.
		{"--permission-mode plan --add-dir /tmp --effort low", "acceptEdits", "high", "--add-dir /tmp --permission-mode acceptEdits --effort high", ""},
		{"--permission-mode=plan --effort=low", "default", "max", "--permission-mode default --effort max", ""},
		// Bypass stays allowed when it was, without forcing it on.
		{"--dangerously-skip-permissions --ide", "default", "", "--ide --permission-mode default --allow-dangerously-skip-permissions", ""},
		{"--allow-dangerously-skip-permissions", "plan", "", "--permission-mode plan --allow-dangerously-skip-permissions", ""},
		{"", "bypassPermissions", "", "--permission-mode bypassPermissions --allow-dangerously-skip-permissions", ""},
		// A session that started in bypass may turn it on again.
		{"--permission-mode bypassPermissions", "default", "", "--permission-mode default --allow-dangerously-skip-permissions", ""},
		{"--permission-mode=bypassPermissions --ide", "plan", "", "--ide --permission-mode plan --allow-dangerously-skip-permissions", ""},
		// Claude Code stops on a mode it does not know, so it is not passed.
		{"--permission-mode plan", "turbo", "", "--permission-mode plan", "--permission-mode turbo"},
	} {
		in := split(c.in)
		out, dropped := WithState(in, c.mode, c.effort)
		if strings.Join(out, " ") != c.out || strings.Join(dropped, " ") != c.dropped {
			t.Errorf("%q %q %q:\n  %q, want %q\n  dropped %q, want %q", c.in, c.mode, c.effort, out, c.out, dropped, c.dropped)
		}
		if strings.Join(in, " ") != c.in {
			t.Errorf("%q: the args given changed to %q", c.in, in)
		}
	}
}
