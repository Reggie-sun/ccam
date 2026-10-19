package usage

import (
	"os"
	"testing"
	"time"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestParseNewFormat(t *testing.T) {
	l, err := ParseLive(fixture(t, "new-format.json"))
	if err != nil {
		t.Fatal(err)
	}
	if l.FiveHour.Percent != 72 || l.Weekly.Percent != 41 || !l.Weekly.Known {
		t.Fatalf("%+v", l)
	}
	if !l.FiveHour.ResetsAt.Equal(time.Date(2026, 10, 18, 18, 10, 0, 299994000, time.UTC)) {
		t.Fatalf("reset %v", l.FiveHour.ResetsAt)
	}
	if len(l.Scoped) != 1 || l.Scoped[0].Name != "Fable" || l.Scoped[0].Window.Percent != 88 {
		t.Fatalf("scoped %+v", l.Scoped)
	}
}

// The real response seen on 2026-10-18 had no weekly_all and a null seven_day.
func TestParseScopedOnly(t *testing.T) {
	l, err := ParseLive(fixture(t, "scoped-only.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !l.FiveHour.Known || l.FiveHour.Percent != 0 || l.Weekly.Known || len(l.Scoped) != 1 {
		t.Fatalf("%+v", l)
	}
}

func TestParseOldFormat(t *testing.T) {
	l, err := ParseLive(fixture(t, "old-format.json"))
	if err != nil {
		t.Fatal(err)
	}
	if l.FiveHour.Percent != 100 || l.Weekly.Percent != 68.4 || len(l.Scoped) != 1 || l.Scoped[0].Name != "Opus" {
		t.Fatalf("%+v", l)
	}
}

func TestParseRejectsUnknownShape(t *testing.T) {
	if _, err := ParseLive([]byte(`{"something_else":1}`)); err == nil {
		t.Fatal("want error")
	}
	if _, err := ParseLive([]byte(`nope`)); err == nil {
		t.Fatal("want error")
	}
}

func TestParseStatusLine(t *testing.T) {
	l, model, ok, err := ParseStatusLine(fixture(t, "statusline-input.json"))
	if err != nil || !ok || model != "Opus" {
		t.Fatalf("%v %v %v", model, ok, err)
	}
	if l.FiveHour.Percent != 23.5 || l.Weekly.Percent != 41 || l.FiveHour.ResetsAt.Unix() != 1792346400 {
		t.Fatalf("%+v", l)
	}
	_, _, ok, err = ParseStatusLine([]byte(`{"model":{"display_name":"Opus"}}`))
	if ok || err != nil {
		t.Fatal("no rate_limits → not ok")
	} // No resets_at: unknown reset time, not 1970.
	l, _, ok, _ = ParseStatusLine([]byte(`{"rate_limits":{"five_hour":{"used_percentage":5}}}`))
	if !ok || !l.FiveHour.ResetsAt.IsZero() {
		t.Fatalf("missing resets_at: %+v", l.FiveHour)
	}
}
