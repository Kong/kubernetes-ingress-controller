package parser

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
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

	for _, httproute := range httpRouteList {
		if err := p.ingressRulesFromHTTPRoute(&result, httproute); err != nil {
			p.registerTranslationFailure(fmt.Sprintf("HTTPRoute can't be routed: %s", err), httproute)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.ReportKubernetesObjectUpdate(httproute)
		}
	}

	return result
}

func (p *Parser) ingressRulesFromHTTPRoute(result *ingressRules, httproute *gatewayv1beta1.HTTPRoute) error {
	if err := validateHTTPRoute(httproute); err != nil {
		return fmt.Errorf("validation failed : %w", err)
	}

	if p.featureEnabledCombinedServiceRoutes {
		return p.ingressRulesFromHTTPRouteWithCombinedServiceRoutes(httproute, result)
	}

	return p.ingressRulesFromHTTPRouteLegacyFallback(httproute, result)
}

func validateHTTPRoute(httproute *gatewayv1beta1.HTTPRoute) error {
	spec := httproute.Spec

	// validation for HTTPRoutes will happen at a higher layer, but in spite of that we run
	// validation at this level as well as a fallback so that if routes are posted which
	// are invalid somehow make it past validation (e.g. the webhook is not enabled) we can
	// at least try to provide a helpful message about the situation in the manager logs.
	if len(spec.Rules) == 0 {
		return errRouteValidationNoRules
	}

	for _, rule := range spec.Rules {
		if len(rule.BackendRefs) == 0 {
			return errRouteValidationMissingBackendRefs
		}
	}
	return nil
}

// ingressRulesFromHTTPRouteWithCombinedServiceRoutes generates a set of proto-Kong routes (ingress rules) from an HTTPRoute.
// If multiple rules in the HTTPRoute use the same Service, it combines them into a single Kong route.
func (p *Parser) ingressRulesFromHTTPRouteWithCombinedServiceRoutes(httproute *gatewayv1beta1.HTTPRoute, result *ingressRules) error {
	for _, kongServiceTranslation := range translators.TranslateHTTPRoute(httproute) {
		// HTTPRoute uses a wrapper HTTPBackendRef to add optional filters to its BackendRefs
		backendRefs := httpBackendRefsToBackendRefs(kongServiceTranslation.BackendRefs)

		serviceName := kongServiceTranslation.Name

		// create a service and attach the routes to it
		service, err := generateKongServiceFromBackendRefWithName(p.logger, p.storer, result, serviceName, httproute, "http", backendRefs...)
		if err != nil {
			return err
		}

		// generate the routes for the service and attach them to the service
		for _, kongRouteTranslation := range kongServiceTranslation.KongRoutes {
			route, err := generateKongRouteFromTranslation(httproute, kongRouteTranslation, p.flagEnabledRegexPathPrefix)
			if err != nil {
				return err
			}
			service.Routes = append(service.Routes, route)
		}

		// cache the service to avoid duplicates in further loop iterations
		result.ServiceNameToServices[*service.Service.Name] = service
	}

	return nil
}

