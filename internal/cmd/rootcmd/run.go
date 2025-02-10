package rootcmd

import (
	"context"
	"fmt"
	"io"
	"os/signal"

	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/health"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/manager"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/telemetry"
)

// Run sets up a default stderr logger and starts the controller manager.
func Run(ctx context.Context, c managercfg.Config, output io.Writer) error {
	logger, err := logging.SetupLoggers(c, output)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	ctx, err = SetupSignalHandler(ctx, c, logger)
	if err != nil {
		return fmt.Errorf("failed to setup signal handler: %w", err)
	}
	defer signal.Ignore(shutdownSignals...)

	// Single manager with the same ID.
	mid := manager.NewRandomID()

	m, err := manager.NewManager(ctx, mid, logger, c)
	if err != nil {
		return fmt.Errorf("failed to create manager: %w", err)
	}

	logger.Info("Starting standalone health check server")
	health.NewHealthCheckServer(
		healthz.Ping, health.NewHealthCheckerFromFunc(m.IsReady),
	).Start(ctx, c.ProbeAddr, logger.WithName("health-check"))

	if c.AnonymousReports {
		stopAnonymousReports, err := telemetry.SetupAnonymousReports(
			ctx,
			m.GetKubeconfig(),
			m.GetClientsManager(),
			telemetry.ReportConfig{
				SplunkEndpoint:                   c.SplunkEndpoint,
				SplunkEndpointInsecureSkipVerify: c.SplunkEndpointInsecureSkipVerify,
				TelemetryPeriod:                  c.TelemetryPeriod,
				ReportValues: telemetry.ReportValues{
					PublishServiceNN:               c.PublishService.OrEmpty(),
					FeatureGates:                   c.FeatureGates,
					MeshDetection:                  len(c.WatchNamespaces) == 0,
					KonnectSyncEnabled:             c.Konnect.ConfigSynchronizationEnabled,
					GatewayServiceDiscoveryEnabled: c.KongAdminSvc.IsPresent(),
				},
			},
			mid,
		)
		if err != nil {
			logger.Error(err, "Failed setting up anonymous reports, continuing without telemetry")
		} else {
			defer stopAnonymousReports()
			logger.Info("Anonymous reports enabled")
		}
	} else {
		logger.Info("Anonymous reports disabled, skipping")
	}

	return m.Run(ctx)
}
