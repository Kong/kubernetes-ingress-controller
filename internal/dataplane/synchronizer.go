package dataplane

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bombsimon/logrusr/v2"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

// -----------------------------------------------------------------------------
// Dataplane Synchronizer - Public Vars
// -----------------------------------------------------------------------------

const (
	// DefaultSyncSeconds indicates the time.Duration (minimum) that will occur between
	// updates to the DataplaneClient.
	//
	// This default was based on local testing wherein it appeared sub-second updates
	// to the Admin API could be problematic (or at least operate differently) based on
	// which storage backend was in use (i.e. "dbless", "postgres"). This is a workaround
	// for improvements we still need to investigate upstream.
	//
	// See Also: https://github.com/Kong/kubernetes-ingress-controller/issues/1398
	DefaultSyncSeconds float32 = 3.0

	DefaultInitWaitPeriod = 5 * time.Second
)

// -----------------------------------------------------------------------------
// Synchronizer - Public Types
// -----------------------------------------------------------------------------

// Synchronizer is a threadsafe object which starts a goroutine to updates
// the data-plane at regular intervals.
type Synchronizer struct {
	logger logr.Logger

	// dataplane client to send updates to the Kong Admin API
	dataplaneClient Client
	dbMode          string

	// server configuration, flow control, channels and utility attributes
	stagger         time.Duration
	syncTicker      *time.Ticker
	configApplied   bool
	isServerRunning bool
	initWaitPeriod  time.Duration

	lock sync.RWMutex
}

type SynchronizerOption func(*Synchronizer)

// WithStagger returns a SynchronizerOption which sets the stagger period.
func WithStagger(period time.Duration) SynchronizerOption {
	return func(s *Synchronizer) {
		s.stagger = period
	}
}

// WithInitWaitPeriod returns a SynchronizerOption which sets the initial wait period.
func WithInitWaitPeriod(period time.Duration) SynchronizerOption {
	return func(s *Synchronizer) {
		s.initWaitPeriod = period
	}
}

// NewSynchronizer will provide a new Synchronizer object with a specified
// stagger time for data-plane updates to occur. Note that this starts some
// background goroutines and the caller is resonsible for marking the provided
// context.Context as "Done()" to shut down the background routines.
func NewSynchronizer(logger logrus.FieldLogger, client Client, opts ...SynchronizerOption) (*Synchronizer, error) {
	synchronizer := &Synchronizer{
		logger:          logrusr.New(logger),
		stagger:         time.Duration(DefaultSyncSeconds),
		initWaitPeriod:  DefaultInitWaitPeriod,
		dataplaneClient: client,
		configApplied:   false,
		dbMode:          client.DBMode(),
	}

	for _, opt := range opts {
		opt(synchronizer)
	}

	synchronizer.dbMode = client.DBMode()

	return synchronizer, nil
}

// -----------------------------------------------------------------------------
// Synchronizer - Public Methods
// -----------------------------------------------------------------------------

// Start starts the goroutine synchronization server that will perform an
// Update() on the provided dataplane.Client according to the provided stagger
// time, or using the DefaultSyncSeconds if not otherwise provided.
//
// To stop the server, the provided context must be Done().
func (p *Synchronizer) Start(ctx context.Context) error {
	select {
	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2249
	// This is a temporary mitigation to allow some time for controllers to
	// populate their dataplaneClient cache.
	case <-time.After(p.initWaitPeriod):
	case <-ctx.Done():
		return fmt.Errorf("Synchronizer Start() interrupted: %w", ctx.Err())
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	if p.isServerRunning {
		return fmt.Errorf("server is already running")
	}

	p.syncTicker = time.NewTicker(p.stagger)
	go p.startUpdateServer(ctx)
	p.isServerRunning = true

	return nil
}

// IsRunning informs the caller whether the synchronization server is running.
func (p *Synchronizer) IsRunning() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.isServerRunning
}

// IsReady indicates whether the synchronizer is actively able to synchronize
// configuration to the dataplane. It's similar to IsRunning() but reports
// on whether configuration can actually be successful and is also used as part
// of a controller-runtime Runnable interface to wait for readiness before
// starting controllers.
func (p *Synchronizer) IsReady() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	// If the proxy is has no database, it is only ready after a successful sync
	// Otherwise, it has no configuration loaded
	if p.dbMode == "off" {
		return p.configApplied
	}
	// If the proxy has a database, it is ready immediately
	// It will load existing configuration from the database
	return true
}

// NeedLeaderElection implements the controller-runtime Runnable interface to
// inform the controller manager whether leadership election is needed, which
// is always true in our case.
func (p *Synchronizer) NeedLeaderElection() bool {
	return true
}

// -----------------------------------------------------------------------------
// Synchronizer - Private Methods - Server Utilities
// -----------------------------------------------------------------------------

// startUpdateServer runs a server in a background goroutine that is responsible for
// updating the kong proxy backend at regular intervals.
func (p *Synchronizer) startUpdateServer(ctx context.Context) {
	var initialConfig sync.Once
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("context done: shutting down the proxy update server")
			if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
				p.logger.Error(err, "context completed with error")
			}
			p.syncTicker.Stop()

			p.lock.Lock()
			defer p.lock.Unlock()
			p.isServerRunning = false
			p.configApplied = false

			return

		case <-p.syncTicker.C:
			if err := p.dataplaneClient.Update(ctx); err != nil {
				p.logger.Error(err, "could not update kong admin")
				continue
			}
			initialConfig.Do(p.markConfigApplied)
		}
	}
}

// -----------------------------------------------------------------------------
// Synchronizer - Private Methods - Helper
// -----------------------------------------------------------------------------

// markConfigApplied marks that config has been applied.
func (p *Synchronizer) markConfigApplied() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.configApplied = true
}
