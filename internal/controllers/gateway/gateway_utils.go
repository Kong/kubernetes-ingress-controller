package gateway

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
)

// -----------------------------------------------------------------------------
// Gateway Utilities
// -----------------------------------------------------------------------------

const (
	// maxConds is the maximum number of status conditions a Gateway can have at one time.
	maxConds = 8
)

// setGatewayCondition sets the condition with specified type in gateway status
// to expected condition in newCondition.
// if the gateway status does not contain a condition with that type, add one more condition.
// if the gateway status contains condition(s) with the type, then replace with the new condition.
func setGatewayCondition(gateway *Gateway, newCondition metav1.Condition) {
	newConditions := []metav1.Condition{}
	for _, condition := range gateway.Status.Conditions {
		if condition.Type != newCondition.Type {
			newConditions = append(newConditions, condition)
		}
	}
	newConditions = append(newConditions, newCondition)
	gateway.Status.Conditions = newConditions
}

// isGatewayScheduled returns boolean whether or not the gateway object was scheduled
// previously by the gateway controller.
func isGatewayScheduled(gateway *Gateway) bool {
	return util.CheckCondition(
		gateway.Status.Conditions,
		util.ConditionType(gatewayv1beta1.GatewayConditionAccepted),
		util.ConditionReason(gatewayv1beta1.GatewayReasonAccepted),
		metav1.ConditionTrue,
		gateway.Generation,
	)
}

// isGatewayReady returns boolean whether the ready condition exists
// for the given gateway object if it matches the currently known generation of that object.
func isGatewayReady(gateway *Gateway) bool {
	return util.CheckCondition(
		gateway.Status.Conditions,
		util.ConditionType(gatewayv1beta1.GatewayConditionReady),
		util.ConditionReason(gatewayv1beta1.GatewayReasonReady),
		metav1.ConditionTrue,
		gateway.Generation,
	)
}

// isObjectUnmanaged returns boolean if the object is configured
// for unmanaged mode.
func isObjectUnmanaged(anns map[string]string) bool {
	annotationValue := annotations.ExtractUnmanagedGatewayClassMode(anns)
	return annotationValue != ""
}

// isGatewayClassControlledAndUnmanaged returns boolean if the GatewayClass
// is controlled by this controller and is configured for unmanaged mode.
func isGatewayClassControlledAndUnmanaged(gatewayClass *GatewayClass) bool {
	return gatewayClass.Spec.ControllerName == GetControllerName() && isObjectUnmanaged(gatewayClass.Annotations)
}

// getRefFromPublishService splits a publish service string in the format namespace/name into a types.NamespacedName
// and verifies the contents producing an error if they don't match namespace/name format.
func getRefFromPublishService(publishService string) (types.NamespacedName, error) {
	publishServiceSplit := strings.SplitN(publishService, "/", 3)
	if len(publishServiceSplit) != 2 {
		return types.NamespacedName{}, fmt.Errorf("--publish-service expected in format 'namespace/name' but got %s", publishService)
	}
	return types.NamespacedName{
		Namespace: publishServiceSplit[0],
		Name:      publishServiceSplit[1],
	}, nil
}

// pruneGatewayStatusConds cleans out old status conditions if the Gateway currently has more
// status conditions set than the 8 maximum allowed by the Kubernetes API.
func pruneGatewayStatusConds(gateway *Gateway) *Gateway {
	if len(gateway.Status.Conditions) > maxConds {
		gateway.Status.Conditions = gateway.Status.Conditions[len(gateway.Status.Conditions)-maxConds:]
	}
	return gateway
}

