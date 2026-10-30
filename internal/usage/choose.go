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
	// WeeklyReset takes the account whose weekly window resets soonest, so
	// its quota is not left unused; on equal resets, the one with the most
	// weekly quota left. Accounts with fewer running sessions come first.
	WeeklyReset = "weekly-reset"
	// Drain stays on one account until it is full: the one whose weekly
	// window resets soonest, running sessions or not, and ccam pick does
	// not spread sessions over the accounts. For quota about to reset.
	Drain = "drain"
	// Balance takes the account whose weekly quota would go unused first:
	// the most left per day until its reset, less for each running session.
	Balance = "balance"
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
	if p.Strategy != "" && p.Strategy != Headroom && p.Strategy != Expiring && p.Strategy != Balance && p.Strategy != WeeklyReset && p.Strategy != Drain {
		return fmt.Errorf("auto_strategy must be %q, %q, %q, %q or %q, not %q", Headroom, Expiring, Balance, WeeklyReset, Drain, p.Strategy)
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

// Full reports whether r is at 100% in a window it is judged by.
func (p Policy) Full(r *Report) bool {
	return r != nil && r.Limits != nil && p.peak(r) >= 100
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
		return p.best(p.uncrowded(free), now)
	}
	return Policy{}.best(over, now)
}

// A crowded account runs at least crowdedSessions sessions and has used more
// than crowdedUsed percent of its five-hour window: one more session there
// would run the window out for all of them.
const (
	crowdedSessions = 3
	crowdedUsed     = 50.0
)

// uncrowded drops the crowded accounts under WeeklyReset and Drain, which
// otherwise pile every session on one account, unless every account is
// crowded: then the strategy picks among them all.
func (p Policy) uncrowded(rs []*Report) []*Report {
	if p.Strategy != WeeklyReset && p.Strategy != Drain {
		return rs
	}
	var roomy []*Report
	for _, r := range rs {
		if r.Sessions < crowdedSessions || r.Limits.FiveHour.Percent <= crowdedUsed {
			roomy = append(roomy, r)
		}
	}
	if len(roomy) == 0 {
		return rs
	}
	return roomy
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
	switch p.Strategy {
	case Expiring:
		ra, rb := resetIn(a, now), resetIn(b, now)
		if ra != rb {
			return ra < rb
		}
	case WeeklyReset, Drain:
		// WeeklyReset puts an account already running a session after idle
		// ones (how many it runs does not count); Drain keeps using it.
		if ia, ib := a.Sessions == 0, b.Sessions == 0; p.Strategy == WeeklyReset && ia != ib {
			return ia
		}
		ra, rb := weeklyResetIn(a, now), weeklyResetIn(b, now)
		if ra != rb {
			return ra < rb
		}
		if pa, pb := weeklyUsed(a), weeklyUsed(b); pa != pb {
			return pa < pb
		}
	case Balance:
		ua, ub := urgency(a, now), urgency(b, now)
		if ua != ub {
			return ua > ub
		}
	}
	return a.Limits.FiveHour.Percent < b.Limits.FiveHour.Percent
}

// sessionCost is how much of the weekly quota per day a running session
// is taken to use, so balance spreads sessions over the accounts.
const sessionCost = 10

// urgency is the share of r's weekly quota left per day until the week
// resets, less sessionCost for each running session. Without a known week,
// a full week is assumed.
func urgency(r *Report, now time.Time) float64 {
	w := r.Limits.Weekly
	left, days := 100.0, 7.0
	if w.Known {
		left = 100 - w.Percent
		if !w.ResetsAt.IsZero() {
			days = max(w.ResetsAt.Sub(now).Hours()/24, 1.0/24)
		}
	}
	return left/days - sessionCost*float64(r.Sessions)
}

// Why says in a few words why Choose took r.
func (p Policy) Why(r *Report, now time.Time) string {
	if r == nil || !usable(r) {
		return ""
	}
	five := fmt.Sprintf("%.0f%% of its 5 hours used", r.Limits.FiveHour.Percent)
	if top := p.peak(r); top >= p.max() {
		return fmt.Sprintf("every account is at %.0f%% or more; it has the lowest use, %s", p.max(), five)
	}
	switch p.Strategy {
	case Expiring:
		if at := r.Limits.FiveHour.ResetsAt; !at.IsZero() {
			return fmt.Sprintf("its 5-hour window resets first, in %s (%s)", Duration(at.Sub(now)), five)
		}
	case WeeklyReset, Drain:
		w := r.Limits.Weekly
		if w.Known && !w.ResetsAt.IsZero() {
			return fmt.Sprintf("its week resets first, in %s (%.0f%% of the week left)", Duration(w.ResetsAt.Sub(now)), 100-w.Percent)
		}
	case Balance:
		w := r.Limits.Weekly
		why := "its weekly use is unknown"
		if w.Known {
			why = fmt.Sprintf("%.0f%% of its week is left", 100-w.Percent)
			if !w.ResetsAt.IsZero() {
				why += ", resetting in " + Duration(w.ResetsAt.Sub(now))
			}
		}
		if r.Sessions > 0 {
			why += fmt.Sprintf("; %d running", r.Sessions)
		}
		return why
	}
	return five + ", the least"
}

// Use says in a few words how far r is by what p goes by, as Why does but
// without comparing it with other accounts: for an account taken while the
// one Choose takes was passed over.
func (p Policy) Use(r *Report, now time.Time) string {
	if r == nil || !usable(r) {
		return ""
	}
	five := fmt.Sprintf("%.0f%% of its 5 hours used", r.Limits.FiveHour.Percent)
	if p.peak(r) < p.max() {
		switch p.Strategy {
		case Expiring:
			if at := r.Limits.FiveHour.ResetsAt; !at.IsZero() {
				return fmt.Sprintf("its 5-hour window resets in %s (%s)", Duration(at.Sub(now)), five)
			}
		case WeeklyReset, Drain, Balance:
			return p.Why(r, now) // it names no other account
		}
	}
	return five
}

// weeklyResetIn is the time until r's weekly window resets; an unknown week
// sorts last.
func weeklyResetIn(r *Report, now time.Time) time.Duration {
	w := r.Limits.Weekly
	if !w.Known || w.ResetsAt.IsZero() {
		return 1<<63 - 1
	}
	return w.ResetsAt.Sub(now)
}

// weeklyUsed is r's weekly usage; an unknown week counts as untouched.
func weeklyUsed(r *Report) float64 {
	if !r.Limits.Weekly.Known {
		return 0
	}
	return r.Limits.Weekly.Percent
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
