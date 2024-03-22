//go:build integration_tests

package integration

import (
	"context"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

var udpMutex sync.Mutex

// coreDNSImage is the image and version of CoreDNS that will be used for UDP
// testing.
const coreDNSImage = "registry.k8s.io/coredns/coredns:v1.8.6"

func TestUDPIngressEssentials(t *testing.T) {
	RunWhenKongExpressionRouter(t)
	t.Parallel()

	// Ensure no other UDP tests run concurrently to avoid fights over the port
	t.Log("locking UDP port")
	udpMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking UDP port")
		udpMutex.Unlock()
	})

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("configuring coredns corefile")
	cfgmap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "coredns"}, Data: map[string]string{"Corefile": corefile}}
	cfgmap, err := env.Cluster().Client().CoreV1().ConfigMaps(ns.Name).Create(ctx, cfgmap, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(cfgmap)

	t.Log("configuring a coredns deployent to deploy for UDP testing")
	container := generators.NewContainer("coredns", coreDNSImage, 53)
	container.Ports[0].Protocol = corev1.ProtocolUDP
	container.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
	container.Args = []string{"-conf", "/etc/coredns/Corefile"}
	deployment := generators.NewDeploymentForContainer(container)

	t.Log("configuring the coredns pod with a custom corefile")
	configVolume := corev1.Volume{
		Name: "config-volume",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: cfgmap.Name},
			Items:                []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}},
		}},
	}
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, configVolume)

	t.Log("deploying coredns")
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Log("exposing DNS service via UDPIngress")
	udp := &kongv1beta1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "minudp",
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.UDPIngressSpec{Rules: []kongv1beta1.UDPIngressRule{
			{
				Port: ktfkong.DefaultUDPServicePort,
				Backend: kongv1beta1.IngressBackend{
					ServiceName: service.Name,
					ServicePort: int(service.Spec.Ports[0].Port),
				},
			},
		}},
	}
	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	assert.NoError(t, err)
	udp, err = gatewayClient.ConfigurationV1beta1().UDPIngresses(ns.Name).Create(ctx, udp, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(udp)

	t.Log("configurating a net.Resolver to resolve DNS via the proxy")
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, proxyUDPURL)
		},
	}

	t.Logf("checking udpingress %s status readiness.", udp.Name)
	ingCli := gatewayClient.ConfigurationV1beta1().UDPIngresses(ns.Name)
	assert.Eventually(t, func() bool {
		curIng, err := ingCli.Get(ctx, udp.Name, metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		ingresses := curIng.Status.LoadBalancer.Ingress
		for _, ingress := range ingresses {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				ip, _, err := net.SplitHostPort(proxyUDPURL)
				if err != nil {
					return false
				}
				if ingress.IP == ip {
					t.Logf("udpingress hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
					return true
				}
			}
		}
		return false
	}, statusWait, waitTick, true)

	t.Logf("checking DNS to resolve via UDPIngress %s", udp.Name)
	assert.Eventually(t, func() bool {
		_, err := resolver.LookupHost(ctx, corednsKnownHostname)
		return err == nil
	}, ingressWait, waitTick)

	t.Logf("tearing down UDPIngress %s and ensuring backends are torn down", udp.Name)
	assert.NoError(t, gatewayClient.ConfigurationV1beta1().UDPIngresses(ns.Name).Delete(ctx, udp.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		_, err := resolver.LookupHost(ctx, corednsKnownHostname)
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				return true
			}
		}
		return false
	}, ingressWait, waitTick)
}

