//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util/builder"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
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

	var gw *gatewayv1beta1.Gateway

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("gathering test data and generating a gateway kubernetes client")
	pubsvc, err := env.Cluster().Client().CoreV1().Services(consts.ControllerNamespace).Get(ctx, "ingress-controller-kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a new gateway")
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gateway)
	err = gatewayHealthCheck(ctx, gatewayClient, gateway.Name, ns.Name)
	require.NoError(t, err)

	t.Log("verifying that the gateway service ref gets provisioned when placeholder is used")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, defaultGatewayName, metav1.GetOptions{})
		require.NoError(t, err)
		return gw.Annotations[annotations.GatewayClassUnmanagedAnnotation] == "kong-system/ingress-controller-kong-proxy"
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway address is populated from the publish service")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		if len(gw.Spec.Addresses) == len(pubsvc.Status.LoadBalancer.Ingress) {
			addrs := make(map[string]bool, len(pubsvc.Status.LoadBalancer.Ingress))
			for _, ing := range pubsvc.Status.LoadBalancer.Ingress {
				// taking a slight shortcut by using the same map for both types. value lookups will still work
				// and the test isn't concerned with the weird case where you've somehow wound up with
				// LB Hostname 10.0.0.1 and GW IP 10.0.0.1. the GW type is also optional, so we don't always know
				addrs[ing.IP] = true
				addrs[ing.Hostname] = true
			}
			for _, addr := range gw.Spec.Addresses {
				if _, ok := addrs[addr.Value]; !ok {
					return false
				}
			}
			return true
		}
		return false
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway status gets updated to match the publish service")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return len(gw.Status.Listeners) == len(gw.Spec.Listeners) && len(gw.Status.Addresses) == len(gw.Spec.Addresses)
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway receives a final ready condition once reconciliation completes")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		ready := util.CheckCondition(
			gw.Status.Conditions,
			util.ConditionType(gatewayv1beta1.GatewayConditionReady),
			util.ConditionReason(gatewayv1beta1.GatewayReasonReady),
			metav1.ConditionTrue,
			gw.Generation,
		)
		return ready
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway listeners reach the ready condition")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, lstatus := range gw.Status.Listeners {
			if listenerReady := util.CheckCondition(
				lstatus.Conditions,
				util.ConditionType(gatewayv1beta1.ListenerConditionReady),
				util.ConditionReason(gatewayv1beta1.ListenerReasonReady),
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

	var gw *gatewayv1beta1.Gateway

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("generating a gateway kubernetes client and gathering test data")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a new Gateway using the default GatewayClass")
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, unmanagedGatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gateway)
	err = gatewayHealthCheck(ctx, gatewayClient, gateway.Name, ns.Name)
	require.NoError(t, err)

	t.Log("verifying that the gateway listeners reach the ready condition")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, defaultGatewayName, metav1.GetOptions{})
		require.NoError(t, err)
		for _, lstatus := range gw.Status.Listeners {
			ready := util.CheckCondition(
				lstatus.Conditions,
				util.ConditionType(gatewayv1beta1.GatewayConditionReady),
				util.ConditionReason(gatewayv1beta1.GatewayReasonReady),
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
		gatewayv1beta1.Listener{
			Name:     "badhttp",
			Protocol: gatewayv1beta1.HTTPProtocolType,
			Port:     gatewayv1beta1.PortNumber(80),
		},
		gatewayv1beta1.Listener{
			Name:     "badudp",
			Protocol: gatewayv1beta1.UDPProtocolType,
			Port:     gatewayv1beta1.PortNumber(80),
		},
	)

	_, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)
	gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)

	t.Log("confirming existing listen becomes unready and conflicted, new HTTP listen has hostname conflict, new UDP listen has proto conflict")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		var badhttpReady, badhttpConflicted, badudpReady, badudpConflicted, httpReady, httpConflicted bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "badudp" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1beta1.ListenerConditionConflicted) && condition.Status == metav1.ConditionTrue {
						badudpConflicted = (condition.Reason == string(gatewayv1beta1.ListenerReasonProtocolConflict))
					}
					if condition.Type == string(gatewayv1beta1.ListenerConditionReady) {
						badudpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "badhttp" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1beta1.ListenerConditionConflicted) && condition.Status == metav1.ConditionTrue {
						// this is a PROTOCOL conflict: although this only conflicts with the existing HTTP listen by
						// hostname, it also conflicts with the new UDP listener by protocol, and the latter takes
						// precedence
						badhttpConflicted = (condition.Reason == string(gatewayv1beta1.ListenerReasonProtocolConflict))
					}
					if condition.Type == string(gatewayv1beta1.ListenerConditionReady) {
						badhttpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "http" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1beta1.ListenerConditionConflicted) {
						httpConflicted = (condition.Status == metav1.ConditionTrue)
					}
					if condition.Type == string(gatewayv1beta1.ListenerConditionReady) {
						httpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
		}
		return !badhttpReady && badhttpConflicted && !badudpReady && badudpConflicted && !httpReady && httpConflicted
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("changing listeners to a set with conflicting hostnames")
	// these both use the empty hostname
	gw.Spec.Listeners = []gatewayv1beta1.Listener{
		{
			Name:     "httpsalpha",
			Protocol: gatewayv1beta1.HTTPSProtocolType,
			Port:     gatewayv1beta1.PortNumber(443),
		},
		{
			Name:     "httpsbravo",
			Protocol: gatewayv1beta1.HTTPSProtocolType,
			Port:     gatewayv1beta1.PortNumber(443),
		},
	}

	_, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("confirming listeners with conflicted hostnames receive appropriate conditions")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		var httpAlphaReady, httpAlphaConflicted, httpBravoReady, httpBravoConflicted bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "httpsalpha" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1beta1.ListenerConditionConflicted) && condition.Status == metav1.ConditionTrue {
						httpAlphaConflicted = (condition.Reason == string(gatewayv1beta1.ListenerReasonHostnameConflict))
					}
					if condition.Type == string(gatewayv1beta1.ListenerConditionReady) {
						httpAlphaReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "httpsbravo" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1beta1.ListenerConditionConflicted) && condition.Status == metav1.ConditionTrue {
						httpBravoConflicted = (condition.Reason == string(gatewayv1beta1.ListenerReasonHostnameConflict))
					}
					if condition.Type == string(gatewayv1beta1.ListenerConditionReady) {
						httpBravoReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
		}
		return !httpAlphaReady && httpAlphaConflicted && !httpBravoReady && httpBravoConflicted
	}, gatewayUpdateWaitTime, time.Second)
	t.Log("swapping out existing listeners with multiple compatible listeners")
	tlsHost := gatewayv1beta1.Hostname("tls.example")
	httpsHost := gatewayv1beta1.Hostname("https.example")
	httphostHost := gatewayv1beta1.Hostname("http.example")

	// this tests compatibility to the extent that we can with Kong listens. it does not support the full range
	// of compatible Gateway Routes. Gateway permits TLS and HTTPS routes to coexist on the same port so long
	// as all use unique hostnames. Kong, however, requires that TLS routes go through a TLS stream listen, so
	// the binds are separate and we cannot combine them. attempting to do so (e.g. setting the tls port to 443 here)
	// will result in ListenerReasonPortUnavailable
	gw.Spec.Listeners = []gatewayv1beta1.Listener{
		{
			Name:     "http",
			Protocol: gatewayv1beta1.HTTPProtocolType,
			Port:     gatewayv1beta1.PortNumber(80),
		},
		{
			Name:     "tls",
			Protocol: gatewayv1beta1.TLSProtocolType,
			Port:     gatewayv1beta1.PortNumber(8899),
			Hostname: &tlsHost,
		},
		{
			Name:     "https",
			Protocol: gatewayv1beta1.HTTPSProtocolType,
			Port:     gatewayv1beta1.PortNumber(443),
			Hostname: &httpsHost,
		},
		{
			Name:     "httphost",
			Protocol: gatewayv1beta1.HTTPProtocolType,
			Port:     gatewayv1beta1.PortNumber(80),
			Hostname: &httphostHost,
		},
	}

	_, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)
	gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)

	t.Log("confirming existing listen remains ready and new listens become ready")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		var httpReady, tlsReady, httpsReady, httphostReady bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "http" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1beta1.ListenerConditionReady) {
						httpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "tls" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1beta1.ListenerConditionReady) {
						tlsReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "https" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1beta1.ListenerConditionReady) {
						httpsReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "httphost" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1beta1.ListenerConditionReady) {
						httphostReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
		}
		return httpReady && tlsReady && httpsReady && httphostReady
	}, gatewayUpdateWaitTime, time.Second)
}

