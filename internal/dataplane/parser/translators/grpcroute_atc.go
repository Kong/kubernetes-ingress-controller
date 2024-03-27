package translators

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

// GenerateKongExpressionRoutesFromGRPCRouteRule generates expression based kong routes
// from a single GRPCRouteRule.
func GenerateKongExpressionRoutesFromGRPCRouteRule(grpcroute *gatewayv1alpha2.GRPCRoute, ruleNumber int) []kongstate.Route {
	if ruleNumber >= len(grpcroute.Spec.Rules) {
		return nil
	}
	rule := grpcroute.Spec.Rules[ruleNumber]

	routes := make([]kongstate.Route, 0, len(rule.Matches))
	// gather the k8s object information and hostnames from the grpcroute
	ingressObjectInfo := util.FromK8sObject(grpcroute)

	// generate a route to match hostnames only if there is no match in the rule.
	if len(rule.Matches) == 0 {
		routeName := fmt.Sprintf(
			"grpcroute.%s.%s.%d.0",
			grpcroute.Namespace,
			grpcroute.Name,
			ruleNumber,
		)
		r := kongstate.Route{
			Ingress: ingressObjectInfo,
			Route: kong.Route{
				Name: kong.String(routeName),
			},
			ExpressionRoutes: true,
		}
		hostnames := getGRPCRouteHostnamesAsSliceOfStrings(grpcroute)
		// assign an empty match to generate matchers by only hostnames and annotations.
		matcher := generateMathcherFromGRPCMatch(gatewayv1alpha2.GRPCRouteMatch{}, hostnames, ingressObjectInfo.Annotations)
		atc.ApplyExpression(&r.Route, matcher, 1)
		return []kongstate.Route{r}
	}

	for matchNumber, match := range rule.Matches {
		routeName := fmt.Sprintf(
			"grpcroute.%s.%s.%d.%d",
			grpcroute.Namespace,
			grpcroute.Name,
			ruleNumber,
			matchNumber,
		)

		r := kongstate.Route{
			Ingress: ingressObjectInfo,
			Route: kong.Route{
				Name: kong.String(routeName),
			},
			ExpressionRoutes: true,
		}

		hostnames := getGRPCRouteHostnamesAsSliceOfStrings(grpcroute)
		matcher := generateMathcherFromGRPCMatch(match, hostnames, ingressObjectInfo.Annotations)

		atc.ApplyExpression(&r.Route, matcher, 1)
		routes = append(routes, r)
	}

	return routes
}

func generateMathcherFromGRPCMatch(match gatewayv1alpha2.GRPCRouteMatch, hostnames []string, metaAnnotations map[string]string) atc.Matcher {
	routeMatcher := atc.And()

	if match.Method != nil {
		methodMatcher := methodMatcherFromGRPCMethodMatch(match.Method)
		routeMatcher.And(methodMatcher)
	}

	if len(match.Headers) > 0 {
		headerMatcher := headerMatcherFromGRPCHeaderMatches(match.Headers)
		routeMatcher.And(headerMatcher)
	}

	if len(hostnames) > 0 {
		hostMatcher := hostMatcherFromHosts(hostnames)
		routeMatcher.And(hostMatcher)
	}

	snis, exist := annotations.ExtractSNIs(metaAnnotations)
	if exist && len(snis) > 0 {
		sniMatcher := sniMatcherFromSNIs(snis)
		routeMatcher.And(sniMatcher)
	}

	return routeMatcher
}

// methodMatcherFromGRPCMethodMatch translates ONE GRPC method match in GRPCRoute to ATC matcher.
func methodMatcherFromGRPCMethodMatch(methodMatch *gatewayv1alpha2.GRPCMethodMatch) atc.Matcher {
	matchType := gatewayv1alpha2.GRPCMethodMatchExact
	if methodMatch.Type != nil {
		matchType = *methodMatch.Type
	}

	switch matchType {
	case gatewayv1alpha2.GRPCMethodMatchExact:
		return methodMatcherFromGRPCExactMethodMatch(methodMatch.Service, methodMatch.Method)
	case gatewayv1alpha2.GRPCMethodMatchRegularExpression:
		return methodMatcherFromGRPCRegexMethodMatch(methodMatch.Service, methodMatch.Method)
	}

	return nil // should be unreachable
}

