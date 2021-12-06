package gateway

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/proxy"
)

// -----------------------------------------------------------------------------
// HTTPRoute Controller - HTTPRouteReconciler
// -----------------------------------------------------------------------------

// HTTPRouteReconciler reconciles an HTTPRoute object
type HTTPRouteReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
	Proxy  proxy.Proxy
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("httproute-controller", mgr, controller.Options{
		Reconciler: r,
		Log:        r.Log,
	})
	if err != nil {
		return err
	}

	// TODO: add watches for changes to referenced GatewayClasses and Gateways.
	//       See: https://github.com/Kong/kubernetes-ingress-controller/issues/2077

	// because of the additional burden of having to manage reference data-plane
	// configurations for HTTPRoute objects in the underlying Kong Gateway, we
	// simply reconcile ALL HTTPRoute objects. This allows us to drop the backend
	// data-plane config for an HTTPRoute if it somehow becomes disconnected from
	// a supported Gateway and GatewayClass.
	return c.Watch(
		&source.Kind{Type: &gatewayv1alpha2.HTTPRoute{}},
		&handler.EnqueueRequestForObject{},
	)
}

// -----------------------------------------------------------------------------
// HTTPRoute Controller - Reconciliation
// -----------------------------------------------------------------------------

//+kubebuilder:rbac:groups=networking.k8s.io,resources=httproutes,verbs=get;list;watch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=httproutes/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *HTTPRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("NetV1Alpha2HTTPRoute", req.NamespacedName)

	httproute := new(gatewayv1alpha2.HTTPRoute)
	if err := r.Get(ctx, req.NamespacedName, httproute); err != nil {
		// if the queued object is no longer present in the proxy cache we need
		// to ensure that if it was ever added to the cache, it gets removed.
		if errors.IsNotFound(err) {
			debug(log, httproute, "object does not exist, ensuring it is not present in the proxy cache")
			httproute.Namespace = req.Namespace
			httproute.Name = req.Name
			return ctrlutils.EnsureProxyDeleteObject(r.Proxy, httproute)
		}

		// for any error other than 404, requeue
		return ctrl.Result{}, err
	}

	// if there's a present deletion timestamp then we need to update the proxy cache
	// to drop all relevant routes from its configuration, regardless of whether or
	// not we can find a valid gateway as that gateway may now be deleted but we still
	// need to ensure removal of the data-plane configuration.
	debug(log, httproute, "checking deletion timestamp")
	if httproute.DeletionTimestamp != nil {
		debug(log, httproute, "httproute is being deleted, re-configuring data-plane")
		if err := r.Proxy.DeleteObject(httproute); err != nil {
			debug(log, httproute, "failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, httproute, "ensured object was removed from the data-plane (if ever present)")
		return ctrl.Result{}, nil
	}

	// we need to pull both the GatewayClass and Gateway parent objects for
	// the HTTPRoute to verify routing behavior and ensure compatibility with
	// Gateway configurations.
	debug(log, httproute, "retrieving GatewayClass and Gateway for route")
	gateway, err := getSupportedGatewayForRoute(ctx, r.Client, httproute)
	if err != nil {
		if err.Error() == unsupportedGW {
			// if the route doesn't have a supported Gateway+GatewayClass associated with
			// it it's possible it became orphaned after becoming queued. In either case
			// ensure that it's removed from the proxy cache to avoid orphaned data-plane
			// configurations.
			return ctrlutils.EnsureProxyDeleteObject(r.Proxy, httproute)
		}
		return ctrl.Result{}, err
	}

	// the referenced gateway object for the HTTPRoute needs to be ready
	// before we'll attempt any configurations of it. If it's not we'll
	// requeue the object and wait until it's ready.
	debug(log, httproute, "checking if the httproute's gateway is ready")
	if !isGatewayReady(gateway) {
		debug(log, httproute, "gateway for route was not ready, waiting")
		return ctrl.Result{Requeue: true}, nil
	}

	// finally if all matching has succeeded and the object is not being deleted,
	// we can configure it in the data-plane.
	debug(log, httproute, "sending httproute information to the data-plane for configuration")
	if err := r.Proxy.UpdateObject(httproute); err != nil {
		debug(log, httproute, "failed to update object in data-plane, requeueing")
		return ctrl.Result{}, err
	}

	// once the data-plane has accepted the HTTPRoute object, we're all set.
	return ctrl.Result{}, nil
}
