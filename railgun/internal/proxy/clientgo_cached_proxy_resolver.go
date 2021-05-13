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

	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/kong/kubernetes-ingress-controller/pkg/parser"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
)

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Public Functions
// -----------------------------------------------------------------------------

// NewCacheBasedProxy will provide a new Proxy object. Note that this starts some background services
// and the caller is thereafter responsible for closing the Proxy.StopCh.
func NewCacheBasedProxy(ctx context.Context,
	logger logrus.FieldLogger,
	k8s client.Client,
	kongConfig sendconfig.Kong,
	ingressClassName string,
	processClasslessIngressV1Beta1, processClasslessIngressV1, processClasslessKongConsumer, enableReverseSync bool,
) (Proxy, error) {
	return NewCacheBasedProxyWithStagger(ctx, logger, k8s, kongConfig, ingressClassName, processClasslessIngressV1Beta1, processClasslessIngressV1, processClasslessKongConsumer, enableReverseSync, DefaultStagger)
}

func NewCacheBasedProxyWithStagger(ctx context.Context,
	logger logrus.FieldLogger,
	k8s client.Client,
	kongConfig sendconfig.Kong,
	ingressClassName string,
	processClasslessIngressV1Beta1, processClasslessIngressV1, processClasslessKongConsumer, enableReverseSync bool,
	stagger time.Duration,
) (Proxy, error) {
	// download the kong root configuration (and validate connectivity to the proxy API)
	root, err := kongConfig.Client.Root(ctx)
	if err != nil {
		return nil, err
	}

	// pull the proxy configuration out of the root config and validate it
	proxyConfig, ok := root["configuration"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid database configuration, expected a string got %t", proxyConfig["database"])
	}

	// validate the database configuration for the proxy
	dbmode, ok := proxyConfig["database"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid database configuration, expected a string got %t", proxyConfig["database"])
	}
	if dbmode == "off" || dbmode == "" {
		kongConfig.InMemory = true
	}

	// check for supported database configurations
	switch dbmode {
	case "off":
		kongConfig.InMemory = true
	case "":
		kongConfig.InMemory = true
	case "postgres":
		kongConfig.InMemory = false
	case "cassandra":
		return nil, fmt.Errorf("Cassandra-backed deployments of Kong managed by the ingress controller are no longer supported; you must migrate to a Postgres-backed or DB-less deployment")
	default:
		return nil, fmt.Errorf("%s is not a supported database backend", dbmode)
	}

	// validate the proxy version
	proxyStringVersion, ok := root["version"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid value found for version in kong root config: %t", root["version"])
	}
	proxySemver, err := ctrlutils.GetSemver(proxyStringVersion)
	if err != nil {
		return nil, err
	}
	kongConfig.Version = proxySemver

	// configure the cache stores
	cache := store.NewCacheStores()

	// configure the proxy
	proxy := &clientgoCachedProxyResolver{
		cache: &cache,

		kongConfig:        kongConfig,
		kongRootConfig:    root,
		kongProxyConfig:   proxyConfig,
		enableReverseSync: enableReverseSync,
		dbmode:            dbmode,
		version:           proxySemver,

		deprecatedLogger: logger,
		logger:           logrusr.NewLogger(logger),

		ingressClassName:               ingressClassName,
		processClasslessIngressV1Beta1: processClasslessIngressV1Beta1,
		processClasslessIngressV1:      processClasslessIngressV1,
		processClasslessKongConsumer:   processClasslessKongConsumer,
		stopCh:                         make(chan struct{}),

		ctx:     ctx,
		update:  make(chan *cachedObject, DefaultObjectBufferSize),
		del:     make(chan *cachedObject, DefaultObjectBufferSize),
		stagger: stagger,
	}

	// start the proxy cache server
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

	// cache store configuration options
	ingressClassName               string
	processClasslessIngressV1Beta1 bool
	processClasslessIngressV1      bool
	processClasslessKongConsumer   bool

	// cache server flow control, channels and utility attributes
	ctx     context.Context
	stagger time.Duration
	stopCh  chan struct{}

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

func (p *clientgoCachedProxyResolver) DBMode() string {
	return p.dbmode
}

func (p *clientgoCachedProxyResolver) Version() semver.Version {
	return p.version
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

	// syncTimer is a regular interval to check for cache updates and resolve the cache to the Kong Admin API
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
	storer := store.New(*p.cache, p.ingressClassName, p.processClasslessIngressV1, p.processClasslessIngressV1Beta1, p.processClasslessKongConsumer, p.deprecatedLogger)
	kongstate, err := parser.Build(p.deprecatedLogger, storer)
	if err != nil {
		return err
	}

	// generate the deck configuration to be applied to the admin API
	targetConfig := deckgen.ToDeckContent(p.ctx, p.deprecatedLogger, kongstate, nil, nil)

	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(p.ctx, 10*time.Second)
	defer cancel()
	p.lastConfigSHA, err = sendconfig.PerformUpdate(timedCtx,
		p.deprecatedLogger, &p.kongConfig,
		p.kongConfig.InMemory, p.enableReverseSync,
		targetConfig, nil, nil, p.lastConfigSHA,
	)
	if err != nil {
		return err
	}

	return nil
}
