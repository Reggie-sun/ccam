package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func completionCmd(e *env) *cobra.Command {
	return &cobra.Command{
		Use:       "completion bash|zsh",
		Short:     "Print a shell completion script",
		Example:   "  ccam completion bash > ~/.local/share/bash-completion/completions/ccam\n  ccam completion zsh > \"${fpath[1]}/_ccam\"",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh"},
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch args[0] {
			case "bash":
				return root.GenBashCompletionV2(e.stdout, true)
			case "zsh":
				return root.GenZshCompletion(e.stdout)
			default:
				return fmt.Errorf("only bash and zsh are supported")
			}
		},
	}
}

// profileNames completes the first argument with profile names.
func profileNames(e *env) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		cfg, err := e.load()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var out []string
		for _, p := range cfg.Profiles {
			out = append(out, p.Name+"\t"+p.Email)
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	}
}
