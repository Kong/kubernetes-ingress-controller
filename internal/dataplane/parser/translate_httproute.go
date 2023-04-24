package parser

import (
	"context"
	"fmt"
	pathlib "path"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
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
func (p *Parser) ingressRulesFromHTTPRoutes(ctx context.Context) ingressRules {
	result := newIngressRules()

	httpRouteList, err := p.storer.ListHTTPRoutes(ctx)
	if err != nil {
		p.logger.WithError(err).Error("failed to list HTTPRoutes")
		return result
	}

	for _, httproute := range httpRouteList {
		if err := p.ingressRulesFromHTTPRoute(ctx, &result, httproute); err != nil {
			p.registerTranslationFailure(fmt.Sprintf("HTTPRoute can't be routed: %s", err), httproute)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.ReportKubernetesObjectUpdate(httproute)
		}
	}

	return result
}

func (p *Parser) ingressRulesFromHTTPRoute(ctx context.Context, result *ingressRules, httproute *gatewayv1beta1.HTTPRoute) error {
	if err := validateHTTPRoute(httproute); err != nil {
		return fmt.Errorf("validation failed : %w", err)
	}

	if p.featureEnabledCombinedServiceRoutes {
		return p.ingressRulesFromHTTPRouteWithCombinedServiceRoutes(ctx, httproute, result)
	}

	return p.ingressRulesFromHTTPRouteLegacyFallback(ctx, httproute, result)
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

	return nil
}

// ingressRulesFromHTTPRouteWithCombinedServiceRoutes generates a set of proto-Kong routes (ingress rules) from an HTTPRoute.
// If multiple rules in the HTTPRoute use the same Service, it combines them into a single Kong route.
func (p *Parser) ingressRulesFromHTTPRouteWithCombinedServiceRoutes(ctx context.Context, httproute *gatewayv1beta1.HTTPRoute, result *ingressRules) error {
	for _, kongServiceTranslation := range translators.TranslateHTTPRoute(httproute) {
		// HTTPRoute uses a wrapper HTTPBackendRef to add optional filters to its BackendRefs
		backendRefs := httpBackendRefsToBackendRefs(kongServiceTranslation.BackendRefs)

		serviceName := kongServiceTranslation.Name

		// create a service and attach the routes to it
		service, err := generateKongServiceFromBackendRefWithName(ctx, p.logger, p.storer, result, serviceName, httproute, "http", backendRefs...)
		if err != nil {
			return err
		}

		// generate the routes for the service and attach them to the service
		for _, kongRouteTranslation := range kongServiceTranslation.KongRoutes {
			routes, err := generateKongRouteFromTranslation(httproute, kongRouteTranslation, p.flagEnabledRegexPathPrefix)
			if err != nil {
				return err
			}
			service.Routes = append(service.Routes, routes...)
		}

		// cache the service to avoid duplicates in further loop iterations
		result.ServiceNameToServices[*service.Service.Name] = service
		result.ServiceNameToParent[serviceName] = httproute
	}

	return nil
}

// ingressRulesFromHTTPRouteLegacyFallback generates a set of proto-Kong routes (ingress rules) from an HTTPRoute.
// It generates a separate route for each rule.
// It is planned for deprecation in favor of ingressRulesFromHTTPRouteWithCombinedServiceRoutes.
func (p *Parser) ingressRulesFromHTTPRouteLegacyFallback(ctx context.Context, httproute *gatewayv1beta1.HTTPRoute, result *ingressRules) error {
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
		service, err := generateKongServiceFromBackendRefWithRuleNumber(ctx, p.logger, p.storer, result, httproute, ruleNumber, "http", backendRefs...)
		if err != nil {
			return err
		}
		service.Routes = append(service.Routes, routes...)

		// cache the service to avoid duplicates in further loop iterations
		result.ServiceNameToServices[*service.Service.Name] = service
		result.ServiceNameToParent[*service.Service.Name] = httproute
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
	tags := util.GenerateTagsForObject(httproute)

	// the HTTPRoute specification upstream specifically defines matches as
	// independent (e.g. each match is an OR with other matches, not an AND).
	// Therefore we treat each match rule as a separate Kong Route, so we iterate through
	// all matches to determine all the routes that will be needed for the services.
	var routes []kongstate.Route

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

			r, err := generateKongRoutesFromHTTPRouteMatches(
				routeName,
				rule.Matches[matchNumber:matchNumber+1],
				rule.Filters,
				objectInfo,
				hostnames,
				addRegexPrefix,
				tags,
			)
			if err != nil {
				return nil, err
			}

			// add the route to the list of routes for the service(s)
			routes = append(routes, r...)
		}
	} else {
		routeName := fmt.Sprintf("httproute.%s.%s.0.0", httproute.Namespace, httproute.Name)
		r, err := generateKongRoutesFromHTTPRouteMatches(routeName,
			rule.Matches,
			rule.Filters,
			objectInfo,
			hostnames,
			addRegexPrefix,
			tags)
		if err != nil {
			return nil, err
		}

		// add the route to the list of routes for the service(s)
		routes = append(routes, r...)
	}

	return routes, nil
}

