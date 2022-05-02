package parser

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// -----------------------------------------------------------------------------
// Translate HTTPRoute - IngressRules Translation
// -----------------------------------------------------------------------------

// ingressRulesFromHTTPRoutes processes a list of HTTPRoute objects and translates
// then into Kong configuration objects.
func (p *Parser) ingressRulesFromHTTPRoutes() ingressRules {
	result := newIngressRules()

	httpRouteList, err := p.storer.ListHTTPRoutes()
	if err != nil {
		p.logger.WithError(err).Error("failed to list HTTPRoutes")
		return result
	}

	var errs []error
	for _, httproute := range httpRouteList {
		if err := ingressRulesFromHTTPRoute(&result, httproute); err != nil {
			err = fmt.Errorf("HTTPRoute %s/%s can't be routed: %w", httproute.Namespace, httproute.Name, err)
			errs = append(errs, err)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.ReportKubernetesObjectUpdate(httproute)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			p.logger.Errorf(err.Error())
		}
	}

	return result
}

func ingressRulesFromHTTPRoute(result *ingressRules, httproute *gatewayv1alpha2.HTTPRoute) error {
	// first we grab the spec and gather some metdata about the object
	spec := httproute.Spec

	// validation for HTTPRoutes will happen at a higher layer, but in spite of that we run
	// validation at this level as well as a fallback so that if routes are posted which
	// are invalid somehow make it past validation (e.g. the webhook is not enabled) we can
	// at least try to provide a helpful message about the situation in the manager logs.
	if len(spec.Rules) == 0 {
		return fmt.Errorf("no rules provided")
	}

	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range spec.Rules {
		// TODO: add this to a generic HTTPRoute validation, and then we should probably
		//       simply be calling validation on each httproute object at the begininning
		//       of the topmost list.
		if len(rule.BackendRefs) == 0 {
			return fmt.Errorf("missing backendRef in rule")
		}

		// determine the routes needed to route traffic to services for this rule
		routes, err := generateKongRoutesFromHTTPRouteRule(httproute, ruleNumber, rule)
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		service, err := generateKongServiceFromHTTPRouteBackendRef(result, httproute, ruleNumber, rule.BackendRefs...)
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
// Translate HTTPRoute - Utils
// -----------------------------------------------------------------------------

// getHTTPRouteHostnamesAsSliceOfStringPointers translates the hostnames defined
// in an HTTPRoute specification into a []*string slice, which is the type required
// by kong.Route{}.
func getHTTPRouteHostnamesAsSliceOfStringPointers(httproute *gatewayv1alpha2.HTTPRoute) []*string {
	hostnames := make([]*string, 0, len(httproute.Spec.Hostnames))
	for _, hostname := range httproute.Spec.Hostnames {
		hostnames = append(hostnames, kong.String(string(hostname)))
	}
	return hostnames
}

// generateKongRoutesFromHTTPRouteRule converts an HTTPRoute rule to one or more
// Kong Route objects to route traffic to services. This function will accept an
// HTTPRoute that does not include any matches as long as it includes hostnames
// to route traffic to the backend service with, in this case it will use the default
// path prefix routing option for that service in addition to hostname routing.
// If an HTTPRoute is provided that has matches that include any unsupported matching
// configurations, this will produce an error and the route is considered invalid.
func generateKongRoutesFromHTTPRouteRule(httproute *gatewayv1alpha2.HTTPRoute, ruleNumber int, rule gatewayv1alpha2.HTTPRouteRule) ([]kongstate.Route, error) {
	// gather the k8s object information and hostnames from the httproute
	objectInfo := util.FromK8sObject(httproute)
	hostnames := getHTTPRouteHostnamesAsSliceOfStringPointers(httproute)

	// the HTTPRoute specification upstream specifically defines matches as
	// independent (e.g. each match is an OR with other matches, not an AND).
	// Therefore we treat each match rule as a separate Kong Route, so we iterate through
	// all matches to determine all the routes that will be needed for the services.
	var routes []kongstate.Route
	if len(rule.Matches) > 0 {
		for matchNumber, match := range rule.Matches {
			// determine the name of the route, identify it as a route that belongs
			// to a Kubernetes HTTPRoute object.
			routeName := kong.String(fmt.Sprintf(
				"httproute.%s.%s.%d.%d",
				httproute.Namespace,
				httproute.Name,
				ruleNumber,
				matchNumber,
			))

			// TODO: implement query param matches
			if len(match.QueryParams) > 0 {
				return nil, fmt.Errorf("query param matches are not yet supported")
			}

			// build the route object using the method and pathing information
			r := kongstate.Route{
				Ingress: objectInfo,
				Route: kong.Route{
					Name:         routeName,
					Protocols:    kong.StringSlice("http", "https"),
					PreserveHost: kong.Bool(true),
				},
			}

			// attach any hostnames associated with the httproute
			if len(hostnames) > 0 {
				r.Hosts = hostnames
			}

			// configure path matching information about the route if paths matching was defined
			// Kong automatically infers whether or not a path is a regular expression and uses a prefix match by
			// default it it is not. For those types, we use the path value as-is and let Kong determine the type.
			// For exact matches, we transform the path into a regular expression that terminates after the value
			if match.Path != nil {
				if *match.Path.Type == gatewayv1alpha2.PathMatchExact {
					terminated := *match.Path.Value + "$"
					r.Route.Paths = []*string{&terminated}
				} else if *match.Path.Type == gatewayv1alpha2.PathMatchRegularExpression || *match.Path.Type == gatewayv1alpha2.PathMatchPathPrefix {
					r.Route.Paths = []*string{match.Path.Value}
				}
			}

			// configure method matching information about the route if method
			// matching was defined.
			if match.Method != nil {
				r.Route.Methods = append(r.Route.Methods, kong.String(string(*match.Method)))
			}

			// convert header matching from HTTPRoute to Route format
			headers, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(match.Headers)
			if err != nil {
				return nil, err
			}
			if len(headers) > 0 {
				r.Route.Headers = headers
			}

			// add the route to the list of routes for the service(s)
			routes = append(routes, r)
		}
	} else {
		// it's acceptable for an HTTPRoute to have no matches in the rulesets,
		// but only backends as long as there are hostnames. In this case, we
		// match all traffic based on the hostname and leave all other routing
		// options default.
		r := kongstate.Route{
			Ingress: objectInfo,
			Route: kong.Route{
				Name:         kong.String(fmt.Sprintf("httproute.%s.%s.0.0", httproute.Namespace, httproute.Name)),
				Protocols:    kong.StringSlice("http", "https"),
				PreserveHost: kong.Bool(true),
			},
		}

		// however in this case there must actually be some present hostnames
		// configured for the HTTPRoute or else it's not valid.
		if len(hostnames) == 0 {
			return nil, fmt.Errorf("no match rules or hostnames specified")
		}

		// otherwise apply the hostnames to the route
		r.Hosts = append(r.Hosts, hostnames...)
		routes = append(routes, r)
	}

	return routes, nil
}

// generateKongServiceFromHTTPRouteBackendRef converts a provided backendRef for an HTTPRoute
// into a kong.Service so that routes for that object can be attached to the Service.
func generateKongServiceFromHTTPRouteBackendRef(
	result *ingressRules,
	httproute *gatewayv1alpha2.HTTPRoute,
	ruleNumber int,
	backendRefs ...gatewayv1alpha2.HTTPBackendRef,
) (kongstate.Service, error) {
	// at least one backendRef must be present
	if len(backendRefs) == 0 {
		return kongstate.Service{}, fmt.Errorf("no backendRefs present for HTTPRoute %s/%s", httproute.Namespace, httproute.Name)
	}

	// create a kongstate backend for each HTTPRoute backendRef
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
	serviceName := fmt.Sprintf("%s.%d", getUniqueKongServiceNameForObject(httproute), ruleNumber)

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
				Protocol:       kong.String("http"),
				Path:           kong.String("/"),
				ConnectTimeout: kong.Int(DefaultServiceTimeout),
				ReadTimeout:    kong.Int(DefaultServiceTimeout),
				WriteTimeout:   kong.Int(DefaultServiceTimeout),
				Retries:        kong.Int(DefaultRetries),
			},
			Namespace: httproute.Namespace,
			Backends:  backends,
		}
	}

	return service, nil
}
