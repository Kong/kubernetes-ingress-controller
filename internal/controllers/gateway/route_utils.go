package gateway

import (
	"context"
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

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
	gateway      *gatewayv1alpha2.Gateway
	condition    metav1.Condition
	listenerName string
}

// parentRefsForRoute provides a list of the parentRefs given a Gateway APIs route object
// (e.g. HTTPRoute, TCPRoute, e.t.c.) which refer to the Gateway resource(s) which manage it.
func parentRefsForRoute(obj client.Object) ([]gatewayv1alpha2.ParentReference, error) {
	switch v := obj.(type) {
	case *gatewayv1alpha2.HTTPRoute:
		return v.Spec.ParentRefs, nil
	case *gatewayv1alpha2.UDPRoute:
		return v.Spec.ParentRefs, nil
	case *gatewayv1alpha2.TCPRoute:
		return v.Spec.ParentRefs, nil
	case *gatewayv1alpha2.TLSRoute:
		return v.Spec.ParentRefs, nil
	default:
		return nil, fmt.Errorf("cant determine parent gateway for unsupported type %s", reflect.TypeOf(obj))
	}
}

// getSupportedGatewayForRoute will retrieve the Gateway and GatewayClass object for any
// Gateway APIs route object (e.g. HTTPRoute, TCPRoute, e.t.c.) from the provided cached
// client if they match this controller. If there are no gateways present for this route
// OR the present gateways are references to missing objects, this will return a unsupportedGW error.
func getSupportedGatewayForRoute(ctx context.Context, mgrc client.Client, obj client.Object) ([]supportedGatewayWithCondition, error) {
	// gather the parentrefs for this route object
	parentRefs, err := parentRefsForRoute(obj)
	if err != nil {
		return nil, err
	}

	// search each parentRef to see if this controller is one of the supported ones
	gateways := make([]supportedGatewayWithCondition, 0)
	for _, parentRef := range parentRefs {
		// gather the namespace/name for the gateway
		namespace := obj.GetNamespace()
		if parentRef.Namespace != nil {
			// TODO: need namespace restrictions implementation done before
			// merging this, need to filter out objects with a disallowed NS.
			// https://github.com/Kong/kubernetes-ingress-controller/issues/2080
			namespace = string(*parentRef.Namespace)
		}
		name := string(parentRef.Name)

		// pull the Gateway object from the cached client
		gateway := gatewayv1alpha2.Gateway{}
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
		gatewayClass := gatewayv1alpha2.GatewayClass{}
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
			// set true if we find any AllowedRoutes. there may be none, in which case any namespace is permitted
			filtered := false
			matchingHostname := metav1.ConditionFalse
			for _, listener := range gateway.Spec.Listeners {
				// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2408
				// This currently only performs a baseline filter to ensure that routes cannot match based on namespace
				// criteria on a listener that cannot possibly handle them (e.g. an HTTPRoute should not be included
				// based on matching a filter for a UDP listener). This needs to be expanded to an allowedRoutes.kind
				// implementation with default allowed kinds when there's no user-specified filter.
				var oneHostnameMatch bool
				switch obj := obj.(type) {
				case *gatewayv1alpha2.HTTPRoute:
					hostnames := obj.Spec.Hostnames
					oneHostnameMatch = listenerHostnameIntersectWithRouteHostnames(listener, hostnames)
					if !(listener.Protocol == gatewayv1alpha2.HTTPProtocolType || listener.Protocol == gatewayv1alpha2.HTTPSProtocolType) {
						continue
					}
				case *gatewayv1alpha2.TCPRoute:
					if listener.Protocol != gatewayv1alpha2.TCPProtocolType {
						continue
					}
				case *gatewayv1alpha2.UDPRoute:
					if listener.Protocol != gatewayv1alpha2.UDPProtocolType {
						continue
					}
				case *gatewayv1alpha2.TLSRoute:
					hostnames := obj.Spec.Hostnames
					oneHostnameMatch = listenerHostnameIntersectWithRouteHostnames(listener, hostnames)
					if listener.Protocol != gatewayv1alpha2.TLSProtocolType {
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
					if *listener.AllowedRoutes.Namespaces.From == gatewayv1alpha2.NamespacesFromAll {
						// we allow "all" by just stuffing the namespace we want to find into the map
						allowedNamespaces[obj.GetNamespace()] = nil
					} else if *listener.AllowedRoutes.Namespaces.From == gatewayv1alpha2.NamespacesFromSame {
						allowedNamespaces[gateway.ObjectMeta.Namespace] = nil
					} else if *listener.AllowedRoutes.Namespaces.From == gatewayv1alpha2.NamespacesFromSelector {
						namespaces := &corev1.NamespaceList{}
						selector, err := metav1.LabelSelectorAsSelector(listener.AllowedRoutes.Namespaces.Selector)
						if err != nil {
							return nil, fmt.Errorf("failed to convert LabelSelector to Selector for gateway %s",
								gateway.ObjectMeta.Name)
						}
						err = mgrc.List(ctx, namespaces,
							&client.ListOptions{LabelSelector: selector})
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

			_, allowedNamespace := allowedNamespaces[obj.GetNamespace()]
			if !filtered || allowedNamespace {
				// if there is no matchingHostname, the gateway Status Condition Accepted must be set to False
				// with reason NoMatchingListenerHostname
				reason := gatewayv1alpha2.RouteReasonAccepted
				if matchingHostname == metav1.ConditionFalse {
					reason = gatewayv1alpha2.RouteReasonNoMatchingListenerHostname
				}
				var listenerName string
				if parentRef.SectionName != nil && *parentRef.SectionName != "" {
					listenerName = string(*parentRef.SectionName)
				}
				gateways = append(gateways, supportedGatewayWithCondition{
					gateway:      &gateway,
					listenerName: listenerName,
					condition: metav1.Condition{
						Type:   string(gatewayv1alpha2.RouteConditionAccepted),
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

func listenerHostnameIntersectWithRouteHostnames(listener gatewayv1alpha2.Listener, hostnames []gatewayv1alpha2.Hostname) bool {
	// if the listener has no hostname, all hostnames automatically intersect
	if listener.Hostname == nil || *listener.Hostname == "" || len(hostnames) == 0 {
		return true
	}

	// iterate over all the hostnames and check that at least one intersect with the listener hostname
	for _, hostname := range hostnames {
		if util.HostnamesIntersect(string(*listener.Hostname), string(hostname)) {
			return true
		}
	}
	return false
}

// filterHostnames accepts a HTTPRoute and returns a version of the same object with only a subset of the
// hostnames, the ones matching with the listeners' hostname.
func filterHostnames(gateways []supportedGatewayWithCondition, httpRoute *gatewayv1alpha2.HTTPRoute) *gatewayv1alpha2.HTTPRoute {
	filteredHostnames := make([]gatewayv1alpha2.Hostname, 0)

	// if no hostnames are specified in the route spec, get all the hostnames from
	// the gateway
	if len(httpRoute.Spec.Hostnames) == 0 {
		for _, gateway := range gateways {
			for _, listener := range gateway.gateway.Spec.Listeners {
				if listenerName := gatewayv1alpha2.SectionName(gateway.listenerName); listenerName == "" || listenerName == listener.Name {
					if listener.Hostname != nil {
						filteredHostnames = append(filteredHostnames, *listener.Hostname)
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
func getMinimumHostnameIntersection(gateways []supportedGatewayWithCondition, hostname gatewayv1alpha2.Hostname) gatewayv1alpha2.Hostname {
	for _, gateway := range gateways {
		for _, listener := range gateway.gateway.Spec.Listeners {
			// if the listenerName is specified and matches the name of the gateway listener proceed
			if gatewayv1alpha2.SectionName(gateway.listenerName) == "" ||
				gatewayv1alpha2.SectionName(gateway.listenerName) == listener.Name {
				if listener.Hostname == nil || *listener.Hostname == "" {
					return hostname
				}
				if util.HostnamesMatch(string(*listener.Hostname), string(hostname)) {
					return hostname
				}
				if util.HostnamesMatch(string(hostname), string(*listener.Hostname)) {
					return *listener.Hostname
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
func isHTTPReferenceGranted(grantSpec gatewayv1alpha2.ReferenceGrantSpec, backendRef gatewayv1alpha2.HTTPBackendRef, fromNamespace string) bool {
	var backendRefGroup gatewayv1alpha2.Group
	var backendRefKind gatewayv1alpha2.Kind

	if backendRef.Group != nil {
		backendRefGroup = *backendRef.Group
	}
	if backendRef.Kind != nil {
		backendRefKind = *backendRef.Kind
	}
	for _, from := range grantSpec.From {
		if from.Group != gatewayv1alpha2.GroupName || from.Kind != "HTTPRoute" || fromNamespace != string(from.Namespace) {
			continue
		}

		for _, to := range grantSpec.To {
			if backendRefGroup == to.Group &&
				backendRefKind == to.Kind &&
				(to.Name == nil || *to.Name == backendRef.Name) {
				return true
			}
		}
	}

	return false
}
