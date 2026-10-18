package e2e

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// TestPackageInitIsCheap guards the launch path against dependencies whose
// init() does I/O (a clipboard package once added ~110 ms to every launch).
func TestPackageInitIsCheap(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	cmd := exec.Command(ccamBin, "--dry-run", "work")
	cmd.Env = append(append([]string{}, h.env...), "GODEBUG=inittrace=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	total := 0.0
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) > 5 && f[0] == "init" && f[5] == "ms" {
			ms, _ := strconv.ParseFloat(f[4], 64)
			total += ms
			if ms > 2 {
				t.Errorf("slow package init: %s", line)
			}
		}
	}
	if total > 5 {
		t.Fatalf("package init takes %.1f ms in total", total)
	}
}
