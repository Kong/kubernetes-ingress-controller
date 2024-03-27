package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/bombsimon/logrusr/v4"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	k8stypes "k8s.io/apimachinery/pkg/types"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

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

	if p.featureFlags.ExpressionRoutes {
		p.ingressRulesFromHTTPRoutesUsingExpressionRoutes(httpRouteList, &result)
		return result
	}

	for _, httproute := range httpRouteList {
		if err := p.ingressRulesFromHTTPRoute(&result, httproute); err != nil {
			p.registerTranslationFailure(fmt.Sprintf("HTTPRoute can't be routed: %s", err), httproute)
		} else {
			// at this point the object has been configured and can be
			// reported as successfully parsed.
			p.registerSuccessfullyParsedObject(httproute)
		}
	}

	return result
}

func (p *Parser) ingressRulesFromHTTPRoute(result *ingressRules, httproute *gatewayv1.HTTPRoute) error {
	if err := validateHTTPRoute(httproute); err != nil {
		return fmt.Errorf("validation failed : %w", err)
	}

	if p.featureFlags.CombinedServiceRoutes {
		return p.ingressRulesFromHTTPRouteWithCombinedServiceRoutes(httproute, result)
	}

	return p.ingressRulesFromHTTPRouteLegacyFallback(httproute, result)
}

func validateHTTPRoute(httproute *gatewayv1.HTTPRoute) error {
	spec := httproute.Spec

	// validation for HTTPRoutes will happen at a higher layer, but in spite of that we run
	// validation at this level as well as a fallback so that if routes are posted which
	// are invalid somehow make it past validation (e.g. the webhook is not enabled) we can
	// at least try to provide a helpful message about the situation in the manager logs.
	if len(spec.Rules) == 0 {
		return translators.ErrRouteValidationNoRules
	}

	return nil
}

// ingressRulesFromHTTPRoutesUsingExpressionRoutes translates HTTPRoutes to expression based routes
// when ExpressionRoutes feature flag is enabled.
// Because we need to assign different priorities based on the hostname and match in the specification of HTTPRoutes,
// We need to split the HTTPRoutes into ones with only one hostname and one match, then assign priority to them
// and finally translate the split HTTPRoutes into Kong services and routes with assigned priorities.
func (p *Parser) ingressRulesFromHTTPRoutesUsingExpressionRoutes(httpRoutes []*gatewayv1.HTTPRoute, result *ingressRules) {
	// first, split HTTPRoutes by hostnames and matches.
	splitHTTPRouteMatches := []translators.SplitHTTPRouteMatch{}
	for _, httproute := range httpRoutes {
		if err := validateHTTPRoute(httproute); err != nil {
			p.registerTranslationFailure(fmt.Sprintf("HTTPRoute can't be routed: %s", err), httproute)
			continue
		}
		splitHTTPRouteMatches = append(splitHTTPRouteMatches, translators.SplitHTTPRoute(httproute)...)
	}
	// assign priorities to split HTTPRoutes.
	splitHTTPRoutesWithPriorities := translators.AssignRoutePriorityToSplitHTTPRouteMatches(logrusr.New(p.logger), splitHTTPRouteMatches)
	httpRouteNameToTranslationFailure := map[k8stypes.NamespacedName][]error{}

	// translate split HTTPRoute matches to ingress rules, including services, routes, upstreams.
	for _, httpRouteWithPriority := range splitHTTPRoutesWithPriorities {
		err := p.ingressRulesFromSplitHTTPRouteMatchWithPriority(result, httpRouteWithPriority)
		if err != nil {
			nsName := k8stypes.NamespacedName{
				Namespace: httpRouteWithPriority.Match.Source.Namespace,
				Name:      httpRouteWithPriority.Match.Source.Name,
			}
			httpRouteNameToTranslationFailure[nsName] = append(httpRouteNameToTranslationFailure[nsName], err)
		}
	}
	// Register successful parsed objects and  translation failures.
	// Because one HTTPRoute may be split into multiple HTTPRoutes, we need to de-duplicate by namespace and name.
	for _, httproute := range httpRoutes {
		nsName := k8stypes.NamespacedName{
			Namespace: httproute.Namespace,
			Name:      httproute.Name,
		}
		if translationFailures, ok := httpRouteNameToTranslationFailure[nsName]; ok {
			p.registerTranslationFailure(
				fmt.Sprintf("HTTPRoute can't be routed: %v", errors.Join(translationFailures...)),
				httproute,
			)
			continue
		}
		p.registerSuccessfullyParsedObject(httproute)
	}
}