// methodMatcherFromGRPCExactMethodMatch translates exact GRPC method match to ATC matcher.
// reference: https://gateway-api.sigs.k8s.io/geps/gep-1016/?h=#method-matchers
func methodMatcherFromGRPCExactMethodMatch(service *string, method *string) atc.Matcher {
	if service == nil && method == nil {
		// should not happen, but we gernerate a catch-all matcher here.
		return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/")
	}
	if service != nil && method == nil {
		// Prefix /${SERVICE}/
		return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, fmt.Sprintf("/%s/", *service))
	}
	if service == nil && method != nil {
		// Suffix /${METHOD}
		return atc.NewPredicateHTTPPath(atc.OpSuffixMatch, fmt.Sprintf("/%s", *method))
	}
	// service and method are both specified
	return atc.NewPredicateHTTPPath(atc.OpEqual, fmt.Sprintf("/%s/%s", *service, *method))
}

// methodMatcherFromGRPCRegexMethodMatch translates regular expression GRPC method match to ATC matcher.
// reference: https://gateway-api.sigs.k8s.io/geps/gep-1016/?h=#type-regularexpression
func methodMatcherFromGRPCRegexMethodMatch(service *string, method *string) atc.Matcher {
	if service == nil && method == nil {
		return atc.NewPredicateHTTPPath(atc.OpPrefixMatch, "/")
	}
	// the regex to match service part and method part match any non-empty string if they are not specified.
	serviceRegex := ".+"
	methodRegex := ".+"

	if service != nil {
		serviceRegex = *service
	}
	if method != nil {
		methodRegex = *method
	}
	return atc.NewPredicateHTTPPath(atc.OpRegexMatch, fmt.Sprintf("^/%s/%s", serviceRegex, methodRegex))
}

func headerMatcherFromGRPCHeaderMatches(headerMatches []gatewayv1alpha2.GRPCHeaderMatch) atc.Matcher {
	// sort headerMatches by names to generate a stable output.
	sort.Slice(headerMatches, func(i, j int) bool {
		return string(headerMatches[i].Name) < string(headerMatches[j].Name)
	})

	matchers := make([]atc.Matcher, 0, len(headerMatches))
	for _, headerMatch := range headerMatches {
		httpHeaderMatch := gatewayv1.HTTPHeaderMatch{
			Type:  headerMatch.Type,
			Name:  gatewayv1.HTTPHeaderName(headerMatch.Name),
			Value: headerMatch.Value,
		}
		matchers = append(matchers, headerMatcherFromHTTPHeaderMatch(httpHeaderMatch))
	}
	return atc.And(matchers...)
}

func getGRPCRouteHostnamesAsSliceOfStrings(grpcroute *gatewayv1alpha2.GRPCRoute) []string {
	return lo.Map(grpcroute.Spec.Hostnames, func(h gatewayv1alpha2.Hostname, _ int) string {
		return string(h)
	})
}

// SplitGRPCRoute splits a GRPCRoute by hostname and match into multiple matches.
// Each split match contains at most 1 hostname, and 1 rule with 1 match.
func SplitGRPCRoute(grpcroute *gatewayv1alpha2.GRPCRoute) []SplitGRPCRouteMatch {
	splitMatches := []SplitGRPCRouteMatch{}
	splitGRPCRouteByMatch := func(hostname string) {
		for ruleIndex, rule := range grpcroute.Spec.Rules {
			// split out a match with only the hostname (non-empty only) when there are no matches in rule.
			if len(rule.Matches) == 0 {
				splitMatches = append(splitMatches, SplitGRPCRouteMatch{
					Source:     grpcroute,
					Hostname:   hostname,
					Match:      gatewayv1alpha2.GRPCRouteMatch{}, // empty grpcRoute match with ALL nil fields
					RuleIndex:  ruleIndex,
					MatchIndex: 0,
				})
			}
			for matchIndex, match := range rule.Matches {
				splitMatches = append(splitMatches, SplitGRPCRouteMatch{
					Source:     grpcroute,
					Hostname:   hostname,
					Match:      *match.DeepCopy(),
					RuleIndex:  ruleIndex,
					MatchIndex: matchIndex,
				})
			}
		}
	}

	// split GRPCRoute by matches if no hostname specified in spec.
	if len(grpcroute.Spec.Hostnames) == 0 {
		splitGRPCRouteByMatch("")
		return splitMatches
	}
	// split by hostname, then split further by rule and match.
	for _, hostname := range grpcroute.Spec.Hostnames {
		splitGRPCRouteByMatch(string(hostname))
	}
	return splitMatches
}

