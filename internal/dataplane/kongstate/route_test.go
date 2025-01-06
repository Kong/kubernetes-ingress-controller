package kongstate

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
)

func TestOverrideRoute(t *testing.T) {
	testCases := []struct {
		name          string
		inRoute       Route
		inAnnotations map[string]string
		expectedRoute Route
	}{
		{
			name: "no annotations",
			inRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			inAnnotations: map[string]string{},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
		},
		{
			name: "override methods",
			inRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			inAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.MethodsKey: "GET,POST",
			},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts:   kong.StringSlice("foo.com", "bar.com"),
					Methods: kong.StringSlice("GET", "POST"),
				},
			},
		},
		{
			name: "override methods case insensitive",
			inRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			inAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.MethodsKey: "GET,post",
			},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts:   kong.StringSlice("foo.com", "bar.com"),
					Methods: kong.StringSlice("GET", "POST"),
				},
			},
		},
		{
			name: "override methods with invalid method",
			inRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			inAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.MethodsKey: "GET,-1",
			},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
		},
		{
			name: "override https redirect status code",
			inRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			inAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.HTTPSRedirectCodeKey: "302",
			},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts:                   kong.StringSlice("foo.com", "bar.com"),
					HTTPSRedirectStatusCode: kong.Int(302),
				},
			},
		},
		{
			name: "override protocols, preserve host, strip path and regex priority",
			inRoute: Route{
				Route: kong.Route{
					Hosts:        kong.StringSlice("foo.com", "bar.com"),
					PreserveHost: kong.Bool(true),
					StripPath:    kong.Bool(true),
				},
			},
			inAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.ProtocolsKey:     "http",
				annotations.AnnotationPrefix + annotations.PreserveHostKey:  "false",
				annotations.AnnotationPrefix + annotations.StripPathKey:     "false",
				annotations.AnnotationPrefix + annotations.RegexPriorityKey: "10",
			},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts:         kong.StringSlice("foo.com", "bar.com"),
					Protocols:     kong.StringSlice("http"),
					PreserveHost:  kong.Bool(false),
					StripPath:     kong.Bool(false),
					RegexPriority: kong.Int(10),
				},
			},
		},
		{
			name: "override headers",
			inRoute: Route{
				Route: kong.Route{
					Hosts:     kong.StringSlice("foo.com"),
					Protocols: kong.StringSlice("http", "https"),
				},
			},
			inAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.HeadersKey + ".foo-header": "bar-value",
			},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts:     kong.StringSlice("foo.com"),
					Protocols: kong.StringSlice("http", "https"),
					Headers: map[string][]string{
						"foo-header": {"bar-value"},
					},
				},
			},
		},
		{
			name: "override protocols",
			inRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com"),
				},
			},
			inAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.ProtocolsKey: "grpc,grpcs",
			},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts:     kong.StringSlice("foo.com"),
					Protocols: kong.StringSlice("grpc", "grpcs"),
					StripPath: nil,
				},
			},
		},
		{
			name: "override path handling",
			inRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com"),
				},
			},
			inAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.PathHandlingKey: "v1",
			},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts:        kong.StringSlice("foo.com"),
					PathHandling: kong.String("v1"),
				},
			},
		},
		{
			name: "override request/response buffering",
			inRoute: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			inAnnotations: map[string]string{
				annotations.AnnotationPrefix + annotations.RequestBuffering:  "true",
				annotations.AnnotationPrefix + annotations.ResponseBuffering: "true",
			},
			expectedRoute: Route{
				Route: kong.Route{
					Hosts:             kong.StringSlice("foo.com", "bar.com"),
					RequestBuffering:  kong.Bool(true),
					ResponseBuffering: kong.Bool(true),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			route := tc.inRoute
			route.Ingress = util.K8sObjectInfo{
				Annotations: tc.inAnnotations,
			}
			route.overrideByAnnotation(zapr.NewLogger(zap.NewNop()))
			require.Equal(t, tc.expectedRoute.Route, route.Route)
		})
	}
}

func TestNilRouteOverrideDoesntPanic(t *testing.T) {
	require.NotPanics(t, func() {
		var nilRoute *Route
		nilRoute.override(zapr.NewLogger(zap.NewNop()))
	})
}

