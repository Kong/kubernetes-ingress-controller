//go:build integration_tests

package integration

import (
	"context"
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

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

func TestTCPRouteReferenceGrant(t *testing.T) {
	ctx := context.Background()
	RunWhenKongExpressionRouter(context.Background(), t)
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
	gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	const gatewayTCPPortName = "tcp"

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	gatewayName := uuid.NewString()
	gateway, err := helpers.DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayapi.Gateway) {
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
		return test.EchoResponds(test.ProtocolTCP, proxyTCPURL, testUUID1) == nil
	}, ingressWait*2, waitTick)
	require.Never(t, func() bool {
		return test.EchoResponds(test.ProtocolTCP, proxyTCPURL, testUUID2) == nil
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
		return test.EchoResponds(test.ProtocolTCP, proxyTCPURL, testUUID1) == nil
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		return test.EchoResponds(test.ProtocolTCP, proxyTCPURL, testUUID2) == nil
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
		return test.EchoResponds(test.ProtocolTCP, proxyTCPURL, testUUID2) == nil
	}, ingressWait*2, waitTick)

	t.Logf("testing incorrect name does not match")
	blueguyName := gatewayapi.ObjectName("blueguy")
	grant.Spec.To[1].Name = &blueguyName
	_, err = gatewayClient.GatewayV1beta1().ReferenceGrants(otherNs.Name).Update(ctx, grant, metav1.UpdateOptions{})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return test.EchoResponds(test.ProtocolTCP, proxyTCPURL, testUUID2) != nil
	}, ingressWait, waitTick)
}
