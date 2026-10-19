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
		Use:   "import [folder]",
		Short: "Register an existing ~/.claude or ~/.claude-N folder as an account, in place",
		Long: "The folder is registered where it is, never moved or copied (copied tokens log each other out).\n" +
			"Without a folder, finds unregistered ~/.claude and ~/.claude-* folders and asks about each.",
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
					e.printf("No folders to import (~/.claude, ~/.claude-*).\n")
					return nil
				}
				for _, c := range cands {
					label := c.Name
					if label == "" {
						label = "(needs a name: ccam import " + c.Path + " --name <name>)"
					}
					q := fmt.Sprintf("Import %s (%s) as %s?", c.Path, orNone(c.Email), label)
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
				e.printf("✓ %s (%s) registered\n", prof.Name, orNone(prof.Email))
				printShortcut(e, prof)
			}
			if len(paths) == 0 {
				return nil
			}
			if err := e.save(cfg); err != nil {
				return err
			}
			if _, err := installShims(e, cfg); err != nil {
				fmt.Fprintln(e.stderr, "Could not create shortcuts:", err)
			}
			e.printf("Imported accounts start with nothing shared. To share: `ccam share <name> --settings on …`\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "account name (default: derived from the folder name)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "import every folder found without asking")
	return cmd
}

func orNone(s string) string {
	if s == "" {
		return "not logged in"
	}
	return s
}
