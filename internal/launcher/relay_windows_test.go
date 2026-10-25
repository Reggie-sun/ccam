package launcher

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/session"
)

// innerPid starts c, which runs TestHelperReexec's inner copy somewhere
// below it, and returns the pid that copy says.
func innerPid(t *testing.T, c *exec.Cmd) int {
	t.Helper()
	out, _ := c.StdoutPipe()
	if err := c.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Process.Kill(); c.Wait() })
	line, _ := bufio.NewReader(out).ReadString('\n')
	pid, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		t.Fatalf("%v said %q", c.Args, line)
	}
	return pid
}

// On Windows too only the relay's own claude counts: its child, that
// child's child when both are the same program, or Claude Code under the
// cmd.exe that runs a batch claude (npm's claude.cmd). What another program
// started does not.
func TestRelayChildOnWindows(t *testing.T) {
	self, _ := os.Executable()
	helper := func(bin, role string, env ...string) *exec.Cmd {
		c := exec.Command(bin, "-test.run=^TestHelperReexec$")
		c.Env = append(append(os.Environ(), "CCAM_TEST_REEXEC="+role), env...)
		return c
	}
	// A copy of this test binary is another program: another file.
	other := filepath.Join(t.TempDir(), "other.exe")
	data, _ := os.ReadFile(self)
	if err := os.WriteFile(other, data, 0o700); err != nil {
		t.Fatal(err)
	}
	child := innerPid(t, helper(other, "inner"))
	if c, ok := RelayChild(child, os.Getpid()); !ok || c != child {
		t.Fatalf("own child: %d %v", c, ok)
	}
	if _, ok := RelayChild(child, os.Getppid()); ok {
		t.Fatal("another process's grandchild")
	}
	// The same program started again.
	outer := helper(self, "outer")
	inner := innerPid(t, outer)
	if c, ok := RelayChild(inner, os.Getpid()); !ok || c != outer.Process.Pid {
		t.Fatalf("started again: %d %v, want %d", c, ok, outer.Process.Pid)
	}
	// A batch claude, run as the relay runs it.
	bat := filepath.Join(t.TempDir(), "claude.cmd")
	os.WriteFile(bat, []byte("@\""+self+"\" %*\r\n"), 0o600)
	shell, err := batchCommand(t.Context(), bat, []string{"-test.run=^TestHelperReexec$"})
	if err != nil {
		t.Fatal(err)
	}
	shell.Env = append(os.Environ(), "CCAM_TEST_REEXEC=inner")
	inner = innerPid(t, shell)
	if c, ok := RelayChild(inner, os.Getpid()); !ok || c != shell.Process.Pid {
		t.Fatalf("under cmd.exe: %d %v, want %d", c, ok, shell.Process.Pid)
	}
	// Another program runs it.
	inner = innerPid(t, helper(other, "outer", "CCAM_TEST_INNER="+self))
	if _, ok := RelayChild(inner, os.Getpid()); ok {
		t.Fatal("another program's child")
	}
}

// Windows has no way to end Claude Code from the hook: Terminate says so
// and leaves even the relay's own claude running.
func TestTerminateEndsNothingOnWindows(t *testing.T) {
	self, _ := os.Executable()
	c := exec.Command(self, "-test.run=^TestHelperReexec$")
	c.Env = append(os.Environ(), "CCAM_TEST_REEXEC=inner")
	child := innerPid(t, c)
	start, _ := session.StartTime(child)
	if err := Terminate(child, os.Getpid(), child, start); err == nil {
		t.Fatal("no error")
	}
	if _, ok := session.StartTime(child); !ok {
		t.Fatal("the child ended")
	}
}
