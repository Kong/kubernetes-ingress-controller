//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

var (
	emptyHeaderSet     = make(map[string]string)
	kongDebugHeaderSet = map[string]string{"kong-debug": "1"}
)

func TestHTTPRouteEssentials(t *testing.T) {
	ns, cleaner := setup(t)
	defer func() {
		if true {
			output, err := cleaner.DumpDiagnostics(ctx, t.Name())
			t.Logf("%s failed, dumped diagnostics to %s", t.Name(), output)
			assert.NoError(t, err)
		}
		assert.NoError(t, cleaner.Cleanup(ctx))
	}()

	t.Log("getting a gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a new gatewayClass")
	gatewayClassName := uuid.NewString()
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("deploying a new gateway")
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container1 := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	container1.Env = []corev1.EnvVar{
		{
			Name:  "GUNICORN_CMD_ARGS",
			Value: "--capture-output --error-logfile - --access-logfile - --access-logformat '%(h)s %(t)s %(r)s %(s)s Host: %({Host}i)s}'",
		},
	}
	deployment1 := generators.NewDeploymentForContainer(container1)
	deployment1, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment1, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("deploying an extra minimal HTTP container deployment to test multiple backendRefs")
	container2 := generators.NewContainer("nginx", "nginx", 80)
	deployment2 := generators.NewDeploymentForContainer(container2)
	deployment2, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment2, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment1)
	cleaner.Add(deployment2)

	t.Logf("exposing deployment %s via service", deployment1.Name)
	service1 := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service1, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment2.Name)
	service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service2, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service1)
	cleaner.Add(service2)

	kongplugin := &kongv1.KongPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns.Name,
			Name:      "correlation",
		},
		PluginName: "correlation-id",
		Config: apiextensionsv1.JSON{
			Raw: []byte(`{"header_name":"reqid", "echo_downstream": true}`),
		},
	}
	require.NoError(t, err)
	pluginClient, err := clientset.NewForConfig(env.Cluster().Config())
	_, err = pluginClient.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, kongplugin, metav1.CreateOptions{})

	t.Logf("creating an httproute to access deployment %s via kong", deployment1.Name)
	httpPort := gatewayv1beta1.PortNumber(80)
	pathMatchPrefix := gatewayv1beta1.PathMatchPathPrefix
	pathMatchRegularExpression := gatewayv1beta1.PathMatchRegularExpression
	pathMatchExact := gatewayv1beta1.PathMatchExact
	headerMatchRegex := gatewayv1beta1.HeaderMatchRegularExpression
	httpRoute := &gatewayv1beta1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.StripPathKey: "true",
				annotations.AnnotationPrefix + annotations.PluginsKey:   "correlation",
			},
		},
		Spec: gatewayv1beta1.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
				ParentRefs: []gatewayv1beta1.ParentReference{{
					Name: gatewayv1beta1.ObjectName(gateway.Name),
				}},
			},
			Rules: []gatewayv1beta1.HTTPRouteRule{{
				Matches: []gatewayv1beta1.HTTPRouteMatch{
					{
						Path: &gatewayv1beta1.HTTPPathMatch{
							Type:  &pathMatchPrefix,
							Value: kong.String("/test-http-route-essentials"),
						},
					},
					{
						Path: &gatewayv1beta1.HTTPPathMatch{
							Type:  &pathMatchRegularExpression,
							Value: kong.String(`/2/test-http-route-essentials/regex/\d{3}`),
						},
					},
					{
						Path: &gatewayv1beta1.HTTPPathMatch{
							Type:  &pathMatchExact,
							Value: kong.String(`/3/exact-test-http-route-essentials`),
						},
					},
				},
				BackendRefs: []gatewayv1beta1.HTTPBackendRef{{
					BackendRef: gatewayv1beta1.BackendRef{
						BackendObjectReference: gatewayv1beta1.BackendObjectReference{
							Name: gatewayv1beta1.ObjectName(service1.Name),
							Port: &httpPort,
							Kind: util.StringToGatewayAPIKindV1Beta1Ptr("Service"),
						},
					},
				}},
			}},
		},
	}
	if versions.GetKongVersion().MajorMinorOnly().GTE(versions.RegexHeaderVersionCutoff) {
		httpRoute.Spec.Rules[0].Matches = append(httpRoute.Spec.Rules[0].Matches, gatewayv1beta1.HTTPRouteMatch{
			Headers: []gatewayv1beta1.HTTPHeaderMatch{
				{
					Type:  &headerMatchRegex,
					Value: "^audio/*",
					Name:  "Content-Type",
				},
			},
		})
	}
	httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Create(ctx, httpRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(httpRoute)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	defer func() {
		req := newRequest(t, http.MethodGet, "2/test-http-route-essentials/regex/999", kongDebugHeaderSet)
		resp, err := httpc.Do(req)
		if err != nil {
			t.Logf("WARNING: http request failed for GET %s/%s: %v", proxyURL, "999", err)
		}
		var head string
		for h, v := range resp.Header {
			head = fmt.Sprintf("%s\n%s: %s", head, h, strings.Join(v, ","))
		}
		defer resp.Body.Close()
		b := new(bytes.Buffer)
		_, err = b.ReadFrom(resp.Body)
		if err != nil {
			t.Logf("WARNING: http request unreadable for GET %s/%s: %v", proxyURL, "999", err)
		}
		t.Logf("TestHTTPRouteEssentials diag: status %v headers===\n%s\nbody===\n%s\n\n", resp.StatusCode, head, b.String())
	}()
	t.Log("waiting for routes from HTTPRoute to become operational")
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "test-http-route-essentials/base64/wqt5b8q7ccK7IGRhbiBib3NocWEgYmlyIGphdm9iaW1peiB5b8q7cWRpci4K",
		http.StatusOK, "«yoʻq» dan boshqa bir javobimiz yoʻqdir.", emptyHeaderSet)
	eventuallyGETPath(t, "2/test-http-route-essentials/regex/999", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "3/exact-test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "3/exact-test-http-route-essentialsNO", http.StatusNotFound, "no Route matched", emptyHeaderSet)

	require.Eventually(t, func() bool {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", proxyURL, "test-http-route-essentials"), nil)
		if err != nil {
			t.Logf("WARNING: failed to create HTTP request: %v", err)
			return false
		}
		resp, err := httpc.Do(req)
		if err != nil {
			t.Logf("WARNING: http request failed for GET %s/%s: %v", proxyURL, "test-http-route-essentials", err)
			return false
		}
		defer resp.Body.Close()
		if _, ok := resp.Header["Reqid"]; ok {
			return true
		}
		return false
	}, ingressWait, waitTick)

	if versions.GetKongVersion().MajorMinorOnly().GTE(versions.RegexHeaderVersionCutoff) {
		t.Log("verifying HTTPRoute header match")
		eventuallyGETPath(t, "", http.StatusOK, "<title>httpbin.org</title>", map[string]string{"Content-Type": "audio/mp3"})
	}

	t.Log("adding an additional backendRef to the HTTPRoute")
	var httpbinWeight int32 = 75
	var nginxWeight int32 = 25
	require.Eventually(t, func() bool {
		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(httpRoute.Namespace).Get(ctx, httpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httpRoute.Spec.Rules[0].BackendRefs = []gatewayv1beta1.HTTPBackendRef{
			{
				BackendRef: gatewayv1beta1.BackendRef{
					BackendObjectReference: gatewayv1beta1.BackendObjectReference{
						Name: gatewayv1beta1.ObjectName(service1.Name),
						Port: &httpPort,
						Kind: util.StringToGatewayAPIKindV1Beta1Ptr("Service"),
					},
					Weight: &httpbinWeight,
				},
			},
			{
				BackendRef: gatewayv1beta1.BackendRef{
					BackendObjectReference: gatewayv1beta1.BackendObjectReference{
						Name: gatewayv1beta1.ObjectName(service2.Name),
						Port: &httpPort,
						Kind: util.StringToGatewayAPIKindV1Beta1Ptr("Service"),
					},
					Weight: &nginxWeight,
				},
			},
		}

		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(httpRoute.Namespace).Update(ctx, httpRoute, metav1.UpdateOptions{})
		if err != nil {
			t.Logf("WARNING: could not update httproute with an additional backendRef: %s (retrying)", err)
			return false
		}

		return true
	}, ingressWait, waitTick)

	t.Log("verifying that both backends are ready to receive traffic")
	httpbinRespContent := "<title>httpbin.org</title>"
	nginxRespContent := "<title>Welcome to nginx!</title>"
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusOK, httpbinRespContent, emptyHeaderSet)
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusOK, nginxRespContent, emptyHeaderSet)

	t.Log("verifying that both backends receive requests according to weighted distribution")
	httpbinRespName := "httpbin-resp"
	nginxRespName := "nginx-resp"
	toleranceDelta := 0.2
	expectedRespRatio := map[string]int{
		httpbinRespName: int(httpbinWeight),
		nginxRespName:   int(nginxWeight),
	}
	weightedLoadBalancingTestConfig := countHTTPResponsesConfig{
		Method:      http.MethodGet,
		Path:        "test-http-route-essentials",
		Headers:     emptyHeaderSet,
		Duration:    5 * time.Second,
		RequestTick: 50 * time.Millisecond,
	}
	respCounter := countHTTPGetResponses(t, weightedLoadBalancingTestConfig,
		matchRespByStatusAndContent(httpbinRespName, http.StatusOK, httpbinRespContent),
		matchRespByStatusAndContent(nginxRespName, http.StatusOK, nginxRespContent),
	)
	assert.InDeltaMapValues(t,
		distributionOfMapValues(respCounter),
		distributionOfMapValues(expectedRespRatio),
		toleranceDelta,
		"Response distribution does not match expected distribution within %f%% delta,"+
			" request-count=%v, expected-ratio=%v",
		toleranceDelta*100, respCounter, expectedRespRatio,
	)

	t.Log("removing the parentrefs from the HTTPRoute")
	oldParentRefs := httpRoute.Spec.ParentRefs
	require.Eventually(t, func() bool {
		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Get(ctx, httpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httpRoute.Spec.ParentRefs = nil
		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Update(ctx, httpRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the parentRefs now removed")
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusNotFound, "", emptyHeaderSet)

	t.Log("putting the parentRefs back")
	require.Eventually(t, func() bool {
		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Get(ctx, httpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httpRoute.Spec.ParentRefs = oldParentRefs
		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Update(ctx, httpRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("deleting the GatewayClass")
	require.NoError(t, gatewayClient.GatewayV1alpha2().GatewayClasses().Delete(ctx, gatewayClassName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)
	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the GatewayClass now removed")
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusNotFound, "", emptyHeaderSet)

	t.Log("putting the GatewayClass back")
	gwc, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of HTTPRoutes and the route becomes available again")
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("deleting the Gateway")
	require.NoError(t, gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the Gateway now removed")
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusNotFound, "", emptyHeaderSet)

	t.Log("putting the Gateway back")
	gateway, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the Gateway again triggers reconciliation of HTTPRoutes and the route becomes available again")
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, gatewayClient.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the HTTPRoute does not get orphaned with the GatewayClass and Gateway gone")
	eventuallyGETPath(t, "test-http-route-essentials", http.StatusNotFound, "", emptyHeaderSet)
}
