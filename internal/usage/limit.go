package usage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/fsutil"
)

// A limit record keeps what Claude Code reported when it refused a request
// at a usage limit: that window is full until its reset. Fetch lays it over
// whatever else is known, so an account that just hit its limit is not
// chosen again before then, however old or behind the saved numbers are
// (state/usage/<name>.limit.json, written by the ccam plugin's hook).
type limitRecord struct {
	At      time.Time            `json:"at"`      // when the last refused request was
	Windows map[string]time.Time `json:"windows"` // window → when it resets
}

// limitWindows are the windows a record keeps, by the names Claude Code's
// rate_limit_info gives them, in the order they are laid over a report.
// The Fable allowance (seven_day_overage_included) and usage credits
// (overage) are no windows ccam chooses by.
var limitWindows = []string{"five_hour", "seven_day", "seven_day_opus", "seven_day_sonnet"}

// LimitWindow reports whether SaveLimit records window.
func LimitWindow(window string) bool { return slices.Contains(limitWindows, window) }

func (f *Fetcher) limitPath(name string) string { return filepath.Join(f.StateDir, name+".limit.json") }

// SaveLimit records that account name's window was full at at, when Claude
// Code refused a request, until resetsAt. It keeps the account's other
// windows that have not reset by then, and records nothing for a limit that
// is already over.
func SaveLimit(stateDir, name, window string, resetsAt, at time.Time) error {
	if !LimitWindow(window) {
		return fmt.Errorf("ccam does not record the %s window", window)
	}
	if !resetsAt.After(at) {
		return nil
	}
	f := Fetcher{StateDir: stateDir}
	rec, _ := f.readLimit(name, at)
	if rec.Windows == nil {
		rec.Windows = map[string]time.Time{}
	}
	rec.Windows[window] = resetsAt.UTC()
	if at.After(rec.At) {
		rec.At = at.UTC()
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return fsutil.AtomicWrite(f.limitPath(name), data, 0o600)
}

// DropLimit forgets account name's record: its limits were another
// login's.
func DropLimit(stateDir, name string) { os.Remove((&Fetcher{StateDir: stateDir}).limitPath(name)) }

// readLimit reads account name's record with only the windows that have
// not reset by now.
func (f *Fetcher) readLimit(name string, now time.Time) (limitRecord, bool) {
	data, err := os.ReadFile(f.limitPath(name))
	if err != nil {
		return limitRecord{}, false
	}
	var rec limitRecord
	if json.Unmarshal(data, &rec) != nil {
		return limitRecord{}, false
	}
	for w, at := range rec.Windows {
		if !at.After(now) {
			delete(rec.Windows, w)
		}
	}
	return rec, len(rec.Windows) > 0
}

// withLimit lays the account's limit record over r: each window in it is
// at 100% until its reset. With nothing else known, the record alone is
// what r says.
func (f *Fetcher) withLimit(a Account, r Report) Report {
	if r.Status == StatusAPI || r.Status == StatusLoggedOut {
		return r
	}
	rec, ok := f.readLimit(a.Name, f.Now())
	if !ok {
		return r
	}
	var l Limits
	if r.Limits != nil {
		l = *r.Limits
		l.Scoped = slices.Clone(l.Scoped)
	} else {
		r.At, r.Source = rec.At, "limit"
	}
	for _, w := range limitWindows {
		at, ok := rec.Windows[w]
		if !ok {
			continue
		}
		full := Window{Percent: 100, ResetsAt: at, Known: true}
		label := ""
		switch w {
		case "five_hour":
			l.FiveHour, label = full, "5-hour"
		case "seven_day":
			l.Weekly, label = full, "Weekly"
		default: // seven_day_opus, seven_day_sonnet
			fam := strings.TrimPrefix(w, "seven_day_")
			i := slices.IndexFunc(l.Scoped, func(s Scoped) bool { return strings.Contains(strings.ToLower(s.Name), fam) })
			if i < 0 {
				l.Scoped = append(l.Scoped, Scoped{Name: strings.ToUpper(fam[:1]) + fam[1:]})
				i = len(l.Scoped) - 1
			}
			l.Scoped[i].Window, label = full, l.Scoped[i].Name
		}
		r.LimitedWindows = append(r.LimitedWindows, label)
	}
	r.Limits, r.Limited = &l, rec.At
	return r
}