func generateKongRouteFromTranslation(
	httproute *gatewayv1beta1.HTTPRoute,
	translation translators.KongRouteTranslation,
	addRegexPrefix bool,
) ([]kongstate.Route, error) {
	// gather the k8s object information and hostnames from the httproute
	objectInfo := util.FromK8sObject(httproute)
	tags := util.GenerateTagsForObject(httproute)

	// get the hostnames from the HTTPRoute
	hostnames := getHTTPRouteHostnamesAsSliceOfStringPointers(httproute)

	return generateKongRoutesFromHTTPRouteMatches(
		translation.Name,
		translation.Matches,
		translation.Filters,
		objectInfo,
		hostnames,
		addRegexPrefix,
		tags,
	)
}

// generateKongRoutesFromHTTPRouteMatches converts an HTTPRouteMatches to a slice of Kong Route objects.
// This function assumes that the HTTPRouteMatches share the query params, headers and methods.
func generateKongRoutesFromHTTPRouteMatches(
	routeName string,
	matches []gatewayv1beta1.HTTPRouteMatch,
	filters []gatewayv1beta1.HTTPRouteFilter,
	ingressObjectInfo util.K8sObjectInfo,
	hostnames []*string,
	addRegexPrefix bool,
	tags []*string,
) ([]kongstate.Route, error) {
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
				Tags:         tags,
			},
		}

		// however in this case there must actually be some present hostnames
		// configured for the HTTPRoute or else it's not valid.
		if len(hostnames) == 0 {
			return []kongstate.Route{}, errRouteValidationNoMatchRulesOrHostnamesSpecified
		}

		// otherwise apply the hostnames to the route
		r.Hosts = append(r.Hosts, hostnames...)

		return []kongstate.Route{r}, nil
	}

	// TODO: implement query param matches (https://github.com/Kong/kubernetes-ingress-controller/issues/2778)
	if len(matches[0].QueryParams) > 0 {
		return []kongstate.Route{}, errRouteValidationQueryParamMatchesUnsupported
	}

	r := generateKongstateHTTPRoute(routeName, ingressObjectInfo, hostnames)
	r.Tags = tags

	// convert header matching from HTTPRoute to Route format
	headers, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(matches[0].Headers)
	if err != nil {
		return []kongstate.Route{}, err
	}
	if len(headers) > 0 {
		r.Route.Headers = headers
	}

	// stripPath needs to be disabled by default to be conformant with the Gateway API
	r.StripPath = kong.Bool(false)

	_, hasRedirectFilter := lo.Find(filters, func(filter gatewayv1beta1.HTTPRouteFilter) bool {
		return filter.Type == gatewayv1beta1.HTTPRouteFilterRequestRedirect
	})

	routes := getRoutesFromMatches(matches, &r, filters, tags, hasRedirectFilter, addRegexPrefix)

	// if the redirect filter has not been set, we still need to set the route plugins
	if !hasRedirectFilter {
		plugins := generatePluginsFromHTTPRouteFilters(filters, "", tags)
		r.Plugins = append(r.Plugins, plugins...)
		routes = []kongstate.Route{r}
	}

	return routes, nil
}

