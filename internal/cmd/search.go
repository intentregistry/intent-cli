package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"
	"os"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/intentregistry/intent-cli/internal/httpclient"
	"github.com/spf13/cobra"
)

func SearchCmd() *cobra.Command {
	var jsonOutput bool
	
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search public intents",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := strings.Join(args, " ")
			cfg := config.Load()
			if apiURLFlag != "" {
				cfg.APIURL = apiURLFlag
			}
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
			
			if jsonOutput {
				jsonData, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(jsonData))
				return nil
			}
			
			if len(resp.Items) == 0 {
				fmt.Println("No intents found matching your query.")
				return nil
			}
			
			// Use tabwriter for better column alignment
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SLUG\tDESCRIPTION\tOWNER")
			fmt.Fprintln(w, "----\t-----------\t-----")
			
			for _, it := range resp.Items {
				// Truncate description to 60 characters
				desc := it.Desc
				if len(desc) > 60 {
					desc = desc[:57] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", it.Slug, desc, it.Owner)
			}
			w.Flush()
			
			return nil
		},
	}
	
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	return cmd
}