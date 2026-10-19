package share

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/profile"
)

// MergeRegistry merges my plugin registry file into the group's, saving the
// group's copy to save first. My paths through any of myDirs are pointed at
// the group's plugins folder before merging, so a plugin both sides have
// stays one entry instead of two that differ only in installPath.
func MergeRegistry(mine, group, save string, myDirs ...string) error {
	if err := fsutil.CopyFile(group, save); err != nil {
		return err
	}
	grp, err := os.ReadFile(group)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(mine)
	if err != nil {
		return err
	}
	var doc any
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("could not read my JSON: %w", err)
	}
	sep := string(filepath.Separator)
	groupDir := filepath.Dir(group) + sep
	doc = profile.RewriteStrings(doc, func(s string) string {
		for _, d := range myDirs {
			if rest, ok := strings.CutPrefix(s, d+sep); ok {
				return groupDir + rest
			}
		}
		return s
	})
	if data, err = json.Marshal(doc); err != nil {
		return err
	}
	merged, err := MergeJSONBytes(grp, data)
	if err != nil {
		return err
	}
	return fsutil.AtomicWrite(group, merged, 0o600)
}

// MergeJSONBytes deep-merges two JSON objects: objects merge recursively,
// arrays become the union (base order first), and on scalar clashes base wins
// so the shared group's behaviour does not change under other accounts.
func MergeJSONBytes(base, overlay []byte) ([]byte, error) {
	var a, b map[string]any
	if err := json.Unmarshal(base, &a); err != nil {
		return nil, fmt.Errorf("could not read the group JSON: %w", err)
	}
	if err := json.Unmarshal(overlay, &b); err != nil {
		return nil, fmt.Errorf("could not read my JSON: %w", err)
	}
	return fsutil.MarshalJSON(mergeValue(a, b))
}

func mergeValue(base, overlay any) any {
	switch bv := base.(type) {
	case map[string]any:
		ov, ok := overlay.(map[string]any)
		if !ok {
			return base
		}
		out := make(map[string]any, len(bv)+len(ov))
		for k, v := range bv {
			out[k] = v
		}
		for k, v := range ov {
			if cur, ok := out[k]; ok {
				out[k] = mergeValue(cur, v)
			} else {
				out[k] = v
			}
		}
		return out
	case []any:
		ov, ok := overlay.([]any)
		if !ok {
			return base
		}
		out := append([]any{}, bv...)
		for _, v := range ov {
			dup := false
			for _, have := range out {
				if reflect.DeepEqual(have, v) {
					dup = true
					break
				}
			}
			if !dup {
				out = append(out, v)
			}
		}
		return out
	default:
		return base
	}
}
