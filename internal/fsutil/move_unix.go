//go:build !windows

package fsutil

import (
	"errors"
	"syscall"
)

// crossDevice reports whether a rename failed because the two paths are on
// different filesystems.
func crossDevice(err error) bool { return errors.Is(err, syscall.EXDEV) }
