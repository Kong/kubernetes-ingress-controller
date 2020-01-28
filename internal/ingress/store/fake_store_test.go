package store

import (
	"github.com/pkg/errors"
	"testing"

	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
				obj: networking.Ingress{
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

func TestFakeStoreIngress(t *testing.T) {
	assert := assert.New(t)

	ingresses := []*networking.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
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
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
				Annotations: map[string]string{
					"kubernetes.io/ingress.class": "not-kong",
				},
			},
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Path: "/bar",
										Backend: networking.IngressBackend{
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
	store, err := NewFakeStore(FakeObjects{Ingresses: ingresses})
	assert.Nil(err)
	assert.NotNil(store)
	assert.Len(store.ListIngresses(), 1)
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
	plugins, err = store.ListGlobalKongPlugins()
	assert.Len(plugins, 0)

	plugins = []*configurationv1.KongPlugin{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Labels: map[string]string{
					"global": "true",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
				Labels: map[string]string{
					"global": "true",
				},
			},
		},
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
	assert.Len(plugins, 2)

	plugin, err := store.GetKongPlugin("default", "bar")
	assert.NotNil(plugin)
	assert.Nil(err)

	plugin, err = store.GetKongPlugin("default", "does-not-exist")
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
	assert.Len(plugins, 0)

	plugins = []*configurationv1.KongClusterPlugin{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
				Labels: map[string]string{
					"global": "true",
				},
			},
		},
		{
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
	assert.Len(plugins, 2)

	plugin, err := store.GetKongClusterPlugin("foo")
	assert.NotNil(plugin)
	assert.Nil(err)

	plugin, err = store.GetKongClusterPlugin("does-not-exist")
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
	assert.Nil(plugin)
}

func TestFakeStoreCredentials(t *testing.T) {
	assert := assert.New(t)

	credentials := []*configurationv1.KongCredential{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "baz",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{KongCredentials: credentials})
	assert.Nil(err)
	assert.NotNil(store)
	credentials = store.ListKongCredentials()
	assert.Len(credentials, 2)
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
