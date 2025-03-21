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

func generateMatcherFromHTTPRouteMatch(match gatewayapi.HTTPRouteMatch, containReplacePrefixMatchURLRewriteFilter bool) atc.Matcher {
	matcher := atc.And()

	if match.Path != nil {
		pathMatcher := pathMatcherFromHTTPPathMatch(match.Path, containReplacePrefixMatchURLRewriteFilter)
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

func pathMatcherFromHTTPPathMatch(pathMatch *gatewayapi.HTTPPathMatch, containReplacePrefixMatchURLRewriteFilter bool) atc.Matcher {
	path := ""
	if pathMatch.Value != nil {
		path = *pathMatch.Value
	}

	switch *pathMatch.Type {
	case gatewayapi.PathMatchExact:
		return atc.NewPredicateHTTPPath(atc.OpEqual, path)
	case gatewayapi.PathMatchPathPrefix:
		// if path ends with /, we should remove the trailing / because it should be ignored.
		// So we normalize the path when the path match type is prefix.
		// https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io/v1.PathMatchType
		path = normalizePath(path)
		pathIsRoot := isPathRoot(path)
		// When the path match type is `prefix` and the rule contains a URLRewrite filtre with `ReplacePrefixMatch`,
		// we should generate a regex matcher with a capture group for the match to let the plugin to extract the path segment to replace.
		if containReplacePrefixMatchURLRewriteFilter {
			exactPrefixPredicate := atc.NewPredicateHTTPPath(atc.OpEqual, path)
			subpathsPredicate := func() atc.Predicate {
				if pathIsRoot {
					// If the path is "/", we don't capture the slash as Kong Route's path has to begin with a slash.
					// If we captured the slash, we'd generate "(/.*)", and it'd be rejected by Kong.
					return atc.NewPredicateHTTPPath(atc.OpRegexMatch, "^/(.*)")
				}
				// If the path is not "/", i.e. it has a prefix, we capture the slash to make it possible to
				// route "/prefix" to "/replacement" and "/prefix/" to "/replacement/" correctly.
				return atc.NewPredicateHTTPPath(atc.OpRegexMatch, fmt.Sprintf("^%s(/.*)", path))
			}()
			return atc.Or(exactPrefixPredicate, subpathsPredicate)
		}
		// Otherwise, we generate an "Or" expression of an exact match and a prefix match, like
		// path == "/prefix" || path ^= "/prefix/".
		if pathIsRoot {
			return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/")
		}

		return atc.Or(
			atc.NewPredicateHTTPPath(atc.OpEqual, path),
			atc.NewPredicateHTTPPath(atc.OpPrefixMatch, path+"/"),
		)
	case gatewayapi.PathMatchRegularExpression:
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
	Source   *gatewayapi.HTTPRoute
	Hostname string
	Match    gatewayapi.HTTPRouteMatch
	// OptionalNamedRouteRule represents RouteName - an optional name
	// of the particular route that can be defined in the K8s HTTPRoute,
	// https://gateway-api.sigs.k8s.io/geps/gep-995/#api.
	OptionalNamedRouteRule string
	RuleIndex              int
	MatchIndex             int
}

// SplitHTTPRoute splits HTTPRoutes into matches with at most one hostname, and one rule
// with exactly one match. It will split one rule with multiple hostnames and multiple matches
// to one hostname and one match per each HTTPRoute.
func SplitHTTPRoute(httproute *gatewayapi.HTTPRoute) []SplitHTTPRouteMatch {
	splitHTTPRouteByMatch := func(hostname string) []SplitHTTPRouteMatch {
		ret := []SplitHTTPRouteMatch{}
		for ruleIndex, rule := range httproute.Spec.Rules {
			optionalNamedRouteRule := string(lo.FromPtr(rule.Name))
			if len(rule.Matches) == 0 {
				ret = append(ret, SplitHTTPRouteMatch{
					Source:                 httproute,
					Hostname:               hostname,
					Match:                  gatewayapi.HTTPRouteMatch{},
					OptionalNamedRouteRule: optionalNamedRouteRule,
					RuleIndex:              ruleIndex,
					MatchIndex:             0,
				})
			}
			for matchIndex, match := range rule.Matches {
				ret = append(ret, SplitHTTPRouteMatch{
					Source:                 httproute,
					Hostname:               hostname,
					Match:                  *match.DeepCopy(),
					OptionalNamedRouteRule: optionalNamedRouteRule,
					RuleIndex:              ruleIndex,
					MatchIndex:             matchIndex,
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

// calculateHTTPRouteMatchPriorityTraits calculates the parts of priority
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
func calculateHTTPRouteMatchPriorityTraits(match SplitHTTPRouteMatch) HTTPRoutePriorityTraits {
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

// assignRoutePriorityToSplitHTTPRouteMatches assigns priority to
// ALL split matches from ALL HTTPRoutes in the cache.
// Firstly assign "fixed" bits by the following fields of the match:
// hostname, path type, path length, method match, number of header matches, number of query param matches.
// If ties exists in the first step, where multiple matches has the same priority
// calculated from the fields, we run a sort for the matches in the tie
// and assign the bits for "relative order" according to the sorting result of these matches.
func assignRoutePriorityToSplitHTTPRouteMatches(
	logger logr.Logger,
	splitHTTPRouteMatches []SplitHTTPRouteMatch,
) []SplitHTTPRouteMatchToKongRoutePriority {
	priorityToSplitHTTPRouteMatches := map[RoutePriorityType][]SplitHTTPRouteMatch{}

	for _, match := range splitHTTPRouteMatches {
		priority := calculateHTTPRouteMatchPriorityTraits(match).EncodeToPriority()
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

// kongExpressionRouteFromHTTPRouteMatchWithPriority translates a split HTTPRoute match into expression
// based kong route with assigned priority.
func kongExpressionRouteFromHTTPRouteMatchWithPriority(
	httpRouteMatchWithPriority SplitHTTPRouteMatchToKongRoutePriority,
	supportRedirectPlugin bool,
) (*kongstate.Route, error) {
	match := httpRouteMatchWithPriority.Match
	httproute := httpRouteMatchWithPriority.Match.Source
	tags := util.GenerateTagsForObject(httproute, util.AdditionalTagsK8sNamedRouteRule(match.OptionalNamedRouteRule)...)

	// Since we split HTTPRoutes by hostname, rule and match, we generate the route name in
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

	containReplacePrefixMatchURLRewriteFilter := lo.ContainsBy(httproute.Spec.Rules[match.RuleIndex].Filters, func(filter gatewayapi.HTTPRouteFilter) bool {
		return filter.Type == gatewayapi.HTTPRouteFilterURLRewrite &&
			filter.URLRewrite.Path != nil && filter.URLRewrite.Path.Type == gatewayapi.PrefixMatchHTTPPathModifier
	})

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
	matchers = append(matchers, generateMatcherFromHTTPRouteMatch(match.Match, containReplacePrefixMatchURLRewriteFilter))

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
		setPluginsOptions := setKongRoutePluginsOptions{
			expressionsRouterEnabled:  true,
			redirectKongPluginEnabled: supportRedirectPlugin,
		}
		if err := setRoutePlugins(r, rule.Filters, path, tags, setPluginsOptions); err != nil {
			return nil, err
		}
	}

	return r, nil
}

// groupHTTPRouteMatchesWithPrioritiesByRule groups split HTTPRoute matches that has priorities assigned by the source HTTPRoute rule,.
func groupHTTPRouteMatchesWithPrioritiesByRule(
	logger logr.Logger, routes []*gatewayapi.HTTPRoute,
) splitHTTPRouteMatchesWithPrioritiesGroupedByRule {
	splitHTTPRouteMatches := []SplitHTTPRouteMatch{}
	for _, route := range routes {
		splitHTTPRouteMatches = append(splitHTTPRouteMatches, SplitHTTPRoute(route)...)
	}
	// assign priorities to split HTTPRoutes.
	splitHTTPRouteMatchesWithPriorities := assignRoutePriorityToSplitHTTPRouteMatches(logger, splitHTTPRouteMatches)

	// group the matches with priorities by its source HTTPRoute rule.
	ruleToSplitMatchesWithPriorities := splitHTTPRouteMatchesWithPrioritiesGroupedByRule{}
	for _, matchWithPriority := range splitHTTPRouteMatchesWithPriorities {
		sourceRoute := matchWithPriority.Match.Source
		ruleKey := fmt.Sprintf("%s/%s.%d", sourceRoute.Namespace, sourceRoute.Name, matchWithPriority.Match.RuleIndex)
		ruleToSplitMatchesWithPriorities[ruleKey] = append(ruleToSplitMatchesWithPriorities[ruleKey], matchWithPriority)
	}
	return ruleToSplitMatchesWithPriorities
}

// translateSplitHTTPRouteMatchesToKongstateRoutesWithExpression translates a list of split HTTPRoute matches with assigned priorities
// that are pointing to the same service to list of kongstate route with expressions.
func translateSplitHTTPRouteMatchesToKongstateRoutesWithExpression(
	matchesWithPriorities []SplitHTTPRouteMatchToKongRoutePriority,
	supportRedirectPlugin bool,
) ([]kongstate.Route, error) {
	routes := make([]kongstate.Route, 0, len(matchesWithPriorities))
	for _, matchWithPriority := range matchesWithPriorities {
		// Since each match is assigned a deterministic priority, we have to generate one route for each split match
		// because every match have a different priority.
		// TODO: update the algorithm to assign priorities to matches to make it possible to consolidate some matches.
		// For example, we can assign the same priority to multiple matches from the same rule if they tie on the priority from the fixed fields:
		// https://github.com/Kong/kubernetes-ingress-controller/issues/6807
		route, err := kongExpressionRouteFromHTTPRouteMatchWithPriority(matchWithPriority, supportRedirectPlugin)
		if err != nil {
			return []kongstate.Route{}, err
		}
		routes = append(routes, *route)
	}
	return routes, nil
}
