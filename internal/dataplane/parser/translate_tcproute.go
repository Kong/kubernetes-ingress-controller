package parser

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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
		routes, err := generateKongRoutesFromTCPRouteRule(tcproute, ruleNumber, rule)
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		service, err := p.generateKongServiceFromBackendRef(result, tcproute, ruleNumber, "tcp", rule.BackendRefs...)
		if err != nil {
			return err
		}
		service.Routes = append(service.Routes, routes...)

		// cache the service to avoid duplicates in further loop iterations
		result.ServiceNameToServices[*service.Service.Name] = service
	}

	return nil
}

// -----------------------------------------------------------------------------
// Translate TCPRoute - Utils
// -----------------------------------------------------------------------------

// generateKongRoutesFromTCPRouteRule converts an TCPRoute rule to one or more
// Kong Route objects to route traffic to services.
func generateKongRoutesFromTCPRouteRule(
	tcproute *gatewayv1alpha2.TCPRoute,
	ruleNumber int,
	rule gatewayv1alpha2.TCPRouteRule,
) ([]kongstate.Route, error) {
	// gather the k8s object information and hostnames from the tcproute
	objectInfo := util.FromK8sObject(tcproute)

	var routes []kongstate.Route
	if len(rule.Matches) > 0 {
		// As of 2022-03-04, matches are supported only in experimental CRDs. if you apply a TCPRoute with matches against
		// the stable CRDs, the matches disappear into the ether (only if doing it via client-go, kubectl rejects them)
		// We do not intend to implement these until they are stable per https://github.com/Kong/kubernetes-ingress-controller/issues/2087#issuecomment-1079053290
		return routes, fmt.Errorf("TCPRoute Matches are not yet supported")
	}

	if len(rule.BackendRefs) == 0 {
		return routes, fmt.Errorf("TCPRoute rules must include at least one backendRef")
	}

	routeName := kong.String(fmt.Sprintf(
		"tcproute.%s.%s.%d.%d",
		tcproute.Namespace,
		tcproute.Name,
		ruleNumber,
		0,
	))

	// for now, TCPRoutes provide no means of specifying a destination port other than the backend target port
	// they will once https://gateway-api.sigs.k8s.io/geps/gep-957/ is stable. in the interim, this always uses the
	var destinations []*kong.CIDRPort
	for _, backendRef := range rule.BackendRefs {
		destinations = append(destinations, &kong.CIDRPort{
			Port: kong.Int(int(*backendRef.Port)),
		})
	}

	r := kongstate.Route{
		Ingress: objectInfo,
		Route: kong.Route{
			Name:         routeName,
			Protocols:    kong.StringSlice("tcp"),
			Destinations: destinations,
		},
	}

	return append(routes, r), nil
}
