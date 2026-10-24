package cli

import (
	"os"
	"os/exec"

	"github.com/charmbracelet/x/term"
	"golang.org/x/sys/windows"
)

// endSignals end ccam during a masked login: Ctrl+C is the only one.
var endSignals = []os.Signal{os.Interrupt}

// noEcho is saved without echo and line input, so each key arrives as it is
// typed. Processed input stays on, so Ctrl+C still works.
func noEcho(saved *term.State) *term.State {
	masked := *saved
	masked.Mode &^= windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT
	return &masked
}

// ownGroup leaves c as it is: Cancel ends the shell only.
func ownGroup(*exec.Cmd) {}

// resignal ends ccam as Ctrl+C would have.
func resignal(os.Signal) { os.Exit(130) }
