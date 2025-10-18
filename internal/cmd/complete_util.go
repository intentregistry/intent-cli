package cmd

import "github.com/spf13/cobra"

// NoFileCompletion disables the fallback to file/dir completion.
func NoFileCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// applyNoFileCompletion recursively sets NoFileCompletion on commands
// that don't already provide their own ValidArgsFunction.
func applyNoFileCompletion(c *cobra.Command) {
	if c.ValidArgsFunction == nil {
		c.ValidArgsFunction = NoFileCompletion
	}
	for _, child := range c.Commands() {
		applyNoFileCompletion(child)
	}
}

// helper to fetch a subcommand by name (optional).
func getCmd(name string) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}