package launcher

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jhyunwoo/ccam/internal/handoff"
	"github.com/jhyunwoo/ccam/internal/paths"
)

// relayable reports whether claude args start an interactive session on a
// terminal, the only kind /ccam:switch can come from.
func relayable(args []string) bool {
	if fi, err := os.Stdin.Stat(); err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	for _, a := range args {
		if a == "-p" || a == "--print" || len(a) > 8 && a[:8] == "--print=" {
			return false
		}
	}
	return true
}

// relay runs claude as a child (relay = true in config.toml). When the
// session asked with /ccam:switch to move on, `ccam next --continue` takes
// its place and carries the conversation on as the other account. The
// arguments claude started with are not carried over.
func (pl *Plan) relay(p paths.Paths) error {
	// The hook writes this into its note, so only a /ccam:switch of this
	// very session is carried on.
	me := strconv.Itoa(os.Getpid())
	pl.Env = append(pl.Env, "CCAM_RELAY="+me)
	start := time.Now()
	code, err := pl.Run(p, os.Stdin, os.Stdout, os.Stderr)
	pl.Done()
	if err != nil {
		return err
	}
	if t, ok := handoff.Load(p.State, pl.Profile.Name); ok && t.Reason == handoff.Switch && t.Relay == me && t.At.After(start) {
		self, err := os.Executable()
		if err != nil {
			return err
		}
		args := []string{"ccam", "next", "--continue", "--note", pl.Profile.Name}
		if t.Target == "" {
			args = append(args, "--auto")
		}
		fmt.Fprintln(os.Stderr, "ccam: carrying the conversation on (ccam next --continue)")
		return execSelf(self, args)
	}
	if code != 0 {
		return ExitCode(code)
	}
	return nil
}
