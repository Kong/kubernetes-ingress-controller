package kongstate

import (
	"reflect"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser/util"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestOverrideRoute(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inRoute        Route
		inKongIngresss configurationv1.KongIngress
		outRoute       Route
	}{
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			configurationv1.KongIngress{},
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
			configurationv1.KongIngress{
				Route: &kong.Route{
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
			configurationv1.KongIngress{
				Route: &kong.Route{
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
			configurationv1.KongIngress{
				Route: &kong.Route{
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
			configurationv1.KongIngress{
				Route: &kong.Route{
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
			configurationv1.KongIngress{
				Route: &kong.Route{
					Protocols:     kong.StringSlice("http"),
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
			configurationv1.KongIngress{
				Route: &kong.Route{
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
			configurationv1.KongIngress{
				Route: &kong.Route{
					Protocols: kong.StringSlice("grpc", "grpcs"),
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
			configurationv1.KongIngress{
				Route: &kong.Route{
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
	}

	for _, testcase := range testTable {
		testcase.inRoute.override(logrus.New(), &testcase.inKongIngresss)
		assert.Equal(testcase.inRoute, testcase.outRoute)
	}

	assert.NotPanics(func() {
		var nilRoute *Route
		nilRoute.override(logrus.New(), nil)
	})
}

func TestOverrideRoutePriority(t *testing.T) {
	assert := assert.New(t)
	var route Route
	route = Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
	}
	kongIngress := configurationv1.KongIngress{
		Route: &kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
	}

	ingMeta := util.K8sObjectInfo{
		Annotations: map[string]string{
			"konghq.com/protocols": "grpc,grpcs",
		},
	}

	route = Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
		Ingress: ingMeta,
	}
	route.override(logrus.New(), &kongIngress)
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
	kongIngress := configurationv1.KongIngress{
		Route: &kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
	}

	route.overrideByKongIngress(logrus.New(), &kongIngress)
	assert.Equal(route.Hosts, kong.StringSlice("foo.com", "bar.com"))
	assert.NotPanics(func() {
		var nilRoute *Route
		nilRoute.override(logrus.New(), nil)
	})
}
func TestOverrideRouteByAnnotation(t *testing.T) {
	assert := assert.New(t)
	var route Route
	route = Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
	}

	ingMeta := util.K8sObjectInfo{
		Annotations: map[string]string{
			"konghq.com/protocols": "grpc,grpcs",
		},
	}

	route = Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
		Ingress: ingMeta,
	}
	route.overrideByAnnotation(logrus.New())
	assert.Equal(route.Hosts, kong.StringSlice("foo.com", "bar.com"))
	assert.Equal(route.Protocols, kong.StringSlice("grpc", "grpcs"))

	assert.NotPanics(func() {
		var nilRoute *Route
		nilRoute.override(logrus.New(), nil)
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
		nilUpstream.override(nil, make(map[string]string))
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

func Test_overrideRouteStripPath(t *testing.T) {
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
					Route: kong.Route{}},
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

func Test_overrideRouteHTTPSRedirectCode(t *testing.T) {
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

func Test_overrideRoutePreserveHost(t *testing.T) {
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

func Test_overrideRouteRegexPriority(t *testing.T) {
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

func Test_overrideRouteMethods(t *testing.T) {
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
			tt.args.route.overrideMethods(logrus.New(), tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overrideRouteMethods() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}

func Test_overrideRouteSNIs(t *testing.T) {
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
			tt.args.route.overrideSNIs(logrus.New(), tt.args.anns)
			if !reflect.DeepEqual(tt.args.route, tt.want) {
				t.Errorf("overrideRouteSNIs() got = %v, want %v", tt.args.route, tt.want)
			}
		})
	}
}
