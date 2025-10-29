package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
			
			reader := bufio.NewReader(os.Stdin)
			
			// Prompt for API URL if not already set
			if cfg.APIURL == "" || cfg.APIURL == "https://api.intentregistry.com" {
				fmt.Print("Enter API URL (default: https://api.intentregistry.com): ")
				apiURL, _ := reader.ReadString('\n')
				apiURL = strings.TrimSpace(apiURL)
				if apiURL != "" {
					cfg.APIURL = apiURL
				} else {
					cfg.APIURL = "https://api.intentregistry.com"
				}
			} else {
				// API URL already set, show it
				fmt.Printf("Using API URL: %s\n", cfg.APIURL)
			}
			
			// Prompt for token (optional for local dev)
			fmt.Print("Enter API token (optional, press Enter to skip): ")
			token, _ := reader.ReadString('\n')
			token = strings.TrimSpace(token)
			// Empty token is OK for local development
			
			// Save both token and api_url
			if err := config.SaveConfig(cfg.APIURL, token); err != nil {
				return err
			}
			
			if token == "" {
				fmt.Printf("✅ Configured for %s (no authentication token)\n", cfg.APIURL)
			} else {
				fmt.Printf("✅ Logged in to %s\n", cfg.APIURL)
			}
			return nil
		},
	}
	return c
}