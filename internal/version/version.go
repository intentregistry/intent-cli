package version

import (
	"fmt"
	"strings"
)

var (
	Version = "0.1.0"
	Commit  = ""
	Date    = ""
)

// GetVersion returns the full version string including dev suffix if applicable
func GetVersion() string {
	if Version == "dev" && Commit != "" {
		return fmt.Sprintf("dev+%s", Commit[:8])
	}
	if Commit != "" && !strings.Contains(Version, "dev") {
		return fmt.Sprintf("%s-dev+%s", Version, Commit[:8])
	}
	return Version
}

// IsDev returns true if this is a development build
func IsDev() bool {
	return Commit != "" && !strings.Contains(Version, "dev")
}