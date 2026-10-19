package cli

import (
	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/doctor"
)

func doctorCmd(e *env) *cobra.Command {
	var fix bool
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check links, permissions, login and Claude Code version (--fix repairs)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			self, _ := ccamPath()
			findings := doctor.Run(e.paths, cfg, doctor.Options{Fix: fix, CcamPath: self, Usage: e.fetcher()})
			for _, f := range findings {
				e.printf("%s\n", f)
			}
			if doctor.Failed(findings) {
				if !fix {
					e.printf("\n`ccam doctor --fix` repairs what it can.\n")
				}
				return exitError(1)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&fix, "fix", false, "repair what is safe to repair (changed files are backed up)")
	return cmd
}
