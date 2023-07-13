package translators

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
)

func TestGenerateKongExpressionRoutesFromHTTPRouteMatches(t *testing.T) {
	testCases := []struct {
		name              string
		routeName         string
		matches           []gatewayv1beta1.HTTPRouteMatch
		filters           []gatewayv1beta1.HTTPRouteFilter
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
			expectedRoutes:    []kongstate.Route{},
			expectedError:     ErrRouteValidationNoMatchRulesOrHostnamesSpecified,
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
						Expression:   kong.String(`(http.host == "foo.com") || (http.host =^ ".bar.com")`),
						Priority:     kong.Int(1),
					},
					ExpressionRoutes: true,
				},
			},
		},
		{
			name:              "single prefix path match",
			routeName:         "prefix_path_match.defualt.0.0",
			ingressObjectInfo: util.K8sObjectInfo{},
			matches: []gatewayv1beta1.HTTPRouteMatch{
				builder.NewHTTPRouteMatch().WithPathPrefix("/prefix").Build(),
			},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("prefix_path_match.defualt.0.0"),
						PreserveHost: kong.Bool(true),
						Expression:   kong.String(`((http.path == "/prefix") || (http.path ^= "/prefix/")) && ((net.protocol == "http") || (net.protocol == "https"))`),
						Priority:     kong.Int(1),
					},
					Plugins:          []kong.Plugin{},
					ExpressionRoutes: true,
				},
			},
		},
		{
			name:              "multiple matches without filters",
			routeName:         "multiple_matches.default.0.0",
			ingressObjectInfo: util.K8sObjectInfo{},
			matches: []gatewayv1beta1.HTTPRouteMatch{
				builder.NewHTTPRouteMatch().WithPathPrefix("/prefix").Build(),
				builder.NewHTTPRouteMatch().WithPathExact("/exact").WithMethod(gatewayv1beta1.HTTPMethodGet).Build(),
			},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("multiple_matches.default.0.0"),
						PreserveHost: kong.Bool(true),
						Expression:   kong.String(`(((http.path == "/prefix") || (http.path ^= "/prefix/")) || ((http.path == "/exact") && (http.method == "GET"))) && ((net.protocol == "http") || (net.protocol == "https"))`),
						Priority:     kong.Int(1),
					},
					Plugins:          []kong.Plugin{},
					ExpressionRoutes: true,
				},
			},
		},
		{
			name:              "multiple matches with request redirect filter",
			routeName:         "request_redirect.default.0.0",
			ingressObjectInfo: util.K8sObjectInfo{},
			matches: []gatewayv1beta1.HTTPRouteMatch{
				builder.NewHTTPRouteMatch().WithPathExact("/exact/0").Build(),
				builder.NewHTTPRouteMatch().WithPathExact("/exact/1").Build(),
			},
			filters: []gatewayv1beta1.HTTPRouteFilter{
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
						Expression:   kong.String(`(http.path == "/exact/0") && ((net.protocol == "http") || (net.protocol == "https"))`),
						Priority:     kong.Int(1),
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
								"add": map[string][]string{
									"headers": {`Location: http://a.foo.com:80/exact/0`},
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
						Expression:   kong.String(`(http.path == "/exact/1") && ((net.protocol == "http") || (net.protocol == "https"))`),
						Priority:     kong.Int(1),
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
								"add": map[string][]string{
									"headers": {`Location: http://a.foo.com:80/exact/1`},
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
			matches: []gatewayv1beta1.HTTPRouteMatch{
				builder.NewHTTPRouteMatch().WithPathExact("/exact/0").Build(),
				builder.NewHTTPRouteMatch().WithPathRegex("/regex/[a-z]+").Build(),
			},
			filters: []gatewayv1beta1.HTTPRouteFilter{
				builder.NewHTTPRouteRequestHeaderModifierFilter().WithRequestHeaderAdd([]gatewayv1beta1.HTTPHeader{
					{Name: "foo", Value: "bar"},
				}).Build(),
			},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("request_header_mod.default.0.0"),
						PreserveHost: kong.Bool(true),
						Expression:   kong.String(`((http.path == "/exact/0") || (http.path ~ "^/regex/[a-z]+")) && ((net.protocol == "http") || (net.protocol == "https"))`),
						Priority:     kong.Int(1),
					},
					Plugins: []kong.Plugin{
						{
							Name: kong.String("request-transformer"),
							Config: kong.Configuration{
								"append": map[string][]string{
									"headers": {"foo:bar"},
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
			matches: []gatewayv1beta1.HTTPRouteMatch{
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
						Expression:   kong.String(`((http.path == "/prefix/0") || (http.path ^= "/prefix/0/")) && (http.host == "a.foo.com") && (net.protocol == "https") && (tls.sni == "a.foo.com")`),
						Priority:     kong.Int(1),
					},
					Plugins:          []kong.Plugin{},
					ExpressionRoutes: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
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
		match      gatewayv1beta1.HTTPRouteMatch
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
				WithMethod(gatewayv1beta1.HTTPMethodGet).
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
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expression, generateMatcherFromHTTPRouteMatch(tc.match).Expression())
		})
	}
}

func TestCalculateHTTPRoutePriorityTraits(t *testing.T) {
	testCases := []struct {
		name           string
		httpRoute      *gatewayv1beta1.HTTPRoute
		expectedTraits HTTPRoutePriorityTraits
	}{
		{
			name: "precise hostname and exact path",
			httpRoute: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "precise-hostname-exact-path",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
					Rules: []gatewayv1beta1.HTTPRouteRule{
						{
							Matches: []gatewayv1beta1.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
							},
						},
					},
				},
			},
			expectedTraits: HTTPRoutePriorityTraits{
				PreciseHostname: true,
				HostnameLength:  len("foo.com"),
				PathType:        gatewayv1beta1.PathMatchExact,
				PathLength:      len("/foo"),
			},
		},
		{
			name: "wildcard hostname and prefix path",
			httpRoute: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "wildcard-hostname-prefix-path",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{"*.foo.com"},
					Rules: []gatewayv1beta1.HTTPRouteRule{
						{
							Matches: []gatewayv1beta1.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathPrefix("/foo/").Build(),
							},
						},
					},
				},
			},
			expectedTraits: HTTPRoutePriorityTraits{
				PreciseHostname: false,
				HostnameLength:  len("*.foo.com"),
				PathType:        gatewayv1beta1.PathMatchPathPrefix,
				PathLength:      len("/foo/"),
			},
		},
		{
			name: "no hostname and regex path, with header matches",
			httpRoute: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "no-hostname-regex-path",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Rules: []gatewayv1beta1.HTTPRouteRule{
						{
							Matches: []gatewayv1beta1.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathRegex("/[a-z0-9]+").
									WithHeader("foo", "bar").Build(),
							},
						},
					},
				},
			},
			expectedTraits: HTTPRoutePriorityTraits{
				PathType:    gatewayv1beta1.PathMatchRegularExpression,
				PathLength:  len("/[a-z0-9]+"),
				HeaderCount: 1,
			},
		},
		{
			name: "precise hostname and method, query param match",
			httpRoute: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "precise-hostname-method-query",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						"foo.com",
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{
						{
							Matches: []gatewayv1beta1.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithMethod("GET").
									WithQueryParam("foo", "bar").Build(),
							},
						},
					},
				},
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			traits := CalculateSplitHTTPRoutePriorityTraits(tc.httpRoute)
			require.Equal(t, tc.expectedTraits, traits)
		})
	}
}

