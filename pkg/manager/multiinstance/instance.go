package multiinstance

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-logr/logr"
)

// instance represents a single manager.Manager instance in the multi-instance manager.
type instance struct {
	logger logr.Logger
	in     ManagerInstance

	stopOnce sync.Once
	stopCh   chan struct{}
}

func newInstance(in ManagerInstance, logger logr.Logger) *instance {
	return &instance{
		logger: logger.WithValues("instanceID", in.ID()),
		in:     in,
		stopCh: make(chan struct{}),
	}
}

// Stop stops the instance. Only its first call has an effect.
func (i *instance) Stop() {
	// Close stopCh only once as otherwise it would panic.
	i.stopOnce.Do(func() {
		close(i.stopCh)
	})
}

// StopChannel returns a channel one can use to wait for the instance to stop.
func (i *instance) StopChannel() <-chan struct{} {
	return i.stopCh
}

// Start runs the instance in a goroutine and blocks until the instance is stopped or the context is done.
func (i *instance) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err := i.in.Start(ctx); err != nil {
			i.logger.Error(err, "Instance exited with an error")
		}
	}()

	defer cancel() // Cancel the context once the parent context is done or the instance is stopped.
	select {
	case <-ctx.Done():
	case <-i.stopCh:
	}
}

// DiagnosticsHandler returns an HTTP handler that exposes diagnostics information for this instance. It can return
// nil if the instance does not expose diagnostics information.
func (i *instance) DiagnosticsHandler() http.Handler {
	return i.in.DiagnosticsHandler()
}

// IsReady returns an error if the instance is not ready.
func (i *instance) IsReady() error {
	return i.in.IsReady()
}
