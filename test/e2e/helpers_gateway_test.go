//go:build e2e_tests

// The file is used for putting functions related to gateway APIs.

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

// deployGateway deploys a gateway with a new created gateway class and a fixed name `kong`.
func deployGateway(ctx context.Context, t *testing.T, env environments.Environment) *gatewayv1.Gateway {
	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a supported gatewayclass to the test cluster")
	supportedGatewayClass := &gatewayv1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				// annotate the gatewayclass to unmanaged.
				annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
			},
		},
		Spec: gatewayv1.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
	}
	supportedGatewayClass, err = gc.GatewayV1beta1().GatewayClasses().Create(ctx, supportedGatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	gw := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong",
		},
		Spec: gatewayv1.GatewaySpec{
			GatewayClassName: gatewayv1.ObjectName(supportedGatewayClass.Name),
			Listeners: []gatewayv1.Listener{{
				Name:     "http",
				Protocol: gatewayv1.HTTPProtocolType,
				Port:     gatewayv1.PortNumber(80),
			}},
		},
	}
	gw, err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)
	return gw
}

// verifyGateway verifies that the gateway `gw` is ready.
func verifyGateway(ctx context.Context, t *testing.T, env environments.Environment, gw *gatewayv1.Gateway) {
	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("verifying that the gateway receives a final programmed condition once reconciliation completes")
	require.Eventually(t, func() bool {
		gw, err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		if ready := util.CheckCondition(
			gw.Status.Conditions,
			util.ConditionType(gatewayv1.GatewayConditionProgrammed),
			util.ConditionReason(gatewayv1.GatewayReasonProgrammed),
			metav1.ConditionTrue,
			gw.Generation,
		); ready {
			return true
		}

		t.Logf("conditions: %v", gw.Status.Conditions)
		return false
	}, gatewayUpdateWaitTime, time.Second)
}

// deployGatewayWithTCPListener deploys a gateway `kong` with a tcp listener to test TCPRoute.
func deployGatewayWithTCPListener(ctx context.Context, t *testing.T, env environments.Environment) *gatewayv1.Gateway {
	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a supported gatewayclass to the test cluster")
	supportedGatewayClass := &gatewayv1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				// annotate the gatewayclass to unmanaged.
				annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
			},
		},
		Spec: gatewayv1.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
	}
	supportedGatewayClass, err = gc.GatewayV1beta1().GatewayClasses().Create(ctx, supportedGatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	gw := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong",
		},
		Spec: gatewayv1.GatewaySpec{
			GatewayClassName: gatewayv1.ObjectName(supportedGatewayClass.Name),
			Listeners: []gatewayv1.Listener{
				{
					Name:     "http",
					Protocol: gatewayv1.HTTPProtocolType,
					Port:     gatewayv1.PortNumber(80),
				},
				{
					Name:     "tcp",
					Protocol: gatewayv1.TCPProtocolType,
					Port:     gatewayv1.PortNumber(tcpListenerPort),
				},
			},
		},
	}
	_, err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Get(ctx, gw.Name, metav1.GetOptions{})
	if err == nil {
		t.Logf("gateway %s exists, delete and re-create it", gw.Name)
		err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Delete(ctx, gw.Name, metav1.DeleteOptions{})
		require.NoError(t, err)
		gw, err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Create(ctx, gw, metav1.CreateOptions{})
		require.NoError(t, err)
	} else {
		require.True(t, apierrors.IsNotFound(err))
		gw, err = gc.GatewayV1beta1().Gateways(corev1.NamespaceDefault).Create(ctx, gw, metav1.CreateOptions{})
		require.NoError(t, err)
	}
	return gw
}

