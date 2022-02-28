package dataplane

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blang/semver/v4"
	"github.com/bombsimon/logrusr/v2"
	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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

	// kong configuration metadata
	dbmode  string
	version semver.Version

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

	// TODO: this initialization needs to move into the dataplane client
	if err := synchronizer.initialize(); err != nil {
		return nil, err
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
	if p.dbmode == "off" {
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

// initialize validates connectivity with the Kong proxy and some of the configuration options thereof
// and populates several local attributes given retrieved configuration data from the proxy root config.
//
// Note: this must be run (and must succeed) in order to successfully start the cache server.
func (p *Synchronizer) initialize() error {
	// FIXME/TODO: technically we don't have any dataplane implementations
	// EXCEPT Kong so this switch is here temporarily while we move the
	// relevant functionality into the dataplane client.
	switch dataplaneClient := p.dataplaneClient.(type) {
	case *KongClient:
		// download the kong root configuration (and validate connectivity to the proxy API)
		root, err := dataplaneClient.RootWithTimeout()
		if err != nil {
			return err
		}

		// pull the proxy configuration out of the root config and validate it
		proxyConfig, ok := root["configuration"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid root configuration, expected a map[string]interface{} got %T", proxyConfig["configuration"])
		}
		// validate the database configuration for the proxy and check for supported database configurations
		dbmode, ok := proxyConfig["database"].(string)
		if !ok {
			return fmt.Errorf("invalid database configuration, expected a string got %t", proxyConfig["database"])
		}
		switch dbmode {
		case "off", "":
			dataplaneClient.KongConfig.InMemory = true
		case "postgres":
			dataplaneClient.KongConfig.InMemory = false
		case "cassandra":
			return fmt.Errorf("Cassandra-backed deployments of Kong managed by the ingress controller are no longer supported; you must migrate to a Postgres-backed or DB-less deployment")
		default:
			return fmt.Errorf("%s is not a supported database backend", dbmode)
		}

		// validate the proxy version
		proxySemver, err := kong.ParseSemanticVersion(kong.VersionFromInfo(root))
		if err != nil {
			return err
		}

		// store the gathered configuration options
		dataplaneClient.KongConfig.Version = proxySemver
		p.dbmode = dbmode
		p.version = proxySemver
	}

	return nil
}

// markConfigApplied marks that config has been applied
func (p *Synchronizer) markConfigApplied() {
	p.configAppliedMutex.Lock()
	defer p.configAppliedMutex.Unlock()
	p.configApplied = true
}

// -----------------------------------------------------------------------------
// Private Helper Functions
// -----------------------------------------------------------------------------

// fetchCustomEntities returns the value of the "config" key from a Secret (identified by a "namespace/secretName"
// string in the store.
func fetchCustomEntities(secretName string, store store.Storer) ([]byte, error) {
	ns, name, err := util.ParseNameNS(secretName)
	if err != nil {
		return nil, fmt.Errorf("parsing kong custom entities secret: %w", err)
	}
	secret, err := store.GetSecret(ns, name)
	if err != nil {
		return nil, fmt.Errorf("fetching secret: %w", err)
	}
	config, ok := secret.Data["config"]
	if !ok {
		return nil, fmt.Errorf("'config' key not found in "+
			"custom entities secret '%v'", secretName)
	}
	return config, nil
}
