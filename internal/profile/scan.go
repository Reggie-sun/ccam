package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// TopLevelValue returns the raw JSON value of one top-level key of an object
// without decoding the rest. It is a byte scanner, several times faster than
// json.Decoder on multi-MB .claude.json files, which matters on the launch
// path. Keys are compared after unescaping only when they contain escapes.
func TopLevelValue(data []byte, key string) ([]byte, bool, error) {
	i := skipWS(data, 0)
	if i >= len(data) || data[i] != '{' {
		return nil, false, fmt.Errorf("JSON 객체가 아닙니다")
	}
	i++
	for {
		i = skipWS(data, i)
		if i >= len(data) {
			return nil, false, io.ErrUnexpectedEOF
		}
		if data[i] == '}' {
			return nil, false, nil
		}
		if data[i] != '"' {
			return nil, false, fmt.Errorf("잘못된 JSON (위치 %d)", i)
		}
		end, err := skipString(data, i)
		if err != nil {
			return nil, false, err
		}
		k := data[i+1 : end-1]
		i = skipWS(data, end)
		if i >= len(data) || data[i] != ':' {
			return nil, false, fmt.Errorf("잘못된 JSON (위치 %d)", i)
		}
		i = skipWS(data, i+1)
		vend, err := skipValue(data, i)
		if err != nil {
			return nil, false, err
		}
		if keyEquals(k, key) {
			return data[i:vend], true, nil
		}
		i = skipWS(data, vend)
		if i < len(data) && data[i] == ',' {
			i++
		}
	}
}

func keyEquals(raw []byte, key string) bool {
	if bytes.IndexByte(raw, '\\') < 0 {
		return string(raw) == key
	}
	var s string
	if err := json.Unmarshal(append(append([]byte{'"'}, raw...), '"'), &s); err != nil {
		return false
	}
	return s == key
}

func skipWS(data []byte, i int) int {
	for i < len(data) {
		switch data[i] {
		case ' ', '\t', '\n', '\r':
			i++
		default:
			return i
		}
	}
	return i
}

// skipString returns the index just past the closing quote of the string at i.
func skipString(data []byte, i int) (int, error) {
	i++
	for {
		j := bytes.IndexAny(data[i:], `"\`)
		if j < 0 {
			return 0, io.ErrUnexpectedEOF
		}
		i += j
		if data[i] == '"' {
			return i + 1, nil
		}
		i += 2 // skip the escaped byte
		if i > len(data) {
			return 0, io.ErrUnexpectedEOF
		}
	}
}

// skipValue returns the index just past the value starting at i.
func skipValue(data []byte, i int) (int, error) {
	if i >= len(data) {
		return 0, io.ErrUnexpectedEOF
	}
	switch data[i] {
	case '"':
		return skipString(data, i)
	case '{', '[':
		depth := 0
		for i < len(data) {
			switch data[i] {
			case '"':
				end, err := skipString(data, i)
				if err != nil {
					return 0, err
				}
				i = end
				continue
			case '{', '[':
				depth++
			case '}', ']':
				depth--
				if depth == 0 {
					return i + 1, nil
				}
			}
			i++
		}
		return 0, io.ErrUnexpectedEOF
	default:
		for i < len(data) {
			switch data[i] {
			case ',', '}', ']', ' ', '\t', '\n', '\r':
				return i, nil
			}
			i++
		}
		return i, nil
	}
}
