package atc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateExpression(t *testing.T) {
	testCases := []struct {
		name       string
		matcher    Matcher
		expression string
	}{
		{
			name:       "simple predicate matching string field: HTTP path",
			matcher:    NewPredicateHTTPPath(OpPrefixMatch, "/foo/"),
			expression: `http.path ^= "/foo/"`,
		},
		{
			name:       "simple predicate matching HTTP header",
			matcher:    NewPredicateHTTPHeader("X-Kong-Test", OpEqual, "test"),
			expression: `http.header.x_kong_test == "test"`,
		},
		{
			name: "simple predicate with lower() transformer",
			matcher: NewPredicate(
				NewTransformerLower(FieldHTTPMethod),
				OpEqual,
				StringLiteral("get"),
			),
			expression: `lower(http.method) == "get"`,
		},
		{
			name: "multiple predicates connected by AND(&&)",
			matcher: And(
				NewPredicateHTTPHeader("X-Kong-Test", OpEqual, "test"),
				NewPrediacteHTTPHost(OpSuffixMatch, ".konghq.com"),
				NewPredicateTLSSNI(OpSuffixMatch, ".konghq.com"),
			),
			expression: `(http.header.x_kong_test == "test") && (http.host =^ ".konghq.com") && (tls.sni =^ ".konghq.com")`,
		},
		{
			name: "multiple predicates connected by OR(||)",
			matcher: Or(
				NewPredicateHTTPPath(OpEqual, "/foo"),
				NewPredicateHTTPPath(OpPrefixMatch, "/foo/"),
			),
			expression: `(http.path == "/foo") || (http.path ^= "/foo/")`,
		},
		{
			name: "multiple predicates connected by complex concatation of AND/OR",
			matcher: And(
				NewPrediacteHTTPHost(OpEqual, "test.konghq.com"),
				Or(
					NewPredicateNetProtocol(OpEqual, "http"),
					NewPredicateNetProtocol(OpEqual, "https"),
				),
			).And(NewPredicateHTTPPath(OpEqual, "/foo")),
			expression: `(http.host == "test.konghq.com") && ((net.protocol == "http") || (net.protocol == "https")) && (http.path == "/foo")`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.matcher.Expression()
			require.Equal(t, tc.expression, exp)
		})
	}
}
