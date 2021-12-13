//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/gateway/versioned"
)

func TestHTTPRouteEssentials(t *testing.T) {
	// TODO: make parallel later
	ns, cleanup := namespace(t)
	defer cleanup()

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
	container := generators.NewContainer("httpbin", httpBinImage, 80)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		if err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		if err := env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Logf("creating an httproute to access deployment %s via kong", deployment.Name)
	httpPort := gatewayv1alpha2.PortNumber(80)
	pathMatchPrefix := gatewayv1alpha2.PathMatchPathPrefix
	httproute := &gatewayv1alpha2.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentRef{{
					Name: gatewayv1alpha2.ObjectName(gw.Name),
				}},
			},
			Rules: []gatewayv1alpha2.HTTPRouteRule{{
				Matches: []gatewayv1alpha2.HTTPRouteMatch{{
					Path: &gatewayv1alpha2.HTTPPathMatch{
						Type:  &pathMatchPrefix,
						Value: kong.String("/httpbin"),
					},
				}},
				BackendRefs: []gatewayv1alpha2.HTTPBackendRef{{
					BackendRef: gatewayv1alpha2.BackendRef{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(service.Name),
							Port: &httpPort,
						},
					},
				}},
			}},
		},
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

	t.Log("waiting for routes from HTTPRoute to become operational")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			// now that the ingress backend is routable, make sure the contents we're getting back are what we expect
			// Expected: "<title>httpbin.org</title>"
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Log("removing the parentrefs from the HTTPRoute")
	oldParentRefs := httproute.Spec.ParentRefs
	httproute.Spec.ParentRefs = nil
	httproute, err = c.GatewayV1alpha2().HTTPRoutes(ns.Name).Update(ctx, httproute, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the parentRefs now removed")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusNotFound
	}, ingressWait, waitTick)

	t.Log("putting the parentRefs back")
	httproute.Spec.ParentRefs = oldParentRefs
	httproute, err = c.GatewayV1alpha2().HTTPRoutes(ns.Name).Update(ctx, httproute, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Log("deleting the GatewayClass")
	oldGWCName := gwc.Name
	require.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the GatewayClass now removed")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusNotFound
	}, ingressWait, waitTick)

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

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of HTTPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Log("deleting the Gateway")
	oldGWName := gw.Name
	require.NoError(t, c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the data-plane configuration from the HTTPRoute gets dropped with the Gateway now removed")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusNotFound
	}, ingressWait, waitTick)

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

	t.Log("verifying that creating the Gateway again triggers reconciliation of HTTPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			b := new(bytes.Buffer)
			n, err := b.ReadFrom(resp.Body)
			require.NoError(t, err)
			require.True(t, n > 0)
			return strings.Contains(b.String(), "<title>httpbin.org</title>")
		}
		return false
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the data-plane configuration from the HTTPRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("%s/httpbin", proxyURL))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", proxyURL, err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusNotFound
	}, ingressWait, waitTick)
}
