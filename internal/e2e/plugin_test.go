package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The plugin's instructions are what Claude tells the person: they must
// match what ccam does. /ccam:switch is how a session moves on, and
// `ccam next --continue` works between accounts that do not share history.
func TestPluginSaysHowToCarryOn(t *testing.T) {
	read := func(path string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(repoRoot(), "plugins", "ccam", path))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	skill := read(filepath.Join("skills", "ccam", "SKILL.md"))
	if !strings.Contains(skill, "/ccam:switch") || !strings.Contains(skill, "ccam next --continue") || strings.Contains(skill, "must share") {
		t.Fatalf("SKILL.md:\n%s", skill)
	}
	// Without the hook, the account /ccam:switch names is unknown to ccam.
	if sw := read(filepath.Join("commands", "switch.md")); !strings.Contains(sw, "/ccam:switch <account>") {
		t.Fatalf("switch.md:\n%s", sw)
	}
}
