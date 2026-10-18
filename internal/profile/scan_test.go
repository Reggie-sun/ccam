package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestTopLevelValue(t *testing.T) {
	doc := `{"projects":{"/a":{"mcpServers":{"nested":1}}},"s":"a\"}b\\","n":-1.5e3,"t":true,"x":null,"mcpServers":{"top":{"command":"c"}},"arr":[{"}":"]"}]}`
	cases := map[string]string{
		"mcpServers": `{"top":{"command":"c"}}`,
		"s":          `"a\"}b\\"`,
		"n":          `-1.5e3`,
		"t":          `true`,
		"x":          `null`,
		"arr":        `[{"}":"]"}]`,
	}
	for key, want := range cases {
		got, ok, err := TopLevelValue([]byte(doc), key)
		if err != nil || !ok || string(got) != want {
			t.Errorf("%s: %q %v %v", key, got, ok, err)
		}
	}
	if _, ok, err := TopLevelValue([]byte(doc), "nested"); ok || err != nil {
		t.Fatal("must not match nested keys")
	}
	if _, ok, _ := TopLevelValue([]byte(`{"mcpServers":1}`), "mcpServers"); !ok {
		t.Fatal("escaped key")
	}
	for _, bad := range []string{``, `[]`, `{"a":`, `{"a":"x`, `{"a" 1}`} {
		if _, ok, err := TopLevelValue([]byte(bad), "zzz"); ok || err == nil {
			t.Errorf("bad %q: ok=%v err=%v", bad, ok, err)
		}
	}
}

// FuzzTopLevelValue checks the scanner against encoding/json on valid input.
func FuzzTopLevelValue(f *testing.F) {
	f.Add(`{"a":1,"mcpServers":{"x":[1,"]",{"y":"\\\""}]}}`, "mcpServers")
	f.Add(`{"k":"v"}`, "k")
	f.Add(` { "a" : [ ] , "b" : { } } `, "b")
	f.Fuzz(func(t *testing.T, doc, key string) {
		var want map[string]json.RawMessage
		if json.Unmarshal([]byte(doc), &want) != nil {
			return
		}
		got, ok, err := TopLevelValue([]byte(doc), key)
		if err != nil {
			t.Fatalf("valid JSON rejected: %v", err)
		}
		w, has := want[key]
		if ok != has {
			t.Fatalf("found=%v want %v", ok, has)
		}
		if !ok {
			return
		}
		var a, b any
		if json.Unmarshal(got, &a) != nil || json.Unmarshal(w, &b) != nil || fmt.Sprint(a) != fmt.Sprint(b) {
			// Duplicate keys: encoding/json keeps the last, we return the first.
			if strings.Count(doc, `"`+key+`"`) < 2 {
				t.Fatalf("got %s want %s", got, w)
			}
		}
	})
}

func bigClaudeJSON() []byte {
	var b bytes.Buffer
	b.WriteString(`{"numStartups":3,"projects":{`)
	for i := 0; i < 400; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, `"/p%d":{"allowedTools":[],"mcpServers":{"n":{"command":"x"}},"history":[`, i)
		for j := 0; j < 20; j++ {
			if j > 0 {
				b.WriteByte(',')
			}
			b.WriteString(`"` + strings.Repeat("x", 200) + `"`)
		}
		b.WriteString(`]}`)
	}
	b.WriteString(`},"oauthAccount":{"emailAddress":"a@b"},"mcpServers":{"a":{"command":"b"}}}`)
	return b.Bytes()
}

func BenchmarkTopLevelValue(b *testing.B) {
	data := bigClaudeJSON()
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		if _, ok, err := TopLevelValue(data, "mcpServers"); !ok || err != nil {
			b.Fatal(ok, err)
		}
	}
}

func BenchmarkDecoderSkip(b *testing.B) {
	data := bigClaudeJSON()
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		var v map[string]json.RawMessage
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.Token()
		for dec.More() {
			k, _ := dec.Token()
			var raw json.RawMessage
			dec.Decode(&raw)
			if k == "mcpServers" {
				_ = json.Unmarshal(raw, &v)
				break
			}
		}
	}
}
