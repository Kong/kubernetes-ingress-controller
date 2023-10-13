package helpers

import (
	"testing"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

func FreePort(t *testing.T) int {
	port, err := freeport.GetFreePort()
	require.NoError(t, err)
	return port
}
