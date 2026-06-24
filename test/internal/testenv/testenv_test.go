package testenv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testenv"
)

func TestIsKongGatewayEnterpriseOnly(t *testing.T) {
	tests := []struct {
		name     string
		kongTag  string
		expected bool
	}{
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
			name:     "at cutoff - four-part",
			kongTag:  "3.15.0.0",
			expected: true,
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
			t.Setenv("TEST_KONG_TAG", tc.kongTag)
			require.Equal(t, tc.expected, testenv.IsKongGatewayVersionEnterpriseOnly())
		})
	}
}
