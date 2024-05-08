package subtranslator

import (
	"errors"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func TestGeneratePluginsFromHTTPRouteFilters(t *testing.T) {
	testCases := []struct {
		name                       string
		filters                    []gatewayapi.HTTPRouteFilter
		path                       string
		expectedPlugins            []kong.Plugin
		expectedRouteModifications kongstate.Route
		expectedErr                error
	}{
		{
			name:            "no filters",
			filters:         []gatewayapi.HTTPRouteFilter{},
			expectedPlugins: nil,
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
						"add": TransformerPluginConfig{
							Headers: []string{
								"header-to-set:bar",
							},
						},
						"append": TransformerPluginConfig{
							Headers: []string{
								"header-to-add:foo",
							},
						},
						"remove": TransformerPluginConfig{
							Headers: []string{
								"header-to-remove",
							},
						},
						"replace": TransformerPluginReplaceConfig{
							Headers: []string{
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
						"add": TransformerPluginConfig{
							Headers: []string{
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
						"add": TransformerPluginConfig{
							Headers: []string{
								"header-to-set:bar",
							},
						},
						"append": TransformerPluginConfig{
							Headers: []string{
								"header-to-add:foo",
							},
						},
						"remove": TransformerPluginConfig{
							Headers: []string{
								"header-to-remove",
							},
						},
						"replace": TransformerPluginReplaceConfig{
							Headers: []string{
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
			expectedRouteModifications: kongstate.Route{
				Ingress: util.K8sObjectInfo{
					Annotations: map[string]string{
						"konghq.com/plugins": "plugin1,plugin2",
					},
				},
			},
			expectedPlugins: nil,
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
			expectedErr: errors.New("plugin wrong.group/KongPlugin unsupported"),
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
			expectedErr: errors.New("plugin configuration.konghq.com/WrongKind unsupported"),
		},
		{
			name: "RequestHeaderModifier and PrefixMatchHTTPPathModifier",
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
					},
				},
				{
					Type: gatewayapi.HTTPRouteFilterURLRewrite,
					URLRewrite: &gatewayapi.HTTPURLRewriteFilter{
						Path: &gatewayapi.HTTPPathModifier{
							Type:               gatewayapi.PrefixMatchHTTPPathModifier,
							ReplacePrefixMatch: lo.ToPtr("/new"),
						},
					},
				},
			},
			path: "/prefix",
			expectedPlugins: []kong.Plugin{
				{
					Name: kong.String("request-transformer"),
					Config: kong.Configuration{
						"add": TransformerPluginConfig{
							Headers: []string{
								"header-to-set:bar",
							},
						},
						"replace": TransformerPluginReplaceConfig{
							URI: "/new$(uri_captures[1])",
							Headers: []string{
								"header-to-set:bar",
							},
						},
					},
				},
			},
			expectedRouteModifications: kongstate.Route{
				Route: kong.Route{
					Paths: []*string{
						lo.ToPtr("~/prefix$"),
						lo.ToPtr("~/prefix(/.*)"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := generatePluginsFromHTTPRouteFilters(tc.filters, tc.path, nil, false)
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedPlugins, result.Plugins)

			route := kongstate.Route{}
			for _, modifier := range result.KongRouteModifiers {
				modifier(&route)
			}
			require.Equal(t, tc.expectedRouteModifications, route)
		})
	}
}

func TestGenerateRequestTransformerForURLRewrite(t *testing.T) {
	testCases := []struct {
		name                          string
		modifier                      *gatewayapi.HTTPURLRewriteFilter
		firstMatchPath                string
		expectedKongRouteModification kongstate.Route
		expected                      []transformerPlugin
		expectedErr                   error
	}{
		{
			name: "valid URLRewriteFilter with ReplaceFullPath",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:            gatewayapi.FullPathHTTPPathModifier,
					ReplaceFullPath: lo.ToPtr("/new-path"),
				},
			},
			expected: []transformerPlugin{{
				Type: transformerPluginTypeRequest,
				Replace: TransformerPluginReplaceConfig{
					URI: "/new-path",
				},
			}},
			expectedErr: nil,
		},
		{
			name: "URLRewriteFilter with non-empty ReplacePrefixMatch",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:               gatewayapi.PrefixMatchHTTPPathModifier,
					ReplacePrefixMatch: lo.ToPtr("/new"),
				},
			},
			firstMatchPath: "/prefix",
			expected: []transformerPlugin{{
				Type: transformerPluginTypeRequest,
				Replace: TransformerPluginReplaceConfig{
					URI: "/new$(uri_captures[1])",
				},
			}},
			expectedKongRouteModification: kongstate.Route{
				Route: kong.Route{
					Paths: []*string{
						lo.ToPtr("~/prefix$"),
						lo.ToPtr("~/prefix(/.*)"),
					},
				},
			},
		},
		{
			name: "URLRewriteFilter with empty ReplacePrefixMatch",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:               gatewayapi.PrefixMatchHTTPPathModifier,
					ReplacePrefixMatch: lo.ToPtr(""),
				},
			},
			firstMatchPath: "/prefix",
			expected: []transformerPlugin{{
				Type: transformerPluginTypeRequest,
				Replace: TransformerPluginReplaceConfig{
					URI: `$(uri_captures[1] == nil and "/" or uri_captures[1])`,
				},
			}},
			expectedKongRouteModification: kongstate.Route{
				Route: kong.Route{
					Paths: []*string{
						lo.ToPtr("~/prefix$"),
						lo.ToPtr("~/prefix(/.*)"),
					},
				},
			},
		},
		{
			name: "URLRewriteFilter with empty firstMatchPath",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:               gatewayapi.PrefixMatchHTTPPathModifier,
					ReplacePrefixMatch: lo.ToPtr("/prefix"),
				},
			},
			firstMatchPath: "",
			expected: []transformerPlugin{{
				Type: transformerPluginTypeRequest,
				Replace: TransformerPluginReplaceConfig{
					URI: `/prefix$(uri_captures[1] == nil and "" or "/" .. uri_captures[1])`,
				},
			}},
			expectedKongRouteModification: kongstate.Route{
				Route: kong.Route{
					Paths: []*string{
						lo.ToPtr("~/$"),
						lo.ToPtr("~/(.*)"),
					},
				},
			},
		},
		{
			name: "URLRewriteFilter with '/' ReplacePrefixPatch",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:               gatewayapi.PrefixMatchHTTPPathModifier,
					ReplacePrefixMatch: lo.ToPtr("/"),
				},
			},
			firstMatchPath: "/prefix",
			expected: []transformerPlugin{{
				Type: transformerPluginTypeRequest,
				Replace: TransformerPluginReplaceConfig{
					URI: `$(uri_captures[1] == nil and "/" or uri_captures[1])`,
				},
			}},
			expectedKongRouteModification: kongstate.Route{
				Route: kong.Route{
					Paths: []*string{
						lo.ToPtr("~/prefix$"),
						lo.ToPtr("~/prefix(/.*)"),
					},
				},
			},
		},
		{
			name: "URLRewriteFilter with '/' firstMatchPath and '/' ReplacePrefixPatch",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:               gatewayapi.PrefixMatchHTTPPathModifier,
					ReplacePrefixMatch: lo.ToPtr("/"),
				},
			},
			firstMatchPath: "/",
			expected: []transformerPlugin{{
				Type: transformerPluginTypeRequest,
				Replace: TransformerPluginReplaceConfig{
					URI: `$(uri_captures[1] == nil and "/" or "/" .. uri_captures[1])`,
				},
			}},
			expectedKongRouteModification: kongstate.Route{
				Route: kong.Route{
					Paths: []*string{
						lo.ToPtr("~/$"),
						lo.ToPtr("~/(.*)"),
					},
				},
			},
		},
		{
			name: "URLRewriteFilter with '/' firstMatchPath and non-empty ReplacePrefixPatch",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:               gatewayapi.PrefixMatchHTTPPathModifier,
					ReplacePrefixMatch: lo.ToPtr("/new-prefix"),
				},
			},
			firstMatchPath: "/",
			expected: []transformerPlugin{{
				Type: transformerPluginTypeRequest,
				Replace: TransformerPluginReplaceConfig{
					URI: `/new-prefix$(uri_captures[1] == nil and "" or "/" .. uri_captures[1])`,
				},
			}},
			expectedKongRouteModification: kongstate.Route{
				Route: kong.Route{
					Paths: []*string{
						lo.ToPtr("~/$"),
						lo.ToPtr("~/(.*)"),
					},
				},
			},
		},
		{
			name: "URLRewriteFilter with firstMatchPath with a trailing slash",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:               gatewayapi.PrefixMatchHTTPPathModifier,
					ReplacePrefixMatch: lo.ToPtr("/new-prefix"),
				},
			},
			firstMatchPath: "/prefix/",
			expected: []transformerPlugin{{
				Type: transformerPluginTypeRequest,
				Replace: TransformerPluginReplaceConfig{
					URI: `/new-prefix$(uri_captures[1])`,
				},
			}},
			expectedKongRouteModification: kongstate.Route{
				Route: kong.Route{
					Paths: []*string{
						lo.ToPtr("~/prefix$"),
						lo.ToPtr("~/prefix(/.*)"),
					},
				},
			},
		},
		{
			name: "URLRewriteFilter with hostname",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Hostname: lo.ToPtr(gatewayapi.PreciseHostname("replaced.host")),
			},
			expected: []transformerPlugin{{
				Type: transformerPluginTypeRequest,
				Replace: TransformerPluginReplaceConfig{
					Headers: []string{
						"host:replaced.host",
					},
				},
				Add: TransformerPluginConfig{
					Headers: []string{
						"host:replaced.host",
					},
				},
			}},
		},
		{
			name: "URLRewriteFilter with hostname and path",
			modifier: &gatewayapi.HTTPURLRewriteFilter{
				Path: &gatewayapi.HTTPPathModifier{
					Type:               gatewayapi.PrefixMatchHTTPPathModifier,
					ReplacePrefixMatch: lo.ToPtr("/new-prefix"),
				},
				Hostname: lo.ToPtr(gatewayapi.PreciseHostname("replaced.host")),
			},
			firstMatchPath: "/prefix",
			expected: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,

					Replace: TransformerPluginReplaceConfig{
						URI: `/new-prefix$(uri_captures[1])`,
					},
				},
				{
					Type: transformerPluginTypeRequest,
					Replace: TransformerPluginReplaceConfig{
						Headers: []string{
							"host:replaced.host",
						},
					},
					Add: TransformerPluginConfig{
						Headers: []string{
							"host:replaced.host",
						},
					},
				},
			},
			expectedKongRouteModification: kongstate.Route{
				Route: kong.Route{
					Paths: []*string{
						lo.ToPtr("~/prefix$"),
						lo.ToPtr("~/prefix(/.*)"),
					},
				},
			},
		},
		{
			name:        "nil URLRewriteFilter",
			modifier:    nil,
			expected:    nil,
			expectedErr: errors.New("URLRewrite is not provided"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plugins, routeModifiers, err := generateRequestTransformerForURLRewrite(tc.modifier, tc.firstMatchPath, false)
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expected, plugins)

			route := kongstate.Route{}
			for _, routeModifier := range routeModifiers {
				routeModifier(&route)
			}
			require.Equal(t, tc.expectedKongRouteModification, route)
		})
	}
}

