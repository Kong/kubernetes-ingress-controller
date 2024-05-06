package subtranslator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/atc"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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
			// stripPath needs to be disabled by default to be conformant with the Gateway API
			StripPath: kong.Bool(false),
			Tags:      tags,
		},
		ExpressionRoutes: true,
	}

	if len(translation.Matches) == 0 {
		if len(hostnames) == 0 {
			r.Expression = kong.String(CatchAllHTTPExpression)
			return []kongstate.Route{r}, nil
		}
		hostMatcher := hostMatcherFromHosts(hostnames)
		atc.ApplyExpression(&r.Route, hostMatcher, 1)
		return []kongstate.Route{r}, nil
	}

	hasRedirectFilter := lo.ContainsBy(translation.Filters, func(filter gatewayapi.HTTPRouteFilter) bool {
		return filter.Type == gatewayapi.HTTPRouteFilterRequestRedirect
	})
	// if the rule has request redirect filter(s), we need to generate a route for each match to
	// attach plugins for the filter.
	if hasRedirectFilter {
		return generateKongExpressionRoutesWithRequestRedirectFilter(translation, ingressObjectInfo, hostnames, tags)
	}

	// if we do not need to generate a kong route for each match, we OR matchers from all matches together.
	routeMatcher := atc.And(atc.Or(generateMatchersFromHTTPRouteMatches(translation.Matches)...))
	// Add matcher from parent httproute (hostnames, SNIs) to be ANDed with the matcher from match.
	matchersFromParent := matchersFromParentHTTPRoute(hostnames, ingressObjectInfo.Annotations)
	for _, matcher := range matchersFromParent {
		routeMatcher.And(matcher)
	}

	atc.ApplyExpression(&r.Route, routeMatcher, 1)
	// generate plugins.
	if err := SetRoutePlugins(&r, translation.Filters, "", tags, true); err != nil {
		return nil, err
	}
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
				StripPath:    kong.Bool(false),
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
		if err := SetRoutePlugins(&matchRoute, translation.Filters, path, tags, true); err != nil {
			return nil, err
		}
		routes = append(routes, matchRoute)
	}
	return routes, nil
}

func generateMatchersFromHTTPRouteMatches(matches []gatewayapi.HTTPRouteMatch) []atc.Matcher {
	ret := make([]atc.Matcher, 0, len(matches))
	for _, match := range matches {
		matcher := generateMatcherFromHTTPRouteMatch(match)
		ret = append(ret, matcher)
	}
	return ret
}

