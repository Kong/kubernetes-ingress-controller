package rootcmd

import (
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

func Run(c *manager.Config) error {
	ctx := SetupSignalHandler(c.TermDelay)

	diag, err := StartDiagnosticsServer(ctx, manager.DiagnosticsPort, c)
	if err != nil {
		return fmt.Errorf("failed to start diagnostics server: %w", err)
	}
	return manager.Run(ctx, c, diag.ConfigDumps)
}
