package translator

import (
	"errors"
	"fmt"

	"github.com/kong/go-kong/kong"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

// -----------------------------------------------------------------------------
// Translate TLSRoute - IngressRules Translation
// -----------------------------------------------------------------------------

// ingressRulesFromTLSRoutes processes a list of TLSRoute objects and translates
// then into Kong configuration objects.
func (t *Translator) ingressRulesFromTLSRoutes() ingressRules {
	result := newIngressRules()

	tlsRouteList, err := t.storer.ListTLSRoutes()
	if err != nil {
		t.logger.Error(err, "Failed to list TLSRoutes")
		return result
	}

	var errs []error
	for _, tlsroute := range tlsRouteList {
		if err := t.ingressRulesFromTLSRoute(&result, tlsroute); err != nil {
			err = fmt.Errorf("TLSRoute %s/%s can't be routed: %w", tlsroute.Namespace, tlsroute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully translated.
			t.registerSuccessfullyTranslatedObject(tlsroute)
		}
	}

	if t.featureFlags.ExpressionRoutes {
		applyExpressionToIngressRules(&result)
	}

	for _, err := range errs {
		t.logger.Error(err, "Could not generate route from TLSRoute")
	}

	return result
}

func (t *Translator) ingressRulesFromTLSRoute(result *ingressRules, tlsroute *gatewayapi.TLSRoute) error {
	spec := tlsroute.Spec

	if len(spec.Hostnames) == 0 {
		return fmt.Errorf("no hostnames provided")
	}
	if len(spec.Rules) == 0 {
		return subtranslator.ErrRouteValidationNoRules
	}

	tlsPassthrough, err := t.isTLSRoutePassthrough(tlsroute)
	if err != nil {
		return err
	}

	// Each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// Determine the routes needed to route traffic to services for this rule.
		// TLSRoute matches based on hostname with Gateway listener thus passing gwPorts is pointless.
		routes, err := generateKongRoutesFromRouteRule(tlsroute, nil, ruleNumber, rule)
		// Change protocols in route to tls_passthrough.
		if tlsPassthrough {
			for i := range routes {
				routes[i].Protocols = kong.StringSlice("tls_passthrough")
			}
		}
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		service, err := generateKongServiceFromBackendRefWithRuleNumber(t.logger, t.storer, result, tlsroute, ruleNumber, "tcp", rule.BackendRefs...)
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
func (t *Translator) isTLSRoutePassthrough(tlsroute *gatewayapi.TLSRoute) (bool, error) {
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

		gateway, err := t.storer.GetGateway(gatewayNamespace, string(parentRef.Name))
		if err != nil {
			if errors.As(err, &store.NotFoundError{}) {
				// log an error if the gateway expected to support the TLSRoute is not found in our cache.
				t.logger.Error(err, "Gateway not found for TLSRoute",
					"gateway_namespace", gatewayNamespace,
					"gateway_name", parentRef.Name,
					"tlsroute_namesapce", tlsroute.Namespace,
					"tlsroute_name", tlsroute.Name)
				continue
			}
			return false, err
		}

		// If any of the gateway's listeners is configured to passthrough
		// TLS requests, we return true.
		for _, listener := range gateway.Spec.Listeners {
			if parentRef.SectionName == nil || listener.Name == *parentRef.SectionName {
				if listener.TLS != nil && listener.TLS.Mode != nil &&
					*listener.TLS.Mode == gatewayapi.TLSModePassthrough {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
