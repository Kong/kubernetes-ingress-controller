//go:build integration_tests
// +build integration_tests

package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
)

const (
	// gatewayWaitTimeToVerifyScheduling is the amount of time to wait during testing to verify
	// whether the sheduling of a new Gateway object has occurred.
	gatewayWaitTimeToVerifyScheduling = time.Second * 30

	// gatewayUpdateWaitTime is the amount of time to wait for updates to the Gateway, or to its
	// parent Service to fully resolve into ready state.
	gatewayUpdateWaitTime = time.Minute * 3

	unmanagedAnnotation = annotations.AnnotationPrefix + annotations.GatewayUnmanagedAnnotation
)

func TestUnmanagedGatewayBasics(t *testing.T) {
	var gw *gatewayv1alpha2.Gateway

	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	t.Log("gathering test data and generating a gateway kubernetes client")
	pubsvc, err := env.Cluster().Client().CoreV1().Services(controllerNamespace).Get(ctx, "ingress-controller-kong-proxy", metav1.GetOptions{})
	require.NoError(t, err)
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a new gateway")
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, managedGatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gateway)
	err = gatewayHealthCheck(ctx, gatewayClient, gateway.Name, ns.Name)
	require.NoError(t, err)

	t.Log("verifying that the gateway service ref gets provisioned when placeholder is used")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, defaultGatewayName, metav1.GetOptions{})
		require.NoError(t, err)
		return gw.Annotations[unmanagedAnnotation] == "kong-system/ingress-controller-kong-proxy"
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway address is populated from the publish service")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
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
		gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		return len(gw.Status.Listeners) == len(gw.Spec.Listeners) && len(gw.Status.Addresses) == len(gw.Spec.Addresses)
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway receives a final ready condition once reconciliation completes")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gw.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayReasonReady) {
				return true
			}
		}
		return false
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("verifying that the gateway listeners reach the ready condition")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, lstatus := range gw.Status.Listeners {
			// we may have several conditions but only care about one, so loop through each and mark ready only if
			// we find the correct one with the correct status
			ready := false
			for _, condition := range lstatus.Conditions {
				if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) && condition.Status == metav1.ConditionTrue {
					ready = true
				}
			}
			if !ready {
				return false
			}
		}
		return true
	}, gatewayUpdateWaitTime, time.Second)
}

