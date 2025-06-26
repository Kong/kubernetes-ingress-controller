package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/telemetry/types"
)

func TestConfigResolve(t *testing.T) {
	t.Run("Admin Token Path", func(t *testing.T) {
		validWithTokenPath := func() managercfg.Config {
			tempDir := t.TempDir()
			tokenFile, err := os.CreateTemp(tempDir, "kong.token")
			require.NoError(t, err)
			_, err = tokenFile.Write([]byte("non-empty-token"))
			require.NoError(t, err)
			return managercfg.Config{
				KongAdminTokenPath: tokenFile.Name(),
			}
		}

		t.Run("admin token path accepted", func(t *testing.T) {
			c := validWithTokenPath()
			require.NoError(t, c.Resolve())
			require.Equal(t, c.KongAdminToken, "non-empty-token")
		})
	})
}

func TestAnonymousReportsFixedPayloadCustomizer(t *testing.T) {
	cfg := managercfg.Config{}

	fixedPayload := types.Payload{
		"v":  "1.2.3",
		"kv": "3.2.1",
		"db": "db",
		"rf": "route_flavor",
		"id": "my-id",
	}

	cfg.AnonymousReportsFixedPayloadCustomizer = func(payload types.Payload) types.Payload {
		if payload == nil {
			payload = make(types.Payload)
		}
		payload["customized"] = true
		delete(payload, "v")
		return payload
	}

	result := cfg.AnonymousReportsFixedPayloadCustomizer(fixedPayload)

	require.NotNil(t, result)
	require.Equal(t, true, result["customized"])
	require.NotContains(t, result, "v")
}