// getRoutesFromMatches converts all the httpRoute matches to the proper set of kong routes.
func getRoutesFromMatches(matches []gatewayv1beta1.HTTPRouteMatch,
	route *kongstate.Route,
	filters []gatewayv1beta1.HTTPRouteFilter,
	tags []*string,
	hasRedirectFilter bool,
	addRegexPrefix bool,
) []kongstate.Route {
	seenMethods := make(map[string]struct{})
	routes := make([]kongstate.Route, 0)

	for _, match := range matches {
		// if the rule specifies the redirectFilter, we cannot put all the paths under the same route,
		// as the kong plugin needs to know the exact path to use to perform redirection.
		if hasRedirectFilter {
			matchRoute := route
			// configure path matching information about the route if paths matching was defined
			// Kong automatically infers whether or not a path is a regular expression and uses a prefix match by
			// default if it is not. For those types, we use the path value as-is and let Kong determine the type.
			// For exact matches, we transform the path into a regular expression that terminates after the value
			if match.Path != nil {
				paths := generateKongRoutePathFromHTTPRouteMatch(match, addRegexPrefix)
				for _, p := range paths {
					matchRoute.Route.Paths = append(matchRoute.Route.Paths, kong.String(p))
				}
			}

			// configure method matching information about the route if method
			// matching was defined.
			if match.Method != nil {
				method := string(*match.Method)
				if _, ok := seenMethods[method]; !ok {
					matchRoute.Route.Methods = append(matchRoute.Route.Methods, kong.String(string(*match.Method)))
					seenMethods[method] = struct{}{}
				}
			}
			path := ""
			if match.Path.Value != nil {
				path = *match.Path.Value
			}

			// generate kong plugins from rule.filters
			plugins := generatePluginsFromHTTPRouteFilters(filters, path, tags)
			matchRoute.Plugins = append(matchRoute.Plugins, plugins...)

			routes = append(routes, *route)
		} else {
			// configure path matching information about the route if paths matching was defined
			// Kong automatically infers whether or not a path is a regular expression and uses a prefix match by
			// default if it is not. For those types, we use the path value as-is and let Kong determine the type.
			// For exact matches, we transform the path into a regular expression that terminates after the value
			if match.Path != nil {
				paths := generateKongRoutePathFromHTTPRouteMatch(match, addRegexPrefix)
				for _, p := range paths {
					route.Route.Paths = append(route.Route.Paths, kong.String(p))
				}
			}

			if match.Method != nil {
				method := string(*match.Method)
				if _, ok := seenMethods[method]; !ok {
					route.Route.Methods = append(route.Route.Methods, kong.String(string(*match.Method)))
					seenMethods[method] = struct{}{}
				}
			}
		}
	}
	return routes
}

func generateKongRoutePathFromHTTPRouteMatch(match gatewayv1beta1.HTTPRouteMatch, addRegexPrefix bool) []string {
	switch *match.Path.Type {
	case gatewayv1beta1.PathMatchExact:
		terminated := *match.Path.Value + "$"
		if addRegexPrefix {
			terminated = translators.KongPathRegexPrefix + terminated
		}
		return []string{terminated}

	case gatewayv1beta1.PathMatchPathPrefix:
		paths := make([]string, 0, 2)
		path := *match.Path.Value
		if addRegexPrefix {
			paths = append(paths, fmt.Sprintf("%s%s$", translators.KongPathRegexPrefix, path))
			if !strings.HasSuffix(path, "/") {
				path = fmt.Sprintf("%s/", path)
			}
		}
		return append(paths, path)

	case gatewayv1beta1.PathMatchRegularExpression:
		path := *match.Path.Value
		if addRegexPrefix {
			path = translators.KongPathRegexPrefix + path
		}
		return []string{path}
	}

	return []string{""} // unreachable code
}

