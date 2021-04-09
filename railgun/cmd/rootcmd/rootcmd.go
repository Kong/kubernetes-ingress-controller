package rootcmd

import (
	"flag"

	"github.com/kong/kubernetes-ingress-controller/railgun/controllers"
	"github.com/kong/kubernetes-ingress-controller/railgun/manager"
	"github.com/spf13/cobra"
)

var config manager.Config

func init() {
	registerFlags(&config)
}

func registerFlags(c *manager.Config) {
	rootCmd.Flags().StringVar(&c.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	rootCmd.Flags().StringVar(&c.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	rootCmd.Flags().BoolVar(&c.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	rootCmd.Flags().StringVar(&c.KongURL, "kong-url", "http://localhost:8001", "TODO")
	rootCmd.Flags().StringVar(&c.FilterTag, "kong-filter-tag", "managed-by-railgun", "TODO")
	rootCmd.Flags().IntVar(&c.Concurrency, "kong-concurrency", 10, "TODO")
	rootCmd.Flags().StringVar(&c.SecretName, "secret-name", "kong-config", "TODO")
	rootCmd.Flags().StringVar(&c.SecretNamespace, "secret-namespace", controllers.DefaultNamespace, "TODO")

	zapFlags := flag.NewFlagSet("", flag.ExitOnError)
	c.ZapOptions.BindFlags(zapFlags)
	rootCmd.Flags().AddGoFlagSet(zapFlags)
}

var rootCmd = &cobra.Command{
	Use: "controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.Run(&config)
	},
	SilenceUsage: true,
}

func Execute() {
	rootCmd.Execute()
}
