package gateway

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	ctrlutils "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	k8sobj "github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/kubernetes/object/status"
)

// -----------------------------------------------------------------------------
// HTTPRoute Controller - HTTPRouteReconciler
// -----------------------------------------------------------------------------

// HTTPRouteReconciler reconciles an HTTPRoute object.
type HTTPRouteReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	DataplaneClient  controllers.DataPlane
	CacheSyncTimeout time.Duration
	StatusQueue      *status.Queue

	// If enableReferenceGrant is true, we will check for ReferenceGrant if backend in another
	// namespace is in backendRefs.
	// If it is false, referencing backend in different namespace will be rejected.
	// It's resolved on SetupWithManager call.
	enableReferenceGrant bool

	// If GatewayNN is set,
	// only resources managed by the specified Gateway are reconciled.
	GatewayNN controllers.OptionalNamespacedName
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("httproute-controller", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}

	// We're verifying whether ReferenceGrant CRD is installed at setup of the HTTPRouteReconciler
	// to decide whether we should run additional ReferenceGrant watch and handle ReferenceGrants
	// when reconciling HTTPRoutes.
	// Once the HTTPRouteReconciler is set up without ReferenceGrant, there's no possibility to enable
	// ReferenceGrant handling again in this reconciler at runtime.
	r.enableReferenceGrant = ctrlutils.CRDExists(mgr.GetRESTMapper(), schema.GroupVersionResource{
		Group:    gatewayv1beta1.GroupVersion.Group,
		Version:  gatewayv1beta1.GroupVersion.Version,
		Resource: "referencegrants",
	})

	// if a GatewayClass updates then we need to enqueue the linked HTTPRoutes to
	// ensure that any route objects that may have been orphaned by that change get
	// removed from data-plane configurations, and any routes that are now supported
	// due to that change get added to data-plane configurations.
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &gatewayapi.GatewayClass{}),
		handler.EnqueueRequestsFromMapFunc(r.listHTTPRoutesForGatewayClass),
		predicate.Funcs{
			GenericFunc: func(_ event.GenericEvent) bool { return false }, // we don't need to enqueue from generic
			CreateFunc:  func(e event.CreateEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
			UpdateFunc:  func(e event.UpdateEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
			DeleteFunc:  func(e event.DeleteEvent) bool { return isGatewayClassEventInClass(r.Log, e) },
		},
	); err != nil {
		return err
	}

	// if a Gateway updates then we need to enqueue the linked HTTPRoutes to
	// ensure that any route objects that may have been orphaned by that change get
	// removed from data-plane configurations, and any routes that are now supported
	// due to that change get added to data-plane configurations.
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &gatewayapi.Gateway{}),
		handler.EnqueueRequestsFromMapFunc(r.listHTTPRoutesForGateway),
	); err != nil {
		return err
	}

	if r.enableReferenceGrant {
		if err := c.Watch(
			source.Kind(mgr.GetCache(), &gatewayapi.ReferenceGrant{}),
			handler.EnqueueRequestsFromMapFunc(r.listReferenceGrantsForHTTPRoute),
			predicate.NewPredicateFuncs(referenceGrantHasHTTPRouteFrom),
		); err != nil {
			return err
		}
	}

	if r.StatusQueue != nil {
		if err := c.Watch(
			&source.Channel{Source: r.StatusQueue.Subscribe(schema.GroupVersionKind{
				Group:   gatewayv1.GroupVersion.Group,
				Version: gatewayv1.GroupVersion.Version,
				Kind:    "HTTPRoute",
			})},
			&handler.EnqueueRequestForObject{},
		); err != nil {
			return err
		}
	}

	// because of the additional burden of having to manage reference data-plane
	// configurations for HTTPRoute objects in the underlying Kong Gateway, we
	// simply reconcile ALL HTTPRoute objects. This allows us to drop the backend
	// data-plane config for an HTTPRoute if it somehow becomes disconnected from
	// a supported Gateway and GatewayClass.
	return c.Watch(
		source.Kind(mgr.GetCache(), &gatewayapi.HTTPRoute{}),
		&handler.EnqueueRequestForObject{},
	)
}

