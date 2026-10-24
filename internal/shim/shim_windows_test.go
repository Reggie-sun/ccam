package shim

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
)

// A session still runs from a shortcut that is an older ccam's file:
// Windows refuses to delete it, so it is renamed out of the way and linked
// again.
func TestRunningShortcutOfAnOlderCcamIsRenewed(t *testing.T) {
	if os.Getenv("CCAM_TEST_SLEEP") != "" {
		time.Sleep(10 * time.Second)
		return
	}
	f := newBin(t)
	self, _ := os.Executable() // has ccamBuild in it, as every ccam has
	link := filepath.Join(f.bin, "cca.exe")
	if err := fsutil.CopyFile(self, link); err != nil {
		t.Fatal(err)
	}
	run := exec.Command(link, "-test.run=^TestRunningShortcutOfAnOlderCcamIsRenewed$")
	run.Env = append(os.Environ(), "CCAM_TEST_SLEEP=1")
	if err := run.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { run.Process.Kill(); run.Wait() })
	time.Sleep(500 * time.Millisecond)
	if os.Remove(link) == nil {
		t.Fatal("Windows deleted a running program's file")
	}
	cfg := &config.Config{Profiles: []config.Profile{{Name: "a", Slot: 1, Shortcut: "cca"}}}
	if created, err := Install(f.p, f.ccam, cfg); err != nil || len(created) != 1 || !Linked(f.bin, f.ccam, "cca") {
		t.Fatalf("%v %v", created, err)
	}
	if !fileExists(link + ".old") {
		t.Fatal("the running file is not beside it")
	}
}
