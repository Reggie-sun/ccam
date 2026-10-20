package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFolderPicksTheAccountWithoutAName(t *testing.T) {
	h := newHarness(t)
	h.addAccounts()
	h.ok("ccam", "use", "personal")
	company := filepath.Join(h.home, "code", "company")
	deep := filepath.Join(company, "api", "src")
	other := filepath.Join(h.home, "code", "companyx")
	for _, d := range []string{deep, other} {
		os.MkdirAll(d, 0o700)
	}
	// Entries are stored with symlinks resolved (macOS: /var is /private/var).
	company, _ = filepath.EvalSymlinks(company)
	deep, _ = filepath.EvalSymlinks(deep)

	r := h.runIn(company, "", nil, "ccam", "folder", "work")
	if r.code != 0 || !strings.Contains(r.stdout, company+" → work") {
		t.Fatalf("set: %+v", r)
	}
	launches := []struct {
		dir, want string
		argv      []string
	}{
		{deep, "work", []string{"ccam", "-p", "hi"}},
		{deep, "work", []string{"ccam", "run"}},
		{deep, "personal", []string{"ccam", "personal", "-p", "hi"}}, // a name still wins
		{other, "personal", []string{"ccam", "-p", "hi"}},            // companyx is not inside company
		{h.home, "personal", []string{"ccam", "-p", "hi"}},
	}
	for _, l := range launches {
		if r := h.runIn(l.dir, "", nil, l.argv...); r.code != 0 {
			t.Fatalf("%v in %s: %+v", l.argv, l.dir, r)
		}
		if dir, _ := h.called(); dir != h.profileDir(l.want) {
			t.Errorf("%v in %s started %s, want %s", l.argv, l.dir, dir, l.want)
		}
	}
	if out := h.runIn(deep, "", nil, "ccam", "--dry-run").stdout; !strings.Contains(out, h.profileDir("work")) {
		t.Errorf("dry run: %s", out)
	}

	r = h.runIn(deep, "", nil, "ccam", "folder")
	if !strings.Contains(r.stdout, company+" → work") || !strings.Contains(r.stdout, "Here ("+deep+"): work, set for "+company) {
		t.Fatalf("list: %+v", r)
	}
	r = h.runIn(deep, "", nil, "ccam", "folder", "--remove")
	if r.code == 0 || !strings.Contains(r.stderr, "ccam folder --remove "+company) {
		t.Fatalf("remove below the folder must name the folder: %+v", r)
	}

	// Rename follows; remove drops the account's folders.
	h.ok("ccam", "rename", "work", "job")
	if out := h.runIn(deep, "", nil, "ccam", "folder").stdout; !strings.Contains(out, company+" → job") {
		t.Fatalf("after rename: %s", out)
	}
	h.ok("ccam", "folder", "--remove", company)
	if out := h.runIn(deep, "", nil, "ccam", "folder").stdout; !strings.Contains(out, "No folders set") || !strings.Contains(out, "personal, the default account") {
		t.Fatalf("after remove: %s", out)
	}
	h.ok("ccam", "folder", "job", company)
	h.ok("ccam", "remove", "job", "--yes")
	if data, _ := os.ReadFile(filepath.Join(h.home, "config", "config.toml")); strings.Contains(string(data), "[[folder]]") {
		t.Fatalf("removed account kept its folder:\n%s", data)
	}

	if r := h.run("", nil, "ccam", "folder", "personal", filepath.Join(h.home, "missing")); r.code == 0 {
		t.Fatal("a missing folder was accepted")
	}
	if r := h.run("", nil, "ccam", "folder", "nobody"); r.code == 0 {
		t.Fatal("an unknown account was accepted")
	}
}
