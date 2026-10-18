package launcher

import (
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"runtime"
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
	err = pl.Exec(p)
	runtime.KeepAlive(held) // a finalizer must not close the lock fds before exec
	return err
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
	if !dry {
		l, err := historyLock(p, pl)
		if err != nil {
			return nil, err
		}
		if l != nil {
			held = append(held, l)
		}
	}
	return held, nil
}

// HistoryLockPath guards one (history group, working dir) pair: two
// accounts sharing conversations in the same folder would write the same
// session files.
func HistoryLockPath(p paths.Paths, group, cwd string) string {
	h := fnv.New64a()
	h.Write([]byte(cwd))
	return filepath.Join(p.Runtime, fmt.Sprintf("history-%s-%016x.lock", group, h.Sum64()))
}

func historyLock(p paths.Paths, pl *Plan) (*lock.Flock, error) {
	group := pl.Profile.Share.Get("history")
	if group == config.Local {
		return nil, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil
	}
	l, err := lock.TryFlock(HistoryLockPath(p, group, cwd))
	if err == nil {
		return l, nil
	}
	if err != lock.ErrBusy {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "ccam: 이 폴더에서 대화 내역을 공유하는 다른 Claude Code(%s 그룹)가 실행 중입니다.\n      같은 대화를 동시에 이어 쓰면 기록이 섞일 수 있습니다.\n", group)
	if !askTTY("그래도 실행할까요?") {
		return nil, fmt.Errorf("취소했습니다")
	}
	return nil, nil
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
