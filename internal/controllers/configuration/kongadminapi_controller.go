package configuration

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	discoveryv1 "k8s.io/api/discovery/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/adminapi"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// KongAdminAPIServiceReconciler reconciles Kong Admin API Service Endpointslices
// and notifies the provided notifier about those.
type KongAdminAPIServiceReconciler struct {
	client.Client

	// ServiceNN is the service NamespacedName to watch EndpointSlices for.
	ServiceNN        types.NamespacedName
	Log              logr.Logger
	CacheSyncTimeout time.Duration
	// EndpointsNotifier is used to notify about Admin API endpoints changes.
	// We're going to call this only with endpoints when they change.
	EndpointsNotifier EndpointsNotifier

	Cache CacheT
}

type CacheT map[types.NamespacedName]sets.Set[string]

type EndpointsNotifier interface {
	Notify(addresses []string)
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
		r.Cache = make(CacheT)
	}

	return c.Watch(
		&source.Kind{Type: &discoveryv1.EndpointSlice{}},
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

//+kubebuilder:rbac:groups="discovery.k8s.io",resources=endpointslices,verbs=get;list;watch

// Reconcile processes the watched objects.
func (r *KongAdminAPIServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var endpoints discoveryv1.EndpointSlice
	if err := r.Get(ctx, req.NamespacedName, &endpoints); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	r.Log.Info("reconciling EndpointSlice", "namespace", req.Namespace, "name", req.Name)

	nn := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Name,
	}

	if !endpoints.DeletionTimestamp.IsZero() {
		r.Log.V(util.DebugLevel).Info("EndpointSlice is being deleted",
			"type", "EndpointSlice", "namespace", req.Namespace, "name", req.Name,
		)

		// If we have an entry for this EndpointSlice...
		if _, ok := r.Cache[nn]; ok {
			// ... remove it and notify about the change.
			delete(r.Cache, nn)
			r.notify()
		}

		return ctrl.Result{}, nil
	}

	cached, ok := r.Cache[nn]
	if !ok {
		// If we don't have an entry for this EndpointSlice then save it and notify
		// about the change.
		r.Cache[nn] = adminapi.AddressesFromEndpointSlice(endpoints)
		r.notify()
		return ctrl.Result{}, nil
	}

	// We do have an entry for this EndpointSlice.
	// Let's check if it's the same that we're already aware of...
	addresses := adminapi.AddressesFromEndpointSlice(endpoints)
	if cached.Equal(addresses) {
		// No change, don't notify
		return ctrl.Result{}, nil
	}

	// ... it's not the same. Store it and notify.
	r.Cache[nn] = addresses
	r.notify()

	return ctrl.Result{}, nil
}

func (r *KongAdminAPIServiceReconciler) notify() {
	addresses := addressesFromAddressesMap(r.Cache)

	r.Log.V(util.DebugLevel).
		Info("notifying about newly detected Admin API addresses", "addresses", addresses)
	r.EndpointsNotifier.Notify(addresses)
}

func addressesFromAddressesMap(cache CacheT) []string {
	addresses := []string{}
	for _, v := range cache {
		addresses = append(addresses, v.UnsortedList()...)
	}
	return addresses
}
