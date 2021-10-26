package version

import (
	"fmt"

	"github.com/spf13/cobra"

	kic_version "github.com/kong/kubernetes-ingress-controller/v2/pkg/version"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Long:  `Print version.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			buildInfo := kic_version.Build

			cmd.Println(fmt.Sprintf("Version:    %s", buildInfo.Version))
			cmd.Println(fmt.Sprintf("Git Tag:    %s", buildInfo.GitTag))
			cmd.Println(fmt.Sprintf("Git Commit: %s", buildInfo.GitCommit))
			cmd.Println(fmt.Sprintf("Build Date: %s", buildInfo.BuildDate))

			return nil
		},
	}
	return cmd
}
