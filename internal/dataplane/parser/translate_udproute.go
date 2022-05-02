package parser

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
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
		p.logger.WithError(err).Errorf("failed to list UDPRoutes")
		return result
	}

	var errs []error
	for _, udproute := range udpRouteList {
		if err := ingressRulesFromUDPRoute(&result, udproute); err != nil {
			err = fmt.Errorf("UDPRoute %s/%s can't be routed: %w", udproute.Namespace, udproute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.ReportKubernetesObjectUpdate(udproute)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			p.logger.Errorf(err.Error())
		}
	}

	return result
}

func ingressRulesFromUDPRoute(result *ingressRules, udproute *gatewayv1alpha2.UDPRoute) error {
	// first we grab the spec and gather some metdata about the object
	spec := udproute.Spec

	// validation for UDPRoutes will happen at a higher layer, but in spite of that we run
	// validation at this level as well as a fallback so that if routes are posted which
	// are invalid somehow make it past validation (e.g. the webhook is not enabled) we can
	// at least try to provide a helpful message about the situation in the manager logs.
	if len(spec.Rules) == 0 {
		return fmt.Errorf("no rules provided")
	}

	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// TODO: add this to a generic UDPRoute validation, and then we should probably
		//       simply be calling validation on each udproute object at the begininning
		//       of the topmost list.
		if len(rule.BackendRefs) == 0 {
			return fmt.Errorf("missing backendRef in rule")
		}

		// determine the routes needed to route traffic to services for this rule
		routes, err := generateKongRoutesFromUDPRouteRule(udproute, ruleNumber, rule)
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		service, err := generateKongServiceFromUDPRouteBackendRef(result, udproute, ruleNumber, rule.BackendRefs...)
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
// Translate UDPRoute - Utils
// -----------------------------------------------------------------------------

// generateKongRoutesFromUDPRouteRule converts an UDPRoute rule to one or more
// Kong Route objects to route traffic to services.
func generateKongRoutesFromUDPRouteRule(udproute *gatewayv1alpha2.UDPRoute, ruleNumber int,
	rule gatewayv1alpha2.UDPRouteRule) ([]kongstate.Route, error) {
	// gather the k8s object information and hostnames from the udproute
	objectInfo := util.FromK8sObject(udproute)

	var routes []kongstate.Route
	if len(rule.Matches) > 0 {
		// As of 2022-03-04, matches are supported only in experimental CRDs. if you apply a UDPRoute with matches against
		// the stable CRDs, the matches disappear into the ether (only if doing it via client-go, kubectl rejects them)
		// We do not intend to implement these until they are stable per https://github.com/Kong/kubernetes-ingress-controller/issues/2087#issuecomment-1079053290
		return routes, fmt.Errorf("UDPRoute Matches are not yet supported")
	}

	if len(rule.BackendRefs) == 0 {
		return routes, fmt.Errorf("UDPRoute rules must include at least one backendRef")
	}

	routeName := kong.String(fmt.Sprintf(
		"udproute.%s.%s.%d.%d",
		udproute.Namespace,
		udproute.Name,
		ruleNumber,
		0,
	))

	// for now, UDPRoutes provide no means of specifying a destination port other than the backend target port
	// they will once https://gateway-api.sigs.k8s.io/geps/gep-957/ is stable. in the interim, this always uses the
	// backend target
	var destinations []*kong.CIDRPort
	for _, backendRef := range rule.BackendRefs {
		destinations = append(destinations, &kong.CIDRPort{Port: kong.Int(int(*backendRef.Port))})
	}

	r := kongstate.Route{
		Ingress: objectInfo,
		Route: kong.Route{
			Name:         routeName,
			Protocols:    kong.StringSlice("udp"),
			Destinations: destinations,
		},
	}

	routes = append(routes, r)

	return routes, nil
}

// generateKongServiceFromUDPRouteBackendRef converts a provided backendRef for an UDPRoute
// into a kong.Service so that routes for that object can be attached to the Service.
// TODO add a generic backendRef handler for all GW routes. HTTPRoute needs a wrapper because it uses a special wrapped
// type with filters. Deferred til after https://github.com/Kong/kubernetes-ingress-controller/issues/2166 though
// we probably shouldn't see much change to the service (just the upstream it references in Host)
func generateKongServiceFromUDPRouteBackendRef(
	result *ingressRules,
	udproute *gatewayv1alpha2.UDPRoute,
	ruleNumber int,
	backendRefs ...gatewayv1alpha2.BackendRef,
) (kongstate.Service, error) {
	// at least one backendRef must be present
	if len(backendRefs) == 0 {
		return kongstate.Service{}, fmt.Errorf("no backendRefs present for UDPRoute: %s/%s", udproute.Namespace, udproute.Name)
	}

	// create a kongstate backend for each UDPRoute backendRef
	backends := make(kongstate.ServiceBackends, 0, len(backendRefs))
	for _, backendRef := range backendRefs {
		// convert each backendRef into a kongstate.ServiceBackend
		backends = append(backends, kongstate.ServiceBackend{
			Name: string(backendRef.Name),
			PortDef: kongstate.PortDef{
				Mode:   kongstate.PortModeByNumber,
				Number: int32(*backendRef.Port),
			},
			Weight: backendRef.Weight,
		})
	}

	// the service name needs to uniquely identify this service given it's list of
	// one or more backends.
	serviceName := fmt.Sprintf("%s.%d", getUniqueKongServiceNameForObject(udproute), ruleNumber)

	// the service host needs to be a resolvable name due to legacy logic so we'll
	// use the anchor backendRef as the basis for the name
	serviceHost := serviceName

	// check if the service is already known, and if not create it
	service, ok := result.ServiceNameToServices[serviceName]
	if !ok {
		service = kongstate.Service{
			Service: kong.Service{
				Name:           kong.String(serviceName),
				Host:           kong.String(serviceHost),
				Protocol:       kong.String("udp"),
				ConnectTimeout: kong.Int(DefaultServiceTimeout),
				ReadTimeout:    kong.Int(DefaultServiceTimeout),
				WriteTimeout:   kong.Int(DefaultServiceTimeout),
				Retries:        kong.Int(DefaultRetries),
			},
			Namespace: udproute.Namespace,
			Backends:  backends,
		}
	}

	return service, nil
}
