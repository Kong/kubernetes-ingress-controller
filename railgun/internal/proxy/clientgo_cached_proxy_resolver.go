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
// and needs to be handled by the statusServer for error checking and reporting.
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
//
// Successes and errors that occur while caching the objects will be reported to the status server.
// An object that can not be cached will be thrown out and not processed for configuration.
//
// If an error occurs when updating the Kong Admin API, all objects being tracked will be reported
// to the Status server with an error condition, but will continue to be tracked until the configuration
// eventually succeeds as new updates arrive.
func (p *clientgoCachedProxyResolver) startCacheServer() {
	p.logger.Info("the proxy cache server has been started")
	tracker := make(objectTracker)

	// syncTimer is a regular interval to check for cache updates and resolve the cache to the Kong Admin API
	syncTicker := time.NewTicker(p.stagger)

	// fullResyncTimer is a regular interval to completely re-seed the cache from the Kubernetes API and resolve it to the Kong Admin API.
	// this is done infrequently but regularly to deal with lost Events (such as due to netsplits)
	fullResyncTicker := time.NewTicker(time.Minute * 10)

	for {
		select {
		case cobj := <-p.update:
			if err := p.cacheUpdate(cobj, tracker); err != nil {
				p.logger.Error(err, "object could not be updated in the cache and will be discarded")
				break
			}
		case cobj := <-p.del:
			if err := p.cacheDelete(cobj, tracker); err != nil {
				p.logger.Error(err, "object could not be deleted from the cache and will be discarded")
				break
			}
		case <-syncTicker.C:
			// if there are no relevant object updates for this cycle, then there's no reason to
			// bother the Kong Proxy with updates from the cache.
			if len(tracker) < 1 {
				break
			}

			// resolve all cached items by converting them to Kong DSL and updating the Admin API configuration.
			// if an error occurs with the Kong Admin API, the objects will be retained for further status reporting.
			newTracker := p.resolveCache(tracker)

			// provide status updates for the proccessed objects
			p.statusUpdates(tracker)

			// refresh the tracker (this will either be empty, or only contain retained objects)
			tracker = newTracker
		case <-fullResyncTicker.C:
			// TODO - interface for resync logic? this may help us change engine parts later
		case <-p.ctx.Done():
			p.logger.Info("the proxy cache server's context is done, shutting down")
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
func (p *clientgoCachedProxyResolver) cacheUpdate(cobj *cachedObject, tracker objectTracker) error {
	cobj.err = p.cache.Add(cobj.runtimeObj)
	tracker[cobj.key] = cobj
	return cobj.err
}

// cacheDelete removes the cache entry the provided object from the proxy cache and the provided object tracker and reports any errors.
func (p *clientgoCachedProxyResolver) cacheDelete(cobj *cachedObject, tracker objectTracker) error {
	cobj.err = p.cache.Delete(cobj.runtimeObj)
	tracker[cobj.key] = cobj
	return cobj.err
}

// resolveCache updates the Kong Admin API and retains tracked objects if a failure occurred so that they can be reported
// on later when a configuration update succeeds.
func (p *clientgoCachedProxyResolver) resolveCache(tracker objectTracker) objectTracker {
	// take the current cache and convert and apply that to the Kong Admin API
	err := p.updateKongAdmin()

	// if any error occurred updating the kong cache, all tracked objects should be retained so that they can later be reported on.
	newTracker := make(objectTracker)
	for key, cobj := range tracker {
		if err != nil {
			cobj.err = err
			newTracker[key] = cobj
		}
	}

	return newTracker
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

// -----------------------------------------------------------------------------
// Client Go Cached Proxy Resolver - Private Methods - Status Updates
// -----------------------------------------------------------------------------

// statusUpdates is responsible for providing status updates for objects which have been resolved by the cacheServer
func (p *clientgoCachedProxyResolver) statusUpdates(tracker objectTracker) {
	return // FIXME - not finished yet, short circuiting
	/*
		wg := sync.WaitGroup{}
		for _, obj := range tracker {
			cobj := *obj.clientObj
			// FIXME - need to get the object's statuses updated, it looks like the tests for controller runtime do this with client-go?
			// NOTE: I vaguely remember being able to use the p.k8s.Status to do this but I need to find an example
			if err := p.k8s.Get(p.ctx, types.NamespacedName{Namespace: cobj.GetNamespace(), Name: cobj.GetName()}, cobj); err != nil {
				panic(err) // FIXME
			}
			p.k8s.Status()
			// FIXME
			wg.Add(1)
			go func() {
				wg.Done()
			}()
		}
		wg.Wait()
	*/
}
