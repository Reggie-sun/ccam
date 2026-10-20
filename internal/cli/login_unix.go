//go:build !windows

package cli

import (
	"os"
	"syscall"

	"github.com/charmbracelet/x/term"
	"golang.org/x/sys/unix"
)

// endSignals end ccam during a masked login.
var endSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGHUP}

// noEcho is saved without echo and line editing, one byte per read.
func noEcho(saved *term.State) *term.State {
	masked := *saved
	masked.Lflag &^= unix.ECHO | unix.ICANON
	masked.Cc[unix.VMIN], masked.Cc[unix.VTIME] = 1, 0
	return &masked
}

// resignal ends ccam with the signal it caught, as if it had not.
func resignal(s os.Signal) { _ = syscall.Kill(os.Getpid(), s.(syscall.Signal)) }
