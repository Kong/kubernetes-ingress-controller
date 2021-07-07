package rootcmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const envKeyPrefix = "CONTROLLER_"

// bindEnvVars, for each flag defined on `cmd` (local or parent persistent), looks up the corresponding environment
// variable and (if the flag is unset) takes that environment variable value as the flag value.
func bindEnvVars(cmd *cobra.Command, _ []string) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		envKey := fmt.Sprintf("%s%s", envKeyPrefix, strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")))

		if f.Changed {
			return // flags take precedence over environment variables
		}

		if envValue, envSet := os.LookupEnv(envKey); envSet {
			// for convenience, any EnablementStatus type variable we need to translate "false" into "disabled"
			if f.Value.Type() == "EnablementStatus" {
				if envValue == "false" {
					envValue = "disabled"
				}
				if envValue == "true" {
					envValue = "enabled"
				}
			}
			cmd.Flags().Set(f.Name, envValue)
		}
	})
}
