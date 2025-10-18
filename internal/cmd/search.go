package cmd

import (
	"fmt"
	"strings"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/intentregistry/intent-cli/internal/httpclient"
	"github.com/spf13/cobra"
)

func SearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search public intents",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := strings.Join(args, " ")
			cfg := config.Load()
			cl := httpclient.NewWithDebug(cfg, Debug())
			var resp struct {
				Items []struct {
					Slug  string `json:"slug"`
					Desc  string `json:"summary"`
					Owner string `json:"owner"`
				} `json:"items"`
			}
			if err := cl.Get("/v1/search?q="+q, &resp); err != nil {
				return err
			}
			for _, it := range resp.Items {
				fmt.Printf("• %s — %s (by %s)\n", it.Slug, it.Desc, it.Owner)
			}
			return nil
		},
	}
}