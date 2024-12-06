package subtranslator

import (
	"fmt"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

func TestGenerateKongExpressionRoutesFromHTTPRouteMatches(t *testing.T) {
	testCases := []struct {
		name              string
		routeName         string
		matches           []gatewayapi.HTTPRouteMatch
		filters           []gatewayapi.HTTPRouteFilter
		ingressObjectInfo util.K8sObjectInfo
		hostnames         []string
		tags              []string
		expectedRoutes    []kongstate.Route
		expectedError     error
	}{
		{
			name:              "no hostnames and no matches",
			routeName:         "empty_route.default.0.0",
			ingressObjectInfo: util.K8sObjectInfo{},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("empty_route.default.0.0"),
						PreserveHost: kong.Bool(true),
						StripPath:    kong.Bool(false),
						Expression:   kong.String(CatchAllHTTPExpression),
					},
					ExpressionRoutes: true,
				},
			},
		},
		{
			name:              "no matches but have hostnames",
			routeName:         "host_only.default.0.0",
			ingressObjectInfo: util.K8sObjectInfo{},
			hostnames:         []string{"foo.com", "*.bar.com"},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("host_only.default.0.0"),
						PreserveHost: kong.Bool(true),
						StripPath:    kong.Bool(false),
						Expression:   kong.String(`(http.host == "foo.com") || (http.host =^ ".bar.com")`),
						Priority:     kong.Uint64(1),
					},
					ExpressionRoutes: true,
				},
			},
		},
		{
			name:              "single prefix path match",
			routeName:         "prefix_path_match.defualt.0.0",
			ingressObjectInfo: util.K8sObjectInfo{},
			matches: []gatewayapi.HTTPRouteMatch{
				builder.NewHTTPRouteMatch().WithPathPrefix("/prefix").Build(),
			},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("prefix_path_match.defualt.0.0"),
						PreserveHost: kong.Bool(true),
						StripPath:    kong.Bool(false),
						Expression:   kong.String(`(http.path == "/prefix") || (http.path ^= "/prefix/")`),
						Priority:     kong.Uint64(1),
					},
					ExpressionRoutes: true,
				},
			},
		},
		{
			name:              "multiple matches without filters",
			routeName:         "multiple_matches.default.0.0",
			ingressObjectInfo: util.K8sObjectInfo{},
			matches: []gatewayapi.HTTPRouteMatch{
				builder.NewHTTPRouteMatch().WithPathPrefix("/prefix").Build(),
				builder.NewHTTPRouteMatch().WithPathExact("/exact").WithMethod(gatewayapi.HTTPMethodGet).Build(),
			},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("multiple_matches.default.0.0"),
						PreserveHost: kong.Bool(true),
						StripPath:    kong.Bool(false),
						Expression:   kong.String(`((http.path == "/prefix") || (http.path ^= "/prefix/")) || ((http.path == "/exact") && (http.method == "GET"))`),
						Priority:     kong.Uint64(1),
					},
					ExpressionRoutes: true,
				},
			},
		},
		{
			name:              "multiple matches with request redirect filter",
			routeName:         "request_redirect.default.0.0",
			ingressObjectInfo: util.K8sObjectInfo{},
			matches: []gatewayapi.HTTPRouteMatch{
				builder.NewHTTPRouteMatch().WithPathExact("/exact/0").Build(),
				builder.NewHTTPRouteMatch().WithPathExact("/exact/1").Build(),
			},
			filters: []gatewayapi.HTTPRouteFilter{
				builder.NewHTTPRouteRequestRedirectFilter().
					WithRequestRedirectScheme("http").
					WithRequestRedirectHost("a.foo.com").
					WithRequestRedirectStatusCode(301).
					Build(),
			},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("request_redirect.default.0.0"),
						PreserveHost: kong.Bool(true),
						StripPath:    kong.Bool(false),
						Expression:   kong.String(`http.path == "/exact/0"`),
						Priority:     kong.Uint64(1),
					},
					Plugins: []kong.Plugin{
						{
							Name: kong.String("request-termination"),
							Config: kong.Configuration{
								"status_code": kong.Int(301),
							},
						},
						{
							Name: kong.String("response-transformer"),
							Config: kong.Configuration{
								"add": TransformerPluginConfig{
									Headers: []string{`Location: http://a.foo.com:80/exact/0`},
								},
							},
						},
					},
					ExpressionRoutes: true,
				},
				{
					Route: kong.Route{
						Name:         kong.String("request_redirect.default.0.0"),
						PreserveHost: kong.Bool(true),
						StripPath:    kong.Bool(false),
						Expression:   kong.String(`http.path == "/exact/1"`),
						Priority:     kong.Uint64(1),
					},
					Plugins: []kong.Plugin{
						{
							Name: kong.String("request-termination"),
							Config: kong.Configuration{
								"status_code": kong.Int(301),
							},
						},
						{
							Name: kong.String("response-transformer"),
							Config: kong.Configuration{
								"add": TransformerPluginConfig{
									Headers: []string{`Location: http://a.foo.com:80/exact/1`},
								},
							},
						},
					},
					ExpressionRoutes: true,
				},
			},
		},
		{
			name:              "multiple matches with request header transformer filter",
			routeName:         "request_header_mod.default.0.0",
			ingressObjectInfo: util.K8sObjectInfo{},
			matches: []gatewayapi.HTTPRouteMatch{
				builder.NewHTTPRouteMatch().WithPathExact("/exact/0").Build(),
				builder.NewHTTPRouteMatch().WithPathRegex("/regex/[a-z]+").Build(),
			},
			filters: []gatewayapi.HTTPRouteFilter{
				builder.NewHTTPRouteRequestHeaderModifierFilter().WithRequestHeaderAdd([]gatewayapi.HTTPHeader{
					{Name: "foo", Value: "bar"},
				}).Build(),
			},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("request_header_mod.default.0.0"),
						PreserveHost: kong.Bool(true),
						StripPath:    kong.Bool(false),
						Expression:   kong.String(`(http.path == "/exact/0") || (http.path ~ "^/regex/[a-z]+")`),
						Priority:     kong.Uint64(1),
					},
					Plugins: []kong.Plugin{
						{
							Name: kong.String("request-transformer"),
							Config: kong.Configuration{
								"append": TransformerPluginConfig{
									Headers: []string{"foo:bar"},
								},
							},
						},
					},
					ExpressionRoutes: true,
				},
			},
		},
		{
			name:      "routes with annotations to set protocols and SNIs",
			routeName: "annotations_protocol_sni.default.0.0",
			ingressObjectInfo: util.K8sObjectInfo{
				Namespace: "default",
				Name:      "httproute-annotations",
				Annotations: map[string]string{
					"konghq.com/protocols": "https",
					"konghq.com/snis":      "a.foo.com",
				},
			},
			hostnames: []string{"a.foo.com"},
			matches: []gatewayapi.HTTPRouteMatch{
				builder.NewHTTPRouteMatch().WithPathPrefix("/prefix/0/").Build(),
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Namespace: "default",
						Name:      "httproute-annotations",
						Annotations: map[string]string{
							"konghq.com/protocols": "https",
							"konghq.com/snis":      "a.foo.com",
						},
					},
					Route: kong.Route{
						Name:         kong.String("annotations_protocol_sni.default.0.0"),
						PreserveHost: kong.Bool(true),
						StripPath:    kong.Bool(false),
						Expression:   kong.String(`((http.path == "/prefix/0") || (http.path ^= "/prefix/0/")) && (http.host == "a.foo.com") && (tls.sni == "a.foo.com")`),
						Priority:     kong.Uint64(1),
					},
					ExpressionRoutes: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			routes, err := GenerateKongExpressionRoutesFromHTTPRouteMatches(
				KongRouteTranslation{
					Name:    tc.routeName,
					Matches: tc.matches,
					Filters: tc.filters,
				},
				tc.ingressObjectInfo,
				tc.hostnames,
				kong.StringSlice(tc.tags...),
			)

			if tc.expectedError != nil {
				require.ErrorAs(t, err, &tc.expectedError)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedRoutes, routes)
		})
	}
}

