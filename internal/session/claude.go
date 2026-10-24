package session

import (
	"cmp"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// InUse reports whether Claude Code runs on the config folder dir: a
// session ccam started (a marker in markers) or any Claude Code process
// that lists itself in the folder (ClaudeRuns).
func InUse(markers, dir string) bool { return Active(markers) || ClaudeRuns(dir) }

// ClaudeRuns reports whether Claude Code lists a running process of its own
// in its config folder dir, however that was started (plain claude, an IDE,
// a background session, the daemon behind those): Claude Code writes
// sessions/<pid>.json for every session and daemon.lock for its daemon. A
// record counts while the process it names runs (see running), so one left
// behind by a crash, or naming a reused PID, does not. The records are
// Claude Code's: they are only read, never removed.
func ClaudeRuns(dir string) bool {
	sessions := filepath.Join(dir, "sessions")
	entries, _ := os.ReadDir(sessions)
	for _, e := range entries {
		// <pid>.json; the <pid>.<hex>.key files beside them are secrets.
		pid, ok := strings.CutSuffix(e.Name(), ".json")
		if _, err := strconv.Atoi(pid); ok && err == nil && running(filepath.Join(sessions, e.Name())) {
			return true
		}
	}
	return running(filepath.Join(dir, "daemon.lock"))
}

// maxRecord bounds what is read of a record. Claude Code's are well under a
// kilobyte; it drops a daemon.lock over 64 KiB itself.
const maxRecord = 64 << 10

// slack is how much later than its record a process may seem to have
// started (clocks move) where only the record's time tells a reused PID.
const slack = time.Minute

// running reports whether the Claude Code record at path names a process
// that still runs. Its PID must run, and the process must have started
// when the record says (procStartFt, else procStart, as Claude Code reads
// them); where ccam cannot compare that, it must not have started after the
// record was written (startedAt, ms since 1970).
func running(path string) bool {
	// Only a small regular file: opening a FIFO would wait for a writer.
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxRecord {
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var r struct {
		PID         int     `json:"pid"`
		ProcStart   string  `json:"procStart"`
		ProcStartFt string  `json:"procStartFt"` // Windows daemon.lock
		StartedAt   float64 `json:"startedAt"`
	}
	// PID 1 is never a session; Claude Code does not count it either.
	if json.Unmarshal(data, &r) != nil || r.PID <= 1 {
		return false
	}
	now, ok := StartTime(r.PID)
	if !ok {
		return false
	}
	if parent, ok := Parent(r.PID); ok && parent == os.Getpid() {
		return false // started by this ccam itself: its usage probe
	}
	if want := cmp.Or(r.ProcStartFt, r.ProcStart); want != "" {
		if same, known := sameStart(want, now); known {
			return same
		}
	}
	at, ok := wallStart(now)
	return !ok || r.StartedAt <= 0 || !at.After(time.UnixMilli(int64(r.StartedAt)).Add(slack))
}