type GRPCRoutePriorityTraits struct {
	// PreciseHostname is set to true if the hostname is non-wildcard.
	PreciseHostname bool
	// HostnameLength is the length of hostname. Max 253.
	HostnameLength int
	// MethodMatchType is the type of method match (if exists).
	// preserve this field since the API specification has not provided priority of type of method match yet.
	// (In normal situation should be higher than length of service/method value).
	// related issue: https://github.com/kubernetes-sigs/gateway-api/issues/2216
	MethodMatchType gatewayv1alpha2.GRPCMethodMatchType
	// ServiceLength is the length of GRPC service name. Max 1024.
	ServiceLength int
	// MethodLength is the length of GRPC method name. Max 1024.
	MethodLength int
	// HeaderCount is the number of header matches in the match. Max 16.
	HeaderCount int
}

// CalculateGRCPRouteMatchPriorityTraits calculates the traits to decide
// priority based on the hostname and match split from source GRPCRoute.
// Specification of priority goes as follow:
// (The following comments are extracted from gateway API specification about GRPCRoute)
//
// Precedence MUST be given to the rule with the largest number of:
//   - Characters in a matching non-wildcard hostname.
//   - Characters in a matching hostname.
//   - Characters in a matching service.
//   - Characters in a matching method.
//   - Header matches.
//
// Method match type is preserved since the specification of GRPCRoute has not provided its priority yet.
func CalculateGRCPRouteMatchPriorityTraits(match SplitGRPCRouteMatch) GRPCRoutePriorityTraits {
	traits := GRPCRoutePriorityTraits{}

	// calculate traits from hostname.
	if len(match.Hostname) > 0 {
		traits.HostnameLength = len(match.Hostname)
		if !strings.HasPrefix(match.Hostname, "*") {
			traits.PreciseHostname = true
		}
	}
	// calculate traits from match.
	if match.Match.Method != nil {
		methodMatch := match.Match.Method
		if methodMatch.Type != nil {
			traits.MethodMatchType = *methodMatch.Type
		}
		if methodMatch.Service != nil {
			traits.ServiceLength = len(*methodMatch.Service)
		}
		if methodMatch.Method != nil {
			traits.MethodLength = len(*methodMatch.Method)
		}
	}
	// extract header count.
	traits.HeaderCount = len(match.Match.Headers)

	return traits
}

// EncodeToPriority turns GRPCRoute priority traits into the integer expressed priority.
//
//		   4                   3                   2                   1
//	 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0
//	+-+---------------+---------------------+---------------------+---------+----------------+
//	|P| host len      | GRPC service length | GRPC method length  |Header No| relative order |
//	+-+---------------+---------------------+---------------------+---------+----------------+
//
// Where:
// P: set to 1 if the hostname is non-wildcard.
// host len: length of hostname.
// GRPC service length: length of `Service` part in method match
// GRPC method length: length of `Method` part in method match
// Header No.: number of header matches.
// relative order: relative order of creation timestamp, namespace and name and internal rule/match order between different (split) GRPCRoutes.
//
// REVIEW: althogh not specified in official docs, do we need to assign a bit for GRPC method match type
// to assign higher priority for method match with `Exact` match?
func (t GRPCRoutePriorityTraits) EncodeToPriority() RoutePriorityType {
	const (
		// PreciseHostnameShiftBits assigns bit 43 for marking if the hostname is non-wildcard.
		PreciseHostnameShiftBits = 43
		// HostnameLengthShiftBits assigns bits 35-42 for the length of hostname.
		HostnameLengthShiftBits = 35
		// ServiceLengthShiftBits assigns bits 24-34 for the length of `Service` in method match.
		ServiceLengthShiftBits = 24
		// MethodLengthShiftBits assigns bits 13-23 for the length of `Method` in method match.
		MethodLengthShiftBits = 13
		// HeaderCountShiftBits assigns bits 8-12 for the number of header matches.
		HeaderCountShiftBits = 8
		// bits 0-7 are used for relative order of creation timestamp, namespace/name, and internal order of rules and matches.
		// the bits are calculated by sorting GRPCRoutes with the same priority calculated from the fields above
		// and start from all 1s, then decrease by one for each GRPCRoute.
	)

	var priority RoutePriorityType
	priority += ResourceKindBitsGRPCRoute << FromResourceKindPriorityShiftBits
	if t.PreciseHostname {
		priority += (1 << PreciseHostnameShiftBits)
	}
	priority += RoutePriorityType(t.HostnameLength) << HostnameLengthShiftBits
	priority += RoutePriorityType(t.ServiceLength) << ServiceLengthShiftBits
	priority += RoutePriorityType(t.MethodLength) << MethodLengthShiftBits
	priority += RoutePriorityType(t.HeaderCount) << HeaderCountShiftBits

	return priority
}

