package rootcmd

import (
	"fmt"
	"io"
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
		Run: func(cmd *cobra.Command, _ []string) {
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

	t.Setenv("CONTROLLER_FLAG_2", "env2")
	t.Setenv("CONTROLLER_FLAG_4", "env4")

	cmd.SetArgs([]string{
		"--flag-3=args3",
		"--flag-4=args4",
	})
	err := cmd.Execute()

	assert.True(t, commandHasRun)
	assert.NoError(t, err)
}

func TestBindEnvVarsSlice(t *testing.T) {
	t.Run("set by flags", func(t *testing.T) {
		cmd := &cobra.Command{
			PreRunE: bindEnvVars,
			Run:     func(_ *cobra.Command, _ []string) {},
		}

		ss := cmd.Flags().StringSlice("flag-string-slice", []string{"default"}, "No description")

		t.Setenv("CONTROLLER_FLAG_STRING_SLICE", "q,w,e,r,t,y")

		cmd.SetArgs([]string{
			"--flag-string-slice=1",
			"--flag-string-slice=2",
			"--flag-string-slice=3",
		})
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, []string{"1", "2", "3"}, *ss)
	})

	t.Run("set by env", func(t *testing.T) {
		cmd := &cobra.Command{
			PreRunE: bindEnvVars,
			Run:     func(_ *cobra.Command, _ []string) {},
		}

		ss := cmd.Flags().StringSlice("flag-string-slice", []string{"default"}, "No description")

		t.Setenv("CONTROLLER_FLAG_STRING_SLICE", "q,w,e,r,t,y")

		cmd.SetArgs([]string{})
		assert.NoError(t, cmd.Execute())
		assert.Equal(t, []string{"q", "w", "e", "r", "t", "y"}, *ss)
	})
}

func TestBindEnvVarsValidation(t *testing.T) {
	cmd := &cobra.Command{
		PreRunE: bindEnvVars,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Flags().Parse(nil)
		},
	}

	cmd.Flags().Var(testvar("validation_test"), "validation_test", "")
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	t.Setenv("CONTROLLER_VALIDATION_TEST", "intentionally_fail")

	err := cmd.Execute()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "variable CONTROLLER_VALIDATION_TEST: bad value for var type testvar"))
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
	return nil
}

func (t testvar) Type() string {
	return "testvar"
}
