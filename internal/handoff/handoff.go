// Package handoff keeps the note a Claude Code session leaves when it stops
// at a usage limit or asks to move to another account (/ccam:switch), so
// that `ccam next --continue` carries on that very conversation. The ccam
// plugin's hooks write it; nothing runs in the background.
package handoff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// Reasons a ticket is written for.
const (
	Limit  = "limit"  // the session stopped at a usage limit
	Switch = "switch" // /ccam:switch asked to move on
)

// Ticket is one account's latest note.
type Ticket struct {
	Account    string    `json:"account"`
	SessionID  string    `json:"session_id"`
	Transcript string    `json:"transcript_path"`
	Cwd        string    `json:"cwd"`
	At         time.Time `json:"at"`
	Reason     string    `json:"reason"`
	Limit      string    `json:"limit,omitempty"`  // session, weekly, model or spend
	Target     string    `json:"target,omitempty"` // the account /ccam:switch named
	Relay      string    `json:"relay,omitempty"`  // the relaying ccam's process, with relay on
}

// MaxAge is how long a ticket stays good.
const MaxAge = 12 * time.Hour

func path(state, account string) string { return filepath.Join(state, "handoff", account+".json") }

// Save writes t as its account's ticket.
func Save(state string, t Ticket) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return fsutil.AtomicWrite(path(state, t.Account), data, 0o600)
}

// Load reads account's ticket.
func Load(state, account string) (*Ticket, bool) {
	data, err := os.ReadFile(path(state, account))
	if err != nil {
		return nil, false
	}
	var t Ticket
	if json.Unmarshal(data, &t) != nil {
		return nil, false
	}
	return &t, true
}

// Remove drops account's ticket.
func Remove(state, account string) { os.Remove(path(state, account)) }

// Paths lists what Save writes for account, for rename and remove.
func Paths(state, account string) []string { return []string{path(state, account)} }

// ForFolder returns the newest good ticket left in the folder cwd.
func ForFolder(state, cwd string, now time.Time) (*Ticket, bool) {
	entries, err := os.ReadDir(filepath.Join(state, "handoff"))
	if err != nil {
		return nil, false
	}
	var best *Ticket
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(state, "handoff", e.Name()))
		if err != nil {
			continue
		}
		var t Ticket
		if json.Unmarshal(data, &t) != nil || !t.Good(now) || !t.In(cwd) {
			continue
		}
		if best == nil || t.At.After(best.At) {
			best = &t
		}
	}
	return best, best != nil
}

// Good reports whether t is recent and its conversation is still there.
func (t *Ticket) Good(now time.Time) bool {
	return now.Sub(t.At) <= MaxAge && fsutil.Exists(t.Transcript)
}

// In reports whether t's conversation belongs to the folder cwd: the
// session started there, or Claude Code keeps the conversation under it.
func (t *Ticket) In(cwd string) bool {
	return samePath(t.Cwd, cwd) || filepath.Base(filepath.Dir(t.Transcript)) == filepath.Base(transcript.ProjectDir("", cwd))
}

func samePath(a, b string) bool {
	if a == b {
		return true
	}
	ra, err1 := filepath.EvalSymlinks(a)
	rb, err2 := filepath.EvalSymlinks(b)
	return err1 == nil && err2 == nil && ra == rb
}

// LimitKind reads which limit stopped a session from Claude Code's message
// ("You've hit your session limit · resets 3:45pm", Claude Code 2.1.304):
// session, weekly, spend or model; "" for anything else, such as a short
// server throttle or an API key's rate limit, which are not usage limits.
func LimitKind(msg string) string {
	m := strings.TrimSpace(strings.ReplaceAll(msg, "’", "'"))
	rest, ok := strings.CutPrefix(m, "You've hit your ")
	if !ok {
		return ""
	}
	head := rest[:min(len(rest), 40)]
	switch {
	case strings.HasPrefix(rest, "session limit"):
		return "session"
	case strings.HasPrefix(rest, "weekly limit"):
		return "weekly"
	case strings.Contains(head, "spend limit"), strings.Contains(head, "usage credit limit"):
		return "spend"
	case strings.Contains(head, " limit"):
		return "model"
	}
	return ""
}

// Reset is the "resets …" part of Claude Code's message, or "".
func Reset(msg string) string {
	_, after, ok := strings.Cut(msg, "resets ")
	if !ok {
		return ""
	}
	after, _, _ = strings.Cut(after, "\n")
	after, _, _ = strings.Cut(after, " ·")
	return strings.TrimRight(strings.TrimSpace(after), ".")
}
