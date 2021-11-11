//go:build integration_tests
// +build integration_tests

package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/gateway/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
)

const (
	gatewayReconcilationWait = time.Second * 10
	unmanagedAnnotation      = annotations.AnnotationPrefix + annotations.GatewayUnmanagedAnnotation
)

func TestUnmanagedGatewayBasics(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("generating a gateway kubernetes client and gathering test data")
	pubsvc, err := env.Cluster().Client().CoreV1().Services(controllerNamespace).Get(ctx, "ingress-controller-kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	c, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a supported gatewayclass to the test cluster")
	supportedGatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	supportedGatewayClass, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, supportedGatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up gatewayclasses")
		assert.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, supportedGatewayClass.Name, metav1.DeleteOptions{}))
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
			GatewayClassName: gatewayv1alpha2.ObjectName(supportedGatewayClass.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the Gateway object gets scheduled by the controller")
	require.Eventually(t, func() bool {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gw.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayReasonScheduled) {
				return true
			}
		}
		return false
	}, time.Minute, time.Second)

	t.Log("verifying that the gateway service ref gets provisioned when placeholder is used")
	require.Eventually(t, func() bool {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return gw.Annotations[unmanagedAnnotation] == "kong-system/ingress-controller-kong-proxy"
	}, time.Minute, time.Second)

	t.Log("verifying that the gateway spec gets updated to match the publish service")
	require.Eventually(t, func() bool {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return len(gw.Spec.Listeners) == len(pubsvc.Spec.Ports) && len(gw.Spec.Addresses) > 0
	}, time.Minute, time.Second)

	t.Log("verifying that the gateway status gets updated to match the publish service")
	require.Eventually(t, func() bool {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return len(gw.Status.Listeners) == len(gw.Spec.Listeners) && len(gw.Status.Addresses) == len(gw.Spec.Addresses)
	}, time.Minute, time.Second)

	t.Log("verifying that the gateway listeners get configured with L7 configurations from the data-plane")
	require.Eventually(t, func() bool {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, listener := range gw.Spec.Listeners {
			if listener.Protocol == gatewayv1alpha2.HTTPProtocolType {
				return true
			}
		}
		return false
	}, time.Minute, time.Second)

	t.Log("verifying that the gateway receives a final ready condition once reconciliation completes")
	require.Eventually(t, func() bool {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gw.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayReasonReady) {
				return true
			}
		}
		return false
	}, time.Minute, time.Second)
}

func TestUnmanagedGatewayServiceUpdates(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("generating a gateway kubernetes client")
	c, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a supported gatewayclass to the test cluster")
	supportedGatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	supportedGatewayClass, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, supportedGatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up gatewayclasses")
		assert.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, supportedGatewayClass.Name, metav1.DeleteOptions{}))
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
			GatewayClassName: gatewayv1alpha2.ObjectName(supportedGatewayClass.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("deploying another gateway to the test cluster using unmanaged gateway mode")
	gw2 := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong2",
			Annotations: map[string]string{
				unmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(supportedGatewayClass.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	gw2, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw2, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the gateways receive a final ready condition once reconciliation completes")
	require.Eventually(t, func() bool {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gw.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayReasonReady) {
				return true
			}
		}
		return false
	}, time.Minute, time.Second)
	require.Eventually(t, func() bool {
		gw2, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw2.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gw2.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayReasonReady) {
				return true
			}
		}
		return false
	}, time.Minute, time.Second)

	t.Log("triggering an update to the gateway service")
	svc, err := env.Cluster().Client().CoreV1().Services(controllerNamespace).Get(ctx, "ingress-controller-kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	svc.Spec.Ports = append(svc.Spec.Ports, corev1.ServicePort{
		Name:     "sanfrancisco",
		Protocol: corev1.ProtocolTCP,
		Port:     65455,
	})
	_, err = env.Cluster().Client().CoreV1().Services(controllerNamespace).Update(ctx, svc, metav1.UpdateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up changes to the gateway service")
		svc, err := env.Cluster().Client().CoreV1().Services(controllerNamespace).Get(ctx, "ingress-controller-kong-proxy", metav1.GetOptions{})
		require.NoError(t, err)
		newPorts := make([]corev1.ServicePort, 0, len(svc.Spec.Ports)-1)
		for _, port := range svc.Spec.Ports {
			if port.Name != "sanfrancisco" {
				newPorts = append(newPorts, port)
			}
		}
		svc.Spec.Ports = newPorts
		_, err = env.Cluster().Client().CoreV1().Services(controllerNamespace).Update(ctx, svc, metav1.UpdateOptions{})
		require.NoError(t, err)
	}()

	t.Log("verifying that the gateway resources get updated with listeners that match the new port")
	require.Eventually(t, func() bool {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, listener := range gw.Spec.Listeners {
			if listener.Name == "sanfrancisco" && listener.Protocol == gatewayv1alpha2.TCPProtocolType && listener.Port == 65455 {
				return true
			}
		}
		return false
	}, time.Minute, time.Second)
	require.Eventually(t, func() bool {
		gw2, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw2.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, listener := range gw.Spec.Listeners {
			if listener.Name == "sanfrancisco" && listener.Protocol == gatewayv1alpha2.TCPProtocolType && listener.Port == 65455 {
				return true
			}
		}
		return false
	}, time.Minute, time.Second)
}

