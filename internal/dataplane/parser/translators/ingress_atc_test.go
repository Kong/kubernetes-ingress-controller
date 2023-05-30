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
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1alpha1"
)

func TestTranslateIngressATC(t *testing.T) {
	testCases := []struct {
		name             string
		ingress          *netv1.Ingress
		expectedServices map[string]kongstate.Service
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
			expectedServices: map[string]kongstate.Service{
				"default.test-ingress.test-service.80": {
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
							Priority:          kong.Int((2 << 41) + (1 << 32) + (1 << 16) + 5),
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
				},
			},
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
			expectedServices: map[string]kongstate.Service{
				"default.test-ingress.test-service.80": {
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
							Priority:          kong.Int((2 << 41) + (1 << 32) + (5)),
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
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			services := TranslateIngresses(
				[]*netv1.Ingress{tc.ingress},
				v1alpha1.IngressClassParametersSpec{},
				TranslateIngressFeatureFlags{
					ExpressionRoutes: true,
					RegexPathPrefix:  false,
					CombinedServices: false,
				},
				noopObjectsCollector{},
			)
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
			name: "regex match with initial ^",
			path: netv1.HTTPIngressPath{
				Path:     "/~^/foo/[a-z]+",
				PathType: &pathTypeImplementationSpecific,
			},
			expression: `http.path ~ "^/foo/[a-z]+"`,
		},
		{
			name: "empty prefix path",
			path: netv1.HTTPIngressPath{
				Path:     "",
				PathType: &pathTypePrefix,
			},
			expression: `http.path ^= "/"`,
		},
		{
			name: "empty exact match",
			path: netv1.HTTPIngressPath{
				Path:     "",
				PathType: &pathTypeExact,
			},
			expression: `http.path == "/"`,
		},
		{
			name: "empty regex match",
			path: netv1.HTTPIngressPath{
				Path:     "/~",
				PathType: &pathTypeImplementationSpecific,
			},
			expression: `http.path ~ "^/"`,
		},
		{
			name: "empty implementation specific (non-regex) match",
			path: netv1.HTTPIngressPath{
				Path:     "",
				PathType: &pathTypeImplementationSpecific,
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
		{
			name: "header 'Host' is skipped and multiple headers",
			headers: map[string][]string{
				"Host":   {"konghq.com"},
				"X-Key1": {"value1"},
				"X-Key2": {"value2"},
			},
			expression: `(http.headers.x_key1 == "value1") && (http.headers.x_key2 == "value2")`,
		},
		{
			name: "single header with multiple values",
			headers: map[string][]string{
				"X-Key1": {"value1", "value2"},
			},
			expression: `(http.headers.x_key1 == "value1") || (http.headers.x_key1 == "value2")`,
		},
		{
			name: "single header with regex value",
			headers: map[string][]string{
				"X-Key1": {"~*[a-z]+"},
			},
			expression: `http.headers.x_key1 ~ "[a-z]+"`,
		},
		{
			name: "empty value",
			headers: map[string][]string{
				"X-Key1": nil,
				"X-Key2": {},
			},
			expression: "",
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

func TestProtocolMatcherFromProtocols(t *testing.T) {
	testCases := []struct {
		name       string
		protocols  []string
		expression string
	}{
		{
			name:       "single protocol",
			protocols:  []string{"https"},
			expression: `net.protocol == "https"`,
		},
		{
			name:       "multiple protocols",
			protocols:  []string{"http", "https"},
			expression: `(net.protocol == "http") || (net.protocol == "https")`,
		},
		{
			name:       "multiple protocols including invalid protocol",
			protocols:  []string{"http", "ppp"},
			expression: `net.protocol == "http"`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		matcher := protocolMatcherFromProtocols(tc.protocols)
		require.Equal(t, tc.expression, matcher.Expression())
	}
}

func TestMethodMatcherFromMethods(t *testing.T) {
	testCases := []struct {
		name       string
		methods    []string
		expression string
	}{
		{
			name:       "single method",
			methods:    []string{"GET"},
			expression: `http.method == "GET"`,
		},
		{
			name:       "multiple methods",
			methods:    []string{"POST", "PUT"},
			expression: `(http.method == "POST") || (http.method == "PUT")`,
		},
		{
			name:       "multiple methods with invalid method",
			methods:    []string{"HEAD", "OPTIONS", "paTch"},
			expression: `(http.method == "HEAD") || (http.method == "OPTIONS")`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		matcher := methodMatcherFromMethods(tc.methods)
		require.Equal(t, tc.expression, matcher.Expression())
	}
}

func TestSNIMatcherFromSNIs(t *testing.T) {
	testCases := []struct {
		name       string
		snis       []string
		expression string
	}{
		{
			name:       "single SNI",
			snis:       []string{"konghq.com"},
			expression: `tls.sni == "konghq.com"`,
		},
		{
			name:       "multiple SNIs",
			snis:       []string{"docs.konghq.com", "apis.konghq.com"},
			expression: `(tls.sni == "docs.konghq.com") || (tls.sni == "apis.konghq.com")`,
		},
		{
			name:       "multiple SNIs with wildcard SNI, which should be omitted",
			snis:       []string{"foo.com", "*.bar.com"},
			expression: `tls.sni == "foo.com"`,
		},
		{
			name:       "multiple SNIs with invalid SNI",
			snis:       []string{"foo.com", "a..bar.com"},
			expression: `tls.sni == "foo.com"`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		matcher := sniMatcherFromSNIs(tc.snis)
		require.Equal(t, tc.expression, matcher.Expression())
	}
}
