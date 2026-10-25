package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/handoff"
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
	// The hook cannot end Claude Code on Windows: the person types /exit
	// there, and the relay carries the conversation on then.
	if !strings.Contains(skill, "on Windows, once they type `/exit`") {
		t.Fatalf("SKILL.md does not say what Windows takes:\n%s", skill)
	}
	// Without the hook, the account /ccam:switch names is unknown to ccam.
	if sw := read(filepath.Join("commands", "switch.md")); !strings.Contains(sw, "/ccam:switch <account>") {
		t.Fatalf("switch.md:\n%s", sw)
	}
	// With auto_switch, a session can start on another account with ccam's
	// fixed first message; the skill says what it means, and the hooks'
	// description names it. A long conversation goes on more cheaply from
	// ccam's summary.
	if !strings.Contains(skill, "auto_switch") || !strings.Contains(skill, handoff.ContinuePrompt) || !strings.Contains(skill, "ccam next --continue --brief") {
		t.Fatalf("SKILL.md does not explain auto_switch and --brief:\n%s", skill)
	}
	if hooks := read(filepath.Join("hooks", "hooks.json")); !strings.Contains(hooks, "auto_switch") {
		t.Fatalf("hooks.json:\n%s", hooks)
	}
}
