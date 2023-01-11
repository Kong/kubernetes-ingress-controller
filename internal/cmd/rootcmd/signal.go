package rootcmd

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
)

var (
	mutex           sync.Mutex
	shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

// SetupSignalHandler registers for SIGTERM and SIGINT. A context is returned
// which is canceled on one of these signals. If a second signal is not caught, the program
// will delay for the configured period of time before terminating. If a second signal is caught,
// the program is terminated with exit code 1.
func SetupSignalHandler(ctx context.Context, cfg *manager.Config, logger logr.Logger) (context.Context, error) {
	// This will prevent multiple signal handlers from being created
	if ok := mutex.TryLock(); !ok {
		return nil, errors.New("signal handler can only be setup once")
	}

	var cancel func()
	ctx, cancel = context.WithCancel(ctx)

	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		sig := <-c
		logger.Info("Signal received, shutting down", "graceful_period", cfg.TermDelay.String(), "signal", sig.String())
		cancel()

		// If code in other places has already exited then code below won't
		// execute, and hence the os.Exit() will not be called.
		// This allows deferred code in other parts of the application to execute.
		if cfg.TermDelay != 0 {
			select {
			case <-time.After(cfg.TermDelay):
				logger.Info("Graceful termination period has passed, exiting immediately", "graceful_period", cfg.TermDelay.String())
			case sig := <-c:
				logger.Info("Signal received during graceful shutdown, exiting immediately", "signal", sig.String())
			}
		} else {
			sig := <-c
			logger.Info("Signal received during graceful shutdown, exiting immediately", "signal", sig.String())
		}

		os.Exit(1)
	}()

	return ctx, nil
}
