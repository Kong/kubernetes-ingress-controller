package subtranslator

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-kong/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
	incubatorv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/incubator/v1alpha1"
)

func TestTranslateIngressATC(t *testing.T) {
	testCases := []struct {
		name               string
		ingress            *netv1.Ingress
		kongServiceFacades []*incubatorv1alpha1.KongServiceFacade
		expectedServices   map[string]kongstate.Service
	}{
		{
			name: "a basic ingress resource with a single rule and prefix path type",
			ingress: &netv1.Ingress{
				TypeMeta: ingressTypeMeta,
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
				"default.test-service.80": {
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-service.80"),
						Host:           kong.String("test-service.default.80.svc"),
						ConnectTimeout: defaultServiceTimeoutInKongFormat(),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    defaultServiceTimeoutInKongFormat(),
						WriteTimeout:   defaultServiceTimeoutInKongFormat(),
					},
					Routes: []kongstate.Route{{
						Ingress: util.K8sObjectInfo{
							Name:             "test-ingress",
							Namespace:        corev1.NamespaceDefault,
							GroupVersionKind: ingressGVK,
						},
						Route: kong.Route{
							Name:       kong.String("default.test-ingress.test-service.konghq.com.80"),
							Expression: kong.String(`(http.host == "konghq.com") && ((http.path == "/api") || (http.path ^= "/api/"))`),
							Priority: kong.Uint64(IngressRoutePriorityTraits{
								MatchFields:   2,
								PlainHostOnly: true,
								MaxPathLength: 5,
								HasRegexPath:  true,
							}.EncodeToPriority()),
							PreserveHost:      kong.Bool(true),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
						ExpressionRoutes: true,
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "a basic ingress resource with a single rule, and only one path results in a single kong service and route",
			ingress: &netv1.Ingress{
				TypeMeta: ingressTypeMeta,
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
				"default.test-service.80": {
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-service.80"),
						Host:           kong.String("test-service.default.80.svc"),
						ConnectTimeout: defaultServiceTimeoutInKongFormat(),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    defaultServiceTimeoutInKongFormat(),
						WriteTimeout:   defaultServiceTimeoutInKongFormat(),
					},
					Routes: []kongstate.Route{{
						Ingress: util.K8sObjectInfo{
							Name:             "test-ingress",
							Namespace:        corev1.NamespaceDefault,
							GroupVersionKind: ingressGVK,
						},
						Route: kong.Route{
							Name:       kong.String("default.test-ingress.test-service.konghq.com.80"),
							Expression: kong.String(`(http.host == "konghq.com") && (http.path ^= "/api/")`),
							Priority: kong.Uint64(IngressRoutePriorityTraits{
								MatchFields:   2,
								PlainHostOnly: true,
								MaxPathLength: 5,
								HasRegexPath:  false,
							}.EncodeToPriority()),
							PreserveHost:      kong.Bool(true),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
						ExpressionRoutes: true,
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "an ingress with method, protocol, and header annotations",
			ingress: &netv1.Ingress{
				TypeMeta: ingressTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress-annotations",
					Namespace: corev1.NamespaceDefault,
					Annotations: map[string]string{
						"konghq.com/methods":     "GET",
						"konghq.com/protocols":   "http",
						"konghq.com/headers.foo": "bar",
					},
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
				"default.test-service.80": {
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-service.80"),
						Host:           kong.String("test-service.default.80.svc"),
						ConnectTimeout: defaultServiceTimeoutInKongFormat(),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    defaultServiceTimeoutInKongFormat(),
						WriteTimeout:   defaultServiceTimeoutInKongFormat(),
					},
					Routes: []kongstate.Route{{
						Ingress: util.K8sObjectInfo{
							Name:      "test-ingress-annotations",
							Namespace: corev1.NamespaceDefault,
							Annotations: map[string]string{
								"konghq.com/methods":     "GET",
								"konghq.com/protocols":   "http",
								"konghq.com/headers.foo": "bar",
							},
							GroupVersionKind: ingressGVK,
						},
						Route: kong.Route{
							Name:       kong.String("default.test-ingress-annotations.test-service.konghq.com.80"),
							Expression: kong.String(`(http.host == "konghq.com") && (http.path ^= "/api/") && (http.headers.foo == "bar") && (http.method == "GET")`),
							Priority: kong.Uint64(IngressRoutePriorityTraits{
								MatchFields:   4,
								PlainHostOnly: true,
								MaxPathLength: 5,
								HasRegexPath:  false,
								HeaderCount:   1,
							}.EncodeToPriority()),
							PreserveHost:      kong.Bool(true),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress-annotations", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
						ExpressionRoutes: true,
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).MustBuild(),
					},
					Parent: &netv1.Ingress{
						TypeMeta: metav1.TypeMeta{Kind: "Ingress", APIVersion: netv1.SchemeGroupVersion.String()},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-ingress-annotations",
							Namespace: corev1.NamespaceDefault,
							Annotations: map[string]string{
								"konghq.com/methods":     "GET",
								"konghq.com/protocols":   "http",
								"konghq.com/headers.foo": "bar",
							},
						},
					},
				},
			},
		},
		{
			name: "KongServiceFacade used as a backend",
			ingress: &netv1.Ingress{
				TypeMeta: ingressTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "default",
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{{
									Path: "/api/",
									Backend: netv1.IngressBackend{
										Resource: &corev1.TypedLocalObjectReference{
											APIGroup: lo.ToPtr(incubatorv1alpha1.GroupVersion.Group),
											Kind:     incubatorv1alpha1.KongServiceFacadeKind,
											Name:     "svc-facade",
										},
									},
								}},
							},
						},
					}},
				},
			},
			kongServiceFacades: []*incubatorv1alpha1.KongServiceFacade{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-facade",
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						Kind:       incubatorv1alpha1.KongServiceFacadeKind,
						APIVersion: incubatorv1alpha1.GroupVersion.String(),
					},
					Spec: incubatorv1alpha1.KongServiceFacadeSpec{
						Backend: incubatorv1alpha1.KongServiceFacadeBackend{
							Name: "svc",
							Port: 8080,
						},
					},
				},
			},
			expectedServices: map[string]kongstate.Service{
				"default.svc-facade.svc.facade": {
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.svc-facade.svc.facade"),
						Host:           kong.String("default.svc-facade.svc.facade"),
						ConnectTimeout: defaultServiceTimeoutInKongFormat(),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    defaultServiceTimeoutInKongFormat(),
						WriteTimeout:   defaultServiceTimeoutInKongFormat(),
					},
					Routes: []kongstate.Route{{
						Ingress: util.K8sObjectInfo{
							Name:             "test-ingress",
							Namespace:        corev1.NamespaceDefault,
							GroupVersionKind: ingressGVK,
						},
						Route: kong.Route{
							Name:       kong.String("default.test-ingress.konghq.com.svc-facade.svc.facade"),
							Expression: kong.String(`(http.host == "konghq.com") && (http.path ^= "/api/")`),
							Priority: kong.Uint64(IngressRoutePriorityTraits{
								MatchFields:   2,
								PlainHostOnly: true,
								MaxPathLength: 5,
								HasRegexPath:  false,
							}.EncodeToPriority()),
							PreserveHost:      kong.Bool(true),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
						ExpressionRoutes: true,
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("svc-facade").
							WithType(kongstate.ServiceBackendTypeKongServiceFacade).
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(8080).
							MustBuild(),
					},
					Parent: &incubatorv1alpha1.KongServiceFacade{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "svc-facade",
							Namespace: corev1.NamespaceDefault,
						},
						TypeMeta: metav1.TypeMeta{
							Kind:       incubatorv1alpha1.KongServiceFacadeKind,
							APIVersion: incubatorv1alpha1.GroupVersion.String(),
						},
					},
				},
			},
		},
		{
			name: "not existing KongServiceFacade used as a backend",
			ingress: &netv1.Ingress{
				TypeMeta: ingressTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: "default",
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{{
									Path: "/api/",
									Backend: netv1.IngressBackend{
										Resource: &corev1.TypedLocalObjectReference{
											APIGroup: lo.ToPtr(incubatorv1alpha1.GroupVersion.Group),
											Kind:     incubatorv1alpha1.KongServiceFacadeKind,
											Name:     "svc-facade",
										},
									},
								}},
							},
						},
					}},
				},
			},
			expectedServices: map[string]kongstate.Service{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			failuresCollector := failures.NewResourceFailuresCollector(logr.Discard())
			storer := lo.Must(store.NewFakeStore(store.FakeObjects{
				KongServiceFacades: tc.kongServiceFacades,
			}))
			services := TranslateIngresses(
				[]*netv1.Ingress{tc.ingress},
				kongv1alpha1.IngressClassParametersSpec{},
				TranslateIngressFeatureFlags{
					ExpressionRoutes:  true,
					KongServiceFacade: true,
				},
				noopObjectsCollector{},
				failuresCollector,
				storer,
			)
			checkOnlyIngressMeta := cmp.Transformer("checkOnlyIngressMeta", func(i *netv1.Ingress) *netv1.Ingress {
				// In the result we only care about ingresses' metadata being equal.
				// We ignore specification to simplify tests.
				return &netv1.Ingress{
					ObjectMeta: i.ObjectMeta,
				}
			})
			checkOnlyKongServiceFacadeMeta := cmp.Transformer("checkOnlyKongServiceFacadeMeta", func(i *incubatorv1alpha1.KongServiceFacade) *incubatorv1alpha1.KongServiceFacade {
				// In the result we only care about KongServiceFacades' metadata being equal.
				// We ignore specification to simplify tests.
				return &incubatorv1alpha1.KongServiceFacade{
					ObjectMeta: i.ObjectMeta,
				}
			})
			compareServiceBackend := cmp.AllowUnexported(kongstate.ServiceBackend{})
			diff := cmp.Diff(tc.expectedServices, services, checkOnlyIngressMeta, checkOnlyKongServiceFacadeMeta, compareServiceBackend)
			require.Empty(t, diff, "expected no difference between expected and translated ingress")
		})
	}
}

