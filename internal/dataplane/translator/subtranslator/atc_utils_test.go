package subtranslator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHostMatcherFromHosts(t *testing.T) {
	testCases := []struct {
		name       string
		hosts      []string
		expression string
	}{
		{
			name:       "simple exact host",
			hosts:      []string{"a.example.com"},
			expression: `http.host == "a.example.com"`,
		},
		{
			name:       "single wildcard host",
			hosts:      []string{"*.example.com"},
			expression: `http.host =^ ".example.com"`,
		},
		{
			name:       "multiple hosts with mixture of exact and wildcard",
			hosts:      []string{"foo.com", "*.bar.com"},
			expression: `(http.host == "foo.com") || (http.host =^ ".bar.com")`,
		},
		{
			name:       "multiple hosts including invalid host",
			hosts:      []string{"foo.com", "a..bar.com"},
			expression: `http.host == "foo.com"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matcher := hostMatcherFromHosts(tc.hosts)
			require.Equal(t, tc.expression, matcher.Expression())
		})
	}
}
