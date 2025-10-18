package cmd

import (
	"fmt"
	"path/filepath"

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
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}
			abs, _ := filepath.Abs(path)
			fmt.Println("📦 Packing:", abs)
			tarball, sha, err := pack.TarGz(abs) // devuelve ruta del .tar.gz y sha256
			if err != nil {
				return err
			}
			fmt.Println("  →", tarball)
			fmt.Println("  sha256:", sha)

			cfg := config.Load()
			cl := httpclient.New(cfg)

			payload := map[string]any{
				"private": isPrivate,
				"tag":     tag,
				"message": message,
				"sha256":  sha,
			}
			// POST multipart: file + payload
			if err := cl.PostMultipart("/v1/packages/publish", payload, "file", tarball, nil); err != nil {
				return err
			}
			fmt.Println("✅ Published")
			return nil
		},
	}
	c.Flags().BoolVar(&isPrivate, "private", false, "publish as private")
	c.Flags().StringVar(&tag, "tag", "", "beta|rc")
	c.Flags().StringVar(&message, "message", "", "release note")
	return c
}