package parser

import (
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
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
		p.logger.Error(err, "Failed to list TCPRoutes")
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
			p.registerSuccessfullyParsedObject(tcproute)
		}
	}

	if p.featureFlags.ExpressionRoutes {
		applyExpressionToIngressRules(&result)
	}

	for _, err := range errs {
		p.logger.Error(err, "could not generate route from TCPRoute")
	}

	return result
}

func (p *Parser) ingressRulesFromTCPRoute(result *ingressRules, tcproute *gatewayapi.TCPRoute) error {
	spec := tcproute.Spec
	if len(spec.Rules) == 0 {
		return translators.ErrRouteValidationNoRules
	}

	gwPorts := p.getGatewayListeningPorts(tcproute.Namespace, gatewayapi.TCPProtocolType, spec.CommonRouteSpec.ParentRefs)

	// Each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {

		// Determine the routes needed to route traffic to services for this rule.
		routes, err := generateKongRoutesFromRouteRule(tcproute, gwPorts, ruleNumber, rule)
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
		result.ServiceNameToParent[*service.Service.Name] = tcproute
	}

	return nil
}
