// Package rootcmd implements the cobra.Command that manages the controller manager lifecycle.
package rootcmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager/metadata"
)

// Execute is the entry point to the controller manager.
func Execute() {
	var (
		cfg     manager.Config
		rootCmd = &cobra.Command{
			PersistentPreRunE: bindEnvVars,
			RunE: func(cmd *cobra.Command, args []string) error {
				return Run(cmd.Context(), &cfg)
			},
			SilenceUsage: true,
		}
		versionCmd = &cobra.Command{
			Use:   "version",
			Short: "Show JSON version information",
			RunE: func(cmd *cobra.Command, args []string) error {
				type Version struct {
					Release string `json:"release"`
					Repo    string `json:"repo"`
					Commit  string `json:"commit"`
				}
				out, err := json.Marshal(Version{
					Release: metadata.Release,
					Repo:    metadata.Repo,
					Commit:  metadata.Commit,
				})
				if err != nil {
					return fmt.Errorf("failed to print version information: %w", err)
				}
				fmt.Printf("%s\n", out)
				return nil
			},
		}
	)
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().AddFlagSet(cfg.FlagSet())
	cobra.CheckErr(rootCmd.Execute())
}
