package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/intentregistry/intent-cli/internal/pack"
	"github.com/spf13/cobra"
)

func PackageCmd() *cobra.Command {
	var (
		path    string
		outDir  string
	)

	c := &cobra.Command{
		Use:   "package [path] [--out directory]",
		Short: "Package an intent directory into a tar.gz archive",
		Long: `Create a tar.gz package from an intent directory.
		
The package includes all files in the directory and generates a SHA256 checksum.
This command creates packages locally without publishing to the registry.

Examples:
  intent package examples/hello.itml
  intent package examples/hello.itml --out dist/
  intent package . --out ./packages`,
		Args: cobra.MaximumNArgs(1),
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

			// Check if path exists
			info, err := os.Stat(abs)
			if os.IsNotExist(err) {
				return fmt.Errorf("path does not exist: %s", abs)
			}
			if err != nil {
				return fmt.Errorf("failed to stat path: %w", err)
			}

			// If path is a file, package the directory containing it
			packageDir := abs
			packageName := filepath.Base(abs)
			if !info.IsDir() {
				packageDir = filepath.Dir(abs)
				// Use file name (without extension) for package name
				ext := filepath.Ext(packageName)
				packageName = packageName[:len(packageName)-len(ext)]
			} else {
				packageName = filepath.Base(packageDir)
			}

			fmt.Println("📦 Packing:", packageDir)

			// Determine output directory
			if outDir == "" {
				outDir = "."
			}
			outDirAbs, err := filepath.Abs(outDir)
			if err != nil {
				return fmt.Errorf("failed to resolve output directory: %w", err)
			}

			// Create output directory if it doesn't exist
			if err := os.MkdirAll(outDirAbs, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Generate package filename
			packageFilename := fmt.Sprintf("%s.tar.gz", packageName)
			packagePath := filepath.Join(outDirAbs, packageFilename)

			// Create the package
			tarball, sha, err := pack.TarGzToPath(packageDir, packagePath)
			if err != nil {
				return fmt.Errorf("failed to create package: %w", err)
			}

			fmt.Println("  →", tarball)
			fmt.Println("  sha256:", sha)
			fmt.Println("✅ Package created successfully")

			return nil
		},
	}

	c.Flags().StringVar(&outDir, "out", "", "Output directory for the package (default: current directory)")

	return c
}