func TestGenerateKongRouteModifierForURLRewritePrefixMatch_ExpressionsRouter(t *testing.T) {
	testCases := []struct {
		name                      string
		path                      string
		expectedRouteModification kongstate.Route
	}{
		{
			name: "root path",
			path: "/",
			expectedRouteModification: kongstate.Route{
				Route: kong.Route{
					Expression: lo.ToPtr(`(http.path == "/") || (http.path ~ "^/(.*)")`),
				},
			},
		},
		{
			name: "prefix path",
			path: "/prefix",
			expectedRouteModification: kongstate.Route{
				Route: kong.Route{
					Expression: lo.ToPtr(`(http.path == "/prefix") || (http.path ~ "^/prefix(/.*)")`),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			modifier := generateKongRouteModifierForURLRewritePrefixMatch(tc.path, true)
			route := kongstate.Route{}
			modifier(&route)
			require.Equal(t, tc.expectedRouteModification, route)
		})
	}
}

func TestMergePluginsOfTheSameType(t *testing.T) {
	testCases := []struct {
		name     string
		plugins  []transformerPlugin
		expected []transformerPlugin
	}{
		{
			name:     "no plugins",
			plugins:  []transformerPlugin{},
			expected: []transformerPlugin{},
		},
		{
			name: "single plugin",
			plugins: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
				},
			},
			expected: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
				},
			},
		},
		{
			name: "multiple plugins of different types",
			plugins: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
				},
				{
					Type: transformerPluginTypeResponse,
				},
			},
			expected: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
				},
				{
					Type: transformerPluginTypeResponse,
				},
			},
		},
		{
			name: "multiple plugins of the same types",
			plugins: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
				},
				{
					Type: transformerPluginTypeResponse,
				},
				{
					Type: transformerPluginTypeRequest,
				},
				{
					Type: transformerPluginTypeResponse,
				},
			},
			expected: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
				},
				{
					Type: transformerPluginTypeResponse,
				},
			},
		},
		{
			name: "multiple plugins of the same types with different configurations - configuration is merged",
			plugins: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
					Add: TransformerPluginConfig{
						Headers: []string{"header1:value1"},
					},
					Replace: TransformerPluginReplaceConfig{
						URI: "path1",
					},
				},
				{
					Type: transformerPluginTypeRequest,
					Add: TransformerPluginConfig{
						Headers: []string{"header2:value2"},
					},
					Replace: TransformerPluginReplaceConfig{
						URI: "path2",
					},
				},
			},
			expected: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
					Add: TransformerPluginConfig{
						Headers: []string{"header1:value1", "header2:value2"},
					},
					Replace: TransformerPluginReplaceConfig{
						URI: "path1",
					},
				},
			},
		},
		{
			name: "multiple plugins of the same types with same configuration keys - configuration is merged and the first wins",
			plugins: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
					Replace: TransformerPluginReplaceConfig{
						URI: "path1",
					},
				},
				{
					Type: transformerPluginTypeRequest,
					Replace: TransformerPluginReplaceConfig{
						URI: "path2",
					},
				},
			},
			expected: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
					Replace: TransformerPluginReplaceConfig{
						URI: "path1",
					},
				},
			},
		},
		{
			name: "multiple plugins of the same types with same configuration keys, first URI empty",
			plugins: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
					Replace: TransformerPluginReplaceConfig{
						URI: "",
					},
				},
				{
					Type: transformerPluginTypeRequest,
					Replace: TransformerPluginReplaceConfig{
						URI: "path2",
					},
				},
			},
			expected: []transformerPlugin{
				{
					Type: transformerPluginTypeRequest,
					Replace: TransformerPluginReplaceConfig{
						URI: "path2",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plugins, err := mergePluginsOfTheSameType(tc.plugins)
			require.NoError(t, err)
			require.Equal(t, tc.expected, plugins)
		})
	}
}