func TestEncodeHTTPRoutePriorityFromTraits(t *testing.T) {
	testCases := []struct {
		name             string
		traits           HTTPRoutePriorityTraits
		expectedPriority int
	}{
		{
			name: "precise hostname and exact path",
			traits: HTTPRoutePriorityTraits{
				PreciseHostname: true,
				HostnameLength:  7,
				PathType:        gatewayv1beta1.PathMatchExact,
				PathLength:      4,
			},
			expectedPriority: (2 << 50) | (1 << 49) | (7 << 41) | (1 << 40) | (3 << 29),
		},
		{
			name: "wildcard hostname and prefix path",
			traits: HTTPRoutePriorityTraits{
				PreciseHostname: false,
				HostnameLength:  7,
				PathType:        gatewayv1beta1.PathMatchPathPrefix,
				PathLength:      5,
			},
			expectedPriority: (2 << 50) | (7 << 41) | (4 << 29),
		},
		{
			name: "no hostname and regex path, with header matches",
			traits: HTTPRoutePriorityTraits{
				PathType:    gatewayv1beta1.PathMatchRegularExpression,
				PathLength:  5,
				HeaderCount: 2,
			},
			expectedPriority: (2 << 50) | (1 << 39) | (4 << 29) | (2 << 23),
		},
		{
			name: "no hostname and exact path, with method match and query parameter matches",
			traits: HTTPRoutePriorityTraits{
				PathType:        gatewayv1beta1.PathMatchExact,
				PathLength:      5,
				HasMethodMatch:  true,
				QueryParamCount: 1,
			},
			expectedPriority: (2 << 50) | (1 << 40) | (4 << 29) | (1 << 28) | (1 << 18),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedPriority, tc.traits.EncodeToPriority())
		})
	}
}