// SplitGRPCRouteMatch is the GRPCRouteMatch split by rule and match from the source GRPCRoute.
// RuleIndex and MatchIndex annotates the place of the match in the source GRPCRoute.
type SplitGRPCRouteMatch struct {
	Source     *gatewayv1alpha2.GRPCRoute
	Hostname   string
	Match      gatewayv1alpha2.GRPCRouteMatch
	RuleIndex  int
	MatchIndex int
}

type SplitGRPCRouteMatchToPriority struct {
	Match    SplitGRPCRouteMatch
	Priority uint64
}

// AssignRoutePriorityToSplitGRPCRouteMatches assigns priority to ALL split GRPCRoute matches
// that are split by hostnames and matches from GRPCRoutes listed from the cache.
// Firstly assign "fixed" bits by the following fields of the matches:
// hostname, GRPC method match, number of header matches.
// If ties exists in the first step, where multiple matches has the same priority
// calculated from the fields, we run a sort for the matches in the tie
// and assign the bits for "relative order" according to the sorting result of these matches.
func AssignRoutePriorityToSplitGRPCRouteMatches(
	logger logr.Logger,
	splitGRPCouteMatches []SplitGRPCRouteMatch,
) []SplitGRPCRouteMatchToPriority {
	priorityToSplitGRPCRouteMatches := map[RoutePriorityType][]SplitGRPCRouteMatch{}
	for _, match := range splitGRPCouteMatches {
		priority := CalculateGRCPRouteMatchPriorityTraits(match).EncodeToPriority()
		priorityToSplitGRPCRouteMatches[priority] = append(priorityToSplitGRPCRouteMatches[priority], match)
	}

	splitGRPCRoutesToPriority := make([]SplitGRPCRouteMatchToPriority, 0, len(splitGRPCouteMatches))

	// Bits 0-7 (8 bits) are assigned for relative order of GRPCRoutes.
	// If multiple GRPCRoutes are assigned to the same priority in the previous step,
	// sort them then starts with 2^8 -1 and decrease by one for each GRPCRoute;
	// If only one GRPCRoute occupies the priority, fill the relative order bits with all 1s.
	const relativeOrderAssignedBits = 8
	const defaultRelativeOrderPriorityBits RoutePriorityType = (1 << relativeOrderAssignedBits) - 1

	for priority, matches := range priorityToSplitGRPCRouteMatches {
		if len(matches) == 1 {
			splitGRPCRoutesToPriority = append(
				splitGRPCRoutesToPriority, SplitGRPCRouteMatchToPriority{
					Match:    matches[0],
					Priority: priority + defaultRelativeOrderPriorityBits,
				})
			continue
		}

		sort.SliceStable(matches, func(i, j int) bool {
			return compareSplitGRPCRouteMatchesRelativePriority(matches[i], matches[j])
		})

		for i, match := range matches {
			relativeOrderBits := defaultRelativeOrderPriorityBits - RoutePriorityType(i)
			// Although it is very unlikely that there are 2^8 = 256 GRPCRoutes
			// should be given priority by their relative order, here we limit the
			// relativeOrderBits to be at least 0.
			if relativeOrderBits <= 0 {
				relativeOrderBits = 0
			}
			splitGRPCRoutesToPriority = append(splitGRPCRoutesToPriority, SplitGRPCRouteMatchToPriority{
				Match:    match,
				Priority: priority + relativeOrderBits,
			})
		}

		// Just in case, log a very unlikely scenario where we have more than 2^8 routes with the same base
		// priority and we have no bit space for them to be deterministically ordered.
		if len(matches) > (1 << 8) {
			logger.Error(nil, "Too many GRPCRoute matches to be deterministically ordered", "grpcroute_number", len(matches))
		}
	}

	return splitGRPCRoutesToPriority
}