// reconcileGatewaysIfClassMatches is a filter function to convert a list of gateways into a list
// of reconciliation requests for those gateways based on which match the given class.
func reconcileGatewaysIfClassMatches(gatewayClass client.Object, gateways []Gateway) (recs []reconcile.Request) {
	for _, gateway := range gateways {
		if string(gateway.Spec.GatewayClassName) == gatewayClass.GetName() {
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

// list namespaced names of secrets referred by the gateway.
func listSecretNamesReferredByGateway(gateway *gatewayv1beta1.Gateway) map[types.NamespacedName]struct{} {
	nsNames := make(map[types.NamespacedName]struct{})

	for _, listener := range gateway.Spec.Listeners {
		if listener.TLS == nil {
			continue
		}

		for _, certRef := range listener.TLS.CertificateRefs {
			if certRef.Group != nil && *certRef.Group != corev1.GroupName {
				continue
			}

			if certRef.Kind != nil && *certRef.Kind != "Secret" {
				continue
			}

			refNamespace := gateway.Namespace
			if certRef.Namespace != nil {
				refNamespace = string(*certRef.Namespace)
			}

			nsNames[types.NamespacedName{
				Namespace: refNamespace,
				Name:      string(certRef.Name),
			}] = struct{}{}
		}
	}
	return nsNames
}

// extractListenerSpecFromGateway returns the spec of the listener with the given name.
// returns nil if the listener with given name is not found.
func extractListenerSpecFromGateway(gateway *gatewayv1beta1.Gateway, listenerName gatewayv1beta1.SectionName) *gatewayv1beta1.Listener {
	for i, l := range gateway.Spec.Listeners {
		if l.Name == listenerName {
			return &gateway.Spec.Listeners[i]
		}
	}
	return nil
}

type (
	protocolPortMap     map[ProtocolType]map[PortNumber]bool
	portProtocolMap     map[PortNumber]ProtocolType
	portHostnameMap     map[PortNumber]map[Hostname]bool
	listenerAttachedMap map[SectionName]int32
)

func buildKongPortMap(listens []Listener) protocolPortMap {
	p := make(map[ProtocolType]map[PortNumber]bool, len(listens))
	for _, listen := range listens {
		_, ok := p[listen.Protocol]
		if !ok {
			p[listen.Protocol] = map[PortNumber]bool{}
		}
		p[listen.Protocol][listen.Port] = true
	}
	return p
}

// initializeListenerMaps takes a Gateway and builds indices used in status updates and conflict detection. It returns
// empty maps from port to protocol to listener name and from port to hostnames, and a populated map from listener name
// to attached route count from their status.
func initializeListenerMaps(gateway *Gateway) (
	portProtocolMap,
	portHostnameMap,
	listenerAttachedMap,
) {
	portToProtocol := make(portProtocolMap, len(gateway.Status.Listeners))
	portToHostname := make(portHostnameMap, len(gateway.Status.Listeners))
	listenerToAttached := make(listenerAttachedMap, len(gateway.Status.Listeners))

	existingStatuses := make(map[SectionName]ListenerStatus,
		len(gateway.Status.Listeners))
	for _, listenerStatus := range gateway.Status.Listeners {
		existingStatuses[listenerStatus.Name] = listenerStatus
	}

	for _, listener := range gateway.Spec.Listeners {
		portToHostname[listener.Port] = make(map[Hostname]bool)
		if existingStatus, ok := existingStatuses[listener.Name]; ok {
			listenerToAttached[listener.Name] = existingStatus.AttachedRoutes
		} else {
			listenerToAttached[listener.Name] = 0
		}
	}
	return portToProtocol, portToHostname, listenerToAttached
}

func canSharePort(requested, existing ProtocolType) bool {
	switch requested {
	// TCP and UDP listeners must always use unique ports
	case gatewayv1beta1.TCPProtocolType, gatewayv1beta1.UDPProtocolType:
		return false
	// HTTPS and TLS Listeners can share ports with others of their type or the other TLS type
	// note that this is not actually possible in Kong: TLS is a stream listen and HTTPS is an http listen
	// however, this section implements the spec ignoring Kong's reality
	case gatewayv1beta1.HTTPSProtocolType:
		if existing == gatewayv1beta1.HTTPSProtocolType ||
			existing == gatewayv1beta1.TLSProtocolType {
			return true
		}
		return false
	case gatewayv1beta1.TLSProtocolType:
		if existing == gatewayv1beta1.HTTPSProtocolType ||
			existing == gatewayv1beta1.TLSProtocolType {
			return true
		}
		return false
	// HTTP Listeners can share ports with others of the same protocol only
	case gatewayv1beta1.HTTPProtocolType:
		if existing == gatewayv1beta1.HTTPProtocolType {
			return true
		}
		return false
	default:
		return false
	}
}

func getListenerStatus(
	ctx context.Context,
	gateway *Gateway,
	kongListens []Listener,
	referenceGrants []gatewayv1beta1.ReferenceGrant,
	client client.Client,
) ([]ListenerStatus, error) {
	statuses := make(map[SectionName]ListenerStatus, len(gateway.Spec.Listeners))
	portToProtocol, portToHostname, listenerToAttached := initializeListenerMaps(gateway)
	kongProtocolsToPort := buildKongPortMap(kongListens)
	conflictedPorts := make(map[PortNumber]bool, len(gateway.Spec.Listeners))
	conflictedHostnames := make(map[PortNumber]map[Hostname]bool, len(gateway.Spec.Listeners))

	// TODO we should check transition time rather than always nowing, which we do throughout the below
	// https://github.com/Kong/kubernetes-ingress-controller/issues/2556
	for _, listener := range gateway.Spec.Listeners {
		var hostname Hostname
		if listener.Hostname != nil {
			hostname = *listener.Hostname
		}
		status := ListenerStatus{
			Name:           listener.Name,
			Conditions:     []metav1.Condition{},
			SupportedKinds: getListenerSupportedRouteKinds(listener),
			// this has been populated by initializeListenerMaps()
			AttachedRoutes: listenerToAttached[listener.Name],
		}
		// TODO this only handles some Listener conditions and reasons as needed to check cross-listener compatibility
		// and unattachability due to missing Kong configuration. There are others available and it may be appropriate
		// for us to add them https://github.com/Kong/kubernetes-ingress-controller/issues/2558
		if _, ok := portToProtocol[listener.Port]; !ok {
			// unoccupied ports are free game
			portToProtocol[listener.Port] = listener.Protocol
			portToHostname[listener.Port][hostname] = true
		} else {
			if !canSharePort(listener.Protocol, portToProtocol[listener.Port]) {
				status.Conditions = append(status.Conditions, metav1.Condition{
					Type:               string(gatewayv1beta1.ListenerConditionConflicted),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1beta1.ListenerReasonProtocolConflict),
				})
				conflictedPorts[listener.Port] = true
			} else {
				// shareable ports determine conflicts by hostname
				// Each Listener within the group specifies a Hostname that is unique within the group.
				// As a special case, one Listener within a group may omit Hostname, in which case this Listener
				// matches when no other Listener matches.

				// TODO this only checks if a hostname is already used on a specific port, which is what the Gateway
				// spec requires. However, Kong does not actually implement HTTP route separation by port: Kong serves
				// all HTTP routes on all HTTP ports. Effectively, if you add an HTTP(S) Listener with hostname
				// example.com on port 8000, and your Kong instance has a proxy_listen with both port 8000 and 8200,
				// you have also added a phantom Listener for hostname example.com and port 8200, because Kong will
				// serve the route on both. See https://github.com/Kong/kubernetes-ingress-controller/issues/2606
				if conflictedHostnames[listener.Port] == nil {
					conflictedHostnames[listener.Port] = map[Hostname]bool{}
				}
				if _, exists := portToHostname[listener.Port][hostname]; !exists {
					portToHostname[listener.Port][hostname] = true
				} else {
					status.Conditions = append(status.Conditions, metav1.Condition{
						Type:               string(gatewayv1beta1.ListenerConditionConflicted),
						Status:             metav1.ConditionTrue,
						ObservedGeneration: gateway.Generation,
						LastTransitionTime: metav1.Now(),
						Reason:             string(gatewayv1beta1.ListenerReasonHostnameConflict),
					})
					conflictedHostnames[listener.Port][hostname] = true
				}
			}
		}

		// independent of conflict detection. for example, two TCP Listeners both requesting the same port that Kong
		// does not provide should be both Conflicted and Detached
		if len(kongProtocolsToPort[listener.Protocol]) == 0 {
			status.Conditions = append(status.Conditions, metav1.Condition{
				Type:               string(gatewayv1beta1.ListenerConditionAccepted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: gateway.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1beta1.ListenerReasonUnsupportedProtocol),
				Message:            "no Kong listen with the requested protocol is configured",
			})
		} else {
			if _, ok := kongProtocolsToPort[listener.Protocol][listener.Port]; !ok {
				status.Conditions = append(status.Conditions, metav1.Condition{
					Type:               string(gatewayv1beta1.ListenerConditionAccepted),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1beta1.ListenerReasonPortUnavailable),
					Message:            "no Kong listen with the requested protocol is configured for the requested port",
				})
			}
		}

		// finalize adding any general conditions
		// TODO these (and really the others too) do not account for the conditions maybe having already been present
		// we simply generate them from scratch each time and mark the current generation the observed generation,
		// whereas we should preserve the original observed generation
		// https://github.com/Kong/kubernetes-ingress-controller/issues/2556
		if len(status.Conditions) == 0 {
			// if we've gotten this far with no conditions, the listener is good to go
			status.Conditions = append(status.Conditions,
				metav1.Condition{
					Type:               string(gatewayv1beta1.ListenerConditionConflicted),
					Status:             metav1.ConditionFalse,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1beta1.ListenerReasonNoConflicts),
				},
				metav1.Condition{
					Type:               string(gatewayv1beta1.ListenerConditionReady),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1beta1.ListenerReasonReady),
					Message:            "the listener is ready and available for routing",
				},
			)
		} else {
			// any conditions we added above will prevent the Listener from becoming ready
			// unsure if we want to add the ready=false condition on a per-failure basis or use this else to just mark
			// it generic unready if we hit anything bad. do any failure conditions block readiness? do we care about
			// having distinct ready false messages, assuming we have more descriptive messages in the other conditions?
			status.Conditions = append(status.Conditions, metav1.Condition{
				Type:               string(gatewayv1beta1.ListenerConditionReady),
				Status:             metav1.ConditionFalse,
				ObservedGeneration: gateway.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayv1beta1.ListenerReasonPending),
				Message:            "the listener is not ready and cannot route requests",
			})
		}

		// consistent sort statuses to allow equality comparisons
		sort.Slice(status.Conditions, func(i, j int) bool {
			a := status.Conditions[i]
			b := status.Conditions[j]
			return fmt.Sprintf("%s%s%s%s", a.Type, a.Status, a.Reason, a.Message) <
				fmt.Sprintf("%s%s%s%s", b.Type, b.Status, b.Reason, b.Message)
		})
		statuses[listener.Name] = status
	}

	// any conflict applies to all listeners sharing the conflicted resource (see
	// https://github.com/Kong/kubernetes-ingress-controller/pull/2555#issuecomment-1154579046 for discussion)
	// if we encountered conflicts, we must strip the ready status we originally set
	for _, listener := range gateway.Spec.Listeners {
		var conflictReason string
		var resolvedRefReason string

		var hostname Hostname
		if listener.Hostname != nil {
			hostname = *listener.Hostname
		}
		// there's no filter for protocols that don't use Hostname, but this won't be populated from earlier for those
		if _, ok := conflictedHostnames[listener.Port][hostname]; ok {
			conflictReason = string(gatewayv1beta1.ListenerReasonHostnameConflict)
		}

		if _, ok := conflictedPorts[listener.Port]; ok {
			conflictReason = string(gatewayv1beta1.ListenerReasonProtocolConflict)
		}

		// If the listener uses TLS, we need to ensure that the gateway is granted to reference
		// all the secrets it references
		if listener.TLS != nil {
			resolvedRefReason = string(gatewayv1beta1.ListenerReasonResolvedRefs)
			for _, certRef := range listener.TLS.CertificateRefs {
				// if the certificate is in the same namespace of the gateway, no ReferenceGrant is needed
				if certRef.Namespace != nil && *certRef.Namespace != (Namespace)(gateway.Namespace) {
					// get the result of the certificate reference. If the returned reason is not successful, the loop
					// must be broken because the secret reference isn't granted
					resolvedRefReason = getReferenceGrantConditionReason(gateway.Namespace, certRef, referenceGrants)
					if resolvedRefReason != string(gatewayv1beta1.ListenerReasonResolvedRefs) {
						break
					}
				}

				// only secrets are supported as certificate references
				if (certRef.Group != nil && (*certRef.Group != "core" && *certRef.Group != "")) ||
					(certRef.Kind != nil && *certRef.Kind != "Secret") {
					resolvedRefReason = string(gatewayv1beta1.ListenerReasonInvalidCertificateRef)
					break
				}
				secret := &corev1.Secret{}
				secretNamespace := gateway.Namespace
				if certRef.Namespace != nil {
					secretNamespace = string(*certRef.Namespace)
				}
				if err := client.Get(ctx, types.NamespacedName{Namespace: secretNamespace, Name: string(certRef.Name)}, secret); err != nil {
					if !apierrors.IsNotFound(err) {
						return nil, err
					}
					resolvedRefReason = string(gatewayv1beta1.ListenerReasonInvalidCertificateRef)
				}
			}
		}

		newConditions := []metav1.Condition{}

		// if resolvedRefReason has any value, it means that the listener references a secret.
		// A ListenerConditionResolvedRefs condition must be set to reflect in the gateway status
		// the outcome of that reference (that means if the gateway is granted to access that secret)
		if resolvedRefReason != "" {
			conditionStatus := metav1.ConditionTrue
			message := "the listener is ready and available for routing"
			if resolvedRefReason != string(gatewayv1beta1.ListenerReasonResolvedRefs) {
				conditionStatus = metav1.ConditionFalse
				message = "the listener is not ready and cannot route requests"
			}
			newConditions = append(newConditions,
				metav1.Condition{
					Type:               string(gatewayv1beta1.ListenerConditionResolvedRefs),
					Status:             conditionStatus,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             resolvedRefReason,
				},
				metav1.Condition{
					Type:               string(gatewayv1beta1.ListenerConditionReady),
					Status:             conditionStatus,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1beta1.ListenerReasonReady),
					Message:            message,
				},
			)
		}

		if len(conflictReason) > 0 {
			for _, cond := range statuses[listener.Name].Conditions {
				// shut up linter, there's a default
				switch gatewayv1alpha2.ListenerConditionType(cond.Type) { //nolint:exhaustive
				case gatewayv1beta1.ListenerConditionReady, gatewayv1beta1.ListenerConditionConflicted:
					continue
				default:
					newConditions = append(newConditions, cond)
				}
			}
			newConditions = append(newConditions,
				metav1.Condition{
					Type:               string(gatewayv1beta1.ListenerConditionConflicted),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             conflictReason,
				},
				metav1.Condition{
					Type:               string(gatewayv1beta1.ListenerConditionReady),
					Status:             metav1.ConditionFalse,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayv1beta1.ListenerReasonReady),
					Message:            "the listener is not ready and cannot route requests",
				},
			)
		}
		if len(newConditions) > 0 {
			status := statuses[listener.Name]
			// consistent sort statuses to allow equality comparisons
			sort.Slice(newConditions, func(i, j int) bool {
				a := newConditions[i]
				b := newConditions[j]
				return fmt.Sprintf("%s%s%s%s", a.Type, a.Status, a.Reason, a.Message) <
					fmt.Sprintf("%s%s%s%s", b.Type, b.Status, b.Reason, b.Message)
			})
			status.Conditions = newConditions
			statuses[listener.Name] = status
		}
	}
	statusArray := []ListenerStatus{}
	for _, status := range statuses {
		statusArray = append(statusArray, status)
	}

	return statusArray, nil
}

