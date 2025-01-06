package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestControllerManagerOptAdditionalWatchNamespace(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
		ns       string
	}{
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{"--watch-namespace=default"},
			ns:       "default",
		},
		{
			name:     "existing namespace",
			input:    []string{"--watch-namespace=default"},
			expected: []string{"--watch-namespace=default"},
			ns:       "default",
		},
		{
			name:     "new namespace",
			input:    []string{"--watch-namespace=default"},
			expected: []string{"--watch-namespace=default,testing"},
			ns:       "testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ControllerManagerOptAdditionalWatchNamespace(tt.ns)(tt.input)

			require.Equal(t, tt.expected, actual)
		})
	}
}
