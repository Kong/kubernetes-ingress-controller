package translators

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
)

func TestGeneratePluginsFromHTTPRouteFilters(t *testing.T) {
	testCases := []struct {
		name            string
		filters         []gatewayapi.HTTPRouteFilter
		path            string
		expectedPlugins []kong.Plugin
	}{
		{
			name:            "no filters",
			filters:         []gatewayapi.HTTPRouteFilter{},
			expectedPlugins: []kong.Plugin{},
		},
		{
			name: "request header modifier",
			filters: []gatewayapi.HTTPRouteFilter{
				{
					Type: gatewayapi.HTTPRouteFilterRequestHeaderModifier,
					RequestHeaderModifier: &gatewayapi.HTTPHeaderFilter{
						Set: []gatewayapi.HTTPHeader{
							{
								Name:  "header-to-set",
								Value: "bar",
							},
						},
						Add: []gatewayapi.HTTPHeader{
							{
								Name:  "header-to-add",
								Value: "foo",
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
						"append": map[string][]string{
							"headers": {
								"header-to-add:foo",
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
			filters: []gatewayapi.HTTPRouteFilter{
				{
					Type: gatewayapi.HTTPRouteFilterRequestRedirect,
					RequestRedirect: &gatewayapi.HTTPRequestRedirectFilter{
						Hostname:   (*gatewayapi.PreciseHostname)(lo.ToPtr("example.org")),
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
		{
			name: "response header modifier",
			filters: []gatewayapi.HTTPRouteFilter{
				{
					Type: gatewayapi.HTTPRouteFilterResponseHeaderModifier,
					ResponseHeaderModifier: &gatewayapi.HTTPHeaderFilter{
						Set: []gatewayapi.HTTPHeader{
							{
								Name:  "header-to-set",
								Value: "bar",
							},
						},
						Add: []gatewayapi.HTTPHeader{
							{
								Name:  "header-to-add",
								Value: "foo",
							},
						},
						Remove: []string{"header-to-remove"},
					},
				},
			},
			expectedPlugins: []kong.Plugin{
				{
					Name: kong.String("response-transformer"),
					Config: kong.Configuration{
						"add": map[string][]string{
							"headers": {
								"header-to-set:bar",
							},
						},
						"append": map[string][]string{
							"headers": {
								"header-to-add:foo",
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
	}

	for _, tc := range testCases {
		tc := tc
		plugins := GeneratePluginsFromHTTPRouteFilters(tc.filters, tc.path, nil)
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedPlugins, plugins)
		})
	}
}
