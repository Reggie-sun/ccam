// Package trace prints per-step timings of the launch path when CCAM_TRACE=1.
package trace

import (
	"fmt"
	"os"
	"time"
)

var (
	enabled = os.Getenv("CCAM_TRACE") == "1"
	start   = time.Now()
	last    = start
)

// Step records the time spent since the previous step.
func Step(name string) {
	if !enabled {
		return
	}
	now := time.Now()
	fmt.Fprintf(os.Stderr, "ccam trace: %-18s %8.3fms (total %.3fms)\n",
		name, float64(now.Sub(last).Microseconds())/1000, float64(now.Sub(start).Microseconds())/1000)
	last = now
}
