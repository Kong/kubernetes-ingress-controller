package translators

import (
	"fmt"
	"sort"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func GenerateKongExpressionRoutesFromGRPCRouteRule(grpcroute *gatewayv1alpha2.GRPCRoute, ruleNumber int, rule gatewayv1alpha2.GRPCRouteRule) []kongstate.Route {
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
// REVIEW(naming): this function actually generates matcher to match HTTP path but not HTTP method. rename to pathMatcher...?
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
