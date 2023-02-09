package gateway

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/reference"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
)

// -----------------------------------------------------------------------------
// Vars & Consts
// -----------------------------------------------------------------------------

var (
	ErrUnmanagedAnnotation = errors.New("invalid unmanaged annotation value")
	gatewayV1beta1Group    = gatewayv1beta1.Group(gatewayv1beta1.GroupName)
)

// -----------------------------------------------------------------------------
// Gateway Controller - GatewayReconciler
// -----------------------------------------------------------------------------

// GatewayReconciler reconciles a Gateway object.
type GatewayReconciler struct { //nolint:revive
	client.Client

	Log             logr.Logger
	Scheme          *runtime.Scheme
	DataplaneClient *dataplane.KongClient

	PublishService  string
	WatchNamespaces []string
	// If EnableReferenceGrant is true, controller will watch ReferenceGrants
	// to invalidate or allow cross-namespace TLSConfigs in gateways.
	EnableReferenceGrant bool
	CacheSyncTimeout     time.Duration

	ReferenceIndexers ctrlref.CacheIndexers

	publishServiceRef types.NamespacedName
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// verify that the PublishService was configured properly
	var err error
	r.publishServiceRef, err = getRefFromPublishService(r.PublishService)
	if err != nil {
		return err
	}

	// generate the controller object and attach it to the manager and link the reconciler object
	c, err := controller.New("gateway-controller", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}

	// watch Gateway objects, filtering out any Gateways which are not configured with
	// a supported GatewayClass controller name.
	if err := c.Watch(
		&source.Kind{Type: &gatewayv1beta1.Gateway{}},
		&handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(r.gatewayHasMatchingGatewayClass),
	); err != nil {
		return err
	}

	// watch for updates to gatewayclasses, if any gateway classes change, enqueue
	// reconciliation for all supported gateway objects which reference it.
	if err := c.Watch(
		&source.Kind{Type: &gatewayv1beta1.GatewayClass{}},
		handler.EnqueueRequestsFromMapFunc(r.listGatewaysForGatewayClass),
		predicate.NewPredicateFuncs(r.gatewayClassMatchesController),
	); err != nil {
		return err
	}

	// if an update to the gateway service occurs, we need to make sure to trigger
	// reconciliation on all Gateway objects referenced by it (in the most common
	// deployments this will be a single Gateway).
	if err := c.Watch(
		&source.Kind{Type: &corev1.Service{}},
		handler.EnqueueRequestsFromMapFunc(r.listGatewaysForService),
		predicate.NewPredicateFuncs(r.isGatewayService),
	); err != nil {
		return err
	}

	// watch ReferenceGrants, which may invalidate or allow cross-namespace TLSConfigs
	if r.EnableReferenceGrant {
		if err := c.Watch(
			&source.Kind{Type: &gatewayv1beta1.ReferenceGrant{}},
			handler.EnqueueRequestsFromMapFunc(r.listReferenceGrantsForGateway),
			predicate.NewPredicateFuncs(referenceGrantHasGatewayFrom),
		); err != nil {
			return err
		}
	}

	// start the required gatewayclass controller as well
	gwcCTRL := &GatewayClassReconciler{
		Client:           r.Client,
		Log:              r.Log.WithName("V1Beta1GatewayClass"),
		Scheme:           r.Scheme,
		CacheSyncTimeout: r.CacheSyncTimeout,
	}

	return gwcCTRL.SetupWithManager(mgr)
}

// -----------------------------------------------------------------------------
// Gateway Controller - Watch Predicates
// -----------------------------------------------------------------------------

