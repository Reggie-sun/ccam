package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestUninstall runs a copy of ccam, since uninstall deletes the binary it
// runs from: every account goes (the managed ones logged out first), and
// so do the shortcuts, ccam's folders and the copy. A folder registered in
// place stays, out of its groups, with the status line it had before ccam.
func TestUninstall(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("a running program cannot delete itself on Windows")
	}
	h := newHarness(t)
	self := filepath.Join(h.home, "inst", "ccam")
	data, err := os.ReadFile(ccamBin)
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Dir(self), 0o755)
	if err := os.WriteFile(self, data, 0o755); err != nil {
		t.Fatal(err)
	}
	if self, err = filepath.EvalSymlinks(self); err != nil { // as ccam prints it (/private/var on macOS)
		t.Fatal(err)
	}
	run := func(stdin string, args ...string) result {
		t.Helper()
		cmd := exec.Command(self, args...)
		cmd.Env = h.env
		cmd.Stdin = strings.NewReader(stdin)
		var out, errb strings.Builder
		cmd.Stdout, cmd.Stderr = &out, &errb
		code := 0
		if err := cmd.Run(); err != nil {
			ee, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatal(err)
			}
			code = ee.ExitCode()
		}
		return result{out.String(), errb.String(), code}
	}
	ok := func(args ...string) string {
		t.Helper()
		r := run("", args...)
		if r.code != 0 {
			t.Fatalf("%v: %+v", args, r)
		}
		return r.stdout
	}

	if r := run("", "add", "work"); r.code != 0 {
		t.Fatalf("add: %+v", r)
	}
	h.writeLogin("work", `{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT"}}`)
	ok("shortcut", "work", "ccwork")
	other := filepath.Join(h.home, "home", ".claude-other")
	os.MkdirAll(other, 0o700)
	os.WriteFile(filepath.Join(other, ".credentials.json"), []byte(`{"claudeAiOauth":{"accessToken":"fake","refreshToken":"FAKE-RT-OTHER"}}`), 0o600)
	os.WriteFile(filepath.Join(other, "settings.json"), []byte(`{"statusLine":{"type":"command","command":"echo 'hi there'","padding":1}}`), 0o600)
	ok("import", other)
	ok("group", "join", "default", "other", "--yes")
	ok("statusline", "install")
	shortcut := filepath.Join(h.home, "bin", "ccwork")
	if _, err := os.Lstat(shortcut); err != nil {
		t.Fatal(err)
	}

	out := ok("--dry-run", "uninstall")
	for _, want := range []string{"work: logged out, then deleted with its folder", "other: " + other + " stays; it leaves its groups",
		"shortcut commands: ccwork", "folder " + filepath.Join(h.home, "data"), "the command " + self} {
		if !strings.Contains(out, want) {
			t.Errorf("dry run lacks %q:\n%s", want, out)
		}
	}
	if r := run("n\n", "uninstall"); r.code != 0 || !strings.Contains(r.stdout, "Cancelled.") || !fileExists(self) || !fileExists(h.profileDir("work")) {
		t.Fatalf("no: %+v", r)
	}

	os.Remove(h.log)
	r := run("y\n", "uninstall")
	if r.code != 0 {
		t.Fatalf("%+v", r)
	}
	if dir, args := h.called(); dir != h.profileDir("work") || strings.Join(args, " ") != "auth logout" {
		t.Fatalf("claude got %s %v", dir, args)
	}
	for _, gone := range []string{self, shortcut, filepath.Join(h.home, "config"), filepath.Join(h.home, "data"), filepath.Join(h.home, "state"), filepath.Join(h.home, "run")} {
		if _, err := os.Lstat(gone); err == nil {
			t.Errorf("%s is still there\n%s", gone, r.stdout)
		}
	}
	if !fileExists(filepath.Join(other, ".credentials.json")) {
		t.Fatal("the imported folder lost its login")
	}
	settings := filepath.Join(other, "settings.json")
	if info, err := os.Lstat(settings); err != nil || info.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("settings.json: %v %v", info, err)
	}
	var doc struct {
		StatusLine map[string]any `json:"statusLine"`
	}
	data, _ = os.ReadFile(settings)
	if err := json.Unmarshal(data, &doc); err != nil || doc.StatusLine["command"] != "echo 'hi there'" || doc.StatusLine["padding"] != 1.0 {
		t.Fatalf("settings.json:\n%s", data)
	}
}
