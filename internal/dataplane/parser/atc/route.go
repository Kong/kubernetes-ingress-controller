package atc

import (
	"github.com/kong/go-kong/kong"
)

func ApplyExpression(r *kong.Route, m Matcher, priority int) {
	r.Expression = kong.String(m.Expression())
	r.Priority = kong.Int(priority)
}