// ingressRulesFromHTTPRouteWithCombinedServiceRoutes generates a set of proto-Kong routes (ingress rules) from an HTTPRoute.
// If multiple rules in the HTTPRoute use the same Service, it combines them into a single Kong route.
func (p *Parser) ingressRulesFromHTTPRouteWithCombinedServiceRoutes(httproute *gatewayv1.HTTPRoute, result *ingressRules) error {
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
			routes, err := GenerateKongRouteFromTranslation(httproute, kongRouteTranslation, p.featureFlags.RegexPathPrefix, p.featureFlags.ExpressionRoutes, p.kongVersion)
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
func (p *Parser) ingressRulesFromHTTPRouteLegacyFallback(httproute *gatewayv1.HTTPRoute, result *ingressRules) error {
	// each rule may represent a different set of backend services that will be accepting
	// traffic, so we make separate routes and Kong services for every present rule.
	for ruleNumber, rule := range httproute.Spec.Rules {
		// determine the routes needed to route traffic to services for this rule
		routes, err := generateKongRoutesFromHTTPRouteRule(httproute, ruleNumber, rule, p.featureFlags.RegexPathPrefix, p.kongVersion)
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
		result.ServiceNameToParent[*service.Service.Name] = httproute
	}
	return nil
}

// -----------------------------------------------------------------------------
// Translate HTTPRoute - Utils
// -----------------------------------------------------------------------------

// getHTTPRouteHostnamesAsSliceOfStrings translates the hostnames defined in an
// HTTPRoute specification into a []*string slice, which is the type required by translating to matchers
// in expression based routes.
func getHTTPRouteHostnamesAsSliceOfStrings(httproute *gatewayv1.HTTPRoute) []string {
	return lo.Map(httproute.Spec.Hostnames, func(h gatewayv1.Hostname, _ int) string {
		return string(h)
	})
}

// getHTTPRouteHostnamesAsSliceOfStringPointers translates the hostnames defined
// in an HTTPRoute specification into a []*string slice, which is the type required
// by kong.Route{}.
func getHTTPRouteHostnamesAsSliceOfStringPointers(httproute *gatewayv1.HTTPRoute) []*string {
	return lo.Map(httproute.Spec.Hostnames, func(h gatewayv1.Hostname, _ int) *string {
		return kong.String(string(h))
	})
}

// generateKongRoutesFromHTTPRouteRule converts an HTTPRoute rule to one or more
// Kong Route objects to route traffic to services. This function will accept an
// HTTPRoute that does not include any matches as long as it includes hostnames
// to route traffic to the backend service with, in this case it will use the default
// path prefix routing option for that service in addition to hostname routing.
// If an HTTPRoute is provided that has matches that include any unsupported matching
// configurations, this will produce an error and the route is considered invalid.
func generateKongRoutesFromHTTPRouteRule(
	httproute *gatewayv1.HTTPRoute,
	ruleNumber int,
	rule gatewayv1.HTTPRouteRule,
	addRegexPrefix bool,
	kongVersion semver.Version,
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
				kongVersion,
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
			tags,
			kongVersion,
		)
		if err != nil {
			return nil, err
		}

		// add the route to the list of routes for the service(s)
		routes = append(routes, r...)
	}

	return routes, nil
}

// GenerateKongRouteFromTranslation generates Kong routes from HTTPRoute
// pointing to a specific backend. It is used for both traditional and expression based routes.
func GenerateKongRouteFromTranslation(
	httproute *gatewayv1.HTTPRoute,
	translation translators.KongRouteTranslation,
	addRegexPrefix bool,
	expressionRoutes bool,
	kongVersion semver.Version,
) ([]kongstate.Route, error) {
	// gather the k8s object information and hostnames from the httproute
	objectInfo := util.FromK8sObject(httproute)
	tags := util.GenerateTagsForObject(httproute)

	// translate to expression based routes when expressionRoutes is enabled.
	if expressionRoutes {
		// get the hostnames from the HTTPRoute
		hostnames := getHTTPRouteHostnamesAsSliceOfStrings(httproute)
		return translators.GenerateKongExpressionRoutesFromHTTPRouteMatches(
			translation,
			objectInfo,
			hostnames,
			tags,
		)
	}

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
		kongVersion,
	)
}

