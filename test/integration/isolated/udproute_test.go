//go:build integration_tests

package isolated

import (
	"context"
	"errors"
	"io"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

func TestUDPRouteEssentials(t *testing.T) {
	// Constants shared in many steps of this test that doesn't change.
	const gatewayUDPPortName = "udp"

	const service1Port = ktfkong.DefaultUDPServicePort
	const service1Name = "udpecho-1"
	test1UUID := uuid.NewString()

	const service2Name = "udpecho-2"
	const service2Port = 8080
	test2UUID := uuid.NewString()

	gatewayClassName := uuid.NewString()
	gatewayName := uuid.NewString()

	// Helpers used in this test.
	requireNoResponse := func(udpGatewayURL string) {
		t.Helper()
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			// For UDP lack of response (a timeout) means that we can't reach a service.
			err := test.EchoResponds(test.ProtocolUDP, udpGatewayURL, "irrelevant")
			assert.True(c, os.IsTimeout(err), "unexpected error: %v", err)
		}, consts.IngressWait, consts.WaitTick)
	}
	requireResponse := func(udpGatewayURL, expectedMsg string) {
		t.Helper()
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.NoError(c, test.EchoResponds(test.ProtocolUDP, udpGatewayURL, expectedMsg))
		}, consts.IngressWait, consts.WaitTick)
	}

	f := features.
		New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindUDPRoute).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		WithSetup("configure UDP Deployments with Services and UDPRoutes", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			gatewayClient, err := gatewayclient.NewForConfig(cfg.Client().RESTConfig())
			assert.NoError(t, err)
			ctx = SetInCtxForT(ctx, t, gatewayClient)

			cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
			namespace := GetNamespaceForT(ctx, t)
			cluster := GetClusterFromCtx(ctx)

			t.Log("deploying a supported gatewayclass to the test cluster")
			gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
			assert.NoError(t, err)
			cleaner.Add(gwc)

			t.Logf("deploying a gateway to the test cluster using unmanaged gateway mode and port %d", ktfkong.DefaultUDPServicePort)
			gateway, err := helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName, func(gw *gatewayapi.Gateway) {
				gw.Name = gatewayName
				gw.Spec.Listeners = []gatewayapi.Listener{{
					Name:     gatewayUDPPortName,
					Protocol: gatewayapi.UDPProtocolType,
					Port:     gatewayapi.PortNumber(ktfkong.DefaultUDPServicePort),
				}}
			})
			assert.NoError(t, err)
			cleaner.Add(gateway)

			t.Log("creating a udpecho pod to test UDPRoute traffic routing")
			container1 := generators.NewContainer(service1Name, test.EchoImage, test.EchoUDPPort)
			// App go-echo sends a "Running on Pod <UUID>." immediately on connecting.
			container1.Env = []corev1.EnvVar{
				{
					Name:  "POD_NAME",
					Value: test1UUID,
				},
			}
			deployment1 := generators.NewDeploymentForContainer(container1)
			deployment1, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment1, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(deployment1)

			t.Log("creating an additional udpecho pod to test UDPRoute multiple backendRef loadbalancing")
			container2 := generators.NewContainer(service2Name, test.EchoImage, test.EchoUDPPort)
			// App go-echo sends a "Running on Pod <UUID>." immediately on connecting.
			container2.Env = []corev1.EnvVar{
				{
					Name:  "POD_NAME",
					Value: test2UUID,
				},
			}
			deployment2 := generators.NewDeploymentForContainer(container2)
			deployment2, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment2, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(deployment2)

			t.Logf("exposing deployment %s/%s via service", deployment1.Namespace, deployment1.Name)
			service1 := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeClusterIP)
			service1.Name = service1Name
			// Use the same port as the default UDP port from the Kong Gateway deployment
			// to the udpecho port, as this is what will be used to route the traffic at the Gateway.
			service1.Spec.Ports = []corev1.ServicePort{{
				Name:       gatewayUDPPortName,
				Protocol:   corev1.ProtocolUDP,
				Port:       service1Port,
				TargetPort: intstr.FromInt(test.EchoUDPPort),
			}}
			service1, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service1, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(service1)

			t.Logf("exposing deployment %s/%s via service", deployment2.Namespace, deployment2.Name)
			service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeClusterIP)
			service2.Name = service2Name
			// Configure service to expose a different port than Gateway's UDP listener port (ktfkong.DefaultUDPServicePort)
			// to check whether traffic will be routed correctly.
			service2.Spec.Ports = []corev1.ServicePort{{
				Name:       gatewayUDPPortName,
				Protocol:   corev1.ProtocolUDP,
				Port:       service2Port,
				TargetPort: intstr.FromInt(test.EchoUDPPort),
			}}
			service2, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service2, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(service2)

			t.Logf("creating a UDPRoute to access deployment %s via kong", deployment1.Name)
			udpRoute := &gatewayapi.UDPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
				},
				Spec: gatewayapi.UDPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name:        gatewayapi.ObjectName(gatewayName),
							SectionName: lo.ToPtr(gatewayapi.SectionName(gatewayUDPPortName)),
						}},
					},
					Rules: []gatewayapi.UDPRouteRule{{
						BackendRefs: []gatewayapi.BackendRef{{
							BackendObjectReference: gatewayapi.BackendObjectReference{
								Name: gatewayapi.ObjectName(service1.Name),
								Port: lo.ToPtr(gatewayapi.PortNumber(service1Port)),
							},
						}},
					}},
				},
			}
			udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Create(ctx, udpRoute, metav1.CreateOptions{})
			assert.NoError(t, err)
			cleaner.Add(udpRoute)
			ctx = SetInCtxForT(ctx, t, udpRoute)

			return ctx
		}).
		Assess("basic test - route status and connectivity", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			t.Log("verifying that the Gateway gets linked to the route via status")
			gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			udpRoute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
			namespace := GetNamespaceForT(ctx, t)
			callback := helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)
			t.Log("verifying that the udproute contains 'Programmed' condition")
			assert.Eventually(t,
				helpers.GetVerifyProgrammedConditionCallback(t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name, metav1.ConditionTrue),
				consts.IngressWait, consts.WaitTick,
			)

			t.Log("verifying that the udpecho is responding properly")
			udpGatewayURL := GetUDPURLFromCtx(ctx)
			requireResponse(udpGatewayURL, test1UUID)

			return ctx
		}).
		Assess("verifying behavior when UDPRoute is modified", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			t.Log("removing the parentrefs from the UDPRoute")
			gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			udpRoute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
			namespace := GetNamespaceForT(ctx, t)
			udpGatewayURL := GetUDPURLFromCtx(ctx)

			oldParentRefs := udpRoute.Spec.ParentRefs
			assert.Eventually(t, func() bool {
				udpRoute, err := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Get(ctx, udpRoute.Name, metav1.GetOptions{})
				assert.NoError(t, err)
				udpRoute.Spec.ParentRefs = nil
				_, err = gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Update(ctx, udpRoute, metav1.UpdateOptions{})
				return err == nil
			}, time.Minute, time.Second)

			t.Log("verifying that the Gateway gets unlinked from the route via status")
			callback := helpers.GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

			t.Log("verifying that the udpecho is no longer responding")
			defer func() {
				if t.Failed() {
					err := test.EchoResponds(test.ProtocolUDP, udpGatewayURL, test1UUID)
					t.Logf("no longer responding check failure state: eof=%v, reset=%v, err=%v",
						errors.Is(err, io.EOF), errors.Is(err, syscall.ECONNRESET), err)
				}
			}()
			requireNoResponse(udpGatewayURL)

			t.Log("putting the parentRefs back")
			assert.Eventually(t, func() bool {
				udpRoute, err := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Get(ctx, udpRoute.Name, metav1.GetOptions{})
				assert.NoError(t, err)
				udpRoute.Spec.ParentRefs = oldParentRefs
				_, err = gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Update(ctx, udpRoute, metav1.UpdateOptions{})
				return err == nil
			}, time.Minute, time.Second)

			t.Log("verifying that the Gateway gets linked to the route via status")
			callback = helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

			t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
			requireResponse(udpGatewayURL, test1UUID)

			return ctx
		}).
		Assess("verifying behavior when Gateway is deleted and recreated", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			namespace := GetNamespaceForT(ctx, t)
			udpRoute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
			udpGatewayURL := GetUDPURLFromCtx(ctx)

			t.Log("deleting the GatewayClass")
			assert.NoError(t, gatewayClient.GatewayV1().GatewayClasses().Delete(ctx, gatewayClassName, metav1.DeleteOptions{}))

			t.Log("verifying that the Gateway gets unlinked from the route via status")
			callback := helpers.GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

			t.Log("verifying that the data-plane configuration from the UDPRoute gets dropped with the GatewayClass now removed")
			requireNoResponse(udpGatewayURL)

			t.Log("putting the GatewayClass back")
			gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
			assert.NoError(t, err)

			t.Log("verifying that the Gateway gets linked to the route via status")
			callback = helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

			t.Log("verifying that creating the GatewayClass again triggers reconciliation of UDPRoutes and the route becomes available again")
			requireResponse(udpGatewayURL, test1UUID)

			t.Log("deleting the Gateway")
			assert.NoError(t, gatewayClient.GatewayV1().Gateways(namespace).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

			t.Log("verifying that the Gateway gets unlinked from the route via status")
			callback = helpers.GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

			t.Log("verifying that the data-plane configuration from the UDPRoute gets dropped with the Gateway now removed")
			requireNoResponse(udpGatewayURL)

			t.Log("putting the Gateway back")
			_, err = helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName, func(gw *gatewayapi.Gateway) {
				gw.Name = gatewayName
				gw.Spec.Listeners = []gatewayapi.Listener{{
					Name:     gatewayUDPPortName,
					Protocol: gatewayapi.UDPProtocolType,
					Port:     gatewayapi.PortNumber(ktfkong.DefaultUDPServicePort),
				}}
			})
			assert.NoError(t, err)

			t.Log("verifying that the Gateway gets linked to the route via status")
			callback = helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

			t.Log("verifying that creating the Gateway again triggers reconciliation of UDPRoutes and the route becomes available again")
			requireResponse(udpGatewayURL, test1UUID)

			t.Log("deleting both GatewayClass and Gateway rapidly")
			assert.NoError(t, gatewayClient.GatewayV1().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
			assert.NoError(t, gatewayClient.GatewayV1().Gateways(namespace).Delete(ctx, gatewayName, metav1.DeleteOptions{}))

			t.Log("verifying that the Gateway gets unlinked from the route via status")
			callback = helpers.GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

			t.Log("verifying that the data-plane configuration from the UDPRoute does not get orphaned with the GatewayClass and Gateway gone")
			requireNoResponse(udpGatewayURL)

			t.Log("putting the GatewayClass back")
			_, err = helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
			assert.NoError(t, err)

			t.Log("putting the Gateway back")
			_, err = helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName, func(gw *gatewayapi.Gateway) {
				gw.Name = gatewayName
				gw.Spec.Listeners = []gatewayapi.Listener{{
					Name:     gatewayUDPPortName,
					Protocol: gatewayapi.UDPProtocolType,
					Port:     gatewayapi.PortNumber(ktfkong.DefaultUDPServicePort),
				}}
			})
			assert.NoError(t, err)

			t.Log("verifying that the Gateway gets linked to the route via status")
			callback = helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udpRoute.Name)
			assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

			t.Log("verifying that creating the Gateway again triggers reconciliation of UDPRoutes and the route becomes available again")
			requireResponse(udpGatewayURL, test1UUID)

			return ctx
		}).
		Assess("verifying behavior with many backends", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
			namespace := GetNamespaceForT(ctx, t)
			udpRoute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
			udpGatewayURL := GetUDPURLFromCtx(ctx)

			t.Log("adding an additional backendRef to the UDPRoute")
			assert.Eventually(t, func() bool {
				udpRoute, err := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Get(ctx, udpRoute.Name, metav1.GetOptions{})
				assert.NoError(t, err)
				udpRoute.Spec.Rules[0].BackendRefs = []gatewayapi.BackendRef{
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Name: gatewayapi.ObjectName(service1Name),
							Port: lo.ToPtr(gatewayapi.PortNumber(service1Port)),
						},
					},
					{
						BackendObjectReference: gatewayapi.BackendObjectReference{
							Name: gatewayapi.ObjectName(service2Name),
							Port: lo.ToPtr(gatewayapi.PortNumber(service2Port)),
						},
					},
				}

				_, err = gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Update(ctx, udpRoute, metav1.UpdateOptions{})
				return err == nil
			}, consts.IngressWait, consts.WaitTick)

			t.Log("verifying that the UDPRoute is now load-balanced between two services")
			requireResponse(udpGatewayURL, test1UUID)
			requireResponse(udpGatewayURL, test2UUID)

			t.Log("testing port matching")
			t.Log("putting the Gateway back")
			_, err := helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName, func(gw *gatewayapi.Gateway) {
				gw.Name = gatewayName
				gw.Spec.Listeners = []gatewayapi.Listener{{
					Name:     gatewayUDPPortName,
					Protocol: gatewayapi.UDPProtocolType,
					Port:     gatewayapi.PortNumber(ktfkong.DefaultUDPServicePort),
				}}
			})
			assert.NoError(t, err)
			t.Log("putting the GatewayClass back")
			_, err = helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
			assert.NoError(t, err)

			t.Log("verifying that the UDPRoute responds before specifying a port not existent in Gateway")
			requireResponse(udpGatewayURL, test1UUID)

			t.Log("setting the port in ParentRef which does not have a matching listener in Gateway")
			assert.Eventually(t, func() bool {
				udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Get(ctx, udpRoute.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				notExistingPort := gatewayapi.PortNumber(81)
				udpRoute.Spec.ParentRefs[0].Port = &notExistingPort
				udpRoute.Spec.ParentRefs[0].Name = gatewayapi.ObjectName(service1Name)
				udpRoute, err = gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Update(ctx, udpRoute, metav1.UpdateOptions{})
				return err == nil
			}, time.Minute, time.Second)

			t.Log("verifying that the UDPRoute does not respond after specifying a port not existent in Gateway")
			requireNoResponse(udpGatewayURL)
			return ctx
		}).
		Teardown(featureTeardown())

	tenv.Test(t, f.Feature())
}
