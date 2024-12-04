//go:build integration_tests

package isolated

import (
	"context"
	"os"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestExampleTCPRoute(t *testing.T) {
	tcpRouteExampleManifests := examplesManifestPath("gateway-tcproute.yaml")

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindTCPRoute).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		Assess("deploying to cluster works and tcp traffic is routed to the service",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				cluster := GetClusterFromCtx(ctx)
				proxyTCPURL := GetTCPURLFromCtx(ctx)

				t.Logf("applying yaml manifest %s", tcpRouteExampleManifests)
				b, err := os.ReadFile(tcpRouteExampleManifests)
				assert.NoError(t, err)

				s := string(b)
				assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, s))
				cleaner.AddManifest(s)

				t.Logf("verifying that the TCPRoute becomes routable")
				assertEventuallyResponseTCP(t, proxyTCPURL, "tcproute-example-manifest")

				return ctx
			}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
