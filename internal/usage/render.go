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

// width measures terminal columns. Box-drawing blocks are "ambiguous width";
// Korean locales would count them as 2, but terminals draw them as 1.
var width = &runewidth.Condition{EastAsianWidth: false}

// Table renders reports as aligned lines: header first, one line per
// account (plus model-specific lines with Verbose). Pending reports render
// as "조회 중…" so a live view can redraw in place.
func Table(reports []*Report, o RenderOptions) []string {
	nameW := 4
	for _, r := range reports {
		if r != nil {
			nameW = max(nameW, width.StringWidth(r.Name))
		}
	}
	lines := []string{fmt.Sprintf(" %s  %s %s  %s %s  %s",
		pad("계정", nameW+2), pad("5시간", barWidth+5), pad("초기화", 10), pad("주간", barWidth+5), pad("초기화", 10), "기준")}
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
		body = "조회 중…"
	case r.Status == StatusAPI:
		body = "한도 없음 (API)"
	case r.Status == StatusLoggedOut:
		body = fmt.Sprintf("로그인 안 됨 → ccam login %s", r.Name)
	case r.Limits == nil && r.Status == StatusExpired:
		body = fmt.Sprintf("로그인 만료, `ccs %s` 를 한 번 실행하면 갱신", r.Name)
	case r.Limits == nil:
		body = "조회 실패: " + r.Err
	default:
		l := r.Limits
		body = fmt.Sprintf("%s %s  %s %s  %s",
			gauge(l.FiveHour, o), pad(until(l.FiveHour, o.Now), 10),
			gauge(l.Weekly, o), pad(until(l.Weekly, o.Now), 10), basis(r, o.Now))
	}
	out := []string{head + "  " + body}
	if o.Verbose && r.Limits != nil {
		for _, s := range r.Limits.Scoped {
			out = append(out, fmt.Sprintf(" %s  %s  %s", pad("", nameW+2), pad("· "+s.Name+" 주간", barWidth+5+1+10),
				gauge(s.Window, o)+" "+until(s.Window, o.Now)))
		}
		if r.Source != "" {
			out = append(out, fmt.Sprintf(" %s  출처: %s", pad("", nameW+2), r.Source))
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

// Duration renders a short Korean duration: "1시간 40분", "3일 6시간", "22분".
func Duration(d time.Duration) string {
	if d <= 0 {
		return "곧"
	}
	mins := int(d.Minutes())
	days, hours, m := mins/1440, mins%1440/60, mins%60
	switch {
	case days > 0:
		return fmt.Sprintf("%d일 %d시간", days, hours)
	case hours > 0:
		return fmt.Sprintf("%d시간 %d분", hours, m)
	default:
		return fmt.Sprintf("%d분", max(m, 1))
	}
}

func basis(r *Report, now time.Time) string {
	ago := now.Sub(r.At)
	s := "방금"
	if ago >= time.Minute {
		s = Duration(ago) + " 전"
	}
	switch {
	case r.Source == "cache" || r.Source == "statusline":
		s += " (캐시)"
	case r.Status == StatusExpired:
		s += " (만료)"
	}
	return s
}

func pad(s string, w int) string {
	return width.FillRight(s, w)
}

// Best returns the account with the most 5-hour headroom, skipping full
// and unknown ones. Live values beat cached ones only through their
// percentage; both count.
func Best(reports []*Report) *Report {
	var best *Report
	for _, r := range reports {
		if r == nil || r.Limits == nil || !r.Limits.FiveHour.Known || r.Limits.FiveHour.Percent >= 100 {
			continue
		}
		if r.Status == StatusLoggedOut {
			continue
		}
		if best == nil || r.Limits.FiveHour.Percent < best.Limits.FiveHour.Percent {
			best = r
		}
	}
	return best
}
