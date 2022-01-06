package rootcmd

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

func Run(ctx context.Context, c *manager.Config) error {
	diag, err := StartDiagnosticsServer(ctx, manager.DiagnosticsPort, c)
	if err != nil {
		return fmt.Errorf("failed to start diagnostics server: %w", err)
	}
	return manager.Run(ctx, c, diag.ConfigDumps)
}
