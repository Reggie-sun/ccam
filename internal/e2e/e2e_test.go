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
	buildOnce sync.Once
	buildErr  error
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "ccam-e2e")
	if err != nil {
		panic(err)
	}
	ccamBin = filepath.Join(dir, "ccam")
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// buildCcam runs inside a test (not TestMain) so the files it reads are
// recorded by the go test cache.
func buildCcam(t *testing.T) {
	buildOnce.Do(func() {
		touchSources()
		build := exec.Command("go", "build", "-o", ccamBin, "./cmd/ccam")
		build.Dir = repoRoot()
		if out, err := build.CombinedOutput(); err != nil {
			buildErr = fmt.Errorf("build failed: %v\n%s", err, out)
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
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

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
	fake := filepath.Join(repoRoot(), "testdata", "fake-claude")
	os.MkdirAll(filepath.Join(home, "config"), 0o700)
	os.WriteFile(filepath.Join(home, "config", "config.toml"), []byte(fmt.Sprintf("version = 1\nclaude_bin = %q\n", fake)), 0o600)
	h.env = append(os.Environ(),
		"CCAM_HOME="+home,
		"FAKE_CLAUDE_LOG="+h.log,
		"PATH="+filepath.Join(home, "bin")+":"+os.Getenv("PATH"),
	)
	return h
}

type result struct {
	stdout, stderr string
	code           int
}

// run executes argv[0] (ccam, ccs, claude2 …) with extra env.
func (h *harness) run(stdin string, extraEnv []string, argv ...string) result {
	h.t.Helper()
	bin := ccamBin
	if argv[0] != "ccam" {
		bin = filepath.Join(h.home, "bin", argv[0])
	}
	cmd := exec.Command(bin, argv[1:]...)
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

func TestAddCreatesShimsAndLogsIn(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	for _, s := range []string{"claude1", "claude2", "ccs"} {
		if _, err := os.Lstat(filepath.Join(h.home, "bin", s)); err != nil {
			t.Fatalf("missing shim %s", s)
		}
	}
	out := h.ok("ccam", "list")
	if !strings.Contains(out, "me@example.com") || !strings.Contains(out, "dev@company.com") {
		t.Fatalf("list:\n%s", out)
	}
	info, _ := os.Stat(h.profileDir("work"))
	if info.Mode().Perm() != 0o700 {
		t.Fatalf("profile dir perm %v", info.Mode())
	}
}

// TestArgumentPassThrough covers every row of the PLAN.md argument table.
func TestArgumentPassThrough(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "use", "personal")
	cases := []struct {
		argv    []string
		profile string
		want    []string
	}{
		{[]string{"ccs", "work", "--resume"}, "work", []string{"--resume"}},
		{[]string{"ccs", "work", "-p", "테스트 고쳐줘", "--output-format", "json"}, "work", []string{"-p", "테스트 고쳐줘", "--output-format", "json"}},
		{[]string{"claude2", "--model", "opus", "--continue"}, "work", []string{"--model", "opus", "--continue"}},
		{[]string{"ccs", "work", "--help"}, "work", []string{"--help"}},
		{[]string{"ccs", "-p", "요약해줘"}, "personal", []string{"-p", "요약해줘"}},
		{[]string{"ccs", "work", "--", "--version"}, "work", []string{"--version"}},
		{[]string{"ccs", "work", "--version"}, "work", []string{"--version"}},
		{[]string{"ccs", "2"}, "work", nil},
		{[]string{"ccam", "run", "work", "--", "-c"}, "work", []string{"-c"}},
		{[]string{"claude1", "it's \"quoted\"", ""}, "personal", []string{"it's \"quoted\"", ""}},
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
	r := h.run("log line\n", []string{"FAKE_CLAUDE_STDIN=" + stdinFile, "FAKE_EXIT=7"}, "ccs", "work", "-p", "원인 찾아줘")
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
	out := h.ok("ccs", "--dry-run", "work", "-c")
	if !strings.Contains(out, "CLAUDE_CONFIG_DIR="+h.profileDir("work")) || !strings.HasSuffix(strings.TrimSpace(out), " -c") {
		t.Fatalf("dry run: %s", out)
	}
	if _, err := os.Stat(h.log); err == nil {
		t.Fatal("dry run must not exec")
	}
	r := h.run("", nil, "ccs", "nobody")
	if r.code == 0 || !strings.Contains(r.stderr, "nobody") {
		t.Fatalf("unknown profile: %+v", r)
	}
	r = h.run("", nil, "ccam", "add", "list")
	if r.code == 0 {
		t.Fatal("reserved name accepted")
	}
}

func TestRemoveAndRename(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "rename", "work", "job")
	h.ok("ccs", "job")
	if dir, _ := h.called(); dir != h.profileDir("job") {
		t.Fatalf("rename dir %s", dir)
	}
	r := h.run("n\n", nil, "ccam", "remove", "job")
	if r.code != 0 || !strings.Contains(r.stdout, "취소") {
		t.Fatalf("declined remove: %+v", r)
	}
	h.run("y\n", nil, "ccam", "remove", "job")
	if _, err := os.Stat(h.profileDir("job")); err == nil {
		t.Fatal("profile dir still there")
	}
	if _, err := os.Lstat(filepath.Join(h.home, "bin", "claude2")); err == nil {
		t.Fatal("claude2 shim still there")
	}
	trash, _ := os.ReadDir(filepath.Join(h.home, "data", "trash"))
	if len(trash) != 1 {
		t.Fatalf("trash %v", trash)
	}
}