func TestSplitHTTPRoutes(t *testing.T) {
	namesToBackendRefs := func(names []string) []gatewayv1beta1.HTTPBackendRef {
		backendRefs := []gatewayv1beta1.HTTPBackendRef{}
		for _, name := range names {
			backendRefs = append(backendRefs,
				gatewayv1beta1.HTTPBackendRef{
					BackendRef: gatewayv1beta1.BackendRef{
						BackendObjectReference: gatewayv1beta1.BackendObjectReference{
							Name: gatewayv1beta1.ObjectName(name),
						},
					},
				},
			)
		}
		return backendRefs
	}

	testCases := []struct {
		name            string
		httpRoute       *gatewayv1beta1.HTTPRoute
		splitHTTPRoutes []*gatewayv1beta1.HTTPRoute
	}{
		{
			name: "no hostname and only one match",
			httpRoute: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "httproute-1",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Rules: []gatewayv1beta1.HTTPRouteRule{
						{
							Matches: []gatewayv1beta1.HTTPRouteMatch{
								{
									Path: &gatewayv1beta1.HTTPPathMatch{
										Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
										Value: lo.ToPtr("/"),
									},
								},
							},
							BackendRefs: namesToBackendRefs([]string{"svc1"}),
						},
					},
				},
			},
			splitHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-1",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									{
										Path: &gatewayv1beta1.HTTPPathMatch{
											Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
											Value: lo.ToPtr("/"),
										},
									},
								},
								BackendRefs: namesToBackendRefs([]string{"svc1"}),
							},
						},
					},
				},
			},
		},
		{
			name: "multiple hostnames with one match",
			httpRoute: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "httproute-2",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						"a.foo.com",
						"b.foo.com",
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{
						{
							Matches: []gatewayv1beta1.HTTPRouteMatch{
								{
									Path: &gatewayv1beta1.HTTPPathMatch{
										Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
										Value: lo.ToPtr("/"),
									},
								},
							},
							BackendRefs: namesToBackendRefs([]string{"svc1", "svc2"}),
						},
					},
				},
			},
			splitHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-2",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{
							"a.foo.com",
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									{
										Path: &gatewayv1beta1.HTTPPathMatch{
											Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
											Value: lo.ToPtr("/"),
										},
									},
								},
								BackendRefs: namesToBackendRefs([]string{"svc1", "svc2"}),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-2",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{
							"b.foo.com",
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									{
										Path: &gatewayv1beta1.HTTPPathMatch{
											Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
											Value: lo.ToPtr("/"),
										},
									},
								},
								BackendRefs: namesToBackendRefs([]string{"svc1", "svc2"}),
							},
						},
					},
				},
			},
		},
		{
			name: "single hostname with multiple rules and matches",
			httpRoute: &gatewayv1beta1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns1",
					Name:      "httproute-3",
				},
				Spec: gatewayv1beta1.HTTPRouteSpec{
					Hostnames: []gatewayv1beta1.Hostname{
						"a.foo.com",
					},
					Rules: []gatewayv1beta1.HTTPRouteRule{
						{
							Matches: []gatewayv1beta1.HTTPRouteMatch{
								{
									Path: &gatewayv1beta1.HTTPPathMatch{
										Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
										Value: lo.ToPtr("/foo"),
									},
								},
								{
									Path: &gatewayv1beta1.HTTPPathMatch{
										Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
										Value: lo.ToPtr("/bar"),
									},
								},
							},
							BackendRefs: namesToBackendRefs([]string{"svc1"}),
						},
						{
							Matches: []gatewayv1beta1.HTTPRouteMatch{
								{
									Path: &gatewayv1beta1.HTTPPathMatch{
										Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
										Value: lo.ToPtr("/v2/foo"),
									},
								},
								{
									Path: &gatewayv1beta1.HTTPPathMatch{
										Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
										Value: lo.ToPtr("/v2/bar"),
									},
								},
							},
							BackendRefs: namesToBackendRefs([]string{"svc2"}),
						},
					},
				},
			},
			splitHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-3",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{
							"a.foo.com",
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									{
										Path: &gatewayv1beta1.HTTPPathMatch{
											Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
											Value: lo.ToPtr("/foo"),
										},
									},
								},
								BackendRefs: namesToBackendRefs([]string{"svc1"}),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-3",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "1",
						},
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{
							"a.foo.com",
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									{
										Path: &gatewayv1beta1.HTTPPathMatch{
											Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
											Value: lo.ToPtr("/bar"),
										},
									},
								},
								BackendRefs: namesToBackendRefs([]string{"svc1"}),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-3",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "1",
							InternalMatchIndexAnnotationKey: "0",
						},
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{
							"a.foo.com",
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									{
										Path: &gatewayv1beta1.HTTPPathMatch{
											Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
											Value: lo.ToPtr("/v2/foo"),
										},
									},
								},
								BackendRefs: namesToBackendRefs([]string{"svc2"}),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns1",
						Name:      "httproute-3",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "1",
							InternalMatchIndexAnnotationKey: "1",
						},
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{
							"a.foo.com",
						},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									{
										Path: &gatewayv1beta1.HTTPPathMatch{
											Type:  lo.ToPtr(gatewayv1beta1.PathMatchExact),
											Value: lo.ToPtr("/v2/bar"),
										},
									},
								},
								BackendRefs: namesToBackendRefs([]string{"svc2"}),
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		splitHTTPRoutes := SplitHTTPRoute(tc.httpRoute)
		require.Len(t, splitHTTPRoutes, len(splitHTTPRoutes), "should have same number of split HTTPRoutes with expected")
		for i, splitHTTPRoute := range tc.splitHTTPRoutes {
			require.True(t, reflect.DeepEqual(splitHTTPRoute, splitHTTPRoutes[i]))
		}
	}
}

