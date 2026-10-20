// Package app is the process entry point. It routes `ccam <name>`, `ccam run`
// and account shortcuts straight to the launcher, and only hands real
// subcommands to the cobra CLI.
package app

import (
	"errors"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/jhyunwoo/ccam/internal/cli"
	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/shim"
)

type routeKind int

const (
	routeLaunch routeKind = iota
	routeCLI
)

type route struct {
	kind     routeKind
	sel      string   // profile selector for routeLaunch; "" = default profile
	args     []string // claude args for routeLaunch, cobra args for routeCLI
	dry      bool
	fallback bool // --fallback: run again as another account at a usage limit
}

// launchFlags takes ccam's own launch flags (--dry-run, --fallback, in any
// order) off the front of args.
func launchFlags(args []string, r *route) []string {
	for len(args) > 0 {
		switch args[0] {
		case "--dry-run":
			r.dry = true
		case "--fallback":
			r.fallback = true
		default:
			return args
		}
		args = args[1:]
	}
	return args
}

// parse routes `ccam …`. Shortcuts (any other argv[0]) are resolved in Main
// against the config first.
func parse(argv []string) route {
	var flags route
	rest := launchFlags(argv[1:], &flags)
	dry, fallback := flags.dry, flags.fallback
	switch {
	case len(rest) == 0:
		if dry || fallback {
			return route{kind: routeLaunch, dry: dry, fallback: fallback}
		}
		return route{kind: routeCLI}
	case rest[0] == "run":
		// ccam's own flags may come right after run (it is a ccam command);
		// anything after the account name goes to claude untouched.
		r := launchFlags(rest[1:], &flags)
		dry, fallback = flags.dry, flags.fallback
		if len(r) > 0 && (r[0] == "-h" || r[0] == "--help") {
			// cobra does not parse flags of run, so ask for its help directly.
			return route{kind: routeCLI, args: []string{"help", "run"}}
		}
		if len(r) > 0 && !strings.HasPrefix(r[0], "-") {
			return route{kind: routeLaunch, sel: r[0], args: r[1:], dry: dry, fallback: fallback}
		}
		return route{kind: routeLaunch, args: r, dry: dry, fallback: fallback}
	case fallback && (rest[0] == "--version" || rest[0] == "-v" || rest[0] == "-h" || rest[0] == "--help" || config.Reserved[rest[0]]):
		return route{kind: routeCLI, args: []string{"--fallback"}} // cobra explains the flag is unknown there
	case rest[0] == "--version" || rest[0] == "-v":
		// The version command prints the same line; a custom cobra version
		// template would pull text/template in and add 2 MB to the binary.
		return route{kind: routeCLI, args: []string{"version"}}
	case rest[0] == "-h" || rest[0] == "--help":
		return route{kind: routeCLI, args: withDry(rest, dry)}
	case strings.HasPrefix(rest[0], "-"):
		return route{kind: routeLaunch, args: rest, dry: dry, fallback: fallback}
	case config.Reserved[rest[0]] || strings.HasPrefix(rest[0], "__"): // __complete is cobra's
		return route{kind: routeCLI, args: withDry(rest, dry)}
	default:
		return route{kind: routeLaunch, sel: rest[0], args: rest[1:], dry: dry, fallback: fallback}
	}
}

func withDry(args []string, dry bool) []string {
	if !dry {
		return args
	}
	return append([]string{"--dry-run"}, args...)
}

// Main runs ccam and returns the process exit code.
func Main(argv []string) int {
	if path := os.Getenv("CCAM_PPROF"); path != "" {
		if stop := startProfile(path); stop != nil {
			defer stop()
			launcher.BeforeExec = stop
		}
	}
	// Run as a shortcut: argv[0] names the account. The config is the one
	// file the launch reads anyway, so it is loaded once and passed on.
	if base := shim.Name(argv[0]); base != "ccam" {
		p := paths.FromEnv()
		if cfg, err := config.Load(p.ConfigFile()); err == nil {
			if sel, ok := shim.Lookup(cfg, base); ok {
				return exitCode(launcher.LaunchWith(p, cfg, sel, argv[1:], false, os.Stdout))
			}
			// A link to ccam in the bin folder that no account claims (its
			// shortcut was changed in config.toml by hand) must not fall
			// through and open the default account.
			if self, err := os.Executable(); err == nil && shim.Linked(p.Bin, self, base) {
				fmt.Fprintf(os.Stderr, "ccam: %s is not the shortcut of any account (see `ccam list`, `ccam shortcut`)\n", base)
				return 1
			}
		}
	}
	r := parse(argv)
	if r.kind == routeCLI {
		return cli.Execute(r.args)
	}
	if r.fallback {
		return cli.Fallback(paths.FromEnv(), r.sel, r.args, r.dry)
	}
	return exitCode(launcher.Launch(paths.FromEnv(), r.sel, r.args, r.dry, os.Stdout))
}

// exitCode turns a launch result into ccam's exit code. Where ccam waits
// for claude (Windows), claude's own code comes back as an ExitCode.
func exitCode(err error) int {
	var code launcher.ExitCode
	switch {
	case err == nil:
		return 0
	case errors.As(err, &code):
		return int(code)
	}
	fmt.Fprintln(os.Stderr, "ccam:", err)
	return 1
}

// startProfile writes a CPU profile to path (CCAM_PPROF) until stop is called.
func startProfile(path string) (stop func()) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ccam: CCAM_PPROF:", err)
		return nil
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		fmt.Fprintln(os.Stderr, "ccam: CCAM_PPROF:", err)
		return nil
	}
	done := false
	return func() {
		if !done {
			done = true
			pprof.StopCPUProfile()
			f.Close()
		}
	}
}
