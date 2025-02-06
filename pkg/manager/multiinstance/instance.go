package multiinstance

import (
	"context"
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

// Run runs the instance in a goroutine and blocks until the instance is stopped or the context is done.
func (i *instance) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err := i.in.Run(ctx); err != nil {
			i.logger.Error(err, "Instance exited with an error")
		}
	}()

	defer cancel() // Cancel the context once the parent context is done or the instance is stopped.
	select {
	case <-ctx.Done():
	case <-i.stopCh:
	}
}

func (i *instance) IsReady() error {
	return i.in.IsReady()
}
