package fsutil

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"
)

// A rename from one drive to another fails with ERROR_NOT_SAME_DEVICE;
// Move then copies and removes the source, as it does for EXDEV elsewhere.
func TestMoveCopiesBetweenDrives(t *testing.T) {
	orig := rename
	rename = func(oldpath, newpath string) error {
		return &os.LinkError{Op: "rename", Old: oldpath, New: newpath, Err: windows.ERROR_NOT_SAME_DEVICE}
	}
	t.Cleanup(func() { rename = orig })
	root := t.TempDir()
	src := filepath.Join(root, "src")
	os.MkdirAll(filepath.Join(src, "sub"), 0o700)
	os.WriteFile(filepath.Join(src, "sub", "f"), []byte("x"), 0o600)
	dst := filepath.Join(root, "dst")
	if err := Move(src, dst); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(filepath.Join(dst, "sub", "f")); err != nil || string(b) != "x" {
		t.Fatalf("copied: %q %v", b, err)
	}
	if Exists(src) {
		t.Fatal("the source is still there")
	}
}
