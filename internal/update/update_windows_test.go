package update

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/fsutil"
)

// A session still runs from ccam.exe.old, which the update before left:
// the next update still goes through.
func TestInstallWhileASessionRunsFromTheOldBinary(t *testing.T) {
	if os.Getenv("CCAM_TEST_SLEEP") != "" {
		time.Sleep(10 * time.Second)
		return
	}
	dir := t.TempDir()
	exe := filepath.Join(dir, "ccam.exe")
	self, _ := os.Executable()
	if err := fsutil.CopyFile(self, exe+".old"); err != nil {
		t.Fatal(err)
	}
	run := exec.Command(exe+".old", "-test.run=^TestInstallWhileASessionRunsFromTheOldBinary$")
	run.Env = append(os.Environ(), "CCAM_TEST_SLEEP=1")
	if err := run.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { run.Process.Kill(); run.Wait() })
	time.Sleep(500 * time.Millisecond)
	os.WriteFile(exe, []byte("v1"), 0o755)
	u := &Updater{Exe: exe, GOOS: "windows"}
	if err := u.Install([]byte("v2")); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(exe); string(b) != "v2" {
		t.Fatalf("exe %q", b)
	}
}
