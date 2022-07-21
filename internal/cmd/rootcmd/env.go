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
//
// Why this is custom-written and not imported from spf13/viper:
// KIC 1.x used `spf13/viper` in order to provide equivalent (in fact, 100% compatible)
// behavior. In KIC 2.x we decided to drop spf13/viper in favor of a hand-written mechanism for the following reasons:
// - viper-cobra integration, per [1] and [2], requires a pflag.FlagSet.VisitAll call that binds viper values to flags
//   one by one,
// - viper is a comparably heavy dependency - it pulls in several additional features that KIC does not need (e.g.
//   configfile support),
// - the viper-cobra integration described in the 1st bullet point required more code than the whole mechanism below.
// 
// [1] https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/
// [2] https://github.com/carolynvs/stingoftheviper
func bindEnvVars(cmd *cobra.Command, _ []string) (err error) {
	var envKey string
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("environment binding failed for variable %s: %v", envKey, r)
		}
	}()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		envKey = fmt.Sprintf("%s%s", envKeyPrefix, strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")))

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
