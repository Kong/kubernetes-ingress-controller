package parser

import (
	"fmt"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

func TestGetKubernetesObjectReferenceForKongRouteName(t *testing.T) {
	for _, tt := range []struct {
		name         string
		routeName    string
		expectedType string
		expectedNSN  types.NamespacedName
		expectedErr  error
	}{
		{
			name:        "a route name without at least 4 sections is invalid",
			routeName:   "invalid.route.name",
			expectedErr: fmt.Errorf("invalid route name invalid.route.name"),
		},
		{
			name:         "a route with 4 sections is valid",
			routeName:    fmt.Sprintf("%s.default.httpbin.FAKEUID", IngressV1RoutePrefix),
			expectedType: IngressV1RoutePrefix,
			expectedNSN: types.NamespacedName{
				Namespace: "default",
				Name:      "httpbin",
			},
		},
		{
			name:         "a route which belongs to an object with periods in the name is valid",
			routeName:    fmt.Sprintf("%s.default.name.with.periods.FAKEUID", HTTPRouteV1Alpha2RoutePrefix),
			expectedType: HTTPRouteV1Alpha2RoutePrefix,
			expectedNSN: types.NamespacedName{
				Namespace: "default",
				Name:      "name.with.periods",
			},
		},
		{
			name:         "an Ingress linked route using the legacy naming convention is valid, but will return an unknown type",
			routeName:    "default.v1beta1ingress-object.00",
			expectedType: "unknown",
			expectedNSN: types.NamespacedName{
				Namespace: "default",
				Name:      "v1beta1ingress-object",
			},
		},
		{
			name:         "a UDPIngress linked route using the legacy naming convention is valid and will return the correct type",
			routeName:    "default.v1beta1udpingress-object.0.udp",
			expectedType: UDPIngressV1Beta1RoutePrefix,
			expectedNSN: types.NamespacedName{
				Namespace: "default",
				Name:      "v1beta1udpingress-object",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			objectType, nsn, err := GetKubernetesObjectReferenceForKongRouteName(tt.routeName)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedType, objectType)
				assert.Equal(t, tt.expectedNSN, nsn)
			}
		})
	}
}

func Test_getUniqIDForRouteConfig(t *testing.T) {
	for _, tt := range []struct {
		name   string
		config interface{}
		id     string
		err    string
	}{
		{
			name:   "a unique ID can not be created for a function",
			config: func() string { return "just you try to serialize me" },
			err:    "could not marshal input for route hash: json: unsupported type: func() string",
		},
		{
			name: "a unique ID can be created for a path match with no hostnames",
			config: httpRouteInput{
				Match: gatewayv1alpha2.HTTPRouteMatch{
					Path: &gatewayv1alpha2.HTTPPathMatch{
						Type:  &httproutePathMatchPrefix,
						Value: kong.String("/httpbin"),
					},
				},
			},
			id: "8CF9GAJDICCJG",
		},
		{
			name: "a unique ID can be created for a path match with hostnames",
			config: httpRouteInput{
				Hostnames: []*string{
					kong.String("konghq.com"),
					kong.String("docs.konghq.com"),
				},
				Match: gatewayv1alpha2.HTTPRouteMatch{
					Path: &gatewayv1alpha2.HTTPPathMatch{
						Type:  &httproutePathMatchPrefix,
						Value: kong.String("/httpbin"),
					},
				},
			},
			id: "HELRE3LKN8H24",
		},
		{
			name: "a difference of one character from a previous match makes a completely unique ID",
			config: httpRouteInput{
				Hostnames: []*string{
					kong.String("konghq.com"),
					kong.String("docs.konghq.com"),
				},
				Match: gatewayv1alpha2.HTTPRouteMatch{
					Path: &gatewayv1alpha2.HTTPPathMatch{
						Type:  &httproutePathMatchPrefix,
						Value: kong.String("/httpbin2"),
					},
				},
			},
			id: "LERGSFU4489F6",
		},
		{
			name: "a match using many hostnames and various matching strategies can be provided a unique ID",
			config: httpRouteInput{
				Hostnames: []*string{
					kong.String("a.konghq.com"),
					kong.String("b.konghq.com"),
					kong.String("c.konghq.com"),
					kong.String("d.konghq.com"),
					kong.String("e.konghq.com"),
					kong.String("f.konghq.com"),
				},
				Match: gatewayv1alpha2.HTTPRouteMatch{
					Path: &gatewayv1alpha2.HTTPPathMatch{
						Type:  &httproutePathMatchPrefix,
						Value: kong.String("/httpbin"),
					},
					Headers: []gatewayv1alpha2.HTTPHeaderMatch{
						{
							Name:  "Content-Type",
							Value: "audio/ogg-vorbis",
						},
						{
							Name:  "Content-Type",
							Value: "audio/mpeg",
						},
					},
				},
			},
			id: "PL8TCCAGTA1M2",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			id, err := getUniqIDForRouteConfig(tt.config)
			if tt.err != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, id)
			}
		})
	}
}
