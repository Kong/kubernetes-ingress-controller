package gateway

import (
	"context"
	"fmt"
	"reflect"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/types"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Route Utilities
// -----------------------------------------------------------------------------

const (
	unsupportedGW = "no supported Gateway found for route"
)

// supportedGatewayWithCondition is a struct that wraps a gateway and some further info
// such as the condition Status condition Accepted of the gateway and the listenerName.
type supportedGatewayWithCondition struct {
	gateway      *Gateway
	condition    metav1.Condition
	listenerName string
}

// parentRefsForRoute provides a list of the parentRefs given a Gateway APIs route object
// (e.g. HTTPRoute, TCPRoute, e.t.c.) which refer to the Gateway resource(s) which manage it.
func parentRefsForRoute[T types.RouteT](route T) ([]ParentReference, error) {
	// Note: Ideally we wouldn't have to do this but it's hard to juggle around types
	// and support ParentReference and gatewayv1alpha2.ParentReference
	// at the same time so we just copy v1alpha2 refs to a new v1beta1 slice.
	convertV1Alpha2ToV1Beta1ParentReference := func(
		refsAlpha []gatewayv1alpha2.ParentReference,
	) []ParentReference {
		ret := make([]ParentReference, len(refsAlpha))
		for i, v := range refsAlpha {
			ret[i] = ParentReference{
				Group:       (*gatewayv1beta1.Group)(v.Group),
				Kind:        (*Kind)(v.Kind),
				Namespace:   (*gatewayv1beta1.Namespace)(v.Namespace),
				Name:        (gatewayv1beta1.ObjectName)(v.Name),
				SectionName: (*SectionName)(v.SectionName),
				Port:        (*PortNumber)(v.Port),
			}
		}
		return ret
	}

	switch r := (interface{})(route).(type) {
	case *gatewayv1beta1.HTTPRoute:
		return r.Spec.ParentRefs, nil
	case *gatewayv1alpha2.UDPRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	case *gatewayv1alpha2.TCPRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	case *gatewayv1alpha2.TLSRoute:
		return convertV1Alpha2ToV1Beta1ParentReference(r.Spec.ParentRefs), nil
	default:
		return nil, fmt.Errorf("cant determine parent gateway for unsupported type %s", reflect.TypeOf(route))
	}
}

const (
	// This reason is used with the "Accepted" condition when the Gateway has no
	// compatible Listeners whose Port matches the route
	// NOTE: This should probably be proposed upstream.
	RouteReasonNoMatchingListenerPort gatewayv1beta1.RouteConditionReason = "NoMatchingListenerPort"
)

// getSupportedGatewayForRoute will retrieve the Gateway and GatewayClass object for any
// Gateway APIs route object (e.g. HTTPRoute, TCPRoute, e.t.c.) from the provided cached
// client if they match this controller. If there are no gateways present for this route
// OR the present gateways are references to missing objects, this will return a unsupportedGW error.
func getSupportedGatewayForRoute[T types.RouteT](ctx context.Context, mgrc client.Client, route T) ([]supportedGatewayWithCondition, error) {
	// gather the parentrefs for this route object
	parentRefs, err := parentRefsForRoute(route)
	if err != nil {
		return nil, err
	}

	// search each parentRef to see if this controller is one of the supported ones
	gateways := make([]supportedGatewayWithCondition, 0)
	for _, parentRef := range parentRefs {
		// gather the namespace/name for the gateway
		namespace := route.GetNamespace()
		if parentRef.Namespace != nil {
			// TODO: need namespace restrictions implementation done before
			// merging this, need to filter out objects with a disallowed NS.
			// https://github.com/Kong/kubernetes-ingress-controller/issues/2080
			namespace = string(*parentRef.Namespace)
		}
		name := string(parentRef.Name)

		// pull the Gateway object from the cached client
		gateway := gatewayv1beta1.Gateway{}
		if err := mgrc.Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}, &gateway); err != nil {
			if errors.IsNotFound(err) {
				// if a configured gateway is not found it's still possible
				// that there's another gateway, so keep searching through the list.
				continue
			}
			return nil, fmt.Errorf("failed to retrieve gateway for route: %w", err)
		}

		// pull the GatewayClass for the Gateway object from the cached client
		gatewayClass := gatewayv1beta1.GatewayClass{}
		if err := mgrc.Get(ctx, client.ObjectKey{
			Name: string(gateway.Spec.GatewayClassName),
		}, &gatewayClass); err != nil {
			if errors.IsNotFound(err) {
				// if a configured gatewayClass is not found it's still possible
				// that there's another properly configured gateway in the parentRefs,
				// so keep searching through the list.
				continue
			}
			return nil, fmt.Errorf("failed to retrieve gatewayclass for gateway: %w", err)
		}

		// if the GatewayClass matches this controller we're all set and this controller
		// should reconcile this object.
		if gatewayClass.Spec.ControllerName == ControllerName {
			allowedNamespaces := make(map[string]interface{})
			var (
				// set true if we find any AllowedRoutes. there may be none, in which case any namespace is permitted
				filtered         = false
				matchingHostname = metav1.ConditionFalse
				// set to true if ParentRef specifies a Port and a listener matches that Port.
				portMatched = false
			)
			for _, listener := range gateway.Spec.Listeners {
				// TODO check listenerStatus.SupportedKinds

				// Check if we already have a matching listener in status.
				if !existsMatchingReadyListenerInStatus(listener, gateway.Status.Listeners) {
					continue
				}

				// Perform the port matching as described in GEP-957.
				if parentRef.Port != nil && *parentRef.Port != listener.Port {
					// This ParentRef has a port specified and it's different than current listener's port.
					continue
				} else if parentRef.Port != nil && *parentRef.Port == listener.Port {
					portMatched = true
				}

				// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2408
				// This currently only performs a baseline filter to ensure that routes cannot match based on namespace
				// criteria on a listener that cannot possibly handle them (e.g. an HTTPRoute should not be included
				// based on matching a filter for a UDP listener). This needs to be expanded to an allowedRoutes.kind
				// implementation with default allowed kinds when there's no user-specified filter.
				var oneHostnameMatch bool
				switch r := (interface{})(route).(type) {
				case *gatewayv1beta1.HTTPRoute:
					hostnames := r.Spec.Hostnames
					oneHostnameMatch = listenerHostnameIntersectWithRouteHostnames(listener, hostnames)
					if !(listener.Protocol == HTTPProtocolType || listener.Protocol == HTTPSProtocolType) {
						continue
					}

				case *gatewayv1alpha2.TCPRoute:
					if listener.Protocol != (TCPProtocolType) {
						continue
					}
				case *gatewayv1alpha2.UDPRoute:
					if listener.Protocol != (UDPProtocolType) {
						continue
					}
				case *gatewayv1alpha2.TLSRoute:
					hostnames := r.Spec.Hostnames
					oneHostnameMatch = listenerHostnameIntersectWithRouteHostnames(listener, hostnames)
					if listener.Protocol != (TLSProtocolType) {
						continue
					}
				default:
					continue
				}
				if oneHostnameMatch {
					matchingHostname = metav1.ConditionTrue
				}
				if listener.AllowedRoutes != nil {
					filtered = true
					if *listener.AllowedRoutes.Namespaces.From == gatewayv1beta1.NamespacesFromAll {
						// we allow "all" by just stuffing the namespace we want to find into the map
						allowedNamespaces[route.GetNamespace()] = nil
					} else if *listener.AllowedRoutes.Namespaces.From == gatewayv1beta1.NamespacesFromSame {
						allowedNamespaces[gateway.ObjectMeta.Namespace] = nil
					} else if *listener.AllowedRoutes.Namespaces.From == gatewayv1beta1.NamespacesFromSelector {
						namespaces := &corev1.NamespaceList{}
						selector, err := metav1.LabelSelectorAsSelector(listener.AllowedRoutes.Namespaces.Selector)
						if err != nil {
							return nil, fmt.Errorf("failed to convert LabelSelector to Selector for gateway %s",
								gateway.ObjectMeta.Name)
						}
						err = mgrc.List(ctx, namespaces, &client.ListOptions{LabelSelector: selector})
						if err != nil {
							return nil, fmt.Errorf("could not fetch allowed namespaces for gateway %s",
								gateway.ObjectMeta.Name)
						}
						for _, allowed := range namespaces.Items {
							allowedNamespaces[allowed.ObjectMeta.Name] = nil
						}
					}
				}
			}

			_, allowedNamespace := allowedNamespaces[route.GetNamespace()]
			if ((parentRef.Port != nil) && !portMatched) ||
				(!filtered || allowedNamespace) {

				reason := gatewayv1beta1.RouteReasonAccepted
				if (parentRef.Port != nil) && !portMatched {
					// If ParentRef specified a Port but none of the listeners matched, the gateway Status
					// Condition Accepted must be set to False with reason NoMatchingListenerPort
					reason = RouteReasonNoMatchingListenerPort
				} else if matchingHostname == metav1.ConditionFalse {
					// If there is no matchingHostname, the gateway Status Condition Accepted must be set to False
					// with reason NoMatchingListenerHostname
					reason = gatewayv1beta1.RouteReasonNoMatchingListenerHostname
				}

				var listenerName string
				if parentRef.SectionName != nil && *parentRef.SectionName != "" {
					listenerName = string(*parentRef.SectionName)
				}

				gateways = append(gateways, supportedGatewayWithCondition{
					gateway:      &gateway,
					listenerName: listenerName,
					condition: metav1.Condition{
						Type:   string(gatewayv1beta1.RouteConditionAccepted),
						Status: matchingHostname,
						Reason: string(reason),
					},
				})
			}
		}
	}

	if len(gateways) == 0 {
		// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2417 separate out various rejected reasons
		// and apply specific statuses for those failures in the Route controllers
		return nil, fmt.Errorf(unsupportedGW)
	}

	return gateways, nil
}

