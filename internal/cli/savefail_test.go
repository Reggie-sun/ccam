package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// unsaveable makes config.toml impossible to save: a line ccam does not
// read (a setting below the last table), which a save would delete.
func unsaveable(t *testing.T, c *testCLI) func() {
	t.Helper()
	path := c.p.ConfigFile()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(append([]byte{}, data...), "\nrelay = true\n"...), 0o600); err != nil {
		t.Fatal(err)
	}
	return func() { os.WriteFile(path, data, 0o600) }
}

// noClaude points claude_bin at a file that does not exist, so that no
// claude runs (where there is no `true`, PATH would find the real one):
// remove's logout then fails, and the removal goes on.
func noClaude(t *testing.T, c *testCLI) {
	t.Helper()
	cfg := c.cfg()
	cfg.ClaudeBin = filepath.Join(t.TempDir(), "no-claude")
	if err := cfg.Save(c.p.ConfigFile()); err != nil {
		t.Fatal(err)
	}
}

// pluginOf records a plugin installed through work's own plugins folder.
func pluginOf(t *testing.T, c *testCLI, name string) string {
	t.Helper()
	dir := c.p.ProfileDir(name)
	path := filepath.Join(dir, "plugins", "cache", "m", "p", "1.0")
	writeFiles(t, filepath.Join(dir, "plugins"), map[string]string{
		"installed_plugins.json": `{"version":2,"plugins":{"p@m":[{"installPath":"` + filepath.ToSlash(path) + `"}]}}`,
	})
	return filepath.ToSlash(path)
}

// When config.toml cannot be saved, rename, remove and add leave the folders
// where the config still has them: before, the account's folder was gone
// (renamed, or in the trash) while the config still named it.
func TestFailedSaveLeavesTheFoldersAsTheConfigHasThem(t *testing.T) {
	c := newTestCLI(t, nil)
	noClaude(t, c)
	work := c.p.ProfileDir("work")
	writeFiles(t, work, map[string]string{".credentials.json": `{"claudeAiOauth":{"accessToken":"t","refreshToken":"r","expiresAt":4102444800000}}`})
	installed := pluginOf(t, c, "work")
	restore := unsaveable(t, c)

	if _, errOut, code := c.run("", "rename", "work", "job"); code == 0 || !strings.Contains(errOut, "not saved") {
		t.Fatalf("rename: %d %s", code, errOut)
	}
	if !fileThere(filepath.Join(work, ".credentials.json")) || fileThere(c.p.ProfileDir("job")) {
		t.Fatal("rename moved the folder although config.toml still names work")
	}
	if got := readFile(t, filepath.Join(work, "plugins", "installed_plugins.json")); !strings.Contains(got, installed) {
		t.Fatalf("plugin paths not put back: %s", got)
	}

	if _, errOut, code := c.run("", "remove", "work", "--yes"); code == 0 || !strings.Contains(errOut, "not saved") {
		t.Fatalf("remove: %d %s", code, errOut)
	}
	if !fileThere(filepath.Join(work, ".credentials.json")) {
		t.Fatal("remove left the folder in the trash although config.toml still lists work")
	}

	if _, errOut, code := c.run("", "add", "extra", "--no-login"); code == 0 || !strings.Contains(errOut, "not saved") {
		t.Fatalf("add: %d %s", code, errOut)
	}
	if fileThere(c.p.ProfileDir("extra")) {
		t.Fatal("add left an unregistered folder, which blocks the next try")
	}

	restore()
	c.ok("rename", "work", "job")
	c.ok("add", "extra", "--no-login")
}

// The same with a config folder that cannot be written.
func TestReadOnlyConfigFolderLeavesTheFolders(t *testing.T) {
	if os.Geteuid() == 0 || runtime.GOOS == "windows" {
		t.Skip("root, and Windows, ignore the read-only folder this test relies on")
	}
	c := newTestCLI(t, nil)
	noClaude(t, c)
	work := c.p.ProfileDir("work")
	installed := pluginOf(t, c, "work")
	if err := os.Chmod(c.p.Config, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(c.p.Config, 0o700) })
	if _, _, code := c.run("", "rename", "work", "job"); code == 0 {
		t.Fatal("rename saved into a read-only folder")
	}
	if !fileThere(work) || fileThere(c.p.ProfileDir("job")) {
		t.Fatal("rename moved the folder although config.toml still names work")
	}
	if got := readFile(t, filepath.Join(work, "plugins", "installed_plugins.json")); !strings.Contains(got, installed) {
		t.Fatalf("plugin paths not put back: %s", got)
	}
	if _, _, code := c.run("", "remove", "work", "--yes"); code == 0 || !fileThere(work) {
		t.Fatal("remove left the folder in the trash")
	}
	if _, _, code := c.run("", "add", "extra", "--no-login"); code == 0 || fileThere(c.p.ProfileDir("extra")) {
		t.Fatal("add left an unregistered folder")
	}
}

// A login takes minutes in the browser: one whose record could not be
// saved stops before it starts.
func TestLoginStopsBeforeTheBrowserWhenTheConfigCannotBeSaved(t *testing.T) {
	c := newTestCLI(t, nil)
	marker := filepath.Join(t.TempDir(), "ran")
	fake := filepath.Join(t.TempDir(), "claude")
	if runtime.GOOS == "windows" {
		t.Skip("the fake claude is a shell script")
	}
	os.WriteFile(fake, []byte("#!/bin/sh\n: > '"+marker+"'\n"), 0o755)
	cfg := c.cfg()
	cfg.ClaudeBin = fake
	cfg.Save(c.p.ConfigFile())
	unsaveable(t, c)
	for _, args := range [][]string{{"login", "work"}, {"login", "--expired"}} {
		if _, errOut, code := c.run("", args...); code == 0 || !strings.Contains(errOut, "relay is read as part of") {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
	}
	if fileThere(marker) {
		t.Fatal("the login started")
	}
}
