//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestGRPCRouteExample(t *testing.T) {
	testGRPC := func(ctx context.Context, t *testing.T, manifestName string, gatewayPort int, hostname string, enableTLS bool) {
		cluster := GetClusterFromCtx(ctx)
		proxyURL := GetProxyURLFromCtx(ctx)
		manifestPath := examplesManifestPath(manifestName)
		t.Logf("applying yaml manifest %s", manifestPath)
		b, err := os.ReadFile(manifestPath)
		assert.NoError(t, err)
		manifest := string(b)
		assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, manifest))

		t.Log("verifying that GRPCRoute becomes routable")
		assert.Eventually(t, func() bool {
			if err := grpcEchoResponds(
				ctx, fmt.Sprintf("%s:%d", proxyURL.Hostname(), gatewayPort), hostname, "kong", enableTLS,
			); err != nil {
				t.Log(err)
				return false
			}
			return true
		}, consts.IngressWait, consts.WaitTick)

		t.Logf("deleting yaml manifest %s", manifestPath)
		assert.NoError(t, clusters.DeleteManifestByYAML(ctx, cluster, manifest))
	}

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindGRPCRoute).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(
				// NOTE: "kong-grpcroute" namespace is chosen for test specificially for the reason
				// to isolate the test from the default namespace and avoid any potential conflicts
				// with other tests. This namespace is hardcoded in the UDPRoute manifest.
				// Usage of "default" namespace has caused flakiness in the past because other example
				// manifests also use "kong" Gateway and those could overwrite each other.
				// We could potentially overwrite the namespace in the client context but that
				// can ony be done in ktf's https://github.com/Kong/kubernetes-testing-framework/blob/2f5b03bcf9c28f5fa11d20d85ba2c24c87650513/pkg/utils/kubernetes/generators/kubeconfig.go#L33-L36
				// which is not currently configurable.
				helpers.ControllerManagerOptAdditionalWatchNamespace("kong-grpcroute"),
			),
			withKongProxyEnvVars(map[string]string{
				"PROXY_LISTEN": `0.0.0.0:8000 http2\, 0.0.0.0:8443 http2 ssl`,
			}),
		)).
		Assess("deploying to cluster works and deployed GRPC via HTTP responds", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			testGRPC(ctx, t, "gateway-grpcroute-via-http.yaml", ktfkong.DefaultProxyHTTPPort, "example-grpc-via-http.com", false)
			return ctx
		}).
		Assess("deploying to cluster works and deployed GRPC via HTTPS responds", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			testGRPC(ctx, t, "gateway-grpcroute-via-https.yaml", ktfkong.DefaultProxyTLSServicePort, "example.com", true)
			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
