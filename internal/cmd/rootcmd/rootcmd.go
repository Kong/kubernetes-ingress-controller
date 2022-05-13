// Package rootcmd implements the cobra.Command that manages the controller manager lifecycle.
package rootcmd

import (
	"github.com/spf13/cobra"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

var cfg manager.Config

func init() {
	rootCmd.Flags().AddFlagSet(cfg.FlagSet())
}

var rootCmd = &cobra.Command{
	PersistentPreRunE: bindEnvVars,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Run(&cfg)
	},
	SilenceUsage: true,
}

// Execute is the entry point to the controller manager.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
