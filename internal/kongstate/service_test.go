package kongstate

import (
	"reflect"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

func TestOverrideService(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inService      Service
		inKongIngresss configurationv1.KongIngress
		outService     Service
		inAnnotation   map[string]string
	}{
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Protocol: kong.String("https"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Retries: kong.Int(0),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
					Retries:  kong.Int(0),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Path: kong.String("/new-path"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/new-path"),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Retries: kong.Int(1),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
					Retries:  kong.Int(1),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					ConnectTimeout: kong.Int(100),
					ReadTimeout:    kong.Int(100),
					WriteTimeout:   kong.Int(100),
				},
			},
			Service{
				Service: kong.Service{
					Host:           kong.String("foo.com"),
					Port:           kong.Int(80),
					Name:           kong.String("foo"),
					Protocol:       kong.String("http"),
					Path:           kong.String("/"),
					ConnectTimeout: kong.Int(100),
					ReadTimeout:    kong.Int(100),
					WriteTimeout:   kong.Int(100),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpc"),
					Path:     nil,
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Protocol: kong.String("grpc"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpc"),
					Path:     nil,
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     nil,
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Protocol: kong.String("grpcs"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpcs"),
					Path:     nil,
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Protocol: kong.String("grpcs"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpcs"),
					Path:     nil,
				},
			},
			map[string]string{"konghq.com/protocol": "grpcs"},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Protocol: kong.String("grpcs"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpc"),
					Path:     nil,
				},
			},
			map[string]string{"konghq.com/protocol": "grpc"},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpcs"),
					Path:     nil,
				},
			},
			map[string]string{"konghq.com/protocol": "grpcs"},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Protocol: kong.String("grpcs"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			map[string]string{"konghq.com/protocol": "https"},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			map[string]string{"konghq.com/protocol": "https"},
		},
	}

	for _, testcase := range testTable {
		testcase.inService.override(&testcase.inKongIngresss, testcase.inAnnotation)
		assert.Equal(testcase.inService, testcase.outService)
	}

	assert.NotPanics(func() {
		var nilService *Service
		nilService.override(nil, nil)
	})
}

func Test_overrideServicePath(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{},
		{name: "basic empty service"},
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/path": "/foo",
				},
			},
			want: Service{
				Service: kong.Service{
					Path: kong.String("/foo"),
				},
			},
		},
		{
			name: "does not set path if doesn't start with /",
			args: args{
				anns: map[string]string{
					"konghq.com/path": "foo",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						Path: kong.String("/foo"),
					},
				},
				anns: map[string]string{
					"konghq.com/path": "/bar",
				},
			},
			want: Service{
				Service: kong.Service{
					Path: kong.String("/bar"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overridePath(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideServicePath() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}
