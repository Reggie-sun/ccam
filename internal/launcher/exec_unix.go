//go:build !windows

package launcher

import (
	"os"
	"syscall"

	"github.com/jhyunwoo/ccam/internal/paths"
)

// passOn lists the signals Run catches; all but Ctrl+C go on to claude.
var passOn = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT}

// Exec replaces the current process with claude. It only returns on error.
func (pl *Plan) Exec(p paths.Paths) error {
	setTitle(pl.Profile.Name)
	MarkUsed(p, pl.Profile.Name)
	if BeforeExec != nil {
		BeforeExec()
	}
	return syscall.Exec(pl.Bin, pl.Argv, pl.Env)
}

// execSelf replaces the current process with ccam run with args.
func execSelf(self string, args []string) error { return syscall.Exec(self, args, os.Environ()) }