func TestCalculateIngressRoutePriorityTraits(t *testing.T) {
	testCases := []struct {
		name               string
		paths              []netv1.HTTPIngressPath
		regexPathPrefix    string
		ingressHost        string
		ingressAnnotations map[string]string
		expectedTraits     IngressRoutePriorityTraits
	}{
		{
			name: "single prefix path with no hosts",
			paths: []netv1.HTTPIngressPath{
				{
					Path:     "/foo/",
					PathType: lo.ToPtr(netv1.PathTypePrefix),
				},
			},
			expectedTraits: IngressRoutePriorityTraits{
				MatchFields:   1,
				MaxPathLength: len("/foo/"),
				HasRegexPath:  true,
			},
		},
		{
			name: "multiple paths with one host",
			paths: []netv1.HTTPIngressPath{
				{
					Path:     "/foo/",
					PathType: lo.ToPtr(netv1.PathTypePrefix),
				},
				{
					Path:     "/foobar/",
					PathType: lo.ToPtr(netv1.PathTypeExact),
				},
			},
			ingressHost: "example.com",
			expectedTraits: IngressRoutePriorityTraits{
				MatchFields:   2,
				PlainHostOnly: true,
				MaxPathLength: len("/foobar/"),
				HasRegexPath:  true,
			},
		},
		{
			name: "multiple paths with headers and hosts",
			paths: []netv1.HTTPIngressPath{
				{
					Path:     "/foo/",
					PathType: lo.ToPtr(netv1.PathTypePrefix),
				},
				{
					Path:     "/foobar/",
					PathType: lo.ToPtr(netv1.PathTypeExact),
				},
			},
			ingressHost: "example.com",
			ingressAnnotations: map[string]string{
				"konghq.com/host-aliases": "*.example.com",
				"konghq.com/headers.key1": "value1",
				"konghq.com/headers.key2": "value2",
			},
			expectedTraits: IngressRoutePriorityTraits{
				MatchFields:   3,
				PlainHostOnly: false,
				HeaderCount:   2,
				MaxPathLength: len("/foobar/"),
				HasRegexPath:  true,
			},
		},
		{
			name: "ImplementationSpecific path with regex",
			paths: []netv1.HTTPIngressPath{
				{
					Path:     "/~/[a-z0-9]{3}/",
					PathType: lo.ToPtr(netv1.PathTypeImplementationSpecific),
				},
			},
			regexPathPrefix: "/~",
			expectedTraits: IngressRoutePriorityTraits{
				MatchFields:   1,
				PlainHostOnly: false,
				MaxPathLength: len("/[a-z0-9]{3}/"),
				HasRegexPath:  true,
			},
		},
		{
			name: "ImplementationSpecific path without regex",
			paths: []netv1.HTTPIngressPath{
				{
					Path:     "/abc/def",
					PathType: lo.ToPtr(netv1.PathTypeImplementationSpecific),
				},
			},
			regexPathPrefix: "/~",
			expectedTraits: IngressRoutePriorityTraits{
				MatchFields:   1,
				PlainHostOnly: false,
				MaxPathLength: len("/abc/def"),
				HasRegexPath:  false,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			traits := calculateIngressRoutePriorityTraits(
				tc.paths, tc.regexPathPrefix, tc.ingressHost, tc.ingressAnnotations,
			)
			require.Equal(t, tc.expectedTraits, traits)
		})
	}
}

