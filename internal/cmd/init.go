package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	var force bool
	
	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new intent project",
		Long: `Initialize a new intent project in the current directory.
Creates a manifest.yaml file with the basic structure for an intent.

Examples:
  intent init                    # Initialize in current directory
  intent init my-awesome-intent  # Initialize with specific name
  intent init --force           # Overwrite existing manifest.yaml`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var name string
			if len(args) > 0 {
				name = args[0]
			} else {
				// Use current directory name as default
				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				name = filepath.Base(dir)
			}
			
			// Validate name
			if name == "" || strings.ContainsAny(name, " \t\n\r/\\") {
				return fmt.Errorf("invalid intent name: %q (must be non-empty and not contain spaces or path separators)", name)
			}
			
			manifestPath := "manifest.yaml"
			
			// Check if manifest already exists
			if _, err := os.Stat(manifestPath); err == nil && !force {
				return fmt.Errorf("manifest.yaml already exists. Use --force to overwrite")
			}
			
			// Create manifest content
			manifest := fmt.Sprintf(`name: %s
version: "1.0.0"
description: "A new intent for %s"
author: ""
license: "MIT"
tags: []
parameters:
  # Define your intent parameters here
  # Example:
  # - name: "query"
  #   type: "string"
  #   description: "The search query"
  #   required: true
  #   default: ""
outputs:
  # Define your intent outputs here
  # Example:
  # - name: "result"
  #   type: "string"
  #   description: "The processed result"
examples:
  # Add usage examples here
  # Example:
  # - input: "search for cats"
  #   output: "Found 5 results about cats"
`, name, name)
			
			// Write manifest file
			if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
				return fmt.Errorf("failed to write manifest.yaml: %w", err)
			}
			
			fmt.Printf("✅ Initialized intent project: %s\n", name)
			fmt.Printf("📄 Created manifest.yaml\n")
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Println("  1. Edit manifest.yaml to define your intent parameters and outputs")
			fmt.Println("  2. Add your intent implementation")
			fmt.Println("  3. Run 'intent publish' to publish your intent")
			
			return nil
		},
	}
	
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing manifest.yaml")
	return cmd
}
