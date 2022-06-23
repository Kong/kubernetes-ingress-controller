package gateway

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane"
)

// -----------------------------------------------------------------------------
// Vars & Consts
// -----------------------------------------------------------------------------

var (
	// ManagedGatewaysUnsupported is an error used whenever a failure occurs
	// due to a Gateway that is not properly configured for unmanaged mode.
	ManagedGatewaysUnsupported = fmt.Errorf("invalid gateway spec: managed gateways are not currently supported") //nolint:revive
	gatewayV1alpha2Group       = gatewayv1alpha2.Group(gatewayv1alpha2.GroupName)
)

// -----------------------------------------------------------------------------
// Gateway Controller - GatewayReconciler
// -----------------------------------------------------------------------------

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct { //nolint:revive
	client.Client

	Log             logr.Logger
	Scheme          *runtime.Scheme
	DataplaneClient *dataplane.KongClient

	PublishService  string
	WatchNamespaces []string

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
	})
	if err != nil {
		return err
	}

	// watch Gateway objects, filtering out any Gateways which are not configured with
	// a supported GatewayClass controller name.
	if err := c.Watch(
		&source.Kind{Type: &gatewayv1alpha2.Gateway{}},
		&handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(r.gatewayHasMatchingGatewayClass),
	); err != nil {
		return err
	}

	// watch for updates to gatewayclasses, if any gateway classes change, enqueue
	// reconciliation for all supported gateway objects which reference it.
	if err := c.Watch(
		&source.Kind{Type: &gatewayv1alpha2.GatewayClass{}},
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

	// start the required gatewayclass controller as well
	gwcCTRL := &GatewayClassReconciler{
		Client: r.Client,
		Log:    r.Log.WithName("V1Alpha2GatewayClass"),
		Scheme: r.Scheme,
	}

	return gwcCTRL.SetupWithManager(mgr)
}

// -----------------------------------------------------------------------------
// Gateway Controller - Watch Predicates
// -----------------------------------------------------------------------------

// gatewayHasMatchingGatewayClass is a watch predicate which filters out reconciliation events for
// gateway objects which aren't supported by this controller.
func (r *GatewayReconciler) gatewayHasMatchingGatewayClass(obj client.Object) bool {
	gateway, ok := obj.(*gatewayv1alpha2.Gateway)
	if !ok {
		r.Log.Error(fmt.Errorf("unexpected object type in gateway watch predicates"), "expected", "*gatewayv1alpha2.Gateway", "found", reflect.TypeOf(obj))
		return false
	}
	gatewayClass := &gatewayv1alpha2.GatewayClass{}
	if err := r.Client.Get(context.Background(), client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, gatewayClass); err != nil {
		r.Log.Error(err, "could not retrieve gatewayclass", "gatewayclass", gateway.Spec.GatewayClassName)
		return false
	}
	return gatewayClass.Spec.ControllerName == ControllerName
}

// gatewayClassMatchesController is a watch predicate which filters out events for gatewayclasses which
// aren't configured with the required ControllerName, e.g. they are not supported by this controller.
func (r *GatewayReconciler) gatewayClassMatchesController(obj client.Object) bool {
	gatewayClass, ok := obj.(*gatewayv1alpha2.GatewayClass)
	if !ok {
		r.Log.Error(fmt.Errorf("unexpected object type in gatewayclass watch predicates"), "expected", "*gatewayv1alpha2.GatewayClass", "found", reflect.TypeOf(obj))
		return false
	}
	return gatewayClass.Spec.ControllerName == ControllerName
}

// listGatewaysForGatewayClass is a watch predicate which finds all the gateway objects reference
// by a gatewayclass to enqueue them for reconciliation. This is generally used when a GatewayClass
// is updated to ensure that idle gateways are initialized when their gatewayclass becomes available.
func (r *GatewayReconciler) listGatewaysForGatewayClass(gatewayClass client.Object) []reconcile.Request {
	gateways := &gatewayv1alpha2.GatewayList{}
	if err := r.Client.List(context.Background(), gateways); err != nil {
		r.Log.Error(err, "failed to list gateways for gatewayclass in watch", "gatewayclass", gatewayClass.GetName())
		return nil
	}
	return reconcileGatewaysIfClassMatches(gatewayClass, gateways.Items)
}

