//go:build integration_tests

package isolated

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestHTTPRouteRewriteExample(t *testing.T) {
	httprouteURLRewritePathFullExampleManifests := examplesManifestPath("gateway-httproute-rewrite-path.yaml")

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindHTTPRoute).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		Assess("deploying to cluster works and HTTP requests get properly rewritten URI",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				cluster := GetClusterFromCtx(ctx)
				proxyURL := GetHTTPURLFromCtx(ctx)

				t.Logf("applying yaml manifest %s", httprouteURLRewritePathFullExampleManifests)
				manifest, err := os.ReadFile(httprouteURLRewritePathFullExampleManifests)
				assert.NoError(t, err)
				assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, string(manifest)))
				cleaner.AddManifest(string(manifest))

				t.Logf("verifying that the UDPIngress routes traffic properly")

				t.Logf("asserting /dummy-random-string path is redirected (as any other path for this HTTPRoute) to /echo?msg=hello from the manifest")
				helpers.EventuallyGETPath(
					t,
					proxyURL,
					proxyURL.Host,
					"/dummy-random-string",
					http.StatusOK,
					"hello",
					nil,
					consts.IngressWait,
					consts.WaitTick,
				)

				return ctx
			}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
