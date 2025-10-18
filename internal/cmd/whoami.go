package cmd

import (
	"fmt"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/intentregistry/intent-cli/internal/httpclient"
	"github.com/spf13/cobra"
)

func WhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show current authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()
			cl := httpclient.New(cfg)
			var resp struct {
				Email string `json:"email"`
				UserID string `json:"userId"`
			}
			// Ajusta el endpoint real
			if err := cl.Get("/me", &resp); err != nil {
				return err
			}
			fmt.Printf("👤 %s (%s)\n", resp.Email, resp.UserID)
			return nil
		},
	}
}