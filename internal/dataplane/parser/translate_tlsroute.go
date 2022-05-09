package parser

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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
		if err := p.ingressRulesFromTLSRoute(&result, tlsroute); err != nil {
			err = fmt.Errorf("TLSRoute %s/%s can't be routed: %w", tlsroute.Namespace, tlsroute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.ReportKubernetesObjectUpdate(tlsroute)
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
		return fmt.Errorf("no rules provided")
	}

	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// determine the routes needed to route traffic to services for this rule
		routes, err := generateKongRoutesFromTLSRouteRule(tlsroute, ruleNumber, rule)
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		service, err := p.generateKongServiceFromBackendRef(result, tlsroute, ruleNumber, "tcp", rule.BackendRefs...)
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
// Translate TLSRoute - Utils
// -----------------------------------------------------------------------------

// generateKongRoutesFromTLSRouteRule converts an TLSRoute rule to one or more
// Kong Route objects to route traffic to services.
func generateKongRoutesFromTLSRouteRule(
	tlsroute *gatewayv1alpha2.TLSRoute,
	ruleNumber int,
	rule gatewayv1alpha2.TLSRouteRule,
) ([]kongstate.Route, error) {
	// gather the k8s object information and hostnames from the tlsroute
	objectInfo := util.FromK8sObject(tlsroute)

	var routes []kongstate.Route

	if len(rule.BackendRefs) == 0 {
		return routes, fmt.Errorf("TLSRoute rules must include at least one backendRef")
	}

	routeName := kong.String(fmt.Sprintf(
		"tlsroute.%s.%s.%d.%d",
		tlsroute.Namespace,
		tlsroute.Name,
		ruleNumber,
		0,
	))

	hostnames := make([]string, len(tlsroute.Spec.Hostnames))
	for i, hostname := range tlsroute.Spec.Hostnames {
		hostnames[i] = string(hostname)
	}

	r := kongstate.Route{
		Ingress: objectInfo,
		Route: kong.Route{
			Name:      routeName,
			Protocols: kong.StringSlice("tls"),
			SNIs:      kong.StringSlice(hostnames...),
		},
	}

	return append(routes, r), nil
}
