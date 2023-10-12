package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
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
		p.logger.Error(err, "failed to list HTTPRoutes")
		return result
	}

	httpRoutesToTranslate := make([]*gatewayapi.HTTPRoute, 0, len(httpRouteList))
	for _, httproute := range httpRouteList {
		// Validate each HTTPRoute before translating and register translation failures if an HTTPRoute is invalid.
		if err := validateHTTPRoute(httproute, p.featureFlags); err != nil {
			p.registerTranslationFailure(fmt.Sprintf("HTTPRoute can't be routed: %v", err), httproute)
			continue
		}
		httpRoutesToTranslate = append(httpRoutesToTranslate, httproute)
	}

	if p.featureFlags.ExpressionRoutes {
		p.ingressRulesFromHTTPRoutesUsingExpressionRoutes(httpRoutesToTranslate, &result)
		return result
	}

	for _, httproute := range httpRoutesToTranslate {
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

// ingressRulesFromHTTPRoute validates and generates a set of proto-Kong routes (ingress rules) from an HTTPRoute.
// If multiple rules in the HTTPRoute use the same Service, it combines them into a single Kong route.
func (p *Parser) ingressRulesFromHTTPRoute(result *ingressRules, httproute *gatewayapi.HTTPRoute) error {
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
			routes, err := GenerateKongRouteFromTranslation(httproute, kongRouteTranslation, p.featureFlags.ExpressionRoutes)
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

func validateHTTPRoute(httproute *gatewayapi.HTTPRoute, featureFlags FeatureFlags) error {
	spec := httproute.Spec

	// validation for HTTPRoutes will happen at a higher layer, but in spite of that we run
	// validation at this level as well as a fallback so that if routes are posted which
	// are invalid somehow make it past validation (e.g. the webhook is not enabled) we can
	// at least try to provide a helpful message about the situation in the manager logs.
	if len(spec.Rules) == 0 {
		return translators.ErrRouteValidationNoRules
	}

	// Kong supports query parameter match only with expression router,
	// so we return error when query param match is specified and expression router is not enabled in the parser.
	if !featureFlags.ExpressionRoutes {
		for _, rule := range spec.Rules {
			for _, match := range rule.Matches {
				if len(match.QueryParams) > 0 {
					return translators.ErrRouteValidationQueryParamMatchesUnsupported
				}
			}
		}
	}

	return nil
}

// ingressRulesFromHTTPRoutesUsingExpressionRoutes translates HTTPRoutes to expression based routes
// when ExpressionRoutes feature flag is enabled.
// Because we need to assign different priorities based on the hostname and match in the specification of HTTPRoutes,
// We need to split the HTTPRoutes into ones with only one hostname and one match, then assign priority to them
// and finally translate the split HTTPRoutes into Kong services and routes with assigned priorities.
func (p *Parser) ingressRulesFromHTTPRoutesUsingExpressionRoutes(httpRoutes []*gatewayapi.HTTPRoute, result *ingressRules) {
	// first, split HTTPRoutes by hostnames and matches.
	splitHTTPRouteMatches := []translators.SplitHTTPRouteMatch{}
	for _, httproute := range httpRoutes {
		splitHTTPRouteMatches = append(splitHTTPRouteMatches, translators.SplitHTTPRoute(httproute)...)
	}
	// assign priorities to split HTTPRoutes.
	splitHTTPRoutesWithPriorities := translators.AssignRoutePriorityToSplitHTTPRouteMatches(p.logger, splitHTTPRouteMatches)
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

// -----------------------------------------------------------------------------
// Translate HTTPRoute - Utils
// -----------------------------------------------------------------------------

// getHTTPRouteHostnamesAsSliceOfStrings translates the hostnames defined in an
// HTTPRoute specification into a []*string slice, which is the type required by translating to matchers
// in expression based routes.
func getHTTPRouteHostnamesAsSliceOfStrings(httproute *gatewayapi.HTTPRoute) []string {
	return lo.Map(httproute.Spec.Hostnames, func(h gatewayapi.Hostname, _ int) string {
		return string(h)
	})
}

// getHTTPRouteHostnamesAsSliceOfStringPointers translates the hostnames defined
// in an HTTPRoute specification into a []*string slice, which is the type required
// by kong.Route{}.
func getHTTPRouteHostnamesAsSliceOfStringPointers(httproute *gatewayapi.HTTPRoute) []*string {
	return lo.Map(httproute.Spec.Hostnames, func(h gatewayapi.Hostname, _ int) *string {
		return kong.String(string(h))
	})
}

// GenerateKongRouteFromTranslation generates Kong routes from HTTPRoute
// pointing to a specific backend. It is used for both traditional and expression based routes.
func GenerateKongRouteFromTranslation(
	httproute *gatewayapi.HTTPRoute,
	translation translators.KongRouteTranslation,
	expressionRoutes bool,
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
		tags,
	)
}

// generateKongRoutesFromHTTPRouteMatches converts an HTTPRouteMatches to a slice of Kong Route objects with traditional routes.
// This function assumes that the HTTPRouteMatches share the query params, headers and methods.
func generateKongRoutesFromHTTPRouteMatches(
	routeName string,
	matches []gatewayapi.HTTPRouteMatch,
	filters []gatewayapi.HTTPRouteFilter,
	ingressObjectInfo util.K8sObjectInfo,
	hostnames []*string,
	tags []*string,
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

	_, hasRedirectFilter := lo.Find(filters, func(filter gatewayapi.HTTPRouteFilter) bool {
		return filter.Type == gatewayapi.HTTPRouteFilterRequestRedirect
	})

	routes := getRoutesFromMatches(matches, &r, filters, tags, hasRedirectFilter)

	// if the redirect filter has not been set, we still need to set the route plugins
	if !hasRedirectFilter {
		translators.ConvertFiltersToPlugins(&r, filters, "", tags)
		routes = []kongstate.Route{r}
	}

	return routes, nil
}

// getRoutesFromMatches converts all the httpRoute matches to the proper set of kong routes.
func getRoutesFromMatches(
	matches []gatewayapi.HTTPRouteMatch,
	route *kongstate.Route,
	filters []gatewayapi.HTTPRouteFilter,
	tags []*string,
	hasRedirectFilter bool,
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
				paths := generateKongRoutePathFromHTTPRouteMatch(match)
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
			translators.ConvertFiltersToPlugins(matchRoute, filters, path, tags)

			routes = append(routes, *route)
		} else {
			// Configure path matching information about the route if paths matching was defined
			// Kong automatically infers whether or not a path is a regular expression and uses a prefix match by
			// default if it is not. For those types, we use the path value as-is and let Kong determine the type.
			// For exact matches, we transform the path into a regular expression that terminates after the value.
			if match.Path != nil {
				for _, path := range generateKongRoutePathFromHTTPRouteMatch(match) {
					route.Route.Paths = append(route.Route.Paths, kong.String(path))
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

func generateKongRoutePathFromHTTPRouteMatch(match gatewayapi.HTTPRouteMatch) []string {
	switch *match.Path.Type {
	case gatewayapi.PathMatchExact:
		return []string{translators.KongPathRegexPrefix + *match.Path.Value + "$"}

	case gatewayapi.PathMatchPathPrefix:
		paths := make([]string, 0, 2)
		path := *match.Path.Value
		paths = append(paths, fmt.Sprintf("%s%s$", translators.KongPathRegexPrefix, path))
		if !strings.HasSuffix(path, "/") {
			path = fmt.Sprintf("%s/", path)
		}
		return append(paths, path)

	case gatewayapi.PathMatchRegularExpression:
		return []string{translators.KongPathRegexPrefix + *match.Path.Value}
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

func httpBackendRefsToBackendRefs(httpBackendRef []gatewayapi.HTTPBackendRef) []gatewayapi.BackendRef {
	backendRefs := make([]gatewayapi.BackendRef, 0, len(httpBackendRef))

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
		p.logger.Error(nil, "split match has rule out of bound of rules in source HTTPRoute",
			"rule_index", match.RuleIndex, "rule_count", len(httpRoute.Spec.Rules))
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
