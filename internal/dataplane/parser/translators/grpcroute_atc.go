package translators

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

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

const (
	InternalRuleIndexAnnotationKey  = "internal-rule-index"
	InternalMatchIndexAnnotationKey = "internal-match-index"
)

func SplitGRPCRoute(grpcroute *gatewayv1alpha2.GRPCRoute) []*gatewayv1alpha2.GRPCRoute {
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

	splittedGRPCRoutes := []*gatewayv1alpha2.GRPCRoute{}
	for _, hostnamedGRPCRoute := range hostnamedGRPCRoutes {
		for i, rule := range hostnamedGRPCRoute.Spec.Rules {
			for j, match := range rule.Matches {
				splittedRoute := hostnamedGRPCRoute.DeepCopy()
				splittedRoute.Spec.Rules = []gatewayv1alpha2.GRPCRouteRule{
					{
						Matches:     []gatewayv1alpha2.GRPCRouteMatch{match},
						Filters:     rule.Filters,
						BackendRefs: rule.BackendRefs,
					},
				}
				if splittedRoute.Annotations == nil {
					splittedRoute.Annotations = map[string]string{}
				}
				splittedRoute.Annotations[InternalRuleIndexAnnotationKey] = strconv.Itoa(i)
				splittedRoute.Annotations[InternalMatchIndexAnnotationKey] = strconv.Itoa(j)
				splittedGRPCRoutes = append(splittedGRPCRoutes, splittedRoute)
			}
		}
	}
	return splittedGRPCRoutes
}

type GRPCRouteWithPriority struct {
	GRPCRoute *gatewayv1alpha2.GRPCRoute
	Priority  int
}