// -----------------------------------------------------------------------------
// HTTPRoute Controller - Event Handlers
// -----------------------------------------------------------------------------

// listReferenceGrantsForHTTPRoute is a watch predicate which finds all HTTPRoutes
// mentioned in a From clause for a ReferenceGrant.
func (r *HTTPRouteReconciler) listReferenceGrantsForHTTPRoute(ctx context.Context, obj client.Object) []reconcile.Request {
	grant, ok := obj.(*gatewayapi.ReferenceGrant)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"Referencegrant watch predicate received unexpected object type",
			"expected", "*gatewayapi.ReferenceGrant", "found", reflect.TypeOf(obj),
		)
		return nil
	}
	httproutes := &gatewayapi.HTTPRouteList{}
	if err := r.Client.List(ctx, httproutes); err != nil {
		r.Log.Error(err, "Failed to list httproutes in watch", "referencegrant", grant.Name)
		return nil
	}
	recs := []reconcile.Request{}
	for _, gateway := range httproutes.Items {
		for _, from := range grant.Spec.From {
			if string(from.Namespace) == gateway.Namespace &&
				from.Kind == ("HTTPRoute") &&
				from.Group == ("gateway.networking.k8s.io") {
				recs = append(recs, reconcile.Request{
					NamespacedName: k8stypes.NamespacedName{
						Namespace: gateway.Namespace,
						Name:      gateway.Name,
					},
				})
			}
		}
	}
	return recs
}

func referenceGrantHasHTTPRouteFrom(obj client.Object) bool {
	grant, ok := obj.(*gatewayapi.ReferenceGrant)
	if !ok {
		return false
	}
	for _, from := range grant.Spec.From {
		if from.Kind == "HTTPRoute" && from.Group == "gateway.networking.k8s.io" {
			return true
		}
	}
	return false
}

// listHTTPRoutesForGatewayClass is a controller-runtime event.Handler which
// produces a list of HTTPRoutes which were bound to a Gateway which is or was
// bound to this GatewayClass. This implementation effectively does a map-reduce
// to determine the HTTProutes as the relationship has to be discovered entirely
// by object reference. This relies heavily on the inherent performance benefits of
// the cached manager client to avoid API overhead.
func (r *HTTPRouteReconciler) listHTTPRoutesForGatewayClass(ctx context.Context, obj client.Object) []reconcile.Request {
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

	// map all HTTPRoute objects
	httprouteList := gatewayapi.HTTPRouteList{}
	if err := r.Client.List(ctx, &httprouteList); err != nil {
		r.Log.Error(err, "Failed to list httproute objects from the cached client")
		return nil
	}

	// reduce for HTTPRoute objects bound to an in-class Gateway
	queue := make([]reconcile.Request, 0)
	for _, httproute := range httprouteList.Items {
		// check the httproute's parentRefs
		for _, parentRef := range httproute.Spec.ParentRefs {
			// determine what namespace the parent gateway is in
			namespace := httproute.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}

			// if the gateway matches one of our previously filtered gateways, enqueue the route
			if gatewaysForNamespace, ok := gateways[namespace]; ok {
				if _, ok := gatewaysForNamespace[string(parentRef.Name)]; ok {
					queue = append(queue, reconcile.Request{
						NamespacedName: k8stypes.NamespacedName{
							Namespace: httproute.Namespace,
							Name:      httproute.Name,
						},
					})
				}
			}
		}
	}

	return queue
}