// generateKongRoutesFromHTTPRouteMatches converts an HTTPRouteMatches to a slice of Kong Route objects.
// This function assumes that the HTTPRouteMatches share the query params, headers and methods.
func generateKongRoutesFromHTTPRouteMatches(
	routeName string,
	matches []gatewayv1.HTTPRouteMatch,
	filters []gatewayv1.HTTPRouteFilter,
	ingressObjectInfo util.K8sObjectInfo,
	hostnames []*string,
	addRegexPrefix bool,
	tags []*string,
	kongVersion semver.Version,
) ([]kongstate.Route, error) {
	if len(matches) == 0 {
		// it's acceptable for an HTTPRoute to have no matches in the rulesets,
		// but only backends as long as there are hostnames. In this case, we
		// match all traffic based on the hostname and leave all other routing
		// options default.
		// for rules with no hostnames, we generate a "catch-all" route for it.
		r := kongstate.Route{
			Ingress: ingressObjectInfo,
			Route: kong.Route{
				Name:         kong.String(routeName),
				Protocols:    kong.StringSlice("http", "https"),
				PreserveHost: kong.Bool(true),
				Tags:         tags,
			},
		}
		r.Hosts = append(r.Hosts, hostnames...)

		return []kongstate.Route{r}, nil
	}

	// TODO: implement query param matches (https://github.com/Kong/kubernetes-ingress-controller/issues/2778)
	if len(matches[0].QueryParams) > 0 {
		return []kongstate.Route{}, translators.ErrRouteValidationQueryParamMatchesUnsupported
	}

	r := generateKongstateHTTPRoute(routeName, ingressObjectInfo, hostnames)
	r.Tags = tags

	// convert header matching from HTTPRoute to Route format
	headers, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(matches[0].Headers, kongVersion)
	if err != nil {
		return []kongstate.Route{}, err
	}
	if len(headers) > 0 {
		r.Route.Headers = headers
	}

	// stripPath needs to be disabled by default to be conformant with the Gateway API
	r.StripPath = kong.Bool(false)

	_, hasRedirectFilter := lo.Find(filters, func(filter gatewayv1.HTTPRouteFilter) bool {
		return filter.Type == gatewayv1.HTTPRouteFilterRequestRedirect
	})

	routes := getRoutesFromMatches(matches, &r, filters, tags, hasRedirectFilter, addRegexPrefix)

	// if the redirect filter has not been set, we still need to set the route plugins
	if !hasRedirectFilter {
		plugins := translators.GeneratePluginsFromHTTPRouteFilters(filters, "", tags)
		r.Plugins = append(r.Plugins, plugins...)
		routes = []kongstate.Route{r}
	}

	return routes, nil
}

// getRoutesFromMatches converts all the httpRoute matches to the proper set of kong routes.
func getRoutesFromMatches(matches []gatewayv1.HTTPRouteMatch,
	route *kongstate.Route,
	filters []gatewayv1.HTTPRouteFilter,
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
			plugins := translators.GeneratePluginsFromHTTPRouteFilters(filters, path, tags)
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

func generateKongRoutePathFromHTTPRouteMatch(match gatewayv1.HTTPRouteMatch, addRegexPrefix bool) []string {
	switch *match.Path.Type {
	case gatewayv1.PathMatchExact:
		terminated := *match.Path.Value + "$"
		if addRegexPrefix {
			terminated = translators.KongPathRegexPrefix + terminated
		}
		return []string{terminated}

	case gatewayv1.PathMatchPathPrefix:
		paths := make([]string, 0, 2)
		path := *match.Path.Value
		if addRegexPrefix {
			paths = append(paths, fmt.Sprintf("%s%s$", translators.KongPathRegexPrefix, path))
			if !strings.HasSuffix(path, "/") {
				path = fmt.Sprintf("%s/", path)
			}
		}
		return append(paths, path)

	case gatewayv1.PathMatchRegularExpression:
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

func httpBackendRefsToBackendRefs(httpBackendRef []gatewayv1.HTTPBackendRef) []gatewayv1.BackendRef {
	backendRefs := make([]gatewayv1.BackendRef, 0, len(httpBackendRef))

	for _, hRef := range httpBackendRef {
		backendRefs = append(backendRefs, hRef.BackendRef)
	}
	return backendRefs
}

// ingressRulesFromSplitHTTPRouteMatchWithPriority translates a single match split from HTTPRoute
// to ingress rule, including Kong service and Kong route.
func (p *Parser) ingressRulesFromSplitHTTPRouteMatchWithPriority(
	rules *ingressRules,
	httpRouteMatchWithPriority translators.SplitHTTPRouteMatchToKongRoutePriority,
) error {
	match := httpRouteMatchWithPriority.Match
	httpRoute := httpRouteMatchWithPriority.Match.Source
	if match.RuleIndex >= len(httpRoute.Spec.Rules) {
		p.logger.Errorf("split match has rule index %d out of bound of rules in source HTTPRoute %d",
			match.RuleIndex, len(httpRoute.Spec.Rules))
		return nil
	}

	rule := httpRoute.Spec.Rules[match.RuleIndex]
	backendRefs := httpBackendRefsToBackendRefs(rule.BackendRefs)
	serviceName := translators.KongServiceNameFromSplitHTTPRouteMatch(httpRouteMatchWithPriority.Match)

	kongService, err := generateKongServiceFromBackendRefWithName(
		p.logger,
		p.storer,
		rules,
		serviceName,
		httpRoute,
		"http",
		backendRefs...,
	)
	if err != nil {
		return err
	}

	kongService.Routes = append(
		kongService.Routes,
		translators.KongExpressionRouteFromHTTPRouteMatchWithPriority(httpRouteMatchWithPriority),
	)
	// cache the service to avoid duplicates in further loop iterations
	rules.ServiceNameToServices[serviceName] = kongService
	rules.ServiceNameToParent[serviceName] = httpRoute
	return nil
}
