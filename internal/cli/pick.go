package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/lock"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// recentPick is how long a picked account goes to the back of the line, so
// workers that start together get different accounts.
const recentPick = 2 * time.Minute

type pickResult struct {
	Name     string        `json:"name"`
	Slot     int           `json:"slot"`
	Email    string        `json:"email,omitempty"`
	FiveHour *usage.Window `json:"five_hour,omitempty"`
	Weekly   *usage.Window `json:"weekly,omitempty"`
	Source   string        `json:"source,omitempty"`
}

func pickCmd(e *env) *cobra.Command {
	var asJSON bool
	var model string
	var choice choiceFlags
	cmd := &cobra.Command{
		Use:   "pick [--json] [--model <model>] [--max-usage N] [--strategy S] [--wait[=D]]",
		Short: "Print the account with the most room, for scripts and parallel workers",
		Long: "Chooses a logged-in account by usage, the way next --auto does, and prints its\n" +
			"name. Calls made within 2 minutes of each other get different accounts while\n" +
			"there are enough accounts with room, so workers started together spread out.\n" +
			"With --model, that model's weekly limit counts too. With no account to pick,\n" +
			"ccam exits with code 75; --wait waits for an account to reset instead.",
		Example: "  ccam \"$(ccam pick)\" -p \"fix the failing tests\"\n" +
			"  for t in a b c; do ccam \"$(ccam pick)\" -p \"task $t\" & done; wait\n" +
			"  ccam pick --json --model opus --max-usage 80",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			p, err := choice.policy(cfg, nil)
			if err != nil {
				return err
			}
			p.Model = model
			d, err := waitFor(choice.wait)
			if err != nil {
				return err
			}
			var profs []*config.Profile
			for i := range cfg.Profiles {
				prof := &cfg.Profiles[i]
				if creds.Has(profile.Dir(e.paths, prof), e.paths.Home) {
					profs = append(profs, prof)
				}
			}
			if len(profs) == 0 {
				return noRoom("no logged-in account to pick")
			}
			var picked *usage.Report
			prof, err := untilRoom(e, p, d, "no account with known usage and headroom (see `ccam usage`)", func() (*config.Profile, []*usage.Report, error) {
				prof, r, reports, err := pickOnce(e, cfg, profs, p)
				picked = r
				return prof, reports, err
			})
			if err != nil {
				return err
			}
			if !asJSON {
				e.printf("%s\n", prof.Name)
				return nil
			}
			out := pickResult{Name: prof.Name, Slot: prof.Slot, Email: prof.Email, Source: picked.Source}
			if l := picked.Limits; l != nil {
				out.FiveHour, out.Weekly = &l.FiveHour, &l.Weekly
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			e.printf("%s\n", data)
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "print the account and its usage as JSON")
	cmd.Flags().StringVar(&model, "model", "", "count this model's weekly limit too, e.g. opus")
	choice.add(cmd)
	return cmd
}

// pickPath records when an account was last picked.
func pickPath(e *env, name string) string { return filepath.Join(e.paths.State, "picks", name) }

// pickOnce chooses under a lock, so concurrent picks see each other's
// choice: accounts picked within recentPick come last.
func pickOnce(e *env, cfg *config.Config, profs []*config.Profile, p usage.Policy) (*config.Profile, *usage.Report, []*usage.Report, error) {
	l, err := waitLock(filepath.Join(e.paths.State, "picks.lock"), 30*time.Second)
	if err != nil {
		return nil, nil, nil, err
	}
	defer l.Unlock()
	reports := reportsOf(e, cfg, profs, usage.Options{})
	now := clock()
	var fresh []*usage.Report
	for _, r := range reports {
		if info, err := os.Stat(pickPath(e, r.Name)); err != nil || now.Sub(info.ModTime()) >= recentPick {
			fresh = append(fresh, r)
		}
	}
	best := p.Choose(fresh, now)
	if best == nil {
		best = p.Choose(reports, now)
	}
	if best == nil {
		return nil, nil, reports, nil
	}
	if !e.dryRun {
		_ = fsutil.WriteFast(pickPath(e, best.Name), []byte(now.UTC().Format(time.RFC3339)))
	}
	prof, _ := cfg.Find(best.Name)
	return prof, best, reports, nil
}

// waitLock takes the flock at path, retrying until timeout.
func waitLock(path string, timeout time.Duration) (*lock.Flock, error) {
	deadline := time.Now().Add(timeout)
	for {
		l, err := lock.TryFlock(path)
		if err == nil || !errors.Is(err, lock.ErrBusy) {
			return l, err
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("another ccam pick holds %s", path)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
