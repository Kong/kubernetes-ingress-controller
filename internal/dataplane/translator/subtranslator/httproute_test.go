package subtranslator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
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

func TestConvertGatewayMatchHeadersToKongRouteMatchHeaders(t *testing.T) {
	t.Log("generating several gateway header matches")
	tests := []struct {
		msg    string
		input  []gatewayapi.HTTPHeaderMatch
		output map[string][]string
		err    error
	}{
		{
			msg: "regex header matches convert correctly",
			input: []gatewayapi.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayapi.HeaderMatchRegularExpression),
				Name:  "Content-Type",
				Value: "^audio/*",
			}},
			output: map[string][]string{
				"Content-Type": {KongHeaderRegexPrefix + "^audio/*"},
			},
		},
		{
			msg: "a single exact header match with no type defaults to exact type and converts properly",
			input: []gatewayapi.HTTPHeaderMatch{{
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
		{
			msg: "a single exact header match with a single value converts properly",
			input: []gatewayapi.HTTPHeaderMatch{{
				Type:  lo.ToPtr(gatewayapi.HeaderMatchExact),
				Name:  "Content-Type",
				Value: "audio/vorbis",
			}},
			output: map[string][]string{
				"Content-Type": {"audio/vorbis"},
			},
		},
		{
			msg: "multiple header matches for the same header are rejected",
			input: []gatewayapi.HTTPHeaderMatch{
				{
					Name:  "Content-Type",
					Value: "audio/vorbis",
				},
				{
					Name:  "Content-Type",
					Value: "audio/flac",
				},
			},
			output: nil,
			err:    fmt.Errorf("multiple header matches for the same header are not allowed: Content-Type"),
		},
		{
			msg: "multiple header matches convert properly",
			input: []gatewayapi.HTTPHeaderMatch{
				{
					Type:  lo.ToPtr(gatewayapi.HeaderMatchExact),
					Name:  "Content-Type",
					Value: "audio/vorbis",
				},
				{
					Name:  "Content-Length",
					Value: "999999999",
				},
			},
			output: map[string][]string{
				"Content-Type":   {"audio/vorbis"},
				"Content-Length": {"999999999"},
			},
		},
		{
			msg:    "an empty list of headers will produce no converted headers",
			output: map[string][]string{},
		},
	}

	t.Log("verifying header match conversions")
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			output, err := convertGatewayMatchHeadersToKongRouteMatchHeaders(tt.input)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.output, output)
		})
	}
}

func TestGetKongServiceNameByBackendRefs(t *testing.T) {
	testHTTPRoute := &gatewayapi.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-httproute",
			Namespace: "default",
		},
	}
	kindService := lo.ToPtr(gatewayapi.Kind("Service"))
	testCases := []struct {
		name                string
		ruleMeta            httpRouteRuleMeta
		expectedServiceName string
	}{
		{
			name: "single backend",
			ruleMeta: httpRouteRuleMeta{
				Rule: gatewayapi.HTTPRouteRule{
					BackendRefs: []gatewayapi.HTTPBackendRef{
						{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Kind: kindService,
									Name: gatewayapi.ObjectName("service-1"),
								},
							},
						},
					},
				},
				RuleNumber:  0,
				parentRoute: testHTTPRoute,
			},
			expectedServiceName: "httproute.default.svc.default.service-1",
		},
		{
			name: "multiple backends",
			ruleMeta: httpRouteRuleMeta{
				Rule: gatewayapi.HTTPRouteRule{
					BackendRefs: []gatewayapi.HTTPBackendRef{
						{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Kind: kindService,
									Name: gatewayapi.ObjectName("service-1"),
									Port: lo.ToPtr(gatewayapi.PortNumber(80)),
								},
							},
						},
						{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Kind: kindService,
									Name: gatewayapi.ObjectName("service-2"),
									Port: lo.ToPtr(gatewayapi.PortNumber(8080)),
								},
							},
						},
						{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Kind:      kindService,
									Name:      gatewayapi.ObjectName("service-2"),
									Namespace: lo.ToPtr(gatewayapi.Namespace("another-namespace")),
									Port:      lo.ToPtr(gatewayapi.PortNumber(80)),
								},
							},
						},
					},
				},
				RuleNumber:  0,
				parentRoute: testHTTPRoute,
			},
			expectedServiceName: "httproute.default.svc.another-namespace.service-2.80_default.service-1.80_default.service-2.8080",
		},
		{
			name: "multiple backends with weights",
			ruleMeta: httpRouteRuleMeta{
				Rule: gatewayapi.HTTPRouteRule{
					BackendRefs: []gatewayapi.HTTPBackendRef{
						{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Kind: kindService,
									Name: gatewayapi.ObjectName("service-1"),
									Port: lo.ToPtr(gatewayapi.PortNumber(80)),
								},
								Weight: lo.ToPtr(int32(75)),
							},
						},
						{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Kind: kindService,
									Name: gatewayapi.ObjectName("service-1"),
									Port: lo.ToPtr(gatewayapi.PortNumber(1080)),
								},
								Weight: lo.ToPtr(int32(25)),
							},
						},
					},
				},
				RuleNumber:  0,
				parentRoute: testHTTPRoute,
			},
			expectedServiceName: "httproute.default.svc.default.service-1.1080.25_default.service-1.80.75",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedServiceName, tc.ruleMeta.getKongServiceNameByBackendRefs())
		})
	}
}

