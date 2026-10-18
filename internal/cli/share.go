package cli

import (
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/profile"
	"github.com/jhyunwoo/ccam/internal/share"
)

// applyShare links a freshly added (empty) profile into its share groups.
func applyShare(e *env, cfg *config.Config, prof *config.Profile) error {
	pl, err := share.Build(profile.Dir(e.paths, prof), e.paths.GroupsRoot(), config.Share{}, prof.Share, share.Options{})
	if err != nil {
		return err
	}
	return pl.Apply()
}

var itemLabels = map[string]string{
	"mcp":        "MCP 서버",
	"settings":   "settings.json",
	"memory":     "전역 CLAUDE.md",
	"extensions": "에이전트·명령·스킬·플러그인",
	"history":    "대화 내역",
}

func shareCmd(e *env) *cobra.Command {
	var fresh, yes bool
	var prefer string
	values := map[string]*string{}
	cmd := &cobra.Command{
		Use:   "share <이름|번호> [--mcp on|off|<그룹>] [--history …]",
		Short: "공유 항목 변경 (데이터 이전 포함)",
		Long: "항목마다 on(=default 그룹), off(=local), 또는 그룹 이름을 줍니다.\n" +
			"항목: --mcp --settings --memory --extensions --history\n" +
			"옵션 없이 실행하면 현재 상태만 보여줍니다.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			prof, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			target := prof.Share
			changed := false
			for _, item := range config.Items {
				if cmd.Flags().Changed(item) {
					if err := target.Set(item, *values[item]); err != nil {
						return err
					}
					changed = true
				}
			}
			if !changed {
				return showShare(e, prof)
			}
			res, err := parsePrefer(prefer)
			if err != nil {
				return err
			}
			return changeShare(e, cfg, prof, target, share.Options{Fresh: fresh}, yes, res)
		},
	}
	for _, item := range config.Items {
		values[item] = new(string)
		cmd.Flags().StringVar(values[item], item, "", itemLabels[item]+": on|off|<그룹>")
	}
	cmd.Flags().BoolVar(&fresh, "fresh", false, "공유를 끌 때 그룹 내용을 복사하지 않고 빈 상태로 시작")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "확인 없이 실행 (충돌은 --prefer 로 처리)")
	cmd.Flags().StringVar(&prefer, "prefer", "", "충돌 처리: group|mine|merge (기본: 물어봄, --yes 면 항목별 기본값)")
	return cmd
}

func parsePrefer(s string) (*share.Resolution, error) {
	var r share.Resolution
	switch s {
	case "":
		return nil, nil
	case "group":
		r = share.KeepGroup
	case "mine":
		r = share.UseMine
	case "merge":
		r = share.MergeJSON
	default:
		return nil, fmt.Errorf("--prefer 는 group, mine, merge 중 하나입니다")
	}
	return &r, nil
}

func showShare(e *env, prof *config.Profile) error {
	tw := tabwriter.NewWriter(e.stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "%s 의 공유 상태\n", prof.Name)
	for _, item := range config.Items {
		g := prof.Share.Get(item)
		if g == config.Local {
			g = "분리"
		} else {
			g = "공유 (" + g + " 그룹)"
		}
		fmt.Fprintf(tw, "  %s\t%s\t--%s\n", itemLabels[item], g, item)
	}
	return tw.Flush()
}

// changeShare prints the plan, asks, resolves conflicts and applies.
func changeShare(e *env, cfg *config.Config, prof *config.Profile, target config.Share, opt share.Options, yes bool, prefer *share.Resolution) error {
	dir := profile.Dir(e.paths, prof)
	pl, err := share.Build(dir, e.paths.GroupsRoot(), prof.Share, target, opt)
	if err != nil {
		return err
	}
	for _, c := range pl.Conflicts {
		if prefer != nil {
			for _, cand := range c.Candidates {
				if cand == *prefer || (*prefer == share.MergeJSON && cand == share.MergeLines) {
					c.Choice = cand
				}
			}
			continue
		}
		if yes {
			continue
		}
		c.Choice = askConflict(e, c)
	}
	e.printf("계획 (%s):\n", prof.Name)
	pl.Describe(e.stdout)
	if e.dryRun {
		return nil
	}
	if !pl.Empty() && !yes {
		if !e.confirm("진행할까요?") {
			e.printf("취소했습니다.\n")
			return nil
		}
	}
	if err := pl.Apply(); err != nil {
		return fmt.Errorf("적용 중 실패 (원본은 %s/%s 에 있습니다): %w", dir, share.BackupDir, err)
	}
	prof.Share = target
	if err := e.save(cfg); err != nil {
		return err
	}
	e.printf("✓ %s 공유 설정을 바꿨습니다\n", prof.Name)
	return nil
}

func askConflict(e *env, c *share.Conflict) share.Resolution {
	e.printf("\n%s 가 그룹 것과 다릅니다.\n  내 것:   %s\n  그룹 것: %s\n", c.Entry.Name, c.Mine, c.Group)
	for i, r := range c.Candidates {
		e.printf("  %d) %s\n", i+1, r)
	}
	for {
		e.printf("선택 [1]: ")
		line, err := e.in.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" || err != nil {
			return c.Candidates[0]
		}
		if n, err := strconv.Atoi(line); err == nil && n >= 1 && n <= len(c.Candidates) {
			return c.Candidates[n-1]
		}
	}
}
