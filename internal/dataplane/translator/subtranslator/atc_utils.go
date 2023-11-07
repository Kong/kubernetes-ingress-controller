package subtranslator

import (
	"strings"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/atc"
)

const (
	// FromResourceKindPriorityShiftBits is the highest 2 bits 45-44 used in priority field of Kong route
	// to note the kind of the resource from which the route is translated.
	// 11 - routes from Ingress.
	// 10 - routes from HTTPRoute.
	// 01 - routes from GRPCRoute.
	FromResourceKindPriorityShiftBits = 44
	// ResourceKindBitsIngress is the value of highest 2 bits for routes from ingresses.
	ResourceKindBitsIngress = 3
	// ResourceKindBitsHTTPRoute is the value of highest 2 bits for routes from HTTPRoutes.
	ResourceKindBitsHTTPRoute = 2
	// ResourceKindBitsGRPCRoute is the value of highest 2 bits for routes from GRPCRoutes.
	ResourceKindBitsGRPCRoute = 1
)

type (
	// RoutePriorityType is the type of priority field of Kong routes.
	RoutePriorityType = uint64
)

const (
	// CatchAllHTTPExpression is the expression to match all HTTP/HTTPS requests.
	// For rules with no matches and no hostnames in its parent HTTPRoute or GRPCRoute,
	// we need to generate a "catch-all" route for the rule:
	// https://github.com/Kong/kubernetes-ingress-controller/issues/4526
	// but Kong does not allow empty expression in expression router,
	// so we define this is we need a route to match all HTTP/HTTPS requests.
	CatchAllHTTPExpression = `(net.protocol == "http") || (net.protocol == "https")`
)

// -----------------------------------------------------------------------------
// Translator - common functions in translating expression(ATC) routes from multiple kinds of k8s objects.
// -----------------------------------------------------------------------------

// hostMatcherFromHosts translates hosts to ATC matcher that matches any of them.
// used in translating hostname matches in ingresses, HTTPRoutes, GRPCRoutes.
// the hostname format includes:
// - wildcard hosts, starting with exactly one *
// - precise hosts, otherwise.
func hostMatcherFromHosts(hosts []string) atc.Matcher {
	matchers := make([]atc.Matcher, 0, len(hosts))
	for _, host := range hosts {
		if !validHosts.MatchString(host) {
			continue
		}

		if strings.HasPrefix(host, "*") {
			// wildcard match on hosts (like *.foo.com), genreate a suffix match.
			matchers = append(matchers, atc.NewPrediacteHTTPHost(atc.OpSuffixMatch, strings.TrimPrefix(host, "*")))
		} else {
			// exact match on hosts, generate an exact match.
			matchers = append(matchers, atc.NewPrediacteHTTPHost(atc.OpEqual, host))
		}
	}
	return atc.Or(matchers...)
}
