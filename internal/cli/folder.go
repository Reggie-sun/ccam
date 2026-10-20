package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/launcher"
)

func folderCmd(e *env) *cobra.Command {
	var remove bool
	cmd := &cobra.Command{
		Use:   "folder [name|number] [folder]",
		Short: "Use an account in a folder when no account is named",
		Long: "Sets which account `ccam` uses in a folder and every folder below it when you\n" +
			"start Claude Code without an account name (ccam -p \"…\", ccam run, ccam\n" +
			"--dry-run). A name, a number or a shortcut still starts that account. Where\n" +
			"no folder is set, the default account (ccam use) is used.\n\n" +
			"  ccam folder                     list the folders and the account used here\n" +
			"  ccam folder <account> [folder]  use the account in folder (default: here)\n" +
			"  ccam folder --remove [folder]   stop using an account there",
		Example: "  ccam folder work ~/code/company\n  ccam folder personal\n  ccam folder --remove ~/code/company",
		Args:    cobra.MaximumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if remove || len(args) == 1 {
				return nil, cobra.ShellCompDirectiveFilterDirs
			}
			return profileNames(e)(cmd, args, toComplete)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := e.load()
			if err != nil {
				return err
			}
			switch {
			case remove:
				if len(args) > 1 {
					return fmt.Errorf("--remove takes one folder at most")
				}
				return removeFolder(e, cfg, args)
			case len(args) == 0:
				return listFolders(e, cfg)
			}
			prof, err := cfg.MustFind(args[0])
			if err != nil {
				return err
			}
			dir, err := folderArg(args[1:])
			if err != nil {
				return err
			}
			if info, err := os.Stat(dir); err != nil || !info.IsDir() {
				return fmt.Errorf("%s is not a folder", dir)
			}
			if e.dryRun {
				e.printf("Would use %s in %s.\n", prof.Name, dir)
				return nil
			}
			cfg.SetFolder(dir, prof.Name)
			if err := e.save(cfg); err != nil {
				return err
			}
			e.printf("✓ %s → %s\n  ccam without an account name starts %s here and in every folder below.\n", dir, prof.Name, prof.Name)
			return nil
		},
	}
	cmd.Flags().BoolVar(&remove, "remove", false, "stop using an account in the folder")
	return cmd
}

// folderArg is the folder args[0] names, or the working directory, as an
// absolute path with symlinks resolved (the form entries are stored in).
func folderArg(args []string) (string, error) {
	dir, err := os.Getwd()
	if len(args) > 0 {
		dir, err = filepath.Abs(args[0])
	}
	if err != nil {
		return "", err
	}
	if real, err := filepath.EvalSymlinks(dir); err == nil {
		return real, nil
	}
	return dir, nil
}

func removeFolder(e *env, cfg *config.Config, args []string) error {
	dir, err := folderArg(args)
	if err != nil {
		return err
	}
	f, ok := cfg.FolderFor(dir)
	if !ok {
		return fmt.Errorf("no account is set for %s", dir)
	}
	if !config.SamePath(f.Path, dir) {
		return fmt.Errorf("no account is set for %s itself; it gets %s from %s (`ccam folder --remove %s`)", dir, f.Account, f.Path, f.Path)
	}
	if e.dryRun {
		e.printf("Would stop using %s in %s.\n", f.Account, f.Path)
		return nil
	}
	account, path := f.Account, f.Path
	cfg.UnsetFolder(path)
	if err := e.save(cfg); err != nil {
		return err
	}
	e.printf("✓ %s no longer uses %s\n", path, account)
	return nil
}

func listFolders(e *env, cfg *config.Config) error {
	if len(cfg.Folders) == 0 {
		e.printf("No folders set. Set one with `ccam folder <account> [folder]`.\n")
	}
	for _, f := range cfg.Folders {
		e.printf("%s → %s\n", f.Path, f.Account)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	if f, ok := launcher.FolderHere(cfg); ok {
		e.printf("\nHere (%s): %s, set for %s\n", cwd, f.Account, f.Path)
	} else if def, err := cfg.DefaultProfile(); err == nil {
		e.printf("\nHere (%s): %s, the default account\n", cwd, def.Name)
	}
	return nil
}
