//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	gatewaypkg "github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

var emptyHeaderSet = make(map[string]string)

func TestHTTPRouteEssentials(t *testing.T) {
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)

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
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1beta1.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

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
	pluginClient, err := clientset.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	kongplugin, err = pluginClient.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, kongplugin, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(kongplugin)

	t.Logf("creating an httproute to access deployment %s via kong", deployment.Name)
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
							Name: gatewayv1beta1.ObjectName(service.Name),
							Port: &httpPort,
							Kind: util.StringToGatewayAPIKindPtr("Service"),
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
					Value: `^audio/.*`,
					Name:  "Content-Type",
				},
			},
		})
	}
	httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Create(ctx, httpRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(httpRoute)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)
	t.Log("verifying that the httproute contains 'Programmed' condition")
	require.Eventually(t,
		GetVerifyProgrammedConditionCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name, metav1.ConditionTrue),
		ingressWait, waitTick,
	)

	t.Log("waiting for routes from HTTPRoute to become operational")
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-essentials/base64/wqt5b8q7ccK7IGRhbiBib3NocWEgYmlyIGphdm9iaW1peiB5b8q7cWRpci4K",
		http.StatusOK, "«yoʻq» dan boshqa bir javobimiz yoʻqdir.", emptyHeaderSet, ingressWait, waitTick)
	helpers.EventuallyGETPath(t, proxyURL, "2/test-http-route-essentials/regex/999", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)
	helpers.EventuallyGETPath(t, proxyURL, "3/exact-test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)
	helpers.EventuallyGETPath(t, proxyURL, "3/exact-test-http-route-essentialsNO", http.StatusNotFound, "no Route matched", emptyHeaderSet, ingressWait, waitTick)

	require.Eventually(t, func() bool {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", proxyURL, "test-http-route-essentials"), nil)
		if err != nil {
			t.Logf("WARNING: failed to create HTTP request: %v", err)
			return false
		}
		resp, err := helpers.DefaultHTTPClient().Do(req)
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
		helpers.EventuallyGETPath(t, proxyURL, "", http.StatusOK, "<title>httpbin.org</title>", map[string]string{"Content-Type": "audio/mp3"}, ingressWait, waitTick)
	}

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
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the parentRefs now removed")
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-essentials", http.StatusNotFound, "", emptyHeaderSet, ingressWait, waitTick)

	t.Log("putting the parentRefs back")
	require.Eventually(t, func() bool {
		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Get(ctx, httpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httpRoute.Spec.ParentRefs = oldParentRefs
		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Update(ctx, httpRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)

	t.Log("deleting the GatewayClass")
	require.NoError(t, gatewayClient.GatewayV1beta1().GatewayClasses().Delete(ctx, gatewayClassName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)
	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the GatewayClass now removed")
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-essentials", http.StatusNotFound, "", emptyHeaderSet, ingressWait, waitTick)

	t.Log("putting the GatewayClass back")
	gwc, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of HTTPRoutes and the route becomes available again")
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)

	t.Log("deleting the Gateway")
	require.NoError(t, gatewayClient.GatewayV1beta1().Gateways(ns.Name).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the Gateway now removed")
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-essentials", http.StatusNotFound, "", emptyHeaderSet, ingressWait, waitTick)

	t.Log("putting the Gateway back")
	gateway, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1beta1.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the Gateway again triggers reconciliation of HTTPRoutes and the route becomes available again")
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-essentials", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, gatewayClient.GatewayV1beta1().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, gatewayClient.GatewayV1beta1().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the HTTPRoute does not get orphaned with the GatewayClass and Gateway gone")
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-essentials", http.StatusNotFound, "", emptyHeaderSet, ingressWait, waitTick)

	t.Log("testing port matching....")
	t.Log("putting the Gateway back")
	_, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1beta1.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)
	t.Log("putting the GatewayClass back")
	_, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)

	t.Log("verifying that the HTTPRoute has the Condition 'Accepted' set to 'True' before specifying a port not existent in Gateway")
	require.Eventually(t, HTTPRouteMatchesAcceptedCallback(t, gatewayClient, httpRoute, true, gatewayv1beta1.RouteReasonAccepted), statusWait, waitTick)

	// Set the Port in ParentRef which does not have a matching listener in Gateway.
	require.Eventually(t, func() bool {
		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Get(ctx, httpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		port81 := gatewayv1beta1.PortNumber(81)
		httpRoute.Spec.ParentRefs[0].Port = &port81
		httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Update(ctx, httpRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the HTTPRoute has the Condition 'Accepted' set to 'False' when it specified a port not existent in Gateway")
	require.Eventually(t, HTTPRouteMatchesAcceptedCallback(t, gatewayClient, httpRoute, false, gatewaypkg.RouteReasonNoMatchingParent), statusWait, waitTick)
}

func TestHTTPRouteMultipleServices(t *testing.T) {
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)

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
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1beta1.Gateway) {
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
	// service3 has an annotation the others don't. we expect the controller to skip rules that have different annotations
	service3 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	service3.Annotations = map[string]string{annotations.AnnotationPrefix + annotations.HostHeaderKey: "example.com"}
	service3.Name = "nginx-host"
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service3, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service1)
	cleaner.Add(service2)
	cleaner.Add(service3)

	t.Log("adding an HTTPRoute with multi-backend rules")
	var httpbinWeight int32 = 75
	var nginxWeight int32 = 25
	httpPort := gatewayv1beta1.PortNumber(80)
	pathMatchPrefix := gatewayv1beta1.PathMatchPathPrefix
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
			Rules: []gatewayv1beta1.HTTPRouteRule{
				{
					Matches: []gatewayv1beta1.HTTPRouteMatch{
						{
							Path: &gatewayv1beta1.HTTPPathMatch{
								Type:  &pathMatchPrefix,
								Value: kong.String("/test-http-route-multiple-services"),
							},
						},
					},
					BackendRefs: []gatewayv1beta1.HTTPBackendRef{
						{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName(service1.Name),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
								Weight: &httpbinWeight,
							},
						},
						{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName(service2.Name),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
								Weight: &nginxWeight,
							},
						},
					},
				},
				{
					Matches: []gatewayv1beta1.HTTPRouteMatch{
						{
							Path: &gatewayv1beta1.HTTPPathMatch{
								Type:  &pathMatchPrefix,
								Value: kong.String("/test-http-route-multiple-services-broken"),
							},
						},
					},
					BackendRefs: []gatewayv1beta1.HTTPBackendRef{
						{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName(service1.Name),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
								Weight: &httpbinWeight,
							},
						},
						{
							BackendRef: gatewayv1beta1.BackendRef{
								BackendObjectReference: gatewayv1beta1.BackendObjectReference{
									Name: gatewayv1beta1.ObjectName(service3.Name),
									Port: &httpPort,
									Kind: util.StringToGatewayAPIKindPtr("Service"),
								},
								Weight: &nginxWeight,
							},
						},
					},
				},
			},
		},
	}
	httpRoute, err = gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name).Create(ctx, httpRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(httpRoute)

	t.Log("verifying that both backends are ready to receive traffic")
	httpbinRespContent := "<title>httpbin.org</title>"
	nginxRespContent := "<title>Welcome to nginx!</title>"
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-multiple-services", http.StatusOK, httpbinRespContent, emptyHeaderSet, ingressWait, waitTick)
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-multiple-services", http.StatusOK, nginxRespContent, emptyHeaderSet, ingressWait, waitTick)

	t.Log("verifying that both backends receive requests according to weighted distribution")
	httpbinRespName := "httpbin-resp"
	nginxRespName := "nginx-resp"
	toleranceDelta := 0.2
	expectedRespRatio := map[string]int{
		httpbinRespName: int(httpbinWeight),
		nginxRespName:   int(nginxWeight),
	}
	weightedLoadBalancingTestConfig := helpers.CountHTTPResponsesConfig{
		Method:      http.MethodGet,
		Path:        "test-http-route-multiple-services",
		Headers:     emptyHeaderSet,
		Duration:    5 * time.Second,
		RequestTick: 50 * time.Millisecond,
	}
	respCounter := helpers.CountHTTPGetResponses(t,
		proxyURL,
		weightedLoadBalancingTestConfig,
		helpers.MatchRespByStatusAndContent(httpbinRespName, http.StatusOK, httpbinRespContent),
		helpers.MatchRespByStatusAndContent(nginxRespName, http.StatusOK, nginxRespContent),
	)
	assert.InDeltaMapValues(t,
		helpers.DistributionOfMapValues(respCounter),
		helpers.DistributionOfMapValues(expectedRespRatio),
		toleranceDelta,
		"Response distribution does not match expected distribution within %f%% delta,"+
			" request-count=%v, expected-ratio=%v",
		toleranceDelta*100, respCounter, expectedRespRatio,
	)

	t.Log("verifying that misconfigured service rules are _not_ routed")
	helpers.EventuallyGETPath(t, proxyURL, "test-http-route-multiple-services-broken", http.StatusNotFound, "", emptyHeaderSet, ingressWait, waitTick)
}

