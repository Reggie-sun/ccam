package next

import (
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
)

func cands(t0 time.Time) []Candidate {
	team := config.Share{History: "team"}
	return []Candidate{
		{Profile: &config.Profile{Name: "a", Slot: 1, Share: team}, LastUsed: t0, LoggedIn: true},
		{Profile: &config.Profile{Name: "b", Slot: 2}, LastUsed: t0.Add(-time.Hour), LoggedIn: true},
		{Profile: &config.Profile{Name: "c", Slot: 3, Share: team}, LoggedIn: true},
		{Profile: &config.Profile{Name: "d", Slot: 4, Share: team}, LoggedIn: false},
	}
}

func TestRotate(t *testing.T) {
	now := time.Now()
	from, to, err := Rotate(cands(now), false)
	if err != nil || from.Name != "a" || to.Name != "b" {
		t.Fatalf("%v %v %v", from, to, err)
	}
	// Continuing needs the same history group: skips b (local).
	_, to, _ = Rotate(cands(now), true)
	if to.Name != "c" {
		t.Fatalf("sameHistory picked %s", to.Name)
	}
	// Wraps around and skips logged-out d.
	cs := cands(now)
	cs[2].LastUsed = now.Add(time.Minute) // c most recent
	_, to, _ = Rotate(cs, false)
	if to.Name != "a" {
		t.Fatalf("wrap picked %s", to.Name)
	}
}

func TestRotateErrors(t *testing.T) {
	if _, _, err := Rotate(nil, false); err == nil {
		t.Fatal("empty")
	}
	one := []Candidate{{Profile: &config.Profile{Name: "a", Slot: 1}, LoggedIn: true}}
	if _, _, err := Rotate(one, false); err == nil {
		t.Fatal("single account has no next")
	}
	if _, _, err := Rotate(cands(time.Now())[1:2], true); err == nil {
		t.Fatal("local history cannot continue")
	}
}
