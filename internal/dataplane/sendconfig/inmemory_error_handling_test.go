package sendconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRawResourceError(t *testing.T) {
	testcases := []struct {
		name        string
		input       rawResourceError
		expected    ResourceError
		expectedErr bool
	}{
		{
			name: "KongClusterPlugin invalid schema - unknown field",
			input: rawResourceError{
				Name: "prometheus",
				ID:   "",
				Tags: []string{
					"k8s-name:default-kong",
					"k8s-kind:KongClusterPlugin",
					"k8s-uid:f9439f18-a1f8-4090-a248-3d0071c234d1",
					"k8s-group:configuration.konghq.com",
					"k8s-version:v1",
				},
				Problems: map[string]string{
					"config.config": "unknown field",
				},
			},
			expected: ResourceError{
				Name:       "default-kong",
				Kind:       "KongClusterPlugin",
				UID:        "f9439f18-a1f8-4090-a248-3d0071c234d1",
				APIVersion: "configuration.konghq.com/v1",
				Namespace:  "",
				Problems: map[string]string{
					"config.config": "unknown field",
				},
			},
		},
		{
			name: "KongPlugin invalid schema - unknown field",
			input: rawResourceError{
				Name: "prometheus",
				ID:   "",
				Tags: []string{
					"k8s-name:default-kong",
					"k8s-namespace:kong",
					"k8s-kind:KongClusterPlugin",
					"k8s-uid:f9439f18-a1f8-4090-a248-3d0071c234d1",
					"k8s-group:configuration.konghq.com",
					"k8s-version:v1",
				},
				Problems: map[string]string{
					"config.config": "unknown field",
				},
			},
			expected: ResourceError{
				Name:       "default-kong",
				Kind:       "KongClusterPlugin",
				UID:        "f9439f18-a1f8-4090-a248-3d0071c234d1",
				APIVersion: "configuration.konghq.com/v1",
				Namespace:  "kong",
				Problems: map[string]string{
					"config.config": "unknown field",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			resErr, err := parseRawResourceError(tc.input)
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, resErr)
			}
		})
	}
}
