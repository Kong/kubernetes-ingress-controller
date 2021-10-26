package parser

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	knative "knative.dev/networking/pkg/apis/networking/v1alpha1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/kongstate"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
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
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networkingv1beta1.Ingress{})
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("simple ingress rule is parsed", func(t *testing.T) {
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networkingv1beta1.Ingress{
			ingressList[0],
		})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])
	})
	t.Run("ingress rule with default backend", func(t *testing.T) {
		parsedInfo := fromIngressV1beta1(logrus.New(),
			[]*networkingv1beta1.Ingress{ingressList[0], ingressList[2]},
		)
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
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networkingv1beta1.Ingress{
			ingressList[1],
		})
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret2"]))
	})
	t.Run("ingress rule with ACME like path has strip_path set to false", func(t *testing.T) {
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networkingv1beta1.Ingress{
			ingressList[3],
		})
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
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networkingv1beta1.Ingress{
			ingressList[4],
		})
		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])
	})
	t.Run("empty Ingress rule doesn't cause a panic", func(t *testing.T) {
		assert.NotPanics(func() {
			fromIngressV1beta1(logrus.New(), []*networkingv1beta1.Ingress{
				ingressList[5],
			})
		})
	})
	t.Run("Ingress rules with multiple ports for one Service use separate hostnames for each port", func(t *testing.T) {
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networkingv1beta1.Ingress{
			ingressList[6],
		})
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal("foo-svc.foo-namespace.8000.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.8000"].Host)
	})
	t.Run("Ingress rule with path containing multiple slashes ('//') is skipped", func(t *testing.T) {
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networkingv1beta1.Ingress{
			ingressList[7],
		})
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
		parsedInfo := fromIngressV1(logrus.New(), []*networkingv1.Ingress{})
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("simple ingress rule is parsed", func(t *testing.T) {
		parsedInfo := fromIngressV1(logrus.New(), []*networkingv1.Ingress{
			ingressList[0],
		})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Hosts[0])
	})
	t.Run("ingress rule with default backend", func(t *testing.T) {
		parsedInfo := fromIngressV1(logrus.New(),
			[]*networkingv1.Ingress{ingressList[0], ingressList[2]},
		)
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
		parsedInfo := fromIngressV1(logrus.New(), []*networkingv1.Ingress{
			ingressList[1],
		})
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret2"]))
	})
	t.Run("ingress rule with ACME like path has strip_path set to false", func(t *testing.T) {
		parsedInfo := fromIngressV1(logrus.New(), []*networkingv1.Ingress{
			ingressList[3],
		})
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
		parsedInfo := fromIngressV1(logrus.New(), []*networkingv1.Ingress{
			ingressList[4],
		})
		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Routes[0].Hosts[0])
	})
	t.Run("empty Ingress rule doesn't cause a panic", func(t *testing.T) {
		assert.NotPanics(func() {
			fromIngressV1(logrus.New(), []*networkingv1.Ingress{
				ingressList[5],
			})
		})
	})
	t.Run("Ingress rules with multiple ports for one Service use separate hostnames for each port", func(t *testing.T) {
		parsedInfo := fromIngressV1(logrus.New(), []*networkingv1.Ingress{
			ingressList[6],
		})
		assert.Equal("foo-svc.foo-namespace.80.svc",
			*parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-80"].Host)
		assert.Equal("foo-svc.foo-namespace.8000.svc",
			*parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pnum-8000"].Host)
	})
	t.Run("Ingress rule with path containing multiple slashes ('//') is skipped", func(t *testing.T) {
		parsedInfo := fromIngressV1(logrus.New(), []*networkingv1.Ingress{
			ingressList[7],
		})
		assert.Empty(parsedInfo.ServiceNameToServices)
	})
	t.Run("Ingress rule with ports defined by name", func(t *testing.T) {
		parsedInfo := fromIngressV1(logrus.New(), []*networkingv1.Ingress{
			ingressList[8],
		})
		_, ok := parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pname-http"]
		assert.True(ok)
		_, ok = parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.pname-ws"]
		assert.True(ok)
	})
}

func TestFromTCPIngressV1beta1(t *testing.T) {
	assert := assert.New(t)
	tcpIngressList := []*configurationv1beta1.TCPIngress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
		},
		// 1
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
			Spec: configurationv1beta1.TCPIngressSpec{
				Rules: []configurationv1beta1.IngressRule{
					{
						Port: 9000,
						Backend: configurationv1beta1.IngressBackend{
							ServiceName: "foo-svc",
							ServicePort: 80,
						},
					},
				},
			},
		},
		// 2
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
			Spec: configurationv1beta1.TCPIngressSpec{
				Rules: []configurationv1beta1.IngressRule{
					{
						Host: "example.com",
						Port: 9000,
						Backend: configurationv1beta1.IngressBackend{
							ServiceName: "foo-svc",
							ServicePort: 80,
						},
					},
				},
			},
		},
		// 3
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
			Spec: configurationv1beta1.TCPIngressSpec{
				TLS: []configurationv1beta1.IngressTLS{
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
			},
		},
	}
	t.Run("no TCPIngress returns empty info", func(t *testing.T) {
		parsedInfo := fromTCPIngressV1beta1(logrus.New(), []*configurationv1beta1.TCPIngress{})
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("empty TCPIngress return empty info", func(t *testing.T) {
		parsedInfo := fromTCPIngressV1beta1(logrus.New(), []*configurationv1beta1.TCPIngress{tcpIngressList[0]})
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("simple TCPIngress rule is parsed", func(t *testing.T) {
		parsedInfo := fromTCPIngressV1beta1(logrus.New(), []*configurationv1beta1.TCPIngress{tcpIngressList[1]})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		svc := parsedInfo.ServiceNameToServices["default.foo-svc.80"]
		assert.Equal("foo-svc.default.80.svc", *svc.Host)
		assert.Equal(80, *svc.Port)
		assert.Equal("tcp", *svc.Protocol)

		assert.Equal(1, len(svc.Routes))
		route := svc.Routes[0]
		assert.Equal(kong.Route{
			Name:      kong.String("default.foo.0"),
			Protocols: kong.StringSlice("tcp", "tls"),
			Destinations: []*kong.CIDRPort{
				{
					Port: kong.Int(9000),
				},
			},
		}, route.Route)
	})
	t.Run("TCPIngress rule with host is parsed", func(t *testing.T) {
		parsedInfo := fromTCPIngressV1beta1(logrus.New(), []*configurationv1beta1.TCPIngress{tcpIngressList[2]})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		svc := parsedInfo.ServiceNameToServices["default.foo-svc.80"]
		assert.Equal("foo-svc.default.80.svc", *svc.Host)
		assert.Equal(80, *svc.Port)
		assert.Equal("tcp", *svc.Protocol)

		assert.Equal(1, len(svc.Routes))
		route := svc.Routes[0]
		assert.Equal(kong.Route{
			Name:      kong.String("default.foo.0"),
			Protocols: kong.StringSlice("tcp", "tls"),
			SNIs:      kong.StringSlice("example.com"),
			Destinations: []*kong.CIDRPort{
				{
					Port: kong.Int(9000),
				},
			},
		}, route.Route)
	})
	t.Run("TCPIngress with TLS", func(t *testing.T) {
		parsedInfo := fromTCPIngressV1beta1(logrus.New(), []*configurationv1beta1.TCPIngress{tcpIngressList[3]})
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["default/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["default/sooper-secret2"]))
	})
}

func TestFromKnativeIngress(t *testing.T) {
	assert := assert.New(t)
	ingressList := []*knative.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{},
				},
			},
		},
		// 1
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{
						Hosts: []string{"my-func.example.com"},
						HTTP: &knative.HTTPIngressRuleValue{
							Paths: []knative.HTTPIngressPath{
								{
									Path: "/",
									AppendHeaders: map[string]string{
										"foo": "bar",
									},
									Splits: []knative.IngressBackendSplit{
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "foo-ns",
												ServiceName:      "foo-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 100,
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
				Name:      "foo",
				Namespace: "foo-namespace",
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{
						Hosts: []string{"my-func.example.com"},
						HTTP: &knative.HTTPIngressRuleValue{
							Paths: []knative.HTTPIngressPath{
								{
									Path: "/",
									AppendHeaders: map[string]string{
										"foo": "bar",
									},
									Splits: []knative.IngressBackendSplit{
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "bar-ns",
												ServiceName:      "bar-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 20,
										},
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "foo-ns",
												ServiceName:      "foo-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 100,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// 3
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
			},
			Spec: knative.IngressSpec{
				Rules: []knative.IngressRule{
					{
						Hosts: []string{"my-func.example.com"},
						HTTP: &knative.HTTPIngressRuleValue{
							Paths: []knative.HTTPIngressPath{
								{
									Path: "/",
									AppendHeaders: map[string]string{
										"foo": "bar",
									},
									Splits: []knative.IngressBackendSplit{
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "bar-ns",
												ServiceName:      "bar-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 20,
										},
										{
											IngressBackend: knative.IngressBackend{
												ServiceNamespace: "foo-ns",
												ServiceName:      "foo-svc",
												ServicePort:      intstr.FromInt(42),
											},
											Percent: 100,
										},
									},
								},
							},
						},
					},
				},
				TLS: []knative.IngressTLS{
					{
						Hosts: []string{
							"foo.example.com",
							"foo1.example.com",
						},
						SecretName: "foo-secret",
					},
					{
						Hosts: []string{
							"bar.example.com",
							"bar1.example.com",
						},
						SecretName: "bar-secret",
					},
				},
			},
		},
	}
	t.Run("no ingress returns empty info", func(t *testing.T) {
		parsedInfo := fromKnativeIngress(logrus.New(), []*knative.Ingress{})
		assert.Equal(map[string]kongstate.Service{}, parsedInfo.ServiceNameToServices)
		assert.Equal(newSecretNameToSNIs(), parsedInfo.SecretNameToSNIs)
	})
	t.Run("empty ingress returns empty info", func(t *testing.T) {
		parsedInfo := fromKnativeIngress(logrus.New(), []*knative.Ingress{ingressList[0]})
		assert.Equal(map[string]kongstate.Service{}, parsedInfo.ServiceNameToServices)
		assert.Equal(newSecretNameToSNIs(), parsedInfo.SecretNameToSNIs)
	})
	t.Run("basic knative Ingress resource is parsed", func(t *testing.T) {
		parsedInfo := fromKnativeIngress(logrus.New(), []*knative.Ingress{ingressList[1]})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		svc := parsedInfo.ServiceNameToServices["foo-ns.foo-svc.42"]
		assert.Equal(kong.Service{
			Name:           kong.String("foo-ns.foo-svc.42"),
			Port:           kong.Int(80),
			Host:           kong.String("foo-svc.foo-ns.42.svc"),
			Path:           kong.String("/"),
			Protocol:       kong.String("http"),
			WriteTimeout:   kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			ConnectTimeout: kong.Int(60000),
			Retries:        kong.Int(5),
		}, svc.Service)
		assert.Equal(kong.Route{
			Name:          kong.String("foo-namespace.foo.00"),
			RegexPriority: kong.Int(0),
			StripPath:     kong.Bool(false),
			Paths:         kong.StringSlice("/"),
			PreserveHost:  kong.Bool(true),
			Protocols:     kong.StringSlice("http", "https"),
			Hosts:         kong.StringSlice("my-func.example.com"),
		}, svc.Routes[0].Route)
		assert.Equal(kong.Plugin{
			Name: kong.String("request-transformer"),
			Config: kong.Configuration{
				"add": map[string]interface{}{
					"headers": []string{"foo:bar"},
				},
			},
		}, svc.Plugins[0])

		assert.Equal(newSecretNameToSNIs(), parsedInfo.SecretNameToSNIs)
	})
	t.Run("knative TLS section is correctly parsed", func(t *testing.T) {
		parsedInfo := fromKnativeIngress(logrus.New(), []*knative.Ingress{ingressList[3]})

		assert.Equal(SecretNameToSNIs(map[string][]string{
			"foo-namespace/bar-secret": {"bar.example.com", "bar1.example.com"},
			"foo-namespace/foo-secret": {"foo.example.com", "foo1.example.com"},
		}), parsedInfo.SecretNameToSNIs)
	})
	t.Run("split knative Ingress resource chooses the highest split", func(t *testing.T) {
		parsedInfo := fromKnativeIngress(logrus.New(), []*knative.Ingress{ingressList[2]})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		svc := parsedInfo.ServiceNameToServices["foo-ns.foo-svc.42"]
		assert.Equal(kong.Service{
			Name:           kong.String("foo-ns.foo-svc.42"),
			Port:           kong.Int(80),
			Host:           kong.String("foo-svc.foo-ns.42.svc"),
			Path:           kong.String("/"),
			Protocol:       kong.String("http"),
			WriteTimeout:   kong.Int(60000),
			ReadTimeout:    kong.Int(60000),
			ConnectTimeout: kong.Int(60000),
			Retries:        kong.Int(5),
		}, svc.Service)
		assert.Equal(kong.Route{
			Name:          kong.String("foo-namespace.foo.00"),
			RegexPriority: kong.Int(0),
			StripPath:     kong.Bool(false),
			Paths:         kong.StringSlice("/"),
			PreserveHost:  kong.Bool(true),
			Protocols:     kong.StringSlice("http", "https"),
			Hosts:         kong.StringSlice("my-func.example.com"),
		}, svc.Routes[0].Route)
		assert.Equal(kong.Plugin{
			Name: kong.String("request-transformer"),
			Config: kong.Configuration{
				"add": map[string]interface{}{
					"headers": []string{"foo:bar"},
				},
			},
		}, svc.Plugins[0])

		assert.Equal(newSecretNameToSNIs(), parsedInfo.SecretNameToSNIs)
	})
}

func TestPathsFromK8s(t *testing.T) {
	for _, tt := range []struct {
		name         string
		path         string
		wantPrefix   []*string
		wantExact    []*string
		wantImplSpec []*string
	}{
		{
			name:         "empty",
			wantPrefix:   kong.StringSlice("/"),
			wantExact:    kong.StringSlice("/$"),
			wantImplSpec: kong.StringSlice("/"),
		},
		{
			name:         "root",
			path:         "/",
			wantPrefix:   kong.StringSlice("/"),
			wantExact:    kong.StringSlice("/$"),
			wantImplSpec: kong.StringSlice("/"),
		},
		{
			name:         "one segment, no trailing slash",
			path:         "/foo",
			wantPrefix:   kong.StringSlice("/foo$", "/foo/"),
			wantExact:    kong.StringSlice("/foo$"),
			wantImplSpec: kong.StringSlice("/foo"),
		},
		{
			name:         "one segment, has trailing slash",
			path:         "/foo/",
			wantPrefix:   kong.StringSlice("/foo$", "/foo/"),
			wantExact:    kong.StringSlice("/foo/$"),
			wantImplSpec: kong.StringSlice("/foo/"),
		},
		{
			name:         "two segments, no trailing slash",
			path:         "/foo/bar",
			wantPrefix:   kong.StringSlice("/foo/bar$", "/foo/bar/"),
			wantExact:    kong.StringSlice("/foo/bar$"),
			wantImplSpec: kong.StringSlice("/foo/bar"),
		},
		{
			name:         "two segments, has trailing slash",
			path:         "/foo/bar/",
			wantPrefix:   kong.StringSlice("/foo/bar$", "/foo/bar/"),
			wantExact:    kong.StringSlice("/foo/bar/$"),
			wantImplSpec: kong.StringSlice("/foo/bar/"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			{
				gotPrefix, gotErr := pathsFromK8s(tt.path, networkingv1.PathTypePrefix)
				require.NoError(t, gotErr)
				require.Equal(t, tt.wantPrefix, gotPrefix, "prefix match")
			}
			{
				gotExact, gotErr := pathsFromK8s(tt.path, networkingv1.PathTypeExact)
				require.NoError(t, gotErr)
				require.Equal(t, tt.wantExact, gotExact, "exact match")
			}
			{
				gotImplSpec, gotErr := pathsFromK8s(tt.path, networkingv1.PathTypeImplementationSpecific)
				require.NoError(t, gotErr)
				require.Equal(t, tt.wantImplSpec, gotImplSpec, "implementation specific match")
			}
		})
	}
}
