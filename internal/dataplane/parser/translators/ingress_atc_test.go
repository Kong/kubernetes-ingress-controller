package translators

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

func TestTranslateIngressATC(t *testing.T) {
	testCases := []struct {
		name             string
		ingress          *netv1.Ingress
		expectedServices []*kongstate.Service
	}{
		{
			name: "a basic ingress resource with a single rule and prefix path type",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{{
									Path:     "/api",
									PathType: &pathTypePrefix,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: "test-service",
											Port: netv1.ServiceBackendPort{
												Name:   "http",
												Number: 80,
											},
										},
									},
								}},
							},
						},
					}},
				},
			},
			expectedServices: []*kongstate.Service{{
				Namespace: corev1.NamespaceDefault,
				Service: kong.Service{
					Name:           kong.String("default.test-ingress.test-service.80"),
					Host:           kong.String("test-service.default.80.svc"),
					ConnectTimeout: kong.Int(int(defaultServiceTimeout.Milliseconds())),
					Path:           kong.String("/"),
					Port:           kong.Int(80),
					Protocol:       kong.String("http"),
					Retries:        kong.Int(defaultRetries),
					ReadTimeout:    kong.Int(int(defaultServiceTimeout.Milliseconds())),
					WriteTimeout:   kong.Int(int(defaultServiceTimeout.Milliseconds())),
				},
				Routes: []kongstate.Route{{
					Ingress: util.K8sObjectInfo{
						Name:      "test-ingress",
						Namespace: corev1.NamespaceDefault,
					},
					Route: kong.Route{
						Name:              kong.String("default.test-ingress.test-service.konghq.com.80"),
						Expression:        kong.String(`(http.host == "konghq.com") && ((http.path == "/api") || (http.path ^= "/api/")) && ((net.protocol == "http") || (net.protocol == "https"))`),
						Priority:          kong.Int(NormalIngressExpressionPriority),
						PreserveHost:      kong.Bool(true),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
						Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default"),
					},
					ExpressionRoutes: true,
				}},
				Backends: []kongstate.ServiceBackend{{
					Name:      "test-service",
					Namespace: corev1.NamespaceDefault,
					PortDef: kongstate.PortDef{
						Mode:   kongstate.PortModeByNumber,
						Number: 80,
					},
				}},
				Parent: expectedParentIngress(),
			}},
		},
		{
			name: "a basic ingress resource with a single rule, and only one path results in a single kong service and route",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{{
									Path: "/api/",
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: "test-service",
											Port: netv1.ServiceBackendPort{
												Name:   "http",
												Number: 80,
											},
										},
									},
								}},
							},
						},
					}},
				},
			},
			expectedServices: []*kongstate.Service{{
				Namespace: corev1.NamespaceDefault,
				Service: kong.Service{
					Name:           kong.String("default.test-ingress.test-service.80"),
					Host:           kong.String("test-service.default.80.svc"),
					ConnectTimeout: kong.Int(int(defaultServiceTimeout.Milliseconds())),
					Path:           kong.String("/"),
					Port:           kong.Int(80),
					Protocol:       kong.String("http"),
					Retries:        kong.Int(defaultRetries),
					ReadTimeout:    kong.Int(int(defaultServiceTimeout.Milliseconds())),
					WriteTimeout:   kong.Int(int(defaultServiceTimeout.Milliseconds())),
				},
				Routes: []kongstate.Route{{
					Ingress: util.K8sObjectInfo{
						Name:      "test-ingress",
						Namespace: corev1.NamespaceDefault,
					},
					Route: kong.Route{
						Name:              kong.String("default.test-ingress.test-service.konghq.com.80"),
						Expression:        kong.String(`(http.host == "konghq.com") && (http.path ^= "/api/") && ((net.protocol == "http") || (net.protocol == "https"))`),
						Priority:          kong.Int(NormalIngressExpressionPriority),
						PreserveHost:      kong.Bool(true),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
						Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default"),
					},
					ExpressionRoutes: true,
				}},
				Backends: []kongstate.ServiceBackend{{
					Name:      "test-service",
					Namespace: corev1.NamespaceDefault,
					PortDef: kongstate.PortDef{
						Mode:   kongstate.PortModeByNumber,
						Number: 80,
					},
				}},
				Parent: expectedParentIngress(),
			}},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			services := TranslateIngress(tc.ingress, false, true)
			checkOnlyObjectMeta := cmp.Transformer("checkOnlyObjectMeta", func(i *netv1.Ingress) *netv1.Ingress {
				// In the result we only care about ingresses' metadata being equal.
				// We ignore specification to simplify tests.
				return &netv1.Ingress{
					ObjectMeta: i.ObjectMeta,
				}
			})
			diff := cmp.Diff(tc.expectedServices, services, checkOnlyObjectMeta)
			require.Empty(t, diff, "expected no difference between expected and translated ingress")
		})
	}
}

func TestPathMatcherFromIngressPath(t *testing.T) {
	testCases := []struct {
		name        string
		path        netv1.HTTPIngressPath
		regexPrefix string
		expression  string
	}{
		{
			name: "simple prefix match",
			path: netv1.HTTPIngressPath{
				Path:     "/v1/api",
				PathType: &pathTypePrefix,
			},
			expression: `(http.path == "/v1/api") || (http.path ^= "/v1/api/")`,
		},
		{
			name: "simple exact match",
			path: netv1.HTTPIngressPath{
				Path:     "/v1/api",
				PathType: &pathTypeExact,
			},
			expression: `http.path == "/v1/api"`,
		},
		{
			name: "regex match",
			path: netv1.HTTPIngressPath{
				Path:     "/~/[a-z]+",
				PathType: &pathTypeImplementationSpecific,
			},
			expression: `http.path ~ "^/[a-z]+"`,
		},
		{
			name: "empty path",
			path: netv1.HTTPIngressPath{
				Path:     "",
				PathType: &pathTypePrefix,
			},
			expression: `http.path ^= "/"`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			regexPrefix := tc.regexPrefix
			if regexPrefix == "" {
				regexPrefix = ControllerPathRegexPrefix
			}
			matcher := pathMatcherFromIngressPath(tc.path, regexPrefix)
			require.Equal(t, tc.expression, matcher.Expression())
		})
	}
}

func TestHeaderMatcherFromHeaders(t *testing.T) {
	testCases := []struct {
		name       string
		headers    map[string][]string
		expression string
	}{
		{
			name: "single header with single value",
			headers: map[string][]string{
				"X-Key1": {"value1"},
			},
			expression: `http.headers.x_key1 == "value1"`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			matcher := headerMatcherFromHeaders(tc.headers)
			require.Equal(t, tc.expression, matcher.Expression())
		})
	}
}