func TestGenerateMatcherFromHTTPRouteMatch(t *testing.T) {
	testCases := []struct {
		name       string
		match      gatewayapi.HTTPRouteMatch
		expression string
	}{
		{
			name:       "empty prefix path match",
			match:      builder.NewHTTPRouteMatch().WithPathPrefix("/").Build(),
			expression: `http.path ^= "/"`,
		},
		{
			name:       "simple non-empty prefix path match",
			match:      builder.NewHTTPRouteMatch().WithPathPrefix("/prefix/0").Build(),
			expression: `(http.path == "/prefix/0") || (http.path ^= "/prefix/0/")`,
		},
		{
			name:       "simple exact path match",
			match:      builder.NewHTTPRouteMatch().WithPathExact("/exact/0/").Build(),
			expression: `http.path == "/exact/0/"`,
		},
		{
			name:       "simple regex match",
			match:      builder.NewHTTPRouteMatch().WithPathRegex("/regex/\\d{1,3}").Build(),
			expression: `http.path ~ "^/regex/\\d{1,3}"`,
		},
		{
			name: "exact path and method and a single header",
			match: builder.NewHTTPRouteMatch().WithPathExact("/exact/0").
				WithMethod(gatewayapi.HTTPMethodGet).
				WithHeader("foo", "bar").
				Build(),
			expression: `(http.path == "/exact/0") && (http.headers.foo == "bar") && (http.method == "GET")`,
		},
		{
			name: "prefix path match and multiple headers",
			match: builder.NewHTTPRouteMatch().WithPathPrefix("/prefix/0").
				WithHeader("X-Foo", "Bar").
				WithHeaderRegex("Hash", "[0-9A-Fa-f]{32}").
				Build(),
			expression: `((http.path == "/prefix/0") || (http.path ^= "/prefix/0/")) && ((http.headers.hash ~ "[0-9A-Fa-f]{32}") && (http.headers.x_foo == "Bar"))`,
		},
		{
			name: "prefix path match and multiple query params",
			match: builder.NewHTTPRouteMatch().WithPathPrefix("/prefix/0").
				WithQueryParam("foo", "bar").
				WithQueryParamRegex("id", "[0-9a-z-]+").
				Build(),
			expression: `((http.path == "/prefix/0") || (http.path ^= "/prefix/0/")) && ((http.queries.foo == "bar") && (http.queries.id ~ "[0-9a-z-]+"))`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expression, generateMatcherFromHTTPRouteMatch(tc.match).Expression())
		})
	}
}

