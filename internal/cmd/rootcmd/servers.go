package rootcmd

import (
	"context"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
)

// StartDiagnosticsServer starts a goroutine that handles requests for the diagnostics server.
func StartDiagnosticsServer(
	ctx context.Context,
	port int,
	c *manager.Config,
	logger logr.Logger,
) (diagnostics.Server, error) {
	if !c.EnableProfiling && !c.EnableConfigDumps {
		logger.Info("Diagnostics server disabled")
		return diagnostics.Server{}, nil
	}
	logger.Info("Starting diagnostics server")

	s := diagnostics.NewServer(logger, diagnostics.ServerConfig{
		ProfilingEnabled:    c.EnableProfiling,
		ConfigDumpsEnabled:  c.EnableConfigDumps,
		DumpSensitiveConfig: c.DumpSensitiveConfig,
	})
	go func() {
		if err := s.Listen(ctx, port); err != nil {
			logger.Error(err, "Unable to start diagnostics server")
		}
	}()
	return s, nil
}
