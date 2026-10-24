package migrate

import (
	"bytes"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// A migrate whose copy fails halfway (a file size limit stands in for a
// full disk, a kill or Ctrl+C) leaves no cut-off file in the group, which a
// rerun would keep as "differs, group copy kept": the rerun copies it whole.
func TestMigrateAfterAFailedCopyCompletesTheFile(t *testing.T) {
	big := bytes.Repeat([]byte("{\"type\":\"user\",\"message\":\"x\"}\n"), 100000) // 3.1 MB
	f := newFixture(t, map[string]string{".claude/projects/p/s1.jsonl": string(big)})
	var old syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_FSIZE, &old); err != nil {
		t.Skip(err)
	}
	lim := old
	lim.Cur = 1 << 20
	if err := syscall.Setrlimit(syscall.RLIMIT_FSIZE, &lim); err != nil {
		t.Skip(err)
	}
	err := f.build(t, "history").Apply()
	syscall.Setrlimit(syscall.RLIMIT_FSIZE, &old)
	if err == nil {
		t.Fatal("the copy did not fail under the size limit")
	}
	f.migrate(t, "history")
	if got, _ := os.ReadFile(filepath.Join(f.group, "projects", "p", "s1.jsonl")); !bytes.Equal(got, big) {
		t.Fatalf("after the rerun the group's copy has %d of %d bytes", len(got), len(big))
	}
}
