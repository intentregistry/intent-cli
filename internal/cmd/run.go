package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/intentregistry/intent-cli/internal/executor"
	"github.com/intentregistry/intent-cli/internal/parser"
	"github.com/spf13/cobra"
)

func RunCmd() *cobra.Command {
	var (
		inputs    []string
		outputDir string
		verbose   bool
	)
	
	c := &cobra.Command{
		Use:   "run FILE.itml [--inputs k=v]",
		Short: "Execute an intent file with optional input parameters",
		Long: `Execute an intent file (.itml) with optional input parameters.

The --inputs flag allows you to pass key-value pairs that will be available
to the intent during execution. Multiple inputs can be provided.

Examples:
  intent run my-intent.itml
  intent run my-intent.itml --inputs name=John --inputs age=30
  intent run my-intent.itml --inputs query="search for cats" --output-dir ./results`,
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// Enable file completion for .itml files
			return nil, cobra.ShellCompDirectiveDefault
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			itmlFile := args[0]
			
			// Validate file exists and has .itml extension
			if !strings.HasSuffix(itmlFile, ".itml") {
				return fmt.Errorf("file must have .itml extension: %s", itmlFile)
			}
			
			if _, err := os.Stat(itmlFile); os.IsNotExist(err) {
				return fmt.Errorf("intent file not found: %s", itmlFile)
			}
			
			// Parse input parameters
			inputParams := make(map[string]string)
			for _, input := range inputs {
				parts := strings.SplitN(input, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid input format '%s', expected 'key=value'", input)
				}
				inputParams[parts[0]] = parts[1]
			}
			
			if verbose {
				fmt.Printf("🔍 Parsing intent file: %s\n", itmlFile)
			}
			
			// Parse the .itml file
			intent, err := parser.ParseITML(itmlFile)
			if err != nil {
				return fmt.Errorf("failed to parse intent file: %w", err)
			}
			
			if verbose {
				fmt.Printf("✅ Intent parsed successfully: %s\n", intent.Name)
				fmt.Printf("📝 Description: %s\n", intent.Description)
				fmt.Printf("🔧 Parameters: %d\n", len(intent.Parameters))
				fmt.Printf("📤 Outputs: %d\n", len(intent.Outputs))
			}
			
			// Validate required parameters
			for _, param := range intent.Parameters {
				if param.Required && inputParams[param.Name] == "" {
					return fmt.Errorf("required parameter '%s' not provided", param.Name)
				}
			}
			
			if verbose {
				fmt.Printf("🚀 Executing intent...\n")
			}
			
			// Execute the intent
			results, err := executor.Execute(intent, inputParams, outputDir)
			if err != nil {
				return fmt.Errorf("execution failed: %w", err)
			}
			
			// Display results
			fmt.Println("✅ Intent executed successfully!")
			fmt.Println()
			
			if len(results) > 0 {
				fmt.Println("📊 Results:")
				for name, value := range results {
					fmt.Printf("  %s: %v\n", name, value)
				}
			}
			
			if outputDir != "" {
				fmt.Printf("📁 Results saved to: %s\n", outputDir)
			}
			
			return nil
		},
	}
	
	c.Flags().StringSliceVar(&inputs, "inputs", []string{}, "Input parameters as key=value pairs (can be used multiple times)")
	c.Flags().StringVar(&outputDir, "output-dir", "", "Directory to save output files")
	c.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	
	return c
}
