package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPort(t *testing.T) {
	assert.True(t, IsValidPort(1))
	assert.True(t, IsValidPort(80))
	assert.True(t, IsValidPort(8080))
	assert.True(t, IsValidPort(65535))
	assert.False(t, IsValidPort(0)) // 0 is a reserved port
	assert.False(t, IsValidPort(65536))
	assert.False(t, IsValidPort(9999999))
}
