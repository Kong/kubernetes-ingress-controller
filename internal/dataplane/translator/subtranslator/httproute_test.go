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
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
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
								"Location: http://example.org/test",
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
			options := setKongRoutePluginsOptions{}
			result, err := generatePluginsFromHTTPRouteFilters(tc.filters, tc.path, nil, options)
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
	const (
		longServiceName = "service-with-a-very-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-long-name-having-251-characters"
		// hashStr is the SHA256 sum of the original generated service name for the test case where a long service name is generated and trimmed.
		hashStr = "d39c14b023c01526d6d7a9b4aaf61dbd8daf53eb7241f933daec622ea59e2da9"
	)
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
			name: "no backends",
			ruleMeta: httpRouteRuleMeta{
				Rule:        gatewayapi.HTTPRouteRule{},
				RuleNumber:  0,
				parentRoute: testHTTPRoute,
			},
			expectedServiceName: "httproute.default.svc._",
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
		{
			name: "multiple backends that generates a name exceeding the length limit",
			ruleMeta: httpRouteRuleMeta{
				Rule: gatewayapi.HTTPRouteRule{
					BackendRefs: []gatewayapi.HTTPBackendRef{
						{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Kind: kindService,
									Name: gatewayapi.ObjectName(longServiceName),
									Port: lo.ToPtr(gatewayapi.PortNumber(80)),
								},
								Weight: lo.ToPtr(int32(75)),
							},
						},
						{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Kind: kindService,
									Name: gatewayapi.ObjectName(longServiceName),
									Port: lo.ToPtr(gatewayapi.PortNumber(8080)),
								},
								Weight: lo.ToPtr(int32(25)),
							},
						},
					},
				},
				RuleNumber:  0,
				parentRoute: testHTTPRoute,
			},
			expectedServiceName: fmt.Sprintf("httproute.default.svc.default.%s.80.75_combined.%s", longServiceName, hashStr),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedServiceName, tc.ruleMeta.getKongServiceNameByBackendRefs())
		})
	}
}

