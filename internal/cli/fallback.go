package cli

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/creds"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/usage"
)

// limitText is what Claude Code prints when an account reaches its usage
// limit ("You've hit your limit · resets 3pm", older "Claude AI usage limit
// reached|<time>"), checked case-insensitively.
var limitText = [][]byte{[]byte("hit your limit"), []byte("usage limit reached")}

// Fallback runs `ccam --fallback [name] -p …`: claude runs as a child with
// its output held back. When it fails because the account reached its
// usage limit, the output is dropped and the same command runs again as
// another account, at most once per account. When every account is at its
// limit it exits with 75, or with wait waits for one to reset first. It
// returns the exit code.
func Fallback(p paths.Paths, sel string, args []string, dry bool, wait string) int {
	e := &env{paths: p, stdin: os.Stdin, in: bufio.NewReader(os.Stdin), stdout: os.Stdout, stderr: os.Stderr, dryRun: dry}
	err := fallback(e, sel, args, wait)
	var code launcher.ExitCode
	var coded codedError
	switch {
	case err == nil:
		return 0
	case errors.As(err, &code):
		return int(code)
	case errors.As(err, &coded):
		fmt.Fprintln(e.stderr, "ccam:", coded.msg)
		return coded.code
	}
	fmt.Fprintln(e.stderr, "ccam:", err)
	return 1
}

func fallback(e *env, sel string, args []string, wait string) error {
	if !printMode(args) {
		return fmt.Errorf("--fallback works for non-interactive runs with -p or --print. When an interactive session reaches its limit, use `ccam next --auto --continue`")
	}
	cfg, err := e.load()
	if err != nil {
		return err
	}
	p, err := choiceFlags{}.policy(cfg, args)
	if err != nil {
		return err
	}
	waitFor, err := waitFor(wait)
	if err != nil {
		return err
	}
	pl, err := launcher.Ready(e.paths, cfg, sel, args, e.dryRun)
	if err != nil {
		return err
	}
	if e.dryRun {
		pl.DryRun(e.stdout)
		fmt.Fprintln(e.stderr, "ccam: if this account reaches its usage limit, the command runs again as the logged-in account with the most room")
		return nil
	}
	// Each attempt gets the same input, so it is read once up front.
	var input []byte
	if !isTerminal(e.stdin) {
		if input, err = io.ReadAll(e.stdin); err != nil {
			return err
		}
	}
	tried := map[string]bool{}
	for {
		tried[pl.Profile.Name] = true
		stdin := e.stdin
		if input != nil {
			stdin = bytes.NewReader(input)
		}
		var out bytes.Buffer
		errTail := &tail{max: 64 << 10}
		code, err := pl.Run(e.paths, stdin, &out, io.MultiWriter(e.stderr, errTail))
		pl.Done()
		if err != nil {
			return err
		}
		if code == 0 || !limited(e, cfg, pl.Profile, out.Bytes(), errTail.buf) {
			_, _ = e.stdout.Write(out.Bytes())
			return exitWith(code)
		}
		next, _ := fallbackTarget(e, cfg, pl.Profile, tried, args, p)
		if next == nil && waitFor > 0 {
			// Every account was tried or is full: wait for one to reset, and
			// then any account may be tried again.
			fmt.Fprintf(e.stderr, "ccam: %s reached its usage limit, and no other logged-in account has room left\n", pl.Profile.Name)
			// The first look only finds out how long to wait: an account that
			// just hit its limit may still show room in saved usage, and must
			// not be run again at once.
			from, first := pl.Profile, true
			next, err = untilRoom(e, p, waitFor, "no logged-in account has room left", func() (*config.Profile, []*usage.Report, error) {
				prof, reports := fallbackTarget(e, cfg, from, map[string]bool{}, args, p)
				if first {
					first = false
					return nil, reports, nil
				}
				return prof, reports, nil
			})
			if err != nil {
				_, _ = e.stdout.Write(out.Bytes())
				return err
			}
			tried = map[string]bool{}
		}
		if next == nil {
			_, _ = e.stdout.Write(out.Bytes())
			fmt.Fprintf(e.stderr, "ccam: %s reached its usage limit, and no other logged-in account has room left\n", pl.Profile.Name)
			return launcher.ExitCode(exitNoRoom)
		}
		fmt.Fprintf(e.stderr, "ccam: %s reached its usage limit; running the command again as %s\n", pl.Profile.Name, next.Name)
		if pl, err = launcher.Ready(e.paths, cfg, next.Name, args, false); err != nil {
			return err
		}
	}
}

