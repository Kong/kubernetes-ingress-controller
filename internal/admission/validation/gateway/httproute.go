package gateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/kong/go-kong/kong"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/admission/validation/kongplugin"
	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
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
	managerClient client.Client,
) (bool, string, error) {
	// Check if route is managed by this controller. If not, we don't need to validate it.
	routeIsManaged, err := ensureHTTPRouteIsManagedByController(ctx, httproute, managerClient)
	if err != nil {
		return false, "", fmt.Errorf("failed to determine whether HTTPRoute is managed by %q controller: %w",
			gatewaycontroller.GetControllerName(), err)
	}
	if !routeIsManaged {
		return true, "", nil
	}

	if err := kongplugin.ValidatePluginUniquenessPerObject(ctx, managerClient, httproute); err != nil {
		return false, fmt.Sprintf("HTTPRoute has invalid KongPlugin annotation: %s", err), nil
	}

	if err := validateHTTPRouteTimeoutBackendRequest(httproute); err != nil {
		return false, fmt.Sprintf("HTTPRoute spec did not pass validation: %s", err), nil
	}

	// Validate that no unsupported features are in use.
	if err := validateHTTPRouteFeatures(httproute, translatorFeatures); err != nil {
		return false, fmt.Sprintf("HTTPRoute spec did not pass validation: %s", err), nil
	}

	// Validate that the route uses only supported annotations.
	if err := validation.ValidateRouteSourceAnnotations(httproute); err != nil {
		return false, fmt.Sprintf("HTTPRoute has invalid Kong annotations: %s", err), nil
	}

	// Validate that the route is valid against Kong Gateway.
	ok, msg := validateWithKongGateway(ctx, routesValidator, translatorFeatures, httproute)
	return ok, msg, nil
}

// -----------------------------------------------------------------------------
// Validation - HTTPRoute - Private Functions
// -----------------------------------------------------------------------------

// parentRefIsGateway returns true if the group/kind of ParentReference is empty or gateway.networking.k8s.io/Gateway.
func parentRefIsGateway(parentRef gatewayapi.ParentReference) bool {
	const KindGateway = gatewayapi.Kind("Gateway")

	return (parentRef.Group == nil || (*parentRef.Group == "" || *parentRef.Group == gatewayapi.V1Group)) &&
		(parentRef.Kind == nil || (*parentRef.Kind == "" || *parentRef.Kind == KindGateway))
}

// ensureHTTPRouteIsManagedByController checks whether the provided HTTPRoute is managed by this controller implementation.
func ensureHTTPRouteIsManagedByController(ctx context.Context, httproute *gatewayapi.HTTPRoute, managerClient client.Client) (bool, error) {
	// In order to be sure whether an HTTPRoute resource is managed by this
	// controller we ignore references to Gateway resources that do not exist.
	for _, parentRef := range httproute.Spec.ParentRefs {
		// Skip the parentRefs that are not Gateways because they cannot refer to the controller.
		// https://github.com/Kong/kubernetes-ingress-controller/issues/5912
		if !parentRefIsGateway(parentRef) {
			continue
		}

		// Determine the namespace of the gateway referenced via parentRef. If no
		// explicit namespace is provided, assume the namespace of the route.
		namespace := httproute.Namespace
		if parentRef.Namespace != nil {
			namespace = string(*parentRef.Namespace)
		}

		// gather the Gateway resource referenced by parentRef and fail validation
		// if there is no such Gateway resource.
		gateway := gatewayapi.Gateway{}
		if err := managerClient.Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      string(parentRef.Name),
		}, &gateway); err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			return false, fmt.Errorf("failed to get Gateway: %w", err)
		}

		// Pull the referenced GatewayClass object from the Gateway.
		gatewayClass := gatewayapi.GatewayClass{}
		if err := managerClient.Get(ctx, client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, &gatewayClass); err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			return false, fmt.Errorf("failed to get GatewayClass: %w", err)
		}

		// Determine ultimately whether the Gateway is managed by this controller implementation.
		if gatewayClass.Spec.ControllerName == gatewaycontroller.GetControllerName() {
			return true, nil
		}
	}

	// If we get here, the HTTPRoute is not managed by this controller.
	return false, nil
}

