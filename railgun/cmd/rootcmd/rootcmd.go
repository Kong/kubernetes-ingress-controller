package rootcmd

import (
	"context"
	"flag"

	"github.com/kong/kubernetes-ingress-controller/railgun/manager"
	"github.com/spf13/cobra"
)

var config manager.Config

func init() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	manager.RegisterFlags(&config, fs)
	rootCmd.Flags().AddGoFlagSet(fs)
}

var rootCmd = &cobra.Command{
	Use: "controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.Run(cmd.Context(), &config)
	},
	SilenceUsage: true,
}

func Execute(ctx context.Context) {
	rootCmd.ExecuteContext(ctx)
}
