//go:build !linux && !darwin && !windows

package session

import "time"

// StartTime is unsupported here; sessions are then not tracked.
func StartTime(pid int) (string, bool) { return "", false }

// sameStart and wallStart are never reached here: no process has a start
// time.
func sameStart(recorded, now string) (same, known bool) { return false, false }
func wallStart(now string) (time.Time, bool)            { return time.Time{}, false }
