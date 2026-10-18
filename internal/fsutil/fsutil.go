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
func Move(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	if err := CopyTree(src, dst); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

// CopyTree copies a file, symlink or directory tree. Existing files in dst are overwritten.
func CopyTree(src, dst string) error {
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
			return os.Symlink(link, target)
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
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// Stamp returns a filesystem-safe timestamp for backup and trash folder names.
func Stamp(t time.Time) string { return t.Format("20060102-150405") }
