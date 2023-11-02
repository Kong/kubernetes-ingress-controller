package parser

import (
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// -----------------------------------------------------------------------------
// Translate UDPRoute - IngressRules Translation
// -----------------------------------------------------------------------------

// ingressRulesFromUDPRoutes processes a list of UDPRoute objects and translates
// then into Kong configuration objects.
func (p *Parser) ingressRulesFromUDPRoutes() ingressRules {
	result := newIngressRules()

	udpRouteList, err := p.storer.ListUDPRoutes()
	if err != nil {
		p.logger.Error(err, "Failed to list UDPRoutes")
		return result
	}

	var errs []error
	for _, udproute := range udpRouteList {
		if err := validateUDPRoute(udproute); err != nil {
			errs = append(errs, err)
			p.registerTranslationFailure(err.Error(), udproute)
			continue
		}

		if err := p.ingressRulesFromUDPRoute(&result, udproute); err != nil {
			err = fmt.Errorf("UDPRoute %s/%s can't be routed: %w", udproute.Namespace, udproute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.registerSuccessfullyParsedObject(udproute)
		}
	}

	// Translate generated Kong Route to expression based route.
	if p.featureFlags.ExpressionRoutes {
		applyExpressionToIngressRules(&result)
	}

	for _, err := range errs {
		p.logger.Error(err, "could not generate route from UDPRoute")
	}

	return result
}

func (p *Parser) ingressRulesFromUDPRoute(result *ingressRules, udproute *gatewayapi.UDPRoute) error {
	spec := udproute.Spec
	if len(spec.Rules) == 0 {
		return translators.ErrRouteValidationNoRules
	}

	gwPorts := p.getGatewayListeningPorts(udproute.Namespace, gatewayapi.UDPProtocolType, spec.CommonRouteSpec.ParentRefs)
	// Each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// Determine the routes needed to route traffic to services for this rule.
		routes, err := generateKongRoutesFromRouteRule(udproute, gwPorts, ruleNumber, rule)
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		service, err := generateKongServiceFromBackendRefWithRuleNumber(p.logger, p.storer, result, udproute, ruleNumber, "udp", rule.BackendRefs...)
		if err != nil {
			return err
		}
		service.Routes = append(service.Routes, routes...)

		// cache the service to avoid duplicates in further loop iterations
		result.ServiceNameToServices[*service.Service.Name] = service
		result.ServiceNameToParent[*service.Service.Name] = udproute
	}

	return nil
}

// validateUDPRoute validates UDPRoute, and return a translation error if the spec is invalid.
// Validation for UDPRoutes will happen at a higher layer, but in spite of that we run
// validation at this level as well as a fallback so that if routes are posted which
// are invalid somehow make it past validation (e.g. the webhook is not enabled) we can
// at least try to provide a helpful message about the situation in the manager logs.
func validateUDPRoute(udproute *gatewayapi.UDPRoute) error {
	if len(udproute.Spec.Rules) == 0 {
		return translators.ErrRouteValidationNoRules
	}
	for _, rule := range udproute.Spec.Rules {
		if len(rule.BackendRefs) == 0 {
			return translators.ErrRotueValidationRuleNoBackendRef
		}
	}
	return nil
}
