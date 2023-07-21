package parser

import (
	"fmt"

	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/bombsimon/logrusr/v2"
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
	// first we grab the spec and gather some metdata about the object
	spec := grpcroute.Spec

	if len(spec.Rules) == 0 {
		return translators.ErrRouteValidationNoRules
	}

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
// Because we need to assign different priorities based on the hostname and match in the specification of HTTPRoutes,
// We need to split the GRPCRoutes into ones with only one hostname and one match, then assign priority to them
// and finally translate the split GRPCRoutes into Kong services and routes with assigned priorities.
func (p *Parser) ingressRulesFromGRPCRoutesUsingExpressionRoutes(grpcRoutes []*gatewayv1alpha2.GRPCRoute, result *ingressRules) {
	// first, split GRPCRoutes by hostname and match.
	splitGRPCRoutes := []*gatewayv1alpha2.GRPCRoute{}
	for _, grpcRoute := range grpcRoutes {
		if len(grpcRoute.Spec.Rules) == 0 {
			p.registerTranslationFailure(
				translators.ErrRouteValidationNoRules.Error(),
				grpcRoute,
			)
			continue
		}
		splitGRPCRoutes = append(splitGRPCRoutes, translators.SplitGRPCRoute(grpcRoute)...)
	}

	// assign priorities to split GRPCRoutes.
	splitGRPCRoutesWithPriorities := translators.AssignRoutePriorityToSplitGRPCRoutes(logrusr.New(p.logger), splitGRPCRoutes)
	// generate Kong service and route from each split GRPC route with its assigned priority of Kong route.
	for _, splitGRPCRouteWithPriority := range splitGRPCRoutesWithPriorities {
		p.ingressRulesFromGRPCRouteWithPriority(result, splitGRPCRouteWithPriority)
	}
}

func (p *Parser) ingressRulesFromGRPCRouteWithPriority(
	rules *ingressRules,
	grpcRouteWithPriority translators.SplitGRPCRouteToKongRoutePriority,
) {
	grpcRoute := grpcRouteWithPriority.GRPCRoute
	if len(grpcRoute.Spec.Rules) != 1 {
		return
	}
	grpcRouteRule := grpcRoute.Spec.Rules[0]
	backendRefs := grpcBackendRefsToBackendRefs(grpcRouteRule.BackendRefs)

	serviceName := translators.KongServiceNameFromSplitGRPCRoute(grpcRoute)

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
		translators.GenerateKongExpressionRouteFromSplitGRPCRouteWithPriority(grpcRouteWithPriority),
	)
	// cache the service to avoid duplicates in further loop iterations
	rules.ServiceNameToServices[serviceName] = kongService
	rules.ServiceNameToParent[serviceName] = grpcRoute
}

func grpcBackendRefsToBackendRefs(grpcBackendRef []gatewayv1alpha2.GRPCBackendRef) []gatewayv1beta1.BackendRef {
	backendRefs := make([]gatewayv1beta1.BackendRef, 0, len(grpcBackendRef))

	for _, hRef := range grpcBackendRef {
		backendRefs = append(backendRefs, hRef.BackendRef)
	}
	return backendRefs
}
