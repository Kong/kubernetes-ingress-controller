package versions_test

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
)

func TestDeckFileFormat(t *testing.T) {
	t.Run("Kong version >= 3.0", func(t *testing.T) {
		actualDeckFormat := versions.DeckFileFormat(semver.Version{Major: 3, Minor: 0})
		require.Equal(t, "3.0", actualDeckFormat)
	})

	t.Run("Kong version >= 3.0 - other than 3.0", func(t *testing.T) {
		actualDeckFormat := versions.DeckFileFormat(semver.Version{Major: 3, Minor: 5})
		require.Equal(t, "3.0", actualDeckFormat)
	})

	t.Run("Kong version < 3.0", func(t *testing.T) {
		actualDeckFormat := versions.DeckFileFormat(semver.Version{Major: 2, Minor: 9})
		require.Equal(t, "1.1", actualDeckFormat)
	})
}
