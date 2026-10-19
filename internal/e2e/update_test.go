package e2e

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// buildRelease builds ccam as a release would, with its version set by ldflags.
func buildRelease(t *testing.T, version, out string) {
	t.Helper()
	touchSources()
	build := exec.Command("go", "build", "-buildvcs=false", "-trimpath",
		"-ldflags", "-s -w -X github.com/jhyunwoo/ccam/internal/cli.Version="+version, "-o", out, "./cmd/ccam")
	build.Dir = repoRoot()
	if b, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build %s: %v\n%s", version, err, b)
	}
}

// TestUpdate installs a release build of 0.2.0, serves 0.3.0 from a fake
// releases page, and checks that `ccam update` swaps the binary in place.
func TestUpdate(t *testing.T) {
	h := newHarness(t)
	tmp := t.TempDir()
	newBin := filepath.Join(tmp, "new", "ccam")
	buildRelease(t, "0.3.0", newBin)
	exe := filepath.Join(h.home, "bin", "ccam")
	buildRelease(t, "0.2.0", exe)
	os.Symlink(exe, filepath.Join(h.home, "bin", "ccwork")) // a shortcut

	data, err := os.ReadFile(newBin)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	tw.WriteHeader(&tar.Header{Name: "ccam", Mode: 0o755, Size: int64(len(data)), Typeflag: tar.TypeReg})
	tw.Write(data)
	tw.Close()
	gz.Close()
	archive := buf.Bytes()
	sum := sha256.Sum256(archive)
	asset := fmt.Sprintf("ccam_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/latest":
			http.Redirect(w, r, "/tag/v0.3.0", http.StatusFound)
		case "/download/v0.3.0/" + asset:
			w.Write(archive)
		case "/download/v0.3.0/checksums.txt":
			fmt.Fprintf(w, "%s  %s\n", hex.EncodeToString(sum[:]), asset)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	run := func(argv ...string) result {
		cmd := exec.Command(exe, argv...)
		cmd.Env = append(append([]string{}, h.env...), "CCAM_RELEASES_URL="+srv.URL)
		var out, errb bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &errb
		err := cmd.Run()
		code := 0
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else if err != nil {
			t.Fatal(err)
		}
		return result{out.String(), errb.String(), code}
	}

	if r := run("update", "--check"); r.code != 0 || r.stdout != "ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.\n" {
		t.Fatalf("check: %+v", r)
	}
	if r := run("--dry-run", "update"); r.code != 0 || !strings.HasPrefix(r.stdout, "Would update ccam 0.2.0 → 0.3.0:\n") {
		t.Fatalf("dry run: %+v", r)
	}
	if r := run("version"); r.stdout != "ccam 0.2.0\n" {
		t.Fatalf("dry run changed the binary: %+v", r)
	}
	r := run("update")
	if r.code != 0 || !strings.HasSuffix(r.stdout, "✓ Updated ccam 0.2.0 → 0.3.0 ("+exe+")\n") {
		t.Fatalf("update: %+v", r)
	}
	if r := run("version"); r.stdout != "ccam 0.3.0\n" {
		t.Fatalf("after update: %+v", r)
	}
	if r := run("update"); r.code != 0 || r.stdout != "ccam 0.3.0 is up to date.\n" {
		t.Fatalf("second update: %+v", r)
	}
	if info, _ := os.Stat(exe); info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("lost the executable bit: %v", info.Mode())
	}
	if got, _ := os.ReadFile(filepath.Join(h.home, "bin", "ccwork")); !bytes.Equal(got, data) {
		t.Fatal("the shortcut link does not reach the new binary")
	}
}
