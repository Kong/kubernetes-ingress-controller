//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	gatewayapi "github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v3/test"
	"github.com/kong/kubernetes-ingress-controller/v3/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v3/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

const (
	// gatewayWaitTimeToVerifyScheduling is the amount of time to wait during testing to verify
	// whether the sheduling of a new Gateway object has occurred.
	gatewayWaitTimeToVerifyScheduling = time.Second * 30

	// gatewayUpdateWaitTime is the amount of time to wait for updates to the Gateway, or to its
	// parent Service to fully resolve into ready state.
	gatewayUpdateWaitTime = time.Minute * 3
)

func TestUnmanagedGatewayBasics(t *testing.T) {
	ctx := context.Background()

	var gw *gatewayapi.Gateway

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("gathering test data and generating a gateway kubernetes client")
	pubsvc, err := env.Cluster().Client().CoreV1().Services(consts.ControllerNamespace).Get(ctx, "ingress-controller-kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	pubsvcUDP, err := env.Cluster().Client().CoreV1().Services(consts.ControllerNamespace).Get(ctx, "ingress-controller-kong-udp-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a new gateway")
	gateway, err := helpers.DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gateway)
	err = gatewayHealthCheck(ctx, gatewayClient, gateway.Name, ns.Name)
	require.NoError(t, err)

	t.Log("verifying that the gateway service ref gets provisioned when placeholder is used")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, defaultGatewayName, metav1.GetOptions{})
		require.NoError(t, err)
		return lo.Contains(
			annotations.ExtractGatewayPublishService(gw.Annotations),
			fmt.Sprintf("%s/%s", pubsvc.Namespace, pubsvc.Name),
		) && lo.Contains(
			annotations.ExtractGatewayPublishService(gw.Annotations),
			fmt.Sprintf("%s/%s", pubsvcUDP.Namespace, pubsvcUDP.Name),
		)
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway address is populated from the publish service")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		addrs := make(map[string]bool, len(pubsvc.Status.LoadBalancer.Ingress))
		for _, ing := range pubsvc.Status.LoadBalancer.Ingress {
			// taking a slight shortcut by using the same map for both types. value lookups will still work
			// and the test isn't concerned with the weird case where you've somehow wound up with
			// LB Hostname 10.0.0.1 and GW IP 10.0.0.1. the GW type is also optional, so we don't always know
			addrs[ing.IP] = true
			addrs[ing.Hostname] = true
		}
		for _, ing := range pubsvcUDP.Status.LoadBalancer.Ingress {
			addrs[ing.IP] = true
			addrs[ing.Hostname] = true
		}
		for _, addr := range gw.Status.Addresses {
			if _, ok := addrs[addr.Value]; !ok {
				return false
			}
		}
		return true
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway listeners status match the spec listeners")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return len(gw.Status.Listeners) == len(gw.Spec.Listeners)
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway receives a final programmed condition once reconciliation completes")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		ready := util.CheckCondition(
			gw.Status.Conditions,
			util.ConditionType(gatewayapi.GatewayConditionProgrammed),
			util.ConditionReason(gatewayapi.GatewayReasonProgrammed),
			metav1.ConditionTrue,
			gw.Generation,
		)
		return ready
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway listeners reach the programmed condition")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, lstatus := range gw.Status.Listeners {
			if listenerReady := util.CheckCondition(
				lstatus.Conditions,
				util.ConditionType(gatewayapi.ListenerConditionProgrammed),
				util.ConditionReason(gatewayapi.ListenerReasonProgrammed),
				metav1.ConditionTrue,
				gw.Generation,
			); !listenerReady {
				return false
			}
		}
		return true
	}, gatewayUpdateWaitTime, time.Second)
}

