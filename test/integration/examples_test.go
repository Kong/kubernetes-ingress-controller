//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
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
		resp, err := httpc.Get(fmt.Sprintf("http://%s/httpbin", gatewayAddr))
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