func TestAssignRoutePriorityToSplitHTTPRoutes(t *testing.T) {
	type splitHTTPRouteIndex struct {
		namespace  string
		name       string
		hostname   string
		ruleIndex  int
		matchIndex int
	}
	now := time.Now()
	const maxRelativeOrderPriorityBits = (1 << 18) - 1

	testCases := []struct {
		name            string
		splitHTTPRoutes []*gatewayv1beta1.HTTPRoute
		// HTTPRoute index -> priority
		priorities map[splitHTTPRouteIndex]int
	}{
		{
			name: "no dupelicated fixed priority",
			splitHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-2",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"*.bar.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
							},
						},
					},
				},
			},
			priorities: map[splitHTTPRouteIndex]int{
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					PathType:        gatewayv1beta1.PathMatchExact,
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
					PathType:        gatewayv1beta1.PathMatchExact,
					PathLength:      len("/bar"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
			},
		},
		{
			name: "break tie by creation timestamp",
			splitHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-2",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-1 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"bar.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
			},
			priorities: map[splitHTTPRouteIndex]int{
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					PathType:        gatewayv1beta1.PathMatchExact,
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
					PathType:        gatewayv1beta1.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits - 1,
			},
		},
		{
			name: "break tie namespace and name",
			splitHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-2",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"bar.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
			},
			priorities: map[splitHTTPRouteIndex]int{
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					PathType:        gatewayv1beta1.PathMatchExact,
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
					PathType:        gatewayv1beta1.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits - 1,
			},
		},
		{
			name: "break tie by internal match index",
			splitHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "1",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
							},
						},
					},
				},
			},
			priorities: map[splitHTTPRouteIndex]int{
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					PathType:        gatewayv1beta1.PathMatchExact,
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
					PathType:        gatewayv1beta1.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits - 1,
			},
		},
		{
			name: "httproutes without rule index and internal match index annotations are omitted",
			splitHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey:  "0",
							InternalMatchIndexAnnotationKey: "0",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"foo.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/foo").Build(),
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-2",
						Annotations: map[string]string{
							InternalRuleIndexAnnotationKey: "0",
						},
						CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"*.bar.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:         "default",
						Name:              "httproute-3",
						CreationTimestamp: metav1.NewTime(now.Add(-10 * time.Second)),
					},
					Spec: gatewayv1beta1.HTTPRouteSpec{
						Hostnames: []gatewayv1beta1.Hostname{"a.bar.com"},
						Rules: []gatewayv1beta1.HTTPRouteRule{
							{
								Matches: []gatewayv1beta1.HTTPRouteMatch{
									builder.NewHTTPRouteMatch().WithPathExact("/bar").Build(),
								},
							},
						},
					},
				},
			},
			priorities: map[splitHTTPRouteIndex]int{
				{
					namespace:  "default",
					name:       "httproute-1",
					hostname:   "foo.com",
					ruleIndex:  0,
					matchIndex: 0,
				}: HTTPRoutePriorityTraits{
					PreciseHostname: true,
					HostnameLength:  len("foo.com"),
					PathType:        gatewayv1beta1.PathMatchExact,
					PathLength:      len("/foo"),
				}.EncodeToPriority() + maxRelativeOrderPriorityBits,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			splitHTTPRoutesWithPriorities := AssignRoutePriorityToSplitHTTPRoutes(logr.Discard(), tc.splitHTTPRoutes)
			require.Equal(t, len(tc.priorities), len(splitHTTPRoutesWithPriorities), "should have required number of results")
			for _, r := range splitHTTPRoutesWithPriorities {
				httpRoute := r.HTTPRoute
				ruleIndex, err := strconv.Atoi(httpRoute.Annotations[InternalRuleIndexAnnotationKey])
				require.NoError(t, err)
				matchIndex, err := strconv.Atoi(httpRoute.Annotations[InternalMatchIndexAnnotationKey])
				require.NoError(t, err)

				require.Equalf(t, tc.priorities[splitHTTPRouteIndex{
					namespace:  httpRoute.Namespace,
					name:       httpRoute.Name,
					hostname:   string(httpRoute.Spec.Hostnames[0]),
					ruleIndex:  ruleIndex,
					matchIndex: matchIndex,
				}], r.Priority, "httproute %s/%s: hostname %s, rule %d match %d",
					httpRoute.Namespace, httpRoute.Name, httpRoute.Spec.Hostnames[0], ruleIndex, matchIndex)
			}
		})
	}
}
