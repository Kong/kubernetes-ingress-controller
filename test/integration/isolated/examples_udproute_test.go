//go:build integration_tests

package isolated

import (
	"context"
	"os"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestExampleUDPRoute(t *testing.T) {
	t.Parallel()

	udpRouteExampleManifests := examplesManifestPath("gateway-udproute.yaml")

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindUDPRoute).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		Assess("deploying to cluster works and udp traffic is routed to the service",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				cluster := GetClusterFromCtx(ctx)
				// GetUDPURLFromCtx returns the URL of the UDP service, but with http prefix
				// http://<IP>:<PORT> (bug in KTF), taking the Host part trims scheme part.
				proxyUDPURL := GetUDPURLFromCtx(ctx).Host

				t.Logf("applying yaml manifest %s", udpRouteExampleManifests)
				b, err := os.ReadFile(udpRouteExampleManifests)
				require.NoError(t, err)

				s := string(b)
				require.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, s))
				cleaner.AddManifest(s)

				t.Logf("verifying that the UDPRoute becomes routable")
				require.EventuallyWithT(t, func(c *assert.CollectT) {
					assert.NoError(
						c, test.EchoResponds(test.ProtocolUDP, proxyUDPURL, "udproute-example-manifest"),
					)
				}, consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
