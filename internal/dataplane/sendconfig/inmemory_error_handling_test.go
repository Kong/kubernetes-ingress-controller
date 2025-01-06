package sendconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRawResourceError(t *testing.T) {
	testcases := []struct {
		name                   string
		input                  rawResourceError
		expected               ResourceError
		expectedErrMsgContains string
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
		{
			name: "name tag is missing",
			input: rawResourceError{
				Name: "",
				ID:   "0450730c-9e16-4650-bde6-24e7afc83926",
				Tags: []string{},
				Problems: map[string]string{
					"certificate:": "required field missing",
				},
			},
			expectedErrMsgContains: "resource error has no name tag",
		},
		{
			name: "namespace tag is missing",
			input: rawResourceError{
				Name: "secret1",
				ID:   "f9439f18-a1f8-4090-a248-3d0071c234d1",
				Tags: []string{
					"k8s-name:secret1",
					"k8s-kind:Secret",
					"k8s-uid:f9439f18-a1f8-4090-a248-3d0071c234d1",
					"k8s-group:configuration.konghq.com",
				},
				Problems: map[string]string{
					"certificate:": "required field missing",
				},
			},
			expectedErrMsgContains: "resource error has no namespace tag, name: secret1",
		},
		{
			name: "namespace tag is missing",
			input: rawResourceError{
				Name: "secret1",
				ID:   "f9439f18-a1f8-4090-a248-3d0071c234d1",
				Tags: []string{
					"k8s-name:secret1",
					"k8s-namespace:default",
				},
				Problems: map[string]string{
					"certificate:": "required field missing",
				},
			},
			expectedErrMsgContains: "resource error has not enough kind, group, version tags, name: secret1",
		},
		{
			name: "namespace tag is missing",
			input: rawResourceError{
				Name: "secret1",
				ID:   "f9439f18-a1f8-4090-a248-3d0071c234d1",
				Tags: []string{
					"k8s-name:secret1",
					"k8s-namespace:default",
					"k8s-kind:Secret",
					"k8s-group:configuration.konghq.com",
				},
				Problems: map[string]string{
					"certificate:": "required field missing",
				},
			},
			expectedErrMsgContains: "resource error has no uid tag, name: secret1",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			resErr, err := parseRawResourceError(tc.input)
			if tc.expectedErrMsgContains != "" {
				require.ErrorContains(t, err, tc.expectedErrMsgContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, resErr)
			}
		})
	}
}