// getReferenceGrantConditionReason gets a certRef belonging to a specific listener and a slice of referenceGrants.
func getReferenceGrantConditionReason(
	gatewayNamespace string,
	certRef gatewayv1beta1.SecretObjectReference,
	referenceGrants []gatewayv1beta1.ReferenceGrant,
) string {
	// no need to have this reference granted
	if certRef.Namespace == nil || *certRef.Namespace == (Namespace)(gatewayNamespace) {
		return string(gatewayv1beta1.ListenerReasonResolvedRefs)
	}

	certRefNamespace := string(*certRef.Namespace)
	for _, grant := range referenceGrants {
		// the grant must exist in the same namespace of the referenced resource
		if grant.Namespace != certRefNamespace {
			continue
		}
		for _, from := range grant.Spec.From {
			// we are interested only in grants for gateways that want to reference secrets
			if from.Group != gatewayV1beta1Group || from.Kind != "Gateway" {
				continue
			}
			if from.Namespace == gatewayv1alpha2.Namespace(gatewayNamespace) {
				for _, to := range grant.Spec.To {
					if (to.Group != "" && to.Group != "core") || to.Kind != "Secret" {
						continue
					}
					// if all the above conditions are satisfied, and the name of the referenced secret matches
					// the granted resource name, then return a reason "ResolvedRefs"
					if to.Name == nil || string(*to.Name) == string(certRef.Name) {
						return string(gatewayv1beta1.ListenerReasonResolvedRefs)
					}
				}
			}
		}
	}
	// if no grants have been found for the reference, return an "InvalidCertificateRef" reason
	return string(gatewayv1beta1.ListenerReasonRefNotPermitted)
}

