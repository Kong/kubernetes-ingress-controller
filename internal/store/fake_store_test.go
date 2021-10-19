package store

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func Test_keyFunc(t *testing.T) {
	type args struct {
		obj interface{}
	}

	type F struct {
		Name      string
		Namespace string
	}
	type B struct {
		F
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			want: "Bar/Foo",
			args: args{
				obj: &F{
					Name:      "Foo",
					Namespace: "Bar",
				},
			},
		},
		{
			want: "Bar/Fu",
			args: args{
				obj: B{
					F: F{

						Name:      "Fu",
						Namespace: "Bar",
					},
				},
			},
		},
		{
			want: "default/foo",
			args: args{
				obj: networkingv1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "default",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := keyFunc(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("keyFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("keyFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeStoreEmpty(t *testing.T) {
	assert := assert.New(t)
	store, err := NewFakeStore(FakeObjects{})
	assert.Nil(err)
	assert.NotNil(store)
}

func TestFakeStoreIngressV1beta1(t *testing.T) {
	assert := assert.New(t)

	ingresses := []*networkingv1beta1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: networkingv1beta1.IngressSpec{
				Rules: []networkingv1beta1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networkingv1beta1.IngressRuleValue{
							HTTP: &networkingv1beta1.HTTPIngressRuleValue{
								Paths: []networkingv1beta1.HTTPIngressPath{
									{
										Path: "/",
										Backend: networkingv1beta1.IngressBackend{
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
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: "not-kong",
				},
			},
			Spec: networkingv1beta1.IngressSpec{
				Rules: []networkingv1beta1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networkingv1beta1.IngressRuleValue{
							HTTP: &networkingv1beta1.HTTPIngressRuleValue{
								Paths: []networkingv1beta1.HTTPIngressPath{
									{
										Path: "/bar",
										Backend: networkingv1beta1.IngressBackend{
											ServiceName: "bar-svc",
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
	}
	store, err := NewFakeStore(FakeObjects{IngressesV1beta1: ingresses})
	assert.Nil(err)
	assert.NotNil(store)
	assert.Len(store.ListIngressesV1beta1(), 1)
	assert.Len(store.ListIngressesV1(), 0)
}

func TestFakeStoreIngressV1(t *testing.T) {
	assert := assert.New(t)

	defaultClass := annotations.DefaultIngressClass
	ingresses := []*networkingv1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path: "/",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
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
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: "not-kong",
				},
			},
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path: "/bar",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "bar-svc",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
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
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: "skip-me-im-not-default",
				},
			},
			Spec: networkingv1.IngressSpec{
				Rules:            []networkingv1.IngressRule{},
				IngressClassName: &defaultClass,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
			},
			Spec: networkingv1.IngressSpec{
				Rules:            []networkingv1.IngressRule{},
				IngressClassName: &defaultClass,
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{IngressesV1: ingresses})
	assert.Nil(err)
	assert.NotNil(store)
	assert.Len(store.ListIngressesV1(), 2)
	assert.Len(store.ListIngressesV1beta1(), 0)
}

func TestFakeStoreListTCPIngress(t *testing.T) {
	assert := assert.New(t)

	ingresses := []*configurationv1beta1.TCPIngress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: configurationv1beta1.TCPIngressSpec{
				Rules: []configurationv1beta1.IngressRule{
					{
						Port: 9000,
						Backend: configurationv1beta1.IngressBackend{
							ServiceName: "foo-svc",
							ServicePort: 80,
						},
					},
				},
			},
		},
		{
			// this TCPIngress should *not* be loaded, as it lacks a class
			ObjectMeta: metav1.ObjectMeta{
				Name:      "baz",
				Namespace: "default",
			},
			Spec: configurationv1beta1.TCPIngressSpec{
				Rules: []configurationv1beta1.IngressRule{
					{
						Port: 9000,
						Backend: configurationv1beta1.IngressBackend{
							ServiceName: "foo-svc",
							ServicePort: 80,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: "not-kong",
				},
			},
			Spec: configurationv1beta1.TCPIngressSpec{
				Rules: []configurationv1beta1.IngressRule{
					{
						Port: 8000,
						Backend: configurationv1beta1.IngressBackend{
							ServiceName: "bar-svc",
							ServicePort: 80,
						},
					},
				},
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{TCPIngresses: ingresses})
	assert.Nil(err)
	assert.NotNil(store)
	ings, err := store.ListTCPIngresses()
	assert.Nil(err)
	assert.Len(ings, 1)
}

func TestFakeStoreListKnativeIngress(t *testing.T) {
	assert := assert.New(t)

	ingresses := []*knative.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					"networking.knative.dev/ingress.class": annotations.DefaultIngressClass,
				},
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{
						Hosts: []string{"example.com"},
						HTTP: &knative.HTTPIngressRuleValue{
							Paths: []knative.HTTPIngressPath{
								{
									Path: "/",
									Splits: []knative.IngressBackendSplit{
										{
											IngressBackend: knative.IngressBackend{
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
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "i-dont-get-processed-because-i-have-no-class-annotation",
				Namespace: "default",
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{
						Hosts: []string{"example.com"},
						HTTP: &knative.HTTPIngressRuleValue{
							Paths: []knative.HTTPIngressPath{
								{
									Path: "/",
									Splits: []knative.IngressBackendSplit{
										{
											IngressBackend: knative.IngressBackend{
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
	store, err := NewFakeStore(FakeObjects{KnativeIngresses: ingresses})
	assert.Nil(err)
	assert.NotNil(store)
	ings, err := store.ListKnativeIngresses()
	assert.Len(ings, 1)
	assert.Nil(err)
}

func TestFakeStoreService(t *testing.T) {
	assert := assert.New(t)

	services := []*apiv1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{Services: services})
	assert.Nil(err)
	assert.NotNil(store)
	service, err := store.GetService("default", "foo")
	assert.NotNil(service)
	assert.Nil(err)

	service, err = store.GetService("default", "does-not-exists")
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
	assert.Nil(service)
}

func TestFakeStoreEndpiont(t *testing.T) {
	assert := assert.New(t)

	endpoints := []*apiv1.Endpoints{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{Endpoints: endpoints})
	assert.Nil(err)
	assert.NotNil(store)
	c, err := store.GetEndpointsForService("default", "foo")
	assert.Nil(err)
	assert.NotNil(c)

	c, err = store.GetEndpointsForService("default", "does-not-exist")
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
	assert.Nil(c)
}

func TestFakeStoreConsumer(t *testing.T) {
	assert := assert.New(t)

	consumers := []*configurationv1.KongConsumer{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{KongConsumers: consumers})
	assert.Nil(err)
	assert.NotNil(store)
	assert.Len(store.ListKongConsumers(), 1)
	c, err := store.GetKongConsumer("default", "foo")
	assert.Nil(err)
	assert.NotNil(c)

	c, err = store.GetKongConsumer("default", "does-not-exist")
	assert.Nil(c)
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
}

func TestFakeStorePlugins(t *testing.T) {
	assert := assert.New(t)

	plugins := []*configurationv1.KongPlugin{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{KongPlugins: plugins})
	assert.Nil(err)
	assert.NotNil(store)

	plugins = []*configurationv1.KongPlugin{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "baz",
				Namespace: "default",
			},
		},
	}
	store, err = NewFakeStore(FakeObjects{KongPlugins: plugins})
	assert.Nil(err)
	assert.NotNil(store)
	plugins, err = store.ListGlobalKongPlugins()
	assert.NoError(err)
	assert.Len(plugins, 0)

	plugin, err := store.GetKongPlugin("default", "does-not-exist")
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
	assert.Nil(plugin)
}

func TestFakeStoreClusterPlugins(t *testing.T) {
	assert := assert.New(t)

	plugins := []*configurationv1.KongClusterPlugin{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{KongClusterPlugins: plugins})
	assert.Nil(err)
	assert.NotNil(store)
	plugins, err = store.ListGlobalKongClusterPlugins()
	assert.NoError(err)
	assert.Len(plugins, 0)

	plugins = []*configurationv1.KongClusterPlugin{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
				Labels: map[string]string{
					"global": "true",
				},
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
		},
		{
			// invalid due to lack of class, not loaded
			ObjectMeta: metav1.ObjectMeta{
				Name: "bar",
				Labels: map[string]string{
					"global": "true",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "baz",
			},
		},
	}
	store, err = NewFakeStore(FakeObjects{KongClusterPlugins: plugins})
	assert.Nil(err)
	assert.NotNil(store)
	plugins, err = store.ListGlobalKongClusterPlugins()
	assert.NoError(err)
	assert.Len(plugins, 1)

	plugin, err := store.GetKongClusterPlugin("foo")
	assert.NotNil(plugin)
	assert.Nil(err)

	plugin, err = store.GetKongClusterPlugin("does-not-exist")
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
	assert.Nil(plugin)
}

func TestFakeStoreSecret(t *testing.T) {
	assert := assert.New(t)

	secrets := []*apiv1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{Secrets: secrets})
	assert.Nil(err)
	assert.NotNil(store)
	secret, err := store.GetSecret("default", "foo")
	assert.Nil(err)
	assert.NotNil(secret)

	secret, err = store.GetSecret("default", "does-not-exist")
	assert.Nil(secret)
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
}

func TestFakeKongIngress(t *testing.T) {
	assert := assert.New(t)

	kongIngresses := []*configurationv1.KongIngress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{KongIngresses: kongIngresses})
	assert.Nil(err)
	assert.NotNil(store)
	kingress, err := store.GetKongIngress("default", "foo")
	assert.Nil(err)
	assert.NotNil(kingress)

	kingress, err = store.GetKongIngress("default", "does-not-exist")
	assert.NotNil(err)
	assert.Nil(kingress)
	assert.True(errors.As(err, &ErrNotFound{}))
}

func TestFakeStore_ListCACerts(t *testing.T) {
	assert := assert.New(t)

	secrets := []*apiv1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{Secrets: secrets})
	assert.Nil(err)
	assert.NotNil(store)
	certs, err := store.ListCACerts()
	assert.Nil(err)
	assert.Len(certs, 0)

	secrets = []*apiv1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Labels: map[string]string{
					"konghq.com/ca-cert": "true",
				},
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo1",
				Namespace: "default",
				Labels: map[string]string{
					"konghq.com/ca-cert": "true",
				},
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
		},
	}
	store, err = NewFakeStore(FakeObjects{Secrets: secrets})
	assert.Nil(err)
	assert.NotNil(store)
	certs, err = store.ListCACerts()
	assert.Nil(err)
	assert.Len(certs, 2, "expect two secrets as CA certificates")
}
