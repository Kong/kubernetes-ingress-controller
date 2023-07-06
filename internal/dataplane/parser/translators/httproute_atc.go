package translators

import (
	"fmt"
	"sort"
	"strconv"
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

func headerMatcherFromHTTPHeaderMatch(headerMatch gatewayv1beta1.HTTPHeaderMatch) atc.Matcher {
	matchType := gatewayv1beta1.HeaderMatchExact
	if headerMatch.Type != nil {
		matchType = *headerMatch.Type
	}
	headerKey := strings.ReplaceAll(strings.ToLower(string(headerMatch.Name)), "-", "_")
	switch matchType {
	case gatewayv1beta1.HeaderMatchExact:
		return atc.NewPredicateHTTPHeader(headerKey, atc.OpEqual, headerMatch.Value)
	case gatewayv1beta1.HeaderMatchRegularExpression:
		return atc.NewPredicateHTTPHeader(headerKey, atc.OpRegexMatch, headerMatch.Value)
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
		matchers = append(matchers, headerMatcherFromHTTPHeaderMatch(headerMatch))
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

const (
	InternalRuleIndexAnnotationKey  = "internal-rule-index"
	InternalMatchIndexAnnotationKey = "internal-match-index"
)

// SplitHTTPRoute split HTTPRoutes into HTTPRoutes with at most one hostname, and at most one rule
// with exactly one match. It will split one rule with multiple hostnames and multiple matches
// to one hostname and one match per each HTTPRoute.
func SplitHTTPRoute(httproute *gatewayv1beta1.HTTPRoute) []*gatewayv1beta1.HTTPRoute {
	hostnamedRoutes := []*gatewayv1beta1.HTTPRoute{}
	if len(httproute.Spec.Hostnames) == 0 {
		hostnamedRoutes = append(hostnamedRoutes, httproute.DeepCopy())
	} else {
		for _, hostname := range httproute.Spec.Hostnames {
			hostNamedRoute := httproute.DeepCopy()
			hostNamedRoute.Spec.Hostnames = []gatewayv1beta1.Hostname{hostname}
			hostnamedRoutes = append(hostnamedRoutes, hostNamedRoute)
		}
	}

	newHTTPRoutes := []*gatewayv1beta1.HTTPRoute{}
	for _, route := range hostnamedRoutes {
		for i, rule := range route.Spec.Rules {
			for j, match := range rule.Matches {
				splittedRoute := route.DeepCopy()
				splittedRoute.Spec.Rules = []gatewayv1beta1.HTTPRouteRule{
					{
						Matches:     []gatewayv1beta1.HTTPRouteMatch{match},
						Filters:     rule.Filters,
						BackendRefs: rule.BackendRefs,
					},
				}
				if splittedRoute.Annotations == nil {
					splittedRoute.Annotations = map[string]string{}
				}
				splittedRoute.Annotations[InternalRuleIndexAnnotationKey] = strconv.Itoa(i)
				splittedRoute.Annotations[InternalMatchIndexAnnotationKey] = strconv.Itoa(j)
				newHTTPRoutes = append(newHTTPRoutes, splittedRoute)
			}
		}
	}

	return newHTTPRoutes
}

type SplittedHTTPRouteToKongRoutePriority struct {
	HTTPRoute *gatewayv1beta1.HTTPRoute
	Priority  int
}

type HTTPRoutePriorityTraits struct {
	PreciseHostname bool
	HostnameLength  int
	PathType        gatewayv1beta1.PathMatchType
	PathLength      int
	HeaderCount     int
	HasMethodMatch  bool
	QueryParamCount int
}

func CalculateHTTPRoutePriorityTraits(httpRoute *gatewayv1beta1.HTTPRoute) HTTPRoutePriorityTraits {
	traits := HTTPRoutePriorityTraits{}
	if len(httpRoute.Spec.Hostnames) != 0 {
		hostname := httpRoute.Spec.Hostnames[0]
		traits.HostnameLength = len(hostname)
		if !strings.HasPrefix(string(hostname), "*") {
			traits.PreciseHostname = true
		}
	}

	if len(httpRoute.Spec.Rules) > 0 && len(httpRoute.Spec.Rules[0].Matches) > 0 {
		match := httpRoute.Spec.Rules[0].Matches[0]
		if match.Path != nil {
			// fill path type.
			if match.Path.Type != nil {
				traits.PathType = *match.Path.Type
			}
			// fill path length.
			if match.Path.Value != nil {
				traits.PathLength = len(*match.Path.Value)
			}
		}

		// fill number of header matches.
		traits.HeaderCount = len(match.Headers)
		// fill method match.
		if match.Method != nil {
			traits.HasMethodMatch = true
		}
		// fill number of query parameters.
		traits.QueryParamCount = len(match.QueryParams)
	}
	return traits
}

func (t HTTPRoutePriorityTraits) EncodeToPriority() int {
	const (
		// PreciseHostnameShiftBits assigns bit 49 for marking if the hostname is non-wildcard.
		PreciseHostnameShiftBits = 49
		// HostnameLengthShiftBits assigns bits 41-48 for the length of hostname.
		HostnameLengthShiftBits = 41
		// ExactPathShiftBits assigns bit 40 to mark if the match is exact path match.
		ExactPathShiftBits = 40
		// RegularExpressionPathShiftBits assigns bit 39 to mark if the match is regex path match.
		RegularExpressionPathShiftBits = 39
		// PathLengthShiftBits assigns bits 29-38 to path length. (max length = 1024, but must start with /)
		PathLengthShiftBits = 29
		// HeaderNumberShiftBits assign bits 24-28 to number of headers. (max number of headers = 16)
		HeaderNumberShiftBits = 24
		// MethodMatchShiftBits assigns bit 23 to mark if method is specified.
		MethodMatchShiftBits = 23
		// QueryParamNumberShiftBits makes bits 18-22 used for number of query params (max number of query params = 16)
		QueryParamNumberShiftBits = 18
		// bits 0-17 are used for relative order of creation timestamp, namespace/name, and internal order of rules and matches.
	)

	var priority int
	if t.PreciseHostname {
		priority += (1 << PreciseHostnameShiftBits)
	}
	priority += t.HostnameLength << HostnameLengthShiftBits

	if t.PathType == gatewayv1beta1.PathMatchExact {
		priority += (1 << ExactPathShiftBits)
	}
	if t.PathType == gatewayv1beta1.PathMatchRegularExpression {
		priority += (1 << RegularExpressionPathShiftBits)
	}

	// max length of path is 1024, but path must start with /, so we use PathLength-1 to fill the bits.
	if t.PathLength > 0 {
		priority += ((t.PathLength - 1) << PathLengthShiftBits)
	}

	priority += (t.HeaderCount << HeaderNumberShiftBits)
	if t.HasMethodMatch {
		priority += (1 << MethodMatchShiftBits)
	}
	priority += (t.QueryParamCount << QueryParamNumberShiftBits)
	priority += (ResourceKindBitsHTTPRoute << FromResourceKindPriorityShiftBits)

	return priority
}

func AssignRoutePriorityToSplittedHTTPRoutes(
	splittedHTTPRoutes []*gatewayv1beta1.HTTPRoute,
) []SplittedHTTPRouteToKongRoutePriority {
	priorityToHTTPRoutes := map[int][]*gatewayv1beta1.HTTPRoute{}

	for _, httpRoute := range splittedHTTPRoutes {
		anns := httpRoute.Annotations
		// skip if HTTPRoute does not contain the annotation, because this means the HTTPRoute is not a splitted one.
		if anns == nil || anns[InternalRuleIndexAnnotationKey] == "" || anns[InternalMatchIndexAnnotationKey] == "" {
			continue
		}

		priority := CalculateHTTPRoutePriorityTraits(httpRoute).EncodeToPriority()
		priorityToHTTPRoutes[priority] = append(priorityToHTTPRoutes[priority], httpRoute)
	}

	httpRoutesToPriorities := make([]SplittedHTTPRouteToKongRoutePriority, 0, len(splittedHTTPRoutes))

	const defaultRelativeOrderPriorityBits = (1 << 18) - 1
	for priority, routes := range priorityToHTTPRoutes {
		if len(routes) == 1 {
			httpRoutesToPriorities = append(httpRoutesToPriorities, SplittedHTTPRouteToKongRoutePriority{
				HTTPRoute: routes[0],
				Priority:  priority + defaultRelativeOrderPriorityBits,
			})
			continue
		}

		sort.Slice(routes, func(i, j int) bool {
			return compareSplittedHTTPRoutesRelativePriority(routes[i], routes[j])
		})

		relativeOrderBits := defaultRelativeOrderPriorityBits
		for i, route := range routes {
			relativeOrderBits = defaultRelativeOrderPriorityBits - i
			httpRoutesToPriorities = append(httpRoutesToPriorities, SplittedHTTPRouteToKongRoutePriority{
				HTTPRoute: route,
				Priority:  priority + relativeOrderBits,
			})
		}
	}

	return httpRoutesToPriorities
}

func compareSplittedHTTPRoutesRelativePriority(route1, route2 *gatewayv1beta1.HTTPRoute) bool {
	// compare by creation timestamp.
	if !route1.CreationTimestamp.Equal(&route2.CreationTimestamp) {
		return route1.CreationTimestamp.Before(&route2.CreationTimestamp)
	}
	// compare by namespace.
	if route1.Namespace != route2.Namespace {
		return route1.Namespace < route2.Namespace
	}
	// compare by name.
	if route1.Name != route2.Name {
		return route1.Name < route2.Name
	}
	// if ties still exist, compare by internal rule order and match order.
	ruleIndex1, _ := strconv.Atoi(route1.Annotations[InternalRuleIndexAnnotationKey])
	ruleIndex2, _ := strconv.Atoi(route2.Annotations[InternalRuleIndexAnnotationKey])
	if ruleIndex1 != ruleIndex2 {
		return ruleIndex1 < ruleIndex2
	}

	matchIndex1, _ := strconv.Atoi(route1.Annotations[InternalMatchIndexAnnotationKey])
	matchIndex2, _ := strconv.Atoi(route2.Annotations[InternalMatchIndexAnnotationKey])
	if matchIndex1 != matchIndex2 {
		return matchIndex1 < matchIndex2
	}

	// should be unreachable.
	return true
}

// getHTTPRouteHostnamesAsSliceOfStrings translates the hostnames defined in an
// HTTPRoute specification into a []*string slice, which is the type required by translating to matchers
// in expression based routes.
func getHTTPRouteHostnamesAsSliceOfStrings(httproute *gatewayv1beta1.HTTPRoute) []string {
	return lo.Map(httproute.Spec.Hostnames, func(h gatewayv1beta1.Hostname, _ int) string {
		return string(h)
	})
}

// KongExpressionRouteFromHTTPRouteWithPriority translates splitted HTTPRoute into expression
// based kong route with assigned priority.
// the HTTPRoute should have at most one hostname, and at most one rule having exactly one match.
func KongExpressionRouteFromHTTPRouteWithPriority(
	httpRouteWithPriority SplittedHTTPRouteToKongRoutePriority,
) kongstate.Route {
	httproute := httpRouteWithPriority.HTTPRoute
	tags := util.GenerateTagsForObject(httproute)
	routeName := fmt.Sprintf("httproute.%s.%s.%s.%s",
		httproute.Namespace,
		httproute.Name,
		httproute.Annotations[InternalRuleIndexAnnotationKey],
		httproute.Annotations[InternalMatchIndexAnnotationKey],
	)

	r := kongstate.Route{
		Route: kong.Route{
			Name:         kong.String(routeName),
			PreserveHost: kong.Bool(true),
			Tags:         tags,
		},
		Ingress:          util.FromK8sObject(httproute),
		ExpressionRoutes: true,
	}

	hostnames := getHTTPRouteHostnamesAsSliceOfStrings(httproute)
	matchers := matchersFromParentHTTPRoute(hostnames, httproute.Annotations)

	if len(httproute.Spec.Rules) > 0 && len(httproute.Spec.Rules[0].Matches) > 0 {
		matchers = append(matchers, generateMatcherFromHTTPRouteMatch(httproute.Spec.Rules[0].Matches[0]))
	}
	atc.ApplyExpression(&r.Route, atc.And(matchers...), httpRouteWithPriority.Priority)

	// translate filters in the rule.
	if len(httproute.Spec.Rules) > 0 {
		rule := httproute.Spec.Rules[0]
		path := ""
		// since we have at most one match per rule, we do not need to generate request redirect for each match.
		if len(rule.Matches) > 0 {
			match := rule.Matches[0]
			if match.Path != nil && match.Path.Value != nil {
				path = *match.Path.Value
			}
		}

		plugins := GeneratePluginsFromHTTPRouteFilters(rule.Filters, path, tags)
		r.Plugins = plugins
	}

	return r
}

func KongServiceNameFromHTTPRouteWithPriority(
	httpRouteWithPriority SplittedHTTPRouteToKongRoutePriority,
) string {
	httproute := httpRouteWithPriority.HTTPRoute
	return fmt.Sprintf("httproute.%s.%s.%s",
		httproute.Namespace,
		httproute.Name,
		httproute.Annotations[InternalRuleIndexAnnotationKey],
	)
}