func exitWith(code int) error {
	if code != 0 {
		return launcher.ExitCode(code)
	}
	return nil
}

// printMode reports whether claude args ask for a non-interactive run.
func printMode(args []string) bool {
	for _, a := range args {
		if a == "--" {
			return false
		}
		if a == "-p" || a == "--print" {
			return true
		}
	}
	return false
}

// resumes reports whether claude args continue an earlier conversation,
// which only an account sharing its history can do.
func resumes(args []string) bool {
	return slices.ContainsFunc(args, func(a string) bool {
		return a == "-c" || a == "--continue" || a == "-r" || a == "--resume"
	})
}

// limited reports whether a failed run hit the account's usage limit: by
// Claude Code's message, or by the account's usage when it can be looked up.
func limited(e *env, cfg *config.Config, prof *config.Profile, out, errOut []byte) bool {
	for _, b := range [][]byte{out, errOut} {
		low := bytes.ToLower(b)
		for _, t := range limitText {
			if bytes.Contains(low, t) {
				return true
			}
		}
	}
	for _, a := range usageAccounts(e, cfg) {
		if a.Name == prof.Name {
			r := e.fetcher().Fetch(context.Background(), a, usage.Options{NoCache: true})
			return full(&r)
		}
	}
	return false
}

func full(r *usage.Report) bool {
	return r.Limits != nil && (r.Limits.FiveHour.Known && r.Limits.FiveHour.Percent >= 100 ||
		r.Limits.Weekly.Known && r.Limits.Weekly.Percent >= 100)
}

// fallbackTarget picks the next account: logged in, not tried yet, sharing
// history when the args resume a conversation. p chooses among those with
// known usage (full ones never); accounts whose usage is unknown (offline,
// lookup refused) come after them in number order. It also returns the
// reports, for waiting.
func fallbackTarget(e *env, cfg *config.Config, from *config.Profile, tried map[string]bool, args []string, p usage.Policy) (*config.Profile, []*usage.Report) {
	var profs []*config.Profile
	for i := range cfg.Profiles {
		prof := &cfg.Profiles[i]
		if tried[prof.Name] || !creds.Has(profile.Dir(e.paths, prof), e.paths.Home) {
			continue
		}
		if resumes(args) {
			if h := from.Share.Get("history"); h == config.Local || prof.Share.Get("history") != h {
				continue
			}
		}
		profs = append(profs, prof)
	}
	if len(profs) == 0 {
		return nil, nil
	}
	reports := reportsOf(e, cfg, profs, usage.Options{})
	if best := p.Choose(reports, clock()); best != nil {
		prof, _ := cfg.Find(best.Name)
		return prof, reports
	}
	var unknown []*usage.Report
	for _, r := range reports {
		if r.Limits == nil || !r.Limits.FiveHour.Known {
			unknown = append(unknown, r)
		}
	}
	if len(unknown) == 0 {
		return nil, reports
	}
	slices.SortFunc(unknown, func(a, b *usage.Report) int { return cmp.Compare(a.Slot, b.Slot) })
	prof, _ := cfg.Find(unknown[0].Name)
	return prof, reports
}

// tail keeps the last max bytes written to it.
type tail struct {
	buf []byte
	max int
}

func (t *tail) Write(p []byte) (int, error) {
	t.buf = append(t.buf, p...)
	if len(t.buf) > t.max {
		t.buf = t.buf[len(t.buf)-t.max:]
	}
	return len(p), nil
}
