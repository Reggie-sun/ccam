// Package fsutil holds the small file helpers every other package relies on.
package fsutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// AtomicWrite writes data via temp file + fsync + rename so readers never see a half file.
// If path is a symlink, the file it points to is written and the link kept.
func AtomicWrite(path string, data []byte, perm fs.FileMode) error {
	path = resolveLink(path)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp) // no-op after a successful rename
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Chmod(perm); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// resolveLink follows symlinks at path (even dangling ones) so a rename
// replaces the target, not the link. Gives up after a few hops.
func resolveLink(path string) string {
	for i := 0; i < 8; i++ {
		target, err := os.Readlink(path)
		if err != nil {
			return path
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(path), target)
		}
		path = target
	}
	return path
}

// MarshalJSON indents like the files Claude Code writes and keeps shell
// characters such as && and > readable (json.Marshal escapes them).
func MarshalJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteFast is for state that is cheap to lose (last_used, caches): no fsync.
func WriteFast(path string, data []byte) error {
	err := os.WriteFile(path, data, 0o600)
	if errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(filepath.Dir(path), 0o700); err == nil {
			err = os.WriteFile(path, data, 0o600)
		}
	}
	return err
}

// Exists reports whether path exists without following a final symlink.
func Exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

// Move renames src to dst, falling back to copy+delete across filesystems.
// Any other rename error (dst is a folder that is not empty, …) is returned
// instead of copying over it; like rename(2), an existing dst file is replaced. Relative links that pointed outside src
// are made absolute so they still resolve from the new place.
func Move(src, dst string) error { return MoveWithin(src, dst, src) }

// MoveWithin is Move for one piece of the folder root whose other pieces end
// up next to dst as well: relative links that stay inside root are kept.
func MoveWithin(src, dst, root string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	err := os.Rename(src, dst)
	if err == nil {
		return fixMovedLinks(src, dst, root)
	}
	if !errors.Is(err, syscall.EXDEV) {
		return err
	}
	if err := copyTree(src, dst, root); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

// UniqueDir returns base, or base-2, base-3 … : the first path that does not
// exist. Backup and trash names have one-second stamps that can repeat.
func UniqueDir(base string) string {
	path := base
	for i := 2; Exists(path); i++ {
		path = base + "-" + strconv.Itoa(i)
	}
	return path
}

// PortableLink returns the target to give a copy of the symlink at path
// (which reads link) when the tree root is copied or moved elsewhere: a
// relative target inside root stays relative, one that leaves root becomes
// absolute so it still points at the same place.
func PortableLink(path, link, root string) string {
	if filepath.IsAbs(link) {
		return link
	}
	abs := filepath.Join(filepath.Dir(path), link)
	if rel, err := filepath.Rel(root, abs); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return link
	}
	return abs
}

// fixMovedLinks rewrites, after src was renamed to dst, the relative links
// in dst whose target lay outside root.
func fixMovedLinks(src, dst, root string) error {
	return filepath.WalkDir(dst, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.Type()&fs.ModeSymlink == 0 {
			return err
		}
		link, err := os.Readlink(path)
		if err != nil || filepath.IsAbs(link) {
			return err
		}
		rel, _ := filepath.Rel(dst, path)
		fixed := PortableLink(filepath.Join(src, rel), link, root)
		if fixed == link {
			return nil
		}
		if err := os.Remove(path); err != nil {
			return err
		}
		return os.Symlink(fixed, path)
	})
}

// CopyTree copies a file, symlink or directory tree. Existing files in dst
// are overwritten. Missing parents of dst are created 0700, not with the
// source's (often 0755) mode.
func CopyTree(src, dst string) error { return copyTree(src, dst, src) }

func copyTree(src, dst, root string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		switch {
		case d.Type()&fs.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			os.Remove(target)
			return os.Symlink(PortableLink(path, link, root), target)
		case d.IsDir():
			info, err := d.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(target, info.Mode().Perm()|0o700)
		default:
			return CopyFile(path, target)
		}
	})
}

// CopyFile copies one regular file keeping its permission bits.
// io.Copy between *os.File uses copy_file_range on Linux.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	// OpenFile only applies perm to a new file; an overwritten one keeps its own.
	if err := out.Chmod(info.Mode().Perm()); err != nil {
		out.Close()
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// SameContent reports whether b is a file with exactly a's content. A
// missing or unreadable b is just different; only an error reading a counts.
func SameContent(a, b string) (bool, error) {
	ia, err := os.Stat(a)
	if err != nil {
		return false, err
	}
	ib, err := os.Stat(b)
	if err != nil || ia.Size() != ib.Size() || ib.IsDir() {
		return false, nil
	}
	x, err := os.ReadFile(a)
	if err != nil {
		return false, err
	}
	y, err := os.ReadFile(b)
	if err != nil {
		return false, nil
	}
	return bytes.Equal(x, y), nil
}

// Stamp returns a filesystem-safe timestamp for backup and trash folder names.
func Stamp(t time.Time) string { return t.Format("20060102-150405") }
