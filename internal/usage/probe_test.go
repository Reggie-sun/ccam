package usage

import (
	"strings"
	"testing"
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
	old := `{"type":"control_response","response":{"subtype":"error","request_id":"ccam","error":"Unknown subtype: get_usage"}}`
	if _, _, err := ParseProbe([]byte(old)); err == nil || !strings.Contains(err.Error(), "Unknown subtype") {
		t.Fatal(err)
	}
	if _, _, err := ParseProbe([]byte("Error: not logged in\n")); err == nil || !strings.Contains(err.Error(), "2.1.270") {
		t.Fatal(err)
	}
}
