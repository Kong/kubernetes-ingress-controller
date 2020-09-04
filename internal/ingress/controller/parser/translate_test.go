package parser

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/annotations"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/controller/parser/kongstate"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestFromIngressV1beta1(t *testing.T) {
	assert := assert.New(t)
	ingressList := []*networking.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Path: "/",
										Backend: networking.IngressBackend{
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
			Spec: networking.IngressSpec{
				TLS: []networking.IngressTLS{
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
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Path: "/",
										Backend: networking.IngressBackend{
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
			Spec: networking.IngressSpec{
				Backend: &networking.IngressBackend{
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
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Path: "/.well-known/acme-challenge/yolo",
										Backend: networking.IngressBackend{
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
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Backend: networking.IngressBackend{
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
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host:             "example.com",
						IngressRuleValue: networking.IngressRuleValue{},
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
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Backend: networking.IngressBackend{
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
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Backend: networking.IngressBackend{
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
			Spec: networking.IngressSpec{
				Rules: []networking.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: networking.IngressRuleValue{
							HTTP: &networking.HTTPIngressRuleValue{
								Paths: []networking.HTTPIngressPath{
									{
										Path: "/foo//bar",
										Backend: networking.IngressBackend{
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
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networking.Ingress{})
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("simple ingress rule is parsed", func(t *testing.T) {
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networking.Ingress{
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
			[]*networking.Ingress{ingressList[0], ingressList[2]},
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
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networking.Ingress{
			ingressList[1],
		})
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret2"]))
	})
	t.Run("ingress rule with ACME like path has strip_path set to false", func(t *testing.T) {
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networking.Ingress{
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
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networking.Ingress{
			ingressList[4],
		})
		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])
	})
	t.Run("empty Ingress rule doesn't cause a panic", func(t *testing.T) {
		assert.NotPanics(func() {
			fromIngressV1beta1(logrus.New(), []*networking.Ingress{
				ingressList[5],
			})
		})
	})
	t.Run("Ingress rules with multiple ports for one Service use separate hostnames for each port", func(t *testing.T) {
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networking.Ingress{
			ingressList[6],
		})
		assert.Equal("foo-svc.foo-namespace.80.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal("foo-svc.foo-namespace.8000.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.8000"].Host)
	})
	t.Run("Ingress rule with path containing multiple slashes ('//') is skipped", func(t *testing.T) {
		parsedInfo := fromIngressV1beta1(logrus.New(), []*networking.Ingress{
			ingressList[7],
		})
		assert.Empty(parsedInfo.ServiceNameToServices)
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
			Spec: configurationv1beta1.IngressSpec{
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
			Spec: configurationv1beta1.IngressSpec{
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
			Spec: configurationv1beta1.IngressSpec{
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
