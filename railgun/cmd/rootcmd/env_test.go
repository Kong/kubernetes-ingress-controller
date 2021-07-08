package rootcmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBindEnvVars(t *testing.T) {
	commandHasRun := false
	cmd := &cobra.Command{
		PreRunE: bindEnvVars,
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

	_ = os.Setenv("CONTROLLER_FLAG_2", "env2")
	_ = os.Setenv("CONTROLLER_FLAG_4", "env4")
	defer func() {
		_ = os.Unsetenv("CONTROLLER_FLAG_2")
		_ = os.Unsetenv("CONTROLLER_FLAG_4")
	}()

	cmd.SetArgs([]string{
		"--flag-3=args3",
		"--flag-4=args4",
	})
	err := cmd.Execute()

	assert.True(t, commandHasRun)
	assert.NoError(t, err)
}

func TestBindEnvVarsValidation(t *testing.T) {
	cmd := &cobra.Command{
		PreRunE: bindEnvVars,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Flags().Parse(nil)
		},
	}

	cmd.Flags().Var(testvar("validation_test"), "validation_test", "")
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	_ = os.Setenv("CONTROLLER_VALIDATION_TEST", "intentionally_fail")
	defer os.Unsetenv("CONTROLLER_VALIDATION_TEST")

	err := cmd.Execute()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "bad value for var type testvar"))
}

// -----------------------------------------------------------------------------
// Test Utilities
// -----------------------------------------------------------------------------

type testvar string

func (t testvar) String() string {
	return string(t)
}

func (t testvar) Set(newstr string) error {
	if newstr == "intentionally_fail" {
		return fmt.Errorf("bad value for var type %s", t.Type())
	}
	t = testvar(newstr)
	return nil
}

func (t testvar) Type() string {
	return "testvar"
}
