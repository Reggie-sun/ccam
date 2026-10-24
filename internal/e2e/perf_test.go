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
// Each package counts with its fastest of a few runs: a busy machine can
// preempt any init once, while real I/O is slow in every run.
func TestPackageInitIsCheap(t *testing.T) {
	if raceBuild {
		t.Skip("a -race build's init is slow by design")
	}
	h := newHarness(t)
	h.addAccounts()
	best := map[string]float64{}
	lines := map[string]string{}
	for range 5 {
		cmd := exec.Command(ccamBin, "--dry-run", "work")
		cmd.Env = append(append([]string{}, h.env...), "GODEBUG=inittrace=1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%v\n%s", err, out)
		}
		for _, line := range strings.Split(string(out), "\n") {
			f := strings.Fields(line)
			if len(f) > 5 && f[0] == "init" && f[5] == "ms" {
				ms, _ := strconv.ParseFloat(f[4], 64)
				if old, ok := best[f[1]]; !ok || ms < old {
					best[f[1]], lines[f[1]] = ms, line
				}
			}
		}
	}
	total := 0.0
	for pkg, ms := range best {
		total += ms
		if ms > initLimit {
			t.Errorf("slow package init (fastest of 5 runs): %s", lines[pkg])
		}
	}
	if total > 2.5*initLimit {
		t.Fatalf("package init takes %.1f ms in total (fastest of 5 runs each)", total)
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
