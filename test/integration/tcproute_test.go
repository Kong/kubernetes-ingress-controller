//go:build integration_tests

package integration

import (
	"context"
	"errors"
	"fmt"
	"io"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

const gatewayTCPPortName = "tcp"

func TestTCPRouteEssentials(t *testing.T) {
	ctx := context.Background()
	RunWhenKongExpressionRouter(t)
	t.Log("locking TCP port")
	tcpMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TCP port")
		tcpMutex.Unlock()
	})

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("getting gateway client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a supported gatewayclass to the test cluster")
	gatewayClassName := uuid.NewString()
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Logf("deploying a gateway to the test cluster using unmanaged gateway mode and port %d", ktfkong.DefaultTCPServicePort)
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayapi.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayapi.Listener{{
			Name:     gatewayTCPPortName,
			Protocol: gatewayapi.TCPProtocolType,
			Port:     gatewayapi.PortNumber(ktfkong.DefaultTCPServicePort),
		}}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("creating a tcpecho pod to test TCPRoute traffic routing")
	container1 := generators.NewContainer("tcpecho-1", test.EchoImage, test.EchoTCPPort)
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
	cleaner.Add(deployment1)

	t.Log("creating an additional tcpecho pod to test TCPRoute multiple backendRef loadbalancing")
	container2 := generators.NewContainer("tcpecho-2", test.EchoImage, test.EchoTCPPort)
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
	cleaner.Add(deployment2)

	t.Logf("exposing deployment %s/%s via service", deployment1.Namespace, deployment1.Name)
	service1 := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeLoadBalancer)
	// Use the same port as the default TCP port from the Kong Gateway deployment
	// to the tcpecho port, as this is what will be used to route the traffic at the Gateway.
	const service1Port = ktfkong.DefaultTCPServicePort
	service1.Spec.Ports = []corev1.ServicePort{{
		Name:       "tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       service1Port,
		TargetPort: intstr.FromInt(test.EchoTCPPort),
	}}
	service1, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service1, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(service1)

	t.Logf("exposing deployment %s/%s via service", deployment2.Namespace, deployment2.Name)
	const service2Port = 8080 // Use a different port that listening on the Gateway for TCP.
	service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	service2.Spec.Ports = []corev1.ServicePort{{
		Name:       "tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       service2Port,
		TargetPort: intstr.FromInt(test.EchoTCPPort),
	}}
	service2, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service2, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(service2)

	t.Logf("creating a TCPRoute to access deployment %s via kong", deployment1.Name)
	tcpRoute := &gatewayapi.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayapi.TCPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name:        gatewayapi.ObjectName(gatewayName),
					SectionName: lo.ToPtr(gatewayapi.SectionName(gatewayTCPPortName)),
				}},
			},
			Rules: []gatewayapi.TCPRouteRule{{
				BackendRefs: []gatewayapi.BackendRef{{
					BackendObjectReference: gatewayapi.BackendObjectReference{
						Name: gatewayapi.ObjectName(service1.Name),
						Port: lo.ToPtr(gatewayapi.PortNumber(service1Port)),
					},
				}},
			}},
		},
	}
	tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Create(ctx, tcpRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(tcpRoute)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)
	t.Log("verifying that the tcproute contains 'Programmed' condition")
	require.Eventually(t,
		GetVerifyProgrammedConditionCallback(t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name, metav1.ConditionTrue),
		ingressWait, waitTick,
	)

	t.Log("verifying that the tcpecho is responding properly")
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
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
	callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the tcpecho is no longer responding")
	defer func() {
		if t.Failed() {
			err := test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
			t.Logf("no longer responding check failure state: eof=%v, reset=%v, err=%v",
				errors.Is(err, io.EOF), errors.Is(err, syscall.ECONNRESET), err)
		}
	}()
	require.Eventually(t, func() bool {
		err := test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET)
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
	callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
	}, ingressWait, waitTick)

	t.Log("deleting the GatewayClass")
	require.NoError(t, gatewayClient.GatewayV1beta1().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TCPRoute gets dropped with the GatewayClass now removed")
	require.Eventually(t, func() bool {
		err := test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET)
	}, ingressWait, waitTick)

	t.Log("putting the GatewayClass back")
	gwc, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of TCPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
	}, ingressWait, waitTick)

	t.Log("deleting the Gateway")
	require.NoError(t, gatewayClient.GatewayV1().Gateways(ns.Name).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TCPRoute gets dropped with the Gateway now removed")
	require.Eventually(t, func() bool {
		err := test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET)
	}, ingressWait, waitTick)

	t.Log("putting the Gateway back")
	gateway, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayapi.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayapi.Listener{{
			Name:     gatewayTCPPortName,
			Protocol: gatewayapi.TCPProtocolType,
			Port:     gatewayapi.PortNumber(ktfkong.DefaultTCPServicePort),
		}}
	})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the Gateway again triggers reconciliation of TCPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, gatewayClient.GatewayV1().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, gatewayClient.GatewayV1().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TCPRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		err := test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET)
	}, ingressWait, waitTick)

	t.Log("putting the GatewayClass back")
	gwc, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)

	t.Log("putting the Gateway back")
	gateway, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayapi.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayapi.Listener{{
			Name:     gatewayTCPPortName,
			Protocol: gatewayapi.TCPProtocolType,
			Port:     gatewayapi.PortNumber(ktfkong.DefaultTCPServicePort),
		}}
	})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that creating the Gateway again triggers reconciliation of TCPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
	}, ingressWait, waitTick)

	t.Log("adding an additional backendRef to the TCPRoute")
	require.Eventually(t, func() bool {
		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Get(ctx, tcpRoute.Name, metav1.GetOptions{})
		require.NoError(t, err)

		tcpRoute.Spec.Rules[0].BackendRefs = []gatewayapi.BackendRef{
			{
				BackendObjectReference: gatewayapi.BackendObjectReference{
					Name: gatewayapi.ObjectName(service1.Name),
					Port: lo.ToPtr(gatewayapi.PortNumber(service1Port)),
				},
			},
			{
				BackendObjectReference: gatewayapi.BackendObjectReference{
					Name: gatewayapi.ObjectName(service2.Name),
					Port: lo.ToPtr(gatewayapi.PortNumber(service2Port)),
				},
			},
		}

		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Update(ctx, tcpRoute, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Log("verifying that the TCPRoute is now load-balanced between two services")
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2) == nil
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2) == nil
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, gatewayClient.GatewayV1().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, gatewayClient.GatewayV1().Gateways(ns.Name).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.TCPProtocolType, ns.Name, tcpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("verifying that the data-plane configuration from the TCPRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		err := test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET)
	}, ingressWait, waitTick)

	t.Log("testing port matching")
	t.Log("putting the Gateway back")
	_, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayapi.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayapi.Listener{{
			Name:     gatewayTCPPortName,
			Protocol: gatewayapi.TCPProtocolType,
			Port:     gatewayapi.PortNumber(ktfkong.DefaultTCPServicePort),
		}}
	})
	require.NoError(t, err)
	t.Log("putting the GatewayClass back")
	_, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)

	t.Log("verifying that the TCPRoute responds before specifying a port not existent in Gateway")
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
	}, ingressWait, waitTick)

	t.Log("setting the port in ParentRef which does not have a matching listener in Gateway")
	require.Eventually(t, func() bool {
		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Get(ctx, tcpRoute.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}
		notExistingPort := gatewayapi.PortNumber(81)
		tcpRoute.Spec.ParentRefs[0].Port = &notExistingPort
		tcpRoute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Update(ctx, tcpRoute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the TCPRoute does not respond after specifying a port not existent in Gateway")
	require.Eventually(t, func() bool {
		err := test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1)
		return errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET)
	}, ingressWait, waitTick)
}

