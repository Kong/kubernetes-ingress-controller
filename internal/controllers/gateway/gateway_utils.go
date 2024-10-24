package gateway

import (
	"context"
	"encoding/pem"
	"fmt"
	"reflect"
	"sort"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
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
func setGatewayCondition(gateway *gatewayapi.Gateway, newCondition metav1.Condition) {
	newConditions := []metav1.Condition{}
	for _, condition := range gateway.Status.Conditions {
		if condition.Type != newCondition.Type {
			newConditions = append(newConditions, condition)
		}
	}
	newConditions = append(newConditions, newCondition)
	gateway.Status.Conditions = newConditions
}

// isGatewayAccepted returns boolean whether or not the gateway object was accepted
// previously by the gateway controller.
func isGatewayAccepted(gateway *gatewayapi.Gateway) bool {
	return util.CheckCondition(
		gateway.Status.Conditions,
		util.ConditionType(gatewayapi.GatewayConditionAccepted),
		util.ConditionReason(gatewayapi.GatewayReasonAccepted),
		metav1.ConditionTrue,
		gateway.Generation,
	)
}

// isGatewayProgrammed returns boolean whether the Programmed condition exists
// for the given Gateway object and if it matches the currently known generation of that object.
func isGatewayProgrammed(gateway *gatewayapi.Gateway) bool {
	return util.CheckCondition(
		gateway.Status.Conditions,
		util.ConditionType(gatewayapi.GatewayConditionProgrammed),
		util.ConditionReason(gatewayapi.GatewayReasonProgrammed),
		metav1.ConditionTrue,
		gateway.Generation,
	)
}

// Warning: this function is used for both GatewayClasses and Gateways.
// The former uses "true" as the value, whereas the latter uses "namespace/service" CSVs for the proxy services.

// isGatewayClassUnmanaged returns boolean if the object is configured
// for unmanaged mode.
func isGatewayClassUnmanaged(anns map[string]string) bool {
	annotationValue := annotations.ExtractUnmanagedGatewayClassMode(anns)
	return annotationValue != ""
}

// isGatewayClassControlled returns boolean if the GatewayClass
// is controlled by this controller and is configured for unmanaged mode.
func isGatewayClassControlled(gatewayClass *gatewayapi.GatewayClass) bool {
	return gatewayClass.Spec.ControllerName == GetControllerName()
}

// pruneGatewayStatusConds cleans out old status conditions if the Gateway currently has more
// status conditions set than the 8 maximum allowed by the Kubernetes API.
func pruneGatewayStatusConds(gateway *gatewayapi.Gateway) *gatewayapi.Gateway {
	if len(gateway.Status.Conditions) > maxConds {
		gateway.Status.Conditions = gateway.Status.Conditions[len(gateway.Status.Conditions)-maxConds:]
	}
	return gateway
}

// reconcileGatewaysIfClassMatches is a filter function to convert a list of gateways into a list
// of reconciliation requests for those gateways based on which match the given class.
func reconcileGatewaysIfClassMatches(gatewayClass client.Object, gateways []gatewayapi.Gateway) (recs []reconcile.Request) {
	for _, gateway := range gateways {
		if string(gateway.Spec.GatewayClassName) == gatewayClass.GetName() {
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

// list namespaced names of secrets referred by the gateway.
func listSecretNamesReferredByGateway(gateway *gatewayapi.Gateway) map[k8stypes.NamespacedName]struct{} {
	nsNames := make(map[k8stypes.NamespacedName]struct{})

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

			nsNames[k8stypes.NamespacedName{
				Namespace: refNamespace,
				Name:      string(certRef.Name),
			}] = struct{}{}
		}
	}
	return nsNames
}

// extractListenerSpecFromGateway returns the spec of the listener with the given name.
// returns nil if the listener with given name is not found.
func extractListenerSpecFromGateway(gateway *gatewayapi.Gateway, listenerName gatewayapi.SectionName) *gatewayapi.Listener {
	for i, l := range gateway.Spec.Listeners {
		if l.Name == listenerName {
			return &gateway.Spec.Listeners[i]
		}
	}
	return nil
}

type (
	protocolPortMap map[gatewayapi.ProtocolType]map[gatewayapi.PortNumber]bool
	portProtocolMap map[gatewayapi.PortNumber]gatewayapi.ProtocolType
	portHostnameMap map[gatewayapi.PortNumber]map[gatewayapi.Hostname]bool
)

func buildKongPortMap(listens []gatewayapi.Listener) protocolPortMap {
	p := make(map[gatewayapi.ProtocolType]map[gatewayapi.PortNumber]bool, len(listens))
	for _, listen := range listens {
		_, ok := p[listen.Protocol]
		if !ok {
			p[listen.Protocol] = map[gatewayapi.PortNumber]bool{}
		}
		p[listen.Protocol][listen.Port] = true
	}
	return p
}

// initializeListenerMaps takes a Gateway and builds indices used in status updates and conflict detection. It returns
// empty maps from port to protocol to listener name and from port to hostnames, and a populated map from listener name
// to attached route count from their status.
func initializeListenerMaps(gateway *gatewayapi.Gateway) (
	portProtocolMap,
	portHostnameMap,
) {
	portToProtocol := make(portProtocolMap, len(gateway.Status.Listeners))
	portToHostname := make(portHostnameMap, len(gateway.Status.Listeners))

	existingStatuses := make(map[gatewayapi.SectionName]gatewayapi.ListenerStatus,
		len(gateway.Status.Listeners))
	for _, listenerStatus := range gateway.Status.Listeners {
		existingStatuses[listenerStatus.Name] = listenerStatus
	}

	for _, listener := range gateway.Spec.Listeners {
		portToHostname[listener.Port] = make(map[gatewayapi.Hostname]bool)
	}
	return portToProtocol, portToHostname
}

func canSharePort(requested, existing gatewayapi.ProtocolType) bool {
	switch requested {
	// TCP and UDP listeners must always use unique ports
	case gatewayapi.TCPProtocolType, gatewayapi.UDPProtocolType:
		return false
	// HTTPS and TLS Listeners can share ports with others of their type or the other TLS type
	// note that this is not actually possible in Kong: TLS is a stream listen and HTTPS is an http listen
	// however, this section implements the spec ignoring Kong's reality
	case gatewayapi.HTTPSProtocolType:
		if existing == gatewayapi.HTTPSProtocolType ||
			existing == gatewayapi.TLSProtocolType {
			return true
		}
		return false
	case gatewayapi.TLSProtocolType:
		if existing == gatewayapi.HTTPSProtocolType ||
			existing == gatewayapi.TLSProtocolType {
			return true
		}
		return false
	// HTTP Listeners can share ports with others of the same protocol only
	case gatewayapi.HTTPProtocolType:
		if existing == gatewayapi.HTTPProtocolType {
			return true
		}
		return false
	default:
		return false
	}
}

func getListenerStatus(
	ctx context.Context,
	gateway *gatewayapi.Gateway,
	kongListens []gatewayapi.Listener,
	referenceGrants []gatewayapi.ReferenceGrant,
	client client.Client,
) ([]gatewayapi.ListenerStatus, error) {
	statuses := make(map[gatewayapi.SectionName]gatewayapi.ListenerStatus, len(gateway.Spec.Listeners))
	portToProtocol, portToHostname := initializeListenerMaps(gateway)
	kongProtocolsToPort := buildKongPortMap(kongListens)
	conflictedPorts := make(map[gatewayapi.PortNumber]bool, len(gateway.Spec.Listeners))
	conflictedHostnames := make(map[gatewayapi.PortNumber]map[gatewayapi.Hostname]bool, len(gateway.Spec.Listeners))

	// TODO we should check transition time rather than always nowing, which we do throughout the below
	// https://github.com/Kong/kubernetes-ingress-controller/issues/2556
	for listenerIndex, listener := range gateway.Spec.Listeners {
		var hostname gatewayapi.Hostname
		if listener.Hostname != nil {
			hostname = *listener.Hostname
		}
		supportedkinds, ResolvedRefsReason := getListenerSupportedRouteKinds(listener)

		// If the listener uses TLS, we need to ensure that the gateway is granted to reference
		// all the secrets it references
		if listener.TLS != nil {
			tlsResolvedRefReason := string(gatewayapi.ListenerReasonResolvedRefs)
			for _, certRef := range listener.TLS.CertificateRefs {
				// if the certificate is in the same namespace of the gateway, no ReferenceGrant is needed
				if certRef.Namespace != nil && *certRef.Namespace != (gatewayapi.Namespace)(gateway.Namespace) {
					// get the result of the certificate reference. If the returned reason is not successful, the loop
					// must be broken because the secret reference isn't granted
					tlsResolvedRefReason = getReferenceGrantConditionReason(gateway.Namespace, certRef, referenceGrants)
					if tlsResolvedRefReason != string(gatewayapi.ListenerReasonResolvedRefs) {
						break
					}
				}

				// only secrets are supported as certificate references
				if (certRef.Group != nil && (*certRef.Group != "core" && *certRef.Group != "")) ||
					(certRef.Kind != nil && *certRef.Kind != "Secret") {
					tlsResolvedRefReason = string(gatewayapi.ListenerReasonInvalidCertificateRef)
					break
				}
				secret := &corev1.Secret{}
				secretNamespace := gateway.Namespace
				if certRef.Namespace != nil {
					secretNamespace = string(*certRef.Namespace)
				}
				if err := client.Get(ctx, k8stypes.NamespacedName{Namespace: secretNamespace, Name: string(certRef.Name)}, secret); err != nil {
					if !apierrors.IsNotFound(err) {
						return nil, err
					}
					tlsResolvedRefReason = string(gatewayapi.ListenerReasonInvalidCertificateRef)
					break
				}
				if !isTLSSecretValid(secret) {
					tlsResolvedRefReason = string(gatewayapi.ListenerReasonInvalidCertificateRef)
				}
			}
			if gatewayapi.ListenerConditionReason(tlsResolvedRefReason) != gatewayapi.ListenerReasonResolvedRefs {
				ResolvedRefsReason = gatewayapi.ListenerConditionReason(tlsResolvedRefReason)
			}
		}

		attachedRoutes, err := getAttachedRoutesForListener(ctx, client, *gateway, listenerIndex)
		if err != nil {
			return nil, err
		}

		status := gatewayapi.ListenerStatus{
			Name:           listener.Name,
			Conditions:     []metav1.Condition{},
			SupportedKinds: supportedkinds,
			AttachedRoutes: attachedRoutes,
		}

		// if the resolvedRefs condition is not successful, append the resolvedRefs condition failed with the proper reason
		if ResolvedRefsReason != gatewayapi.ListenerReasonResolvedRefs {
			status.Conditions = append(status.Conditions, metav1.Condition{
				Type:               string(gatewayapi.ListenerConditionResolvedRefs),
				Reason:             string(ResolvedRefsReason),
				Status:             metav1.ConditionFalse,
				LastTransitionTime: metav1.Now(),
				ObservedGeneration: gateway.Generation,
			})
		}

		if _, ok := portToProtocol[listener.Port]; !ok {
			// unoccupied ports are free game
			portToProtocol[listener.Port] = listener.Protocol
			portToHostname[listener.Port][hostname] = true
		} else {
			if !canSharePort(listener.Protocol, portToProtocol[listener.Port]) {
				status.Conditions = append(status.Conditions, metav1.Condition{
					Type:               string(gatewayapi.ListenerConditionConflicted),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayapi.ListenerReasonProtocolConflict),
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
					conflictedHostnames[listener.Port] = map[gatewayapi.Hostname]bool{}
				}
				if _, exists := portToHostname[listener.Port][hostname]; !exists {
					portToHostname[listener.Port][hostname] = true
				} else {
					status.Conditions = append(status.Conditions, metav1.Condition{
						Type:               string(gatewayapi.ListenerConditionConflicted),
						Status:             metav1.ConditionTrue,
						ObservedGeneration: gateway.Generation,
						LastTransitionTime: metav1.Now(),
						Reason:             string(gatewayapi.ListenerReasonHostnameConflict),
					})
					conflictedHostnames[listener.Port][hostname] = true
				}
			}
		}

		// independent of conflict detection. for example, two TCP Listeners both requesting the same port that Kong
		// does not provide should be both Conflicted and Detached
		if len(kongProtocolsToPort[listener.Protocol]) == 0 {
			status.Conditions = append(status.Conditions, metav1.Condition{
				Type:               string(gatewayapi.ListenerConditionAccepted),
				Status:             metav1.ConditionFalse,
				ObservedGeneration: gateway.Generation,
				LastTransitionTime: metav1.Now(),
				Reason:             string(gatewayapi.ListenerReasonUnsupportedProtocol),
				Message:            "no Kong listen with the requested protocol is configured",
			})
		} else {
			if _, ok := kongProtocolsToPort[listener.Protocol][listener.Port]; !ok {
				status.Conditions = append(status.Conditions, metav1.Condition{
					Type:               string(gatewayapi.ListenerConditionAccepted),
					Status:             metav1.ConditionFalse,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayapi.ListenerReasonPortUnavailable),
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
					Type:               string(gatewayapi.ListenerConditionConflicted),
					Status:             metav1.ConditionFalse,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayapi.ListenerReasonNoConflicts),
				},
				metav1.Condition{
					Type:               string(gatewayapi.ListenerConditionResolvedRefs),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayapi.ListenerReasonResolvedRefs),
				},
				metav1.Condition{
					Type:               string(gatewayapi.ListenerConditionProgrammed),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayapi.ListenerReasonProgrammed),
				},
			)
		} else {
			// Any conditions we added above will prevent the Listener from becoming programmed.
			status.Conditions = append(status.Conditions,
				metav1.Condition{
					Type:               string(gatewayapi.ListenerConditionProgrammed),
					Status:             metav1.ConditionFalse,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayapi.ListenerReasonInvalid),
				},
			)
		}

		// If the we did not set the "Accepted" status of the listener previously,
		// which means that there are no conflicts, we set `Accepted` condition to true.
		// (while the listener can have other problems to prevent it being programmed)
		if !lo.ContainsBy(status.Conditions,
			func(condition metav1.Condition) bool {
				return condition.Type == string(gatewayapi.ListenerConditionAccepted)
			},
		) {
			status.Conditions = append(status.Conditions,
				metav1.Condition{
					Type:               string(gatewayapi.ListenerConditionAccepted),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayapi.ListenerReasonAccepted),
				},
			)
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

		var hostname gatewayapi.Hostname
		if listener.Hostname != nil {
			hostname = *listener.Hostname
		}
		// there's no filter for protocols that don't use Hostname, but this won't be populated from earlier for those
		if _, ok := conflictedHostnames[listener.Port][hostname]; ok {
			conflictReason = string(gatewayapi.ListenerReasonHostnameConflict)
		}

		if _, ok := conflictedPorts[listener.Port]; ok {
			conflictReason = string(gatewayapi.ListenerReasonProtocolConflict)
		}

		newConditions := []metav1.Condition{}

		if len(conflictReason) > 0 {
			for _, cond := range statuses[listener.Name].Conditions {
				// shut up linter, there's a default
				switch gatewayapi.ListenerConditionType(cond.Type) {
				case gatewayapi.ListenerConditionProgrammed, gatewayapi.ListenerConditionConflicted:
					continue
				default:
					newConditions = append(newConditions, cond)
				}
			}
			newConditions = append(newConditions,
				metav1.Condition{
					Type:               string(gatewayapi.ListenerConditionConflicted),
					Status:             metav1.ConditionTrue,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             conflictReason,
				},
				metav1.Condition{
					Type:               string(gatewayapi.ListenerConditionProgrammed),
					Status:             metav1.ConditionFalse,
					ObservedGeneration: gateway.Generation,
					LastTransitionTime: metav1.Now(),
					Reason:             string(gatewayapi.ListenerReasonInvalid),
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
	statusArray := []gatewayapi.ListenerStatus{}
	for _, status := range statuses {
		statusArray = append(statusArray, status)
	}

	return statusArray, nil
}

// getReferenceGrantConditionReason gets a certRef belonging to a specific listener and a slice of referenceGrants.
func getReferenceGrantConditionReason(
	gatewayNamespace string,
	certRef gatewayapi.SecretObjectReference,
	referenceGrants []gatewayapi.ReferenceGrant,
) string {
	// no need to have this reference granted
	if certRef.Namespace == nil || *certRef.Namespace == (gatewayapi.Namespace)(gatewayNamespace) {
		return string(gatewayapi.ListenerReasonResolvedRefs)
	}

	certRefNamespace := string(*certRef.Namespace)
	for _, grant := range referenceGrants {
		// the grant must exist in the same namespace of the referenced resource
		if grant.Namespace != certRefNamespace {
			continue
		}
		for _, from := range grant.Spec.From {
			// we are interested only in grants for gateways that want to reference secrets
			if from.Group != gatewayapi.V1Group || from.Kind != "Gateway" {
				continue
			}
			if from.Namespace == gatewayapi.Namespace(gatewayNamespace) {
				for _, to := range grant.Spec.To {
					if (to.Group != "" && to.Group != "core") || to.Kind != "Secret" {
						continue
					}
					// if all the above conditions are satisfied, and the name of the referenced secret matches
					// the granted resource name, then return a reason "ResolvedRefs"
					if to.Name == nil || string(*to.Name) == string(certRef.Name) {
						return string(gatewayapi.ListenerReasonResolvedRefs)
					}
				}
			}
		}
	}
	// if no grants have been found for the reference, return an "InvalidCertificateRef" reason
	return string(gatewayapi.ListenerReasonRefNotPermitted)
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
		log.Error(fmt.Errorf("invalid type"), "Received invalid event type in event handlers", "found", reflect.TypeOf(watchEvent))
		return false
	}

	for _, obj := range objs {
		gwc, ok := obj.(*gatewayapi.GatewayClass)
		if !ok {
			log.Error(fmt.Errorf("invalid type"), "Received invalid object type in event handlers", "expected", "GatewayClass", "found", reflect.TypeOf(obj))
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
func getListenerSupportedRouteKinds(l gatewayapi.Listener) ([]gatewayapi.RouteGroupKind, gatewayapi.ListenerConditionReason) {
	if l.AllowedRoutes == nil || len(l.AllowedRoutes.Kinds) == 0 {
		switch l.Protocol {
		case gatewayapi.HTTPProtocolType, gatewayapi.HTTPSProtocolType:
			return []gatewayapi.RouteGroupKind{
				builder.NewRouteGroupKind().HTTPRoute().Build(),
				builder.NewRouteGroupKind().GRPCRoute().Build(),
			}, gatewayapi.ListenerReasonResolvedRefs
		case gatewayapi.TCPProtocolType:
			return builder.NewRouteGroupKind().TCPRoute().IntoSlice(), gatewayapi.ListenerReasonResolvedRefs
		case gatewayapi.UDPProtocolType:
			return builder.NewRouteGroupKind().UDPRoute().IntoSlice(), gatewayapi.ListenerReasonResolvedRefs
		case gatewayapi.TLSProtocolType:
			return builder.NewRouteGroupKind().TLSRoute().IntoSlice(), gatewayapi.ListenerReasonResolvedRefs
		}
	}

	var (
		supportedRGK = []gatewayapi.RouteGroupKind{}
		reason       = gatewayapi.ListenerReasonResolvedRefs
	)
	for _, gk := range l.AllowedRoutes.Kinds {
		if gk.Group != nil && *gk.Group == gatewayv1.GroupName {
			_, ok := lo.Find(supportedKinds, func(k gatewayapi.Kind) bool {
				return gk.Kind == k
			})
			if ok {
				supportedRGK = append(supportedRGK, gk)
				continue
			}
			reason = gatewayapi.ListenerReasonInvalidRouteKinds
		} else {
			reason = gatewayapi.ListenerReasonInvalidRouteKinds
		}
	}

	return supportedRGK, reason
}

func isTLSSecretValid(secret *corev1.Secret) bool {
	var ok bool
	var crt, key []byte
	if crt, ok = secret.Data["tls.crt"]; !ok {
		return false
	}
	if key, ok = secret.Data["tls.key"]; !ok {
		return false
	}
	if p, _ := pem.Decode(crt); p == nil {
		return false
	}
	if p, _ := pem.Decode(key); p == nil {
		return false
	}
	return true
}

// routeAcceptedByGateways finds all the Gateways the route has been accepted by
// and returns them in the form of a NamespacedName slice.
func routeAcceptedByGateways(route *gatewayapi.HTTPRoute,
) []k8stypes.NamespacedName {
	gateways := []k8stypes.NamespacedName{}
	for _, routeParentStatus := range getRouteStatusParents(route) {
		gatewayNamespace := route.GetNamespace()
		parentRef := routeParentStatus.ParentRef
		if (parentRef.Group != nil && *parentRef.Group != gatewayapi.V1Group) ||
			(parentRef.Kind != nil && *parentRef.Kind != "Gateway") {
			continue
		}
		if parentRef.Namespace != nil {
			gatewayNamespace = string(*parentRef.Namespace)
		}

		gateways = append(gateways,
			k8stypes.NamespacedName{
				Namespace: gatewayNamespace,
				Name:      string(parentRef.Name),
			},
		)
	}
	return gateways
}

// getAttachedRoutesForListener returns the number of all the routes that are attached
// to the provided Gateway.
//
// NOTE: At this point we take into account HTTPRoutes only, as they are the
// only routes in GA.
func getAttachedRoutesForListener(ctx context.Context, mgrc client.Client, gateway gatewayapi.Gateway, listenerIndex int) (int32, error) {
	httpRouteList := gatewayapi.HTTPRouteList{}
	if err := mgrc.List(ctx, &httpRouteList); err != nil {
		return 0, err
	}

	var attachedRoutes int32
	for _, route := range httpRouteList.Items {
		acceptedByGateway := lo.ContainsBy(route.Status.Parents, func(parentStatus gatewayapi.RouteParentStatus) bool {
			parentRef := parentStatus.ParentRef
			if parentRef.Group != nil && *parentRef.Group != gatewayapi.V1Group {
				return false
			}
			if parentRef.Kind != nil && *parentRef.Kind != "Gateway" {
				return false
			}
			gatewayNamespace := route.Namespace
			if parentRef.Namespace != nil {
				gatewayNamespace = string(*parentRef.Namespace)
			}
			return gateway.Namespace == gatewayNamespace && gateway.Name == string(parentRef.Name)
		})
		if !acceptedByGateway {
			continue
		}

		for _, parentRef := range route.Spec.ParentRefs {
			accepted, err := isRouteAcceptedByListener(
				ctx,
				mgrc,
				&route,
				gateway,
				listenerIndex,
				parentRef,
			)
			if err != nil {
				return 0, err
			}
			if accepted {
				attachedRoutes++
			}
		}
	}
	return attachedRoutes, nil
}