func TestTranslateHTTPRoutesToKongstateServices(t *testing.T) {
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
		name             string
		k8sServices      []*corev1.Service
		referenceGrants  []*gatewayapi.ReferenceGrant
		httpRoutes       []*gatewayapi.HTTPRoute
		expectedServices map[string]kongstate.Service
	}{
		{
			name: "multiple rules in one HTTPRoute sharing the same backends",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
				},
			},
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
							},
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedServices: map[string]kongstate.Service{
				"httproute.default.svc.default.service-1.80": {
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-1.80"),
						Host: kong.String("httproute.default.svc.default.service-1.80"),
					},
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
				},
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-2",
						Namespace: "default",
					},
				},
			},
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
									builder.NewHTTPBackendRef("service-2").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedServices: map[string]kongstate.Service{
				"httproute.default.svc.default.service-1.80_default.service-2.80": {
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-1.80_default.service-2.80"),
						Host: kong.String("httproute.default.svc.default.service-1.80_default.service-2.80"),
					},
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
		},
		{
			name: "multiple HTTPRoutes with the same same backends in the same namespace",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
				},
			},
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-1",
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-2",
						Namespace: "default",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedServices: map[string]kongstate.Service{
				"httproute.default.svc.default.service-1.80": {
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-1.80"),
						Host: kong.String("httproute.default.svc.default.service-1.80"),
					},
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
		},
		{
			name: "multiple HTTPRoutes with the same same backends in the same namespace with correct referenceGrant",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
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
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-1",
						Namespace: "another-namespace",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithNamespace("default").WithPort(80).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-2",
						Namespace: "another-namespace",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithNamespace("default").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedServices: map[string]kongstate.Service{
				"httproute.another-namespace.svc.default.service-1.80": {
					Service: kong.Service{
						Name: kong.String("httproute.another-namespace.svc.default.service-1.80"),
						Host: kong.String("httproute.another-namespace.svc.default.service-1.80"),
					},
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
		},
		{
			name: "multiple HTTPRoutes with the same same backends in the same namespace without referenceGrant",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
				},
			},
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-1",
						Namespace: "another-namespace",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithNamespace("default").WithPort(80).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httproute-2",
						Namespace: "another-namespace",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithNamespace("default").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedServices: map[string]kongstate.Service{
				"httproute.another-namespace.svc.default.service-1.80": {
					Service: kong.Service{
						Name: kong.String("httproute.another-namespace.svc.default.service-1.80"),
						Host: kong.String("httproute.another-namespace.svc.default.service-1.80"),
					},
					Backends: []kongstate.ServiceBackend{},
				},
			},
		},
		{
			name: "multiple rules with different backends in the same HTTPRoute",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
				},
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-2",
						Namespace: "default",
					},
				},
			},
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
							},
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-2").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedServices: map[string]kongstate.Service{
				"httproute.default.svc.default.service-1.80": {
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-1.80"),
						Host: kong.String("httproute.default.svc.default.service-1.80"),
					},
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
				"httproute.default.svc.default.service-2.80": {
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-2.80"),
						Host: kong.String("httproute.default.svc.default.service-2.80"),
					},
					Backends: []kongstate.ServiceBackend{
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
		},
		{
			name: "rules sharing backends from multiple HTTPRoutes in the different namespaces",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
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
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
							},
						},
					},
				},
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "another-namespace",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithNamespace("default").WithPort(80).Build(),
								},
							},
						},
					},
				},
			},
			expectedServices: map[string]kongstate.Service{
				"httproute.default.svc.default.service-1.80": {
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-1.80"),
						Host: kong.String("httproute.default.svc.default.service-1.80"),
					},
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
				"httproute.another-namespace.svc.default.service-1.80": {
					Service: kong.Service{
						Name: kong.String("httproute.another-namespace.svc.default.service-1.80"),
						Host: kong.String("httproute.another-namespace.svc.default.service-1.80"),
					},
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
		},
		{
			name: "HTTPRoute with ExtensionRef plugin",
			k8sServices: []*corev1.Service{
				{
					TypeMeta: serviceTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
				},
			},
			httpRoutes: []*gatewayapi.HTTPRoute{
				{
					TypeMeta: httpRouteTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "httproute-1",
					},
					Spec: gatewayapi.HTTPRouteSpec{
						CommonRouteSpec: commonRouteSpecMock("fake-gateway-1"),
						Rules: []gatewayapi.HTTPRouteRule{
							{
								BackendRefs: []gatewayapi.HTTPBackendRef{
									builder.NewHTTPBackendRef("service-1").WithPort(80).Build(),
								},
								Matches: []gatewayapi.HTTPRouteMatch{
									{
										Path: &gatewayapi.HTTPPathMatch{
											Type:  lo.ToPtr(gatewayapi.PathMatchExact),
											Value: lo.ToPtr("/foo"),
										},
									},
								},
								Filters: []gatewayapi.HTTPRouteFilter{
									{
										Type: gatewayapi.HTTPRouteFilterExtensionRef,
										ExtensionRef: &gatewayapi.LocalObjectReference{
											Name:  "plugin-1",
											Kind:  "KongPlugin",
											Group: "configuration.konghq.com",
										},
									},
								},
							},
						},
					},
				},
			},
			expectedServices: map[string]kongstate.Service{
				"httproute.default.svc.default.service-1.80": {
					Service: kong.Service{
						Name: kong.String("httproute.default.svc.default.service-1.80"),
						Host: kong.String("httproute.default.svc.default.service-1.80"),
					},
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

			oldHTTPRoutes := make([]*gatewayapi.HTTPRoute, 0, len(tc.httpRoutes))
			for _, r := range tc.httpRoutes {
				oldHTTPRoutes = append(oldHTTPRoutes, r.DeepCopy())
			}

			translateOptions := TranslateHTTPRouteToKongstateServiceOptions{
				CombinedServicesFromDifferentHTTPRoutes: true,
				ExpressionRoutes:                        false,
				SupportRedirectPlugin:                   false,
			}
			translationResult := TranslateHTTPRoutesToKongstateServices(logger, fakestore, tc.httpRoutes, translateOptions)
			require.Len(t, translationResult.HTTPRouteNameToTranslationErrors, 0, "Should not get translation errors in translating")

			kongstateServices := translationResult.ServiceNameToKongstateService
			require.Len(t, kongstateServices, len(tc.expectedServices))
			for serviceName, expectedService := range tc.expectedServices {
				s, ok := kongstateServices[serviceName]
				require.Truef(t, ok, "Should find service %s in translated services", serviceName)
				// compare the name and host of the translated service.
				require.Equal(t, *expectedService.Name, *s.Name, "Service %s should have expected name inside", serviceName)
				require.Equal(t, *expectedService.Host, *s.Host, "Service %s should have expected host", serviceName)
				// compare backends.
				require.Lenf(t, s.Backends, len(expectedService.Backends), "Service %s should have expected number of backends")
				backendCompareMsg := `Service %s backend %d should have expected %s` // service name, backend index, field name
				for i, expectedBackend := range expectedService.Backends {
					require.Equalf(t, s.Backends[i].Namespace(), expectedBackend.Namespace(), backendCompareMsg, serviceName, i, "namespace")
					require.Equal(t, s.Backends[i].Name(), expectedBackend.Name(), backendCompareMsg, serviceName, i, "name")
					require.Equal(t, s.Backends[i].PortDef(), expectedBackend.PortDef(), backendCompareMsg, serviceName, i, "port")
					require.Equal(t, s.Backends[i].Weight(), expectedBackend.Weight(), backendCompareMsg, serviceName, i, "weight")
				}
			}

			require.Equal(t, oldHTTPRoutes, tc.httpRoutes, "HTTPRoutes should not be modified")
		})
	}
}