func TestGatewayListenerConflicts(t *testing.T) {
	ctx := context.Background()

	var gw *gatewayapi.Gateway

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("generating a gateway kubernetes client and gathering test data")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("adding a test certificate")
	cert, key := certificate.MustGenerateSelfSignedCertPEMFormat(certificate.WithCommonName(ns.Name + ".example.com"))
	certName := "cert"
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      certName,
			Namespace: ns.Name,
		},
		Data: map[string][]byte{
			"tls.crt": cert,
			"tls.key": key,
		},
	}

	_, err = env.Cluster().Client().CoreV1().Secrets(ns.Name).Create(ctx, secret, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("deploying a new Gateway using the default GatewayClass")
	gateway, err := helpers.DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gateway)
	err = gatewayHealthCheck(ctx, gatewayClient, gateway.Name, ns.Name)
	require.NoError(t, err)

	t.Log("verifying that the gateway listeners reach the programmed condition")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, defaultGatewayName, metav1.GetOptions{})
		require.NoError(t, err)
		for _, lstatus := range gw.Status.Listeners {
			ready := util.CheckCondition(
				lstatus.Conditions,
				util.ConditionType(gatewayapi.GatewayConditionProgrammed),
				util.ConditionReason(gatewayapi.GatewayReasonProgrammed),
				metav1.ConditionTrue,
				gw.Generation,
			)
			if !ready {
				return false
			}
		}
		return true
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("adding conflicting listeners")
	gw.Spec.Listeners = append(gw.Spec.Listeners,
		gatewayapi.Listener{
			Name:     "badudp",
			Protocol: gatewayapi.UDPProtocolType,
			Port:     gatewayapi.PortNumber(80),
		},
	)

	_, err = gatewayClient.GatewayV1().Gateways(ns.Name).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)
	gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)

	t.Log("confirming existing listen becomes unready and conflicted, new UDP listen has proto conflict")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		var badudpReady, badudpConflicted, httpReady, httpConflicted bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "badudp" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayapi.ListenerConditionConflicted) && condition.Status == metav1.ConditionTrue {
						badudpConflicted = (condition.Reason == string(gatewayapi.ListenerReasonProtocolConflict))
					}
					if condition.Type == string(gatewayapi.ListenerConditionProgrammed) {
						badudpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "http" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayapi.ListenerConditionConflicted) {
						httpConflicted = (condition.Status == metav1.ConditionTrue)
					}
					if condition.Type == string(gatewayapi.ListenerConditionProgrammed) {
						httpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
		}
		return !badudpReady && badudpConflicted && !httpReady && httpConflicted
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("swapping out existing listeners with multiple compatible listeners")
	tlsHost := gatewayapi.Hostname("tls.example")
	httpsHost := gatewayapi.Hostname("https.example")
	httphostHost := gatewayapi.Hostname("http.example")

	// This tests compatibility to the extent that we can with Kong listens. it does not support the full range
	// of compatible Gateway Routes. Gateway permits TLS and HTTPS routes to coexist on the same port so long
	// as all use unique hostnames. Kong, however, requires that TLS routes go through a TLS stream listen, so
	// the binds are separate and we cannot combine them. attempting to do so (e.g. setting the tls port to 443 here)
	// will result in ListenerReasonPortUnavailable.
	gw.Spec.Listeners = []gatewayapi.Listener{
		{
			Name:     "http",
			Protocol: gatewayapi.HTTPProtocolType,
			Port:     gatewayapi.PortNumber(ktfkong.DefaultProxyHTTPPort),
		},
		{
			Name:     "tls",
			Protocol: gatewayapi.TLSProtocolType,
			Port:     gatewayapi.PortNumber(ktfkong.DefaultTLSServicePort),
			Hostname: &tlsHost,
			TLS: &gatewayapi.GatewayTLSConfig{
				CertificateRefs: []gatewayapi.SecretObjectReference{
					{
						Name: gatewayapi.ObjectName(certName),
					},
				},
			},
		},
		{
			Name:     "https",
			Protocol: gatewayapi.HTTPSProtocolType,
			Port:     gatewayapi.PortNumber(ktfkong.DefaultProxyTLSServicePort),
			Hostname: &httpsHost,
			TLS: &gatewayapi.GatewayTLSConfig{
				CertificateRefs: []gatewayapi.SecretObjectReference{
					{
						Name: gatewayapi.ObjectName(certName),
					},
				},
			},
		},
		{
			Name:     "httphost",
			Protocol: gatewayapi.HTTPProtocolType,
			Port:     gatewayapi.PortNumber(ktfkong.DefaultProxyHTTPPort),
			Hostname: &httphostHost,
		},
	}

	_, err = gatewayClient.GatewayV1().Gateways(ns.Name).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)
	gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)

	t.Log("confirming existing listen remains ready and new listens become ready")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		var httpReady, tlsReady, httpsReady, httphostReady bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "http" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayapi.ListenerConditionProgrammed) {
						httpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "tls" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayapi.ListenerConditionProgrammed) {
						tlsReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "https" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayapi.ListenerConditionProgrammed) {
						httpsReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "httphost" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayapi.ListenerConditionProgrammed) {
						httphostReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
		}
		return httpReady && tlsReady && httpsReady && httphostReady
	}, gatewayUpdateWaitTime, time.Second)
}