// gatewayHasMatchingGatewayClass is a watch predicate which filters out reconciliation events for
// gateway objects which aren't supported by this controller or not using an unmanaged GatewayClass.
func (r *GatewayReconciler) gatewayHasMatchingGatewayClass(obj client.Object) bool {
	gateway, ok := obj.(*gatewayv1beta1.Gateway)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"gateway watch predicate received unexpected object type",
			"expected", "*gatewayv1beta1.Gateway", "found", reflect.TypeOf(obj),
		)
		return false
	}
	gatewayClass := &gatewayv1beta1.GatewayClass{}
	if err := r.Client.Get(context.Background(), client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, gatewayClass); err != nil {
		r.Log.Error(err, "could not retrieve gatewayclass", "gatewayclass", gateway.Spec.GatewayClassName)
		return false
	}
	return isGatewayClassControlledAndUnmanaged(gatewayClass)
}

// gatewayClassMatchesController is a watch predicate which filters out events for gatewayclasses which
// aren't configured with the required ControllerName or not annotated as unmanaged.
func (r *GatewayReconciler) gatewayClassMatchesController(obj client.Object) bool {
	gatewayClass, ok := obj.(*gatewayv1beta1.GatewayClass)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"gatewayclass watch predicate received unexpected object type",
			"expected", "*gatewayv1beta1.GatewayClass", "found", reflect.TypeOf(obj),
		)
		return false
	}
	return isGatewayClassControlledAndUnmanaged(gatewayClass)
}

// listGatewaysForGatewayClass is a watch predicate which finds all the gateway objects reference
// by a gatewayclass to enqueue them for reconciliation. This is generally used when a GatewayClass
// is updated to ensure that idle gateways are initialized when their gatewayclass becomes available.
func (r *GatewayReconciler) listGatewaysForGatewayClass(gatewayClass client.Object) []reconcile.Request {
	gateways := &gatewayv1beta1.GatewayList{}
	if err := r.Client.List(context.Background(), gateways); err != nil {
		r.Log.Error(err, "failed to list gateways for gatewayclass in watch", "gatewayclass", gatewayClass.GetName())
		return nil
	}
	return reconcileGatewaysIfClassMatches(gatewayClass, gateways.Items)
}

// listReferenceGrantsForGateway is a watch predicate which finds all Gateways mentioned in a From clause for a
// ReferenceGrant.
func (r *GatewayReconciler) listReferenceGrantsForGateway(obj client.Object) []reconcile.Request {
	grant, ok := obj.(*gatewayv1beta1.ReferenceGrant)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"referencegrant watch predicate received unexpected object type",
			"expected", "*gatewayv1beta1.ReferenceGrant", "found", reflect.TypeOf(obj),
		)
		return nil
	}
	gateways := &gatewayv1beta1.GatewayList{}
	if err := r.Client.List(context.Background(), gateways); err != nil {
		r.Log.Error(err, "failed to list gateways in watch", "referencegrant", grant.Name)
		return nil
	}
	recs := []reconcile.Request{}
	for _, gateway := range gateways.Items {
		for _, from := range grant.Spec.From {
			if string(from.Namespace) == gateway.Namespace &&
				from.Kind == gatewayv1alpha2.Kind("Gateway") &&
				from.Group == gatewayv1alpha2.Group("gateway.networking.k8s.io") {
				recs = append(recs, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: gateway.Namespace,
						Name:      gateway.Name,
					},
				})
			}
		}
	}
	return recs
}

// listGatewaysForService is a watch predicate which finds all the gateway objects which use
// GatewayClasses supported by this controller and are configured for the same service via
// unmanaged mode and enqueues them for reconciliation. This is generally used to ensure
// all gateways are updated when the service gets updated with new listeners.
func (r *GatewayReconciler) listGatewaysForService(svc client.Object) (recs []reconcile.Request) {
	gateways := &gatewayv1beta1.GatewayList{}
	if err := r.Client.List(context.Background(), gateways); err != nil {
		r.Log.Error(err, "failed to list gateways for service in watch predicates", "service", svc)
		return
	}
	for _, gateway := range gateways.Items {
		gatewayClass := &gatewayv1beta1.GatewayClass{}
		if err := r.Client.Get(context.Background(), types.NamespacedName{Name: string(gateway.Spec.GatewayClassName)}, gatewayClass); err != nil {
			r.Log.Error(err, "failed to retrieve gateway class in watch predicates", "gatewayclass", gateway.Spec.GatewayClassName)
			return
		}
		if isGatewayClassControlledAndUnmanaged(gatewayClass) {
			recs = append(recs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: gateway.Namespace,
					Name:      gateway.Name,
				},
			})
		}
	}
	return
}

