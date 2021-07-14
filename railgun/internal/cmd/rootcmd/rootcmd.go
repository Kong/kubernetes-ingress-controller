// Package rootcmd implements the cobra.Command that manages the controller manager lifecycle.
package rootcmd

import (
	"context"

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/manager"
	"github.com/spf13/cobra"
)

var cfg manager.Config

func init() {
	rootCmd.Flags().AddFlagSet(cfg.FlagSet())
}

var rootCmd = &cobra.Command{
	PersistentPreRunE: bindEnvVars,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Run(cmd.Context(), &cfg)
	},
	SilenceUsage: true,
}

// Execute is the entry point to the controller manager.
func Execute(ctx context.Context) {
	rootCmd.ExecuteContext(ctx)
}
