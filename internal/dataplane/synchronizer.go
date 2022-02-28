package dataplane

import (
	"context"
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
	// This 1s default was based on local testing wherein it appeared sub-second updates
	// to the Admin API could be problematic (or at least operate differently) based on
	// which storage backend was in use (i.e. "dbless", "postgres"). This is a workaround
	// for improvements we still need to investigate upstream.
	//
	// See Also: https://github.com/Kong/kubernetes-ingress-controller/issues/1398
	DefaultSyncSeconds float32 = 3.0
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

	// server configuration, flow control, channels and utility attributes
	stagger            time.Duration
	syncTicker         *time.Ticker
	stopCh             chan struct{}
	configApplied      bool
	configAppliedMutex sync.RWMutex
}

// NewSynchronizer will provide a new Synchronizer object. Note that this
// starts some background goroutines and the caller is resonsible for marking
// the provided context.Context as "Done()" to shut down the background routines
func NewSynchronizer(logger logrus.FieldLogger, dataplaneClient Client) (*Synchronizer, error) {
	stagger, err := time.ParseDuration(fmt.Sprintf("%gs", DefaultSyncSeconds))
	if err != nil {
		return nil, err
	}
	return NewSynchronizerWithStagger(logger, dataplaneClient, stagger)
}

// NewSynchronizer will provide a new Synchronizer object with a specified
// stagger time for data-plane updates to occur. Note that this starts some
// background goroutines and the caller is resonsible for marking the provided
// context.Context as "Done()" to shut down the background routines
func NewSynchronizerWithStagger(logger logrus.FieldLogger, dataplaneClient Client, stagger time.Duration) (*Synchronizer, error) {
	synchronizer := &Synchronizer{
		logger:          logrusr.New(logger),
		dataplaneClient: dataplaneClient,
		stagger:         stagger,
		stopCh:          make(chan struct{}),
		syncTicker:      time.NewTicker(stagger),
		configApplied:   false,
	}

	return synchronizer, nil
}

// -----------------------------------------------------------------------------
// Synchronizer - Public Methods
// -----------------------------------------------------------------------------

func (p *Synchronizer) Start(ctx context.Context) error {
	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2249
	// This is a temporary mitigation to allow some time for controllers to
	// populate their dataplaneClient cache.
	time.Sleep(time.Second * 5)
	go p.startUpdateServer(ctx)
	return nil
}

func (p *Synchronizer) IsReady() bool {
	// If the proxy is has no database, it is only ready after a successful sync
	// Otherwise, it has no configuration loaded
	if p.dataplaneClient.DBMode() == "off" {
		p.configAppliedMutex.RLock()
		defer p.configAppliedMutex.RUnlock()
		return p.configApplied
	}
	// If the proxy has a database, it is ready immediately
	// It will load existing configuration from the database
	return true
}

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
			if err := ctx.Err(); err != nil {
				p.logger.Error(err, "context completed with error")
			}
			p.syncTicker.Stop()
			return
		case <-p.syncTicker.C:
			if err := p.dataplaneClient.Update(ctx); err != nil {
				p.logger.Error(err, "could not update kong admin")
				break
			}
			initialConfig.Do(p.markConfigApplied)
		}
	}
}

// -----------------------------------------------------------------------------
// Synchronizer - Private Methods - Helper
// -----------------------------------------------------------------------------

// markConfigApplied marks that config has been applied
func (p *Synchronizer) markConfigApplied() {
	p.configAppliedMutex.Lock()
	defer p.configAppliedMutex.Unlock()
	p.configApplied = true
}
