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
		if err := p.ingressRulesFromHTTPRoute(&result, httproute); err != nil {
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

func (p *Parser) ingressRulesFromHTTPRoute(result *ingressRules, httproute *gatewayv1alpha2.HTTPRoute) error {
	// first we grab the spec and gather some metadata about the object
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
		routes, err := generateKongRoutesFromHTTPRouteRule(httproute, ruleNumber, rule, p.flagEnabledRegexPathPrefix)
		if err != nil {
			return err
		}

		// create a service and attach the routes to it
		var backendRefs []gatewayv1alpha2.BackendRef
		// HTTPRoute uses a wrapper HTTPBackendRef to add optional filters to its BackendRefs
		for _, hRef := range rule.BackendRefs {
			backendRefs = append(backendRefs, hRef.BackendRef)
		}
		service, err := p.generateKongServiceFromBackendRef(result, httproute, ruleNumber, "http", backendRefs...)
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
func generateKongRoutesFromHTTPRouteRule(
	httproute *gatewayv1alpha2.HTTPRoute,
	ruleNumber int,
	rule gatewayv1alpha2.HTTPRouteRule,
	addRegexPrefix bool,
) ([]kongstate.Route, error) {
	// gather the k8s object information and hostnames from the httproute
	objectInfo := util.FromK8sObject(httproute)
	hostnames := getHTTPRouteHostnamesAsSliceOfStringPointers(httproute)

	// the HTTPRoute specification upstream specifically defines matches as
	// independent (e.g. each match is an OR with other matches, not an AND).
	// Therefore we treat each match rule as a separate Kong Route, so we iterate through
	// all matches to determine all the routes that will be needed for the services.
	var routes []kongstate.Route

	// generate kong plugins from rule.filters
	plugins := generatePluginsFromHTTPRouteRuleFilters(rule)
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

			// TODO: implement query param matches (https://github.com/Kong/kubernetes-ingress-controller/issues/2778)
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
					if addRegexPrefix {
						terminated = kongPathRegexPrefix + terminated
					}
					r.Route.Paths = []*string{&terminated}
				} else if *match.Path.Type == gatewayv1alpha2.PathMatchRegularExpression || *match.Path.Type == gatewayv1alpha2.PathMatchPathPrefix {
					path := *match.Path.Value
					if addRegexPrefix {
						path = kongPathRegexPrefix + path
					}
					r.Route.Paths = []*string{&path}
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

			// stripPath needs to be disabled by default to be conformant with the Gateway API
			r.StripPath = kong.Bool(false)

			// attach the plugins to be applied to the given route
			if len(plugins) != 0 {
				if r.Plugins == nil {
					r.Plugins = make([]kong.Plugin, 0, len(plugins))
				}
				r.Plugins = append(r.Plugins, plugins...)
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
		// attach the plugins to be applied to the given route
		if len(plugins) != 0 {
			if r.Plugins == nil {
				r.Plugins = make([]kong.Plugin, 0, len(plugins))
			}
			r.Plugins = append(r.Plugins, plugins...)
		}
		routes = append(routes, r)
	}

	return routes, nil
}

// generatePluginsFromHTTPRouteRuleFilters accepts a rule as argument and converts
// HttpRouteRule.Filters into Kong filters.
func generatePluginsFromHTTPRouteRuleFilters(rule gatewayv1alpha2.HTTPRouteRule) []kong.Plugin {
	kongPlugins := make([]kong.Plugin, 0)
	if rule.Filters == nil {
		return kongPlugins
	}

	for _, filter := range rule.Filters {
		if filter.Type == gatewayv1alpha2.HTTPRouteFilterRequestHeaderModifier {
			kongPlugins = append(kongPlugins, generateRequestHeaderModifierKongPlugin(filter.RequestHeaderModifier))
		}
		// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/2793
	}

	return kongPlugins
}

// generateRequestHeaderModifierKongPlugin converts a gatewayv1alpha2.HTTPRequestHeaderFilter into a
// kong.Plugin of type request-transformer.
func generateRequestHeaderModifierKongPlugin(modifier *gatewayv1alpha2.HTTPRequestHeaderFilter) kong.Plugin {
	plugin := kong.Plugin{
		Name:   kong.String("request-transformer"),
		Config: make(kong.Configuration),
	}

	// modifier.Set is converted to a pair composed of "replace" and "add"
	if modifier.Set != nil {
		setModifiers := make([]string, 0, len(modifier.Set))
		for _, s := range modifier.Set {
			setModifiers = append(setModifiers, kongHeaderFormatter(s))
		}
		plugin.Config["replace"] = map[string][]string{
			"headers": setModifiers,
		}
		plugin.Config["add"] = map[string][]string{
			"headers": setModifiers,
		}
	}

	// modifier.Add is converted to "append"
	if modifier.Add != nil {
		appendModifiers := make([]string, 0, len(modifier.Add))
		for _, a := range modifier.Add {
			appendModifiers = append(appendModifiers, kongHeaderFormatter(a))
		}
		plugin.Config["append"] = map[string][]string{
			"headers": appendModifiers,
		}
	}

	if modifier.Remove != nil {
		plugin.Config["remove"] = map[string][]string{
			"headers": modifier.Remove,
		}
	}

	return plugin
}

func kongHeaderFormatter(header gatewayv1alpha2.HTTPHeader) string {
	return fmt.Sprintf("%s:%s", header.Name, header.Value)
}
