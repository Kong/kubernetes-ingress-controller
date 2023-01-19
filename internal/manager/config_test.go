package manager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	t.Run("--gateway-api-controller-name", func(t *testing.T) {
		t.Run("valid config", func(t *testing.T) {
			var c Config
			require.NoError(t, c.FlagSet().Parse(
				[]string{
					os.Args[0],
					`--gateway-api-controller-name`, `example.com/controller-name`,
				},
			))
			require.NoError(t, c.Validate())
		})

		t.Run("invalid config", func(t *testing.T) {
			var c Config
			require.NoError(t, c.FlagSet().Parse(
				[]string{
					os.Args[0],
					`--gateway-api-controller-name`, `%invalid_controller_name$`,
				},
			))
			require.Error(t, c.Validate())
		})
	})

	t.Run("--publish-service", func(t *testing.T) {
		t.Run("valid config", func(t *testing.T) {
			var c Config
			require.NoError(t, c.FlagSet().Parse(
				[]string{
					os.Args[0],
					`--publish-service`, `namespace/servicename`,
				},
			))
			require.NoError(t, c.Validate())
		})

		t.Run("invalid config", func(t *testing.T) {
			var c Config
			require.Error(t, c.FlagSet().Parse(
				[]string{
					os.Args[0],
					`--publish-service`, `servicename`,
				},
			))
			// publish service is validated through FlagNamespacedName validation logic.
			require.NoError(t, c.Validate())
		})
	})
}