func TestTranslateHTTPRouteRulesMetaToKongstateRoutes(t *testing.T) {
	httpRouteTypeMeta := metav1.TypeMeta{
		APIVersion: "gateway.networking.k8s.io/v1",
		Kind:       "HTTPRoute",
	}
	backendRefList := []gatewayapi.HTTPBackendRef{
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
	}
	httpRouteWithoutHost := &gatewayapi.HTTPRoute{
		TypeMeta: httpRouteTypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "httproute-1",
		},
	}

	testCases := []struct {
		name           string
		rulesMeta      []httpRouteRuleMeta
		expectedRoutes []kongstate.Route
		expectError    bool
	}{
		{
			name: "multiple rules with the same(empty) filter from the same HTTPRoute",
			rulesMeta: []httpRouteRuleMeta{
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: backendRefList,
						Matches: []gatewayapi.HTTPRouteMatch{
							{
								Path: &gatewayapi.HTTPPathMatch{
									Type:  lo.ToPtr(gatewayapi.PathMatchExact),
									Value: lo.ToPtr("/foo"),
								},
							},
							{
								Path: &gatewayapi.HTTPPathMatch{
									Type:  lo.ToPtr(gatewayapi.PathMatchExact),
									Value: lo.ToPtr("/bar"),
								},
							},
						},
					},
					RuleNumber:  0,
					parentRoute: httpRouteWithoutHost,
				},
				{
					Rule: gatewayapi.HTTPRouteRule{
						BackendRefs: backendRefList,
						Matches: []gatewayapi.HTTPRouteMatch{
							{
								Path: &gatewayapi.HTTPPathMatch{
									Type:  lo.ToPtr(gatewayapi.PathMatchExact),
									Value: lo.ToPtr("/baz"),
								},
							},
						},
					},
					RuleNumber:  1,
					parentRoute: httpRouteWithoutHost,
				},
			},
			expectedRoutes: []kongstate.Route{
				{
					Route: kong.Route{
						Name:         kong.String("httproute.default.httproute-1.0.0"),
						Paths:        kong.StringSlice("~/foo$", "~/bar$", "~/baz$"),
						PreserveHost: kong.Bool(true),
						StripPath:    kong.Bool(false),
						Protocols: []*string{
							kong.String("http"),
							kong.String("https"),
						},
						Tags: []*string{
							kong.String("k8s-name:httproute-1"),
							kong.String("k8s-namespace:default"),
							kong.String("k8s-kind:HTTPRoute"),
							kong.String("k8s-group:gateway.networking.k8s.io"),
							kong.String("k8s-version:v1"),
						},
					},
					Ingress: util.FromK8sObject(httpRouteWithoutHost),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generateOptions := TranslateHTTPRouteRulesToKongRouteOptions{
				ExpressionRoutes:      false,
				SupportRedirectPlugin: false,
			}
			routes, err := translateHTTPRouteRulesMetaToKongstateRoutes(tc.rulesMeta, generateOptions)
			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, routes, len(tc.expectedRoutes))
			for i, expectedRoute := range tc.expectedRoutes {
				require.Equal(t, expectedRoute, routes[i])
			}
		})
	}
}

