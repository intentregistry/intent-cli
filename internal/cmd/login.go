package cmd

import (
	"fmt"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/spf13/cobra"
)

func LoginCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "login",
		Short: "Authenticate against IntentRegistry",
		RunE: func(cmd *cobra.Command, args []string) error {
			// MVP: pedir token y guardarlo
			var token string
			fmt.Print("Enter API token: ")
			_, err := fmt.Scanln(&token)
			if err != nil {
				return err
			}
			if err := config.SaveToken(token); err != nil {
				return err
			}
			fmt.Println("✅ Logged in")
			return nil
		},
	}
	return c
}