func TestCalculateHTTPRoutePriorityTraits(t *testing.T) {
	testCases := []struct {
		name           string
		match          SplitHTTPRouteMatch
		expectedTraits HTTPRoutePriorityTraits
	}{
		{
			name: "precise hostname and exact path",
			match: SplitHTTPRouteMatch{
				Source: &gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "precise-hostname-exact-path",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Hostnames: []gatewayapi.Hostname{"foo.com"},
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
				Hostname: "foo.com",
				Match:    builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
			},
			expectedTraits: HTTPRoutePriorityTraits{
				PreciseHostname: true,
				HostnameLength:  len("foo.com"),
				PathType:        gatewayapi.PathMatchExact,
				PathLength:      len("/foo"),
			},
		},
		{
			name: "wildcard hostname and prefix path",
			match: SplitHTTPRouteMatch{
				Source: &gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "wildcard-hostname-prefix-path",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Hostnames: []gatewayapi.Hostname{"*.foo.com"},
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathPrefix("/foo/").Build(),
								},
							},
						},
					},
				},
				Hostname: "*.foo.com",
				Match:    builder.NewHTTPRouteMatch().WithPathPrefix("/foo/").Build(),
			},
			expectedTraits: HTTPRoutePriorityTraits{
				PreciseHostname: false,
				HostnameLength:  len("*.foo.com"),
				PathType:        gatewayapi.PathMatchPathPrefix,
				PathLength:      len("/foo/"),
			},
		},
		{
			name: "no hostname and regex path, with header matches",
			match: SplitHTTPRouteMatch{
				Source: &gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "no-hostname-regex-path",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathRegex("/[a-z0-9]+").
										WithHeader("foo", "bar").Build(),
								},
							},
						},
					},
				},
				Match: builder.NewHTTPRouteMatch().WithPathRegex("/[a-z0-9]+").
					WithHeader("foo", "bar").Build(),
			},
			expectedTraits: HTTPRoutePriorityTraits{
				PathType:    gatewayapi.PathMatchRegularExpression,
				PathLength:  len("/[a-z0-9]+"),
				HeaderCount: 1,
			},
		},
		{
			name: "precise hostname and method, query param match",
			match: SplitHTTPRouteMatch{
				Source: &gatewayapi.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "precise-hostname-method-query",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							"foo.com",
						},
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithMethod("GET").
										WithQueryParam("foo", "bar").Build(),
								},
							},
						},
					},
				},
				Hostname: "foo.com",
				Match: builder.NewHTTPRouteMatch().WithMethod("GET").
					WithQueryParam("foo", "bar").Build(),
			},
			expectedTraits: HTTPRoutePriorityTraits{
				PreciseHostname: true,
				HostnameLength:  len("foo.com"),
				HasMethodMatch:  true,
				QueryParamCount: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			traits := CalculateHTTPRouteMatchPriorityTraits(tc.match)
			require.Equal(t, tc.expectedTraits, traits)
		})
	}
}