// ingressRulesFromHTTPRouteLegacyFallback generates a set of proto-Kong routes (ingress rules) from an HTTPRoute.
// It generates a separate route for each rule.
// It is planned for deprecation in favor of ingressRulesFromHTTPRouteWithCombinedServiceRoutes.
func (p *Parser) ingressRulesFromHTTPRouteLegacyFallback(httproute *gatewayv1beta1.HTTPRoute, result *ingressRules) error {
	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range httproute.Spec.Rules {
		// determine the routes needed to route traffic to services for this rule
		routes, err := generateKongRoutesFromHTTPRouteRule(httproute, ruleNumber, rule, p.flagEnabledRegexPathPrefix)
		if err != nil {
			return err
		}

		// HTTPRoute uses a wrapper HTTPBackendRef to add optional filters to its BackendRefs
		backendRefs := httpBackendRefsToBackendRefs(rule.BackendRefs)

		// create a service and attach the routes to it
		service, err := generateKongServiceFromBackendRefWithRuleNumber(p.logger, p.storer, result, httproute, ruleNumber, "http", backendRefs...)
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
func getHTTPRouteHostnamesAsSliceOfStringPointers(httproute *gatewayv1beta1.HTTPRoute) []*string {
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
	httproute *gatewayv1beta1.HTTPRoute,
	ruleNumber int,
	rule gatewayv1beta1.HTTPRouteRule,
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
	plugins := generatePluginsFromHTTPRouteFilters(rule.Filters)

	if len(rule.Matches) > 0 {
		for matchNumber := range rule.Matches {
			// determine the name of the route, identify it as a route that belongs
			// to a Kubernetes HTTPRoute object.
			routeName := fmt.Sprintf(
				"httproute.%s.%s.%d.%d",
				httproute.Namespace,
				httproute.Name,
				ruleNumber,
				matchNumber,
			)

			r, err := generateKongRouteFromHTTPRouteMatches(
				routeName,
				rule.Matches[matchNumber:matchNumber+1],
				objectInfo,
				hostnames,
				plugins,
				addRegexPrefix,
			)
			if err != nil {
				return nil, err
			}

			// add the route to the list of routes for the service(s)
			routes = append(routes, r)
		}
	} else {
		routeName := fmt.Sprintf("httproute.%s.%s.0.0", httproute.Namespace, httproute.Name)
		r, err := generateKongRouteFromHTTPRouteMatches(routeName, rule.Matches, objectInfo, hostnames, plugins, addRegexPrefix)
		if err != nil {
			return nil, err
		}

		// add the route to the list of routes for the service(s)
		routes = append(routes, r)
	}

	return routes, nil
}

func generateKongRouteFromTranslation(
	httproute *gatewayv1beta1.HTTPRoute,
	translation translators.KongRouteTranslation,
	addRegexPrefix bool,
) (kongstate.Route, error) {
	// gather the k8s object information and hostnames from the httproute
	objectInfo := util.FromK8sObject(httproute)

	// get the hostnames from the HTTPRoute
	hostnames := getHTTPRouteHostnamesAsSliceOfStringPointers(httproute)

	// generate kong plugins from rule.filters
	plugins := generatePluginsFromHTTPRouteFilters(translation.Filters)

	return generateKongRouteFromHTTPRouteMatches(
		translation.Name,
		translation.Matches,
		objectInfo,
		hostnames,
		plugins,
		addRegexPrefix,
	)
}

// generateKongRouteFromHTTPRouteMatches converts an HTTPRouteMatches to a Kong Route object.
// This function assumes that the HTTPRouteMatches share the query params, headers and methods.
func generateKongRouteFromHTTPRouteMatches(
	routeName string,
	matches []gatewayv1beta1.HTTPRouteMatch,
	ingressObjectInfo util.K8sObjectInfo,
	hostnames []*string,
	plugins []kong.Plugin,
	addRegexPrefix bool,
) (kongstate.Route, error) {
	if len(matches) == 0 {
		// it's acceptable for an HTTPRoute to have no matches in the rulesets,
		// but only backends as long as there are hostnames. In this case, we
		// match all traffic based on the hostname and leave all other routing
		// options default.
		r := kongstate.Route{
			Ingress: ingressObjectInfo,
			Route: kong.Route{
				Name:         kong.String(routeName),
				Protocols:    kong.StringSlice("http", "https"),
				PreserveHost: kong.Bool(true),
			},
		}

		// however in this case there must actually be some present hostnames
		// configured for the HTTPRoute or else it's not valid.
		if len(hostnames) == 0 {
			return kongstate.Route{}, errRouteValidationNoMatchRulesOrHostnamesSpecified
		}

		// otherwise apply the hostnames to the route
		r.Hosts = append(r.Hosts, hostnames...)

		// attach the plugins to be applied to the given route
		r.Plugins = append(r.Plugins, plugins...)

		return r, nil
	}

	// TODO: implement query param matches (https://github.com/Kong/kubernetes-ingress-controller/issues/2778)
	if len(matches[0].QueryParams) > 0 {
		return kongstate.Route{}, errRouteValidationQueryParamMatchesUnsupported
	}

	r := generateKongstateRoute(routeName, ingressObjectInfo, hostnames)

	// convert header matching from HTTPRoute to Route format
	headers, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(matches[0].Headers)
	if err != nil {
		return kongstate.Route{}, err
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

	seenMethods := make(map[string]struct{})

	for _, match := range matches {
		// configure path matching information about the route if paths matching was defined
		// Kong automatically infers whether or not a path is a regular expression and uses a prefix match by
		// default if it is not. For those types, we use the path value as-is and let Kong determine the type.
		// For exact matches, we transform the path into a regular expression that terminates after the value
		if match.Path != nil {
			path := generateKongRoutePathFromHTTPRouteMatch(match, addRegexPrefix)
			r.Route.Paths = append(r.Route.Paths, kong.String(path))
		}

		// configure method matching information about the route if method
		// matching was defined.
		if match.Method != nil {
			method := string(*match.Method)
			if _, ok := seenMethods[method]; !ok {
				r.Route.Methods = append(r.Route.Methods, kong.String(string(*match.Method)))
				seenMethods[method] = struct{}{}
			}
		}
	}

	return r, nil
}

func generateKongRoutePathFromHTTPRouteMatch(match gatewayv1beta1.HTTPRouteMatch, addRegexPrefix bool) string {
	switch *match.Path.Type {
	case gatewayv1beta1.PathMatchExact:
		terminated := *match.Path.Value + "$"
		if addRegexPrefix {
			terminated = translators.KongPathRegexPrefix + terminated
		}
		return terminated

	case gatewayv1beta1.PathMatchPathPrefix:
		path := *match.Path.Value
		return path

	case gatewayv1beta1.PathMatchRegularExpression:
		path := *match.Path.Value
		if addRegexPrefix {
			path = translators.KongPathRegexPrefix + path
		}
		return path
	}

	return "" // unreachable code
}

func generateKongstateRoute(routeName string, ingressObjectInfo util.K8sObjectInfo, hostnames []*string) kongstate.Route {
	// build the route object using the method and pathing information
	r := kongstate.Route{
		Ingress: ingressObjectInfo,
		Route: kong.Route{
			Name:         kong.String(routeName),
			Protocols:    kong.StringSlice("http", "https"),
			PreserveHost: kong.Bool(true),
		},
	}

	// attach any hostnames associated with the httproute
	if len(hostnames) > 0 {
		r.Hosts = hostnames
	}

	return r
}

// generatePluginsFromHTTPRouteFilters  converts HTTPRouteFilter into Kong filters.
func generatePluginsFromHTTPRouteFilters(filters []gatewayv1beta1.HTTPRouteFilter) []kong.Plugin {
	kongPlugins := make([]kong.Plugin, 0)
	if len(filters) == 0 {
		return kongPlugins
	}

	for _, filter := range filters {
		if filter.Type == gatewayv1beta1.HTTPRouteFilterRequestHeaderModifier {
			kongPlugins = append(kongPlugins, generateRequestHeaderModifierKongPlugin(filter.RequestHeaderModifier))
		}
		// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/2793
	}

	return kongPlugins
}

// generateRequestHeaderModifierKongPlugin converts a gatewayv1beta1.HTTPRequestHeaderFilter into a
// kong.Plugin of type request-transformer.
func generateRequestHeaderModifierKongPlugin(modifier *gatewayv1beta1.HTTPHeaderFilter) kong.Plugin {
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

func kongHeaderFormatter(header gatewayv1beta1.HTTPHeader) string {
	return fmt.Sprintf("%s:%s", header.Name, header.Value)
}

func httpBackendRefsToBackendRefs(httpBackendRef []gatewayv1beta1.HTTPBackendRef) []gatewayv1beta1.BackendRef {
	backendRefs := make([]gatewayv1beta1.BackendRef, 0, len(httpBackendRef))

	for _, hRef := range httpBackendRef {
		backendRefs = append(backendRefs, hRef.BackendRef)
	}
	return backendRefs
}
