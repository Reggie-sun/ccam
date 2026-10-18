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
		fmt.Fprintf(os.Stderr, "ccam: 이 폴더에서 대화 내역을 공유하는 다른 Claude Code(%s 그룹)가 실행 중입니다.\n      같은 대화를 동시에 이어 쓰면 기록이 섞일 수 있습니다.\n", group)
		if !askTTY("그래도 실행할까요?") {
			return fmt.Errorf("취소했습니다")
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
