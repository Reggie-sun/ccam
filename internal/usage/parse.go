// Package usage reads each account's 5-hour and weekly limits.
//
// All response-format knowledge lives in this file, so a format change in
// the undocumented endpoint is a one-file fix (docs/claude-layout.md).
package usage

import (
	"encoding/json"
	"fmt"
	"time"
)

// Window is one rate-limit window.
type Window struct {
	Percent  float64   `json:"percent"`
	ResetsAt time.Time `json:"resets_at"`
	Known    bool      `json:"known"`
}

// Scoped is a model-specific weekly window (shown with -v).
type Scoped struct {
	Name   string `json:"name"`
	Window Window `json:"window"`
}

// Limits is what ccam shows for one account.
type Limits struct {
	FiveHour Window   `json:"five_hour"`
	Weekly   Window   `json:"weekly"`
	Scoped   []Scoped `json:"scoped,omitempty"`
}

type apiWindow struct {
	Utilization *float64 `json:"utilization"`
	ResetsAt    string   `json:"resets_at"`
}

type apiLimit struct {
	Kind     string  `json:"kind"`
	Percent  float64 `json:"percent"`
	ResetsAt string  `json:"resets_at"`
	Scope    *struct {
		Model *struct {
			ID          *string `json:"id"`
			DisplayName string  `json:"display_name"`
		} `json:"model"`
	} `json:"scope"`
}

type apiResponse struct {
	FiveHour     *apiWindow `json:"five_hour"`
	SevenDay     *apiWindow `json:"seven_day"`
	SevenDayOpus *apiWindow `json:"seven_day_opus"`
	SevenDaySon  *apiWindow `json:"seven_day_sonnet"`
	Limits       []apiLimit `json:"limits"`
}

// ParseLive reads the /api/oauth/usage body. The new `limits` array wins;
// the old five_hour/seven_day objects fill whatever it does not cover.
func ParseLive(body []byte) (Limits, error) {
	var r apiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return Limits{}, fmt.Errorf("could not read the usage response: %w", err)
	}
	var l Limits
	for _, lim := range r.Limits {
		w := Window{Percent: lim.Percent, ResetsAt: parseTime(lim.ResetsAt), Known: true}
		switch lim.Kind {
		case "session":
			l.FiveHour = w
		case "weekly_all":
			l.Weekly = w
		case "weekly_scoped":
			name := "per-model"
			if lim.Scope != nil && lim.Scope.Model != nil {
				name = lim.Scope.Model.DisplayName
			}
			l.Scoped = append(l.Scoped, Scoped{Name: name, Window: w})
		}
	}
	if !l.FiveHour.Known {
		l.FiveHour = fromOld(r.FiveHour)
	}
	if !l.Weekly.Known {
		l.Weekly = fromOld(r.SevenDay)
	}
	if len(r.Limits) == 0 {
		for name, w := range map[string]*apiWindow{"Opus": r.SevenDayOpus, "Sonnet": r.SevenDaySon} {
			if ow := fromOld(w); ow.Known {
				l.Scoped = append(l.Scoped, Scoped{Name: name, Window: ow})
			}
		}
	}
	if !l.FiveHour.Known && !l.Weekly.Known && len(l.Scoped) == 0 {
		return l, fmt.Errorf("the usage response has no limits (the format may have changed)")
	}
	return l, nil
}

func fromOld(w *apiWindow) Window {
	if w == nil || w.Utilization == nil {
		return Window{}
	}
	return Window{Percent: *w.Utilization, ResetsAt: parseTime(w.ResetsAt), Known: true}
}

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// statusInput is the part of Claude Code's status line JSON we read.
type statusInput struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	RateLimits *struct {
		FiveHour *statusWindow `json:"five_hour"`
		SevenDay *statusWindow `json:"seven_day"`
	} `json:"rate_limits"`
}

type statusWindow struct {
	UsedPercentage *float64 `json:"used_percentage"`
	ResetsAt       int64    `json:"resets_at"` // Unix seconds
}

func (w *statusWindow) window() Window {
	if w == nil || w.UsedPercentage == nil {
		return Window{}
	}
	return Window{Percent: *w.UsedPercentage, ResetsAt: time.Unix(w.ResetsAt, 0), Known: true}
}

// ParseStatusLine reads the JSON Claude Code pipes to the status line
// command. ok is false when it carries no rate limits (API accounts, or
// before the session's first response).
func ParseStatusLine(body []byte) (l Limits, model string, ok bool, err error) {
	var in statusInput
	if err := json.Unmarshal(body, &in); err != nil {
		return Limits{}, "", false, err
	}
	model = in.Model.DisplayName
	if in.RateLimits == nil {
		return Limits{}, model, false, nil
	}
	l = Limits{FiveHour: in.RateLimits.FiveHour.window(), Weekly: in.RateLimits.SevenDay.window()}
	return l, model, l.FiveHour.Known || l.Weekly.Known, nil
}
