//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

const (
	// defaultGatewayName is the default name for the Gateways created during tests.
	defaultGatewayName = "kong"
	// unmanagedGatewayClassName is the name of the default GatewayClass created during the test environment setup.
	unmanagedGatewayClassName = "kong-unmanaged"
	// unsupportedControllerName is the name of the controller used for those gateways that are not supported
	// by an actual controller (i.e., they won't be scheduled).
	unsupportedControllerName gatewayv1.GatewayController = "example.com/unsupported-gateway-controller"
	// kongRouterFlavorExpressions is the value used in router_flavor of kong configuration
	// to enable expression based router of kong.
	kongRouterFlavorExpressions = "expressions"
)

// DeployGateway creates a default gatewayClass, accepts a variadic set of options,
// and finally deploys it on the Kubernetes cluster by means of the gateway client given as arg.
func DeployGatewayClass(ctx context.Context, client *gatewayclient.Clientset, gatewayClassName string, opts ...func(*gatewayv1.GatewayClass)) (*gatewayv1.GatewayClass, error) {
	gwc := &gatewayv1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: gatewayClassName,
			Annotations: map[string]string{
				annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
			},
		},
		Spec: gatewayv1.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
	}

	for _, opt := range opts {
		opt(gwc)
	}

	result, err := client.GatewayV1beta1().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		err = client.GatewayV1beta1().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{})
		if err != nil {
			return result, err
		}
		result, err = client.GatewayV1beta1().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
	}
	return result, err
}

// DeployGateway creates a default gateway, accepts a variadic set of options,
// and finally deploys it on the Kubernetes cluster by means of the gateway client given as arg.
func DeployGateway(ctx context.Context, client *gatewayclient.Clientset, namespace, gatewayClassName string, opts ...func(*gatewayv1.Gateway)) (*gatewayv1.Gateway, error) {
	// create a default gateway with a listener set to port 80 for HTTP traffic
	gw := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultGatewayName,
		},
		Spec: gatewayv1.GatewaySpec{
			GatewayClassName: gatewayv1.ObjectName(gatewayClassName),
			Listeners: []gatewayv1.Listener{
				{
					Name:     "http",
					Protocol: gatewayv1.HTTPProtocolType,
					Port:     gatewayv1.PortNumber(ktfkong.DefaultProxyTCPServicePort),
				},
			},
		},
	}

	// call all the modifiers passed as args
	for _, opt := range opts {
		opt(gw)
	}

	result, err := client.GatewayV1beta1().Gateways(namespace).Create(ctx, gw, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		err = client.GatewayV1beta1().Gateways(namespace).Delete(ctx, gw.Name, metav1.DeleteOptions{})
		if err != nil {
			return result, err
		}
		result, err = client.GatewayV1beta1().Gateways(namespace).Create(ctx, gw, metav1.CreateOptions{})
	}
	return result, err
}

// gatewayHealthCheck checks the gateway has been scheduled by KIC. This function is inspired by
// assert.eventually (https://github.com/stretchr/testify/blob/v1.7.5/assert/assertions.go#L1669-L1700)
// and patched with custom behavior to determine the health of the gateway and to return an error
// instead of failing (at the time of its call, we don't have any test instance to make fail yet).
func gatewayHealthCheck(ctx context.Context, client *gatewayclient.Clientset, gatewayName, namespace string) error {
	ch := make(chan bool, 1)

	timer := time.NewTimer(gatewayWaitTimeToVerifyScheduling)
	defer timer.Stop()

	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	for tick := ticker.C; ; {
		select {
		case <-timer.C:
			return fmt.Errorf("the gateway %s/%s did not become scheduled after %s", namespace, defaultGatewayName, gatewayWaitTimeToVerifyScheduling)
		case <-tick:
			tick = nil
			ch <- func() bool {
				gw, err := client.GatewayV1beta1().Gateways(namespace).Get(ctx, gatewayName, metav1.GetOptions{})
				exitOnErr(ctx, err)
				ok := util.CheckCondition(
					gw.Status.Conditions,
					util.ConditionType(gatewayv1.GatewayConditionProgrammed),
					util.ConditionReason(gatewayv1.GatewayReasonProgrammed),
					metav1.ConditionTrue,
					gw.Generation,
				)
				return ok
			}()
		case v := <-ch:
			if v {
				return nil
			}
			tick = ticker.C
		}
	}
}

// HTTPRouteMatchesAcceptedCallback is a testing helper functions that returns a callback which
// checks if the provided HTTPRoute has an Accepted condition with:
// - Status matching the provided 'accepted' boolean argument.
// - Reason matching the provided 'reason' string argument.
func HTTPRouteMatchesAcceptedCallback(t *testing.T, c *gatewayclient.Clientset, httpRoute *gatewayv1.HTTPRoute, accepted bool, reason gatewayv1.RouteConditionReason) func() bool {
	return func() bool {
		return httpRouteAcceptedConditionMatches(t, c, httpRoute, accepted, reason)
	}
}