func TestHTTPRouteFilterHosts(t *testing.T) {
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)

	listenerHostname := gatewayv1beta1.Hostname("test.specific.io")

	t.Log("getting a gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a new gatewayClass")
	gatewayClassName := uuid.NewString()
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("deploying a new gateway with specified hostname")
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1beta1.Gateway) {
		gw.Name = gatewayName
		for i := range gw.Spec.Listeners {
			gw.Spec.Listeners[i].Hostname = &listenerHostname
		}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("creating an httproute with a same hostname and another unmatched hostname")
	httpRoute := &gatewayv1beta1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.StripPathKey: "true",
			},
		},
		Spec: gatewayv1beta1.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
				ParentRefs: []gatewayv1beta1.ParentReference{{
					Name: gatewayv1beta1.ObjectName(gateway.Name),
				}},
			},
			Hostnames: []gatewayv1beta1.Hostname{
				gatewayv1beta1.Hostname("test.specific.io"),
				gatewayv1beta1.Hostname("another.specific.io"),
			},
			Rules: []gatewayv1beta1.HTTPRouteRule{{
				Matches: []gatewayv1beta1.HTTPRouteMatch{
					builder.NewHTTPRouteMatch().WithPathPrefix("/test-http-route-filter-hosts").Build(),
				},
				BackendRefs: []gatewayv1beta1.HTTPBackendRef{
					builder.NewHTTPBackendRef(service.Name).WithPort(80).Build(),
				},
			}},
		},
	}
	hClient := gatewayClient.GatewayV1beta1().HTTPRoutes(ns.Name)
	httpRoute, err = hClient.Create(ctx, httpRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(httpRoute)

	// testGetByHost tries to get the test path with specified host in request,
	// and returns true if 200 returned.
	testGetByHost := func(t *testing.T, host string) bool {
		req := helpers.MustHTTPRequest(t, "GET", proxyURL, "/test-http-route-filter-hosts", nil)
		req.Host = host
		resp, err := helpers.DefaultHTTPClient().Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}

	t.Logf("test host matched hostname in listeners")
	require.Eventually(t, func() bool {
		return testGetByHost(t, "test.specific.io")
	}, ingressWait, waitTick)
	t.Logf("test host matched in httproute, but not in listeners")
	require.False(t, testGetByHost(t, "another.specific.io"))

	t.Logf("update hostnames in httproute to wildcard")
	require.Eventually(t, func() bool {
		httpRoute, err = hClient.Get(ctx, httpRoute.Name, metav1.GetOptions{})
		if err != nil {
			t.Logf("failed getting the HTTPRoute %s: %v", httpRoute.Name, err)
			return false
		}
		httpRoute.Spec.Hostnames = []gatewayv1beta1.Hostname{
			gatewayv1beta1.Hostname("*.specific.io"),
		}
		httpRoute, err = hClient.Update(ctx, httpRoute, metav1.UpdateOptions{})
		if err != nil {
			t.Logf("failed updating the HTTPRoute %s: %v", httpRoute.Name, err)
			return false
		}
		return true
	}, 10*time.Second, 100*time.Millisecond)
	t.Logf("test host matched hostname in listeners")
	require.Eventually(t, func() bool {
		return testGetByHost(t, "test.specific.io")
	}, ingressWait, waitTick)
	t.Logf("test host matched in httproute, but not in listeners")
	require.False(t, testGetByHost(t, "another2.specific.io"))

	t.Logf("update hostname in httproute to an unmatched host")
	httpRoute, err = hClient.Get(ctx, httpRoute.Name, metav1.GetOptions{})
	require.NoError(t, err)
	httpRoute.Spec.Hostnames = []gatewayv1beta1.Hostname{
		gatewayv1beta1.Hostname("another.specific.io"),
	}
	httpRoute, err = hClient.Update(ctx, httpRoute, metav1.UpdateOptions{})
	require.NoError(t, err)
	t.Logf("status of httproute should contain an 'Accepted' condition with 'False' status")
	require.Eventuallyf(t, func() bool {
		currentHTTPRoute, err := hClient.Get(ctx, httpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, parent := range currentHTTPRoute.Status.Parents {
			for _, condition := range parent.Conditions {
				if condition.Type == string(gatewayv1beta1.RouteReasonAccepted) && condition.Status == metav1.ConditionFalse {
					return true
				}
			}
		}
		return false
	}, ingressWait, waitTick,
		func() string {
			currentHTTPRoute, err := hClient.Get(ctx, httpRoute.Name, metav1.GetOptions{})
			if err != nil {
				return err.Error()
			}
			return fmt.Sprintf("current status of HTTPRoute %s/%s:%v", httpRoute.Namespace, httpRoute.Name, currentHTTPRoute.Status)
		}())
	t.Logf("test host matched in httproute, but not in listeners")
	require.False(t, testGetByHost(t, "another.specific.io"))
}
