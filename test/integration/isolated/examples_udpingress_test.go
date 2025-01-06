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

	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestUDPIngressExample(t *testing.T) {
	udpIngressExampleManifests := examplesManifestPath("udpingress.yaml")

	replaceIngressClassAnnotationInManifests := func(manifests string, ingressClass string) string {
		const ingressClassTemplate = `kubernetes.io/ingress.class: "%s"`
		return strings.ReplaceAll(manifests, fmt.Sprintf(ingressClassTemplate, "kong"), fmt.Sprintf(ingressClassTemplate, ingressClass))
	}

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		WithLabel(testlabels.Kind, testlabels.KindKongUDPIngress).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		Assess("deploying to cluster works and UDP traffic is routed to the service",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				cluster := GetClusterFromCtx(ctx)
				proxyUDPURL := GetUDPURLFromCtx(ctx)

				t.Logf("applying yaml manifest %s", udpIngressExampleManifests)
				b, err := os.ReadFile(udpIngressExampleManifests)
				assert.NoError(t, err)
				manifest := replaceIngressClassAnnotationInManifests(string(b), GetIngressClassFromCtx(ctx))
				assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, manifest))
				cleaner.AddManifest(manifest)

				t.Logf("verifying that the UDPIngress routes traffic properly")
				assert.EventuallyWithT(t, func(c *assert.CollectT) {
					assert.NoError(
						c, test.EchoResponds(test.ProtocolUDP, proxyUDPURL, "udpingress-example-manifest"),
					)
				}, consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
