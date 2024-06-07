//go:build integration_tests

package isolated

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager/featuregates"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v3/pkg/clientset"
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
			withControllerManagerFeatureGates(map[string]string{featuregates.FallbackConfiguration: "true"}),
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

func TestHTTPRouteUseLastValidConfigWithBrokenPluginFallback(t *testing.T) {
	httprouteExampleManifest := examplesManifestPath("gateway-httproute.yaml")
	const (
		namespace                    = "default"
		additionalRouteName          = "httproute-testing-additional"
		additionalRoutePath          = "/additional-route"
		additionalRouteServiceTarget = "echo-1"

		additionalHeaderKey   = "X-Additional-Header"
		additionalHeaderValue = "additional-header-value"
	)
	testAdditionalRoute := func(t *testing.T, proxyURL *url.URL) {
		t.Helper()
		t.Log("verifying that routing to additional route works and header added by plugin is returned")
		helpers.EventuallyGETPath(
			t,
			proxyURL,
			proxyURL.Host,
			additionalRoutePath,
			http.StatusOK,
			additionalRouteServiceTarget,
			map[string]string{
				additionalHeaderKey: additionalHeaderValue,
			},
			consts.IngressWait,
			consts.WaitTick,
		)
	}

	f := features.
		New("example").
		WithLabel(testlabels.Example, testlabels.ExampleTrue).
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindHTTPRoute).
		WithSetup("deploy kong addon into cluster", featureSetup(
			withControllerManagerOpts(
				helpers.ControllerManagerOptAdditionalWatchNamespace(namespace),
				helpers.ControllerManagerOptFlagUseLastValidConfigForFallback(),
			),
			withControllerManagerFeatureGates(map[string]string{featuregates.FallbackConfiguration: "true"}),
		)).
		Assess("deploying to cluster works and HTTP requests are routed properly", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			runHTTPRouteExampleTestScenario(httprouteExampleManifest)(ctx, t, c)

			clusterCfg := GetClusterFromCtx(ctx).Config()
			t.Log("getting a gateway client")
			gatewayClient, err := gatewayclient.NewForConfig(clusterCfg)
			assert.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, gatewayClient)

			t.Log("adding additional properly configured route with plugin")
			_, err = gatewayClient.GatewayV1().HTTPRoutes(namespace).Create(ctx, &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      additionalRouteName,
					Annotations: map[string]string{
						annotations.AnnotationPrefix + annotations.StripPathKey: "true",
					},
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{
							{
								Name: "kong",
							},
						},
					},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								{
									Path: &gatewayapi.HTTPPathMatch{
										Type:  lo.ToPtr(gatewayapi.PathMatchPathPrefix),
										Value: lo.ToPtr(additionalRoutePath),
									},
								},
							},
							BackendRefs: []gatewayapi.HTTPBackendRef{
								{
									BackendRef: gatewayapi.BackendRef{
										BackendObjectReference: gatewayapi.BackendObjectReference{
											Kind: lo.ToPtr(gatewayapi.Kind("Service")),
											Name: additionalRouteServiceTarget,
											Port: lo.ToPtr(gatewayapi.PortNumber(80)),
										},
									},
								},
							},
						},
					},
				},
			}, metav1.CreateOptions{})
			assert.NoError(t, err)
			client, err := clientset.NewForConfig(clusterCfg)
			require.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, client)

			workingPlugin := &kongv1.KongPlugin{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "response-transformer",
				},
				PluginName: "response-transformer",
				// Misconfigured on purpose.
				Config: apiextensionsv1.JSON{
					Raw: []byte(fmt.Sprintf(`
						{
							"config": {"add": {"headers": ["%s:%s"]}}
						}`, additionalHeaderKey, additionalHeaderValue),
					),
				},
			}
			_, err = client.ConfigurationV1().KongPlugins(namespace).Create(ctx, workingPlugin, metav1.CreateOptions{})
			require.NoError(t, err)

			t.Logf("verifying that routing to %s works and header added by plugin is returned", additionalRoutePath)
			testAdditionalRoute(t, GetHTTPURLFromCtx(ctx))

			return ctx
		}).
		Assess("attach broken plugin to a working route", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			client := GetFromCtxForT[*clientset.Clientset](ctx, t)

			plugin, err := client.ConfigurationV1().KongPlugins(namespace).Get(ctx, "response-transformer", metav1.GetOptions{})
			require.NoError(t, err)
			plugin.Config = apiextensionsv1.JSON{
				Raw: []byte(`{"test": "test"}`),
			}
			_, err = client.ConfigurationV1().KongPlugins(namespace).Update(ctx, plugin, metav1.UpdateOptions{})
			require.NoError(t, err)

			t.Log("getting a gateway client")
			gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			assert.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, gatewayClient)

			route, err := gatewayClient.GatewayV1().HTTPRoutes(namespace).Get(ctx, additionalRouteName, metav1.GetOptions{})
			assert.NoError(t, err)
			route.Annotations[annotations.AnnotationPrefix+annotations.PluginsKey] = plugin.Name
			_, err = gatewayClient.GatewayV1().HTTPRoutes(namespace).Update(ctx, route, metav1.UpdateOptions{})
			assert.NoError(t, err)

			return ctx
		}).
		Assess("verify that route with misconfigured plugin operates with previous config", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			testAdditionalRoute(t, GetHTTPURLFromCtx(ctx))
			return ctx
		}).
		Assess("modify working route /httproute-testing to /new-route", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			gatewayclient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			route, err := gatewayclient.GatewayV1().HTTPRoutes(namespace).Get(ctx, "httproute-testing", metav1.GetOptions{})
			require.NoError(t, err)
			route.Spec.Rules[0].Matches[0].Path.Value = lo.ToPtr("/new-route")
			_, err = gatewayclient.GatewayV1().HTTPRoutes(namespace).Update(ctx, route, metav1.UpdateOptions{})
			require.NoError(t, err)

			return ctx
		}).
		Assess("verify that all routes are operational", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			proxyURL := GetHTTPURLFromCtx(ctx)
			testAdditionalRoute(t, proxyURL)

			const newRoute = "/new-route"
			t.Logf("verifying that /httproute-testing is no longer operational")
			helpers.EventuallyGETPath(
				t,
				proxyURL,
				proxyURL.Host,
				"/httproute-testing",
				http.StatusNotFound,
				"no Route matched with those values",
				nil,
				consts.IngressWait,
				consts.WaitTick,
			)

			t.Logf("verifying that /new-route is operational and loadbalanced")
			helpers.EventuallyGETPath(
				t,
				proxyURL,
				proxyURL.Host,
				newRoute,
				http.StatusOK,
				"echo-1",
				nil,
				consts.IngressWait,
				consts.WaitTick,
			)
			helpers.EventuallyGETPath(
				t,
				proxyURL,
				proxyURL.Host,
				newRoute,
				http.StatusOK,
				"echo-2",
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