func httpRouteAcceptedConditionMatches(t *testing.T, c *gatewayclient.Clientset, httpRoute *gatewayv1.HTTPRoute, accepted bool, reason gatewayv1.RouteConditionReason) bool {
	var err error
	httpRoute, err = c.GatewayV1beta1().HTTPRoutes(httpRoute.Namespace).Get(context.Background(), httpRoute.Name, metav1.GetOptions{})
	require.NoError(t, err)

	if len(httpRoute.Status.Parents) <= 0 {
		return false
	}

	var expectedStatus metav1.ConditionStatus = "False"
	if accepted {
		expectedStatus = "True"
	}

	for _, cond := range httpRoute.Status.Parents[0].Conditions {
		if cond.Type == string(gatewayv1.RouteConditionAccepted) &&
			cond.Status == expectedStatus &&
			cond.Reason == string(reason) {
			return true
		}
	}

	return false
}

// ListenersHaveNAttachedRoutesCallback checks that every listener has a given number of attachedRoutes.
// The attachedRoutesByListener parameter contains the number of expected acceptedRoutes for
// each listener's name.
func ListenersHaveNAttachedRoutesCallback(t *testing.T, c *gatewayclient.Clientset, namespace, gatewayName string, attachedRoutesByListener map[string]int32) func() bool {
	return func() bool {
		gateway, err := c.GatewayV1beta1().Gateways(namespace).Get(context.Background(), gatewayName, metav1.GetOptions{})
		assert.NoError(t, err)

		for _, listenerStatus := range gateway.Status.Listeners {
			if attachedRoutesByListener[string(listenerStatus.Name)] != listenerStatus.AttachedRoutes {
				return false
			}
		}
		return true
	}
}

// GetGatewayIsLinkedCallback returns a callback that checks if the specific Route (HTTP, TCP, TLS, or UDP)
// is correctly linked to a supported gateway.
func GetGatewayIsLinkedCallback(
	ctx context.Context,
	t *testing.T,
	c *gatewayclient.Clientset,
	protocolType gatewayv1.ProtocolType,
	namespace,
	name string,
) func() bool {
	return func() bool {
		return gatewayLinkStatusMatches(ctx, t, c, true, protocolType, namespace, name)
	}
}

// GetGatewayIsUnlinkedCallback returns a callback that checks if the specific Route (HTTP, TCP, TLS, or UDP)
// is correctly unlinked from a supported gateway.
func GetGatewayIsUnlinkedCallback(
	ctx context.Context,
	t *testing.T,
	c *gatewayclient.Clientset,
	protocolType gatewayv1.ProtocolType,
	namespace,
	name string,
) func() bool {
	return func() bool {
		return gatewayLinkStatusMatches(ctx, t, c, false, protocolType, namespace, name)
	}
}

type routeParents struct {
	parents []gatewayv1.RouteParentStatus
}

func newRouteParentsStatus(parents []gatewayv1.RouteParentStatus) routeParents {
	return routeParents{
		parents: parents,
	}
}

func (rp routeParents) check(verifyLinked bool, controllerName string) bool {
	for _, ps := range rp.parents {
		if string(ps.ControllerName) == controllerName {
			// supported Gateway link was found, hence if we want to ensure
			// the link existence return true
			return verifyLinked
		}
	}

	// supported Gateway link was not found, hence if we want to ensure
	// the link existence return false
	return !verifyLinked
}

