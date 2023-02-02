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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	k8sobj "github.com/kong/kubernetes-ingress-controller/v2/internal/util/kubernetes/object"
)

// -----------------------------------------------------------------------------
// UDPRoute Controller - UDPRouteReconciler
// -----------------------------------------------------------------------------

// UDPRouteReconciler reconciles an UDPRoute object.
type UDPRouteReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	DataplaneClient  *dataplane.KongClient
	CacheSyncTimeout time.Duration
}

// SetupWithManager sets up the controller with the Manager.
func (r *UDPRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("udproute-controller", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}

	// if a GatewayClass updates then we need to enqueue the linked UDPRoutes to
	// ensure that any route objects that may have been orphaned by that change get
	// removed from data-plane configurations, and any routes that are now supported
	// due to that change get added to data-plane configurations.
	if err := c.Watch(
		&source.Kind{Type: &gatewayv1beta1.GatewayClass{}},
		handler.EnqueueRequestsFromMapFunc(r.listUDPRoutesForGatewayClass),
		predicate.Funcs{
			GenericFunc: func(e event.GenericEvent) bool { return false }, // we don't need to enqueue from generic
			CreateFunc:  func(e event.CreateEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
			UpdateFunc:  func(e event.UpdateEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
			DeleteFunc:  func(e event.DeleteEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
		},
	); err != nil {
		return err
	}

	// if a Gateway updates then we need to enqueue the linked UDPRoutes to
	// ensure that any route objects that may have been orphaned by that change get
	// removed from data-plane configurations, and any routes that are now supported
	// due to that change get added to data-plane configurations.
	if err := c.Watch(
		&source.Kind{Type: &gatewayv1beta1.Gateway{}},
		handler.EnqueueRequestsFromMapFunc(r.listUDPRoutesForGateway),
	); err != nil {
		return err
	}

	// because of the additional burden of having to manage reference data-plane
	// configurations for UDPRoute objects in the underlying Kong Gateway, we
	// simply reconcile ALL UDPRoute objects. This allows us to drop the backend
	// data-plane config for an UDPRoute if it somehow becomes disconnected from
	// a supported Gateway and GatewayClass.
	return c.Watch(
		&source.Kind{Type: &gatewayv1alpha2.UDPRoute{}},
		&handler.EnqueueRequestForObject{},
	)
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
func (r *UDPRouteReconciler) listUDPRoutesForGatewayClass(obj client.Object) []reconcile.Request {
	// verify that the object is a GatewayClass
	gwc, ok := obj.(*gatewayv1beta1.GatewayClass)
	if !ok {
		r.Log.Error(fmt.Errorf("invalid type"), "found invalid type in event handlers", "expected", "GatewayClass", "found", reflect.TypeOf(obj))
		return nil
	}

	// map all Gateway objects
	gatewayList := gatewayv1beta1.GatewayList{}
	if err := r.Client.List(context.Background(), &gatewayList); err != nil {
		r.Log.Error(err, "failed to list gateway objects from the cached client")
		return nil
	}

	// reduce for in-class Gateway objects
	gateways := make(map[string]map[string]struct{})
	for _, gateway := range gatewayList.Items {
		if string(gateway.Spec.GatewayClassName) == gwc.Name {
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
	udprouteList := gatewayv1alpha2.UDPRouteList{}
	if err := r.Client.List(context.Background(), &udprouteList); err != nil {
		r.Log.Error(err, "failed to list udproute objects from the cached client")
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
						NamespacedName: types.NamespacedName{
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
func (r *UDPRouteReconciler) listUDPRoutesForGateway(obj client.Object) []reconcile.Request {
	// verify that the object is a Gateway
	gw, ok := obj.(*gatewayv1beta1.Gateway)
	if !ok {
		r.Log.Error(fmt.Errorf("invalid type"), "found invalid type in event handlers", "expected", "Gateway", "found", reflect.TypeOf(obj))
		return nil
	}

	// map all UDPRoute objects
	udprouteList := gatewayv1alpha2.UDPRouteList{}
	if err := r.Client.List(context.Background(), &udprouteList); err != nil {
		r.Log.Error(err, "failed to list udproute objects from the cached client")
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
					NamespacedName: types.NamespacedName{
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

	udproute := new(gatewayv1alpha2.UDPRoute)
	if err := r.Get(ctx, req.NamespacedName, udproute); err != nil {
		// if the queued object is no longer present in the proxy cache we need
		// to ensure that if it was ever added to the cache, it gets removed.
		if apierrors.IsNotFound(err) {
			debug(log, udproute, "object does not exist, ensuring it is not present in the proxy cache")
			udproute.Namespace = req.Namespace
			udproute.Name = req.Name
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(udproute)
		}

		// for any error other than 404, requeue
		return ctrl.Result{}, err
	}
	debug(log, udproute, "processing udproute")

	// if there's a present deletion timestamp then we need to update the proxy cache
	// to drop all relevant routes from its configuration, regardless of whether or
	// not we can find a valid gateway as that gateway may now be deleted but we still
	// need to ensure removal of the data-plane configuration.
	debug(log, udproute, "checking deletion timestamp")
	if udproute.DeletionTimestamp != nil {
		debug(log, udproute, "udproute is being deleted, re-configuring data-plane")
		if err := r.DataplaneClient.DeleteObject(udproute); err != nil {
			debug(log, udproute, "failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, udproute, "ensured object was removed from the data-plane (if ever present)")
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(udproute)
	}

	// we need to pull the Gateway parent objects for the UDPRoute to verify
	// routing behavior and ensure compatibility with Gateway configurations.
	debug(log, udproute, "retrieving GatewayClass and Gateway for route")
	gateways, err := getSupportedGatewayForRoute(ctx, r.Client, udproute)
	if err != nil {
		if err.Error() == unsupportedGW {
			debug(log, udproute, "unsupported route found, processing to verify whether it was ever supported")
			// if there's no supported Gateway then this route could have been previously
			// supported by this controller. As such we ensure that no supported Gateway
			// references exist in the object status any longer.
			statusUpdated, err := r.ensureGatewayReferenceStatusRemoved(ctx, udproute)
			if err != nil {
				// some failure happened so we need to retry to avoid orphaned statuses
				return ctrl.Result{}, err
			}
			if statusUpdated {
				// the status did in fact needed to be updated, so no need to requeue
				// as the status update will trigger a requeue.
				debug(log, udproute, "unsupported route was previously supported, status was updated")
				return ctrl.Result{}, nil
			}

			// if the route doesn't have a supported Gateway+GatewayClass associated with
			// it it's possible it became orphaned after becoming queued. In either case
			// ensure that it's removed from the proxy cache to avoid orphaned data-plane
			// configurations.
			debug(log, udproute, "ensuring that dataplane is updated to remove unsupported route (if applicable)")
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(udproute)
		}
		return ctrl.Result{}, err
	}

	// the referenced gateway object(s) for the UDPRoute needs to be ready
	// before we'll attempt any configurations of it. If it's not we'll
	// requeue the object and wait until all supported gateways are ready.
	debug(log, udproute, "checking if the udproute's gateways are ready")
	for _, gateway := range gateways {
		if !isGatewayReady(gateway.gateway) {
			debug(log, udproute, "gateway for route was not ready, waiting")
			return ctrl.Result{Requeue: true}, nil
		}
	}

	// TODO: UDPRoute should be 'Accepted' before proceeding further. Get here back once
	// https://github.com/Kong/kubernetes-ingress-controller/issues/2544 is implemented.

	// if the gateways are ready, and the UDPRoute is destined for them, ensure that
	// the object is pushed to the dataplane.
	if err := r.DataplaneClient.UpdateObject(udproute); err != nil {
		debug(log, udproute, "failed to update object in data-plane, requeueing")
		return ctrl.Result{}, err
	}

	// now that the object has been successfully configured for in the dataplane
	// we can update the object status to indicate that it's now properly linked
	// to the configured Gateways.
	debug(log, udproute, "ensuring status contains Gateway associations")
	statusUpdated, err := r.ensureGatewayReferenceStatusAdded(ctx, udproute, gateways...)
	if err != nil {
		// don't proceed until the statuses can be updated appropriately
		return ctrl.Result{}, err
	}
	if statusUpdated {
		// if the status was updated it will trigger a follow-up reconciliation
		// so we don't need to do anything further here.
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
			debug(log, udproute, "udproute not configured,requeueing")
			return ctrl.Result{Requeue: true}, nil
		}

		if configurationStatus == k8sobj.ConfigurationStatusFailed {
			debug(log, udproute, "tcproute configuration failed")
			statusUpdated, err := r.ensureParentsProgrammedCondition(ctx, udproute, gateways, metav1.ConditionFalse, ConditionReasonTranslationError, "")
			if err != nil {
				// don't proceed until the statuses can be updated appropriately
				debug(log, udproute, "failed to update programmed condition")
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: !statusUpdated}, nil
		}

		statusUpdated, err := r.ensureParentsProgrammedCondition(ctx, udproute, gateways, metav1.ConditionTrue, ConditionReasonConfiguredInGateway, "")
		if err != nil {
			// don't proceed until the statuses can be updated appropriately
			debug(log, udproute, "failed to update programmed condition")
			return ctrl.Result{}, err
		}
		if statusUpdated {
			// if the status was updated it will trigger a follow-up reconciliation
			// so we don't need to do anything further here.
			debug(log, udproute, "programmed condition updated")
			return ctrl.Result{}, nil
		}
	}

	// once the data-plane has accepted the UDPRoute object, we're all set.
	info(log, udproute, "udproute has been configured on the data-plane")
	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// UDPRouteReconciler - Status Helpers
// -----------------------------------------------------------------------------

// udprouteParentKind indicates the only object KIND that this UDPRoute
// implementation supports for route object parent references.
var udprouteParentKind = "Gateway"

// ensureGatewayReferenceStatus takes any number of Gateways that should be
// considered "attached" to a given UDPRoute and ensures that the status
// for the UDPRoute is updated appropriately.
func (r *UDPRouteReconciler) ensureGatewayReferenceStatusAdded(ctx context.Context, udproute *gatewayv1alpha2.UDPRoute, gateways ...supportedGatewayWithCondition) (bool, error) {
	// map the existing parentStatues to avoid duplications
	parentStatuses := getParentStatuses(udproute, udproute.Status.Parents)

	// overlay the parent ref statuses for all new gateway references
	statusChangesWereMade := false
	for _, gateway := range gateways {
		// build a new status for the parent Gateway
		gatewayParentStatus := &gatewayv1alpha2.RouteParentStatus{
			ParentRef: gatewayv1alpha2.ParentReference{
				Group:     (*gatewayv1alpha2.Group)(&gatewayv1alpha2.GroupVersion.Group),
				Kind:      (*gatewayv1alpha2.Kind)(util.StringToGatewayAPIKindPtr(udprouteParentKind)),
				Namespace: (*gatewayv1alpha2.Namespace)(&gateway.gateway.Namespace),
				Name:      gatewayv1alpha2.ObjectName(gateway.gateway.Name),
			},
			ControllerName: (gatewayv1alpha2.GatewayController)(GetControllerName()),
			Conditions: []metav1.Condition{{
				Type:               string(gatewayv1alpha2.RouteConditionAccepted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: udproute.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1alpha2.RouteReasonAccepted),
			}},
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
		ObservedGeneration: udproute.Generation,
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

	// update the udproute status with the new status references
	udproute.Status.Parents = make([]gatewayv1alpha2.RouteParentStatus, 0, len(parentStatuses))
	for _, parent := range parentStatuses {
		udproute.Status.Parents = append(udproute.Status.Parents, *parent)
	}

	// update the object status in the API
	if err := r.Status().Update(ctx, udproute); err != nil {
		return false, err
	}

	// the status needed an update and it was updated successfully
	return true, nil
}

// ensureGatewayReferenceStatusRemoved uses the ControllerName provided by the Gateway
// implementation to prune status references to Gateways supported by this controller
// in the provided UDPRoute object.
func (r *UDPRouteReconciler) ensureGatewayReferenceStatusRemoved(ctx context.Context, udproute *gatewayv1alpha2.UDPRoute) (bool, error) {
	// drop all status references to supported Gateway objects
	newStatuses := make([]gatewayv1alpha2.RouteParentStatus, 0)
	for _, status := range udproute.Status.Parents {
		if status.ControllerName != (gatewayv1alpha2.GatewayController)(GetControllerName()) {
			newStatuses = append(newStatuses, status)
		}
	}

	// if the new list of statuses is the same length as the old
	// nothing has changed and we're all done.
	if len(newStatuses) == len(udproute.Status.Parents) {
		return false, nil
	}

	// update the object status in the API
	udproute.Status.Parents = newStatuses
	if err := r.Status().Update(ctx, udproute); err != nil {
		return false, err
	}

	// the status needed an update and it was updated successfully
	return true, nil
}

func (r *UDPRouteReconciler) ensureParentsProgrammedCondition(
	ctx context.Context,
	udproute *gatewayv1alpha2.UDPRoute,
	gateways []supportedGatewayWithCondition,
	conditionStatus metav1.ConditionStatus,
	conditionReason gatewayv1beta1.RouteConditionReason,
	conditionMessage string,
) (bool, error) {
	// map the existing parentStatues to avoid duplications
	parentStatuses := getParentStatuses(udproute, udproute.Status.Parents)

	programmedCondition := metav1.Condition{
		Type:               ConditionTypeProgrammed,
		Status:             conditionStatus,
		Reason:             string(conditionReason),
		ObservedGeneration: udproute.Generation,
		Message:            conditionMessage,
		LastTransitionTime: metav1.Now(),
	}
	statusChanged := false
	for _, g := range gateways {
		gateway := g.gateway
		parentRefKey := fmt.Sprintf("%s/%s", gateway.Namespace, gateway.Name)
		parentStatus, ok := parentStatuses[parentRefKey]
		if ok {
			// update existing parent in status.
			changed := setRouteParentStatusCondition(parentStatus, programmedCondition)
			statusChanged = statusChanged || changed
		} else {
			// add a new parent if the parent is not found in status.
			newParentStatus := &gatewayv1alpha2.RouteParentStatus{
				ParentRef: gatewayv1alpha2.ParentReference{
					Namespace:   lo.ToPtr(gatewayv1alpha2.Namespace(gateway.Namespace)),
					Name:        gatewayv1alpha2.ObjectName(gateway.Name),
					SectionName: lo.ToPtr(gatewayv1alpha2.SectionName(g.listenerName)),
					// TODO: set port after gateway port matching implemented: https://github.com/Kong/kubernetes-ingress-controller/issues/3016
				},
				Conditions: []metav1.Condition{
					programmedCondition,
				},
			}
			udproute.Status.Parents = append(udproute.Status.Parents, *newParentStatus)
			parentStatuses[parentRefKey] = newParentStatus
			statusChanged = true
		}
	}

	// update status if needed.
	if statusChanged {
		if err := r.Status().Update(ctx, udproute); err != nil {
			return false, err
		}
		return true, nil
	}
	// no need to update if no status is changed.
	return false, nil
}
