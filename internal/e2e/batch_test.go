package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// Claude Code installed with npm is claude.cmd, a batch file. ccam runs
// that very file, not a claude.* that cmd.exe would find in the working
// folder, and the arguments reach Claude Code as typed: cmd.exe must not
// read & | > or %NAME% in them as its own.
func TestBatchClaudeGetsTheArgumentsAsTyped(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("batch files are Windows'")
	}
	h := newHarness(t)
	shim := filepath.Join(h.home, "npm", "claude.cmd")
	os.MkdirAll(filepath.Dir(shim), 0o700)
	// The last line is the one npm's cmd-shim writes, with the fake claude
	// for node and its cli.js.
	os.WriteFile(shim, []byte("@ECHO off\r\nSETLOCAL\r\nSET \"_prog="+fakeBin+"\"\r\n"+
		"endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & \"%_prog%\" %*\r\n"), 0o600)
	os.WriteFile(filepath.Join(h.home, "config", "config.toml"), fmt.Appendf(nil, "claude_bin = %q\n", shim), 0o600)
	h.addAccounts() // `claude auth login --email …` through the batch file
	proj := t.TempDir()
	planted := filepath.Join(proj, "planted.txt")
	os.WriteFile(filepath.Join(proj, "claude.bat"), []byte("@echo PLANTED> \""+planted+"\"\r\n"), 0o600)
	os.WriteFile(filepath.Join(proj, "claude.exe"), []byte("not a program"), 0o600)
	args := []string{"-p", `fix A & echo INJECTED> injected.txt`, `Say "hi&echo INJECTED>injected.txt&rem" now`,
		`100% %PATH% %%cd:~,% !OS! %1 %*`, `a"b`, `a\"b`, `C:\dir with space\`, "", `^c|p<l>g(x)`, `"`, `""`,
		`\\server\share`, "--x=1", "it's", "é"}
	// node.exe reads its command line by the C runtime's rules.
	r := h.runIn(proj, "", []string{"FAKE_ARGV=ucrt"}, append([]string{"ccam", "work"}, args...)...)
	if r.code != 0 {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	dir, got := h.called()
	if dir != h.profileDir("work") || !reflect.DeepEqual(got, args) {
		t.Fatalf("%s got\n%q\nwant\n%q", dir, got, args)
	}
	for _, f := range []string{planted, filepath.Join(proj, "injected.txt")} {
		if _, err := os.Stat(f); err == nil {
			t.Fatalf("cmd.exe ran something else: %s exists", f)
		}
	}
	// A line break would end the command in cmd.exe: refused, claude does
	// not run.
	r = h.runIn(proj, "", nil, "ccam", "work", "-p", "two\nlines & echo INJECTED> injected.txt")
	if r.code == 0 || !strings.Contains(r.stderr, "line break") {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	if _, err := os.Stat(h.log); err == nil {
		t.Fatal("claude ran with a line break")
	}
}
