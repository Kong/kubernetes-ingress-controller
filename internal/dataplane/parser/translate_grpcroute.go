package parser

import (
	"fmt"

	"github.com/bombsimon/logrusr/v4"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
)

// -----------------------------------------------------------------------------
// Translate GRPCRoute - IngressRules Translation
// -----------------------------------------------------------------------------

// ingressRulesFromGRPCRoutes processes a list of GRPCRoute objects and translates
// then into Kong configuration objects.
func (p *Parser) ingressRulesFromGRPCRoutes() ingressRules {
	result := newIngressRules()

	grpcRouteList, err := p.storer.ListGRPCRoutes()
	if err != nil {
		p.logger.WithError(err).Error("failed to list GRPCRoutes")
		return result
	}

	if p.featureFlags.ExpressionRoutes {
		p.ingressRulesFromGRPCRoutesUsingExpressionRoutes(grpcRouteList, &result)
		return result
	}

	var errs []error
	for _, grpcroute := range grpcRouteList {
		if err := p.ingressRulesFromGRPCRoute(&result, grpcroute); err != nil {
			err = fmt.Errorf("GRPCRoute %s/%s can't be routed: %w", grpcroute.Namespace, grpcroute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.registerSuccessfullyParsedObject(grpcroute)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			p.logger.Errorf(err.Error())
		}
	}

	return result
}

func (p *Parser) ingressRulesFromGRPCRoute(result *ingressRules, grpcroute *gatewayv1alpha2.GRPCRoute) error {
	// validate the grpcRoute before it gets translated
	if err := validateGRPCRoute(grpcroute); err != nil {
		return err
	}
	// first we grab the spec and gather some metdata about the object
	spec := grpcroute.Spec

	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// determine the routes needed to route traffic to services for this rule
		var routes []kongstate.Route
		if p.featureFlags.ExpressionRoutes {
			routes = translators.GenerateKongExpressionRoutesFromGRPCRouteRule(grpcroute, ruleNumber)
		} else {
			routes = translators.GenerateKongRoutesFromGRPCRouteRule(grpcroute, ruleNumber, p.featureFlags.RegexPathPrefix)
		}

		// create a service and attach the routes to it
		service, err := generateKongServiceFromBackendRefWithRuleNumber(p.logger, p.storer, result, grpcroute, ruleNumber, "grpcs", grpcBackendRefsToBackendRefs(rule.BackendRefs)...)
		if err != nil {
			return err
		}
		service.Routes = append(service.Routes, routes...)

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
func (p *Parser) ingressRulesFromGRPCRoutesUsingExpressionRoutes(grpcRoutes []*gatewayv1alpha2.GRPCRoute, result *ingressRules) {
	// first, split GRPCRoutes by hostname and match.
	splitGRPCRouteMatches := []translators.SplitGRPCRouteMatch{}
	// record GRPCRoutes passing the validation and get translated.
	// after they are translated, register the success event in the parser.
	translatedGRPCRoutes := []*gatewayv1alpha2.GRPCRoute{}
	for _, grpcRoute := range grpcRoutes {
		// validate the GRPCRoute before it gets split by hostnames and matches.
		if err := validateGRPCRoute(grpcRoute); err != nil {
			p.registerTranslationFailure(err.Error(), grpcRoute)
			continue
		}
		splitGRPCRouteMatches = append(splitGRPCRouteMatches, translators.SplitGRPCRoute(grpcRoute)...)
		translatedGRPCRoutes = append(translatedGRPCRoutes, grpcRoute)
	}

	// assign priorities to split GRPCRoutes.
	splitGRPCRouteMatchesWithPriorities := translators.AssignRoutePriorityToSplitGRPCRouteMatches(logrusr.New(p.logger), splitGRPCRouteMatches)
	// generate Kong service and route from each split GRPC route with its assigned priority of Kong route.
	for _, splitGRPCRouteMatchWithPriority := range splitGRPCRouteMatchesWithPriorities {
		p.ingressRulesFromGRPCRouteWithPriority(result, splitGRPCRouteMatchWithPriority)
	}

	// register successful parses of GRPCRoutes.
	for _, grpcRoute := range translatedGRPCRoutes {
		p.registerSuccessfullyParsedObject(grpcRoute)
	}
}

func (p *Parser) ingressRulesFromGRPCRouteWithPriority(
	rules *ingressRules,
	splitGRPCRouteMatchWithPriority translators.SplitGRPCRouteMatchToPriority,
) {
	match := splitGRPCRouteMatchWithPriority.Match
	grpcRoute := splitGRPCRouteMatchWithPriority.Match.Source
	// (very unlikely that) the rule index split from the source GRPCRoute is larger then length of original rules.
	if len(grpcRoute.Spec.Rules) <= match.RuleIndex {
		p.logger.Infof("WARN: split rule index %d is larger then the length of rules in source GRPCRoute %d",
			match.RuleIndex, len(grpcRoute.Spec.Rules))
		return
	}
	grpcRouteRule := grpcRoute.Spec.Rules[match.RuleIndex]
	backendRefs := grpcBackendRefsToBackendRefs(grpcRouteRule.BackendRefs)

	serviceName := translators.KongServiceNameFromSplitGRPCRouteMatch(match)

	kongService, _ := generateKongServiceFromBackendRefWithName(
		p.logger,
		p.storer,
		rules,
		serviceName,
		grpcRoute,
		"grpcs",
		backendRefs...,
	)
	kongService.Routes = append(
		kongService.Routes,
		translators.KongExpressionRouteFromSplitGRPCRouteMatchWithPriority(splitGRPCRouteMatchWithPriority),
	)
	// cache the service to avoid duplicates in further loop iterations
	rules.ServiceNameToServices[serviceName] = kongService
	rules.ServiceNameToParent[serviceName] = grpcRoute
}

func grpcBackendRefsToBackendRefs(grpcBackendRef []gatewayv1alpha2.GRPCBackendRef) []gatewayv1.BackendRef {
	backendRefs := make([]gatewayv1.BackendRef, 0, len(grpcBackendRef))

	for _, hRef := range grpcBackendRef {
		backendRefs = append(backendRefs, hRef.BackendRef)
	}
	return backendRefs
}

func validateGRPCRoute(grpcRoute *gatewayv1alpha2.GRPCRoute) error {
	if len(grpcRoute.Spec.Hostnames) == 0 {
		if len(grpcRoute.Spec.Rules) == 0 {
			return translators.ErrRouteValidationNoRules
		}
	}
	return nil
}
