package versions_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
)

func TestDeckFileFormat(t *testing.T) {
	t.Run("Kong Getaway version supported by KIC >=3.0.0 should has 3.0 config format", func(t *testing.T) {
		require.Equal(t, "3.0", versions.DeckFileFormat(versions.KICv3VersionCutoff))
	})
}
