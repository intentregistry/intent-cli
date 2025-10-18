package cmd

import (
	"github.com/spf13/cobra"
	"github.com/intentregistry/intent-cli/internal/version"
)

func RootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "intent",
		Short: "IntentRegistry CLI",
		Long:  "Publish & install AI Intents from intentregistry.com",
	}
	c.Version = version.GetVersion()
	c.SetVersionTemplate("intent {{.Version}}\n")
	// Flags globales (si quieres)
	c.PersistentFlags().String("api-url", "", "Override API base URL (default from config)")
	c.PersistentFlags().String("token", "", "API token (default from config)")
	return c
}