func TestUnmanagedGatewayControllerSupport(t *testing.T) {
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("generating a gateway kubernetes client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying an unsupported gatewayclass to the test cluster")
	unsupportedGatewayClass, err := DeployGatewayClass(ctx, gatewayClient, uuid.NewString(), func(gc *gatewayv1beta1.GatewayClass) {
		gc.Spec.ControllerName = unsupportedControllerName
	})
	require.NoError(t, err)
	cleaner.Add(unsupportedGatewayClass)

	t.Log("deploying a gateway using the unsupported gateway class")
	unsupportedGateway, err := DeployGateway(ctx, gatewayClient, ns.Name, unsupportedGatewayClass.Name)
	require.NoError(t, err)
	cleaner.Add(unsupportedGateway)

	t.Log("verifying that the unsupported Gateway object does not get scheduled by the controller")
	timeout := time.Now().Add(gatewayWaitTimeToVerifyScheduling)
	for timeout.After(time.Now()) {
		unsupportedGateway, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, unsupportedGateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		require.Len(t, unsupportedGateway.Status.Conditions, 1)
		//lint:ignore SA1019 the default upstream reason is still NotReconciled https://github.com/kubernetes-sigs/gateway-api/pull/1701
		//nolint:staticcheck
		require.Equal(t, string(gatewayv1beta1.GatewayReasonNotReconciled), unsupportedGateway.Status.Conditions[0].Reason)
	}
}

