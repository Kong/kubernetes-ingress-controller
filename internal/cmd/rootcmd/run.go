package rootcmd

import (
	"context"
	"fmt"
	"io"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

// Run sets up a default stderr logger and starts the controller manager.
func Run(ctx context.Context, c *manager.Config, output io.Writer) error {
	logger, err := manager.SetupLoggers(c, output)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	ctx = ctrl.LoggerInto(ctx, logger)

	return RunWithLogger(ctx, c, logger)
}

// RunWithLogger starts the controller manager with a provided logger.
func RunWithLogger(ctx context.Context, c *manager.Config, logger logr.Logger) error {
	ctx, err := SetupSignalHandler(ctx, c, logger)
	if err != nil {
		return fmt.Errorf("failed to setup signal handler: %w", err)
	}

	if err := c.Validate(); err != nil {
		return fmt.Errorf("config invalid: %w", err)
	}

	diag, err := StartDiagnosticsServer(ctx, c.DiagnosticServerPort, c, logger)
	if err != nil {
		return fmt.Errorf("failed to start diagnostics server: %w", err)
	}

	return manager.Run(ctx, c, diag.ConfigDumps, logger)
}
