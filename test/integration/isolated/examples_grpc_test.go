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
		t.Helper()
		cluster := GetClusterFromCtx(ctx)
		proxyURL := GetHTTPURLFromCtx(ctx)
		manifestPath := examplesManifestPath(manifestName)
		t.Logf("applying yaml manifest %s", manifestPath)
		b, err := os.ReadFile(manifestPath)
		assert.NoError(t, err)
		manifest := string(b)
		assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, manifest))

		t.Log("verifying that GRPCRoute becomes routable")
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			err := grpcEchoResponds(
				ctx, fmt.Sprintf("%s:%d", proxyURL.Hostname(), gatewayPort), hostname, "kong", enableTLS,
			)
			assert.NoError(c, err)
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
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
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
