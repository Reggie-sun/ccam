package handoff

import (
	"testing"
	"time"
)

// The messages Claude Code 2.1.306 stops a turn with (its usage-limit
// prefixes, the window names and the "requires usage credits" form), and
// the 429s that are not usage limits.
func TestClassify(t *testing.T) {
	cases := map[string]string{
		"You've hit your limit · resets 3pm":                                                        Account,
		"You've hit your session limit · resets 3:45pm (Asia/Seoul)":                                Account,
		"You’ve hit your weekly limit · resets Oct 23, 9am (Asia/Seoul)":                            Account,
		"You've hit your usage limit · resets 3pm":                                                  Account,
		"You've hit your usage limit · contact your admin to increase it":                           Account,
		"You've hit your session limit · resets 3:45pm · progress saved":                            Account,
		"Claude AI usage limit reached|1760000000":                                                  Account,
		"You've hit your Opus limit · resets Oct 24, 9am (Asia/Seoul)":                              Model,
		"You've hit your Sonnet limit · resets Oct 24, 9am":                                         Model,
		"You've hit your Fable limit · resets Oct 24, 9am":                                          Model,
		"You've reached your Fable limit. Switch to another model to continue.":                     Model,
		"Fable 5 requires usage credits · /model to set it up":                                      Model,
		"Opus 5.5 (1M context) requires usage credits. Switch to another model to continue.":        Model,
		"You've hit your monthly spend limit · raise it at https://claude.ai/settings/usage":        Credits,
		"You've hit your individual spend limit · ask your admin for a higher limit":                Credits,
		"You've hit your usage credit limit":                                                        Credits,
		"You're out of usage credits · resets 3pm":                                                  Credits,
		"You're out of extra usage":                                                                 Credits,
		"Your seat type doesn't include usage credits":                                              Credits,
		"You've hit your org's monthly spend limit · ask your admin to raise it":                    Credits,
		"You've hit your team's shared budget · raise it at https://claude.ai/admin-settings/usage": Credits,
		"You've hit your channel's monthly spend limit.":                                            Credits,
		"You've hit your org's monthly usage limit · resets Nov 1, 12am":                            Credits,
		"You've hit your individual usage limit":                                                    Credits,
		"Your org is out of usage · contact your admin":                                             Org,
		"Your group's usage limit is set to $0 · ask your admin for a higher limit":                 Org,
		"Your usage allocation has been disabled by your admin · ask your admin for a higher limit": Org,
		"This service is disabled for your org":                                                     Org,
		"API Error: Server is temporarily limiting requests (not your usage limit)":                 "",
		"API Error: Request rejected (429) · Number of request tokens has exceeded your rate limit": "",
		"You've used 92% of your session limit · resets 3pm":                                        "",
		"You're now using usage credits · Your session limit resets 3pm":                            "",
		"": "",
	}
	for in, want := range cases {
		if got := Classify(in); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
}

func TestLimitKindWrapsClassify(t *testing.T) {
	cases := map[string]string{
		"You've hit your session limit · resets 3:45pm":                             "session",
		"You’ve hit your weekly limit · resets Oct 23, 9am":                         "weekly",
		"You've hit your limit · resets 3pm":                                        "usage",
		"You've hit your usage limit · resets 3pm":                                  "usage",
		"You've hit your Opus limit · resets Oct 23":                                "model",
		"You've reached your Fable limit.":                                          "model",
		"You're out of usage credits":                                               "spend",
		"Your org is out of usage · contact your admin":                             "org",
		"API Error: Server is temporarily limiting requests (not your usage limit)": "",
	}
	for in, want := range cases {
		if got := LimitKind(in); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
	// The zone Claude Code adds is the local one, so it is left out.
	if got := Reset("You've hit your session limit · resets 3:45pm (Asia/Seoul)"); got != "3:45pm" {
		t.Fatal(got)
	}
	if got := Reset("You've hit your weekly limit · resets Oct 24, 9am (Asia/Seoul) · progress saved"); got != "Oct 24, 9am" {
		t.Fatal(got)
	}
}

// Reset times read as in Claude Code: the time within a day, else the day
// and time, with the year only when it is another year.
func TestFormatReset(t *testing.T) {
	now := time.Date(2026, 10, 22, 13, 0, 0, 0, time.Local)
	cases := map[time.Time]string{
		now.Add(2*time.Hour + 45*time.Minute):              "3:45pm",
		time.Date(2026, 10, 22, 21, 0, 0, 0, time.Local):   "9pm",
		time.Date(2026, 10, 23, 9, 30, 0, 0, time.Local):   "9:30am",
		time.Date(2026, 10, 24, 9, 0, 0, 0, time.Local):    "Oct 24, 9am",
		time.Date(2026, 10, 28, 18, 15, 0, 0, time.Local):  "Oct 28, 6:15pm",
		time.Date(2027, 1, 2, 0, 0, 0, 0, time.Local):      "Jan 2, 2027, 12am",
		now.Add(24 * time.Hour).In(time.FixedZone("x", 0)): now.Add(24 * time.Hour).Format("3pm"),
	}
	for at, want := range cases {
		if got := FormatReset(at, now); got != want {
			t.Errorf("%v: %q, want %q", at, got, want)
		}
	}
}