func TestOverrideExpressionRoute(t *testing.T) {
	testCases := []struct {
		name     string
		inRoute  Route
		outRoute Route
	}{
		{
			name: "protocols should be overridden, but hosts, method, headers, snis should not",
			inRoute: Route{
				Route: kong.Route{
					Name:       kong.String("expression-route-1"),
					Expression: kong.String(`(http.host == "foo.com") && (http.path ^= "/v1/api")`),
				},
				Ingress: util.K8sObjectInfo{
					Annotations: map[string]string{
						"konghq.com/protocols":    "https",
						"konghq.com/method":       "GET",
						"konghq.com/host-aliases": "bar.com",
						"konghq.com/headers.foo":  "bar",
						"kohghq.com/snis":         "foo.com,bar.com",
					},
				},
				ExpressionRoutes: true,
			},
			outRoute: Route{
				Route: kong.Route{
					Name:       kong.String("expression-route-1"),
					Expression: kong.String(`(http.host == "foo.com") && (http.path ^= "/v1/api")`),
					Protocols:  kong.StringSlice("https"),
				},
				Ingress: util.K8sObjectInfo{
					Annotations: map[string]string{
						"konghq.com/protocols":    "https",
						"konghq.com/method":       "GET",
						"konghq.com/host-aliases": "bar.com",
						"konghq.com/headers.foo":  "bar",
						"kohghq.com/snis":         "foo.com,bar.com",
					},
				},
				ExpressionRoutes: true,
			},
		},
		{
			name: "request_buffering, response_buffering should be overridden, but regex_priority, path_handling should not",
			inRoute: Route{
				Route: kong.Route{
					Name:       kong.String("expression-route-2"),
					Expression: kong.String(`(http.host == "foo.com") && (http.path ^= "/v1/api")`),
				},
				Ingress: util.K8sObjectInfo{
					Annotations: map[string]string{
						"konghq.com/request-buffering":  "true",
						"konghq.com/response-buffering": "true",
						"konghq.com/regex-priority":     "100",
						"konghq.com/path-handling":      "v1",
					},
				},
				ExpressionRoutes: true,
			},
			outRoute: Route{
				Route: kong.Route{
					Name:              kong.String("expression-route-2"),
					Expression:        kong.String(`(http.host == "foo.com") && (http.path ^= "/v1/api")`),
					RequestBuffering:  kong.Bool(true),
					ResponseBuffering: kong.Bool(true),
				},
				Ingress: util.K8sObjectInfo{
					Annotations: map[string]string{
						"konghq.com/request-buffering":  "true",
						"konghq.com/response-buffering": "true",
						"konghq.com/regex-priority":     "100",
						"konghq.com/path-handling":      "v1",
					},
				},
				ExpressionRoutes: true,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			tc.inRoute.override(zapr.NewLogger(zap.NewNop()))
			assert.Equal(t, tc.outRoute, tc.inRoute, "should be the same as expected after overriding")
		})
	}
}

func TestNormalizeProtocols(t *testing.T) {
	assert := assert.New(t)
	testTable := []struct {
		inRoute  Route
		outRoute Route
	}{
		{
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("grpc", "grpcs"),
				},
			},
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("grpc", "grpcs"),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("http", "https"),
				},
			},
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("http", "https"),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("grpc", "https"),
				},
			},
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("http", "https"),
				},
			},
		},
	}

	for _, testcase := range testTable {
		testcase.inRoute.normalizeProtocols()
		assert.Equal(testcase.inRoute.Protocols, testcase.outRoute.Protocols)
	}

	assert.NotPanics(func() {
		var nilUpstream *Upstream
		nilUpstream.override(nil, nil)
	})
}

func TestUseSSLProtocol(t *testing.T) {
	assert := assert.New(t)
	testTable := []struct {
		inRoute  Route
		outRoute kong.Route
	}{
		{
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("grpc", "grpcs"),
				},
			},
			kong.Route{
				Protocols: kong.StringSlice("grpcs"),
			},
		},
		{
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("http", "https"),
				},
			},
			kong.Route{
				Protocols: kong.StringSlice("https"),
			},
		},
		{
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("grpcs", "https"),
				},
			},

			kong.Route{
				Protocols: kong.StringSlice("grpcs", "https"),
			},
		},
		{
			Route{
				Route: kong.Route{
					Protocols: kong.StringSlice("grpc", "http"),
				},
			},
			kong.Route{
				Protocols: kong.StringSlice("grpcs", "https"),
			},
		},
		{
			Route{
				Route: kong.Route{
					Protocols: []*string{},
				},
			},
			kong.Route{
				Protocols: kong.StringSlice("https"),
			},
		},
	}

	for _, testcase := range testTable {
		testcase.inRoute.useSSLProtocol()
		assert.Equal(testcase.inRoute.Protocols, testcase.outRoute.Protocols)
	}
}

