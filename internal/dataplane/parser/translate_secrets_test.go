package parser

import (
	"testing"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"

	"github.com/kong/go-kong/kong"

	"github.com/stretchr/testify/require"
)

func TestGetPluginsAssociatedWithCACertSecret(t *testing.T) {
	secretID := "8a3753e0-093b-43d9-9d39-27985c987d92"
	plugins := []kongstate.Plugin{
		{
			Plugin: kong.Plugin{
				Name: kong.String("associated-plugin"),
				Config: map[string]interface{}{
					"ca_certificates": []string{secretID},
				},
			},
		},
		{
			Plugin: kong.Plugin{
				Name: kong.String("another-associated-plugin"),
				Config: map[string]interface{}{
					"ca_certificates": []string{secretID},
				},
			},
		},
		{
			Plugin: kong.Plugin{
				Name: kong.String("non-associated-plugin"),
			},
		},
	}

	associatedPlugins := getPluginsAssociatedWithCACertSecret(plugins, secretID)
	require.ElementsMatch(t, []string{"associated-plugin", "another-associated-plugin"}, associatedPlugins)
}
