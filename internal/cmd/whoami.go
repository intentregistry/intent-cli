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
			if apiURLFlag != "" {
				cfg.APIURL = apiURLFlag
			}
			
			// Always show API URL
			if cfg.APIURL == "" {
				fmt.Println("❌ API URL not configured")
				fmt.Println("   Set INTENT_API_URL or run 'intent login'")
				return nil
			}
			
			fmt.Printf("🔗 API URL: %s\n", cfg.APIURL)
			
			// Try to get user info if token exists
			if cfg.Token == "" {
				fmt.Println("⚠️  No authentication token configured")
				fmt.Println("   Authentication is optional for local development")
				fmt.Println("   Run 'intent login' if your API requires authentication")
				return nil
			}
			
			// Try to fetch user info
			cl := httpclient.NewWithDebug(cfg, Debug())
			var resp struct {
				Email  string `json:"email"`
				UserID string `json:"userId"`
				Username string `json:"username"`
			}
			
			// Try common endpoints
			endpoints := []string{"/v1/users/me", "/me", "/v1/whoami"}
			var err error
			for _, endpoint := range endpoints {
				err = cl.Get(endpoint, &resp)
				if err == nil {
					break
				}
			}
			
			if err != nil {
				fmt.Println("⚠️  Authentication token configured but API call failed")
				fmt.Printf("   Error: %v\n", err)
				fmt.Println("   This is OK for local development if your API doesn't require auth")
				return nil
			}
			
			// Show user info
			if resp.Email != "" {
				fmt.Printf("👤 Email: %s\n", resp.Email)
			}
			if resp.Username != "" {
				fmt.Printf("👤 Username: %s\n", resp.Username)
			}
			if resp.UserID != "" {
				fmt.Printf("🆔 User ID: %s\n", resp.UserID)
			}
			
			return nil
		},
	}
}