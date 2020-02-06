package parser

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hbagdi/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/store"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type TLSPair struct {
	Key, Cert string
}

var (
	tlsPairs = []TLSPair{
		{
			Cert: `-----BEGIN CERTIFICATE-----
MIIC2DCCAcACCQC32eFOsWpKojANBgkqhkiG9w0BAQsFADAuMRcwFQYDVQQDDA5z
ZWN1cmUtZm9vLWJhcjETMBEGA1UECgwKa29uZ2hxLm9yZzAeFw0xODEyMTgyMTI4
MDBaFw0xOTEyMTgyMTI4MDBaMC4xFzAVBgNVBAMMDnNlY3VyZS1mb28tYmFyMRMw
EQYDVQQKDAprb25naHEub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAqhl/HSwV6PbMv+cMFU9X+HuM7QbNNPh39GKa4pkxzFgiAnuuJ4jw9V/bzsEy
S+ZIyjzo+QKB1LzmgdcX4vkdI22BjxUd9HPHdZxtv3XilbNmSk9UOl2Hh1fORJoS
7YH+VbvVwiz5lo7qKRepbg/jcKkbs6AUE0YWFygtDLTvhP2qkphQkxZ0m8qroW91
CWgI73Ar6U2W/YQBRI3+LwtsKo0p2ASDijvqxElQBgBIiyGIr0RZc5pkCJ1eQdDB
2F6XaMfpeEyBj0MxypNL4S9HHfchOt55J1KOzYnUPkQnSoxp6oEjef4Q/ZCj5BRL
EGZnTb3tbwzHZCxGtgl9KqO9pQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAKQ5BX
kkBL+alERL31hsOgWgRiUMw+sPDtRS96ozUlPtVvAg9XFdpY4ldtWkxFcmBnhKzp
UewjrHkf9rR16NISwUTjlGIwaJu/ACQrY15v+r301Crq2DV+GjiUJFVuT495dp/l
0LZbt2Sh/uD+r3UNTcJpJ7jb1V0UP7FWXFj8oafsoFSgmxAPjpKQySTC54JK4AYb
QSnWu1nQLyohnrB9qLZhe2+jOQZnkKuCcWJQ5njvU6SxT3SOKE5XaOZCezEQ6IVL
U47YCCXsq+7wKWXBhKl4H2Ztk6x3HOC56l0noXWezsMfrou/kjwGuuViGnrjqelS
WQ7uVeNCUBY+l+qY
-----END CERTIFICATE-----`,
			Key: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCqGX8dLBXo9sy/
5wwVT1f4e4ztBs00+Hf0YprimTHMWCICe64niPD1X9vOwTJL5kjKPOj5AoHUvOaB
1xfi+R0jbYGPFR30c8d1nG2/deKVs2ZKT1Q6XYeHV85EmhLtgf5Vu9XCLPmWjuop
F6luD+NwqRuzoBQTRhYXKC0MtO+E/aqSmFCTFnSbyquhb3UJaAjvcCvpTZb9hAFE
jf4vC2wqjSnYBIOKO+rESVAGAEiLIYivRFlzmmQInV5B0MHYXpdox+l4TIGPQzHK
k0vhL0cd9yE63nknUo7NidQ+RCdKjGnqgSN5/hD9kKPkFEsQZmdNve1vDMdkLEa2
CX0qo72lAgMBAAECggEADxMTYNJ3Xp4Ap0EioQDXGv5YDul7ZiZe+xmCAHLzJtjo
qq+rT3WjZRuJr1kPzAosiT+8pdTDDMdw5jDZvRO2sV0TDksgzHk2RAYI897OpdWw
SwWcwU9oo2X0sb+1zbang5GR8BNsSxt/RQUDzu05itJx0gltvgeIDaVR2L5wO6ja
USa8OVuj/92XtIIve9OtyK9jAzgR6LQOTFrCCEv89/vmy5Bykv4Uz8s8swZmTs3v
XJmAmruHGuSLMfXk8lBRp/gVyNTi3uMsdph5AJbVKnra5TZLguEozZKbLdNUYk0p
+aAc7rxDcH2sPqa/7DwRvei9dvd5oB3VJlxGVgC8AQKBgQDfznRSSKAD15hoSDzt
cKNyhLgWAL+MD0jhHKUy3x+Z9OCvf0DVnmru5HfQKq5UfT0t8VTRPGKmOtAMD4cf
LYjIurvMvpVzQGSJfhtHQuULZTh3dfsM7xivMqSV+9txklMAakM7vGQlOQxhrScM
21Mp5LWDU6+e2pFCrQPop0IPkQKBgQDCkVE+dou2yFuJx3uytCH1yKPSy9tkdhQH
dGF12B5dq8MZZozAz5P9YN/COa9WjsNKDqWbEgLEksEQUq4t8SBjHnSV/D3x7rEF
qgwii0GETYxax6gms8nueIqWZQf+0NbX7Gc5mTqeVb7v3TrhsKr0VNMFRXXQwP2E
M/pxJq8q1QKBgQC3rH7oXLP+Ez0AMHDYSL3LKULOw/RvpMeh/9lQA6+ysTaIsP3r
kuSdhCEUVULXEiVYhBug0FcBp3jAvSmem8cLPb0Mjkim2mzoLfeDJ1JEZODPoaLU
fZEbj4tlj9oLvhOiXpMo/jaOGeCgdPN8aK86zXlt+wtBao0WVFnF4SalEQKBgQC1
uLfi2SGgs/0a8B/ORoO5ZY3s4c2lRMtsMvyb7iBeaIAuByPLKZUVABe89deXxnsL
fiaacPX41wBO2IoqCp2vNdC6DP9mKQNZQPtYgCvPAAbo+rVIgH9HpXn7AZ24FyGy
RfAbUcv3+in9KelGxZTF4zu8HqXtNXMSuOFeMT1FiQKBgF0R+IFDGHhD4nudAQvo
hncXsgyzK6QUzak6HmFji/CMZ6EU9q6A67JkiEWrYoKqIAKZ2Og8+Eucr/rDdGWc
kqlmLPBJAJeUsP/9KidBjTE5mIbn/2n089VPMBvnlt2xIcuB6+zrf2NjvlcZEyKS
Gn+T2uCyOP4a1DTUoPyoNJXo
-----END PRIVATE KEY-----`,
		},
	}
)

func TestGlobalPlugin(t *testing.T) {
	assert := assert.New(t)
	t.Run("global plugins are processed correctly", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			KongPlugins: []*configurationv1.KongPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-plugin",
						Namespace: "default",
						Labels: map[string]string{
							"global": "true",
						},
					},
					PluginName: "key-auth",
					Protocols:  []string{"grpc"},
					Config: configurationv1.Configuration{
						"foo": "bar",
					},
				},
			},
			KongClusterPlugins: []*configurationv1.KongClusterPlugin{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-plugin",
						Labels: map[string]string{
							"global": "true",
						},
					},
					Protocols:  []string{"http"},
					PluginName: "basic-auth",
					Config: configurationv1.Configuration{
						"foo1": "bar1",
					},
				},
			},
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(2, len(state.Plugins),
			"expected one plugin to be rendered")

		sort.SliceStable(state.Plugins, func(i, j int) bool {
			return strings.Compare(*state.Plugins[i].Name, *state.Plugins[j].Name) > 0
		})
		assert.Equal("key-auth", *state.Plugins[0].Name)
		assert.Equal(1, len(state.Plugins[0].Protocols))
		assert.Equal(kong.Configuration{"foo": "bar"}, state.Plugins[0].Config)

		assert.Equal("basic-auth", *state.Plugins[1].Name)
		assert.Equal(kong.Configuration{"foo1": "bar1"}, state.Plugins[1].Config)
	})
}

func TestServiceClientCertificate(t *testing.T) {
	assert := assert.New(t)
	t.Run("valid client-cert annotation", func(t *testing.T) {
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
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					UID:       types.UID("7428fb98-180b-4702-a91f-61351a33c6e4"),
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
		}
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"configuration.konghq.com/client-cert": "secret1",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses: ingresses,
			Secrets:   secrets,
			Services:  services,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Certificates),
			"expected one certificates to be rendered")
		assert.Equal("7428fb98-180b-4702-a91f-61351a33c6e4",
			*state.Certificates[0].ID)

		assert.Equal(1, len(state.Services))
		assert.Equal("7428fb98-180b-4702-a91f-61351a33c6e4",
			*state.Services[0].ClientCertificate.ID)
	})
	t.Run("client-cert secret doesn't exist", func(t *testing.T) {
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
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"configuration.konghq.com/client-cert": "secret1",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses: ingresses,
			Services:  services,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(0, len(state.Certificates),
			"expected no certificates to be rendered")

		assert.Equal(1, len(state.Services))
		assert.Nil(state.Services[0].ClientCertificate)
	})
}
func TestDefaultBackend(t *testing.T) {
	assert := assert.New(t)
	t.Run("default backend is processed correctly", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ing-with-default-backend",
					Namespace: "default",
				},
				Spec: networking.IngressSpec{
					Backend: &networking.IngressBackend{
						ServiceName: "default-svc",
						ServicePort: intstr.FromInt(80),
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "default-svc",
					Namespace: "default",
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses: ingresses,
			Services:  services,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Services),
			"expected one service to be rendered")
		assert.Equal("default.default-svc.80", *state.Services[0].Name)
		assert.Equal("default-svc.default.80.svc", *state.Services[0].Host)
		assert.Equal(1, len(state.Services[0].Routes),
			"expected one routes to be rendered")
		assert.Equal("default.ing-with-default-backend", *state.Services[0].Routes[0].Name)
		assert.Equal("/", *state.Services[0].Routes[0].Paths[0])
	})

	t.Run("client-cert secret doesn't exist", func(t *testing.T) {
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
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-svc",
					Namespace: "default",
					Annotations: map[string]string{
						"configuration.konghq.com/client-cert": "secret1",
					},
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses: ingresses,
			Services:  services,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(0, len(state.Certificates),
			"expected no certificates to be rendered")

		assert.Equal(1, len(state.Services))
		assert.Nil(state.Services[0].ClientCertificate)
	})
}

func TestParserSecret(t *testing.T) {
	assert := assert.New(t)
	t.Run("invalid TLS secret", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "default",
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"bar.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(""),
					"tls.key": []byte(""),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses: ingresses,
			Secrets:   secrets,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(0, len(state.Certificates),
			"expected no certificates to be rendered with empty secret")
	})
	t.Run("duplicate certificates", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "ns1",
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret2",
							Hosts:      []string{"bar.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret2",
					Namespace: "ns1",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses: ingresses,
			Secrets:   secrets,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Certificates),
			"certificates are de-duplicated")
	})
	t.Run("duplicate SNIs", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret1",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bar",
					Namespace: "ns1",
				},
				Spec: networking.IngressSpec{
					TLS: []networking.IngressTLS{
						{
							SecretName: "secret2",
							Hosts:      []string{"foo.com"},
						},
					},
				},
			},
		}

		secrets := []*corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret2",
					Namespace: "ns1",
				},
				Data: map[string][]byte{
					"tls.crt": []byte(tlsPairs[0].Cert),
					"tls.key": []byte(tlsPairs[0].Key),
				},
			},
		}
		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses: ingresses,
			Secrets:   secrets,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Certificates),
			"SNIs are de-duplicated")
	})
}

func TestPluginAnnotations(t *testing.T) {
	assert := assert.New(t)
	t.Run("simple association", func(t *testing.T) {
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "foo-svc",
					Namespace:   "default",
					Annotations: map[string]string{},
				},
			},
		}
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						"plugins.konghq.com": "foo-plugin",
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
		}
		plugins := []*configurationv1.KongPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "key-auth",
				Protocols:  []string{"grpc"},
				Config: configurationv1.Configuration{
					"foo": "bar",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses:   ingresses,
			Services:    services,
			KongPlugins: plugins,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
		assert.Equal("key-auth", *state.Plugins[0].Name)
		assert.Equal("grpc", *state.Plugins[0].Protocols[0])
	})
	t.Run("KongPlugin takes precedence over KongPlugin", func(t *testing.T) {
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "foo-svc",
					Namespace:   "default",
					Annotations: map[string]string{},
				},
			},
		}
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						"plugins.konghq.com": "foo-plugin",
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
		}
		clusterPlugins := []*configurationv1.KongClusterPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "basic-auth",
				Protocols:  []string{"grpc"},
				Config: configurationv1.Configuration{
					"foo": "bar",
				},
			},
		}
		plugins := []*configurationv1.KongPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "key-auth",
				Protocols:  []string{"grpc"},
				Config: configurationv1.Configuration{
					"foo": "bar",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses:          ingresses,
			Services:           services,
			KongPlugins:        plugins,
			KongClusterPlugins: clusterPlugins,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
		assert.Equal("key-auth", *state.Plugins[0].Name)
		assert.Equal("grpc", *state.Plugins[0].Protocols[0])
	})
	t.Run("KongClusterPlugin association", func(t *testing.T) {
		services := []*corev1.Service{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "foo-svc",
					Namespace:   "default",
					Annotations: map[string]string{},
				},
			},
		}
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						"plugins.konghq.com": "foo-plugin",
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
		}
		clusterPlugins := []*configurationv1.KongClusterPlugin{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo-plugin",
					Namespace: "default",
				},
				PluginName: "basic-auth",
				Protocols:  []string{"grpc"},
				Config: configurationv1.Configuration{
					"foo": "bar",
				},
			},
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses:          ingresses,
			Services:           services,
			KongClusterPlugins: clusterPlugins,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(1, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
		assert.Equal("basic-auth", *state.Plugins[0].Name)
		assert.Equal("grpc", *state.Plugins[0].Protocols[0])
	})
	t.Run("missing plugin", func(t *testing.T) {
		ingresses := []*networking.Ingress{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						"plugins.konghq.com": "does-not-exist",
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
		}

		store, err := store.NewFakeStore(store.FakeObjects{
			Ingresses: ingresses,
		})
		assert.Nil(err)
		parser := New(store)
		state, err := parser.Build()
		assert.Nil(err)
		assert.NotNil(state)
		assert.Equal(0, len(state.Plugins),
			"expected no plugins to be rendered with missing plugin")
	})
}

func TestParseIngressRules(t *testing.T) {
	assert := assert.New(t)
	p := Parser{}
	ingressList := []*networking.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
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
		// 1
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ing-with-tls",
				Namespace: "bar-namespace",
			},
			Spec: networking.IngressSpec{
				TLS: []networking.IngressTLS{
					{
						Hosts: []string{
							"1.example.com",
							"2.example.com",
						},
						SecretName: "sooper-secret",
					},
					{
						Hosts: []string{
							"3.example.com",
							"4.example.com",
						},
						SecretName: "sooper-secret2",
					},
				},
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
		// 2
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ing-with-default-backend",
				Namespace: "bar-namespace",
			},
			Spec: networking.IngressSpec{
				Backend: &networking.IngressBackend{
					ServiceName: "default-svc",
					ServicePort: intstr.FromInt(80),
				},
			},
		},
		// 3
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
			},
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Path: "/.well-known/acme-challenge/yolo",
										Backend: networking.IngressBackend{
											ServiceName: "cert-manager-solver-pod",
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
		// 4
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
			},
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
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
		// 5
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "baz",
				Namespace: "foo-namespace",
			},
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host:             "example.com",
						IngressRuleValue: networking.IngressRuleValue{},
					},
				},
			},
		},
		// 6
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
			},
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Backend: networking.IngressBackend{
											ServiceName: "foo-svc",
											ServicePort: intstr.FromInt(80),
										},
									},
								},
							},
						},
					},
					{
						Host: "example.net",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Backend: networking.IngressBackend{
											ServiceName: "foo-svc",
											ServicePort: intstr.FromInt(8000),
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
	tcpIngressList := []*configurationv1beta1.TCPIngress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
		// 1
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
			Spec: configurationv1beta1.IngressSpec{
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
		// 2
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
			Spec: configurationv1beta1.IngressSpec{
				Rules: []configurationv1beta1.IngressRule{
					{
						Host: "example.com",
						Port: 9000,
						Backend: configurationv1beta1.IngressBackend{
							ServiceName: "foo-svc",
							ServicePort: 80,
						},
					},
				},
			},
		},
		// 3
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
			Spec: configurationv1beta1.IngressSpec{
				TLS: []configurationv1beta1.IngressTLS{
					{
						Hosts: []string{
							"1.example.com",
							"2.example.com",
						},
						SecretName: "sooper-secret",
					},
					{
						Hosts: []string{
							"3.example.com",
							"4.example.com",
						},
						SecretName: "sooper-secret2",
					},
				},
			},
		},
	}
	t.Run("no ingress returns empty info", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{},
			[]*configurationv1beta1.TCPIngress{})
		assert.Equal(&parsedIngressRules{
			ServiceNameToServices: make(map[string]Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
		assert.Nil(err)
	})
	t.Run("empty TCPIngress return empty info", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{},
			[]*configurationv1beta1.TCPIngress{tcpIngressList[0]})
		assert.Equal(&parsedIngressRules{
			ServiceNameToServices: make(map[string]Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
		assert.Nil(err)
	})
	t.Run("simple ingress rule is parsed", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{
			ingressList[0],
		}, []*configurationv1beta1.TCPIngress{})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])
		assert.Nil(err)
	})
	t.Run("simple TCPIngress rule is parsed", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{},
			[]*configurationv1beta1.TCPIngress{tcpIngressList[1]})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		svc := parsedInfo.ServiceNameToServices["default.foo-svc.80"]
		assert.Equal("foo-svc.default.80.svc", *svc.Host)
		assert.Equal(80, *svc.Port)
		assert.Equal("tcp", *svc.Protocol)

		assert.Equal(1, len(svc.Routes))
		route := svc.Routes[0]
		assert.Equal(kong.Route{
			Name:      kong.String("default.foo.0"),
			Protocols: kong.StringSlice("tcp", "tls"),
			Destinations: []*kong.CIDRPort{
				{
					Port: kong.Int(9000),
				},
			},
		}, route.Route)
		assert.Nil(err)
	})
	t.Run("TCPIngress rule with host is parsed", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{},
			[]*configurationv1beta1.TCPIngress{tcpIngressList[2]})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		svc := parsedInfo.ServiceNameToServices["default.foo-svc.80"]
		assert.Equal("foo-svc.default.80.svc", *svc.Host)
		assert.Equal(80, *svc.Port)
		assert.Equal("tcp", *svc.Protocol)

		assert.Equal(1, len(svc.Routes))
		route := svc.Routes[0]
		assert.Equal(kong.Route{
			Name:      kong.String("default.foo.0"),
			Protocols: kong.StringSlice("tcp", "tls"),
			SNIs:      kong.StringSlice("example.com"),
			Destinations: []*kong.CIDRPort{
				{
					Port: kong.Int(9000),
				},
			},
		}, route.Route)
		assert.Nil(err)
	})
	t.Run("ingress rule with default backend", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{
			ingressList[0],
			ingressList[2],
		}, []*configurationv1beta1.TCPIngress{})
		assert.Equal(2, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])

		assert.Equal(1, len(parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes))
		assert.Equal("/", *parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes[0].Paths[0])
		assert.Equal(0, len(parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes[0].Hosts))
		assert.Nil(err)
	})
	t.Run("ingress rule with TLS", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{
			ingressList[1],
		}, []*configurationv1beta1.TCPIngress{})
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret2"]))

		assert.Nil(err)
	})
	t.Run("TCPIngress with TLS", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{},
			[]*configurationv1beta1.TCPIngress{tcpIngressList[3]})
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["default/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["default/sooper-secret2"]))

		assert.Nil(err)
	})
	t.Run("ingress rule with ACME like path has strip_path set to false", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{
			ingressList[3],
		}, []*configurationv1beta1.TCPIngress{})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("cert-manager-solver-pod.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Port)

		assert.Equal("/.well-known/acme-challenge/yolo", *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Routes[0].Hosts[0])
		assert.False(*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Routes[0].StripPath)

		assert.Nil(err)
	})
	t.Run("ingress with empty path is correctly parsed", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{
			ingressList[4],
		}, []*configurationv1beta1.TCPIngress{})
		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])

		assert.Nil(err)
	})
	t.Run("empty Ingress rule doesn't cause a panic", func(t *testing.T) {
		assert.NotPanics(func() {
			_, err := p.parseIngressRules([]*networking.Ingress{
				ingressList[5],
			}, []*configurationv1beta1.TCPIngress{})
			assert.Nil(err)
		})
	})
	t.Run("Ingress rules with multiple ports for one Service use separate hostnames for each port", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*networking.Ingress{
			ingressList[6],
		}, []*configurationv1beta1.TCPIngress{})
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal("foo-svc.foo-namespace.8000.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.8000"].Host)
		assert.Nil(err)
	})
}

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
			map[string]string{"configuration.konghq.com/protocol": "grpcs"},
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
			map[string]string{"configuration.konghq.com/protocol": "grpc"},
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
			map[string]string{"configuration.konghq.com/protocol": "grpcs"},
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
			map[string]string{"configuration.konghq.com/protocol": "https"},
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
			map[string]string{"configuration.konghq.com/protocol": "https"},
		},
	}

	for _, testcase := range testTable {
		overrideService(&testcase.inService, &testcase.inKongIngresss, testcase.inAnnotation)
		assert.Equal(testcase.inService, testcase.outService)
	}

	assert.NotPanics(func() {
		overrideService(nil, nil, nil)
	})
}

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
	}

	for _, testcase := range testTable {
		overrideRoute(&testcase.inRoute, &testcase.inKongIngresss)
		assert.Equal(testcase.inRoute, testcase.outRoute)
	}

	assert.NotPanics(func() {
		overrideRoute(nil, nil)
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
	var kongIngress configurationv1.KongIngress
	kongIngress = configurationv1.KongIngress{
		Route: &kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
	}

	var netIngress networking.Ingress

	netIngress = networking.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"configuration.konghq.com/protocols": "grpc,grpcs",
			},
		},
	}

	route = Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
		Ingress: netIngress,
	}
	overrideRoute(&route, &kongIngress)
	assert.Equal(route.Hosts, kong.StringSlice("foo.com", "bar.com"))
	assert.Equal(route.Protocols, kong.StringSlice("grpc", "grpcs"))
}

func TestOverrideRouteByKongIngress(t *testing.T) {
	assert := assert.New(t)
	var route Route
	route = Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
	}
	var kongIngress configurationv1.KongIngress
	kongIngress = configurationv1.KongIngress{
		Route: &kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
	}

	overrideRouteByKongIngress(&route, &kongIngress)
	assert.Equal(route.Hosts, kong.StringSlice("foo.com", "bar.com"))
	assert.NotPanics(func() {
		overrideRoute(nil, nil)
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

	var netIngress networking.Ingress

	netIngress = networking.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"configuration.konghq.com/protocols": "grpc,grpcs",
			},
		},
	}

	route = Route{
		Route: kong.Route{
			Hosts: kong.StringSlice("foo.com", "bar.com"),
		},
		Ingress: netIngress,
	}
	overrideRouteByAnnotation(&route, route.Ingress.GetAnnotations())
	assert.Equal(route.Hosts, kong.StringSlice("foo.com", "bar.com"))
	assert.Equal(route.Protocols, kong.StringSlice("grpc", "grpcs"))

	assert.NotPanics(func() {
		overrideRoute(nil, nil)
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
		normalizeProtocols(&testcase.inRoute)
		assert.Equal(testcase.inRoute.Protocols, testcase.outRoute.Protocols)
	}

	assert.NotPanics(func() {
		overrideUpstream(nil, nil)
	})
}

func TestValidateProtocol(t *testing.T) {
	assert := assert.New(t)
	testTable := []struct {
		input  string
		result bool
	}{
		{"http", true},
		{"https", true},
		{"grpc", true},
		{"grpcs", true},
		{"grcpsfdsafdsfafdshttp", false},
	}
	for _, testcase := range testTable {
		isMatch := validateProtocol(testcase.input)
		assert.Equal(isMatch, testcase.result)
	}

	assert.NotPanics(func() {
		overrideUpstream(nil, nil)
	})
}
func TestOverrideUpstream(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inUpstream     Upstream
		inKongIngresss configurationv1.KongIngress
		outUpstream    Upstream
	}{
		{
			Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			configurationv1.KongIngress{
				Upstream: &kong.Upstream{},
			},
			Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
		},
		{
			Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			configurationv1.KongIngress{
				Upstream: &kong.Upstream{
					Name:               kong.String("wrong.com"),
					HashOn:             kong.String("HashOn"),
					HashOnCookie:       kong.String("HashOnCookie"),
					HashOnCookiePath:   kong.String("HashOnCookiePath"),
					HashOnHeader:       kong.String("HashOnHeader"),
					HashFallback:       kong.String("HashFallback"),
					HashFallbackHeader: kong.String("HashFallbackHeader"),
					Slots:              kong.Int(42),
				},
			},
			Upstream{
				Upstream: kong.Upstream{
					Name:               kong.String("foo.com"),
					HashOn:             kong.String("HashOn"),
					HashOnCookie:       kong.String("HashOnCookie"),
					HashOnCookiePath:   kong.String("HashOnCookiePath"),
					HashOnHeader:       kong.String("HashOnHeader"),
					HashFallback:       kong.String("HashFallback"),
					HashFallbackHeader: kong.String("HashFallbackHeader"),
					Slots:              kong.Int(42),
				},
			},
		},
	}

	for _, testcase := range testTable {
		overrideUpstream(&testcase.inUpstream, &testcase.inKongIngresss)
		assert.Equal(testcase.inUpstream, testcase.outUpstream)
	}

	assert.NotPanics(func() {
		overrideUpstream(nil, nil)
	})
}

func TestGetEndpoints(t *testing.T) {
	tests := []struct {
		name   string
		svc    *corev1.Service
		port   *corev1.ServicePort
		proto  corev1.Protocol
		fn     func(string, string) (*corev1.Endpoints, error)
		result []utils.Endpoint
	}{
		{
			"no service should return 0 endpoints",
			nil,
			nil,
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, nil
			},
			[]utils.Endpoint{},
		},
		{
			"no service port should return 0 endpoints",
			&corev1.Service{},
			nil,
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, nil
			},
			[]utils.Endpoint{},
		},
		{
			"a service without endpoints should return 0 endpoints",
			&corev1.Service{},
			&corev1.ServicePort{Name: "default"},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"a service type ServiceTypeExternalName service with an invalid port should return 0 endpoints",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeExternalName,
				},
			},
			&corev1.ServicePort{Name: "default"},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"a service type ServiceTypeExternalName with a valid port should return one endpoint",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:         corev1.ServiceTypeExternalName,
					ExternalName: "10.0.0.1.xip.io",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{
				{
					Address: "10.0.0.1.xip.io",
					Port:    "80",
				},
			},
		},
		{
			"a service with ingress.kubernetes.io/service-upstream annotation should return one endpoint",
			&corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						"ingress.kubernetes.io/service-upstream": "true",
					},
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeClusterIP,
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(2080),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{
				{
					Address: "foo.bar.svc",
					Port:    "2080",
				},
			},
		},
		{
			"should return no endpoints when there is an error searching for endpoints",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			[]utils.Endpoint{},
		},
		{
			"should return no endpoints when the protocol does not match",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolUDP,
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"should return no endpoints when there is no ready Addresses",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							NotReadyAddresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolUDP,
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"should return no endpoints when the name of the port name do not match any port in the endpoint Subsets",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolTCP,
									Port:     int32(80),
									Name:     "another-name",
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"should return one endpoint when the name of the port name match a port in the endpoint Subsets",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolTCP,
									Port:     int32(80),
									Name:     "default",
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{
				{
					Address: "1.1.1.1",
					Port:    "80",
				},
			},
		},
		{
			"should return one endpoint when the name of the port name match more than one port in the endpoint Subsets",
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
			&corev1.ServicePort{
				Name:       "port-1",
				TargetPort: intstr.FromString("port-1"),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Name:     "port-1",
									Protocol: corev1.ProtocolTCP,
									Port:     80,
								},
								{
									Name:     "port-1",
									Protocol: corev1.ProtocolTCP,
									Port:     80,
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{
				{
					Address: "1.1.1.1",
					Port:    "80",
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := getEndpoints(testCase.svc, testCase.port, testCase.proto, testCase.fn)
			if len(testCase.result) != len(result) {
				t.Errorf("expected %v Endpoints but got %v", testCase.result, len(result))
			}
		})
	}
}

func Test_processCredential(t *testing.T) {
	type args struct {
		credType   string
		consumer   *Consumer
		credConfig interface{}
	}
	tests := []struct {
		name    string
		args    args
		result  *Consumer
		wantErr bool
	}{
		{
			name: "invalid cred type errors",
			args: args{
				credType:   "invalid-type",
				consumer:   &Consumer{},
				credConfig: nil,
			},
			result:  &Consumer{},
			wantErr: true,
		},
		{
			name: "key-auth",
			args: args{
				credType:   "key-auth",
				consumer:   &Consumer{},
				credConfig: map[string]string{"key": "foo"},
			},
			result: &Consumer{
				KeyAuths: []*kong.KeyAuth{
					{
						Key: kong.String("foo"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "keyauth_credential",
			args: args{
				credType:   "keyauth_credential",
				consumer:   &Consumer{},
				credConfig: map[string]string{"key": "foo"},
			},
			result: &Consumer{
				KeyAuths: []*kong.KeyAuth{
					{
						Key: kong.String("foo"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "basic-auth",
			args: args{
				credType: "basic-auth",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"username": "foo",
					"password": "bar",
				},
			},
			result: &Consumer{
				BasicAuths: []*kong.BasicAuth{
					{
						Username: kong.String("foo"),
						Password: kong.String("bar"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "basicauth_credential",
			args: args{
				credType: "basicauth_credential",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"username": "foo",
					"password": "bar",
				},
			},
			result: &Consumer{
				BasicAuths: []*kong.BasicAuth{
					{
						Username: kong.String("foo"),
						Password: kong.String("bar"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "hmac-auth",
			args: args{
				credType: "hmac-auth",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"username": "foo",
					"secret":   "bar",
				},
			},
			result: &Consumer{
				HMACAuths: []*kong.HMACAuth{
					{
						Username: kong.String("foo"),
						Secret:   kong.String("bar"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "hmacauth_credential",
			args: args{
				credType: "hmacauth_credential",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"username": "foo",
					"secret":   "bar",
				},
			},
			result: &Consumer{
				HMACAuths: []*kong.HMACAuth{
					{
						Username: kong.String("foo"),
						Secret:   kong.String("bar"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "oauth2",
			args: args{
				credType: "oauth2",
				consumer: &Consumer{},
				credConfig: map[string]interface{}{
					"name":          "foo",
					"client_id":     "bar",
					"client_secret": "baz",
					"redirect_uris": []string{"example.com"},
				},
			},
			result: &Consumer{
				Oauth2Creds: []*kong.Oauth2Credential{
					{
						Name:         kong.String("foo"),
						ClientID:     kong.String("bar"),
						ClientSecret: kong.String("baz"),
						RedirectURIs: kong.StringSlice("example.com"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "jwt",
			args: args{
				credType: "jwt",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"key":            "foo",
					"rsa_public_key": "bar",
					"secret":         "baz",
				},
			},
			result: &Consumer{
				JWTAuths: []*kong.JWTAuth{
					{
						Key:          kong.String("foo"),
						RSAPublicKey: kong.String("bar"),
						Secret:       kong.String("baz"),
						// set by default
						Algorithm: kong.String("HS256"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "jwt_secret",
			args: args{
				credType: "jwt_secret",
				consumer: &Consumer{},
				credConfig: map[string]string{
					"key":            "foo",
					"rsa_public_key": "bar",
					"secret":         "baz",
				},
			},
			result: &Consumer{
				JWTAuths: []*kong.JWTAuth{
					{
						Key:          kong.String("foo"),
						RSAPublicKey: kong.String("bar"),
						Secret:       kong.String("baz"),
						// set by default
						Algorithm: kong.String("HS256"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "acl",
			args: args{
				credType:   "acl",
				consumer:   &Consumer{},
				credConfig: map[string]string{"group": "group-foo"},
			},
			result: &Consumer{
				ACLGroups: []*kong.ACLGroup{
					{
						Group: kong.String("group-foo"),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := processCredential(tt.args.credType, tt.args.consumer,
				tt.args.credConfig); (err != nil) != tt.wantErr {
				t.Errorf("processCredential() error = %v, wantErr %v",
					err, tt.wantErr)
			}
			assert.Equal(t, tt.result, tt.args.consumer)
		})
	}
}

func Test_getPluginRelations(t *testing.T) {
	type args struct {
		state KongState
	}
	tests := []struct {
		name string
		args args
		want map[string]foreignRelations
	}{
		{
			name: "empty state",
			want: map[string]foreignRelations{},
		},
		{
			name: "single consumer annotation",
			args: args{
				state: KongState{
					Consumers: []Consumer{
						{
							Consumer: kong.Consumer{
								Username: kong.String("foo-consumer"),
							},
							k8sKongConsumer: configurationv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										"plugins.konghq.com": "foo,bar",
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
				"ns1:foo": {Consumer: []string{"foo-consumer"}},
				"ns1:bar": {Consumer: []string{"foo-consumer"}},
			},
		},
		{
			name: "single service annotation",
			args: args{
				state: KongState{
					Services: []Service{
						{
							Service: kong.Service{
								Name: kong.String("foo-service"),
							},
							K8sService: corev1.Service{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										"plugins.konghq.com": "foo,bar",
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
				"ns1:foo": {Service: []string{"foo-service"}},
				"ns1:bar": {Service: []string{"foo-service"}},
			},
		},
		{
			name: "single Ingress annotation",
			args: args{
				state: KongState{
					Services: []Service{
						{
							Service: kong.Service{
								Name: kong.String("foo-service"),
							},
							Routes: []Route{
								{
									Route: kong.Route{
										Name: kong.String("foo-route"),
									},
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												"plugins.konghq.com": "foo,bar",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
				"ns2:foo": {Route: []string{"foo-route"}},
				"ns2:bar": {Route: []string{"foo-route"}},
			},
		},
		{
			name: "multiple routes with annotation",
			args: args{
				state: KongState{
					Services: []Service{
						{
							Service: kong.Service{
								Name: kong.String("foo-service"),
							},
							Routes: []Route{
								{
									Route: kong.Route{
										Name: kong.String("foo-route"),
									},
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												"plugins.konghq.com": "foo,bar",
											},
										},
									},
								},
								{
									Route: kong.Route{
										Name: kong.String("bar-route"),
									},
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												"plugins.konghq.com": "bar,baz",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
				"ns2:foo": {Route: []string{"foo-route"}},
				"ns2:bar": {Route: []string{"foo-route", "bar-route"}},
				"ns2:baz": {Route: []string{"bar-route"}},
			},
		},
		{
			name: "multiple consumers, routes and services",
			args: args{
				state: KongState{
					Consumers: []Consumer{
						{
							Consumer: kong.Consumer{
								Username: kong.String("foo-consumer"),
							},
							k8sKongConsumer: configurationv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										"plugins.konghq.com": "foo,bar",
									},
								},
							},
						},
						{
							Consumer: kong.Consumer{
								Username: kong.String("foo-consumer"),
							},
							k8sKongConsumer: configurationv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns2",
									Annotations: map[string]string{
										"plugins.konghq.com": "foo,bar",
									},
								},
							},
						},
						{
							Consumer: kong.Consumer{
								Username: kong.String("bar-consumer"),
							},
							k8sKongConsumer: configurationv1.KongConsumer{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										"plugins.konghq.com": "foobar",
									},
								},
							},
						},
					},
					Services: []Service{
						{
							Service: kong.Service{
								Name: kong.String("foo-service"),
							},
							K8sService: corev1.Service{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "ns1",
									Annotations: map[string]string{
										"plugins.konghq.com": "foo,bar",
									},
								},
							},
							Routes: []Route{
								{
									Route: kong.Route{
										Name: kong.String("foo-route"),
									},
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												"plugins.konghq.com": "foo,bar",
											},
										},
									},
								},
								{
									Route: kong.Route{
										Name: kong.String("bar-route"),
									},
									Ingress: networking.Ingress{
										ObjectMeta: metav1.ObjectMeta{
											Name:      "some-ingress",
											Namespace: "ns2",
											Annotations: map[string]string{
												"plugins.konghq.com": "bar,baz",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]foreignRelations{
				"ns1:foo":    {Consumer: []string{"foo-consumer"}, Service: []string{"foo-service"}},
				"ns1:bar":    {Consumer: []string{"foo-consumer"}, Service: []string{"foo-service"}},
				"ns1:foobar": {Consumer: []string{"bar-consumer"}},
				"ns2:foo":    {Consumer: []string{"foo-consumer"}, Route: []string{"foo-route"}},
				"ns2:bar":    {Consumer: []string{"foo-consumer"}, Route: []string{"foo-route", "bar-route"}},
				"ns2:baz":    {Route: []string{"bar-route"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPluginRelations(tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPluginRelations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCombinations(t *testing.T) {
	type args struct {
		relations foreignRelations
	}
	tests := []struct {
		name string
		args args
		want []rel
	}{
		{
			name: "empty",
			args: args{
				relations: foreignRelations{},
			},
			want: nil,
		},
		{
			name: "plugins on consumer only",
			args: args{
				relations: foreignRelations{
					Consumer: []string{"foo", "bar"},
				},
			},
			want: []rel{
				{
					Consumer: "foo",
				},
				{
					Consumer: "bar",
				},
			},
		},
		{
			name: "plugins on service only",
			args: args{
				relations: foreignRelations{
					Service: []string{"foo", "bar"},
				},
			},
			want: []rel{
				{
					Service: "foo",
				},
				{
					Service: "bar",
				},
			},
		},
		{
			name: "plugins on routes only",
			args: args{
				relations: foreignRelations{
					Route: []string{"foo", "bar"},
				},
			},
			want: []rel{
				{
					Route: "foo",
				},
				{
					Route: "bar",
				},
			},
		},
		{
			name: "plugins on service and routes only",
			args: args{
				relations: foreignRelations{
					Route:   []string{"foo", "bar"},
					Service: []string{"foo", "bar"},
				},
			},
			want: []rel{
				{
					Service: "foo",
				},
				{
					Service: "bar",
				},
				{
					Route: "foo",
				},
				{
					Route: "bar",
				},
			},
		},
		{
			name: "plugins on combination of route and consumer",
			args: args{
				relations: foreignRelations{
					Route:    []string{"foo", "bar"},
					Consumer: []string{"foo", "bar"},
				},
			},
			want: []rel{
				{
					Consumer: "foo",
					Route:    "foo",
				},
				{
					Consumer: "bar",
					Route:    "foo",
				},
				{
					Consumer: "foo",
					Route:    "bar",
				},
				{
					Consumer: "bar",
					Route:    "bar",
				},
			},
		},
		{
			name: "plugins on combination of service and consumer",
			args: args{
				relations: foreignRelations{
					Service:  []string{"foo", "bar"},
					Consumer: []string{"foo", "bar"},
				},
			},
			want: []rel{
				{
					Consumer: "foo",
					Service:  "foo",
				},
				{
					Consumer: "bar",
					Service:  "foo",
				},
				{
					Consumer: "foo",
					Service:  "bar",
				},
				{
					Consumer: "bar",
					Service:  "bar",
				},
			},
		},
		{
			name: "plugins on combination of service,route and consumer",
			args: args{
				relations: foreignRelations{
					Consumer: []string{"c1", "c2"},
					Route:    []string{"r1", "r2"},
					Service:  []string{"s1", "s2"},
				},
			},
			want: []rel{
				{
					Consumer: "c1",
					Service:  "s1",
				},
				{
					Consumer: "c2",
					Service:  "s1",
				},
				{
					Consumer: "c1",
					Service:  "s2",
				},
				{
					Consumer: "c2",
					Service:  "s2",
				},
				{
					Consumer: "c1",
					Route:    "r1",
				},
				{
					Consumer: "c2",
					Route:    "r1",
				},
				{
					Consumer: "c1",
					Route:    "r2",
				},
				{
					Consumer: "c2",
					Route:    "r2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCombinations(tt.args.relations); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCombinations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processTLSSections(t *testing.T) {
	type args struct {
		tlsSections []networking.IngressTLS
		namespace   string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			args: args{
				tlsSections: []networking.IngressTLS{
					{
						Hosts: []string{
							"1.example.com",
							"2.example.com",
						},
						SecretName: "sooper-secret",
					},
					{
						Hosts: []string{
							"3.example.com",
							"4.example.com",
						},
						SecretName: "sooper-secret2",
					},
				},
				namespace: "foo",
			},
			want: map[string][]string{
				"foo/sooper-secret":  {"1.example.com", "2.example.com"},
				"foo/sooper-secret2": {"3.example.com", "4.example.com"},
			},
		},
		{
			args: args{
				tlsSections: []networking.IngressTLS{
					{
						Hosts: []string{
							"1.example.com",
						},
						SecretName: "sooper-secret",
					},
					{
						Hosts: []string{
							"3.example.com",
							"1.example.com",
							"4.example.com",
						},
						SecretName: "sooper-secret2",
					},
				},
				namespace: "foo",
			},
			want: map[string][]string{
				"foo/sooper-secret":  {"1.example.com"},
				"foo/sooper-secret2": {"3.example.com", "4.example.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := map[string][]string{}
			processTLSSections(tt.args.tlsSections, tt.args.namespace, got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processTLSSections() = %v, want %v", got, tt.want)
			}
		})
	}
}
