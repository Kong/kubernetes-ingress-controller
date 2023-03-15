package atc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateExpression(t *testing.T) {
	p := NewPredicate(
		FieldHTTPHeader{headerName: "X-Kong-Test"},
		OpEqual,
		StringLiteral("test"),
	)
	exp := p.Expression()
	t.Log(exp)
	q := NewPredicate(
		FieldHTTPHost{},
		OpPrefixMatch,
		StringLiteral(".konghq.com"),
	)

	p2 := NewPredicate(
		TransformLower{
			inner: FieldHTTPPath{},
		},
		OpPrefixMatch,
		StringLiteral("/abc/def/"),
	)
	q2 := NewPredicate(
		TransformLower{
			inner: FieldHTTPPath{},
		},
		OpEqual,
		StringLiteral("/abc/def"),
	)

	require.Equal(t,
		`( http.header.x_kong_test == "test" ) && ( http.host ^= ".konghq.com" ) && ( ( lower(http.path) ^= "/abc/def/" ) || ( lower(http.path) == "/abc/def" ) )`,
		And(p, q, Or(p2, q2)).Expression())
}
