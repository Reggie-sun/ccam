package launcher

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/x/term"

	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// relayProtocol starts the CCAM_RELAY value the relay gives its claude:
// "2:<the relaying ccam's process>". Hooks end Claude Code only under a
// relay of this form; a relay of an older ccam would not carry their notes
// on.
const relayProtocol = "2:"

// RelayPID is the relaying ccam's process in a CCAM_RELAY value of this
// protocol.
func RelayPID(v string) (int, bool) {
	s, ok := strings.CutPrefix(v, relayProtocol)
	if !ok {
		return 0, false
	}
	pid, err := strconv.Atoi(s)
	if err != nil || pid <= 0 {
		return 0, false
	}
	return pid, true
}

// RelayChild is the relay's own child for the Claude Code process claude
// (its CLAUDE_PID): claude itself when relay started it, or its parent when
// Claude Code started itself again (as gateway, Bedrock and Vertex settings
// make it) and that parent is the same program, or the cmd.exe that runs a
// batch claude (npm's claude.cmd) on Windows. Any other claude, such as one
// Claude ran in its shell with CCAM_RELAY inherited, is not the relay's.
func RelayChild(claude, relay int) (int, bool) {
	parent, ok := session.Parent(claude)
	if !ok {
		return 0, false
	}
	if parent == relay {
		return claude, true
	}
	if grand, ok := session.Parent(parent); !ok || grand != relay {
		return 0, false
	}
	a, ok1 := session.Exe(parent)
	b, ok2 := session.Exe(claude)
	if !ok1 || !ok2 {
		return 0, false
	}
	fa, err1 := os.Stat(a)
	fb, err2 := os.Stat(b)
	if err1 != nil || err2 != nil || !os.SameFile(fa, fb) && !batchShell(fa) {
		return 0, false
	}
	return parent, true
}

// Terminate ends the Claude Code process claude (its CLAUDE_PID) with
// SIGTERM for the stop-failure hook, but only while it is still the process
// the hook looked at: the relaying ccam's own claude (RelayChild(claude,
// relay) gives child) that started at start (session.StartTime). On Linux
// os.FindProcess holds the process by a pidfd taken before those checks,
// so a pid given to another process in between is never signaled;
// elsewhere the checks come right before the kill. Windows cannot end
// Claude Code from here: an error. A variable so tests can stand in for it.
var Terminate = func(claude, relay, child int, start string) error {
	if runtime.GOOS == "windows" {
		return errors.New("ccam cannot end Claude Code on Windows")
	}
	proc, err := os.FindProcess(claude)
	if err != nil {
		return err
	}
	defer func() { _ = proc.Release() }()
	if c, ok := RelayChild(claude, relay); !ok || c != child {
		return errors.New("the Claude Code process is no longer the relay's")
	}
	// An ended claude the relay has not reaped yet still has the relay as
	// its parent, but no start time.
	if s, ok := session.StartTime(claude); !ok || s != start {
		return errors.New("the Claude Code process has ended")
	}
	return proc.Signal(syscall.SIGTERM)
}

// relayable reports whether claude args start an interactive session on a
// terminal, the only kind /ccam:switch can come from. Not `claude
// remote-control`: the sessions it serves are its own.
func relayable(args []string) bool {
	if !term.IsTerminal(os.Stdin.Fd()) || claudeVerb(args) == verbRemoteControl {
		return false
	}
	for _, a := range args {
		if a == "-p" || a == "--print" || len(a) > 8 && a[:8] == "--print=" {
			return false
		}
	}
	return true
}

// relay runs claude as a child (relay = true in config.toml, or
// auto_switch = true). When the session asked to move on (/ccam:switch) or
// the hook chose an account after a usage limit (auto_switch), and claude
// ended the way its note expects (hopOK), `ccam next --continue` takes this
// process's place with that note and carries the conversation on as the
// other account, with the arguments claude started with that resuming does
// not restore (CarryArgs, WithState).
func (pl *Plan) relay(p paths.Paths) error {
	// The hook writes this into its note, so only this very session is
	// carried on.
	me := relayProtocol + strconv.Itoa(os.Getpid())
	pl.Env = append(pl.Env, "CCAM_RELAY="+me)
	start := time.Now()
	code, err := pl.Run(p, os.Stdin, os.Stdout, os.Stderr)
	pl.Done()
	if err != nil {
		return err
	}
	if t, ok := handoff.ForRelay(p.State, pl.Profile.Name, me, pl.run.child, start); ok {
		h := hop{note: t, child: pl.run.child, code: code, start: start, signaled: pl.run.signaled.Load(), terminal: term.IsTerminal(os.Stdin.Fd())}
		if ns := pl.run.interrupt.Load(); ns != 0 {
			h.interrupt = time.Unix(0, ns)
		}
		ok, why := hopOK(h)
		if why != "" {
			fmt.Fprintln(os.Stderr, why)
		}
		if ok {
			return pl.hop(t)
		}
	}
	if code != 0 {
		return ExitCode(code)
	}
	return nil
}

