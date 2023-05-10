package translators

import (
	"sort"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// GenerateKongExpressionRoutesFromHTTPRouteMatches generates Kong routes from HTTPRouteRule
// pointing to a specific backend.
func GenerateKongExpressionRoutesFromHTTPRouteMatches(
	translation KongRouteTranslation,
	ingressObjectInfo util.K8sObjectInfo,
	hostnames []string,
	tags []*string,
) ([]kongstate.Route, error) {
	// initialize the route with route name, preserve_host, and tags.
	r := kongstate.Route{
		Ingress: ingressObjectInfo,
		Route: kong.Route{
			Name:         kong.String(translation.Name),
			PreserveHost: kong.Bool(true),
			Tags:         tags,
		},
		ExpressionRoutes: true,
	}

	if len(translation.Matches) == 0 {
		if len(hostnames) == 0 {
			return []kongstate.Route{}, ErrRouteValidationNoMatchRulesOrHostnamesSpecified
		}

		hostMatcher := hostMatcherFromHosts(hostnames)
		atc.ApplyExpression(&r.Route, hostMatcher, 1)
		return []kongstate.Route{r}, nil
	}

	_, hasRedirectFilter := lo.Find(translation.Filters, func(filter gatewayv1beta1.HTTPRouteFilter) bool {
		return filter.Type == gatewayv1beta1.HTTPRouteFilterRequestRedirect
	})
	// if the rule has request redirect filter(s), we need to generate a route for each match to
	// attach plugins for the filter.
	if hasRedirectFilter {
		return generateKongExpressionRoutesWithRequestRedirectFilter(translation, ingressObjectInfo, hostnames, tags)
	}

	// if we do not need to generate a kong route for each match, we OR matchers from all matches together.
	routeMatcher := atc.And(atc.Or(generateMatchersFromHTTPRouteMatches(translation.Matches)...))
	// add matcher from parent httproute (hostnames, protocols, SNIs) to be ANDed with the matcher from match.
	matchersFromParent := matchersFromParentHTTPRoute(hostnames, ingressObjectInfo.Annotations)
	for _, matcher := range matchersFromParent {
		routeMatcher.And(matcher)
	}

	atc.ApplyExpression(&r.Route, routeMatcher, 1)
	// generate plugins.
	plugins := GeneratePluginsFromHTTPRouteFilters(translation.Filters, "", tags)
	r.Plugins = plugins
	return []kongstate.Route{r}, nil
}

func generateKongExpressionRoutesWithRequestRedirectFilter(
	translation KongRouteTranslation,
	ingressObjectInfo util.K8sObjectInfo,
	hostnames []string,
	tags []*string,
) ([]kongstate.Route, error) {
	routes := make([]kongstate.Route, 0, len(translation.Matches))
	for _, match := range translation.Matches {
		matchRoute := kongstate.Route{
			Ingress: ingressObjectInfo,
			Route: kong.Route{
				Name:         kong.String(translation.Name),
				PreserveHost: kong.Bool(true),
				Tags:         tags,
			},
			ExpressionRoutes: true,
		}
		// generate matcher for this HTTPRoute Match.
		matcher := atc.And(generateMatcherFromHTTPRouteMatch(match))

		// add matcher from parent httproute (hostnames, protocols, SNIs) to be ANDed with the matcher from match.
		matchersFromParent := matchersFromParentHTTPRoute(hostnames, ingressObjectInfo.Annotations)
		for _, m := range matchersFromParent {
			matcher.And(m)
		}
		atc.ApplyExpression(&matchRoute.Route, matcher, 1)

		// we need to extract the path to configure redirect path of the plugins for request redirect filter.
		path := ""
		if match.Path != nil && match.Path.Value != nil {
			path = *match.Path.Value
		}
		plugins := GeneratePluginsFromHTTPRouteFilters(translation.Filters, path, tags)
		matchRoute.Plugins = plugins

		routes = append(routes, matchRoute)
	}
	return routes, nil
}

func generateMatchersFromHTTPRouteMatches(matches []gatewayv1beta1.HTTPRouteMatch) []atc.Matcher {
	ret := make([]atc.Matcher, 0, len(matches))
	for _, match := range matches {
		matcher := generateMatcherFromHTTPRouteMatch(match)
		ret = append(ret, matcher)
	}
	return ret
}

func generateMatcherFromHTTPRouteMatch(match gatewayv1beta1.HTTPRouteMatch) atc.Matcher {
	matcher := atc.And()

	if match.Path != nil {
		pathMatcher := pathMatcherFromHTTPPathMatch(match.Path)
		matcher.And(pathMatcher)
	}

	if len(match.Headers) > 0 {
		headerMatcher := headerMatcherFromHTTPHeaderMatches(match.Headers)
		matcher.And(headerMatcher)
	}

	if match.Method != nil {
		method := *match.Method
		methodMatcher := methodMatcherFromMethods([]string{string(method)})
		matcher.And(methodMatcher)
	}
	return matcher
}

func appendRegexBeginIfNotExist(regex string) string {
	if !strings.HasPrefix(regex, "^") {
		return "^" + regex
	}
	return regex
}

func pathMatcherFromHTTPPathMatch(pathMatch *gatewayv1beta1.HTTPPathMatch) atc.Matcher {
	path := ""
	if pathMatch.Value != nil {
		path = *pathMatch.Value
	}
	switch *pathMatch.Type {
	case gatewayv1beta1.PathMatchExact:
		return atc.NewPredicateHTTPPath(atc.OpEqual, path)
	case gatewayv1beta1.PathMatchPathPrefix:
		if path == "" || path == "/" {
			return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/")
		}
		// if path ends with /, we should remove the trailing / because it should be ignored:
		// https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io/v1beta1.PathMatchType
		path = strings.TrimSuffix(path, "/")
		return atc.Or(
			atc.NewPredicateHTTPPath(atc.OpEqual, path),
			atc.NewPredicateHTTPPath(atc.OpPrefixMatch, path+"/"),
		)
	case gatewayv1beta1.PathMatchRegularExpression:
		// TODO: for compatibility with kong traditional routes, here we append the ^ prefix to match the path from beginning.
		// Could we allow the regex to match any part of the path?
		// https://github.com/Kong/kubernetes-ingress-controller/issues/3983
		return atc.NewPredicateHTTPPath(atc.OpRegexMatch, appendRegexBeginIfNotExist(path))
	}

	return nil // should be unreachable
}

func headerMatcherFromHTTPHeaderMatches(headerMatches []gatewayv1beta1.HTTPHeaderMatch) atc.Matcher {
	// sort headerMatches by names to generate a stable output.
	sort.Slice(headerMatches, func(i, j int) bool {
		return string(headerMatches[i].Name) < string(headerMatches[j].Name)
	})

	matchers := make([]atc.Matcher, 0, len(headerMatches))
	for _, headerMatch := range headerMatches {
		matchType := gatewayv1beta1.HeaderMatchExact
		if headerMatch.Type != nil {
			matchType = *headerMatch.Type
		}
		headerKey := strings.ReplaceAll(strings.ToLower(string(headerMatch.Name)), "-", "_")
		switch matchType {
		case gatewayv1beta1.HeaderMatchExact:
			matchers = append(matchers, atc.NewPredicateHTTPHeader(headerKey, atc.OpEqual, headerMatch.Value))
		case gatewayv1beta1.HeaderMatchRegularExpression:
			matchers = append(matchers, atc.NewPredicateHTTPHeader(headerKey, atc.OpRegexMatch, headerMatch.Value))
		}
	}
	return atc.And(matchers...)
}

func matchersFromParentHTTPRoute(hostnames []string, metaAnnotations map[string]string) []atc.Matcher {
	// translate hostnames.
	ret := []atc.Matcher{}
	if len(hostnames) > 0 {
		hostMatcher := hostMatcherFromHosts(hostnames)
		ret = append(ret, hostMatcher)
	}

	// translate protocols.
	protocols := []string{"http", "https"}
	// override from "protocols" key in annotations.
	annonationProtocols := annotations.ExtractProtocolNames(metaAnnotations)
	if len(annonationProtocols) > 0 {
		protocols = annonationProtocols
	}
	protocolMatcher := protocolMatcherFromProtocols(protocols)
	ret = append(ret, protocolMatcher)

	// translate SNIs.
	snis, exist := annotations.ExtractSNIs(metaAnnotations)
	if exist && len(snis) > 0 {
		sniMatcher := sniMatcherFromSNIs(snis)
		ret = append(ret, sniMatcher)
	}
	return ret
}
