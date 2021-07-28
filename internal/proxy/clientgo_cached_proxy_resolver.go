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
	"github.com/kong/kubernetes-ingress-controller/internal/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/internal/store"
	"github.com/kong/kubernetes-ingress-controller/internal/util"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	timeout time.Duration,
) (Proxy, error) {
	stagger, err := time.ParseDuration(fmt.Sprintf("%gs", DefaultSyncSeconds))
	if err != nil {
		return nil, err
	}
	return NewCacheBasedProxyWithStagger(ctx, logger, k8s, kongConfig, ingressClassName, enableReverseSync, stagger, timeout, kongUpdater)
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
	timeout time.Duration,
	kongUpdater KongUpdater,
) (Proxy, error) {
	// configure the cachestores and the proxy instance
	cache := store.NewCacheStores()
	proxy := &clientgoCachedProxyResolver{
		cache: &cache,

		kongConfig:        kongConfig,
		kongUpdater:       kongUpdater,
		enableReverseSync: enableReverseSync,

		deprecatedLogger: logger,
		logger:           logrusr.NewLogger(logger),

		ingressClassName: ingressClassName,
		stopCh:           make(chan struct{}),

		ctx:        ctx,
		update:     make(chan *cachedObject, DefaultObjectBufferSize),
		del:        make(chan *cachedObject, DefaultObjectBufferSize),
		stagger:    stagger,
		timeout:    timeout,
		syncTicker: time.NewTicker(stagger),

		internalCacheLock: &sync.RWMutex{},
	}

	// initialize the proxy which validates connectivity with the Admin API and
	// checks several proxy server attributes such as version and dbmode.
	if err := proxy.initialize(); err != nil {
		return nil, err
	}

	// start the proxy cache server in the background
	go proxy.startCacheServer()

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
	kongUpdater KongUpdater

	// cache server configuration, flow control, channels and utility attributes
	ingressClassName string
	ctx              context.Context
	stagger          time.Duration
	timeout          time.Duration
	syncTicker       *time.Ticker
	stopCh           chan struct{}

	// New code should log using "logger". "deprecatedLogger" is here for compatibility with legacy code that relies
	// on the logrus API.
	deprecatedLogger logrus.FieldLogger
	logger           logr.Logger

	// channels
	update chan *cachedObject
	del    chan *cachedObject

	// locks
	internalCacheLock *sync.RWMutex
}

// cacheAction indicates what caching action (update, delete) was taken for any particular runtime.Object.
type cacheAction string

var (
	// updated indicates that this object was either newly added OR updated in the cache (no distinction made)
	updated cacheAction = "updated"

	// deleted indicates that this object was removed from the cache
	deleted cacheAction = "deleted"
)

// cachedObject represents an object that has been processed by the cacheServer
type cachedObject struct {
	action cacheAction
	err    error

	key        string
	runtimeObj runtime.Object
}

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Public Methods - Interface Implementation
// -----------------------------------------------------------------------------

func (p *clientgoCachedProxyResolver) UpdateObject(obj client.Object) error {
	cobj := &cachedObject{action: updated, key: p.clientObjectKey(obj), runtimeObj: obj.DeepCopyObject()}
	select {
	case p.update <- cobj:
		return nil
	default:
		return fmt.Errorf("the proxy is too busy to accept requests at this time, try again later")
	}
}

func (p *clientgoCachedProxyResolver) DeleteObject(obj client.Object) error {
	cobj := &cachedObject{action: deleted, key: p.clientObjectKey(obj), runtimeObj: obj.DeepCopyObject()}
	select {
	case p.del <- cobj:
		return nil
	default:
		return fmt.Errorf("the proxy is too busy to accept requests at this time, try again later")
	}
}

func (p *clientgoCachedProxyResolver) ObjectExists(obj client.Object) (bool, error) {
	p.internalCacheLock.RLock()
	defer p.internalCacheLock.RUnlock()
	_, exists, err := p.cache.Get(obj)
	return exists, err
}

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Private Methods - Cache Server
// -----------------------------------------------------------------------------

// startCacheServer runs a server in a background goroutine that is responsible for:
//
//   1. processing kubernetes object updates (add, replace)
//   2. processing kubernetes object deletes
//   3. regularly synchronizing configuration to the Kong Admin API (staggered)
//
// While processing objects the cacheServer will (synchronously) convert the objects to kong DSL and
// submit POST updates to the Kong Admin API with the new configuration.
func (p *clientgoCachedProxyResolver) startCacheServer() {
	backendNeedsSync := false
	p.logger.Info("the proxy cache server has been started")
	for {
		select {
		case cobj := <-p.update:
			if err := p.cacheUpdate(cobj); err != nil {
				p.logger.Error(err, "object could not be updated in the cache and will be discarded")
				break
			}
			backendNeedsSync = true
		case cobj := <-p.del:
			if err := p.cacheDelete(cobj); err != nil {
				p.logger.Error(err, "object could not be deleted from the cache and will be discarded")
				break
			}
			backendNeedsSync = true
		case <-p.syncTicker.C:
			if !p.enableReverseSync && !backendNeedsSync {
				break
			}

			updateConfigSHA, err := p.kongUpdater(p.ctx, p.lastConfigSHA, p.cache, p.ingressClassName, p.deprecatedLogger, p.kongConfig, p.enableReverseSync)
			if err != nil {
				p.logger.Error(err, "could not update kong admin")
				break
			}
			p.lastConfigSHA = updateConfigSHA
			backendNeedsSync = false
		case <-p.ctx.Done():
			p.logger.Info("the proxy cache server's context is done, shutting down")
			if err := p.ctx.Err(); err != nil {
				p.logger.Error(err, "context completed with error")
			}
			p.syncTicker.Stop()
			return
		}
	}
}

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Private Methods - Cache Server Utils
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

	return nil
}

// objectKey provides a key unique to the cacheServer for tracking objects that are being resolved.
func (p *clientgoCachedProxyResolver) clientObjectKey(obj client.Object) string {
	return fmt.Sprintf("%s/%s/%s", obj.GetObjectKind().GroupVersionKind().String(), obj.GetNamespace(), obj.GetName())
}

// cacheUpdate caches the provided object to the proxy cache and the provided object tracker and reports any errors.
func (p *clientgoCachedProxyResolver) cacheUpdate(cobj *cachedObject) error {
	cobj.err = p.cache.Add(cobj.runtimeObj)
	return cobj.err
}

// cacheDelete removes the cache entry the provided object from the proxy cache and the provided object tracker and reports any errors.
func (p *clientgoCachedProxyResolver) cacheDelete(cobj *cachedObject) error {
	cobj.err = p.cache.Delete(cobj.runtimeObj)
	return cobj.err
}

// kongRootWithTimeout provides the root configuration from Kong, but uses a configurable timeout to avoid long waits if the Admin API
// is not yet ready to respond. If a timeout error occurs, the caller is responsible for providing a retry mechanism.
func (p *clientgoCachedProxyResolver) kongRootWithTimeout() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(p.ctx, p.timeout)
	defer cancel()
	return p.kongConfig.Client.Root(ctx)
}

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