func TestTCPRouteReferenceGrant(t *testing.T) {
	ctx := context.Background()
	RunWhenKongExpressionRouter(t)
	t.Log("locking TCP port")
	tcpMutex.Lock()
	t.Cleanup(func() {
		t.Log("unlocking TCP port")
		tcpMutex.Unlock()
	})

	ns, cleaner := helpers.Setup(ctx, t, env)

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
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayapi.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayapi.Listener{{
			Name:     gatewayTCPPortName,
			Protocol: gatewayapi.TCPProtocolType,
			Port:     gatewayapi.PortNumber(ktfkong.DefaultTCPServicePort),
		}}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("creating a tcpecho pod to test TCPRoute traffic routing")
	container1 := generators.NewContainer("tcpecho-1", test.EchoImage, test.EchoTCPPort)
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
	cleaner.Add(deployment1)

	t.Log("creating an additional tcpecho pod to test TCPRoute multiple backendRef loadbalancing")
	container2 := generators.NewContainer("tcpecho-2", test.EchoImage, test.EchoTCPPort)
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
	cleaner.Add(deployment2)

	t.Logf("exposing deployment %s/%s via service", deployment1.Namespace, deployment1.Name)
	service1 := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeLoadBalancer)
	// Use the same port as the default TCP port from the Kong Gateway deployment
	// to the tcpecho port, as this is what will be used to route the traffic at the Gateway.
	const service1Port = ktfkong.DefaultTCPServicePort
	service1.Spec.Ports = []corev1.ServicePort{{
		Name:       "tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       service1Port,
		TargetPort: intstr.FromInt(test.EchoTCPPort),
	}}
	service1, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service1, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(service1)

	t.Logf("exposing deployment %s/%s via service", deployment2.Namespace, deployment2.Name)
	const service2Port = 8080 // Use a different port that listening on the Gateway for TCP.
	service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	service2.Spec.Ports = []corev1.ServicePort{{
		Name:       "tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       service2Port,
		TargetPort: intstr.FromInt(test.EchoTCPPort),
	}}
	service2, err = env.Cluster().Client().CoreV1().Services(otherNs.Name).Create(ctx, service2, metav1.CreateOptions{})
	assert.NoError(t, err)
	cleaner.Add(service2)

	t.Logf("creating a tcproute to access deployment %s via kong", deployment1.Name)
	remoteNamespace := gatewayapi.Namespace(otherNs.Name)
	tcproute := &gatewayapi.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayapi.TCPRouteSpec{
			CommonRouteSpec: gatewayapi.CommonRouteSpec{
				ParentRefs: []gatewayapi.ParentReference{{
					Name:        gatewayapi.ObjectName(gatewayName),
					SectionName: lo.ToPtr(gatewayapi.SectionName(gatewayTCPPortName)),
				}},
			},
			Rules: []gatewayapi.TCPRouteRule{{
				BackendRefs: []gatewayapi.BackendRef{
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Name: gatewayapi.ObjectName(service1.Name),
							Port: lo.ToPtr(gatewayapi.PortNumber(service1Port)),
						},
					},
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Name:      gatewayapi.ObjectName(service2.Name),
							Namespace: &remoteNamespace,
							Port:      lo.ToPtr(gatewayapi.PortNumber(service2Port)),
						},
					},
				},
			}},
		},
	}
	tcproute, err = gatewayClient.GatewayV1alpha2().TCPRoutes(ns.Name).Create(ctx, tcproute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(tcproute)

	t.Log("verifying that only the local tcpecho is responding without a ReferenceGrant")
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
	}, ingressWait*2, waitTick)
	require.Never(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2) == nil
	}, time.Second*10, time.Second)

	t.Logf("creating a ReferenceGrant that permits tcproute access from %s to services in %s", ns.Name, otherNs.Name)
	grant := &gatewayapi.ReferenceGrant{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayapi.ReferenceGrantSpec{
			From: []gatewayapi.ReferenceGrantFrom{
				{
					// this isn't actually used, it's just a dummy extra from to confirm we handle multiple fine
					Group:     gatewayapi.Group("gateway.networking.k8s.io"),
					Kind:      gatewayapi.Kind("TCPRoute"),
					Namespace: gatewayapi.Namespace("garbage"),
				},
				{
					Group:     gatewayapi.Group("gateway.networking.k8s.io"),
					Kind:      gatewayapi.Kind("TCPRoute"),
					Namespace: gatewayapi.Namespace(tcproute.Namespace),
				},
			},
			To: []gatewayapi.ReferenceGrantTo{
				// also a dummy
				{
					Group: gatewayapi.Group(""),
					Kind:  gatewayapi.Kind("Pterodactyl"),
				},
				{
					Group: gatewayapi.Group(""),
					Kind:  gatewayapi.Kind("Service"),
				},
			},
		},
	}

	grant, err = gatewayClient.GatewayV1beta1().ReferenceGrants(otherNs.Name).Create(ctx, grant, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that requests reach both the local and remote namespace echo instances")
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID1) == nil
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2) == nil
	}, ingressWait, waitTick)

	t.Logf("testing specific name references")
	serviceName := gatewayapi.ObjectName(service2.ObjectMeta.Name)
	grant.Spec.To[1] = gatewayapi.ReferenceGrantTo{
		Kind:  gatewayapi.Kind("Service"),
		Group: gatewayapi.Group(""),
		Name:  &serviceName,
	}

	grant, err = gatewayClient.GatewayV1beta1().ReferenceGrants(otherNs.Name).Update(ctx, grant, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2) == nil
	}, ingressWait*2, waitTick)

	t.Logf("testing incorrect name does not match")
	blueguyName := gatewayapi.ObjectName("blueguy")
	grant.Spec.To[1].Name = &blueguyName
	_, err = gatewayClient.GatewayV1beta1().ReferenceGrants(otherNs.Name).Update(ctx, grant, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return test.TCPEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID2) != nil
	}, ingressWait, waitTick)
}
