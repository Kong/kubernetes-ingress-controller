package subtranslator

import (
	"errors"
	"sort"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-kong/kong"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	incubatorv1alpha1 "github.com/kong/kubernetes-configuration/api/incubator/v1alpha1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/failures"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
)

var (
	pathTypeExact                  = netv1.PathTypeExact
	pathTypeImplementationSpecific = netv1.PathTypeImplementationSpecific
	pathTypePrefix                 = netv1.PathTypePrefix
	ingressTypeMeta                = metav1.TypeMeta{
		APIVersion: netv1.SchemeGroupVersion.Group + "/" + netv1.SchemeGroupVersion.Version,
		Kind:       "Ingress",
	}
	ingressGVK = schema.GroupVersionKind{
		Group:   netv1.SchemeGroupVersion.Group,
		Version: netv1.SchemeGroupVersion.Version,
		Kind:    "Ingress",
	}
)

func expectedParentIngress() *netv1.Ingress {
	return &netv1.Ingress{
		TypeMeta:   metav1.TypeMeta{Kind: "Ingress", APIVersion: netv1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{Name: "test-ingress", Namespace: corev1.NamespaceDefault},
	}
}

type noopObjectsCollector struct{}

func (noopObjectsCollector) Add(client.Object) {}

func TestTranslateIngress(t *testing.T) {
	tts := []struct {
		name               string
		ingresses          []*netv1.Ingress
		expected           map[string]kongstate.Service
		kongServiceFacades []*incubatorv1alpha1.KongServiceFacade
	}{
		{
			name: "a basic ingress resource with a single rule and prefix path type",
			ingresses: []*netv1.Ingress{{
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
			}},
			expected: map[string]kongstate.Service{
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
							Name:              kong.String("default.test-ingress.test-service.konghq.com.80"),
							Hosts:             kong.StringSlice("konghq.com"),
							Paths:             kong.StringSlice("/api/", "~/api$"),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "a basic ingress resource with a single rule, and only one path results in a single kong service and route",
			ingresses: []*netv1.Ingress{{
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
			}},
			expected: map[string]kongstate.Service{
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
							Name:              kong.String("default.test-ingress.test-service.konghq.com.80"),
							Hosts:             kong.StringSlice("konghq.com"),
							Paths:             kong.StringSlice("/api/"), // default ImplementationSpecific
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "a basic ingress resource with a single rule and prefix path type",
			ingresses: []*netv1.Ingress{{
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
			}},
			expected: map[string]kongstate.Service{
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
							Name:              kong.String("default.test-ingress.test-service.konghq.com.80"),
							Hosts:             kong.StringSlice("konghq.com"),
							Paths:             kong.StringSlice("/api/", "~/api$"),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "an ingress with path type exact gets a kong route with an exact path match",
			ingresses: []*netv1.Ingress{{
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
									PathType: &pathTypeExact,
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
			}},
			expected: map[string]kongstate.Service{
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
							Name:              kong.String("default.test-ingress.test-service.konghq.com.80"),
							Hosts:             kong.StringSlice("konghq.com"),
							Paths:             kong.StringSlice("~/api$"),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "an Ingress resource with implementation specific path type doesn't modify the path",
			ingresses: []*netv1.Ingress{{
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
									PathType: &pathTypeImplementationSpecific,
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
			}},
			expected: map[string]kongstate.Service{
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
							Name:              kong.String("default.test-ingress.test-service.konghq.com.80"),
							Hosts:             kong.StringSlice("konghq.com"),
							Paths:             kong.StringSlice("/api"), // No path mods
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "an Ingress resource with paths with double /'s gets flattened",
			ingresses: []*netv1.Ingress{{
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
									Path: "/v1//api///",
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
			}},
			expected: map[string]kongstate.Service{
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
							Name:              kong.String("default.test-ingress.test-service.konghq.com.80"),
							Hosts:             kong.StringSlice("konghq.com"),
							Paths:             kong.StringSlice("/v1/api/"),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "empty paths get treated as '/'",
			ingresses: []*netv1.Ingress{{
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
									Path: "",
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
			}},
			expected: map[string]kongstate.Service{
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
							Name:              kong.String("default.test-ingress.test-service.konghq.com.80"),
							Hosts:             kong.StringSlice("konghq.com"),
							Paths:             kong.StringSlice("/"),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "multiple and various paths get compiled together properly",
			ingresses: []*netv1.Ingress{{
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
								Paths: []netv1.HTTPIngressPath{
									{
										Path: "/v1/api",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v2/api",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v3/api",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "/other/path/1",
										PathType: &pathTypeExact,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "/other/path/2",
										PathType: &pathTypeImplementationSpecific,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "",
										PathType: &pathTypeImplementationSpecific,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
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
					}},
				},
			}},
			expected: map[string]kongstate.Service{
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
							Name:  kong.String("default.test-ingress.test-service.konghq.com.80"),
							Hosts: kong.StringSlice("konghq.com"),
							Paths: kong.StringSlice(
								"/v1/api",
								"/v2/api",
								"/v3/api",
								"~/other/path/1$",
								"/other/path/2",
								"/",
							),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "when no host is provided, all hosts are matched",
			ingresses: []*netv1.Ingress{{
				TypeMeta: ingressTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{{
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path: "/v1/api",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v2/api",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v3/api",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "/other/path/1",
										PathType: &pathTypeExact,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "/other/path/2",
										PathType: &pathTypeImplementationSpecific,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "",
										PathType: &pathTypeImplementationSpecific,
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service",
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
					}},
				},
			}},
			expected: map[string]kongstate.Service{
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
							Name: kong.String("default.test-ingress.test-service..80"),
							Paths: kong.StringSlice(
								"/v1/api",
								"/v2/api",
								"/v3/api",
								"~/other/path/1$",
								"/other/path/2",
								"/",
							),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "when there are multiple backends services, paths wont be combined and separate kong services will be provided",
			ingresses: []*netv1.Ingress{{
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
								Paths: []netv1.HTTPIngressPath{
									{
										Path: "/v1/api",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service1",
												Port: netv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v2/api",
										Backend: netv1.IngressBackend{
											Service: &netv1.IngressServiceBackend{
												Name: "test-service2",
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
					}},
				},
			}},
			expected: map[string]kongstate.Service{
				"default.test-service1.80": {
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-service1.80"),
						Host:           kong.String("test-service1.default.80.svc"),
						ConnectTimeout: defaultServiceTimeoutInKongFormat(),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    defaultServiceTimeoutInKongFormat(),
						WriteTimeout:   defaultServiceTimeoutInKongFormat(),
					},
					Routes: []kongstate.Route{
						{
							Ingress: util.K8sObjectInfo{
								Name:             "test-ingress",
								Namespace:        corev1.NamespaceDefault,
								GroupVersionKind: ingressGVK,
							},
							Route: kong.Route{
								Name:              kong.String("default.test-ingress.test-service1.konghq.com.80"),
								Hosts:             kong.StringSlice("konghq.com"),
								Paths:             kong.StringSlice("/v1/api"),
								PreserveHost:      kong.Bool(true),
								Protocols:         kong.StringSlice("http", "https"),
								RegexPriority:     kong.Int(0),
								StripPath:         kong.Bool(false),
								ResponseBuffering: kong.Bool(true),
								RequestBuffering:  kong.Bool(true),
								Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
							},
						},
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service1").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
				"default.test-service2.80": {
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-service2.80"),
						Host:           kong.String("test-service2.default.80.svc"),
						ConnectTimeout: defaultServiceTimeoutInKongFormat(),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    defaultServiceTimeoutInKongFormat(),
						WriteTimeout:   defaultServiceTimeoutInKongFormat(),
					},
					Routes: []kongstate.Route{
						{
							Ingress: util.K8sObjectInfo{
								Name:             "test-ingress",
								Namespace:        corev1.NamespaceDefault,
								GroupVersionKind: ingressGVK,
							},
							Route: kong.Route{
								Name:              kong.String("default.test-ingress.test-service2.konghq.com.80"),
								Hosts:             kong.StringSlice("konghq.com"),
								Paths:             kong.StringSlice("/v2/api"),
								PreserveHost:      kong.Bool(true),
								Protocols:         kong.StringSlice("http", "https"),
								RegexPriority:     kong.Int(0),
								StripPath:         kong.Bool(false),
								ResponseBuffering: kong.Bool(true),
								RequestBuffering:  kong.Bool(true),
								Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
							},
						},
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service2").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "when there are multiple ingress rules with overlapping host and service, separate kong services will be provided",
			ingresses: []*netv1.Ingress{{
				TypeMeta: ingressTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{
						{
							Host: "konghq.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/v1/api",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "ad-service",
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
						{
							Host: "konghq.co",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{
										{
											Path: "/v1/api",
											Backend: netv1.IngressBackend{
												Service: &netv1.IngressServiceBackend{
													Name: "mad-service",
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
			}},
			expected: map[string]kongstate.Service{
				"default.ad-service.80": {
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.ad-service.80"),
						Host:           kong.String("ad-service.default.80.svc"),
						ConnectTimeout: defaultServiceTimeoutInKongFormat(),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    defaultServiceTimeoutInKongFormat(),
						WriteTimeout:   defaultServiceTimeoutInKongFormat(),
					},
					Routes: []kongstate.Route{
						{
							Ingress: util.K8sObjectInfo{
								Name:             "test-ingress",
								Namespace:        corev1.NamespaceDefault,
								GroupVersionKind: ingressGVK,
							},
							Route: kong.Route{
								Name:              kong.String("default.test-ingress.ad-service.konghq.com.80"),
								Hosts:             kong.StringSlice("konghq.com"),
								Paths:             kong.StringSlice("/v1/api"),
								PreserveHost:      kong.Bool(true),
								Protocols:         kong.StringSlice("http", "https"),
								RegexPriority:     kong.Int(0),
								StripPath:         kong.Bool(false),
								ResponseBuffering: kong.Bool(true),
								RequestBuffering:  kong.Bool(true),
								Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
							},
						},
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("ad-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
				"default.mad-service.80": {
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.mad-service.80"),
						Host:           kong.String("mad-service.default.80.svc"),
						ConnectTimeout: defaultServiceTimeoutInKongFormat(),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    defaultServiceTimeoutInKongFormat(),
						WriteTimeout:   defaultServiceTimeoutInKongFormat(),
					},
					Routes: []kongstate.Route{
						{
							Ingress: util.K8sObjectInfo{
								Name:             "test-ingress",
								Namespace:        corev1.NamespaceDefault,
								GroupVersionKind: ingressGVK,
							},
							Route: kong.Route{
								Name:              kong.String("default.test-ingress.mad-service.konghq.co.80"),
								Hosts:             kong.StringSlice("konghq.co"),
								Paths:             kong.StringSlice("/v1/api"),
								PreserveHost:      kong.Bool(true),
								Protocols:         kong.StringSlice("http", "https"),
								RegexPriority:     kong.Int(0),
								StripPath:         kong.Bool(false),
								ResponseBuffering: kong.Bool(true),
								RequestBuffering:  kong.Bool(true),
								Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
							},
						},
					},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("mad-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "* in host is replaced to _",
			ingresses: []*netv1.Ingress{{
				TypeMeta: ingressTypeMeta,
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: netv1.IngressSpec{
					Rules: []netv1.IngressRule{{
						Host: "*.konghq.com",
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
			}},
			expected: map[string]kongstate.Service{
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
							Name:              kong.String("default.test-ingress.test-service._.konghq.com.80"),
							Hosts:             kong.StringSlice("*.konghq.com"),
							Paths:             kong.StringSlice("/api/"), // default ImplementationSpecific
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortNumber(80).
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "use port name when service port number is not provided",
			ingresses: []*netv1.Ingress{{
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
												Name: "http",
											},
										},
									},
								}},
							},
						},
					}},
				},
			}},
			expected: map[string]kongstate.Service{
				"default.test-service.http": {
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-service.http"),
						Host:           kong.String("test-service.default.http.svc"),
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
							Name:              kong.String("default.test-ingress.test-service.konghq.com.http"),
							Hosts:             kong.StringSlice("konghq.com"),
							Paths:             kong.StringSlice("/api/"), // default ImplementationSpecific
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags: kong.StringSlice(
								"k8s-name:test-ingress",
								"k8s-namespace:default",
								"k8s-kind:Ingress",
								"k8s-group:networking.k8s.io",
								"k8s-version:v1",
							),
						},
					}},
					Backends: []kongstate.ServiceBackend{
						builder.NewKongstateServiceBackend("test-service").
							WithNamespace(corev1.NamespaceDefault).
							WithPortName("http").
							MustBuild(),
					},
					Parent: expectedParentIngress(),
				},
			},
		},
		{
			name: "KongServiceFacade used as a backend",
			ingresses: []*netv1.Ingress{{
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
			}},
			kongServiceFacades: []*incubatorv1alpha1.KongServiceFacade{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-facade",
						Namespace: "default",
					},
					Spec: incubatorv1alpha1.KongServiceFacadeSpec{
						Backend: incubatorv1alpha1.KongServiceFacadeBackend{
							Name: "svc",
							Port: 8080,
						},
					},
				},
			},
			expected: map[string]kongstate.Service{
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
							Name:              kong.String("default.test-ingress.konghq.com.svc-facade.svc.facade"),
							Hosts:             kong.StringSlice("konghq.com"),
							Paths:             kong.StringSlice("/api/"),
							PreserveHost:      kong.Bool(true),
							Protocols:         kong.StringSlice("http", "https"),
							RegexPriority:     kong.Int(0),
							StripPath:         kong.Bool(false),
							ResponseBuffering: kong.Bool(true),
							RequestBuffering:  kong.Bool(true),
							Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
						},
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
							Namespace: "default",
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
			ingresses: []*netv1.Ingress{{
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
			}},
			expected: map[string]kongstate.Service{},
		},
		{
			name: "KongServiceFacade used in multiple Ingresses",
			ingresses: []*netv1.Ingress{
				{
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
										Path: "/ingress-1/",
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
				{
					TypeMeta: ingressTypeMeta,
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-ingress-2",
						Namespace: "default",
					},
					Spec: netv1.IngressSpec{
						Rules: []netv1.IngressRule{{
							Host: "konghq.com",
							IngressRuleValue: netv1.IngressRuleValue{
								HTTP: &netv1.HTTPIngressRuleValue{
									Paths: []netv1.HTTPIngressPath{{
										Path: "/ingress-2/",
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
			},
			kongServiceFacades: []*incubatorv1alpha1.KongServiceFacade{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "svc-facade",
						Namespace: "default",
					},
					Spec: incubatorv1alpha1.KongServiceFacadeSpec{
						Backend: incubatorv1alpha1.KongServiceFacadeBackend{
							Name: "svc",
							Port: 8080,
						},
					},
				},
			},
			expected: map[string]kongstate.Service{
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
					Routes: []kongstate.Route{
						{
							Ingress: util.K8sObjectInfo{
								Name:             "test-ingress",
								Namespace:        corev1.NamespaceDefault,
								GroupVersionKind: ingressGVK,
							},
							Route: kong.Route{
								Name:              kong.String("default.test-ingress.konghq.com.svc-facade.svc.facade"),
								Hosts:             kong.StringSlice("konghq.com"),
								Paths:             kong.StringSlice("/ingress-1/"),
								PreserveHost:      kong.Bool(true),
								Protocols:         kong.StringSlice("http", "https"),
								RegexPriority:     kong.Int(0),
								StripPath:         kong.Bool(false),
								ResponseBuffering: kong.Bool(true),
								RequestBuffering:  kong.Bool(true),
								Tags:              kong.StringSlice("k8s-name:test-ingress", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
							},
						},
						{
							Ingress: util.K8sObjectInfo{
								Name:             "test-ingress-2",
								Namespace:        corev1.NamespaceDefault,
								GroupVersionKind: ingressGVK,
							},
							Route: kong.Route{
								Name:              kong.String("default.test-ingress-2.konghq.com.svc-facade.svc.facade"),
								Hosts:             kong.StringSlice("konghq.com"),
								Paths:             kong.StringSlice("/ingress-2/"),
								PreserveHost:      kong.Bool(true),
								Protocols:         kong.StringSlice("http", "https"),
								RegexPriority:     kong.Int(0),
								StripPath:         kong.Bool(false),
								ResponseBuffering: kong.Bool(true),
								RequestBuffering:  kong.Bool(true),
								Tags:              kong.StringSlice("k8s-name:test-ingress-2", "k8s-namespace:default", "k8s-kind:Ingress", "k8s-group:networking.k8s.io", "k8s-version:v1"),
							},
						},
					},
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
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
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

			failuresCollector := failures.NewResourceFailuresCollector(logr.Discard())
			storer := lo.Must(store.NewFakeStore(store.FakeObjects{
				KongServiceFacades: tt.kongServiceFacades,
			}))
			translatedServices := TranslateIngresses(
				tt.ingresses,
				kongv1alpha1.IngressClassParametersSpec{},
				TranslateIngressFeatureFlags{
					ExpressionRoutes:  false,
					KongServiceFacade: true,
				},
				noopObjectsCollector{},
				failuresCollector,
				storer,
			)

			// Sort routes to make the test deterministic. Not doing this in the code itself as the deterministic
			// order is not required on this level of translation and that would be an unnecessary performance hit.
			for _, service := range translatedServices {
				sort.Slice(service.Routes, func(i, j int) bool {
					return *service.Routes[i].Route.Name > *service.Routes[j].Route.Name
				})
			}

			compareServiceBackend := cmp.AllowUnexported(kongstate.ServiceBackend{})
			diff := cmp.Diff(tt.expected, translatedServices, checkOnlyIngressMeta, checkOnlyKongServiceFacadeMeta, compareServiceBackend)
			require.Empty(t, diff, "expected no difference between expected and translated ingress")
		})
	}
}

func TestTranslateIngress_KongServiceFacadeFailures(t *testing.T) {
	testCases := []struct {
		name                   string
		ingress                *netv1.Ingress
		storerObjects          store.FakeObjects
		serviceFacadeFeatureOn bool
		expectedFailures       []string
	}{
		{
			name: "KongServiceFacade used as backend with no feature flag on",
			ingress: builder.NewIngress("ingress", "kong").
				WithNamespace("default").
				WithRules(netv1.IngressRule{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{{
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
				}).Build(),
			serviceFacadeFeatureOn: false,
			expectedFailures:       []string{`failed to get backend for ingress path "/": KongServiceFacade is not enabled, please set the "KongServiceFacade" feature gate to 'true' to enable it`},
		},
		{
			name: "KongServiceFacade used as backend but not existing",
			ingress: builder.NewIngress("ingress", "kong").
				WithNamespace("default").
				WithRules(netv1.IngressRule{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{{
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
				}).Build(),
			serviceFacadeFeatureOn: true,
			expectedFailures:       []string{`failed to get backend for ingress path "/": failed to get KongServiceFacade "svc-facade": KongServiceFacade default/svc-facade not found`},
		},
		{
			name: "wrong API group used for KongServiceFacade",
			ingress: builder.NewIngress("ingress", "kong").
				WithNamespace("default").
				WithRules(netv1.IngressRule{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{{
								Backend: netv1.IngressBackend{
									Resource: &corev1.TypedLocalObjectReference{
										APIGroup: lo.ToPtr("wrong-group"),
										Kind:     incubatorv1alpha1.KongServiceFacadeKind,
										Name:     "svc-facade",
									},
								},
							}},
						},
					},
				}).Build(),
			serviceFacadeFeatureOn: true,
			expectedFailures:       []string{`failed to get backend for ingress path "/": unknown resource type wrong-group/KongServiceFacade`},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			failuresCollector := failures.NewResourceFailuresCollector(logr.Discard())
			storer := lo.Must(store.NewFakeStore(tc.storerObjects))
			result := TranslateIngresses(
				[]*netv1.Ingress{tc.ingress},
				kongv1alpha1.IngressClassParametersSpec{},
				TranslateIngressFeatureFlags{
					KongServiceFacade: tc.serviceFacadeFeatureOn,
				},
				noopObjectsCollector{},
				failuresCollector,
				storer,
			)
			require.Empty(t, result)

			collectedFailures := failuresCollector.PopResourceFailures()
			require.Len(t, collectedFailures, len(tc.expectedFailures))
			for _, failure := range collectedFailures {
				require.Contains(t, tc.expectedFailures, failure.Message())
			}
		})
	}
}

func TestFlattenMultipleSlashes(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "a path with no slashes gets left alone",
			in:   "api",
			out:  "api",
		},
		{
			name: "a normal path gets left alone",
			in:   "/v1/api/packages/",
			out:  "/v1/api/packages/",
		},
		{
			name: "a path with two starting slashes gets flattened",
			in:   "//",
			out:  "/",
		},
		{
			name: "a path with many slashes gets flattened",
			in:   "/////////////////",
			out:  "/",
		},
		{
			name: "a path with multiple groups of double slashes gets flattened",
			in:   "//api//packages//",
			out:  "/api/packages/",
		},
		{
			name: "a path with multiple groups of various sized groups of slashes gets flattened",
			in:   "////////v1////api//packages///and/stuff//////////////////",
			out:  "/v1/api/packages/and/stuff/",
		},
		{
			name: "a path with multiple slashes but none at the end gets flattened",
			in:   "////////v1////api//packages///and/stuff",
			out:  "/v1/api/packages/and/stuff",
		},
		{
			name: "a path with multiple slashes but none at the beginning gets flattened",
			in:   "v1/////api//packages///and/stuff//////",
			out:  "v1/api/packages/and/stuff/",
		},
		{
			name: "a path with multiple slashes but none at the beginning or end gets flattened",
			in:   "v1/////api//packages//////////////////////////and/stuff",
			out:  "v1/api/packages/and/stuff",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, flattenMultipleSlashes(tt.in))
		})
	}
}

func TestPathsFromIngressPaths(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   netv1.HTTPIngressPath
		out  []*string
	}{
		{
			name: "path type prefix will expand the match to a trailing slash if not provided",
			in: netv1.HTTPIngressPath{
				Path:     "/v1/api/packages",
				PathType: &pathTypePrefix,
			},
			out: kong.StringSlice(
				"/v1/api/packages/",
				"~/v1/api/packages$",
			),
		},
		{
			name: "path type prefix will expand the match with a literal match if a slash is provided",
			in: netv1.HTTPIngressPath{
				Path:     "/v1/api/packages/",
				PathType: &pathTypePrefix,
			},
			out: kong.StringSlice(
				"/v1/api/packages/",
				"~/v1/api/packages$",
			),
		},
		{
			name: "path type prefix will provide a default when no path is provided",
			in: netv1.HTTPIngressPath{
				Path:     "",
				PathType: &pathTypePrefix,
			},
			out: kong.StringSlice("/"),
		},
		{
			name: "path type exact will cause an exact matching path on a regular path",
			in: netv1.HTTPIngressPath{
				Path:     "/v1/api/packages",
				PathType: &pathTypeExact,
			},
			out: kong.StringSlice("~/v1/api/packages$"),
		},
		{
			name: "path type exact will cause an exact matching path on a regular path with a / suffix",
			in: netv1.HTTPIngressPath{
				Path:     "/v1/api/packages/",
				PathType: &pathTypeExact,
			},
			out: kong.StringSlice("~/v1/api/packages/$"),
		},
		{
			name: "path type exact will supply a default if no path is provided",
			in: netv1.HTTPIngressPath{
				Path:     "",
				PathType: &pathTypeExact,
			},
			out: kong.StringSlice("~/$"),
		},
		{
			name: "path type implementation-specific will leave the path alone",
			in: netv1.HTTPIngressPath{
				Path:     "/asdfasd9jhf09432$",
				PathType: &pathTypeImplementationSpecific,
			},
			out: kong.StringSlice("/asdfasd9jhf09432$"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, PathsFromIngressPaths(tt.in))
		})
	}
}

func TestMaybePrependRegexPrefix(t *testing.T) {
	testCases := []struct {
		name                 string
		path                 string
		controllerPrefix     string
		applyLegacyHeuristic bool
		expected             string
	}{
		{
			name:             "default controller regex prefix - prefixed path",
			path:             "/~/v1/api/packages",
			controllerPrefix: ControllerPathRegexPrefix,
			expected:         "~/v1/api/packages",
		},
		{
			name:             "default controller regex prefix - not prefixed path",
			path:             "/v1/api/packages",
			controllerPrefix: ControllerPathRegexPrefix,
			expected:         "/v1/api/packages",
		},
		{
			name:             "custom controller regex prefix - prefixed path",
			path:             "##/v1/api/packages",
			controllerPrefix: "##",
			expected:         "~/v1/api/packages",
		},
		{
			name:             "custom controller regex prefix - not prefixed path",
			path:             "/v1/api/packages",
			controllerPrefix: "##",
			expected:         "/v1/api/packages",
		},
		{
			name:                 "default controller regex prefix - path not prefixed, but legacy heuristic is applied",
			path:                 "/v1/api/resource/\\d+/",
			controllerPrefix:     ControllerPathRegexPrefix,
			applyLegacyHeuristic: true,
			expected:             "~/v1/api/resource/\\d+/",
		},
		{
			name:                 "default controller regex prefix - path not prefixed, no legacy heuristic is applied",
			path:                 "/v1/api/resource/\\d+/",
			controllerPrefix:     ControllerPathRegexPrefix,
			applyLegacyHeuristic: false,
			expected:             "/v1/api/resource/\\d+/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MaybePrependRegexPrefix(tc.path, tc.controllerPrefix, tc.applyLegacyHeuristic)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestMaybePrependRegexPrefixForIngressV1Fn(t *testing.T) {
	// Let it be const as the heuristic logic is already tested in TestMaybePrependRegexPrefix.
	const applyLegacyHeuristic = true

	t.Run("ingress with a custom regex prefix generates fn with a custom prefix", func(t *testing.T) {
		ingress := &netv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					annotations.AnnotationPrefix + annotations.RegexPrefixKey: "##",
				},
			},
		}

		generatedMaybePrependRegexPrefixFn := MaybePrependRegexPrefixForIngressV1Fn(ingress, applyLegacyHeuristic)
		result := generatedMaybePrependRegexPrefixFn("##/v1/api/packages")
		require.Equal(t, "~/v1/api/packages", *result)
	})

	t.Run("ingress with no custom regex prefix generates fn with default prefix", func(t *testing.T) {
		generatedMaybePrependRegexPrefixFn := MaybePrependRegexPrefixForIngressV1Fn(&netv1.Ingress{}, applyLegacyHeuristic)
		result := generatedMaybePrependRegexPrefixFn("/~/v1/api/packages")
		require.Equal(t, "~/v1/api/packages", *result)
	})
}

func TestGenerateRewriteURIConfig(t *testing.T) {
	testCases := []struct {
		name          string
		uri           string
		expectedError error
		expectedURI   string
	}{
		{
			name:        "no capture group",
			uri:         "/bar/xx/yy",
			expectedURI: "/bar/xx/yy",
		},
		{
			name:        "valid single digit capture group",
			uri:         "/bar/$1xx/yy",
			expectedURI: "/bar/$(uri_captures[1])xx/yy",
		},
		{
			name:        "valid multiple digits capture group",
			uri:         "/bar/$12xx/yy",
			expectedURI: "/bar/$(uri_captures[12])xx/yy",
		},
		{
			name:        "valid multiple capture groups",
			uri:         "/bar/$1xx/$12yy",
			expectedURI: "/bar/$(uri_captures[1])xx/$(uri_captures[12])yy",
		},
		{
			name:        "valid multiple capture groups (end with capture group)",
			uri:         "/bar/$1xx/$12",
			expectedURI: "/bar/$(uri_captures[1])xx/$(uri_captures[12])",
		},
		{
			name:        "valid multiple capture groups (adjacent capture groups)",
			uri:         "/bar/$11$12",
			expectedURI: "/bar/$(uri_captures[11])$(uri_captures[12])",
		},
		{
			name:          "left brace following $",
			uri:           "/bar/${}11",
			expectedError: errors.New("unexpected { at pos 6"),
			expectedURI:   "",
		},

		{
			name:        "digits without $",
			uri:         "/bar/123$12",
			expectedURI: "/bar/123$(uri_captures[12])",
		},
		{
			name:          "$ at end",
			uri:           "/bar/xxxx$",
			expectedError: errors.New("unexpected end of string"),
			expectedURI:   "",
		},
		{
			name:        "escaped $",
			uri:         "/bar/xxxx\\$12",
			expectedURI: "/bar/xxxx$12",
		},
		{
			name:        "escaped $ at end",
			uri:         "/bar/xxxx\\$",
			expectedURI: "/bar/xxxx$",
		},
		{
			name:        "escaped $ after capture group",
			uri:         "/bar/xxxx$13\\$x",
			expectedURI: "/bar/xxxx$(uri_captures[13])$x",
		},
		{
			name:        "mixed with escaped and unescaped $",
			uri:         "/bar/xxxx\\$12/fd$33",
			expectedURI: "/bar/xxxx$12/fd$(uri_captures[33])",
		},
		{
			name:          "non $ after \\",
			uri:           "/bar/xxxx\\n",
			expectedError: errors.New("unexpected n at pos 10"),
			expectedURI:   "",
		},
		{
			name:          "\\ at end",
			uri:           "/bar/xxxx\\",
			expectedError: errors.New("unexpected end of string"),
			expectedURI:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uri, err := generateRewriteURIConfig(tc.uri)
			require.Equal(t, tc.expectedError, err)
			require.Equal(t, tc.expectedURI, uri)
		})
	}
}

func TestMaybeRewriteURI(t *testing.T) {
	testCases := []struct {
		name            string
		service         kongstate.Service
		expectedError   error
		expectedPlugins []kong.Plugin
	}{
		{
			name: "konghq.com/rewrite annotation is not exist",
			service: kongstate.Service{
				Parent: &netv1.Ingress{},
			},
			expectedError:   nil,
			expectedPlugins: nil,
		},
		{
			name: "konghq.com/rewrite annotation is empty",
			service: kongstate.Service{
				Parent: &netv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.RewriteURIKey: "",
						},
					},
				},
			},
			expectedError: nil,
			expectedPlugins: []kong.Plugin{
				{
					Name: kong.String("request-transformer"),
					Config: kong.Configuration{
						"replace": map[string]string{
							"uri": "/",
						},
					},
				},
			},
		},
		{
			name: "valid rewrite uri",
			service: kongstate.Service{
				Parent: &netv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							annotations.AnnotationPrefix + annotations.RewriteURIKey: "/xxx$11yy/",
						},
					},
				},
			},
			expectedError: nil,
			expectedPlugins: []kong.Plugin{
				{
					Name: kong.String("request-transformer"),
					Config: kong.Configuration{
						"replace": map[string]string{
							"uri": "/xxx$(uri_captures[11])yy/",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := MaybeRewriteURI(&tc.service, true)
			require.Equal(t, tc.expectedError, err)
			for _, route := range tc.service.Routes {
				require.Equal(t, tc.expectedPlugins, route.Plugins)
			}
		})
	}
}