// listHTTPRoutesForGateway is a controller-runtime event.Handler which enqueues HTTPRoute
// objects for changes to Gateway objects. The relationship between HTTPRoutes and their
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
// implementation enqueues ALL HTTPRoute objects for reconciliation every time a Gateway
// changes. This is not ideal, but after communicating with other members of the
// community this appears to be a standard approach across multiple implementations at
// the moment for v1alpha2. As future releases of Gateway come out we'll need to
// continue iterating on this and perhaps advocating for upstream changes to help avoid
// this kind of problem without having to enqueue extra objects.
func (r *HTTPRouteReconciler) listHTTPRoutesForGateway(ctx context.Context, obj client.Object) []reconcile.Request {
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

	// map all HTTPRoute objects
	httprouteList := gatewayapi.HTTPRouteList{}
	if err := r.Client.List(ctx, &httprouteList); err != nil {
		r.Log.Error(err, "Failed to list httproute objects from the cached client")
		return nil
	}

	// reduce for HTTPRoute objects bound to the Gateway
	queue := make([]reconcile.Request, 0)
	for _, httproute := range httprouteList.Items {
		for _, parentRef := range httproute.Spec.ParentRefs {
			namespace := httproute.Namespace
			if parentRef.Namespace != nil {
				namespace = string(*parentRef.Namespace)
			}
			if namespace == gw.Namespace && string(parentRef.Name) == gw.Name {
				queue = append(queue, reconcile.Request{
					NamespacedName: k8stypes.NamespacedName{
						Namespace: httproute.Namespace,
						Name:      httproute.Name,
					},
				})
			}
		}
	}

	return queue
}

