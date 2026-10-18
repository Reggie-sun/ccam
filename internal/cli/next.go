package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/profile"
)

func nextCmd(e *env) *cobra.Command {
	var cont, auto bool
	cmd := &cobra.Command{
		Use:   "next [--auto] [--continue] [-- claude 인자]",
		Short: "다음 계정으로 실행 (--continue: 마지막 대화 이어가기)",
		Long: "마지막으로 쓴 계정의 다음 번호 계정으로 Claude Code를 엽니다.\n" +
			"--continue 는 대화 내역을 공유하는 계정 중에서 고르고 claude --continue 로 엽니다.\n" +
			"--auto 는 5시간 사용률이 가장 낮은 계정을 고릅니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			cands := candidates(e, cfg)
			var to *config.Profile
			var from *config.Profile
			if auto {
				from, to, err = pickAuto(e, cfg, cands, cont)
			} else {
				from, to, err = next.Rotate(cands, cont)
			}
			if err != nil {
				return err
			}
			if cont {
				args = append([]string{"--continue"}, args...)
			}
			if from != nil {
				fmt.Fprintf(e.stderr, "→ %s → %s\n", from.Name, to.Name)
			}
			return launcher.LaunchWith(e.paths, cfg, to.Name, args, e.dryRun, e.stdout)
		},
	}
	cmd.Flags().BoolVarP(&cont, "continue", "c", false, "대화 내역을 공유하는 다음 계정에서 마지막 대화 이어가기")
	cmd.Flags().BoolVar(&auto, "auto", false, "5시간 사용률이 가장 낮은 계정 고르기")
	return cmd
}

func candidates(e *env, cfg *config.Config) []next.Candidate {
	out := make([]next.Candidate, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		out = append(out, next.Candidate{
			Profile:  p,
			LastUsed: launcher.LastUsed(e.paths, p.Name),
			LoggedIn: creds.Has(profile.Dir(e.paths, p), e.paths.Home),
		})
	}
	return out
}
