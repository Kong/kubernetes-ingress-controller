package store

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
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
		{
			name: "correctly transformers from extensions to networking group",
			args: args{
				obj: &extensions.Ingress{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "extensions/v1beta1",
						Kind:       "Ingress",
					},
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
					Status: extensions.IngressStatus{
						LoadBalancer: corev1.LoadBalancerStatus{
							Ingress: []corev1.LoadBalancerIngress{{IP: "1.2.3.4"}},
						},
					},
				},
			},
			want: &networking.Ingress{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "networking.k8s.io/v1beta1",
					Kind:       "Ingress",
				},
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
				Status: networking.IngressStatus{
					LoadBalancer: corev1.LoadBalancerStatus{
						Ingress: []corev1.LoadBalancerIngress{{IP: "1.2.3.4"}},
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

func TestCacheStoresGet(t *testing.T) {
	t.Log("configuring some yaml objects to store in the cache")
	svcYAML := []byte(`---
apiVersion: v1
kind: Service
metadata:
  name: httpbin-deployment
  namespace: default
  labels:
    app: httpbin
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: httpbin
  type: ClusterIP
`)
	ingYAML := []byte(`---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: httpbin-ingress
  namespace: default
  annotations:
    httpbin.ingress.kubernetes.io/rewrite-target: /
    kubernetes.io/ingress.class: "kong"
spec:
  rules:
  - http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: httpbin-deployment
            port:
              number: 80
`)

	t.Log("creating a new cache store from object yaml files")
	cs, err := NewCacheStoresFromObjYAML(svcYAML, ingYAML)
	require.NoError(t, err)

	t.Log("verifying that the cache store doesnt try to retrieve unsupported object types")
	_, exists, err := cs.Get(new(appsv1.Deployment))
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "Deployment is not a supported cache object type"))
	assert.False(t, exists)

	t.Log("verifying the integrity of the cache store")
	assert.Len(t, cs.IngressV1.List(), 1)
	assert.Len(t, cs.Service.List(), 1)
	assert.Len(t, cs.IngressV1beta1.List(), 0)
	assert.Len(t, cs.KongIngress.List(), 0)
	_, exists, err = cs.Get(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: "doesntexist", Name: "doesntexist"}})
	assert.NoError(t, err)
	assert.False(t, exists)

	t.Log("ensuring that we can Get() the objects back out of the cache store")
	svc := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "httpbin-deployment"}}
	ing := &netv1.Ingress{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "httpbin-ingress"}}
	_, exists, err = cs.Get(svc)
	assert.NoError(t, err)
	assert.True(t, exists)
	_, exists, err = cs.Get(ing)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func Test_getIngressClassHandling(t *testing.T) {
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
			if got := s.(Store).getIngressClassHandling(); got != tt.want {
				t.Errorf("s.getIngressClassHandling() = %v, want %v", got, tt.want)
			}
		})
	}
}
