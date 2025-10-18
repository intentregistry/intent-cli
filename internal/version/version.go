package version

// These are injected by GoReleaser via -ldflags at build time.
var (
	Version = "dev" // e.g. "0.2.8" on release; "dev" locally
	Commit  = ""    // full git SHA
	Date    = ""    // RFC3339 build time
)

// Short returns just the semantic version (no suffixes).
func Short() string {
	return Version
}

// Long returns a human-friendly long string (for "version --long").
func Long() string {
	shortSHA := Commit
	if len(shortSHA) > 7 {
		shortSHA = shortSHA[:7]
	}
	if shortSHA == "" && Date == "" {
		return Version
	}
	if shortSHA == "" {
		return Version + " (built " + Date + ")"
	}
	if Date == "" {
		return Version + " (commit " + shortSHA + ")"
	}
	return Version + " (commit " + shortSHA + ", built " + Date + ")"
}