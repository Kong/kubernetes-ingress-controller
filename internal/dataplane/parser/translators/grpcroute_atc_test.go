package translators

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var grpcRouteGVK = schema.GroupVersionKind{
	Group:   "gateway.networking.k8s.io",
	Version: "v1alpha2",
	Kind:    "GRPCRoute",
}

func TestGenerateKongExpressionRoutesFromGRPCRouteRule(t *testing.T) {
	makeGRPCRoute := func(
		name string, namespace string, annotations map[string]string,
		hostnames []string,
		rules []gatewayv1alpha2.GRPCRouteRule,
	) *gatewayv1alpha2.GRPCRoute {
		return &gatewayv1alpha2.GRPCRoute{
			TypeMeta: metav1.TypeMeta{
				Kind:       "GRPCRoute",
				APIVersion: "gateway.networking.k8s.io/v1alpha2",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Namespace:   namespace,
				Annotations: annotations,
			},
			Spec: gatewayv1alpha2.GRPCRouteSpec{
				Hostnames: lo.Map(hostnames, func(h string, _ int) gatewayv1beta1.Hostname {
					return gatewayv1beta1.Hostname(h)
				}),
				Rules: rules,
			},
		}
	}
	testCases := []struct {
		name           string
		objectName     string
		annotations    map[string]string
		hostnames      []string
		rule           gatewayv1alpha2.GRPCRouteRule
		expectedRoutes []kongstate.Route
	}{
		{
			name:        "single match without hostname",
			objectName:  "single-match",
			annotations: map[string]string{},
			rule: gatewayv1alpha2.GRPCRouteRule{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("service0"),
							Method:  nil,
						},
						Headers: []gatewayv1alpha2.GRPCHeaderMatch{
							{
								Name:  gatewayv1alpha2.GRPCHeaderName("X-Foo"),
								Value: "Bar",
							},
						},
					},
				},
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "single-match",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.single-match.0.0"),
						Expression: kong.String(`(http.path ^= "/service0/") && (http.headers.x_foo == "Bar")`),
						Priority:   kong.Int(1),
					},
				},
			},
		},
		{
			name:        "single match with hostname",
			objectName:  "single-match-with-hostname",
			annotations: map[string]string{},
			hostnames:   []string{"foo.com", "*.foo.com"},
			rule: gatewayv1alpha2.GRPCRouteRule{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("service0"),
							Method:  lo.ToPtr("method0"),
						},
					},
				},
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "single-match-with-hostname",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.single-match-with-hostname.0.0"),
						Expression: kong.String(`(http.path == "/service0/method0") && ((http.host == "foo.com") || (http.host =^ ".foo.com"))`),
						Priority:   kong.Int(1),
					},
				},
			},
		},
		{
			name:       "multuple matches without hostname",
			objectName: "multiple-matches",
			rule: gatewayv1alpha2.GRPCRouteRule{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: nil,
							Method:  lo.ToPtr("method0"),
						},
						Headers: []gatewayv1alpha2.GRPCHeaderMatch{
							{
								Name:  "Version",
								Value: "2",
							},
							{
								Name:  "Client",
								Value: "kong-test",
							},
						},
					},
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchRegularExpression),
							Service: lo.ToPtr("v[012]"),
						},
					},
				},
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "multiple-matches",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.multiple-matches.0.0"),
						Expression: kong.String(`(http.path =^ "/method0") && ((http.headers.client == "kong-test") && (http.headers.version == "2"))`),
						Priority:   kong.Int(1),
					},
				},
				{
					Ingress: util.K8sObjectInfo{
						Name:             "multiple-matches",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.multiple-matches.0.1"),
						Expression: kong.String(`http.path ~ "^/v[012]/.+"`),
						Priority:   kong.Int(1),
					},
				},
			},
		},
		{
			name:       "single match with annotations",
			objectName: "single-match-with-annotations",
			annotations: map[string]string{
				"konghq.com/methods":   "POST,GET",
				"konghq.com/protocols": "https",
				"konghq.com/snis":      "kong.foo.com",
			},
			rule: gatewayv1alpha2.GRPCRouteRule{
				Matches: []gatewayv1alpha2.GRPCRouteMatch{
					{
						Method: &gatewayv1alpha2.GRPCMethodMatch{
							Service: lo.ToPtr("service0"),
							Method:  lo.ToPtr("method0"),
						},
					},
				},
			},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:      "single-match-with-annotations",
						Namespace: "default",
						Annotations: map[string]string{
							"konghq.com/methods":   "POST,GET",
							"konghq.com/protocols": "https",
							"konghq.com/snis":      "kong.foo.com",
						},
						GroupVersionKind: grpcRouteGVK,
					},
					ExpressionRoutes: true,
					Route: kong.Route{
						Name:       kong.String("grpcroute.default.single-match-with-annotations.0.0"),
						Expression: kong.String(`(http.path == "/service0/method0") && (net.protocol == "https") && (tls.sni == "kong.foo.com")`),
						Priority:   kong.Int(1),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			grpcroute := makeGRPCRoute(tc.objectName, "default", tc.annotations, tc.hostnames, []gatewayv1alpha2.GRPCRouteRule{tc.rule})
			routes := GenerateKongExpressionRoutesFromGRPCRouteRule(grpcroute, 0, tc.rule)
			require.Equal(t, tc.expectedRoutes, routes)
		})
	}
}

func TestMethodMatcherFromGRPCMethodMatch(t *testing.T) {
	testCases := []struct {
		name       string
		match      gatewayv1alpha2.GRPCMethodMatch
		expression string
	}{
		{
			name: "exact method match",
			match: gatewayv1alpha2.GRPCMethodMatch{
				Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
				Service: lo.ToPtr("service0"),
				Method:  lo.ToPtr("method0"),
			},
			expression: `http.path == "/service0/method0"`,
		},
		{
			name: "exact match with service unspecified",
			match: gatewayv1alpha2.GRPCMethodMatch{
				Type:   lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchExact),
				Method: lo.ToPtr("method0"),
			},
			expression: `http.path =^ "/method0"`,
		},
		{
			name: "regex method match",
			match: gatewayv1alpha2.GRPCMethodMatch{
				Type:    lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchRegularExpression),
				Service: lo.ToPtr("auth-v[012]"),
				Method:  lo.ToPtr("[a-z]+"),
			},
			expression: `http.path ~ "^/auth-v[012]/[a-z]+"`,
		},
		{
			name: "empty regex match",
			match: gatewayv1alpha2.GRPCMethodMatch{
				Type: lo.ToPtr(gatewayv1alpha2.GRPCMethodMatchRegularExpression),
			},
			expression: `http.path ^= "/"`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expression, methodMatcherFromGRPCMethodMatch(&tc.match).Expression())
		})
	}
}
