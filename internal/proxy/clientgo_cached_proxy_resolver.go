package proxy

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blang/semver/v4"
	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Public Functions
// -----------------------------------------------------------------------------

// NewCacheBasedProxy will provide a new Proxy object. Note that this starts some background goroutines and the
// caller is resonsible for marking the provided context.Context as "Done()" to shut down the background routines.
func NewCacheBasedProxy(ctx context.Context,
	logger logrus.FieldLogger,
	k8s client.Client,
	kongConfig sendconfig.Kong,
	ingressClassName string,
	enableReverseSync bool,
	kongUpdater KongUpdater,
	diagnostic util.ConfigDumpDiagnostic,
	proxyRequestTimeout time.Duration,
) (Proxy, error) {
	stagger, err := time.ParseDuration(fmt.Sprintf("%gs", DefaultSyncSeconds))
	if err != nil {
		return nil, err
	}
	return NewCacheBasedProxyWithStagger(ctx, logger, k8s, kongConfig, ingressClassName, enableReverseSync, stagger, proxyRequestTimeout, diagnostic, kongUpdater)
}

// NewCacheBasedProxy will provide a new Proxy object. Note that this starts some background goroutines and the caller
// is resonsible for marking the provided context.Context as "Done()" to shut down the background routines. A "stagger"
// time duration is provided to indicate how often the background routines will sync configuration to the Kong Admin API.
func NewCacheBasedProxyWithStagger(ctx context.Context,
	logger logrus.FieldLogger,
	k8s client.Client,
	kongConfig sendconfig.Kong,
	ingressClassName string,
	enableReverseSync bool,
	stagger time.Duration,
	proxyRequestTimeout time.Duration,
	diagnostic util.ConfigDumpDiagnostic,
	kongUpdater KongUpdater,
) (Proxy, error) {
	// configure the cachestores and the proxy instance
	cache := store.NewCacheStores()
	proxy := &clientgoCachedProxyResolver{
		cache: &cache,

		kongConfig:        kongConfig,
		kongUpdater:       kongUpdater,
		diagnostic:        diagnostic,
		enableReverseSync: enableReverseSync,

		deprecatedLogger: logger,
		logger:           logrusr.NewLogger(logger),

		ingressClassName: ingressClassName,
		stopCh:           make(chan struct{}),

		ctx:                 ctx,
		stagger:             stagger,
		proxyRequestTimeout: proxyRequestTimeout,
		syncTicker:          time.NewTicker(stagger),

		configApplied: false,
	}

	// initialize the proxy which validates connectivity with the Admin API and
	// checks several proxy server attributes such as version and dbmode.
	if err := proxy.initialize(); err != nil {
		return nil, err
	}

	// start the proxy update server in the background
	go proxy.startProxyUpdateServer()

	return proxy, nil
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
	// kubernetes configuration
	cache *store.CacheStores

	// lastConfigSHA indicates the last SHA sum for the last configuration
	// updated in the Kong Proxy and is used to avoid making unnecessary updates.
	lastConfigSHA []byte

	// configApplied is true if config has been applied at least once
	configApplied      bool
	configAppliedMutex sync.RWMutex

	// kong configuration
	kongConfig        sendconfig.Kong
	enableReverseSync bool
	dbmode            string
	version           semver.Version

	// KongCustomEntitiesSecret is a "namespace/name" Secret locator for a Secret
	// that contains raw YAML custom entities, for use with DB-less mode
	KongCustomEntitiesSecret string

	// kongUpdater is the function that will be used by the cache server to ultimately make the API
	// call to resolve the current cache to the Kong Admin API configuration endpoint.
	// It may ship diagnostic information through diagnostic
	kongUpdater KongUpdater
	diagnostic  util.ConfigDumpDiagnostic
	promMetrics *metrics.CtrlFuncMetrics

	// server configuration, flow control, channels and utility attributes
	ingressClassName    string
	ctx                 context.Context
	stagger             time.Duration
	proxyRequestTimeout time.Duration
	syncTicker          *time.Ticker
	stopCh              chan struct{}

	// New code should log using "logger". "deprecatedLogger" is here for compatibility with legacy code that relies
	// on the logrus API.
	deprecatedLogger logrus.FieldLogger
	logger           logr.Logger
}

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Public Methods - Interface Implementation
// -----------------------------------------------------------------------------

func (p *clientgoCachedProxyResolver) UpdateObject(obj client.Object) error {
	return p.cache.Add(obj)
}

func (p *clientgoCachedProxyResolver) DeleteObject(obj client.Object) error {
	return p.cache.Delete(obj)
}

func (p *clientgoCachedProxyResolver) ObjectExists(obj client.Object) (bool, error) {
	_, exists, err := p.cache.Get(obj)
	return exists, err
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

func (p *clientgoCachedProxyResolver) Listeners(ctx context.Context) ([]kong.ProxyListener, []kong.StreamListener, error) {
	return p.kongConfig.Client.Listeners(ctx)
}

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Private Methods - Servers
// -----------------------------------------------------------------------------

// startProxyUpdateServer runs a server in a background goroutine that is responsible for
// updating the kong proxy backend at regular intervals.
func (p *clientgoCachedProxyResolver) startProxyUpdateServer() {
	var initialConfig sync.Once
	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info("context done: shutting down the proxy update server")
			if err := p.ctx.Err(); err != nil {
				p.logger.Error(err, "context completed with error")
			}
			p.syncTicker.Stop()
			return
		case <-p.syncTicker.C:
			updateConfigSHA, err := p.kongUpdater(p.ctx, p.lastConfigSHA, p.cache,
				p.ingressClassName, p.deprecatedLogger, p.kongConfig, p.enableReverseSync, p.diagnostic, p.proxyRequestTimeout, p.promMetrics)
			if err != nil {
				p.logger.Error(err, "could not update kong admin")
				break
			}
			p.lastConfigSHA = updateConfigSHA
			initialConfig.Do(p.markConfigApplied)
		}
	}
}

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Private Methods - Server Utils
// -----------------------------------------------------------------------------

// initialize validates connectivity with the Kong proxy and some of the configuration options thereof
// and populates several local attributes given retrieved configuration data from the proxy root config.
//
// Note: this must be run (and must succeed) in order to successfully start the cache server.
func (p *clientgoCachedProxyResolver) initialize() error {
	// download the kong root configuration (and validate connectivity to the proxy API)
	root, err := p.kongRootWithTimeout()
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
		p.kongConfig.InMemory = true
	case "postgres":
		p.kongConfig.InMemory = false
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
	p.kongConfig.Version = proxySemver
	p.dbmode = dbmode
	p.version = proxySemver
	p.promMetrics = metrics.NewCtrlFuncMetrics()

	return nil
}

// kongRootWithTimeout provides the root configuration from Kong, but uses a configurable timeout to avoid long waits if the Admin API
// is not yet ready to respond. If a timeout error occurs, the caller is responsible for providing a retry mechanism.
func (p *clientgoCachedProxyResolver) kongRootWithTimeout() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(p.ctx, p.proxyRequestTimeout)
	defer cancel()
	return p.kongConfig.Client.Root(ctx)
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
