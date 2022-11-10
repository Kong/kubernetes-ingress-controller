package rootcmd

import (
	"context"
	"sync"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	// DiagnosticConfigBufferDepth is the size of the channel buffer for receiving diagnostic
	// config dumps from the proxy sync loop. The chosen size is essentially arbitrary: we don't
	// expect that the receive end will get backlogged (it only assigns the value to a local
	// variable) but do want a small amount of leeway to account for goroutine scheduling, so it
	// is not zero.
	DiagnosticConfigBufferDepth = 3
)

// StartDiagnosticsServer starts a goroutine that handles requests for the diagnostics server.
func StartDiagnosticsServer(
	ctx context.Context,
	port int,
	c *manager.Config,
	logger logr.Logger,
) (diagnostics.Server, error) {
	if !c.EnableProfiling && !c.EnableConfigDumps {
		logger.Info("diagnostics server disabled")
		return diagnostics.Server{}, nil
	}
	logger.Info("starting diagnostics server")

	s := diagnostics.Server{
		Logger:           logger,
		ProfilingEnabled: c.EnableProfiling,
		ConfigLock:       &sync.RWMutex{},
	}
	if c.EnableConfigDumps {
		s.ConfigDumps = util.ConfigDumpDiagnostic{
			DumpsIncludeSensitive: c.DumpSensitiveConfig,
			Configs:               make(chan util.ConfigDump, DiagnosticConfigBufferDepth),
		}
	}
	go func() {
		if err := s.Listen(ctx, port); err != nil {
			logger.Error(err, "unable to start diagnostics server")
		}
	}()
	return s, nil
}
