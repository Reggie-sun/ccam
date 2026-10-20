package fsutil

import (
	"io/fs"
	"os"
	"path/filepath"
)

// LinkTarget returns where the link at path points: a symbolic link, or on
// Windows a directory junction, which Lstat reports as an irregular file.
func LinkTarget(path string) (string, bool) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", false
	}
	return linkTarget(path, info.Mode())
}

// LinkTargetOf is LinkTarget for a path already Lstat-ed into info.
func LinkTargetOf(path string, info fs.FileInfo) (string, bool) {
	return linkTarget(path, info.Mode())
}

// IsLink is LinkTarget for a directory walk entry.
func IsLink(path string, d fs.DirEntry) bool {
	_, ok := linkTarget(path, d.Type())
	return ok
}

func linkTarget(path string, mode fs.FileMode) (string, bool) {
	if mode&fs.ModeSymlink == 0 && !(junctions && mode&fs.ModeIrregular != 0) {
		return "", false
	}
	target, err := os.Readlink(path)
	return target, err == nil
}

// CopyLink recreates the link at path, pointing to link, at target.
func CopyLink(link, target string) error {
	os.Remove(target)
	resolved := link
	if !filepath.IsAbs(link) {
		resolved = filepath.Join(filepath.Dir(target), link)
	}
	info, err := os.Stat(resolved)
	return Link(link, target, err == nil && info.IsDir())
}
