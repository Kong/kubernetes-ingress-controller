package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func TestIntersection(t *testing.T) {
	for _, tt := range []struct {
		name     string
		pattern  string
		value    string
		expected bool
	}{
		{
			name:     "same hostname",
			pattern:  "test.com",
			value:    "test.com",
			expected: true,
		},
		{
			name:     "different hostname suffix",
			pattern:  "test.com",
			value:    "test.net",
			expected: false,
		},
		{
			name:     "different hostname prefix",
			pattern:  "foo.com",
			value:    "bar.com",
			expected: false,
		},
		{
			name:     "different hostname lengths with only suffix matching",
			pattern:  "foo.test.com",
			value:    "test.com",
			expected: false,
		},
		{
			name:     "different hostname lengths with only prefix matching",
			pattern:  "foo.test.com",
			value:    "foo.net",
			expected: false,
		},
		{
			name:     "valid wildcard for one element",
			pattern:  "*.test.com",
			value:    "foo.test.com",
			expected: true,
		},
		{
			name:     "valid wildcard for many elements",
			pattern:  "*.example.com",
			value:    "so.many.names.example.com",
			expected: true,
		},
		{
			name:     "not matching wildcard",
			pattern:  "*.example.com",
			value:    "example.com",
			expected: false,
		},
		{
			name:     "double matching wildcard",
			pattern:  "*.example.com",
			value:    "*.example.com",
			expected: true,
		},
		{
			name:     "double not matching wildcard",
			pattern:  "*.example.com",
			value:    "*.example.net",
			expected: false,
		},
		{
			name:     "double matching wildcard with different sizes",
			pattern:  "*.test.example.com",
			value:    "*.example.com",
			expected: true,
		},
		{
			name:     "double wildcard with different sizes and not matching suffix",
			pattern:  "*.test.example.com",
			value:    "*.test.net",
			expected: false,
		},
		{
			name:     "double wildcard with different sizes and not matching prefix",
			pattern:  "*.test.example.com",
			value:    "*.example.net",
			expected: false,
		},
	} {
		// Test that the functions behave in the same way even swapping the parameters
		assert.Equal(t, tt.expected, util.HostnamesIntersect(tt.pattern, tt.value), tt.name)
		assert.Equal(t, tt.expected, util.HostnamesIntersect(tt.value, tt.pattern), tt.name)
	}
}
