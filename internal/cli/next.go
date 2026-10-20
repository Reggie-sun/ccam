package cli

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/next"
	"github.com/jhyunwoo/ccam/internal/profile"
)

func nextCmd(e *env) *cobra.Command {
	var cont, auto bool
	var noteOf string
	var choice choiceFlags
	cmd := &cobra.Command{
		Use:   "next [--auto] [--continue] [-- claude args]",
		Short: "Run the next account (--continue: pick up the last conversation)",
		Long: "Opens Claude Code as the account numbered after the last one used.\n" +
			"--continue carries on this folder's last conversation: with claude --resume <id> when the account\n" +
			"shares its history, else Claude Code copies it over (--resume <file> --fork-session).\n" +
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
			var from, to *config.Profile
			var conv *conversation
			if cont {
				// Continue from the account that has the conversation.
				cwd, _ := os.Getwd()
				var own *handoff.Ticket
				if noteOf != "" {
					// The relay names its session's note; a stale one is no note.
					if t, ok := handoff.Load(e.paths.State, noteOf); ok && t.Good(time.Now()) {
						own = t
					}
				}
				if conv, err = lastConversation(e, cfg, cands, cwd, own); err != nil {
					return err
				}
				from = conv.owner
			}
			why := "" // why --auto chose to
			var t *handoff.Ticket
			if conv != nil {
				t = conv.ticket
			}
			if t != nil && t.Target != "" {
				// /ccam:switch named the account.
				if p, ok := cfg.Find(t.Target); ok {
					to = p
				} else {
					fmt.Fprintf(e.stderr, "ccam: the note names %s, which is no longer an account; choosing as usual\n", t.Target)
				}
			}
			switch {
			case to != nil:
			case auto:
				p, perr := choice.policy(cfg, args)
				if perr != nil {
					return perr
				}
				if t != nil {
					// The session asked to move on, or hit a limit its usage
					// may not show yet.
					cands = slices.DeleteFunc(slices.Clone(cands), func(c next.Candidate) bool { return c.Profile.Name == t.Account })
				}
				from, to, why, err = pickAuto(e, cfg, cands, from, p, choice.wait, choice.crossOrg)
			default:
				var fits func(*config.Profile) bool
				name := ""
				if from != nil {
					org := newOrgs(e, cfg, false)
					fits = func(p *config.Profile) bool { return org.fits(from, p) }
					name = from.Name
				}
				from, to, err = next.Rotate(cands, name, fits)
			}
			if err != nil {
				return err
			}
			switch {
			case from != nil && why != "":
				fmt.Fprintf(e.stderr, "→ %s → %s (%s)\n", from.Name, to.Name, why)
			case from != nil:
				fmt.Fprintf(e.stderr, "→ %s → %s\n", from.Name, to.Name)
			}
			if cont {
				var note string
				args, note = continueArgs(e, cfg, conv, to, args)
				fmt.Fprintf(e.stderr, "  %s\n", note)
			}
			var used func()
			if t != nil {
				// Used up once claude is about to start, not before: a launch
				// that fails leaves it for the next try.
				used = func() { handoff.Remove(e.paths.State, t.Account) }
			}
			return launcher.LaunchThen(e.paths, cfg, to.Name, args, e.dryRun, e.stdout, used)
		},
	}
	cmd.Flags().BoolVarP(&cont, "continue", "c", false, "continue this folder's last conversation on the next account")
	cmd.Flags().StringVar(&noteOf, "note", "", "with --continue: carry on the note of this account's session (the relay uses it)")
	_ = cmd.Flags().MarkHidden("note")
	cmd.Flags().BoolVar(&auto, "auto", false, "pick the account by usage (most room)")
	choice.add(cmd)
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		for _, f := range []string{"max-usage", "strategy", "wait", "cross-org"} {
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
