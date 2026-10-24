package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf8"
)

// TopLevelValue returns the raw JSON value of one top-level key of an object
// without decoding the rest. It is a byte scanner, several times faster than
// json.Decoder on multi-MB .claude.json files, which matters on the launch
// path. Keys are compared after unescaping only when they contain escapes.
func TopLevelValue(data []byte, key string) ([]byte, bool, error) {
	var found []byte
	err := members(data, func(k, v []byte) bool {
		if keyEquals(k, key) {
			found = v
			return false
		}
		return true
	})
	return found, found != nil, err
}

// EachValue calls fn with each key (unescaped) and raw value of a JSON
// object, in order, with the scanner of TopLevelValue: no value is
// decoded. It stops at the first error fn returns.
func EachValue(data []byte, fn func(key string, value []byte) error) error {
	var ferr error
	err := members(data, func(k, v []byte) bool {
		key := string(k)
		if bytes.IndexByte(k, '\\') >= 0 || !utf8.Valid(k) { // as encoding/json reads it
			if ferr = json.Unmarshal(append(append([]byte{'"'}, k...), '"'), &key); ferr != nil {
				return false
			}
		}
		ferr = fn(key, v)
		return ferr == nil
	})
	if ferr != nil {
		return ferr
	}
	return err
}

// members scans the members of the object in data: key bytes as written
// (escapes kept) and raw value, while fn says to go on.
func members(data []byte, fn func(key, value []byte) bool) error {
	i := skipWS(data, 0)
	if i >= len(data) || data[i] != '{' {
		return fmt.Errorf("not a JSON object")
	}
	i++
	for {
		i = skipWS(data, i)
		if i >= len(data) {
			return io.ErrUnexpectedEOF
		}
		if data[i] == '}' {
			return nil
		}
		if data[i] != '"' {
			return fmt.Errorf("invalid JSON (offset %d)", i)
		}
		end, err := skipString(data, i)
		if err != nil {
			return err
		}
		k := data[i+1 : end-1]
		i = skipWS(data, end)
		if i >= len(data) || data[i] != ':' {
			return fmt.Errorf("invalid JSON (offset %d)", i)
		}
		i = skipWS(data, i+1)
		vend, err := skipValue(data, i)
		if err != nil {
			return err
		}
		if !fn(k, data[i:vend]) {
			return nil
		}
		i = skipWS(data, vend)
		if i < len(data) && data[i] == ',' {
			i++
		}
	}
}

func keyEquals(raw []byte, key string) bool {
	if bytes.IndexByte(raw, '\\') < 0 && utf8.Valid(raw) {
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

// skipString returns the index just past the closing quote of the string at
// i: the first quote after it that an odd run of backslashes does not
// escape. bytes.IndexByte, unlike a search for either of two bytes, scans
// many bytes at a time.
func skipString(data []byte, i int) (int, error) {
	i++
	for {
		q := bytes.IndexByte(data[i:], '"')
		if q < 0 {
			return 0, io.ErrUnexpectedEOF
		}
		end := i + q
		n := 0
		for j := end - 1; j >= i && data[j] == '\\'; j-- {
			n++
		}
		if n%2 == 0 {
			return end + 1, nil
		}
		i = end + 1
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
