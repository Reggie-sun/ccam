// Package app is the process entry point. It routes `ccs <name>`, `ccs run`
// and claudeN shims straight to the launcher, and only hands real
// subcommands to the cobra CLI.
package app

import (
	"fmt"
	"os"
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

func parse(argv []string) route {
	if sel, ok := shim.Selector(argv[0]); ok {
		return route{kind: routeLaunch, sel: sel, args: argv[1:]}
	}
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
		if len(r) > 0 && !strings.HasPrefix(r[0], "-") {
			return route{kind: routeLaunch, sel: r[0], args: r[1:], dry: dry}
		}
		return route{kind: routeLaunch, args: r, dry: dry}
	case rest[0] == "-h" || rest[0] == "--help" || rest[0] == "--version":
		return route{kind: routeCLI, args: withDry(rest, dry)}
	case strings.HasPrefix(rest[0], "-"):
		return route{kind: routeLaunch, args: rest, dry: dry}
	case config.Reserved[rest[0]]:
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
