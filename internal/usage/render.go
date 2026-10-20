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

const barWidth = 8

// resetWidth fits the longest Duration below a week, "23h 59m".
const resetWidth = 8

// width measures terminal columns. Box-drawing blocks are "ambiguous width";
// Korean locales would count them as 2, but terminals draw them as 1.
var width = &runewidth.Condition{EastAsianWidth: false}

// Table renders reports as aligned lines: header first, one line per
// account (plus model-specific lines with Verbose). Pending reports render
// as "fetching…" so a live view can redraw in place.
func Table(reports []*Report, o RenderOptions) []string {
	nameW := 5 // fits the "Account" header with the slot column
	for _, r := range reports {
		if r != nil {
			nameW = max(nameW, width.StringWidth(r.Name))
		}
	}
	lines := []string{fmt.Sprintf(" %s  %s %s  %s %s  %s",
		pad("Account", nameW+2), pad("5-hour", barWidth+5), pad("Resets", resetWidth), pad("Weekly", barWidth+5), pad("Resets", resetWidth), "As of")}
	for _, r := range reports {
		if r == nil {
			continue
		}
		lines = append(lines, row(r, nameW, o)...)
	}
	return lines
}

func row(r *Report, nameW int, o RenderOptions) []string {
	head := fmt.Sprintf(" %d %s", r.Slot, pad(r.Name, nameW))
	var body string
	switch {
	case r.Status == StatusPending:
		body = "fetching…"
	case r.Status == StatusAPI:
		body = "no limits (API)"
	case r.Status == StatusLoggedOut:
		body = fmt.Sprintf("not logged in → ccam login %s", r.Name)
	case r.Limits == nil && r.Status == StatusExpired:
		body = fmt.Sprintf("login expired, run `ccam %s` once to refresh", r.Name)
	case r.Limits == nil:
		body = "lookup failed: " + r.Err
	default:
		l := r.Limits
		body = fmt.Sprintf("%s %s  %s %s  %s",
			gauge(l.FiveHour, o), pad(until(l.FiveHour, o.Now), resetWidth),
			gauge(l.Weekly, o), pad(until(l.Weekly, o.Now), resetWidth), basis(r, o.Now))
	}
	out := []string{head + "  " + body}
	if o.Verbose && r.Limits != nil {
		labelW := barWidth + 5 + 1 + resetWidth // the 5-hour and its reset column
		for _, s := range r.Limits.Scoped {
			label := width.Truncate("· "+s.Name+" weekly", labelW, "…")
			out = append(out, fmt.Sprintf(" %s  %s  %s", pad("", nameW+2), pad(label, labelW),
				gauge(s.Window, o)+" "+until(s.Window, o.Now)))
		}
		// "(cached)" already says the numbers are saved ones; the line adds
		// where they came from, which "(cached)" does not tell.
		switch r.Source {
		case "live":
			out = append(out, fmt.Sprintf(" %s  source: live, from the usage endpoint", pad("", nameW+2)))
		case "claude":
			out = append(out, fmt.Sprintf(" %s  source: live, from Claude Code (usage_source = \"claude\")", pad("", nameW+2)))
		case "statusline":
			out = append(out, fmt.Sprintf(" %s  source: saved by the status line", pad("", nameW+2)))
		case "cache":
			out = append(out, fmt.Sprintf(" %s  source: saved from the usage endpoint", pad("", nameW+2)))
		}
	}
	return out
}

func gauge(w Window, o RenderOptions) string {
	if !w.Known {
		return pad("-", barWidth+5)
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
