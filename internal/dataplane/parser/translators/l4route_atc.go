package translators

import (
	"github.com/samber/lo"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/atc"
)

// ApplyExpressionToL4KongRoute convert route flavor from traditional to expressions
// against protocols, snis and dest ports.
func ApplyExpressionToL4KongRoute(r *kongstate.Route) {
	matchers := []atc.Matcher{}

	sniMatcher := sniMatcherFromSNIs(lo.Map(r.Route.SNIs, func(item *string, _ int) string { return *item }))
	matchers = append(matchers, sniMatcher)

	// TODO(rodman10): replace with helper function.
	portMatchers := make([]atc.Matcher, 0, len(r.Destinations))
	for _, dst := range r.Destinations {
		portMatchers = append(portMatchers, atc.NewPredicate(atc.FieldNetDstPort, atc.OpEqual, atc.IntLiteral(*dst.Port)))
	}
	matchers = append(matchers, atc.Or(portMatchers...))

	r.ExpressionRoutes = true
	atc.ApplyExpression(&r.Route, atc.And(matchers...), 1)
}