func TestUnmanagedGatewayClass(t *testing.T) {
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("generating a gateway kubernetes client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a gateway to the test cluster using unmanaged mode, but with no valid gatewayclass yet")
	gatewayClassName := uuid.NewString()
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1beta1.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("verifying that the Gateway object does not get scheduled by the controller due to missing its GatewayClass")
	timeout := time.Now().Add(gatewayWaitTimeToVerifyScheduling)
	for timeout.After(time.Now()) {
		gateway, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		require.Len(t, gateway.Status.Conditions, 1)
		//lint:ignore SA1019 the default upstream reason is still NotReconciled https://github.com/kubernetes-sigs/gateway-api/pull/1701
		//nolint:staticcheck
		require.Equal(t, string(gatewayv1beta1.GatewayReasonNotReconciled), gateway.Status.Conditions[0].Reason)
	}

	t.Log("deploying the missing gatewayclass to the test cluster")
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("now that the gatewayclass exists, verifying that the gateway resource gets resolved")
	require.Eventually(t, func() bool {
		gateway, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gateway.Status.Conditions {
			if cond.Reason == string(gatewayv1beta1.GatewayReasonReady) {
				return true
			}
		}
		return false
	}, gatewayUpdateWaitTime, time.Second)
}

func TestManagedGatewayClass(t *testing.T) {
	ctx := context.Background()

	ns, cleaner := helpers.Setup(ctx, t, env)

	t.Log("generating a gateway kubernetes client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a gateway to the test cluster, but with no valid gatewayclass yet")
	gatewayClassName := uuid.NewString()
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1beta1.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("verifying that the Gateway object does not get scheduled by the controller due to missing its GatewayClass")
	timeout := time.Now().Add(gatewayWaitTimeToVerifyScheduling)
	for timeout.After(time.Now()) {
		gateway, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		require.Len(t, gateway.Status.Conditions, 1)
		//lint:ignore SA1019 the default upstream reason is still NotReconciled https://github.com/kubernetes-sigs/gateway-api/pull/1701
		//nolint:staticcheck
		require.Equal(t, string(gatewayv1beta1.GatewayReasonNotReconciled), gateway.Status.Conditions[0].Reason)
	}

	t.Log("deploying a missing managed gatewayclass to the test cluster")
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName, func(gc *gatewayv1beta1.GatewayClass) {
		gc.Annotations = nil
	})
	require.NoError(t, err)
	cleaner.Add(gwc)

	finished := make(chan struct{})

	// Let's wait for one minute and check that the Gateway hasn't reconciled by the operator. It should never get ready.
	t.Log("the Gateway must not be reconciled as it is using a managed GatewayClass")
	time.AfterFunc(time.Minute, func() {
		defer close(finished)
		gateway, err = gatewayClient.GatewayV1beta1().Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gateway.Status.Conditions {
			if cond.Type == string(gatewayv1beta1.GatewayConditionReady) {
				require.Equal(t, cond.Status, metav1.ConditionFalse)
			}
		}
	})
	<-finished
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
	gateway, err := DeployGateway(ctx, gwClientSet, ns.Name, unmanagedGatewayClassName, func(gw *gatewayv1beta1.Gateway) {
		gw.Name = uuid.NewString()
		gw.Spec.Listeners = []gatewayv1beta1.Listener{
			builder.NewListener("http").HTTP().WithPort(80).
				WithAllowedRoutes(builder.NewAllowedRoutesFromAllNamespaces()).Build(),
			builder.NewListener("https").HTTPS().WithPort(443).
				WithAllowedRoutes(builder.NewAllowedRoutesFromAllNamespaces()).Build(),
		}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	client := env.Cluster().Client()

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
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
	HTTPRoute := func() *gatewayv1beta1.HTTPRoute {
		return &gatewayv1beta1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name: uuid.NewString(),
				Annotations: map[string]string{
					annotations.AnnotationPrefix + annotations.StripPathKey: "true",
				},
			},
			Spec: gatewayv1beta1.HTTPRouteSpec{
				CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
					ParentRefs: []gatewayv1beta1.ParentReference{{
						Name:      gatewayv1beta1.ObjectName(gateway.Name),
						Namespace: lo.ToPtr(gatewayv1beta1.Namespace(gateway.Namespace)),
					}},
				},
				Rules: []gatewayv1beta1.HTTPRouteRule{{
					Matches: []gatewayv1beta1.HTTPRouteMatch{
						builder.NewHTTPRouteMatch().WithPathPrefix("/test_gateway_filters").Build(),
					},
					BackendRefs: []gatewayv1beta1.HTTPBackendRef{
						builder.NewHTTPBackendRef(service.Name).WithPort(80).Build(),
					},
				}},
			},
		}
	}

	gatewayClient := gwClientSet.GatewayV1beta1()

	httpRoute, err := gatewayClient.HTTPRoutes(ns.Name).Create(ctx, HTTPRoute(), metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(httpRoute)

	otherRoute := HTTPRoute()
	otherRoute.Spec.Rules[0].Matches[0].Path.Value = kong.String("/other_test_gateway_filters")
	otherRoute, err = gatewayClient.HTTPRoutes(other.Name).Create(ctx, otherRoute, metav1.CreateOptions{})
	require.NoError(t, err)
	cleaner.Add(otherRoute)

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(t, gwClientSet, gatewayv1beta1.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("waiting for routes from HTTPRoute to become operational")
	helpers.EventuallyGETPath(t, proxyURL, "test_gateway_filters", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)
	t.Log("waiting for routes from HTTPRoute in other namespace to become operational")
	helpers.EventuallyGETPath(t, proxyURL, "other_test_gateway_filters", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)

	t.Log("changing to the same namespace filter")
	gateway, err = gatewayClient.Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
	require.NoError(t, err)
	fromSame := gatewayv1beta1.NamespacesFromSame
	gateway.Spec.Listeners = []gatewayv1beta1.Listener{
		{
			Name:     "http",
			Protocol: gatewayv1beta1.HTTPProtocolType,
			Port:     gatewayv1beta1.PortNumber(80),
			AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
				Namespaces: &gatewayv1beta1.RouteNamespaces{
					From: &fromSame,
				},
			},
		},
		{
			Name:     "https",
			Protocol: gatewayv1beta1.HTTPSProtocolType,
			Port:     gatewayv1beta1.PortNumber(443),
			AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
				Namespaces: &gatewayv1beta1.RouteNamespaces{
					From: &fromSame,
				},
			},
		},
	}
	_, err = gatewayClient.Gateways(ns.Name).Update(ctx, gateway, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("confirming other namespace route becomes inaccessible")
	helpers.EventuallyGETPath(t, proxyURL, "other_test_gateway_filters", http.StatusNotFound, "no Route matched", emptyHeaderSet, ingressWait, waitTick)
	t.Log("confirming same namespace route still operational")
	helpers.EventuallyGETPath(t, proxyURL, "test_gateway_filters", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)

	t.Log("changing to a selector filter")
	gateway, err = gatewayClient.Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
	require.NoError(t, err)
	fromSelector := gatewayv1beta1.NamespacesFromSelector
	gateway.Spec.Listeners = []gatewayv1beta1.Listener{
		{
			Name:     "http",
			Protocol: gatewayv1beta1.HTTPProtocolType,
			Port:     gatewayv1beta1.PortNumber(80),
			AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
				Namespaces: &gatewayv1beta1.RouteNamespaces{
					From: &fromSelector,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							clusters.TestResourceLabel: t.Name() + "other",
						},
					},
				},
			},
		},
		{
			Name:     "https",
			Protocol: gatewayv1beta1.HTTPSProtocolType,
			Port:     gatewayv1beta1.PortNumber(443),
			AllowedRoutes: &gatewayv1beta1.AllowedRoutes{
				Namespaces: &gatewayv1beta1.RouteNamespaces{
					From: &fromSelector,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							clusters.TestResourceLabel: t.Name() + "other",
						},
					},
				},
			},
		},
	}

	_, err = gatewayClient.Gateways(ns.Name).Update(ctx, gateway, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("confirming wrong selector namespace route becomes inaccessible")
	helpers.EventuallyGETPath(t, proxyURL, "test_gateway_filters", http.StatusNotFound, "no Route matched", emptyHeaderSet, ingressWait, waitTick)
	t.Log("confirming right selector namespace route becomes operational")
	helpers.EventuallyGETPath(t, proxyURL, "other_test_gateway_filters", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet, ingressWait, waitTick)
}
