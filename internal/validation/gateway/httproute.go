package gateway

import (
	"fmt"

	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// -----------------------------------------------------------------------------
// Validation - HTTPRoute - Public Functions
// -----------------------------------------------------------------------------

// ValidateHTTPRoute provides a suite of validation for a given HTTPRoute and
// any number of Gateway resources it's attached to that the caller wants to
// have it validated against.
func ValidateHTTPRoute(httproute *gatewayv1alpha2.HTTPRoute, attachedGateways ...*gatewayv1alpha2.Gateway) (bool, string, error) {
	// perform Gateway validations for the HTTPRoute (e.g. listener validation, namespace validation, e.t.c.)
	for _, gateway := range attachedGateways {
		// TODO: validate that the namespace is supported by the linked Gateway objects
		//       See: https://github.com/Kong/kubernetes-ingress-controller/issues/2080

		// determine the parentRef for this gateway
		parentRef, err := getParentRefForHTTPRouteGateway(httproute, gateway)
		if err != nil {
			return false, "couldn't determine parentRefs for httproute", err
		}

		// gather the relevant gateway listeners for the httproute
		listeners, err := getListenersForHTTPRouteValidation(parentRef.SectionName, gateway)
		if err != nil {
			return false, "couldn't find gateway listeners for httproute", err
		}

		// perform validation of this route against it's linked gateway listeners
		for _, listener := range listeners {
			if err := validateHTTPRouteListener(listener); err != nil {
				return false, "httproute linked gateway listeners did not pass validation", err
			}
		}
	}

	// validate that no unsupported features are in use
	if err := validateHTTPRouteFeatures(httproute); err != nil {
		return false, "httproute spec did not pass validation", err
	}

	return true, "", nil
}

// -----------------------------------------------------------------------------
// Validation - HTTPRoute - Private Functions
// -----------------------------------------------------------------------------

// validateHTTPRouteListener verifies that a given HTTPRoute is configured properly
// for a given gateway listener which it is linked to.
func validateHTTPRouteListener(listener *gatewayv1alpha2.Listener) error {
	// verify that the listener supports HTTPRoute objects
	if listener.AllowedRoutes != nil && // if there are no allowed routes, assume all are allowed
		len(listener.AllowedRoutes.Kinds) > 0 { // if there are no allowed kinds, assume all are allowed
		// search each of the allowedRoutes in the listener to verify that HTTPRoute is supported
		supported := false
		for _, allowedKind := range listener.AllowedRoutes.Kinds {
			if allowedKind.Kind == "HTTPRoute" {
				supported = true
			}
		}

		// verify that we found a supported kind
		if !supported {
			return fmt.Errorf("HTTPRoute not supported by listener %s", listener.Name)
		}
	}

	return nil
}

// validateHTTPRouteFeatures checks for features that are not supported by this
// HTTPRoute implementation and validates that the provided object is not using
// any of those unsupported features.
func validateHTTPRouteFeatures(httproute *gatewayv1alpha2.HTTPRoute) error {
	for _, rule := range httproute.Spec.Rules {
		for _, match := range rule.Matches {
			// we don't support queryparam matching rules
			// See: https://github.com/Kong/kubernetes-ingress-controller/issues/2152
			if len(match.QueryParams) != 0 {
				return fmt.Errorf("queryparam matching is not yet supported for httproute")
			}

			// we don't support regex path matching rules
			// See: https://github.com/Kong/kubernetes-ingress-controller/issues/2153
			if match.Path != nil && match.Path.Type != nil && *match.Path.Type == gatewayv1alpha2.PathMatchRegularExpression {
				return fmt.Errorf("regex path matching is not yet supported for httproute")
			}

			// we don't support regex header matching rules
			// See: https://github.com/Kong/kubernetes-ingress-controller/issues/2154
			for _, hdr := range match.Headers {
				if hdr.Type != nil && *hdr.Type == gatewayv1alpha2.HeaderMatchRegularExpression {
					return fmt.Errorf("regex header matching is not yet supported for httproute")
				}
			}
		}

		// we don't currently support multiple backendRefs
		// See: https://github.com/Kong/kubernetes-ingress-controller/issues/2166
		if len(rule.BackendRefs) > 1 {
			return fmt.Errorf("multiple backendRefs is not yet supported for httproute")
		}

		// we don't support any backendRef types except Kubernetes Services
		for _, ref := range rule.BackendRefs {
			if ref.BackendRef.Group != nil && *ref.BackendRef.Group != "core" {
				return fmt.Errorf("%s is not a supported group for httproute backendRefs, only core is supported", *ref.BackendRef.Group)
			}
			if ref.BackendRef.Kind != nil && *ref.BackendRef.Kind != "Service" {
				return fmt.Errorf("%s is not a supported kind for httproute backendRefs, only Service is supported", *ref.BackendRef.Kind)
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// Validation - HTTPRoute - Private Utility Functions
// -----------------------------------------------------------------------------

// getParentRefForHTTPRouteGateway extracts an existing parentRef from an HTTPRoute
// which links to the provided Gateway if available. If the provided Gateway is not
// actually referenced by parentRef in the provided HTTPRoute this is considered
// invalid input and will produce an error.
func getParentRefForHTTPRouteGateway(httproute *gatewayv1alpha2.HTTPRoute, gateway *gatewayv1alpha2.Gateway) (*gatewayv1alpha2.ParentRef, error) {
	// search all the parentRefs on the HTTPRoute to find one that matches the Gateway
	for _, ref := range httproute.Spec.ParentRefs {
		// determine the namespace for the gateway reference
		namespace := httproute.Namespace
		if ref.Namespace != nil {
			namespace = string(*ref.Namespace)
		}

		// match the gateway with its parentRef
		if gateway.Namespace == namespace && gateway.Name == string(ref.Name) {
			copyRef := ref
			return &copyRef, nil
		}
	}

	// if no matches could be found then the input is invalid
	return nil, fmt.Errorf("no parentRef matched gateway %s/%s", gateway.Namespace, gateway.Name)
}

// getListenersForHTTPRouteValidation determines if ALL http listeners should be used for validation
// or if only a select listener should be considered.
func getListenersForHTTPRouteValidation(sectionName *gatewayv1alpha2.SectionName, gateway *gatewayv1alpha2.Gateway) ([]*gatewayv1alpha2.Listener, error) {
	var listenersForValidation []*gatewayv1alpha2.Listener
	if sectionName != nil {
		// only one specified listener is in use, only need to validate the
		// route against that listener.
		for _, listener := range gateway.Spec.Listeners {
			if listener.Name == *sectionName {
				listenerCopy := listener
				listenersForValidation = append(listenersForValidation, &listenerCopy)
			}
		}

		// if the sectionName isn't empty, we need to verify that we actually found
		// a listener which matched it, otherwise the object is invalid.
		if len(listenersForValidation) == 0 {
			return nil, fmt.Errorf("sectionname referenced listener %s was not found on gateway %s/%s", *sectionName, gateway.Namespace, gateway.Name)
		}
	} else {
		// no specific listener was chosen, so we'll simply validate against
		// all HTTP listeners on the Gateway.
		for _, listener := range gateway.Spec.Listeners {
			if listener.Protocol == gatewayv1alpha2.HTTPProtocolType || listener.Protocol == gatewayv1alpha2.HTTPSProtocolType {
				listenerCopy := listener
				listenersForValidation = append(listenersForValidation, &listenerCopy)
			}
		}
	}

	// if for some reason the gateway has no listeners (it may be under active provisioning)
	// the HTTPRoute fails validation because it has no listeners that can be used.
	if len(listenersForValidation) == 0 {
		return nil, fmt.Errorf("no listeners could be found for gateway %s/%s", gateway.Namespace, gateway.Name)
	}

	return listenersForValidation, nil
}
