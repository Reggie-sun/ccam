// Package update replaces the running ccam binary with a GitHub release.
// `ccam update` uses it, and the dashboard asks Newer each time it opens.
package update

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jhyunwoo/ccam/internal/fsutil"
)

const (
	// DefaultReleases is the GitHub releases page; CCAM_RELEASES_URL overrides
	// it with anything laid out the same way (<url>/latest redirects to
	// <url>/tag/<tag>, assets live at <url>/download/<tag>/<asset>).
	DefaultReleases = "https://github.com/jhyunwoo/ccam/releases"
	InstallScript   = "curl -fsSL https://ccam.dev/install.sh | sh"
	// InstallScriptWindows is the PowerShell install command.
	InstallScriptWindows = "irm https://ccam.dev/install.ps1 | iex"

	lookupTimeout   = 10 * time.Second
	downloadTimeout = 2 * time.Minute
	maxDownload     = 256 << 20
)

var errNotFound = errors.New("not found")

// Updater holds everything Run touches, so tests can point it elsewhere.
type Updater struct {
	Releases     string // see DefaultReleases
	Current      string // running version, such as 0.2.0 (goreleaser drops the v)
	Exe          string // running binary with symlinks resolved
	GOOS, GOARCH string
	Client       *http.Client
	SystemDirs   []string               // binaries here belong to a package manager
	PackageOwned func(path string) bool // reports whether dpkg installed path
}

// Options are the command's flags.
type Options struct {
	Check   bool   // only report
	Version string // install this tag instead of the latest
	DryRun  bool
}

// New returns an Updater for the running binary.
func New(current, exe, goos, goarch string) *Updater {
	releases := strings.TrimRight(os.Getenv("CCAM_RELEASES_URL"), "/")
	if releases == "" {
		releases = DefaultReleases
	}
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: lookupTimeout}).DialContext,
		TLSHandshakeTimeout:   lookupTimeout,
		ResponseHeaderTimeout: lookupTimeout,
		ForceAttemptHTTP2:     true,
	}
	u := &Updater{
		Releases:     releases,
		Current:      current,
		Exe:          exe,
		GOOS:         goos,
		GOARCH:       goarch,
		Client:       &http.Client{Transport: tr},
		SystemDirs:   []string{"/usr/bin", "/bin", "/usr/sbin", "/sbin"},
		PackageOwned: dpkgOwned,
	}
	if goos == "windows" {
		u.SystemDirs, u.PackageOwned = nil, nil
	}
	return u
}

// Asset is the release archive for one platform (.goreleaser.yaml).
func Asset(goos, goarch string) string {
	if goos == "windows" {
		return "ccam_windows_" + goarch + ".zip"
	}
	return "ccam_" + goos + "_" + goarch + ".tar.gz"
}

// installScript is the command that installs ccam on goos.
func installScript(goos string) string {
	if goos == "windows" {
		return InstallScriptWindows
	}
	return InstallScript
}

