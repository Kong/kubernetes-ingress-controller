package gateway

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/samber/mo"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers"
	ctrlref "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/reference"
	ctrlutils "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/utils"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// -----------------------------------------------------------------------------
// Gateway Controller - GatewayReconciler
// -----------------------------------------------------------------------------

// GatewayReconciler reconciles a Gateway object.
type GatewayReconciler struct { //nolint:revive
	client.Client

	Log             logr.Logger
	Scheme          *runtime.Scheme
	DataplaneClient controllers.DataPlane

	WatchNamespaces  []string
	CacheSyncTimeout time.Duration

	ReferenceIndexers ctrlref.CacheIndexers

	PublishServiceRef    k8stypes.NamespacedName
	PublishServiceUDPRef mo.Option[k8stypes.NamespacedName]

	// AddressOverrides are addresses to use in Gateway status instead of the PublishServiceRef addresses.
	AddressOverrides []string
	// AddressOverridesUDP are addresses to use in Gateway status instead of the PublishServiceUDPRef addresses.
	AddressOverridesUDP []string

	// If enableReferenceGrant is true, controller will watch ReferenceGrants
	// to invalidate or allow cross-namespace TLSConfigs in gateways.
	// It's resolved on SetupWithManager call.
	enableReferenceGrant bool

	// If GatewayNN is set,
	// only resources managed by the specified Gateway are reconciled.
	GatewayNN controllers.OptionalNamespacedName
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// verify that the PublishService was configured properly
	if r.PublishServiceRef.Name == "" || r.PublishServiceRef.Namespace == "" {
		return fmt.Errorf("publish service must be configured")
	}

	// We're verifying whether ReferenceGrant CRD is installed at setup of the GatewayReconciler
	// to decide whether we should run additional ReferenceGrant watch and handle ReferenceGrants
	// when reconciling Gateways.
	// Once the GatewayReconciler is set up without ReferenceGrant, there's no possibility to enable
	// ReferenceGrant handling again in this reconciler at runtime.
	r.enableReferenceGrant = ctrlutils.CRDExists(mgr.GetRESTMapper(), schema.GroupVersionResource{
		Group:    gatewayv1beta1.GroupVersion.Group,
		Version:  gatewayv1beta1.GroupVersion.Version,
		Resource: "referencegrants",
	})

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
		source.Kind(mgr.GetCache(), &gatewayapi.Gateway{}),
		&handler.EnqueueRequestForObject{},
		predicate.NewPredicateFuncs(r.gatewayHasMatchingGatewayClass),
	); err != nil {
		return err
	}

	// watch for updates to gatewayclasses, if any gateway classes change, enqueue
	// reconciliation for all supported gateway objects which reference it.
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &gatewayapi.GatewayClass{}),
		handler.EnqueueRequestsFromMapFunc(r.listGatewaysForGatewayClass),
		predicate.NewPredicateFuncs(r.gatewayClassMatchesController),
	); err != nil {
		return err
	}

	// if an update to the gateway service occurs, we need to make sure to trigger
	// reconciliation on all Gateway objects referenced by it (in the most common
	// deployments this will be a single Gateway).
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &corev1.Service{}),
		handler.EnqueueRequestsFromMapFunc(r.listGatewaysForService),
		predicate.NewPredicateFuncs(r.isGatewayService),
	); err != nil {
		return err
	}

	// if a HTTPRoute gets accepted by a Gateway, we need to make sure to trigger
	// reconciliation on the gateway, as we need to update the number of attachedRoutes.
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &gatewayapi.HTTPRoute{}),
		handler.EnqueueRequestsFromMapFunc(r.listGatewaysForHTTPRoute),
	); err != nil {
		return err
	}

	// watch ReferenceGrants, which may invalidate or allow cross-namespace TLSConfigs
	if r.enableReferenceGrant {
		if err := c.Watch(
			source.Kind(mgr.GetCache(), &gatewayapi.ReferenceGrant{}),
			handler.EnqueueRequestsFromMapFunc(r.listReferenceGrantsForGateway),
			predicate.NewPredicateFuncs(referenceGrantHasGatewayFrom),
		); err != nil {
			return err
		}
	}

	// start the required gatewayclass controller as well
	gwcCTRL := &GatewayClassReconciler{
		Client:           r.Client,
		Log:              r.Log.WithName(strings.ToUpper(gatewayapi.V1GroupVersion) + "GatewayClass"),
		Scheme:           r.Scheme,
		CacheSyncTimeout: r.CacheSyncTimeout,
	}

	return gwcCTRL.SetupWithManager(mgr)
}

