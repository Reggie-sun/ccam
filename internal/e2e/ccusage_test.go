package e2e

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"
)

// The fake claude stands in for ccusage: it records CLAUDE_CONFIG_DIR and
// its arguments, and ends with FAKE_EXIT.
func TestCcusageRunsOverTheAccounts(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	for _, n := range []string{"personal", "work"} {
		if err := os.MkdirAll(filepath.Join(h.profileDir(n), "projects"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(fakeBin)
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Join(h.home, "bin"), 0o700)
	if err := os.WriteFile(filepath.Join(h.home, "bin", "ccusage"+exe), data, 0o755); err != nil {
		t.Fatal(err)
	}

	r := h.run("", []string{"FAKE_EXIT=3", "FAKE_STDOUT=the report"}, "ccam", "ccusage", "daily", "--json")
	if r.code != 3 || !strings.Contains(r.stdout, "the report") {
		t.Fatalf("ccusage's exit code or output is not ccam's: %+v", r)
	}
	dirs, args := h.called()
	if dirs != h.profileDir("personal")+","+h.profileDir("work") || !slices.Equal(args, []string{"daily", "--json"}) {
		t.Fatalf("CLAUDE_CONFIG_DIR=%s %q", dirs, args)
	}

	h.ok("ccam", "ccusage", "-a", "work", "monthly", "-a")
	dirs, args = h.called()
	if dirs != h.profileDir("work") || !slices.Equal(args, []string{"monthly", "-a"}) {
		t.Fatalf("CLAUDE_CONFIG_DIR=%s %q", dirs, args)
	}
}

// npm installs ccusage as ccusage.cmd, and npx is npx.cmd: ccam starts the
// batch file through cmd.exe with the arguments as typed, as it does an npm
// claude.cmd.
func TestBatchCcusageGetsTheArgumentsAsTyped(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("batch files are Windows'")
	}
	h := newHarness(t)
	h.addAccounts()
	os.MkdirAll(filepath.Join(h.profileDir("work"), "projects"), 0o700)
	shim := filepath.Join(h.home, "bin", "ccusage.cmd")
	os.MkdirAll(filepath.Dir(shim), 0o700)
	// The last line is the one npm's cmd-shim writes, with the fake claude
	// for node and its script.
	os.WriteFile(shim, []byte("@ECHO off\r\nSETLOCAL\r\nSET \"_prog="+fakeBin+"\"\r\n"+
		"endLocal & goto #_undefined_# 2>NUL || title %COMSPEC% & \"%_prog%\" %*\r\n"), 0o600)
	proj := t.TempDir()
	args := []string{"daily", `a & echo INJECTED> injected.txt`, `100% %PATH% !OS!`, `a"b`, `C:\dir with space\`, "é"}
	r := h.runIn(proj, "", []string{"FAKE_ARGV=ucrt"}, append([]string{"ccam", "ccusage"}, args...)...)
	if r.code != 0 {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	dir, got := h.called()
	if dir != h.profileDir("work") || !reflect.DeepEqual(got, args) {
		t.Fatalf("%s got\n%q\nwant\n%q", dir, got, args)
	}
	if _, err := os.Stat(filepath.Join(proj, "injected.txt")); err == nil {
		t.Fatal("cmd.exe ran something else")
	}
	// A line break would end the command in cmd.exe: refused, ccusage does
	// not run.
	r = h.runIn(proj, "", nil, "ccam", "ccusage", "daily", "two\nlines & echo INJECTED> injected.txt")
	if r.code == 0 || !strings.Contains(r.stderr, "line break") {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	if _, err := os.Stat(h.log); err == nil {
		t.Fatal("ccusage ran with a line break")
	}
	// The dry run says the same, not the advice for a claude.cmd.
	r = h.runIn(proj, "", nil, "ccam", "--dry-run", "ccusage", "daily", "two\nlines")
	if r.code == 0 || !strings.Contains(r.stderr, "ccusage.cmd is a batch file") || strings.Contains(r.stderr, "claude") {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
}
