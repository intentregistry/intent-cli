package cmd

import (
	"github.com/spf13/cobra"
	"github.com/intentregistry/intent-cli/internal/version"
)

var (
	root  = &cobra.Command{
		Use:     "intent",
		Short:   "IntentRegistry CLI",
		Long:    "Publish & install AI Intents from intentregistry.com",
		Version: version.Short(),
	}
	debug bool
)

func init() {
	root.SetVersionTemplate("intent {{.Version}}\n")
	root.PersistentFlags().BoolVar(&debug, "debug", false, "Enable verbose debug output")

	// attach all subcommands here (example)
	root.AddCommand(
		LoginCmd(),
		PublishCmd(),
		InstallCmd(),
		WhoamiCmd(),
		SearchCmd(),
		VersionCmd(),
		CompletionCmd(),
	)
}

func RootCmd() *cobra.Command { return root }
func Execute() error          { return root.Execute() }
func Debug() bool             { return debug }