func generateKongstateHTTPRoute(routeName string, ingressObjectInfo util.K8sObjectInfo, hostnames []*string) kongstate.Route {
	// build the route object using the method and pathing information
	r := kongstate.Route{
		Ingress: ingressObjectInfo,
		Route: kong.Route{
			Name:         kong.String(routeName),
			Protocols:    kong.StringSlice("http", "https"),
			PreserveHost: kong.Bool(true),
			// metadata tags aren't added here, they're added by the caller
		},
	}

	// attach any hostnames associated with the httproute
	if len(hostnames) > 0 {
		r.Hosts = hostnames
	}

	return r
}

// generatePluginsFromHTTPRouteFilters converts HTTPRouteFilter into Kong plugins.
// path is the parameter to be used by the redirect plugin, to perform redirection.
func generatePluginsFromHTTPRouteFilters(filters []gatewayv1beta1.HTTPRouteFilter, path string, tags []*string) []kong.Plugin {
	kongPlugins := make([]kong.Plugin, 0)
	if len(filters) == 0 {
		return kongPlugins
	}

	for _, filter := range filters {
		switch filter.Type {
		case gatewayv1beta1.HTTPRouteFilterRequestHeaderModifier:
			kongPlugins = append(kongPlugins, generateRequestHeaderModifierKongPlugin(filter.RequestHeaderModifier))

		case gatewayv1beta1.HTTPRouteFilterRequestRedirect:
			kongPlugins = append(kongPlugins, generateRequestRedirectKongPlugin(filter.RequestRedirect, path)...)

		case gatewayv1beta1.HTTPRouteFilterExtensionRef,
			gatewayv1beta1.HTTPRouteFilterRequestMirror,
			gatewayv1beta1.HTTPRouteFilterResponseHeaderModifier,
			gatewayv1beta1.HTTPRouteFilterURLRewrite:
			// not supported
		}
	}
	for _, p := range kongPlugins {
		// This plugin is derived from an HTTPRoute filter, not a KongPlugin, so we apply tags indicating that
		// HTTPRoute as the parent Kubernetes resource for these generated plugins.
		p.Tags = tags
	}

	return kongPlugins
}

func generateRequestRedirectKongPlugin(modifier *gatewayv1beta1.HTTPRequestRedirectFilter, path string) []kong.Plugin {
	plugins := make([]kong.Plugin, 2)
	plugins[0] = kong.Plugin{
		Name: kong.String("request-termination"),
		Config: kong.Configuration{
			"status_code": modifier.StatusCode,
		},
	}

	var locationHeader string
	scheme := "http"
	port := 80

	if modifier.Scheme != nil {
		scheme = *modifier.Scheme
	}
	if modifier.Port != nil {
		port = int(*modifier.Port)
	}
	if modifier.Path != nil && modifier.Path.Type == gatewayv1beta1.FullPathHTTPPathModifier && modifier.Path.ReplaceFullPath != nil {
		// only ReplaceFullPath currently supported
		path = *modifier.Path.ReplaceFullPath
	}
	if modifier.Hostname != nil {
		locationHeader = fmt.Sprintf("Location: %s://%s", scheme, pathlib.Join(fmt.Sprintf("%s:%d", *modifier.Hostname, port), path))
	} else {
		locationHeader = fmt.Sprintf("Location: %s", path)
	}

	plugins[1] = kong.Plugin{
		Name: kong.String("response-transformer"),
		Config: kong.Configuration{
			"add": map[string][]string{
				"headers": {locationHeader},
			},
		},
	}

	return plugins
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
