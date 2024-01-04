package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateProtocol(t *testing.T) {
	assert := assert.New(t)
	testTable := []struct {
		input  string
		result bool
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
	for _, testcase := range testTable {
		isMatch := ValidateProtocol(testcase.input)
		assert.Equal(isMatch, testcase.result)
	}
}
