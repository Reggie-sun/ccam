package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/profile"
)

func importCmd(e *env) *cobra.Command {
	var name string
	var yes bool
	cmd := &cobra.Command{
		Use:   "import [폴더]",
		Short: "기존 ~/.claude 나 ~/.claude-N 폴더를 그대로 계정으로 등록",
		Long: "폴더를 옮기거나 복사하지 않고 제자리에 등록합니다 (토큰을 복사하면 서로 로그아웃시킵니다).\n" +
			"폴더를 주지 않으면 ~/.claude 와 ~/.claude-* 중 등록되지 않은 것을 찾아 하나씩 묻습니다.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			var paths []string
			names := map[string]string{}
			if len(args) == 1 {
				paths = []string{args[0]}
				names[args[0]] = name
			} else {
				cands := profile.FindCandidates(e.paths, cfg)
				if len(cands) == 0 {
					e.printf("가져올 폴더가 없습니다 (~/.claude, ~/.claude-*).\n")
					return nil
				}
				for _, c := range cands {
					label := c.Name
					if label == "" {
						label = "(이름 필요: ccam import " + c.Path + " --name <이름>)"
					}
					q := fmt.Sprintf("%s (%s) → %s 로 가져올까요?", c.Path, orNone(c.Email), label)
					if c.Name != "" && (yes || e.confirm(q)) {
						paths = append(paths, c.Path)
					} else if c.Name == "" {
						e.printf("%s %s\n", c.Path, label)
					}
				}
			}
			for _, path := range paths {
				prof, notes, err := profile.Import(e.paths, cfg, path, names[path])
				if err != nil {
					return err
				}
				for _, n := range notes {
					e.printf("· %s\n", n)
				}
				e.printf("✓ %s (%s) 등록 — 단축 명령: claude%d\n", prof.Name, orNone(prof.Email), prof.Slot)
			}
			if len(paths) == 0 {
				return nil
			}
			if err := e.save(cfg); err != nil {
				return err
			}
			if _, err := installShims(e, cfg); err != nil {
				fmt.Fprintln(e.stderr, "단축 명령을 만들지 못했습니다:", err)
			}
			e.printf("가져온 계정은 공유 없이 시작합니다. 공유하려면 `ccam share <이름> --settings on …`\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "계정 이름 (기본: 폴더 이름에서 만듦)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "찾은 폴더를 묻지 않고 모두 가져오기")
	return cmd
}

func orNone(s string) string {
	if s == "" {
		return "로그인 정보 없음"
	}
	return s
}
