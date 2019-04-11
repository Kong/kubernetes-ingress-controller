package dbless

import (
	"reflect"
	"testing"

	"github.com/hbagdi/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser"
)

func TestKongNativeState(t *testing.T) {
	type args struct {
		k8sState *parser.KongState
	}
	tests := []struct {
		name string
		args args
		want *KongDeclarativeConfig
	}{
		{
			"empty state",
			args{},
			&KongDeclarativeConfig{
				FormatVersion: "1.1",
			},
		},
		{
			"routes and services",
			args{
				&parser.KongState{
					Services: []parser.Service{
						{
							Service: kong.Service{
								Name:     kong.String("foo"),
								Host:     kong.String("example.com"),
								Protocol: kong.String("https"),
								Port:     kong.Int(443),
							},
							Routes: []parser.Route{
								{
									Route: kong.Route{
										Name:    kong.String("bar"),
										Methods: kong.StringSlice("GET", "POST"),
										Paths:   kong.StringSlice("/bar"),
									},
								},
							},
						},
						{
							Service: kong.Service{
								Name:     kong.String("bar"),
								Host:     kong.String("1.example.com"),
								Protocol: kong.String("http"),
								Port:     kong.Int(80),
							},
							Routes: []parser.Route{
								{
									Route: kong.Route{
										Name:    kong.String("baz"),
										Methods: kong.StringSlice("POST"),
										Paths:   kong.StringSlice("/baz"),
									},
								},
							},
						},
					},
				},
			},
			&KongDeclarativeConfig{
				FormatVersion: "1.1",
				Services: []Service{
					{
						Service: kong.Service{
							Name:     kong.String("foo"),
							Host:     kong.String("example.com"),
							Protocol: kong.String("https"),
							Port:     kong.Int(443),
						},
						Routes: []Route{
							{
								Route: kong.Route{
									Name:    kong.String("bar"),
									Methods: kong.StringSlice("GET", "POST"),
									Paths:   kong.StringSlice("/bar"),
								},
							},
						},
					},
					{
						Service: kong.Service{
							Name:     kong.String("bar"),
							Host:     kong.String("1.example.com"),
							Protocol: kong.String("http"),
							Port:     kong.Int(80),
						},
						Routes: []Route{
							{
								Route: kong.Route{
									Name:    kong.String("baz"),
									Methods: kong.StringSlice("POST"),
									Paths:   kong.StringSlice("/baz"),
								},
							},
						},
					},
				},
			},
		},
		{
			"service-level plugins",
			args{
				&parser.KongState{
					Services: []parser.Service{
						{
							Service: kong.Service{
								Name:     kong.String("foo"),
								Host:     kong.String("example.com"),
								Protocol: kong.String("https"),
								Port:     kong.Int(443),
							},
							Plugins: []kong.Plugin{
								{
									Name: kong.String("key-auth"),
									Config: kong.Configuration{
										"property1": "value1",
									},
								},
								{
									Name:    kong.String("rate-limiting"),
									Enabled: kong.Bool(false),
								},
							},
						},
					},
				},
			},
			&KongDeclarativeConfig{
				FormatVersion: "1.1",
				Services: []Service{
					{
						Service: kong.Service{
							Name:     kong.String("foo"),
							Host:     kong.String("example.com"),
							Protocol: kong.String("https"),
							Port:     kong.Int(443),
						},
						Plugins: []kong.Plugin{
							{
								Name: kong.String("key-auth"),
								Config: kong.Configuration{
									"property1": "value1",
								},
							},
							{
								Name:    kong.String("rate-limiting"),
								Enabled: kong.Bool(false),
							},
						},
					},
				},
			},
		},
		{
			"route-level plugins",
			args{
				&parser.KongState{
					Services: []parser.Service{
						{
							Service: kong.Service{
								Name:     kong.String("foo"),
								Host:     kong.String("example.com"),
								Protocol: kong.String("https"),
								Port:     kong.Int(443),
							},
							Routes: []parser.Route{
								{
									Route: kong.Route{
										Name:    kong.String("baz"),
										Methods: kong.StringSlice("POST"),
										Paths:   kong.StringSlice("/baz"),
									},
									Plugins: []kong.Plugin{
										{
											Name: kong.String("key-auth"),
											Config: kong.Configuration{
												"property1": "value1",
											},
										},
										{
											Name:    kong.String("rate-limiting"),
											Enabled: kong.Bool(false),
										},
									},
								},
								{
									Route: kong.Route{
										Name:    kong.String("bar"),
										Methods: kong.StringSlice("GET", "POST"),
										Paths:   kong.StringSlice("/bar"),
									},
									Plugins: []kong.Plugin{
										{
											Name: kong.String("basic-auth"),
											Config: kong.Configuration{
												"property1": "value1",
											},
										},
										{
											Name:    kong.String("rate-limiting"),
											Enabled: kong.Bool(false),
										},
									},
								},
							},
						},
					},
				},
			},
			&KongDeclarativeConfig{
				FormatVersion: "1.1",
				Services: []Service{
					{
						Service: kong.Service{
							Name:     kong.String("foo"),
							Host:     kong.String("example.com"),
							Protocol: kong.String("https"),
							Port:     kong.Int(443),
						},
						Routes: []Route{
							{
								Route: kong.Route{
									Name:    kong.String("baz"),
									Methods: kong.StringSlice("POST"),
									Paths:   kong.StringSlice("/baz"),
								},
								Plugins: []kong.Plugin{
									{
										Name: kong.String("key-auth"),
										Config: kong.Configuration{
											"property1": "value1",
										},
									},
									{
										Name:    kong.String("rate-limiting"),
										Enabled: kong.Bool(false),
									},
								},
							},
							{
								Route: kong.Route{
									Name:    kong.String("bar"),
									Methods: kong.StringSlice("GET", "POST"),
									Paths:   kong.StringSlice("/bar"),
								},
								Plugins: []kong.Plugin{
									{
										Name: kong.String("basic-auth"),
										Config: kong.Configuration{
											"property1": "value1",
										},
									},
									{
										Name:    kong.String("rate-limiting"),
										Enabled: kong.Bool(false),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"upstream and targets",
			args{
				&parser.KongState{
					Upstreams: []parser.Upstream{
						{
							Upstream: kong.Upstream{
								Name: kong.String("foo"),
							},
							Targets: []parser.Target{
								{
									Target: kong.Target{
										Target: kong.String("10.0.1.42"),
									},
								},
								{
									Target: kong.Target{
										Target: kong.String("10.0.1.43"),
									},
								},
								{
									Target: kong.Target{
										Target: kong.String("10.0.1.44"),
									},
								},
							},
						},
						{
							Upstream: kong.Upstream{
								Name: kong.String("bar"),
							},
							Targets: []parser.Target{
								{
									Target: kong.Target{
										Target: kong.String("10.0.2.42"),
									},
								},
								{
									Target: kong.Target{
										Target: kong.String("10.0.2.43"),
									},
								},
								{
									Target: kong.Target{
										Target: kong.String("10.0.2.44"),
									},
								},
							},
						},
					},
				},
			},
			&KongDeclarativeConfig{
				FormatVersion: "1.1",
				Upstreams: []Upstream{
					{
						Upstream: kong.Upstream{
							Name: kong.String("foo"),
						},
						Targets: []kong.Target{

							{
								Target: kong.String("10.0.1.42"),
							},
							{
								Target: kong.String("10.0.1.43"),
							},
							{
								Target: kong.String("10.0.1.44"),
							},
						},
					},
					{
						Upstream: kong.Upstream{
							Name: kong.String("bar"),
						},
						Targets: []kong.Target{
							{
								Target: kong.String("10.0.2.42"),
							},
							{
								Target: kong.String("10.0.2.43"),
							},
							{
								Target: kong.String("10.0.2.44"),
							},
						},
					},
				},
			},
		},
		{
			"global plugins",
			args{
				&parser.KongState{
					GlobalPlugins: []parser.Plugin{
						{
							Plugin: kong.Plugin{

								Name: kong.String("basic-auth"),
								Config: kong.Configuration{
									"property1": "value1",
								},
							},
						},
						{
							Plugin: kong.Plugin{

								Name:    kong.String("rate-limiting"),
								Enabled: kong.Bool(false),
							},
						},
					},
				},
			},
			&KongDeclarativeConfig{
				FormatVersion: "1.1",
				Plugins: []kong.Plugin{
					{
						Name: kong.String("basic-auth"),
						Config: kong.Configuration{
							"property1": "value1",
						},
					},
					{
						Name:    kong.String("rate-limiting"),
						Enabled: kong.Bool(false),
					},
				},
			},
		},
		{
			"certificate and SNIs",
			args{
				&parser.KongState{
					Certificates: []parser.Certificate{
						{
							Certificate: kong.Certificate{
								Cert: kong.String("foo"),
								Key:  kong.String("bar"),
								SNIs: kong.StringSlice("1.example.com",
									"2.example.com"),
							},
						},
						{
							Certificate: kong.Certificate{
								Cert: kong.String("fuz"),
								Key:  kong.String("baz"),
								SNIs: kong.StringSlice("3.example.com",
									"4.example.com"),
							},
						},
					},
				},
			},
			&KongDeclarativeConfig{
				FormatVersion: "1.1",
				Certificates: []Certificate{
					{
						Cert: kong.String("foo"),
						Key:  kong.String("bar"),
						SNIs: []kong.SNI{
							{
								Name: kong.String("1.example.com"),
							},
							{
								Name: kong.String("2.example.com"),
							},
						},
					},
					{
						Cert: kong.String("fuz"),
						Key:  kong.String("baz"),
						SNIs: []kong.SNI{
							{
								Name: kong.String("3.example.com"),
							},
							{
								Name: kong.String("4.example.com"),
							},
						},
					},
				},
			},
		},
		{
			"consumers",
			args{
				&parser.KongState{
					Consumers: []parser.Consumer{
						{
							Consumer: kong.Consumer{
								Username: kong.String("foo"),
							},
							Credentials: map[string][]map[string]interface{}{
								"key-auth": {
									{
										"apikey": "secret-api-key",
									},
								},
							},
						},
						{
							Consumer: kong.Consumer{
								Username: kong.String("bar"),
							},
						},
					},
				},
			},
			&KongDeclarativeConfig{
				FormatVersion: "1.1",
				Consumers: []Consumer{
					{
						Consumer: kong.Consumer{
							Username: kong.String("foo"),
						},
						Credentials: map[string][]map[string]interface{}{
							"key-auth": {
								{
									"apikey": "secret-api-key",
								},
							},
						},
					},
					{
						Consumer: kong.Consumer{
							Username: kong.String("bar"),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KongNativeState(tt.args.k8sState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KongNativeState() = %v, want %v", got, tt.want)
			}
		})
	}
}