func TestEncodeIngressRoutePriorityFromTraits(t *testing.T) {
	testCases := []struct {
		name             string
		traits           IngressRoutePriorityTraits
		expectedPriority RoutePriorityType
	}{
		{
			name: "plain host true regex path false",
			traits: IngressRoutePriorityTraits{
				MatchFields:   2,
				PlainHostOnly: true,
				MaxPathLength: 5,
				HasRegexPath:  false,
			},
			expectedPriority: (3 << 44) | (2 << 41) | (1 << 32) | 5,
		},
		{
			name: "plain host false regex path true",
			traits: IngressRoutePriorityTraits{
				MatchFields:   2,
				PlainHostOnly: false,
				HeaderCount:   2,
				MaxPathLength: 5,
				HasRegexPath:  true,
			},
			expectedPriority: (3 << 44) | (2 << 41) | (2 << 33) | (1 << 16) | 5,
		},
		{
			name: "header number exceed limit",
			traits: IngressRoutePriorityTraits{
				MatchFields:   2,
				PlainHostOnly: false,
				HeaderCount:   256,
				MaxPathLength: 5,
				HasRegexPath:  true,
			},
			expectedPriority: (3 << 44) | (2 << 41) | (255 << 33) | (1 << 16) | 5,
		},
		{
			name: "path length exceed limit",
			traits: IngressRoutePriorityTraits{
				MatchFields:   2,
				PlainHostOnly: true,
				MaxPathLength: 100000,
			},
			expectedPriority: (3 << 44) | (2 << 41) | (1 << 32) | 65535,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			priority := tc.traits.EncodeToPriority()
			require.Equal(t, tc.expectedPriority, priority)
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