// isGatewayService is a watch predicate that filters out events for objects that aren't
// the gateway service referenced by --publish-service.
func (r *GatewayReconciler) isGatewayService(obj client.Object) bool {
	return fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName()) == r.PublishService
}

func referenceGrantHasGatewayFrom(obj client.Object) bool {
	grant, ok := obj.(*gatewayv1beta1.ReferenceGrant)
	if !ok {
		return false
	}
	for _, from := range grant.Spec.From {
		if from.Kind == "Gateway" && from.Group == "gateway.networking.k8s.io" {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------
// Gateway Controller - Reconciliation
// -----------------------------------------------------------------------------

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/status,verbs=get;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("GatewayV1Beta1Gateway", req.NamespacedName)

	// gather the gateway object based on the reconciliation trigger. It's possible for the object
	// to be gone at this point in which case it will be ignored.
	gateway := new(gatewayv1beta1.Gateway)
	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		if apierrors.IsNotFound(err) {
			gateway.Namespace = req.Namespace
			gateway.Name = req.Name
			// delete reference relationships where the gateway is the referrer.
			err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, gateway)
			if err != nil {
				return ctrl.Result{}, err
			}
			debug(log, gateway, "reconciliation triggered but gateway does not exist, deleting it in dataplane")
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(gateway)
		}
		return ctrl.Result{Requeue: true}, err
	}
	debug(log, gateway, "processing gateway")

	// though our watch configuration eliminates reconciliation of unsupported gateways it's
	// technically possible for the gatewayclass configuration of a gateway to change in
	// the interim while the object has been queued for reconciliation. This double check
	// reduces reconciliation operations that would occur on old information.
	debug(log, gateway, "verifying gatewayclass")
	gwc := &gatewayv1beta1.GatewayClass{}
	if err := r.Client.Get(ctx, client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, gwc); err != nil {
		debug(log, gateway, "could not retrieve gatewayclass for gateway", "gatewayclass", string(gateway.Spec.GatewayClassName))
		// delete reference relationships where the gateway is the referrer, as we will not process the gateway.
		err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, gateway)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err := r.DataplaneClient.DeleteObject(gateway); err != nil {
			debug(log, gateway, "failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, gateway, "ensured gateway was removed from the data-plane (if ever present)")
		return ctrl.Result{}, nil
	}
	if gwc.Spec.ControllerName != GetControllerName() {
		debug(log, gateway, "unsupported gatewayclass controllername, ignoring", "gatewayclass", gwc.Name, "controllername", gwc.Spec.ControllerName)
		// delete reference relationships where the gateway is the referrer, as we will not process the gateway.
		err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, gateway)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err := r.DataplaneClient.DeleteObject(gateway); err != nil {
			debug(log, gateway, "failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, gateway, "ensured gateway was removed from the data-plane (if ever present)")
		return ctrl.Result{}, nil
	}

	// if there's any deletion timestamp on the object, we can simply ignore it. At this point
	// with unmanaged mode being the only option supported there are no finalizers and
	// the object should be cleaned up by GC promptly.
	debug(log, gateway, "checking deletion timestamp")
	if gateway.DeletionTimestamp != nil {
		debug(log, gateway, "gateway is being deleted, ignoring")
		return ctrl.Result{Requeue: false}, nil
	}

	// ensure that the GatewayClass matches the ControllerName and is unmanaged.
	// This check has already been performed by predicates, but we need to ensure this condition
	// as the reconciliation loop may be triggered by objects in which predicates we
	// cannot check the ControllerName and the unmanaged mode (e.g., ReferenceGrants).
	if !isGatewayClassControlledAndUnmanaged(gwc) {
		return reconcile.Result{}, nil
	}

	// reconciliation assumes unmanaged mode, in the future we may have a slot here for
	// other gateway management modes.
	result, err := r.reconcileUnmanagedGateway(ctx, log, gateway)
	// reconcileUnmanagedGateway has side effects and modifies the referenced gateway object. dataplane updates must
	// happen afterwards
	if err == nil {
		if err := r.DataplaneClient.UpdateObject(gateway); err != nil {
			debug(log, gateway, "failed to update object in data-plane, requeueing")
			return result, err
		}

		referredSecretNames := listSecretNamesReferredByGateway(gateway)
		if err := ctrlref.UpdateReferencesToSecret(
			ctx, r.Client, r.ReferenceIndexers, r.DataplaneClient,
			gateway, referredSecretNames); err != nil {
			if apierrors.IsNotFound(err) {
				result.Requeue = true
				return result, nil
			}
			return result, err
		}
	}
	return result, err
}

// reconcileUnmanagedGateway reconciles a Gateway that is configured for unmanaged mode,
// this mode will extract the Addresses and Listeners for the Gateway from the Kubernetes Service
// used for the Kong Gateway in the pre-existing deployment.
func (r *GatewayReconciler) reconcileUnmanagedGateway(ctx context.Context, log logr.Logger, gateway *gatewayv1beta1.Gateway) (ctrl.Result, error) {
	// currently this controller supports only unmanaged gateway mode, we need to verify
	// any Gateway object that comes to us is configured appropriately, and if not reject it
	// with a clear status condition and message.
	debug(log, gateway, "validating management mode for gateway") // this will also be done by the validating webhook, this is a fallback

	// enforce the service reference as the annotation value for the key UnmanagedGateway.
	debug(log, gateway, "initializing admin service annotation if unset")
	if !isObjectUnmanaged(gateway.GetAnnotations()) {
		debug(log, gateway, fmt.Sprintf("a placeholder value was provided for %s, adding the default service ref %s", annotations.GatewayClassUnmanagedAnnotation, r.PublishService))
		if gateway.Annotations == nil {
			gateway.Annotations = map[string]string{}
		}
		annotations.UpdateUnmanagedAnnotation(gateway.Annotations, r.PublishService)
		return ctrl.Result{}, r.Update(ctx, gateway)
	}

	serviceRef := annotations.ExtractUnmanagedGatewayClassMode(gateway.Annotations)
	// validation check of the Gateway to ensure that the publish service is actually available
	// in the cluster. If it is not the object will be requeued until it exists (or is otherwise retrievable).
	debug(log, gateway, "gathering the gateway publish service") // this will also be done by the validating webhook, this is a fallback
	svc, err := r.determineServiceForGateway(ctx, serviceRef)
	if err != nil {
		log.Error(err, "could not determine service for gateway", "namespace", gateway.Namespace, "name", gateway.Name)
		return ctrl.Result{Requeue: true}, err
	}

	// set the Gateway as scheduled to indicate that validation is complete and reconciliation work
	// on the object is ready to begin.
	if !isGatewayScheduled(gateway) {
		info(log, gateway, "marking gateway as accepted")
		scheduledCondition := metav1.Condition{
			Type:               string(gatewayv1beta1.GatewayConditionAccepted),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: gateway.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             string(gatewayv1beta1.GatewayReasonAccepted),
			Message:            "this unmanaged gateway has been picked up by the controller and will be processed",
		}
		setGatewayCondition(gateway, scheduledCondition)
		return ctrl.Result{}, r.Status().Update(ctx, pruneGatewayStatusConds(gateway))
	}

	// When deployed on Kubernetes Kong can not be relied on for the address data needed for Gateway because
	// it's commonly deployed in a way agnostic to the container network (e.g. it's simply configured with
	// 0.0.0.0 as the address for its listeners internally). In order to get addresses we have to derive them
	// from the Kubernetes Service which will also give us all the L4 information about the proxy. From there
	// we can use that L4 information to derive the higher level TLS and HTTP,GRPC, e.t.c. information from
	// the data-plane's metadata.
	debug(log, gateway, "determining listener configurations from publish service")
	kongAddresses, kongListeners, err := r.determineL4ListenersFromService(log, svc)
	if err != nil {
		return ctrl.Result{}, err
	}
	debug(log, gateway, "determining listener configurations from Kong data-plane")
	kongListeners, err = r.determineListenersFromDataPlane(ctx, svc, kongListeners)
	if err != nil {
		return ctrl.Result{}, err
	}

	if !reflect.DeepEqual(gateway.Spec.Addresses, kongAddresses) {
		debug(log, gateway, "updating addresses to match Kong proxy Service")
		gateway.Spec.Addresses = kongAddresses
		if err := r.Update(ctx, gateway); err != nil {
			if apierrors.IsConflict(err) {
				// if there's a conflict that's normal just requeue to retry, no need to make noise.
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil // dont requeue here because spec update will trigger new reconciliation
	}

	// the ReferenceGrants need to be retrieved to ensure that all gateway listeners reference
	// TLS secrets they are granted for
	referenceGrantList := &gatewayv1beta1.ReferenceGrantList{}
	if r.EnableReferenceGrant {
		if err := r.Client.List(ctx, referenceGrantList); err != nil {
			return ctrl.Result{}, err
		}
	}

	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2559 check cross-Gateway compatibility
	// When all Listeners on all Gateways were derived from Kong's configuration, they were guaranteed to be compatible
	// because they were all identical, though there may have been some ambiguity re de facto merged Gateways that
	// used different allowedRoutes. We only merged allowedRoutes within a single Gateway, but merged all Gateways into
	// a single set of shared listens. We lack knowledge of whether this is compatible with user intent, and it may
	// be incompatible with the spec, so we should consider evaluating cross-Gateway compatibility and raising error
	// conditions in the event of a problem
	listenerStatuses, err := getListenerStatus(ctx, gateway, kongListeners, referenceGrantList.Items, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// once specification matches the reference Service, all that's left to do is ensure that the
	// Gateway status reflects the spec. As the status is simply a mirror of the Service, this is
	// a given and we can simply update spec to status.
	debug(log, gateway, "updating the gateway status if necessary")
	isChanged, err := r.updateAddressesAndListenersStatus(ctx, gateway, listenerStatuses)
	if err != nil {
		if apierrors.IsConflict(err) {
			// if there's a conflict that's normal just requeue to retry, no need to make noise.
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	if isChanged {
		debug(log, gateway, "gateways status updated")
		return ctrl.Result{}, nil
	}

	info(log, gateway, "gateway provisioning complete")
	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// Gateway Controller - Listener Derivation Methods
// -----------------------------------------------------------------------------

var (
	// supportedKinds indicates which gateway kinds are supported by this implementation.
	supportedKinds = []Kind{
		Kind("HTTPRoute"),
		Kind("TCPRoute"),
		Kind("UDPRoute"),
		Kind("TLSRoute"),
	}

	// supportedRouteGroupKinds indicates the full kinds with GVK that are supported by this implementation.
	supportedRouteGroupKinds []gatewayv1beta1.RouteGroupKind
)

func init() {
	// gather the supported RouteGroupKinds for the Gateway listeners
	group := gatewayv1beta1.Group(gatewayv1beta1.GroupName)
	for _, supportedKind := range supportedKinds {
		supportedRouteGroupKinds = append(supportedRouteGroupKinds, gatewayv1beta1.RouteGroupKind{
			Group: &group,
			Kind:  supportedKind,
		})
	}
}

// determineServiceForGateway provides the "publish service" (aka the proxy Service) object which
// will be used to populate unmanaged gateways.
func (r *GatewayReconciler) determineServiceForGateway(ctx context.Context, ref string) (*corev1.Service, error) {
	// currently the gateway controller ONLY supports service references that correspond with the --publish-service
	// provided to the controller manager via flags when operating on unmanaged gateways. This constraint may
	// be loosened in later iterations if there is need.
	if ref != r.PublishService {
		return nil, fmt.Errorf("service ref %s did not match controller manager ref %s", ref, r.PublishService)
	}

	// retrieve the service for the kong gateway
	svc := &corev1.Service{}
	return svc, r.Client.Get(ctx, r.publishServiceRef, svc)
}

// determineL4ListenersFromService generates L4 addresses and listeners for a
// unmanaged Gateway provided the service to reference from.
func (r *GatewayReconciler) determineL4ListenersFromService(
	log logr.Logger,
	svc *corev1.Service,
) (
	[]GatewayAddress,
	[]Listener,
	error,
) {
	// if there are no clusterIPs available yet then this service
	// is still being provisioned so we will need to wait.
	if len(svc.Spec.ClusterIPs) < 1 {
		return nil, nil, fmt.Errorf("gateway service %s/%s is not yet ready (no cluster IPs provisioned)", svc.Namespace, svc.Name)
	}

	// take var copies of the address types so we can take pointers to them
	gatewayIPAddrType := gatewayv1beta1.IPAddressType
	gatewayHostAddrType := gatewayv1beta1.HostnameAddressType

	// for all service types we're going to capture the ClusterIP
	addresses := make([]GatewayAddress, 0, len(svc.Spec.ClusterIPs))
	listeners := make([]Listener, 0, len(svc.Spec.Ports))
	protocolToRouteGroupKind := map[corev1.Protocol]gatewayv1beta1.RouteGroupKind{
		corev1.ProtocolTCP: {Group: &gatewayV1beta1Group, Kind: Kind("TCPRoute")},
		corev1.ProtocolUDP: {Group: &gatewayV1beta1Group, Kind: Kind("UDPRoute")},
	}

	for _, port := range svc.Spec.Ports {
		listeners = append(listeners, Listener{
			Name:     (SectionName)(port.Name),
			Protocol: (ProtocolType)(port.Protocol),
			Port:     (PortNumber)(port.Port),
			AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
				Kinds: []gatewayv1beta1.RouteGroupKind{
					protocolToRouteGroupKind[port.Protocol],
				},
			},
		})
	}

	// for LoadBalancer service types we'll also capture the LB IP or Host
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		// if the loadbalancer IPs/Hosts haven't been provisioned yet
		// the service is not ready and we'll need to wait.
		if len(svc.Status.LoadBalancer.Ingress) < 1 {
			info(log, svc, "gateway service is type LoadBalancer but has not yet been provisioned: LoadBalancer IPs can not be added to the Gateway's addresses until this is resolved")
			return addresses, listeners, nil
		}

		// otherwise gather any IPs or Hosts provisioned for the LoadBalancer
		// and record them as Gateway Addresses. The LoadBalancer addresses
		// are pre-pended to the address list to make them prominent, as they
		// are often the most common address used for traffic.
		for _, ingress := range svc.Status.LoadBalancer.Ingress {
			if ingress.IP != "" {
				addresses = append([]GatewayAddress{{
					Type:  &gatewayIPAddrType,
					Value: ingress.IP,
				}}, addresses...)
			}
			if ingress.Hostname != "" {
				addresses = append([]GatewayAddress{{
					Type:  &gatewayHostAddrType,
					Value: ingress.Hostname,
				}}, addresses...)
			}
		}
	}

	// the API server transforms a Gateway with a zero-length address slice into a Gateway with a nil address slice
	// the value we return here needs to match the transformed Gateway, as otherwise the controller will always see
	// that the Gateway needs a status update and will never mark it ready
	if len(addresses) == 0 {
		addresses = nil
	}

	return addresses, listeners, nil
}

// determineListenersFromDataPlane takes a list of Gateway listeners and references
// them against the data-plane to determine any higher level protocol (TLS, HTTP)
// configured for them.
func (r *GatewayReconciler) determineListenersFromDataPlane(
	ctx context.Context,
	svc *corev1.Service,
	listeners []Listener,
) ([]Listener, error) {
	// gather the proxy and stream listeners from the data-plane and map them
	// to their respective ports (which will be the targetPorts of the proxy
	// Service in Kubernetes).
	proxyListeners, streamListeners, err := r.DataplaneClient.Listeners(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve listeners from the data-plane: %w", err)
	}
	proxyListenersMap := make(map[int]kong.ProxyListener)
	for _, listener := range proxyListeners {
		proxyListenersMap[listener.Port] = listener
	}
	streamListenersMap := make(map[int]kong.StreamListener)
	for _, listener := range streamListeners {
		streamListenersMap[listener.Port] = listener
	}

	// map servicePorts to targetPorts in order to identify the data-plane
	// proxy_listeners and stream_listeners for the service by their ports.
	portMapper := make(map[int]int)
	for _, port := range svc.Spec.Ports {
		portMapper[int(port.Port)] = port.TargetPort.IntValue()
	}

	// upgrade existing L4 listeners with any higher level protocols that
	// are configured for them in the data-plane.
	upgradedListeners := make([]Listener, 0, len(listeners))
	for _, listener := range listeners {
		if streamListener, ok := streamListenersMap[portMapper[int(listener.Port)]]; ok {
			if streamListener.SSL {
				listener.Protocol = gatewayv1beta1.TLSProtocolType
				listener.AllowedRoutes = &gatewayv1beta1.AllowedRoutes{
					Kinds: []gatewayv1beta1.RouteGroupKind{
						{Group: &gatewayV1beta1Group, Kind: (Kind)("TLSRoute")},
					},
				}
			}
		}
		if proxyListener, ok := proxyListenersMap[portMapper[int(listener.Port)]]; ok {
			if proxyListener.SSL {
				listener.Protocol = gatewayv1beta1.HTTPSProtocolType
				listener.AllowedRoutes = &gatewayv1beta1.AllowedRoutes{
					Kinds: []gatewayv1beta1.RouteGroupKind{
						{Group: &gatewayV1beta1Group, Kind: (Kind)("HTTPRoute")},
					},
				}
			} else {
				listener.Protocol = gatewayv1beta1.HTTPProtocolType
				listener.AllowedRoutes = &gatewayv1beta1.AllowedRoutes{
					Kinds: []gatewayv1beta1.RouteGroupKind{
						{Group: &gatewayV1beta1Group, Kind: (Kind)("HTTPRoute")},
					},
				}
			}
		}
		upgradedListeners = append(upgradedListeners, listener)
	}

	return upgradedListeners, nil
}

// -----------------------------------------------------------------------------
// Gateway Controller - Private Object Update Methods
// -----------------------------------------------------------------------------

// updateAddressesAndListenersStatus updates a unmanaged gateway's status with new addresses and listeners.
// If the addresses and listeners provided are the same as what exists, it is assumed that reconciliation is complete and a Ready condition is posted.
func (r *GatewayReconciler) updateAddressesAndListenersStatus(
	ctx context.Context,
	gateway *gatewayv1beta1.Gateway,
	listenerStatuses []gatewayv1beta1.ListenerStatus,
) (bool, error) {
	if !isGatewayReady(gateway) {
		gateway.Status.Listeners = listenerStatuses
		gateway.Status.Addresses = gateway.Spec.Addresses
		readyCondition := metav1.Condition{
			Type:               string(gatewayv1beta1.GatewayConditionReady),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: gateway.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             string(gatewayv1beta1.GatewayReasonReady),
			Message:            "addresses and listeners for the Gateway resource were successfully updated",
		}
		setGatewayCondition(gateway, readyCondition)
		return true, r.Status().Update(ctx, pruneGatewayStatusConds(gateway))
	}
	if !reflect.DeepEqual(gateway.Status.Listeners, listenerStatuses) {
		gateway.Status.Listeners = listenerStatuses
		return true, r.Status().Update(ctx, gateway)
	}
	return false, nil
}
