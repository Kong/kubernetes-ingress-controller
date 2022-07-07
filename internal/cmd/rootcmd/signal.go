package rootcmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bombsimon/logrusr/v2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var (
	mutex           sync.Mutex
	shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

// SetupSignalHandler registers for SIGTERM and SIGINT. A context is returned
// which is canceled on one of these signals. If a second signal is not caught, the program
// will delay for the configured period of time before terminating. If a second signal is caught,
// the program is terminated with exit code 1.
func SetupSignalHandler(cfg *manager.Config) (context.Context, error) {
	// This will prevent multiple signal handlers from being created
	if ok := mutex.TryLock(); !ok {
		return nil, errors.New("signal handler can only be setup once")
	}

	deprecatedLogger, err := util.MakeLogger(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		return nil, err
	}
	logger := logrusr.New(deprecatedLogger)

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		logger.Info("Signal received, shutting down", "timeout", fmt.Sprint(cfg.TermDelay))

		select {
		case <-time.After(cfg.TermDelay):
			cancel()
		case <-c:
			logger.Info("Signal received during termination delay, exiting immediately")
			os.Exit(1) // second signal. Exit directly.
		}

		<-c
		logger.Info("Signal received during graceful shutdown, exiting immediately")
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx, nil
}