func existsMatchingReadyListenerInStatus(listener Listener, lss []ListenerStatus) bool {
	// Find listener's status...
	listenerStatus, ok := lo.Find(lss, func(ls gatewayv1beta1.ListenerStatus) bool {
		return ls.Name == listener.Name
	})
	if !ok {
		return false // Listener's status not found
	}
	// ... and verify if it's ready.
	lReadyCond, ok := lo.Find(listenerStatus.Conditions, func(c metav1.Condition) bool {
		return c.Type == string(gatewayv1beta1.ListenerConditionReady)
	})
	if !ok {
		return false
	}
	if lReadyCond.Status != "True" {
		return false // Listener is not ready yet.
	}
	return true
}

func listenerHostnameIntersectWithRouteHostnames[H types.HostnameT, L types.ListenerT](listener L, hostnames []H) bool {
	if len(hostnames) == 0 {
		return true
	}

	// if the listener has no hostname, all hostnames automatically intersect
	switch l := (interface{})(listener).(type) {
	case gatewayv1alpha2.Listener:
		if l.Hostname == nil || *l.Hostname == "" {
			return true
		}

		// iterate over all the hostnames and check that at least one intersect with the listener hostname
		for _, hostname := range hostnames {
			if util.HostnamesIntersect(*l.Hostname, hostname) {
				return true
			}
		}
	case Listener:
		if l.Hostname == nil || *l.Hostname == "" {
			return true
		}

		// iterate over all the hostnames and check that at least one intersect with the listener hostname
		for _, hostname := range hostnames {
			if util.HostnamesIntersect(*l.Hostname, hostname) {
				return true
			}
		}
	}

	return false
}

