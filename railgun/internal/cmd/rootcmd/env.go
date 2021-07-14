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
func bindEnvVars(cmd *cobra.Command, _ []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("environment binding failed: %w", r)
		}
	}()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		envKey := fmt.Sprintf("%s%s", envKeyPrefix, strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")))

		if f.Changed {
			return // flags take precedence over environment variables
		}

		if envValue, envSet := os.LookupEnv(envKey); envSet {
			if err := f.Value.Set(envValue); err != nil {
				panic(err)
			}
		}
	})

	return
}
