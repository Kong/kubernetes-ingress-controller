package proxy

import (
	"context"
	"fmt"
	"time"

	"github.com/blang/semver/v4"
	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/kong/kubernetes-ingress-controller/pkg/parser"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
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
) (Proxy, error) {
	stagger, err := time.ParseDuration(fmt.Sprintf("%gs", DefaultSyncSeconds))
	if err != nil {
		return nil, err
	}
	return NewCacheBasedProxyWithStagger(ctx, logger, k8s, kongConfig, ingressClassName, enableReverseSync, stagger)
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
) (Proxy, error) {
	// configure the cachestores and the proxy instance
	cache := store.NewCacheStores()
	proxy := &clientgoCachedProxyResolver{
		cache: &cache,

		kongConfig:        kongConfig,
		enableReverseSync: enableReverseSync,

		deprecatedLogger: logger,
		logger:           logrusr.NewLogger(logger),

		ingressClassName: ingressClassName,
		stopCh:           make(chan struct{}),

		ctx:     ctx,
		update:  make(chan *cachedObject, DefaultObjectBufferSize),
		del:     make(chan *cachedObject, DefaultObjectBufferSize),
		stagger: stagger,
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
	k8s   client.Client
	cache *store.CacheStores

	// lastConfigSHA indicates the last SHA sum for the last configuration
	// updated in the Kong Proxy and is used to avoid making unnecessary updates.
	lastConfigSHA []byte

	// kong configuration
	kongConfig        sendconfig.Kong
	kongRootConfig    map[string]interface{}
	kongProxyConfig   map[string]interface{}
	enableReverseSync bool
	dbmode            string
	version           semver.Version

	// KongCustomEntitiesSecret is a "namespace/name" Secret locator for a Secret
	// that contains raw YAML custom entities, for use with DB-less mode
	KongCustomEntitiesSecret string

	// cache server configuration, flow control, channels and utility attributes
	ingressClassName string
	ctx              context.Context
	stagger          time.Duration
	stopCh           chan struct{}

	// New code should log using "logger". "deprecatedLogger" is here for compatibility with legacy code that relies
	// on the logrus API.
	deprecatedLogger logrus.FieldLogger
	logger           logr.Logger

	// channels
	update chan *cachedObject
	del    chan *cachedObject
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

// objectTracker is a secondary cache used to track objects that have been updated/deleted between successful updates of the Kong Admin API.
type objectTracker map[string]*cachedObject

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
	p.logger.Info("the proxy cache server has been started")

	// syncTicker is a regular interval to check for cache updates and resolve the cache to the Kong Admin API
	syncTicker := time.NewTicker(p.stagger)

	// updates tracks whether any updates/deletes were tracked this cycle
	for {
		select {
		case cobj := <-p.update:
			if err := p.cacheUpdate(cobj); err != nil {
				p.logger.Error(err, "object could not be updated in the cache and will be discarded")
				break
			}
		case cobj := <-p.del:
			if err := p.cacheDelete(cobj); err != nil {
				p.logger.Error(err, "object could not be deleted from the cache and will be discarded")
				break
			}
		case <-syncTicker.C:
			if err := p.updateKongAdmin(); err != nil {
				p.logger.Error(err, "could not update kong admin")
			}
		case <-p.ctx.Done():
			p.logger.Info("the proxy cache server's context is done, shutting down")
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
		return fmt.Errorf("invalid database configuration, expected a string got %T", proxyConfig["database"])
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

// updateKongAdmin will take whatever the current state of the Proxy.cache is an convert that to Kong DSL
// and apply the resulting configuration to the Kong Admin API.
func (p *clientgoCachedProxyResolver) updateKongAdmin() error {
	// build the kongstate object from the Kubernetes objects in the storer
	storer := store.New(*p.cache, p.ingressClassName, false, false, false, p.deprecatedLogger)
	kongstate, err := parser.Build(p.deprecatedLogger, storer)
	if err != nil {
		return err
	}

	// generate the deck configuration to be applied to the admin API
	targetConfig := deckgen.ToDeckContent(p.ctx, p.deprecatedLogger, kongstate, nil, nil)

	// retrieve custom entities
	var customEntities []byte
	if p.kongConfig.InMemory && p.KongCustomEntitiesSecret != "" {
		customEntities, err = fetchCustomEntities(p.KongCustomEntitiesSecret, storer)
		if err != nil {
			// failure to fetch custom entities shouldn't block updates
			p.logger.Error(err, "failed to fetch custom entities")
		}
	}

	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(p.ctx, 10*time.Second)
	defer cancel()
	p.lastConfigSHA, err = sendconfig.PerformUpdate(timedCtx,
		p.deprecatedLogger, &p.kongConfig,
		p.kongConfig.InMemory, p.enableReverseSync,
		targetConfig, nil, customEntities, p.lastConfigSHA,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *clientgoCachedProxyResolver) kongRootWithTimeout() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(p.ctx, 3*time.Second)
	defer cancel()
	return p.kongConfig.Client.Root(ctx)
}

// fetchCustomEntities returns the value of the "config" key from a Secret (identified by a "namespace/secretName"
// string in the store.
func fetchCustomEntities(secret string, store store.Storer) ([]byte, error) {
	ns, name, err := util.ParseNameNS(secret)
	if err != nil {
		return nil, fmt.Errorf("parsing kong custom entities secret: %w", err)
	}
	kSecret, err := store.GetSecret(ns, name)
	if err != nil {
		return nil, fmt.Errorf("fetching secret: %w", err)
	}
	config, ok := kSecret.Data["config"]
	if !ok {
		return nil, fmt.Errorf("'config' key not found in "+
			"custom entities secret '%v'", secret)
	}
	return config, nil
}