// filterHostnames accepts a HTTPRoute and returns a version of the same object with only a subset of the
// hostnames, the ones matching with the listeners' hostname.
func filterHostnames(gateways []supportedGatewayWithCondition, httpRoute *gatewayv1beta1.HTTPRoute) *gatewayv1beta1.HTTPRoute {
	filteredHostnames := make([]gatewayv1beta1.Hostname, 0)

	// if no hostnames are specified in the route spec, get all the hostnames from
	// the gateway
	if len(httpRoute.Spec.Hostnames) == 0 {
		for _, gateway := range gateways {
			for _, listener := range gateway.gateway.Spec.Listeners {
				if listenerName := gateway.listenerName; listenerName == "" || listenerName == string(listener.Name) {
					if listener.Hostname != nil {
						filteredHostnames = append(filteredHostnames, (*listener.Hostname))
					}
				}
			}
		}
	} else {
		for _, hostname := range httpRoute.Spec.Hostnames {
			if hostnameMatching := getMinimumHostnameIntersection(gateways, hostname); hostnameMatching != "" {
				filteredHostnames = append(filteredHostnames, hostnameMatching)
			}
		}
	}

	httpRoute.Spec.Hostnames = filteredHostnames
	return httpRoute
}

// getMinimumHostnameIntersection returns the minimum intersecting hostname, in the sense that:
//
// - if the listener hostname is empty, return the httpRoute hostname
// - if the listener hostname acts as a wildcard for the httpRoute hostname, return the httpRoute hostname
// - if the httpRoute hostname acts as a wildcard for the listener hostname, return the listener hostname
// - if the httpRoute hostname is the same of the listener hostname, return it
// - if none of the above is true, return an empty string.
func getMinimumHostnameIntersection(gateways []supportedGatewayWithCondition, hostname gatewayv1beta1.Hostname) gatewayv1beta1.Hostname {
	for _, gateway := range gateways {
		for _, listener := range gateway.gateway.Spec.Listeners {
			// if the listenerName is specified and matches the name of the gateway listener proceed
			if (SectionName)(gateway.listenerName) == "" ||
				(SectionName)(gateway.listenerName) == (listener.Name) {
				if listener.Hostname == nil || *listener.Hostname == "" {
					return hostname
				}
				if util.HostnamesMatch(string(*listener.Hostname), string(hostname)) {
					return hostname
				}
				if util.HostnamesMatch(string(hostname), string(*listener.Hostname)) {
					return (*listener.Hostname)
				}
			}
		}
	}
	return ""
}

