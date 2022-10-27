package parser

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

func TestFromIngressV1beta1(t *testing.T) {
	assert := assert.New(t)
	ingressList := []*netv1beta1.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: netv1beta1.IngressSpec{
				Rules: []netv1beta1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1beta1.IngressRuleValue{
							HTTP: &netv1beta1.HTTPIngressRuleValue{
								Paths: []netv1beta1.HTTPIngressPath{
									{
										Path: "/",
										Backend: netv1beta1.IngressBackend{
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
			Spec: netv1beta1.IngressSpec{
				TLS: []netv1beta1.IngressTLS{
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
				Rules: []netv1beta1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1beta1.IngressRuleValue{
							HTTP: &netv1beta1.HTTPIngressRuleValue{
								Paths: []netv1beta1.HTTPIngressPath{
									{
										Path: "/",
										Backend: netv1beta1.IngressBackend{
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
			Spec: netv1beta1.IngressSpec{
				Backend: &netv1beta1.IngressBackend{
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
			Spec: netv1beta1.IngressSpec{
				Rules: []netv1beta1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1beta1.IngressRuleValue{
							HTTP: &netv1beta1.HTTPIngressRuleValue{
								Paths: []netv1beta1.HTTPIngressPath{
									{
										Path: "/.well-known/acme-challenge/yolo",
										Backend: netv1beta1.IngressBackend{
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
			Spec: netv1beta1.IngressSpec{
				Rules: []netv1beta1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1beta1.IngressRuleValue{
							HTTP: &netv1beta1.HTTPIngressRuleValue{
								Paths: []netv1beta1.HTTPIngressPath{
									{
										Backend: netv1beta1.IngressBackend{
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
			Spec: netv1beta1.IngressSpec{
				Rules: []netv1beta1.IngressRule{
					{
						Host:             "example.com",
						IngressRuleValue: netv1beta1.IngressRuleValue{},
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
			Spec: netv1beta1.IngressSpec{
				Rules: []netv1beta1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1beta1.IngressRuleValue{
							HTTP: &netv1beta1.HTTPIngressRuleValue{
								Paths: []netv1beta1.HTTPIngressPath{
									{
										Backend: netv1beta1.IngressBackend{
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
						IngressRuleValue: netv1beta1.IngressRuleValue{
							HTTP: &netv1beta1.HTTPIngressRuleValue{
								Paths: []netv1beta1.HTTPIngressPath{
									{
										Backend: netv1beta1.IngressBackend{
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
			Spec: netv1beta1.IngressSpec{
				Rules: []netv1beta1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1beta1.IngressRuleValue{
							HTTP: &netv1beta1.HTTPIngressRuleValue{
								Paths: []netv1beta1.HTTPIngressPath{
									{
										Path: "/foo//bar",
										Backend: netv1beta1.IngressBackend{
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
		// 8
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "regex-prefix",
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: netv1beta1.IngressSpec{
				Rules: []netv1beta1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1beta1.IngressRuleValue{
							HTTP: &netv1beta1.HTTPIngressRuleValue{
								Paths: []netv1beta1.HTTPIngressPath{
									{
										Path: translators.ControllerPathRegexPrefix + "/foo/\\d{3}",
										Backend: netv1beta1.IngressBackend{
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
			IngressesV1beta1: []*netv1beta1.Ingress{},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("simple ingress rule is parsed", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*netv1beta1.Ingress{
				ingressList[0],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])
	})
	t.Run("ingress rule with default backend", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*netv1beta1.Ingress{ingressList[0], ingressList[2]},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

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
			IngressesV1beta1: []*netv1beta1.Ingress{
				ingressList[1],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret2"]))
	})
	t.Run("ingress rule with ACME like path has strip_path set to false", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*netv1beta1.Ingress{
				ingressList[3],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

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
			IngressesV1beta1: []*netv1beta1.Ingress{
				ingressList[4],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])
	})
	t.Run("empty Ingress rule doesn't cause a panic", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*netv1beta1.Ingress{
				ingressList[5],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		assert.NotPanics(func() {
			p.ingressRulesFromIngressV1beta1()
		})
	})
	t.Run("Ingress rules with multiple ports for one Service use separate hostnames for each port", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*netv1beta1.Ingress{
				ingressList[6],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal("foo-svc.foo-namespace.8000.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.8000"].Host)
	})
	t.Run("Ingress rule with path containing multiple slashes ('//') is skipped", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*netv1beta1.Ingress{
				ingressList[7],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Empty(parsedInfo.ServiceNameToServices)
	})
	t.Run("Ingress rule with regex prefixed path creates route with Kong regex prefix", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1beta1: []*netv1beta1.Ingress{
				ingressList[8],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1beta1()
		assert.Equal(translators.KongPathRegexPrefix+"/foo/\\d{3}", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
	})
}

func TestFromIngressV1(t *testing.T) {
	assert := assert.New(t)
	ingressList := []*netv1.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
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
												Port: netv1.ServiceBackendPort{Number: 80},
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
			Spec: netv1.IngressSpec{
				TLS: []netv1.IngressTLS{
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
												Port: netv1.ServiceBackendPort{Number: 80},
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
			Spec: netv1.IngressSpec{
				DefaultBackend: &netv1.IngressBackend{
					Service: &netv1.IngressServiceBackend{
						Name: "default-svc",
						Port: netv1.ServiceBackendPort{Number: 80},
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
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path: "/.well-known/acme-challenge/yolo",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "cert-manager-solver-pod",
												Port: netv1.ServiceBackendPort{Number: 80},
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
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: netv1.ServiceBackendPort{Number: 80},
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
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host:             "example.com",
						IngressRuleValue: netv1.IngressRuleValue{},
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
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: netv1.ServiceBackendPort{Number: 80},
											},
										},
									},
								},
							},
						},
					},
					{
						Host: "example.net",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: netv1.ServiceBackendPort{Number: 8000},
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
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path: "/foo//bar",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: netv1.ServiceBackendPort{Number: 80},
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
												Port: netv1.ServiceBackendPort{Name: "http"},
											},
										},
									},
									{
										Path: "/ws",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: netv1.ServiceBackendPort{Name: "ws"},
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
		// 9
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "regex-prefix",
				Namespace: "foo-namespace",
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
										Path: translators.ControllerPathRegexPrefix + "/foo/\\d{3}",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: netv1.ServiceBackendPort{Number: 80},
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
			IngressesV1: []*netv1.Ingress{},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("simple ingress rule is parsed", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				ingressList[0],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Hosts[0])
	})
	t.Run("ingress rule with default backend", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				ingressList[0],
				ingressList[2],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

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
			IngressesV1: []*netv1.Ingress{
				ingressList[1],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret2"]))
	})
	t.Run("ingress rule with ACME like path has strip_path set to false", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				ingressList[3],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

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
			IngressesV1: []*netv1.Ingress{
				ingressList[4],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Hosts[0])
	})
	t.Run("empty Ingress rule doesn't cause a panic", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				ingressList[5],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		assert.NotPanics(func() {
			p.ingressRulesFromIngressV1()
		})
	})
	t.Run("Ingress rules with multiple ports for one Service use separate hostnames for each port", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				ingressList[6],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal("foo-svc.foo-namespace.80.svc",
			*parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Host)
		assert.Equal("foo-svc.foo-namespace.8000.svc",
			*parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-8000"].Host)
	})
	t.Run("Ingress rule with path containing multiple slashes ('//') is skipped", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				ingressList[7],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Empty(parsedInfo.ServiceNameToServices)
	})
	t.Run("Ingress rule with ports defined by name", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				ingressList[9],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1()
		_, ok := parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"]
		assert.True(ok)
	})
	t.Run("Ingress rule with regex prefixed path creates route with Kong regex prefix", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				ingressList[9],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(translators.KongPathRegexPrefix+"/foo/\\d{3}", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Paths[0])
	})
}

func TestFromIngressV1_RegexPrefix(t *testing.T) {
	assert := assert.New(t)
	pathTypeExact := netv1.PathTypeExact
	ingressList := []*netv1.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
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
										Path:     "/whatever",
										PathType: &pathTypeExact,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "foo-svc",
												Port: netv1.ServiceBackendPort{Number: 80},
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
	t.Run("exact rule results in prefixed regex", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				ingressList[0],
			},
		})
		require.NoError(t, err)
		p, err := NewParser(logrus.New(), store)
		require.NoError(t, err)

		p.EnableRegexPathPrefix()

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal("~/whatever$", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Paths[0])
	})
}