// Run checks for a newer release and, unless o.Check or o.DryRun, installs it.
func Run(ctx context.Context, u *Updater, o Options, out io.Writer) error {
	_, release := parseVersion(u.Current)
	var tag string
	if o.Version != "" {
		tag = "v" + strings.TrimPrefix(o.Version, "v")
		if _, ok := parseVersion(tag); !ok {
			return fmt.Errorf("--version takes a release tag such as v0.3.0, not %q", o.Version)
		}
	} else {
		if !release {
			return fmt.Errorf("this ccam is a development build (version %s), so there is no release to compare it with.\n"+
				"Install a release with `ccam update --version vX.Y.Z` (see %s)", orDev(u.Current), u.Releases)
		}
		var err error
		if tag, err = u.Latest(ctx); err != nil {
			return err
		}
	}
	want := strings.TrimPrefix(tag, "v")
	switch {
	case o.Version != "" && want == u.Current:
		fmt.Fprintf(out, "ccam %s is already installed.\n", want)
		return nil
	case o.Version == "" && compareVersions(u.Current, want) >= 0:
		fmt.Fprintf(out, "ccam %s is up to date.\n", u.Current)
		return nil
	case o.Check:
		fmt.Fprintf(out, "ccam %s is available (you have %s). Run `ccam update`.\n", want, orDev(u.Current))
		return nil
	}
	if u.Managed() {
		return u.managedError(tag)
	}
	if err := writable(filepath.Dir(u.Exe)); err != nil {
		return fmt.Errorf("cannot replace %s: %w.\nRerun the install script, which installs ccam in ~/.local/bin:\n  %s", u.Exe, err, installScript(u.GOOS))
	}
	base := u.Releases + "/download/" + tag + "/"
	asset := Asset(u.GOOS, u.GOARCH)
	if o.DryRun {
		fmt.Fprintf(out, "Would update ccam %s → %s:\n  download  %s\n  verify    with %s\n  replace   %s\n",
			orDev(u.Current), want, base+asset, base+"checksums.txt", u.Exe)
		return nil
	}
	fmt.Fprintf(out, "→ downloading %s\n", base+asset)
	bin, err := u.Fetch(ctx, tag)
	if err != nil {
		return err
	}
	if err := u.Install(bin); err != nil {
		return fmt.Errorf("could not replace %s: %w", u.Exe, err)
	}
	fmt.Fprintf(out, "✓ Updated ccam %s → %s (%s)\n", orDev(u.Current), want, u.Exe)
	return nil
}

func orDev(v string) string {
	if v == "" {
		return "dev"
	}
	return v
}

// Newer returns the latest release tag when it is newer than the running
// version, and "" otherwise. A development build has nothing to compare
// with, so it gets "" without a request.
func (u *Updater) Newer(ctx context.Context) (string, error) {
	if _, release := parseVersion(u.Current); !release {
		return "", nil
	}
	tag, err := u.Latest(ctx)
	if err != nil || compareVersions(u.Current, strings.TrimPrefix(tag, "v")) >= 0 {
		return "", err
	}
	return tag, nil
}

// Latest asks the releases page which tag is the latest. GitHub answers
// with a redirect to /tag/<tag>, which avoids the API's rate limit.
func (u *Updater) Latest(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, lookupTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.Releases+"/latest", nil)
	if err != nil {
		return "", err
	}
	c := *u.Client
	c.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := c.Do(req)
	if err != nil {
		return "", networkError(req.URL, err)
	}
	resp.Body.Close()
	if resp.StatusCode/100 != 3 {
		return "", statusError(req.URL, resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	i := strings.LastIndex(loc, "/tag/")
	if i < 0 {
		return "", fmt.Errorf("could not tell the latest release from %s (redirected to %q)", req.URL, loc)
	}
	tag, err := url.PathUnescape(loc[i+len("/tag/"):])
	if _, ok := parseVersion(tag); err != nil || !ok {
		return "", fmt.Errorf("the latest release has an unexpected tag %q", loc[i+len("/tag/"):])
	}
	return tag, nil
}

// Fetch downloads the archive and checksums.txt of tag in parallel, checks
// the archive's SHA-256 and returns the ccam binary inside it.
func (u *Updater) Fetch(ctx context.Context, tag string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, downloadTimeout)
	defer cancel()
	base := u.Releases + "/download/" + tag + "/"
	asset := Asset(u.GOOS, u.GOARCH)
	var sums, archive []byte
	var sumsErr, archiveErr error
	var wg sync.WaitGroup
	wg.Go(func() { sums, sumsErr = u.get(ctx, base+"checksums.txt") })
	wg.Go(func() { archive, archiveErr = u.get(ctx, base+asset) })
	wg.Wait()
	switch {
	case errors.Is(sumsErr, errNotFound) && errors.Is(archiveErr, errNotFound):
		return nil, fmt.Errorf("there is no release %s (see %s)", tag, u.Releases)
	case errors.Is(archiveErr, errNotFound):
		return nil, fmt.Errorf("release %s has no %s for %s/%s (see %s/tag/%s)", tag, asset, u.GOOS, u.GOARCH, u.Releases, tag)
	case archiveErr != nil:
		return nil, archiveErr
	case sumsErr != nil:
		return nil, sumsErr
	}
	want, ok := checksumFor(sums, asset)
	if !ok {
		return nil, fmt.Errorf("checksums.txt of %s does not list %s; nothing was changed", tag, asset)
	}
	sum := sha256.Sum256(archive)
	if got := hex.EncodeToString(sum[:]); got != want {
		return nil, fmt.Errorf("checksum mismatch for %s (expected %s, got %s); nothing was changed", asset, want, got)
	}
	if u.GOOS == "windows" {
		return extractZip(archive)
	}
	return extract(archive)
}

