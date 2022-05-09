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
)

// -----------------------------------------------------------------------------
// Route Utilities
// -----------------------------------------------------------------------------

const unsupportedGW = "no supported Gateway found for route"

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
func getSupportedGatewayForRoute(ctx context.Context, mgrc client.Client, obj client.Object) ([]*gatewayv1alpha2.Gateway, error) {
	// gather the parentrefs for this route object
	parentRefs, err := parentRefsForRoute(obj)
	if err != nil {
		return nil, err
	}

	// search each parentRef to see if this controller is one of the supported ones
	gateways := make([]*gatewayv1alpha2.Gateway, 0)
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
			for _, listener := range gateway.Spec.Listeners {
				// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2408
				// This currently only performs a baseline filter to ensure that routes cannot match based on namespace
				// criteria on a listener that cannot possibly handle them (e.g. an HTTPRoute should not be included
				// based on matching a filter for a UDP listener). This needs to be expanded to an allowedRoutes.kind
				// implementation with default allowed kinds when there's no user-specified filter.
				switch obj.(type) {
				case *gatewayv1alpha2.HTTPRoute:
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
					if listener.Protocol != gatewayv1alpha2.TLSProtocolType {
						continue
					}
				default:
					continue
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
				gateways = append(gateways, &gateway)
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
