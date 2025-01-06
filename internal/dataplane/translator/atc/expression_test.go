package atc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mustNewPredicate(t *testing.T, lhs LHS, op BinaryOperator, rhs Literal) Predicate {
	t.Helper()

	predicate, err := NewPredicate(lhs, op, rhs)
	require.NoError(t, err)
	return predicate
}

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
			expression: `http.headers.x_kong_test == "test"`,
		},
		{
			name: "simple predicate with lower() transformer",
			matcher: mustNewPredicate(
				t,
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
			expression: `(http.headers.x_kong_test == "test") && (http.host =^ ".konghq.com") && (tls.sni =^ ".konghq.com")`,
		},
		{
			name: "HTTP query match predicates",
			matcher: And(
				NewPredicateHTTPQuery("foo", OpEqual, "1"),
				NewPredicateHTTPQuery("bar", OpEqual, "2"),
			),
			expression: `(http.queries.foo == "1") && (http.queries.bar == "2")`,
		},
		{
			name: "HTTP path segments match",
			matcher: And(
				// matches /api/namespaces/*/services/** .
				NewPredicateHTTPPathSegmentLength(OpGreaterThan, 4),
				NewPredicateHTTPPathSegmentInterval(0, 1, OpEqual, "api/namespaces"),
				NewPredicateHTTPPathSingleSegment(3, OpEqual, "services"),
			),
			expression: `(http.path.segments.len > 4) && (http.path.segments.0_1 == "api/namespaces") && (http.path.segments.3 == "services")`,
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
			name: "negate a single predicate using NOT",
			matcher: Not(
				NewPredicateNetProtocol(OpEqual, "http"),
			),
			expression: `!(net.protocol == "http")`,
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
		{
			name: "NOT matcher in submatchers of AND/OR matcher",
			matcher: And(
				Not(
					Or(
						NewPredicateHTTPMethod(OpEqual, "CONNECT"),
						NewPredicateHTTPMethod(OpEqual, "OPTIONS"),
					),
				),
				NewPredicateHTTPPath(OpEqual, "/foo/bar"),
			),
			expression: `(!((http.method == "CONNECT") || (http.method == "OPTIONS"))) && (http.path == "/foo/bar")`,
		},
		{
			name: "negate complex matcher with AND/OR using NOT",
			matcher: Not(
				Or(
					And(
						NewPrediacteHTTPHost(OpEqual, "foo.konghq.com"),
						NewPredicateHTTPPath(OpEqual, "/index.html"),
					),
					And(
						NewPrediacteHTTPHost(OpEqual, "bar.konghq.com"),
						NewPredicateHTTPPath(OpEqual, "/index"),
					),
				),
			),
			expression: `!(((http.host == "foo.konghq.com") && (http.path == "/index.html")) || ((http.host == "bar.konghq.com") && (http.path == "/index")))`,
		},
		{
			name:       "empty Or matcher",
			matcher:    Or(),
			expression: "",
		},
		{
			name: "single element And matcher",
			matcher: And(
				NewPrediacteHTTPHost(OpSuffixMatch, ".konghq.com"),
			),
			expression: `http.host =^ ".konghq.com"`,
		},
		{
			name: "call Or() with nils",
			matcher: Or(
				nil,
				NewPredicateHTTPMethod(OpEqual, "GET"),
				NewPredicateHTTPMethod(OpEqual, "POST"),
			).Or(nil).Or(NewPredicateHTTPMethod(OpEqual, "DELETE")),
			expression: `(http.method == "GET") || (http.method == "POST") || (http.method == "DELETE")`,
		},
		{
			name: "call And() with nils",
			matcher: And(
				NewPredicateHTTPHeader("X-Header-1", OpEqual, "v1"),
				nil,
			).And(NewPredicateHTTPHeader("X-Header-2", OpEqual, "v2")).And(nil),
			expression: `(http.headers.x_header_1 == "v1") && (http.headers.x_header_2 == "v2")`,
		},
		{
			name: "empty expression in Or",
			matcher: Or(
				And(),
				NewPredicateHTTPPath(OpEqual, "/foo"),
			),
			expression: `http.path == "/foo"`,
		},
		{
			name: "empty expression in And",
			matcher: And(
				And(),
				NewPredicateHTTPPath(OpEqual, "/foo"),
				NewPredicateHTTPMethod(OpEqual, "GET"),
			),
			expression: `(http.path == "/foo") && (http.method == "GET")`,
		},
		{
			name:       "nil Or",
			matcher:    (*OrMatcher)(nil),
			expression: "",
		},
		{
			name:       "nil Not",
			matcher:    (*NotMatcher)(nil),
			expression: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.matcher.Expression()
			require.Equal(t, tc.expression, exp)
		})
	}
}
