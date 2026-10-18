package launcher

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/lock"
	"github.com/jhyunwoo/ccam/internal/mcp"
	"github.com/jhyunwoo/ccam/internal/paths"
	"github.com/jhyunwoo/ccam/internal/trace"
)

// Launch is the whole `ccs <name>` sequence: load config, prepare, check
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
	held, err := prelaunch(p, cfg, pl, dry)
	if err != nil {
		return err
	}
	trace.Step("prelaunch")
	if dry {
		pl.DryRun(stdout)
		return nil
	}
	for _, l := range held {
		if err := l.KeepAcrossExec(); err != nil {
			return err
		}
	}
	return pl.Exec(p)
}

// ProfileLockPath is held by a running Claude Code session of that profile.
func ProfileLockPath(p paths.Paths, name string) string {
	return filepath.Join(p.Runtime, "profile-"+name+".lock")
}

// Running reports whether a Claude Code session of the profile is running.
func Running(p paths.Paths, name string) bool { return lock.Held(ProfileLockPath(p, name)) }

// prelaunch takes the locks that the claude process will hold and applies
// MCP sharing. It returns the locks to hand over across exec.
func prelaunch(p paths.Paths, cfg *config.Config, pl *Plan, dry bool) ([]*lock.Flock, error) {
	var held []*lock.Flock
	running := false
	if !dry {
		l, err := lock.TryFlock(ProfileLockPath(p, pl.Profile.Name))
		switch {
		case err == lock.ErrBusy:
			running = true
		case err != nil:
			return nil, err
		default:
			held = append(held, l)
		}
	}
	if err := applyMCP(p, cfg, pl, dry, running); err != nil {
		return nil, err
	}
	return held, nil
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
			fmt.Fprintf(os.Stderr, "ccam: %s 계정이 이미 실행 중이라 MCP 동기화를 건너뜁니다\n", pl.Profile.Name)
		default:
			if err := mcp.Sync(pl.Dir, gdir, state); err != nil {
				fmt.Fprintf(os.Stderr, "ccam: MCP 동기화 실패 (그대로 실행합니다): %v\n", err)
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
		fmt.Fprintf(os.Stderr, "ccam: %s 계정에만 있는 MCP 서버: %s\n      다른 계정과 공유하려면: ccam mcp import %s %s\n",
			pl.Profile.Name, strings.Join(fresh, ", "), pl.Profile.Name, strings.Join(fresh, " "))
	}
	return nil
}
