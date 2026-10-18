package cli

import (
	"context"
	"fmt"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// pickAuto picks the logged-in account with the lowest 5-hour usage,
// skipping full ones; with sameHistory only the current history group.
func pickAuto(e *env, cfg *config.Config, cands []next.Candidate, sameHistory bool) (*config.Profile, *config.Profile, error) {
	from, _, _ := next.Rotate(cands, false)
	var accts []usage.Account
	byName := map[string]*config.Profile{}
	all := usageAccounts(e, cfg)
	for i, c := range cands {
		if !c.LoggedIn {
			continue
		}
		if sameHistory && from != nil {
			h := from.Share.Get("history")
			if h == config.Local || c.Profile.Share.Get("history") != h {
				continue
			}
		}
		accts = append(accts, all[i])
		byName[c.Profile.Name] = c.Profile
	}
	if len(accts) == 0 {
		return nil, nil, fmt.Errorf("고를 수 있는 로그인된 계정이 없습니다")
	}
	reports := newFetcher(usageStateDir(e)).FetchAll(context.Background(), accts, usage.Options{}, nil)
	view := make([]*usage.Report, len(reports))
	for i := range reports {
		view[i] = &reports[i]
	}
	best := usage.Best(view)
	if best == nil {
		return nil, nil, fmt.Errorf("사용량을 알 수 있는 여유 계정이 없습니다 (`ccam usage` 로 확인)")
	}
	return from, byName[best.Name], nil
}
