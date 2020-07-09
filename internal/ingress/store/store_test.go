package store

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Test_networkingIngressV1Beta1(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want *networking.Ingress
	}{
		{
			name: "networking.Ingress is returned as is",
			args: args{
				obj: &networking.Ingress{},
			},
			want: &networking.Ingress{},
		},
		{
			name: "returns nil if a non-ingress object is passed in",
			args: args{
				&core.Service{
					Spec: core.ServiceSpec{
						Type:      core.ServiceTypeClusterIP,
						ClusterIP: "1.1.1.1",
						Ports: []core.ServicePort{
							{
								Name:       "default",
								TargetPort: intstr.FromString("port-1"),
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "correctly transformers from extensions to networking group",
			args: args{
				obj: &extensions.Ingress{
					Spec: extensions.IngressSpec{
						Rules: []extensions.IngressRule{
							{
								Host: "example.com",
								IngressRuleValue: extensions.IngressRuleValue{
									HTTP: &extensions.HTTPIngressRuleValue{
										Paths: []extensions.HTTPIngressPath{
											{
												Path: "/",
												Backend: extensions.IngressBackend{
													ServiceName: "foo-svc",
													ServicePort: intstr.FromInt(80),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: &networking.Ingress{
				Spec: networking.IngressSpec{
					Rules: []networking.IngressRule{
						{
							Host: "example.com",
							IngressRuleValue: networking.IngressRuleValue{
								HTTP: &networking.HTTPIngressRuleValue{
									Paths: []networking.HTTPIngressPath{
										{
											Path: "/",
											Backend: networking.IngressBackend{
												ServiceName: "foo-svc",
												ServicePort: intstr.FromInt(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Store{
				logger: logrus.New(),
			}
			if got := s.networkingIngressV1Beta1(tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("networkingIngressV1Beta1() = %v, want %v", got, tt.want)
			}
		})
	}
}
