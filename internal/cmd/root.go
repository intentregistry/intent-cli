package cmd

import (
	"github.com/intentregistry/intent-cli/internal/version"
	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
		Use:     "intent",
		Short:   "IntentRegistry CLI",
		Long:    "Publish & install AI Intents from intentregistry.com",
		Version: version.Short(),
	}
	debug         bool
	apiURLFlag    string
	telemetryFlag bool
)

func init() {
	// Version output template
	root.SetVersionTemplate("intent {{.Version}}\n")

	// Global flags
	root.PersistentFlags().BoolVar(&debug, "debug", false, "Enable verbose debug output")
	root.PersistentFlags().StringVar(&apiURLFlag, "api-url", "", "Override API base URL (env INTENT_API_URL)")
	root.PersistentFlags().BoolVar(&telemetryFlag, "telemetry", false, "Enable telemetry (env INTENT_TELEMETRY)")

	// (Optional) completion command UX options
	root.CompletionOptions.DisableDefaultCmd = false
	root.CompletionOptions.DisableNoDescFlag = false
	root.CompletionOptions.HiddenDefaultCmd = true

	// Attach subcommands
	root.AddCommand(
		InitCmd(),
		DoctorCmd(),
		LoginCmd(),
		RunCmd(),
		PublishCmd(),
		InstallCmd(),
		WhoamiCmd(),
		SearchCmd(),
		VersionCmd(),
		CompletionCmd(),
	)

	// Disable file fallback for all commands by default
	applyNoFileCompletion(root)

	// If you WANT file completion for a specific command, re-enable it here:
	// if pc := getCmd("publish"); pc != nil { pc.ValidArgsFunction = nil }
}

func RootCmd() *cobra.Command { return root }
func Execute() error          { return root.Execute() }
func Debug() bool             { return debug }
func Telemetry() bool         { return telemetryFlag }