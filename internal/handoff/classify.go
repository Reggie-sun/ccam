package handoff

import (
	"strings"
	"time"
	"unicode/utf8"
)

// Kinds of usage limit, from Classify.
const (
	Account = "account" // the account's plan: a 5-hour or weekly window, or its usage limit as a whole
	Model   = "model"   // one model's window or allowance (Opus, Sonnet, Fable), or a model that needs usage credits
	Credits = "credits" // usage credits or a spend limit (the account's, its team's, organization's or channel's), or a seat without credits
	Org     = "org"     // the organization stopped the account: out of usage, a group set to $0, an allocation or the service turned off
)

// Classify tells which kind of usage limit stopped a turn from the message
// Claude Code (2.1.306) stopped it with: Account, Model, Credits or Org.
// It is "" for anything else, such as the server's short throttle ("Server
// is temporarily limiting requests") or an API key's rate limit ("Request
// rejected (429)"), which are no usage limits. The messages are those
// Claude Code itself counts as a usage limit.
func Classify(msg string) string {
	m := strings.TrimSpace(strings.ReplaceAll(msg, "’", "'"))
	if rest, ok := strings.CutPrefix(m, "You've hit your "); ok {
		return hitYour(rest)
	}
	switch {
	case strings.HasPrefix(m, "You've reached your "), strings.HasPrefix(m, "Fable 5 requires usage credits"), needsCredits(m):
		return Model
	case hasPrefix(m, "You're out of usage credits", "You're out of extra usage", "Your seat type doesn't include"):
		return Credits
	case hasPrefix(m, "Your org is out of usage", "Your group's usage limit is set to $0",
		"Your usage allocation has been disabled by your admin", "This service is disabled for your org"):
		return Org
	case strings.HasPrefix(m, "Claude AI usage limit reached"): // Claude Code before 1.0
		return Account
	}
	return ""
}

// hitYour classifies "You've hit your <rest>": the plan's limit (session,
// weekly, usage or just "limit"), a model's ("Opus limit"), or a spend
// limit or budget of the account, its team, organization or channel.
func hitYour(rest string) string {
	name, _, _ := strings.Cut(rest, " ·")
	if strings.Contains(name, "spend") || strings.Contains(name, "credit") || strings.Contains(name, "budget") ||
		hasPrefix(name, "individual ", "monthly ", "org's ", "team's ", "channel's ") {
		return Credits
	}
	// Window names are lowercase ("session limit"); model names are not.
	if first, _, ok := strings.Cut(name, " limit"); ok && first != "" && !strings.Contains(first, " ") {
		if r, _ := utf8.DecodeRuneInString(first); r >= 'A' && r <= 'Z' {
			return Model
		}
	}
	return Account
}

// needsCredits matches Claude Code's "<model> requires usage credits." (as
// in "Fable 5 requires usage credits." and "Opus 5.5 (1M context) requires
// usage credits."): a model family, then up to 40 characters without "·".
func needsCredits(m string) bool {
	head, _, ok := strings.Cut(m, " requires usage credits.")
	if !ok {
		return false
	}
	for _, family := range []string{"Fable", "Haiku", "Mythos", "Opus", "Sonnet"} {
		if rest, ok := strings.CutPrefix(head, family); ok {
			if rest == "" {
				return true
			}
			v, ok := strings.CutPrefix(rest, " ")
			n := utf8.RuneCountInString(v)
			return ok && n >= 1 && n <= 40 && !strings.ContainsAny(v, "·\n")
		}
	}
	return false
}

func hasPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// LimitKind names which limit stopped a session for its note, from Claude
// Code's message: session, weekly or usage (the plan's limit), model, spend
// (credits or a spend limit) or org; "" when Classify finds no usage limit.
func LimitKind(msg string) string {
	switch Classify(msg) {
	case Account:
		rest := strings.TrimPrefix(strings.ReplaceAll(strings.TrimSpace(msg), "’", "'"), "You've hit your ")
		switch {
		case strings.HasPrefix(rest, "session limit"):
			return "session"
		case strings.HasPrefix(rest, "weekly limit"):
			return "weekly"
		}
		return "usage"
	case Model:
		return "model"
	case Credits:
		return "spend"
	case Org:
		return "org"
	}
	return ""
}

// FormatReset writes a reset time as Claude Code does, in the local time
// zone: 3:45pm or 9am within a day of now, else Oct 24, 9am, with the year
// when it is not now's.
func FormatReset(at, now time.Time) string {
	at = at.Local()
	clock := at.Format("3:04pm")
	if at.Minute() == 0 {
		clock = at.Format("3pm")
	}
	if at.Sub(now) <= 24*time.Hour {
		return clock
	}
	day := at.Format("Jan 2")
	if at.Year() != now.Local().Year() {
		day += at.Format(", 2006")
	}
	return day + ", " + clock
}
