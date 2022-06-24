//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

const examplesDIR = "../../examples"

var httprouteExampleManifests = fmt.Sprintf("%s/gateway-httproute.yaml", examplesDIR)

func TestHTTPRouteExample(t *testing.T) {
	t.Logf("configuring test and setting up API clients")
	gwc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Logf("applying yaml manifest %s", strings.TrimPrefix(httprouteExampleManifests, examplesDIR))
	b, err := os.ReadFile(httprouteExampleManifests)
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyYAML(ctx, env.Cluster(), string(b)))

	defer func() {
		require.NoError(t, clusters.DeleteYAML(ctx, env.Cluster(), string(b)))
	}()

	t.Logf("verifying that the Gateway receives listen addresses")
	var gatewayAddr string
	require.Eventually(t, func() bool {
		obj, err := gwc.GatewayV1alpha2().Gateways(corev1.NamespaceDefault).Get(ctx, "kong", metav1.GetOptions{})
		if err != nil {
			return false
		}

		for _, addr := range obj.Status.Addresses {
			if addr.Type != nil && *addr.Type == gatewayv1alpha2.IPAddressType {
				gatewayAddr = addr.Value
				return true
			}
		}

		return false
	}, gatewayUpdateWaitTime, waitTick)

	t.Logf("verifying that the HTTPRoute becomes routable")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("http://%s/httproute-testing", gatewayAddr))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Logf("verifying that the backendRefs are being loadbalanced")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("http://%s/httproute-testing", gatewayAddr))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>Welcome to nginx!</title>")
		}
		return false
	}, ingressWait, waitTick)
}

var udpRouteExampleManifests = fmt.Sprintf("%s/gateway-udproute.yaml", examplesDIR)

func TestUDPRouteExample(t *testing.T) {
	t.Log("locking Gateway UDP ports")
	udpMutex.Lock()
	defer udpMutex.Unlock()

	t.Logf("applying yaml manifest %s", strings.TrimPrefix(udpRouteExampleManifests, examplesDIR))
	b, err := os.ReadFile(udpRouteExampleManifests)
	// TODO as of 2022-04-01, UDPRoute does not support using a different inbound port than the outbound
	// destination service port. Once parentRef port functionality is stable, we should remove this
	s := string(b)
	s = strings.ReplaceAll(s, "port: 53", "port: 9999")
	b = []byte(s)
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyYAML(ctx, env.Cluster(), string(b)))

	defer func() {
		require.NoError(t, clusters.DeleteYAML(ctx, env.Cluster(), string(b)))
	}()

	t.Logf("configuring test and setting up API clients")
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:%d", proxyUDPURL.Hostname(), 9999))
		},
	}

	t.Logf("verifying that the UDPRoute becomes routable")
	require.Eventually(t, func() bool {
		_, err := resolver.LookupHost(ctx, "kernel.org")
		return err == nil
	}, ingressWait, waitTick)
}

var tcprouteExampleManifests = fmt.Sprintf("%s/gateway-tcproute.yaml", examplesDIR)

func TestTCPRouteExample(t *testing.T) {
	t.Log("locking Gateway TCP ports")
	tcpMutex.Lock()
	defer tcpMutex.Unlock()

	t.Logf("applying yaml manifest %s", tcprouteExampleManifests)
	b, err := os.ReadFile(tcprouteExampleManifests)
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyYAML(ctx, env.Cluster(), string(b)))

	defer func() {
		t.Logf("deleting tcproute example")
		require.NoError(t, clusters.DeleteYAML(ctx, env.Cluster(), string(b)))
	}()

	t.Log("verifying that TCPRoute becomes routable")
	require.Eventually(t, func() bool {
		responds, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), "tcproute-example-manifest")
		return err == nil && responds
	}, ingressWait, waitTick)
}

var tlsrouteExampleManifests = fmt.Sprintf("%s/gateway-tlsroute.yaml", examplesDIR)

func TestTLSRouteExample(t *testing.T) {
	t.Log("locking Gateway TLS ports")
	tlsMutex.Lock()
	defer func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	}()

	t.Logf("applying yaml manifest %s", tlsrouteExampleManifests)
	b, err := os.ReadFile(tlsrouteExampleManifests)
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyYAML(ctx, env.Cluster(), string(b)))

	defer func() {
		t.Logf("deleting tlsroute example")
		require.NoError(t, clusters.DeleteYAML(ctx, env.Cluster(), string(b)))
	}()

	t.Log("verifying that TLSRoute becomes routable")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			"tlsroute-example-manifest", "tlsroute.kong.example", "tlsroute.kong.example")
		return err == nil && responded
	}, ingressWait, waitTick)
}

var ingressExampleManifests = fmt.Sprintf("%s/ingress.yaml", examplesDIR)

func TestIngressExample(t *testing.T) {
	t.Logf("applying yaml manifest %s", strings.TrimPrefix(ingressExampleManifests, examplesDIR))
	b, err := os.ReadFile(ingressExampleManifests)
	require.NoError(t, err)
	manifests := replaceIngressClassInManifests(string(b))
	require.NoError(t, clusters.ApplyYAML(ctx, env.Cluster(), manifests))

	defer func() {
		require.NoError(t, clusters.DeleteYAML(ctx, env.Cluster(), manifests))
	}()

	t.Log("waiting for ingress resource to have an address")
	var ingAddr string
	require.Eventually(t, func() bool {
		ing, err := env.Cluster().Client().NetworkingV1().Ingresses(corev1.NamespaceDefault).Get(ctx, "httpbin-ingress", metav1.GetOptions{})
		if err != nil {
			return false
		}

		for _, lbing := range ing.Status.LoadBalancer.Ingress {
			if lbing.IP != "" {
				ingAddr = lbing.IP
				return true
			}
		}

		return false
	}, ingressWait, waitTick)

	t.Logf("verifying that the Ingress resource becomes routable")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("http://%s/", ingAddr))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)
}

var udpingressExampleManifests = fmt.Sprintf("%s/udpingress.yaml", examplesDIR)

func TestUDPIngressExample(t *testing.T) {
	t.Log("locking Gateway UDP ports")
	udpMutex.Lock()
	defer udpMutex.Unlock()

	t.Logf("applying yaml manifest %s", strings.TrimPrefix(udpingressExampleManifests, examplesDIR))
	b, err := os.ReadFile(udpingressExampleManifests)
	require.NoError(t, err)
	manifests := replaceIngressClassInManifests(string(b))
	require.NoError(t, clusters.ApplyYAML(ctx, env.Cluster(), manifests))

	defer func() {
		require.NoError(t, clusters.DeleteYAML(ctx, env.Cluster(), manifests))
	}()

	t.Log("building a DNS query to request of CoreDNS")
	query := new(dns.Msg)
	query.Id = dns.Id()
	query.Question = make([]dns.Question, 1)
	query.Question[0] = dns.Question{Name: "kernel.org.", Qtype: dns.TypeA, Qclass: dns.ClassINET}

	t.Log("verifying that the UDPIngress resource becomes routable")
	dnsUDPClient := new(dns.Client)
	assert.Eventually(t, func() bool {
		_, _, err := dnsUDPClient.Exchange(query, fmt.Sprintf("%s:9999", proxyUDPURL.Hostname()))
		return err == nil
	}, ingressWait, waitTick)
}

func replaceIngressClassInManifests(manifests string) string {
	return strings.ReplaceAll(manifests, `kubernetes.io/ingress.class: "kong"`, fmt.Sprintf(`kubernetes.io/ingress.class: "%s"`, ingressClass))
}