// gatewayLinkStatusMatches checks if the specific Route (HTTP, TCP, TLS, or UDP)
// is correctly linked to (or unlinked from) a supported gateway. In order to assert
// that the route must be linked to the gateway, or unlinked from the gateway, the
// verifyLinked boolean arg must be set accordingly.
func gatewayLinkStatusMatches(
	ctx context.Context,
	t *testing.T,
	c *gatewayclient.Clientset,
	verifyLinked bool,
	protocolType gatewayv1.ProtocolType,
	namespace, name string,
) bool {
	// gather a fresh copy of the route, given the specific protocol type
	switch protocolType { //nolint:exhaustive
	case gatewayv1.HTTPProtocolType:
		route, err := c.GatewayV1beta1().HTTPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		groute, gerr := c.GatewayV1alpha2().GRPCRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil && gerr != nil {
			t.Logf("error getting http route: %v", err)
			t.Logf("error getting grpc route: %v", err)
		} else {
			if err == nil {
				return newRouteParentsStatus(route.Status.Parents).
					check(verifyLinked, string(gateway.GetControllerName()))
			}
			if gerr == nil {
				return newRouteParentsStatus(groute.Status.Parents).
					check(verifyLinked, string(gateway.GetControllerName()))
			}
		}
	case gatewayv1.TCPProtocolType:
		route, err := c.GatewayV1alpha2().TCPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting tcp route: %v", err)
		} else {
			return newRouteParentsStatus(route.Status.Parents).
				check(verifyLinked, string(gateway.GetControllerName()))
		}
	case gatewayv1.UDPProtocolType:
		route, err := c.GatewayV1alpha2().UDPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting udp route: %v", err)
		} else {
			return newRouteParentsStatus(route.Status.Parents).
				check(verifyLinked, string(gateway.GetControllerName()))
		}
	case gatewayv1.TLSProtocolType:
		route, err := c.GatewayV1alpha2().TLSRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting tls route: %v", err)
		} else {
			return newRouteParentsStatus(route.Status.Parents).
				check(verifyLinked, string(gateway.GetControllerName()))
		}
	default:
		t.Fatalf("protocol %s not supported", string(protocolType))
	}

	t.Fatal("this should not happen")
	return false
}

func parentStatusContainsProgrammedCondition(
	parentStatuses []gatewayv1.RouteParentStatus,
	controllerName gatewayv1.GatewayController,
	expectedStatus metav1.ConditionStatus,
) bool {
	var conditions []metav1.Condition
	parentFound := false
	for _, parentStatus := range parentStatuses {
		if parentStatus.ControllerName == controllerName {
			conditions = parentStatus.Conditions
			parentFound = true
		}

		if parentFound {
			break
		}
	}

	if !parentFound {
		return false
	}
	return lo.ContainsBy(conditions, func(cond metav1.Condition) bool {
		return cond.Type == "Programmed" && cond.Status == expectedStatus
	})
}

func verifyProgrammedConditionStatus(t *testing.T,
	c *gatewayclient.Clientset,
	protocolType gatewayv1.ProtocolType,
	namespace, name string,
	expectedStatus metav1.ConditionStatus,
) bool {
	ctx := context.Background()

	// gather a fresh copy of the route, given the specific protocol type
	switch protocolType { //nolint:exhaustive
	case gatewayv1.HTTPProtocolType:
		route, err := c.GatewayV1beta1().HTTPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		groute, gerr := c.GatewayV1alpha2().GRPCRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil && gerr != nil {
			t.Logf("error getting http route: %v", err)
			t.Logf("error getting grpc route: %v", err)
		} else {
			if err == nil {
				return parentStatusContainsProgrammedCondition(route.Status.Parents, gateway.GetControllerName(), expectedStatus)
			}
			if gerr == nil {
				return parentStatusContainsProgrammedCondition(groute.Status.Parents, gateway.GetControllerName(), expectedStatus)
			}
		}
	case gateway.TCPProtocolType:
		route, err := c.GatewayV1alpha2().TCPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting tcp route: %v", err)
		} else {
			return parentStatusContainsProgrammedCondition(route.Status.Parents, gateway.GetControllerName(), expectedStatus)
		}
	case gateway.TLSProtocolType:
		route, err := c.GatewayV1alpha2().TLSRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting tls route: %v", err)
		} else {
			return parentStatusContainsProgrammedCondition(route.Status.Parents, gateway.GetControllerName(), expectedStatus)
		}
	case gateway.UDPProtocolType:
		route, err := c.GatewayV1alpha2().UDPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting udp route: %v", err)
		} else {
			return parentStatusContainsProgrammedCondition(route.Status.Parents, gateway.GetControllerName(), expectedStatus)
		}
	default:
		t.Fatalf("protocol %s not supported", string(protocolType))
	}

	return false
}

func GetVerifyProgrammedConditionCallback(t *testing.T,
	c *gatewayclient.Clientset,
	protocolType gatewayv1.ProtocolType,
	namespace, name string,
	expectedStatus metav1.ConditionStatus,
) func() bool {
	return func() bool {
		return verifyProgrammedConditionStatus(t, c, protocolType, namespace, name, expectedStatus)
	}
}

// setIngressClassNameWithRetry changes Ingress.Spec.IngressClassName to specified value
// and retries if update conflict happens.
func setIngressClassNameWithRetry(ctx context.Context, namespace string, ingress *netv1.Ingress, ingressClassName *string) error {
	ingressClient := env.Cluster().Client().NetworkingV1().Ingresses(namespace)
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		ingress, err := ingressClient.Get(ctx, ingress.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		ingress.Spec.IngressClassName = ingressClassName
		_, err = ingressClient.Update(ctx, ingress, metav1.UpdateOptions{})
		return err
	})
}
