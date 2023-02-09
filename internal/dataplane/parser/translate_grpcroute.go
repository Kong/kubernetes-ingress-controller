package parser

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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

	var errs []error
	for _, grpcroute := range grpcRouteList {
		if err := p.ingressRulesFromGRPCRoute(&result, grpcroute); err != nil {
			err = fmt.Errorf("GRPCRoute %s/%s can't be routed: %w", grpcroute.Namespace, grpcroute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.ReportKubernetesObjectUpdate(grpcroute)
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
		return errRouteValidationNoRules
	}

	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// determine the routes needed to route traffic to services for this rule
		routes := generateKongRoutesFromGRPCRouteRule(grpcroute, ruleNumber, rule)

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

func generateKongRoutesFromGRPCRouteRule(grpcroute *gatewayv1alpha2.GRPCRoute, ruleNumber int, rule gatewayv1alpha2.GRPCRouteRule) []kongstate.Route {
	routes := make([]kongstate.Route, 0, len(rule.Matches))

	// gather the k8s object information and hostnames from the httproute
	ingressObjectInfo := util.FromK8sObject(grpcroute)

	for matchNumber := range rule.Matches {
		routeName := fmt.Sprintf(
			"httproute.%s.%s.%d.%d",
			grpcroute.Namespace,
			grpcroute.Name,
			ruleNumber,
			matchNumber,
		)

		r := kongstate.Route{
			Ingress: ingressObjectInfo,
			Route: kong.Route{
				Name:      kong.String(routeName),
				Protocols: kong.StringSlice("grpc", "grpcs"),
			},
		}

		if len(grpcroute.Spec.Hostnames) > 0 {
			r.Hosts = getGRPCRouteHostnamesAsSliceOfStringPointers(grpcroute)
		}

		routes = append(routes, r)
	}

	return routes
}

func grpcBackendRefsToBackendRefs(grpcBackendRef []gatewayv1alpha2.GRPCBackendRef) []gatewayv1beta1.BackendRef {
	backendRefs := make([]gatewayv1beta1.BackendRef, 0, len(grpcBackendRef))

	for _, hRef := range grpcBackendRef {
		backendRefs = append(backendRefs, hRef.BackendRef)
	}
	return backendRefs
}

// -----------------------------------------------------------------------------
// Translate HTTPRoute - Utils
// -----------------------------------------------------------------------------

// getGRPCRouteHostnamesAsSliceOfStringPointers translates the hostnames defined
// in an GRPCRoute specification into a []*string slice, which is the type required
// by kong.Route{}.
func getGRPCRouteHostnamesAsSliceOfStringPointers(grpcroute *gatewayv1alpha2.GRPCRoute) []*string {
	return lo.Map(grpcroute.Spec.Hostnames, func(h gatewayv1beta1.Hostname, _ int) *string {
		return lo.ToPtr(string(h))
	})
}
