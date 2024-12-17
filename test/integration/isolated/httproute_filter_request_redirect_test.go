//go:build integration_tests

package isolated

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestHTTPRouteFilterRequestRedirect(t *testing.T) {
	f := features.
		New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindHTTPRoute).
		WithSetup("deploy Kong addon", featureSetup(
			withKongProxyEnvVars(map[string]string{
				"PROXY_LISTEN": `0.0.0.0:8000 http2\, 0.0.0.0:8443 http2 ssl`,
			}),
		)).
		WithSetup("deploying a gateway and a backend `httpbin` service", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			cluster := GetClusterFromCtx(ctx)
			namespace := GetNamespaceForT(ctx, t)

			t.Log("getting a gateway client")
			gatewayClient, err := gatewayclient.NewForConfig(cluster.Config())
			assert.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, gatewayClient)

			t.Log("deploying a new gatewayClass")
			gatewayClassName := uuid.NewString()
			gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
			assert.NoError(t, err)
			cleaner.Add(gwc)

			t.Log("deploying a new gateway")
			gateway, err := helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName, func(gw *gatewayapi.Gateway) {
				gw.Spec.Listeners = builder.NewListener("http").
					HTTP().
					WithPort(ktfkong.DefaultProxyHTTPPort).
					IntoSlice()
			})
			assert.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, gateway)
			cleaner.Add(gateway)

			return ctx
		}).
		Assess("Create an HTTPRoute redirecting to the httpbin service and verify if the request is redirected", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			gateway := GetFromCtxForT[*gatewayapi.Gateway](ctx, t)
			GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			namespace := GetNamespaceForT(ctx, t)
			gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			proxyHTTPURL := GetHTTPURLFromCtx(ctx)
			proxyAdminURL := GetAdminURLFromCtx(ctx)

			httpRoute := &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: "httproute-request-redirect-test",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: gatewayapi.ObjectName(gateway.Name),
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{
						{
							Matches: []gatewayapi.HTTPRouteMatch{
								builder.NewHTTPRouteMatch().WithPathPrefix("/").Build(),
							},
							Filters: []gatewayapi.HTTPRouteFilter{
								builder.NewHTTPRouteRequestRedirectFilter().
									WithRequestRedirectStatusCode(http.StatusMovedPermanently).
									WithRequestRedirectHost("example.com").WithRequestRedirectPathModifier(
									gatewayapi.FullPathHTTPPathModifier,
									"/redirect-target",
								).Build(),
							},
						},
					},
				},
			}

			httpRoute, err := gatewayClient.GatewayV1().HTTPRoutes(namespace).Create(ctx, httpRoute, metav1.CreateOptions{})
			assert.NoError(t, err)

			t.Logf("Verify that replaceFullPath works")
			httpClientIgnoreRedirect := &http.Client{
				CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			assert.Eventually(
				t, func() bool {
					resp, err := httpClientIgnoreRedirect.Get(proxyHTTPURL.String())
					if err != nil {
						return false
					}
					defer resp.Body.Close()
					return resp.StatusCode == http.StatusMovedPermanently &&
						resp.Header.Get("Location") == "http://example.com/redirect-target"
				},
				consts.IngressWait, consts.WaitTick,
			)

			redirectPluginVersionRange, err := kong.NewRange(">=3.9.0")
			assert.NoError(t, err)
			kongVersion, err := helpers.GetKongVersion(ctx, proxyAdminURL, "password")
			assert.NoError(t, err)
			supportRedirectPlugin := redirectPluginVersionRange(kongVersion)
			if !supportRedirectPlugin {
				t.Logf("Skipping case for preserving path in request redirect because Kong version is %s which does not support redirect plugin", kongVersion)
				return ctx
			}

			t.Log("Update RequestRedirect filter of HTTPRoute to test preserving path in request redirect")
			httpRoute, err = gatewayClient.GatewayV1().HTTPRoutes(namespace).Get(ctx, httpRoute.Name, metav1.GetOptions{})
			assert.NoError(t, err)
			httpRoute.Spec.Rules = []gatewayapi.HTTPRouteRule{
				{
					Matches: []gatewayapi.HTTPRouteMatch{
						builder.NewHTTPRouteMatch().WithPathPrefix("/").Build(),
					},
					Filters: []gatewayapi.HTTPRouteFilter{
						builder.NewHTTPRouteRequestRedirectFilter().
							WithRequestRedirectStatusCode(http.StatusMovedPermanently).
							WithRequestRedirectScheme("http").
							WithRequestRedirectHost("example.com").Build(),
					},
				},
			}
			_, err = gatewayClient.GatewayV1().HTTPRoutes(namespace).Update(ctx, httpRoute, metav1.UpdateOptions{})
			assert.NoError(t, err)
			t.Logf("Verify that no path modifier preserves path in the request")
			assert.Eventually(
				t, func() bool {
					resp, err := httpClientIgnoreRedirect.Get(proxyHTTPURL.String() + "/redirect-target")
					if err != nil {
						return false
					}
					defer resp.Body.Close()
					return resp.StatusCode == http.StatusMovedPermanently &&
						resp.Header.Get("Location") == "http://example.com/redirect-target"
				},
				consts.IngressWait, consts.WaitTick,
			)

			return ctx
		}).Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