// -----------------------------------------------------------------------------
// Gateway Utils - Watch Predicate Helpers
// -----------------------------------------------------------------------------

// isGatewayClassEventInClass produces a boolean whether or not a given event which contains
// one or more GatewayClass objects is supported by this controller according to those
// objects ControllerName.
func isGatewayClassEventInClass(log logr.Logger, watchEvent interface{}) bool {
	objs := make([]client.Object, 0, 2)
	switch e := watchEvent.(type) {
	case event.CreateEvent:
		objs = append(objs, e.Object)
	case event.DeleteEvent:
		objs = append(objs, e.Object)
	case event.GenericEvent:
		objs = append(objs, e.Object)
	case event.UpdateEvent:
		objs = append(objs, e.ObjectOld)
		objs = append(objs, e.ObjectNew)
	default:
		log.Error(fmt.Errorf("invalid type"), "received invalid event type in event handlers", "found", reflect.TypeOf(watchEvent))
		return false
	}

	for _, obj := range objs {
		gwc, ok := obj.(*GatewayClass)
		if !ok {
			log.Error(fmt.Errorf("invalid type"), "received invalid object type in event handlers", "expected", "GatewayClass", "found", reflect.TypeOf(obj))
			continue
		}
		if gwc.Spec.ControllerName == GetControllerName() {
			return true
		}
	}

	return false
}

