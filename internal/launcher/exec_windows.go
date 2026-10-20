package launcher

import (
	"errors"
	"os"
	"os/exec"

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

// execSelf runs ccam with args and ends with its exit code, as Windows has
// no exec.
func execSelf(self string, args []string) error {
	c := exec.Command(self, args[1:]...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := c.Run()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ExitCode(ee.ExitCode())
	}
	return err
}
