package main

import (
	"fmt"
	"os"

	"github.com/intentregistry/intent-cli/internal/cmd"
)

func main() {
	root := cmd.RootCmd()
	// Registrar subcomandos
	root.AddCommand(
		cmd.LoginCmd(),
		cmd.PublishCmd(),
		cmd.InstallCmd(),
		cmd.WhoamiCmd(),
		cmd.SearchCmd(),
	)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}