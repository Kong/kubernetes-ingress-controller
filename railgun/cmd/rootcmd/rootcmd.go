// Package rootcmd implements the cobra.Command that manages the controller manager lifecycle.
package rootcmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/config"
)

var cfg config.Config

func init() {
	rootCmd.Flags().AddFlagSet(cfg.FlagSet())
}

var rootCmd = &cobra.Command{
	PersistentPreRun: bindEnvVars,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Run(cmd.Context(), &cfg)
	},
	SilenceUsage: true,
}

// Execute is the entry point to the controller manager.
func Execute(ctx context.Context) {
	rootCmd.ExecuteContext(ctx)
}
