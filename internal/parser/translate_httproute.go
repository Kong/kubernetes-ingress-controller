package parser

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/sirupsen/logrus"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// fromHTTPRoutes processes all the HTTPRoute objects present in the cache and translates
// them to Kong Gateway configurations.
func fromHTTPRoutes(log logrus.FieldLogger, httpRouteList []*gatewayv1alpha2.HTTPRoute) ingressRules {
	result := newIngressRules()

	for _, httproute := range httpRouteList {
		// first we grab the spec and gather some metdata about the object
		objectInfo := util.FromK8sObject(httproute)
		spec := httproute.Spec

		// gather the hostnames that will be used (globally) for route matching
		hostnames := make([]*string, 0, len(spec.Hostnames))
		for _, hostname := range spec.Hostnames {
			hostnames = append(hostnames, kong.String(string(hostname)))
		}

		// validation for HTTPRoutes will happen at a higher layer, but in spite of that we run
		// validation at this level as well as a fallback so that if routes are posted which
		// are invalid somehow make it past validation (e.g. the webhook is not enabled) we can
		// at least try to provide a helpful message about the situation in the manager logs.
		if len(spec.Rules) < 1 {
			log.Errorf("HTTPRoute %s/%s can't be routed: no rules provided", httproute.Namespace, httproute.Name)
			continue
		}

		// each rule may represent a different set of backend services that will be accepting
		// traffic, so we make separate routes and Kong services for every present rule.
		for _, rule := range spec.Rules {
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
						"httproute.%s.%s.%d",
						httproute.Namespace,
						httproute.Name,
						matchNumber, // TODO: avoid route thrash from re-ordering?
					))

					// TODO: implement query param matches
					if len(match.QueryParams) > 0 {
						errmsg := "query param matches are not yet supported"
						log.Errorf("HTTPRoute %s/%s can't be routed for match %+v: %s", errmsg)
						continue
					}

					// TODO: implement regex path matches
					if *match.Path.Type == gatewayv1alpha2.PathMatchRegularExpression {
						errmsg := "regular expression path matches are not yet supported"
						log.Errorf("HTTPRoute %s/%s can't be routed for match %+v: %s", errmsg)
						continue
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

					// configure path matching information about the route if paths
					// matching was defined.
					if match.Path != nil {
						// determine the path match values
						r.Route.Paths = []*string{match.Path.Value}

						// determine whether path stripping needs to be enabled
						r.Route.StripPath = kong.Bool(match.Path.Type == nil || *match.Path.Type == gatewayv1alpha2.PathMatchPathPrefix)
					}

					// configure method matching information about the route if method
					// matching was defined.
					if match.Method != nil {
						r.Route.Methods = append(r.Route.Methods, kong.String(string(*match.Method)))
					}

					// convert header matching from HTTPRoute to Route format
					headers, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(match.Headers)
					if err != nil {
						log.Errorf("HTTPRoute %s/%s can't be routed for match %+v: %w", err)
						continue
					}
					if len(headers) > 0 {
						r.Route.Headers = headers
					}

					// add the route to the list of routes for the service(s)
					routes = append(routes, r)
					log.Debugf("generated route %s for HTTPRoute %s/%s", routeName, httproute.Namespace, httproute.Name)
				}
			} else {
				// it's acceptable for an HTTPRoute to have no matches in the rulesets,
				// but only backends as long as there are hostnames. In this case, we
				// match all traffic based on the hostname and leave all other routing
				// options default.
				r := kongstate.Route{
					Ingress: objectInfo,
					Route: kong.Route{
						Name:         kong.String(fmt.Sprintf("httproute.%s.%s.0", httproute.Namespace, httproute.Name)),
						Protocols:    kong.StringSlice("http", "https"),
						PreserveHost: kong.Bool(true),
					},
				}

				// however in this case there must actually be some present hostnames
				// configured for the HTTPRoute or else it's not valid.
				if len(hostnames) == 0 {
					log.Errorf("HTTPRoute %s/%s can't be routed: no match rules or hostnames specified", httproute.Namespace, httproute.Name)
					continue
				}

				// otherwise apply the hostnames to the route
				r.Hosts = append(r.Hosts, hostnames...)
				routes = append(routes, r)
			}

			// if there are no routes (this can happen if only unsupported match options are
			// provided) then we'll avoid doing anything with the backendRefs.
			if len(routes) == 0 {
				log.Errorf("HTTPRoute %s/%s can't be routed: no valid rules provided", httproute.Namespace, httproute.Name)
				continue
			}
			log.Debugf("HTTPRoute %s/%s required %d kong routes to configure", httproute.Namespace, httproute.Name, len(routes))

			// once all routes have been determined based on matching rules
			// we determine the Services they actually route to.
			for _, backendRef := range rule.BackendRefs {
				// determine the namespace for the service, or default to the same namespace
				// as the HTTPRoute object.
				//
				// TODO: need to add validation to restrict namespaces in backendRefs
				namespace := httproute.Namespace
				if backendRef.Namespace != nil {
					namespace = string(*backendRef.Namespace)
				}

				// determine the name of the Service
				serviceName := fmt.Sprintf("%s.%s.%d", namespace, backendRef.Name, *backendRef.Port)

				// determine the Service port
				port := kongstate.PortDef{
					Mode:   kongstate.PortModeByNumber,
					Number: int32(*backendRef.Port),
				}

				// check if the service is already known, and if not create it
				service, ok := result.ServiceNameToServices[serviceName]
				if !ok {
					service = kongstate.Service{
						Service: kong.Service{
							Name:           kong.String(serviceName),
							Host:           kong.String(fmt.Sprintf("%s.%s.%s.svc", backendRef.Name, namespace, port.CanonicalString())),
							Port:           kong.Int(int(*backendRef.Port)),
							Protocol:       kong.String("http"),
							Path:           kong.String("/"),
							ConnectTimeout: kong.Int(60000),
							ReadTimeout:    kong.Int(60000),
							WriteTimeout:   kong.Int(60000),
							Retries:        kong.Int(5),
						},
						Namespace: httproute.Namespace,
						Backend: kongstate.ServiceBackend{
							Name: string(backendRef.Name),
							Port: port,
						},
					}
					log.Debugf("generated kong service %s for HTTPRoute %s/%s", serviceName, httproute.Namespace, httproute.Name)
				}

				// add all generated routes to this service
				service.Routes = append(service.Routes, routes...)

				// cache the service to avoid duplicates in further loop iterations
				result.ServiceNameToServices[serviceName] = service
			}
		}
	}

	return result
}
