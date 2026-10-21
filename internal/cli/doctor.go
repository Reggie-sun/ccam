package cli

import (
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/doctor"
	"github.com/jhyunwoo/ccam/internal/tui"
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
			fix := fix && !e.dryRun // a dry run reports what --fix would repair
			self, _ := ccamPath()
			findings := doctor.Run(e.paths, cfg, doctor.Options{Fix: fix, CcamPath: self, Usage: e.fetcher()})
			// On a terminal the boxes fit its width; otherwise each box is as
			// wide as its longest line, so no message is cut or wrapped.
			width, tty := 0, isTerminal(e.stdout)
			if f, ok := e.stdout.(*os.File); ok && tty {
				if w, _, err := term.GetSize(f.Fd()); err == nil && w > 1 {
					width = w - 1
				}
			}
			for _, l := range tui.HealthReport(findings, width, tty && os.Getenv("NO_COLOR") == "") {
				e.printf("%s\n", l)
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
