package parser

import (
	"fmt"

	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// -----------------------------------------------------------------------------
// Translate TCPRoute - IngressRules Translation
// -----------------------------------------------------------------------------

// ingressRulesFromTCPRoutes processes a list of TCPRoute objects and translates
// then into Kong configuration objects.
func (p *Parser) ingressRulesFromTCPRoutes() ingressRules {
	result := newIngressRules()

	tcpRouteList, err := p.storer.ListTCPRoutes()
	if err != nil {
		p.logger.WithError(err).Error("failed to list TCPRoutes")
		return result
	}

	var errs []error
	for _, tcproute := range tcpRouteList {
		if err := p.ingressRulesFromTCPRoute(&result, tcproute); err != nil {
			err = fmt.Errorf("TCPRoute %s/%s can't be routed: %w", tcproute.Namespace, tcproute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.ReportKubernetesObjectUpdate(tcproute)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			p.logger.Errorf(err.Error())
		}
	}

	return result
}

func (p *Parser) ingressRulesFromTCPRoute(result *ingressRules, tcproute *gatewayv1alpha2.TCPRoute) error {
	// first we grab the spec and gather some metdata about the object
	spec := tcproute.Spec

	// validation for TCPRoutes will happen at a higher layer, but in spite of that we run
	// validation at this level as well as a fallback so that if routes are posted which
	// are invalid somehow make it past validation (e.g. the webhook is not enabled) we can
	// at least try to provide a helpful message about the situation in the manager logs.
	if len(spec.Rules) == 0 {
		return fmt.Errorf("no rules provided")
	}

	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// TODO: add this to a generic TCPRoute validation, and then we should probably
		//       simply be calling validation on each tcproute object at the begininning
		//       of the topmost list.
		if len(rule.BackendRefs) == 0 {
			return fmt.Errorf("missing backendRef in rule")
		}

		// determine the routes needed to route traffic to services for this rule
		routes, err := generateKongRoutesFromRouteRule(tcproute, ruleNumber, rule)
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		service, err := generateKongServiceFromBackendRefWithRuleNumber(p.logger, p.storer, result, tcproute, ruleNumber, "tcp", rule.BackendRefs...)
		if err != nil {
			return err
		}
		service.Routes = append(service.Routes, routes...)

		// cache the service to avoid duplicates in further loop iterations
		result.ServiceNameToServices[*service.Service.Name] = service
	}

	return nil
}
