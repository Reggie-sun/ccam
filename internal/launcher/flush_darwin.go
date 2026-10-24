package launcher

import (
	"os"

	"golang.org/x/sys/unix"
)

// FIONREAD and FREAD (sys/filio.h, sys/fcntl.h), which x/sys/unix lacks.
const (
	fionread = 0x4004667f
	fread    = 0x1
)

// discardInput drops what was typed on the terminal f and not read yet, and
// reports whether there was a whole line of it (FIONREAD counts no partial
// line while the terminal is in line mode; those bytes go too).
func discardInput(f *os.File) bool {
	fd := int(f.Fd())
	n, err := unix.IoctlGetInt(fd, fionread)
	if err != nil || unix.IoctlSetPointerInt(fd, unix.TIOCFLUSH, fread) != nil {
		return false // no terminal (a pipe has bytes waiting too, but no flush)
	}
	return n > 0
}
