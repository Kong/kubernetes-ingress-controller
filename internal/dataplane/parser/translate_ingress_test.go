package parser

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

func TestFromIngressV1beta1(t *testing.T) {
	assert := assert.New(t)
	ingressList := []*networkingv1beta1.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
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
		// 1
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ing-with-tls",
				Namespace: "bar-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: networkingv1beta1.IngressSpec{
				TLS: []networkingv1beta1.IngressTLS{
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
		// 2
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ing-with-default-backend",
				Namespace: "bar-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: networkingv1beta1.IngressSpec{
				Backend: &networkingv1beta1.IngressBackend{
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
										Path: "/.well-known/acme-challenge/yolo",
										Backend: networkingv1beta1.IngressBackend{
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
		// 5
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "baz",
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: networkingv1beta1.IngressSpec{
				Rules: []networkingv1beta1.IngressRule{
					{
						Host:             "example.com",
						IngressRuleValue: networkingv1beta1.IngressRuleValue{},
					},
				},
			},
		},
		// 6
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
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
										Backend: networkingv1beta1.IngressBackend{
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
						IngressRuleValue: networkingv1beta1.IngressRuleValue{
							HTTP: &networkingv1beta1.HTTPIngressRuleValue{
								Paths: []networkingv1beta1.HTTPIngressPath{
									{
										Backend: networkingv1beta1.IngressBackend{
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
		// 7
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "invalid-path",
				Namespace: "foo-namespace",
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
										Path: "/foo//bar",
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
	}

	t.Run("no ingress returns empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*networkingv1beta1.Ingress{},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("simple ingress rule is parsed", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*networkingv1beta1.Ingress{
				ingressList[0],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])
	})
	t.Run("ingress rule with default backend", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*networkingv1beta1.Ingress{ingressList[0], ingressList[2]},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal(2, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])

		assert.Equal(1, len(parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes))
		assert.Equal("/", *parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes[0].Paths[0])
		assert.Equal(0, len(parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes[0].Hosts))
	})
	t.Run("ingress rule with TLS", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*networkingv1beta1.Ingress{
				ingressList[1],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret2"]))
	})
	t.Run("ingress rule with ACME like path has strip_path set to false", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*networkingv1beta1.Ingress{
				ingressList[3],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("cert-manager-solver-pod.foo-namespace.80.svc",
			*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Port)

		assert.Equal("/.well-known/acme-challenge/yolo",
			*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Routes[0].Paths[0])
		assert.Equal("example.com",
			*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Routes[0].Hosts[0])
		assert.False(*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Routes[0].StripPath)
	})
	t.Run("ingress with empty path is correctly parsed", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*networkingv1beta1.Ingress{
				ingressList[4],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])
	})
	t.Run("empty Ingress rule doesn't cause a panic", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*networkingv1beta1.Ingress{
				ingressList[5],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		assert.NotPanics(func() {
			p.ingressRulesFromIngressV1beta1()
		})
	})
	t.Run("Ingress rules with multiple ports for one Service use separate hostnames for each port", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*networkingv1beta1.Ingress{
				ingressList[6],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal("foo-svc.foo-namespace.8000.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.8000"].Host)
	})
	t.Run("Ingress rule with path containing multiple slashes ('//') is skipped", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*networkingv1beta1.Ingress{
				ingressList[7],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Empty(parsedInfo.ServiceNameToServices)
	})
}

func TestFromIngressV1(t *testing.T) {
	assert := assert.New(t)
	ingressList := []*networkingv1.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
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
												Port: networkingv1.ServiceBackendPort{Number: 80},
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
		// 1
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ing-with-tls",
				Namespace: "bar-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: networkingv1.IngressSpec{
				TLS: []networkingv1.IngressTLS{
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
												Port: networkingv1.ServiceBackendPort{Number: 80},
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
		// 2
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ing-with-default-backend",
				Namespace: "bar-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: networkingv1.IngressSpec{
				DefaultBackend: &networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: "default-svc",
						Port: networkingv1.ServiceBackendPort{Number: 80},
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
										Path: "/.well-known/acme-challenge/yolo",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "cert-manager-solver-pod",
												Port: networkingv1.ServiceBackendPort{Number: 80},
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
		// 4
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
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
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: networkingv1.ServiceBackendPort{Number: 80},
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
		// 5
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "baz",
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						Host:             "example.com",
						IngressRuleValue: networkingv1.IngressRuleValue{},
					},
				},
			},
		},
		// 6
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
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
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: networkingv1.ServiceBackendPort{Number: 80},
											},
										},
									},
								},
							},
						},
					},
					{
						Host: "example.net",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: networkingv1.ServiceBackendPort{Number: 8000},
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
		// 7
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "invalid-path",
				Namespace: "foo-namespace",
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
										Path: "/foo//bar",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: networkingv1.ServiceBackendPort{Number: 80},
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
		// 8
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
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
												Port: networkingv1.ServiceBackendPort{Name: "http"},
											},
										},
									},
									{
										Path: "/ws",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: networkingv1.ServiceBackendPort{Name: "ws"},
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

	t.Run("no ingress returns empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("simple ingress rule is parsed", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{
				ingressList[0],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Hosts[0])
	})
	t.Run("ingress rule with default backend", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{
				ingressList[0],
				ingressList[2],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(2, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Hosts[0])

		assert.Equal(1, len(parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes))
		assert.Equal("/", *parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes[0].Paths[0])
		assert.Equal(0, len(parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes[0].Hosts))
	})
	t.Run("ingress rule with TLS", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{
				ingressList[1],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret2"]))
	})
	t.Run("ingress rule with ACME like path has strip_path set to false", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{
				ingressList[3],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("cert-manager-solver-pod.foo-namespace.80.svc",
			*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.pnum-80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.pnum-80"].Port)

		assert.Equal("/.well-known/acme-challenge/yolo",
			*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.pnum-80"].Routes[0].Paths[0])
		assert.Equal("example.com",
			*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.pnum-80"].Routes[0].Hosts[0])
		assert.False(*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.pnum-80"].Routes[0].StripPath)
	})
	t.Run("ingress with empty path is correctly parsed", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{
				ingressList[4],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Hosts[0])
	})
	t.Run("empty Ingress rule doesn't cause a panic", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{
				ingressList[5],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		assert.NotPanics(func() {
			p.ingressRulesFromIngressV1()
		})
	})
	t.Run("Ingress rules with multiple ports for one Service use separate hostnames for each port", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{
				ingressList[6],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal("foo-svc.foo-namespace.80.svc",
			*parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Host)
		assert.Equal("foo-svc.foo-namespace.8000.svc",
			*parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-8000"].Host)
	})
	t.Run("Ingress rule with path containing multiple slashes ('//') is skipped", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{
				ingressList[7],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Empty(parsedInfo.ServiceNameToServices)
	})
	t.Run("Ingress rule with ports defined by name", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*networkingv1.Ingress{
				ingressList[8],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromIngressV1()
		_, ok := parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pname-http"]
		assert.True(ok)
		_, ok = parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pname-ws"]
		assert.True(ok)
	})
}
