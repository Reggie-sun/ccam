package cli

import (
	"errors"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/update"
)

func updateCmd(e *env) *cobra.Command {
	var o update.Options
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update ccam to the latest release",
		Long: "Downloads the latest ccam release from GitHub, checks it against the release's\n" +
			"checksums.txt (SHA-256) and replaces this ccam binary in one rename. Shortcut\n" +
			"commands link to ccam, so they use the new version too.\n\n" +
			"ccam never updates or checks for updates on its own; only this command does.\n" +
			"A ccam installed from the .deb package is left to apt: the command prints how\n" +
			"to install the new .deb instead.",
		Example: "  ccam update                    # install the latest release\n" +
			"  ccam update --check            # only say whether a newer release exists\n" +
			"  ccam update --version v0.2.0   # install that release (also an older one)\n" +
			"  ccam --dry-run update          # show what would be downloaded and replaced",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.Check && o.Version != "" {
				return errors.New("--check and --version cannot be used together")
			}
			exe, err := ccamPath()
			if err != nil {
				return err
			}
			o.DryRun = e.dryRun
			u := update.New(Version, exe, runtime.GOOS, runtime.GOARCH)
			return update.Run(cmd.Context(), u, o, e.stdout)
		},
	}
	cmd.Flags().BoolVar(&o.Check, "check", false, "only report whether a newer release exists")
	cmd.Flags().StringVar(&o.Version, "version", "", "install this release tag, such as v0.3.0")
	return cmd
}
