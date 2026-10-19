// Package app is the process entry point. It routes `ccam <name>`, `ccam run`
// and account shortcuts straight to the launcher, and only hands real
// subcommands to the cobra CLI.
package app

import (
	"fmt"
	"os"
	"path/filepath"
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
	kind routeKind
	sel  string   // profile selector for routeLaunch; "" = default profile
	args []string // claude args for routeLaunch, cobra args for routeCLI
	dry  bool
}

// parse routes `ccam …`. Shortcuts (any other argv[0]) are resolved in Main
// against the config first.
func parse(argv []string) route {
	rest := argv[1:]
	dry := false
	if len(rest) > 0 && rest[0] == "--dry-run" {
		dry, rest = true, rest[1:]
	}
	switch {
	case len(rest) == 0:
		if dry {
			return route{kind: routeLaunch, dry: true}
		}
		return route{kind: routeCLI}
	case rest[0] == "run":
		r := rest[1:]
		// ccam's own flags may come right after run (it is a ccam command);
		// anything after the account name goes to claude untouched.
		if len(r) > 0 && r[0] == "--dry-run" {
			dry, r = true, r[1:]
		}
		if len(r) > 0 && (r[0] == "-h" || r[0] == "--help") {
			// cobra does not parse flags of run, so ask for its help directly.
			return route{kind: routeCLI, args: []string{"help", "run"}}
		}
		if len(r) > 0 && !strings.HasPrefix(r[0], "-") {
			return route{kind: routeLaunch, sel: r[0], args: r[1:], dry: dry}
		}
		return route{kind: routeLaunch, args: r, dry: dry}
	case rest[0] == "--version" || rest[0] == "-v":
		// The version command prints the same line; a custom cobra version
		// template would pull text/template in and add 2 MB to the binary.
		return route{kind: routeCLI, args: []string{"version"}}
	case rest[0] == "-h" || rest[0] == "--help":
		return route{kind: routeCLI, args: withDry(rest, dry)}
	case strings.HasPrefix(rest[0], "-"):
		return route{kind: routeLaunch, args: rest, dry: dry}
	case config.Reserved[rest[0]] || strings.HasPrefix(rest[0], "__"): // __complete is cobra's
		return route{kind: routeCLI, args: withDry(rest, dry)}
	default:
		return route{kind: routeLaunch, sel: rest[0], args: rest[1:], dry: dry}
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
	if base := filepath.Base(argv[0]); base != "ccam" {
		p := paths.FromEnv()
		if cfg, err := config.Load(p.ConfigFile()); err == nil {
			if sel, ok := shim.Lookup(cfg, base); ok {
				if err := launcher.LaunchWith(p, cfg, sel, argv[1:], false, os.Stdout); err != nil {
					fmt.Fprintln(os.Stderr, "ccam:", err)
					return 1
				}
				return 0
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
	if err := launcher.Launch(paths.FromEnv(), r.sel, r.args, r.dry, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "ccam:", err)
		return 1
	}
	return 0
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