func TestGatewayListenerConflicts(t *testing.T) {
	var gw *gatewayv1alpha2.Gateway

	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	t.Log("generating a gateway kubernetes client and gathering test data")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a new gateway using the default gateway")
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, managedGatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gateway)
	err = gatewayHealthCheck(ctx, gatewayClient, gateway.Name, ns.Name)
	require.NoError(t, err)

	t.Log("verifying that the gateway listeners reach the ready condition")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, defaultGatewayName, metav1.GetOptions{})
		require.NoError(t, err)
		for _, lstatus := range gw.Status.Listeners {
			// we may have several conditions but only care about one, so loop through each and mark ready only if
			// we find the correct one with the correct status
			ready := false
			for _, condition := range lstatus.Conditions {
				if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) && condition.Status == metav1.ConditionTrue {
					ready = true
				}
			}
			if !ready {
				return false
			}
		}
		return true
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("adding conflicting listeners")
	gw.Spec.Listeners = append(gw.Spec.Listeners,
		gatewayv1alpha2.Listener{
			Name:     "badhttp",
			Protocol: gatewayv1alpha2.HTTPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(80),
		},
		gatewayv1alpha2.Listener{
			Name:     "badudp",
			Protocol: gatewayv1alpha2.UDPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(80),
		},
	)

	_, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)
	gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)

	t.Log("confirming existing listen becomes unready and conflicted, new HTTP listen has hostname conflict, new UDP listen has proto conflict")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		var badhttpReady, badhttpConflicted, badudpReady, badudpConflicted, httpReady, httpConflicted bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "badudp" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionConflicted) && condition.Status == metav1.ConditionTrue {
						badudpConflicted = (condition.Reason == string(gatewayv1alpha2.ListenerReasonProtocolConflict))
					}
					if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) {
						badudpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "badhttp" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionConflicted) && condition.Status == metav1.ConditionTrue {
						// this is a PROTOCOL conflict: although this only conflicts with the existing HTTP listen by
						// hostname, it also conflicts with the new UDP listener by protocol, and the latter takes
						// precedence
						badhttpConflicted = (condition.Reason == string(gatewayv1alpha2.ListenerReasonProtocolConflict))
					}
					if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) {
						badhttpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "http" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionConflicted) {
						httpConflicted = (condition.Status == metav1.ConditionTrue)
					}
					if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) {
						httpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
		}
		return !badhttpReady && badhttpConflicted && !badudpReady && badudpConflicted && !httpReady && httpConflicted
	}, gatewayUpdateWaitTime, time.Second)

	t.Log("changing listeners to a set with conflicting hostnames")
	// these both use the empty hostname
	gw.Spec.Listeners = []gatewayv1alpha2.Listener{
		{
			Name:     "httpsalpha",
			Protocol: gatewayv1alpha2.HTTPSProtocolType,
			Port:     gatewayv1alpha2.PortNumber(443),
		},
		{
			Name:     "httpsbravo",
			Protocol: gatewayv1alpha2.HTTPSProtocolType,
			Port:     gatewayv1alpha2.PortNumber(443),
		},
	}

	_, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("confirming listeners with conflicted hostnames receive appropriate conditions")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		var httpAlphaReady, httpAlphaConflicted, httpBravoReady, httpBravoConflicted bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "httpsalpha" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionConflicted) && condition.Status == metav1.ConditionTrue {
						httpAlphaConflicted = (condition.Reason == string(gatewayv1alpha2.ListenerReasonHostnameConflict))
					}
					if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) {
						httpAlphaReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "httpsbravo" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionConflicted) && condition.Status == metav1.ConditionTrue {
						httpBravoConflicted = (condition.Reason == string(gatewayv1alpha2.ListenerReasonHostnameConflict))
					}
					if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) {
						httpBravoReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
		}
		return !httpAlphaReady && httpAlphaConflicted && !httpBravoReady && httpBravoConflicted
	}, gatewayUpdateWaitTime, time.Second)
	t.Log("swapping out existing listeners with multiple compatible listeners")
	tlsHost := gatewayv1alpha2.Hostname("tls.example")
	httpsHost := gatewayv1alpha2.Hostname("https.example")
	httphostHost := gatewayv1alpha2.Hostname("http.example")

	// this tests compatibility to the extent that we can with Kong listens. it does not support the full range
	// of compatible Gateway Routes. Gateway permits TLS and HTTPS routes to coexist on the same port so long
	// as all use unique hostnames. Kong, however, requires that TLS routes go through a TLS stream listen, so
	// the binds are separate and we cannot combine them. attempting to do so (e.g. setting the tls port to 443 here)
	// will result in ListenerReasonPortUnavailable
	gw.Spec.Listeners = []gatewayv1alpha2.Listener{
		{
			Name:     "http",
			Protocol: gatewayv1alpha2.HTTPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(80),
		},
		{
			Name:     "tls",
			Protocol: gatewayv1alpha2.TLSProtocolType,
			Port:     gatewayv1alpha2.PortNumber(8899),
			Hostname: &tlsHost,
		},
		{
			Name:     "https",
			Protocol: gatewayv1alpha2.HTTPSProtocolType,
			Port:     gatewayv1alpha2.PortNumber(443),
			Hostname: &httpsHost,
		},
		{
			Name:     "httphost",
			Protocol: gatewayv1alpha2.HTTPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(80),
			Hostname: &httphostHost,
		},
	}

	_, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Update(ctx, gw, metav1.UpdateOptions{})
	require.NoError(t, err)
	gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
	require.NoError(t, err)

	t.Log("confirming existing listen remains ready and new listens become ready")
	require.Eventually(t, func() bool {
		gw, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gw.Name, metav1.GetOptions{})
		require.NoError(t, err)
		var httpReady, tlsReady, httpsReady, httphostReady bool
		for _, lstatus := range gw.Status.Listeners {
			if lstatus.Name == "http" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) {
						httpReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "tls" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) {
						tlsReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "https" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) {
						httpsReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
			if lstatus.Name == "httphost" {
				for _, condition := range lstatus.Conditions {
					if condition.Type == string(gatewayv1alpha2.ListenerConditionReady) {
						httphostReady = (condition.Status == metav1.ConditionTrue)
					}
				}
			}
		}
		return httpReady && tlsReady && httpsReady && httphostReady
	}, gatewayUpdateWaitTime, time.Second)
}

