package version

import (
	"fmt"
	"runtime"
)

// These variables are set at build time via ldflags
var (
	Version   = "dev"
	Commit    = "unknown"
	Date      = "unknown"
	GoVersion = runtime.Version()
)

// GetVersion returns the version string
func GetVersion() string {
	return Version
}

// GetFullVersion returns the full version information
func GetFullVersion() string {
	return fmt.Sprintf("%s (commit: %s, built: %s, %s)", Version, Commit, Date, GoVersion)
}

// GetBuildInfo returns build information as a map
func GetBuildInfo() map[string]string {
	return map[string]string{
		"version":    Version,
		"commit":     Commit,
		"date":       Date,
		"go_version": GoVersion,
	}
}
