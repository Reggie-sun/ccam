package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// exitNoRoom is EX_TEMPFAIL from sysexits.h: no account has room now, so a
// script may try again later.
const exitNoRoom = 75

// codedError ends ccam with its message and a specific exit code.
type codedError struct {
	code int
	msg  string
}

func (e codedError) Error() string { return e.msg }

func noRoom(msg string) error { return codedError{exitNoRoom, msg} }

// choiceFlags are the flags that override how an account is chosen.
type choiceFlags struct {
	maxUsage float64
	strategy string
	wait     string
}

// DefaultWait is how long a bare --wait waits at most: one 5-hour window.
const DefaultWait = "5h"

func (c *choiceFlags) add(cmd *cobra.Command) {
	cmd.Flags().Float64Var(&c.maxUsage, "max-usage", 0, "prefer accounts below this usage, 1-100 (config: auto_max_usage)")
	cmd.Flags().StringVar(&c.strategy, "strategy", "", "headroom (lowest 5-hour usage) or expiring (soonest reset) (config: auto_strategy)")
	cmd.Flags().StringVar(&c.wait, "wait", "", "when no account has room, wait for one, at most this long (default "+DefaultWait+")")
	cmd.Flags().Lookup("wait").NoOptDefVal = DefaultWait
}

// policy is the config's choice rule with the flags and claude args on top.
func (c choiceFlags) policy(cfg *config.Config, args []string) (usage.Policy, error) {
	p := usage.Policy{MaxUsage: cfg.AutoMaxUsage, Strategy: cfg.AutoStrategy, Model: modelArg(args)}
	if c.maxUsage != 0 {
		p.MaxUsage = c.maxUsage
	}
	if c.strategy != "" {
		p.Strategy = c.strategy
	}
	if err := p.Check(); err != nil {
		return p, err
	}
	return p, nil
}

// waitFor parses a --wait value; "" means do not wait.
func waitFor(v string) (time.Duration, error) {
	if v == "" {
		return 0, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("--wait takes a duration such as 30m or 2h, not %q", v)
	}
	return d, nil
}

// modelArg returns the --model value among claude args.
func modelArg(args []string) string {
	for i, a := range args {
		if a == "--" {
			break
		}
		if v, ok := strings.CutPrefix(a, "--model="); ok {
			return v
		}
		if a == "--model" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

// Clock and sleep, replaced by tests.
var (
	clock = time.Now
	sleep = time.Sleep
)

// afterReset is how long ccam waits past a reset before looking again: the
// endpoint takes a moment to show the new window.
const afterReset = 30 * time.Second

// lookAgain is how long ccam waits when no full account says when it resets.
const lookAgain = 10 * time.Minute

// untilRoom calls try until it chooses an account. When try finds none and
// wait is set, it sleeps until the first full account resets (per p.FreeAt)
// and tries again, until wait runs out; then, or without wait, it fails with
// exit code 75 and the message none.
func untilRoom(e *env, p usage.Policy, wait time.Duration, none string, try func() (*config.Profile, []*usage.Report, error)) (*config.Profile, error) {
	deadline := clock().Add(wait)
	for {
		prof, reports, err := try()
		if err != nil || prof != nil {
			return prof, err
		}
		now := clock()
		if wait <= 0 {
			return nil, noRoom(none)
		}
		at := p.FreeAt(reports, now)
		if at.IsZero() {
			at = now.Add(lookAgain)
		} else {
			at = at.Add(afterReset)
		}
		if at.After(deadline) {
			return nil, noRoom(none + fmt.Sprintf(" (and none gets room before --wait ends at %s)", deadline.Local().Format("15:04")))
		}
		fmt.Fprintf(e.stderr, "ccam: no account has room; waiting until %s\n", at.Local().Format("15:04"))
		sleep(at.Sub(now))
	}
}

// reportsOf fetches the usage of profs at once, in their order.
func reportsOf(e *env, cfg *config.Config, profs []*config.Profile, opt usage.Options) []*usage.Report {
	all := usageAccounts(e, cfg)
	var accts []usage.Account
	for _, p := range profs {
		for _, a := range all {
			if a.Name == p.Name {
				accts = append(accts, a)
			}
		}
	}
	reports := e.fetcher().FetchAll(context.Background(), accts, opt, nil)
	out := make([]*usage.Report, len(reports))
	for i := range reports {
		out[i] = &reports[i]
	}
	return out
}
