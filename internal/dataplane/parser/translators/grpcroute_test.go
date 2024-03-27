package translators

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var grpcRouteGVK = schema.GroupVersionKind{
	Group:   "gateway.networking.k8s.io",
	Version: "v1alpha2",
	Kind:    "GRPCRoute",
}

var grpcRouteTypeMeta = metav1.TypeMeta{
	Kind:       "GRPCRoute",
	APIVersion: "gateway.networking.k8s.io/v1alpha2",
}

func makeTestGRPCRoute(
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
			Hostnames: lo.Map(hostnames, func(h string, _ int) gatewayv1.Hostname {
				return gatewayv1.Hostname(h)
			}),
			Rules: rules,
		},
	}
}

func TestGenerateKongRoutesFromGRPCRouteRule(t *testing.T) {
	testCases := []struct {
		name               string
		objectName         string
		annotations        map[string]string
		hostnames          []string
		rule               gatewayv1alpha2.GRPCRouteRule
		prependRegexPrefix bool
		expectedRoutes     []kongstate.Route
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
			prependRegexPrefix: true,
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "single-match",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					Route: kong.Route{
						Name:      kong.String("grpcroute.default.single-match.0.0"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Paths:     kong.StringSlice("~/service0/"),
						Headers: map[string][]string{
							"X-Foo": {"Bar"},
						},
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
			prependRegexPrefix: true,
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "single-match-with-hostname",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},

					Route: kong.Route{
						Name:      kong.String("grpcroute.default.single-match-with-hostname.0.0"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Paths:     kong.StringSlice("~/service0/method0"),
						Hosts:     kong.StringSlice("foo.com", "*.foo.com"),
						Headers:   map[string][]string{},
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
			prependRegexPrefix: true,
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "multiple-matches",
						Namespace:        "default",
						GroupVersionKind: grpcRouteGVK,
					},
					Route: kong.Route{
						Name:      kong.String("grpcroute.default.multiple-matches.0.0"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Paths:     kong.StringSlice("~/.+/method0"),
						Headers: map[string][]string{
							"Version": {"2"},
							"Client":  {"kong-test"},
						},
					},
				},
				{
					Ingress: util.K8sObjectInfo{
						Name:             "multiple-matches",
						Namespace:        "default",
						GroupVersionKind: grpcRouteGVK,
					},
					Route: kong.Route{
						Name:      kong.String("grpcroute.default.multiple-matches.0.1"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Paths:     kong.StringSlice("~/v[012]/.+"),
						Headers:   map[string][]string{},
					},
				},
			},
		},
		{
			name:        "multiple matches with hostname, prependRegexPrefix off",
			objectName:  "multiple-matches-with-hostname",
			annotations: map[string]string{},
			hostnames:   []string{"foo.com", "*.foo.com"},
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
			prependRegexPrefix: false,
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "multiple-matches-with-hostname",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					Route: kong.Route{
						Name:      kong.String("grpcroute.default.multiple-matches-with-hostname.0.0"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Paths:     kong.StringSlice("/.+/method0"),
						Hosts:     kong.StringSlice("foo.com", "*.foo.com"),
						Headers: map[string][]string{
							"Version": {"2"},
							"Client":  {"kong-test"},
						},
					},
				},
				{
					Ingress: util.K8sObjectInfo{
						Name:             "multiple-matches-with-hostname",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					Route: kong.Route{
						Name:      kong.String("grpcroute.default.multiple-matches-with-hostname.0.1"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Paths:     kong.StringSlice("/v[012]/.+"),
						Hosts:     kong.StringSlice("foo.com", "*.foo.com"),
						Headers:   map[string][]string{},
					},
				},
			},
		},
		{
			name:               "single hostname, no matches",
			objectName:         "hostname-only",
			annotations:        map[string]string{},
			hostnames:          []string{"foo.com"},
			rule:               gatewayv1alpha2.GRPCRouteRule{},
			prependRegexPrefix: true,
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "hostname-only",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					Route: kong.Route{
						Name:      kong.String("grpcroute.default.hostname-only.0.0"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Hosts:     kong.StringSlice("foo.com"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			grpcroute := makeTestGRPCRoute(tc.objectName, "default", tc.annotations, tc.hostnames, []gatewayv1alpha2.GRPCRouteRule{tc.rule})
			routes := GenerateKongRoutesFromGRPCRouteRule(grpcroute, 0, tc.prependRegexPrefix)
			require.Equal(t, tc.expectedRoutes, routes)
		})
	}
}
