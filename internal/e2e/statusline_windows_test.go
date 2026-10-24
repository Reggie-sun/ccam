package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The status line command install writes runs in the shell Claude Code
// uses on Windows (Git Bash, else PowerShell), also with ccam in a folder
// whose name needs quotes, and a wrapped command keeps its own quotes.
func TestStatusLineCommandRunsInTheShell(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	data, err := os.ReadFile(ccamBin)
	if err != nil {
		t.Fatal(err)
	}
	self := filepath.Join(t.TempDir(), "my tools $HOME it's", "ccam.exe")
	os.MkdirAll(filepath.Dir(self), 0o755)
	os.WriteFile(self, data, 0o755)
	settings := filepath.Join(h.profileDir("work"), "settings.json")
	os.WriteFile(settings, []byte(`{"statusLine":{"type":"command","command":"echo \"a  b\" 'it''s'"}}`), 0o600)
	install := exec.Command(self, "statusline", "install")
	install.Env = h.env
	if out, err := install.CombinedOutput(); err != nil {
		t.Fatalf("%v %s", err, out)
	}
	var doc struct{ StatusLine struct{ Command string } }
	raw, _ := os.ReadFile(settings)
	json.Unmarshal(raw, &doc)
	sh, flag := `C:\Program Files\Git\bin\bash.exe`, "-c"
	if _, err := os.Stat(sh); err != nil {
		sh, flag = filepath.Join(os.Getenv("SystemRoot"), "System32", "WindowsPowerShell", "v1.0", "powershell.exe"), "-Command"
	}
	run := exec.Command(sh, flag, doc.StatusLine.Command)
	run.Env = append(append([]string{}, h.env...), "CCAM_PROFILE=work")
	run.Stdin = strings.NewReader("{}")
	out, err := run.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "[work]") || !strings.Contains(string(out), "a  b") {
		t.Fatalf("%s %s: %v\n%s", sh, doc.StatusLine.Command, err, out)
	}
}
