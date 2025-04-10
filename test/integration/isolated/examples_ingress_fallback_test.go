//go:build integration_tests

package isolated

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestIngressWithBrokenPluginFallback(t *testing.T) {
	ingressWithBrokenPluginFallback := examplesManifestPath("ingress-broken-plugin-fallback.yaml")

	replaceIngressClassAnnotationInManifests := func(manifests string, ingressClass string) string {
		return strings.ReplaceAll(manifests, "ingressClassName: kong", fmt.Sprintf("ingressClassName: %s", ingressClass))
	}

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		WithLabel(testlabels.Kind, testlabels.KindIngress).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(
				helpers.ControllerManagerOptAdditionalWatchNamespace("default"),
			),
			withControllerManagerFeatureGates(map[string]string{managercfg.FallbackConfigurationFeature: "true"}),
		)).
		Assess("deploying to cluster works and HTTP traffic is routed to the service",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				cluster := GetClusterFromCtx(ctx)
				proxyURL := GetHTTPURLFromCtx(ctx)

				t.Logf("applying yaml manifest %s", ingressWithBrokenPluginFallback)
				b, err := os.ReadFile(ingressWithBrokenPluginFallback)
				assert.NoError(t, err)
				manifest := replaceIngressClassAnnotationInManifests(string(b), GetIngressClassFromCtx(ctx))
				assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, manifest))
				cleaner.AddManifest(manifest)

				t.Logf("verifying that the Ingress routes traffic properly to the /ingress-testing path")
				helpers.EventuallyGETPath(
					t,
					proxyURL,
					proxyURL.Host,
					"/ingress-testing",
					nil,
					http.StatusOK,
					"Running on Pod",
					nil,
					consts.IngressWait,
					consts.WaitTick,
				)

				return ctx
			}).
		Assess("verify that valid Ingress's status contains LoadBalancer information", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)
			// Verify Ingress status contains LoadBalancer information
			require.Eventually(t, func() bool {
				ingress, err := cluster.Client().NetworkingV1().Ingresses("default").Get(ctx, "working-ingress", metav1.GetOptions{})
				if err != nil {
					t.Logf("Failed to get Ingress: %v", err)
					return false
				}
				lbStatus := ingress.Status.LoadBalancer.Ingress
				return len(lbStatus) > 0
			}, consts.IngressWait, consts.WaitTick)
			return ctx
		}).
		Assess("verify that route with misconfigured plugin is not operational", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			proxyURL := GetHTTPURLFromCtx(ctx)
			t.Logf("verifying that Kong gateway response is returned instead of desired site")

			helpers.EventuallyGETPath(
				t,
				proxyURL,
				proxyURL.Host,
				"/for-auth-users",
				nil,
				http.StatusNotFound,
				"no Route matched with those values",
				nil,
				consts.IngressWait,
				consts.WaitTick,
			)
			return ctx
		}).
		Assess("verify that invalid Ingress's status doesn't contain LoadBalancer information", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cluster := GetClusterFromCtx(ctx)
			// Verify Ingress status contains LoadBalancer information
			require.Eventually(t, func() bool {
				ingress, err := cluster.Client().NetworkingV1().Ingresses("default").Get(ctx, "for-auth", metav1.GetOptions{})
				if err != nil {
					t.Logf("Failed to get Ingress: %v", err)
					return false
				}
				lbStatus := ingress.Status.LoadBalancer.Ingress
				return len(lbStatus) == 0
			}, consts.IngressWait, consts.WaitTick)
			return ctx
		}).
		Assess("verify diagnostic server fallback info", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			diagURL := GetDiagURLFromCtx(ctx)
			t.Logf("verifying diag available")
			require.EventuallyWithT(t, func(c *assert.CollectT) {
				cl := helpers.DefaultHTTPClient()
				resp, err := cl.Do(helpers.MustHTTPRequest(t, http.MethodGet, diagURL.Host, "/debug/config/fallback", nil))
				if !assert.NoError(c, err) {
					return
				}
				defer resp.Body.Close()

				if !assert.Equal(c, http.StatusOK, resp.StatusCode) {
					return
				}

				response := diagnostics.FallbackResponse{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				if !assert.NoError(c, err) {
					return
				}
				assert.Equal(t, response.Status, diagnostics.FallbackStatusTriggered)
				assert.NotEmpty(t, response.ExcludedObjects)
				contains := lo.ContainsBy(response.ExcludedObjects, func(obj diagnostics.FallbackAffectedObjectMeta) bool {
					return obj.Group == "configuration.konghq.com" && obj.Kind == "KongPlugin" && obj.Name == "key-auth"
				})
				assert.Truef(t, contains, "expected to find KongPlugin key-auth in excluded objects, got: %v", response.ExcludedObjects)
			}, consts.IngressWait, consts.WaitTick)
			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
