package parser

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

func TestFromKnativeIngress(t *testing.T) {
	assert := assert.New(t)
	ingressList := []*knative.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.KnativeIngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{},
				},
			},
		},
		// 1
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.KnativeIngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{
						Hosts: []string{"my-func.example.com"},
						HTTP: &knative.HTTPIngressRuleValue{
							Paths: []knative.HTTPIngressPath{
								{
									Path: "/",
									AppendHeaders: map[string]string{
										"foo": "bar",
									},
									Splits: []knative.IngressBackendSplit{
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "foo-ns",
												ServiceName:      "foo-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 100,
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
				Name:      "foo",
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.KnativeIngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{
						Hosts: []string{"my-func.example.com"},
						HTTP: &knative.HTTPIngressRuleValue{
							Paths: []knative.HTTPIngressPath{
								{
									Path: "/",
									AppendHeaders: map[string]string{
										"foo": "bar",
									},
									Splits: []knative.IngressBackendSplit{
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "bar-ns",
												ServiceName:      "bar-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 20,
										},
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "foo-ns",
												ServiceName:      "foo-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 100,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// 3
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.KnativeIngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{
						Hosts: []string{"my-func.example.com"},
						HTTP: &knative.HTTPIngressRuleValue{
							Paths: []knative.HTTPIngressPath{
								{
									Path: "/",
									AppendHeaders: map[string]string{
										"foo": "bar",
									},
									Splits: []knative.IngressBackendSplit{
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "bar-ns",
												ServiceName:      "bar-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 20,
										},
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "foo-ns",
												ServiceName:      "foo-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 100,
										},
									},
								},
							},
						},
					},
				},
				TLS: []knative.IngressTLS{
					{
						Hosts: []string{
							"foo.example.com",
							"foo1.example.com",
						},
						SecretName: "foo-secret",
					},
					{
						Hosts: []string{
							"bar.example.com",
							"bar1.example.com",
						},
						SecretName: "bar-secret",
					},
				},
			},
		},
	}
	t.Run("no ingress returns empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: []*knative.Ingress{},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromKnativeIngress()
		assert.Equal(map[string]kongstate.Service{}, parsedInfo.ServiceNameToServices)
		assert.Equal(newSecretNameToSNIs(), parsedInfo.SecretNameToSNIs)
	})
	t.Run("empty ingress returns empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: []*knative.Ingress{
				ingressList[0],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromKnativeIngress()
		assert.Equal(map[string]kongstate.Service{}, parsedInfo.ServiceNameToServices)
		assert.Equal(newSecretNameToSNIs(), parsedInfo.SecretNameToSNIs)
	})
	t.Run("basic knative Ingress resource is parsed", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: []*knative.Ingress{
				ingressList[1],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromKnativeIngress()
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		svc := parsedInfo.ServiceNameToServices["foo-ns.foo-svc.42"]
		assert.Equal(kong.Service{
			Name:           kong.String("foo-ns.foo-svc.42"),
			Port:           kong.Int(80),
			Host:           kong.String("foo-svc.foo-ns.42.svc"),
			Path:           kong.String("/"),
			Protocol:       kong.String("http"),
			WriteTimeout:   kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			ConnectTimeout: kong.Int(60000),
			Retries:        kong.Int(5),
		}, svc.Service)
		assert.Equal(kong.Route{
			Name:              kong.String("foo-namespace.foo.00"),
			RegexPriority:     kong.Int(0),
			StripPath:         kong.Bool(false),
			Paths:             kong.StringSlice("/"),
			PreserveHost:      kong.Bool(true),
			Protocols:         kong.StringSlice("http", "https"),
			Hosts:             kong.StringSlice("my-func.example.com"),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
		}, svc.Routes[0].Route)
		assert.Equal(kong.Plugin{
			Name: kong.String("request-transformer"),
			Config: kong.Configuration{
				"add": map[string]interface{}{
					"headers": []string{"foo:bar"},
				},
			},
		}, svc.Plugins[0])

		assert.Equal(newSecretNameToSNIs(), parsedInfo.SecretNameToSNIs)
	})
	t.Run("knative TLS section is correctly parsed", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: []*knative.Ingress{
				ingressList[3],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromKnativeIngress()
		assert.Equal(SecretNameToSNIs(map[string][]string{
			"foo-namespace/bar-secret": {"bar.example.com", "bar1.example.com"},
			"foo-namespace/foo-secret": {"foo.example.com", "foo1.example.com"},
		}), parsedInfo.SecretNameToSNIs)
	})
	t.Run("split knative Ingress resource chooses the highest split", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			KnativeIngresses: []*knative.Ingress{
				ingressList[2],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromKnativeIngress()
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		svc := parsedInfo.ServiceNameToServices["foo-ns.foo-svc.42"]
		assert.Equal(kong.Service{
			Name:           kong.String("foo-ns.foo-svc.42"),
			Port:           kong.Int(80),
			Host:           kong.String("foo-svc.foo-ns.42.svc"),
			Path:           kong.String("/"),
			Protocol:       kong.String("http"),
			WriteTimeout:   kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			ConnectTimeout: kong.Int(60000),
			Retries:        kong.Int(5),
		}, svc.Service)
		assert.Equal(kong.Route{
			Name:              kong.String("foo-namespace.foo.00"),
			RegexPriority:     kong.Int(0),
			StripPath:         kong.Bool(false),
			Paths:             kong.StringSlice("/"),
			PreserveHost:      kong.Bool(true),
			Protocols:         kong.StringSlice("http", "https"),
			Hosts:             kong.StringSlice("my-func.example.com"),
			ResponseBuffering: kong.Bool(true),
			RequestBuffering:  kong.Bool(true),
		}, svc.Routes[0].Route)
		assert.Equal(kong.Plugin{
			Name: kong.String("request-transformer"),
			Config: kong.Configuration{
				"add": map[string]interface{}{
					"headers": []string{"foo:bar"},
				},
			},
		}, svc.Plugins[0])

		assert.Equal(newSecretNameToSNIs(), parsedInfo.SecretNameToSNIs)
	})
}