func AssignPrioritiesToSplittedGRPCRoutes(grpcRoutes []*gatewayv1alpha2.GRPCRoute) []GRPCRouteWithPriority {
	const (
		/*
			From gateway API specification on type GRPCRouteRule:

			Precedence MUST be given to the rule with the largest number of:

			- Characters in a matching non-wildcard hostname.
			- Characters in a matching hostname.
			- Characters in a matching service.
			- Characters in a matching method.
			- Header matches.
		*/
		// PreciseHostnameShiftBits assigns bit 51 for marking if the hostname is non-wildcard.
		PreciseHostnameShiftBits = 51
		// HostnameLengthShiftBits assigns bits 43-50 for the length of hostname.
		HostnameLengthShiftBits = 43
		// ServiceLengthShiftBits assigns bits 42-32 for the length of service field of GRPC method match.
		ServiceLengthShiftBits = 32
		// MethodLengthShiftBits assigns bits 31-21 for length of method field of GRPC method match.
		MethodLengthShiftBits = 21
		// HeaderNumberShiftBits assigns bits 20-16 for number of header matches.
		HeaderNumberShiftBits = 16
		// remaining bits 15-0 are used for relative order of creation timestamp, namespace/name, and internal order of rules and matches.
	)

	priorityToGRPCRoutes := map[int][]*gatewayv1alpha2.GRPCRoute{}
	for _, grpcRoute := range grpcRoutes {

		anns := grpcRoute.Annotations
		if anns == nil || anns[InternalRuleIndexAnnotationKey] == "" || anns[InternalMatchIndexAnnotationKey] == "" {
			continue
		}

		var priority int
		// assign priority bits for hostname.
		if len(grpcRoute.Spec.Hostnames) > 0 {
			hostname := grpcRoute.Spec.Hostnames[0]
			priority += len(hostname) << HostnameLengthShiftBits
			if !strings.HasPrefix(string(hostname), "*") {
				priority += (1 << PreciseHostnameShiftBits)
			}
		}

		if len(grpcRoute.Spec.Rules) > 0 && len(grpcRoute.Spec.Rules[0].Matches) > 0 {
			match := grpcRoute.Spec.Rules[0].Matches[0]
			// assign priority bits for GRPC method match.
			if match.Method != nil {
				if match.Method.Service != nil {
					priority += len(*match.Method.Service) << ServiceLengthShiftBits
				}
				if match.Method.Method != nil {
					priority += len(*match.Method.Method) << MethodLengthShiftBits
				}
			}
			priority += len(match.Headers) << HeaderNumberShiftBits
		}
		priorityToGRPCRoutes[priority] = append(priorityToGRPCRoutes[priority], grpcRoute)
	}

	grpcRoutesWithPriorities := make([]GRPCRouteWithPriority, 0, len(grpcRoutes))
	for priority, routes := range priorityToGRPCRoutes {
		if len(routes) == 1 {
			grpcRoutesWithPriorities = append(grpcRoutesWithPriorities, GRPCRouteWithPriority{
				GRPCRoute: routes[0],
				Priority:  priority + ((1 << 16) - 1),
			})
			continue
		}

		sort.Slice(routes, func(i, j int) bool {
			// compare by creation timestamp.
			if !routes[i].CreationTimestamp.Equal(&routes[j].CreationTimestamp) {
				return routes[i].CreationTimestamp.Before(&routes[j].CreationTimestamp)
			}
			// compare by namespace.
			if routes[i].Namespace != routes[j].Namespace {
				return routes[i].Namespace < routes[j].Namespace
			}
			// compare by name.
			if routes[i].Name != routes[j].Name {
				return routes[i].Name < routes[j].Name
			}
			// compare by internal rule order.
			ruleIndexi, _ := strconv.Atoi(routes[i].Annotations[InternalRuleIndexAnnotationKey])
			ruleIndexj, _ := strconv.Atoi(routes[j].Annotations[InternalRuleIndexAnnotationKey])
			if ruleIndexi != ruleIndexj {
				return ruleIndexi < ruleIndexj
			}
			// compare by match order.
			matchIndexi, _ := strconv.Atoi(routes[i].Annotations[InternalMatchIndexAnnotationKey])
			matchIndexj, _ := strconv.Atoi(routes[j].Annotations[InternalMatchIndexAnnotationKey])
			if matchIndexi != matchIndexj {
				return matchIndexi < matchIndexj
			}
			return i < j
		})

		relativeOrderPriority := ((1 << 16) - 1)
		for i, route := range routes {
			grpcRoutesWithPriorities = append(grpcRoutesWithPriorities, GRPCRouteWithPriority{
				GRPCRoute: route,
				Priority:  priority + relativeOrderPriority - i,
			})
		}
	}
	return grpcRoutesWithPriorities
}

func KongServiceNameFromGRPCRouteWithPriority(grpcRouteWithPriority GRPCRouteWithPriority) string {
	grpcRoute := grpcRouteWithPriority.GRPCRoute
	return fmt.Sprintf("grpcroute.%s.%s.%s",
		grpcRoute.Namespace,
		grpcRoute.Name,
		grpcRoute.Annotations[InternalRuleIndexAnnotationKey],
	)
}

func KongExpressionRouteFromGRPCRouteWithPriority(
	grpcRouteWithPriority GRPCRouteWithPriority,
) kongstate.Route {
	grpcRoute := grpcRouteWithPriority.GRPCRoute
	tags := util.GenerateTagsForObject(grpcRoute)
	routeName := fmt.Sprintf("grpcroute.%s.%s.%s.%s",
		grpcRoute.Namespace,
		grpcRoute.Name,
		grpcRoute.Annotations[InternalRuleIndexAnnotationKey],
		grpcRoute.Annotations[InternalMatchIndexAnnotationKey],
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

	hostnames := getGRPCRouteHostnamesAsSliceOfStrings(grpcRoute)
	matcher := generateMathcherFromGRPCMatch(grpcRoute.Spec.Rules[0].Matches[0], hostnames, grpcRoute.Annotations)
	atc.ApplyExpression(&r.Route, matcher, grpcRouteWithPriority.Priority)

	return r
}
