package next

import (
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
)

func cands(t0 time.Time) []Candidate {
	return []Candidate{
		{Profile: &config.Profile{Name: "a", Slot: 1}, LastUsed: t0, LoggedIn: true},
		{Profile: &config.Profile{Name: "b", Slot: 2}, LastUsed: t0.Add(-time.Hour), LoggedIn: true},
		{Profile: &config.Profile{Name: "c", Slot: 3}, LoggedIn: true},
		{Profile: &config.Profile{Name: "d", Slot: 4}, LoggedIn: false},
	}
}

func TestRotate(t *testing.T) {
	now := time.Now()
	from, to, err := Rotate(cands(now), "", nil)
	if err != nil || from.Name != "a" || to.Name != "b" {
		t.Fatalf("%v %v %v", from, to, err)
	}
	// fits leaves b out.
	_, to, _ = Rotate(cands(now), "", func(p *config.Profile) bool { return p.Name != "b" })
	if to.Name != "c" {
		t.Fatalf("fits: picked %s", to.Name)
	}
	// From a given account; wraps around and skips logged-out d.
	from, to, _ = Rotate(cands(now), "c", nil)
	if from.Name != "c" || to.Name != "a" {
		t.Fatalf("from c: %s → %s", from.Name, to.Name)
	}
	cs := cands(now)
	cs[2].LastUsed = now.Add(time.Minute) // c most recent
	_, to, _ = Rotate(cs, "", nil)
	if to.Name != "a" {
		t.Fatalf("wrap picked %s", to.Name)
	}
}

func TestRotateErrors(t *testing.T) {
	if _, _, err := Rotate(nil, "", nil); err == nil {
		t.Fatal("empty")
	}
	one := []Candidate{{Profile: &config.Profile{Name: "a", Slot: 1}, LoggedIn: true}}
	if _, _, err := Rotate(one, "", nil); err == nil {
		t.Fatal("single account has no next")
	}
}
