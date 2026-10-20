package usage

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jhyunwoo/ccam/internal/creds"
)

// ProbeRequest asks Claude Code itself for the account's plan usage, with
// `claude -p --input-format stream-json --output-format stream-json
// --verbose`: the get_usage control request (Claude Code 2.1.270 and
// later; marked experimental there). Claude Code reads its own login, so
// ccam never touches the token.
const ProbeRequest = `{"type":"control_request","request_id":"ccam","request":{"subtype":"get_usage","skip_behaviors":true}}` + "\n"

// probeTimeout bounds one account's probe: a whole Claude Code start, then
// its own request to the usage endpoint.
const probeTimeout = 20 * time.Second

// ParseProbe reads Claude Code's answer to ProbeRequest. ok is false, with
// no error, when the login has no plan limits.
func ParseProbe(out []byte) (l Limits, ok bool, err error) {
	for _, line := range bytes.Split(out, []byte("\n")) {
		var m struct {
			Type     string `json:"type"`
			Response struct {
				Subtype   string `json:"subtype"`
				RequestID string `json:"request_id"`
				Error     string `json:"error"`
				Response  *struct {
					Available  *bool `json:"rate_limits_available"`
					RateLimits *struct {
						FiveHour    *apiWindow        `json:"five_hour"`
						SevenDay    *apiWindow        `json:"seven_day"`
						Opus        *apiWindow        `json:"seven_day_opus"`
						Sonnet      *apiWindow        `json:"seven_day_sonnet"`
						ModelScoped []json.RawMessage `json:"model_scoped"`
					} `json:"rate_limits"`
				} `json:"response"`
			} `json:"response"`
		}
		if json.Unmarshal(line, &m) != nil || m.Type != "control_response" || m.Response.RequestID != "ccam" {
			continue
		}
		if m.Response.Subtype != "success" {
			return Limits{}, false, fmt.Errorf("the usage request failed in Claude Code: %s", cmp.Or(m.Response.Error, m.Response.Subtype))
		}
		resp := m.Response.Response
		if resp == nil || resp.RateLimits == nil || resp.Available != nil && !*resp.Available {
			return Limits{}, false, nil
		}
		rl := resp.RateLimits
		l = Limits{FiveHour: fromOld(rl.FiveHour), Weekly: fromOld(rl.SevenDay)}
		for _, s := range []struct {
			name string
			w    *apiWindow
		}{{"Opus", rl.Opus}, {"Sonnet", rl.Sonnet}} {
			if w := fromOld(s.w); w.Known {
				l.Scoped = append(l.Scoped, Scoped{Name: s.name, Window: w})
			}
		}
		for _, raw := range rl.ModelScoped {
			var s struct {
				Name        string `json:"name"`
				DisplayName string `json:"display_name"`
				apiWindow
			}
			if json.Unmarshal(raw, &s) == nil {
				if w := fromOld(&s.apiWindow); w.Known && cmp.Or(s.DisplayName, s.Name) != "" {
					l.Scoped = append(l.Scoped, Scoped{Name: cmp.Or(s.DisplayName, s.Name), Window: w})
				}
			}
		}
		if !l.FiveHour.Known && !l.Weekly.Known && len(l.Scoped) == 0 {
			return Limits{}, false, nil
		}
		return l, true, nil
	}
	return Limits{}, false, fmt.Errorf("%w: Claude Code did not answer the usage request (it needs Claude Code 2.1.270 or later)", errFormat)
}

// fetchClaude is Fetch with Claude Code doing the lookup (usage_source =
// "claude"). It takes longer (a Claude Code start per account) but leaves
// the login to Claude Code alone.
func (f *Fetcher) fetchClaude(ctx context.Context, a Account, opt Options) Report {
	r := Report{Name: a.Name, Slot: a.Slot}
	if !opt.NoCache {
		if c, ok := f.readCache(a.Name, "claude"); ok && f.Now().Sub(c.At) < CacheTTL {
			r.Status, r.Limits, r.Source, r.At = StatusOK, c.Limits, "claude", c.At
			return r
		}
	}
	if !creds.Has(a.Dir, f.Home) {
		r.Status = StatusLoggedOut
		if hasAPIKey(a.Dir) {
			r.Status = StatusAPI
		}
		return r
	}
	if f.Probe == nil {
		return f.fallback(a, r, StatusError, "no way to run Claude Code for the usage")
	}
	ctx, cancel := context.WithTimeout(ctx, probeTimeout)
	defer cancel()
	out, err := f.Probe(ctx, a)
	l, ok, perr := ParseProbe(out)
	switch {
	case perr != nil && err != nil:
		return f.fallback(a, r, StatusError, err.Error())
	case perr != nil:
		r.FormatChanged = errors.Is(perr, errFormat)
		return f.fallback(a, r, StatusError, perr.Error())
	case !ok:
		return f.fallback(a, r, StatusError, "Claude Code reports no plan limits for this login")
	}
	r.Status, r.Limits, r.Source, r.At = StatusOK, &l, "claude", f.Now()
	f.writeCache(a.Name, "claude", cached{At: r.At, Limits: &l})
	return r
}
