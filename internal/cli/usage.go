package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/usage"
)

func usageCmd(e *env) *cobra.Command {
	var opt usage.Options
	var asJSON, verbose bool
	var watch int
	cmd := &cobra.Command{
		Use:     "usage",
		Aliases: []string{"u"},
		Short:   "모든 계정의 5시간·주간 사용률과 초기화 시각",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opt.Source != "" && opt.Source != "live" && opt.Source != "cache" {
				return fmt.Errorf("--source 는 live 또는 cache 입니다")
			}
			cfg, err := e.load()
			if err != nil {
				return err
			}
			if len(cfg.Profiles) == 0 {
				e.printf("등록된 계정이 없습니다. `ccam add <이름>` 으로 추가하세요.\n")
				return nil
			}
			f := newFetcher(usageStateDir(e))
			if !cmd.Flags().Changed("watch") {
				return showUsage(e, f, cfg, opt, asJSON, verbose, false)
			}
			if watch <= 0 {
				watch = 60
			}
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()
			for {
				if err := showUsage(e, f, cfg, opt, asJSON, verbose, true); err != nil {
					return err
				}
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(time.Duration(watch) * time.Second):
				}
			}
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "JSON 으로 출력")
	cmd.Flags().IntVar(&watch, "watch", 60, "N초마다 새로 고침")
	cmd.Flags().Lookup("watch").NoOptDefVal = "60"
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "모델별 한도와 출처 표시")
	cmd.Flags().StringVar(&opt.Source, "source", "", "live|cache")
	cmd.Flags().BoolVar(&opt.NoCache, "no-cache", false, "60초 캐시 무시")
	cmd.Flags().BoolVar(&opt.Refresh, "refresh", false, "만료된 토큰 갱신 (실행 중이 아닌 계정만)")
	return cmd
}

// newFetcher is replaced in tests to point at a fake server.
var newFetcher = usage.NewFetcher

func usageStateDir(e *env) string { return filepath.Join(e.paths.State, "usage") }

func usageAccounts(e *env, cfg *config.Config) []usage.Account {
	out := make([]usage.Account, 0, len(cfg.Profiles))
	for i := range cfg.Profiles {
		p := &cfg.Profiles[i]
		out = append(out, usage.Account{
			Name: p.Name, Slot: p.Slot, Dir: profile.Dir(e.paths, p),
			LockPath: launcher.ProfileLockPath(e.paths, p.Name),
			Imported: p.Dir != "",
		})
	}
	return out
}

func showUsage(e *env, f *usage.Fetcher, cfg *config.Config, opt usage.Options, asJSON, verbose, clear bool) error {
	accts := usageAccounts(e, cfg)
	if asJSON {
		reports := f.FetchAll(context.Background(), accts, opt, nil)
		enc := json.NewEncoder(e.stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(reports)
	}
	tty := isTerminal(e.stdout)
	ro := usage.RenderOptions{Color: tty && os.Getenv("NO_COLOR") == "", Verbose: verbose}
	view := make([]*usage.Report, len(accts))
	for i, a := range accts {
		view[i] = &usage.Report{Name: a.Name, Slot: a.Slot, Status: usage.StatusPending}
	}
	if clear && tty {
		fmt.Fprint(e.stdout, "\x1b[H\x1b[2J")
	}
	// On a terminal, draw placeholders and redraw as each account answers.
	drawn := 0
	draw := func() {
		ro.Now = time.Now()
		lines := usage.Table(view, ro)
		if drawn > 0 {
			fmt.Fprintf(e.stdout, "\x1b[%dA", drawn)
		}
		for _, l := range lines {
			if tty {
				fmt.Fprint(e.stdout, "\x1b[2K")
			}
			fmt.Fprintln(e.stdout, l)
		}
		drawn = len(lines)
	}
	if tty {
		draw()
	}
	f.FetchAll(context.Background(), accts, opt, func(i int, r usage.Report) {
		view[i] = &r
		if tty {
			draw()
		}
	})
	if !tty {
		draw()
	}
	if best := usage.Best(view); best != nil && len(view) > 1 {
		e.printf("→ 여유가 가장 많은 계정: %s   (ccs next --auto 로 바로 열기)\n", best.Name)
	}
	return nil
}

func isTerminal(w any) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

func statuslineCmd(e *env) *cobra.Command {
	var wrap string
	cmd := &cobra.Command{
		Use:   "statusline",
		Short: "Claude Code 상태줄: 계정 이름 표시와 사용량 수집",
		Long:  "Claude Code가 상태줄 명령으로 실행합니다. 설치는 `ccam statusline install`.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusline(e, wrap)
		},
	}
	cmd.Flags().StringVar(&wrap, "wrap", "", "기존 상태줄 명령 (출력 앞에 계정 이름을 붙임)")
	cmd.AddCommand(statuslineInstallCmd(e))
	return cmd
}

func runStatusline(e *env, wrap string) error {
	input, err := readAll(e.in)
	if err != nil {
		return err
	}
	name := os.Getenv("CCAM_PROFILE")
	if name == "" {
		name = profileByConfigDir(e, os.Getenv("CLAUDE_CONFIG_DIR"))
	}
	limits, model, ok, _ := usage.ParseStatusLine(input)
	if ok && name != "" {
		_ = usage.SaveStatusLine(usageStateDir(e), name, limits, time.Now())
	}
	label := name
	if label == "" {
		label = "claude"
	}
	if wrap != "" {
		out, _ := runShell(wrap, input)
		first, _, _ := strings.Cut(strings.TrimRight(out, "\n"), "\n")
		e.printf("[%s] %s\n", label, first)
		return nil
	}
	parts := []string{"[" + label + "]"}
	if model != "" {
		parts = append(parts, model)
	}
	if limits.FiveHour.Known {
		parts = append(parts, fmt.Sprintf("5h %.0f%%", limits.FiveHour.Percent))
	}
	if limits.Weekly.Known {
		parts = append(parts, fmt.Sprintf("7d %.0f%%", limits.Weekly.Percent))
	}
	e.printf("%s\n", strings.Join(parts, " · "))
	return nil
}

func profileByConfigDir(e *env, dir string) string {
	if dir == "" {
		return ""
	}
	cfg, err := e.load()
	if err != nil {
		return ""
	}
	for i := range cfg.Profiles {
		if profile.Dir(e.paths, &cfg.Profiles[i]) == dir {
			return cfg.Profiles[i].Name
		}
	}
	return ""
}
