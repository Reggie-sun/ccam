package launcher

import (
	"os"

	"golang.org/x/sys/unix"
)

// discardInput drops what was typed on the terminal f and not read yet, and
// reports whether there was a whole line of it (TIOCINQ counts no partial
// line while the terminal is in line mode; those bytes go too).
func discardInput(f *os.File) bool {
	fd := int(f.Fd())
	n, err := unix.IoctlGetInt(fd, unix.TIOCINQ)
	if err != nil || unix.IoctlSetInt(fd, unix.TCFLSH, unix.TCIFLUSH) != nil {
		return false // no terminal (a pipe has bytes waiting too, but no flush)
	}
	return n > 0
}
