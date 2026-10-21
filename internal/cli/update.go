package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/shim"
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
			"The dashboard (ccam on a terminal) also asks for the latest release each time\n" +
			"it opens and offers this update; ccam never updates without your say-so.\n" +
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
			o.DryRun = e.dryRun
			return e.runUpdate(cmd.Context(), o, e.stdout)
		},
	}
	cmd.Flags().BoolVar(&o.Check, "check", false, "only report whether a newer release exists")
	cmd.Flags().StringVar(&o.Version, "version", "", "install this release tag, such as v0.3.0")
	return cmd
}

// runUpdate is `ccam update`, also run by the dashboard's U.
func (e *env) runUpdate(ctx context.Context, o update.Options, out io.Writer) error {
	exe, err := ccamPath()
	if err != nil {
		return err
	}
	u := update.New(Version, exe, runtime.GOOS, runtime.GOARCH)
	if err := update.Run(ctx, u, o, out); err != nil {
		return err
	}
	// Hard-linked shortcuts (Windows) still run the replaced binary.
	if cfg, err := e.load(); err == nil && runtime.GOOS == "windows" && !o.DryRun {
		for _, name := range shim.Relink(e.paths.Bin, exe, update.OldExe(exe), cfg) {
			fmt.Fprintf(out, "✓ shortcut %s renewed\n", name)
		}
	}
	return nil
}

// newerRelease is the latest release tag when it is newer than this ccam,
// else "" (also offline, or for a development build).
func newerRelease(ctx context.Context) string {
	exe, err := ccamPath()
	if err != nil {
		return ""
	}
	tag, _ := update.New(Version, exe, runtime.GOOS, runtime.GOARCH).Newer(ctx)
	return tag
}
