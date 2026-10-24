package fsutil

import (
	"errors"

	"golang.org/x/sys/windows"
)

// crossDevice reports whether a rename failed because the two paths are on
// different volumes. Windows says ERROR_NOT_SAME_DEVICE; syscall.EXDEV is
// only Go's stand-in there, which no call returns.
func crossDevice(err error) bool { return errors.Is(err, windows.ERROR_NOT_SAME_DEVICE) }