func TestGatewayFilters(t *testing.T) {
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)

	other, err := clusters.GenerateNamespace(ctx, env.Cluster(), t.Name()+"other")
	require.NoError(t, err)
	defer func(t *testing.T) {
		assert.NoError(t, clusters.CleanupGeneratedResources(ctx, env.Cluster(), t.Name()+"other"))
	}(t)

	t.Log("deploying a supported gatewayclass to the test cluster")
	gwClientSet, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a gateway that allows routes in all namespaces")
	gateway, err := helpers.DeployGateway(ctx, gwClientSet, ns.Name, unmanagedGatewayClassName, func(gw *gatewayapi.Gateway) {
		gw.Name = uuid.NewString()
		gw.Spec.Listeners = []gatewayapi.Listener{
			builder.NewListener("http").HTTP().WithPort(80).
				WithAllowedRoutes(builder.NewAllowedRoutesFromAllNamespaces()).Build(),
		}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	client := env.Cluster().Client()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, test.HTTPBinPort)
	deploymentTemplate := generators.NewDeploymentForContainer(container)
	deployment, err := client.AppsV1().Deployments(ns.Name).Create(ctx, deploymentTemplate, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(deployment)
	otherDeployment, err := client.AppsV1().Deployments(other.Name).Create(ctx, deploymentTemplate, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(otherDeployment)

	t.Logf("exposing deployment %s/%s via service", ns.Name, deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = client.CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(service)

	t.Logf("exposing deployment %s/%s via service", other.Name, deployment.Name)
	otherService := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	otherService, err = client.CoreV1().Services(other.Name).Create(ctx, otherService, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(otherService)

	t.Logf("creating an httproute to access deployment %s via kong", deployment.Name)
	HTTPRoute := func() *gatewayapi.HTTPRoute {
		return &gatewayapi.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name: uuid.NewString(),
				Annotations: map[string]string{
					annotations.AnnotationPrefix + annotations.StripPathKey: "true",
				},
			},
			Spec: gatewayapi.HTTPRouteSpec{
				CommonRouteSpec: gatewayapi.CommonRouteSpec{
					ParentRefs: []gatewayapi.ParentReference{{
						Name:      gatewayapi.ObjectName(gateway.Name),
						Namespace: lo.ToPtr(gatewayapi.Namespace(gateway.Namespace)),
					}},
				},
				Rules: []gatewayapi.HTTPRouteRule{{
					Matches: []gatewayapi.HTTPRouteMatch{
						builder.NewHTTPRouteMatch().WithPathPrefix("/test_gateway_filters").Build(),
					},
					BackendRefs: []gatewayapi.HTTPBackendRef{
						builder.NewHTTPBackendRef(service.Name).WithPort(80).Build(),
					},
				}},
			},
		}
	}

	gatewayClient := gwClientSet.GatewayV1()

	httpRoute, err := gatewayClient.HTTPRoutes(ns.Name).Create(ctx, HTTPRoute(), metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(httpRoute)

	otherRoute := HTTPRoute()
	otherRoute.Spec.Rules[0].Matches[0].Path.Value = kong.String("/other_test_gateway_filters")
	otherRoute, err = gatewayClient.HTTPRoutes(other.Name).Create(ctx, otherRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(otherRoute)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := helpers.GetGatewayIsLinkedCallback(ctx, t, gwClientSet, gatewayapi.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("waiting for routes from HTTPRoute to become operational")
	helpers.EventuallyGETPath(t, proxyHTTPURL, proxyHTTPURL.Host, "test_gateway_filters", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)
	t.Log("waiting for routes from HTTPRoute in other namespace to become operational")
	helpers.EventuallyGETPath(t, proxyHTTPURL, proxyHTTPURL.Host, "other_test_gateway_filters", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)

	t.Log("changing to the same namespace filter")
	require.Eventually(t, func() bool {
		gateway, err = gatewayClient.Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting gateway %s: %v", gateway.Name, err)
			return false
		}

		gateway.Spec.Listeners = []gatewayapi.Listener{
			builder.NewListener("http").HTTP().WithPort(80).
				WithAllowedRoutes(builder.NewAllowedRoutesFromSameNamespaces()).Build(),
		}
		_, err = gatewayClient.Gateways(ns.Name).Update(ctx, gateway, metav1.UpdateOptions{})
		if err != nil {
			t.Logf("error updating gateway %s: %v", gateway.Name, err)
			return false
		}
		return true
	}, ingressWait, waitTick)

	t.Log("confirming other namespace route becomes inaccessible")
	helpers.EventuallyGETPath(t, proxyHTTPURL, proxyHTTPURL.Host, "other_test_gateway_filters", http.StatusNotFound, "no Route matched", emptyHeaderSet, ingressWait, waitTick)
	t.Log("confirming same namespace route still operational")
	helpers.EventuallyGETPath(t, proxyHTTPURL, proxyHTTPURL.Host, "test_gateway_filters", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)

	t.Log("changing to a selector filter")
	require.Eventually(t, func() bool {
		gateway, err = gatewayClient.Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting gateway %s: %v", gateway.Name, err)
			return false
		}

		fromSelector := builder.NewAllowedRoutesFromSelectorNamespace(
			&metav1.LabelSelector{
				MatchLabels: map[string]string{
					clusters.TestResourceLabel: t.Name() + "other",
				},
			},
		)
		gateway.Spec.Listeners = []gatewayapi.Listener{
			builder.NewListener("http").HTTP().WithPort(80).WithAllowedRoutes(fromSelector).Build(),
		}
		_, err = gatewayClient.Gateways(ns.Name).Update(ctx, gateway, metav1.UpdateOptions{})
		if err != nil {
			t.Logf("error updating gateway %s: %v", gateway.Name, err)
			return false
		}
		return true
	}, ingressWait, waitTick)

	t.Log("confirming wrong selector namespace route becomes inaccessible")
	helpers.EventuallyGETPath(t, proxyHTTPURL, proxyHTTPURL.Host, "test_gateway_filters", http.StatusNotFound, "no Route matched", emptyHeaderSet, ingressWait, waitTick)
	t.Log("confirming right selector namespace route becomes operational")
	helpers.EventuallyGETPath(t, proxyHTTPURL, proxyHTTPURL.Host, "other_test_gateway_filters", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)
}