func TestEncodeHTTPRoutePriorityFromTraits(t *testing.T) {
	testCases := []struct {
		name             string
		traits           HTTPRoutePriorityTraits
		expectedPriority RoutePriorityType
	}{
		{
			name: "precise hostname and exact path",
			traits: HTTPRoutePriorityTraits{
				PreciseHostname: true,
				HostnameLength:  7,
				PathType:        gatewayapi.PathMatchExact,
				PathLength:      4,
			},
			expectedPriority: (2 << 44) | (1 << 43) | (7 << 35) | (1 << 34) | (3 << 23),
		},
		{
			name: "wildcard hostname and prefix path",
			traits: HTTPRoutePriorityTraits{
				PreciseHostname: false,
				HostnameLength:  7,
				PathType:        gatewayapi.PathMatchPathPrefix,
				PathLength:      5,
			},
			expectedPriority: (2 << 44) | (7 << 35) | (4 << 23),
		},
		{
			name: "no hostname and regex path, with header matches",
			traits: HTTPRoutePriorityTraits{
				PathType:    gatewayapi.PathMatchRegularExpression,
				PathLength:  5,
				HeaderCount: 2,
			},
			expectedPriority: (2 << 44) | (1 << 33) | (4 << 23) | (2 << 17),
		},
		{
			name: "no hostname and exact path, with method match and query parameter matches",
			traits: HTTPRoutePriorityTraits{
				PathType:        gatewayapi.PathMatchExact,
				PathLength:      5,
				HasMethodMatch:  true,
				QueryParamCount: 1,
			},
			expectedPriority: (2 << 44) | (1 << 34) | (4 << 23) | (1 << 22) | (1 << 12),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedPriority, tc.traits.EncodeToPriority())
		})
	}
}

func TestSplitHTTPRoutes(t *testing.T) {
	namesToBackendRefs := func(names []string) []gatewayapi.HTTPBackendRef {
		backendRefs := []gatewayapi.HTTPBackendRef{}
		for _, name := range names {
			backendRefs = append(backendRefs,
				gatewayapi.HTTPBackendRef{
					BackendRef: gatewayapi.BackendRef{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Name: gatewayapi.ObjectName(name),
						},
					},
				},
			)
		}
		return backendRefs
	}

	testCases := []struct {
		name                 string
		httpRoute            *gatewayapi.HTTPRoute
		expectedSplitMatches []SplitHTTPRouteMatch
	}{
		{
			name: "no hostname and only one match",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "httproute-1",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Matches:     builder.NewHTTPRouteMatch().WithPathExact("/").ToSlice(),
							BackendRefs: namesToBackendRefs([]string{"svc1"}),
						},
					},
				},
			},
			expectedSplitMatches: []SplitHTTPRouteMatch{
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns1",
							Name:      "httproute-1",
						},
					},
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
		},
		{
			name: "multiple hostnames with one match",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "httproute-2",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						"a.foo.com",
						"b.foo.com",
					},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Matches:     builder.NewHTTPRouteMatch().WithPathExact("/").ToSlice(),
							BackendRefs: namesToBackendRefs([]string{"svc1", "svc2"}),
						},
					},
				},
			},
			expectedSplitMatches: []SplitHTTPRouteMatch{
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns1",
							Name:      "httproute-2",
						},
					},
					Hostname:   "a.foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns1",
							Name:      "httproute-2",
						},
					},
					Hostname:   "b.foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
		},
		{
			name: "single hostname with multiple rules and matches",
			httpRoute: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "httproute-3",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					Hostnames: []gatewayapi.Hostname{
						"a.foo.com",
					},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
							},
							BackendRefs: namesToBackendRefs([]string{"svc1"}),
						},
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathExact("/v2/foo").Build(),
								builder.NewHTTPRouteMatch().WithPathExact("/v2/bar").Build(),
							},
							BackendRefs: namesToBackendRefs([]string{"svc2"}),
						},
					},
				},
			},
			expectedSplitMatches: []SplitHTTPRouteMatch{
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns1",
							Name:      "httproute-3",
						},
						// Spec omitted, we do not check for spec for this test
					},
					Hostname:   "a.foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns1",
							Name:      "httproute-3",
						},
					},
					Hostname:   "a.foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
					RuleIndex:  0,
					MatchIndex: 1,
				},
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns1",
							Name:      "httproute-3",
						},
					},
					Hostname:   "a.foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/v2/foo").Build(),
					RuleIndex:  1,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "ns1",
							Name:      "httproute-3",
						},
					},
					Hostname:   "a.foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/v2/bar").Build(),
					RuleIndex:  1,
					MatchIndex: 1,
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			splitHTTPRouteMatches := SplitHTTPRoute(tc.httpRoute)
			require.Len(t, splitHTTPRouteMatches, len(tc.expectedSplitMatches), "should have same number of split matched with expected")
			for i, expectedMatch := range tc.expectedSplitMatches {
				assert.Equal(t, expectedMatch.Source.Name, splitHTTPRouteMatches[i].Source.Name)
				assert.Equal(t, expectedMatch.Match, splitHTTPRouteMatches[i].Match)
				assert.Equal(t, expectedMatch.Hostname, splitHTTPRouteMatches[i].Hostname)
				assert.Equal(t, expectedMatch.RuleIndex, splitHTTPRouteMatches[i].RuleIndex)
				assert.Equal(t, expectedMatch.MatchIndex, splitHTTPRouteMatches[i].MatchIndex)
			}
		})
	}
}