// validateHTTPRouteFeatures checks for features that are not supported by this
// HTTPRoute implementation and validates that the provided object is not using
// any of those unsupported features.
func validateHTTPRouteFeatures(httproute *gatewayapi.HTTPRoute, translatorFeatures translator.FeatureFlags) error {
	unsupportedFilterMap := map[gatewayapi.HTTPRouteFilterType]struct{}{
		gatewayapi.HTTPRouteFilterRequestMirror: {},
	}
	const (
		KindService = gatewayapi.Kind("Service")
	)

	for ruleIndex, rule := range httproute.Spec.Rules {
		// Filter RequestMirror is not supported.

		// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/3686
		// For URLRewrite, only FullPathHTTPPathModifier is supported.
		for filterIndex, filter := range rule.Filters {
			if _, unsupported := unsupportedFilterMap[filter.Type]; unsupported {
				return fmt.Errorf("rules[%d].filters[%d]: filter type %s is unsupported",
					ruleIndex, filterIndex, filter.Type)
			}

			if filter.Type == gatewayapi.HTTPRouteFilterURLRewrite && filter.URLRewrite != nil {
				// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/3685
				if filter.URLRewrite.Hostname != nil {
					return fmt.Errorf("rules[%d].filters[%d]: filter type %s (with hostname replace) is unsupported",
						ruleIndex, filterIndex, filter.Type)
				}
			}
		}

		for refIndex, ref := range rule.BackendRefs {
			// Specifying filters in backendRef is not supported.
			if len(ref.Filters) != 0 {
				return fmt.Errorf("rules[%d].backendRefs[%d]: filters in backendRef is unsupported",
					ruleIndex, refIndex)
			}

			// We don't support any backendRef types except Kubernetes Services.
			if ref.BackendRef.Group != nil && *ref.BackendRef.Group != "core" && *ref.BackendRef.Group != "" {
				return fmt.Errorf("rules[%d].backendRefs[%d]: %s is not a supported group for httproute backendRefs, only core is supported",
					ruleIndex, refIndex, *ref.BackendRef.Group)
			}
			if ref.BackendRef.Kind != nil && *ref.BackendRef.Kind != KindService {
				return fmt.Errorf("rules[%d].backendRefs[%d]: %s is not a supported kind for httproute backendRefs, only %s is supported",
					ruleIndex, refIndex, *ref.BackendRef.Kind, KindService)
			}
		}

		for matchIndex, match := range rule.Matches {
			// We support query parameters matching rules only with expression router.
			if len(match.QueryParams) != 0 {
				if !translatorFeatures.ExpressionRoutes {
					return fmt.Errorf("rules[%d].matches[%d]: queryparam matching is supported with expression router only",
						ruleIndex, matchIndex)
				}
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// Validation - HTTPRoute - Private Utility Functions
// -----------------------------------------------------------------------------

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

func validateHTTPRouteTimeoutBackendRequest(httproute *gatewayapi.HTTPRoute) error {
	// TODO: remove the validate after we figure out how to handle granular timeout settings
	// (allowing setting timeouts per rule and not enforcing the same timeout for every HTTPRoute's rule).
	// https://github.com/Kong/kubernetes-ingress-controller/issues/5451

	var firstTimeoutFound *gatewayapi.Duration
	for _, rule := range httproute.Spec.Rules {
		if firstTimeoutFound != nil {
			if rule.Timeouts == nil {
				return fmt.Errorf("timeout is set for one of the rules, but not set for another")
			}
			if rule.Timeouts != nil && *rule.Timeouts.BackendRequest != *firstTimeoutFound {
				return fmt.Errorf("timeout is set for one of the rules, but a different value is set in another rule")
			}
		} else if rule.Timeouts != nil && rule.Timeouts.BackendRequest != nil {
			firstTimeoutFound = rule.Timeouts.BackendRequest
		}
	}

	return nil
}
