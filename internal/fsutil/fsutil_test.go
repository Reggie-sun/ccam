package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicWriteCreatesDirsAndPerm(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "f.toml")
	if err := AtomicWrite(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := AtomicWrite(path, []byte("y"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	info, _ := os.Stat(path)
	if string(got) != "y" || info.Mode().Perm() != 0o600 {
		t.Fatalf("got %q %v", got, info.Mode())
	}
	entries, _ := os.ReadDir(filepath.Dir(path))
	if len(entries) != 1 {
		t.Fatalf("temp files left behind: %v", entries)
	}
}

func TestCopyTreeKeepsSymlinks(t *testing.T) {
	src := filepath.Join(t.TempDir(), "src")
	os.MkdirAll(filepath.Join(src, "sub"), 0o700)
	os.WriteFile(filepath.Join(src, "sub", "f"), []byte("data"), 0o640)
	os.Symlink("sub/f", filepath.Join(src, "link"))
	dst := filepath.Join(t.TempDir(), "dst")
	if err := CopyTree(src, dst); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(dst, "sub", "f")); string(b) != "data" {
		t.Fatalf("copy content %q", b)
	}
	if l, _ := os.Readlink(filepath.Join(dst, "link")); l != "sub/f" {
		t.Fatalf("symlink %q", l)
	}
}

func TestMove(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a"), []byte("1"), 0o600)
	if err := Move(filepath.Join(dir, "a"), filepath.Join(dir, "x", "b")); err != nil {
		t.Fatal(err)
	}
	if Exists(filepath.Join(dir, "a")) || !Exists(filepath.Join(dir, "x", "b")) {
		t.Fatal("move did not happen")
	}
}

func TestAtomicWriteKeepsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "real.json")
	link := filepath.Join(dir, "link.json")
	os.WriteFile(target, []byte("old"), 0o600)
	os.Symlink(target, link)
	if err := AtomicWrite(link, []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Lstat(link); info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("link replaced by a regular file")
	}
	if b, _ := os.ReadFile(target); string(b) != "new" {
		t.Fatalf("target %q", b)
	}
	dangling := filepath.Join(dir, "d.json")
	os.Symlink("sub/missing.json", dangling) // relative, target dir missing
	if err := AtomicWrite(dangling, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "sub", "missing.json")); string(b) != "x" {
		t.Fatal("dangling relative link not followed")
	}
}

func TestMarshalJSONKeepsShellCharacters(t *testing.T) {
	out, err := MarshalJSON(map[string]string{"command": "a && b > c"})
	if err != nil || string(out) != "{\n  \"command\": \"a && b > c\"\n}\n" {
		t.Fatalf("%q %v", out, err)
	}
}
