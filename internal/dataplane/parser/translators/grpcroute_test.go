package translators

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
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
		name               string
		objectName         string
		annotations        map[string]string
		hostnames          []string
		rule               gatewayapi.GRPCRouteRule
		prependRegexPrefix bool
		expectedRoutes     []kongstate.Route
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
			rule:               gatewayapi.GRPCRouteRule{},
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
			grpcroute := makeTestGRPCRoute(tc.objectName, "default", tc.annotations, tc.hostnames, []gatewayapi.GRPCRouteRule{tc.rule})
			routes := GenerateKongRoutesFromGRPCRouteRule(grpcroute, 0, tc.prependRegexPrefix)
			require.Equal(t, tc.expectedRoutes, routes)
		})
	}
}
