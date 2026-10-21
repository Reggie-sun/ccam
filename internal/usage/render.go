package usage

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

// RenderOptions controls the table.
type RenderOptions struct {
	Color   bool
	Verbose bool
	Now     time.Time
}

const barWidth = 20

// labelWidth fits "5-hour", "Weekly" and most model names ("Fable", "Opus").
const labelWidth = 8

// width measures terminal columns. Box-drawing blocks are "ambiguous width";
// Korean locales would count them as 2, but terminals draw them as 1.
var width = &runewidth.Condition{EastAsianWidth: false}

// Table renders reports as one block per account: a line with the slot,
// name and age of the numbers, then one bar per limit (5-hour, weekly and
// each model-specific weekly window). Pending reports render as
// "fetching…" on the account line so a live view can redraw in place.
func Table(reports []*Report, o RenderOptions) []string {
	nameW := 0
	for _, r := range reports {
		if r != nil {
			nameW = max(nameW, width.StringWidth(r.Name))
		}
	}
	var lines []string
	for _, r := range reports {
		if r == nil {
			continue
		}
		lines = append(lines, block(r, nameW, o)...)
	}
	return lines
}

func block(r *Report, nameW int, o RenderOptions) []string {
	head := fmt.Sprintf(" %d %s  ", r.Slot, pad(r.Name, nameW))
	switch {
	case r.Status == StatusPending:
		return []string{head + "fetching…"}
	case r.Status == StatusAPI:
		return []string{head + "no limits (API)"}
	case r.Status == StatusLoggedOut:
		return []string{head + fmt.Sprintf("not logged in → ccam login %s", r.Name)}
	case r.Limits == nil && r.Status == StatusExpired:
		return []string{head + fmt.Sprintf("login expired, run `ccam %s` once to refresh", r.Name)}
	case r.Limits == nil:
		return []string{head + "lookup failed: " + r.Err}
	}
	l := r.Limits
	out := []string{head + basis(r, o.Now),
		limitLine("5-hour", l.FiveHour, o),
		limitLine("Weekly", l.Weekly, o)}
	for _, s := range l.Scoped {
		out = append(out, limitLine(s.Name, s.Window, o))
	}
	if o.Verbose {
		// "(cached)" already says the numbers are saved ones; the line adds
		// where they came from, which "(cached)" does not tell.
		switch r.Source {
		case "live":
			out = append(out, "   source: live, from the usage endpoint")
		case "claude":
			out = append(out, "   source: live, from Claude Code (usage_source = \"claude\")")
		case "statusline":
			out = append(out, "   source: saved by the status line")
		case "cache":
			out = append(out, "   source: saved from the usage endpoint")
		}
	}
	return out
}

// limitLine is one window: label, bar, percent and time to reset.
func limitLine(label string, w Window, o RenderOptions) string {
	s := "   " + pad(width.Truncate(label, labelWidth, "…"), labelWidth) + " " + gauge(w, o)
	if w.Known && !w.ResetsAt.IsZero() {
		s += "  resets in " + until(w, o.Now)
	}
	return s
}

func gauge(w Window, o RenderOptions) string {
	if !w.Known {
		return "-"
	}
	filled := int(math.Round(math.Min(w.Percent, 100) / 100 * barWidth))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	s := fmt.Sprintf("%s %3.0f%%", bar, w.Percent)
	if o.Color {
		switch {
		case w.Percent >= 100:
			s = "\x1b[31m" + s + "\x1b[0m"
		case w.Percent >= 80:
			s = "\x1b[33m" + s + "\x1b[0m"
		}
	}
	return s
}

func until(w Window, now time.Time) string {
	if !w.Known || w.ResetsAt.IsZero() {
		return "-"
	}
	return Duration(w.ResetsAt.Sub(now))
}

// Duration renders a short duration: "1h 40m", "3d 6h", "22m".
func Duration(d time.Duration) string {
	if d <= 0 {
		return "soon"
	}
	mins := int(d.Minutes())
	days, hours, m := mins/1440, mins%1440/60, mins%60
	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, m)
	default:
		return fmt.Sprintf("%dm", max(m, 1))
	}
}

func basis(r *Report, now time.Time) string {
	ago := now.Sub(r.At)
	s := "just now"
	if ago >= time.Minute {
		s = Duration(ago) + " ago"
	}
	switch {
	case r.Source == "cache" || r.Source == "statusline":
		s += " (cached)"
	case r.Status == StatusExpired:
		s += " (expired)"
	}
	return s
}

func pad(s string, w int) string {
	return width.FillRight(s, w)
}
