package cli

import (
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/transcript"
)

// conversation is the one `ccam next --continue` carries on.
type conversation struct {
	path   string          // the .jsonl file, links resolved
	owner  *config.Profile // the account it belongs to (the last used of those that see it)
	cwd    string
	ticket *handoff.Ticket // the note the session left, if it is what continues
}

// lastConversation finds the most recent conversation of the folder cwd in
// any account. Accounts that share history see the same files, and the one
// used last counts as their owner. note, when not nil, is the note to carry
// on whatever else is newer (the relay passes its session's own).
func lastConversation(e *env, cfg *config.Config, cands []next.Candidate, cwd string, note *handoff.Ticket) (*conversation, error) {
	order := slices.Clone(cands)
	slices.SortStableFunc(order, func(a, b next.Candidate) int { return b.LastUsed.Compare(a.LastUsed) })
	seen := map[string]bool{}
	var best *conversation
	var at int64
	for _, c := range order {
		dir, err := filepath.EvalSymlinks(transcript.ProjectDir(profile.Dir(e.paths, c.Profile), cwd))
		if err != nil || seen[dir] {
			continue
		}
		seen[dir] = true
		if path, mod, ok := transcript.Newest(dir); ok && (best == nil || mod.UnixNano() > at) {
			best, at = &conversation{path: path, owner: c.Profile, cwd: cwd}, mod.UnixNano()
		}
	}
	// A note left by the session (a usage limit, /ccam:switch) names the
	// exact conversation, unless a newer one has been written since.
	t, ok := note, note != nil
	if !ok {
		t, ok = handoff.ForFolder(e.paths.State, cwd, time.Now())
	}
	if ok {
		if p, found := cfg.Find(t.Account); found && (note != nil || best == nil || t.At.UnixNano() >= at-int64(2*time.Minute)) {
			path := t.Transcript
			if real, err := filepath.EvalSymlinks(path); err == nil {
				path = real
			}
			return &conversation{path: path, owner: p, cwd: cwd, ticket: t}, nil
		}
	}
	if best == nil {
		return nil, fmt.Errorf("no account has a conversation in %s to continue", cwd)
	}
	return best, nil
}

// continueArgs are the claude arguments that carry conv on as to, and a
// line about it: --resume <id> when to sees the conversation's folder
// (--continue could open a newer conversation there); else Claude Code
// copies the conversation into to (--fork-session), leaving the original as
// it is.
func continueArgs(e *env, cfg *config.Config, conv *conversation, to *config.Profile, args []string) ([]string, string) {
	what := "the last conversation"
	if title := transcript.Title(conv.path); title != "" {
		what = fmt.Sprintf("%q", title)
	}
	if n := transcript.Size(conv.path); n > 0 {
		what += fmt.Sprintf(" (about %s tokens)", tokens(n))
		if org := newOrgs(e, cfg, false); !config.SamePool(org.of(conv.owner), org.of(to)) {
			what += "; prompt caches are not shared between organizations, so the first reply reads it all again"
		}
	}
	mine, err := filepath.EvalSymlinks(transcript.ProjectDir(profile.Dir(e.paths, to), conv.cwd))
	if err == nil && mine == filepath.Dir(conv.path) {
		return append([]string{"--resume", transcript.ID(conv.path)}, args...), "continuing " + what
	}
	note := fmt.Sprintf("continuing %s; %s does not share %s's history, so Claude Code copies it into %s", what, to.Name, conv.owner.Name, to.Name)
	return append([]string{"--resume", conv.path, "--fork-session"}, args...), note
}

// tokens writes n as 950, 12k or 1.2M.
func tokens(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	case n >= 1000:
		return fmt.Sprintf("%dk", (n+500)/1000)
	}
	return fmt.Sprint(n)
}