func generateMatcherFromHTTPRouteMatch(match gatewayapi.HTTPRouteMatch) atc.Matcher {
	matcher := atc.And()

	if match.Path != nil {
		pathMatcher := pathMatcherFromHTTPPathMatch(match.Path)
		matcher.And(pathMatcher)
	}

	if len(match.Headers) > 0 {
		headerMatcher := headerMatcherFromHTTPHeaderMatches(match.Headers)
		matcher.And(headerMatcher)
	}

	if len(match.QueryParams) > 0 {
		queryMatcher := queryParamMatcherFromHTTPQueryParamMatches(match.QueryParams)
		matcher.And(queryMatcher)
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

func pathMatcherFromHTTPPathMatch(pathMatch *gatewayapi.HTTPPathMatch) atc.Matcher {
	path := ""
	if pathMatch.Value != nil {
		path = *pathMatch.Value
	}
	switch *pathMatch.Type {
	case gatewayapi.PathMatchExact:
		return atc.NewPredicateHTTPPath(atc.OpEqual, path)
	case gatewayapi.PathMatchPathPrefix:
		if path == "" || path == "/" {
			return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/")
		}
		// if path ends with /, we should remove the trailing / because it should be ignored:
		// https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io/v1.PathMatchType
		path = strings.TrimSuffix(path, "/")
		return atc.Or(
			atc.NewPredicateHTTPPath(atc.OpEqual, path),
			atc.NewPredicateHTTPPath(atc.OpPrefixMatch, path+"/"),
		)
	case gatewayapi.PathMatchRegularExpression:
		// TODO: for compatibility with kong traditional routes, here we append the ^ prefix to match the path from beginning.
		// Could we allow the regex to match any part of the path?
		// https://github.com/Kong/kubernetes-ingress-controller/issues/3983
		return atc.NewPredicateHTTPPath(atc.OpRegexMatch, appendRegexBeginIfNotExist(path))
	}

	return nil // should be unreachable
}

func headerMatcherFromHTTPHeaderMatch(headerMatch gatewayapi.HTTPHeaderMatch) atc.Matcher {
	matchType := gatewayapi.HeaderMatchExact
	if headerMatch.Type != nil {
		matchType = *headerMatch.Type
	}
	headerKey := strings.ReplaceAll(strings.ToLower(string(headerMatch.Name)), "-", "_")
	switch matchType {
	case gatewayapi.HeaderMatchExact:
		return atc.NewPredicateHTTPHeader(headerKey, atc.OpEqual, headerMatch.Value)
	case gatewayapi.HeaderMatchRegularExpression:
		return atc.NewPredicateHTTPHeader(headerKey, atc.OpRegexMatch, headerMatch.Value)
	}
	return nil // should be unreachable
}

func headerMatcherFromHTTPHeaderMatches(headerMatches []gatewayapi.HTTPHeaderMatch) atc.Matcher {
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

func queryParamMatcherFromHTTPQueryParamMatch(queryParamMatch gatewayapi.HTTPQueryParamMatch) atc.Matcher {
	matchType := gatewayapi.QueryParamMatchExact
	if queryParamMatch.Type != nil {
		matchType = *queryParamMatch.Type
	}
	switch matchType {
	case gatewayapi.QueryParamMatchExact:
		return atc.NewPredicateHTTPQuery(string(queryParamMatch.Name), atc.OpEqual, queryParamMatch.Value)
	case gatewayapi.QueryParamMatchRegularExpression:
		return atc.NewPredicateHTTPQuery(string(queryParamMatch.Name), atc.OpRegexMatch, queryParamMatch.Value)
	}
	return nil // should be unreachable
}

func queryParamMatcherFromHTTPQueryParamMatches(queryParamMatches []gatewayapi.HTTPQueryParamMatch) atc.Matcher {
	// Sort queryParamMatches by names to generate a stable output.
	sort.Slice(queryParamMatches, func(i, j int) bool {
		return string(queryParamMatches[i].Name) < string(queryParamMatches[j].Name)
	})

	matchers := make([]atc.Matcher, 0, len(queryParamMatches))
	for _, queryParamMatch := range queryParamMatches {
		matchers = append(matchers, queryParamMatcherFromHTTPQueryParamMatch(queryParamMatch))
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

	// translate SNIs.
	snis, exist := annotations.ExtractSNIs(metaAnnotations)
	if exist && len(snis) > 0 {
		sniMatcher := sniMatcherFromSNIs(snis)
		ret = append(ret, sniMatcher)
	}
	return ret
}

type SplitHTTPRouteMatch struct {
	Source     *gatewayapi.HTTPRoute
	Hostname   string
	Match      gatewayapi.HTTPRouteMatch
	RuleIndex  int
	MatchIndex int
}

// SplitHTTPRoute splits HTTPRoutes into matches with at most one hostname, and one rule
// with exactly one match. It will split one rule with multiple hostnames and multiple matches
// to one hostname and one match per each HTTPRoute.
func SplitHTTPRoute(httproute *gatewayapi.HTTPRoute) []SplitHTTPRouteMatch {
	splitHTTPRouteByMatch := func(hostname string) []SplitHTTPRouteMatch {
		ret := []SplitHTTPRouteMatch{}
		for ruleIndex, rule := range httproute.Spec.Rules {
			if len(rule.Matches) == 0 {
				ret = append(ret, SplitHTTPRouteMatch{
					Source:     httproute,
					Hostname:   hostname,
					Match:      gatewayapi.HTTPRouteMatch{},
					RuleIndex:  ruleIndex,
					MatchIndex: 0,
				})
			}
			for matchIndex, match := range rule.Matches {
				ret = append(ret, SplitHTTPRouteMatch{
					Source:     httproute,
					Hostname:   hostname,
					Match:      *match.DeepCopy(),
					RuleIndex:  ruleIndex,
					MatchIndex: matchIndex,
				})
			}
		}
		return ret
	}
	// HTTPRoute has no hostnames. Split by rule and match, with an empty hostname.
	if len(httproute.Spec.Hostnames) == 0 {
		return splitHTTPRouteByMatch("")
	}
	// HTTPRoute has at least one hostname, split by hostname first.
	splitMatches := []SplitHTTPRouteMatch{}
	for _, hostname := range httproute.Spec.Hostnames {
		splitMatches = append(splitMatches, splitHTTPRouteByMatch(string(hostname))...)
	}
	return splitMatches
}

type SplitHTTPRouteMatchToKongRoutePriority struct {
	Match    SplitHTTPRouteMatch
	Priority RoutePriorityType
}

type HTTPRoutePriorityTraits struct {
	PreciseHostname bool
	HostnameLength  int
	PathType        gatewayapi.PathMatchType
	PathLength      int
	HeaderCount     int
	HasMethodMatch  bool
	QueryParamCount int
}

// CalculateHTTPRouteMatchPriorityTraits calculates the parts of priority
// that can be decided by the fields in spec of the match split from HTTPRoute.
// Specification of priority goes as follow:
// (The following comments are extracted from gateway API specification about HTTPRoute)
//
// In the event that multiple HTTPRoutes specify intersecting hostnames,
// precedence must be given to rules from the HTTPRoute with the largest number of:
//
//   - Characters in a matching non-wildcard hostname.
//   - Characters in a matching hostname.
//
// If ties exist across multiple Routes, the matching precedence rules for HTTPRouteMatches takes over.
//
// Proxy or Load Balancer routing configuration generated from HTTPRoutes MUST prioritize matches based on the following criteria, continuing on ties.
// Across all rules specified on applicable Routes, precedence must be given to the match having:
//
//   - "Exact” path match.
//   - "Prefix" path match with largest number of characters.
//   - Method match.
//   - Largest number of header matches.
//   - Largest number of query param matches.
func CalculateHTTPRouteMatchPriorityTraits(match SplitHTTPRouteMatch) HTTPRoutePriorityTraits {
	traits := HTTPRoutePriorityTraits{}
	// fill traits from hostname.
	if len(match.Hostname) != 0 {
		traits.HostnameLength = len(match.Hostname)
		// if the hostname does not start with *, the split HTTPRoute should have precise hostname.
		if !strings.HasPrefix(match.Hostname, "*") {
			traits.PreciseHostname = true
		}
	}

	// fill traits from path match.
	if match.Match.Path != nil {
		pathMatch := match.Match.Path
		// fill path type.
		if pathMatch.Type != nil {
			traits.PathType = *pathMatch.Type
		}
		// fill path length.
		if pathMatch.Value != nil {
			traits.PathLength = len(*pathMatch.Value)
		}
	}

	// fill method match.
	if match.Match.Method != nil {
		traits.HasMethodMatch = true
	}

	// fill number of header matches.
	traits.HeaderCount = len(match.Match.Headers)

	// fill number of query parameters.
	traits.QueryParamCount = len(match.Match.QueryParams)

	return traits
}

// EncodeToPriority turns HTTPRoute priority traits into the integer expressed priority.
//
//		   4                   3                   2                   1
//	 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0
//	+-+---------------+-+-+-------------------+-+---------+---------+-----------------------+
//	|P| host len      |E|R|  Path length      |M|Header No|Query No.| relative order        |
//	+-+---------------+-+-+-------------------+-+---------+-------- +-----------------------+
//
// Where:
// P: set to 1 if the hostname is non-wildcard.
// host len: host length of hostname.
// E: set to 1 if the path type is `Exact`.
// R: set to 1 if the path type in `RegularExpression`.
// Path length: length of `path.Value`.
// M: set to 1 if Method match is specified.
// Header No.: number of header matches.
// Query No.: number of query parameter matches.
// relative order: relative order of creation timestamp, namespace and name and internal rule/match order between different (split) HTTPRoutes.
func (t HTTPRoutePriorityTraits) EncodeToPriority() RoutePriorityType {
	const (
		// PreciseHostnameShiftBits assigns bit 43 for marking if the hostname is non-wildcard.
		PreciseHostnameShiftBits = 43
		// HostnameLengthShiftBits assigns bits 35-42 for the length of hostname.
		HostnameLengthShiftBits = 35
		// ExactPathShiftBits assigns bit 34 to mark if the match is exact path match.
		ExactPathShiftBits = 34
		// RegularExpressionPathShiftBits assigns bit 33 to mark if the match is regex path match.
		RegularExpressionPathShiftBits = 33
		// PathLengthShiftBits assigns bits 23-32 to path length. (max length = 1024, but must start with /)
		PathLengthShiftBits = 23
		// MethodMatchShiftBits assigns bit 22 to mark if method is specified.
		MethodMatchShiftBits = 22
		// HeaderNumberShiftBits assign bits 17-21 to number of headers. (max number of headers = 16)
		HeaderNumberShiftBits = 17
		// QueryParamNumberShiftBits makes bits 12-16 used for number of query params (max number of query params = 16)
		QueryParamNumberShiftBits = 12
		// bits 0-11 are used for relative order of creation timestamp, namespace/name, and internal order of rules and matches.
		// the bits are calculated by sorting HTTPRoutes with the same priority calculated from the fields above
		// and start from all 1s, then decrease by one for each HTTPRoute.
	)

	var priority RoutePriorityType
	if t.PreciseHostname {
		priority += (1 << PreciseHostnameShiftBits)
	}
	priority += RoutePriorityType(t.HostnameLength << HostnameLengthShiftBits)

	if t.PathType == gatewayapi.PathMatchExact {
		priority += (1 << ExactPathShiftBits)
	}
	if t.PathType == gatewayapi.PathMatchRegularExpression {
		priority += (1 << RegularExpressionPathShiftBits)
	}

	// max length of path is 1024, but path must start with /, so we use PathLength-1 to fill the bits.
	if t.PathLength > 0 {
		priority += RoutePriorityType(((t.PathLength - 1) << PathLengthShiftBits))
	}

	priority += RoutePriorityType(t.HeaderCount << HeaderNumberShiftBits)
	if t.HasMethodMatch {
		priority += (1 << MethodMatchShiftBits)
	}
	priority += RoutePriorityType(t.QueryParamCount << QueryParamNumberShiftBits)
	priority += RoutePriorityType(ResourceKindBitsHTTPRoute << FromResourceKindPriorityShiftBits)

	return priority
}

// AssignRoutePriorityToSplitHTTPRouteMatches assigns priority to
// ALL split matches from ALL HTTPRoutes in the cache.
// Firstly assign "fixed" bits by the following fields of the match:
// hostname, path type, path length, method match, number of header matches, number of query param matches.
// If ties exists in the first step, where multiple matches has the same priority
// calculated from the fields, we run a sort for the matches in the tie
// and assign the bits for "relative order" according to the sorting result of these matches.
func AssignRoutePriorityToSplitHTTPRouteMatches(
	logger logr.Logger,
	splitHTTPRouteMatches []SplitHTTPRouteMatch,
) []SplitHTTPRouteMatchToKongRoutePriority {
	priorityToSplitHTTPRouteMatches := map[RoutePriorityType][]SplitHTTPRouteMatch{}

	for _, match := range splitHTTPRouteMatches {
		priority := CalculateHTTPRouteMatchPriorityTraits(match).EncodeToPriority()
		priorityToSplitHTTPRouteMatches[priority] = append(priorityToSplitHTTPRouteMatches[priority], match)
	}

	httpRouteMatchesToPriorities := make([]SplitHTTPRouteMatchToKongRoutePriority, 0, len(splitHTTPRouteMatches))

	// Bits 0-11 (12 bits) are assigned for relative order of matches.
	// If multiple matches are assigned to the same priority in the previous step,
	// sort them then starts with 2^12 -1 and decrease by one for each HTTPRoute;
	// If only one match occupies the priority, fill the relative order bits with all 1s.
	const RelativeOrderAssignedBits = 12
	const defaultRelativeOrderPriorityBits = (uint64(1) << RelativeOrderAssignedBits) - 1
	for priority, matches := range priorityToSplitHTTPRouteMatches {
		if len(matches) == 1 {
			httpRouteMatchesToPriorities = append(httpRouteMatchesToPriorities, SplitHTTPRouteMatchToKongRoutePriority{
				Match:    matches[0],
				Priority: priority + defaultRelativeOrderPriorityBits,
			})
			continue
		}

		sort.SliceStable(matches, func(i, j int) bool {
			return compareSplitHTTPRouteMatchesRelativePriority(matches[i], matches[j])
		})

		for i, match := range matches {
			relativeOrderBits := defaultRelativeOrderPriorityBits - RoutePriorityType(i)
			// Although it is very unlikely that there are 2^12 = 4096 HTTPRoutes
			// should be given priority by their relative order, here we limit the
			// relativeOrderBits to be at least 0.
			if relativeOrderBits <= 0 {
				relativeOrderBits = 0
			}
			httpRouteMatchesToPriorities = append(httpRouteMatchesToPriorities, SplitHTTPRouteMatchToKongRoutePriority{
				Match:    match,
				Priority: priority + relativeOrderBits,
			})
		}
		// Just in case, log a very unlikely scenario where we have more than 2^12 matches with the same base
		// priority and we have no bit space for them to be deterministically ordered.
		if len(matches) > (1 << 12) {
			logger.Error(nil, "Too many HTTPRoute matches to be deterministically ordered", "match_number", len(matches))
		}
	}

	return httpRouteMatchesToPriorities
}

func compareSplitHTTPRouteMatchesRelativePriority(match1, match2 SplitHTTPRouteMatch) bool {
	route1 := match1.Source
	route2 := match2.Source
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
	if match1.RuleIndex != match2.RuleIndex {
		return match1.RuleIndex < match2.RuleIndex
	}

	if match1.MatchIndex != match2.MatchIndex {
		return match1.MatchIndex < match2.MatchIndex
	}

	// should be unreachable.
	return true
}

// KongExpressionRouteFromHTTPRouteMatchWithPriority translates a split HTTPRoute match into expression
// based kong route with assigned priority.
func KongExpressionRouteFromHTTPRouteMatchWithPriority(
	httpRouteMatchWithPriority SplitHTTPRouteMatchToKongRoutePriority,
) (*kongstate.Route, error) {
	match := httpRouteMatchWithPriority.Match
	httproute := httpRouteMatchWithPriority.Match.Source
	tags := util.GenerateTagsForObject(httproute)
	// since we split HTTPRoutes by hostname, rule and match, we generate the route name in
	// httproute.<namespace>.<name>.<hostname>.<rule index>.<match index> format.
	hostnameInRouteName := "_"
	if len(match.Hostname) > 0 {
		hostnameInRouteName = strings.ReplaceAll(match.Hostname, "*", "_")
	}

	routeName := fmt.Sprintf("httproute.%s.%s.%s.%d.%d",
		httproute.Namespace,
		httproute.Name,
		hostnameInRouteName,
		match.RuleIndex,
		match.MatchIndex,
	)

	r := &kongstate.Route{
		Route: kong.Route{
			Name:         kong.String(routeName),
			PreserveHost: kong.Bool(true),
			// stripPath needs to be disabled by default to be conformant with the Gateway API
			StripPath: kong.Bool(false),
			Tags:      tags,
		},
		Ingress:          util.FromK8sObject(httproute),
		ExpressionRoutes: true,
	}
	// generate ATC matcher from hostname in the match and annotations of parent HTTPRoute.
	hostnames := []string{match.Hostname}
	matchers := matchersFromParentHTTPRoute(hostnames, httproute.Annotations)
	// generate ATC matcher from split HTTPRouteMatch itself.
	matchers = append(matchers, generateMatcherFromHTTPRouteMatch(match.Match))

	atc.ApplyExpression(&r.Route, atc.And(matchers...), httpRouteMatchWithPriority.Priority)

	// generate a "catch-all" route if the generated expression is empty.
	if r.Expression == nil || len(*r.Expression) == 0 {
		r.Expression = kong.String(CatchAllHTTPExpression)
		r.Priority = kong.Uint64(httpRouteMatchWithPriority.Priority)
	}

	// translate filters in the rule.
	if match.RuleIndex < len(httproute.Spec.Rules) {
		rule := httproute.Spec.Rules[match.RuleIndex]
		path := ""
		// since we have one match to translate, we do not need to generate request redirect for each match.
		if match.Match.Path != nil && match.Match.Path.Value != nil {
			path = *match.Match.Path.Value
		}

		if err := SetRoutePlugins(r, rule.Filters, path, tags, true); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// KongServiceNameFromSplitHTTPRouteMatch generates service name from split HTTPRoute match.
// since one HTTPRoute may be split by hostname and rule, the service name will be generated
// in the format "httproute.<namespace>.<name>.<hostname>.<rule index>".
// For example: `httproute.default.example.foo.com.0`.
func KongServiceNameFromSplitHTTPRouteMatch(match SplitHTTPRouteMatch) string {
	httproute := match.Source
	hostname := "_"
	if len(match.Hostname) > 0 {
		hostname = strings.ReplaceAll(match.Hostname, "*", "_")
	}
	return fmt.Sprintf("httproute.%s.%s.%s.%d",
		httproute.Namespace,
		httproute.Name,
		hostname,
		match.RuleIndex,
	)
}
