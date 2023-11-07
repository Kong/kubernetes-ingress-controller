package gateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/kong/go-kong/kong"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

type routeValidator interface {
	Validate(context.Context, *kong.Route) (bool, string, error)
}

// -----------------------------------------------------------------------------
// Validation - HTTPRoute - Public Functions
// -----------------------------------------------------------------------------

// ValidateHTTPRoute provides a suite of validation for a given HTTPRoute and
// any number of Gateway resources it's attached to that the caller wants to
// have it validated against. It checks supported features, linked objects,
// and uses provided routesValidator to validate the route against Kong Gateway
// validation endpoint.
func ValidateHTTPRoute(
	ctx context.Context,
	routesValidator routeValidator,
	translatorFeatures translator.FeatureFlags,
	httproute *gatewayapi.HTTPRoute,
	attachedGateways ...*gatewayapi.Gateway,
) (bool, string, error) {
	// validate that no unsupported features are in use
	if err := validateHTTPRouteFeatures(httproute, translatorFeatures); err != nil {
		return false, fmt.Sprintf("HTTPRoute spec did not pass validation: %s", err), nil
	}

	// perform Gateway validations for the HTTPRoute (e.g. listener validation, namespace validation, e.t.c.)
	for _, gateway := range attachedGateways {
		// TODO: validate that the namespace is supported by the linked Gateway objects
		//       See: https://github.com/Kong/kubernetes-ingress-controller/issues/2080

		// determine the parentRef for this gateway
		parentRef, err := getParentRefForHTTPRouteGateway(httproute, gateway)
		if err != nil {
			return false, fmt.Sprintf("Couldn't determine parentRefs for httproute: %s", err), nil
		}

		// gather the relevant gateway listeners for the httproute
		listeners, err := getListenersForHTTPRouteValidation(parentRef.SectionName, gateway)
		if err != nil {
			return false, fmt.Sprintf("Couldn't find gateway listeners for httproute: %s", err), nil
		}

		// perform validation of this route against it's linked gateway listeners
		for _, listener := range listeners {
			if err := validateHTTPRouteListener(listener); err != nil {
				return false, fmt.Sprintf("HTTPRoute linked Gateway listeners did not pass validation: %s", err), nil
			}
		}
	}

	ok, msg := validateWithKongGateway(ctx, routesValidator, translatorFeatures, httproute)
	return ok, msg, nil
}

// -----------------------------------------------------------------------------
// Validation - HTTPRoute - Private Functions
// -----------------------------------------------------------------------------

// validateHTTPRouteListener verifies that a given HTTPRoute is configured properly
// for a given gateway listener which it is linked to.
func validateHTTPRouteListener(listener *gatewayapi.Listener) error {
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
func validateHTTPRouteFeatures(httproute *gatewayapi.HTTPRoute, translatorFeatures translator.FeatureFlags) error {
	for _, rule := range httproute.Spec.Rules {
		for _, match := range rule.Matches {
			// We support query parameters matching rules only with expression router.
			if len(match.QueryParams) != 0 {
				if !translatorFeatures.ExpressionRoutes {
					return fmt.Errorf("queryparam matching is supported with expression router only")
				}
			}
		}
		// We don't support any backendRef types except Kubernetes Services.
		for _, ref := range rule.BackendRefs {
			if ref.BackendRef.Group != nil && *ref.BackendRef.Group != "core" && *ref.BackendRef.Group != "" {
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
func getParentRefForHTTPRouteGateway(httproute *gatewayapi.HTTPRoute, gateway *gatewayapi.Gateway) (*gatewayapi.ParentReference, error) {
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
func getListenersForHTTPRouteValidation(sectionName *gatewayapi.SectionName, gateway *gatewayapi.Gateway) ([]*gatewayapi.Listener, error) {
	var listenersForValidation []*gatewayapi.Listener
	if sectionName != nil {
		// only one specified listener is in use, only need to validate the
		// route against that listener.
		for _, listener := range gateway.Spec.Listeners {
			if string(listener.Name) == string(*sectionName) {
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
			if (listener.Protocol) == gatewayapi.HTTPProtocolType ||
				(listener.Protocol) == gatewayapi.HTTPSProtocolType {
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

func validateWithKongGateway(
	ctx context.Context, routesValidator routeValidator, translatorFeatures translator.FeatureFlags, httproute *gatewayapi.HTTPRoute,
) (bool, string) {
	// Translate HTTPRoute to Kong Route object(s) that can be sent directly to the Admin API for validation.
	// Use KIC translator that works both for traditional and expressions based routes.
	var kongRoutes []kong.Route
	var errMsgs []string
	for _, rule := range httproute.Spec.Rules {
		translation := subtranslator.KongRouteTranslation{
			Name:    "validation-attempt",
			Matches: rule.Matches,
			Filters: rule.Filters,
		}
		routes, err := translator.GenerateKongRouteFromTranslation(
			httproute, translation, translatorFeatures.ExpressionRoutes,
		)
		if err != nil {
			errMsgs = append(errMsgs, err.Error())
			continue
		}
		for _, r := range routes {
			kongRoutes = append(kongRoutes, r.Route)
		}
	}
	if len(errMsgs) > 0 {
		return false, validationMsg(errMsgs)
	}
	// Validate by using feature of Kong Gateway.
	for _, kg := range kongRoutes {
		kg := kg
		ok, msg, err := routesValidator.Validate(ctx, &kg)
		if err != nil {
			return false, fmt.Sprintf("Unable to validate HTTPRoute schema: %s", err.Error())
		}
		if !ok {
			errMsgs = append(errMsgs, msg)
		}
	}
	if len(errMsgs) > 0 {
		return false, validationMsg(errMsgs)
	}
	return true, ""
}

func validationMsg(errMsgs []string) string {
	return fmt.Sprintf("HTTPRoute failed schema validation: %s", strings.Join(errMsgs, ", "))
}
