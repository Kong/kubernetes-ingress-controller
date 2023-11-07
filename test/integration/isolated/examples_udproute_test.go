//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestExampleUDPRoute(t *testing.T) {
	t.Parallel()

	udpRouteExampleManifests := fmt.Sprintf("%s/gateway-udproute.yaml", examplesDIR)

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindUDPRoute).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup(
			helpers.ControllerManagerOptAdditionalWatchNamespace("default"),
		)).
		WithSetup("inject Klient into test context", injectKlient()).
		Assess("deploying to cluster works and deployed coredns responds to UDP queries",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				cluster := GetClusterFromCtx(ctx)
				proxyUDPURL := GetUDPURLFromCtx(ctx)

				t.Logf("applying yaml manifest %s", udpRouteExampleManifests)
				b, err := os.ReadFile(udpRouteExampleManifests)
				assert.NoError(t, err)

				// TODO as of 2022-04-01, UDPRoute does not support using a different inbound port than the outbound
				// destination service port. Once parentRef port functionality is stable, we should remove this
				s := string(b)
				s = strings.ReplaceAll(s, "port: 53", "port: 9999")
				assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, s))
				cleaner.AddManifest(s)

				t.Logf("verifying that the UDPRoute becomes routable")
				assert.Eventually(t, resolveFn(ctx, proxyUDPURL), 100*consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
