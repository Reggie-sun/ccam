package e2e

import (
	"archive/zip"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// install.ps1 over an older ccam: a shortcut that is a hard link (Windows
// without Developer Mode) runs the new ccam afterwards, and a reinstall
// goes through while a session still runs from ccam.exe.old. The install
// folder is on $env:Path, so the script leaves the user PATH alone.
func TestInstallScriptOverAnOlderCcam(t *testing.T) {
	h := newHarness(t)
	tmp := t.TempDir()
	zips := map[string]string{}
	for _, v := range []string{"0.2.0", "0.3.0"} {
		bin := filepath.Join(tmp, v, "ccam.exe")
		buildRelease(t, v, bin)
		data, err := os.ReadFile(bin)
		if err != nil {
			t.Fatal(err)
		}
		zips[v] = filepath.Join(tmp, "ccam-"+v+".zip")
		f, _ := os.Create(zips[v])
		zw := zip.NewWriter(f)
		w, _ := zw.Create("ccam.exe")
		w.Write(data)
		zw.Close()
		f.Close()
	}
	dir := filepath.Join(h.home, "bin")
	exe := filepath.Join(dir, "ccam.exe")
	powershell := filepath.Join(os.Getenv("SystemRoot"), "System32", "WindowsPowerShell", "v1.0", "powershell.exe")
	install := func(v string) string {
		t.Helper()
		c := exec.Command(powershell, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command",
			"$env:Path = '"+dir+";' + $env:Path; & '"+filepath.Join(repoRoot(), "install.ps1")+"'")
		c.Env = append(append([]string{}, h.env...), "CCAM_ZIP="+zips[v], "CCAM_INSTALL_DIR="+dir)
		out, err := c.CombinedOutput()
		if err != nil || !strings.Contains(string(out), "OK ccam "+v) {
			t.Fatalf("install %s: %v\n%s", v, err, out)
		}
		return string(out)
	}
	ccam := func(args ...string) {
		t.Helper()
		c := exec.Command(exe, args...)
		c.Env = h.env
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
	}
	same := func(a, b string) bool {
		fa, err1 := os.Stat(a)
		fb, err2 := os.Stat(b)
		return err1 == nil && err2 == nil && os.SameFile(fa, fb)
	}

	install("0.2.0")
	ccam("add", "work", "--no-login")
	ccam("shortcut", "work", "ccwork")
	link := filepath.Join(dir, "ccwork.exe")
	os.Remove(link)
	if err := os.Link(exe, link); err != nil {
		t.Fatal(err)
	}
	// A session of 0.2.0 that is still open after the upgrade.
	session := exec.Command(exe, "work")
	session.Env = append(append([]string{}, h.env...), "FAKE_SLEEP=60")
	if err := session.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(session.Process.Pid)).Run()
		session.Wait()
	})
	waitFile(t, h.log)

	if out := install("0.3.0"); !strings.Contains(out, "OK 1 shortcuts now run the new ccam") {
		t.Fatalf("shortcut not renewed:\n%s", out)
	}
	if !same(link, exe) {
		t.Fatal("ccwork still runs the old ccam")
	}
	// The 0.2.0 session runs from ccam.exe.old now.
	install("0.3.0")
	if aside, _ := filepath.Glob(exe + ".old-*"); len(aside) != 1 || !same(link, exe) {
		t.Fatalf("aside: %v", aside)
	}
}
