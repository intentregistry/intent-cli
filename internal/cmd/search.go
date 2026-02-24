package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/intentregistry/intent-cli/internal/httpclient"
	"github.com/spf13/cobra"
)

func SearchCmd() *cobra.Command {
	var (
		jsonOutput bool
		limit      int
		owner      string
		sortBy     string
	)

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search public intents",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit < 0 {
				return fmt.Errorf("--limit must be >= 0")
			}

			sortBy = strings.ToLower(strings.TrimSpace(sortBy))
			switch sortBy {
			case "", "relevance", "slug", "owner":
			default:
				return fmt.Errorf("invalid --sort value %q (valid: relevance, slug, owner)", sortBy)
			}

			q := strings.Join(args, " ")
			cfg := config.Load()
			if apiURLFlag != "" {
				cfg.APIURL = apiURLFlag
			}
			cl := httpclient.NewWithDebug(cfg, Debug())
			var resp searchAPIResponse
			if err := cl.Get("/v1/search?q="+url.QueryEscape(q), &resp); err != nil {
				return err
			}

			items := normalizeSearchItems(resp)
			items = filterSearchItems(items, owner)
			items = sortSearchItems(items, sortBy)
			items = limitSearchItems(items, limit)

			if jsonOutput {
				out := map[string]any{
					"count": len(items),
					"items": items,
				}
				jsonData, err := json.MarshalIndent(out, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(jsonData))
				return nil
			}

			if len(items) == 0 {
				fmt.Println("No intents found matching your query.")
				return nil
			}

			// Use tabwriter for better column alignment
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			if _, err := fmt.Fprintln(w, "SLUG\tDESCRIPTION\tOWNER"); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			if _, err := fmt.Fprintln(w, "----\t-----------\t-----"); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}

			for _, it := range items {
				// Truncate description to 60 characters
				desc := it.Summary
				if len(desc) > 60 {
					desc = desc[:57] + "..."
				}
				if _, err := fmt.Fprintf(w, "%s\t%s\t%s\n", it.Slug, desc, it.Owner); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("failed to flush output: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	cmd.Flags().IntVar(&limit, "limit", 20, "Maximum number of results (0 = no limit)")
	cmd.Flags().StringVar(&owner, "owner", "", "Filter results by owner (case-insensitive)")
	cmd.Flags().StringVar(&sortBy, "sort", "relevance", "Sort order: relevance|slug|owner")
	return cmd
}

type searchAPIItem struct {
	Slug    string `json:"slug"`
	Summary string `json:"summary"`
	Owner   string `json:"owner"`
}

type searchAPIPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type searchAPIResponse struct {
	Items    []searchAPIItem    `json:"items"`
	Packages []searchAPIPackage `json:"packages"`
}

func normalizeSearchItems(resp searchAPIResponse) []searchAPIItem {
	items := make([]searchAPIItem, 0, len(resp.Items)+len(resp.Packages))
	items = append(items, resp.Items...)

	for _, pkg := range resp.Packages {
		summary := ""
		if pkg.Version != "" {
			summary = "version " + pkg.Version
		}
		items = append(items, searchAPIItem{
			Slug:    pkg.Name,
			Summary: summary,
			Owner:   ownerFromSlug(pkg.Name),
		})
	}
	return items
}

func filterSearchItems(items []searchAPIItem, owner string) []searchAPIItem {
	owner = strings.TrimSpace(strings.TrimPrefix(owner, "@"))
	if owner == "" {
		return items
	}

	out := make([]searchAPIItem, 0, len(items))
	for _, it := range items {
		if strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(it.Owner, "@")), owner) {
			out = append(out, it)
		}
	}
	return out
}

func sortSearchItems(items []searchAPIItem, sortBy string) []searchAPIItem {
	if sortBy == "" || sortBy == "relevance" {
		return items
	}
	cloned := append([]searchAPIItem(nil), items...)
	switch sortBy {
	case "slug":
		sort.SliceStable(cloned, func(i, j int) bool {
			return strings.ToLower(cloned[i].Slug) < strings.ToLower(cloned[j].Slug)
		})
	case "owner":
		sort.SliceStable(cloned, func(i, j int) bool {
			oi := strings.ToLower(cloned[i].Owner)
			oj := strings.ToLower(cloned[j].Owner)
			if oi == oj {
				return strings.ToLower(cloned[i].Slug) < strings.ToLower(cloned[j].Slug)
			}
			return oi < oj
		})
	}
	return cloned
}

func limitSearchItems(items []searchAPIItem, limit int) []searchAPIItem {
	if limit <= 0 || len(items) <= limit {
		return items
	}
	return items[:limit]
}

func ownerFromSlug(slug string) string {
	if strings.HasPrefix(slug, "@") {
		parts := strings.SplitN(strings.TrimPrefix(slug, "@"), "/", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) != "" {
			return parts[0]
		}
	}
	return ""
}
