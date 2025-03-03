//go:build integration_tests

package isolated

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"os"
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

func TestExampleTLSRoute(t *testing.T) {
	tlsrouteExampleManifests := examplesManifestPath("gateway-tlsroute.yaml")

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindTLSRoute).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		Assess("deploying to cluster works and tls traffic is routed to the service",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				cluster := GetClusterFromCtx(ctx)
				proxyTLSURL := GetTLSURLFromCtx(ctx)

				t.Logf("applying yaml manifest %s", tlsrouteExampleManifests)
				b, err := os.ReadFile(tlsrouteExampleManifests)
				assert.NoError(t, err)
				assert.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, string(b)))
				cleaner.AddManifest(string(b))
				// Copy pasted cert from gateway-tlsroute.yaml to use it in certPool for checking validity.
				const caForCert = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVMekNDQXhlZ0F3SUJBZ0lVVkdBQWlrd3Fid3VIRFBpd092a1hwM0hpUHlRd0RRWUpLb1pJaHZjTkFRRUwKQlFBd2daQXhDekFKQmdOVkJBWVRBbFZUTVJNd0VRWURWUVFJREFwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSApEQTFUWVc0Z1JuSmhibU5wYzJOdk1SSXdFQVlEVlFRS0RBbExiMjVuSUVsdVl5NHhHREFXQmdOVkJBc01EMVJsCllXMGdTM1ZpWlhKdVpYUmxjekVtTUNRR0ExVUVBd3dkWlhoaGJYQnNaUzEwYkhOeWIzVjBaUzVyYjI1bkxtVjQKWVcxd2JHVXdJQmNOTWpVd016QXpNVGt6T0RBMVdoZ1BNakV5TlRBeU1EY3hPVE00TURWYU1JR1FNUXN3Q1FZRApWUVFHRXdKVlV6RVRNQkVHQTFVRUNBd0tRMkZzYVdadmNtNXBZVEVXTUJRR0ExVUVCd3dOVTJGdUlFWnlZVzVqCmFYTmpiekVTTUJBR0ExVUVDZ3dKUzI5dVp5QkpibU11TVJnd0ZnWURWUVFMREE5VVpXRnRJRXQxWW1WeWJtVjAKWlhNeEpqQWtCZ05WQkFNTUhXVjRZVzF3YkdVdGRHeHpjbTkxZEdVdWEyOXVaeTVsZUdGdGNHeGxNSUlCSWpBTgpCZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUE2QS9kMGRBMVZTM2RQZGV0emRIaGlSRjFwWGtWCnJmM1IzbTBucnUxa05Kd1d1TDdKYnBTd0MraVRvNXRINDYyNzVzdWFQRXpBYm5JZTVsWWVrdGJpOTFLeU9XUGcKQkw0emZiZCtSNXhna3JiTWU4dy9mbGZIeHhwZ0I2d0xXOHdoQS9Ec3hNUUhhVThGY2JSMTZCd3M3czVLNm9YcgpDbngwMEtnMEo1SHRUeWhoSStsQjJRQzJKZXpqbEQxaGkxYjk1ekg3S2Y0bTVydlFYMmxtSndVVm9VQkNlQnNaClRkUzl1bUdBSVcxVDEwczJDaWdHbjBmUVFFR0RCUnFScDJ4NmNQVDhicXNndEZDcTBNZG1qZlBabW9HVmlsNUgKNlh6VnVKdFBZSFJOV0RPK2dydFBBbDhwOUUyM2VtY2gzTUZwN0RVYUFmZGNwZE45czZLOWptZmNjd0lEQVFBQgpvMzB3ZXpBZEJnTlZIUTRFRmdRVW9Jb05BN1BpUExsSTYrRU93Wk4wVTdmVFlrNHdId1lEVlIwakJCZ3dGb0FVCm9Jb05BN1BpUExsSTYrRU93Wk4wVTdmVFlrNHdEd1lEVlIwVEFRSC9CQVV3QXdFQi96QW9CZ05WSFJFRUlUQWYKZ2gxbGVHRnRjR3hsTFhSc2MzSnZkWFJsTG10dmJtY3VaWGhoYlhCc1pUQU5CZ2txaGtpRzl3MEJBUXNGQUFPQwpBUUVBZzZlQ0JaMDRYRVRzd1o2WUkwWjRBZHk5ZEJJRU5NcHBVK050Njh4VWZqR2RuckRBNVByUzBmcUFXcGswClQ1Wm9PUnovZVNnZDZJaUJWZU4yYlJ3VFJIT092dlJWVC9JQ0d5anhjYUJFT0gzVDBoSUwzMDhPbTEvTUZMQnEKNm5wZ0lpcEwzN2FPL015cmFjcUVqOGtyTHJKeitIdlFOUGRjU2JJdE5Pd3RMaFE5VjJhbG53cGNwQlY5bnFDbgpLUllMc05jTTBkaVZ4a2tuTC96Z0oxQy9jdWJQeGlXRVpRY05FV3ppN0cyLzRwdEdqeG1DT2lLUm91c0l5ZU1RCldLS1Bidlg1dExBY29kR1ZTTEFqVGhGYlVFUnF3eHdxOHJ5anhmQlVWdHZzUngrRTRacDZYVitpdWNRdFNNRHAKVmhKZHVGOUliZllnSzVOMlY2TGVJZkplUXc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
				certPool := x509.NewCertPool()
				decodedCert, err := base64.StdEncoding.DecodeString(caForCert)
				assert.NoError(t, err)
				assert.True(t, certPool.AppendCertsFromPEM(decodedCert))

				t.Log("verifying that TLSRoute becomes routable")
				tlsOpt := test.WithTLSOption("example-tlsroute.kong.example", certPool, true) // URL as in gateway-tlsroute.yaml.
				assert.EventuallyWithT(t, func(c *assert.CollectT) {
					err := test.EchoResponds(
						test.ProtocolTLS, proxyTLSURL, "example-tlsroute-manifest", tlsOpt,
					)
					assert.NoError(c, err)
				}, consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
