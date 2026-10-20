package fsutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// perms is whether the platform has Unix permission bits (Windows does not).
var perms = runtime.GOOS != "windows"

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
	if string(got) != "y" || perms && info.Mode().Perm() != 0o600 {
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
	if l, _ := os.Readlink(filepath.Join(dst, "link")); l != filepath.FromSlash("sub/f") {
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

func TestMoveRefusesToOverwriteExistingDir(t *testing.T) {
	root := t.TempDir()
	src, dst := filepath.Join(root, "src"), filepath.Join(root, "dst")
	os.MkdirAll(src, 0o700)
	os.WriteFile(filepath.Join(src, "f"), []byte("new"), 0o600)
	os.MkdirAll(dst, 0o700)
	os.WriteFile(filepath.Join(dst, "f"), []byte("old"), 0o600)
	if err := Move(src, dst); err == nil {
		t.Fatal("moved over a non-empty folder")
	}
	if b, _ := os.ReadFile(filepath.Join(dst, "f")); string(b) != "old" {
		t.Fatalf("dst overwritten: %q", b)
	}
	if !Exists(filepath.Join(src, "f")) {
		t.Fatal("src deleted")
	}
}

func TestUniqueDir(t *testing.T) {
	base := filepath.Join(t.TempDir(), "x")
	if got := UniqueDir(base); got != base {
		t.Fatalf("free name: %s", got)
	}
	os.Mkdir(base, 0o700)
	os.Mkdir(base+"-2", 0o700)
	if got := UniqueDir(base); got != base+"-3" {
		t.Fatalf("taken names: %s", got)
	}
}

// Relative links that leave the moved or copied tree (dotfiles managed by
// stow) must still resolve from the new place; links inside stay relative.
func TestMoveAndCopyKeepOutsideLinksWorking(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "dotfiles", "skill"), 0o700)
	os.WriteFile(filepath.Join(root, "dotfiles", "skill", "SKILL.md"), []byte("s"), 0o600)
	for _, move := range []bool{true, false} {
		src := filepath.Join(root, "home", ".claude", "skills")
		os.MkdirAll(src, 0o700)
		os.Symlink("../../../dotfiles/skill", filepath.Join(src, "out"))
		os.WriteFile(filepath.Join(src, "in.md"), []byte("i"), 0o600)
		os.Symlink("in.md", filepath.Join(src, "inside"))
		dst := filepath.Join(root, "data", "groups", "default", "skills")
		var err error
		if move {
			err = Move(src, dst)
		} else {
			err = CopyTree(src, dst)
		}
		if err != nil {
			t.Fatal(err)
		}
		if b, err := os.ReadFile(filepath.Join(dst, "out", "SKILL.md")); err != nil || string(b) != "s" {
			t.Fatalf("move=%v: outside link broken: %v", move, err)
		}
		if l, _ := os.Readlink(filepath.Join(dst, "inside")); l != "in.md" {
			t.Fatalf("move=%v: inside link changed to %q", move, l)
		}
		os.RemoveAll(filepath.Join(root, "data"))
		os.RemoveAll(filepath.Join(root, "home"))
	}
}

func TestCopyTreeParentsArePrivate(t *testing.T) {
	if !perms {
		t.Skip("no permission bits")
	}
	root := t.TempDir()
	src := filepath.Join(root, "src")
	os.MkdirAll(src, 0o755)
	os.WriteFile(filepath.Join(src, "f"), []byte("x"), 0o644)
	dst := filepath.Join(root, "data", "groups", "default", "agents")
	if err := CopyTree(src, dst); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"data", "data/groups", "data/groups/default"} {
		if info, _ := os.Stat(filepath.Join(root, d)); info.Mode().Perm() != 0o700 {
			t.Errorf("%s is %v", d, info.Mode().Perm())
		}
	}
}

func TestCopyFileSetsPermOnExistingFile(t *testing.T) {
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "src"), filepath.Join(dir, "dst")
	os.WriteFile(src, []byte("x"), 0o600)
	os.WriteFile(dst, []byte("old"), 0o644)
	if err := CopyFile(src, dst); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(dst); perms && info.Mode().Perm() != 0o600 {
		t.Fatalf("perm %v", info.Mode().Perm())
	}
}
