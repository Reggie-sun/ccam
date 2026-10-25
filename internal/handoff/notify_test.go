package handoff

import "testing"

// A notification is the escape sequence the terminal understands, with no
// control character of the message's that could end it early.
func TestNotification(t *testing.T) {
	for _, c := range []struct {
		kitty, program, term string
		want                 string
	}{
		{"1", "", "xterm-kitty", "\x1b]99;;ccam: a b\x1b\\"},
		{"", "ghostty", "xterm-ghostty", "\x1b]777;notify;ccam;ccam: a b\x07"},
		{"", "WarpTerminal", "xterm-256color", "\x1b]777;notify;ccam;ccam: a b\x07"},
		{"", "", "rxvt-unicode-256color", "\x1b]777;notify;ccam;ccam: a b\x07"},
		{"", "iTerm.app", "xterm-256color", "\x1b]9;ccam: a b\x07"},
		{"", "", "xterm-256color", "\x1b]9;ccam: a b\x07"},
	} {
		t.Setenv("KITTY_WINDOW_ID", c.kitty)
		t.Setenv("TERM_PROGRAM", c.program)
		t.Setenv("TERM", c.term)
		if got := Notification("ccam: a\x07b"); got != c.want {
			t.Errorf("%+v: %q", c, got)
		}
	}
}

func TestHitLimitAndWindowLabel(t *testing.T) {
	for window, want := range map[string]string{
		"five_hour": "5-hour", "seven_day": "weekly", "seven_day_opus": "Opus", "seven_day_sonnet": "Sonnet",
		"seven_day_overage_included": "Fable", "overage": "", "": "",
	} {
		if got := WindowLabel(window); got != want {
			t.Errorf("%q: %q, want %q", window, got, want)
		}
	}
	if got := HitLimit("p1", "5-hour", "3:45pm"); got != "ccam: p1 hit its 5-hour limit (resets 3:45pm)" {
		t.Error(got)
	}
	if got := HitLimit("p1", "weekly", ""); got != "ccam: p1 hit its weekly limit" {
		t.Error(got)
	}
}
