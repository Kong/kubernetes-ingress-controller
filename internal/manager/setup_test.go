package manager

import (
	"testing"
)

func Test_validateVersionForEndpointSlices(t *testing.T) {
	tests := []struct {
		name        string
		major       string
		minor       string
		expectError bool
	}{
		{
			name:        "ensure a valid kubernetes version for using endpoint slices does not return an error",
			major:       "1",
			minor:       "18",
			expectError: false,
		},
		{
			name:        "ensure an invalid kubernetes version for using endpoint slices returns an error",
			major:       "1",
			minor:       "16",
			expectError: true,
		},
		{
			name:        "ensure a future kubernetes version for using endpoint slices does not return an error",
			major:       "2",
			minor:       "0",
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVersionForEndpointSlices(tt.major, tt.minor)
			hasError := err != nil
			if hasError != tt.expectError {
				t.Errorf("validateVersionForEndpointSlices(%s, %s); hasError %s; expectError %v",
					tt.major, tt.minor, err, tt.expectError)
			}
		})
	}
}
