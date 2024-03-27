package translators

import (
	"fmt"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func getGRPCMatchDefaults() (
	map[gatewayv1alpha2.GRPCMethodMatchType]string,
	map[gatewayv1alpha2.GRPCMethodMatchType]string,
) {
	// Kong routes derived from a GRPCRoute use a path composed of the match's gRPC service and method
	// If either the service or method is omitted, there is a default regex determined by the match type
	// https://gateway-api.sigs.k8s.io/geps/gep-1016/#matcher-types describes the defaults

	// default path components for the GRPC service
	return map[gatewayv1alpha2.GRPCMethodMatchType]string{
			gatewayv1alpha2.GRPCMethodMatchType(""):          ".+",
			gatewayv1alpha2.GRPCMethodMatchExact:             ".+",
			gatewayv1alpha2.GRPCMethodMatchRegularExpression: ".+",
		},
		// default path components for the GRPC method
		map[gatewayv1alpha2.GRPCMethodMatchType]string{
			gatewayv1alpha2.GRPCMethodMatchType(""):          "",
			gatewayv1alpha2.GRPCMethodMatchExact:             "",
			gatewayv1alpha2.GRPCMethodMatchRegularExpression: ".+",
		}
}

func GenerateKongRoutesFromGRPCRouteRule(
	grpcroute *gatewayv1alpha2.GRPCRoute,
	ruleNumber int,
	prependRegexPrefix bool,
) []kongstate.Route {
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
				Name:      kong.String(routeName),
				Protocols: kong.StringSlice("grpc", "grpcs"),
			},
		}
		r.Hosts = getGRPCRouteHostnamesAsSliceOfStringPointers(grpcroute)
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
				Name:      kong.String(routeName),
				Protocols: kong.StringSlice("grpc", "grpcs"),
			},
		}

		if match.Method != nil {
			serviceMap, methodMap := getGRPCMatchDefaults()
			var method, service string
			matchMethod := match.Method.Method
			matchService := match.Method.Service
			var matchType gatewayv1alpha2.GRPCMethodMatchType
			if match.Method.Type == nil {
				matchType = gatewayv1alpha2.GRPCMethodMatchExact
			} else {
				matchType = *match.Method.Type
			}
			if matchMethod == nil {
				method = methodMap[matchType]
			} else {
				method = *matchMethod
			}
			if matchService == nil {
				service = serviceMap[matchType]
			} else {
				service = *matchService
			}
			// Kong prior to 3.0 does not accept paths starting with ~,
			// so we should only add the path regex prefix (~) only for Kong 3.0+.
			path := fmt.Sprintf("/%s/%s", service, method)
			if prependRegexPrefix {
				path = KongPathRegexPrefix + path
			}
			r.Paths = append(r.Paths, kong.String(path))
		}

		if len(grpcroute.Spec.Hostnames) > 0 {
			r.Hosts = getGRPCRouteHostnamesAsSliceOfStringPointers(grpcroute)
		}

		r.Headers = map[string][]string{}
		for _, hmatch := range match.Headers {
			name := string(hmatch.Name)
			r.Headers[name] = append(r.Headers[name], hmatch.Value)
		}

		routes = append(routes, r)
	}

	return routes
}

// -----------------------------------------------------------------------------
// Translate GRPCRoute - Utils
// -----------------------------------------------------------------------------

// getGRPCRouteHostnamesAsSliceOfStringPointers translates the hostnames defined
// in an GRPCRoute specification into a []*string slice, which is the type required
// by kong.Route{}.
func getGRPCRouteHostnamesAsSliceOfStringPointers(grpcroute *gatewayv1alpha2.GRPCRoute) []*string {
	return lo.Map(grpcroute.Spec.Hostnames, func(h gatewayv1.Hostname, _ int) *string {
		return lo.ToPtr(string(h))
	})
}
