package gateway

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	k8sobj "github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
)

// -----------------------------------------------------------------------------
// UDPRoute Controller - UDPRouteReconciler
// -----------------------------------------------------------------------------

// UDPRouteReconciler reconciles an UDPRoute object.
type UDPRouteReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	DataplaneClient  controllers.DataPlane
	CacheSyncTimeout time.Duration
	StatusQueue      *status.Queue

	// If GatewayNN is set,
	// only resources managed by the specified Gateway are reconciled.
	GatewayNN controllers.OptionalNamespacedName
}

// SetupWithManager sets up the controller with the Manager.
func (r *UDPRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	blder := ctrl.NewControllerManagedBy(mgr).
		Named("udproute-controller").
		WithOptions(controller.Options{
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
		}).
		// if a GatewayClass updates then we need to enqueue the linked UDPRoutes to
		// ensure that any route objects that may have been orphaned by that change get
		// removed from data-plane configurations, and any routes that are now supported
		// due to that change get added to data-plane configurations.
		Watches(&gatewayapi.GatewayClass{},
			handler.EnqueueRequestsFromMapFunc(r.listUDPRoutesForGatewayClass),
			builder.WithPredicates(predicate.Funcs{
				GenericFunc: func(_ event.GenericEvent) bool { return false }, // we don't need to enqueue from generic
				CreateFunc:  func(e event.CreateEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
				UpdateFunc:  func(e event.UpdateEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
				DeleteFunc:  func(e event.DeleteEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
			}),
		).
		// if a Gateway updates then we need to enqueue the linked UDPRoutes to
		// ensure that any route objects that may have been orphaned by that change get
		// removed from data-plane configurations, and any routes that are now supported
		// due to that change get added to data-plane configurations.
		Watches(&gatewayapi.Gateway{},
			handler.EnqueueRequestsFromMapFunc(r.listUDPRoutesForGateway),
		)

	if r.StatusQueue != nil {
		blder.WatchesRawSource(
			source.Channel(
				r.StatusQueue.Subscribe(schema.GroupVersionKind{
					Group:   gatewayv1alpha2.GroupVersion.Group,
					Version: gatewayv1alpha2.GroupVersion.Version,
					Kind:    "UDPRoute",
				}),
				&handler.EnqueueRequestForObject{},
			),
		)
	}

	// We enqueue only routes that are:
	// - attached during creation or deletion
	// - have been attached or detached to a reconciled Gateway.
	// This allows us to drop the backend data-plane config for a route if
	// it somehow becomes disconnected from a supported Gateway and GatewayClass.
	return blder.
		For(&gatewayapi.UDPRoute{},
			builder.WithPredicates(
				IsRouteAttachedToReconciledGatewayPredicate[*gatewayapi.UDPRoute](r.Client, mgr.GetLogger(), r.GatewayNN),
			),
		).
		Complete(r)
}

// -----------------------------------------------------------------------------
// UDPRoute Controller - Event Handlers
// -----------------------------------------------------------------------------

// listUDPRoutesForGatewayClass is a controller-runtime event.Handler which
// produces a list of UDPRoutes which were bound to a Gateway which is or was
// bound to this GatewayClass. This implementation effectively does a map-reduce
// to determine the UDPRoutes as the relationship has to be discovered entirely
// by object reference. This relies heavily on the inherent performance benefits of
// the cached manager client to avoid API overhead.
func (r *UDPRouteReconciler) listUDPRoutesForGatewayClass(ctx context.Context, obj client.Object) []reconcile.Request {
	// verify that the object is a GatewayClass
	gwc, ok := obj.(*gatewayapi.GatewayClass)
	if !ok {
		r.Log.Error(fmt.Errorf("invalid type"), "Found invalid type in event handlers", "expected", "GatewayClass", "found", reflect.TypeOf(obj))
		return nil
	}

	// map all Gateway objects
	gatewayList := gatewayapi.GatewayList{}
	if err := r.Client.List(ctx, &gatewayList); err != nil {
		r.Log.Error(err, "Failed to list gateway objects from the cached client")
		return nil
	}

	// reduce for in-class Gateway objects
	gateways := make(map[string]map[string]struct{})
	for _, gateway := range gatewayList.Items {
		if string(gateway.Spec.GatewayClassName) == gwc.Name {
			// If the flag `--gateway-to-reconcile` is set, KIC will only reconcile the specified gateway.
			// https://github.com/Kong/kubernetes-ingress-controller/issues/5322
			if !r.GatewayNN.Matches(&gateway) {
				continue
			}

			_, ok := gateways[gateway.Namespace]
			if !ok {
				gateways[gateway.Namespace] = make(map[string]struct{})
			}
			gateways[gateway.Namespace][gateway.Name] = struct{}{}
		}
	}

	// if there are no Gateways associated with this GatewayClass we can stop
	if len(gateways) == 0 {
		return nil
	}

	// map all UDPRoute objects
	udprouteList := gatewayapi.UDPRouteList{}
	if err := r.Client.List(ctx, &udprouteList); err != nil {
		r.Log.Error(err, "Failed to list udproute objects from the cached client")
		return nil
	}

	// reduce for UDPRoute objects bound to an in-class Gateway
	queue := make([]reconcile.Request, 0)
	for _, udproute := range udprouteList.Items {
		// check the udproute's parentRefs
		for _, parentRef := range udproute.Spec.ParentRefs {
			// determine what namespace the parent gateway is in
			namespace := udproute.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}

			// if the gateway matches one of our previously filtered gateways, enqueue the route
			if gatewaysForNamespace, ok := gateways[namespace]; ok {
				if _, ok := gatewaysForNamespace[string(parentRef.Name)]; ok {
					queue = append(queue, reconcile.Request{
						NamespacedName: k8stypes.NamespacedName{
							Namespace: udproute.Namespace,
							Name:      udproute.Name,
						},
					})
				}
			}
		}
	}

	return queue
}

// listUDPRoutesForGateway is a controller-runtime event.Handler which enqueues UDPRoute
// objects for changes to Gateway objects. The relationship between UDPRoutes and their
// Gateways (by way of .Spec.ParentRefs) must be discovered by object relation, so this
// implementation effectively does a map reduce to determine inclusion. This relies heavily
// on the inherent performance benefits of the cached manager client to avoid API overhead.
//
// NOTE:
// due to a race condition where a Gateway and a GatewayClass may be updated at the
// same time and could cause a changed Gateway object to look like it wasn't in-class
// while in reality it may still have active data-plane configurations because it was
// recently in-class, we can't reliably filter Gateway objects based on class as we
// can't verify that didn't change since we received the object. As such the current
// implementation enqueues ALL UDPRoute objects for reconciliation every time a Gateway
// changes. This is not ideal, but after communicating with other members of the
// community this appears to be a standard approach across multiple implementations at
// the moment for v1alpha2. As future releases of Gateway come out we'll need to
// continue iterating on this and perhaps advocating for upstream changes to help avoid
// this kind of problem without having to enqueue extra objects.
func (r *UDPRouteReconciler) listUDPRoutesForGateway(ctx context.Context, obj client.Object) []reconcile.Request {
	// verify that the object is a Gateway
	gw, ok := obj.(*gatewayapi.Gateway)
	if !ok {
		r.Log.Error(fmt.Errorf("invalid type"), "Found invalid type in event handlers", "expected", "Gateway", "found", reflect.TypeOf(obj))
		return nil
	}

	// If the flag `--gateway-to-reconcile` is set, KIC will only reconcile the specified gateway.
	// https://github.com/Kong/kubernetes-ingress-controller/issues/5322
	if !r.GatewayNN.Matches(gw) {
		return nil
	}

	// map all UDPRoute objects
	udprouteList := gatewayapi.UDPRouteList{}
	if err := r.Client.List(ctx, &udprouteList); err != nil {
		r.Log.Error(err, "Failed to list udproute objects from the cached client")
		return nil
	}

	// reduce for UDPRoute objects bound to the Gateway
	queue := make([]reconcile.Request, 0)
	for _, udproute := range udprouteList.Items {
		for _, parentRef := range udproute.Spec.ParentRefs {
			namespace := udproute.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}
			if namespace == gw.Namespace && string(parentRef.Name) == gw.Name {
				queue = append(queue, reconcile.Request{
					NamespacedName: k8stypes.NamespacedName{
						Namespace: udproute.Namespace,
						Name:      udproute.Name,
					},
				})
			}
		}
	}

	return queue
}

// -----------------------------------------------------------------------------
// UDPRoute Controller - Reconciliation
// -----------------------------------------------------------------------------

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=udproutes,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=udproutes/status,verbs=get;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *UDPRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("GatewayV1Alpha2UDPRoute", req.NamespacedName)

	udproute := new(gatewayapi.UDPRoute)
	if err := r.Get(ctx, req.NamespacedName, udproute); err != nil {
		// if the queued object is no longer present in the proxy cache we need
		// to ensure that if it was ever added to the cache, it gets removed.
		if apierrors.IsNotFound(err) {
			debug(log, udproute, "Object does not exist, ensuring it is not present in the proxy cache")
			udproute.Namespace = req.Namespace
			udproute.Name = req.Name
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(udproute)
		}

		// for any error other than 404, requeue
		return ctrl.Result{}, err
	}

	debug(log, udproute, "Processing udproute")

	// if there's a present deletion timestamp then we need to update the proxy cache
	// to drop all relevant routes from its configuration, regardless of whether or
	// not we can find a valid gateway as that gateway may now be deleted but we still
	// need to ensure removal of the data-plane configuration.
	debug(log, udproute, "Checking deletion timestamp")
	if udproute.DeletionTimestamp != nil {
		debug(log, udproute, "udproute is being deleted, re-configuring data-plane")
		if err := r.DataplaneClient.DeleteObject(udproute); err != nil {
			debug(log, udproute, "Failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, udproute, "Ensured object was removed from the data-plane (if ever present)")
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(udproute)
	}

	// we need to pull the Gateway parent objects for the UDPRoute to verify
	// routing behavior and ensure compatibility with Gateway configurations.
	debug(log, udproute, "Retrieving GatewayClass and Gateway for route")
	gateways, err := getSupportedGatewayForRoute(ctx, log, r.Client, udproute, r.GatewayNN)
	if err != nil {
		if errors.Is(err, ErrNoSupportedGateway) {
			// if there's no supported Gateway then this route could have been previously
			// supported by this controller. As such we ensure that no supported Gateway
			// references exist in the object status any longer.
			_, err := ensureGatewayReferenceStatusRemoved(ctx, r.Client, log, udproute)
			if err != nil {
				// some failure happened so we need to retry to avoid orphaned statuses
				return ctrl.Result{}, err
			}

			// if the route doesn't have a supported Gateway+GatewayClass associated with
			// it it's possible it became orphaned after becoming queued. In either case
			// ensure that it's removed from the proxy cache to avoid orphaned data-plane
			// configurations.
			debug(log, udproute, "Ensuring that dataplane is updated to remove unsupported route (if applicable)")
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(udproute)
		}
		return ctrl.Result{}, err
	}

	// the referenced gateway object(s) for the UDPRoute needs to be ready
	// before we'll attempt any configurations of it. If it's not we'll
	// requeue the object and wait until all supported gateways are ready.
	debug(log, udproute, "Checking if the udproute's gateways are ready")
	for _, gateway := range gateways {
		if !isGatewayProgrammed(gateway.gateway) {
			debug(log, udproute, "Gateway for route was not ready, waiting")
			return ctrl.Result{Requeue: true}, nil
		}
	}

	if isRouteAccepted(gateways) {
		// if the gateways are ready, and the UDPRoute is destined for them, ensure that
		// the object is pushed to the dataplane.
		if err := r.DataplaneClient.UpdateObject(udproute); err != nil {
			debug(log, udproute, "Failed to update object in data-plane, requeueing")
			return ctrl.Result{}, err
		}
		if r.DataplaneClient.AreKubernetesObjectReportsEnabled() {
			// if the dataplane client has reporting enabled (this is the default and is
			// tied in with status updates being enabled in the controller manager) then
			// we will wait until the object is reported as successfully configured before
			// moving on to status updates.
			if !r.DataplaneClient.KubernetesObjectIsConfigured(udproute) {
				return ctrl.Result{Requeue: true}, nil
			}
		}
	} else {
		// route is not accepted, remove it from kong store
		if err := r.DataplaneClient.DeleteObject(udproute); err != nil {
			debug(log, udproute, "Failed to delete object in data-plane, requeueing")
			return ctrl.Result{}, err
		}
	}

	// now that the object has been successfully configured for in the dataplane
	// we can update the object status to indicate that it's now properly linked
	// to the configured Gateways.
	debug(log, udproute, "Ensuring status contains Gateway associations")
	updated, res, err := r.ensureGatewayReferenceStatusAdded(ctx, udproute, gateways...)
	if err != nil {
		// don't proceed until the statuses can be updated appropriately
		return ctrl.Result{}, err
	}
	if !res.IsZero() {
		return res, nil
	}
	if updated {
		// if the status was updated it will trigger a follow-up reconciliation
		return ctrl.Result{}, nil
	}

	if r.DataplaneClient.AreKubernetesObjectReportsEnabled() {
		// if the dataplane client has reporting enabled (this is the default and is
		// tied in with status updates being enabled in the controller manager) then
		// we will wait until the object is reported as successfully configured before
		// moving on to status updates.
		configurationStatus := r.DataplaneClient.KubernetesObjectConfigurationStatus(udproute)
		if configurationStatus == k8sobj.ConfigurationStatusUnknown {
			// requeue until udproute is configured.
			debug(log, udproute, "UDPRoute not configured, requeueing")
			return ctrl.Result{Requeue: true}, nil
		}

		if configurationStatus == k8sobj.ConfigurationStatusFailed {
			debug(log, udproute, "UDPRoute configuration failed")
			statusUpdated, err := ensureParentsProgrammedCondition(ctx, r.Status(), udproute, udproute.Status.Parents, gateways, metav1.Condition{
				Status: metav1.ConditionFalse,
				Reason: string(ConditionReasonTranslationError),
			})
			if err != nil {
				// don't proceed until the statuses can be updated appropriately
				debug(log, udproute, "Failed to update programmed condition")
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: !statusUpdated}, nil
		}

		statusUpdated, err := ensureParentsProgrammedCondition(ctx, r.Status(), udproute, udproute.Status.Parents, gateways, metav1.Condition{
			Status: metav1.ConditionTrue,
			Reason: string(ConditionReasonConfiguredInGateway),
		})
		if err != nil {
			// don't proceed until the statuses can be updated appropriately
			debug(log, udproute, "Failed to update programmed condition")
			return ctrl.Result{}, err
		}
		if statusUpdated {
			// if the status was updated it will trigger a follow-up reconciliation
			// so we don't need to do anything further here.
			debug(log, udproute, "Programmed condition updated")
			return ctrl.Result{}, nil
		}
	}

	// once the data-plane has accepted the UDPRoute object, we're all set.
	info(log, udproute, "UDPRoute has been configured on the data-plane")
	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// UDPRouteReconciler - Status Helpers
// -----------------------------------------------------------------------------

// ensureGatewayReferenceStatus takes any number of Gateways that should be
// considered "attached" to a given UDPRoute and ensures that the status
// for the UDPRoute is updated appropriately.
// It returns true if controller should requeue the object. Either because
// the status update resulted in a conflict or because the status was updated.
func (r *UDPRouteReconciler) ensureGatewayReferenceStatusAdded(
	ctx context.Context, udproute *gatewayapi.UDPRoute, gateways ...supportedGatewayWithCondition,
) (bool, ctrl.Result, error) {
	parentStatuses, statusChangesWereMade := parentStatusesForRoute(
		udproute,
		udproute.Status.Parents,
		gateways...,
	)

	programmedConditionChanged := initializeParentStatusesWithProgrammedCondition(udproute, parentStatuses)

	// if we didn't have to actually make any changes, no status update is needed
	if !statusChangesWereMade && !programmedConditionChanged {
		return false, ctrl.Result{}, nil
	}

	// update the udproute status with the new status references
	udproute.Status.Parents = make([]gatewayapi.RouteParentStatus, 0, len(parentStatuses))
	for _, parent := range parentStatuses {
		udproute.Status.Parents = append(udproute.Status.Parents, *parent)
	}

	// update the object status in the API
	res, err := handleUpdateError(r.Status().Update(ctx, udproute), r.Log, udproute)
	if err != nil {
		return false, ctrl.Result{}, err
	}
	if !res.IsZero() {
		return false, res, nil
	}

	// the status needed an update and it was updated successfully
	return true, ctrl.Result{}, nil
}
