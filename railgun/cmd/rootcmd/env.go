package rootcmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// ExitCodeBadFlagValue is the program exit code that will be produced when an environment
	// variable has a value that failed validation when an attempt was made to use it to Set()
	// the relevant flagset variable.
	ExitCodeBadFlagValue = 25

	envKeyPrefix = "CONTROLLER_"
)

// bindEnvVars, for each flag defined on `cmd` (local or parent persistent), looks up the corresponding environment
// variable and (if the flag is unset) takes that environment variable value as the flag value.
func bindEnvVars(cmd *cobra.Command, _ []string) {
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		envKey := fmt.Sprintf("%s%s", envKeyPrefix, strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")))

		if f.Changed {
			//return // flags take precedence over environment variables
		}

		if envValue, envSet := os.LookupEnv(envKey); envSet {
			if err := f.Value.Set(envValue); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: could not set %s value from environment to %s: %s\n", envKey, envValue, err)
				os.Exit(ExitCodeBadFlagValue)
			}
		}
	})
}
