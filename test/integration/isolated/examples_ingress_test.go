//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"net/http"
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

func TestIngressExample(t *testing.T) {
	ingressExampleManifests := examplesManifestPath("ingress.yaml")

	replaceIngressClassAnnotationInManifests := func(manifests string, ingressClass string) string {
		return strings.ReplaceAll(manifests, "ingressClassName: kong", fmt.Sprintf("ingressClassName: %s", ingressClass))
	}

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		WithLabel(testlabels.Kind, testlabels.KindIngress).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		Assess("deploying to cluster works and HTTP traffic is routed to the service",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				cluster := GetClusterFromCtx(ctx)
				proxyURL := GetHTTPURLFromCtx(ctx)

				t.Logf("applying yaml manifest %s", ingressExampleManifests)
				b, err := os.ReadFile(ingressExampleManifests)
				assert.NoError(t, err)
				manifest := replaceIngressClassAnnotationInManifests(string(b), GetIngressClassFromCtx(ctx))
				assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, manifest))
				cleaner.AddManifest(manifest)

				t.Logf("verifying that the Ingress routes traffic properly")
				helpers.EventuallyGETPath(
					t,
					proxyURL,
					proxyURL.Host,
					"/",
					nil,
					http.StatusOK,
					"<title>httpbin.org</title>",
					nil,
					consts.IngressWait,
					consts.WaitTick,
				)

				return ctx
			}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
