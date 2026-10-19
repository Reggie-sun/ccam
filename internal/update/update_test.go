package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// fakeRelease serves a GitHub-like releases page whose latest tag is latest.
// Every tag in bins has an archive with that ccam binary and a checksums.txt.
type fakeRelease struct {
	latest string
	bins   map[string]string // tag → binary content
	sums   map[string]string // tag → checksums.txt override

	mu   sync.Mutex // the archive and checksums are fetched concurrently
	hits []string
}

func (f *fakeRelease) requests() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.hits...)
}

func (f *fakeRelease) serve(t *testing.T) string {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.hits = append(f.hits, r.URL.Path)
		f.mu.Unlock()
		if r.URL.Path == "/latest" {
			http.Redirect(w, r, "/tag/"+f.latest, http.StatusFound)
			return
		}
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
		if len(parts) != 3 || parts[0] != "download" {
			http.NotFound(w, r)
			return
		}
		tag, file := parts[1], parts[2]
		bin, ok := f.bins[tag]
		if !ok {
			http.NotFound(w, r)
			return
		}
		archive := tarball(t, bin)
		switch file {
		case Asset("linux", "amd64"):
			w.Write(archive)
		case "checksums.txt":
			if s, ok := f.sums[tag]; ok {
				fmt.Fprint(w, s)
				return
			}
			sum := sha256.Sum256(archive)
			fmt.Fprintf(w, "%s  %s\n%s  ccam_%s_linux_amd64.deb\n", hex.EncodeToString(sum[:]), Asset("linux", "amd64"), strings.Repeat("0", 64), strings.TrimPrefix(tag, "v"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

// tarball builds a goreleaser-like archive.
func tarball(t *testing.T, bin string) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, f := range []struct{ name, body string }{{"LICENSE", "MIT"}, {"ccam", bin}, {"README.md", "readme"}} {
		if err := tw.WriteHeader(&tar.Header{Name: f.name, Mode: 0o755, Size: int64(len(f.body)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		tw.Write([]byte(f.body))
	}
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

// setup installs an "old" binary with mode 0751 in a temp bin folder.
func setup(t *testing.T, current string, f *fakeRelease) *Updater {
	dir := filepath.Join(t.TempDir(), "bin")
	os.Mkdir(dir, 0o755)
	exe := filepath.Join(dir, "ccam")
	os.WriteFile(exe, []byte("old"), 0o751)
	os.Chmod(exe, 0o751)
	u := New(current, exe, "linux", "amd64")
	u.Releases = f.serve(t)
	u.SystemDirs = nil
	u.PackageOwned = func(string) bool { return false }
	return u
}

func run(t *testing.T, u *Updater, o Options) (string, error) {
	var out bytes.Buffer
	err := Run(context.Background(), u, o, &out)
	return out.String(), err
}

func binary(t *testing.T, u *Updater) (string, os.FileMode) {
	data, err := os.ReadFile(u.Exe)
	if err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(u.Exe)
	return string(data), info.Mode().Perm()
}

// leftovers lists files next to the binary other than ccam itself.
func leftovers(t *testing.T, u *Updater) []string {
	entries, _ := os.ReadDir(filepath.Dir(u.Exe))
	var out []string
	for _, e := range entries {
		if e.Name() != "ccam" {
			out = append(out, e.Name())
		}
	}
	return out
}

func release() *fakeRelease {
	return &fakeRelease{latest: "v0.3.0", bins: map[string]string{"v0.2.0": "bin-0.2.0", "v0.3.0": "bin-0.3.0"}, sums: map[string]string{}}
}

func TestUpToDate(t *testing.T) {
	for _, current := range []string{"0.3.0", "0.4.0"} {
		u := setup(t, current, release())
		out, err := run(t, u, Options{})
		if err != nil || out != "ccam "+current+" is up to date.\n" {
			t.Fatalf("%s: %q %v", current, out, err)
		}
		if bin, _ := binary(t, u); bin != "old" {
			t.Fatal("binary changed")
		}
	}
}

func TestCheckOnlyReports(t *testing.T) {
	f := release()
	u := setup(t, "0.2.0", f)
	out, err := run(t, u, Options{Check: true})
	if err != nil || out != "ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.\n" {
		t.Fatalf("%q %v", out, err)
	}
	if bin, _ := binary(t, u); bin != "old" {
		t.Fatal("binary changed")
	}
	if len(f.requests()) != 1 {
		t.Fatalf("--check downloaded something: %v", f.requests())
	}
	u.Current = "0.3.0"
	if out, _ := run(t, u, Options{Check: true}); out != "ccam 0.3.0 is up to date.\n" {
		t.Fatalf("%q", out)
	}
}

func TestUpdateReplacesBinary(t *testing.T) {
	u := setup(t, "0.2.0", release())
	out, err := run(t, u, Options{})
	if err != nil {
		t.Fatal(err)
	}
	want := "→ downloading " + u.Releases + "/download/v0.3.0/ccam_linux_amd64.tar.gz\n✓ Updated ccam 0.2.0 → 0.3.0 (" + u.Exe + ")\n"
	if out != want {
		t.Fatalf("got %q\nwant %q", out, want)
	}
	if bin, mode := binary(t, u); bin != "bin-0.3.0" || mode != 0o751 {
		t.Fatalf("binary %q mode %o", bin, mode)
	}
	if l := leftovers(t, u); len(l) != 0 {
		t.Fatalf("temp files left: %v", l)
	}
}

// A shortcut is a link to the binary; replacing the file by rename keeps it working.
func TestUpdateKeepsLinks(t *testing.T) {
	u := setup(t, "0.2.0", release())
	link := filepath.Join(t.TempDir(), "ccwork")
	os.Symlink(u.Exe, link)
	if _, err := run(t, u, Options{}); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(link); err != nil || string(data) != "bin-0.3.0" {
		t.Fatalf("%q %v", data, err)
	}
}

func TestChecksumMismatchChangesNothing(t *testing.T) {
	f := release()
	f.sums["v0.3.0"] = strings.Repeat("ab", 32) + "  " + Asset("linux", "amd64") + "\n"
	u := setup(t, "0.2.0", f)
	_, err := run(t, u, Options{})
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") || !strings.Contains(err.Error(), "nothing was changed") {
		t.Fatalf("err = %v", err)
	}
	if bin, mode := binary(t, u); bin != "old" || mode != 0o751 {
		t.Fatalf("binary %q mode %o", bin, mode)
	}
	if l := leftovers(t, u); len(l) != 0 {
		t.Fatalf("temp files left: %v", l)
	}
}

func TestAssetNotInChecksums(t *testing.T) {
	f := release()
	f.sums["v0.3.0"] = strings.Repeat("ab", 32) + "  ccam_darwin_arm64.tar.gz\n"
	u := setup(t, "0.2.0", f)
	if _, err := run(t, u, Options{}); err == nil || !strings.Contains(err.Error(), "does not list ccam_linux_amd64.tar.gz") {
		t.Fatalf("err = %v", err)
	}
	if bin, _ := binary(t, u); bin != "old" {
		t.Fatal("binary changed")
	}
}

func TestMissingAsset(t *testing.T) {
	u := setup(t, "0.2.0", release())
	u.GOARCH = "riscv64"
	_, err := run(t, u, Options{})
	if err == nil || !strings.Contains(err.Error(), "release v0.3.0 has no ccam_linux_riscv64.tar.gz for linux/riscv64") {
		t.Fatalf("err = %v", err)
	}
	if bin, _ := binary(t, u); bin != "old" {
		t.Fatal("binary changed")
	}
}

func TestMissingRelease(t *testing.T) {
	u := setup(t, "0.2.0", release())
	_, err := run(t, u, Options{Version: "v9.9.9"})
	if err == nil || !strings.Contains(err.Error(), "there is no release v9.9.9") {
		t.Fatalf("err = %v", err)
	}
}

func TestDevBuildRefused(t *testing.T) {
	for _, current := range []string{"dev", ""} {
		f := release()
		u := setup(t, current, f)
		for _, o := range []Options{{}, {Check: true}} {
			_, err := run(t, u, o)
			if err == nil || !strings.Contains(err.Error(), "development build (version dev)") || !strings.Contains(err.Error(), "ccam update --version vX.Y.Z") {
				t.Fatalf("%q %+v: err = %v", current, o, err)
			}
		}
		if len(f.requests()) != 0 {
			t.Fatalf("dev build went online: %v", f.requests())
		}
		// With --version a development build installs that release.
		out, err := run(t, u, Options{Version: "0.3.0"})
		if err != nil || !strings.HasSuffix(out, "✓ Updated ccam dev → 0.3.0 ("+u.Exe+")\n") {
			t.Fatalf("%q %v", out, err)
		}
	}
}

func TestVersionFlag(t *testing.T) {
	f := release()
	u := setup(t, "0.3.0", f)
	out, err := run(t, u, Options{Version: "v0.2.0"}) // going back
	if err != nil || !strings.HasSuffix(out, "✓ Updated ccam 0.3.0 → 0.2.0 ("+u.Exe+")\n") {
		t.Fatalf("%q %v", out, err)
	}
	if bin, _ := binary(t, u); bin != "bin-0.2.0" {
		t.Fatalf("binary %q", bin)
	}
	for _, h := range f.requests() {
		if h == "/latest" {
			t.Fatal("--version looked up the latest release")
		}
	}
	u.Current = "0.2.0"
	if out, err := run(t, u, Options{Version: "v0.2.0"}); err != nil || out != "ccam 0.2.0 is already installed.\n" {
		t.Fatalf("%q %v", out, err)
	}
	if _, err := run(t, u, Options{Version: "latest"}); err == nil || !strings.Contains(err.Error(), "--version takes a release tag such as v0.3.0") {
		t.Fatalf("err = %v", err)
	}
}

func TestDryRunChangesNothing(t *testing.T) {
	f := release()
	u := setup(t, "0.2.0", f)
	out, err := run(t, u, Options{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	base := u.Releases + "/download/v0.3.0/"
	want := "Would update ccam 0.2.0 → 0.3.0:\n  download  " + base + "ccam_linux_amd64.tar.gz\n  verify    with " + base + "checksums.txt\n  replace   " + u.Exe + "\n"
	if out != want {
		t.Fatalf("got %q\nwant %q", out, want)
	}
	if bin, _ := binary(t, u); bin != "old" {
		t.Fatal("binary changed")
	}
	if len(f.requests()) != 1 || len(leftovers(t, u)) != 0 {
		t.Fatalf("dry run downloaded or wrote: %v %v", f.requests(), leftovers(t, u))
	}
}

func TestNotWritableDir(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root can write anywhere")
	}
	u := setup(t, "0.2.0", release())
	dir := filepath.Dir(u.Exe)
	os.Chmod(dir, 0o555)
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
	_, err := run(t, u, Options{})
	if err == nil || !strings.Contains(err.Error(), "is not writable") || !strings.Contains(err.Error(), InstallScript) {
		t.Fatalf("err = %v", err)
	}
	if bin, _ := binary(t, u); bin != "old" {
		t.Fatal("binary changed")
	}
}

func TestPackageInstallRefused(t *testing.T) {
	for _, viaDpkg := range []bool{false, true} {
		u := setup(t, "0.2.0", release())
		if viaDpkg {
			u.PackageOwned = func(p string) bool { return p == u.Exe }
		} else {
			u.SystemDirs = []string{filepath.Dir(u.Exe)}
		}
		_, err := run(t, u, Options{})
		if err == nil {
			t.Fatal("no error")
		}
		for _, want := range []string{
			u.Exe + " was installed by a package",
			"download ccam_0.3.0_linux_amd64.deb from " + u.Releases + "/tag/v0.3.0",
			"sudo apt install ./ccam_0.3.0_linux_amd64.deb",
			InstallScript,
		} {
			if !strings.Contains(err.Error(), want) {
				t.Fatalf("missing %q in:\n%v", want, err)
			}
		}
		if bin, _ := binary(t, u); bin != "old" {
			t.Fatal("binary changed")
		}
		// Reporting still works.
		if out, err := run(t, u, Options{Check: true}); err != nil || !strings.Contains(out, "0.3.0 is available") {
			t.Fatalf("%q %v", out, err)
		}
	}
}

func TestOffline(t *testing.T) {
	u := setup(t, "0.2.0", release())
	u.Releases = "http://127.0.0.1:1"
	_, err := run(t, u, Options{})
	if err == nil || !strings.Contains(err.Error(), "could not reach 127.0.0.1:1") || !strings.Contains(err.Error(), "check your connection") {
		t.Fatalf("err = %v", err)
	}
}

func TestRateLimited(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()
	u := setup(t, "0.2.0", release())
	u.Releases = srv.URL
	if _, err := run(t, u, Options{}); err == nil || !strings.Contains(err.Error(), "HTTP 429, probably a rate limit") {
		t.Fatalf("err = %v", err)
	}
}

func TestCompareVersions(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want int
	}{
		{"0.2.0", "0.3.0", -1}, {"0.10.0", "0.9.9", 1}, {"1.0.0", "v1.0.0", 0},
		{"1.0.0-rc.1", "1.0.0", -1}, {"1.0.0", "1.0.0-rc.1", 1}, {"1.0.0-rc.1", "1.0.0-rc.2", -1},
	} {
		if got := compareVersions(c.a, c.b); got != c.want {
			t.Errorf("compare(%s, %s) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
	for _, v := range []string{"dev", "", "1.2", "1.2.x", "latest", "+1.2.3", "v1.2.3-%2F..%2Fx", "v1.2.3-a b", "v1.2.3-$(id)", "v1.2.3-../x"} {
		if _, ok := parseVersion(v); ok {
			t.Errorf("parseVersion(%q) ok", v)
		}
	}
}
