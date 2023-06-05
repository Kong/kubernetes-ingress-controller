package translators

import (
	"reflect"
	"testing"

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
		name               string
		httpRoute          *gatewayv1beta1.HTTPRoute
		splittedHTTPRoutes []*gatewayv1beta1.HTTPRoute
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
			splittedHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
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
			splittedHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
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
			splittedHTTPRoutes: []*gatewayv1beta1.HTTPRoute{
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
		splittedHTTPRoutes := SplitHTTPRoute(tc.httpRoute)
		require.Len(t, splittedHTTPRoutes, len(splittedHTTPRoutes), "should have same number of splitted HTTPRoutes with expected")
		for i, splittedHTTPRoute := range tc.splittedHTTPRoutes {
			require.True(t, reflect.DeepEqual(splittedHTTPRoute, splittedHTTPRoutes[i]))
		}
	}
}
