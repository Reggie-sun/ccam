package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/transcript"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// liveSession is what Claude Code writes about a running session into
// <config dir>/sessions/<pid>.json (2.1.311).
type liveSession struct {
	PID             int    `json:"pid"`
	SessionID       string `json:"sessionId"`
	Cwd             string `json:"cwd"`
	Kind            string `json:"kind"`   // interactive, bg
	Status          string `json:"status"` // busy, idle, waiting
	StatusUpdatedAt int64  `json:"statusUpdatedAt"`
}

// quiet reports whether s is a session a person sits at that has waited for
// a prompt for at least idle: moving it then interrupts nothing.
func (s liveSession) quiet(now time.Time, idle time.Duration) bool {
	return s.Kind == "interactive" && s.Status == "idle" && s.StatusUpdatedAt > 0 &&
		now.Sub(time.UnixMilli(s.StatusUpdatedAt)) >= idle
}

func gatherCmd(e *env) *cobra.Command {
	var to string
	var idle time.Duration
	cmd := &cobra.Command{
		Use:   "gather [--to <account>] [--idle 3m]",
		Short: "Move idle sessions of other accounts to the account drain uses",
		Long: "With auto_strategy = \"drain\", moves every session of another account that has\n" +
			"been idle (waiting for a prompt) for --idle to the account drain chooses, the way\n" +
			"/ccam:switch does: the conversation goes on there. Busy sessions, sessions waiting\n" +
			"for an answer, background sessions and Remote Control are left alone, and so is a\n" +
			"session in a folder the account would first ask to trust (the home folder).\n" +
			"--to names the account instead, whatever auto_strategy says.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			target, err := gatherTarget(e, cfg, to)
			if err != nil || target == nil {
				return err
			}
			now := clock()
			tdir := profile.Dir(e.paths, target)
			for i := range cfg.Profiles {
				from := &cfg.Profiles[i]
				if from.Name == target.Name {
					continue
				}
				dir := profile.Dir(e.paths, from)
				metas, _ := filepath.Glob(filepath.Join(dir, "sessions", "*.json"))
				for _, meta := range metas {
					var s liveSession
					if data, err := os.ReadFile(meta); err != nil || json.Unmarshal(data, &s) != nil || !s.quiet(now, idle) {
						continue
					}
					if why := gatherOne(e, from, target, dir, tdir, s); why != "" {
						e.printf("%s %s (%s): %s\n", from.Name, s.SessionID[:min(8, len(s.SessionID))], s.Cwd, why)
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "move them to this account instead of the one drain chooses")
	cmd.Flags().DurationVar(&idle, "idle", 3*time.Minute, "how long a session must have waited for a prompt")
	return cmd
}

// gatherTarget is the account gather moves sessions to: --to, or the one
// drain chooses; nil (nothing to do) under another strategy.
func gatherTarget(e *env, cfg *config.Config, to string) (*config.Profile, error) {
	if to != "" {
		return cfg.MustFind(to)
	}
	if cfg.AutoStrategy != usage.Drain {
		return nil, nil
	}
	var profs []*config.Profile
	for i := range cfg.Profiles {
		if creds.Has(profile.Dir(e.paths, &cfg.Profiles[i]), e.paths.Home) {
			profs = append(profs, &cfg.Profiles[i])
		}
	}
	p := usage.Policy{MaxUsage: cfg.AutoMaxUsage, Strategy: usage.Drain}
	best := p.Choose(reportsOf(e, cfg, profs, usage.Options{}), clock())
	if best == nil {
		return nil, nil
	}
	prof, _ := cfg.Find(best.Name)
	return prof, nil
}

// gatherOne hands the session s of from to target through its relay, as
// /ccam:switch does, and says what happened ("" for nothing worth saying).
func gatherOne(e *env, from, target *config.Profile, dir, tdir string, s liveSession) string {
	relayVar := processEnv(s.PID, "CCAM_RELAY")
	relay, ok := launcher.RelayPID(relayVar)
	if !ok {
		return "" // not started by this ccam: only /exit moves it
	}
	child, ok := launcher.RelayChild(s.PID, relay)
	if !ok {
		return ""
	}
	started, ok := session.StartTime(s.PID)
	if !ok {
		return ""
	}
	if processStopped(s.PID) {
		// A stopped Claude Code (Ctrl+Z) handles no SIGTERM until it goes on.
		return "left here: Claude Code is stopped (Ctrl+Z)"
	}
	path := filepath.Join(transcript.ProjectDir(dir, s.Cwd), s.SessionID+".jsonl")
	if !fsutil.Exists(path) || transcript.RemoteControl(path) {
		return ""
	}
	if !mcp.SessionTrusts(tdir, s.Cwd) {
		if holdsHome(s.Cwd, e.paths.Home) {
			return "left here: " + target.Name + " would first ask to trust this folder"
		}
		if e.dryRun {
			return "would trust this folder for " + target.Name + " and move it there"
		}
		if err := mcp.Trust(tdir, s.Cwd); err != nil {
			return "left here: " + err.Error()
		}
	}
	if e.dryRun {
		return "would move it to " + target.Name
	}
	t := handoff.Ticket{Account: from.Name, SessionID: s.SessionID, Transcript: path, Cwd: s.Cwd, At: time.Now(),
		Reason: handoff.Switch, Target: target.Name, Relay: relayVar, Child: child, Stopped: true}
	if err := handoff.Save(e.paths.State, t); err != nil {
		return "left here: " + err.Error()
	}
	ticket := handoff.Path(e.paths.State, from.Name, s.SessionID)
	if err := launcher.Terminate(s.PID, relay, child, started); err != nil {
		handoff.Remove(ticket)
		return "left here: " + err.Error()
	}
	// Only a Claude Code that ended is carried on by its relay: one that did
	// not keeps its account, and its note must not wait for a later exit.
	for deadline := time.Now().Add(gatherExitWait); time.Now().Before(deadline); time.Sleep(100 * time.Millisecond) {
		if now, ok := session.StartTime(s.PID); !ok || now != started {
			return "moved to " + target.Name
		}
	}
	handoff.Remove(ticket)
	return "left here: Claude Code did not end"
}

// gatherExitWait is how long gather waits for Claude Code to end after
// SIGTERM. Tests change it.
var gatherExitWait = 5 * time.Second

// processStopped reports whether process pid is stopped (Linux: state T or
// t in /proc/<pid>/stat).
func processStopped(pid int) bool {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return false
	}
	_, rest, ok := strings.Cut(string(data), ") ")
	return ok && len(rest) > 0 && (rest[0] == 'T' || rest[0] == 't')
}

// processEnv is the value of the environment variable name of process pid
// (Linux); "" when it cannot be read.
func processEnv(pid int, name string) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid))
	if err != nil {
		return ""
	}
	for _, kv := range strings.Split(string(data), "\x00") {
		if v, ok := strings.CutPrefix(kv, name+"="); ok {
			return v
		}
	}
	return ""
}