func TestAssignRoutePriorityToSplitHTTPRouteMatches(t *testing.T) {
	type splitHTTPRouteIndex struct {
		namespace  string
		name       string
		hostname   string
		ruleIndex  int
		matchIndex int
	}
	now := time.Now()
	const maxRelativeOrderPriorityBits = (1 << 12) - 1

	testCases := []struct {
		name    string
		matches []SplitHTTPRouteMatch
		// HTTPRoute index -> priority
		priorities map[splitHTTPRouteIndex]RoutePriorityType
	}{
		{
			name: "no dupelicated fixed priority",
			matches: []SplitHTTPRouteMatch{
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "httproute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.HTTPRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									Matches: builder.NewHTTPRouteMatch().WithPathExact("/foo").ToSlice(),
								},
							},
						},
					},
					Hostname:   "foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "httproute-2",
							CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
						},
						Spec: gatewayapi.HTTPRouteSpec{
							Hostnames: []gatewayapi.Hostname{"*.bar.com"},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									Matches: builder.NewHTTPRouteMatch().WithPathExact("/bar").ToSlice(),
								},
							},
						},
					},
					Hostname:   "*.bar.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
			priorities: map[splitHTTPRouteIndex]RoutePriorityType{
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					PathType:        gatewayapi.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
				{
					namespace:  "default",
					name:       "httproute-2",
					hostname:   "*.bar.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: false,
					HostnameLength:  len("*.bar.com"),
					PathType:        gatewayapi.PathMatchExact,
					PathLength:      len("/bar"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
			},
		},
		{
			name: "break tie by creation timestamp",
			matches: []SplitHTTPRouteMatch{
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "httproute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.HTTPRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									Matches: []gatewayapi.HTTPRouteMatch{
										builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
									},
								},
							},
						},
					},
					Hostname:   "foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "httproute-2",
							CreationTimestamp: metav1.NewTime(now.Add(-1 * time.Second)),
						},
						Spec: gatewayapi.HTTPRouteSpec{
							Hostnames: []gatewayapi.Hostname{"bar.com"},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									Matches: []gatewayapi.HTTPRouteMatch{
										builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
									},
								},
							},
						},
					},
					Hostname:   "bar.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
			priorities: map[splitHTTPRouteIndex]RoutePriorityType{
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					PathType:        gatewayapi.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
				{
					namespace:  "default",
					name:       "httproute-2",
					hostname:   "bar.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("bar.com"),
					PathType:        gatewayapi.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits - 1,
			},
		},
		{
			name: "break tie by namespace and name",
			matches: []SplitHTTPRouteMatch{
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "httproute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.HTTPRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									Matches: builder.NewHTTPRouteMatch().WithPathExact("/foo").ToSlice(),
								},
							},
						},
					},
					Hostname:   "foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "httproute-2",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.HTTPRouteSpec{
							Hostnames: []gatewayapi.Hostname{"bar.com"},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									Matches: builder.NewHTTPRouteMatch().WithPathExact("/foo").ToSlice(),
								},
							},
						},
					},
					Hostname:   "bar.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
			},
			priorities: map[splitHTTPRouteIndex]RoutePriorityType{
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					PathType:        gatewayapi.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
				{
					namespace:  "default",
					name:       "httproute-2",
					hostname:   "bar.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("bar.com"),
					PathType:        gatewayapi.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits - 1,
			},
		},
		{
			name: "break tie by internal match index",
			matches: []SplitHTTPRouteMatch{
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "httproute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.HTTPRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									Matches: []gatewayapi.HTTPRouteMatch{
										builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
										builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
									},
								},
							},
						},
					},
					Hostname:   "foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
					RuleIndex:  0,
					MatchIndex: 0,
				},
				{
					Source: &gatewayapi.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Namespace:         "default",
							Name:              "httproute-1",
							CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
						},
						Spec: gatewayapi.HTTPRouteSpec{
							Hostnames: []gatewayapi.Hostname{"foo.com"},
							Rules: []gatewayapi.HTTPRouteRule{
								{
									Matches: []gatewayapi.HTTPRouteMatch{
										builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
										builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
									},
								},
							},
						},
					},
					Hostname:   "foo.com",
					Match:      builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
					RuleIndex:  0,
					MatchIndex: 1,
				},
			},
			priorities: map[splitHTTPRouteIndex]RoutePriorityType{
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					PathType:        gatewayapi.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 1,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("bar.com"),
					PathType:        gatewayapi.PathMatchExact,
					PathLength:      len("/bar"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits - 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			splitHTTPRoutesWithPriorities := assignRoutePriorityToSplitHTTPRouteMatches(logr.Discard(), tc.matches)
			require.Equal(t, len(tc.priorities), len(splitHTTPRoutesWithPriorities), "should have required number of results")
			for _, r := range splitHTTPRoutesWithPriorities {
				httpRoute := r.Match.Source

				require.Equalf(t, tc.priorities[splitHTTPRouteIndex{
					namespace:  httpRoute.Namespace,
					name:       httpRoute.Name,
					hostname:   string(httpRoute.Spec.Hostnames[0]),
					ruleIndex:  r.Match.RuleIndex,
					matchIndex: r.Match.MatchIndex,
				}], r.Priority, "httproute %s/%s: hostname %s, rule %d match %d",
					httpRoute.Namespace, httpRoute.Name, httpRoute.Spec.Hostnames[0], r.Match.RuleIndex, r.Match.MatchIndex)
			}
		})
	}
}

