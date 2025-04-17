package version_test

import (
	"runtime"
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/version"
	"github.com/stretchr/testify/assert"
)

func TestBuildInfo(t *testing.T) {
	// Save original values and restore them after test
	origVersion := version.Version
	origCommit := version.Commit
	origDate := version.Date
	defer func() {
		version.Version = origVersion
		version.Commit = origCommit
		version.Date = origDate
	}()

	tests := []struct {
		name     string
		version  string
		commit   string
		date     string
		expected string
	}{
		{
			name:     "Full information",
			version:  "1.2.3",
			commit:   "abc123",
			date:     "2025-04-17",
			expected: "1.2.3 (abc123, 2025-04-17)",
		},
		{
			name:     "Missing commit and date",
			version:  "1.2.3",
			commit:   "",
			date:     "",
			expected: "1.2.3",
		},
		{
			name:     "Only commit is empty",
			version:  "1.2.3",
			commit:   "",
			date:     "2025-04-17",
			expected: "1.2.3",
		},
		{
			name:     "Only date is empty",
			version:  "1.2.3",
			commit:   "abc123",
			date:     "",
			expected: "1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test values
			version.Version = tt.version
			version.Commit = tt.commit
			version.Date = tt.date

			// Execute function and verify results
			assert.Equal(t, tt.expected, version.BuildInfo())
		})
	}
}

func TestInfo(t *testing.T) {
	// Save original values and restore them after test
	origVersion := version.Version
	origCommit := version.Commit
	origDate := version.Date
	defer func() {
		version.Version = origVersion
		version.Commit = origCommit
		version.Date = origDate
	}()

	// Set test values
	version.Version = "2.0.0"
	version.Commit = "def456"
	version.Date = "2025-05-01"

	// Get the info map
	info := version.Info()

	// Verify all values in the map
	assert.Equal(t, "2.0.0", info["version"])
	assert.Equal(t, "def456", info["commit"])
	assert.Equal(t, "2025-05-01", info["date"])
	assert.Equal(t, runtime.Version(), info["goVersion"])
	assert.Equal(t, runtime.GOOS+"/"+runtime.GOARCH, info["platform"])
}
