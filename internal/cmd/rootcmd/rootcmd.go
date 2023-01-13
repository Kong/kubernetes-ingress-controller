// Package rootcmd implements the cobra.Command that manages the controller manager lifecycle.
package rootcmd

import (
	"github.com/spf13/cobra"
	cobra1 "github.com/spf13/cobra"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

// Execute is the entry point to the controller manager.
func Execute() {
	var (
		cfg     manager.Config
		rootCmd = &cobra1.Command{
			PersistentPreRunE: bindEnvVars,
			RunE: func(cmd *cobra.Command, args []string) error {
				return Run(cmd.Context(), &cfg)
			},
			SilenceUsage: true,
		}
	)
	rootCmd.Flags().AddFlagSet(cfg.FlagSet())
	// cobra.CheckErr(rootCmd.Execute())
	rootCmd.Execute()
}
