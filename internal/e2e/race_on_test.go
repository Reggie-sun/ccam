//go:build race

package e2e

import (
	"testing"
	"time"
)

// raceBuild: `go test -race` builds ccam with -race too, so the flows the
// tests drive are race-checked, not only the test process.
const raceBuild = true

// A -race ccam exits at once (TestMain sets GORACE), not after the race
// runtime's one-second wait, which over the package's hundreds of runs
// took it past go test's 10-minute limit.
func TestRaceBuildExitsAtOnce(t *testing.T) {
	h := newHarness(t)
	fastest := time.Hour
	for range 3 {
		start := time.Now()
		h.ok("ccam", "--version")
		fastest = min(fastest, time.Since(start))
	}
	if fastest > 500*time.Millisecond {
		t.Fatalf("ccam --version took %v (fastest of 3 runs); a -race binary waits 1 s before it exits unless GORACE has atexit_sleep_ms=0", fastest)
	}
}