func TestUDPIngressTCPIngressCollision(t *testing.T) {
	RunWhenKongExpressionRouter(t)
	t.Parallel()

	t.Log("locking TCP and UDP ports")
	udpMutex.Lock()
	tcpMutex.Lock()
	t.Cleanup(func() {
		udpMutex.Unlock()
		tcpMutex.Unlock()
	})

	ctx := context.Background()
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("configuring coredns corefile")
	cfgmap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "coredns"}, Data: map[string]string{"Corefile": corefile}}
	cfgmap, err := env.Cluster().Client().CoreV1().ConfigMaps(ns.Name).Create(ctx, cfgmap, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(cfgmap)

	t.Log("configuring a coredns deployent to deploy for UDP testing")
	container := generators.NewContainer("coredns", coreDNSImage, 53)
	container.Ports[0].Protocol = corev1.ProtocolUDP
	container.Ports[0].Name = "dnsudp"
	container.Ports = append(container.Ports, corev1.ContainerPort{Name: "dnstcp", ContainerPort: 53, Protocol: corev1.ProtocolTCP})
	container.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
	container.Args = []string{"-conf", "/etc/coredns/Corefile"}
	deployment := generators.NewDeploymentForContainer(container)

	t.Log("configuring the coredns pod with a custom corefile")
	configVolume := corev1.Volume{
		Name: "config-volume",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: cfgmap.Name},
			Items:                []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}},
		}},
	}
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, configVolume)

	t.Log("deploying coredns")
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeNodePort)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	// Test initial configuration with UDP only
	t.Log("exposing DNS service via UDPIngress")
	udp := &kongv1beta1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "minudp",
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.UDPIngressSpec{Rules: []kongv1beta1.UDPIngressRule{
			{
				Port: ktfkong.DefaultUDPServicePort,
				Backend: kongv1beta1.IngressBackend{
					ServiceName: service.Name,
					ServicePort: int(service.Spec.Ports[0].Port),
				},
			},
		}},
	}
	gatewayClient, err := clientset.NewForConfig(env.Cluster().Config())
	assert.NoError(t, err)
	udp, err = gatewayClient.ConfigurationV1beta1().UDPIngresses(ns.Name).Create(ctx, udp, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(udp)

	t.Log("configurating a dns query and clients")
	query := new(dns.Msg)
	query.Id = dns.Id()
	query.Question = make([]dns.Question, 1)
	query.Question[0] = dns.Question{Name: corednsKnownHostname, Qtype: dns.TypeA, Qclass: dns.ClassINET}
	dnsUDPClient := new(dns.Client)
	dnsTCPClient := dns.Client{Net: "tcp"}

	t.Logf("checking udpingress %s status readiness.", udp.Name)
	ingCli := gatewayClient.ConfigurationV1beta1().UDPIngresses(ns.Name)
	assert.Eventually(t, func() bool {
		curIng, err := ingCli.Get(ctx, udp.Name, metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		ingresses := curIng.Status.LoadBalancer.Ingress
		for _, ingress := range ingresses {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				t.Logf("udpingress hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
				return true
			}
		}
		return false
	}, statusWait, waitTick, true)

	t.Logf("checking DNS to resolve via UDPIngress %s", udp.Name)
	assert.Eventually(t, func() bool {
		_, _, err := dnsUDPClient.Exchange(query, proxyUDPURL)
		return err == nil
	}, ingressWait, waitTick)

	// Add a TCPIngress with the same port integer as the TCPIngress, pointing to the same Service
	t.Log("exposing DNS service via TCPIngress")
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mintcp",
			Annotations: map[string]string{
				annotations.IngressClassKey: consts.IngressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{Rules: []kongv1beta1.IngressRule{
			{
				Port: ktfkong.DefaultTCPServicePort,
				Backend: kongv1beta1.IngressBackend{
					ServiceName: service.Name,
					ServicePort: int(service.Spec.Ports[1].Port),
				},
			},
		}},
	}
	assert.NoError(t, err)
	tcp, err = gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Create(ctx, tcp, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(tcp)

	t.Logf("checking tcpingress %s status readiness.", tcp.Name)
	tcpIngCli := gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name)
	assert.Eventually(t, func() bool {
		curIng, err := tcpIngCli.Get(ctx, tcp.Name, metav1.GetOptions{})
		if err != nil || curIng == nil {
			return false
		}
		ingresses := curIng.Status.LoadBalancer.Ingress
		for _, ingress := range ingresses {
			if len(ingress.Hostname) > 0 || len(ingress.IP) > 0 {
				t.Logf("tcpingress hostname %s or ip %s is ready to redirect traffic.", ingress.Hostname, ingress.IP)
				return true
			}
		}
		return false
	}, statusWait, waitTick, true)

	t.Logf("checking DNS to resolve via TCPIngress %s", tcp.Name)
	assert.Eventually(t, func() bool {
		_, _, err := dnsTCPClient.Exchange(query, proxyTCPURL)
		return err == nil
	}, ingressWait, waitTick)

	t.Logf("checking DNS to resolve via UDPIngress %s still works also", udp.Name)
	assert.Eventually(t, func() bool {
		_, _, err := dnsUDPClient.Exchange(query, proxyUDPURL)
		return err == nil
	}, ingressWait, waitTick)

	// Cleanup
	t.Logf("tearing down UDPIngress %s and ensuring backends are torn down", udp.Name)
	assert.NoError(t, gatewayClient.ConfigurationV1beta1().UDPIngresses(ns.Name).Delete(ctx, udp.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		_, _, err := dnsUDPClient.Exchange(query, proxyUDPURL)
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				return true
			}
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("tearing down TCPIngress %s and ensuring backends are torn down", tcp.Name)
	assert.NoError(t, gatewayClient.ConfigurationV1beta1().TCPIngresses(ns.Name).Delete(ctx, tcp.Name, metav1.DeleteOptions{}))
	assert.Eventually(t, func() bool {
		_, _, err := dnsTCPClient.Exchange(query, proxyTCPURL)
		if err != nil {
			if strings.Contains(err.Error(), "connection reset by peer") {
				return true
			}
		}
		return false
	}, ingressWait, waitTick)
}

const (
	corefile = `
.:53 {
    errors
    health {
       lameduck 5s
    }
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
       pods insecure
       fallthrough in-addr.arpa ip6.arpa
       ttl 30
    }
    forward . /etc/resolv.conf {
       max_concurrent 1000
    }
    cache 30
    loop
    reload
    loadbalance
    hosts {
      10.0.0.1 konghq.com
      fallthrough
    }
}
.:9999 {
    errors
    health {
       lameduck 5s
    }
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
       pods insecure
       fallthrough in-addr.arpa ip6.arpa
       ttl 30
    }
    forward . /etc/resolv.conf {
       max_concurrent 1000
    }
    cache 30
    loop
    reload
    loadbalance
    hosts {
      10.0.0.1 konghq.com
      fallthrough
    }
}
`
	// Querying this hostname should save coredns querying external DNS.
	corednsKnownHostname = "konghq.com."
)
