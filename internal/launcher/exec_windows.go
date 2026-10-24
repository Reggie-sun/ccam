package launcher

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"unicode"

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
// no exec. args[0] stays the child's argv[0]: run through a shortcut, self
// is the shortcut's file, whose name would make the child that account.
// While it waits, Ctrl+C must not end this ccam: the console sends it to
// every process on it, and Windows ends one that has no handler, which
// would give the shell its prompt back while claude still runs.
func execSelf(self string, args []string) error {
	c := &exec.Cmd{Path: self, Args: args, Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)
	err := c.Run()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ExitCode(ee.ExitCode())
	}
	return err
}

// foreground is always true: a Windows console has no background jobs.
func foreground(uintptr) bool { return true }

// shellCommand runs command in cmd.exe as typed there. Go would quote it
// as a program's argument (\" for a quote), which cmd.exe does not read
// back, so the line is given as it is: /s takes off just the outer quotes,
// /d skips AutoRun commands.
func shellCommand(ctx context.Context, command string) *exec.Cmd {
	c := exec.CommandContext(ctx, cmdExe())
	c.SysProcAttr = &syscall.SysProcAttr{CmdLine: `cmd /d /s /c "` + command + `"`}
	return c
}

// cmdExe is Windows' cmd.exe, never one in the working folder.
func cmdExe() string {
	return filepath.Join(cmp.Or(os.Getenv("SystemRoot"), `C:\Windows`), "System32", "cmd.exe")
}

// commandLine is the program and the command line a launch of bin with
// args starts: Go's quoting for a program, cmd.exe's line for a batch file
// (or why a batch file cannot be started with these args).
func commandLine(bin string, args []string) (string, string, error) {
	if isBatch(bin) {
		c, err := batchCommand(context.Background(), bin, args)
		if err != nil {
			return "", "", err
		}
		return c.Path, strings.TrimPrefix(c.SysProcAttr.CmdLine, "cmd "), nil
	}
	quoted := make([]string, len(args))
	for i, a := range args {
		quoted[i] = syscall.EscapeArg(a)
	}
	return bin, strings.Join(quoted, " "), nil
}

// isBatch reports whether bin is a batch file, which Windows runs in cmd.exe.
func isBatch(bin string) bool {
	ext := filepath.Ext(bin)
	return strings.EqualFold(ext, ".cmd") || strings.EqualFold(ext, ".bat")
}

// batchShell reports whether fi is the cmd.exe that batchCommand runs a
// batch claude in, which runs Claude Code as its child.
func batchShell(fi os.FileInfo) bool {
	cmd, err := os.Stat(cmdExe())
	return err == nil && os.SameFile(fi, cmd)
}

// batchCommand runs the batch file bin (claude.cmd, as npm installs it)
// with args. Started as it is, Windows hands the command line to cmd.exe,
// which looks argv[0] up again by itself, the working folder first, and
// reads & | < > ^ and %VAR% in the arguments as its own syntax. So cmd.exe
// gets a line built here: the file's full path, and each argument quoted as
// Rust's standard library does since CVE-2024-24576 (" doubled, each %
// written so that no variable is expanded). A line break would end the
// command and cannot be quoted, so such an argument is refused.
func batchCommand(ctx context.Context, bin string, args []string) (*exec.Cmd, error) {
	var b strings.Builder
	b.WriteString(`/e:ON /v:OFF /d /c ""`)
	writeBatchText(&b, bin)
	b.WriteByte('"')
	for _, a := range args {
		if strings.ContainsAny(a, "\r\n\x00") {
			return nil, fmt.Errorf("claude is a batch file (%s), and cmd.exe cannot be given an argument with a line break. Install the native Claude Code (`claude install`), or set claude_bin to a claude.exe", bin)
		}
		b.WriteByte(' ')
		writeBatchArg(&b, a)
	}
	b.WriteByte('"')
	c := exec.CommandContext(ctx, cmdExe())
	c.SysProcAttr = &syscall.SysProcAttr{CmdLine: `cmd ` + b.String()}
	return c, nil
}

// writeBatchArg quotes a for cmd.exe and the program the batch file passes
// it on to (%*): in quotes unless it is plain, a " doubled and the
// backslashes before it too.
func writeBatchArg(b *strings.Builder, a string) {
	quote := a == "" || strings.HasSuffix(a, `\`)
	for _, r := range a {
		plain := r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune(`#$*+-./:?@\_`, r)
		if r < 0x80 && !plain || unicode.IsControl(r) {
			quote = true
		}
	}
	if quote {
		b.WriteByte('"')
	}
	backslashes := 0
	for _, r := range a {
		switch r {
		case '\\':
			backslashes++
			b.WriteRune(r)
			continue
		case '"':
			b.WriteString(strings.Repeat(`\`, backslashes))
			b.WriteString(`""`)
		default:
			writeBatchText(b, string(r))
		}
		backslashes = 0
	}
	if quote {
		b.WriteString(strings.Repeat(`\`, backslashes))
		b.WriteByte('"')
	}
}

// writeBatchText writes s with each % as %%cd:~,%, which cmd.exe reads as a
// % and an empty part of %cd%, so no %NAME% in s is expanded.
func writeBatchText(b *strings.Builder, s string) {
	b.WriteString(strings.ReplaceAll(s, "%", "%%cd:~,%"))
}
