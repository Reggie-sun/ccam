package usage

import (
	"fmt"
	"strings"
	"time"
)

// Strategies for Policy.Strategy.
const (
	Headroom = "headroom" // the lowest 5-hour usage (default)
	Expiring = "expiring" // the soonest 5-hour reset, so its quota is not left unused
)

// Policy is how `ccam next --auto`, `--fallback` and `ccam pick` choose an
// account from its usage. An account is judged by its 5-hour and weekly
// windows, and by the weekly window of the model named with --model.
type Policy struct {
	// MaxUsage is the usage (1-100) up to which an account is preferred; 0
	// means 100. Accounts at or over it are only chosen when no account is
	// below it, and accounts at 100% never.
	MaxUsage float64
	Strategy string // Headroom (also "") or Expiring
	Model    string // the --model argument, if any
}

// Check reports a policy that cannot be used, such as a typo in config.toml.
func (p Policy) Check() error {
	if p.MaxUsage < 0 || p.MaxUsage > 100 {
		return fmt.Errorf("auto_max_usage must be between 1 and 100, not %g", p.MaxUsage)
	}
	if p.Strategy != "" && p.Strategy != Headroom && p.Strategy != Expiring {
		return fmt.Errorf("auto_strategy must be %q or %q, not %q", Headroom, Expiring, p.Strategy)
	}
	return nil
}

func (p Policy) max() float64 {
	if p.MaxUsage <= 0 {
		return 100
	}
	return p.MaxUsage
}

// windows lists the known windows r is judged by.
func (p Policy) windows(r *Report) []Window {
	var out []Window
	for _, w := range []Window{r.Limits.FiveHour, r.Limits.Weekly} {
		if w.Known {
			out = append(out, w)
		}
	}
	if fam := family(p.Model); fam != "" {
		for _, s := range r.Limits.Scoped {
			if s.Window.Known && strings.Contains(strings.ToLower(s.Name), fam) {
				out = append(out, s.Window)
			}
		}
	}
	return out
}

// family is the model family a --model value names: opus, sonnet, haiku
// or fable, as aliases and full IDs both contain it.
func family(model string) string {
	m := strings.ToLower(model)
	for _, f := range []string{"opus", "sonnet", "haiku", "fable"} {
		if strings.Contains(m, f) {
			return f
		}
	}
	return ""
}

// usable reports whether r has known 5-hour usage and is not logged out.
func usable(r *Report) bool {
	return r != nil && r.Limits != nil && r.Limits.FiveHour.Known && r.Status != StatusLoggedOut
}

// peak is r's highest usage among the windows it is judged by.
func (p Policy) peak(r *Report) float64 {
	top := 0.0
	for _, w := range p.windows(r) {
		top = max(top, w.Percent)
	}
	return top
}

// Choose returns the account to use, or nil when every account is full or
// its usage is unknown. Below MaxUsage the strategy decides; when every
// account is at or over it, the one with the lowest usage is taken.
func (p Policy) Choose(reports []*Report, now time.Time) *Report {
	var free, over []*Report
	for _, r := range reports {
		if !usable(r) {
			continue
		}
		switch top := p.peak(r); {
		case top >= 100:
		case top >= p.max():
			over = append(over, r)
		default:
			free = append(free, r)
		}
	}
	if len(free) > 0 {
		return p.best(free, now)
	}
	return Policy{}.best(over, now)
}

func (p Policy) best(rs []*Report, now time.Time) *Report {
	var best *Report
	for _, r := range rs {
		if best == nil || p.before(r, best, now) {
			best = r
		}
	}
	return best
}

// before reports whether a is a better choice than b.
func (p Policy) before(a, b *Report, now time.Time) bool {
	if p.Strategy == Expiring {
		ra, rb := resetIn(a, now), resetIn(b, now)
		if ra != rb {
			return ra < rb
		}
	}
	return a.Limits.FiveHour.Percent < b.Limits.FiveHour.Percent
}

// resetIn is the time until r's 5-hour window resets; a window that is not
// running (no reset time) sorts last.
func resetIn(r *Report, now time.Time) time.Duration {
	at := r.Limits.FiveHour.ResetsAt
	if at.IsZero() {
		return 1<<63 - 1
	}
	return at.Sub(now)
}

// FreeAt returns when the first full account gets room again: the reset of
// its latest window at 100%. It is zero when no full account has a known
// reset time.
func (p Policy) FreeAt(reports []*Report, now time.Time) time.Time {
	var first time.Time
	for _, r := range reports {
		if !usable(r) {
			continue
		}
		var at time.Time
		for _, w := range p.windows(r) {
			if w.Percent >= 100 && w.ResetsAt.After(at) {
				at = w.ResetsAt
			}
		}
		if at.After(now) && (first.IsZero() || at.Before(first)) {
			first = at
		}
	}
	return first
}
