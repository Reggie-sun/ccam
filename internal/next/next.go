// Package next picks the account `ccam next` switches to.
package next

import (
	"fmt"
	"sort"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
)

// Candidate is a profile with the facts the choice depends on.
type Candidate struct {
	Profile  *config.Profile
	LastUsed time.Time
	LoggedIn bool
}

// Rotate returns the profile after the most recently used one, in slot
// order, wrapping around. Profiles that are not logged in are skipped.
// With sameHistory, only profiles sharing the current one's history group
// qualify, so `--continue` can find the conversation.
func Rotate(cands []Candidate, sameHistory bool) (from, to *config.Profile, err error) {
	if len(cands) == 0 {
		return nil, nil, fmt.Errorf("no accounts yet")
	}
	sorted := append([]Candidate(nil), cands...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Profile.Slot < sorted[j].Profile.Slot })
	cur := 0
	for i, c := range sorted {
		if c.LastUsed.After(sorted[cur].LastUsed) {
			cur = i
		}
	}
	from = sorted[cur].Profile
	hist := from.Share.Get("history")
	for k := 1; k < len(sorted); k++ {
		c := sorted[(cur+k)%len(sorted)]
		if !c.LoggedIn {
			continue
		}
		if sameHistory && (hist == config.Local || c.Profile.Share.Get("history") != hist) {
			continue
		}
		return from, c.Profile, nil
	}
	if sameHistory {
		return from, nil, fmt.Errorf("no other account shares conversation history with %s. Turn it on with `ccam share <name> --history on`", from.Name)
	}
	return from, nil, fmt.Errorf("no logged-in account other than %s", from.Name)
}
