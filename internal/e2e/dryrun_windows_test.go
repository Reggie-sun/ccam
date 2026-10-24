package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// The --dry-run line runs in Windows PowerShell 5.1 (the one Windows has)
// with the arguments as given, and leaves the shell's variables as they
// were.
func TestDryRunLineRunsInPowerShell(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	args := []string{"-p", `say "hi" now`, "a b", "", `C:\dir with space\`, "$HOME", "@x", "it's", "don\u2019t \u2018x\u2019", "--x=1", `a\"b`}
	r := h.run("", nil, append([]string{"ccam", "--dry-run", "work"}, args...)...)
	if r.code != 0 {
		t.Fatalf("%+v", r)
	}
	script := filepath.Join(t.TempDir(), "line.ps1")
	// With a BOM: Windows PowerShell 5.1 reads a .ps1 without one in the
	// ANSI code page, which garbles the typographic quotes of the arguments.
	os.WriteFile(script, []byte("\ufeff$env:CLAUDE_CONFIG_DIR = 'before'\r\n"+strings.TrimSpace(r.stdout)+"\r\n"+
		"'after=' + $env:CLAUDE_CONFIG_DIR + ' profile=' + $env:CCAM_PROFILE + ' code=' + $LASTEXITCODE\r\n"), 0o600)
	powershell := filepath.Join(os.Getenv("SystemRoot"), "System32", "WindowsPowerShell", "v1.0", "powershell.exe")
	c := exec.Command(powershell, "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", script)
	c.Env = h.env
	out, err := c.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "after=before profile= code=0") {
		t.Fatalf("%v\n%s\nline: %s", err, out, r.stdout)
	}
	dir, got := h.called()
	if dir != h.profileDir("work") || !reflect.DeepEqual(got, args) {
		t.Fatalf("%s got %q\nwant %q\nline: %s", dir, got, args, r.stdout)
	}
}
