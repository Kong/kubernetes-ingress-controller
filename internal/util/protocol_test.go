package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateProtocol(t *testing.T) {
	testTable := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"http", true},
		{"https", true},
		{"grpc", true},
		{"grpcs", true},
		{"ws", true},
		{"wss", true},
		{"tls", true},
		{"tcp", true},
		{"tls_passthrough", true},
		{"grcpsfdsafdsfafdshttp", false},
	}
	for _, tc := range testTable {
		t.Run(tc.input, func(t *testing.T) {
			isMatch := ValidateProtocol(tc.input)
			assert.Equal(t, tc.expected, isMatch)
		})
	}
}

func BenchmarkValidateProtocol(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ValidateProtocol("https")
		_ = ValidateProtocol("tcp")
		_ = ValidateProtocol("tls")
		_ = ValidateProtocol("xxxxxxxxx")
	}
}
