package proxy

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/kong/kubernetes-ingress-controller/pkg/parser"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/store"
)

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Public Functions
// -----------------------------------------------------------------------------

// NewCacheBasedProxy will provide a new Proxy object. Note that this starts some background services
// and the caller is thereafter responsible for closing the Proxy.StopCh.
func NewCacheBasedProxy(ctx context.Context, logger logr.Logger, k8s client.Client, kongConfig sendconfig.Kong, ingressClassName string, processClasslessIngressV1Beta1 bool, processClasslessIngressV1 bool, processClasslessKongConsumer bool) Proxy {
	return NewCacheBasedProxyWithStagger(ctx, logger, k8s, kongConfig, ingressClassName, processClasslessIngressV1Beta1, processClasslessIngressV1, processClasslessKongConsumer, DefaultStagger)
}

func NewCacheBasedProxyWithStagger(ctx context.Context, logger logr.Logger, k8s client.Client, kongConfig sendconfig.Kong, ingressClassName string, processClasslessIngressV1Beta1 bool, processClasslessIngressV1 bool, processClasslessKongConsumer bool, stagger time.Duration) Proxy {
	cache := store.NewCacheStores()
	proxy := &clientgoCachedProxyResolver{
		kongConfig: kongConfig,
		cache:      &cache,
		logger:     logr.Discard(),

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
	go proxy.startCacheServer()
	return proxy
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

	// kong configuration
	kongConfig sendconfig.Kong

	// cache store configuration options
	ingressClassName               string
	processClasslessIngressV1Beta1 bool
	processClasslessIngressV1      bool
	processClasslessKongConsumer   bool

	// cache server flow control, channels and utility attributes
	ctx     context.Context
	stagger time.Duration
	logger  logr.Logger
	stopCh  chan struct{}

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

	// syncTimer is a regular interval to check for cache updates and resolve the cache to the Kong Admin API
	syncTicker := time.NewTicker(p.stagger)

	// updates tracks whether any updates/deletes were tracked this cycle
	updates := false
	for {
		select {
		case cobj := <-p.update:
			if err := p.cacheUpdate(cobj); err != nil {
				p.logger.Error(err, "object could not be updated in the cache and will be discarded")
				break
			}
			updates = true
		case cobj := <-p.del:
			if err := p.cacheDelete(cobj); err != nil {
				p.logger.Error(err, "object could not be deleted from the cache and will be discarded")
				break
			}
			updates = true
		case <-syncTicker.C:
			// if there are no relevant object updates for this cycle, then there's no reason to
			// bother the Kong Proxy with updates from the cache.
			if !updates {
				break
			}
			if err := p.updateKongAdmin(); err != nil {
				p.logger.Error(err, "could not update kong admin")
			}
			updates = false
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
	storer := store.New(*p.cache, p.ingressClassName, p.processClasslessIngressV1, p.processClasslessIngressV1Beta1, p.processClasslessKongConsumer, logrus.StandardLogger())
	kongstate, err := parser.Build(logrus.StandardLogger(), storer)
	if err != nil {
		return err
	}

	// generate the deck configuration to be applied to the admin API
	targetConfig := deckgen.ToDeckContent(p.ctx, logrus.StandardLogger(), kongstate, nil, nil)

	// apply the configuration update in Kong
	timedCtx, cancel := context.WithTimeout(p.ctx, 10*time.Second)
	defer cancel()
	_, err = sendconfig.PerformUpdate(timedCtx, logrus.StandardLogger(), &p.kongConfig, true, false, targetConfig, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