// -----------------------------------------------------------------------------
// HTTPRoute Controller - Reconciliation
// -----------------------------------------------------------------------------

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes/status,verbs=get;update
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=list;watch;get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *HTTPRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("GatewayV1HTTPRoute", req.NamespacedName)

	httproute := new(gatewayapi.HTTPRoute)
	if err := r.Get(ctx, req.NamespacedName, httproute); err != nil {
		// if the queued object is no longer present in the proxy cache we need
		// to ensure that if it was ever added to the cache, it gets removed.
		if apierrors.IsNotFound(err) {
			debug(log, httproute, "Object does not exist, ensuring it is not present in the proxy cache")
			httproute.Namespace = req.Namespace
			httproute.Name = req.Name
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(httproute)
		}

		// for any error other than 404, requeue
		return ctrl.Result{}, err
	}

	debug(log, httproute, "Processing httproute")

	// if there's a present deletion timestamp then we need to update the proxy cache
	// to drop all relevant routes from its configuration, regardless of whether or
	// not we can find a valid gateway as that gateway may now be deleted but we still
	// need to ensure removal of the data-plane configuration.
	debug(log, httproute, "Checking deletion timestamp")
	if httproute.DeletionTimestamp != nil {
		debug(log, httproute, "HTTPRoute is being deleted, re-configuring data-plane")
		if err := r.DataplaneClient.DeleteObject(httproute); err != nil {
			debug(log, httproute, "Failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, httproute, "Ensured object was removed from the data-plane (if ever present)")
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(httproute)
	}

	// we need to pull the Gateway parent objects for the HTTPRoute to verify
	// routing behavior and ensure compatibility with Gateway configurations.
	debug(log, httproute, "Retrieving GatewayClass and Gateway for route")
	gateways, err := getSupportedGatewayForRoute(ctx, log, r.Client, httproute, r.GatewayNN)
	if err != nil {
		if err.Error() == unsupportedGW {
			debug(log, httproute, "Unsupported route found, processing to verify whether it was ever supported")
			// if there's no supported Gateway then this route could have been previously
			// supported by this controller. As such we ensure that no supported Gateway
			// references exist in the object status any longer.
			statusUpdated, err := r.ensureGatewayReferenceStatusRemoved(ctx, httproute)
			if err != nil {
				// some failure happened so we need to retry to avoid orphaned statuses
				return ctrl.Result{}, err
			}
			if statusUpdated {
				// the status did in fact needed to be updated, so no need to requeue
				// as the status update will trigger a requeue.
				debug(log, httproute, "Unsupported route was previously supported, status was updated")
				return ctrl.Result{}, nil
			}

			// if the route doesn't have a supported Gateway+GatewayClass associated with
			// it it's possible it became orphaned after becoming queued. In either case
			// ensure that it's removed from the proxy cache to avoid orphaned data-plane
			// configurations.
			debug(log, httproute, "Ensuring that dataplane is updated to remove unsupported route (if applicable)")
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(httproute)
		}
		return ctrl.Result{}, err
	}

	// Ensure we have no status for no-longer defined parentRefs.
	if wasAnyStatusRemoved := ensureNoStaleParentStatus(httproute); wasAnyStatusRemoved {
		if err := r.Status().Update(ctx, httproute); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// the referenced gateway object(s) for the HTTPRoute needs to be ready
	// before we'll attempt any configurations of it. If it's not we'll
	// requeue the object and wait until all supported gateways are ready.
	debug(log, httproute, "Checking if the httproute's gateways are ready")
	for _, gateway := range gateways {
		if !isGatewayProgrammed(gateway.gateway) {
			debug(log, httproute, "Gateway for route was not ready, waiting")
			return ctrl.Result{Requeue: true}, nil
		}
	}

	// if there is no matched hosts in listeners for the httproute, the httproute should not be accepted
	// and have an "Accepted" condition with status false.
	// https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io/v1.HTTPRoute
	filteredHTTPRoute, err := filterHostnames(gateways, httproute.DeepCopy())
	// perform operations on the kong store only if the route is in accepted status and there is hostname matching
	if isRouteAccepted(gateways) && err == nil {
		// if the gateways are ready, and the HTTPRoute is destined for them, ensure that
		// the object is pushed to the dataplane.
		if err := r.DataplaneClient.UpdateObject(filteredHTTPRoute); err != nil {
			debug(log, httproute, "Failed to update object in data-plane, requeueing")
			return ctrl.Result{}, err
		}
	} else {
		// route is not accepted, remove it from kong store
		if err := r.DataplaneClient.DeleteObject(httproute); err != nil {
			debug(log, httproute, "Failed to delete object in data-plane, requeueing")
			return ctrl.Result{}, err
		}
	}

	// now that the object has been successfully configured for in the dataplane
	// we can update the object status to indicate that it's now properly linked
	// to the configured Gateways.
	debug(log, httproute, "Ensuring status contains Gateway associations")
	statusUpdated, err := r.ensureGatewayReferenceStatusAdded(ctx, httproute, gateways...)
	if err != nil {
		// don't proceed until the statuses can be updated appropriately
		return ctrl.Result{}, err
	}
	if statusUpdated {
		// if the status was updated it will trigger a follow-up reconciliation
		// so we don't need to do anything further here.
		return ctrl.Result{}, nil
	}

	// update "Programmed" condition if HTTPRoute is translated to Kong configuration.
	// if the HTTPRoute is not configured in the dataplane, leave it unchanged and requeue.
	// if it is successfully configured, update its "Programmed" condition to True.
	// if translation failure happens, update its "Programmed" condition to False.
	debug(log, httproute, "Ensuring status contains Programmed condition")
	if r.DataplaneClient.AreKubernetesObjectReportsEnabled() {
		// if the dataplane client has reporting enabled (this is the default and is
		// tied in with status updates being enabled in the controller manager) then
		// we will wait until the object is reported as successfully configured before
		// moving on to status updates.
		configurationStatus := r.DataplaneClient.KubernetesObjectConfigurationStatus(httproute)
		if configurationStatus == k8sobj.ConfigurationStatusUnknown {
			// requeue until httproute is configured.
			debug(log, httproute, "HTTPRoute not configured, requeueing")
			return ctrl.Result{Requeue: true}, nil
		}

		if configurationStatus == k8sobj.ConfigurationStatusFailed {
			debug(log, httproute, "HTTPRoute configuration failed")
			statusUpdated, err := ensureParentsProgrammedCondition(ctx, r.Status(), httproute, httproute.Status.Parents, gateways, metav1.Condition{
				Status: metav1.ConditionFalse,
				Reason: string(ConditionReasonTranslationError),
			})
			if err != nil {
				// don't proceed until the statuses can be updated appropriately
				debug(log, httproute, "Failed to update programmed condition")
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: !statusUpdated}, nil
		}

		statusUpdated, err := ensureParentsProgrammedCondition(ctx, r.Status(), httproute, httproute.Status.Parents, gateways, metav1.Condition{
			Status: metav1.ConditionTrue,
			Reason: string(ConditionReasonConfiguredInGateway),
		})
		if err != nil {
			// don't proceed until the statuses can be updated appropriately
			debug(log, httproute, "Failed to update programmed condition")
			return ctrl.Result{}, err
		}
		if statusUpdated {
			// if the status was updated it will trigger a follow-up reconciliation
			// so we don't need to do anything further here.
			debug(log, httproute, "Programmed condition updated")
			return ctrl.Result{}, nil
		}
	}

	// once the data-plane has accepted the HTTPRoute object, we're all set.
	info(log, httproute, "HTTPRoute has been configured on the data-plane")

	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// HTTPRouteReconciler - Status Helpers
// -----------------------------------------------------------------------------

// httprouteParentKind indicates the only object KIND that this HTTPRoute
// implementation supports for route object parent references.
var httprouteParentKind = "Gateway"

// ensureGatewayReferenceStatus takes any number of Gateways that should be
// considered "attached" to a given HTTPRoute and ensures that the status
// for the HTTPRoute is updated appropriately.
func (r *HTTPRouteReconciler) ensureGatewayReferenceStatusAdded(ctx context.Context, httproute *gatewayapi.HTTPRoute, gateways ...supportedGatewayWithCondition) (bool, error) {
	// map the existing parentStatues to avoid duplications
	parentStatuses := getParentStatuses(httproute, httproute.Status.Parents)

	// overlay the parent ref statuses for all new gateway references
	statusChangesWereMade := false
	for _, gateway := range gateways {
		gateway := gateway
		// build a new status for the parent Gateway
		gatewayParentStatus := &gatewayapi.RouteParentStatus{
			ParentRef: gatewayapi.ParentReference{
				Group:     (*gatewayapi.Group)(&gatewayv1.GroupVersion.Group),
				Kind:      util.StringToGatewayAPIKindPtr(httprouteParentKind),
				Namespace: (*gatewayapi.Namespace)(&gateway.gateway.Namespace),
				Name:      gatewayapi.ObjectName(gateway.gateway.Name),
			},
			ControllerName: GetControllerName(),
			Conditions: []metav1.Condition{{
				Type:               gateway.condition.Type,
				Status:             gateway.condition.Status,
				ObservedGeneration: httproute.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             gateway.condition.Reason,
			}},
		}
		if gateway.listenerName != "" {
			gatewayParentStatus.ParentRef.SectionName = lo.ToPtr(gatewayapi.SectionName(gateway.listenerName))
		}

		key := fmt.Sprintf("%s/%s/%s", gateway.gateway.Namespace, gateway.gateway.Name, gateway.listenerName)

		// if the reference already exists and doesn't require any changes
		// then just leave it alone.
		if existingGatewayParentStatus, exists := parentStatuses[key]; exists {
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
		parentStatuses[key] = gatewayParentStatus
		statusChangesWereMade = true
	}

	parentStatuses, resolvedRefsChanged, err := r.setRouteConditionResolvedRefsCondition(ctx, httproute, parentStatuses)
	if err != nil {
		return false, err
	}

	// initialize "programmed" condition to Unknown.
	// do not update the condition If a "Programmed" condition is already present.
	programmedConditionChanged := false
	programmedConditionUnknown := metav1.Condition{
		Type:               ConditionTypeProgrammed,
		Status:             metav1.ConditionUnknown,
		Reason:             string(ConditionReasonProgrammedUnknown),
		ObservedGeneration: httproute.Generation,
		LastTransitionTime: metav1.Now(),
	}
	for _, parentStatus := range parentStatuses {
		if !parentStatusHasProgrammedCondition(parentStatus) {
			programmedConditionChanged = true
			parentStatus.Conditions = append(parentStatus.Conditions, programmedConditionUnknown)
		}
	}

	// if we didn't have to actually make any changes, no status update is needed
	if !statusChangesWereMade && !resolvedRefsChanged && !programmedConditionChanged {
		return false, nil
	}

	// update the httproute status with the new status references
	httproute.Status.Parents = make([]gatewayapi.RouteParentStatus, 0, len(parentStatuses))
	for _, parent := range parentStatuses {
		httproute.Status.Parents = append(httproute.Status.Parents, *parent)
	}

	// update the object status in the API
	if err := r.Status().Update(ctx, httproute); err != nil {
		return false, err
	}

	// the status needed an update and it was updated successfully
	return true, nil
}

// ensureGatewayReferenceStatusRemoved uses the ControllerName provided by the Gateway
// implementation to prune status references to Gateways supported by this controller
// in the provided HTTPRoute object.
func (r *HTTPRouteReconciler) ensureGatewayReferenceStatusRemoved(ctx context.Context, httproute *gatewayapi.HTTPRoute) (bool, error) {
	// drop all status references to supported Gateway objects
	newStatuses := make([]gatewayapi.RouteParentStatus, 0)
	for _, status := range httproute.Status.Parents {
		if status.ControllerName != GetControllerName() {
			newStatuses = append(newStatuses, status)
		}
	}

	// if the new list of statuses is the same length as the old
	// nothing has changed and we're all done.
	if len(newStatuses) == len(httproute.Status.Parents) {
		return false, nil
	}

	// update the object status in the API
	httproute.Status.Parents = newStatuses
	if err := r.Status().Update(ctx, httproute); err != nil {
		return false, err
	}

	// the status needed an update and it was updated successfully
	return true, nil
}

// setRouteConditionResolvedRefsCondition sets a condition of type ResolvedRefs on the route status.
func (r *HTTPRouteReconciler) setRouteConditionResolvedRefsCondition(
	ctx context.Context,
	httpRoute *gatewayapi.HTTPRoute,
	parentStatuses map[string]*gatewayapi.RouteParentStatus,
) (map[string]*gatewayapi.RouteParentStatus, bool, error) {
	var changed bool
	resolvedRefsStatus := metav1.ConditionFalse
	reason, err := r.getHTTPRouteRuleReason(ctx, *httpRoute)
	if err != nil {
		return nil, false, err
	}
	if reason == gatewayapi.RouteReasonResolvedRefs {
		resolvedRefsStatus = metav1.ConditionTrue
	}

	// iterate over all the parentStatuses conditions, and if no RouteConditionResolvedRefs is found,
	// or if the condition is found but has to be changed, update the status and mark it to be updated
	resolvedRefsCondition := metav1.Condition{
		Type:               string(gatewayapi.RouteConditionResolvedRefs),
		Status:             resolvedRefsStatus,
		ObservedGeneration: httpRoute.Generation,
		LastTransitionTime: metav1.Now(),
		Reason:             string(reason),
	}
	for _, parentStatus := range parentStatuses {
		var conditionFound bool
		for i, cond := range parentStatus.Conditions {
			if cond.Type == string(gatewayapi.RouteConditionResolvedRefs) {
				if !(cond.Status == resolvedRefsStatus &&
					cond.Reason == string(reason)) {
					parentStatus.Conditions[i] = resolvedRefsCondition
					changed = true
				}
				conditionFound = true
				break
			}
		}
		if !conditionFound {
			parentStatus.Conditions = append(parentStatus.Conditions, resolvedRefsCondition)
			changed = true
		}
	}

	return parentStatuses, changed, nil
}

func (r *HTTPRouteReconciler) getHTTPRouteRuleReason(ctx context.Context, httpRoute gatewayapi.HTTPRoute) (gatewayapi.RouteConditionReason, error) {
	for _, rule := range httpRoute.Spec.Rules {
		for _, backendRef := range rule.BackendRefs {
			backendNamespace := httpRoute.Namespace
			if backendRef.Namespace != nil && *backendRef.Namespace != "" {
				backendNamespace = string(*backendRef.Namespace)
			}

			// Check if the BackendRef GroupKind is supported
			if !util.IsBackendRefGroupKindSupported(backendRef.Group, backendRef.Kind) {
				return gatewayapi.RouteReasonInvalidKind, nil
			}

			// Check if all the objects referenced actually exist
			// Only services are currently supported as BackendRef objects
			service := &corev1.Service{}
			err := r.Client.Get(ctx, k8stypes.NamespacedName{Namespace: backendNamespace, Name: string(backendRef.Name)}, service)
			if err != nil {
				if !apierrors.IsNotFound(err) {
					return "", err
				}
				return gatewayapi.RouteReasonBackendNotFound, nil
			}

			// Check if the object referenced is in another namespace,
			// and if there is grant for that reference
			if httpRoute.Namespace != backendNamespace {
				if !r.enableReferenceGrant {
					return gatewayapi.RouteReasonRefNotPermitted, nil
				}

				referenceGrantList := &gatewayapi.ReferenceGrantList{}
				if err := r.Client.List(ctx, referenceGrantList, client.InNamespace(backendNamespace)); err != nil {
					return "", err
				}
				if len(referenceGrantList.Items) == 0 {
					return gatewayapi.RouteReasonRefNotPermitted, nil
				}
				var isGranted bool
				for _, grant := range referenceGrantList.Items {
					if isHTTPReferenceGranted(grant.Spec, backendRef, httpRoute.Namespace) {
						isGranted = true
						break
					}
				}
				if !isGranted {
					return gatewayapi.RouteReasonRefNotPermitted, nil
				}
			}
		}
	}
	return gatewayapi.RouteReasonResolvedRefs, nil
}

// SetLogger sets the logger.
func (r *HTTPRouteReconciler) SetLogger(l logr.Logger) {
	r.Log = l
}

// ensureNoStaleParentStatus removes any status references to Gateways that are no longer in the HTTPRoute's parentRefs
// and returns true if any status was removed.
func ensureNoStaleParentStatus(httproute *gatewayapi.HTTPRoute) (wasAnyStatusRemoved bool) {
	// Create a map of currently defined parentRefs for fast lookup.
	currentlyDefinedParentRefs := make(map[string]struct{})
	for _, parentRef := range httproute.Spec.ParentRefs {
		currentlyDefinedParentRefs[parentReferenceKey(httproute.Namespace, parentRef)] = struct{}{}
	}

	for parentIdx, parentStatus := range httproute.Status.Parents {
		// Don't touch statuses from other controllers.
		if parentStatus.ControllerName != GetControllerName() {
			continue
		}
		// Remove the status if the parentRef is no longer defined.
		if _, ok := currentlyDefinedParentRefs[parentReferenceKey(httproute.Namespace, parentStatus.ParentRef)]; !ok {
			httproute.Status.Parents = slices.Delete(httproute.Status.Parents, parentIdx, parentIdx+1)
			wasAnyStatusRemoved = true
		}
	}
	return wasAnyStatusRemoved
}

// parentReferenceKey returns a string key for a parentRef of a route. It can be used for indexing a map.
func parentReferenceKey(routeNamespace string, parentRef gatewayapi.ParentReference) string {
	namespace := routeNamespace
	if parentRef.Namespace != nil {
		namespace = string(*parentRef.Namespace)
	}
	sectionName := ""
	if parentRef.SectionName != nil {
		sectionName = string(*parentRef.SectionName)
	}
	portNumber := ""
	if parentRef.Port != nil {
		portNumber = strconv.Itoa(int(*parentRef.Port))
	}

	// We intentionally do not take into account Kind and Group here as we only support Gateways
	// and that's the only kind we should be getting here thanks to the admission webhook validation.
	return fmt.Sprintf("%s/%s/%s/%s", namespace, parentRef.Name, sectionName, portNumber)
}
