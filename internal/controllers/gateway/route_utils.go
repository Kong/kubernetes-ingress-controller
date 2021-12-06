package gateway

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// -----------------------------------------------------------------------------
// Route Utilities
// -----------------------------------------------------------------------------

const unsupportedGW = "no supported Gateway found for route"

// parentRefsForRoute provides a list of the parentRefs given a Gateway APIs route object
// (e.g. HTTPRoute, TCPRoute, e.t.c.) which refer to the Gateway resource(s) which manage it.
func parentRefsForRoute(obj client.Object) ([]gatewayv1alpha2.ParentRef, error) {
	switch v := obj.(type) {
	case *gatewayv1alpha2.HTTPRoute:
		return v.Spec.ParentRefs, nil
	default:
		return nil, fmt.Errorf("cant determine parent gateway for unsupported type %s", reflect.TypeOf(obj))
	}
}

// getSupportedGatewayForRoute will retrieve the Gateway and GatewayClass object for any
// Gateway APIs route object (e.g. HTTPRoute, TCPRoute, e.t.c.) from the provided cached
// client if they match this controller. If there are no gateways present for this route
// OR the present gateways are references to missing objects, this will return a unsupportedGW error.
func getSupportedGatewayForRoute(ctx context.Context, mgrc client.Client, obj client.Object) (*gatewayv1alpha2.Gateway, error) {
	// gather the parentrefs for this route object
	parentRefs, err := parentRefsForRoute(obj)
	if err != nil {
		return nil, err
	}

	// search each parentRef to see if this controller is one of the supported ones
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
			return &gateway, nil
		}
	}

	// this is the "if all else false" fallback. If we reach this point either:
	//
	//  a) there are no gateways configured
	//  b) the gateways configured are not present in the API
	//  c) combination of a & b
	//
	// we provide a specific error for this condition rather than making the
	// caller check for nil values in the gateway and class.
	return nil, fmt.Errorf(unsupportedGW)
}
