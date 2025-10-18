package cmd

import (
	"fmt"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/intentregistry/intent-cli/internal/httpclient"
	"github.com/spf13/cobra"
)

func InstallCmd() *cobra.Command {
	var dest string
	c := &cobra.Command{
		Use:   "install <@scope/name[@version]>",
		Short: "Install an intent package to local project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			spec := args[0]
			cfg := config.Load()
			cl := httpclient.New(cfg)
			// GET metadata → download artifact → validate → extract into dest
			fmt.Println("⬇️  Installing", spec, "into", dest)
			_ = cl // implementar cuando tengas endpoints
			fmt.Println("✅ Installed")
			return nil
		},
	}
	c.Flags().StringVar(&dest, "dest", "intents", "destination folder")
	return c
}