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

func TestHTTPRouteExample(t *testing.T) {
	httprouteExampleManifest := examplesManifestPath("gateway-httproute.yaml")

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindHTTPRoute).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		Assess("deploying to cluster works and HTTP requests are routed properly",
			runHTTPRouteExampleTestScenario(httprouteExampleManifest),
		).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}

func TestHTTPRouteWithBrokenPluginFallback(t *testing.T) {
	httprouteWithBrokenPluginFallback := examplesManifestPath("gateway-httproute-broken-plugin-fallback.yaml")

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindHTTPRoute).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(
				helpers.ControllerManagerOptAdditionalWatchNamespace("default"),
			),
			withControllerManagerFeatureGates(map[string]string{"FallbackConfiguration": "true"}),
		)).
		Assess("deploying to cluster works and HTTP requests are routed properly",
			runHTTPRouteExampleTestScenario(httprouteWithBrokenPluginFallback),
		).
		Assess("verify that route with misconfigured plugin is not operational", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			proxyURL := GetHTTPURLFromCtx(ctx)
			t.Logf("verifying that Kong gateway response in returned instead of desired site")
			helpers.EventuallyGETPath(
				t,
				proxyURL,
				proxyURL.Host,
				"/for-auth-users",
				http.StatusNotFound,
				"no Route matched with those values",
				nil,
				consts.IngressWait,
				consts.WaitTick,
			)
			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}

func runHTTPRouteExampleTestScenario(manifestToUse string) func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
	return func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
		cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
		cluster := GetClusterFromCtx(ctx)
		proxyURL := GetHTTPURLFromCtx(ctx)

		t.Logf("applying yaml manifest %s", manifestToUse)
		manifest, err := os.ReadFile(manifestToUse)
		assert.NoError(t, err)
		assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, string(manifest)))
		cleaner.AddManifest(string(manifest))

		t.Logf("verifying that traffic is routed properly")

		t.Logf("verifying that the HTTPRoute becomes routable")
		helpers.EventuallyGETPath(
			t,
			proxyURL,
			proxyURL.Host,
			"/httproute-testing",
			http.StatusOK,
			"echo-1",
			nil,
			consts.IngressWait,
			consts.WaitTick,
		)

		t.Logf("verifying that the backendRefs are being loadbalanced")
		helpers.EventuallyGETPath(
			t,
			proxyURL,
			proxyURL.Host,
			"/httproute-testing",
			http.StatusOK,
			"echo-2",
			nil,
			consts.IngressWait,
			consts.WaitTick,
		)

		return ctx
	}
}