// hop replaces this process with `ccam next --continue` for the note t.
func (pl *Plan) hop(t *handoff.Ticket) error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	args, dropped := CarryArgs(pl.UserArgs)
	args, more := WithState(args, transcript.PermissionMode(t.Transcript), t.Effort)
	if dropped = append(dropped, more...); len(dropped) > 0 {
		fmt.Fprintf(os.Stderr, "ccam: not carried over: %s\n", strings.Join(dropped, ", "))
	}
	argv := []string{"ccam", "next", "--continue", "--note", t.File()}
	if t.Reason == handoff.Switch && t.Target == "" {
		argv = append(argv, "--auto")
	}
	if t.Reason == handoff.Limit {
		// The hook ended Claude Code without a word (auto_switch): say it here.
		notice, line := limitHop(t, time.Now())
		if term.IsTerminal(os.Stderr.Fd()) {
			fmt.Fprint(os.Stderr, handoff.Notification(notice))
		}
		fmt.Fprintln(os.Stderr, line)
	} else {
		fmt.Fprintln(os.Stderr, "ccam: carrying the conversation on (ccam next --continue)")
	}
	return execSelf(self, append(append(argv, "--"), args...))
}

// limitHop is what the relay says when it carries the conversation of the
// limit note t on: a desktop notification, and the line before the hop.
func limitHop(t *handoff.Ticket, now time.Time) (notice, line string) {
	label, reset := handoff.WindowLabel(t.Window), ""
	if label == "" {
		label = "usage"
	}
	if !t.ResetsAt.IsZero() {
		reset = handoff.FormatReset(t.ResetsAt, now)
	}
	notice = handoff.HitLimit(t.Account, label, reset) + ". Moving this conversation to " + t.Chosen
	if t.Why != "" {
		notice += " (" + t.Why + ")"
	}
	return notice + ".", handoff.HitLimit(t.Account, label, "") + "; continuing this conversation as " + t.Chosen
}

// hop is what the relay knows once claude has ended and its session left
// a note for this relay.
type hop struct {
	note      *handoff.Ticket
	child     int       // the claude process the relay ran
	code      int       // its exit code
	start     time.Time // when it started
	signaled  bool      // the relay itself got SIGHUP, SIGTERM or SIGQUIT
	interrupt time.Time // the last Ctrl+C the relay got; zero for none
	terminal  bool      // stdin is still a terminal
}

// hopOK reports whether the relay carries the conversation on: the note
// asks for it (/ccam:switch, or a usage limit with an account chosen) and
// is this claude's, claude ended as the note expects (exit code 143 after
// the hook's SIGTERM when the note says it stopped claude, else 0 after
// /exit), and nothing called the switch off since: no signal to the relay,
// no Ctrl+C after the note (or the limit), the terminal still there. why
// tells the person, on a terminal still there, why the switch did not
// happen.
func hopOK(h hop) (ok bool, why string) {
	t := h.note
	if t.Reason != handoff.Switch && (t.Reason != handoff.Limit || t.Chosen == "") {
		return false, ""
	}
	if t.At.Before(h.start) || t.Child != h.child && (t.Child != 0 || t.Stopped) || h.signaled || !h.terminal {
		return false, ""
	}
	since := t.At
	if !t.LimitAt.IsZero() {
		since = t.LimitAt
	}
	if h.interrupt.After(since) {
		return false, "ccam: switch cancelled. `ccam next --continue` carries the conversation on later."
	}
	if t.Stopped && h.code == 143 || !t.Stopped && h.code == 0 {
		return true, ""
	}
	return false, "ccam: Claude Code closed before ccam moved the conversation, so it stays on " + t.Account + ". `ccam next --continue` carries it on."
}
