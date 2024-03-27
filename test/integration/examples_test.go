//go:build integration_tests

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
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

const examplesDIR = "../../examples"

func TestHTTPRouteExample(t *testing.T) {
	var (
		httprouteExampleManifests = fmt.Sprintf("%s/gateway-httproute.yaml", examplesDIR)
		ctx                       = context.Background()
	)

	_, cleaner := helpers.Setup(ctx, t, env)

	t.Logf("configuring test and setting up API clients")
	gwc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Logf("applying yaml manifest %s", httprouteExampleManifests)
	b, err := os.ReadFile(httprouteExampleManifests)
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyManifestByYAML(ctx, env.Cluster(), string(b)))
	cleaner.AddManifest(string(b))

	t.Logf("verifying that the Gateway receives listen addresses")
	var gatewayAddr string
	require.Eventually(t, func() bool {
		obj, err := gwc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Get(ctx, "kong", metav1.GetOptions{})
		if err != nil {
			return false
		}

		for _, addr := range obj.Status.Addresses {
			if addr.Type != nil && *addr.Type == gatewayv1.IPAddressType {
				gatewayAddr = addr.Value
				return true
			}
		}

		return false
	}, gatewayUpdateWaitTime, waitTick)

	t.Logf("verifying that the HTTPRoute becomes routable")
	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("http://%s/httproute-testing", gatewayAddr))
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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("http://%s/httproute-testing", gatewayAddr))
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

func TestUDPRouteExample(t *testing.T) {
	RunWhenKongExpressionRouterWithVersion(t, ">=3.4.0")
	t.Log("locking UDP port")
	udpMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking UDP port")
		udpMutex.Unlock()
	})

	var (
		ctx                      = context.Background()
		udpRouteExampleManifests = fmt.Sprintf("%s/gateway-udproute.yaml", examplesDIR)
	)

	_, cleaner := helpers.Setup(ctx, t, env)

	t.Logf("applying yaml manifest %s", udpRouteExampleManifests)
	b, err := os.ReadFile(udpRouteExampleManifests)
	// TODO as of 2022-04-01, UDPRoute does not support using a different inbound port than the outbound
	// destination service port. Once parentRef port functionality is stable, we should remove this
	s := string(b)
	s = strings.ReplaceAll(s, "port: 53", "port: 9999")
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyManifestByYAML(ctx, env.Cluster(), s))
	cleaner.AddManifest(s)

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
		if err != nil {
			t.Logf("failed resolving kernel.org: %v", err)
			return false
		}
		return true
	}, ingressWait, waitTick)
}

func TestTCPRouteExample(t *testing.T) {
	RunWhenKongExpressionRouterWithVersion(t, ">=3.4.0")
	t.Log("locking TCP port")
	tcpMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TCP port")
		tcpMutex.Unlock()
	})

	var (
		ctx                      = context.Background()
		tcprouteExampleManifests = fmt.Sprintf("%s/gateway-tcproute.yaml", examplesDIR)
	)
	_, cleaner := helpers.Setup(ctx, t, env)

	t.Logf("applying yaml manifest %s", tcprouteExampleManifests)
	b, err := os.ReadFile(tcprouteExampleManifests)
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyManifestByYAML(ctx, env.Cluster(), string(b)))
	cleaner.AddManifest(string(b))

	t.Log("verifying that TCPRoute becomes routable")
	require.Eventually(t, func() bool {
		responds, err := test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), "tcproute-example-manifest")
		return err == nil && responds
	}, ingressWait, waitTick)
}

func TestTLSRouteExample(t *testing.T) {
	skipTestForExpressionRouter(t)
	t.Log("locking Gateway TLS ports")
	tlsMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	})

	var (
		tlsrouteExampleManifests = fmt.Sprintf("%s/gateway-tlsroute.yaml", examplesDIR)
		ctx                      = context.Background()
	)
	_, cleaner := helpers.Setup(ctx, t, env)

	t.Logf("applying yaml manifest %s", tlsrouteExampleManifests)
	b, err := os.ReadFile(tlsrouteExampleManifests)
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyManifestByYAML(ctx, env.Cluster(), string(b)))
	cleaner.AddManifest(string(b))

	t.Log("verifying that TLSRoute becomes routable")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			"tlsroute-example-manifest", "tlsroute.kong.example", "tlsroute.kong.example", true)
		return err == nil && responded
	}, ingressWait, waitTick)
}

func TestGRPCRouteExample(t *testing.T) {
	var (
		grpcrouteExampleManifests = fmt.Sprintf("%s/gateway-grpcroute.yaml", examplesDIR)
		ctx                       = context.Background()
	)
	_, cleaner := helpers.Setup(ctx, t, env)

	t.Logf("applying yaml manifest %s", grpcrouteExampleManifests)
	b, err := os.ReadFile(grpcrouteExampleManifests)
	require.NoError(t, err)
	require.NoError(t, clusters.ApplyManifestByYAML(ctx, env.Cluster(), string(b)))
	cleaner.AddManifest(string(b))

	t.Log("verifying that GRPCRoute becomes routable")
	require.Eventually(t, func() bool {
		err := grpcEchoResponds(ctx, fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultProxyTLSServicePort), "example.com", "kong")
		if err != nil {
			t.Log(err)
		}
		return err == nil
	}, ingressWait, waitTick)
}

func TestIngressExample(t *testing.T) {
	var (
		ingressExampleManifests = fmt.Sprintf("%s/ingress.yaml", examplesDIR)
		ctx                     = context.Background()
	)

	_, cleaner := helpers.Setup(ctx, t, env)

	t.Logf("applying yaml manifest %s", ingressExampleManifests)
	b, err := os.ReadFile(ingressExampleManifests)
	require.NoError(t, err)
	manifests := replaceIngressClassSpecFieldInManifests(string(b))
	require.NoError(t, clusters.ApplyManifestByYAML(ctx, env.Cluster(), manifests))
	cleaner.AddManifest(string(b))

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
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("http://%s/", ingAddr))
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

func TestUDPIngressExample(t *testing.T) {
	skipTestForExpressionRouter(t)
	t.Log("locking UDP port")
	udpMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking UDP port")
		udpMutex.Unlock()
	})

	var (
		udpingressExampleManifests = fmt.Sprintf("%s/udpingress.yaml", examplesDIR)
		ctx                        = context.Background()
	)
	_, cleaner := helpers.Setup(ctx, t, env)

	t.Logf("applying yaml manifest %s", udpingressExampleManifests)
	b, err := os.ReadFile(udpingressExampleManifests)
	require.NoError(t, err)
	manifests := replaceIngressClassAnnotationInManifests(string(b))
	require.NoError(t, clusters.ApplyManifestByYAML(ctx, env.Cluster(), manifests))
	cleaner.AddManifest(string(b))

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

func replaceIngressClassAnnotationInManifests(manifests string) string {
	return strings.ReplaceAll(manifests, `kubernetes.io/ingress.class: "kong"`, fmt.Sprintf(`kubernetes.io/ingress.class: "%s"`, consts.IngressClass))
}

func replaceIngressClassSpecFieldInManifests(manifests string) string {
	return strings.ReplaceAll(manifests, `ingressClassName: kong`, fmt.Sprintf(`ingressClassName: %s`, consts.IngressClass))
}