// Install replaces the binary in one rename, keeping its permission bits.
// Shortcut commands link to this path, so they run the new version.
//
// Windows does not let a running program's file be replaced or deleted,
// only renamed: the old binary becomes <exe>.old (or <exe>.old-<time>
// while a session still runs from an earlier one), removed by a later
// update, and the new one takes its name.
func (u *Updater) Install(bin []byte) error {
	info, err := os.Stat(u.Exe)
	if err != nil {
		return err
	}
	if u.GOOS != "windows" {
		return fsutil.AtomicWrite(u.Exe, bin, info.Mode().Perm())
	}
	tmp := u.Exe + ".new"
	if err := os.WriteFile(tmp, bin, 0o755); err != nil {
		os.Remove(tmp)
		return err
	}
	old, err := fsutil.MoveAside(u.Exe)
	if err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, u.Exe); err != nil {
		_ = os.Rename(old, u.Exe)
		os.Remove(tmp)
		return err
	}
	return nil
}

func (u *Updater) get(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, networkError(req.URL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, statusError(req.URL, resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxDownload+1))
	if err != nil {
		return nil, fmt.Errorf("download of %s failed: %w", req.URL, err)
	}
	if len(data) > maxDownload {
		return nil, fmt.Errorf("%s is larger than %d MB", req.URL, maxDownload>>20)
	}
	return data, nil
}

func networkError(u *url.URL, err error) error {
	var ue *url.Error
	if errors.As(err, &ue) {
		err = ue.Err
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("no answer from %s in time; check your connection and try again", u.Host)
	}
	return fmt.Errorf("could not reach %s (%v); check your connection and try again", u.Host, err)
}

func statusError(u *url.URL, code int) error {
	switch code {
	case http.StatusNotFound:
		return fmt.Errorf("%s: %w (HTTP 404)", u, errNotFound)
	case http.StatusForbidden, http.StatusTooManyRequests:
		return fmt.Errorf("%s refused the request (HTTP %d, probably a rate limit); try again later", u.Host, code)
	}
	return fmt.Errorf("unexpected answer from %s (HTTP %d)", u, code)
}

// checksumFor finds asset in goreleaser's "<sha256>  <file>" lines.
func checksumFor(sums []byte, asset string) (string, bool) {
	sc := bufio.NewScanner(bytes.NewReader(sums))
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) == 2 && strings.TrimPrefix(f[1], "*") == asset {
			return strings.ToLower(f[0]), true
		}
	}
	return "", false
}

// extractZip returns ccam.exe from a .zip (the Windows archives).
func extractZip(archive []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, fmt.Errorf("could not read the archive: %w", err)
	}
	for _, f := range zr.File {
		if path.Clean(f.Name) != "ccam.exe" || !f.Mode().IsRegular() {
			continue
		}
		if f.UncompressedSize64 > maxDownload {
			return nil, fmt.Errorf("the ccam binary in the archive is larger than %d MB", maxDownload>>20)
		}
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("could not read the archive: %w", err)
		}
		defer rc.Close()
		bin, err := io.ReadAll(io.LimitReader(rc, maxDownload+1))
		if err != nil {
			return nil, fmt.Errorf("could not read the archive: %w", err)
		}
		if len(bin) > maxDownload {
			return nil, fmt.Errorf("the ccam binary in the archive is larger than %d MB", maxDownload>>20)
		}
		return bin, nil
	}
	return nil, errors.New("the archive has no ccam.exe")
}

