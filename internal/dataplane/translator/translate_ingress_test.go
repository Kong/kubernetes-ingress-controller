package translator

import (
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/translator/subtranslator"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestFromIngressV1(t *testing.T) {
	t.Run("no ingress returns empty info", func(t *testing.T) {
		s, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{},
		})
		require.NoError(t, err)

		translatedInfo := mustNewTranslator(t, s).ingressRulesFromIngressV1()
		assert.Equal(t, ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			ServiceNameToParent:   make(map[string]client.Object),
			SecretNameToSNIs:      newSecretNameToSNIs(),
		}, translatedInfo)
	})

	t.Run("empty Ingress rule doesn't cause a panic", func(t *testing.T) {
		s, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
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
			},
		})
		require.NoError(t, err)
		p := mustNewTranslator(t, s)

		assert.NotPanics(t, func() {
			p.ingressRulesFromIngressV1()
		})
	})

	t.Run("Ingress with default backend targeting the same Service doesn't overwrite routes", func(t *testing.T) {
		s, err := store.NewFakeStore(store.FakeObjects{
			Services: []*corev1.Service{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc",
						Namespace: "my-ns",
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{
							{
								Name: "http",
								Port: 80,
							},
						},
					},
				},
			},
			IngressesV1: []*netv1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "baz",
						Namespace: "my-ns",
						Annotations: map[string]string{
							annotations.IngressClassKey: annotations.DefaultIngressClass,
						},
					},
					Spec: netv1.IngressSpec{
						DefaultBackend: &netv1.IngressBackend{
							Service: &netv1.IngressServiceBackend{
								Name: "svc",
								Port: netv1.ServiceBackendPort{
									Number: 80,
								},
							},
						},
						Rules: []netv1.IngressRule{
							{
								IngressRuleValue: netv1.IngressRuleValue{
									HTTP: &netv1.HTTPIngressRuleValue{
										Paths: []netv1.HTTPIngressPath{
											{
												Path:     "/foo",
												PathType: lo.ToPtr(netv1.PathTypePrefix),
												Backend: netv1.IngressBackend{
													Service: &netv1.IngressServiceBackend{
														Name: "svc",
														Port: netv1.ServiceBackendPort{
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
			},
		})
		require.NoError(t, err)
		p := mustNewTranslator(t, s)

		result := p.ingressRulesFromIngressV1()
		require.Contains(t, result.ServiceNameToServices, "my-ns.svc.80")
		snts := result.ServiceNameToServices["my-ns.svc.80"]
		require.Len(t, snts.Routes, 2, "ServiceNameToServices should have 2 routes: 1 for default and 1 for /foo")
		if assert.Len(t, snts.Routes[0].Paths, 2) {
			assert.Equal(t, *snts.Routes[0].Paths[0], "/foo/")
			assert.Equal(t, *snts.Routes[0].Paths[1], "~/foo$")
		}
		if assert.Len(t, snts.Routes[1].Paths, 1) {
			assert.Equal(t, *snts.Routes[1].Paths[0], "/")
		}

		require.Len(t, snts.Backends, 1)
		assert.Equal(t, snts.Backends[0].Name, "svc")

		require.Contains(t, result.ServiceNameToParent, "my-ns.svc.80")
		assert.Equal(t, result.ServiceNameToParent["my-ns.svc.80"].GetName(), "baz")
	})
}

func TestGetDefaultBackendService(t *testing.T) {
	someIngress := func(creationTimestamp time.Time, serviceName string) netv1.Ingress {
		return netv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "foo",
				Namespace:         "foo-namespace",
				CreationTimestamp: metav1.NewTime(creationTimestamp),
			},
			Spec: netv1.IngressSpec{
				DefaultBackend: &netv1.IngressBackend{
					Service: &netv1.IngressServiceBackend{
						Name: serviceName,
						Port: netv1.ServiceBackendPort{Number: 80},
					},
				},
			},
		}
	}

	now := time.Now()
	testCases := []struct {
		name                       string
		ingresses                  []netv1.Ingress
		expressionRoutes           bool
		expectedHaveBackendService bool
		expectedServiceName        string
		expectedServiceHost        string
	}{
		{
			name:                       "no ingresses",
			ingresses:                  []netv1.Ingress{},
			expressionRoutes:           false,
			expectedHaveBackendService: false,
		},
		{
			name:                       "no ingresses with expression routes",
			ingresses:                  []netv1.Ingress{},
			expressionRoutes:           true,
			expectedHaveBackendService: false,
		},
		{
			name:                       "one ingress with default backend",
			ingresses:                  []netv1.Ingress{someIngress(now, "foo-svc")},
			expressionRoutes:           false,
			expectedHaveBackendService: true,
			expectedServiceName:        "foo-namespace.foo-svc.80",
			expectedServiceHost:        "foo-svc.foo-namespace.80.svc",
		},
		{
			name:                       "one ingress with default backend and expression routes enabled",
			ingresses:                  []netv1.Ingress{someIngress(now, "foo-svc")},
			expressionRoutes:           true,
			expectedHaveBackendService: true,
			expectedServiceName:        "foo-namespace.foo-svc.80",
			expectedServiceHost:        "foo-svc.foo-namespace.80.svc",
		},
		{
			name: "multiple ingresses with default backend",
			ingresses: []netv1.Ingress{
				someIngress(now.Add(time.Second), "newer"),
				someIngress(now, "older"),
			},
			expressionRoutes:           false,
			expectedHaveBackendService: true,
			expectedServiceName:        "foo-namespace.older.80",
			expectedServiceHost:        "older.foo-namespace.80.svc",
		},
		{
			name: "multiple ingresses with default backend and expression routes enabled",
			ingresses: []netv1.Ingress{
				someIngress(now.Add(time.Second), "newer"),
				someIngress(now, "older"),
			},
			expressionRoutes:           true,
			expectedHaveBackendService: true,
			expectedServiceName:        "foo-namespace.older.80",
			expectedServiceHost:        "older.foo-namespace.80.svc",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc, ok := getDefaultBackendService(tc.ingresses, tc.expressionRoutes)
			require.Equal(t, tc.expectedHaveBackendService, ok)
			if tc.expectedHaveBackendService {
				require.Equal(t, tc.expectedServiceName, *svc.Name)
				require.Equal(t, tc.expectedServiceHost, *svc.Host)
				require.Len(t, svc.Routes, 1)
				route := svc.Routes[0]
				if tc.expressionRoutes {
					require.Equal(t, `(http.path ^= "/") && ((net.protocol == "http") || (net.protocol == "https"))`, *route.Expression)
					require.Equal(t, subtranslator.IngressDefaultBackendPriority, *route.Priority)
				} else {
					require.Len(t, route.Paths, 1)
					require.Equal(t, *route.Paths[0], "/")
				}
			}
		})
	}
}

func TestRewriteURIAnnotation(t *testing.T) {
	someIngress := func(name, rewriteURI string) netv1.Ingress {
		return netv1.Ingress{
			TypeMeta: metav1.TypeMeta{Kind: "Ingress"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey:                              annotations.DefaultIngressClass,
					annotations.AnnotationPrefix + annotations.RewriteURIKey: rewriteURI,
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
										Path:     "/~/api/(.*)",
										PathType: lo.ToPtr(netv1.PathTypePrefix),
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: name,
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
		}
	}

	validRewriteURIIngress := someIngress("valid_annotation_svc", "/rewrite/$1/xx")
	invalidRewriteURIIngress := someIngress("invalid_annotation_svc", "/rewrite/$/xx")
	t.Run("Ingress rule with rewrite annotation enabled", func(t *testing.T) {
		s, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				&invalidRewriteURIIngress,
				&validRewriteURIIngress,
			},
		})
		require.NoError(t, err)

		p := mustNewTranslator(t, s)
		p.featureFlags.RewriteURIs = true

		rules := p.ingressRulesFromIngressV1().ServiceNameToServices
		require.Len(t, rules, 1)

		for _, svc := range rules {
			for _, route := range svc.Routes {
				require.Equal(t, []kong.Plugin{
					{
						Name: kong.String("request-transformer"),
						Config: kong.Configuration{
							"replace": map[string]string{
								"uri": "/rewrite/$(uri_captures[1])/xx",
							},
						},
					},
				}, route.Plugins)
			}
		}

		errs := p.failuresCollector.PopResourceFailures()
		require.Len(t, errs, 1)
		require.Equal(t, "unexpected / at pos 10", errs[0].Message())
	})

	t.Run("Ingress rule with rewrite annotation disabled", func(t *testing.T) {
		emptyRewriteURIIngress := someIngress("empty_annotation_svc", "")
		delete(emptyRewriteURIIngress.ObjectMeta.Annotations, annotations.AnnotationPrefix+annotations.RewriteURIKey)

		s, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{
				&invalidRewriteURIIngress,
				&validRewriteURIIngress,
				&emptyRewriteURIIngress,
			},
		})
		require.NoError(t, err)

		p := mustNewTranslator(t, s)

		rules := p.ingressRulesFromIngressV1().ServiceNameToServices
		require.Len(t, rules, 1)

		for _, svc := range rules {
			require.Nil(t, svc.Plugins)
		}

		errs := p.failuresCollector.PopResourceFailures()
		require.Len(t, errs, 2)

		for _, err := range errs {
			require.Equal(t, "konghq.com/rewrite annotation not supported when rewrite uris disabled", err.Message())
		}
	})
}
