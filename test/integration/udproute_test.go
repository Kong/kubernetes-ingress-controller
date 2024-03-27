//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
)

const testdomain = "konghq.com"

func TestUDPRouteEssentials(t *testing.T) {
	ctx := context.Background()
	RunWhenKongExpressionRouterWithVersion(t, ">=3.4.0")
	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("locking UDP port")
	udpMutex.Lock()
	defer func() {
		// Free up the UDP port
		t.Log("unlocking UDP port")
		udpMutex.Unlock()
	}()

	t.Log("deploying a supported gatewayclass to the test cluster")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a supported gatewayclass to the test cluster")
	gatewayClassName := uuid.NewString()
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode and port 9999")
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1.Listener{{
			Name:     "udp",
			Protocol: gatewayv1.UDPProtocolType,
			Port:     gatewayv1.PortNumber(ktfkong.DefaultUDPServicePort),
		}}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("configuring coredns corefile")
	cfgmap1 := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "coredns"}, Data: map[string]string{"Corefile": corefile}}
	cfgmap1, err = env.Cluster().Client().CoreV1().ConfigMaps(ns.Name).Create(ctx, cfgmap1, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("configuring alternative coredns corefile for load-balanced setup")
	alternativeCorefile := strings.Replace(corefile, "10.0.0.1 konghq.com", "10.0.0.2 konghq.com", -1)
	cfgmap2 := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "coredns2"}, Data: map[string]string{"Corefile": alternativeCorefile}}
	cfgmap2, err = env.Cluster().Client().CoreV1().ConfigMaps(ns.Name).Create(ctx, cfgmap2, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the coredns corefiles %s/%s and %s/%s", cfgmap1.Namespace, cfgmap1.Name, cfgmap2.Namespace, cfgmap2.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().ConfigMaps(ns.Name).Delete(ctx, cfgmap1.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().CoreV1().ConfigMaps(ns.Name).Delete(ctx, cfgmap2.Name, metav1.DeleteOptions{}))
	}()

	t.Log("configuring a coredns deployent to deploy for UDP testing")
	container1 := generators.NewContainer("coredns", coreDNSImage, ktfkong.DefaultUDPServicePort)
	container1.Ports[0].Protocol = corev1.ProtocolUDP
	container1.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
	container1.Args = []string{"-conf", "/etc/coredns/Corefile"}
	deployment1 := generators.NewDeploymentForContainer(container1)

	t.Log("configuring the coredns pod with a custom corefile")
	configVolume1 := corev1.Volume{
		Name: "config-volume",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: cfgmap1.Name},
			Items:                []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}},
		}},
	}
	deployment1.Spec.Template.Spec.Volumes = append(deployment1.Spec.Template.Spec.Volumes, configVolume1)

	t.Log("deploying coredns")
	deployment1, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment1, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("configuring alternative coredns deployent for load-balanced UDP testing")
	container2 := generators.NewContainer("coredns2", coreDNSImage, ktfkong.DefaultUDPServicePort)
	container2.Ports[0].Protocol = corev1.ProtocolUDP
	container2.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
	container2.Args = []string{"-conf", "/etc/coredns/Corefile"}
	deployment2 := generators.NewDeploymentForContainer(container2)
	deployment2.ObjectMeta.Name = "coredns2"

	t.Log("configuring the coredns pod with a custom corefile")
	configVolume2 := corev1.Volume{
		Name: "config-volume",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: cfgmap2.Name},
			Items:                []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}},
		}},
	}
	deployment2.Spec.Template.Spec.Volumes = append(deployment2.Spec.Template.Spec.Volumes, configVolume2)

	t.Log("deploying alternative coredns")
	deployment2, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment2, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up deployments %s/%s and %s/%s", deployment1.Namespace, deployment1.Name, deployment2.Namespace, deployment2.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment1.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment2.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s/%s via service", deployment1.Namespace, deployment1.Name)
	service := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing alternative deployment %s/%s via service", deployment1.Namespace, deployment1.Name)
	service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	service2, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service2, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up services %s/%s and %s/%s", service.Namespace, service.Name, service2.Namespace, service2.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service2.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating a udproute to access deployment %s via kong", deployment1.Name)
	udpPortDefault := gatewayv1alpha2.PortNumber(ktfkong.DefaultUDPServicePort)
	udpRoute := &gatewayv1alpha2.UDPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.UDPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gatewayName),
				}},
			},
			Rules: []gatewayv1alpha2.UDPRouteRule{{
				BackendRefs: []gatewayv1alpha2.BackendRef{{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service.Name),
						Port: &udpPortDefault,
					},
				}},
			}},
		},
	}

	t.Log("configurating a net.Resolver to resolve DNS via the proxy")
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:%d", proxyUDPURL.Hostname(), ktfkong.DefaultUDPServicePort))
		},
	}

	udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Create(ctx, udpRoute, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the udproute %s", udpRoute.Name)
		if err := gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Delete(ctx, udpRoute.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayv1.UDPProtocolType, ns.Name, udpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)
	t.Log("verifying that the udproute contains 'Programmed' condition")
	require.Eventually(t,
		GetVerifyProgrammedConditionCallback(t, gatewayClient, gatewayv1.UDPProtocolType, ns.Name, udpRoute.Name, metav1.ConditionTrue),
		ingressWait, waitTick,
	)

	t.Logf("checking DNS to resolve via UDPIngress %s", udpRoute.Name)
	require.Eventually(t, func() bool {
		_, err := resolver.LookupHost(ctx, "kernel.org")
		return err == nil
	}, ingressWait, waitTick)

	t.Run("removing parentRefs", func(t *testing.T) {
		t.Log("removing the parentrefs from the UDPRoute")
		oldParentRefs := udpRoute.Spec.ParentRefs
		require.Eventually(t, func() bool {
			udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Get(ctx, udpRoute.Name, metav1.GetOptions{})
			require.NoError(t, err)
			udpRoute.Spec.ParentRefs = nil
			udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Update(ctx, udpRoute, metav1.UpdateOptions{})
			return err == nil
		}, time.Minute, time.Second)

		t.Log("verifying that the Gateway gets unlinked from the route via status")
		callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayv1.UDPProtocolType, ns.Name, udpRoute.Name)
		require.Eventually(t, callback, ingressWait, waitTick)

		t.Log("verifying that the data-plane configuration from the UDPRoute gets dropped with the parentRefs now removed")
		// negative checks for these tests check that DNS queries eventually start to fail, presumably because they time
		// out. we assume there shouldn't be unrelated failure reasons because they always follow a test that confirm
		// resolution was working before. we can't use never here because there may be some delay in deleting the route
		require.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			return err != nil
		}, ingressWait, waitTick)

		t.Log("putting the parentRefs back")
		require.Eventually(t, func() bool {
			udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Get(ctx, udpRoute.Name, metav1.GetOptions{})
			require.NoError(t, err)
			udpRoute.Spec.ParentRefs = oldParentRefs
			udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Update(ctx, udpRoute, metav1.UpdateOptions{})
			return err == nil
		}, time.Minute, time.Second)

		t.Log("verifying that the Gateway gets linked to the route via status")
		callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayv1.UDPProtocolType, ns.Name, udpRoute.Name)
		require.Eventually(t, callback, ingressWait, waitTick)

		t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
		t.Logf("checking DNS to resolve via UDPIngress %s", udpRoute.Name)
		require.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			return err == nil
		}, ingressWait, waitTick)
	})

	t.Run("removing GatewayClass", func(t *testing.T) {
		t.Log("deleting the GatewayClass")
		require.NoError(t, gatewayClient.GatewayV1beta1().GatewayClasses().Delete(ctx, gatewayClassName, metav1.DeleteOptions{}))

		t.Log("verifying that the Gateway gets unlinked from the route via status")
		callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayv1.UDPProtocolType, ns.Name, udpRoute.Name)
		require.Eventually(t, callback, ingressWait, waitTick)

		t.Log("verifying that the data-plane configuration from the UDPRoute gets dropped with the GatewayClass now removed")
		require.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			return err != nil
		}, ingressWait, waitTick)

		t.Log("putting the GatewayClass back")
		_, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
		require.NoError(t, err)

		t.Log("verifying that the Gateway gets linked to the route via status")
		callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayv1.UDPProtocolType, ns.Name, udpRoute.Name)
		require.Eventually(t, callback, ingressWait, waitTick)

		t.Log("verifying that creating the GatewayClass again triggers reconciliation of UDPRoutes and the route becomes available again")
		t.Logf("checking DNS to resolve via UDPIngress %s", udpRoute.Name)
		require.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			return err == nil
		}, ingressWait, waitTick)
	})

	t.Run("removing Gateway", func(t *testing.T) {
		t.Log("deleting the Gateway")
		require.NoError(t, gatewayClient.GatewayV1beta1().Gateways(ns.Name).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

		t.Log("verifying that the Gateway gets unlinked from the route via status")
		callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayv1.UDPProtocolType, ns.Name, udpRoute.Name)
		require.Eventually(t, callback, ingressWait, waitTick)

		t.Log("verifying that the data-plane configuration from the UDPRoute gets dropped with the Gateway now removed")
		require.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			return err != nil
		}, ingressWait, waitTick)

		t.Log("putting the Gateway back")
		_, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1.Gateway) {
			gw.Name = gatewayName
			gw.Spec.Listeners = []gatewayv1.Listener{{
				Name:     "udp",
				Protocol: gatewayv1.UDPProtocolType,
				Port:     gatewayv1.PortNumber(ktfkong.DefaultUDPServicePort),
			}}
		})
		require.NoError(t, err)

		t.Log("verifying that the Gateway gets linked to the route via status")
		callback = GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayv1.UDPProtocolType, ns.Name, udpRoute.Name)
		require.Eventually(t, callback, ingressWait, waitTick)

		t.Log("verifying that creating the Gateway again triggers reconciliation of UDPRoutes and the route becomes available again")
		require.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			return err == nil
		}, ingressWait, waitTick)
	})

	t.Run("removing Gateway and GatewayClass simultaneously", func(t *testing.T) {
		t.Log("deleting both GatewayClass and Gateway")
		require.NoError(t, gatewayClient.GatewayV1beta1().GatewayClasses().Delete(ctx, gatewayClassName, metav1.DeleteOptions{}))
		require.NoError(t, gatewayClient.GatewayV1beta1().Gateways(ns.Name).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

		t.Log("verifying that the Gateway gets unlinked from the route via status")
		callback = GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayv1.UDPProtocolType, ns.Name, udpRoute.Name)
		require.Eventually(t, callback, ingressWait, waitTick)

		t.Log("verifying that the data-plane configuration from the UDPRoute does not get orphaned with the GatewayClass and Gateway gone")
		require.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			return err != nil
		}, ingressWait, waitTick)

		t.Log("putting the Gateway back")
		_, err = DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1.Gateway) {
			gw.Name = gatewayName
			gw.Spec.Listeners = []gatewayv1.Listener{{
				Name:     "udp",
				Protocol: gatewayv1.UDPProtocolType,
				Port:     gatewayv1.PortNumber(ktfkong.DefaultUDPServicePort),
			}}
		})
		require.NoError(t, err)
		t.Log("putting the GatewayClass back")
		_, err = DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
		require.NoError(t, err)

		t.Log("verifying that the UDPRoute responds")
		require.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			return err == nil
		}, ingressWait, waitTick)
	})

	t.Run("multiple backends", func(t *testing.T) {
		t.Log("adding another backendRef to load-balance the DNS between multiple CoreDNS pods")
		require.Eventually(t, func() bool {
			udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Get(ctx, udpRoute.Name, metav1.GetOptions{})
			require.NoError(t, err)

			udpRoute.Spec.Rules[0].BackendRefs = []gatewayv1alpha2.BackendRef{
				{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service.Name),
						Port: &udpPortDefault,
					},
				},
				{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service2.Name),
						Port: &udpPortDefault,
					},
				},
			}

			udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Update(ctx, udpRoute, metav1.UpdateOptions{})
			return err == nil
		}, ingressWait, waitTick)

		t.Log("verifying that DNS queries are being load-balanced between multiple CoreDNS pods")
		require.Eventually(t, func() bool { return isDNSResolverReturningExpectedResult(ctx, resolver, testdomain, "10.0.0.1") }, ingressWait, waitTick)
		require.Eventually(t, func() bool { return isDNSResolverReturningExpectedResult(ctx, resolver, testdomain, "10.0.0.2") }, ingressWait, waitTick)
		require.Eventually(t, func() bool { return isDNSResolverReturningExpectedResult(ctx, resolver, testdomain, "10.0.0.1") }, ingressWait, waitTick)
		require.Eventually(t, func() bool { return isDNSResolverReturningExpectedResult(ctx, resolver, testdomain, "10.0.0.2") }, ingressWait, waitTick)
	})

	t.Run("port matching", func(t *testing.T) {
		t.Log("updating UDPRoute parentRef to use a port not in the Gateway Listeners")
		require.Eventually(t, func() bool {
			udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Get(ctx, udpRoute.Name, metav1.GetOptions{})
			if err != nil {
				return false
			}
			notExistingPort := gatewayv1alpha2.PortNumber(81)
			udpRoute.Spec.ParentRefs[0].Port = &notExistingPort
			udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(ns.Name).Update(ctx, udpRoute, metav1.UpdateOptions{})
			return err == nil
		}, time.Minute, time.Second)

		t.Log("verifying that the UDPRoute becomes inactive")
		require.Eventually(t, func() bool {
			_, err := resolver.LookupHost(ctx, "kernel.org")
			return err != nil
		}, ingressWait, waitTick)
	})
}

func isDNSResolverReturningExpectedResult(ctx context.Context, resolver *net.Resolver, host, addr string) bool { //nolint:unparam
	addrs, err := resolver.LookupHost(ctx, host)
	if err != nil {
		return false
	}
	if len(addrs) != 1 {
		return false
	}
	return addrs[0] == addr
}
