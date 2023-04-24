package store

import (
	"context"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
)

func TestNetworkingIngressV1Beta1(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want *netv1beta1.Ingress
	}{
		{
			name: "networking.Ingress is returned as is",
			args: args{
				obj: &netv1beta1.Ingress{},
			},
			want: &netv1beta1.Ingress{},
		},
		{
			name: "returns nil if a non-ingress object is passed in",
			args: args{
				&corev1.Service{
					Spec: corev1.ServiceSpec{
						Type:      corev1.ServiceTypeClusterIP,
						ClusterIP: "1.1.1.1",
						Ports: []corev1.ServicePort{
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

func TestGetIngressClassHandling(t *testing.T) {
	tests := []struct {
		name string
		objs FakeObjects
		want annotations.ClassMatching
	}{
		{
			name: "does not exist",
			objs: FakeObjects{},
			want: annotations.ExactClassMatch,
		},
		{
			name: "not default",
			objs: FakeObjects{
				IngressClassesV1: []*netv1.IngressClass{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: annotations.DefaultIngressClass,
						},
						Spec: netv1.IngressClassSpec{
							Controller: IngressClassKongController,
						},
					},
				},
			},
			want: annotations.ExactClassMatch,
		},
		{
			name: "default",
			objs: FakeObjects{
				IngressClassesV1: []*netv1.IngressClass{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: annotations.DefaultIngressClass,
							Annotations: map[string]string{
								"ingressclass.kubernetes.io/is-default-class": "true",
							},
						},
						Spec: netv1.IngressClassSpec{
							Controller: IngressClassKongController,
						},
					},
				},
			},
			want: annotations.ExactOrEmptyClassMatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewFakeStore(tt.objs)
			require.NoError(t, err)
			if got := s.(Store).getIngressClassHandling(context.TODO()); got != tt.want {
				t.Errorf("s.getIngressClassHandling() = %v, want %v", got, tt.want)
			}
		})
	}
}
