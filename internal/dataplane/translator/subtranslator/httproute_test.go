package subtranslator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

func TestGeneratePluginsFromHTTPRouteFilters(t *testing.T) {
	testCases := []struct {
		name                      string
		filters                   []gatewayapi.HTTPRouteFilter
		path                      string
		expectedPlugins           []kong.Plugin
		expectedPluginsAnnotation string
		expectedErr               error
	}{
		{
			name:            "no filters",
			filters:         []gatewayapi.HTTPRouteFilter{},
			expectedPlugins: []kong.Plugin{},
		},
		{
			name: "request header modifier filter",
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
			name: "request redirect modifier filter",
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
			name: "response header modifier filter",
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
		{
			name: "valid extensionrefs filters",
			filters: []gatewayapi.HTTPRouteFilter{
				{
					Type: gatewayapi.HTTPRouteFilterExtensionRef,
					ExtensionRef: &gatewayapi.LocalObjectReference{
						Group: gatewayapi.Group("configuration.konghq.com"),
						Kind:  gatewayapi.Kind("KongPlugin"),
						Name:  "plugin1",
					},
				},
				{
					Type: gatewayapi.HTTPRouteFilterExtensionRef,
					ExtensionRef: &gatewayapi.LocalObjectReference{
						Group: gatewayapi.Group("configuration.konghq.com"),
						Kind:  gatewayapi.Kind("KongPlugin"),
						Name:  "plugin2",
					},
				},
			},
			expectedPluginsAnnotation: "plugin1,plugin2",
			expectedPlugins:           []kong.Plugin{},
		},
		{
			name: "invalid extensionrefs filter group",
			filters: []gatewayapi.HTTPRouteFilter{
				{
					Type: gatewayapi.HTTPRouteFilterExtensionRef,
					ExtensionRef: &gatewayapi.LocalObjectReference{
						Group: gatewayapi.Group("wrong.group"),
						Kind:  gatewayapi.Kind("KongPlugin"),
						Name:  "plugin1",
					},
				},
			},
			expectedPluginsAnnotation: "",
			expectedErr:               errors.New("plugin wrong.group/KongPlugin unsupported"),
		},
		{
			name: "invalid extensionrefs filter kind",
			filters: []gatewayapi.HTTPRouteFilter{
				{
					Type: gatewayapi.HTTPRouteFilterExtensionRef,
					ExtensionRef: &gatewayapi.LocalObjectReference{
						Group: gatewayapi.Group("configuration.konghq.com"),
						Kind:  gatewayapi.Kind("WrongKind"),
						Name:  "plugin1",
					},
				},
			},
			expectedPluginsAnnotation: "",
			expectedErr:               errors.New("plugin configuration.konghq.com/WrongKind unsupported"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			plugins, pluginsAnnotation, err := generatePluginsFromHTTPRouteFilters(tc.filters, tc.path, nil)
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedPlugins, plugins)
			require.Equal(t, tc.expectedPluginsAnnotation, pluginsAnnotation)
		})
	}
}

func TestGenerateRequestTransformerForURLRewrite(t *testing.T) {
	testCases := []struct {
		name        string
		modifier    *gatewayapi.HTTPURLRewriteFilter
		expected    kong.Plugin
		expectedErr error
	}{
		{
			name: "valid URLRewriteFilter with ReplaceFullPath",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:            gatewayapi.FullPathHTTPPathModifier,
					ReplaceFullPath: lo.ToPtr("/new-path"),
				},
			},
			expected: kong.Plugin{
				Name: lo.ToPtr("request-transformer"),
				Config: kong.Configuration{
					"replace": map[string]string{
						"uri": "/new-path",
					},
				},
			},
			expectedErr: nil,
		},
		// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/3686
		{
			name: "valid URLRewriteFilter with unsupported",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:               gatewayapi.PrefixMatchHTTPPathModifier,
					ReplacePrefixMatch: lo.ToPtr("/new"),
				},
			},
			expected:    kong.Plugin{},
			expectedErr: fmt.Errorf("%s unsupported for %s", gatewayapi.PrefixMatchHTTPPathModifier, gatewayapi.HTTPRouteFilterURLRewrite),
		},
		// TODO: https://github.com/Kong/kubernetes-ingress-controller/issues/3685
		{
			name: "valid URLRewriteFilter with unsupported",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Hostname: lo.ToPtr(gatewayapi.PreciseHostname("hostname")),
			},
			expected:    kong.Plugin{},
			expectedErr: fmt.Errorf("unsupported hostname replace for %s", gatewayapi.HTTPRouteFilterURLRewrite),
		},
		{
			name:        "nil URLRewriteFilter",
			modifier:    nil,
			expected:    kong.Plugin{},
			expectedErr: errors.New("URLRewrite is not provided"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plugin, err := generateRequestTransformerForURLRewrite(tc.modifier)
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expected, plugin)
		})
	}
}
