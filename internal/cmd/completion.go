package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func CompletionCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long:  `Generate shell completion scripts for intent (bash, zsh, fish, powershell).`,
		Args:  cobra.MaximumNArgs(1), // allow zero args
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := ""
			if len(args) == 1 {
				shell = strings.ToLower(args[0])
			} else {
				// try to detect from $SHELL
				if sh := os.Getenv("SHELL"); sh != "" {
					shell = strings.ToLower(filepath.Base(sh))
				}
			}

			switch shell {
			case "bash":
				return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "pwsh", "powershell":
				return cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
			default:
				// fallback: show helpful message
				return fmt.Errorf("unknown or missing shell. Use one of: bash|zsh|fish|powershell")
			}
		},
	}

	// Keep helpful usage hints
	c.ValidArgs = []string{"bash", "zsh", "fish", "powershell"}
	c.DisableFlagsInUseLine = true
	return c
}