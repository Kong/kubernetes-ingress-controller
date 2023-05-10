package translators

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func TestGeneratePluginsFromHTTPRouteFilters(t *testing.T) {
	testCases := []struct {
		name            string
		filters         []gatewayv1beta1.HTTPRouteFilter
		path            string
		expectedPlugins []kong.Plugin
	}{
		{
			name:            "no filters",
			filters:         []gatewayv1beta1.HTTPRouteFilter{},
			expectedPlugins: []kong.Plugin{},
		},
		{
			name: "request header modifier",
			filters: []gatewayv1beta1.HTTPRouteFilter{
				{
					Type: gatewayv1beta1.HTTPRouteFilterRequestHeaderModifier,
					RequestHeaderModifier: &gatewayv1beta1.HTTPHeaderFilter{
						Set: []gatewayv1beta1.HTTPHeader{
							{
								Name:  "header-to-set",
								Value: "bar",
							},
						},
						Remove: []string{"header-to-remove"},
					},
				},
			},
			expectedPlugins: []kong.Plugin{
				{
					Name: kong.String("request-transformer"),
					Config: kong.Configuration{
						"add": map[string][]string{
							"headers": {
								"header-to-set:bar",
							},
						},
						"remove": map[string][]string{
							"headers": {
								"header-to-remove",
							},
						},
						"replace": map[string][]string{
							"headers": {
								"header-to-set:bar",
							},
						},
					},
				},
			},
		},
		{
			name: "request redirect modifier",
			filters: []gatewayv1beta1.HTTPRouteFilter{
				{
					Type: gatewayv1beta1.HTTPRouteFilterRequestRedirect,
					RequestRedirect: &gatewayv1beta1.HTTPRequestRedirectFilter{
						Hostname:   (*gatewayv1beta1.PreciseHostname)(lo.ToPtr("example.org")),
						StatusCode: lo.ToPtr(302),
					},
				},
			},
			path: "/test",
			expectedPlugins: []kong.Plugin{
				{
					Name: kong.String("request-termination"),
					Config: kong.Configuration{
						"status_code": lo.ToPtr(302),
					},
				},
				{
					Name: kong.String("response-transformer"),
					Config: kong.Configuration{
						"add": map[string][]string{
							"headers": {
								"Location: http://example.org:80/test",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		plugins := GeneratePluginsFromHTTPRouteFilters(tc.filters, tc.path, nil)
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedPlugins, plugins)
		})
	}
}
