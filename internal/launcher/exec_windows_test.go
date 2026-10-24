package launcher

import (
	"bytes"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

// The hop's ccam is "ccam" whatever file runs it: run through a shortcut
// (a link named after an account), the name would make it that account.
func TestExecSelfKeepsTheName(t *testing.T) {
	if out := os.Getenv("CCAM_TEST_ARGV0"); out != "" {
		os.WriteFile(out, []byte(os.Args[0]), 0o600)
		return
	}
	out := filepath.Join(t.TempDir(), "argv0")
	t.Setenv("CCAM_TEST_ARGV0", out)
	self, _ := os.Executable()
	if err := execSelf(self, []string{"ccam", "-test.run=^TestExecSelfKeepsTheName$"}); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(out); string(got) != "ccam" {
		t.Fatalf("argv[0] %q", got)
	}
}

// Ctrl+C reaches every process on the console. The ccam that waits for the
// hop's ccam must not end on it: the shell would take the console back
// while claude still runs, and the exit code would be lost.
func TestExecSelfOutlivesCtrlC(t *testing.T) {
	dir := os.Getenv("CCAM_TEST_CTRLC")
	self, _ := os.Executable()
	switch os.Getenv("CCAM_TEST_ROLE") {
	case "parent": // the relay's ccam, after its claude ended
		os.Setenv("CCAM_TEST_ROLE", "child")
		err := execSelf(self, []string{"ccam", "-test.run=^TestExecSelfOutlivesCtrlC$"})
		os.WriteFile(filepath.Join(dir, "returned"), []byte(err.Error()), 0o600)
		return
	case "child": // the hop's ccam, whose claude takes Ctrl+C
		signal.Notify(make(chan os.Signal, 1), os.Interrupt)
		os.WriteFile(filepath.Join(dir, "ready"), nil, 0o600)
		time.Sleep(2 * time.Second)
		os.WriteFile(filepath.Join(dir, "done"), nil, 0o600)
		os.Exit(3)
	}
	dir = t.TempDir()
	parent := exec.Command(self, "-test.run=^TestExecSelfOutlivesCtrlC$")
	parent.Env = append(os.Environ(), "CCAM_TEST_CTRLC="+dir, "CCAM_TEST_ROLE=parent")
	parent.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
	if err := parent.Start(); err != nil {
		t.Fatal(err)
	}
	for deadline := time.Now().Add(20 * time.Second); ; time.Sleep(50 * time.Millisecond) {
		if _, err := os.Stat(filepath.Join(dir, "ready")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			parent.Process.Kill()
			t.Fatal("the child did not start")
		}
	}
	// Go reads Ctrl+Break as Ctrl+C; only Ctrl+Break can be sent to a group.
	if err := windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, uint32(parent.Process.Pid)); err != nil {
		parent.Process.Kill()
		t.Fatal(err)
	}
	if err := parent.Wait(); err != nil {
		t.Fatalf("the parent ended on Ctrl+C: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "done")); err != nil {
		t.Fatal("the parent ended before the child")
	}
	if got, _ := os.ReadFile(filepath.Join(dir, "returned")); string(got) != ExitCode(3).Error() {
		t.Fatalf("execSelf returned %q", got)
	}
}

// An API key command runs as typed in cmd.exe: quotes stay as they are.
func TestAPIKeyCommandKeepsItsQuotes(t *testing.T) {
	file := filepath.Join(t.TempDir(), "with space", "key.txt")
	os.MkdirAll(filepath.Dir(file), 0o700)
	os.WriteFile(file, []byte("sk-from-file\r\nsecond line\r\n"), 0o600)
	for command, want := range map[string]string{
		`echo sk-plain`:                     "sk-plain",
		`echo "sk-quoted"`:                  `"sk-quoted"`,
		`type "` + file + `"`:               "sk-from-file",
		`type "` + file + `" | findstr sk-`: "sk-from-file",
		`powershell -NoProfile -Command "Write-Output sk-from-ps"`: "sk-from-ps",
	} {
		got, err := ReadAPIKey("command:" + command)
		if err != nil || got != want {
			t.Errorf("%s: %q, %v; want %q", command, got, err, want)
		}
	}
}

// A batch file gets a command line built for cmd.exe; a line break, which
// would end the command there, is refused.
func TestBatchCommandRefusesALineBreak(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "claude.cmd")
	if _, err := ClaudeCommand(t.Context(), bin, "-p", "two\nlines"); err == nil || !strings.Contains(err.Error(), "line break") {
		t.Fatalf("got %v", err)
	}
	c, err := ClaudeCommand(t.Context(), bin, "-p", `a "b" & c`, "100%")
	if err != nil {
		t.Fatal(err)
	}
	want := `cmd /e:ON /v:OFF /d /c ""` + bin + `" -p "a ""b"" & c" "100%%cd:~,%""`
	if c.SysProcAttr.CmdLine != want || !strings.EqualFold(filepath.Base(c.Path), "cmd.exe") {
		t.Fatalf("%s: %s\nwant %s", c.Path, c.SysProcAttr.CmdLine, want)
	}
}

// A dry run of a batch claude with an argument the launch refuses fails
// as the launch does, instead of printing a line that runs it unquoted.
func TestDryRunOfARefusedBatchLaunchFails(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "claude.cmd")
	pl := &Plan{Dir: `C:\d`, Bin: bin, Argv: []string{"claude", "-p", "two\nlines"}, Env: buildEnv(nil, `C:\d`, "w", false)}
	var buf bytes.Buffer
	if err := pl.dryRun(&buf, nil); err == nil || !strings.Contains(err.Error(), "line break") || buf.Len() != 0 {
		t.Fatalf("%v %q", err, buf.String())
	}
}

// The empty CLAUDE_SECURESTORAGE_CONFIG_DIR of an imported ~/.claude is
// one Windows PowerShell 5.1 cannot set: the dry run says so.
func TestDryRunSaysWhatPowerShellCannotSet(t *testing.T) {
	pl := &Plan{Dir: `C:\h\.claude`, Bin: `C:\claude.exe`, Argv: []string{"claude"}, Env: buildEnv(nil, `C:\h\.claude`, "main", true)}
	var buf bytes.Buffer
	if err := pl.dryRun(&buf, nil); err != nil || !strings.Contains(buf.String(), "\n# The launch sets CLAUDE_SECURESTORAGE_CONFIG_DIR to an empty value") {
		t.Fatalf("%v %s", err, buf.String())
	}
}
