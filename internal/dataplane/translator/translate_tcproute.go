package translator

import (
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

// -----------------------------------------------------------------------------
// Translate TCPRoute - IngressRules Translation
// -----------------------------------------------------------------------------

// ingressRulesFromTCPRoutes processes a list of TCPRoute objects and translates
// then into Kong configuration objects.
func (t *Translator) ingressRulesFromTCPRoutes() ingressRules {
	result := newIngressRules()

	tcpRouteList, err := t.storer.ListTCPRoutes()
	if err != nil {
		t.logger.Error(err, "Failed to list TCPRoutes")
		return result
	}

	var errs []error
	for _, tcproute := range tcpRouteList {
		if err := t.ingressRulesFromTCPRoute(&result, tcproute); err != nil {
			err = fmt.Errorf("TCPRoute %s/%s can't be routed: %w", tcproute.Namespace, tcproute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully translated.
			t.registerSuccessfullyTranslatedObject(tcproute)
		}
	}

	if t.featureFlags.ExpressionRoutes {
		applyExpressionToIngressRules(&result)
	}

	for _, err := range errs {
		t.logger.Error(err, "could not generate route from TCPRoute")
	}

	return result
}

func (t *Translator) ingressRulesFromTCPRoute(result *ingressRules, tcproute *gatewayapi.TCPRoute) error {
	spec := tcproute.Spec
	if len(spec.Rules) == 0 {
		return subtranslator.ErrRouteValidationNoRules
	}

	gwPorts := t.getGatewayListeningPorts(tcproute.Namespace, gatewayapi.TCPProtocolType, spec.CommonRouteSpec.ParentRefs)

	// Each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {

		// Determine the routes needed to route traffic to services for this rule.
		routes, err := generateKongRoutesFromRouteRule(tcproute, gwPorts, ruleNumber, rule)
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		service, err := generateKongServiceFromBackendRefWithRuleNumber(t.logger, t.storer, result, tcproute, ruleNumber, "tcp", rule.BackendRefs...)
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
