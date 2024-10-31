package translator

import (
	"testing"

	"github.com/kong/go-kong/kong"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/store"
)

func TestFromTCPIngressV1beta1(t *testing.T) {
	assert := assert.New(t)
	tcpIngressList := []*kongv1beta1.TCPIngress{
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
			Spec: kongv1beta1.TCPIngressSpec{
				Rules: []kongv1beta1.IngressRule{
					{
						Port: 9000,
						Backend: kongv1beta1.IngressBackend{
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
			Spec: kongv1beta1.TCPIngressSpec{
				Rules: []kongv1beta1.IngressRule{
					{
						Host: "example.com",
						Port: 9000,
						Backend: kongv1beta1.IngressBackend{
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
			Spec: kongv1beta1.TCPIngressSpec{
				TLS: []kongv1beta1.IngressTLS{
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
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*kongv1beta1.TCPIngress{},
		})
		assert.NoError(err)

		translatedInfo := mustNewTranslator(t, store).ingressRulesFromTCPIngressV1beta1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			ServiceNameToParent:   make(map[string]client.Object),
			SecretNameToSNIs:      newSecretNameToSNIs(),
		}, translatedInfo)
	})
	t.Run("empty TCPIngress return empty info", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*kongv1beta1.TCPIngress{
				tcpIngressList[0],
			},
		})
		assert.NoError(err)

		translatedInfo := mustNewTranslator(t, store).ingressRulesFromTCPIngressV1beta1()
		assert.Equal(ingressRules{
			ServiceNameToServices: make(map[string]kongstate.Service),
			ServiceNameToParent:   make(map[string]client.Object),
			SecretNameToSNIs:      newSecretNameToSNIs(),
		}, translatedInfo)
	})
	t.Run("simple TCPIngress rule is translated", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*kongv1beta1.TCPIngress{
				tcpIngressList[1],
			},
		})
		assert.NoError(err)

		translatedInfo := mustNewTranslator(t, store).ingressRulesFromTCPIngressV1beta1()
		assert.Equal(1, len(translatedInfo.ServiceNameToServices))
		svc := translatedInfo.ServiceNameToServices["default.foo-svc.80"]
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
			Tags: []*string{
				kong.String("k8s-name:foo"),
				kong.String("k8s-namespace:default"),
			},
		}, route.Route)
	})
	t.Run("TCPIngress rule with host is translated", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*kongv1beta1.TCPIngress{
				tcpIngressList[2],
			},
		})
		assert.NoError(err)

		translatedInfo := mustNewTranslator(t, store).ingressRulesFromTCPIngressV1beta1()
		assert.Equal(1, len(translatedInfo.ServiceNameToServices))
		svc := translatedInfo.ServiceNameToServices["default.foo-svc.80"]
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
			Tags: []*string{
				kong.String("k8s-name:foo"),
				kong.String("k8s-namespace:default"),
			},
		}, route.Route)
	})
	t.Run("TCPIngress with TLS", func(t *testing.T) {
		store, err := store.NewFakeStore(store.FakeObjects{
			TCPIngresses: []*kongv1beta1.TCPIngress{
				tcpIngressList[3],
			},
		})
		assert.NoError(err)

		translatedInfo := mustNewTranslator(t, store).ingressRulesFromTCPIngressV1beta1()
		assert.Equal(2, len(translatedInfo.SecretNameToSNIs.Hosts("default/sooper-secret")))
		assert.Equal(2, len(translatedInfo.SecretNameToSNIs.Hosts("default/sooper-secret2")))
	})
}