func TestOverrideRouteStripPath(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want kong.Route
	}{
		{},
		{
			name: "basic empty route",
			args: args{
				route: Route{Route: kong.Route{}},
			},
			want: kong.Route{},
		},
		{
			name: "set to false",
			args: args{
				route: Route{
					Route: kong.Route{},
				},
				anns: map[string]string{
					"konghq.com/strip-path": "false",
				},
			},
			want: kong.Route{
				StripPath: kong.Bool(false),
			},
		},
		{
			name: "set to true and case insensitive",
			args: args{
				route: Route{
					Route: kong.Route{},
				},
				anns: map[string]string{
					"konghq.com/strip-path": "truE",
				},
			},
			want: kong.Route{
				StripPath: kong.Bool(true),
			},
		},
		{
			name: "overrides any other value",
			args: args{
				route: Route{
					Route: kong.Route{
						StripPath: kong.Bool(false),
					},
				},
				anns: map[string]string{
					"konghq.com/strip-path": "truE",
				},
			},
			want: kong.Route{
				StripPath: kong.Bool(true),
			},
		},
		{
			name: "random value",
			args: args{
				route: Route{
					Route: kong.Route{
						StripPath: kong.Bool(false),
					},
				},
				anns: map[string]string{
					"konghq.com/strip-path": "42",
				},
			},
			want: kong.Route{
				StripPath: kong.Bool(false),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overrideStripPath(tt.args.anns)
			if !reflect.DeepEqual(tt.args.route.Route, tt.want) {
				t.Errorf("overrideRouteStripPath() got = %v, want %v", &tt.args.route.Route, tt.want)
			}
		})
	}
}

