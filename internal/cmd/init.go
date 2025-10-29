package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intentregistry/intent-cli/internal/pack"
	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	var (
		force  bool
		app    bool
		scope  string
	)

	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new intent project",
		Long: `Initialize a new intent project with all required files and directories.

Creates a complete project structure including:
- itpkg.json manifest
- intents/ directory with example hello.itml
- policies/ directory with security policy
- README.md with project documentation

Examples:
  intent init                    # Initialize in current directory
  intent init my-awesome-intent  # Create new directory and initialize
  intent init --app              # Create an app package (with entrypoint)
  intent init --scope @acme      # Set package scope to @acme`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			var projectDir string

			if len(args) > 0 {
				projectName = args[0]
				projectDir = projectName
				
				// Check if directory already exists
				if _, err := os.Stat(projectDir); err == nil {
					if !force {
						return fmt.Errorf("directory %s already exists. Use --force to overwrite", projectDir)
					}
				} else {
					// Create directory
					if err := os.MkdirAll(projectDir, 0755); err != nil {
						return fmt.Errorf("failed to create directory: %w", err)
					}
				}
			} else {
				// Use current directory
				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				projectDir = "."
				projectName = filepath.Base(dir)
			}

			// Validate project name
			if projectName == "" || strings.ContainsAny(projectName, " \t\n\r/\\") {
				return fmt.Errorf("invalid project name: %q (must be non-empty and not contain spaces or path separators)", projectName)
			}

			// Check if project already exists
			manifestPath := filepath.Join(projectDir, "itpkg.json")
			if _, err := os.Stat(manifestPath); err == nil && !force {
				return fmt.Errorf("itpkg.json already exists in %s. Use --force to overwrite", projectDir)
			}

			fmt.Printf("🚀 Initializing intent project: %s\n", projectName)
			fmt.Println()

			// Create required directories
			dirs := []string{"intents", "policies"}
			for _, dir := range dirs {
				dirPath := filepath.Join(projectDir, dir)
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					return fmt.Errorf("failed to create %s directory: %w", dir, err)
				}
				fmt.Printf("📁 Created directory: %s/\n", dir)
			}

			// Determine package name
			pkgName := projectName
			if scope != "" {
				scope = strings.TrimPrefix(scope, "@")
				pkgName = fmt.Sprintf("@%s/%s", scope, projectName)
			} else {
				pkgName = fmt.Sprintf("@scope/%s", projectName)
			}

			// Create itpkg.json
			manifest := pack.ItpkgManifest{
				Name:        pkgName,
				Version:     "0.1.0",
				Description: fmt.Sprintf("Intent package for %s", projectName),
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

			if app {
				manifest.Type = "app"
				manifest.Entry = "project.app.itml"
				manifest.Capabilities = []string{"ui.render", "http.outbound"}
			} else {
				manifest.Type = "lib"
			}

			manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal manifest: %w", err)
			}

			if err := os.WriteFile(manifestPath, manifestJSON, 0644); err != nil {
				return fmt.Errorf("failed to write itpkg.json: %w", err)
			}
			fmt.Printf("📄 Created itpkg.json\n")

			// Create example hello.itml
			helloPath := filepath.Join(projectDir, "intents", "hello.itml")
			helloContent := `intent "Hello World"
inputs:
  - name (string) default="World"
workflow:
  → log("Hello {name}!")
  → return(status="ok", message="Hello {name}!")
`
			if err := os.WriteFile(helloPath, []byte(helloContent), 0644); err != nil {
				return fmt.Errorf("failed to write hello.itml: %w", err)
			}
			fmt.Printf("📝 Created intents/hello.itml\n")

			// Create example security policy
			securityPath := filepath.Join(projectDir, "policies", "security.itml")
			securityContent := `security:
  network:
    outbound:
      deny: ["*"]
      allow: []
  filesystem:
    read: ["intents/**", "policies/**"]
    write: []
`
			if err := os.WriteFile(securityPath, []byte(securityContent), 0644); err != nil {
				return fmt.Errorf("failed to write security policy: %w", err)
			}
			fmt.Printf("🔒 Created policies/security.itml\n")

			// Create project.app.itml for app packages
			if app {
				appPath := filepath.Join(projectDir, "project.app.itml")
				appContent := `app "` + projectName + `"
version: "0.1.0"
description: "A new intent application"

routes:
  - path: "/"
    intent: "hello"
`
				if err := os.WriteFile(appPath, []byte(appContent), 0644); err != nil {
					return fmt.Errorf("failed to write project.app.itml: %w", err)
				}
				fmt.Printf("📱 Created project.app.itml\n")
			}

			// Create README.md
			readmePath := filepath.Join(projectDir, "README.md")
			appEntry := ""
			if app {
				appEntry = "├── project.app.itml   # Application entrypoint\n"
			}
			readmeContent := fmt.Sprintf(`# %s

An Intent package built with Intent CLI.

## Project Structure

%s/
├── itpkg.json          # Package manifest
├── intents/            # Intent definitions
│   └── hello.itml     # Example intent
├── policies/           # Security and privacy policies
│   └── security.itml  # Network and filesystem policies
%s

## Getting Started

### Run an Intent

intent run intents/hello.itml --inputs name=World

### Package

intent package . --out dist/

### Publish

intent publish . --tag beta

## Documentation

- [ITML Format](https://docs.intentregistry.com/itml)
- [Package Format](https://docs.intentregistry.com/itpkg)
`, projectName, projectName, appEntry)
			if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
				return fmt.Errorf("failed to write README.md: %w", err)
			}
			fmt.Printf("📖 Created README.md\n")

			fmt.Println()
			fmt.Println("✅ Project initialized successfully!")
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Printf("  1. cd %s\n", projectDir)
			fmt.Println("  2. Edit intents/hello.itml or create new intents")
			fmt.Println("  3. Run: intent run intents/hello.itml --inputs name=World")
			fmt.Println("  4. Package: intent package . --out dist/")
			fmt.Println("  5. Publish: intent publish .")

			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")
	cmd.Flags().BoolVar(&app, "app", false, "Create an app package (with project.app.itml entrypoint)")
	cmd.Flags().StringVar(&scope, "scope", "", "Set package scope (e.g., @acme)")

	return cmd
}
