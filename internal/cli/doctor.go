package cli

import (
	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/doctor"
)

func doctorCmd(e *env) *cobra.Command {
	var fix bool
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "링크·권한·로그인·Claude Code 버전 점검 (--fix 로 자동 수리)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			self, _ := ccamPath()
			findings := doctor.Run(e.paths, cfg, doctor.Options{Fix: fix, CcamPath: self, Usage: newFetcher(usageStateDir(e))})
			for _, f := range findings {
				e.printf("%s\n", f)
			}
			if doctor.Failed(findings) {
				if !fix {
					e.printf("\n`ccam doctor --fix` 로 고칠 수 있는 것은 고칩니다.\n")
				}
				return exitError(1)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&fix, "fix", false, "안전하게 고칠 수 있는 문제를 고침 (바뀌는 파일은 백업)")
	return cmd
}
