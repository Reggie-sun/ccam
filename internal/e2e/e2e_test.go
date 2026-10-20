// Package e2e builds the real ccam binary and drives it against a fake
// claude inside a temporary CCAM_HOME.
package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	ccamBin   string
	fakeBin   string // testdata/fakeclaude, the stand-in for claude
	buildOnce sync.Once
	buildErr  error
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "ccam-e2e")
	if err != nil {
		panic(err)
	}
	ccamBin = filepath.Join(dir, "ccam"+exe)
	fakeBin = filepath.Join(dir, "fakeclaude"+exe)
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// buildCcam runs inside a test (not TestMain) so the files it reads are
// recorded by the go test cache.
func buildCcam(t *testing.T) {
	buildOnce.Do(func() {
		// Binaries built elsewhere, for machines without Go (a Windows host
		// driven from WSL): CCAM_E2E_BIN, CCAM_E2E_FAKE and CCAM_E2E_ROOT.
		if bin, fake := os.Getenv("CCAM_E2E_BIN"), os.Getenv("CCAM_E2E_FAKE"); bin != "" && fake != "" {
			ccamBin, fakeBin = bin, fake
			return
		}
		touchSources()
		for bin, pkg := range map[string]string{ccamBin: "./cmd/ccam", fakeBin: "./testdata/fakeclaude"} {
			// -buildvcs=false: CI containers may refuse git on the checkout (ownership check).
			build := exec.Command("go", "build", "-buildvcs=false", "-o", bin, pkg)
			build.Dir = repoRoot()
			if out, err := build.CombinedOutput(); err != nil {
				buildErr = fmt.Errorf("build of %s failed: %v\n%s", pkg, err, out)
				return
			}
		}
	})
	if buildErr != nil {
		t.Fatal(buildErr)
	}
}

// touchSources opens every Go source file so the go test cache, which tracks
// files a test reads, invalidates this package when the binary would change.
func touchSources() {
	filepath.WalkDir(repoRoot(), func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && (strings.HasSuffix(path, ".go") || d.Name() == "go.mod" || d.Name() == "go.sum") {
			if f, err := os.Open(path); err == nil {
				f.Close()
			}
		}
		return nil
	})
}

func repoRoot() string {
	if root := os.Getenv("CCAM_E2E_ROOT"); root != "" {
		return root
	}
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

// systemPath is the PATH after the test's bin folder: system folders only.
var systemPath = map[bool]string{true: `C:\Windows\System32;C:\Windows`, false: "/usr/bin:/bin"}[runtime.GOOS == "windows"]

// perms is whether the platform has Unix permission bits (Windows does not).
var perms = runtime.GOOS != "windows"

// exe is the file name suffix of executables.
var exe = map[bool]string{true: ".exe"}[runtime.GOOS == "windows"]

// harness is one isolated ccam installation.
type harness struct {
	t    *testing.T
	home string
	log  string
	env  []string
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	buildCcam(t)
	home := t.TempDir()
	h := &harness{t: t, home: home, log: filepath.Join(home, "claude.log")}
	os.MkdirAll(filepath.Join(home, "config"), 0o700)
	os.WriteFile(filepath.Join(home, "config", "config.toml"), []byte(fmt.Sprintf("claude_bin = %q\n", fakeBin)), 0o600)
	h.env = append(os.Environ(),
		"CCAM_HOME="+home,
		"FAKE_CLAUDE_LOG="+h.log,
		// Only system folders after bin: shortcut checks look for clashes on
		// PATH, and the developer's own ~/.local/bin must not leak in.
		"PATH="+filepath.Join(home, "bin")+string(os.PathListSeparator)+systemPath,
	)
	return h
}

type result struct {
	stdout, stderr string
	code           int
}

// run executes argv[0] (ccam, ccwork …) with extra env.
func (h *harness) run(stdin string, extraEnv []string, argv ...string) result {
	h.t.Helper()
	return h.runIn("", stdin, extraEnv, argv...)
}

// runIn is run in the working directory dir ("" = the test's).
func (h *harness) runIn(dir, stdin string, extraEnv []string, argv ...string) result {
	h.t.Helper()
	bin := ccamBin
	if argv[0] != "ccam" {
		bin = filepath.Join(h.home, "bin", argv[0]+exe)
	}
	cmd := exec.Command(bin, argv[1:]...)
	cmd.Dir = dir
	cmd.Env = append(append([]string{}, h.env...), extraEnv...)
	cmd.Stdin = strings.NewReader(stdin)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		h.t.Fatalf("%v: %v", argv, err)
	}
	return result{out.String(), errb.String(), code}
}

