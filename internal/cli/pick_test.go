package cli

import (
	"encoding/json"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestPickSpreadsWorkersOverAccounts(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 30, "work": 10, "side": 50})
	var got []string
	for range 4 {
		out, errOut, code := c.run("", "pick")
		if code != 0 {
			t.Fatalf("%d %s", code, errOut)
		}
		got = append(got, strings.TrimSpace(out))
	}
	// Lowest usage first, then the next ones; when every account was
	// picked recently, the lowest again.
	if want := []string{"work", "personal", "side", "work"}; !slices.Equal(got, want) {
		t.Fatalf("picked %v, want %v", got, want)
	}
}

func TestPickConcurrentCallsGetDifferentAccounts(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 30, "work": 10, "side": 50})
	var mu sync.Mutex
	var got []string
	var wg sync.WaitGroup
	for range 3 {
		wg.Go(func() {
			out, errOut, code := c.run("", "pick")
			if code != 0 {
				t.Errorf("%d %s", code, errOut)
			}
			mu.Lock()
			got = append(got, strings.TrimSpace(out))
			mu.Unlock()
		})
	}
	wg.Wait()
	slices.Sort(got)
	if want := []string{"personal", "side", "work"}; !slices.Equal(got, want) {
		t.Fatalf("picked %v", got)
	}
}

func TestPickJSONDryRunAndNoRoom(t *testing.T) {
	c := newTestCLI(t, map[string]int{"personal": 30, "work": 10})
	out, _, code := c.run("", "--dry-run", "pick", "--json")
	var r pickResult
	if code != 0 || json.Unmarshal([]byte(out), &r) != nil || r.Name != "work" || r.FiveHour == nil || r.FiveHour.Percent != 10 {
		t.Fatalf("%d %s", code, out)
	}
	// A dry run does not count as a pick.
	if out, _, _ := c.run("", "pick"); strings.TrimSpace(out) != "work" {
		t.Fatalf("after a dry run: %q", out)
	}
	if out, _, _ := c.run("", "pick", "--max-usage", "20"); strings.TrimSpace(out) != "personal" {
		t.Fatalf("over the threshold still overflows: %q", out)
	}

	full := newTestCLI(t, map[string]int{"personal": 100})
	if _, errOut, code := full.run("", "pick"); code != exitNoRoom || !strings.Contains(errOut, "headroom") {
		t.Fatalf("full: %d %q", code, errOut)
	}
	none := newTestCLI(t, nil)
	if _, errOut, code := none.run("", "pick"); code != exitNoRoom || !strings.Contains(errOut, "logged-in") {
		t.Fatalf("nobody logged in: %d %q", code, errOut)
	}
}
