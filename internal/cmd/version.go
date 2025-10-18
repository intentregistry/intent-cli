// internal/cmd/version.go
package cmd

import (
	"fmt"

	"github.com/intentregistry/intent-cli/internal/version"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show detailed version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("intent", version.Long())
		},
	}
}