// SetLogger sets the logger.
func (r *GatewayReconciler) SetLogger(l logr.Logger) {
	r.Log = l
}

// -----------------------------------------------------------------------------
// Gateway Controller - Watch Predicates
// -----------------------------------------------------------------------------

// gatewayHasMatchingGatewayClass is a watch predicate which filters out reconciliation events for
// gateway objects which aren't supported by this controller or not using an unmanaged GatewayClass.
func (r *GatewayReconciler) gatewayHasMatchingGatewayClass(obj client.Object) bool {
	gateway, ok := obj.(*gatewayapi.Gateway)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"Gateway watch predicate received unexpected object type",
			"expected", "*gatewayapi.Gateway", "found", reflect.TypeOf(obj),
		)
		return false
	}

	// If the flag `--gateway-to-reconcile` is set, KIC will only reconcile the specified gateway.
	// https://github.com/Kong/kubernetes-ingress-controller/issues/5322
	if gatewayToReconcile, ok := r.GatewayNN.Get(); ok {
		if gatewayToReconcile.Namespace != gateway.Namespace || gatewayToReconcile.Name != gateway.Name {
			return false
		}
	}

	gatewayClass := &gatewayapi.GatewayClass{}
	if err := r.Client.Get(context.Background(), client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, gatewayClass); err != nil {
		r.Log.Error(err, "Could not retrieve gatewayclass", "gatewayclass", gateway.Spec.GatewayClassName)
		return false
	}
	return isGatewayClassControlled(gatewayClass)
}

// gatewayClassMatchesController is a watch predicate which filters out events for gatewayclasses which
// aren't configured with the required ControllerName or not annotated as unmanaged.
func (r *GatewayReconciler) gatewayClassMatchesController(obj client.Object) bool {
	gatewayClass, ok := obj.(*gatewayapi.GatewayClass)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"Gatewayclass watch predicate received unexpected object type",
			"expected", "*gatewayapi.GatewayClass", "found", reflect.TypeOf(obj),
		)
		return false
	}
	return isGatewayClassControlled(gatewayClass)
}

// listGatewaysForGatewayClass is a watch predicate which finds all the gateway objects reference
// by a gatewayclass to enqueue them for reconciliation. This is generally used when a GatewayClass
// is updated to ensure that idle gateways are initialized when their gatewayclass becomes available.
func (r *GatewayReconciler) listGatewaysForGatewayClass(ctx context.Context, gatewayClass client.Object) []reconcile.Request {
	gateways := &gatewayapi.GatewayList{}
	if err := r.Client.List(ctx, gateways); err != nil {
		r.Log.Error(err, "Failed to list gateways for gatewayclass in watch", "gatewayclass", gatewayClass.GetName())
		return nil
	}
	return reconcileGatewaysIfClassMatches(gatewayClass, gateways.Items)
}

