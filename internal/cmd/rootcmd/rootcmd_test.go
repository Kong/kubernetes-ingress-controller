package rootcmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

func TestRootCmd(t *testing.T) {
	t.Run("root command succeeds by default", func(t *testing.T) {
		var cfg manager.Config
		rootCmd := GetRootCmd(&cfg)
		require.NoError(t, rootCmd.PersistentPreRunE(rootCmd, os.Args[:0]))
	})

	t.Run("root command succeeds when correct flags where provided", func(t *testing.T) {
		var cfg manager.Config
		rootCmd := GetRootCmd(&cfg)
		require.NoError(t, rootCmd.PersistentPreRunE(rootCmd,
			append(os.Args[:0],
				"--publish-service", "namespace/servicename",
			),
		))
	})

	t.Run("binding environment variables succeeds when flag validation passes", func(t *testing.T) {
		t.Setenv("CONTROLLER_PUBLISH_SERVICE", "namespace/servicename")
		var cfg manager.Config
		rootCmd := GetRootCmd(&cfg)
		require.NoError(t, rootCmd.PersistentPreRunE(rootCmd, os.Args[:0]))
	})

	t.Run("binding environment variables fails when flag validation fails", func(t *testing.T) {
		t.Setenv("CONTROLLER_PUBLISH_SERVICE", "servicename")
		var cfg manager.Config
		rootCmd := GetRootCmd(&cfg)
		require.Error(t, rootCmd.PersistentPreRunE(rootCmd, os.Args[:0]),
			"binding env vars should fail because a non namespaced name of publish service was provided",
		)
	})
}
