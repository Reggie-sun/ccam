// Package cli holds the cobra commands. It is only reached for real
// subcommands; launching claude goes through app/launcher directly.
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/paths"
)

// Version is set at build time with -ldflags "-X .../cli.Version=...".
var Version = "dev"

// env is what every command works on.
type env struct {
	paths  paths.Paths
	stdin  io.Reader
	in     *bufio.Reader // the one buffered reader over stdin, shared by all prompts
	stdout io.Writer
	stderr io.Writer
	dryRun bool
}

func (e *env) load() (*config.Config, error) { return config.Load(e.paths.ConfigFile()) }

func (e *env) save(cfg *config.Config) error { return cfg.Save(e.paths.ConfigFile()) }

func (e *env) printf(format string, a ...any) { fmt.Fprintf(e.stdout, format, a...) }

// confirm asks a yes/no question; anything but y/yes is no.
func (e *env) confirm(question string) bool {
	fmt.Fprintf(e.stdout, "%s [y/N] ", question)
	line, _ := e.in.ReadString('\n')
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "y" || ans == "yes"
}

// ccamPath is the real ccam binary, which every shim links to.
func ccamPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exe)
}

func newRoot(e *env) *cobra.Command {
	root := &cobra.Command{
		Use:           "ccam",
		Short:         "Claude Code multi-account manager",
		Long:          "ccam adds and switches between several Claude Code accounts.\n\n  ccam <name|number> [claude args]  run Claude Code as that account\n  <shortcut> [claude args]          run by a shortcut command (ccam shortcut)",
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHome(e)
		},
	}
	root.SetIn(e.stdin)
	root.SetOut(e.stdout)
	root.SetErr(e.stderr)
	root.PersistentFlags().BoolVar(&e.dryRun, "dry-run", false, "print what would happen without doing it")
	root.CompletionOptions.DisableDefaultCmd = true
	root.ValidArgsFunction = profileNames(e) // `ccam <tab>` lists accounts
	root.AddCommand(
		addCmd(e), listCmd(e), runCmd(e), useCmd(e), removeCmd(e), renameCmd(e),
		loginCmd(e), logoutCmd(e), shimsCmd(e), versionCmd(e), shareCmd(e), doctorCmd(e), mcpCmd(e), nextCmd(e), usageCmd(e), statuslineCmd(e), importCmd(e), migrateCmd(e), shortcutCmd(e), completionCmd(e),
	)
	for _, c := range root.Commands() {
		switch c.Name() {
		case "run", "use", "remove", "rename", "login", "logout", "share", "shortcut":
			c.ValidArgsFunction = profileNames(e)
		}
	}
	return root
}

func versionCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		Run:   func(cmd *cobra.Command, args []string) { e.printf("ccam %s\n", Version) },
	}
}

// Execute runs the CLI with args (without argv[0]) and returns the exit code.
func Execute(args []string) int {
	return ExecuteWith(args, paths.FromEnv(), os.Stdin, os.Stdout, os.Stderr)
}

func ExecuteWith(args []string, p paths.Paths, stdin io.Reader, stdout, stderr io.Writer) int {
	e := &env{paths: p, stdin: stdin, in: bufio.NewReader(stdin), stdout: stdout, stderr: stderr}
	root := newRoot(e)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		var exit exitError
		if errors.As(err, &exit) {
			return int(exit)
		}
		fmt.Fprintln(stderr, "ccam:", err)
		return 1
	}
	return 0
}

// exitError lets a command end with a specific code and no message.
type exitError int

func (e exitError) Error() string { return fmt.Sprintf("exit %d", int(e)) }