func TestGroupHTTPRouteMatchesWithPrioritiesByRule(t *testing.T) {
	testCases := []struct {
		name                 string
		routes               []*gatewayapi.HTTPRoute
		expectedSplitMatches splitHTTPRouteMatchesWithPrioritiesGroupedByRule
	}{
		{
			name: "single HTTPRoute with single rule, no hostname and single match",
			routes: []*gatewayapi.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches:     builder.NewHTTPRouteMatch().WithPathExact("/").ToSlice(),
								BackendRefs: builder.NewHTTPBackendRef("svc1").ToSlice(),
							},
						},
					},
				},
			},
			expectedSplitMatches: splitHTTPRouteMatchesWithPrioritiesGroupedByRule{
				"ns1/httproute-1.0": []SplitHTTPRouteMatchToKongRoutePriority{
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "",
							Match:    builder.NewHTTPRouteMatch().WithPathExact("/").Build(),
						},
					},
				},
			},
		},
		{
			name: "single HTTPRoute with single rule and multiple hostnames and multiple matches",
			routes: []*gatewayapi.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							gatewayapi.Hostname("a.foo.com"),
							gatewayapi.Hostname("b.bar.com"),
						},
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
								BackendRefs: builder.NewHTTPBackendRef("svc1").ToSlice(),
							},
						},
					},
				},
			},
			expectedSplitMatches: splitHTTPRouteMatchesWithPrioritiesGroupedByRule{
				"ns1/httproute-1.0": []SplitHTTPRouteMatchToKongRoutePriority{
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "a.foo.com",
							Match:    builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
						},
					},
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "a.foo.com",
							Match:    builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
						},
					},
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "b.bar.com",
							Match:    builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
						},
					},
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "b.bar.com",
							Match:    builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
						},
					},
				},
			},
		},
		{
			name: "single HTTPRoute with multiple rules where one of them does not contain a match",
			routes: []*gatewayapi.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							gatewayapi.Hostname("foo.com"),
						},
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
								BackendRefs: builder.NewHTTPBackendRef("svc1").ToSlice(),
							},
							{
								// No matches.
								BackendRefs: builder.NewHTTPBackendRef("svc1").ToSlice(),
							},
						},
					},
				},
			},
			expectedSplitMatches: splitHTTPRouteMatchesWithPrioritiesGroupedByRule{
				"ns1/httproute-1.0": []SplitHTTPRouteMatchToKongRoutePriority{
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "foo.com",
							Match:    builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
						},
					},
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "foo.com",
							Match:    builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
						},
					},
				},
				"ns1/httproute-1.1": []SplitHTTPRouteMatchToKongRoutePriority{
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "foo.com",
						},
					},
				},
			},
		},
		{
			name: "multiple HTTPRoutes where one of them does not have hostnames and matches",
			routes: []*gatewayapi.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Hostnames: []gatewayapi.Hostname{
							gatewayapi.Hostname("foo.com"),
						},
						Rules: []gatewayapi.HTTPRouteRule{
							{
								Matches: []gatewayapi.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
								BackendRefs: builder.NewHTTPBackendRef("svc1").ToSlice(),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-2",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: builder.NewHTTPBackendRef("svc1").ToSlice(),
							},
						},
					},
				},
			},
			expectedSplitMatches: splitHTTPRouteMatchesWithPrioritiesGroupedByRule{
				"ns1/httproute-1.0": []SplitHTTPRouteMatchToKongRoutePriority{
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "foo.com",
							Match:    builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
						},
					},
					{
						Match: SplitHTTPRouteMatch{
							Hostname: "foo.com",
							Match:    builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
						},
					},
				},
				"ns1/httproute-2.0": []SplitHTTPRouteMatchToKongRoutePriority{
					{
						Match: SplitHTTPRouteMatch{},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rulesToMatchesWithPriorities := groupHTTPRouteMatchesWithPrioritiesByRule(logr.Discard(), tc.routes)
			require.Len(t, rulesToMatchesWithPriorities, len(tc.expectedSplitMatches))
			for _, route := range tc.routes {
				for ruleIndex := range route.Spec.Rules {
					ruleKey := fmt.Sprintf("%s/%s.%d", route.Namespace, route.Name, ruleIndex)
					expectedMatches := tc.expectedSplitMatches[ruleKey]
					actualMatches, ok := rulesToMatchesWithPriorities[ruleKey]
					require.Truef(t, ok, "Should find matches from rule %s in HTTPRoute %s/%s", ruleIndex, route.Name, route.Name)
					require.Len(t, actualMatches, len(expectedMatches), "Rule %d in HTTPRoute %s/%s should have expected number of split matches",
						ruleIndex, route.Namespace, route.Name)
					// Sort the matches from the same rule by hostname and match index for comparing.
					// We need to sort the matches because their order were shuffled during the `assignRoutePriorityToSplitHTTPRouteMatches`
					// which takes matches from a map.
					sort.Slice(actualMatches, func(i, j int) bool {
						if actualMatches[i].Match.Hostname != actualMatches[j].Match.Hostname {
							return actualMatches[i].Match.Hostname < actualMatches[j].Match.Hostname
						}
						return actualMatches[i].Match.MatchIndex < actualMatches[j].Match.MatchIndex
					})
					for i, expectedMatch := range expectedMatches {
						require.Equalf(t, expectedMatch.Match.Hostname, actualMatches[i].Match.Hostname, "Split match %d in rule %s, HTTPRoute %s/%s should have expected hostname",
							i, ruleIndex, route.Namespace, route.Name)
						require.Equalf(t, expectedMatch.Match.Match, actualMatches[i].Match.Match, "Split match %d should have expected HTTPRoute match",
							i, ruleIndex, route.Namespace, route.Name)
					}
				}
			}
		})
	}
}