func TestUnmanagedGatewayControllerSupport(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("generating a gateway kubernetes client")
	c, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a supported gatewayclass to the test cluster")
	supportedGatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	supportedGatewayClass, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, supportedGatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("deploying an unsupported gatewayclass to the test cluster")
	unsupportedGatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: "unsupported.acme.com/gateway-controller",
		},
	}
	unsupportedGatewayClass, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, unsupportedGatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up gatewayclasses")
		assert.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, supportedGatewayClass.Name, metav1.DeleteOptions{}))
		assert.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, unsupportedGatewayClass.Name, metav1.DeleteOptions{}))
	}()

	t.Log("deploying a gateway using the unsupported gateway class")
	unsupportedGateway := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "acme",
			Annotations: map[string]string{
				unmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(unsupportedGatewayClass.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	unsupportedGateway, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, unsupportedGateway, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the unsupported Gateway object does not get scheduled by the controller")
	timeout := time.Now().Add(gatewayReconcilationWait)
	for timeout.After(time.Now()) {
		unsupportedGateway, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, unsupportedGateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		require.Len(t, unsupportedGateway.Status.Conditions, 1)
		require.Equal(t, string(gatewayv1alpha2.GatewayReasonNotReconciled), unsupportedGateway.Status.Conditions[0].Reason)
	}

	t.Log("deploying a gateway that is not configured for unmanaged mode, but is using a supported class")
	managedGateway := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong-managed",
			// missing annotation should cause failure
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(supportedGatewayClass.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	managedGateway, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, managedGateway, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the managed Gateway object does not get scheduled due to lack of support")
	require.Eventually(t, func() bool {
		managedGateway, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, managedGateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range managedGateway.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayReasonNoResources) {
				return true
			}
		}
		return false
	}, time.Minute, time.Second)
}

func TestUnmanagedGatewayClass(t *testing.T) {
	t.Parallel()
	ns, cleanup := namespace(t)
	defer cleanup()

	t.Log("generating a gateway kubernetes client")
	c, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	className := uuid.NewString()

	t.Log("deploying a gateway to the test cluster using unmanaged mode, but with no valid gatewayclass yet")
	gw := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong",
			Annotations: map[string]string{
				unmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(className),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
			}},
		},
	}
	gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the Gateway object does not get scheduled by the controller due to missing its GatewayClass")
	timeout := time.Now().Add(gatewayReconcilationWait)
	for timeout.After(time.Now()) {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		require.Len(t, gw.Status.Conditions, 1)
		require.Equal(t, string(gatewayv1alpha2.GatewayReasonNotReconciled), gw.Status.Conditions[0].Reason)
	}

	t.Log("deploying the missing gatewayclass to the test cluster")
	supportedGatewayClass := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: className,
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	_, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, supportedGatewayClass, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("now that the gatewayclass exists, verifying that the gateway resource gets resolved")
	require.Eventually(t, func() bool {
		gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gw.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayReasonReady) {
				return true
			}
		}
		return false
	}, time.Minute, time.Second)
}