func TestSchemeHostPortFromHTTPPathModifier(t *testing.T) {
	testCases := []struct {
		name             string
		modifier         *gatewayapi.HTTPRequestRedirectFilter
		expectedScheme   string
		expectedHostPort string
	}{
		{
			name: "Has scheme, host and port",
			modifier: &gatewayapi.HTTPRequestRedirectFilter{
				Scheme:   lo.ToPtr("https"),
				Hostname: lo.ToPtr(gatewayapi.PreciseHostname("a.com")),
				Port:     lo.ToPtr(gatewayapi.PortNumber(8443)),
			},
			expectedScheme:   "https",
			expectedHostPort: "a.com:8443",
		},
		{
			name: "no scheme, http scheme should be returned",
			modifier: &gatewayapi.HTTPRequestRedirectFilter{
				Hostname: lo.ToPtr(gatewayapi.PreciseHostname("a.com")),
			},
			expectedScheme:   "http",
			expectedHostPort: "a.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			urlScheme, hostPort := schemeHostPortFromHTTPPathModifier(tc.modifier)
			require.Equal(t, tc.expectedScheme, urlScheme)
			require.Equal(t, tc.expectedHostPort, hostPort)
		})
	}
}

func TestGenerateRequestRedirectUsingRedirectKongPlugin(t *testing.T) {
	testCases := []struct {
		name               string
		modifier           *gatewayapi.HTTPRequestRedirectFilter
		expectedKongPlugin kong.Plugin
	}{
		{
			name: "full path replace",
			modifier: &gatewayapi.HTTPRequestRedirectFilter{
				StatusCode: lo.ToPtr(301),
				Hostname:   lo.ToPtr(gatewayapi.PreciseHostname("a.com")),
				Path: &gatewayapi.HTTPPathModifier{
					Type:            gatewayapi.FullPathHTTPPathModifier,
					ReplaceFullPath: lo.ToPtr("/foo"),
				},
			},
			expectedKongPlugin: kong.Plugin{
				Name: kong.String("redirect"),
				Config: kong.Configuration{
					"status_code":        lo.ToPtr(301),
					"location":           lo.ToPtr("http://a.com/foo"),
					"keep_incoming_path": lo.ToPtr(false),
				},
			},
		},
		{
			name: "no path replace",
			modifier: &gatewayapi.HTTPRequestRedirectFilter{
				StatusCode: lo.ToPtr(301),
				Hostname:   lo.ToPtr(gatewayapi.PreciseHostname("a.com")),
				Scheme:     lo.ToPtr("http"),
			},
			expectedKongPlugin: kong.Plugin{
				Name: kong.String("redirect"),
				Config: kong.Configuration{
					"status_code":        lo.ToPtr(301),
					"location":           lo.ToPtr("http://a.com"),
					"keep_incoming_path": lo.ToPtr(true),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedKongPlugin, generateRequestRedirectUsingRedirectKongPlugin(tc.modifier))
		})
	}
}
