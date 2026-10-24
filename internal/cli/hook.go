package cli

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// hookInput is the part of a Claude Code hook's JSON that ccam reads.
type hookInput struct {
	SessionID  string `json:"session_id"`
	Transcript string `json:"transcript_path"`
	Cwd        string `json:"cwd"`
	AgentID    string `json:"agent_id"` // set when a subagent's loop stopped
	// Effort is the session's effort level, which a resumed conversation
	// does not get back on its own (absent for models without one).
	Effort struct {
		Level string `json:"level"`
	} `json:"effort"`
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
			"  stop-failure  a turn stopped at a usage limit: note the conversation, count the\n" +
			"                account as full until the limit resets, and show a terminal\n" +
			"                notification, so `ccam next --continue` carries it on\n" +
			"  switch        /ccam:switch [account]: note the conversation and end the turn\n\n" +
			"In a session that ccam did not start, stop-failure does nothing and switch says why it cannot move it.\n" +
			"stop-failure also does nothing for a subagent or a session nobody sits at (-p, background).",
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

// limitNames name the limits of handoff.LimitKind, windowNames the windows
// of Claude Code's limit row, for the notification.
var (
	limitNames  = map[string]string{"session": "5-hour", "weekly": "weekly", "usage": "usage", "model": "model", "spend": "spend", "org": "organization's usage"}
	windowNames = map[string]string{"five_hour": "5-hour", "seven_day": "weekly", "seven_day_opus": "Opus", "seven_day_sonnet": "Sonnet", "seven_day_overage_included": "Fable"}
)

// limitPoll is how long stop-failure waits for the refused request to show
// in the conversation: Claude Code writes it out every 100 ms and does not
// wait for the hook, which may start first. Tests change it.
var limitPoll = 1500 * time.Millisecond

// hookStopFailure notes a session that stopped at a usage limit, and when
// Claude Code wrote the limit into the conversation, records the account as
// full until the limit resets (usage.SaveLimit), so that no automatic
// choice takes it again before then. A short server throttle or an API
// key's rate limit stops a turn the same way, but is not a usage limit,
// and is left alone; so are a subagent's loop, which the main thread
// reports too, and sessions nobody sits at (-p, background).
func hookStopFailure(e *env, in hookInput, account string) error {
	if in.AgentID != "" {
		return nil
	}
	if account == "" || in.Error != "rate_limit" || unattended() {
		return nil
	}
	start := time.Now()
	var lim transcript.Limit
	note := func(kind string) error {
		return handoff.Save(e.paths.State, handoff.Ticket{Account: account, SessionID: in.SessionID, Transcript: in.Transcript,
			Cwd: in.folder(), At: time.Now(), Reason: handoff.Limit, Limit: kind, Effort: cmp.Or(in.Effort.Level, os.Getenv("CLAUDE_EFFORT")),
			Window: lim.Window, ResetsAt: lim.ResetsAt, LimitAt: lim.At})
	}
	// The note comes first when the message says enough: Claude Code may
	// close before the limit shows in the conversation, and end the hook.
	kind := handoff.LimitKind(in.LastMessage)
	if kind != "" {
		if err := note(kind); err != nil {
			return err
		}
	}
	lim, row := waitLimit(in.Transcript, start, in.LastMessage)
	switch {
	case row:
		// Again, now with the window and when it resets.
		kind = cmp.Or(kind, "usage")
		if err := note(kind); err != nil {
			return err
		}
	case kind == "":
		return nil
	}
	// An organization's stop (out of usage, a group at $0, an allocation
	// turned off) is not recorded: its admins can lift it at any time.
	own := row && handoff.Classify(in.LastMessage) != handoff.Org
	if own && usage.LimitWindow(lim.Window) {
		_ = usage.SaveLimit(usageStateDir(e), account, lim.Window, lim.ResetsAt, lim.At)
	}
	name, reset := limitNames[kind], handoff.Reset(in.LastMessage)
	if n, ok := windowNames[lim.Window]; ok && own {
		name = n
		if !lim.ResetsAt.IsZero() {
			reset = handoff.FormatReset(lim.ResetsAt, time.Now())
		}
	}
	msg := fmt.Sprintf("ccam: %s hit its %s limit", account, name)
	if reset != "" {
		msg += " (resets " + reset + ")"
	}
	msg += ". `ccam next --continue` carries on as another account."
	// Claude Code shows nothing else a StopFailure hook prints.
	return hookOutput(e, map[string]any{"terminalSequence": notification(msg)})
}

// unattended reports whether nobody sits at the session: Claude Code sets
// CLAUDE_CODE_SESSION_ATTENDED=0 for -p, background and daemon sessions,
// and an "sdk-" entrypoint for -p ones (also before it had that variable).
func unattended() bool {
	return os.Getenv("CLAUDE_CODE_SESSION_ATTENDED") == "0" || strings.HasPrefix(os.Getenv("CLAUDE_CODE_ENTRYPOINT"), "sdk-")
}

// waitLimit reads the usage limit Claude Code wrote into the conversation
// at path for the request it refused with msg, giving it limitPoll to show
// up. A limit with another message is an earlier one, perhaps another
// account's in a shared history, not this stop's.
func waitLimit(path string, start time.Time, msg string) (transcript.Limit, bool) {
	if path == "" {
		return transcript.Limit{}, false
	}
	for {
		// The row is made before the hook starts and written soon after.
		l, ok := transcript.LastLimit(path, start.Add(-time.Minute))
		if ok && (msg == "" || strings.TrimSpace(l.Text) == strings.TrimSpace(msg)) {
			return l, true
		}
		if time.Since(start) >= limitPoll {
			return transcript.Limit{}, false
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// endClaude sends SIGTERM to the Claude Code process pid after a moment,
// once Claude Code has read the hook's answer. A variable so tests can run
// the hook under a relay.
var endClaude = func(pid int) error {
	return exec.Command("sh", "-c", "sleep 0.3; kill -TERM "+strconv.Itoa(pid)).Start()
}

// hookSwitch notes the session for /ccam:switch and ends the turn before
// any request is made. Under a relay of this ccam it also ends Claude Code,
// and the relay carries the conversation on as the other account.
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
	if !fsutil.Exists(in.Transcript) {
		// Nothing written yet: there is no conversation to carry on.
		return stop("ccam: this session has no conversation yet, so there is nothing to carry on. To work as another account, type /exit and run `ccam " +
			cmp.Or(target, "next") + "`.")
	}
	t := handoff.Ticket{Account: account, SessionID: in.SessionID, Transcript: in.Transcript, Cwd: in.folder(),
		At: time.Now(), Reason: handoff.Switch, Target: target, Effort: cmp.Or(in.Effort.Level, os.Getenv("CLAUDE_EFFORT"))}
	exit := "ccam: type /exit, then run `ccam next --continue` to carry this conversation on as " + cmp.Or(target, "the next account") + "."
	// Without a target the relay goes on with `ccam next --continue --auto`.
	who := cmp.Or(target, "the account with the most room")
	save := func(msg string) error {
		if err := handoff.Save(e.paths.State, t); err != nil {
			return stop("ccam: " + err.Error())
		}
		return stop(msg)
	}
	// Only a relay of this ccam carries the note on; an older one, started
	// before an update, would not, so its session is left to /exit. Nor
	// is a session nobody sits at (-p, background, a teammate) the one it
	// runs.
	relay, ok := launcher.RelayPID(os.Getenv("CCAM_RELAY"))
	if !ok || unattended() {
		return save(exit)
	}
	claude, _ := strconv.Atoi(os.Getenv("CLAUDE_PID"))
	if runtime.GOOS == "windows" || claude <= 0 {
		// Claude Code cannot be ended from here (Windows, or Claude Code
		// before 2.1.214 without CLAUDE_PID): the relay carries on after /exit.
		t.Relay = os.Getenv("CCAM_RELAY")
		return save("ccam: type /exit; ccam then carries this conversation on as " + who + ".")
	}
	child, ok := launcher.RelayChild(claude, relay)
	if !ok {
		// A claude below the relayed one (started in its shell, say) has
		// CCAM_RELAY too, but is not what the relay runs.
		return save(exit)
	}
	t.Relay, t.Child, t.Stopped = os.Getenv("CCAM_RELAY"), child, true
	if err := handoff.Save(e.paths.State, t); err != nil {
		return stop("ccam: " + err.Error())
	}
	if err := endClaude(claude); err != nil {
		t.Stopped = false
		return save("ccam: type /exit; ccam then carries this conversation on as " + who + ".")
	}
	return stop("ccam: moving this conversation to " + who + "…")
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
