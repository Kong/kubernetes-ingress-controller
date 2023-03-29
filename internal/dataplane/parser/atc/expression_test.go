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
			expression: `http.headers.x_kong_test == "test"`,
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
			expression: `(http.headers.x_kong_test == "test") && (http.host =^ ".konghq.com") && (tls.sni =^ ".konghq.com")`,
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
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exp := tc.matcher.Expression()
			require.Equal(t, tc.expression, exp)
		})
	}
}
