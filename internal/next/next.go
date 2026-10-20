// Package next picks the account `ccam next` switches to.
package next

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
)

// Candidate is a profile with the facts the choice depends on.
type Candidate struct {
	Profile  *config.Profile
	LastUsed time.Time
	LoggedIn bool
}

// Rotate returns the profile after from (the most recently used one when
// from is ""), in slot order, wrapping around. Profiles that are not logged
// in are skipped, and so are those fits refuses (nil: none).
func Rotate(cands []Candidate, from string, fits func(*config.Profile) bool) (cur, to *config.Profile, err error) {
	if len(cands) == 0 {
		return nil, nil, fmt.Errorf("no accounts yet")
	}
	sorted := append([]Candidate(nil), cands...)
	slices.SortFunc(sorted, func(a, b Candidate) int { return cmp.Compare(a.Profile.Slot, b.Profile.Slot) })
	at := 0
	for i, c := range sorted {
		if from == "" && c.LastUsed.After(sorted[at].LastUsed) || c.Profile.Name == from {
			at = i
		}
	}
	cur = sorted[at].Profile
	for k := 1; k < len(sorted); k++ {
		c := sorted[(at+k)%len(sorted)]
		if !c.LoggedIn || fits != nil && !fits(c.Profile) {
			continue
		}
		return cur, c.Profile, nil
	}
	return cur, nil, fmt.Errorf("no logged-in account other than %s", cur.Name)
}