func TestOverrideRouteHTTPSRedirectCode(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Route
	}{
		{name: "basic empty route"},
		{
			name: "basic sanity",
			args: args{
				anns: map[string]string{
					"konghq.com/https-redirect-status-code": "301",
				},
			},
			want: Route{
				Route: kong.Route{
					HTTPSRedirectStatusCode: kong.Int(301),
				},
			},
		},
		{
			name: "random integer value",
			args: args{
				anns: map[string]string{
					"konghq.com/https-redirect-status-code": "42",
				},
			},
		},
		{
			name: "random string",
			args: args{
				anns: map[string]string{
					"konghq.com/https-redirect-status-code": "foo",
				},
			},
		},
		{
			name: "force ssl annotation set to true and protocols is not set",
			args: args{
				anns: map[string]string{
					"ingress.kubernetes.io/force-ssl-redirect": "true",
				},
			},
			want: Route{
				Route: kong.Route{
					HTTPSRedirectStatusCode: kong.Int(302),
					Protocols:               []*string{kong.String("https")},
				},
			},
		},
		{
			name: "force ssl annotation set to true and protocol is set to grpc",
			args: args{
				route: Route{
					Route: kong.Route{
						Protocols: []*string{kong.String("grpc")},
					},
				},
				anns: map[string]string{
					"ingress.kubernetes.io/force-ssl-redirect": "true",
					"konghq.com/protocols":                     "grpc",
				},
			},
			want: Route{
				Route: kong.Route{
					HTTPSRedirectStatusCode: kong.Int(302),
					Protocols:               []*string{kong.String("grpcs")},
				},
			},
		},
		{
			name: "force ssl annotation set to false",
			args: args{
				anns: map[string]string{
					"ingress.kubernetes.io/force-ssl-redirect": "false",
				},
			},
		},
		{
			name: "force ssl annotation set to true and HTTPS redirect code set to 307",
			args: args{
				anns: map[string]string{
					"ingress.kubernetes.io/force-ssl-redirect": "true",
					"konghq.com/https-redirect-status-code":    "307",
				},
			},
			want: Route{
				Route: kong.Route{
					HTTPSRedirectStatusCode: kong.Int(307),
					Protocols:               []*string{kong.String("https")},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overrideHTTPSRedirectCode(tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overrideRouteHTTPSRedirectCode() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}

func TestOverrideRoutePreserveHost(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Route
	}{
		{name: "basic empty route"},
		{
			name: "basic sanity",
			args: args{
				anns: map[string]string{
					"konghq.com/preserve-host": "true",
				},
			},
			want: Route{
				Route: kong.Route{
					PreserveHost: kong.Bool(true),
				},
			},
		},
		{
			name: "case insensitive",
			args: args{
				anns: map[string]string{
					"konghq.com/preserve-host": "faLSe",
				},
			},
			want: Route{
				Route: kong.Route{
					PreserveHost: kong.Bool(false),
				},
			},
		},
		{
			name: "random integer value",
			args: args{
				anns: map[string]string{
					"konghq.com/https-redirect-status-code": "42",
				},
			},
		},
		{
			name: "random string",
			args: args{
				anns: map[string]string{
					"konghq.com/https-redirect-status-code": "foo",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overridePreserveHost(tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overrideRoutePreserveHost() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}

func TestOverrideRouteRegexPriority(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Route
	}{
		{name: "basic empty route"},
		{
			name: "basic sanity",
			args: args{
				anns: map[string]string{
					"konghq.com/regex-priority": "10",
				},
			},
			want: Route{
				Route: kong.Route{
					RegexPriority: kong.Int(10),
				},
			},
		},
		{
			name: "negative integer",
			args: args{
				anns: map[string]string{
					"konghq.com/regex-priority": "-10",
				},
			},
			want: Route{
				Route: kong.Route{
					RegexPriority: kong.Int(-10),
				},
			},
		},
		{
			name: "random float value",
			args: args{
				anns: map[string]string{
					"konghq.com/regex-priority": "42.42",
				},
			},
		},
		{
			name: "random string",
			args: args{
				anns: map[string]string{
					"konghq.com/regex-priority": "foo",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overrideRegexPriority(tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overrideRouteRegexPriority() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}

func TestOverrideRouteMethods(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Route
	}{
		{name: "basic empty route"},
		{
			name: "basic sanity",
			args: args{
				anns: map[string]string{
					"konghq.com/methods": "POST,GET",
				},
			},
			want: Route{
				Route: kong.Route{
					Methods: kong.StringSlice("POST", "GET"),
				},
			},
		},
		{
			name: "non-string",
			args: args{
				anns: map[string]string{
					"konghq.com/methods": "-10,GET",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overrideMethods(zapr.NewLogger(zap.NewNop()), tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overrideRouteMethods() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}

func TestOverrideRouteSNIs(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Route
	}{
		{name: "basic empty route"},
		{
			name: "basic sanity, with strippable space",
			args: args{
				anns: map[string]string{
					"konghq.com/snis": "hrodna.kong.example, katowice.kong.example",
				},
			},
			want: Route{
				Route: kong.Route{
					SNIs: kong.StringSlice("hrodna.kong.example", "katowice.kong.example"),
				},
			},
		},
		{
			name: "not hostnames at all",
			args: args{
				anns: map[string]string{
					"konghq.com/snis": "-10,GET",
				},
			},
		},
		{
			name: "wildcard hostname, not valid for SNI",
			args: args{
				anns: map[string]string{
					"konghq.com/snis": "*.example.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overrideSNIs(zapr.NewLogger(zap.NewNop()), tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overrideRouteSNIs() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}

func TestOverrideRequestBuffering(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want kong.Route
	}{
		{},
		{
			name: "basic empty route",
			args: args{
				route: Route{Route: kong.Route{}},
			},
			want: kong.Route{},
		},
		{
			name: "set to false",
			args: args{
				route: Route{
					Route: kong.Route{},
				},
				anns: map[string]string{
					"konghq.com/request-buffering": "false",
				},
			},
			want: kong.Route{
				RequestBuffering: kong.Bool(false),
			},
		},
		{
			name: "set to true and case insensitive",
			args: args{
				route: Route{
					Route: kong.Route{},
				},
				anns: map[string]string{
					"konghq.com/request-buffering": "tRuE",
				},
			},
			want: kong.Route{
				RequestBuffering: kong.Bool(true),
			},
		},
		{
			name: "overrides any other value",
			args: args{
				route: Route{
					Route: kong.Route{
						RequestBuffering: kong.Bool(false),
					},
				},
				anns: map[string]string{
					"konghq.com/request-buffering": "tRuE",
				},
			},
			want: kong.Route{
				RequestBuffering: kong.Bool(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overrideRequestBuffering(zapr.NewLogger(zap.NewNop()), tt.args.anns)
			if !reflect.DeepEqual(tt.args.route.Route, tt.want) {
				t.Errorf("overrideRequestBuffering() got = %v, want %v", &tt.args.route.Route, tt.want)
			}
		})
	}
}

func TestOverrideResponseBuffering(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want kong.Route
	}{
		{},
		{
			name: "basic empty route",
			args: args{
				route: Route{Route: kong.Route{}},
			},
			want: kong.Route{},
		},
		{
			name: "set to false",
			args: args{
				route: Route{
					Route: kong.Route{},
				},
				anns: map[string]string{
					"konghq.com/response-buffering": "false",
				},
			},
			want: kong.Route{
				ResponseBuffering: kong.Bool(false),
			},
		},
		{
			name: "set to true and case insensitive",
			args: args{
				route: Route{
					Route: kong.Route{},
				},
				anns: map[string]string{
					"konghq.com/response-buffering": "tRuE",
				},
			},
			want: kong.Route{
				ResponseBuffering: kong.Bool(true),
			},
		},
		{
			name: "overrides any other value",
			args: args{
				route: Route{
					Route: kong.Route{
						ResponseBuffering: kong.Bool(false),
					},
				},
				anns: map[string]string{
					"konghq.com/response-buffering": "tRuE",
				},
			},
			want: kong.Route{
				ResponseBuffering: kong.Bool(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overrideResponseBuffering(zapr.NewLogger(zap.NewNop()), tt.args.anns)
			if !reflect.DeepEqual(tt.args.route.Route, tt.want) {
				t.Errorf("overrideResponseBuffering() got = %v, want %v", &tt.args.route.Route, tt.want)
			}
		})
	}
}

func TestOverrideHosts(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Route
	}{
		{name: "basic empty route"},
		{
			name: "basic sanity",
			args: args{
				anns: map[string]string{
					"konghq.com/host-aliases": "illustration.com, example.com, *.example.com, example.*, *.illustration.*",
				},
			},
			want: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("illustration.com", "example.com", "*.example.com", "example.*", "*.illustration.*"),
				},
			},
		},
		{
			name: "ignore duplicates",
			args: args{
				anns: map[string]string{
					"konghq.com/host-aliases": "example.com, example.com",
				},
			},
			want: Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("example.com"),
				},
			},
		},
		{
			name: "not hostnames",
			args: args{
				anns: map[string]string{
					"konghq.com/host-aliases": "-10,GET",
				},
			},
		},
		{
			name: "wildcard not allowed in the domain name",
			args: args{
				anns: map[string]string{
					"konghq.com/host-aliases": "kong.*.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overrideHosts(zapr.NewLogger(zap.NewNop()), tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overrideHosts() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}

func TestOverrideHeaders(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Route
	}{
		{
			name: "basic empty route",
		},
		{
			name: "single header single value",
			args: args{
				anns: map[string]string{
					"konghq.com/headers.x-example": "example",
				},
			},
			want: Route{
				Route: kong.Route{
					Headers: map[string][]string{"x-example": {"example"}},
				},
			},
		},
		{
			name: "single header multi value",
			args: args{
				anns: map[string]string{
					"konghq.com/headers.x-example": "foo,bar",
				},
			},
			want: Route{
				Route: kong.Route{
					Headers: map[string][]string{"x-example": {"foo", "bar"}},
				},
			},
		},
		{
			name: "multi header single value",
			args: args{
				anns: map[string]string{
					"konghq.com/headers.x-foo": "example",
					"konghq.com/headers.x-bar": "example",
				},
			},
			want: Route{
				Route: kong.Route{
					Headers: map[string][]string{
						"x-foo": {"example"},
						"x-bar": {"example"},
					},
				},
			},
		},
		{
			name: "multi header multi value",
			args: args{
				anns: map[string]string{
					"konghq.com/headers.x-foo": "foo,bar",
					"konghq.com/headers.x-bar": "bar,baz",
				},
			},
			want: Route{
				Route: kong.Route{
					Headers: map[string][]string{
						"x-foo": {"foo", "bar"},
						"x-bar": {"bar", "baz"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overrideHeaders(tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overrideHeaders() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}

func TestOverridePathHandling(t *testing.T) {
	type args struct {
		route Route
		anns  map[string]string
	}
	tests := []struct {
		name string
		args args
		want Route
	}{
		{name: "basic empty route"},
		{
			name: "expected value",
			args: args{
				anns: map[string]string{
					"konghq.com/path-handling": "v1",
				},
			},
			want: Route{
				Route: kong.Route{
					PathHandling: kong.String("v1"),
				},
			},
		},
		{
			name: "invalid value",
			args: args{
				anns: map[string]string{
					"konghq.com/path-handling": "vA",
				},
			},
			want: Route{
				Route: kong.Route{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.route.overridePathHandling(zapr.NewLogger(zap.NewNop()), tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overridePathHandling() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}
