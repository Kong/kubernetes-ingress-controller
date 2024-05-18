package translator

import (
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// -----------------------------------------------------------------------------
// Translate GRPCRoute - IngressRules Translation
// -----------------------------------------------------------------------------

// ingressRulesFromGRPCRoutes processes a list of GRPCRoute objects and translates
// then into Kong configuration objects.
func (t *Translator) ingressRulesFromGRPCRoutes() ingressRules {
	result := newIngressRules()

	grpcRouteList, err := t.storer.ListGRPCRoutes()
	if err != nil {
		t.logger.Error(err, "Failed to list GRPCRoutes")
		return result
	}

	if t.featureFlags.ExpressionRoutes {
		t.ingressRulesFromGRPCRoutesUsingExpressionRoutes(grpcRouteList, &result)
		return result
	}

	var errs []error
	for _, grpcRoute := range grpcRouteList {
		if err := t.ingressRulesFromGRPCRoute(&result, grpcRoute); err != nil {
			err = fmt.Errorf("GRPCRoute %s/%s can't be routed: %w", grpcRoute.Namespace, grpcRoute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully translated.
			t.registerSuccessfullyTranslatedObject(grpcRoute)
		}
	}

	for _, err := range errs {
		t.logger.Error(err, "Could not generate route from GRPCRoute")
	}

	return result
}

func (t *Translator) ingressRulesFromGRPCRoute(result *ingressRules, grpcroute *gatewayapi.GRPCRoute) error {
	// first we grab the spec and gather some metadata about the object
	spec := grpcroute.Spec

	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// Create a service and attach the routes to it.
		service, err := generateKongServiceFromBackendRefWithRuleNumber(
			t.logger, t.storer, result, grpcroute, ruleNumber, t.getProtocolForKongService(grpcroute), grpcBackendRefsToBackendRefs(rule.BackendRefs)...,
		)
		if err != nil {
			return err
		}
		service.Routes = append(service.Routes, subtranslator.GenerateKongRoutesFromGRPCRouteRule(grpcroute, ruleNumber)...)

		// cache the service to avoid duplicates in further loop iterations
		result.ServiceNameToServices[*service.Service.Name] = service
	}

	return nil
}

// ingressRulesFromGRPCRoutesUsingExpressionRoutes translates GRPCRoutes to expression based routes
// when ExpressionRoutes feature flag is enabled.
// Because we need to assign different priorities based on the hostname and match in the specification of GRPCRoutes,
// We need to split the GRPCRoutes into ones with only one hostname and one match, then assign priority to them
// and finally translate the split GRPCRoutes into Kong services and routes with assigned priorities.
func (t *Translator) ingressRulesFromGRPCRoutesUsingExpressionRoutes(grpcRoutes []*gatewayapi.GRPCRoute, result *ingressRules) {
	// first, split GRPCRoutes by hostname and match.
	splitGRPCRouteMatches := []subtranslator.SplitGRPCRouteMatch{}
	// record GRPCRoutes passing the validation and get translated.
	// after they are translated, register the success event in the translator.
	translatedGRPCRoutes := []*gatewayapi.GRPCRoute{}
	for _, grpcRoute := range grpcRoutes {
		splitGRPCRouteMatches = append(splitGRPCRouteMatches, subtranslator.SplitGRPCRoute(grpcRoute)...)
		translatedGRPCRoutes = append(translatedGRPCRoutes, grpcRoute)
	}

	// assign priorities to split GRPCRoutes.
	splitGRPCRouteMatchesWithPriorities := subtranslator.AssignRoutePriorityToSplitGRPCRouteMatches(t.logger, splitGRPCRouteMatches)
	// generate Kong service and route from each split GRPC route with its assigned priority of Kong route.
	for _, splitGRPCRouteMatchWithPriority := range splitGRPCRouteMatchesWithPriorities {
		t.ingressRulesFromGRPCRouteWithPriority(result, splitGRPCRouteMatchWithPriority)
	}

	// register successful translation of GRPCRoutes.
	for _, grpcRoute := range translatedGRPCRoutes {
		t.registerSuccessfullyTranslatedObject(grpcRoute)
	}
}

func (t *Translator) ingressRulesFromGRPCRouteWithPriority(
	rules *ingressRules,
	splitGRPCRouteMatchWithPriority subtranslator.SplitGRPCRouteMatchToPriority,
) {
	match := splitGRPCRouteMatchWithPriority.Match
	grpcRoute := splitGRPCRouteMatchWithPriority.Match.Source
	// (very unlikely that) the rule index split from the source GRPCRoute is larger then length of original rules.
	if len(grpcRoute.Spec.Rules) <= match.RuleIndex {
		t.logger.Error(nil, "Split rule index is greater than the length of rules in source GRPCRoute",
			"rule_index", match.RuleIndex,
			"rule_count", len(grpcRoute.Spec.Rules))
		return
	}
	grpcRouteRule := grpcRoute.Spec.Rules[match.RuleIndex]

	serviceName := subtranslator.KongServiceNameFromSplitGRPCRouteMatch(match)

	// Create a service and attach the routes to it.
	kongService, _ := generateKongServiceFromBackendRefWithName(
		t.logger,
		t.storer,
		rules,
		serviceName,
		grpcRoute,
		t.getProtocolForKongService(grpcRoute),
		grpcBackendRefsToBackendRefs(grpcRouteRule.BackendRefs)...,
	)
	kongService.Routes = append(
		kongService.Routes,
		subtranslator.KongExpressionRouteFromSplitGRPCRouteMatchWithPriority(splitGRPCRouteMatchWithPriority),
	)
	// cache the service to avoid duplicates in further loop iterations
	rules.ServiceNameToServices[serviceName] = kongService
	rules.ServiceNameToParent[serviceName] = grpcRoute
}

func grpcBackendRefsToBackendRefs(grpcBackendRef []gatewayapi.GRPCBackendRef) []gatewayapi.BackendRef {
	backendRefs := make([]gatewayapi.BackendRef, 0, len(grpcBackendRef))

	for _, hRef := range grpcBackendRef {
		backendRefs = append(backendRefs, hRef.BackendRef)
	}
	return backendRefs
}

// getProtocolForKongService returns the protocol for the Kong service based on the Gateway listening ports
func (t *Translator) getProtocolForKongService(grpcRoute *gatewayapi.GRPCRoute) string {
	// When Gateway listens on HTTP use "grpc" protocol for the service. Otherwise for HTTPS use "grpcs".
	if len(t.getGatewayListeningPorts(grpcRoute.Namespace, gatewayapi.HTTPProtocolType, grpcRoute.Spec.ParentRefs)) > 0 {
		return "grpc"
	}
	return "grpcs"
}