func TestUnmanagedGatewayControllerSupport(t *testing.T) {
	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	t.Log("generating a gateway kubernetes client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying an unsupported gatewayclass to the test cluster")
	unsupportedGatewayClass, err := DeployGatewayClass(ctx, gatewayClient, uuid.NewString(), func(gc *gatewayv1alpha2.GatewayClass) {
		gc.Spec.ControllerName = unmanagedControllerName
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
		unsupportedGateway, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, unsupportedGateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		require.Len(t, unsupportedGateway.Status.Conditions, 1)
		require.Equal(t, string(gatewayv1alpha2.GatewayReasonNotReconciled), unsupportedGateway.Status.Conditions[0].Reason)
	}
}

func TestUnmanagedGatewayClass(t *testing.T) {
	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	t.Log("generating a gateway kubernetes client")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a gateway to the test cluster using unmanaged mode, but with no valid gatewayclass yet")
	gatewayClassName := uuid.NewString()
	gatewayName := uuid.NewString()
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, gatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("verifying that the Gateway object does not get scheduled by the controller due to missing its GatewayClass")
	timeout := time.Now().Add(gatewayWaitTimeToVerifyScheduling)
	for timeout.After(time.Now()) {
		gateway, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		require.Len(t, gateway.Status.Conditions, 1)
		require.Equal(t, string(gatewayv1alpha2.GatewayReasonNotReconciled), gateway.Status.Conditions[0].Reason)
	}

	t.Log("deploying the missing gatewayclass to the test cluster")
	gwc, err := DeployGatewayClass(ctx, gatewayClient, gatewayClassName)
	require.NoError(t, err)
	cleaner.Add(gwc)

	t.Log("now that the gatewayclass exists, verifying that the gateway resource gets resolved")
	require.Eventually(t, func() bool {
		gateway, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
		require.NoError(t, err)
		for _, cond := range gateway.Status.Conditions {
			if cond.Reason == string(gatewayv1alpha2.GatewayReasonReady) {
				return true
			}
		}
		return false
	}, gatewayUpdateWaitTime, time.Second)
}

func TestGatewayFilters(t *testing.T) {
	ns, cleaner := setup(t)
	defer func() { assert.NoError(t, cleaner.Cleanup(ctx)) }()

	other, err := clusters.GenerateNamespace(ctx, env.Cluster(), t.Name()+"other")
	require.NoError(t, err)
	defer func(t *testing.T) {
		assert.NoError(t, clusters.CleanupGeneratedResources(ctx, env.Cluster(), t.Name()+"other"))
	}(t)

	t.Log("deploying a supported gatewayclass to the test cluster")
	gatewayClient, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)

	t.Log("deploying a gateway that allows routes in all namespaces")
	gatewayName := uuid.NewString()
	fromAll := gatewayv1alpha2.NamespacesFromAll
	gateway, err := DeployGateway(ctx, gatewayClient, ns.Name, managedGatewayClassName, func(gw *gatewayv1alpha2.Gateway) {
		gw.Name = gatewayName
		gw.Spec.Listeners = []gatewayv1alpha2.Listener{
			{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(80),
				AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
					Namespaces: &gatewayv1alpha2.RouteNamespaces{
						From: &fromAll,
					},
				},
			},
			{
				Name:     "https",
				Protocol: gatewayv1alpha2.HTTPSProtocolType,
				Port:     gatewayv1alpha2.PortNumber(443),
				AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
					Namespaces: &gatewayv1alpha2.RouteNamespaces{
						From: &fromAll,
					},
				},
			},
		}
	})
	require.NoError(t, err)
	cleaner.Add(gateway)

	t.Log("deploying a minimal HTTP container deployment to test Ingress routes")
	container := generators.NewContainer("httpbin", test.HTTPBinImage, 80)
	deploymentTemplate := generators.NewDeploymentForContainer(container)
	deployment, err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deploymentTemplate, metav1.CreateOptions{})
	require.NoError(t, err)
	otherDeployment, err := env.Cluster().Client().AppsV1().Deployments(other.Name).Create(ctx, deploymentTemplate, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up deployments %s/%s and %s/%s", ns.Name, deployment.Name, other.Name, otherDeployment.Name)
		if err := env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
		if err := env.Cluster().Client().AppsV1().Deployments(other.Name).Delete(ctx, otherDeployment.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)
	_, err = env.Cluster().Client().CoreV1().Services(other.Name).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the services %s/%s and %s/%s", ns.Name, service.Name, other.Name, service.Name)
		if err := env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
		if err := env.Cluster().Client().CoreV1().Services(other.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Logf("creating an httproute to access deployment %s via kong", deployment.Name)
	httpPort := gatewayv1alpha2.PortNumber(80)
	pathMatchPrefix := gatewayv1alpha2.PathMatchPathPrefix
	refNamespace := gatewayv1alpha2.Namespace(gateway.Namespace)
	httprouteTemplate := &gatewayv1alpha2.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.AnnotationPrefix + annotations.StripPathKey: "true",
			},
		},
		Spec: gatewayv1alpha2.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name:      gatewayv1alpha2.ObjectName(gateway.Name),
					Namespace: &refNamespace,
				}},
			},
			Rules: []gatewayv1alpha2.HTTPRouteRule{{
				Matches: []gatewayv1alpha2.HTTPRouteMatch{
					{
						Path: &gatewayv1alpha2.HTTPPathMatch{
							Type:  &pathMatchPrefix,
							Value: kong.String("/httpbin"),
						},
					},
				},
				BackendRefs: []gatewayv1alpha2.HTTPBackendRef{{
					BackendRef: gatewayv1alpha2.BackendRef{
						BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
							Name: gatewayv1alpha2.ObjectName(service.Name),
							Port: &httpPort,
						},
					},
				}},
			}},
		},
	}
	httpRoute, err := gatewayClient.GatewayV1alpha2().HTTPRoutes(ns.Name).Create(ctx, httprouteTemplate, metav1.CreateOptions{})
	require.NoError(t, err)

	otherRoute, err := gatewayClient.GatewayV1alpha2().HTTPRoutes(other.Name).Create(ctx, httprouteTemplate, metav1.CreateOptions{})
	require.NoError(t, err)
	otherRoute.Spec.Rules[0].Matches[0].Path.Value = kong.String("/otherbin")
	_, err = gatewayClient.GatewayV1alpha2().HTTPRoutes(other.Name).Update(ctx, otherRoute, metav1.UpdateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the httproute %s", httpRoute.Name)
		if err := gatewayClient.GatewayV1alpha2().HTTPRoutes(ns.Name).Delete(ctx, httpRoute.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
		if err := gatewayClient.GatewayV1alpha2().HTTPRoutes(other.Name).Delete(ctx, httpRoute.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("verifying that the Gateway gets linked to the route via status")
	callback := GetGatewayIsLinkedCallback(t, gatewayClient, gatewayv1alpha2.HTTPProtocolType, ns.Name, httpRoute.Name)
	require.Eventually(t, callback, ingressWait, waitTick)

	t.Log("waiting for routes from HTTPRoute to become operational")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
	t.Log("waiting for routes from HTTPRoute in other namespace to become operational")
	eventuallyGETPath(t, "otherbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("changing to the same namespace filter")
	gateway, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
	require.NoError(t, err)
	fromSame := gatewayv1alpha2.NamespacesFromSame
	gateway.Spec.Listeners = []gatewayv1alpha2.Listener{
		{
			Name:     "http",
			Protocol: gatewayv1alpha2.HTTPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(80),
			AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
				Namespaces: &gatewayv1alpha2.RouteNamespaces{
					From: &fromSame,
				},
			},
		},
		{
			Name:     "https",
			Protocol: gatewayv1alpha2.HTTPSProtocolType,
			Port:     gatewayv1alpha2.PortNumber(443),
			AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
				Namespaces: &gatewayv1alpha2.RouteNamespaces{
					From: &fromSame,
				},
			},
		},
	}
	_, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Update(ctx, gateway, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("confirming other namespace route becomes inaccessible")
	eventuallyGETPath(t, "otherbin", http.StatusNotFound, "no Route matched", emptyHeaderSet)
	t.Log("confirming same namespace route still operational")
	eventuallyGETPath(t, "httpbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)

	t.Log("changing to a selector filter")
	gateway, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Get(ctx, gateway.Name, metav1.GetOptions{})
	require.NoError(t, err)
	fromSelector := gatewayv1alpha2.NamespacesFromSelector
	gateway.Spec.Listeners = []gatewayv1alpha2.Listener{
		{
			Name:     "http",
			Protocol: gatewayv1alpha2.HTTPProtocolType,
			Port:     gatewayv1alpha2.PortNumber(80),
			AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
				Namespaces: &gatewayv1alpha2.RouteNamespaces{
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
			Protocol: gatewayv1alpha2.HTTPSProtocolType,
			Port:     gatewayv1alpha2.PortNumber(443),
			AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
				Namespaces: &gatewayv1alpha2.RouteNamespaces{
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

	_, err = gatewayClient.GatewayV1alpha2().Gateways(ns.Name).Update(ctx, gateway, metav1.UpdateOptions{})
	require.NoError(t, err)

	t.Log("confirming wrong selector namespace route becomes inaccessible")
	eventuallyGETPath(t, "httpbin", http.StatusNotFound, "no Route matched", emptyHeaderSet)
	t.Log("confirming right selector namespace route becomes operational")
	eventuallyGETPath(t, "otherbin", http.StatusOK, "<title>httpbin.org</title>", emptyHeaderSet)
}
