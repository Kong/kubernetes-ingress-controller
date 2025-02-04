package rootcmd

import (
	"context"
	"fmt"
	"io"
	"os/signal"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

// Run sets up a default stderr logger and starts the controller manager.
func Run(ctx context.Context, c managercfg.Config, output io.Writer) error {
	logger, err := manager.SetupLoggers(c, output)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	ctx, err = SetupSignalHandler(ctx, c, logger)
	if err != nil {
		return fmt.Errorf("failed to setup signal handler: %w", err)
	}
	defer signal.Ignore(shutdownSignals...)

	return manager.Run(ctx, c, logger)
}