func (h *harness) ok(argv ...string) string {
	h.t.Helper()
	r := h.run("", nil, argv...)
	if r.code != 0 {
		h.t.Fatalf("%v exited %d\nstdout: %s\nstderr: %s", argv, r.code, r.stdout, r.stderr)
	}
	return r.stdout
}

// called returns the env lines and args the fake claude received.
func (h *harness) called() (configDir string, args []string) {
	h.t.Helper()
	data, err := os.ReadFile(h.log)
	if err != nil {
		h.t.Fatalf("claude was not called: %v", err)
	}
	os.Remove(h.log)
	parts := strings.Split(strings.TrimSuffix(string(data), "\x00"), "\x00")
	configDir = strings.TrimPrefix(parts[0], "CLAUDE_CONFIG_DIR=")
	return configDir, parts[2:]
}

func (h *harness) profileDir(name string) string {
	return filepath.Join(h.home, "data", "profiles", name)
}

func (h *harness) addAccounts() {
	h.t.Helper()
	for _, a := range []struct{ name, email string }{{"personal", "me@example.com"}, {"work", "dev@company.com"}} {
		r := h.run("", []string{"FAKE_EMAIL=" + a.email}, "ccam", "add", a.name)
		if r.code != 0 || !strings.Contains(r.stdout, a.email) {
			h.t.Fatalf("add %s: %+v", a.name, r)
		}
	}
}

func TestAddLogsInWithoutShortcutByDefault(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	if entries, _ := os.ReadDir(filepath.Join(h.home, "bin")); len(entries) != 0 {
		t.Fatalf("shortcuts are opt-in, found %v", entries)
	}
	out := h.ok("ccam", "list")
	if !strings.Contains(out, "me@example.com") || !strings.Contains(out, "dev@company.com") {
		t.Fatalf("list:\n%s", out)
	}
	info, _ := os.Stat(h.profileDir("work"))
	if perms && info.Mode().Perm() != 0o700 {
		t.Fatalf("profile dir perm %v", info.Mode())
	}
}

// TestArgumentPassThrough covers every row of the PLAN.md argument table.
func TestArgumentPassThrough(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "use", "personal")
	h.ok("ccam", "shortcut", "personal", "ccpersonal")
	h.ok("ccam", "shortcut", "work", "ccwork")
	cases := []struct {
		argv    []string
		profile string
		want    []string
	}{
		{[]string{"ccam", "work", "--resume"}, "work", []string{"--resume"}},
		{[]string{"ccam", "work", "-p", "corrige les tests — ü", "--output-format", "json"}, "work", []string{"-p", "corrige les tests — ü", "--output-format", "json"}},
		{[]string{"ccwork", "--model", "opus", "--continue"}, "work", []string{"--model", "opus", "--continue"}},
		{[]string{"ccam", "work", "--help"}, "work", []string{"--help"}},
		{[]string{"ccam", "-p", "résume ça"}, "personal", []string{"-p", "résume ça"}},
		{[]string{"ccam", "work", "--", "--version"}, "work", []string{"--version"}},
		{[]string{"ccam", "work", "--version"}, "work", []string{"--version"}},
		{[]string{"ccam", "2"}, "work", nil},
		{[]string{"ccam", "run", "work", "--", "-c"}, "work", []string{"-c"}},
		{[]string{"ccpersonal", "it's \"quoted\"", ""}, "personal", []string{"it's \"quoted\"", ""}},
	}
	for _, c := range cases {
		h.ok(c.argv...)
		dir, args := h.called()
		if dir != h.profileDir(c.profile) {
			t.Errorf("%q: config dir %s", c.argv, dir)
		}
		if !(len(args) == 0 && len(c.want) == 0) && !reflect.DeepEqual(args, c.want) {
			t.Errorf("%q: args %q want %q", c.argv, args, c.want)
		}
	}
}