func isRouteAccepted(gateways []supportedGatewayWithCondition) bool {
	for _, gateway := range gateways {
		if gateway.condition.Type == string(gatewayv1alpha2.RouteConditionAccepted) && gateway.condition.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}

// isHTTPReferenceGranted checks that the backendRef referenced by the HTTPRoute is granted by a ReferenceGrant.
func isHTTPReferenceGranted(grantSpec gatewayv1alpha2.ReferenceGrantSpec, backendRef gatewayv1beta1.HTTPBackendRef, fromNamespace string) bool {
	var backendRefGroup gatewayv1beta1.Group
	var backendRefKind Kind

	if backendRef.Group != nil {
		backendRefGroup = *backendRef.Group
	}
	if backendRef.Kind != nil {
		backendRefKind = *backendRef.Kind
	}
	for _, from := range grantSpec.From {
		if from.Group != gatewayv1beta1.GroupName || from.Kind != "HTTPRoute" || fromNamespace != string(from.Namespace) {
			continue
		}

		for _, to := range grantSpec.To {
			if backendRefGroup == (gatewayv1beta1.Group)(to.Group) &&
				backendRefKind == (Kind)(to.Kind) &&
				(to.Name == nil || (gatewayv1beta1.ObjectName)(*to.Name) == backendRef.Name) {
				return true
			}
		}
	}

	return false
}
