package launcher

import (
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/session"
	"github.com/jhyunwoo/ccam/internal/trace"
)

// Launch is the whole `ccam <name>` sequence: load config, prepare, check
// links, apply MCP sharing, take locks, then exec (or print when dry).
func Launch(p paths.Paths, sel string, args []string, dry bool, stdout io.Writer) error {
	trace.Step("start")
	cfg, err := config.Load(p.ConfigFile())
	if err != nil {
		return err
	}
	trace.Step("load config")
	return LaunchWith(p, cfg, sel, args, dry, stdout)
}

// LaunchWith is Launch with an already loaded config.
func LaunchWith(p paths.Paths, cfg *config.Config, sel string, args []string, dry bool, stdout io.Writer) error {
	pl, err := Prepare(p, cfg, sel, args)
	if err != nil {
		return err
	}
	trace.Step("prepare")
	if err := checkLinks(p, pl); err != nil {
		return err
	}
	trace.Step("check links")
	if err := prelaunch(p, cfg, pl, dry); err != nil {
		return err
	}
	trace.Step("prelaunch")
	if dry {
		pl.DryRun(stdout)
		return nil
	}
	return pl.Exec(p)
}

// SessionDir holds markers of running sessions of a profile.
func SessionDir(p paths.Paths, name string) string {
	return filepath.Join(p.Runtime, "profile-"+name)
}

// Running reports whether a Claude Code session of the profile, started by
// ccam, is running.
func Running(p paths.Paths, name string) bool { return session.Active(SessionDir(p, name)) }

// prelaunch registers the session (the marker then names the exec'd claude
// process) and applies MCP sharing.
func prelaunch(p paths.Paths, cfg *config.Config, pl *Plan, dry bool) error {
	running := Running(p, pl.Profile.Name)
	if err := applyMCP(p, cfg, pl, dry, running); err != nil {
		return err
	}
	if dry {
		return nil
	}
	if err := historyCheck(p, pl); err != nil {
		return err
	}
	return session.Register(SessionDir(p, pl.Profile.Name))
}

// HistoryDir guards one (history group, working dir) pair: two accounts
// sharing conversations in the same folder would write the same session files.
func HistoryDir(p paths.Paths, group, cwd string) string {
	h := fnv.New64a()
	h.Write([]byte(cwd))
	return filepath.Join(p.Runtime, fmt.Sprintf("history-%s-%016x", group, h.Sum64()))
}

func historyCheck(p paths.Paths, pl *Plan) error {
	group := pl.Profile.Share.Get("history")
	if group == config.Local {
		return nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	dir := HistoryDir(p, group, cwd)
	if session.Active(dir) {
		fmt.Fprintf(os.Stderr, "ccam: another Claude Code sharing conversation history (group %s) is running in this folder.\n      Continuing the same conversation from both at once can mix up the history.\n", group)
		if !askTTY("Run anyway?") {
			return fmt.Errorf("cancelled")
		}
	}
	return session.Register(dir)
}

// askTTY asks on the terminal; without one it proceeds (the warning is shown).
func askTTY(question string) bool {
	fi, err := os.Stdin.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		return true
	}
	fmt.Fprintf(os.Stderr, "%s [y/N] ", question)
	var buf [16]byte
	n, _ := os.Stdin.Read(buf[:])
	ans := strings.ToLower(strings.TrimSpace(string(buf[:n])))
	return ans == "y" || ans == "yes"
}

func applyMCP(p paths.Paths, cfg *config.Config, pl *Plan, dry, running bool) error {
	group := pl.Profile.Share.Get("mcp")
	if group == config.Local {
		return nil
	}
	gdir := p.GroupDir(group)
	state := filepath.Join(p.State, "mcp", pl.Profile.Name)
	if cfg.Strategy() == config.MCPSync {
		switch {
		case dry:
		case running:
			fmt.Fprintf(os.Stderr, "ccam: account %s is already running, skipping MCP sync\n", pl.Profile.Name)
		default:
			if err := mcp.Sync(pl.Dir, gdir, state); err != nil {
				fmt.Fprintf(os.Stderr, "ccam: MCP sync failed (launching anyway): %v\n", err)
			}
		}
		return nil
	}
	pl.InsertArgs(mcp.FlagArgs(gdir)...)
	if dry {
		return nil
	}
	fresh, err := mcp.NewLocal(pl.Dir, mcp.Path(gdir), filepath.Join(state, "seen"))
	if err == nil && len(fresh) > 0 {
		fmt.Fprintf(os.Stderr, "ccam: MCP servers only account %s has: %s\n      To share them with other accounts: ccam mcp import %s %s\n",
			pl.Profile.Name, strings.Join(fresh, ", "), pl.Profile.Name, strings.Join(fresh, " "))
	}
	return nil
}
