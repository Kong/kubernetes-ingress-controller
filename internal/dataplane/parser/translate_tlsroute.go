package parser

import (
	"errors"
	"fmt"

	"github.com/kong/go-kong/kong"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

// -----------------------------------------------------------------------------
// Translate TLSRoute - IngressRules Translation
// -----------------------------------------------------------------------------

// ingressRulesFromTLSRoutes processes a list of TLSRoute objects and translates
// then into Kong configuration objects.
func (p *Parser) ingressRulesFromTLSRoutes() ingressRules {
	result := newIngressRules()

	tlsRouteList, err := p.storer.ListTLSRoutes()
	if err != nil {
		p.logger.WithError(err).Error("failed to list TLSRoutes")
		return result
	}

	var errs []error
	for _, tlsroute := range tlsRouteList {
		if p.featureFlags.ExpressionRoutes {
			p.registerResourceFailureNotSupportedForExpressionRoutes(tlsroute)
			continue
		}

		if err := p.ingressRulesFromTLSRoute(&result, tlsroute); err != nil {
			err = fmt.Errorf("TLSRoute %s/%s can't be routed: %w", tlsroute.Namespace, tlsroute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.registerSuccessfullyParsedObject(tlsroute)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			p.logger.Errorf(err.Error())
		}
	}

	return result
}

func (p *Parser) ingressRulesFromTLSRoute(result *ingressRules, tlsroute *gatewayv1alpha2.TLSRoute) error {
	// first we grab the spec and gather some metdata about the object
	spec := tlsroute.Spec

	if len(spec.Hostnames) == 0 {
		return fmt.Errorf("no hostnames provided")
	}
	if len(spec.Rules) == 0 {
		return translators.ErrRouteValidationNoRules
	}

	tlsPassthrough, err := p.isTLSRoutePassthrough(tlsroute)
	if err != nil {
		return err
	}

	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// determine the routes needed to route traffic to services for this rule
		routes, err := generateKongRoutesFromRouteRule(tlsroute, ruleNumber, rule)
		// change protocols in route to tls_passthrough.
		if tlsPassthrough {
			for i := range routes {
				routes[i].Protocols = kong.StringSlice("tls_passthrough")
			}
		}
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		service, err := generateKongServiceFromBackendRefWithRuleNumber(p.logger, p.storer, result, tlsroute, ruleNumber, "tcp", rule.BackendRefs...)
		if err != nil {
			return err
		}
		service.Routes = append(service.Routes, routes...)

		// cache the service to avoid duplicates in further loop iterations
		result.ServiceNameToServices[*service.Service.Name] = service
		result.ServiceNameToParent[*service.Service.Name] = tlsroute
	}

	return nil
}

// isTLSRoutePassthrough returns true if we need to configure TLS passthrough to kong
// for the tlsroute object.
// returns a non-nil error if we failed to get the supported gateway.
func (p *Parser) isTLSRoutePassthrough(tlsroute *gatewayv1alpha2.TLSRoute) (bool, error) {
	// reconcile loop will push TLSRoute object with updated status when
	// gateway is ready and TLSRoute object becomes stable.
	// so we get the supported gateways from status.parents.
	for _, parentStatus := range tlsroute.Status.Parents {
		parentRef := parentStatus.ParentRef

		if parentRef.Group != nil && string(*parentRef.Group) != gatewayv1.GroupName {
			continue
		}

		if parentRef.Kind != nil && *parentRef.Kind != KindGateway {
			continue
		}

		gatewayNamespace := tlsroute.Namespace
		if parentRef.Namespace != nil {
			gatewayNamespace = string(*parentRef.Namespace)
		}

		gateway, err := p.storer.GetGateway(gatewayNamespace, string(parentRef.Name))
		if err != nil {
			if errors.As(err, &store.ErrNotFound{}) {
				// log an error if the gateway expected to support the TLSRoute is not found in our cache.
				p.logger.WithError(err).Errorf("gateway %s/%s not found for TLSRoute %s/%s",
					gatewayNamespace, parentRef.Name, tlsroute.Namespace, tlsroute.Name)
				continue
			}
			return false, err
		}

		// If any of the gateway's listeners is configured to passthrough
		// TLS requests, we return true.
		for _, listener := range gateway.Spec.Listeners {
			if parentRef.SectionName == nil || listener.Name == *parentRef.SectionName {
				if listener.TLS != nil && listener.TLS.Mode != nil &&
					*listener.TLS.Mode == gatewayv1.TLSModePassthrough {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
