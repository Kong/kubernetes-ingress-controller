package parser

import (
	"fmt"
	"testing"
	"time"

	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser/translators"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
)

func TestFromIngressV1(t *testing.T) {
	t.Run("no ingress returns empty info", func(t *testing.T) {
		s, err := store.NewFakeStore(store.FakeObjects{
			IngressesV1: []*netv1.Ingress{},
		})
		require.NoError(t, err)
		p := mustNewParser(t, s)

		parsedInfo := p.ingressRulesFromIngressV1()
		assert.Equal(t, ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			ServiceNameToParent:   make(map[string]client.Object),
			SecretNameToSNIs:      newSecretNameToSNIs(),
		}, parsedInfo)
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
		p := mustNewParser(t, s)

		assert.NotPanics(t, func() {
			p.ingressRulesFromIngressV1()
		})
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
					require.Equal(t, translators.IngressDefaultBackendPriority, *route.Priority)
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

		p := mustNewParser(t, s)
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

		p := mustNewParser(t, s)

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

func TestIngressGeneratingKongRoutesWithConflictingNames(t *testing.T) {
	createIngress := func(noRules, noPaths int) *netv1.Ingress {
		var ingresRules []netv1.IngressRule
		for i := 0; i <= noRules; i++ {
			var paths []netv1.HTTPIngressPath
			for j := 0; j <= noPaths; j++ {
				paths = append(paths, netv1.HTTPIngressPath{
					Path:     "/~/api/(.*)",
					PathType: lo.ToPtr(netv1.PathTypePrefix),
					Backend: netv1.IngressBackend{
						Service: &netv1.IngressServiceBackend{
							Name: "foo-svc",
							Port: netv1.ServiceBackendPort{Number: 80},
						},
					},
				},
				)
			}
			ingresRules = append(ingresRules, netv1.IngressRule{
				Host: "example.com",
				IngressRuleValue: netv1.IngressRuleValue{
					HTTP: &netv1.HTTPIngressRuleValue{
						Paths: paths,
					},
				},
			})
		}
		return &netv1.Ingress{
			TypeMeta: metav1.TypeMeta{Kind: "Ingress"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("conflicting-one-%d-%d", noRules, noPaths),
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: netv1.IngressSpec{
				Rules: ingresRules,
			},
		}
	}
	s, err := store.NewFakeStore(store.FakeObjects{
		IngressesV1: []*netv1.Ingress{
			createIngress(1, 1),
			createIngress(7, 11),
			createIngress(11, 5),
			// The below generates Kong Routes that are conflicting with each other.
			createIngress(11, 11),
			createIngress(22, 23),
		},
	})
	require.NoError(t, err)
	p := mustNewParser(t, s)

	t.Run("more than 11 rules and paths in Ingress and CombinedRoutes=false warns about conflicts", func(t *testing.T) {
		p.featureFlags.CombinedServiceRoutes = false
		_ = p.ingressRulesFromIngressV1()
		errs := p.failuresCollector.PopResourceFailures()
		require.Len(t, errs, 30)
		errMsgs := lo.Map(errs, func(err failures.ResourceFailure, _ int) string {
			return err.Message()
		})

		// Just check couple of them, other follows the same pattern.
		const expectedErrMsg = "ERROR: Kong route with conflicting name: %s use feature gate CombinedRoutes=true or update Kong Kubernetes Ingress Controller version to 3.0.0 or above (both remediation changes naming schema of Kong routes)"
		for _, conflictingName := range []string{
			"foo-namespace.conflicting-one-11-11.110",
			"foo-namespace.conflicting-one-11-11.111",
			"foo-namespace.conflicting-one-22-23.114",
			"foo-namespace.conflicting-one-22-23.222",
		} {
			require.Contains(t, errMsgs, fmt.Sprintf(expectedErrMsg, conflictingName))
		}
	})

	t.Run("more than 11 rules and paths in Ingress and CombinedRoutes=true gives no conflicts and no warnings", func(t *testing.T) {
		p.featureFlags.CombinedServiceRoutes = true
		_ = p.ingressRulesFromIngressV1()
		errs := p.failuresCollector.PopResourceFailures()
		require.Len(t, errs, 0)
	})
}
