package fsutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

// withoutSymlinks makes Link behave as on Windows without Developer Mode.
func withoutSymlinks(t *testing.T) {
	orig := symlink
	symlink = func(string, string) error { return &os.LinkError{Op: "symlink", Err: windows.ERROR_PRIVILEGE_NOT_HELD} }
	t.Cleanup(func() { symlink = orig })
}

func TestLinkFallsBackToAJunctionForFolders(t *testing.T) {
	withoutSymlinks(t)
	root := t.TempDir()
	group := filepath.Join(root, "group", "skills")
	os.MkdirAll(group, 0o700)
	os.WriteFile(filepath.Join(group, "s.md"), []byte("skill"), 0o600)
	link := filepath.Join(root, "profile", "skills")
	os.MkdirAll(filepath.Dir(link), 0o700)
	if err := Link(group, link, true); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(filepath.Join(link, "s.md")); err != nil || string(b) != "skill" {
		t.Fatalf("through the junction: %q %v", b, err)
	}
	if target, ok := LinkTarget(link); !ok || target != group {
		t.Fatalf("LinkTarget %q %v", target, ok)
	}
	// Removing the profile folder must not reach into the group.
	if err := os.RemoveAll(filepath.Join(root, "profile")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(group, "s.md")); err != nil {
		t.Fatalf("RemoveAll followed the junction: %v", err)
	}
}

func TestLinkRefusesFilesWithoutDeveloperMode(t *testing.T) {
	withoutSymlinks(t)
	root := t.TempDir()
	target := filepath.Join(root, "settings.json")
	os.WriteFile(target, []byte("{}"), 0o600)
	err := Link(target, filepath.Join(root, "link.json"), false)
	if err == nil || !strings.Contains(err.Error(), "Developer Mode") {
		t.Fatalf("%v", err)
	}
}

func TestCanLinkFiles(t *testing.T) {
	if err := CanLinkFiles(t.TempDir()); err != nil && !strings.Contains(err.Error(), "Developer Mode") {
		t.Fatal(err)
	}
	withoutSymlinks(t)
	dir := t.TempDir()
	if err := CanLinkFiles(dir); err == nil || !strings.Contains(err.Error(), "Developer Mode") {
		t.Fatal(err)
	}
	if left, _ := os.ReadDir(dir); len(left) != 0 {
		t.Fatalf("left %v behind", left)
	}
}

func TestCopyTreeKeepsJunctions(t *testing.T) {
	withoutSymlinks(t)
	root := t.TempDir()
	target := filepath.Join(root, "target")
	os.MkdirAll(target, 0o700)
	src := filepath.Join(root, "src")
	os.MkdirAll(src, 0o700)
	if err := Link(target, filepath.Join(src, "j"), true); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(root, "dst")
	if err := CopyTree(src, dst); err != nil {
		t.Fatal(err)
	}
	if got, ok := LinkTarget(filepath.Join(dst, "j")); !ok || got != target {
		t.Fatalf("%q %v", got, ok)
	}
}
