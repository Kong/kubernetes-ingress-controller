package helpers

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/require"
)

func TestGetKongImageVersion(t *testing.T) {
	testCases := []struct {
		name             string
		effectiveVersion string
		tag              string
		expected         semver.Version
		expectError      bool
	}{
		{
			name:     "no image override",
			expected: semver.Version{},
		},
		{
			name:             "effective version",
			effectiveVersion: "3.15.0",
			tag:              "nightly",
			expected:         semver.MustParse("3.15.0"),
		},
		{
			name:     "four component image tag",
			tag:      "3.15.0.0",
			expected: semver.MustParse("3.15.0"),
		},
		{
			name:        "invalid image tag",
			tag:         "nightly",
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TEST_KONG_EFFECTIVE_VERSION", tt.effectiveVersion)
			t.Setenv("TEST_KONG_TAG", tt.tag)

			actual, err := GetKongImageVersion()
			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}
