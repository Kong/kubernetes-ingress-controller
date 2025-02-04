package rootcmd

import (
	"context"
	"fmt"
	"io"
	"os/signal"

	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/health"
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

	mid, err := manager.NewID("kic")
	if err != nil {
		return fmt.Errorf("failed to create manager ID: %w", err)
	}
	m, err := manager.NewManager(ctx, mid, logger, c)
	if err != nil {
		return fmt.Errorf("failed to create manager: %w", err)
	}

	logger.Info("Starting standalone health check server")
	health.NewHealthCheckServer(
		healthz.Ping, health.NewHealthCheckerFromFunc(m.IsReady),
	).Start(ctx, c.ProbeAddr, logger.WithName("health-check"))

	return m.Run(ctx)
}
