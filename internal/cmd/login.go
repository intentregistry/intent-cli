package cmd

import (
	"fmt"
	"os"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/spf13/cobra"
)

func LoginCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "login",
		Short: "Authenticate against IntentRegistry",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get current config to preserve existing api_url
			cfg := config.Load()
			
			// Prompt for API URL if not already set
			if cfg.APIURL == "" || cfg.APIURL == "https://api.intentregistry.com" {
				fmt.Print("Enter API URL (default: https://api.intentregistry.com): ")
				var apiURL string
				_, err := fmt.Scanln(&apiURL)
				if err != nil && err.Error() != "unexpected newline" {
					return err
				}
				if apiURL != "" {
					cfg.APIURL = apiURL
				} else {
					cfg.APIURL = "https://api.intentregistry.com"
				}
			}
			
			// Prompt for token
			var token string
			fmt.Print("Enter API token: ")
			_, err := fmt.Scanln(&token)
			if err != nil {
				return err
			}
			
			// Save both token and api_url
			if err := config.SaveConfig(cfg.APIURL, token); err != nil {
				return err
			}
			
			fmt.Printf("✅ Logged in to %s\n", cfg.APIURL)
			return nil
		},
	}
	return c
}