// listReferenceGrantsForGateway is a watch predicate which finds all Gateways mentioned in a From clause for a
// ReferenceGrant.
func (r *GatewayReconciler) listReferenceGrantsForGateway(ctx context.Context, obj client.Object) []reconcile.Request {
	grant, ok := obj.(*gatewayapi.ReferenceGrant)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"Referencegrant watch predicate received unexpected object type",
			"expected", "*gatewayapi.ReferenceGrant", "found", reflect.TypeOf(obj),
		)
		return nil
	}
	gateways := &gatewayapi.GatewayList{}
	if err := r.Client.List(ctx, gateways); err != nil {
		r.Log.Error(err, "Failed to list gateways in watch", "referencegrant", grant.Name)
		return nil
	}
	recs := []reconcile.Request{}
	for _, gateway := range gateways.Items {
		for _, from := range grant.Spec.From {
			if string(from.Namespace) == gateway.Namespace &&
				from.Kind == gatewayapi.Kind("Gateway") &&
				from.Group == gatewayapi.Group("gateway.networking.k8s.io") {
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

// listGatewaysForService is a watch predicate which finds all the gateway objects which use
// GatewayClasses supported by this controller and are configured for the same service via
// unmanaged mode and enqueues them for reconciliation. This is generally used to ensure
// all gateways are updated when the service gets updated with new listeners.
func (r *GatewayReconciler) listGatewaysForService(ctx context.Context, svc client.Object) (recs []reconcile.Request) {
	gateways := &gatewayapi.GatewayList{}
	if err := r.Client.List(ctx, gateways); err != nil {
		r.Log.Error(err, "Failed to list gateways for service in watch predicates", "service", svc)
		return
	}
	for _, gateway := range gateways.Items {
		gatewayClass := &gatewayapi.GatewayClass{}
		if err := r.Client.Get(ctx, k8stypes.NamespacedName{Name: string(gateway.Spec.GatewayClassName)}, gatewayClass); err != nil {
			r.Log.Error(err, "Failed to retrieve gateway class in watch predicates", "gatewayclass", gateway.Spec.GatewayClassName)
			return
		}
		if isGatewayClassControlled(gatewayClass) {
			recs = append(recs, reconcile.Request{
				NamespacedName: k8stypes.NamespacedName{
					Namespace: gateway.Namespace,
					Name:      gateway.Name,
				},
			})
		}
	}
	return
}

// listGatewaysForHTTPRoute retrieves all the gateways referenced as parents by the HTTPRoute.
func (r *GatewayReconciler) listGatewaysForHTTPRoute(_ context.Context, obj client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gatewayapi.HTTPRoute)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"HTTPRoute watch predicate received unexpected object type",
			"expected", "*gatewayapi.HTTPRoute", "found", reflect.TypeOf(obj),
		)
		return nil
	}
	recs := []reconcile.Request{}
	for _, gateway := range routeAcceptedByGateways(httpRoute.Namespace, httpRoute.Status.Parents) {
		recs = append(recs, reconcile.Request{
			NamespacedName: gateway,
		})
	}

	return recs
}

// isGatewayService is a watch predicate that filters out events for objects that aren't
// the gateway service referenced by --publish-service or --publish-service-udp.
func (r *GatewayReconciler) isGatewayService(obj client.Object) bool {
	isPublishService := fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName()) == r.PublishServiceRef.String()
	isUDPPublishService := r.PublishServiceUDPRef.IsPresent() &&
		fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName()) == r.PublishServiceUDPRef.MustGet().String()

	return isPublishService || isUDPPublishService
}