// compareSplitGRPCRouteMatchesRelativePriority compares the "relative order" of two matches split from GRPCRoutes.
// When it returns true, match1 will take precedence over match2 if priority came from host and match ties.
// The order should be decided by: (extracted from API specification in gateway API documents):
//
// If ties still exist across multiple Routes, matching precedence MUST be determined in order of the following criteria, continuing on ties:
// - The oldest Route based on creation timestamp.
// - The Route appearing first in alphabetical order by “{namespace}/{name}”.
// If ties still exist within the Route that has been given precedence,
// matching precedence MUST be granted to the first matching rule meeting the above criteria.
func compareSplitGRPCRouteMatchesRelativePriority(match1, match2 SplitGRPCRouteMatch) bool {
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

// KongExpressionRouteFromSplitGRPCRouteMatchWithPriority generates expression based
// Kong route from split GRPCRoute match which contains one or no hostname, and a GRPCRoute match,
// with its priority is beforehand.
func KongExpressionRouteFromSplitGRPCRouteMatchWithPriority(
	matchWithPriority SplitGRPCRouteMatchToPriority,
) kongstate.Route {
	grpcRoute := matchWithPriority.Match.Source
	tags := util.GenerateTagsForObject(grpcRoute)
	// since we split GRPCRoute by hostname, rule and match, we generate the route name in
	// grpcroute.<namespace>.<name>.<hostname>.<rule index>.<match index> format.
	hostname := matchWithPriority.Match.Hostname
	hostnameInRouteName := "_"
	if len(hostname) > 0 {
		hostnameInRouteName = strings.ReplaceAll(hostname, "*", "_")
	}
	routeName := fmt.Sprintf("grpcroute.%s.%s.%s.%d.%d",
		grpcRoute.Namespace,
		grpcRoute.Name,
		hostnameInRouteName,
		matchWithPriority.Match.RuleIndex,
		matchWithPriority.Match.MatchIndex,
	)

	r := kongstate.Route{
		Route: kong.Route{
			Name:         kong.String(routeName),
			PreserveHost: kong.Bool(true),
			Tags:         tags,
		},
		Ingress:          util.FromK8sObject(grpcRoute),
		ExpressionRoutes: true,
	}

	grpcMatch := matchWithPriority.Match.Match
	matcher := generateMathcherFromGRPCMatch(
		grpcMatch,
		[]string{hostname},
		grpcRoute.Annotations,
	)
	atc.ApplyExpression(&r.Route, matcher, matchWithPriority.Priority)
	if r.Expression == nil || len(*r.Expression) == 0 {
		r.Expression = kong.String(CatchAllHTTPExpression)
		r.Priority = &matchWithPriority.Priority
	}

	return r
}

// KongServiceNameFromSplitGRPCRouteMatch generates the name of translated Kong service
// from split GRPCRoute match with the source GRPCRoute and rule index.
func KongServiceNameFromSplitGRPCRouteMatch(match SplitGRPCRouteMatch) string {
	// fill hostname.
	hostname := "_"
	if len(match.Hostname) > 0 {
		hostname = strings.ReplaceAll(match.Hostname, "*", "_")
	}
	// pattern grpcroute.<namespace>.<name>.<hostname>.<rule index>.
	grpcRoute := match.Source
	return fmt.Sprintf("grpcroute.%s.%s.%s.%d",
		grpcRoute.Namespace,
		grpcRoute.Name,
		hostname,
		match.RuleIndex,
	)
}
