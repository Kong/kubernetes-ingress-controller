// Package rootcmd implements the cobra.Command that manages the controller manager lifecycle.
package rootcmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/kong/kubernetes-ingress-controller/railgun/manager"
)

var config manager.Config

func init() {
	rootCmd.Flags().AddFlagSet(config.FlagSet())
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.Run(cmd.Context(), &config)
	},
	SilenceUsage: true,
}

// Execute is the entry point to the controller manager.
func Execute(ctx context.Context) {
	rootCmd.ExecuteContext(ctx)
}
