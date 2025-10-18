package cmd

import (
    "github.com/spf13/cobra"
    "github.com/intentregistry/intent-cli/internal/version"
)

var root = &cobra.Command{
    Use:     "intent",
    Short:   "IntentRegistry CLI",
    Long:    "Publish & install AI Intents from intentregistry.com",
    Version: version.Short(),
}

func init() {
    root.SetVersionTemplate("intent {{.Version}}\n")

    // attach all subcommands here once at init()
    root.AddCommand(
        LoginCmd(),
        PublishCmd(),
        InstallCmd(),
        WhoamiCmd(),
        SearchCmd(),
        VersionCmd(),
    )
}

func RootCmd() *cobra.Command { return root }

func Execute() error { return root.Execute() }