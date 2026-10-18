package cli

import (
	"fmt"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/next"
)

// pickAuto chooses by usage; filled in with the usage package.
func pickAuto(e *env, cfg *config.Config, cands []next.Candidate, sameHistory bool) (*config.Profile, *config.Profile, error) {
	return nil, nil, fmt.Errorf("--auto 는 아직 지원하지 않습니다")
}