func referenceGrantHasGatewayFrom(obj client.Object) bool {
	grant, ok := obj.(*gatewayapi.ReferenceGrant)
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
	log := r.Log.WithValues("GatewayV1Gateway", req.NamespacedName)

	if gatewayToReconcile, ok := r.GatewayNN.Get(); ok {
		if req.Namespace != gatewayToReconcile.Namespace || req.Name != gatewayToReconcile.Name {
			r.Log.V(util.DebugLevel).Info("The request does not match the specified Gateway and will be skipped.", "gateway", gatewayToReconcile.String())
			return ctrl.Result{}, nil
		}
	}

	// gather the gateway object based on the reconciliation trigger. It's possible for the object
	// to be gone at this point in which case it will be ignored.
	gateway := new(gatewayapi.Gateway)
	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		if apierrors.IsNotFound(err) {
			gateway.Namespace = req.Namespace
			gateway.Name = req.Name
			// delete reference relationships where the gateway is the referrer.
			err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, gateway)
			if err != nil {
				return ctrl.Result{}, err
			}
			debug(log, gateway, "Reconciliation triggered but gateway does not exist, deleting it in dataplane")
			return ctrl.Result{}, r.DataplaneClient.DeleteObject(gateway)
		}
		return ctrl.Result{}, err
	}

	debug(log, gateway, "Processing gateway")

	// though our watch configuration eliminates reconciliation of unsupported gateways it's
	// technically possible for the gatewayclass configuration of a gateway to change in
	// the interim while the object has been queued for reconciliation. This double check
	// reduces reconciliation operations that would occur on old information.
	debug(log, gateway, "Verifying gatewayclass")
	gwc := &gatewayapi.GatewayClass{}
	if err := r.Client.Get(ctx, client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, gwc); err != nil {
		debug(log, gateway, "Could not retrieve gatewayclass for gateway", "gatewayclass", string(gateway.Spec.GatewayClassName))
		// delete reference relationships where the gateway is the referrer, as we will not process the gateway.
		err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, gateway)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err := r.DataplaneClient.DeleteObject(gateway); err != nil {
			debug(log, gateway, "Failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, gateway, "Ensured gateway was removed from the data-plane (if ever present)")
		return ctrl.Result{}, nil
	}
	if gwc.Spec.ControllerName != GetControllerName() {
		debug(log, gateway, "Unsupported gatewayclass controllername, ignoring", "gatewayclass", gwc.Name, "controllername", gwc.Spec.ControllerName)
		// delete reference relationships where the gateway is the referrer, as we will not process the gateway.
		err := ctrlref.DeleteReferencesByReferrer(r.ReferenceIndexers, r.DataplaneClient, gateway)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err := r.DataplaneClient.DeleteObject(gateway); err != nil {
			debug(log, gateway, "Failed to delete object from data-plane, requeuing")
			return ctrl.Result{}, err
		}
		debug(log, gateway, "Ensured gateway was removed from the data-plane (if ever present)")
		return ctrl.Result{}, nil
	}

	// if there's any deletion timestamp on the object, we can simply ignore it. At this point
	// with unmanaged mode being the only option supported there are no finalizers and
	// the object should be cleaned up by GC promptly.
	debug(log, gateway, "Checking deletion timestamp")
	if gateway.DeletionTimestamp != nil {
		debug(log, gateway, "Gateway is being deleted, ignoring")
		return ctrl.Result{Requeue: false}, nil
	}

	// ensure that the GatewayClass matches the ControllerName.
	// This check has already been performed by predicates, but we need to ensure this condition
	// as the reconciliation loop may be triggered by objects in which predicates we
	// cannot check the ControllerName (e.g., ReferenceGrants).
	if !isGatewayClassControlled(gwc) {
		return reconcile.Result{}, nil
	}

	if isGatewayClassUnmanaged(gwc.Annotations) {
		// The Gateway has to be reconciled by KIC only if it is unmanaged.
		if result, err := r.reconcileUnmanagedGateway(ctx, log, gateway); err != nil {
			return result, err
		}
	}

	// If the Gateway has been accepted (by KIC or the managing controller), the dataplane update must be performed.
	if isGatewayAccepted(gateway) {
		if err := r.DataplaneClient.UpdateObject(gateway); err != nil {
			debug(log, gateway, "Failed to update object in data-plane, requeueing")
			return ctrl.Result{}, err
		}

		referredSecretNames := listSecretNamesReferredByGateway(gateway)
		if err := ctrlref.UpdateReferencesToSecret(
			ctx, r.Client, r.ReferenceIndexers, r.DataplaneClient,
			gateway, referredSecretNames); err != nil {
			if apierrors.IsNotFound(err) {
				return ctrl.Result{Requeue: true}, nil
			}
		}
	}
	return ctrl.Result{}, nil
}

