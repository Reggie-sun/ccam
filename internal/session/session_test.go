package session

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRegisterAndActive(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "s")
	if Active(dir) {
		t.Fatal("empty dir is not active")
	}
	if err := Register(dir); err != nil {
		t.Fatal(err)
	}
	if !Active(dir) {
		t.Fatal("this process should count as running")
	}
}

func TestStaleMarkersAreRemoved(t *testing.T) {
	dir := t.TempDir()
	// A finished process.
	self, _ := os.Executable()
	c := exec.Command(self, "-test.run=^$") // this test binary, running no tests
	if err := c.Run(); err != nil {
		t.Fatal(err)
	}
	dead := filepath.Join(dir, strconv.Itoa(c.Process.Pid))
	os.WriteFile(dead, []byte("1"), 0o600)
	// A live PID with the wrong start time: the PID was reused.
	reused := filepath.Join(dir, strconv.Itoa(os.Getpid()))
	os.WriteFile(reused, []byte("0"), 0o600)
	if Active(dir) {
		t.Fatal("stale markers must not count")
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatalf("stale markers left: %v", entries)
	}
}

func TestStartTimeHandlesSpacesInName(t *testing.T) {
	if _, ok := startTime(os.Getpid()); !ok {
		t.Fatal("own start time")
	}
	if _, ok := startTime(-1); ok {
		t.Fatal("bogus pid")
	}
}
