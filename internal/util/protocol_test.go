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
		{"http", true},
		{"https", true},
		{"grpc", true},
		{"grpcs", true},
		{"grcpsfdsafdsfafdshttp", false},
	}
	for _, testcase := range testTable {
		isMatch := ValidateProtocol(testcase.input)
		assert.Equal(isMatch, testcase.result)
	}
}
