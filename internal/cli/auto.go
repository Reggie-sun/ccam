package cli

import (
	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// pickAuto chooses the logged-in account by usage (usage.Policy: below
// auto_max_usage first, by auto_strategy; full ones never). With
// sameHistory only accounts in the current history group qualify. With
// wait it waits for a full account to reset.
func pickAuto(e *env, cfg *config.Config, cands []next.Candidate, sameHistory bool, p usage.Policy, wait string) (*config.Profile, *config.Profile, error) {
	from, _, _ := next.Rotate(cands, false)
	var profs []*config.Profile
	for _, c := range cands {
		if !c.LoggedIn {
			continue
		}
		if sameHistory && from != nil {
			h := from.Share.Get("history")
			if h == config.Local || c.Profile.Share.Get("history") != h {
				continue
			}
		}
		profs = append(profs, c.Profile)
	}
	if len(profs) == 0 {
		return nil, nil, noRoom("no logged-in account to pick")
	}
	d, err := waitFor(wait)
	if err != nil {
		return nil, nil, err
	}
	to, err := untilRoom(e, p, d, "no account with known usage and headroom (see `ccam usage`)", func() (*config.Profile, []*usage.Report, error) {
		reports := reportsOf(e, cfg, profs, usage.Options{})
		if best := p.Choose(reports, clock()); best != nil {
			prof, _ := cfg.Find(best.Name)
			return prof, reports, nil
		}
		return nil, reports, nil
	})
	return from, to, err
}
