package translator

import (
	"errors"
	"fmt"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

// -----------------------------------------------------------------------------
// Translate HTTPRoute - IngressRules Translation
// -----------------------------------------------------------------------------

// ingressRulesFromHTTPRoutes processes a list of HTTPRoute objects and translates
// then into Kong configuration objects.
func (t *Translator) ingressRulesFromHTTPRoutes() ingressRules {
	result := newIngressRules()

	httpRouteList, err := t.storer.ListHTTPRoutes()
	if err != nil {
		t.logger.Error(err, "Failed to list HTTPRoutes")
		return result
	}

	httpRoutesToTranslate := make([]*gatewayapi.HTTPRoute, 0, len(httpRouteList))
	for _, httproute := range httpRouteList {
		// Validate each HTTPRoute before translating and register translation failures if an HTTPRoute is invalid.
		if err := validateHTTPRoute(httproute, t.featureFlags); err != nil {
			t.registerTranslationFailure(fmt.Sprintf("HTTPRoute can't be routed: %v", err), httproute)
			continue
		}
		httpRoutesToTranslate = append(httpRoutesToTranslate, httproute)
	}

	t.ingressRulesFromHTTPRoutesWithCombinedService(httpRoutesToTranslate, &result)
	return result
}

func validateHTTPRoute(httproute *gatewayapi.HTTPRoute, featureFlags FeatureFlags) error {
	spec := httproute.Spec

	// validation for HTTPRoutes will happen at a higher layer, but in spite of that we run
	// validation at this level as well as a fallback so that if routes are posted which
	// are invalid somehow make it past validation (e.g. the webhook is not enabled) we can
	// at least try to provide a helpful message about the situation in the manager logs.
	if len(spec.Rules) == 0 {
		return subtranslator.ErrRouteValidationNoRules
	}

	// Kong supports query parameter match only with expression router,
	// so we return error when query param match is specified and expression router is not enabled in the translator.
	if !featureFlags.ExpressionRoutes {
		for _, rule := range spec.Rules {
			for _, match := range rule.Matches {
				if len(match.QueryParams) > 0 {
					return subtranslator.ErrRouteValidationQueryParamMatchesUnsupported
				}
			}
		}
	}

	return nil
}

// ingressRulesFromHTTPRoutesWithCombinedService translates HTTPRoutes to ingress rules with combined Kong services
// of rules cross different HTTPRoutes (but still in the same namespace).
// TODO: handle the case of using expression based routes: https://github.com/Kong/kubernetes-ingress-controller/issues/6727.
func (t *Translator) ingressRulesFromHTTPRoutesWithCombinedService(httpRoutes []*gatewayapi.HTTPRoute, result *ingressRules) {
	kongstateServices, routeTranslationFailures := subtranslator.TranslateHTTPRoutesToKongstateServices(
		t.logger,
		t.storer,
		httpRoutes,
		t.featureFlags.ExpressionRoutes,
	)
	for serviceName, service := range kongstateServices {
		result.ServiceNameToServices[serviceName] = service
		result.ServiceNameToParent[serviceName] = service.Parent
	}
	for _, httproute := range httpRoutes {
		namespacedName := k8stypes.NamespacedName{
			Namespace: httproute.Namespace,
			Name:      httproute.Name,
		}
		translationFailures, hasError := routeTranslationFailures[namespacedName]
		if hasError && len(translationFailures) > 0 {
			t.failuresCollector.PushResourceFailure(
				fmt.Sprintf("HTTPRoute can't be routed: %v", errors.Join(translationFailures...)),
				httproute,
			)
			continue
		}
		t.registerSuccessfullyTranslatedObject(httproute)
	}
}

// -----------------------------------------------------------------------------
// Translate HTTPRoute - Utils
// -----------------------------------------------------------------------------

// getHTTPRouteHostnamesAsSliceOfStrings translates the hostnames defined in an
// HTTPRoute specification into a []*string slice, which is the type required by translating to matchers
// in expression based routes.
func getHTTPRouteHostnamesAsSliceOfStrings(httproute *gatewayapi.HTTPRoute) []string {
	return lo.Map(httproute.Spec.Hostnames, func(h gatewayapi.Hostname, _ int) string {
		return string(h)
	})
}

// getHTTPRouteHostnamesAsSliceOfStringPointers translates the hostnames defined
// in an HTTPRoute specification into a []*string slice, which is the type required
// by kong.Route{}.
func getHTTPRouteHostnamesAsSliceOfStringPointers(httproute *gatewayapi.HTTPRoute) []*string {
	return lo.Map(httproute.Spec.Hostnames, func(h gatewayapi.Hostname, _ int) *string {
		return kong.String(string(h))
	})
}

// GenerateKongRouteFromTranslation generates Kong routes from HTTPRoute
// pointing to a specific backend. It is used for both traditional and expression based routes.
func GenerateKongRouteFromTranslation(
	httproute *gatewayapi.HTTPRoute,
	translation subtranslator.KongRouteTranslation,
	expressionRoutes bool,
) ([]kongstate.Route, error) {
	// gather the k8s object information and hostnames from the httproute
	objectInfo := util.FromK8sObject(httproute)
	tags := util.GenerateTagsForObject(httproute)

	// translate to expression based routes when expressionRoutes is enabled.
	if expressionRoutes {
		// get the hostnames from the HTTPRoute
		hostnames := getHTTPRouteHostnamesAsSliceOfStrings(httproute)
		return subtranslator.GenerateKongExpressionRoutesFromHTTPRouteMatches(
			translation,
			objectInfo,
			hostnames,
			tags,
		)
	}

	// get the hostnames from the HTTPRoute
	hostnames := getHTTPRouteHostnamesAsSliceOfStringPointers(httproute)
	return subtranslator.GenerateKongRoutesFromHTTPRouteMatches(
		translation.Name,
		translation.Matches,
		translation.Filters,
		objectInfo,
		hostnames,
		tags,
	)
}

func httpBackendRefsToBackendRefs(httpBackendRef []gatewayapi.HTTPBackendRef) []gatewayapi.BackendRef {
	backendRefs := make([]gatewayapi.BackendRef, 0, len(httpBackendRef))

	for _, hRef := range httpBackendRef {
		backendRefs = append(backendRefs, hRef.BackendRef)
	}
	return backendRefs
}

// ingressRulesFromSplitHTTPRouteMatchWithPriority translates a single match split from HTTPRoute
// to ingress rule, including Kong service and Kong route.
func (t *Translator) ingressRulesFromSplitHTTPRouteMatchWithPriority(
	rules *ingressRules,
	httpRouteMatchWithPriority subtranslator.SplitHTTPRouteMatchToKongRoutePriority,
) error {
	match := httpRouteMatchWithPriority.Match
	httpRoute := httpRouteMatchWithPriority.Match.Source
	if match.RuleIndex >= len(httpRoute.Spec.Rules) {
		t.logger.Error(nil, "Split match has rule out of bound of rules in source HTTPRoute",
			"rule_index", match.RuleIndex, "rule_count", len(httpRoute.Spec.Rules))
		return nil
	}

	rule := httpRoute.Spec.Rules[match.RuleIndex]
	backendRefs := httpBackendRefsToBackendRefs(rule.BackendRefs)
	serviceName := subtranslator.KongServiceNameFromSplitHTTPRouteMatch(httpRouteMatchWithPriority.Match)

	kongService, err := generateKongServiceFromBackendRefWithName(
		t.logger,
		t.storer,
		rules,
		serviceName,
		httpRoute,
		"http",
		backendRefs...,
	)
	if err != nil {
		return err
	}

	additionalRoutes, err := subtranslator.KongExpressionRouteFromHTTPRouteMatchWithPriority(httpRouteMatchWithPriority)
	if err != nil {
		return err
	}

	kongService.Routes = append(
		kongService.Routes,
		*additionalRoutes,
	)
	// cache the service to avoid duplicates in further loop iterations
	rules.ServiceNameToServices[serviceName] = kongService
	rules.ServiceNameToParent[serviceName] = httpRoute
	return nil
}
