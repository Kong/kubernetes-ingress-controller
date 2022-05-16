package rootcmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bombsimon/logrusr/v2"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registers for SIGTERM and SIGINT. A context is returned
// which is canceled on one of these signals. If a second signal is not caught, the program
// will delay for the configured period of time before terminating. If a second signal is caught,
// the program is terminated with exit code 1.
func SetupSignalHandler(cfg *manager.Config) (context.Context, error) {

	deprecatedLogger, err := util.MakeLogger(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		return nil, err
	}
	logger := logrusr.New(deprecatedLogger)

	close(onlyOneSignalHandler) // panics when called twice

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		logger.Info("Signal received, shutting down", "timeout", fmt.Sprint(cfg.TermDelay))
		defer cancel()

		select {
		case <-time.After(cfg.TermDelay):
			os.Exit(0)
		case <-c:
			logger.Info("Second signal received, exiting immediately")
			os.Exit(1) // second signal. Exit directly.
		}

	}()

	return ctx, nil
}
