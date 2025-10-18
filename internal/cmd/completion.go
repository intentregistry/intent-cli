package cmd

import (
	"github.com/spf13/cobra"
)

func CompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
$ source <(intent completion bash)

# To load completions for each session, execute once:
Linux:
  $ intent completion bash > /etc/bash_completion.d/intent
MacOS:
  $ intent completion bash > /usr/local/etc/bash_completion.d/intent

Zsh:
# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:
$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ intent completion zsh > "${fpath[1]}/_intent"

# You will need to start a new shell for this setup to take effect.

Fish:
$ intent completion fish | source

# To load completions for each session, execute once:
$ intent completion fish > ~/.config/fish/completions/intent.fish

PowerShell:
PS> intent completion powershell | Out-String | Invoke-Expression

# To load completions for each session, execute once:
PS> intent completion powershell > intent.ps1
# and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
			}
		},
	}
}
