package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetHashFromPseudoVersion(t *testing.T) {
	tests := []struct {
		name          string
		pseudoVersion string
		expectedHash  string
		expectedOK    bool
	}{
		{
			name:          "valid pseudo version",
			pseudoVersion: "v1.1.1-20250217181409-44e5ddce290d",
			expectedHash:  "44e5ddce290d",
			expectedOK:    true,
		},
		{
			name:          "valid pseudo version with different format",
			pseudoVersion: "v0.2.3-20210101120000-abcdef123456",
			expectedHash:  "abcdef123456",
			expectedOK:    true,
		},
		{
			name:          "regular version tag",
			pseudoVersion: "v1.2.3",
			expectedHash:  "",
			expectedOK:    false,
		},
		{
			name:          "invalid version format",
			pseudoVersion: "invalid-version",
			expectedHash:  "",
			expectedOK:    false,
		},
		{
			name:          "empty string",
			pseudoVersion: "",
			expectedHash:  "",
			expectedOK:    false,
		},
		{
			name:          "pseudo version with multiple hyphens",
			pseudoVersion: "v1.0.0-rc1-0.20250101000000-abc123def456",
			expectedHash:  "abc123def456",
			expectedOK:    true,
		},
		{
			name:          "semver alpha",
			pseudoVersion: "v1.0.0-alpha.0",
			expectedHash:  "",
			expectedOK:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, ok := GetHashFromPseudoVersion(tt.pseudoVersion)
			require.Equal(t, tt.expectedOK, ok)
			require.Equal(t, tt.expectedHash, hash)
		})
	}
}
