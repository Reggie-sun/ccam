package cli

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
)

// hookInput is the part of a Claude Code hook's JSON that ccam reads.
type hookInput struct {
	SessionID   string `json:"session_id"`
	Transcript  string `json:"transcript_path"`
	Cwd         string `json:"cwd"`
	Error       string `json:"error"`
	LastMessage string `json:"last_assistant_message"`
	CommandArgs string `json:"command_args"`
}

// folder is where the session started. The hook's cwd follows a cd in
// Claude Code's shell; the conversation stays under the start folder.
func (in hookInput) folder() string { return cmp.Or(os.Getenv("CLAUDE_PROJECT_DIR"), in.Cwd) }

func hookCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "hook stop-failure|switch",
		Short: "Run by the ccam plugin's hooks inside Claude Code (reads the hook's JSON on stdin)",
		Long: "The ccam plugin for Claude Code runs these; you do not need to.\n\n" +
			"  stop-failure  a turn stopped at a usage limit: note the conversation and show a\n" +
			"                terminal notification, so `ccam next --continue` carries it on\n" +
			"  switch        /ccam:switch [account]: note the conversation and end the turn\n\n" +
			"Both do nothing in a session that ccam did not start.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"stop-failure", "switch"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var in hookInput
			_ = json.NewDecoder(io.LimitReader(e.stdin, 1<<20)).Decode(&in)
			account := os.Getenv("CCAM_PROFILE")
			switch args[0] {
			case "stop-failure":
				return hookStopFailure(e, in, account)
			case "switch":
				return hookSwitch(e, in, account)
			}
			return fmt.Errorf("unknown hook %q (stop-failure or switch)", args[0])
		},
	}
}

// hookOutput writes the hook's JSON answer.
func hookOutput(e *env, v map[string]any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	e.printf("%s\n", data)
	return nil
}

var limitNames = map[string]string{"session": "5-hour", "weekly": "weekly", "model": "model", "spend": "spend"}

// hookStopFailure notes a session that stopped at a usage limit. A short
// server throttle or an API key's rate limit stops a turn the same way, but
// is not a usage limit, and is left alone.
func hookStopFailure(e *env, in hookInput, account string) error {
	kind := handoff.LimitKind(in.LastMessage)
	if account == "" || in.Error != "rate_limit" || kind == "" {
		return nil
	}
	t := handoff.Ticket{Account: account, SessionID: in.SessionID, Transcript: in.Transcript, Cwd: in.folder(),
		At: time.Now(), Reason: handoff.Limit, Limit: kind}
	if err := handoff.Save(e.paths.State, t); err != nil {
		return err
	}
	msg := fmt.Sprintf("ccam: %s hit its %s limit", account, limitNames[kind])
	if r := handoff.Reset(in.LastMessage); r != "" {
		msg += " (resets " + r + ")"
	}
	msg += ". `ccam next --continue` carries on as another account."
	// Claude Code shows nothing else a StopFailure hook prints.
	return hookOutput(e, map[string]any{"terminalSequence": notification(msg)})
}

// hookSwitch notes the session for /ccam:switch and ends the turn before
// any request is made. With relay on, it also ends Claude Code, and the
// ccam that started it carries the conversation on.
func hookSwitch(e *env, in hookInput, account string) error {
	stop := func(msg string) error {
		return hookOutput(e, map[string]any{"continue": false, "stopReason": msg})
	}
	if account == "" {
		return stop("ccam: this session was not started by ccam, so ccam cannot move it. Start Claude Code with `ccam <account>`.")
	}
	target := ""
	if f := strings.Fields(in.CommandArgs); len(f) > 0 {
		target = f[0]
	}
	cfg, err := e.load()
	if err != nil {
		return stop("ccam: " + err.Error())
	}
	if target != "" {
		p, ok := cfg.Find(target)
		if !ok {
			return stop(fmt.Sprintf("ccam: no account %q (see `ccam list`).", target))
		}
		if target = p.Name; target == account {
			return stop("ccam: this session already runs as " + account + ".")
		}
	}
	relay := os.Getenv("CCAM_RELAY")
	t := handoff.Ticket{Account: account, SessionID: in.SessionID, Transcript: in.Transcript, Cwd: in.folder(),
		At: time.Now(), Reason: handoff.Switch, Target: target, Relay: relay}
	if err := handoff.Save(e.paths.State, t); err != nil {
		return stop("ccam: " + err.Error())
	}
	who := cmp.Or(target, "the next account")
	if relay != "" && runtime.GOOS != "windows" {
		if pid := os.Getenv("CLAUDE_PID"); pid != "" {
			// After the answer is read: Claude Code then exits, and the
			// relaying ccam starts the conversation as the other account.
			_ = exec.Command("sh", "-c", "sleep 0.3; kill -TERM "+launcher.ShellQuote(pid)).Start()
			return stop("ccam: moving this conversation to " + who + "…")
		}
	}
	if relay != "" {
		return stop("ccam: type /exit; ccam then carries this conversation on as " + who + ".")
	}
	return stop("ccam: type /exit, then run `ccam next --continue` to carry this conversation on as " + who + ".")
}

// notification is a desktop notification in the escape sequence the
// terminal understands: OSC 99 for kitty, OSC 777 for Ghostty, urxvt and
// Warp, OSC 9 for the rest (iTerm2, WezTerm, Windows Terminal …).
func notification(msg string) string {
	msg = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, msg)
	tp := os.Getenv("TERM_PROGRAM")
	switch {
	case os.Getenv("KITTY_WINDOW_ID") != "":
		return "\x1b]99;;" + msg + "\x1b\\"
	case tp == "ghostty" || tp == "WarpTerminal" || strings.HasPrefix(os.Getenv("TERM"), "rxvt"):
		return "\x1b]777;notify;ccam;" + msg + "\x07"
	}
	return "\x1b]9;" + msg + "\x07"
}
