package cmd

import (
	"github.com/spf13/cobra"
	"github.com/intentregistry/intent-cli/internal/version"
)

func RootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "intent",
		Short:   "IntentRegistry CLI",
		Long:    "Publish & install AI Intents from intentregistry.com",
		Version: version.Short(), // prints only "0.x.y" (or "dev" locally)
	}
	c.SetVersionTemplate("intent {{.Version}}\n")
	return c
}