func TestStdinAndExitCodePassThrough(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	stdinFile := filepath.Join(h.home, "stdin.txt")
	r := h.run("log line\n", []string{"FAKE_CLAUDE_STDIN=" + stdinFile, "FAKE_EXIT=7"}, "ccam", "work", "-p", "trouve la cause")
	if r.code != 7 {
		t.Fatalf("exit code %d", r.code)
	}
	if b, _ := os.ReadFile(stdinFile); string(b) != "log line\n" {
		t.Fatalf("stdin %q", b)
	}
}

func TestDryRunAndErrors(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	os.Remove(h.log) // left behind by `auth login`
	out := h.ok("ccam", "--dry-run", "work", "-c")
	want := "CLAUDE_CONFIG_DIR=" + h.profileDir("work")
	if runtime.GOOS == "windows" {
		want = "$env:CLAUDE_CONFIG_DIR='" + h.profileDir("work") + "';"
	}
	if !strings.Contains(out, want) || !strings.HasSuffix(strings.TrimSpace(out), " -c") {
		t.Fatalf("dry run: %s", out)
	}
	if _, err := os.Stat(h.log); err == nil {
		t.Fatal("dry run must not exec")
	}
	r := h.run("", nil, "ccam", "nobody")
	if r.code == 0 || !strings.Contains(r.stderr, "nobody") {
		t.Fatalf("unknown profile: %+v", r)
	}
	r = h.run("", nil, "ccam", "add", "list")
	if r.code == 0 {
		t.Fatal("reserved name accepted")
	}
}

// TestDryRunNeverChangesState: a command that cannot honour --dry-run must
// refuse it rather than do the real thing.
func TestDryRunNeverChangesState(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	cfgFile := filepath.Join(h.home, "config", "config.toml")
	before, _ := os.ReadFile(cfgFile)
	for _, argv := range [][]string{
		{"remove", "work", "--yes"},
		{"add", "ghost", "--no-login"},
		{"rename", "work", "job"},
		{"use", "work"},
		{"logout", "work"},
		{"login", "work"},
		{"mcp", "add", "fs", "--", "npx", "server"},
		{"import", "--yes"},
		{"shims", "install"},
		{"statusline", "install"},
		{"usage", "--refresh"},
	} {
		r := h.run("", nil, append([]string{"ccam", "--dry-run"}, argv...)...)
		if r.code == 0 || !strings.Contains(r.stderr, "dry-run") {
			t.Errorf("--dry-run %v: %+v", argv, r)
		}
	}
	if after, _ := os.ReadFile(cfgFile); !bytes.Equal(before, after) {
		t.Fatalf("config changed:\n%s\n→\n%s", before, after)
	}
	if _, err := os.Stat(filepath.Join(h.home, "data", "groups", "default", "mcp.json")); err == nil {
		t.Fatal("mcp.json written")
	}
	if _, err := os.Stat(h.profileDir("work")); err != nil {
		t.Fatal("profile folder gone")
	}
	// Read-only commands still accept the flag; doctor --fix only reports.
	h.ok("ccam", "--dry-run", "list")
	h.ok("ccam", "--dry-run", "doctor", "--fix")
}

func TestRemoveAndRename(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "shortcut", "work", "ccwork")
	h.ok("ccam", "rename", "work", "job")
	h.ok("ccwork")
	if dir, _ := h.called(); dir != h.profileDir("job") {
		t.Fatalf("shortcut after rename: %s", dir)
	}
	h.ok("ccam", "job")
	if dir, _ := h.called(); dir != h.profileDir("job") {
		t.Fatalf("rename dir %s", dir)
	}
	r := h.run("n\n", nil, "ccam", "remove", "job")
	if r.code != 0 || !strings.Contains(r.stdout, "Cancelled") {
		t.Fatalf("declined remove: %+v", r)
	}
	h.run("y\n", nil, "ccam", "remove", "job")
	if _, err := os.Stat(h.profileDir("job")); err == nil {
		t.Fatal("profile dir still there")
	}
	if _, err := os.Lstat(filepath.Join(h.home, "bin", "ccwork"+exe)); err == nil {
		t.Fatal("ccwork shortcut still there")
	}
	trash, _ := os.ReadDir(filepath.Join(h.home, "data", "trash"))
	if len(trash) != 1 {
		t.Fatalf("trash %v", trash)
	}
}
