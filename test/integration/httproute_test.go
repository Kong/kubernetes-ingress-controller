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
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
	"github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset"
)

var emptyHeaderSet = make(map[string]string)

func TestHTTPRouteEssentials(t *testing.T) {
	ns, cleanup := namespace(t)
	defer cleanup()

	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2461
	t.Log("deploying a supported gatewayclass to the test cluster")
	c, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	gwc := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	gwc, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up gatewayclasses")
		if err := c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	gw := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong",
			Annotations: map[string]string{
				unmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(gwc.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up gateways")
		if err := c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container1 := generators.NewContainer("httpbin", httpBinImage, 80)
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
	kongplugin, err = pluginClient.ConfigurationV1().KongPlugins(ns.Name).Create(ctx, kongplugin, metav1.CreateOptions{})

	t.Logf("creating an httproute to access deployment %s via kong", deployment1.Name)
	httpPort := gatewayv1alpha2.PortNumber(80)
	pathMatchPrefix := gatewayv1alpha2.PathMatchPathPrefix
	pathMatchRegularExpression := gatewayv1alpha2.PathMatchRegularExpression
	pathMatchExact := gatewayv1alpha2.PathMatchExact
	headerMatchRegex := gatewayv1alpha2.HeaderMatchRegularExpression
	httproute := &gatewayv1alpha2.HTTPRoute{
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
					Name: gatewayv1alpha2.ObjectName(gw.Name),
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
		httproute.Spec.Rules[0].Matches = append(httproute.Spec.Rules[0].Matches, gatewayv1alpha2.HTTPRouteMatch{
			Headers: []gatewayv1alpha2.HTTPHeaderMatch{
				{
					Type:  &headerMatchRegex,
					Value: "^audio/*",
					Name:  "Content-Type",
				},
			},
		})
	}
	httproute, err = c.GatewayV1alpha2().HTTPRoutes(ns.Name).Create(ctx, httproute, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the httproute %s", httproute.Name)
		if err := c.GatewayV1alpha2().HTTPRoutes(ns.Name).Delete(ctx, httproute.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("verifying that the Gateway gets linked to the route via status")
	eventuallyGatewayIsLinkedInStatus(t, c, ns.Name, httproute.Name)

	t.Log("waiting for routes from HTTPRoute to become operational")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "httpbin/base64/wqt5b8q7ccK7IGRhbiBib3NocWEgYmlyIGphdm9iaW1peiB5b8q7cWRpci4K",
		http.StatusOK, "«yoʻq» dan boshqa bir javobimiz yoʻqdir.", emptyHeaderSet)
	eventuallyGETPath(t, "regex-123-httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "exact-httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "exact-httpbina", http.StatusNotFound, "no Route matched", emptyHeaderSet)

	require.Eventually(t, func() bool {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", proxyURL, "httpbin"), nil)
		resp, err := httpc.Do(req)
		if err != nil {
			t.Logf("WARNING: http request failed for GET %s/%s: %v", proxyURL, "httpbin", err)
			return false
		}
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
		httproute, err = c.GatewayV1alpha2().HTTPRoutes(httproute.Namespace).Get(ctx, httproute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httproute.Spec.Rules[0].BackendRefs = []gatewayv1alpha2.HTTPBackendRef{
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

		httproute, err = c.GatewayV1alpha2().HTTPRoutes(httproute.Namespace).Update(ctx, httproute, metav1.UpdateOptions{})
		if err != nil {
			t.Logf("WARNING: could not update httproute with an additional backendRef: %s (retrying)", err)
			return false
		}

		return true
	}, ingressWait, waitTick)

	// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2452 need to verify weight distribution
	t.Log("verifying that both backends receive requests")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>Welcome to nginx!</title>", emptyHeaderSet)

	t.Log("removing the parentrefs from the HTTPRoute")
	oldParentRefs := httproute.Spec.ParentRefs
	require.Eventually(t, func() bool {
		httproute, err = c.GatewayV1alpha2().HTTPRoutes(ns.Name).Get(ctx, httproute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httproute.Spec.ParentRefs = nil
		httproute, err = c.GatewayV1alpha2().HTTPRoutes(ns.Name).Update(ctx, httproute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	eventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, httproute.Name)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the parentRefs now removed")
	eventuallyGETPath(t, "httpbin", http.StatusNotFound, "", emptyHeaderSet)

	t.Log("putting the parentRefs back")
	require.Eventually(t, func() bool {
		httproute, err = c.GatewayV1alpha2().HTTPRoutes(ns.Name).Get(ctx, httproute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		httproute.Spec.ParentRefs = oldParentRefs
		httproute, err = c.GatewayV1alpha2().HTTPRoutes(ns.Name).Update(ctx, httproute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets linked to the route via status")
	eventuallyGatewayIsLinkedInStatus(t, c, ns.Name, httproute.Name)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("deleting the GatewayClass")
	oldGWCName := gwc.Name
	require.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	eventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, httproute.Name)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the GatewayClass now removed")
	eventuallyGETPath(t, "httpbin", http.StatusNotFound, "", emptyHeaderSet)

	t.Log("putting the GatewayClass back")
	gwc = &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: oldGWCName,
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	gwc, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	eventuallyGatewayIsLinkedInStatus(t, c, ns.Name, httproute.Name)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of HTTPRoutes and the route becomes available again")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("deleting the Gateway")
	oldGWName := gw.Name
	require.NoError(t, c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	eventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, httproute.Name)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the Gateway now removed")
	eventuallyGETPath(t, "httpbin", http.StatusNotFound, "", emptyHeaderSet)

	t.Log("putting the Gateway back")
	gw = &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: oldGWName,
			Annotations: map[string]string{
				unmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(gwc.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	eventuallyGatewayIsLinkedInStatus(t, c, ns.Name, httproute.Name)

	t.Log("verifying that creating the Gateway again triggers reconciliation of HTTPRoutes and the route becomes available again")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	eventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, httproute.Name)

	t.Log("verifying that the data-plane configuration from the HTTPRoute does not get orphaned with the GatewayClass and Gateway gone")
	eventuallyGETPath(t, "httpbin", http.StatusNotFound, "", emptyHeaderSet)
}

// -----------------------------------------------------------------------------
// HTTPRoute Tests - Status Utilities
// -----------------------------------------------------------------------------

// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2461
func eventuallyGatewayIsLinkedInStatus(t *testing.T, c *gatewayclient.Clientset, namespace, name string) {
	require.Eventually(t, func() bool {
		// gather a fresh copy of the HTTPRoute
		httproute, err := c.GatewayV1alpha2().HTTPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		require.NoError(t, err)

		// determine if there is a link to a supported Gateway
		for _, parentStatus := range httproute.Status.Parents {
			if parentStatus.ControllerName == gateway.ControllerName {
				// supported Gateway link was found
				return true
			}
		}

		// if no link was found yet retry
		return false
	}, ingressWait, waitTick)
}

// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2461
func eventuallyGatewayIsUnlinkedInStatus(t *testing.T, c *gatewayclient.Clientset, namespace, name string) {
	require.Eventually(t, func() bool {
		// gather a fresh copy of the HTTPRoute
		httproute, err := c.GatewayV1alpha2().HTTPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		require.NoError(t, err)

		// determine if there is a link to a supported Gateway
		for _, parentStatus := range httproute.Status.Parents {
			if parentStatus.ControllerName == gateway.ControllerName {
				// a supported Gateway link was found retry
				return false
			}
		}

		// linked gateway is not present, all set
		return true
	}, ingressWait, waitTick)
}
