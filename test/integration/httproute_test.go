//go:build integration_tests
// +build integration_tests

package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

var emptyHeaderSet = make(map[string]string)

func TestHTTPRouteEssentials(t *testing.T) {
	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

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
	deployment1 := generators.NewDeploymentForContainer(container1)
	deployment1, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment1, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("deploying an extra minimal HTTP container deployment to test multiple backendRefs")
	container2 := generators.NewContainer("nginx", "nginx", 80)
	deployment2 := generators.NewDeploymentForContainer(container2)
	deployment2, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment2, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up the deployments")
		if err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment1.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
		if err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment2.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Logf("exposing deployment %s via service", deployment1.Name)
	service1 := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service1, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment2.Name)
	service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service2, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service1.Name)
		if err := env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service1.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
		if err := env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service2.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

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
	httpPort := gatewayv1alpha2.PortNumber(80)
	pathMatchPrefix := gatewayv1alpha2.PathMatchPathPrefix
	pathMatchRegularExpression := gatewayv1alpha2.PathMatchRegularExpression
	pathMatchExact := gatewayv1alpha2.PathMatchExact
	headerMatchRegex := gatewayv1alpha2.HeaderMatchRegularExpression
	httpRoute := &gatewayv1alpha2.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.StripPathKey: "true",
				annotations.AnnotationPrefix + annotations.PluginsKey:   "correlation",
			},
		},
		Spec: gatewayv1alpha2.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gateway.Name),
				}},
			},
			Rules: []gatewayv1alpha2.HTTPRouteRule{{
				Matches: []gatewayv1alpha2.HTTPRouteMatch{
					{
						Path: &gatewayv1alpha2.HTTPPathMatch{
							Type:  &pathMatchPrefix,
							Value: kong.String("/httpbin"),
						},
					},
					{
						Path: &gatewayv1alpha2.HTTPPathMatch{
							Type:  &pathMatchRegularExpression,
							Value: kong.String(`/regex-\d{3}-httpbin`),
						},
					},
					{
						Path: &gatewayv1alpha2.HTTPPathMatch{
							Type:  &pathMatchExact,
							Value: kong.String(`/exact-httpbin`),
						},
					},
				},
				BackendRefs: []gatewayv1alpha2.HTTPBackendRef{{
					BackendRef: gatewayv1alpha2.BackendRef{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(service1.Name),
							Port: &httpPort,
						},
					},
				}},
			}},
		},
	}
	if util.GetKongVersion().GTE(parser.MinRegexHeaderKongVersion) {
		httpRoute.Spec.Rules[0].Matches = append(httpRoute.Spec.Rules[0].Matches, gatewayv1alpha2.HTTPRouteMatch{
			Headers: []gatewayv1alpha2.HTTPHeaderMatch{
				{
					Type:  &headerMatchRegex,
					Value: "^audio/*",
					Name:  "Content-Type",
				},
			},
		})
	}
	httpRoute, err = gatewayClient.GatewayV1alpha2().HTTPRoutes(ns.Name).Create(ctx, httpRoute, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the httproute %s", httpRoute.Name)
		if err := gatewayClient.GatewayV1alpha2().HTTPRoutes(ns.Name).Delete(ctx, httpRoute.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("waiting for routes from HTTPRoute to become operational")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "httpbin/base64/wqt5b8q7ccK7IGRhbiBib3NocWEgYmlyIGphdm9iaW1peiB5b8q7cWRpci4K",
		http.StatusOK, "«yoʻq» dan boshqa bir javobimiz yoʻqdir.", emptyHeaderSet)
	eventuallyGETPath(t, "regex-123-httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "exact-httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "exact-httpbina", http.StatusNotFound, "no Route matched", emptyHeaderSet)

	require.Eventually(t, func() bool {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", proxyURL, "httpbin"), nil)
		if err != nil {
			t.Logf("WARNING: failed to create HTTP request: %v", err)
			return false
		}
		resp, err := httpc.Do(req)
		if err != nil {
			t.Logf("WARNING: http request failed for GET %s/%s: %v", proxyURL, "httpbin", err)
			return false
		}
		defer resp.Body.Close()
		if _, ok := resp.Header["Reqid"]; ok {
			return true
		}
		return false
	}, ingressWait, waitTick)

	if util.GetKongVersion().GTE(parser.MinRegexHeaderKongVersion) {
		t.Log("verifying HTTPRoute header match")
		eventuallyGETPath(t, "", http.StatusOK, "<title>httpbin.org</title>", map[string]string{"Content-Type": "audio/mp3"})
	}

	t.Log("adding an additional backendRef to the HTTPRoute")
	var httpbinWeight int32 = 75
	var nginxWeight int32 = 25
	require.Eventually(t, func() bool {
		httpRoute, err = gatewayClient.GatewayV1alpha2().HTTPRoutes(httpRoute.Namespace).Get(ctx, httpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httpRoute.Spec.Rules[0].BackendRefs = []gatewayv1alpha2.HTTPBackendRef{
			{
				BackendRef: gatewayv1alpha2.BackendRef{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service1.Name),
						Port: &httpPort,
					},
					Weight: &httpbinWeight,
				},
			},
			{
				BackendRef: gatewayv1alpha2.BackendRef{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service2.Name),
						Port: &httpPort,
					},
					Weight: &nginxWeight,
				},
			},
		}

		httpRoute, err = gatewayClient.GatewayV1alpha2().HTTPRoutes(httpRoute.Namespace).Update(ctx, httpRoute, metav1.UpdateOptions{})
		if err != nil {
			t.Logf("WARNING: could not update httproute with an additional backendRef: %s (retrying)", err)
			return false
		}

		return true
	}, ingressWait, waitTick)

	t.Log("verifying that both backends are ready to receive traffic")
	httpbinRespContent := "<title>httpbin.org</title>"
	nginxRespContent := "<title>Welcome to nginx!</title>"
	eventuallyGETPath(t, "httpbin", http.StatusOK, httpbinRespContent, emptyHeaderSet)
	eventuallyGETPath(t, "httpbin", http.StatusOK, nginxRespContent, emptyHeaderSet)

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
		Path:        "httpbin",
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
		httpRoute, err = gatewayClient.GatewayV1alpha2().HTTPRoutes(ns.Name).Get(ctx, httpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httpRoute.Spec.ParentRefs = nil
		httpRoute, err = gatewayClient.GatewayV1alpha2().HTTPRoutes(ns.Name).Update(ctx, httpRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the parentRefs now removed")
	eventuallyGETPath(t, "httpbin", http.StatusNotFound, "", emptyHeaderSet)

	t.Log("putting the parentRefs back")
	require.Eventually(t, func() bool {
		httpRoute, err = gatewayClient.GatewayV1alpha2().HTTPRoutes(ns.Name).Get(ctx, httpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httpRoute.Spec.ParentRefs = oldParentRefs
		httpRoute, err = gatewayClient.GatewayV1alpha2().HTTPRoutes(ns.Name).Update(ctx, httpRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("deleting the GatewayClass")
	require.NoError(t, gatewayClient.GatewayV1alpha2().GatewayClasses().Delete(ctx, gatewayClassName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)
	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the GatewayClass now removed")
	eventuallyGETPath(t, "httpbin", http.StatusNotFound, "", emptyHeaderSet)

	t.Log("putting the GatewayClass back")
	gwc, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of HTTPRoutes and the route becomes available again")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("deleting the Gateway")
	require.NoError(t, gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the Gateway now removed")
	eventuallyGETPath(t, "httpbin", http.StatusNotFound, "", emptyHeaderSet)

	t.Log("putting the Gateway back")
	gateway, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the Gateway again triggers reconciliation of HTTPRoutes and the route becomes available again")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, gatewayClient.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the HTTPRoute does not get orphaned with the GatewayClass and Gateway gone")
	eventuallyGETPath(t, "httpbin", http.StatusNotFound, "", emptyHeaderSet)
}
