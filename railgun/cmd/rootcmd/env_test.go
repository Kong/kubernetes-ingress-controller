package rootcmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBindEnvVars(t *testing.T) {
	commandHasRun := false
	cmd := &cobra.Command{
		PreRun: bindEnvVars,
		Run: func(cmd *cobra.Command, args []string) {
			got1, _ := cmd.Flags().GetString("flag-1")
			got2, _ := cmd.Flags().GetString("flag-2")
			got3, _ := cmd.Flags().GetString("flag-3")
			got4, _ := cmd.Flags().GetString("flag-4")
			require.Equal(t, "default1", got1) // env not set, arg not set
			require.Equal(t, "env2", got2)     // env set, arg not set
			require.Equal(t, "args3", got3)    // env not set, arg set
			require.Equal(t, "args4", got4)    // env set, arg set
			commandHasRun = true
		},
	}

	cmd.Flags().String("flag-1", "default1", "Not set")
	cmd.Flags().String("flag-2", "default2", "Set by env only")
	cmd.Flags().String("flag-3", "default3", "Set by args only")
	cmd.Flags().String("flag-4", "default4", "Set by both env and args")

	os.Setenv("CONTROLLER_FLAG_2", "env2")
	os.Setenv("CONTROLLER_FLAG_4", "env4")

	cmd.SetArgs([]string{
		"--flag-3=args3",
		"--flag-4=args4",
	})
	err := cmd.Execute()

	assert.True(t, commandHasRun)
	assert.NoError(t, err)
}
