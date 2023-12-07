package subtranslator

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
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
	rules []gatewayapi.GRPCRouteRule,
) *gatewayapi.GRPCRoute {
	return &gatewayapi.GRPCRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GRPCRoute",
			APIVersion: "gateway.networking.k8s.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: gatewayapi.GRPCRouteSpec{
			Hostnames: lo.Map(hostnames, func(h string, _ int) gatewayapi.Hostname {
				return gatewayapi.Hostname(h)
			}),
			Rules: rules,
		},
	}
}

func TestGenerateKongRoutesFromGRPCRouteRule(t *testing.T) {
	testCases := []struct {
		name           string
		objectName     string
		annotations    map[string]string
		hostnames      []string
		rule           gatewayapi.GRPCRouteRule
		expectedRoutes []kongstate.Route
	}{
		{
			name:        "single match without hostname",
			objectName:  "single-match",
			annotations: map[string]string{},
			rule: gatewayapi.GRPCRouteRule{
				Matches: []gatewayapi.GRPCRouteMatch{
					{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: lo.ToPtr("service0"),
							Method:  nil,
						},
						Headers: []gatewayapi.GRPCHeaderMatch{
							{
								Name:  gatewayapi.GRPCHeaderName("X-Foo"),
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
					Route: kong.Route{
						Name:      kong.String("grpcroute.default.single-match.0.0"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Paths:     kong.StringSlice("~/service0/"),
						Headers: map[string][]string{
							"X-Foo": {"Bar"},
						},
						Tags: kong.StringSlice(
							"k8s-name:single-match",
							"k8s-namespace:default",
							"k8s-kind:GRPCRoute",
							"k8s-group:gateway.networking.k8s.io",
							"k8s-version:v1alpha2",
						),
					},
				},
			},
		},
		{
			name:        "single match with hostname",
			objectName:  "single-match-with-hostname",
			annotations: map[string]string{},
			hostnames:   []string{"foo.com", "*.foo.com"},
			rule: gatewayapi.GRPCRouteRule{
				Matches: []gatewayapi.GRPCRouteMatch{
					{
						Method: &gatewayapi.GRPCMethodMatch{
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

					Route: kong.Route{
						Name:      kong.String("grpcroute.default.single-match-with-hostname.0.0"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Paths:     kong.StringSlice("~/service0/method0"),
						Hosts:     kong.StringSlice("foo.com", "*.foo.com"),
						Headers:   map[string][]string{},
						Tags: kong.StringSlice(
							"k8s-name:single-match-with-hostname",
							"k8s-namespace:default",
							"k8s-kind:GRPCRoute",
							"k8s-group:gateway.networking.k8s.io",
							"k8s-version:v1alpha2",
						),
					},
				},
			},
		},
		{
			name:       "multiple matches without hostname",
			objectName: "multiple-matches",
			rule: gatewayapi.GRPCRouteRule{
				Matches: []gatewayapi.GRPCRouteMatch{
					{
						Method: &gatewayapi.GRPCMethodMatch{
							Service: nil,
							Method:  lo.ToPtr("method0"),
						},
						Headers: []gatewayapi.GRPCHeaderMatch{
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
						Method: &gatewayapi.GRPCMethodMatch{
							Type:    lo.ToPtr(gatewayapi.GRPCMethodMatchRegularExpression),
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
						Tags: kong.StringSlice(
							"k8s-name:multiple-matches",
							"k8s-namespace:default",
							"k8s-kind:GRPCRoute",
							"k8s-group:gateway.networking.k8s.io",
							"k8s-version:v1alpha2",
						),
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
						Tags: kong.StringSlice(
							"k8s-name:multiple-matches",
							"k8s-namespace:default",
							"k8s-kind:GRPCRoute",
							"k8s-group:gateway.networking.k8s.io",
							"k8s-version:v1alpha2",
						),
					},
				},
			},
		},
		{
			name:        "single hostname, no matches",
			objectName:  "hostname-only",
			annotations: map[string]string{},
			hostnames:   []string{"foo.com"},
			rule:        gatewayapi.GRPCRouteRule{},
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
						Tags: kong.StringSlice(
							"k8s-name:hostname-only",
							"k8s-namespace:default",
							"k8s-kind:GRPCRoute",
							"k8s-group:gateway.networking.k8s.io",
							"k8s-version:v1alpha2",
						),
					},
				},
			},
		},
		{
			name:        "single no hostnames, no matches",
			objectName:  "catch-all",
			annotations: map[string]string{},
			rule:        gatewayapi.GRPCRouteRule{},
			expectedRoutes: []kongstate.Route{
				{
					Ingress: util.K8sObjectInfo{
						Name:             "catch-all",
						Namespace:        "default",
						Annotations:      map[string]string{},
						GroupVersionKind: grpcRouteGVK,
					},
					Route: kong.Route{
						Name:      kong.String("grpcroute.default.catch-all.0.0"),
						Protocols: kong.StringSlice("grpc", "grpcs"),
						Tags: kong.StringSlice(
							"k8s-name:catch-all",
							"k8s-namespace:default",
							"k8s-kind:GRPCRoute",
							"k8s-group:gateway.networking.k8s.io",
							"k8s-version:v1alpha2",
						),
						Paths: kong.StringSlice("/"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			grpcroute := makeTestGRPCRoute(tc.objectName, "default", tc.annotations, tc.hostnames, []gatewayapi.GRPCRouteRule{tc.rule})
			routes := GenerateKongRoutesFromGRPCRouteRule(grpcroute, 0)
			require.Equal(t, tc.expectedRoutes, routes)
		})
	}
}

func TestGetGRPCRouteHostnamesAsSliceOfStringPointers(t *testing.T) {
	for _, tC := range []struct {
		name      string
		grpcroute *gatewayapi.GRPCRoute
		expected  []*string
	}{
		{
			name: "single hostname",
			grpcroute: &gatewayapi.GRPCRoute{
				Spec: gatewayapi.GRPCRouteSpec{
					Hostnames: []gatewayapi.Hostname{"example.com"},
				},
			},
			expected: []*string{
				lo.ToPtr("example.com"),
			},
		},
		{
			name: "multiple hostnames",
			grpcroute: &gatewayapi.GRPCRoute{
				Spec: gatewayapi.GRPCRouteSpec{
					Hostnames: []gatewayapi.Hostname{"example.com", "api.example.com"},
				},
			},
			expected: []*string{
				lo.ToPtr("example.com"),
				lo.ToPtr("api.example.com"),
			},
		},
		{
			name:      "nil hostnames",
			grpcroute: &gatewayapi.GRPCRoute{},
			expected:  nil,
		},
	} {
		t.Run(tC.name, func(t *testing.T) {
			result := getGRPCRouteHostnamesAsSliceOfStringPointers(tC.grpcroute)
			require.Equal(t, tC.expected, result)
		})
	}
}