// extract returns the ccam file from a .tar.gz.
func extract(archive []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, fmt.Errorf("could not read the archive: %w", err)
	}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil, errors.New("the archive has no ccam binary")
		}
		if err != nil {
			return nil, fmt.Errorf("could not read the archive: %w", err)
		}
		if h.Typeflag == tar.TypeReg && path.Clean(h.Name) == "ccam" {
			bin, err := io.ReadAll(io.LimitReader(tr, maxDownload+1))
			if err != nil {
				return nil, fmt.Errorf("could not read the archive: %w", err)
			}
			if len(bin) > maxDownload {
				return nil, fmt.Errorf("the ccam binary in the archive is larger than %d MB", maxDownload>>20)
			}
			return bin, nil
		}
	}
}

// Managed reports whether a package manager owns the binary.
func (u *Updater) Managed() bool {
	dir := filepath.Dir(u.Exe)
	for _, d := range u.SystemDirs {
		if dir == d {
			return true
		}
	}
	return u.PackageOwned != nil && u.PackageOwned(u.Exe)
}

func (u *Updater) managedError(tag string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "%s was installed by a package, so `ccam update` leaves it alone.\n", u.Exe)
	if u.GOOS == "linux" {
		deb := "ccam_" + strings.TrimPrefix(tag, "v") + "_linux_" + u.GOARCH + ".deb"
		fmt.Fprintf(&b, "To update the .deb, download %s from %s/tag/%s and run:\n  sudo apt install ./%s\n", deb, u.Releases, tag, deb)
		b.WriteString("Or install a copy in ~/.local/bin with the install script:\n")
	} else {
		b.WriteString("Install a copy in ~/.local/bin with the install script:\n")
	}
	b.WriteString("  " + installScript(u.GOOS))
	return errors.New(b.String())
}

func dpkgOwned(path string) bool {
	if _, err := exec.LookPath("dpkg"); err != nil {
		return false
	}
	return exec.Command("dpkg", "-S", path).Run() == nil
}

// writable creates and removes a temp file in dir.
func writable(dir string) error {
	f, err := os.CreateTemp(dir, ".ccam-update-*")
	if err != nil {
		var pe *os.PathError
		if errors.As(err, &pe) {
			err = pe.Err
		}
		return fmt.Errorf("%s is not writable (%v)", dir, err)
	}
	f.Close()
	os.Remove(f.Name())
	return nil
}

// versionRe is MAJOR.MINOR.PATCH with an optional v, -prerelease and +build.
// A tag ends up in download URLs and in a printed shell command, so nothing
// else (no '/', '%', spaces or shell characters) gets through. It is
// compiled on first use, not at init of every ccam process.
var versionRe = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z.-]+))?(?:\+[0-9A-Za-z.-]+)?$`)
})

// parseVersion reads a version; ok is false for anything else, such as "dev".
func parseVersion(v string) (version, bool) {
	m := versionRe().FindStringSubmatch(v)
	if m == nil {
		return version{}, false
	}
	var out version
	for i := range out.nums {
		n, err := strconv.Atoi(m[i+1])
		if err != nil {
			return version{}, false
		}
		out.nums[i] = n
	}
	out.pre = m[4]
	return out, true
}

type version struct {
	nums [3]int
	pre  string
}

// compareVersions orders release versions; a prerelease sorts before its
// release. Prereleases compare as plain strings, which is enough here.
func compareVersions(a, b string) int {
	va, _ := parseVersion(a)
	vb, _ := parseVersion(b)
	for i := range va.nums {
		if va.nums[i] != vb.nums[i] {
			if va.nums[i] < vb.nums[i] {
				return -1
			}
			return 1
		}
	}
	switch {
	case va.pre == vb.pre:
		return 0
	case va.pre == "":
		return 1
	case vb.pre == "":
		return -1
	case va.pre < vb.pre:
		return -1
	}
	return 1
}
