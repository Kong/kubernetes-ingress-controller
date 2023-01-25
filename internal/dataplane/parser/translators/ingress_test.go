package translators

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var (
	pathTypeExact                  = netv1.PathTypeExact
	pathTypeImplementationSpecific = netv1.PathTypeImplementationSpecific
	pathTypePrefix                 = netv1.PathTypePrefix
)

func expectedParentIngress() *netv1.Ingress {
	return &netv1.Ingress{
		TypeMeta:   metav1.TypeMeta{Kind: "Ingress", APIVersion: netv1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{Name: "test-ingress", Namespace: corev1.NamespaceDefault},
	}
}

func TestTranslateIngress(t *testing.T) {
	tts := []struct {
		name           string
		ingress        *netv1.Ingress
		addRegexPrefix bool
		expected       []*kongstate.Service
	}{
		{
			name:           "a basic ingress resource with a single rule and prefix path type",
			addRegexPrefix: true,
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
			expected: []*kongstate.Service{{
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
						Hosts:             kong.StringSlice("konghq.com"),
						Paths:             kong.StringSlice("/api/", "~/api$"),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			expected: []*kongstate.Service{{
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
						Hosts:             kong.StringSlice("konghq.com"),
						Paths:             kong.StringSlice("/api/"), // default ImplementationSpecific
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			expected: []*kongstate.Service{{
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
						Hosts:             kong.StringSlice("konghq.com"),
						Paths:             kong.StringSlice("/api/", "/api$"),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			name: "an ingress with path type exact gets a kong route with an exact path match",
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
			},
			expected: []*kongstate.Service{{
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
						Hosts:             kong.StringSlice("konghq.com"),
						Paths:             kong.StringSlice("/api$"), // No Prefix Pathing
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			name: "an Ingress resource with implementation specific path type doesn't modify the path",
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
			},
			expected: []*kongstate.Service{{
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
						Hosts:             kong.StringSlice("konghq.com"),
						Paths:             kong.StringSlice("/api"), // No path mods
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			name: "an Ingress resource with paths with double /'s gets flattened",
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
			},
			expected: []*kongstate.Service{{
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
						Hosts:             kong.StringSlice("konghq.com"),
						Paths:             kong.StringSlice("/v1/api/"),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			name: "empty paths get treated as '/'",
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
			},
			expected: []*kongstate.Service{{
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
						Hosts:             kong.StringSlice("konghq.com"),
						Paths:             kong.StringSlice("/"),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			name: "multiple and various paths get compiled together properly",
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
			},
			expected: []*kongstate.Service{{
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
						Name:  kong.String("default.test-ingress.test-service.konghq.com.80"),
						Hosts: kong.StringSlice("konghq.com"),
						Paths: kong.StringSlice(
							"/v1/api",
							"/v2/api",
							"/v3/api",
							"/other/path/1$",
							"/other/path/2",
							"/",
						),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			name: "when no host is provided, all hosts are matched",
			ingress: &netv1.Ingress{
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
			},
			expected: []*kongstate.Service{{
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
						Name: kong.String("default.test-ingress.test-service..80"),
						Paths: kong.StringSlice(
							"/v1/api",
							"/v2/api",
							"/v3/api",
							"/other/path/1$",
							"/other/path/2",
							"/",
						),
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			name: "when there are multiple backends services, paths wont be combined and separate kong services will be provided",
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
			},
			expected: []*kongstate.Service{
				{
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-ingress.test-service1.80"),
						Host:           kong.String("test-service1.default.80.svc"),
						ConnectTimeout: kong.Int(int(defaultServiceTimeout.Milliseconds())),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    kong.Int(int(defaultServiceTimeout.Milliseconds())),
						WriteTimeout:   kong.Int(int(defaultServiceTimeout.Milliseconds())),
					},
					Routes: []kongstate.Route{
						{
							Ingress: util.K8sObjectInfo{
								Name:      "test-ingress",
								Namespace: corev1.NamespaceDefault,
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
							},
						},
					},
					Backends: []kongstate.ServiceBackend{{
						Name:      "test-service1",
						Namespace: corev1.NamespaceDefault,
						PortDef: kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: 80,
						},
					}},
					Parent: expectedParentIngress(),
				},
				{
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-ingress.test-service2.80"),
						Host:           kong.String("test-service2.default.80.svc"),
						ConnectTimeout: kong.Int(int(defaultServiceTimeout.Milliseconds())),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    kong.Int(int(defaultServiceTimeout.Milliseconds())),
						WriteTimeout:   kong.Int(int(defaultServiceTimeout.Milliseconds())),
					},
					Routes: []kongstate.Route{
						{
							Ingress: util.K8sObjectInfo{
								Name:      "test-ingress",
								Namespace: corev1.NamespaceDefault,
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
							},
						},
					},
					Backends: []kongstate.ServiceBackend{{
						Name:      "test-service2",
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
			name: "when there are multiple ingress rules with overlapping host and service, separate kong services will be provided",
			ingress: &netv1.Ingress{
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
			},
			expected: []*kongstate.Service{
				{
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-ingress.ad-service.80"),
						Host:           kong.String("ad-service.default.80.svc"),
						ConnectTimeout: kong.Int(int(defaultServiceTimeout.Milliseconds())),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    kong.Int(int(defaultServiceTimeout.Milliseconds())),
						WriteTimeout:   kong.Int(int(defaultServiceTimeout.Milliseconds())),
					},
					Routes: []kongstate.Route{
						{
							Ingress: util.K8sObjectInfo{
								Name:      "test-ingress",
								Namespace: corev1.NamespaceDefault,
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
							},
						},
					},
					Backends: []kongstate.ServiceBackend{{
						Name:      "ad-service",
						Namespace: corev1.NamespaceDefault,
						PortDef: kongstate.PortDef{
							Mode:   kongstate.PortModeByNumber,
							Number: 80,
						},
					}},
					Parent: expectedParentIngress(),
				},
				{
					Namespace: corev1.NamespaceDefault,
					Service: kong.Service{
						Name:           kong.String("default.test-ingress.mad-service.80"),
						Host:           kong.String("mad-service.default.80.svc"),
						ConnectTimeout: kong.Int(int(defaultServiceTimeout.Milliseconds())),
						Path:           kong.String("/"),
						Port:           kong.Int(80),
						Protocol:       kong.String("http"),
						Retries:        kong.Int(defaultRetries),
						ReadTimeout:    kong.Int(int(defaultServiceTimeout.Milliseconds())),
						WriteTimeout:   kong.Int(int(defaultServiceTimeout.Milliseconds())),
					},
					Routes: []kongstate.Route{
						{
							Ingress: util.K8sObjectInfo{
								Name:      "test-ingress",
								Namespace: corev1.NamespaceDefault,
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
							},
						},
					},
					Backends: []kongstate.ServiceBackend{{
						Name:      "mad-service",
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
			name: "* in host is replaced to _",
			ingress: &netv1.Ingress{
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
			},
			expected: []*kongstate.Service{{
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
						Name:              kong.String("default.test-ingress.test-service._.konghq.com.80"),
						Hosts:             kong.StringSlice("*.konghq.com"),
						Paths:             kong.StringSlice("/api/"), // default ImplementationSpecific
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
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
			name: "use port name when service port number is not provided",
			ingress: &netv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "networking.k8s.io/v1",
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
			},
			expected: []*kongstate.Service{{
				Namespace: corev1.NamespaceDefault,
				Service: kong.Service{
					Name:           kong.String("default.test-ingress.test-service.http"),
					Host:           kong.String("test-service.default.http.svc"),
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
						Name:              kong.String("default.test-ingress.test-service.konghq.com.http"),
						Hosts:             kong.StringSlice("konghq.com"),
						Paths:             kong.StringSlice("/api/"), // default ImplementationSpecific
						PreserveHost:      kong.Bool(true),
						Protocols:         kong.StringSlice("http", "https"),
						RegexPriority:     kong.Int(0),
						StripPath:         kong.Bool(false),
						ResponseBuffering: kong.Bool(true),
						RequestBuffering:  kong.Bool(true),
					},
				}},
				Backends: []kongstate.ServiceBackend{{
					Name:      "test-service",
					Namespace: corev1.NamespaceDefault,
					PortDef: kongstate.PortDef{
						Mode: kongstate.PortModeByName,
						Name: "http",
					},
				}},
				Parent: expectedParentIngress(),
			}},
		},
	}

	for _, tt := range tts {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			checkOnlyObjectMeta := cmp.Transformer("checkOnlyObjectMeta", func(i *netv1.Ingress) *netv1.Ingress {
				// In the result we only care about ingresses' metadata being equal.
				// We ignore specification to simplify tests.
				return &netv1.Ingress{
					ObjectMeta: i.ObjectMeta,
				}
			})
			diff := cmp.Diff(tt.expected, TranslateIngress(tt.ingress, tt.addRegexPrefix), checkOnlyObjectMeta)
			require.Empty(t, diff, "expected no difference between expected and translated ingress")
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
				"/v1/api/packages$",
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
				"/v1/api/packages$",
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
			out: kong.StringSlice("/v1/api/packages$"),
		},
		{
			name: "path type exact will cause an exact matching path on a regular path with a / suffix",
			in: netv1.HTTPIngressPath{
				Path:     "/v1/api/packages/",
				PathType: &pathTypeExact,
			},
			out: kong.StringSlice("/v1/api/packages/$"),
		},
		{
			name: "path type exact will supply a default if no path is provided",
			in: netv1.HTTPIngressPath{
				Path:     "",
				PathType: &pathTypeExact,
			},
			out: kong.StringSlice("/$"),
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
			// TODO split test cases to handle regex
			assert.Equal(t, tt.out, PathsFromIngressPaths(tt.in, false))
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

func TestPathsFromIngressPathsRegexPrefix(t *testing.T) {
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
			assert.Equal(t, tt.out, PathsFromIngressPaths(tt.in, true))
		})
	}
}
