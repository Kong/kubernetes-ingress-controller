//go:build integration_tests

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	"github.com/kong/kubernetes-ingress-controller/v3/test/internal/helpers"
)

const (
	// defaultGatewayName is the default name for the Gateways created during tests.
	defaultGatewayName = helpers.DefaultGatewayName
	// unmanagedGatewayClassName is the name of the default GatewayClass created during the test environment setup.
	unmanagedGatewayClassName = "kong-unmanaged"
	// unsupportedControllerName is the name of the controller used for those gateways that are not supported
	// by an actual controller (i.e., they won't be scheduled).
	unsupportedControllerName gatewayapi.GatewayController = "example.com/unsupported-gateway-controller"
	// kongRouterFlavorExpressions is the value used in router_flavor of kong configuration
	// to enable expression based router of kong.
	kongRouterFlavorExpressions = "expressions"
)

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
				gw, err := client.GatewayV1().Gateways(namespace).Get(ctx, gatewayName, metav1.GetOptions{})
				helpers.ExitOnErr(ctx, err)
				ok := util.CheckCondition(
					gw.Status.Conditions,
					util.ConditionType(gatewayapi.GatewayConditionProgrammed),
					util.ConditionReason(gatewayapi.GatewayReasonProgrammed),
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
func HTTPRouteMatchesAcceptedCallback(t *testing.T, c *gatewayclient.Clientset, httpRoute *gatewayapi.HTTPRoute, accepted bool, reason gatewayapi.RouteConditionReason) func() bool {
	return func() bool {
		return httpRouteAcceptedConditionMatches(t, c, httpRoute, accepted, reason)
	}
}

func httpRouteAcceptedConditionMatches(t *testing.T, c *gatewayclient.Clientset, httpRoute *gatewayapi.HTTPRoute, accepted bool, reason gatewayapi.RouteConditionReason) bool {
	var err error
	httpRoute, err = c.GatewayV1().HTTPRoutes(httpRoute.Namespace).Get(context.Background(), httpRoute.Name, metav1.GetOptions{})
	require.NoError(t, err)

	if len(httpRoute.Status.Parents) == 0 {
		return false
	}

	var expectedStatus metav1.ConditionStatus = "False"
	if accepted {
		expectedStatus = "True"
	}

	for _, cond := range httpRoute.Status.Parents[0].Conditions {
		if cond.Type == string(gatewayapi.RouteConditionAccepted) &&
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
		gateway, err := c.GatewayV1().Gateways(namespace).Get(context.Background(), gatewayName, metav1.GetOptions{})
		assert.NoError(t, err)

		for _, listenerStatus := range gateway.Status.Listeners {
			if attachedRoutesByListener[string(listenerStatus.Name)] != listenerStatus.AttachedRoutes {
				return false
			}
		}
		return true
	}
}
