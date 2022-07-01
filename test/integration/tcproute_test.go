//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

const (
	tcpEchoPort = 1025
)

func TestTCPRouteEssentials(t *testing.T) {
	t.Log("locking TCP port")
	tcpMutex.Lock()
	defer func() {
		// Free up the TCP port
		t.Log("unlocking TCP port")
		tcpMutex.Unlock()
	}()

	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	t.Log("getting gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a supported gatewayclass to the test cluster")
	gatewayClassName := uuid.NewString()
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode and port 8888")
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1alpha2.Listener{{
			Name:     "tcp",
			Protocol: gatewayv1alpha2.TCPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTCPServicePort),
		}}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("creating a tcpecho pod to test TCPRoute traffic routing")
	container1 := generators.NewContainer("tcpecho-1", test.TCPEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID1 := uuid.NewString()
	container1.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID1,
		},
	}
	deployment1 := generators.NewDeploymentForContainer(container1)
	deployment1, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment1, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("creating an additional tcpecho pod to test TCPRoute multiple backendRef loadbalancing")
	container2 := generators.NewContainer("tcpecho-2", test.TCPEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID2 := uuid.NewString()
	container2.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID2,
		},
	}
	deployment2 := generators.NewDeploymentForContainer(container2)
	deployment2, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment2, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployments %s/%s and %s/%s", deployment1.Namespace, deployment1.Name, deployment2.Namespace, deployment2.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment1.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment2.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s/%s via service", deployment1.Namespace, deployment1.Name)
	service1 := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeLoadBalancer)
	// we have to override the ports so that we can map the default TCP port from
	// the Kong Gateway deployment to the tcpecho port, as this is what will be
	// used to route the traffic at the Gateway (at the time of writing, the
	// Kong Gateway doesn't support an API for dynamically adding these ports. The
	// ports must be added manually to the config or ENV).
	service1.Spec.Ports = []corev1.ServicePort{{
		Name:       "tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       ktfkong.DefaultTCPServicePort,
		TargetPort: intstr.FromInt(tcpEchoPort),
	}}
	service1, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service1, metav1.CreateOptions{})
	assert.NoError(t, err)

	t.Logf("exposing deployment %s/%s via service", deployment2.Namespace, deployment2.Name)
	service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	// we have to override the ports so that we can map the default TCP port from
	// the Kong Gateway deployment to the tcpecho port, as this is what will be
	// used to route the traffic at the Gateway (at the time of writing, the
	// Kong Gateway doesn't support an API for dynamically adding these ports. The
	// ports must be added manually to the config or ENV).
	service2.Spec.Ports = []corev1.ServicePort{{
		Name:       "tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       ktfkong.DefaultTCPServicePort,
		TargetPort: intstr.FromInt(tcpEchoPort),
	}}
	service2, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service2, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service1.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service1.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service2.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating a tcproute to access deployment %s via kong", deployment1.Name)
	tcpPortDefault := gatewayv1alpha2.PortNumber(ktfkong.DefaultTCPServicePort)
	tcpRoute := &gatewayv1alpha2.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:        uuid.NewString(),
			Annotations: map[string]string{},
		},
		Spec: gatewayv1alpha2.TCPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gatewayName),
				}},
			},
			Rules: []gatewayv1alpha2.TCPRouteRule{{
				BackendRefs: []gatewayv1alpha2.BackendRef{{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service1.Name),
						Port: &tcpPortDefault,
					},
				}},
			}},
		},
	}
	tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Create(ctx, tcpRoute, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the tcproute %s", tcpRoute.Name)
		if err := gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Delete(ctx, tcpRoute.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the tcpecho is responding properly")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("removing the parentrefs from the TCPRoute")
	oldParentRefs := tcpRoute.Spec.ParentRefs
	require.Eventually(t, func() bool {
		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Get(ctx, tcpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		tcpRoute.Spec.ParentRefs = nil
		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Update(ctx, tcpRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the tcpecho is no longer responding")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the parentRefs back")
	require.Eventually(t, func() bool {
		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Get(ctx, tcpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		tcpRoute.Spec.ParentRefs = oldParentRefs
		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Update(ctx, tcpRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the GatewayClass")
	require.NoError(t, gatewayClient.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TCPRoute gets dropped with the GatewayClass now removed")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the GatewayClass back")
	gwc, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of TCPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the Gateway")
	require.NoError(t, gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TCPRoute gets dropped with the Gateway now removed")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the Gateway back")
	gateway, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1alpha2.Listener{{
			Name:     "tcp",
			Protocol: gatewayv1alpha2.TCPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTCPServicePort),
		}}
	})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the Gateway again triggers reconciliation of TCPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, gatewayClient.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TCPRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the GatewayClass back")
	gwc, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)

	t.Log("putting the Gateway back")
	gateway, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1alpha2.Listener{{
			Name:     "tcp",
			Protocol: gatewayv1alpha2.TCPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTCPServicePort),
		}}
	})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the Gateway again triggers reconciliation of TCPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("adding an additional backendRef to the TCPRoute")
	require.Eventually(t, func() bool {
		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Get(ctx, tcpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)

		tcpRoute.Spec.Rules[0].BackendRefs = []gatewayv1alpha2.BackendRef{
			{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name: gatewayv1alpha2.ObjectName(service1.Name),
					Port: &tcpPortDefault,
				},
			},
			{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name: gatewayv1alpha2.ObjectName(service2.Name),
					Port: &tcpPortDefault,
				},
			},
		}

		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Update(ctx, tcpRoute, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Log("verifying that the TCPRoute is now load-balanced between two services")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return err == nil && responded == true
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2)
		return err == nil && responded == true
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return err == nil && responded == true
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, gatewayClient.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(t, gatewayClient, gatewayv1alpha2.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TCPRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)
}

func TestTCPRouteReferencePolicy(t *testing.T) {
	t.Log("locking TCP port")
	tcpMutex.Lock()
	defer func() {
		t.Log("unlocking TCP port")
		tcpMutex.Unlock()
	}()

	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	otherNs, err := clusters.GenerateNamespace(ctx, env.Cluster(), t.Name())
	require.NoError(t, err)

	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a gatewayclass to the test cluster")
	gatewayClassName := uuid.NewString()
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1alpha2.Listener{{
			Name:     "tcp",
			Protocol: gatewayv1alpha2.TCPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTCPServicePort),
		}}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("creating a tcpecho pod to test TCPRoute traffic routing")
	container1 := generators.NewContainer("tcpecho-1", test.TCPEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID1 := uuid.NewString()
	container1.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID1,
		},
	}
	deployment1 := generators.NewDeploymentForContainer(container1)
	deployment1, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment1, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("creating an additional tcpecho pod to test TCPRoute multiple backendRef loadbalancing")
	container2 := generators.NewContainer("tcpecho-2", test.TCPEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID2 := uuid.NewString()
	container2.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID2,
		},
	}
	deployment2 := generators.NewDeploymentForContainer(container2)
	deployment2, err = env.Cluster().Client().AppsV1().Deployments(otherNs.Name).Create(ctx, deployment2, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployments %s/%s and %s/%s", deployment1.Namespace, deployment1.Name, deployment2.Namespace, deployment2.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment1.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(otherNs.Name).Delete(ctx, deployment2.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s/%s via service", deployment1.Namespace, deployment1.Name)
	service1 := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeLoadBalancer)
	// we have to override the ports so that we can map the default TCP port from
	// the Kong Gateway deployment to the tcpecho port, as this is what will be
	// used to route the traffic at the Gateway (at the time of writing, the
	// Kong Gateway doesn't support an API for dynamically adding these ports. The
	// ports must be added manually to the config or ENV).
	service1.Spec.Ports = []corev1.ServicePort{{
		Name:       "tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       ktfkong.DefaultTCPServicePort,
		TargetPort: intstr.FromInt(tcpEchoPort),
	}}
	service1, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service1, metav1.CreateOptions{})
	assert.NoError(t, err)

	t.Logf("exposing deployment %s/%s via service", deployment2.Namespace, deployment2.Name)
	service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	service2.Spec.Ports = []corev1.ServicePort{{
		Name:       "tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       ktfkong.DefaultTCPServicePort,
		TargetPort: intstr.FromInt(tcpEchoPort),
	}}
	service2, err = env.Cluster().Client().CoreV1().Services(otherNs.Name).Create(ctx, service2, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service1.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service1.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(otherNs.Name).Delete(ctx, service2.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating a tcproute to access deployment %s via kong", deployment1.Name)
	tcpPortDefault := gatewayv1alpha2.PortNumber(ktfkong.DefaultTCPServicePort)
	remoteNamespace := gatewayv1alpha2.Namespace(otherNs.Name)
	tcproute := &gatewayv1alpha2.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:        uuid.NewString(),
			Annotations: map[string]string{},
		},
		Spec: gatewayv1alpha2.TCPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gatewayName),
				}},
			},
			Rules: []gatewayv1alpha2.TCPRouteRule{{
				BackendRefs: []gatewayv1alpha2.BackendRef{
					{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(service1.Name),
							Port: &tcpPortDefault,
						},
					},
					{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name:      gatewayv1alpha2.ObjectName(service2.Name),
							Namespace: &remoteNamespace,
							Port:      &tcpPortDefault,
						},
					},
				},
			}},
		},
	}
	tcproute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Create(ctx, tcproute, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the tcproute %s", tcproute.Name)
		if err := gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Delete(ctx, tcproute.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("verifying that only the local tcpecho is responding without a ReferencePolicy")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return err == nil && responded == true
	}, ingressWait*2, waitTick)
	require.Never(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2)
		return err == nil && responded == true
	}, time.Second*10, time.Second)

	t.Logf("creating a reference policy that permits tcproute access from %s to services in %s", ns.Name, otherNs.Name)
	policy := &gatewayv1alpha2.ReferencePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        uuid.NewString(),
			Annotations: map[string]string{},
		},
		Spec: gatewayv1alpha2.ReferencePolicySpec{
			From: []gatewayv1alpha2.ReferencePolicyFrom{
				{
					// this isn't actually used, it's just a dummy extra from to confirm we handle multiple fine
					Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
					Kind:      gatewayv1alpha2.Kind("TCPRoute"),
					Namespace: gatewayv1alpha2.Namespace("garbage"),
				},
				{
					Group:     gatewayv1alpha2.Group("gateway.networking.k8s.io"),
					Kind:      gatewayv1alpha2.Kind("TCPRoute"),
					Namespace: gatewayv1alpha2.Namespace(tcproute.Namespace),
				},
			},
			To: []gatewayv1alpha2.ReferencePolicyTo{
				// also a dummy
				{
					Group: gatewayv1alpha2.Group(""),
					Kind:  gatewayv1alpha2.Kind("Pterodactyl"),
				},
				{
					Group: gatewayv1alpha2.Group(""),
					Kind:  gatewayv1alpha2.Kind("Service"),
				},
			},
		},
	}

	policy, err = gatewayClient.GatewayV1alpha2().ReferencePolicies(otherNs.Name).Create(ctx, policy, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that requests reach both the local and remote namespace echo instances")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return err == nil && responded == true
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Logf("testing specific name references")
	serviceName := gatewayv1alpha2.ObjectName(service2.ObjectMeta.Name)
	policy.Spec.To[1] = gatewayv1alpha2.ReferencePolicyTo{
		Kind:  gatewayv1alpha2.Kind("Service"),
		Group: gatewayv1alpha2.Group(""),
		Name:  &serviceName,
	}

	policy, err = gatewayClient.GatewayV1alpha2().ReferencePolicies(otherNs.Name).Update(ctx, policy, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2)
		return err == nil && responded == true
	}, ingressWait*2, waitTick)

	t.Logf("testing incorrect name does not match")
	blueguyName := gatewayv1alpha2.ObjectName("blueguy")
	policy.Spec.To[1].Name = &blueguyName
	_, err = gatewayClient.GatewayV1alpha2().ReferencePolicies(otherNs.Name).Update(ctx, policy, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2)
		return err != nil && responded == false
	}, ingressWait, waitTick)

}

