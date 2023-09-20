package kongstate

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

func TestOverrideRoute(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inRoute        Route
		inKongIngresss kongv1.KongIngress
		outRoute       Route
	}{
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			kongv1.KongIngress{},
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{
					Methods: kong.StringSlice("GET", "POST"),
				},
			},
			Route{
				Route: kong.Route{
					Hosts:   kong.StringSlice("foo.com", "bar.com"),
					Methods: kong.StringSlice("GET", "POST"),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{
					Methods: kong.StringSlice("GET   ", "post"),
				},
			},
			Route{
				Route: kong.Route{
					Hosts:   kong.StringSlice("foo.com", "bar.com"),
					Methods: kong.StringSlice("GET", "POST"),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{
					Methods: kong.StringSlice("GET", "-1"),
				},
			},
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{
					HTTPSRedirectStatusCode: kong.Int(302),
				},
			},
			Route{
				Route: kong.Route{
					Hosts:                   kong.StringSlice("foo.com", "bar.com"),
					HTTPSRedirectStatusCode: kong.Int(302),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts:        kong.StringSlice("foo.com", "bar.com"),
					PreserveHost: kong.Bool(true),
					StripPath:    kong.Bool(true),
				},
			},
			kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{
					Protocols:     kongv1.ProtocolSlice("http"),
					PreserveHost:  kong.Bool(false),
					StripPath:     kong.Bool(false),
					RegexPriority: kong.Int(10),
				},
			},
			Route{
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
			Route{
				Route: kong.Route{
					Hosts:     kong.StringSlice("foo.com"),
					Protocols: kong.StringSlice("http", "https"),
				},
			},
			kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{
					Headers: map[string][]string{
						"foo-header": {"bar-value"},
					},
				},
			},
			Route{
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
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com"),
				},
			},
			kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{
					Protocols: kongv1.ProtocolSlice("grpc", "grpcs"),
				},
			},
			Route{
				Route: kong.Route{
					Hosts:     kong.StringSlice("foo.com"),
					Protocols: kong.StringSlice("grpc", "grpcs"),
					StripPath: nil,
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com"),
				},
			},
			kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{
					PathHandling: kong.String("v1"),
				},
			},
			Route{
				Route: kong.Route{
					Hosts:        kong.StringSlice("foo.com"),
					PathHandling: kong.String("v1"),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			kongv1.KongIngress{
				Route: &kongv1.KongIngressRoute{
					RequestBuffering:  kong.Bool(true),
					ResponseBuffering: kong.Bool(true),
				},
			},
			Route{
				Route: kong.Route{
					Hosts:             kong.StringSlice("foo.com", "bar.com"),
					RequestBuffering:  kong.Bool(true),
					ResponseBuffering: kong.Bool(true),
				},
			},
		},
	}

	for _, testcase := range testTable {
		testcase := testcase
		testcase.inRoute.override(zapr.NewLogger(zap.NewNop()), &testcase.inKongIngresss)
		assert.Equal(testcase.inRoute, testcase.outRoute)
	}

	assert.NotPanics(func() {
		var nilRoute *Route
		nilRoute.override(zapr.NewLogger(zap.NewNop()), nil)
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
		indexStr := strconv.Itoa(i)
		tc := tc
		t.Run(indexStr+"-"+tc.name, func(t *testing.T) {
			tc.inRoute.override(zapr.NewLogger(zap.NewNop()), nil)
			assert.Equal(t, tc.outRoute, tc.inRoute, "should be the same as expected after overriding")
		})
	}
}

func TestOverrideRoutePriority(t *testing.T) {
	assert := assert.New(t)

	kongIngress := kongv1.KongIngress{
		Route: &kongv1.KongIngressRoute{
			Protocols: kongv1.ProtocolSlice("http"),
		},
	}

	ingMeta := util.K8sObjectInfo{
		Annotations: map[string]string{
			"konghq.com/protocols": "grpc,grpcs",
		},
	}

	route := Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
		Ingress: ingMeta,
	}
	route.override(zapr.NewLogger(zap.NewNop()), &kongIngress)
	assert.Equal(route.Hosts, kong.StringSlice("foo.com", "bar.com"))
	assert.Equal(route.Protocols, kong.StringSlice("grpc", "grpcs"))
}

func TestOverrideRouteByKongIngress(t *testing.T) {
	assert := assert.New(t)
	route := Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
	}
	kongIngress := kongv1.KongIngress{
		Route: &kongv1.KongIngressRoute{
			Protocols: kongv1.ProtocolSlice("http"),
		},
	}

	route.overrideByKongIngress(zapr.NewLogger(zap.NewNop()), &kongIngress)
	assert.Equal(route.Protocols, kong.StringSlice("http"))
	assert.NotPanics(func() {
		var nilRoute *Route
		nilRoute.override(zapr.NewLogger(zap.NewNop()), nil)
	})
}

func TestOverrideRouteByAnnotation(t *testing.T) {
	assert := assert.New(t)

	ingMeta := util.K8sObjectInfo{
		Annotations: map[string]string{
			"konghq.com/protocols": "grpc,grpcs",
		},
	}

	route := Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
		Ingress: ingMeta,
	}
	route.overrideByAnnotation(zapr.NewLogger(zap.NewNop()))
	assert.Equal(route.Hosts, kong.StringSlice("foo.com", "bar.com"))
	assert.Equal(route.Protocols, kong.StringSlice("grpc", "grpcs"))

	assert.NotPanics(func() {
		var nilRoute *Route
		nilRoute.override(zapr.NewLogger(zap.NewNop()), nil)
	})
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
		tt := tt
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
		tt := tt
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
		tt := tt
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
