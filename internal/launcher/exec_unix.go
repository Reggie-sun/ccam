//go:build !windows

package launcher

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"

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
	err := syscall.Exec(pl.Bin, pl.Argv, pl.Env)
	return &os.PathError{Op: "exec", Path: pl.Bin, Err: err}
}

// execSelf replaces the current process with ccam run with args.
func execSelf(self string, args []string) error { return syscall.Exec(self, args, os.Environ()) }

// ExecProgram replaces the current process with another program, such as
// ccusage: bin is its path, argv starts with its name. It only returns on
// error.
func ExecProgram(bin string, argv, env []string) error {
	if BeforeExec != nil {
		BeforeExec()
	}
	err := syscall.Exec(bin, argv, env)
	return &os.PathError{Op: "exec", Path: bin, Err: err}
}

// foreground reports whether this process's group is the foreground one of
// the terminal fd: only that group may read it.
func foreground(fd uintptr) bool {
	pgrp, err := unix.IoctlGetInt(int(fd), unix.TIOCGPGRP)
	return err == nil && pgrp == unix.Getpgrp()
}

// shellCommand runs command in sh.
func shellCommand(ctx context.Context, command string) *exec.Cmd {
	return exec.CommandContext(ctx, "sh", "-c", command)
}

// commandLine is for Windows' dry runs.
func commandLine(bin string, args []string) (string, string, error) {
	return bin, strings.Join(args, " "), nil
}

// isBatch is for Windows' batch files.
func isBatch(string) bool { return false }

// batchShell is for Windows' batch files.
func batchShell(os.FileInfo) bool { return false }

func batchCommand(ctx context.Context, bin string, args []string) (*exec.Cmd, error) {
	return exec.CommandContext(ctx, bin, args...), nil
}
