//go:build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
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
	var gatewayIP string
	require.Eventually(t, func() bool {
		obj, err := gwc.GatewayV1().Gateways(corev1.NamespaceDefault).Get(ctx, "kong", metav1.GetOptions{})
		if err != nil {
			return false
		}

		for _, addr := range obj.Status.Addresses {
			if addr.Type != nil && *addr.Type == gatewayapi.IPAddressType {
				gatewayIP = addr.Value
				return true
			}
		}

		return false
	}, gatewayUpdateWaitTime, waitTick)

	require.NoError(t, err)
	t.Logf("verifying that the HTTPRoute becomes routable")
	helpers.EventuallyGETPath(
		t, nil, gatewayIP, "/httproute-testing", http.StatusOK, "<title>httpbin.org</title>", nil, ingressWait, waitTick,
	)

	t.Logf("verifying that the backendRefs are being loadbalanced")
	helpers.EventuallyGETPath(
		t, nil, gatewayIP, "/httproute-testing", http.StatusOK, "<title>Welcome to nginx!</title>", nil, ingressWait, waitTick,
	)
}

func TestTCPRouteExample(t *testing.T) {
	RunWhenKongExpressionRouter(t)
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
		return test.EchoResponds(test.ProtocolTCP, fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), "tcproute-example-manifest") == nil
	}, ingressWait, waitTick)
}

func TestTLSRouteExample(t *testing.T) {
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
		return tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			"tlsroute-example-manifest", "tlsroute.kong.example", "tlsroute.kong.example", true) == nil
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
