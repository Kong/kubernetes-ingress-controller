package subtranslator

import (
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/atc"
)

// ApplyExpressionToL4KongRoute convert route flavor from traditional to expressions
// against protocols, snis and dest ports.
func ApplyExpressionToL4KongRoute(r *kongstate.Route) {
	matchers := []atc.Matcher{}

	sniMatcher := sniMatcherFromSNIs(lo.Map(r.Route.SNIs, func(item *string, _ int) string { return *item }))
	matchers = append(matchers, sniMatcher)

	// TODO(rodman10): replace with helper function.
	portMatchers := make([]atc.Matcher, 0, len(r.Destinations))
	// (and (or sources) (or destinations))

	// Kong route sources and destinations support IP criteria, but Gateway API routes do not (Listeners apply to all IPs
	// assigned to a Gateway) and neither do our TCPIngress and UDPIngress CRs (we simply never added an IP field).
	// If we multiplex multiple Gateways (with different assigned IPs) onto a single Kong instance, we would need to add
	// IP criteria for full compliance. We already break this rule for HTTP Listeners, since Kong HTTP routes do not
	// support sources and destinations.
	//
	// Neither GWAPI Routes nor TCP/UDPIngress support sources either, so this only adds destinations.
	for _, dst := range r.Destinations {
		portMatcher, _ := atc.NewPredicate(atc.FieldNetDstPort, atc.OpEqual, atc.IntLiteral(*dst.Port))
		portMatchers = append(portMatchers, portMatcher)
	}
	matchers = append(matchers, atc.Or(portMatchers...))

	r.ExpressionRoutes = true
	atc.ApplyExpression(&r.Route, atc.And(matchers...), 1)
}