// reconcileUnmanagedGateway reconciles a Gateway that is configured for unmanaged mode,
// this mode will extract the Addresses and Listeners for the Gateway from the Kubernetes Service
// used for the Kong Gateway in the pre-existing deployment.
func (r *GatewayReconciler) reconcileUnmanagedGateway(ctx context.Context, log logr.Logger, gateway *gatewayapi.Gateway) (ctrl.Result, error) {
	// currently this controller supports only unmanaged gateway mode, we need to verify
	// any Gateway object that comes to us is configured appropriately, and if not reject it
	// with a clear status condition and message.
	debug(log, gateway, "Validating management mode for gateway") // this will also be done by the validating webhook, this is a fallback

	// enforce the service reference as the annotation value for the key UnmanagedGateway.
	debug(log, gateway, "Initializing admin service annotation if unset")
	if len(annotations.ExtractGatewayPublishService(gateway.Annotations)) == 0 {
		services := []string{r.PublishServiceRef.String()}

		// UDP service is optional.
		if udpRef, ok := r.PublishServiceUDPRef.Get(); ok {
			services = append(services, udpRef.String())
		}

		debug(log, gateway, fmt.Sprintf("No publish service annotation, setting it to proxy services %s", services))
		if gateway.Annotations == nil {
			gateway.Annotations = map[string]string{}
		}
		annotations.UpdateGatewayPublishService(gateway.Annotations, services)
		return ctrl.Result{}, r.Update(ctx, gateway)
	}

	serviceRefs := annotations.ExtractGatewayPublishService(gateway.Annotations)
	// Validation check of the Gateway to ensure that the ingress service is actually available
	// in the cluster. If it is not the object will be requeued until it exists (or is otherwise retrievable).
	debug(log, gateway, "Gathering the gateway publish service") // this will also be done by the validating webhook, this is a fallback
	var gatewayServices []*corev1.Service
	for _, ref := range serviceRefs {
		r.Log.V(util.DebugLevel).Info("Determining service for ref", "ref", ref)
		svc, err := r.determineServiceForGateway(ctx, ref)
		if err != nil {
			const annotation = annotations.AnnotationPrefix + annotations.GatewayPublishServiceKey
			log.Error(
				err,
				fmt.Sprintf("One of publish services defined in Gateway's %q annotation didn't match controller manager's configuration", annotation),
				"namespace", gateway.Namespace,
				"name", gateway.Name,
				"service", ref,
			)
			return ctrl.Result{}, err
		}
		if svc != nil {
			gatewayServices = append(gatewayServices, svc)
		}
	}

	// set the Gateway as scheduled to indicate that validation is complete and reconciliation work
	// on the object is ready to begin.
	if !isGatewayAccepted(gateway) {
		info(log, gateway, "Marking gateway as accepted")
		acceptedCondition := metav1.Condition{
			Type:               string(gatewayapi.GatewayConditionAccepted),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: gateway.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             string(gatewayapi.GatewayReasonAccepted),
			Message:            "this unmanaged gateway has been picked up by the controller and will be processed",
		}
		setGatewayCondition(gateway, acceptedCondition)
		programmedCondition := metav1.Condition{
			Type:               string(gatewayapi.GatewayConditionProgrammed),
			Status:             metav1.ConditionFalse,
			ObservedGeneration: gateway.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             string(gatewayapi.GatewayReasonPending),
		}
		setGatewayCondition(gateway, programmedCondition)
		return ctrl.Result{}, r.Status().Update(ctx, pruneGatewayStatusConds(gateway))
	}

	// When deployed on Kubernetes Kong can not be relied on for the address data needed for Gateway because
	// it's commonly deployed in a way agnostic to the container network (e.g. it's simply configured with
	// 0.0.0.0 as the address for its listeners internally). In order to get addresses we have to derive them
	// from the Kubernetes Service which will also give us all the L4 information about the proxy. From there
	// we can use that L4 information to derive the higher level TLS and HTTP,GRPC, e.t.c. information from
	// the data-plane's metadata.
	debug(log, gateway, "Determining listener configurations from publish services")
	var combinedAddresses []gatewayapi.GatewayStatusAddress
	var combinedListeners []gatewayapi.Listener
	for _, svc := range gatewayServices {
		kongAddresses, kongListeners, err := r.determineL4ListenersFromService(log, svc)
		if err != nil {
			return ctrl.Result{}, err
		}
		debug(log, gateway, "Determining listener configurations from Kong data-plane")
		kongListeners, err = r.determineListenersFromDataPlane(ctx, svc, kongListeners)
		if err != nil {
			return ctrl.Result{}, err
		}
		combinedAddresses = append(combinedAddresses, kongAddresses...)
		combinedListeners = append(combinedListeners, kongListeners...)
	}

	// This handles PublishStatusAddress(UDP) override config support, which allows users to set an arbitrary string to
	// use in place of the proxy Service addresses, usually because there's another proxy in front of Kong and the
	// addresses associated with the proxy Service aren't actually where you want to direct external clients.
	if len(r.AddressOverrides)+len(r.AddressOverridesUDP) > 0 {
		combinedOverrideAddresses := append(r.AddressOverrides, r.AddressOverridesUDP...)
		overrides := make([]gatewayapi.GatewayStatusAddress, len(combinedOverrideAddresses))
		for i, stringAddr := range combinedOverrideAddresses {
			addr := gatewayapi.GatewayStatusAddress{
				Value: stringAddr,
				Type:  lo.ToPtr(gatewayapi.HostnameAddressType),
			}
			if ip := net.ParseIP(stringAddr); ip != nil {
				addr.Type = lo.ToPtr(gatewayapi.IPAddressType)
			}
			overrides[i] = addr
		}
		combinedAddresses = overrides
	}

	// the ReferenceGrants need to be retrieved to ensure that all gateway listeners reference
	// TLS secrets they are granted for
	referenceGrantList := &gatewayapi.ReferenceGrantList{}
	if r.enableReferenceGrant {
		if err := r.Client.List(ctx, referenceGrantList); err != nil {
			return ctrl.Result{}, err
		}
	}

	listenerStatuses, err := getListenerStatus(ctx, gateway, combinedListeners, referenceGrantList.Items, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// once specification matches the reference Service, all that's left to do is ensure that the
	// Gateway status reflects the spec. As the status is simply a mirror of the Service, this is
	// a given and we can simply update spec to status.
	debug(log, gateway, "Updating the gateway status if necessary")
	isChanged, err := r.updateAddressesAndListenersStatus(ctx, gateway, listenerStatuses, combinedAddresses)
	if err != nil {
		if apierrors.IsConflict(err) {
			// if there's a conflict that's normal just requeue to retry, no need to make noise.
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	if isChanged {
		debug(log, gateway, "Gateway status updated")
		return ctrl.Result{}, nil
	}

	info(log, gateway, "Gateway provisioning complete")
	return ctrl.Result{}, nil
}

// -----------------------------------------------------------------------------
// Gateway Controller - Listener Derivation Methods
// -----------------------------------------------------------------------------

var (
	// supportedKinds indicates which gateway kinds are supported by this implementation.
	supportedKinds = []gatewayapi.Kind{
		gatewayapi.Kind("HTTPRoute"),
		gatewayapi.Kind("TCPRoute"),
		gatewayapi.Kind("UDPRoute"),
		gatewayapi.Kind("TLSRoute"),
	}

	// supportedRouteGroupKinds indicates the full kinds with GVK that are supported by this implementation.
	supportedRouteGroupKinds []gatewayapi.RouteGroupKind
)

func init() {
	// gather the supported RouteGroupKinds for the Gateway listeners
	group := gatewayapi.Group(gatewayv1beta1.GroupName)
	for _, supportedKind := range supportedKinds {
		supportedRouteGroupKinds = append(supportedRouteGroupKinds, gatewayapi.RouteGroupKind{
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

	var name k8stypes.NamespacedName
	switch {
	case ref == r.PublishServiceRef.String():
		name = r.PublishServiceRef
	case r.PublishServiceUDPRef.IsPresent() && ref == r.PublishServiceUDPRef.MustGet().String():
		name = r.PublishServiceUDPRef.MustGet()
	default:
		configuredServiceRefs := []string{fmt.Sprintf("%q", r.PublishServiceRef)}
		if udpRef, ok := r.PublishServiceUDPRef.Get(); ok {
			configuredServiceRefs = append(configuredServiceRefs, fmt.Sprintf("%q [udp]", udpRef))
		}
		return nil, fmt.Errorf("publish service reference %q from Gateway's annotations did not match configured controller manager's publish services (%s)",
			ref, strings.Join(configuredServiceRefs, ", "))
	}

	// retrieve the service for the kong gateway
	svc := &corev1.Service{}
	if name.Name == "" && name.Namespace == "" {
		r.Log.V(util.DebugLevel).Info("Service not configured, discarding it", "ref", ref)
		return nil, nil
	}
	if err := r.Client.Get(ctx, name, svc); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("publish service %q couldn't be found: %w", name, err)
		}
		return nil, fmt.Errorf("publish service %q couldn't be retrieved: %w", name, err)
	}
	return svc, nil
}

// determineL4ListenersFromService generates L4 addresses and listeners for a
// unmanaged Gateway provided the service to reference from.
func (r *GatewayReconciler) determineL4ListenersFromService(
	log logr.Logger,
	svc *corev1.Service,
) (
	[]gatewayapi.GatewayStatusAddress,
	[]gatewayapi.Listener,
	error,
) {
	// if there are no clusterIPs available yet then this service
	// is still being provisioned so we will need to wait.
	if len(svc.Spec.ClusterIPs) < 1 {
		return nil, nil, fmt.Errorf("gateway service %s/%s is not yet ready (no cluster IPs provisioned)", svc.Namespace, svc.Name)
	}

	// take var copies of the address types so we can take pointers to them
	gatewayIPAddrType := gatewayapi.IPAddressType
	gatewayHostAddrType := gatewayapi.HostnameAddressType

	// for all service types we're going to capture the ClusterIP
	addresses := make([]gatewayapi.GatewayStatusAddress, 0, len(svc.Spec.ClusterIPs))
	listeners := make([]gatewayapi.Listener, 0, len(svc.Spec.Ports))
	protocolToRouteGroupKind := map[corev1.Protocol]gatewayapi.RouteGroupKind{
		corev1.ProtocolTCP: {Group: lo.ToPtr(gatewayapi.V1Group), Kind: gatewayapi.Kind("TCPRoute")},
		corev1.ProtocolUDP: {Group: lo.ToPtr(gatewayapi.V1Group), Kind: gatewayapi.Kind("UDPRoute")},
	}

	for _, port := range svc.Spec.Ports {
		listeners = append(listeners, gatewayapi.Listener{
			Name:     (gatewayapi.SectionName)(port.Name),
			Protocol: (gatewayapi.ProtocolType)(port.Protocol),
			Port:     (gatewayapi.PortNumber)(port.Port),
			AllowedRoutes: &gatewayapi.AllowedRoutes{
				Kinds: []gatewayapi.RouteGroupKind{
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
			info(log, svc, "Gateway service is type LoadBalancer but has not yet been provisioned: LoadBalancer IPs can not be added to the Gateway's addresses until this is resolved")
			return addresses, listeners, nil
		}

		// otherwise gather any IPs or Hosts provisioned for the LoadBalancer
		// and record them as Gateway Addresses. The LoadBalancer addresses
		// are pre-pended to the address list to make them prominent, as they
		// are often the most common address used for traffic.
		for _, ingress := range svc.Status.LoadBalancer.Ingress {
			if ingress.IP != "" {
				addresses = append([]gatewayapi.GatewayStatusAddress{{
					Type:  &gatewayIPAddrType,
					Value: ingress.IP,
				}}, addresses...)
			}
			if ingress.Hostname != "" {
				addresses = append([]gatewayapi.GatewayStatusAddress{{
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
	listeners []gatewayapi.Listener,
) ([]gatewayapi.Listener, error) {
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
	upgradedListeners := make([]gatewayapi.Listener, 0, len(listeners))
	for _, listener := range listeners {
		if streamListener, ok := streamListenersMap[portMapper[int(listener.Port)]]; ok {
			if streamListener.SSL {
				listener.Protocol = gatewayapi.TLSProtocolType
				listener.AllowedRoutes = &gatewayapi.AllowedRoutes{
					Kinds: []gatewayapi.RouteGroupKind{
						{Group: lo.ToPtr(gatewayapi.V1Group), Kind: (gatewayapi.Kind)("TLSRoute")},
					},
				}
			}
		}
		if proxyListener, ok := proxyListenersMap[portMapper[int(listener.Port)]]; ok {
			if proxyListener.SSL {
				listener.Protocol = gatewayapi.HTTPSProtocolType
				listener.AllowedRoutes = &gatewayapi.AllowedRoutes{
					Kinds: []gatewayapi.RouteGroupKind{
						{Group: lo.ToPtr(gatewayapi.V1Group), Kind: (gatewayapi.Kind)("HTTPRoute")},
					},
				}
			} else {
				listener.Protocol = gatewayapi.HTTPProtocolType
				listener.AllowedRoutes = &gatewayapi.AllowedRoutes{
					Kinds: []gatewayapi.RouteGroupKind{
						{Group: lo.ToPtr(gatewayapi.V1Group), Kind: (gatewayapi.Kind)("HTTPRoute")},
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
// If the addresses and listeners provided are the same as what exists, it is assumed that reconciliation is complete and a Programmed condition is posted.
func (r *GatewayReconciler) updateAddressesAndListenersStatus(
	ctx context.Context,
	gateway *gatewayapi.Gateway,
	listenerStatuses []gatewayapi.ListenerStatus,
	addresses []gatewayapi.GatewayStatusAddress,
) (bool, error) {
	if !isGatewayProgrammed(gateway) {
		gateway.Status.Listeners = listenerStatuses
		gateway.Status.Addresses = addresses
		programmedCondition := metav1.Condition{
			Type:               string(gatewayapi.GatewayConditionProgrammed),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: gateway.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             string(gatewayapi.GatewayReasonProgrammed),
		}
		setGatewayCondition(gateway, programmedCondition)
		return true, r.Status().Update(ctx, pruneGatewayStatusConds(gateway))
	}
	if !reflect.DeepEqual(gateway.Status.Listeners, listenerStatuses) {
		gateway.Status.Listeners = listenerStatuses
		return true, r.Status().Update(ctx, gateway)
	}
	return false, nil
}
