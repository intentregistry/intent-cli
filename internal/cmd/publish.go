package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/intentregistry/intent-cli/internal/httpclient"
	"github.com/intentregistry/intent-cli/internal/pack"
	"github.com/spf13/cobra"
)

func PublishCmd() *cobra.Command {
	var (
		path     string
		isPrivate bool
		tag      string
		message  string
	)
		c := &cobra.Command{
		Use:   "publish [path]",
		Short: "Publish an intent package",
		Long: `Publish an intent package to the registry.

If path is a .itpkg file, it will be published directly.
If path is a directory, it will be packaged first, then published.

Examples:
  intent publish dist/package-0.1.0.itpkg    # Publish existing .itpkg file
  intent publish .                           # Package and publish directory
  intent publish . --tag beta                 # Publish as beta release
  intent publish . --private                 # Publish as private package`,
		Args: cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// Enable file completion for .itpkg files and directories
			return nil, cobra.ShellCompDirectiveDefault
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}

			abs, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}

			// Check if path is a .itpkg file
			var tarball string
			var sha string

			if strings.HasSuffix(strings.ToLower(abs), ".itpkg") {
				// It's already a package file - use it directly
				info, err := os.Stat(abs)
				if err != nil {
					return fmt.Errorf("package file not found: %s", abs)
				}
				if info.IsDir() {
					return fmt.Errorf("expected .itpkg file, but %s is a directory", abs)
				}
				tarball = abs
				fmt.Println("📦 Publishing package:", abs)

				// Calculate SHA256 of the file
				file, err := os.Open(abs)
				if err != nil {
					return fmt.Errorf("failed to open package file: %w", err)
				}
				defer file.Close()

				hasher := sha256.New()
				if _, err := io.Copy(hasher, file); err != nil {
					return fmt.Errorf("failed to calculate checksum: %w", err)
				}
				sha = hex.EncodeToString(hasher.Sum(nil))
				fmt.Println("  sha256:", sha)
			} else {
				// It's a directory - package it first
				fmt.Println("📦 Packing:", abs)
				tarball, sha, err = pack.TarGz(abs)
				if err != nil {
					return fmt.Errorf("failed to create package: %w", err)
				}
				fmt.Println("  →", tarball)
				fmt.Println("  sha256:", sha)
			}

			cfg := config.Load()
			if apiURLFlag != "" {
				cfg.APIURL = apiURLFlag
			}

			// Check authentication before attempting to publish
			if cfg.Token == "" {
				return fmt.Errorf("authentication required\n\nTo publish packages, you need to authenticate:\n  1. Run: intent login\n  2. Enter your API token\n  3. Then try publishing again\n\nOr set INTENT_TOKEN environment variable")
			}

			if cfg.APIURL == "" {
				return fmt.Errorf("API URL not configured\n\nSet INTENT_API_URL environment variable or run 'intent login'")
			}

			cl := httpclient.NewWithDebug(cfg, Debug())

			payload := map[string]any{
				"private": isPrivate,
				"tag":     tag,
				"message": message,
				"sha256":  sha,
			}

			fmt.Printf("📤 Publishing to: %s\n", cfg.APIURL)
			// POST multipart: file + payload
			if err := cl.PostMultipart("/v1/packages/publish", payload, "file", tarball, nil); err != nil {
				return fmt.Errorf("failed to publish: %w\n\nMake sure:\n  - You're authenticated (run 'intent login')\n  - The API endpoint is correct (set --api-url or INTENT_API_URL)\n  - Your token has publish permissions", err)
			}
			fmt.Println("✅ Published successfully")
			return nil
		},
	}
	c.Flags().BoolVar(&isPrivate, "private", false, "publish as private")
	c.Flags().StringVar(&tag, "tag", "", "beta|rc")
	c.Flags().StringVar(&message, "message", "", "release note")
	return c
}