// listGatewaysForService is a watch predicate which finds all the gateway objects which use
// GatewayClasses supported by this controller and are configured for the same service via
// unmanaged mode and enqueues them for reconciliation. This is generally used to ensure
// all gateways are updated when the service gets updated with new listeners.
func (r *GatewayReconciler) listGatewaysForService(svc client.Object) (recs []reconcile.Request) {
	gateways := &gatewayv1alpha2.GatewayList{}
	if err := r.Client.List(context.Background(), gateways); err != nil {
		r.Log.Error(err, "failed to list gateways for service in watch predicates", "service")
		return
	}
	for _, gateway := range gateways.Items {
		gatewayClass := &gatewayv1alpha2.GatewayClass{}
		if err := r.Client.Get(context.Background(), types.NamespacedName{Name: string(gateway.Spec.GatewayClassName)}, gatewayClass); err != nil {
			r.Log.Error(err, "failed to retrieve gateway class in watch predicates", "gatewayclass", gateway.Spec.GatewayClassName)
			return
		}
		if isGatewayInClassAndUnmanaged(gatewayClass, gateway) {
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

// -----------------------------------------------------------------------------
// Gateway Controller - Reconciliation
// -----------------------------------------------------------------------------

//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways,verbs=get;list;watch;update
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/status,verbs=get;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("V1Alpha2Gateway", req.NamespacedName)

	// gather the gateway object based on the reconciliation trigger. It's possible for the object
	// to be gone at this point in which case it will be ignored.
	gateway := new(gatewayv1alpha2.Gateway)
	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		if errors.IsNotFound(err) {
			debug(log, gateway, "reconciliation triggered but gateway does not exist, ignoring")
			return ctrl.Result{Requeue: false}, nil
		}
		return ctrl.Result{Requeue: true}, r.DataplaneClient.DeleteObject(gateway)
	}
	debug(log, gateway, "processing gateway")

	// though our watch configuration eliminates reconciliation of unsupported gateways it's
	// technically possible for the gatewayclass configuration of a gateway to change in
	// the interim while the object has been queued for reconciliation. This double check
	// reduces reconciliation operations that would occur on old information.
	debug(log, gateway, "verifying gatewayclass")
	gwc := &gatewayv1alpha2.GatewayClass{}
	if err := r.Client.Get(ctx, client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, gwc); err != nil {
		debug(log, gateway, "could not retrieve gatewayclass for gateway", "gatewayclass", string(gateway.Spec.GatewayClassName))
		if err := r.DataplaneClient.DeleteObject(gateway); err != nil {
			debug(log, gateway, "failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, gateway, "ensured object was removed from the data-plane (if ever present)")
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(gateway)
	}
	if gwc.Spec.ControllerName != ControllerName {
		debug(log, gateway, "unsupported gatewayclass controllername, ignoring", "gatewayclass", gwc.Name, "controllername", gwc.Spec.ControllerName)
		if err := r.DataplaneClient.DeleteObject(gateway); err != nil {
			debug(log, gateway, "failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, gateway, "ensured object was removed from the data-plane (if ever present)")
		return ctrl.Result{}, r.DataplaneClient.DeleteObject(gateway)
	}

	// if there's any deletion timestamp on the object, we can simply ignore it. At this point
	// with unmanaged mode being the only option supported there are no finalizers and
	// the object should be cleaned up by GC promptly.
	debug(log, gateway, "checking deletion timestamp")
	if gateway.DeletionTimestamp != nil {
		debug(log, gateway, "gateway is being deleted, ignoring")
		return ctrl.Result{Requeue: false}, nil
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
	}
	return result, err
}

// reconcileUnmanagedGateway reconciles a Gateway that is configured for unmanaged mode,
// this mode will extract the Addresses and Listeners for the Gateway from the Kubernetes Service
// used for the Kong Gateway in the pre-existing deployment.
func (r *GatewayReconciler) reconcileUnmanagedGateway(ctx context.Context, log logr.Logger, gateway *gatewayv1alpha2.Gateway) (ctrl.Result, error) {
	// currently this controller supports only unmanaged gateway mode, we need to verify
	// any Gateway object that comes to us is configured appropriately, and if not reject it
	// with a clear status condition and message.
	debug(log, gateway, "validating management mode for gateway") // this will also be done by the validating webhook, this is a fallback
	unmanagedAnnotation := annotations.AnnotationPrefix + annotations.GatewayUnmanagedAnnotation
	existingGatewayEnabled, ok := annotations.ExtractUnmanagedGatewayMode(gateway.GetAnnotations())

	// allow for Gateway resources to be configured with "true" in place of the publish service
	// reference as a placeholder to automatically populate the annotation with the namespace/name
	// that was provided to the controller manager via --publish-service.
	debug(log, gateway, "initializing admin service annotation if unset")
	if !ok || existingGatewayEnabled == "true" { // true is a placeholder which triggers auto-initialization of the ref
		debug(log, gateway, fmt.Sprintf("a placeholder value was provided for %s, adding the default service ref %s", unmanagedAnnotation, r.PublishService))
		if gateway.Annotations == nil {
			gateway.Annotations = make(map[string]string)
		}
		gateway.Annotations[unmanagedAnnotation] = r.PublishService
		return ctrl.Result{}, r.Update(ctx, gateway)
	}

	// validation check of the Gateway to ensure that the publish service is actually available
	// in the cluster. If it is not the object will be requeued until it exists (or is otherwise retrievable).
	debug(log, gateway, "gathering the gateway publish service") // this will also be done by the validating webhook, this is a fallback
	svc, err := r.determineServiceForGateway(ctx, existingGatewayEnabled)
	if err != nil {
		log.Error(err, "could not determine service for gateway", "namespace", gateway.Namespace, "name", gateway.Name)
		return ctrl.Result{Requeue: true}, err
	}

	// set the Gateway as scheduled to indicate that validation is complete and reconciliation work
	// on the object is ready to begin.
	info(log, gateway, "marking gateway as scheduled")
	if !isGatewayScheduled(gateway) {
		gateway.Status.Conditions = append(gateway.Status.Conditions, metav1.Condition{
			Type:               string(gatewayv1alpha2.GatewayConditionScheduled),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: gateway.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             string(gatewayv1alpha2.GatewayReasonScheduled),
			Message:            "this unmanaged gateway has been picked up by the controller and will be processed",
		})
		return ctrl.Result{}, r.Status().Update(ctx, pruneGatewayStatusConds(gateway))
	}

	// When deployed on Kubernetes Kong can not be relied on for the address data needed for Gateway because
	// it's commonly deployed in a way agnostic to the container network (e.g. it's simply configured with
	// 0.0.0.0 as the address for its listeners internally). In order to get addresses we have to derive them
	// from the Kubernetes Service which will also give us all the L4 information about the proxy. From there
	// we can use that L4 information to derive the higher level TLS and HTTP,GRPC, e.t.c. information from
	// the data-plane's // metadata.
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
			if errors.IsConflict(err) {
				// if there's a conflict that's normal just requeue to retry, no need to make noise.
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil // dont requeue here because spec update will trigger new reconciliation
	}

	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2559 check cross-Gateway compatibility
	// When all Listeners on all Gateways were derived from Kong's configuration, they were guaranteed to be compatible
	// because they were all identical, though there may have been some ambiguity re de facto merged Gateways that
	// used different allowedRoutes. We only merged allowedRoutes within a single Gateway, but merged all Gateways into
	// a single set of shared listens. We lack knowledge of whether this is compatible with user intent, and it may
	// be incompatible with the spec, so we should consider evaluating cross-Gateway compatibility and raising error
	// conditions in the event of a problem
	listenerStatuses := getListenerStatus(gateway, kongListeners)

	// once specification matches the reference Service, all that's left to do is ensure that the
	// Gateway status reflects the spec. As the status is simply a mirror of the Service, this is
	// a given and we can simply update spec to status.
	debug(log, gateway, "updating the gateway status if necessary")
	isChanged, err := r.updateAddressesAndListenersStatus(ctx, gateway, listenerStatuses)
	if err != nil {
		if errors.IsConflict(err) {
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
	// supportedKinds indicates which gateway kinds are supported by this implementation
	supportedKinds = []gatewayv1alpha2.Kind{
		gatewayv1alpha2.Kind("HTTPRoute"),
	}

	// supportedRouteGroupKinds indicates the full kinds with GVK that are supported by this implementation
	supportedRouteGroupKinds []gatewayv1alpha2.RouteGroupKind
)

func init() {
	// gather the supported RouteGroupKinds for the Gateway listeners
	group := gatewayv1alpha2.Group(gatewayv1alpha2.GroupName)
	for _, supportedKind := range supportedKinds {
		supportedRouteGroupKinds = append(supportedRouteGroupKinds, gatewayv1alpha2.RouteGroupKind{
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
	[]gatewayv1alpha2.GatewayAddress,
	[]gatewayv1alpha2.Listener,
	error,
) {
	// if there are no clusterIPs available yet then this service
	// is still being provisioned so we will need to wait.
	if len(svc.Spec.ClusterIPs) < 1 {
		return nil, nil, fmt.Errorf("gateway service %s/%s is not yet ready (no cluster IPs provisioned)", svc.Namespace, svc.Name)
	}

	// take var copies of the address types so we can take pointers to them
	gatewayIPAddrType := gatewayv1alpha2.IPAddressType
	gatewayHostAddrType := gatewayv1alpha2.HostnameAddressType

	// for all service types we're going to capture the ClusterIP
	addresses := make([]gatewayv1alpha2.GatewayAddress, 0, len(svc.Spec.ClusterIPs))
	listeners := make([]gatewayv1alpha2.Listener, 0, len(svc.Spec.Ports))
	protocolToRouteGroupKind := map[corev1.Protocol]gatewayv1alpha2.RouteGroupKind{
		corev1.ProtocolTCP: {Group: &gatewayV1alpha2Group, Kind: gatewayv1alpha2.Kind("TCPRoute")},
		corev1.ProtocolUDP: {Group: &gatewayV1alpha2Group, Kind: gatewayv1alpha2.Kind("UDPRoute")},
	}

	for _, port := range svc.Spec.Ports {
		listeners = append(listeners, gatewayv1alpha2.Listener{
			Name:     gatewayv1alpha2.SectionName(port.Name),
			Protocol: gatewayv1alpha2.ProtocolType(port.Protocol),
			Port:     gatewayv1alpha2.PortNumber(port.Port),
			AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
				Kinds: []gatewayv1alpha2.RouteGroupKind{
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
				addresses = append([]gatewayv1alpha2.GatewayAddress{{
					Type:  &gatewayIPAddrType,
					Value: ingress.IP,
				}}, addresses...)
			}
			if ingress.Hostname != "" {
				addresses = append([]gatewayv1alpha2.GatewayAddress{{
					Type:  &gatewayHostAddrType,
					Value: ingress.Hostname,
				}}, addresses...)
			}
		}
	}

	return addresses, listeners, nil
}

// determineListenersFromDataPlane takes a list of Gateway listeners and references
// them against the data-plane to determine any higher level protocol (TLS, HTTP)
// configured for them.
func (r *GatewayReconciler) determineListenersFromDataPlane(ctx context.Context, svc *corev1.Service, listeners []gatewayv1alpha2.Listener) ([]gatewayv1alpha2.Listener, error) {
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
	upgradedListeners := make([]gatewayv1alpha2.Listener, 0, len(listeners))
	for _, listener := range listeners {
		if streamListener, ok := streamListenersMap[portMapper[int(listener.Port)]]; ok {
			if streamListener.SSL {
				listener.Protocol = gatewayv1alpha2.TLSProtocolType
				listener.AllowedRoutes = &gatewayv1alpha2.AllowedRoutes{
					Kinds: []gatewayv1alpha2.RouteGroupKind{
						{Group: &gatewayV1alpha2Group, Kind: gatewayv1alpha2.Kind("TLSRoute")},
					},
				}
			}
		}
		if proxyListener, ok := proxyListenersMap[portMapper[int(listener.Port)]]; ok {
			if proxyListener.SSL {
				listener.Protocol = gatewayv1alpha2.HTTPSProtocolType
				listener.AllowedRoutes = &gatewayv1alpha2.AllowedRoutes{
					Kinds: []gatewayv1alpha2.RouteGroupKind{
						{Group: &gatewayV1alpha2Group, Kind: gatewayv1alpha2.Kind("HTTPRoute")},
					},
				}
			} else {
				listener.Protocol = gatewayv1alpha2.HTTPProtocolType
				listener.AllowedRoutes = &gatewayv1alpha2.AllowedRoutes{
					Kinds: []gatewayv1alpha2.RouteGroupKind{
						{Group: &gatewayV1alpha2Group, Kind: gatewayv1alpha2.Kind("HTTPRoute")},
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
	gateway *gatewayv1alpha2.Gateway,
	listenerStatuses []gatewayv1alpha2.ListenerStatus,
) (bool, error) {
	if !isGatewayReady(gateway) {
		gateway.Status.Listeners = listenerStatuses
		gateway.Status.Addresses = gateway.Spec.Addresses
		gateway.Status.Conditions = append(gateway.Status.Conditions, metav1.Condition{
			Type:               string(gatewayv1alpha2.GatewayConditionReady),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: gateway.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             string(gatewayv1alpha2.GatewayReasonReady),
			Message:            "addresses and listeners for the Gateway resource were successfully updated",
		})
		return true, r.Status().Update(ctx, pruneGatewayStatusConds(gateway))
	}
	return false, nil
}

// areAllowedRoutesConsistentByProtocol returns an error if a set of listeners includes multiple listeners for the same
// protocol that do not use the same AllowedRoutes filters. Kong does not support limiting routes to a specific listen:
// all routes are always served on all listens compatible with their protocol. As such, while we can filter the routes
// we ingest, if we ingest routes from two incompatible filters, we will combine them into a single proxy configuration
// It may be possible to write a new Kong plugin that checks the inbound port/address to de facto apply listen-based
// filters in the future.
func areAllowedRoutesConsistentByProtocol(listeners []gatewayv1alpha2.Listener) bool {
	allowedByProtocol := make(map[gatewayv1alpha2.ProtocolType]gatewayv1alpha2.AllowedRoutes)
	for _, listener := range listeners {
		var allowed gatewayv1alpha2.AllowedRoutes
		var exists bool
		if allowed, exists = allowedByProtocol[listener.Protocol]; !exists {
			allowedByProtocol[listener.Protocol] = *listener.AllowedRoutes
		} else {
			if !reflect.DeepEqual(allowed, *listener.AllowedRoutes) {
				return false
			}
		}
	}
	return true
}
