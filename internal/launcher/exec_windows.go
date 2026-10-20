package launcher

import (
	"os"

	"github.com/jhyunwoo/ccam/internal/paths"
)

// passOn lists the signals Run catches. Windows has only Ctrl+C, which the
// console delivers to claude itself.
var passOn = []os.Signal{os.Interrupt}

// Exec runs claude and waits for it: Windows has no exec, so ccam stays
// in between and ends with claude's exit code (an ExitCode error).
func (pl *Plan) Exec(p paths.Paths) error {
	if BeforeExec != nil {
		BeforeExec()
	}
	code, err := pl.Run(p, os.Stdin, os.Stdout, os.Stderr)
	pl.Done()
	if err != nil {
		return err
	}
	if code != 0 {
		return ExitCode(code)
	}
	return nil
}
