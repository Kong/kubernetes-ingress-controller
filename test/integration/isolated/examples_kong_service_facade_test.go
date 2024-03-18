//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestKongServiceFacadeExample(t *testing.T) {
	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyIngress).
		WithLabel(testlabels.Kind, testlabels.KindKongServiceFacade).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(helpers.ControllerManagerOptAdditionalWatchNamespace("default")),
		)).
		WithSetup("deploy example manifest", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			manifestPath := examplesManifestPath("kong-service-facade.yaml")

			b, err := os.ReadFile(manifestPath)
			require.NoError(t, err)
			manifest := string(b)

			ingressClass := GetIngressClassFromCtx(ctx)

			t.Logf("replacing kong ingress class in yaml manifest to %s", ingressClass)
			manifest = strings.ReplaceAll(
				manifest,
				"kubernetes.io/ingress.class: kong",
				fmt.Sprintf("kubernetes.io/ingress.class: %s", ingressClass),
			)
			manifest = strings.ReplaceAll(
				manifest,
				"ingressClassName: kong",
				fmt.Sprintf("ingressClassName: %s", ingressClass),
			)

			t.Logf("applying yaml manifest %s", manifestPath)
			cluster := GetClusterFromCtx(ctx)
			require.NoError(t, clusters.ApplyManifestByYAML(ctx, cluster, manifest))

			t.Cleanup(func() {
				t.Logf("deleting yaml manifest %s", manifestPath)
				assert.NoError(t, clusters.DeleteManifestByYAML(ctx, cluster, manifest))
			})

			t.Log("waiting for httpbin deployment to be ready")
			helpers.WaitForDeploymentRollout(ctx, t, cluster, "default", "httpbin-deployment")
			return ctx
		}).
		Assess("basic-auth and key-auth plugins are applied to KongServiceFacades as expected", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			var (
				proxyURL = GetProxyURLFromCtx(ctx)

				keyAuthSecuredEndpoint = fmt.Sprintf("%s%s", proxyURL, "/alpha")
				keyAuthValidKey        = "alice-password"

				basicAuthSecuredEndpoint = fmt.Sprintf("%s%s", proxyURL, "/beta")
				basicAuthValidUsername   = "bob"
				basicAuthValidPassword   = "bob-password"
			)

			httpClient := helpers.DefaultHTTPClient()
			respondsWithExpectedStatusCode := func(t *testing.T, req *http.Request, expectedStatusCode int) {
				require.Eventually(t, func() bool {
					res, err := httpClient.Do(req)
					if err != nil {
						t.Logf("%s request returned an error: %s", req.URL, err)
						return false
					}
					defer res.Body.Close()

					if res.StatusCode != expectedStatusCode {
						t.Logf("%s request returned unexpected status code %d instead of %d", req.URL, res.StatusCode, expectedStatusCode)
						return false
					}

					return true
				}, time.Minute, 250*time.Millisecond)
			}

			newRequest := func(endpoint string, modFn func(*http.Request)) *http.Request {
				req := lo.Must(http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil))
				modFn(req)
				return req
			}

			t.Run("key-auth endpoint responses", func(t *testing.T) {
				t.Log("ensuring key-auth endpoint allows a valid key")
				validKeyAuthReq := newRequest(keyAuthSecuredEndpoint, func(r *http.Request) {
					r.Header.Set("key", keyAuthValidKey)
				})
				respondsWithExpectedStatusCode(t, validKeyAuthReq, http.StatusOK)

				t.Log("ensuring key-auth endpoint doesn't allow an invalid key")
				invalidKeyAuthReq := newRequest(keyAuthSecuredEndpoint, func(r *http.Request) {
					r.Header.Set("key", "invalid-pass")
				})
				respondsWithExpectedStatusCode(t, invalidKeyAuthReq, http.StatusUnauthorized)

				t.Log("ensuring key-auth endpoint doesn't allow valid basic-auth credentials")
				invalidKeyAuthUsingBasicAuthReq := newRequest(keyAuthSecuredEndpoint, func(r *http.Request) {
					r.SetBasicAuth(basicAuthValidUsername, basicAuthValidPassword)
				})
				respondsWithExpectedStatusCode(t, invalidKeyAuthUsingBasicAuthReq, http.StatusUnauthorized)
			})

			t.Run("basic-auth endpoint responses", func(t *testing.T) {
				t.Log("ensuring basic-auth endpoint allows valid credentials")
				validBasicAuthReq := newRequest(basicAuthSecuredEndpoint, func(r *http.Request) {
					r.SetBasicAuth(basicAuthValidUsername, basicAuthValidPassword)
				})
				respondsWithExpectedStatusCode(t, validBasicAuthReq, http.StatusOK)

				t.Log("ensuring basic-auth endpoint doesn't allow invalid credentials")
				invalidBasicAuthReq := newRequest(basicAuthSecuredEndpoint, func(r *http.Request) {
					r.SetBasicAuth(basicAuthValidUsername, "invalid-pass")
				})
				respondsWithExpectedStatusCode(t, invalidBasicAuthReq, http.StatusUnauthorized)

				t.Log("ensuring basic-auth endpoint doesn't allow a valid key-auth key")
				invalidBasicAuthUsingKeyAuthReq := newRequest(basicAuthSecuredEndpoint, func(r *http.Request) {
					r.Header.Set("key", keyAuthValidKey)
				})
				respondsWithExpectedStatusCode(t, invalidBasicAuthUsingKeyAuthReq, http.StatusUnauthorized)
			})

			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