// tcpEchoResponds takes a TCP address URL and a Pod name and checks if a
// go-echo instance is running on that Pod at that address. It compares an
// expected message and its length against an expected message, returning true
// if it is and false and an error explanation if it is not
func tcpEchoResponds(url string, podName string) (bool, error) {
	dialer := net.Dialer{Timeout: time.Second * 10}
	conn, err := dialer.Dial("tcp", url)
	if err != nil {
		return false, err
	}

	header := []byte(fmt.Sprintf("Running on Pod %s.", podName))
	message := []byte("testing tcproute")

	wrote, err := conn.Write(message)
	if err != nil {
		return false, err
	}

	if wrote != len(message) {
		return false, fmt.Errorf("wrote message of size %d, expected %d", wrote, len(message))
	}

	if err := conn.SetDeadline(time.Now().Add(time.Second * 5)); err != nil {
		return false, err
	}

	headerResponse := make([]byte, len(header)+1)
	read, err := conn.Read(headerResponse)
	if err != nil {
		return false, err
	}

	if read != len(header)+1 { // add 1 for newline
		return false, fmt.Errorf("read %d bytes but expected %d", read, len(header)+1)
	}

	if !bytes.Contains(headerResponse, header) {
		return false, fmt.Errorf(`expected header response "%s", received: "%s"`, string(header), string(headerResponse))
	}

	messageResponse := make([]byte, wrote+1)
	read, err = conn.Read(messageResponse)
	if err != nil {
		return false, err
	}

	if read != len(message) {
		return false, fmt.Errorf("read %d bytes but expected %d", read, len(message))
	}

	if !bytes.Contains(messageResponse, message) {
		return false, fmt.Errorf(`expected message response "%s", received: "%s"`, string(message), string(messageResponse))
	}

	return true, nil
}
