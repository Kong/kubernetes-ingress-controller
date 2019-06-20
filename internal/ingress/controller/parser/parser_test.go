package parser

import (
	"fmt"
	"testing"

	"github.com/hbagdi/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/internal/ingress/utils"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestParseIngressRules(t *testing.T) {
	assert := assert.New(t)
	p := Parser{}
	ingressList := []*extensions.Ingress{
		// 0
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "foo-namespace",
			},
			Spec: extensions.IngressSpec{
				Rules: []extensions.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: extensions.IngressRuleValue{
							HTTP: &extensions.HTTPIngressRuleValue{
								Paths: []extensions.HTTPIngressPath{
									{
										Path: "/",
										Backend: extensions.IngressBackend{
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
			Spec: extensions.IngressSpec{
				TLS: []extensions.IngressTLS{
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
				Rules: []extensions.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: extensions.IngressRuleValue{
							HTTP: &extensions.HTTPIngressRuleValue{
								Paths: []extensions.HTTPIngressPath{
									{
										Path: "/",
										Backend: extensions.IngressBackend{
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
			Spec: extensions.IngressSpec{
				Backend: &extensions.IngressBackend{
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
			Spec: extensions.IngressSpec{
				Rules: []extensions.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: extensions.IngressRuleValue{
							HTTP: &extensions.HTTPIngressRuleValue{
								Paths: []extensions.HTTPIngressPath{
									{
										Path: "/.well-known/acme-challenge/yolo",
										Backend: extensions.IngressBackend{
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
			Spec: extensions.IngressSpec{
				Rules: []extensions.IngressRule{
					{
						Host: "example.com",
						IngressRuleValue: extensions.IngressRuleValue{
							HTTP: &extensions.HTTPIngressRuleValue{
								Paths: []extensions.HTTPIngressPath{
									{
										Backend: extensions.IngressBackend{
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
		parsedInfo, err := p.parseIngressRules([]*extensions.Ingress{})
		assert.Equal(&parsedIngressRules{
			ServiceNameToServices: make(map[string]Service),
			SecretNameToSNIs:      make(map[string][]string),
		}, parsedInfo)
		assert.Nil(err)
	})
	t.Run("simple ingress rule is parsed", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*extensions.Ingress{
			ingressList[0],
		})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])
		assert.Nil(err)
	})
	t.Run("ingress rule with default backend", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*extensions.Ingress{
			ingressList[0],
			ingressList[2],
		})
		assert.Equal(2, len(parsedInfo.ServiceNameToServices))
		assert.Equal("foo-svc.foo-namespace.svc", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Port)

		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])

		assert.Equal(1, len(parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes))
		assert.Equal("/", *parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes[0].Paths[0])
		assert.Equal(0, len(parsedInfo.ServiceNameToServices["bar-namespace.default-svc.80"].Routes[0].Hosts))
		assert.Nil(err)
	})
	t.Run("ingress rule with TLS", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*extensions.Ingress{
			ingressList[1],
		})
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret"]))
		assert.Equal(2, len(parsedInfo.SecretNameToSNIs["bar-namespace/sooper-secret2"]))

		assert.Nil(err)
	})
	t.Run("ingress rule with ACME like path has strip_path set to false", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*extensions.Ingress{
			ingressList[3],
		})
		assert.Equal(1, len(parsedInfo.ServiceNameToServices))
		assert.Equal("cert-manager-solver-pod.foo-namespace.svc", *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Host)
		assert.Equal(80, *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Port)

		assert.Equal("/.well-known/acme-challenge/yolo", *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Routes[0].Hosts[0])
		assert.False(*parsedInfo.ServiceNameToServices["foo-namespace.cert-manager-solver-pod.80"].Routes[0].StripPath)

		assert.Nil(err)
	})
	t.Run("ingress with empty path is correctly parsed", func(t *testing.T) {
		parsedInfo, err := p.parseIngressRules([]*extensions.Ingress{
			ingressList[4],
		})
		assert.Equal("/", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Paths[0])
		assert.Equal("example.com", *parsedInfo.ServiceNameToServices["foo-namespace.foo-svc.80"].Routes[0].Hosts[0])

		assert.Nil(err)
	})
}

func TestOverrideService(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inService      Service
		inKongIngresss configurationv1.KongIngress
		outService     Service
	}{
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Protocol: kong.String("https"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Retries: kong.Int(0),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
					Retries:  kong.Int(0),
				},
			},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Path: kong.String("/new-path"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/new-path"),
				},
			},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					Retries: kong.Int(1),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
					Retries:  kong.Int(1),
				},
			},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			configurationv1.KongIngress{
				Proxy: &kong.Service{
					ConnectTimeout: kong.Int(100),
					ReadTimeout:    kong.Int(100),
					WriteTimeout:   kong.Int(100),
				},
			},
			Service{
				Service: kong.Service{
					Host:           kong.String("foo.com"),
					Port:           kong.Int(80),
					Name:           kong.String("foo"),
					Protocol:       kong.String("http"),
					Path:           kong.String("/"),
					ConnectTimeout: kong.Int(100),
					ReadTimeout:    kong.Int(100),
					WriteTimeout:   kong.Int(100),
				},
			},
		},
	}

	for _, testcase := range testTable {
		overrideService(&testcase.inService, &testcase.inKongIngresss)
		assert.Equal(testcase.inService, testcase.outService)
	}

	assert.NotPanics(func() {
		overrideService(nil, nil)
	})
}

func TestOverrideRoute(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inRoute        Route
		inKongIngresss configurationv1.KongIngress
		outRoute       Route
	}{
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			configurationv1.KongIngress{},
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			configurationv1.KongIngress{
				Route: &kong.Route{
					Methods: kong.StringSlice("GET", "POST"),
				},
			},
			Route{
				Route: kong.Route{
					Hosts:   kong.StringSlice("foo.com", "bar.com"),
					Methods: kong.StringSlice("GET", "POST"),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts: kong.StringSlice("foo.com", "bar.com"),
				},
			},
			configurationv1.KongIngress{
				Route: &kong.Route{
					HTTPSRedirectStatusCode: kong.Int(302),
				},
			},
			Route{
				Route: kong.Route{
					Hosts:                   kong.StringSlice("foo.com", "bar.com"),
					HTTPSRedirectStatusCode: kong.Int(302),
				},
			},
		},
		{
			Route{
				Route: kong.Route{
					Hosts:        kong.StringSlice("foo.com", "bar.com"),
					PreserveHost: kong.Bool(true),
					StripPath:    kong.Bool(true),
				},
			},
			configurationv1.KongIngress{
				Route: &kong.Route{
					Protocols:     kong.StringSlice("http"),
					PreserveHost:  kong.Bool(false),
					StripPath:     kong.Bool(false),
					RegexPriority: kong.Int(10),
				},
			},
			Route{
				Route: kong.Route{
					Hosts:         kong.StringSlice("foo.com", "bar.com"),
					Protocols:     kong.StringSlice("http"),
					PreserveHost:  kong.Bool(false),
					StripPath:     kong.Bool(false),
					RegexPriority: kong.Int(10),
				},
			},
		},
	}

	for _, testcase := range testTable {
		overrideRoute(&testcase.inRoute, &testcase.inKongIngresss)
		assert.Equal(testcase.inRoute, testcase.outRoute)
	}

	assert.NotPanics(func() {
		overrideRoute(nil, nil)
	})
}

func TestOverrideUpstream(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inUpstream     Upstream
		inKongIngresss configurationv1.KongIngress
		outUpstream    Upstream
	}{
		{
			Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			configurationv1.KongIngress{
				Upstream: &kong.Upstream{},
			},
			Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
		},
		{
			Upstream{
				Upstream: kong.Upstream{
					Name: kong.String("foo.com"),
				},
			},
			configurationv1.KongIngress{
				Upstream: &kong.Upstream{
					Name:               kong.String("wrong.com"),
					HashOn:             kong.String("HashOn"),
					HashOnCookie:       kong.String("HashOnCookie"),
					HashOnCookiePath:   kong.String("HashOnCookiePath"),
					HashOnHeader:       kong.String("HashOnHeader"),
					HashFallback:       kong.String("HashFallback"),
					HashFallbackHeader: kong.String("HashFallbackHeader"),
					Slots:              kong.Int(42),
				},
			},
			Upstream{
				Upstream: kong.Upstream{
					Name:               kong.String("foo.com"),
					HashOn:             kong.String("HashOn"),
					HashOnCookie:       kong.String("HashOnCookie"),
					HashOnCookiePath:   kong.String("HashOnCookiePath"),
					HashOnHeader:       kong.String("HashOnHeader"),
					HashFallback:       kong.String("HashFallback"),
					HashFallbackHeader: kong.String("HashFallbackHeader"),
					Slots:              kong.Int(42),
				},
			},
		},
	}

	for _, testcase := range testTable {
		overrideUpstream(&testcase.inUpstream, &testcase.inKongIngresss)
		assert.Equal(testcase.inUpstream, testcase.outUpstream)
	}

	assert.NotPanics(func() {
		overrideUpstream(nil, nil)
	})
}

func TestGetEndpoints(t *testing.T) {
	tests := []struct {
		name   string
		svc    *corev1.Service
		port   *corev1.ServicePort
		proto  corev1.Protocol
		fn     func(string, string) (*corev1.Endpoints, error)
		result []utils.Endpoint
	}{
		{
			"no service should return 0 endpoints",
			nil,
			nil,
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, nil
			},
			[]utils.Endpoint{},
		},
		{
			"no service port should return 0 endpoints",
			&corev1.Service{},
			nil,
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, nil
			},
			[]utils.Endpoint{},
		},
		{
			"a service without endpoints should return 0 endpoints",
			&corev1.Service{},
			&corev1.ServicePort{Name: "default"},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"a service type ServiceTypeExternalName service with an invalid port should return 0 endpoints",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeExternalName,
				},
			},
			&corev1.ServicePort{Name: "default"},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"a service type ServiceTypeExternalName with a valid port should return one endpoint",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:         corev1.ServiceTypeExternalName,
					ExternalName: "10.0.0.1.xip.io",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return &corev1.Endpoints{}, nil
			},
			[]utils.Endpoint{
				{
					Address: "10.0.0.1.xip.io",
					Port:    "80",
				},
			},
		},
		{
			"should return no endpoints when there is an error searching for endpoints",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				return nil, fmt.Errorf("unexpected error")
			},
			[]utils.Endpoint{},
		},
		{
			"should return no endpoints when the protocol does not match",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolUDP,
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"should return no endpoints when there is no ready Addresses",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							NotReadyAddresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolUDP,
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"should return no endpoints when the name of the port name do not match any port in the endpoint Subsets",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolTCP,
									Port:     int32(80),
									Name:     "another-name",
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{},
		},
		{
			"should return one endpoint when the name of the port name match a port in the endpoint Subsets",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromInt(80),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "default",
				TargetPort: intstr.FromInt(80),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Protocol: corev1.ProtocolTCP,
									Port:     int32(80),
									Name:     "default",
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{
				{
					Address: "1.1.1.1",
					Port:    "80",
				},
			},
		},
		{
			"should return one endpoint when the name of the port name match more than one port in the endpoint Subsets",
			&corev1.Service{
				Spec: corev1.ServiceSpec{
					Type:      corev1.ServiceTypeClusterIP,
					ClusterIP: "1.1.1.1",
					Ports: []corev1.ServicePort{
						{
							Name:       "default",
							TargetPort: intstr.FromString("port-1"),
						},
					},
				},
			},
			&corev1.ServicePort{
				Name:       "port-1",
				TargetPort: intstr.FromString("port-1"),
			},
			corev1.ProtocolTCP,
			func(string, string) (*corev1.Endpoints, error) {
				nodeName := "dummy"
				return &corev1.Endpoints{
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP:       "1.1.1.1",
									NodeName: &nodeName,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Name:     "port-1",
									Protocol: corev1.ProtocolTCP,
									Port:     80,
								},
								{
									Name:     "port-1",
									Protocol: corev1.ProtocolTCP,
									Port:     80,
								},
							},
						},
					},
				}, nil
			},
			[]utils.Endpoint{
				{
					Address: "1.1.1.1",
					Port:    "80",
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := getEndpoints(testCase.svc, testCase.port, testCase.proto, testCase.fn)
			if len(testCase.result) != len(result) {
				t.Errorf("expected %v Endpoints but got %v", testCase.result, len(result))
			}
		})
	}
}
