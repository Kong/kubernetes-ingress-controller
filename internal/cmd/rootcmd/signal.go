package rootcmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registers for SIGTERM and SIGINT. A context is returned
// which is canceled on one of these signals. If a second signal is not caught, the program
// will delay for the configured period of time before terminating. If a second signal is caught,
// the program is terminated with exit code 1.
func SetupSignalHandler(delay time.Duration) context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		cancel()
		select {
		case <-time.After(delay):
		case <-c:
		}
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}
