package rootcmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmd(t *testing.T) {
	t.Run("root command succeeds by default", func(t *testing.T) {
		rootCmd := GetRootCmd()
		require.NoError(t, rootCmd.PersistentPreRunE(rootCmd, os.Args[:0]))
	})

	t.Run("root command succeeds when correct flags where provided", func(t *testing.T) {
		rootCmd := GetRootCmd()
		require.NoError(t, rootCmd.PersistentPreRunE(rootCmd,
			append(os.Args[:0],
				"--publish-service", "namespace/servicename",
			),
		))
	})

	t.Run("binding environment variables succeeds when flag validation passes", func(t *testing.T) {
		t.Setenv("CONTROLLER_PUBLISH_SERVICE", "namespace/servicename")
		rootCmd := GetRootCmd()
		require.NoError(t, rootCmd.PersistentPreRunE(rootCmd, os.Args[:0]))
	})

	t.Run("binding environment variables fails when flag validation fails", func(t *testing.T) {
		t.Setenv("CONTROLLER_PUBLISH_SERVICE", "servicename")
		rootCmd := GetRootCmd()
		require.Error(t, rootCmd.PersistentPreRunE(rootCmd, os.Args[:0]),
			"binding env vars should fail because a non namespaced name of publish service was provided",
		)
	})
}
