package e2e

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// initLimit is the most one package init may take, in ms. On Windows net's
// init loads Winsock and takes 2 to 7 ms depending on load, so the limit
// there only catches real I/O such as the 110 ms clipboard probe; the launch
// path still never uses the network.
var initLimit = map[bool]float64{true: 10, false: 2}[runtime.GOOS == "windows"]

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
			if ms > initLimit {
				t.Errorf("slow package init: %s", line)
			}
		}
	}
	if total > 2.5*initLimit {
		t.Fatalf("package init takes %.1f ms in total", total)
	}
}

// TestNoTextTemplate guards the binary size: text/template's reflective
// method calls keep every exported method alive, which cost 2 MB (+20%)
// when a custom cobra version template was set.
func TestNoTextTemplate(t *testing.T) {
	buildCcam(t)
	out, err := exec.Command("go", "tool", "nm", ccamBin).Output()
	if err != nil {
		t.Skipf("go tool nm: %v", err)
	}
	if strings.Contains(string(out), " text/template.(*state).evalCall") {
		t.Fatal("text/template is linked in; avoid cobra templates (SetVersionTemplate, SetUsageTemplate, …)")
	}
}
