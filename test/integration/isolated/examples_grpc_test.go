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
		Assess("deploying to cluster works and deployed GRPC via HTTP responds", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			testGRPC(ctx, t, manifestPath("gateway-grpcroute-via-http.yaml"), ktfkong.DefaultProxyHTTPPort, false)
			return ctx
		}).
		Assess("deploying to cluster works and deployed GRPC via HTTPS responds", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			testGRPC(ctx, t, manifestPath("gateway-grpcroute-via-https.yaml"), ktfkong.DefaultProxyTLSServicePort, true)
			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}

func testGRPC(ctx context.Context, t *testing.T, manifestPath string, gatewayPort int, enableTLS bool) {
	cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
	cluster := GetClusterFromCtx(ctx)
	proxyURL := GetProxyURLFromCtx(ctx)
	t.Logf("applying yaml manifest %s", manifestPath)
	b, err := os.ReadFile(manifestPath)
	assert.NoError(t, err)
	assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, string(b)))
	cleaner.AddManifest(string(b))

	t.Log("verifying that GRPCRoute becomes routable")
	assert.Eventually(t, func() bool {
		if err := grpcEchoResponds(
			ctx, fmt.Sprintf("%s:%d", proxyURL.Hostname(), gatewayPort), "example.com", "kong", enableTLS,
		); err != nil {
			t.Log(err)
			return false
		}
		return true
	}, consts.IngressWait, consts.WaitTick)
}
