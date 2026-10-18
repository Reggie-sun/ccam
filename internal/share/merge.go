package share

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jhyunwoo/ccam/internal/fsutil"
)

// MergeJSONBytes deep-merges two JSON objects: objects merge recursively,
// arrays become the union (base order first), and on scalar clashes base wins
// so the shared group's behaviour does not change under other accounts.
func MergeJSONBytes(base, overlay []byte) ([]byte, error) {
	var a, b map[string]any
	if err := json.Unmarshal(base, &a); err != nil {
		return nil, fmt.Errorf("그룹 JSON 을 읽지 못했습니다: %w", err)
	}
	if err := json.Unmarshal(overlay, &b); err != nil {
		return nil, fmt.Errorf("내 JSON 을 읽지 못했습니다: %w", err)
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
