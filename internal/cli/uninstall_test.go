package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The status line uninstall leaves is the one ccam wrapped, or none.
func TestRemoveStatusLineRestoresTheWrappedOne(t *testing.T) {
	for command, want := range map[string]string{
		"'/opt/my tools/ccam' statusline":                             "",
		`~/.local/bin/ccam statusline --wrap 'it'\''s "mine"'`:        `it's "mine"`,
		"& 'C:/tools/ccam.exe' statusline --wrap-base64 ZWNobyBoaQ==": "echo hi",
		"ccam statusline --wrap-base64 not-base64":                    "ccam statusline --wrap-base64 not-base64",
		"ccam statusline --wrap 'open":                                "ccam statusline --wrap 'open",
		"echo mine":                                                   "echo mine",
	} {
		path := filepath.Join(t.TempDir(), "settings.json")
		data, _ := json.Marshal(map[string]any{"model": "opus", "statusLine": map[string]any{"type": "command", "command": command, "padding": 2}})
		os.WriteFile(path, data, 0o600)
		changed, err := removeStatusLine(path)
		if err != nil {
			t.Fatal(err)
		}
		var doc struct {
			Model      string
			StatusLine *struct {
				Command string
				Padding int
			} `json:"statusLine"`
		}
		data, _ = os.ReadFile(path)
		json.Unmarshal(data, &doc)
		switch {
		case doc.Model != "opus":
			t.Errorf("%s: other settings lost:\n%s", command, data)
		case want == "" && doc.StatusLine != nil:
			t.Errorf("%s: status line left: %s", command, data)
		case want != "" && (doc.StatusLine == nil || doc.StatusLine.Command != want || doc.StatusLine.Padding != 2):
			t.Errorf("%s: got %s, want %q", command, data, want)
		case changed != (want != command):
			t.Errorf("%s: changed %v", command, changed)
		}
	}
	if changed, err := removeStatusLine(filepath.Join(t.TempDir(), "settings.json")); changed || err != nil {
		t.Fatalf("missing file: %v %v", changed, err)
	}
}