func TestTranslateHTTPRouteRulesMetaToKongstateService(t *testing.T) {
	serviceTypeMeta := metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	}

	refernceGrantTypeMeta := metav1.TypeMeta{
		APIVersion: "gateway.networking.k8s.io/v1beta1",
		Kind:       "ReferenceGrant",
	}

	httpRouteTypeMeta := metav1.TypeMeta{
		APIVersion: "gateway.networking.k8s.io/v1",
		Kind:       "HTTPRoute",
	}

	mustNewKongstateServiceBackend := func(
		typ kongstate.ServiceBackendType,
		nn k8stypes.NamespacedName,
		portDef kongstate.PortDef,
		weight *int32,
	) kongstate.ServiceBackend {
		b, err := kongstate.NewServiceBackend(
			typ,
			nn,
			portDef,
		)
		require.NoError(t, err)
		if weight != nil {
			b.SetWeight(*weight)
		}
		return b
	}

	testCases := []struct {
		name            string
		k8sServices     []*corev1.Service
		referenceGrants []*gatewayapi.ReferenceGrant
		rulesMeta       []httpRouteRuleMeta
		expectedService kongstate.Service
	}{
		{
			name: "Multiple rules in one HTTPRoute",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: int32(80),
							},
						},
					},
				},
			},
			rulesMeta: []httpRouteRuleMeta{
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
										Name: gatewayapi.ObjectName("service-1"),
										Port: lo.ToPtr(gatewayapi.PortNumber(80)),
									},
								},
							},
						},
					},
					RuleNumber: 0,
					parentRoute: &gatewayapi.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-1",
						},
					},
				},
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
										Name: gatewayapi.ObjectName("service-1"),
										Port: lo.ToPtr(gatewayapi.PortNumber(80)),
									},
								},
							},
						},
					},
					RuleNumber: 1,
					parentRoute: &gatewayapi.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-1",
						},
					},
				},
			},
			expectedService: kongstate.Service{
				Backends: []kongstate.ServiceBackend{
					mustNewKongstateServiceBackend(
						kongstate.ServiceBackendTypeKubernetesService,
						k8stypes.NamespacedName{
							Name:      "service-1",
							Namespace: "default",
						},
						kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: 80,
						},
						nil,
					),
				},
			},
		},
		{
			name: "multiple backends in one rule of one HTTPRoute",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: int32(80),
							},
						},
					},
				},
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-2",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: int32(80),
							},
						},
					},
				},
			},
			rulesMeta: []httpRouteRuleMeta{
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
										Name: gatewayapi.ObjectName("service-1"),
										Port: lo.ToPtr(gatewayapi.PortNumber(80)),
									},
								},
							},
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
										Name: gatewayapi.ObjectName("service-2"),
										Port: lo.ToPtr(gatewayapi.PortNumber(80)),
									},
								},
							},
						},
					},
					RuleNumber: 0,
					parentRoute: &gatewayapi.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-1",
						},
					},
				},
			},
			expectedService: kongstate.Service{
				Backends: []kongstate.ServiceBackend{
					mustNewKongstateServiceBackend(
						kongstate.ServiceBackendTypeKubernetesService,
						k8stypes.NamespacedName{
							Name:      "service-1",
							Namespace: "default",
						},
						kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: 80,
						},
						nil,
					),
					mustNewKongstateServiceBackend(
						kongstate.ServiceBackendTypeKubernetesService,
						k8stypes.NamespacedName{
							Name:      "service-2",
							Namespace: "default",
						},
						kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: 80,
						},
						nil,
					),
				},
			},
		},
		{
			name: "rules from multiple HTTPRoutes in the same namespace",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: int32(80),
							},
						},
					},
				},
			},
			rulesMeta: []httpRouteRuleMeta{
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
										Name: gatewayapi.ObjectName("service-1"),
										Port: lo.ToPtr(gatewayapi.PortNumber(80)),
									},
								},
							},
						},
					},
					RuleNumber: 0,
					parentRoute: &gatewayapi.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-1",
						},
					},
				},
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind: lo.ToPtr(gatewayapi.Kind("Service")),
										Name: gatewayapi.ObjectName("service-1"),
										Port: lo.ToPtr(gatewayapi.PortNumber(80)),
									},
								},
							},
						},
					},
					RuleNumber: 1,
					parentRoute: &gatewayapi.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "httproute-2",
						},
					},
				},
			},
			expectedService: kongstate.Service{
				Backends: []kongstate.ServiceBackend{
					mustNewKongstateServiceBackend(
						kongstate.ServiceBackendTypeKubernetesService,
						k8stypes.NamespacedName{
							Name:      "service-1",
							Namespace: "default",
						},
						kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: 80,
						},
						nil,
					),
				},
			},
		},
		{
			name: "rules from multiple HTTPRoutes in the same namespace with correct referenceGrant",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: int32(80),
							},
						},
					},
				},
			},
			referenceGrants: []*gatewayapi.ReferenceGrant{
				{
					TypeMeta: refernceGrantTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "grant-from-httproute-to-service",
					},
					Spec: gatewayapi.ReferenceGrantSpec{
						From: []gatewayapi.ReferenceGrantFrom{
							{
								Namespace: "another-namespace",
								Kind:      "HTTPRoute",
								Group:     gatewayapi.V1Group,
							},
						},
						To: []gatewayapi.ReferenceGrantTo{
							{
								Group: corev1.GroupName,
								Kind:  "Service",
							},
						},
					},
				},
			},
			rulesMeta: []httpRouteRuleMeta{
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
										Name:      gatewayapi.ObjectName("service-1"),
										Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
										Port:      lo.ToPtr(gatewayapi.PortNumber(80)),
									},
								},
							},
						},
					},
					RuleNumber: 0,
					parentRoute: &gatewayapi.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "another-namespace",
							Name:      "httproute-1",
						},
					},
				},
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
										Name:      gatewayapi.ObjectName("service-1"),
										Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
										Port:      lo.ToPtr(gatewayapi.PortNumber(80)),
									},
								},
							},
						},
					},
					RuleNumber: 0,
					parentRoute: &gatewayapi.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "another-namespace",
							Name:      "httproute-2",
						},
					},
				},
			},
			expectedService: kongstate.Service{
				Backends: []kongstate.ServiceBackend{
					mustNewKongstateServiceBackend(
						kongstate.ServiceBackendTypeKubernetesService,
						k8stypes.NamespacedName{
							Name:      "service-1",
							Namespace: "default",
						},
						kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: 80,
						},
						nil,
					),
				},
			},
		},
		{
			name: "rules from multiple HTTPRoutes in the same namespace without referenceGrant",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: int32(80),
							},
						},
					},
				},
			},
			rulesMeta: []httpRouteRuleMeta{
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
										Name:      gatewayapi.ObjectName("service-1"),
										Namespace: lo.ToPtr(gatewayapi.Namespace("default")),
										Port:      lo.ToPtr(gatewayapi.PortNumber(80)),
									},
								},
							},
						},
					},
					RuleNumber: 0,
					parentRoute: &gatewayapi.HTTPRoute{
						TypeMeta: httpRouteTypeMeta,
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "another-namespace",
							Name:      "httproute-1",
						},
					},
				},
			},
			// No backends should be generated.
			expectedService: kongstate.Service{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logr.Discard()
			fakestore, err := store.NewFakeStore(store.FakeObjects{
				Services:        tc.k8sServices,
				ReferenceGrants: tc.referenceGrants,
			})
			require.NoError(t, err)

			serviceName := tc.rulesMeta[0].getKongServiceNameByBackendRefs()

			kongService, err := translateHTTPRouteRulesMetaToKongstateService(logger, fakestore, serviceName, tc.rulesMeta)
			require.NoError(t, err)
			require.Equal(t, serviceName, *kongService.Name)
			require.Len(t, kongService.Backends, len(tc.expectedService.Backends))
			for i, expectedBackend := range tc.expectedService.Backends {
				require.Equal(t, expectedBackend.Namespace(), kongService.Backends[i].Namespace())
				require.Equal(t, expectedBackend.Name(), kongService.Backends[i].Name())
				require.Equal(t, expectedBackend.PortDef(), kongService.Backends[i].PortDef())
				require.Equal(t, expectedBackend.Weight(), kongService.Backends[i].Weight())
			}
		})
	}
}
