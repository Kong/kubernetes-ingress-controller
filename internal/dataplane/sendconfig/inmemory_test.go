package sendconfig

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/require"
)

func Test_shouldUseFlattenedErrors(t *testing.T) {
	testcases := []struct {
		version  string
		expected bool
	}{
		{
			version:  "3.1.0",
			expected: false,
		},
		{
			version:  "3.1.9",
			expected: false,
		},
		{
			version:  "3.2.0",
			expected: true,
		},
		{
			version:  "3.2.1",
			expected: true,
		},
		{
			version:  "3.5.0",
			expected: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.version, func(t *testing.T) {
			v, err := semver.Parse(tc.version)
			require.NoError(t, err)
			if tc.expected {
				require.True(t, shouldUseFlattenedErrors(v))
			} else {
				require.False(t, shouldUseFlattenedErrors(v))
			}
		})
	}
}
