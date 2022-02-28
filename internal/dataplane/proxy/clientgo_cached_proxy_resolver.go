package proxy

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

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Public Functions
// -----------------------------------------------------------------------------

// NewCacheBasedProxy will provide a new Proxy object. Note that this starts some background goroutines and the
// caller is resonsible for marking the provided context.Context as "Done()" to shut down the background routines.
func NewCacheBasedProxy(logger logrus.FieldLogger, dataplaneClient dataplane.Client) (Proxy, error) {
	stagger, err := time.ParseDuration(fmt.Sprintf("%gs", DefaultSyncSeconds))
	if err != nil {
		return nil, err
	}
	return NewCacheBasedProxyWithStagger(logger, dataplaneClient, stagger)
}

// NewCacheBasedProxy will provide a new Proxy object. Note that this starts some background goroutines and the caller
// is resonsible for marking the provided context.Context as "Done()" to shut down the background routines. A "stagger"
// time duration is provided to indicate how often the background routines will sync configuration to the Kong Admin API.
func NewCacheBasedProxyWithStagger(logger logrus.FieldLogger, dataplaneClient dataplane.Client, stagger time.Duration) (Proxy, error) {
	// configure the cachestores and the proxy instance
	proxy := &clientgoCachedProxyResolver{
		deprecatedLogger: logger,
		logger:           logrusr.New(logger),

		dataplaneClient: dataplaneClient,
		stagger:         stagger,
		stopCh:          make(chan struct{}),
		syncTicker:      time.NewTicker(stagger),
		configApplied:   false,
	}

	// initialize the proxy which validates connectivity with the Admin API and
	// checks several proxy server attributes such as version and dbmode.
	if err := proxy.initialize(); err != nil {
		return nil, err
	}

	return proxy, nil
}

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Public Methods - Interface Implementation
// -----------------------------------------------------------------------------

func (p *clientgoCachedProxyResolver) NeedLeaderElection() bool {
	return true
}

func (p *clientgoCachedProxyResolver) Start(ctx context.Context) error {
	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2249
	// This is a temporary mitigation to allow some time for controllers to populate the proxy cache
	time.Sleep(time.Second * 5)
	go p.startProxyUpdateServer(ctx)
	return nil
}

func (p *clientgoCachedProxyResolver) IsReady() bool {
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

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Private Types
// -----------------------------------------------------------------------------

// clientgoCachedProxyResolver represents the cached objects and Kong DSL configuration.
//
// This implements the Proxy interface to provide asynchronous, non-blocking updates to
// the Kong Admin API for controller-runtime based controller managers.
//
// This object's attributes are immutable (private), and it is threadsafe.
type clientgoCachedProxyResolver struct {
	// New code should log using "logger". "deprecatedLogger" is here for compatibility with legacy code that relies
	// on the logrus API.
	deprecatedLogger logrus.FieldLogger
	logger           logr.Logger

	// dataplane client to send updates to the Kong Admin API
	dataplaneClient dataplane.Client

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

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Private Methods - Server Utilities
// -----------------------------------------------------------------------------

// startProxyUpdateServer runs a server in a background goroutine that is responsible for
// updating the kong proxy backend at regular intervals.
func (p *clientgoCachedProxyResolver) startProxyUpdateServer(ctx context.Context) {
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

// initialize validates connectivity with the Kong proxy and some of the configuration options thereof
// and populates several local attributes given retrieved configuration data from the proxy root config.
//
// Note: this must be run (and must succeed) in order to successfully start the cache server.
func (p *clientgoCachedProxyResolver) initialize() error {
	// FIXME/TODO: technically we don't have any dataplane implementations
	// EXCEPT Kong so this switch is here temporarily while we move the
	// relevant functionality into the dataplane client.
	switch dataplaneClient := p.dataplaneClient.(type) {
	case *dataplane.KongClient:
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
func (p *clientgoCachedProxyResolver) markConfigApplied() {
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
