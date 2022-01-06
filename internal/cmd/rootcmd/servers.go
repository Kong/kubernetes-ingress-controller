package rootcmd

import (
	"context"
	"sync"

	"github.com/bombsimon/logrusr"

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

func StartDiagnosticsServer(ctx context.Context, port int, c *manager.Config) (diagnostics.Server, error) {
	deprecatedLogger, err := util.MakeLogger(c.LogLevel, c.LogFormat)
	if err != nil {
		return diagnostics.Server{}, err
	}
	logger := logrusr.NewLogger(deprecatedLogger)

	if !c.EnableProfiling && !c.EnableConfigDumps {
		logger.Info("diagnostics server disabled")
		return diagnostics.Server{}, nil
	}

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
