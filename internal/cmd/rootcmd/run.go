package rootcmd

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/internal/manager"
)

func Run(ctx context.Context, c *manager.Config) error {
	if err := StartAdmissionServer(ctx, c); err != nil {
		return fmt.Errorf("StartAdmissionServer: %w", err)
	}
	diag, err := StartDiagnosticsServer(ctx, manager.DiagnosticsPort, c)
	if err != nil {
		return fmt.Errorf("StartDiagnosticsServer: %w", err)
	}
	return manager.Run(ctx, c, diag.ConfigDumps)
}
