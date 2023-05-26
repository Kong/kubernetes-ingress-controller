package configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	discoveryv1 "k8s.io/api/discovery/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	cfgtypes "github.com/kong/kubernetes-ingress-controller/v2/internal/manager/config/types"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// KongAdminAPIServiceReconciler reconciles Kong Admin API Service Endpointslices
// and notifies the provided notifier about those.
type KongAdminAPIServiceReconciler struct {
	client.Client

	// ServiceNN is the service NamespacedName to watch EndpointSlices for.
	ServiceNN k8stypes.NamespacedName
	// PortNames is the set of port names that Admin API Service ports will be
	// matched against.
	PortNames        sets.Set[string]
	Log              logr.Logger
	CacheSyncTimeout time.Duration
	// EndpointsNotifier is used to notify about Admin API endpoints changes.
	// We're going to call this only with endpoints when they change.
	EndpointsNotifier EndpointsNotifier
	DNSStrategy       cfgtypes.DNSStrategy

	Cache DiscoveredAdminAPIsCache
}

type DiscoveredAdminAPIsCache map[k8stypes.NamespacedName]sets.Set[adminapi.DiscoveredAdminAPI]

type EndpointsNotifier interface {
	Notify(adminAPIs []adminapi.DiscoveredAdminAPI)
}

// SetupWithManager sets up the controller with the Manager.
func (r *KongAdminAPIServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("KongAdminAPIEndpoints", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}

	if r.Cache == nil {
		r.Cache = make(DiscoveredAdminAPIsCache)
	}

	return c.Watch(
		source.Kind(mgr.GetCache(), &discoveryv1.EndpointSlice{}),
		&handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(r.shouldReconcileEndpointSlice),
	)
}

func (r *KongAdminAPIServiceReconciler) shouldReconcileEndpointSlice(obj client.Object) bool {
	endpoints, ok := obj.(*discoveryv1.EndpointSlice)
	if !ok {
		return false
	}

	if endpoints.Namespace != r.ServiceNN.Namespace {
		return false
	}

	if !lo.ContainsBy(endpoints.OwnerReferences, func(ref metav1.OwnerReference) bool {
		return ref.Kind == "Service" && ref.Name == r.ServiceNN.Name
	}) {
		return false
	}

	return true
}

// +kubebuilder:rbac:groups="discovery.k8s.io",resources=endpointslices,verbs=get;list;watch

// Reconcile processes the watched objects.
func (r *KongAdminAPIServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var endpoints discoveryv1.EndpointSlice
	if err := r.Get(ctx, req.NamespacedName, &endpoints); err != nil {
		if apierrors.IsNotFound(err) {
			// If we have an entry for this EndpointSlice, remove it and notify about the change.
			if _, ok := r.Cache[req.NamespacedName]; ok {
				delete(r.Cache, req.NamespacedName)
				r.notify()
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	r.Log.Info("reconciling Admin API EndpointSlice", "namespace", req.Namespace, "name", req.Name)

	if !endpoints.DeletionTimestamp.IsZero() {
		r.Log.V(util.DebugLevel).Info("EndpointSlice is being deleted",
			"type", "EndpointSlice", "namespace", req.Namespace, "name", req.Name,
		)

		// If we have an entry for this EndpointSlice, remove it and notify about the change.
		if _, ok := r.Cache[req.NamespacedName]; ok {
			delete(r.Cache, req.NamespacedName)
			r.notify()
		}

		return ctrl.Result{}, nil
	}

	cached, ok := r.Cache[req.NamespacedName]
	if !ok {
		// If we don't have an entry for this EndpointSlice then save it and notify
		// about the change.
		var err error
		r.Cache[req.NamespacedName], err = adminapi.AdminAPIsFromEndpointSlice(endpoints, r.PortNames, r.DNSStrategy)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf(
				"failed getting Admin API from endpoints: %s/%s: %w", endpoints.Namespace, endpoints.Name, err,
			)
		}
		r.notify()
		return ctrl.Result{}, nil
	}

	// We do have an entry for this EndpointSlice.
	// If the address set is the same, do nothing.
	// If the address set has changed, update the cache and send a notification.
	addresses, err := adminapi.AdminAPIsFromEndpointSlice(endpoints, r.PortNames, r.DNSStrategy)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf(
			"failed getting Admin API from endpoints: %s/%s: %w", endpoints.Namespace, endpoints.Name, err,
		)
	}
	if cached.Equal(addresses) {
		// No change, don't notify
		return ctrl.Result{}, nil
	}

	r.Cache[req.NamespacedName] = addresses
	r.notify()

	return ctrl.Result{}, nil
}

func (r *KongAdminAPIServiceReconciler) notify() {
	discovered := flattenDiscoveredAdminAPIs(r.Cache)
	r.Log.V(util.DebugLevel).
		Info("notifying about newly detected Admin APIs", "admin_apis", discovered)
	r.EndpointsNotifier.Notify(discovered)
}

func flattenDiscoveredAdminAPIs(cache DiscoveredAdminAPIsCache) []adminapi.DiscoveredAdminAPI {
	var adminAPIs []adminapi.DiscoveredAdminAPI
	for _, v := range cache {
		adminAPIs = append(adminAPIs, v.UnsortedList()...)
	}
	return adminAPIs
}
