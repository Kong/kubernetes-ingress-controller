package parser

import (
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/store"
	configurationv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func TestFromTCPIngressV1beta1(t *testing.T) {
	assert := assert.New(t)
	tcpIngressList := []*configurationv1beta1.TCPIngress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
		},
		// 1
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
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
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
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
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
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
		// 4
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: configurationv1beta1.TCPIngressSpec{
				Rules: []configurationv1beta1.IngressRule{
					{
						Port: 9000,
						Backend: configurationv1beta1.IngressBackend{
							ServiceName: "",
							ServicePort: 80,
						},
					},
				},
			},
		},
		// 5
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: configurationv1beta1.TCPIngressSpec{
				Rules: []configurationv1beta1.IngressRule{
					{
						Port: 0,
						Backend: configurationv1beta1.IngressBackend{
							ServiceName: "foo-svc",
							ServicePort: 80,
						},
					},
				},
			},
		},
		// 6
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
				Annotations: map[string]string{
					annotations.IngressClassKey: annotations.DefaultIngressClass,
				},
			},
			Spec: configurationv1beta1.TCPIngressSpec{
				Rules: []configurationv1beta1.IngressRule{
					{
						Port: 9000,
						Backend: configurationv1beta1.IngressBackend{
							ServiceName: "foo-svc",
							ServicePort: 0,
						},
					},
				},
			},
		},
	}
	t.Run("no TCPIngress returns empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*configurationv1beta1.TCPIngress{},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromTCPIngressV1beta1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("empty TCPIngress return empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*configurationv1beta1.TCPIngress{
				tcpIngressList[0],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromTCPIngressV1beta1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("simple TCPIngress rule is parsed", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*configurationv1beta1.TCPIngress{
				tcpIngressList[1],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromTCPIngressV1beta1()
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
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*configurationv1beta1.TCPIngress{
				tcpIngressList[2],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromTCPIngressV1beta1()
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
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*configurationv1beta1.TCPIngress{
				tcpIngressList[3],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromTCPIngressV1beta1()
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["default/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["default/sooper-secret2"]))
	})
	t.Run("TCPIngress without service name returns empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*configurationv1beta1.TCPIngress{
				tcpIngressList[4],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromTCPIngressV1beta1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("TCPIngress with invalid port returns empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*configurationv1beta1.TCPIngress{
				tcpIngressList[5],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromTCPIngressV1beta1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
	t.Run("empty TCPIngress with invalid service port returns empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*configurationv1beta1.TCPIngress{
				tcpIngressList[6],
			},
		})
		assert.NoError(err)
		p := NewParser(logrus.New(), store)

		parsedInfo := p.ingressRulesFromTCPIngressV1beta1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
	})
}
