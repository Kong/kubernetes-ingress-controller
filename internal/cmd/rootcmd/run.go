package rootcmd

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

// Run starts the controller manager.
func Run(c *manager.Config) error {
	ctx, err := SetupSignalHandler(c)
	if err != nil {
		return fmt.Errorf("failed to setup signal handler: %w", err)
	}

	diag, err := StartDiagnosticsServer(ctx, manager.DiagnosticsPort, c)
	if err != nil {
		return fmt.Errorf("failed to start diagnostics server: %w", err)
	}
	return manager.Run(ctx, c, diag.ConfigDumps)
}

// RunWithLogger starts the controller manager.
// This function is intended for use in tests, where the logger can be injected.
func RunWithLogger(c *manager.Config, l logrus.FieldLogger) error {
	ctx, err := SetupSignalHandler(c)
	if err != nil {
		return fmt.Errorf("failed to setup signal handler: %w", err)
	}

	diag, err := StartDiagnosticsServer(ctx, manager.DiagnosticsPort, c)
	if err != nil {
		return fmt.Errorf("failed to start diagnostics server: %w", err)
	}
	return manager.RunWithLogger(ctx, c, diag.ConfigDumps, l)
}
