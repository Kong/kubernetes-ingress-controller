package gateway

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
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
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	k8sobj "github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
)

// -----------------------------------------------------------------------------
// GRPCRoute Controller - GRPCRouteReconciler
// -----------------------------------------------------------------------------

// GRPCRouteReconciler reconciles an GRPCRoute object.
type GRPCRouteReconciler struct {
	client.Client

	Log             logr.Logger
	Scheme          *runtime.Scheme
	DataplaneClient controllers.DataPlane
	StatusQueue     *status.Queue
	// If EnableReferenceGrant is true, we will check for ReferenceGrant if backend in another
	// namespace is in backendRefs.
	// If it is false, referencing backend in different namespace will be rejected.
	EnableReferenceGrant bool
	CacheSyncTimeout     time.Duration

	// If GatewayNN is set,
	// only resources managed by the specified Gateway are reconciled.
	GatewayNN controllers.OptionalNamespacedName
}

// SetupWithManager sets up the controller with the Manager.
func (r *GRPCRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	blder := ctrl.NewControllerManagedBy(mgr).
		// set the controller name
		Named("grpcroute-controller").
		WithOptions(controller.Options{
			LogConstructor: func(_ *reconcile.Request) logr.Logger {
				return r.Log
			},
			CacheSyncTimeout: r.CacheSyncTimeout,
		}).
		// if a GatewayClass updates then we need to enqueue the linked GRPCRoutes to
		// ensure that any route objects that may have been orphaned by that change get
		// removed from data-plane configurations, and any routes that are now supported
		// due to that change get added to data-plane configurations.
		Watches(&gatewayapi.GatewayClass{},
			handler.EnqueueRequestsFromMapFunc(r.listGRPCRoutesForGatewayClass),
			builder.WithPredicates(predicate.Funcs{
				GenericFunc: func(_ event.GenericEvent) bool { return false }, // we don't need to enqueue from generic
				CreateFunc:  func(e event.CreateEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
				UpdateFunc:  func(e event.UpdateEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
				DeleteFunc:  func(e event.DeleteEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
			}),
		).
		// if a Gateway updates then we need to enqueue the linked GRPCRoutes to
		// ensure that any route objects that may have been orphaned by that change get
		// removed from data-plane configurations, and any routes that are now supported
		// due to that change get added to data-plane configurations.
		Watches(&gatewayapi.Gateway{},
			handler.EnqueueRequestsFromMapFunc(r.listGRPCRoutesForGateway),
		)

	if r.StatusQueue != nil {
		blder.WatchesRawSource(
			source.Channel(
				r.StatusQueue.Subscribe(schema.GroupVersionKind{
					Group:   gatewayv1alpha2.GroupVersion.Group,
					Version: gatewayv1alpha2.GroupVersion.Version,
					Kind:    "GRPCRoute",
				}),
				&handler.EnqueueRequestForObject{},
			),
		)
	}

	// because of the additional burden of having to manage reference data-plane
	// configurations for GRPCRoute objects in the underlying Kong Gateway, we
	// simply reconcile ALL GRPCRoute objects. This allows us to drop the backend
	// data-plane config for an GRPCRoute if it somehow becomes disconnected from
	// a supported Gateway and GatewayClass.
	return blder.For(&gatewayapi.GRPCRoute{}).
		Complete(r)
}

// -----------------------------------------------------------------------------
// GRPCRoute Controller - Event Handlers
// -----------------------------------------------------------------------------

// listGRPCRoutesForGatewayClass is a controller-runtime event.Handler which
// produces a list of GRPCRoutes which were bound to a Gateway which is or was
// bound to this GatewayClass. This implementation effectively does a map-reduce
// to determine the GRPCRoutes as the relationship has to be discovered entirely
// by object reference. This relies heavily on the inherent performance benefits of
// the cached manager client to avoid API overhead.
func (r *GRPCRouteReconciler) listGRPCRoutesForGatewayClass(ctx context.Context, obj client.Object) []reconcile.Request {
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
			if gatewayToReconcile, ok := r.GatewayNN.Get(); ok {
				if gatewayToReconcile.Namespace != gateway.Namespace || gatewayToReconcile.Name != gateway.Name {
					continue
				}
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

	// map all GRPCRoute objects
	grpcrouteList := gatewayapi.GRPCRouteList{}
	if err := r.Client.List(ctx, &grpcrouteList); err != nil {
		r.Log.Error(err, "Failed to list grpcroute objects from the cached client")
		return nil
	}

	// reduce for GRPCRoute objects bound to an in-class Gateway
	queue := make([]reconcile.Request, 0)
	for _, grpcroute := range grpcrouteList.Items {
		// check the grpcroute's parentRefs
		for _, parentRef := range grpcroute.Spec.ParentRefs {
			// determine what namespace the parent gateway is in
			namespace := grpcroute.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}

			// if the gateway matches one of our previously filtered gateways, enqueue the route
			if gatewaysForNamespace, ok := gateways[namespace]; ok {
				if _, ok := gatewaysForNamespace[string(parentRef.Name)]; ok {
					queue = append(queue, reconcile.Request{
						NamespacedName: k8stypes.NamespacedName{
							Namespace: grpcroute.Namespace,
							Name:      grpcroute.Name,
						},
					})
				}
			}
		}
	}

	return queue
}

// listGRPCRoutesForGateway is a controller-runtime event.Handler which enqueues GRPCRoute
// objects for changes to Gateway objects. The relationship between GRPCRoutes and their
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
// implementation enqueues ALL GRPCRoute objects for reconciliation every time a Gateway
// changes. This is not ideal, but after communicating with other members of the
// community this appears to be a standard approach across multiple implementations at
// the moment for v1alpha2. As future releases of Gateway come out we'll need to
// continue iterating on this and perhaps advocating for upstream changes to help avoid
// this kind of problem without having to enqueue extra objects.
func (r *GRPCRouteReconciler) listGRPCRoutesForGateway(ctx context.Context, obj client.Object) []reconcile.Request {
	// verify that the object is a Gateway
	gw, ok := obj.(*gatewayapi.Gateway)
	if !ok {
		r.Log.Error(fmt.Errorf("invalid type"), "Found invalid type in event handlers", "expected", "Gateway", "found", reflect.TypeOf(obj))
		return nil
	}

	// If the flag `--gateway-to-reconcile` is set, KIC will only reconcile the specified gateway.
	// https://github.com/Kong/kubernetes-ingress-controller/issues/5322
	if gatewayToReconcile, ok := r.GatewayNN.Get(); ok {
		if gatewayToReconcile.Namespace != gw.Namespace || gatewayToReconcile.Name != gw.Name {
			return nil
		}
	}

	// map all GRPCRoute objects
	grpcrouteList := gatewayapi.GRPCRouteList{}
	if err := r.Client.List(ctx, &grpcrouteList); err != nil {
		r.Log.Error(err, "Failed to list grpcroute objects from the cached client")
		return nil
	}

	// reduce for GRPCRoute objects bound to the Gateway
	queue := make([]reconcile.Request, 0)
	for _, grpcroute := range grpcrouteList.Items {
		for _, parentRef := range grpcroute.Spec.ParentRefs {
			namespace := grpcroute.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}
			if namespace == gw.Namespace && string(parentRef.Name) == gw.Name {
				queue = append(queue, reconcile.Request{
					NamespacedName: k8stypes.NamespacedName{
						Namespace: grpcroute.Namespace,
						Name:      grpcroute.Name,
					},
				})
			}
		}
	}

	return queue
}

// -----------------------------------------------------------------------------
// GRPCRoute Controller - Reconciliation
// -----------------------------------------------------------------------------

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=grpcroutes,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=grpcroutes/status,verbs=get;patch;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *GRPCRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("GatewayV1Alpha2GRPCRoute", req.NamespacedName)

	grpcroute := new(gatewayapi.GRPCRoute)
	if err := r.Get(ctx, req.NamespacedName, grpcroute); err != nil {
		// if the queued object is no longer present in the proxy cache we need
		// to ensure that if it was ever added to the cache, it gets removed.
		if apierrors.IsNotFound(err) {
			debug(log, grpcroute, "Object does not exist, ensuring it is not present in the proxy cache")
			grpcroute.Namespace = req.Namespace
			grpcroute.Name = req.Name
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(grpcroute)
		}

		// for any error other than 404, requeue
		return ctrl.Result{}, err
	}

	debug(log, grpcroute, "Processing grpcroute")

	// if there's a present deletion timestamp then we need to update the proxy cache
	// to drop all relevant routes from its configuration, regardless of whether or
	// not we can find a valid gateway as that gateway may now be deleted but we still
	// need to ensure removal of the data-plane configuration.
	debug(log, grpcroute, "Checking deletion timestamp")
	if grpcroute.DeletionTimestamp != nil {
		debug(log, grpcroute, "grpcroute is being deleted, re-configuring data-plane")
		if err := r.DataplaneClient.DeleteObject(grpcroute); err != nil {
			debug(log, grpcroute, "Failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, grpcroute, "Ensured object was removed from the data-plane (if ever present)")
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(grpcroute)
	}

	// we need to pull the Gateway parent objects for the grpcroute to verify
	// routing behavior and ensure compatibility with Gateway configurations.
	debug(log, grpcroute, "Retrieving GatewayClass and Gateway for route")
	gateways, err := getSupportedGatewayForRoute(ctx, log, r.Client, grpcroute, r.GatewayNN)
	if err != nil {
		if err.Error() == unsupportedGW {
			debug(log, grpcroute, "Unsupported route found, processing to verify whether it was ever supported")
			// if there's no supported Gateway then this route could have been previously
			// supported by this controller. As such we ensure that no supported Gateway
			// references exist in the object status any longer.
			statusUpdated, err := r.ensureGatewayReferenceStatusRemoved(ctx, grpcroute)
			if err != nil {
				// some failure happened so we need to retry to avoid orphaned statuses
				return ctrl.Result{}, err
			}
			if statusUpdated {
				// the status did in fact needed to be updated, so no need to requeue
				// as the status update will trigger a requeue.
				debug(log, grpcroute, "Unsupported route was previously supported, status was updated")
				return ctrl.Result{}, nil
			}

			// if the route doesn't have a supported Gateway+GatewayClass associated with
			// it it's possible it became orphaned after becoming queued. In either case
			// ensure that it's removed from the proxy cache to avoid orphaned data-plane
			// configurations.
			debug(log, grpcroute, "Ensuring that dataplane is updated to remove unsupported route (if applicable)")
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(grpcroute)
		}
		return ctrl.Result{}, err
	}

	// the referenced gateway object(s) for the grpcroute needs to be ready
	// before we'll attempt any configurations of it. If it's not we'll
	// requeue the object and wait until all supported gateways are ready.
	debug(log, grpcroute, "Checking if the grpcroute's gateways are ready")
	for _, gateway := range gateways {
		if !isGatewayProgrammed(gateway.gateway) {
			debug(log, grpcroute, "Gateway for route was not ready, waiting")
			return ctrl.Result{Requeue: true}, nil
		}
	}

	if isRouteAccepted(gateways) {
		// if the gateways are ready, and the GRPCRoute is destined for them, ensure that
		// the object is pushed to the dataplane.
		if err := r.DataplaneClient.UpdateObject(grpcroute); err != nil {
			debug(log, grpcroute, "Failed to update object in data-plane, requeueing")
			return ctrl.Result{}, err
		}
	} else {
		// route is not accepted, remove it from kong store
		if err := r.DataplaneClient.DeleteObject(grpcroute); err != nil {
			debug(log, grpcroute, "Failed to delete object in data-plane, requeueing")
			return ctrl.Result{}, err
		}
	}

	// now that the object has been successfully configured for in the dataplane
	// we can update the object status to indicate that it's now properly linked
	// to the configured Gateways.
	debug(log, grpcroute, "Ensuring status contains Gateway associations")
	statusUpdated, err := r.ensureGatewayReferenceStatusAdded(ctx, grpcroute, gateways...)
	if err != nil {
		// don't proceed until the statuses can be updated appropriately
		return ctrl.Result{}, err
	}
	if statusUpdated {
		// if the status was updated it will trigger a follow-up reconciliation
		// so we don't need to do anything further here.
		return ctrl.Result{}, nil
	}

	// update "Programmed" condition if the GRPCRoute is translated to Kong configuration.
	// if the GRPCRoute is not configured in the dataplane, leave it unchanged and requeue.
	// if it is successfully configured, update its "Programmed" condition to True.
	// if translation failure happens, update its "Programmed" condition to False.
	if r.DataplaneClient.AreKubernetesObjectReportsEnabled() {
		// if the dataplane client has reporting enabled (this is the default and is
		// tied in with status updates being enabled in the controller manager) then
		// we will wait until the object is reported as successfully configured before
		// moving on to status updates.
		configurationStatus := r.DataplaneClient.KubernetesObjectConfigurationStatus(grpcroute)
		if configurationStatus == k8sobj.ConfigurationStatusUnknown {
			// requeue until grpcroute is configured.
			debug(log, grpcroute, "GRPCRoute not configured, requeueing")
			return ctrl.Result{Requeue: true}, nil
		}

		if configurationStatus == k8sobj.ConfigurationStatusFailed {
			debug(log, grpcroute, "GRPCRoute configuration failed")
			statusUpdated, err := ensureParentsProgrammedCondition(ctx, r.Status(), grpcroute, grpcroute.Status.Parents, gateways, metav1.Condition{
				Status: metav1.ConditionFalse,
				Reason: string(ConditionReasonTranslationError),
			})
			if err != nil {
				// don't proceed until the statuses can be updated appropriately
				debug(log, grpcroute, "Failed to update programmed condition")
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: !statusUpdated}, nil
		}

		statusUpdated, err := ensureParentsProgrammedCondition(ctx, r.Status(), grpcroute, grpcroute.Status.Parents, gateways, metav1.Condition{
			Status: metav1.ConditionTrue,
			Reason: string(ConditionReasonConfiguredInGateway),
		})
		if err != nil {
			// don't proceed until the statuses can be updated appropriately
			debug(log, grpcroute, "Failed to update programmed condition")
			return ctrl.Result{}, err
		}
		if statusUpdated {
			// if the status was updated it will trigger a follow-up reconciliation
			// so we don't need to do anything further here.
			debug(log, grpcroute, "Programmed condition updated")
			return ctrl.Result{}, nil
		}
	}

	// once the data-plane has accepted the GRPCRoute object, we're all set.
	info(log, grpcroute, "GRPCRoute has been configured on the data-plane")

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// GRPCRouteReconciler - Status Helpers
// -----------------------------------------------------------------------------

// grpcrouteParentKind indicates the only object KIND that this GRPCRoute
// implementation supports for route object parent references.
var grpcrouteParentKind = "Gateway"

// ensureGatewayReferenceStatus takes any number of Gateways that should be
// considered "attached" to a given GRPCRoute and ensures that the status
// for the GRPCRoute is updated appropriately.
func (r *GRPCRouteReconciler) ensureGatewayReferenceStatusAdded(ctx context.Context, grpcroute *gatewayapi.GRPCRoute, gateways ...supportedGatewayWithCondition) (bool, error) {
	// map the existing parentStatues to avoid duplications
	parentStatuses := getParentStatuses(grpcroute, grpcroute.Status.Parents)

	// overlay the parent ref statuses for all new gateway references
	statusChangesWereMade := false
	for _, gateway := range gateways {
		gateway := gateway
		// build a new status for the parent Gateway
		gatewayParentStatus := &gatewayapi.RouteParentStatus{
			ParentRef: gatewayapi.ParentReference{
				Group:     (*gatewayapi.Group)(&gatewayv1alpha2.GroupVersion.Group),
				Kind:      util.StringToGatewayAPIKindPtr(grpcrouteParentKind),
				Namespace: (*gatewayapi.Namespace)(&gateway.gateway.Namespace),
				Name:      gatewayapi.ObjectName(gateway.gateway.Name),
			},
			ControllerName: GetControllerName(),
			Conditions: []metav1.Condition{{
				Type:               string(gatewayapi.RouteConditionAccepted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: grpcroute.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayapi.RouteReasonAccepted),
			}},
		}

		if gateway.listenerName != "" {
			sectionName := gatewayapi.SectionName(gateway.listenerName)
			gatewayParentStatus.ParentRef.SectionName = &sectionName
		}

		// if the reference already exists and doesn't require any changes
		// then just leave it alone.
		parentRefKey := gateway.gateway.Namespace + "/" + gateway.gateway.Name
		if existingGatewayParentStatus, exists := parentStatuses[parentRefKey]; exists {
			//  check if the parentRef and controllerName are equal, and whether the new condition is present in existing conditions
			if reflect.DeepEqual(existingGatewayParentStatus.ParentRef, gatewayParentStatus.ParentRef) &&
				existingGatewayParentStatus.ControllerName == gatewayParentStatus.ControllerName &&
				lo.ContainsBy(existingGatewayParentStatus.Conditions, func(condition metav1.Condition) bool {
					return sameCondition(gatewayParentStatus.Conditions[0], condition)
				}) {
				continue
			}
		}

		// otherwise overlay the new status on top the list of parentStatuses
		parentStatuses[parentRefKey] = gatewayParentStatus
		statusChangesWereMade = true
	}

	// initialize "programmed" condition to Unknown.
	// do not update the condition If a "Programmed" condition is already present.
	programmedConditionChanged := false
	programmedConditionUnknown := metav1.Condition{
		Type:               ConditionTypeProgrammed,
		Status:             metav1.ConditionUnknown,
		Reason:             string(ConditionReasonProgrammedUnknown),
		ObservedGeneration: grpcroute.Generation,
		LastTransitionTime: metav1.Now(),
	}
	for _, parentStatus := range parentStatuses {
		if !parentStatusHasProgrammedCondition(parentStatus) {
			programmedConditionChanged = true
			parentStatus.Conditions = append(parentStatus.Conditions, programmedConditionUnknown)
		}
	}

	// if we didn't have to actually make any changes, no status update is needed
	if !statusChangesWereMade && !programmedConditionChanged {
		return false, nil
	}

	// update the grpcroute status with the new status references
	grpcroute.Status.Parents = make([]gatewayapi.RouteParentStatus, 0, len(parentStatuses))
	for _, parent := range parentStatuses {
		grpcroute.Status.Parents = append(grpcroute.Status.Parents, *parent)
	}

	// update the object status in the API
	if err := r.Status().Update(ctx, grpcroute); err != nil {
		return false, err
	}

	// the status needed an update and it was updated successfully
	return true, nil
}

// ensureGatewayReferenceStatusRemoved uses the ControllerName provided by the Gateway
// implementation to prune status references to Gateways supported by this controller
// in the provided GRPCRoute object.
func (r *GRPCRouteReconciler) ensureGatewayReferenceStatusRemoved(ctx context.Context, grpcroute *gatewayapi.GRPCRoute) (bool, error) {
	// drop all status references to supported Gateway objects
	newStatuses := make([]gatewayapi.RouteParentStatus, 0)
	for _, status := range grpcroute.Status.Parents {
		if status.ControllerName != GetControllerName() {
			newStatuses = append(newStatuses, status)
		}
	}

	// if the new list of statuses is the same length as the old
	// nothing has changed and we're all done.
	if len(newStatuses) == len(grpcroute.Status.Parents) {
		return false, nil
	}

	// update the object status in the API
	grpcroute.Status.Parents = newStatuses
	if err := r.Status().Update(ctx, grpcroute); err != nil {
		return false, err
	}

	// the status needed an update and it was updated successfully
	return true, nil
}
