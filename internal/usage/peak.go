package usage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Peak is the highest percent of the 5-hour and weekly windows in the
// newest numbers saved for account name (a lookup's or the status line's),
// as they stand at now. ok is false when nothing usable is saved.
func Peak(stateDir, name string, now time.Time) (percent float64, ok bool) {
	var newest cached
	for _, src := range []string{name + ".json", name + ".statusline.json"} {
		data, err := os.ReadFile(filepath.Join(stateDir, src))
		var c cached
		if err != nil || json.Unmarshal(data, &c) != nil || c.Limits == nil {
			continue
		}
		if newest.Limits == nil || c.At.After(newest.At) {
			newest = c
		}
	}
	if newest.Limits == nil {
		return 0, false
	}
	l := newest.Limits.as(now)
	return max(l.FiveHour.Percent, l.Weekly.Percent), true
}
