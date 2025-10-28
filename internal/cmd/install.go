package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/intentregistry/intent-cli/internal/config"
	"github.com/intentregistry/intent-cli/internal/httpclient"
	"github.com/intentregistry/intent-cli/internal/pack"
	"github.com/spf13/cobra"
)

func InstallCmd() *cobra.Command {
	var dest string
	c := &cobra.Command{
		Use:   "install <@scope/name[@version]>",
		Short: "Install an intent package to local project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
            spec := args[0]
            cfg := config.Load()
            if apiURLFlag != "" {
                cfg.APIURL = apiURLFlag
            }
            cl := httpclient.NewWithDebug(cfg, Debug())

            // Resolve package metadata
            fmt.Println("🔎 Resolving", spec)
            // Fetch metadata as a generic map to avoid decoding issues
            var metaMap map[string]string
            if err := cl.Get("/v1/packages/resolve?spec="+spec, &metaMap); err != nil {
                return err
            }
            name := metaMap["name"]
            version := metaMap["version"]
            tarball := metaMap["tarball"]
            sha256sum := metaMap["sha256"]
            if tarball == "" || sha256sum == "" {
                return fmt.Errorf("invalid metadata received for %s", spec)
            }

            // Prepare download
            dlDir := filepath.Join(".intent-cache", "downloads")
            dlPath := filepath.Join(dlDir, sanitizeFilename(name+"-"+version)+".tar.gz")
            if err := cl.Download(tarball, dlPath); err != nil {
                return err
            }

            // Verify checksum
            sum, err := cl.SHA256(dlPath)
            if err != nil { return err }
            if !strings.EqualFold(sum, sha256sum) {
                return fmt.Errorf("checksum mismatch: got %s expected %s", sum, sha256sum)
            }

            // Extract
            targetDir := filepath.Join(dest, sanitizeFilename(name))
            fmt.Println("📦 Extracting to", targetDir)
            if err := pack.UntarGz(dlPath, targetDir); err != nil {
                return err
            }

            // Write install manifest
            installed := struct {
                Name    string `json:"name"`
                Version string `json:"version"`
                Source  string `json:"source"`
                Sha256  string `json:"sha256"`
            }{Name: name, Version: version, Source: tarball, Sha256: sum}
            manPath := filepath.Join(targetDir, ".installed.json")
            b, _ := json.MarshalIndent(installed, "", "  ")
            _ = os.WriteFile(manPath, b, 0o644)

            fmt.Printf("✅ Installed %s@%s into %s\n", name, version, targetDir)
            return nil
		},
	}
	c.Flags().StringVar(&dest, "dest", "intents", "destination folder")
	return c
}

var filenameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9._@-]+`)

func sanitizeFilename(s string) string {
    return filenameSanitizer.ReplaceAllString(s, "-")
}