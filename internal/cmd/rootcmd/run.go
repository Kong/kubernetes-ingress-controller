package rootcmd

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

// Run sets up a default stderr logger and starts the controller manager.
func Run(ctx context.Context, c *manager.Config) error {
	deprecatedLogger, logger, err := manager.SetupLoggers(c, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	return RunWithLogger(ctx, c, deprecatedLogger, logger)
}

// RunWithLogger starts the controller manager with a provided logger.
func RunWithLogger(ctx context.Context, c *manager.Config, deprecatedLogger logrus.FieldLogger, logger logr.Logger) error {
	ctx, err := SetupSignalHandler(ctx, c, logger)
	if err != nil {
		return fmt.Errorf("failed to setup signal handler: %w", err)
	}

	if err := c.Validate(); err != nil {
		return fmt.Errorf("config invalid: %w", err)
	}

	diag, err := StartDiagnosticsServer(ctx, manager.DiagnosticsPort, c, logger)
	if err != nil {
		return fmt.Errorf("failed to start diagnostics server: %w", err)
	}

	return manager.Run(ctx, c, diag.ConfigDumps, deprecatedLogger)
}