// deployHTTPRoute creates an `HTTPRoute` and related backend deployment/service.
// it matches the specified path `/httpbin` by prefix, so we can access the backend service by `http://$PROXY_IP/httpbin`.
func deployHTTPRoute(ctx context.Context, t *testing.T, env environments.Environment, gw *gatewayv1.Gateway) {
	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	t.Log("deploying an HTTP service to test the ingress controller and proxy")
	container := generators.NewContainer("httpbin-httproute", test.HTTPBinImage, test.HTTPBinPort)
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("creating an HTTPRoute for service %s with Gateway %s", service.Name, gw.Name)
	pathMatchPrefix := gatewayv1.PathMatchPathPrefix
	path := "/httpbin"
	httpPort := gatewayv1.PortNumber(80)
	httproute := &gatewayv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.StripPathKey: "true",
			},
		},
		Spec: gatewayv1.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: []gatewayv1.ParentReference{{
					Name: gatewayv1.ObjectName(gw.Name),
				}},
			},
			Rules: []gatewayv1.HTTPRouteRule{{
				Matches: []gatewayv1.HTTPRouteMatch{{
					Path: &gatewayv1.HTTPPathMatch{
						Type:  &pathMatchPrefix,
						Value: &path,
					},
				}},
				BackendRefs: []gatewayv1.HTTPBackendRef{{
					BackendRef: gatewayv1.BackendRef{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name: gatewayv1.ObjectName(service.Name),
							Port: &httpPort,
						},
					},
				}},
			}},
		},
	}
	_, err = gc.GatewayV1beta1().HTTPRoutes(corev1.NamespaceDefault).Create(ctx, httproute, metav1.CreateOptions{})
	require.NoError(t, err)
}

// verifyHTTPRoute verifies an HTTPRoute exposes a route at /httpbin
// TODO this is not actually specific to HTTPRoutes. It is verifyIngress with the KongIngress removed
// Once we support HTTPMethod HTTPRouteMatch handling, we can combine the two into a single generic function.
func verifyHTTPRoute(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Log("finding the kong proxy service ip")
	proxyIP := getKongProxyIP(ctx, t, env)

	t.Logf("waiting for route from Ingress to be operational at http://%s/httpbin", proxyIP)

	require.Eventually(t, func() bool {
		resp, err := helpers.DefaultHTTPClient().Get(fmt.Sprintf("http://%s/httpbin", proxyIP))
		if err != nil {
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
	}, ingressWait, time.Second)
}

// deployTCPRoute creates a `TCPRoute` and related backend deployment/service.
func deployTCPRoute(ctx context.Context, t *testing.T, env environments.Environment, gw *gatewayv1.Gateway) {
	gc, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	t.Log("deploying a TCP service to test the ingress controller and proxy")
	container := generators.NewContainer("tcpecho-tcproute", test.EchoImage, test.EchoTCPPort)
	container.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: "tcpecho-tcproute",
		},
	}
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "echo",
			Protocol:   corev1.ProtocolTCP,
			Port:       tcpListenerPort,
			TargetPort: intstr.FromInt(test.EchoTCPPort),
		},
	}
	_, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("creating a TCPRoute for service %s with Gateway %s", service.Name, gw.Name)
	portNumber := gatewayv1alpha2.PortNumber(tcpListenerPort)
	tcpRoute := &gatewayv1alpha2.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.TCPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gw.Name),
				}},
			},
			Rules: []gatewayv1alpha2.TCPRouteRule{
				{
					BackendRefs: []gatewayv1alpha2.BackendRef{
						{
							BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
								Name: gatewayv1alpha2.ObjectName(service.Name),
								Port: &portNumber,
							},
						},
					},
				},
			},
		},
	}
	_, err = gc.GatewayV1alpha2().TCPRoutes(corev1.NamespaceDefault).Create(ctx, tcpRoute, metav1.CreateOptions{})
	require.NoError(t, err)
}

// verifyTCPRoute checks whether the traffic is routed to the backend tcp-echo service,
// using eventually testing helper.
func verifyTCPRoute(ctx context.Context, t *testing.T, env environments.Environment) {
	t.Log("finding the kong proxy service ip")
	proxyIP := getKongProxyIP(ctx, t, env)

	t.Logf("waiting for route from TCPRoute to be operational at %s:%d", proxyIP, tcpListenerPort)
	require.Eventually(t, func() bool {
		ok, err := test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyIP, tcpListenerPort), "tcpecho-tcproute")
		if err != nil {
			t.Logf("failed to connect to %s:%d, error %v", proxyIP, tcpListenerPort, err)
			return false
		}
		return ok
	}, ingressWait, 5*time.Second,
	)
}
