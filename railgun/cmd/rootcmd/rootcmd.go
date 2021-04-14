// Package rootcmd implements the cobra.Command that manages the controller manager lifecycle.
package rootcmd

import (
	"context"
	"flag"

	"github.com/kong/kubernetes-ingress-controller/railgun/manager"
	"github.com/spf13/cobra"
)

var config manager.Config

func bindFlags(cmd *cobra.Command, args []string) {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	manager.RegisterFlags(&config, fs)
	cmd.Flags().AddGoFlagSet(fs)
}

var rootCmd = &cobra.Command{
	Use:    "controller",
	PreRun: bindFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.Run(cmd.Context(), &config)
	},
	SilenceUsage: true,
}

// Execute is the entry point to the controller manager.
func Execute(ctx context.Context) {
	rootCmd.ExecuteContext(ctx)
}