// getListenerSupportedRouteKinds determines what RouteGroupKinds are supported by the Listener.
// If no AllowedRoutes.Kinds are specified for the Listener, the supported RouteGroupKind is derived directly
// from the Listener's Protocol.
// Otherwise, user specified AllowedRoutes.Kinds are used, filtered by the global Gateway supported kinds.
//
// TODO: Make ListenerStatus.SupportedKinds compliant with GW API specification
// https://github.com/Kong/kubernetes-ingress-controller/issues/3228
func getListenerSupportedRouteKinds(l gatewayv1beta1.Listener) []gatewayv1beta1.RouteGroupKind {
	if l.AllowedRoutes == nil || len(l.AllowedRoutes.Kinds) == 0 {
		switch string(l.Protocol) {
		case string(gatewayv1beta1.HTTPProtocolType), string(gatewayv1beta1.HTTPSProtocolType):
			return builder.NewRouteGroupKind().HTTPRoute().IntoSlice()
		case string(gatewayv1beta1.TCPProtocolType):
			return builder.NewRouteGroupKind().TCPRoute().IntoSlice()
		case string(gatewayv1beta1.UDPProtocolType):
			return builder.NewRouteGroupKind().UDPRoute().IntoSlice()
		case string(gatewayv1beta1.TLSProtocolType):
			return builder.NewRouteGroupKind().TLSRoute().IntoSlice()
		}
	}

	var supportedRGK []gatewayv1beta1.RouteGroupKind
	for _, gk := range l.AllowedRoutes.Kinds {
		if gk.Group != nil && *gk.Group == gatewayv1beta1.GroupName {
			_, ok := lo.Find(supportedKinds, func(k Kind) bool {
				return gk.Kind == k
			})
			if ok {
				supportedRGK = append(supportedRGK, gk)
			}
		}
	}

	return supportedRGK
}
