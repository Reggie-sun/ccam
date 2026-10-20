//go:build !windows

package fsutil

import "os"

// junctions is whether the platform has directory junctions.
const junctions = false

// Link makes path a symbolic link to target.
func Link(target, path string, dir bool) error { return os.Symlink(target, path) }
