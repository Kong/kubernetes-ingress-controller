package rootcmd

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

// Run sets up a default stderr logger and starts the controller manager.
func Run(c *manager.Config) error {
	deprecatedLogger, logger, err := manager.SetupLoggers(c, os.Stderr)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	ctx, err := SetupSignalHandler(c, logger)
	if err != nil {
		return fmt.Errorf("failed to setup signal handler: %w", err)
	}

	diag, err := StartDiagnosticsServer(ctx, manager.DiagnosticsPort, c, logger)
	if err != nil {
		return fmt.Errorf("failed to start diagnostics server: %w", err)
	}
	return manager.Run(ctx, c, diag.ConfigDumps, deprecatedLogger)
}

// RunWithLogger starts the controller manager with a provided logger.
func RunWithLogger(c *manager.Config, deprecatedLogger logrus.FieldLogger, logger logr.Logger) error {
	ctx, err := SetupSignalHandler(c, logger)
	if err != nil {
		return fmt.Errorf("failed to setup signal handler: %w", err)
	}

	diag, err := StartDiagnosticsServer(ctx, manager.DiagnosticsPort, c, logger)
	if err != nil {
		return fmt.Errorf("failed to start diagnostics server: %w", err)
	}
	return manager.Run(ctx, c, diag.ConfigDumps, deprecatedLogger)
}
