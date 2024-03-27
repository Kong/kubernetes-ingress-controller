package store

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func TestKeyFunc(t *testing.T) {
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
				obj: netv1.Ingress{
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

func TestFakeStoreIngressV1(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	defaultClass := annotations.DefaultIngressClass
	ingresses := []*netv1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path: "/",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: netv1.ServiceBackendPort{
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
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path: "/bar",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "bar-svc",
												Port: netv1.ServiceBackendPort{
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
			Spec: netv1.IngressSpec{
				Rules:            []netv1.IngressRule{},
				IngressClassName: &defaultClass,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "default",
			},
			Spec: netv1.IngressSpec{
				Rules:            []netv1.IngressRule{},
				IngressClassName: &defaultClass,
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{IngressesV1: ingresses})
	require.Nil(err)
	require.NotNil(store)
	assert.Len(store.ListIngressesV1(), 2)
}

func TestFakeStoreIngressClassV1(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	classes := []*netv1.IngressClass{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
			Spec: netv1.IngressClassSpec{
				Controller: IngressClassKongController,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bar",
			},
			Spec: netv1.IngressClassSpec{
				Controller: IngressClassKongController,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "baz",
			},
			Spec: netv1.IngressClassSpec{
				Controller: "some-other-controller.example.com/controller",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{IngressClassesV1: classes})
	require.Nil(err)
	require.NotNil(store)
	assert.Len(store.ListIngressClassesV1(), 2)
}

func TestFakeStoreListTCPIngress(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	ingresses := []*kongv1beta1.TCPIngress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: kongv1beta1.TCPIngressSpec{
				Rules: []kongv1beta1.IngressRule{
					{
						Port: 9000,
						Backend: kongv1beta1.IngressBackend{
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
			Spec: kongv1beta1.TCPIngressSpec{
				Rules: []kongv1beta1.IngressRule{
					{
						Port: 9000,
						Backend: kongv1beta1.IngressBackend{
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
			Spec: kongv1beta1.TCPIngressSpec{
				Rules: []kongv1beta1.IngressRule{
					{
						Port: 8000,
						Backend: kongv1beta1.IngressBackend{
							ServiceName: "bar-svc",
							ServicePort: 80,
						},
					},
				},
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{TCPIngresses: ingresses})
	require.Nil(err)
	require.NotNil(store)
	ings, err := store.ListTCPIngresses()
	assert.Nil(err)
	assert.Len(ings, 1)
}

func TestFakeStoreListKnativeIngress(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	ingresses := []*knative.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					"networking.knative.dev/ingress-class": annotations.DefaultIngressClass,
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
	require.Nil(err)
	require.NotNil(store)
	ings, err := store.ListKnativeIngresses()
	assert.Len(ings, 1)
	assert.Nil(err)
}

func TestFakeStoreService(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	services := []*corev1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{Services: services})
	require.Nil(err)
	require.NotNil(store)
	service, err := store.GetService("default", "foo")
	assert.NotNil(service)
	assert.Nil(err)

	service, err = store.GetService("default", "does-not-exists")
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
	assert.Nil(service)
}

func TestFakeStoreEndpointSlice(t *testing.T) {
	t.Parallel()
	endpoints := []*discoveryv1.EndpointSlice{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo-1",
				Namespace: "default",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "foo",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo-1",
				Namespace: "bar",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "foo",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo-2",
				Namespace: "bar",
				Labels: map[string]string{
					discoveryv1.LabelServiceName: "foo",
				},
			},
		},
	}

	store, err := NewFakeStore(FakeObjects{EndpointSlices: endpoints})
	require.Nil(t, err)
	require.NotNil(t, store)

	t.Run("Get EndpointSlices for Service with single EndpointSlice", func(t *testing.T) {
		c, err := store.GetEndpointSlicesForService("default", "foo")
		require.Nil(t, err)
		require.Len(t, c, 1)
	})

	t.Run("Get EndpointSlices for Service with multiple EndpointSlices", func(t *testing.T) {
		c, err := store.GetEndpointSlicesForService("bar", "foo")
		require.Nil(t, err)
		require.Len(t, c, 2)
	})

	t.Run("Get EndpointSlices for non-existing Service", func(t *testing.T) {
		c, err := store.GetEndpointSlicesForService("default", "does-not-exist")
		require.ErrorAs(t, err, &ErrNotFound{})
		require.Nil(t, c)
	})
}

func TestFakeStoreConsumer(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	consumers := []*kongv1.KongConsumer{
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
	require.Nil(err)
	require.NotNil(store)
	assert.Len(store.ListKongConsumers(), 1)
	c, err := store.GetKongConsumer("default", "foo")
	assert.Nil(err)
	assert.NotNil(c)

	c, err = store.GetKongConsumer("default", "does-not-exist")
	assert.Nil(c)
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
}

func TestFakeStoreConsumerGroup(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	consumerGroups := []*kongv1beta1.KongConsumerGroup{
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
	store, err := NewFakeStore(FakeObjects{KongConsumerGroups: consumerGroups})
	require.Nil(err)
	require.NotNil(store)
	assert.Len(store.ListKongConsumerGroups(), 1)
	c, err := store.GetKongConsumerGroup("default", "foo")
	assert.Nil(err)
	assert.NotNil(c)

	c, err = store.GetKongConsumerGroup("default", "does-not-exist")
	assert.Nil(c)
	assert.NotNil(err)
	assert.True(errors.As(err, &ErrNotFound{}))
}

func TestFakeStorePlugins(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	plugins := []*kongv1.KongPlugin{
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

	plugins = []*kongv1.KongPlugin{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "baz",
				Namespace: "default",
			},
		},
	}
	store, err = NewFakeStore(FakeObjects{KongPlugins: plugins})
	require.Nil(err)
	require.NotNil(store)
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
	require := require.New(t)

	plugins := []*kongv1.KongClusterPlugin{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{KongClusterPlugins: plugins})
	require.Nil(err)
	require.NotNil(store)
	plugins, err = store.ListGlobalKongClusterPlugins()
	assert.NoError(err)
	assert.Len(plugins, 0)

	plugins = []*kongv1.KongClusterPlugin{
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
	require.Nil(err)
	require.NotNil(store)
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
	require := require.New(t)

	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{Secrets: secrets})
	require.Nil(err)
	require.NotNil(store)
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
	require := require.New(t)

	kongIngresses := []*kongv1.KongIngress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{KongIngresses: kongIngresses})
	require.Nil(err)
	require.NotNil(store)
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
	require := require.New(t)

	secrets := []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
	}
	store, err := NewFakeStore(FakeObjects{Secrets: secrets})
	require.Nil(err)
	require.NotNil(store)
	certs, err := store.ListCACerts()
	assert.Nil(err)
	assert.Len(certs, 0)

	secrets = []*corev1.Secret{
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
	require.Nil(err)
	require.NotNil(store)
	certs, err = store.ListCACerts()
	assert.Nil(err)
	assert.Len(certs, 2, "expect two secrets as CA certificates")
}

func TestFakeStoreHTTPRoute(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	classes := []*gatewayv1.HTTPRoute{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
			Spec: gatewayv1.HTTPRouteSpec{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bar",
			},
			Spec: gatewayv1.HTTPRouteSpec{},
		},
	}
	store, err := NewFakeStore(FakeObjects{HTTPRoutes: classes})
	require.Nil(err)
	require.NotNil(store)
	routes, err := store.ListHTTPRoutes()
	assert.Nil(err)
	assert.Len(routes, 2, "expect two HTTPRoutes")
}

func TestFakeStoreUDPRoute(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	classes := []*gatewayv1alpha2.UDPRoute{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
			Spec: gatewayv1alpha2.UDPRouteSpec{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bar",
			},
			Spec: gatewayv1alpha2.UDPRouteSpec{},
		},
	}
	store, err := NewFakeStore(FakeObjects{UDPRoutes: classes})
	require.Nil(err)
	require.NotNil(store)
	routes, err := store.ListUDPRoutes()
	assert.Nil(err)
	assert.Len(routes, 2, "expect two UDPRoutes")
}

func TestFakeStoreTCPRoute(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	classes := []*gatewayv1alpha2.TCPRoute{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
			Spec: gatewayv1alpha2.TCPRouteSpec{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bar",
			},
			Spec: gatewayv1alpha2.TCPRouteSpec{},
		},
	}
	store, err := NewFakeStore(FakeObjects{TCPRoutes: classes})
	require.Nil(err)
	require.NotNil(store)
	routes, err := store.ListTCPRoutes()
	assert.Nil(err)
	assert.Len(routes, 2, "expect two TCPRoutes")
}

func TestFakeStoreTLSRoute(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	classes := []*gatewayv1alpha2.TLSRoute{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
			Spec: gatewayv1alpha2.TLSRouteSpec{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bar",
			},
			Spec: gatewayv1alpha2.TLSRouteSpec{},
		},
	}
	store, err := NewFakeStore(FakeObjects{TLSRoutes: classes})
	require.Nil(err)
	require.NotNil(store)
	routes, err := store.ListTLSRoutes()
	assert.Nil(err)
	assert.Len(routes, 2, "expect two TLSRoutes")
}

func TestFakeStoreReferenceGrant(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	grants := []*gatewayv1beta1.ReferenceGrant{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
			Spec: gatewayv1beta1.ReferenceGrantSpec{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bar",
			},
			Spec: gatewayv1beta1.ReferenceGrantSpec{},
		},
	}
	store, err := NewFakeStore(FakeObjects{ReferenceGrants: grants})
	require.Nil(err)
	require.NotNil(store)
	routes, err := store.ListReferenceGrants()
	assert.Nil(err)
	assert.Len(routes, 2, "expect two ReferenceGrants")
}

func TestFakeStoreGateway(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	grants := []*gatewayv1.Gateway{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
			Spec: gatewayv1.GatewaySpec{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bar",
			},
			Spec: gatewayv1.GatewaySpec{},
		},
	}
	store, err := NewFakeStore(FakeObjects{Gateways: grants})
	require.Nil(err)
	require.NotNil(store)
	routes, err := store.ListGateways()
	assert.Nil(err)
	assert.Len(routes, 2, "expect two Gateways")
}
