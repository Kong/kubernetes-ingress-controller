// Package rootcmd implements the cobra.Command that manages the controller manager lifecycle.
package rootcmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/metadata"
)

// Execute is the entry point to the controller manager.
func Execute() {
	var (
		rootCmd    = GetRootCmd()
		versionCmd = GetVersionCmd()
	)
	rootCmd.AddCommand(versionCmd)
	cobra.CheckErr(rootCmd.Execute())
}

func GetRootCmd() *cobra.Command {
	cliCfg := config.NewCLIConfig()

	cmd := &cobra.Command{
		PersistentPreRunE: bindEnvVars,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd.Context(), *cliCfg.Config, os.Stderr)
		},
		SilenceUsage: true,
		// We can silence the errors because cobra.CheckErr below will print
		// the returned error and set the exit code to 1.
		SilenceErrors: true,
	}
	cmd.Flags().AddFlagSet(cliCfg.FlagSet())
	return cmd
}

func GetVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show JSON version information",
		RunE: func(_ *cobra.Command, _ []string) error {
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
}
