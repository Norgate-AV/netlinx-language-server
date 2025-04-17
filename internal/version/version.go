package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the semantic version
	Version = "dev"

	// Commit is the git commit SHA
	Commit = "none"

	// Date is the build date
	Date = "unknown"
)

func BuildInfo() string {
	if Commit != "" && Date != "" {
		return fmt.Sprintf("%s (%s, %s)", Version, Commit, Date)
	}

	return Version
}

func Info() map[string]string {
	return map[string]string{
		"version":   Version,
		"commit":    Commit,
		"date":      Date,
		"goVersion": runtime.Version(),
		"platform":  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
