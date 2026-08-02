package testenv_test

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

func TestKongImageVersion(t *testing.T) {
	tests := []struct {
		name             string
		effectiveVersion string
		kongTag          string
		expected         semver.Version
		expectError      bool
	}{
		{
			name:     "no image override",
			expected: semver.Version{},
		},
		{
			name:             "effective version takes precedence over nightly tag",
			effectiveVersion: "3.15.0",
			kongTag:          "nightly",
			expected:         semver.MustParse("3.15.0"),
		},
		{
			name:     "three-part tag",
			kongTag:  "3.14.0",
			expected: semver.MustParse("3.14.0"),
		},
		{
			name:     "four-part tag",
			kongTag:  "3.15.0.0-rc.6",
			expected: semver.MustParse("3.15.0"),
		},
		{
			name:             "invalid effective version",
			effectiveVersion: "invalid",
			kongTag:          "3.15.0.0",
			expectError:      true,
		},
		{
			name:        "invalid tag",
			kongTag:     "nightly",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TEST_KONG_EFFECTIVE_VERSION", tc.effectiveVersion)
			t.Setenv("TEST_KONG_TAG", tc.kongTag)

			actual, err := testenv.KongImageVersion()
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestIsKongGatewayEnterpriseOnly(t *testing.T) {
	tests := []struct {
		name             string
		effectiveVersion string
		kongTag          string
		expected         bool
	}{
		{
			name:     "no image override",
			expected: false,
		},
		{
			name:     "below cutoff - three-part",
			kongTag:  "3.14.0",
			expected: false,
		},
		{
			name:     "below cutoff - four-part",
			kongTag:  "3.14.0.0",
			expected: false,
		},
		{
			name:             "at cutoff - effective version for nightly",
			effectiveVersion: "3.15.0",
			kongTag:          "nightly",
			expected:         true,
		},
		{
			name:     "at cutoff - four-part with rc pre-release",
			kongTag:  "3.15.0.0-rc.6",
			expected: true,
		},
		{
			name:     "above cutoff - four-part",
			kongTag:  "3.16.0.0",
			expected: true,
		},
		{
			name:     "above cutoff - higher major",
			kongTag:  "4.0.0.0",
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TEST_KONG_EFFECTIVE_VERSION", tc.effectiveVersion)
			t.Setenv("TEST_KONG_TAG", tc.kongTag)

			actual, err := testenv.IsKongGatewayVersionEnterpriseOnly()
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
