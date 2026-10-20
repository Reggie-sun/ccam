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
	var choice choiceFlags
	cmd := &cobra.Command{
		Use:   "next [--auto] [--continue] [-- claude args]",
		Short: "Run the next account (--continue: pick up the last conversation)",
		Long: "Opens Claude Code as the account numbered after the last one used.\n" +
			"--continue picks among accounts that share conversation history and opens it with claude --continue.\n" +
			"--auto picks by usage: the lowest 5-hour usage among accounts below auto_max_usage\n" +
			"(config.toml, default 100), skipping full ones; --max-usage and --strategy override\n" +
			"the config, and --wait waits for an account to reset when none has room. With no\n" +
			"account to pick, ccam exits with code 75.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			cands := candidates(e, cfg)
			var to *config.Profile
			var from *config.Profile
			if auto {
				p, perr := choice.policy(cfg, args)
				if perr != nil {
					return perr
				}
				from, to, err = pickAuto(e, cfg, cands, cont, p, choice.wait)
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
	cmd.Flags().BoolVarP(&cont, "continue", "c", false, "continue the last conversation on the next account that shares history")
	cmd.Flags().BoolVar(&auto, "auto", false, "pick the account by usage (most room)")
	choice.add(cmd)
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		for _, f := range []string{"max-usage", "strategy", "wait"} {
			if cmd.Flags().Changed(f) && !auto {
				return fmt.Errorf("--%s works with --auto", f)
			}
		}
		return nil
	}
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
