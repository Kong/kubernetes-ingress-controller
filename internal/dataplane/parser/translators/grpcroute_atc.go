package translators

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

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

	// override protocols from annotations.
	// Because Kong expression based router extracts net.protocol field from scheme of request,
	// GRPC over HTTP/2 requests could not be matched if protocol is set to grpc/grpcs since protocol could only be http or https.
	// So we do not AND a protocol matcher if no protocol is specified in annotations.
	protocols := annotations.ExtractProtocolNames(metaAnnotations)
	if len(protocols) > 0 {
		protocolMatcher := protocolMatcherFromProtocols(protocols)
		routeMatcher.And(protocolMatcher)
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
		httpHeaderMatch := gatewayv1beta1.HTTPHeaderMatch{
			Type:  headerMatch.Type,
			Name:  gatewayv1beta1.HTTPHeaderName(headerMatch.Name),
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

// SplitGRPCRoute splits a GRPCRoute by hostname and match into multiple GRPCRoutes.
// Each split GRPCRoute contains at most 1 hostname, and 1 rule with 1 match.
func SplitGRPCRoute(grpcroute *gatewayv1alpha2.GRPCRoute) []*gatewayv1alpha2.GRPCRoute {
	// split GRPCRoute by hostname.
	hostnamedGRPCRoutes := make([]*gatewayv1alpha2.GRPCRoute, 0, len(grpcroute.Spec.Hostnames))
	if len(grpcroute.Spec.Hostnames) == 0 {
		hostnamedGRPCRoute := grpcroute.DeepCopy()
		hostnamedGRPCRoutes = append(hostnamedGRPCRoutes, hostnamedGRPCRoute)
	} else {
		for _, hostname := range grpcroute.Spec.Hostnames {
			hostnamedGRPCRoute := grpcroute.DeepCopy()
			hostnamedGRPCRoute.Spec.Hostnames = []gatewayv1beta1.Hostname{hostname}
			hostnamedGRPCRoutes = append(hostnamedGRPCRoutes, hostnamedGRPCRoute)
		}
	}

	splitGRPCRoutes := []*gatewayv1alpha2.GRPCRoute{}
	for _, hostnamedGRPCRoute := range hostnamedGRPCRoutes {
		// split the GRPCRoutes already split once by hostnames
		// into GRPCRoutes with one match for each.
		for i, rule := range hostnamedGRPCRoute.Spec.Rules {
			for j, match := range rule.Matches {
				splitRoute := hostnamedGRPCRoute.DeepCopy()
				splitRoute.Spec.Rules = []gatewayv1alpha2.GRPCRouteRule{
					{
						Matches:     []gatewayv1alpha2.GRPCRouteMatch{match},
						Filters:     rule.Filters,
						BackendRefs: rule.BackendRefs,
					},
				}
				if splitRoute.Annotations == nil {
					splitRoute.Annotations = map[string]string{}
				}
				splitRoute.Annotations[InternalRuleIndexAnnotationKey] = strconv.Itoa(i)
				splitRoute.Annotations[InternalMatchIndexAnnotationKey] = strconv.Itoa(j)
				splitGRPCRoutes = append(splitGRPCRoutes, splitRoute)
			}
		}
	}
	return splitGRPCRoutes
}

type GRPCRoutePriorityTraits struct {
	// PreciseHostname is set to true if the hostname is non-wildcard.
	PreciseHostname bool
	// HostnameLength is the length of hostname. Max 253.
	HostnameLength int
	// ServiceLength is the length of GRPC service name. Max 1024.
	ServiceLength int
	// MethodLength is the length of GRPC method name. Max 1024.
	MethodLength int
	// HeaderCount is the number of header matches in the match. Max 16.
	HeaderCount int
}

// CalculateSplitGRCPRoutePriorityTraits calculates the traits to decide
// priority based on the spec of the GRPCRoute that is already split by hostname and match.
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
// REVIEW: do we need to include GRPC method match type (`Exact`,`RegularExpression` or others) into traits?
func CalculateSplitGRCPRoutePriorityTraits(grpcRoute *gatewayv1alpha2.GRPCRoute) GRPCRoutePriorityTraits {
	traits := GRPCRoutePriorityTraits{}

	// calculate traits from hostname.
	// The split GRPCRoute has at most one hostname in its spec.
	if len(grpcRoute.Spec.Hostnames) > 0 {
		hostname := grpcRoute.Spec.Hostnames[0]
		if !strings.HasPrefix(string(hostname), "*") {
			traits.PreciseHostname = true
		}
		traits.HostnameLength = len(hostname)
	}
	// calculate traits from match.
	// The split GRPCRoute has only 1 rule including 1 match.
	if len(grpcRoute.Spec.Rules) > 0 && len(grpcRoute.Spec.Rules[0].Matches) > 0 {
		match := grpcRoute.Spec.Rules[0].Matches[0]
		// extract length of GRPC service and method.
		if match.Method != nil {
			if match.Method.Service != nil {
				traits.ServiceLength = len(*match.Method.Service)
			}
			if match.Method.Method != nil {
				traits.MethodLength = len(*match.Method.Method)
			}
		}
		// extract header count.
		traits.HeaderCount = len(match.Headers)
	}
	return traits
}

// EncodeToPriority turns GRPCRoute priority traits into the integer expressed priority.
//
//					   4                   3                   2                   1
//	 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0
//	+-+---------------+---------------------+---------------------+---------+---------------------------+
//	|P| host len      | GRPC service length | GRPC method length  |Header No| relative order            |
//	+-+---------------+---------------------+---------------------+---------+---------------------------+
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
func (t GRPCRoutePriorityTraits) EncodeToPriority() int {
	const (
		// PreciseHostnameShiftBits assigns bit 49 for marking if the hostname is non-wildcard.
		PreciseHostnameShiftBits = 49
		// HostnameLengthShiftBits assigns bits 41-48 for the length of hostname.
		HostnameLengthShiftBits = 41
		// ServiceLengthShiftBits assigns bits 30-40 for the length of `Service` in method match.
		ServiceLengthShiftBits = 30
		// MethodLengthShiftBits assigns bits 19-29 for the length of `Method` in method match.
		MethodLengthShiftBits = 19
		// HeaderCountShiftBits assigns bits 14-18 for the number of header matches.
		HeaderCountShiftBits = 14
		// bits 0-13 are used for relative order of creation timestamp, namespace/name, and internal order of rules and matches.
		// the bits are calculated by sorting GRPCRoutes with the same priority calculated from the fields above
		// and start from all 1s, then decrease by one for each GRPCRoute.
	)

	var priority int
	priority += ResourceKindBitsGRPCRoute << FromResourceKindPriorityShiftBits
	if t.PreciseHostname {
		priority += (1 << PreciseHostnameShiftBits)
	}
	priority += t.HostnameLength << HostnameLengthShiftBits
	priority += t.ServiceLength << ServiceLengthShiftBits
	priority += t.MethodLength << MethodLengthShiftBits
	priority += t.HeaderCount << HeaderCountShiftBits

	return priority
}

type SplitGRPCRouteToKongRoutePriority struct {
	GRPCRoute *gatewayv1alpha2.GRPCRoute
	Priority  int
}

// AssignRoutePriorityToSplitGRPCRoutes assigns priority to ALL split GRPCRoutes
// that are split by hostnames and matches from GRPCRoutes listed from the cache.
// Firstly assign "fixed" bits by the following fields of the GRPCRoute:
// hostname, GRPC method match, number of header matches.
// If ties exists in the first step, where multiple GRPCRoutes has the same priority
// calculated from the fields, we run a sort for the GRPCRoutes in the tie
// and assign the bits for "relative order" according to the sorting result of these GRPCRoutes.
func AssignRoutePriorityToSplitGRPCRoutes(
	logger logr.Logger,
	splitGRPCoutes []*gatewayv1alpha2.GRPCRoute,
) []SplitGRPCRouteToKongRoutePriority {
	priorityToSplitGRPCRoutes := map[int][]*gatewayv1alpha2.GRPCRoute{}
	for _, grpcRoute := range splitGRPCoutes {
		// skip if GRPCRoute does not contain the annotation, because this means the GRPCRoute is not a split one.
		anns := grpcRoute.Annotations
		if anns == nil || anns[InternalRuleIndexAnnotationKey] == "" || anns[InternalMatchIndexAnnotationKey] == "" {
			continue
		}

		priority := CalculateSplitGRCPRoutePriorityTraits(grpcRoute).EncodeToPriority()
		priorityToSplitGRPCRoutes[priority] = append(priorityToSplitGRPCRoutes[priority], grpcRoute)
	}

	splitGRPCRoutesToPriority := make([]SplitGRPCRouteToKongRoutePriority, 0, len(splitGRPCoutes))

	// Bits 0-13 (14 bits) are assigned for relative order of GRPCRoutes.
	// If multiple GRPCRoutes are assigned to the same priority in the previous step,
	// sort them then starts with 2^14 -1 and decrease by one for each GRPCRoute;
	// If only one GRPCRoute occupies the priority, fill the relative order bits with all 1s.
	const relativeOrderAssignedBits = 14
	const defaultRelativeOrderPriorityBits = (1 << relativeOrderAssignedBits) - 1

	for priority, routes := range priorityToSplitGRPCRoutes {
		if len(routes) == 1 {
			splitGRPCRoutesToPriority = append(
				splitGRPCRoutesToPriority, SplitGRPCRouteToKongRoutePriority{
					GRPCRoute: routes[0],
					Priority:  priority + defaultRelativeOrderPriorityBits,
				})
			continue
		}

		sort.SliceStable(routes, func(i, j int) bool {
			return compareSplitGRPCRoutesRelativePriority(routes[i], routes[j])
		})

		for i, route := range routes {
			relativeOrderBits := defaultRelativeOrderPriorityBits - i
			// Although it is very unlikely that there are 2^14 = 16384 GRPCRoutes
			// should be given priority by their relative order, here we limit the
			// relativeOrderBits to be at least 0.
			if relativeOrderBits <= 0 {
				relativeOrderBits = 0
			}
			splitGRPCRoutesToPriority = append(splitGRPCRoutesToPriority, SplitGRPCRouteToKongRoutePriority{
				GRPCRoute: route,
				Priority:  priority + relativeOrderBits,
			})
		}

		// Just in case, log a very unlikely scenario where we have more than 2^14 routes with the same base
		// priority and we have no bit space for them to be deterministically ordered.
		if len(routes) > (1 << 14) {
			logger.V(util.WarnLevel).Info("Too many GRPCRoutes to be deterministically ordered", "grpcroute_number", len(routes))
		}
	}

	return splitGRPCRoutesToPriority
}

// compareSplitGRPCRoutesRelativePriority compares the "relative order" of two split GRPCRoutes.
// When it returns true, route1 will take precedence over route2 if all fields in spec ties.
// The order should be decided by: (extracted from API specification in gateway API documents):
//
// If ties still exist across multiple Routes, matching precedence MUST be determined in order of the following criteria, continuing on ties:
// - The oldest Route based on creation timestamp.
// - The Route appearing first in alphabetical order by “{namespace}/{name}”.
// If ties still exist within the Route that has been given precedence,
// matching precedence MUST be granted to the first matching rule meeting the above criteria.
func compareSplitGRPCRoutesRelativePriority(route1, route2 *gatewayv1alpha2.GRPCRoute) bool {
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

// GenerateKongExpressionRouteFromSplitGRPCRouteWithPriority generates expression based
// Kong route from split GRPCRoute which contains at most one hostname and one rule with one match,
// and its priority is calculated beforehand.
func GenerateKongExpressionRouteFromSplitGRPCRouteWithPriority(
	grpcRouteWithPriority SplitGRPCRouteToKongRoutePriority,
) kongstate.Route {
	grpcRoute := grpcRouteWithPriority.GRPCRoute
	tags := util.GenerateTagsForObject(grpcRoute)
	// since we split HTTPRoutes by hostname, rule and match, we generate the route name in
	// grpcroute.<namespace>.<name>.<hostname>.<rule index>.<match index> format.
	hostnameInRouteName := "_"
	if len(grpcRoute.Spec.Hostnames) > 0 {
		hostnameInRouteName = string(grpcRoute.Spec.Hostnames[0])
		hostnameInRouteName = strings.ReplaceAll(hostnameInRouteName, "*", "_")
	}
	routeName := fmt.Sprintf("grpcroute.%s.%s.%s.%s.%s",
		grpcRoute.Namespace,
		grpcRoute.Name,
		hostnameInRouteName,
		grpcRoute.Annotations[InternalRuleIndexAnnotationKey],
		grpcRoute.Annotations[InternalMatchIndexAnnotationKey],
	)

	r := kongstate.Route{
		Route: kong.Route{
			Name:         kong.String(routeName),
			PreserveHost: kong.Bool(true),
			// TODO: add `StripPath: false` here?
			Tags: tags,
		},
		Ingress:          util.FromK8sObject(grpcRoute),
		ExpressionRoutes: true,
	}

	if len(grpcRoute.Spec.Rules) > 0 && len(grpcRoute.Spec.Rules[0].Matches) > 0 {
		match := grpcRoute.Spec.Rules[0].Matches[0]
		matcher := generateMathcherFromGRPCMatch(
			match,
			getGRPCRouteHostnamesAsSliceOfStrings(grpcRoute),
			grpcRoute.Annotations,
		)
		atc.ApplyExpression(&r.Route, matcher, grpcRouteWithPriority.Priority)
	}
	return r
}

// KongServiceNameFromSplitGRPCRoute generates the name of translated Kong service
// from split GRPCRoute that contains at most 1 hostname and 1 rule.
func KongServiceNameFromSplitGRPCRoute(splitGRPCRoute *gatewayv1alpha2.GRPCRoute) string {
	hostname := "_"
	if len(splitGRPCRoute.Spec.Hostnames) > 0 {
		hostname = string(splitGRPCRoute.Spec.Hostnames[0])
		hostname = strings.ReplaceAll(hostname, "*", "_")
	}
	// pattern grpcroute.<namespace>.<name>.<hostname>.<rule index>.
	return fmt.Sprintf("grpcroute.%s.%s.%s.%s",
		splitGRPCRoute.Namespace,
		splitGRPCRoute.Name,
		hostname,
		splitGRPCRoute.Annotations[InternalRuleIndexAnnotationKey],
	)
}
