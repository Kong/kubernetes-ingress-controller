//go:build integration_tests

package isolated

import (
	"context"
	"net"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/test/integration/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/testlabels"
)

const (
	// coreDNSImage is the image and version of CoreDNS that will be used for UDP testing.
	coreDNSImage = "registry.k8s.io/coredns/coredns:v1.8.6"
)

func TestUDPRoute(t *testing.T) {
	const (
		corefile = `
	.:53 {
		errors
		health
		ready
		kubernetes cluster.local in-addr.arpa ip6.arpa {
		   pods insecure
		   fallthrough in-addr.arpa ip6.arpa
		   ttl 5
		}
		forward . /etc/resolv.conf {
		   max_concurrent 1000
		}
		cache 1
		loop
		reload
		loadbalance
		hosts {
		  10.0.0.1 konghq.com
		  fallthrough
		}
	}
	.:9999 {
		errors
		health
		ready
		kubernetes cluster.local in-addr.arpa ip6.arpa {
		   pods insecure
		   fallthrough in-addr.arpa ip6.arpa
		   ttl 5
		}
		forward . /etc/resolv.conf {
		   max_concurrent 1000
		}
		cache 1
		loop
		reload
		loadbalance
		hosts {
		  10.0.0.1 konghq.com
		  fallthrough
		}
	}
	`
		testdomain = "konghq.com"
	)

	t.Parallel()

	var udprouteParentRefs []gatewayapi.ParentReference

	fEssentials := features.
		New("essentials").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindUDPRoute).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		WithSetup("prepare Gateway and GatewayClass",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				t.Log("creating Gateway API client")
				gatewayClient, err := gatewayclient.NewForConfig(cfg.Client().RESTConfig())
				assert.NoError(t, err, "failed creating Gateway API client")
				ctx = SetInCtxForT(ctx, t, gatewayClient)

				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)

				gatewayClassName := uuid.NewString()
				t.Logf("deploying a supported GatewayClass %s to the test cluster", gatewayClassName)
				gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
				assert.NoError(t, err)
				cleaner.Add(gwc)
				ctx = SetInCtxForT(ctx, t, gwc)

				gatewayName := uuid.NewString()
				t.Logf("deploying a Gateway %s to the test cluster using unmanaged gateway mode and port %d", gatewayName, ktfkong.DefaultUDPServicePort)
				namespace := GetNamespaceForT(ctx, t)
				gateway, err := helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName, func(gw *gatewayapi.Gateway) {
					gw.Name = gatewayName
					gw.Spec.Listeners = builder.NewListener("udp").
						UDP().
						WithPort(ktfkong.DefaultUDPServicePort).
						IntoSlice()
				})
				assert.NoError(t, err)
				cleaner.Add(gateway)
				ctx = SetInCtxForT(ctx, t, gateway)

				return ctx
			}).
		WithSetup("prepare coredns deployments",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				namespace := GetNamespaceForT(ctx, t)
				cl := cfg.Client().Resources()
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)

				t.Log("configuring coredns corefile")
				cfgmap1 := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "coredns",
						Namespace: namespace,
					},
					Data: map[string]string{"Corefile": corefile},
				}
				assert.NoError(t, cl.Create(ctx, cfgmap1))
				cleaner.Add(cfgmap1)

				t.Log("configuring alternative coredns corefile for load-balanced setup")
				alternativeCorefile := strings.Replace(corefile, "10.0.0.1 konghq.com", "10.0.0.2 konghq.com", -1)
				cfgmap2 := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "coredns2",
						Namespace: namespace,
					},
					Data: map[string]string{"Corefile": alternativeCorefile},
				}
				assert.NoError(t, cl.Create(ctx, cfgmap2))
				cleaner.Add(cfgmap2)

				t.Log("configuring a coredns deployent to deploy for UDP testing")
				container1 := generators.NewContainer("coredns", coreDNSImage, ktfkong.DefaultUDPServicePort)
				container1.Ports[0].Protocol = corev1.ProtocolUDP
				container1.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
				container1.Args = []string{"-conf", "/etc/coredns/Corefile"}
				deployment1 := generators.NewDeploymentForContainer(container1)

				t.Log("configuring the coredns pod with a custom corefile")
				deployment1.Spec.Template.Spec.Volumes = append(deployment1.Spec.Template.Spec.Volumes,
					corev1.Volume{
						Name: "config-volume",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: cfgmap1.Name,
							},
							Items: []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}},
						}},
					})
				deployment1.Namespace = namespace

				t.Logf("deploying coredns deployment %q", deployment1.Name)
				assert.NoError(t, cl.Create(ctx, deployment1))
				cleaner.Add(deployment1)

				t.Logf("exposing deployment %s/%s via service", deployment1.Namespace, deployment1.Name)
				service1 := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeLoadBalancer)
				service1.Namespace = namespace
				service1.Labels = map[string]string{"app": "coredns"}
				assert.NoError(t, cl.Create(ctx, service1))
				cleaner.Add(service1)

				t.Log("configuring alternative coredns deployent for load-balanced UDP testing")
				container2 := generators.NewContainer("coredns2", coreDNSImage, ktfkong.DefaultUDPServicePort)
				container2.Ports[0].Protocol = corev1.ProtocolUDP
				container2.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
				container2.Args = []string{"-conf", "/etc/coredns/Corefile"}
				deployment2 := generators.NewDeploymentForContainer(container2)
				deployment2.Name = "coredns2"

				t.Log("configuring the coredns pod with a custom corefile")
				deployment2.Spec.Template.Spec.Volumes = append(deployment2.Spec.Template.Spec.Volumes,
					corev1.Volume{
						Name: "config-volume",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: cfgmap2.Name,
							},
							Items: []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}},
						}},
					})
				deployment2.Namespace = namespace

				t.Logf("deploying coredns deployment %q", deployment2.Name)
				assert.NoError(t, cl.Create(ctx, deployment2))
				cleaner.Add(deployment2)

				t.Logf("exposing alternative deployment %s/%s via service", deployment2.Namespace, deployment2.Name)
				service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
				service2.Namespace = namespace
				service2.Labels = map[string]string{"app": "coredns"}
				assert.NoError(t, cl.Create(ctx, service2))
				cleaner.Add(service2)

				t.Logf("creating a UDPRoute to access deployment %s via kong", deployment1.Name)
				gateway := GetFromCtxForT[*gatewayapi.Gateway](ctx, t)
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				udproute := &gatewayapi.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: namespace,
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{{
								Name: gatewayapi.ObjectName(gateway.Name),
							}},
						},
						Rules: []gatewayapi.UDPRouteRule{{
							BackendRefs: builder.NewBackendRef(service1.Name).WithPort(ktfkong.DefaultUDPServicePort).ToSlice(),
						}},
					},
				}
				udproute, err := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Create(ctx, udproute, metav1.CreateOptions{})
				assert.NoError(t, err)
				cleaner.Add(udproute)
				ctx = SetInCtxForT(ctx, t, udproute)

				return ctx
			}).
		Assess("Gateway gets linked to the UDPRoute via status",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				namespace := GetNamespaceForT(ctx, t)
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)

				callback := helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udproute.Name)
				assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)
				t.Log("verifying that the UDPRoute contains 'Programmed' condition")
				assert.Eventually(t,
					helpers.GetVerifyProgrammedConditionCallback(t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udproute.Name, metav1.ConditionTrue),
					consts.IngressWait, consts.WaitTick,
				)

				return ctx
			}).
		Assess("DNS lookups work",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
				proxyUDPURL := GetUDPURLFromCtx(ctx)

				t.Logf("checking DNS to resolve via UDPIngress %s", udproute.Name)
				assert.Eventually(t, urlResolvesSuccessfullyFn(ctx, proxyUDPURL), consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Assess("removing UDPRoute parentRef removes the configuration from the Gateway",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				namespace := GetNamespaceForT(ctx, t)
				udprouteClient := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace)
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
				proxyUDPURL := GetUDPURLFromCtx(ctx)

				t.Log("removing the parentRefs from the UDPRoute")
				udprouteParentRefs = udproute.Spec.ParentRefs
				assert.Eventually(t, func() bool {
					udproute, err := udprouteClient.Get(ctx, udproute.Name, metav1.GetOptions{})
					if err != nil {
						return false
					}
					udproute.Spec.ParentRefs = nil
					_, err = udprouteClient.Update(ctx, udproute, metav1.UpdateOptions{})
					return err == nil
				}, consts.IngressWait, consts.WaitTick)
				ctx = SetInCtxForT(ctx, t, udproute)

				t.Log("verifying that the Gateway gets unlinked from the route via status")
				callback := helpers.GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udproute.Name)
				assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

				t.Log("verifying that the data-plane configuration from the UDPRoute gets dropped with the parentRefs now removed")
				// negative checks for these tests check that DNS queries eventually start to fail, presumably because they time
				// out. we assume there shouldn't be unrelated failure reasons because they always follow a test that confirm
				// resolution was working before. we can't use never here because there may be some delay in deleting the route
				assert.Eventually(t, not(urlResolvesSuccessfullyFn(ctx, proxyUDPURL)), consts.IngressWait, consts.WaitTick)
				return ctx
			}).
		Assess("restoring UDPRoute parentRef bring the configuration back",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				namespace := GetNamespaceForT(ctx, t)
				udprouteClient := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace)
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
				proxyUDPURL := GetUDPURLFromCtx(ctx)

				t.Log("putting the parentRefs back")
				assert.Eventually(t, func() bool {
					udproute, err := udprouteClient.Get(ctx, udproute.Name, metav1.GetOptions{})
					if err != nil {
						return false
					}
					udproute.Spec.ParentRefs = udprouteParentRefs
					_, err = udprouteClient.Update(ctx, udproute, metav1.UpdateOptions{})
					return err == nil
				}, consts.IngressWait, consts.WaitTick)

				t.Log("verifying that the Gateway gets linked to the route via status")
				callback := helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udproute.Name)
				assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

				t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
				t.Logf("checking DNS to resolve via UDPRoute %s", udproute.Name)
				assert.Eventually(t, urlResolvesSuccessfullyFn(ctx, proxyUDPURL), consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Assess("removing the GatewayClass unlinks the UDPRoute",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
				proxyUDPURL := GetUDPURLFromCtx(ctx)
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				namespace := GetNamespaceForT(ctx, t)
				gatewayclass := GetFromCtxForT[*gatewayapi.GatewayClass](ctx, t)

				t.Logf("deleting the GatewayClass %s", gatewayclass.Name)
				assert.NoError(t, gatewayClient.GatewayV1().GatewayClasses().Delete(ctx, gatewayclass.Name, metav1.DeleteOptions{}))

				t.Log("verifying that the Gateway gets unlinked from the route via status")
				callback := helpers.GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udproute.Name)
				assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

				t.Log("verifying that the data-plane configuration from the UDPRoute gets dropped with the GatewayClass now removed")
				assert.Eventually(t, not(urlResolvesSuccessfullyFn(ctx, proxyUDPURL)), consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Assess("putting back the GatewayClass restores the configuration",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
				proxyUDPURL := GetUDPURLFromCtx(ctx)
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				namespace := GetNamespaceForT(ctx, t)
				gatewayclass := GetFromCtxForT[*gatewayapi.GatewayClass](ctx, t)
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)

				t.Logf("putting the GatewayClass %s back", gatewayclass.Name)
				gatewayclass, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayclass.Name)
				assert.NoError(t, err)
				cleaner.Add(gatewayclass)
				ctx = SetInCtxForT(ctx, t, gatewayclass)

				t.Log("verifying that the Gateway gets linked to the route via status")
				callback := helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udproute.Name)
				assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

				t.Log("verifying that creating the GatewayClass again triggers reconciliation of UDPRoutes and the route becomes available again")
				t.Logf("checking DNS to resolve via UDPRoute %s", udproute.Name)
				assert.Eventually(t, urlResolvesSuccessfullyFn(ctx, proxyUDPURL), consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Assess("removing the Gateway removes the link to UDPRoute",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
				proxyUDPURL := GetUDPURLFromCtx(ctx)
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				namespace := GetNamespaceForT(ctx, t)
				gateway := GetFromCtxForT[*gatewayapi.Gateway](ctx, t)

				t.Log("deleting the Gateway")
				assert.NoError(t, gatewayClient.GatewayV1().Gateways(namespace).Delete(ctx, gateway.Name, metav1.DeleteOptions{}))

				t.Log("verifying that the Gateway gets unlinked from the route via status")
				callback := helpers.GetGatewayIsUnlinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udproute.Name)
				assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

				t.Log("verifying that the data-plane configuration from the UDPRoute gets dropped with the Gateway now removed")
				assert.Eventually(t, not(urlResolvesSuccessfullyFn(ctx, proxyUDPURL)), consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Assess("putting the Gateway back brings back the configuration",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				namespace := GetNamespaceForT(ctx, t)
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
				udprouteClient := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace)
				proxyUDPURL := GetUDPURLFromCtx(ctx)
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				gatewayclass := GetFromCtxForT[*gatewayapi.GatewayClass](ctx, t)

				t.Log("putting the Gateway back")
				gateway, err := helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayclass.Name, func(gw *gatewayapi.Gateway) {
					gw.Name = uuid.NewString()
					gw.Spec.Listeners = builder.NewListener("udp").
						UDP().
						WithPort(ktfkong.DefaultUDPServicePort).
						IntoSlice()
				})
				assert.NoError(t, err)
				cleaner.Add(gateway)
				ctx = SetInCtxForT(ctx, t, gateway)

				t.Log("update the UDPRoute with new Gateway ref")
				udproute, err = udprouteClient.Get(ctx, udproute.Name, metav1.GetOptions{})
				assert.NoError(t, err)
				udproute.Spec.CommonRouteSpec = gatewayapi.CommonRouteSpec{
					ParentRefs: []gatewayapi.ParentReference{{
						Name: gatewayapi.ObjectName(gateway.Name),
					}},
				}
				udproute, err = udprouteClient.Update(ctx, udproute, metav1.UpdateOptions{})
				assert.NoError(t, err)
				ctx = SetInCtxForT(ctx, t, udproute)

				t.Log("verifying that the Gateway gets linked to the route via status")
				callback := helpers.GetGatewayIsLinkedCallback(ctx, t, gatewayClient, gatewayapi.UDPProtocolType, namespace, udproute.Name)
				assert.Eventually(t, callback, consts.IngressWait, consts.WaitTick)

				t.Log("verifying that creating the Gateway again triggers reconciliation of UDPRoutes and the route becomes available again")
				assert.Eventually(t, urlResolvesSuccessfullyFn(ctx, proxyUDPURL), consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Assess("multiple backends load balance requests",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				namespace := GetNamespaceForT(ctx, t)
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
				udprouteClient := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace)
				proxyUDPURL := GetUDPURLFromCtx(ctx)

				t.Log("adding another backendRef to load-balance the DNS between multiple CoreDNS pods")
				var services corev1.ServiceList
				assert.NoError(t, cfg.Client().Resources(namespace).List(ctx, &services, func(lo *metav1.ListOptions) { lo.LabelSelector = "app=coredns" }))
				if !assert.Len(t, services.Items, 2) {
					return ctx
				}

				assert.Eventually(t, func() bool {
					udproute, err := udprouteClient.Get(ctx, udproute.Name, metav1.GetOptions{})
					if err != nil {
						return false
					}

					udproute.Spec.Rules[0].BackendRefs = nil
					for _, svc := range services.Items {
						svc := svc
						udproute.Spec.Rules[0].BackendRefs = append(udproute.Spec.Rules[0].BackendRefs,
							builder.NewBackendRef(svc.Name).WithPort(ktfkong.DefaultUDPServicePort).Build(),
						)
					}

					_, err = udprouteClient.Update(ctx, udproute, metav1.UpdateOptions{})
					return err == nil
				}, consts.IngressWait, consts.WaitTick)

				resolver := createResolver(proxyUDPURL)
				t.Log("verifying that DNS queries are being load-balanced between multiple CoreDNS pods")
				assert.Eventually(t, func() bool { return isDNSResolverReturningExpectedResult(ctx, resolver, testdomain, "10.0.0.1") }, consts.IngressWait, consts.WaitTick)
				assert.Eventually(t, func() bool { return isDNSResolverReturningExpectedResult(ctx, resolver, testdomain, "10.0.0.2") }, consts.IngressWait, consts.WaitTick)
				assert.Eventually(t, func() bool { return isDNSResolverReturningExpectedResult(ctx, resolver, testdomain, "10.0.0.1") }, consts.IngressWait, consts.WaitTick)
				assert.Eventually(t, func() bool { return isDNSResolverReturningExpectedResult(ctx, resolver, testdomain, "10.0.0.2") }, consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Teardown(featureTeardown()).
		Feature()

	fPortMatching := features.
		New("port matching").
		WithLabel(testlabels.NetworkingFamily, testlabels.NetworkingFamilyGatewayAPI).
		WithLabel(testlabels.Kind, testlabels.KindUDPRoute).
		Setup(SkipIfRouterNotExpressions).
		WithSetup("deploy kong addon into cluster", featureSetup()).
		WithSetup("prepare Gateway and GatewayClass",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				t.Log("creating Gateway API client")
				gatewayClient, err := gatewayclient.NewForConfig(cfg.Client().RESTConfig())
				assert.NoError(t, err, "failed creating Gateway API client")
				ctx = SetInCtxForT(ctx, t, gatewayClient)
				namespace := GetNamespaceForT(ctx, t)

				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)

				gatewayClassName := uuid.NewString()
				t.Logf("deploying a supported GatewayClass %s to the test cluster", gatewayClassName)
				gwc, err := helpers.DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
				assert.NoError(t, err)
				cleaner.Add(gwc)
				ctx = SetInCtxForT(ctx, t, gwc)

				gatewayName := uuid.NewString()
				t.Logf("deploying a Gateway %s to the test cluster using unmanaged gateway mode and port %d", gatewayName, ktfkong.DefaultUDPServicePort)
				gateway, err := helpers.DeployGateway(ctx, gatewayClient, namespace, gatewayClassName, func(gw *gatewayapi.Gateway) {
					gw.Name = gatewayName
					gw.Spec.Listeners = builder.NewListener("udp").
						UDP().
						WithPort(ktfkong.DefaultUDPServicePort).
						IntoSlice()
				})
				assert.NoError(t, err)
				cleaner.Add(gateway)
				ctx = SetInCtxForT(ctx, t, gateway)

				return ctx
			}).
		WithSetup("prepare coredns deployment",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				namespace := GetNamespaceForT(ctx, t)
				cl := cfg.Client().Resources()
				cleaner := GetFromCtxForT[*clusters.Cleaner](ctx, t)
				gateway := GetFromCtxForT[*gatewayapi.Gateway](ctx, t)
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)

				t.Log("configuring coredns corefile")
				cfgmap1 := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "coredns",
						Namespace: namespace,
					},
					Data: map[string]string{"Corefile": corefile},
				}
				assert.NoError(t, cl.Create(ctx, cfgmap1))
				cleaner.Add(cfgmap1)

				t.Log("configuring a coredns deployent to deploy for UDP testing")
				container := generators.NewContainer("coredns", coreDNSImage, ktfkong.DefaultUDPServicePort)
				container.Ports[0].Protocol = corev1.ProtocolUDP
				container.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
				container.Args = []string{"-conf", "/etc/coredns/Corefile"}
				deployment := generators.NewDeploymentForContainer(container)

				t.Log("configuring the coredns pod with a custom corefile")
				deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes,
					corev1.Volume{
						Name: "config-volume",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: cfgmap1.Name,
							},
							Items: []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}},
						}},
					})
				deployment.Namespace = namespace

				t.Logf("deploying coredns deployment %q", deployment.Name)
				assert.NoError(t, cl.Create(ctx, deployment))
				cleaner.Add(deployment)

				t.Logf("exposing deployment %s/%s via service", deployment.Namespace, deployment.Name)
				service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
				service.Namespace = namespace
				service.Labels = map[string]string{"app": "coredns"}
				assert.NoError(t, cl.Create(ctx, service))
				cleaner.Add(service)

				t.Logf("creating a UDPRoute to access deployment %s via kong", deployment.Name)
				udproute := &gatewayapi.UDPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      uuid.NewString(),
						Namespace: namespace,
					},
					Spec: gatewayapi.UDPRouteSpec{
						CommonRouteSpec: gatewayapi.CommonRouteSpec{
							ParentRefs: []gatewayapi.ParentReference{{
								Name: gatewayapi.ObjectName(gateway.Name),
							}},
						},
						Rules: []gatewayapi.UDPRouteRule{{
							BackendRefs: builder.NewBackendRef(service.Name).WithPort(ktfkong.DefaultUDPServicePort).ToSlice(),
						}},
					},
				}
				udproute, err := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace).Create(ctx, udproute, metav1.CreateOptions{})
				assert.NoError(t, err)
				cleaner.Add(udproute)
				ctx = SetInCtxForT(ctx, t, udproute)

				return ctx
			}).
		Assess("using a port in UDPRoute not define in Gateway Listeners does not get the UDPRoute active",
			func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
				namespace := GetNamespaceForT(ctx, t)
				gatewayClient := GetFromCtxForT[*gatewayclient.Clientset](ctx, t)
				udproute := GetFromCtxForT[*gatewayapi.UDPRoute](ctx, t)
				udprouteClient := gatewayClient.GatewayV1alpha2().UDPRoutes(namespace)
				proxyUDPURL := GetUDPURLFromCtx(ctx)

				t.Log("updating UDPRoute parentRef to use a port not in the Gateway Listeners")
				assert.Eventually(t, func() bool {
					udproute, err := udprouteClient.Get(ctx, udproute.Name, metav1.GetOptions{})
					if err != nil {
						return false
					}
					notExistingPort := gatewayapi.PortNumber(81)
					udproute.Spec.ParentRefs[0].Port = &notExistingPort
					_, err = udprouteClient.Update(ctx, udproute, metav1.UpdateOptions{})
					return err == nil
				}, consts.IngressWait, consts.WaitTick)

				t.Log("verifying that the UDPRoute does not get active")
				assert.Eventually(t, not(urlResolvesSuccessfullyFn(ctx, proxyUDPURL)), consts.IngressWait, consts.WaitTick)

				return ctx
			}).
		Teardown(featureTeardown()).
		Feature()

	tenv.TestInParallel(t, fEssentials, fPortMatching)
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
