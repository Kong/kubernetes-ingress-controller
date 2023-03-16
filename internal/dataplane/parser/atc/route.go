package atc

import (
	"github.com/kong/go-kong/kong"
)

func ApplyExpression(r *kong.Route, m Matcher, priority int) {
	exp := m.Expression()
	r.Expression = &exp
	r.Priority = &priority
}
