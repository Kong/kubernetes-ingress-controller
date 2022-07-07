package translators

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

var (
	pathTypeExact                  = networkingv1.PathTypeExact
	pathTypeImplementationSpecific = networkingv1.PathTypeImplementationSpecific
	pathTypePrefix                 = networkingv1.PathTypePrefix
)

func TestTranslateIngress(t *testing.T) {
	tts := []struct {
		name     string
		ingress  *networkingv1.Ingress
		expected []*kongstate.Service
	}{
		{
			name: "a basic ingress resource with a single rule, and only one path results in a single kong service and route",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{{
									Path: "/api",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "test-service",
											Port: networkingv1.ServiceBackendPort{
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
						Paths:             kong.StringSlice("/api$", "/api/"), // Prefix pathing is the default behavior when no pathtype is defined
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
			}},
		},

		{
			name: "an ingress with path type exact gets a kong route with an exact path match",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{{
									Path:     "/api",
									PathType: &pathTypeExact,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "test-service",
											Port: networkingv1.ServiceBackendPort{
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
			}},
		},

		{
			name: "an Ingress resource with implementation specific path type doesn't modify the path",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{{
									Path:     "/api",
									PathType: &pathTypeImplementationSpecific,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "test-service",
											Port: networkingv1.ServiceBackendPort{
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
			}},
		},

		{
			name: "an Ingress resource with paths with double /'s gets flattened",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{{
									Path: "/v1//api///",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "test-service",
											Port: networkingv1.ServiceBackendPort{
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
						Paths:             kong.StringSlice("/v1/api$", "/v1/api/"),
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
			}},
		},

		{
			name: "empty paths get treated as '/'",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{{
									Path: "",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "test-service",
											Port: networkingv1.ServiceBackendPort{
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
			}},
		},

		{
			name: "multiple and various paths get compiled together properly",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path: "/v1/api",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v2/api",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v3/api",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "/other/path/1",
										PathType: &pathTypeExact,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "/other/path/2",
										PathType: &pathTypeImplementationSpecific,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "",
										PathType: &pathTypeImplementationSpecific,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
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
							"/v1/api$", "/v1/api/",
							"/v2/api$", "/v2/api/",
							"/v3/api$", "/v3/api/",
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
			}},
		},

		{
			name: "when no host is provided, all hosts are matched",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{{
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path: "/v1/api",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v2/api",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v3/api",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "/other/path/1",
										PathType: &pathTypeExact,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "/other/path/2",
										PathType: &pathTypeImplementationSpecific,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path:     "",
										PathType: &pathTypeImplementationSpecific,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
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
							"/v1/api$", "/v1/api/",
							"/v2/api$", "/v2/api/",
							"/v3/api$", "/v3/api/",
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
			}},
		},

		{
			name: "when there are multiple backends services, paths wont be combined and separate kong services will be provided",
			ingress: &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ingress",
					Namespace: corev1.NamespaceDefault,
				},
				Spec: networkingv1.IngressSpec{
					Rules: []networkingv1.IngressRule{{
						Host: "konghq.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path: "/v1/api",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service1",
												Port: networkingv1.ServiceBackendPort{
													Name:   "http",
													Number: 80,
												},
											},
										},
									},
									{
										Path: "/v2/api",
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service2",
												Port: networkingv1.ServiceBackendPort{
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
								Paths:             kong.StringSlice("/v1/api$", "/v1/api/"),
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
								Paths:             kong.StringSlice("/v2/api$", "/v2/api/"),
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
				},
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, TranslateIngress(tt.ingress), tt.expected)
		})
	}
}

func Test_pathsFromIngressPaths(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   networkingv1.HTTPIngressPath
		out  []*string
	}{
		{
			name: "path type prefix will expand the match to a trailing slash if not provided",
			in: networkingv1.HTTPIngressPath{
				Path:     "/v1/api/packages",
				PathType: &pathTypePrefix,
			},
			out: kong.StringSlice(
				"/v1/api/packages$",
				"/v1/api/packages/",
			),
		},
		{
			name: "path type prefix will expand the match with a literal match if a slash is provided",
			in: networkingv1.HTTPIngressPath{
				Path:     "/v1/api/packages/",
				PathType: &pathTypePrefix,
			},
			out: kong.StringSlice(
				"/v1/api/packages$",
				"/v1/api/packages/",
			),
		},
		{
			name: "path type prefix will provide a default when no path is provided",
			in: networkingv1.HTTPIngressPath{
				Path:     "",
				PathType: &pathTypePrefix,
			},
			out: kong.StringSlice("/"),
		},
		{
			name: "path type exact will cause an exact matching path on a regular path",
			in: networkingv1.HTTPIngressPath{
				Path:     "/v1/api/packages",
				PathType: &pathTypeExact,
			},
			out: kong.StringSlice("/v1/api/packages$"),
		},
		{
			name: "path type exact will cause an exact matching path on a regular path with a / suffix",
			in: networkingv1.HTTPIngressPath{
				Path:     "/v1/api/packages/",
				PathType: &pathTypeExact,
			},
			out: kong.StringSlice("/v1/api/packages/$"),
		},
		{
			name: "path type exact will supply a default if no path is provided",
			in: networkingv1.HTTPIngressPath{
				Path:     "",
				PathType: &pathTypeExact,
			},
			out: kong.StringSlice("/"),
		},
		{
			name: "path type implementation-specific will leave the path alone",
			in: networkingv1.HTTPIngressPath{
				Path:     "/asdfasd9jhf09432$",
				PathType: &pathTypeImplementationSpecific,
			},
			out: kong.StringSlice("/asdfasd9jhf09432$"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, pathsFromIngressPaths(tt.in))
		})
	}
}

func Test_flattenMultipleSlashes(t *testing.T) {
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
