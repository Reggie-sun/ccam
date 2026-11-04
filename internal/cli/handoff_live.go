package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jhyunwoo/ccam/internal/usage"
)

// defaultHandoffAt is the percent of a window at which a session is asked
// to start keeping a handoff (config handoff_at).
const defaultHandoffAt = 85

// remindEvery is how many tool calls pass between two reminders once the
// session has been asked; each prompt of the person asks again.
const remindEvery = 6

// livePath is the detailed handoff a session keeps once its account runs
// low. ccam next --continue --brief names it to the account taking over.
func livePath(e *env, sessionID string) string {
	return filepath.Join(e.paths.State, "live", sessionID+".md")
}

const liveFirst = `ccam: this account's usage is at %.0f%% and may reach its limit soon. The conversation would then go on in a new session on another account that reads only a handoff file, not this transcript. Keep that file at %s with the Write tool (replace the whole file each time) and update it after every meaningful step until the session ends. Do not mention this to the user. Write it so a fresh session can carry on without the transcript:
1. Goal, and the user's requirements and constraints, in their own words.
2. Current state: what is done, what is in progress, and the exact next step.
3. Decisions made and why; approaches tried and rejected.
4. Files touched (path and what changed) and the key facts about the code.
5. Commands run, test and build results, errors seen.
6. Open questions, risks, and what must not be redone.
7. How to resume: working folder, branch, commands.
Write it now.`

const liveAgain = `ccam: usage is at %.0f%%. If anything changed since you last wrote the handoff at %s (progress, decisions, files, results, next step), rewrite it now.`

// hookHandoff asks a session to keep a detailed handoff file up to date
// once its account's usage passes handoff_at, so that the account taking
// over (ccam next --continue --brief) gets more than ccam's summary. It
// answers on UserPromptSubmit and PostToolUse; a subagent's tool calls and
// unattended sessions are left alone.
func hookHandoff(e *env, in hookInput, account string) error {
	if in.AgentID != "" || account == "" || in.SessionID == "" || unattended() {
		return nil
	}
	cfg, err := e.load()
	if err != nil {
		return nil
	}
	at := cfg.HandoffAt
	if at == 0 {
		at = defaultHandoffAt
	}
	percent, ok := usage.Peak(usageStateDir(e), account, clock())
	if !ok || percent < at {
		return nil
	}
	path := livePath(e, in.SessionID)
	if in.Event == "PostToolUse" && in.ToolInput.FilePath == path {
		return nil // the handoff being written is not a step to write about
	}
	counter := strings.TrimSuffix(path, ".md") + ".n"
	n := 0
	if data, err := os.ReadFile(counter); err == nil {
		n, _ = strconv.Atoi(strings.TrimSpace(string(data)))
	}
	n++
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil
	}
	_ = os.WriteFile(counter, []byte(strconv.Itoa(n)), 0o600)
	pruneLive(filepath.Dir(path))
	var text string
	switch {
	case n == 1:
		text = fmt.Sprintf(liveFirst, percent, path)
	case in.Event == "UserPromptSubmit" || n%remindEvery == 0:
		text = fmt.Sprintf(liveAgain, percent, path)
	default:
		return nil
	}
	return hookOutput(e, map[string]any{"hookSpecificOutput": map[string]any{"hookEventName": in.Event, "additionalContext": text}})
}

// pruneLive removes handoffs and counters that nobody has touched for a
// week, as pruneBriefs does for summaries.
func pruneLive(dir string) {
	entries, _ := os.ReadDir(dir)
	now := clock()
	for _, de := range entries {
		if info, err := de.Info(); err == nil && !de.IsDir() && now.Sub(info.ModTime()) > briefKeep {
			os.Remove(filepath.Join(dir, de.Name()))
		}
	}
}
