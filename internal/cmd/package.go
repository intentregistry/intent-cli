package cmd

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intentregistry/intent-cli/internal/pack"
	"github.com/spf13/cobra"
)

func PackageCmd() *cobra.Command {
	var (
		path       string
		outDir     string
		unsigned   bool
		signKeyPath string
		scaffold   bool
	)

	c := &cobra.Command{
		Use:   "package [path] [--out directory]",
		Short: "Package an intent directory into a signed .itpkg archive",
		Long: `Create a signed .itpkg package from an intent directory.

The package must contain an itpkg.json manifest with name, version, entry, and policies.
The directory structure must include intents/ and policies/ directories.

Default behavior requires a valid ed25519 signing key. Use --unsigned to create
unsigned packages (not recommended for production).

Examples:
  intent package examples/hello
  intent package . --out dist/
  intent package . --scaffold  # Generate itpkg.json if missing
  intent package . --sign-key ~/.ssh/intent_sign_key`,
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

			// Check for itpkg.json
			manifestPath := filepath.Join(packageDir, "itpkg.json")
			if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
				if scaffold {
					fmt.Printf("📝 Scaffolding itpkg.json...\n")
					if err := scaffoldItpkgJSON(packageDir, packageName); err != nil {
						return fmt.Errorf("failed to scaffold itpkg.json: %w", err)
					}
					fmt.Println("✅ Created itpkg.json")
				} else {
					return fmt.Errorf("itpkg.json not found in %s\n\nTo generate it automatically, use:\n  intent package . --scaffold --unsigned\n\nOr create itpkg.json manually with name, version, itmlVersion, and policies.", packageDir)
				}
			} else if scaffold {
				// Even if manifest exists, ensure required directories exist when scaffold flag is set
				requiredDirs := []string{"intents", "policies"}
				for _, dirName := range requiredDirs {
					dirPath := filepath.Join(packageDir, dirName)
					if _, err := os.Stat(dirPath); os.IsNotExist(err) {
						if err := os.MkdirAll(dirPath, 0755); err != nil {
							return fmt.Errorf("failed to create directory %s: %w", dirName, err)
						}
						fmt.Printf("📁 Created directory: %s/\n", dirName)
					}
				}
			}

			fmt.Println("📦 Packing:", packageDir)

			// Load signing key
			var signKey ed25519.PrivateKey
			if !unsigned {
				if signKeyPath == "" {
					signKeyPath = os.Getenv("INTENT_SIGN_KEY")
				}
				if signKeyPath == "" {
					return fmt.Errorf("signing key required (use --sign-key, INTENT_SIGN_KEY env, or --unsigned)")
				}
				var err error
				signKey, err = loadEd25519Key(signKeyPath)
				if err != nil {
					return fmt.Errorf("failed to load signing key: %w", err)
				}
			}

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
			manifest, err := pack.ReadItpkgManifest(manifestPath)
			if err != nil {
				return fmt.Errorf("failed to read manifest: %w", err)
			}
			
			// Use name from manifest if available, otherwise use directory name
			pkgName := manifest.Name
			if pkgName == "" {
				pkgName = packageName
			}
			// Sanitize package name for filename
			pkgName = sanitizePackageName(pkgName)
			packageFilename := fmt.Sprintf("%s-%s.itpkg", pkgName, manifest.Version)
			packagePath := filepath.Join(outDirAbs, packageFilename)

			// Create the package
			itpkg, err := pack.CreateItpkg(packageDir, packagePath, signKey, unsigned)
			if err != nil {
				return fmt.Errorf("failed to create .itpkg: %w", err)
			}
			fmt.Println("  →", itpkg)
			fmt.Println("✅ Package created successfully")

			return nil
		},
	}

	c.Flags().StringVar(&outDir, "out", "", "Output directory for the package (default: current directory)")
	c.Flags().BoolVar(&unsigned, "unsigned", false, "Allow creating unsigned .itpkg (not recommended)")
	c.Flags().StringVar(&signKeyPath, "sign-key", "", "Path to ed25519 private key file (defaults to env INTENT_SIGN_KEY)")
	c.Flags().BoolVar(&scaffold, "scaffold", false, "Generate itpkg.json if missing")

	return c
}

// scaffoldItpkgJSON creates a minimal itpkg.json file and required directories
func scaffoldItpkgJSON(dir, name string) error {
	// Create required directories if they don't exist
	requiredDirs := []string{"intents", "policies"}
	for _, dirName := range requiredDirs {
		dirPath := filepath.Join(dir, dirName)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dirName, err)
			}
		}
	}

	// Check if entrypoint exists
	entryPoint := "project.app.itml"
	entryPath := filepath.Join(dir, entryPoint)
	if _, err := os.Stat(entryPath); os.IsNotExist(err) {
		entryPoint = "" // Will be lib type
	}

	manifest := pack.ItpkgManifest{
		Name:        fmt.Sprintf("@scope/%s", name),
		Version:     "0.1.0",
		Description: fmt.Sprintf("Intent package for %s", name),
		ItmlVersion: "0.1",
		Capabilities: []string{},
		Policies: map[string]interface{}{
			"security": map[string]interface{}{
				"network": map[string]interface{}{
					"outbound": map[string]interface{}{
						"deny": []string{"*"},
					},
				},
			},
			"privacy": map[string]interface{}{
				"pii": map[string]interface{}{
					"export": "deny",
				},
			},
			"energy": map[string]interface{}{
				"mode": "balanced",
			},
		},
	}

	if entryPoint == "" {
		manifest.Type = "lib"
	} else {
		manifest.Entry = entryPoint
		manifest.Type = "app"
		manifest.Capabilities = []string{"ui.render", "http.outbound"}
	}

	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	manifestPath := filepath.Join(dir, "itpkg.json")
	return os.WriteFile(manifestPath, manifestJSON, 0644)
}

// loadEd25519Key loads an ed25519 private key from a file (hex or PEM format)
func loadEd25519Key(path string) (ed25519.PrivateKey, error) {
	// Expand ~ and environment variables in path
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}
	path = os.ExpandEnv(path) // Expand $VAR and ${VAR}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file %s: %w", path, err)
	}

	// Try hex format first (64 bytes = 128 hex chars)
	if len(data) == 128 || len(data) == 129 { // 128 or 128+\n
		hexKey := strings.TrimSpace(string(data))
		keyBytes, err := hex.DecodeString(hexKey)
		if err == nil && len(keyBytes) == ed25519.PrivateKeySize {
			return ed25519.PrivateKey(keyBytes), nil
		}
	}

	// TODO: Add PEM format support if needed
	return nil, fmt.Errorf("unsupported key format; expected hex-encoded ed25519 private key (%d bytes)", ed25519.PrivateKeySize*2)
}

// sanitizePackageName sanitizes a package name for use in filenames
func sanitizePackageName(name string) string {
	// Replace @scope/name with scope-name
	name = strings.ReplaceAll(name, "@", "")
	name = strings.ReplaceAll(name, "/", "-")
	// Remove any other invalid characters
	var sanitized strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			sanitized.WriteRune(r)
		}
	}
	result := sanitized.String()
	if result == "" {
		result = "package"
	}
	return result
}
