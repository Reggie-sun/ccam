package usage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseProbe(t *testing.T) {
	ok := strings.Join([]string{
		`{"type":"system","subtype":"init"}`,
		`{"type":"control_response","response":{"subtype":"success","request_id":"other","response":{}}}`,
		`{"type":"control_response","response":{"subtype":"success","request_id":"ccam","response":{"subscription_type":"max","rate_limits_available":true,` +
			`"rate_limits":{"five_hour":{"utilization":42,"resets_at":"2026-10-20T15:00:00Z"},"seven_day":{"utilization":7.5,"resets_at":"2026-10-24T09:00:00Z"},` +
			`"seven_day_opus":{"utilization":90,"resets_at":"2026-10-24T09:00:00Z"},"seven_day_sonnet":null,` +
			`"model_scoped":[{"display_name":"Fable","utilization":12,"resets_at":"2026-10-24T09:00:00Z"},{"weird":true}]}}}}`,
	}, "\n")
	l, got, err := ParseProbe([]byte(ok))
	if err != nil || !got || l.FiveHour.Percent != 42 || !l.FiveHour.Known || l.Weekly.Percent != 7.5 || l.FiveHour.ResetsAt.IsZero() {
		t.Fatalf("%+v %v %v", l, got, err)
	}
	if len(l.Scoped) != 2 || l.Scoped[0].Name != "Opus" || l.Scoped[1].Name != "Fable" {
		t.Fatalf("scoped %+v", l.Scoped)
	}
	none := `{"type":"control_response","response":{"subtype":"success","request_id":"ccam","response":{"rate_limits_available":false,"rate_limits":null}}}`
	if _, got, err := ParseProbe([]byte(none)); got || err != nil {
		t.Fatalf("no plan limits: %v %v", got, err)
	}
	if _, got, err := ParseProbe([]byte(strings.Replace(none, `"rate_limits_available":false,`, "", 1))); got || err != nil {
		t.Fatalf("no plan limits, without saying: %v %v", got, err)
	}
	old := `{"type":"control_response","response":{"subtype":"error","request_id":"ccam","error":"Unknown subtype: get_usage"}}`
	if _, _, err := ParseProbe([]byte(old)); err == nil || !strings.Contains(err.Error(), "Unknown subtype") {
		t.Fatal(err)
	}
	if _, _, err := ParseProbe([]byte("Error: not logged in\n")); err == nil || !strings.Contains(err.Error(), "2.1.270") {
		t.Fatal(err)
	}
}

// Claude Code 2.1.306 answers rate_limits null with rate_limits_available
// true when its own lookup failed (the endpoint refused or did not answer):
// that is an error, not a login without plan limits.
func TestProbeLookupFailedInClaudeCode(t *testing.T) {
	answer := `{"type":"control_response","response":{"subtype":"success","request_id":"ccam","response":{"subscription_type":"team","rate_limits_available":true,"rate_limits":null}}}`
	_, got, err := ParseProbe([]byte(answer))
	if got || err == nil || !strings.Contains(err.Error(), "gave Claude Code no usage") || errors.Is(err, errFormat) {
		t.Fatalf("%v %v", got, err)
	}
	e := newEnv(t)
	e.f.Source = "claude"
	e.f.Probe = func(context.Context, Account) ([]byte, error) { return []byte(answer), nil }
	r := e.f.Fetch(context.Background(), e.account("a", "ok-new-format.json", future), Options{})
	if r.Status != StatusError || strings.Contains(r.Err, "no plan limits") || !strings.Contains(r.Err, "gave Claude Code no usage") || r.FormatChanged {
		t.Fatalf("%+v", r)
	}
}

// Claude Code renews a login whose access token has expired, or soon
// expires, as it starts; cutting the probe off then may lose the login, so
// such a probe gets the time Claude Code's renewal may take. A login Claude
// Code has cleared starts no probe.
func TestProbeWaitsForAClaudeCodeRenewal(t *testing.T) {
	e := newEnv(t)
	e.f.Source = "claude"
	var left []time.Duration
	e.f.Probe = func(ctx context.Context, a Account) ([]byte, error) {
		d, _ := ctx.Deadline()
		left = append(left, time.Until(d))
		return nil, errors.New("no answer")
	}
	now := time.Now()
	for _, expires := range []time.Time{now.Add(time.Hour), now.Add(-time.Hour), now.Add(2 * time.Minute)} {
		e.f.Fetch(context.Background(), e.account(fmt.Sprint(expires.Unix()), "ok-new-format.json", expires.UnixMilli()), Options{NoCache: true})
	}
	if len(left) != 3 || left[0] > probeTimeout || left[1] < 40*time.Second || left[2] < 40*time.Second {
		t.Fatalf("time left: %v", left)
	}
	cleared := e.account("cleared", "", 0)
	os.WriteFile(filepath.Join(cleared.Dir, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"","refreshToken":"","expiresAt":0}}`), 0o600)
	// It expired, as doctor and the live source say: it is not "not logged in".
	if r := e.f.Fetch(context.Background(), cleared, Options{}); r.Status != StatusExpired || r.Err != "login expired → `ccam login cleared`" || len(left) != 3 {
		t.Fatalf("%+v %d", r, len(left))
	}
}
