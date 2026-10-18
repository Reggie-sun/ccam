package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhyunwoo/ccam/internal/fsutil"
)

// pluginFiles store absolute paths under the config dir that installed them
// (docs/claude-layout.md), so they break when that folder moves.
var PluginFiles = []string{"installed_plugins.json", "known_marketplaces.json"}

// RewritePluginPaths replaces oldPrefix with newPrefix in the plugin
// registry files of pluginsDir and reports how many paths changed.
func RewritePluginPaths(pluginsDir, oldPrefix, newPrefix string) (int, error) {
	total := 0
	for _, name := range PluginFiles {
		file := filepath.Join(pluginsDir, name)
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var doc any
		if err := json.Unmarshal(data, &doc); err != nil {
			return total, err
		}
		n := 0
		doc = RewriteStrings(doc, func(s string) string {
			if strings.HasPrefix(s, oldPrefix) {
				n++
				return newPrefix + strings.TrimPrefix(s, oldPrefix)
			}
			return s
		})
		if n == 0 {
			continue
		}
		out, err := fsutil.MarshalJSON(doc)
		if err != nil {
			return total, err
		}
		if err := fsutil.AtomicWrite(file, out, 0o600); err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

// RewriteStrings applies fn to every string value in a decoded JSON document.
func RewriteStrings(v any, fn func(string) string) any {
	switch x := v.(type) {
	case string:
		return fn(x)
	case []any:
		for i := range x {
			x[i] = RewriteStrings(x[i], fn)
		}
	case map[string]any:
		for k := range x {
			x[k] = RewriteStrings(x[k], fn)
		}
	}
	return v
}
