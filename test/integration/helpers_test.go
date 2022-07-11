//go:build integration_tests
// +build integration_tests

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
)

const (
	// defaultGatewayName is the default name for the Gateways created during tests
	defaultGatewayName = "kong"
	// managedGatewayClassName is the name of the default GatewayClass created during the test environment setup
	managedGatewayClassName = "kong-managed"
	// unmanagedControllerName is the name of the controller used for those gateways that are not supported
	// by an actual controller (i.e., they won't be scheduled)
	unmanagedControllerName gatewayv1alpha2.GatewayController = "example.com/unmanaged-gateway-controller"
)

// DeployGateway creates a default gatewayClass, accepts a variadic set of options,
// and finally deploys it on the Kubernetes cluster by means of the gateway client given as arg.
func DeployGatewayClass(ctx context.Context, client *gatewayclient.Clientset, gatewayClassName string, opts ...func(*gatewayv1alpha2.GatewayClass)) (*gatewayv1alpha2.GatewayClass, error) {
	gwc := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: gatewayClassName,
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}

	for _, opt := range opts {
		opt(gwc)
	}

	return client.GatewayV1alpha2().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
}

// DeployGateway creates a default gateway, accepts a variadic set of options,
// and finally deploys it on the Kubernetes cluster by means of the gateway client given as arg.
func DeployGateway(ctx context.Context, client *gatewayclient.Clientset, namespace, gatewayClassName string, opts ...func(*gatewayv1alpha2.Gateway)) (*gatewayv1alpha2.Gateway, error) {
	// create a default gateway with a listener set to port 80 for HTTP traffic
	gw := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultGatewayName,
			Annotations: map[string]string{
				unmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(gatewayClassName),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "http",
				Protocol: gatewayv1alpha2.HTTPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultProxyTCPServicePort),
			}},
		},
	}

	// call all the modifiers passed as args
	for _, opt := range opts {
		opt(gw)
	}

	return client.GatewayV1alpha2().Gateways(namespace).Create(ctx, gw, metav1.CreateOptions{})
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
				gw, err := client.GatewayV1alpha2().Gateways(namespace).Get(ctx, gatewayName, metav1.GetOptions{})
				exitOnErr(err)
				for _, cond := range gw.Status.Conditions {
					if cond.Reason == string(gatewayv1alpha2.GatewayReasonReady) {
						return true
					}
				}
				return false
			}()
		case v := <-ch:
			if v {
				return nil
			}
			tick = ticker.C
		}
	}
}

// GetGatewayIsLinkedCallback returns a callback that checks if the specific Route (HTTP, TCP, TLS, or UDP)
// is correctly linked to a supported gateway
func GetGatewayIsLinkedCallback(t *testing.T, c *gatewayclient.Clientset, protocolType gatewayv1alpha2.ProtocolType, namespace, name string) func() bool {
	return func() bool {
		return gatewayLinkStatusMatches(t, c, true, protocolType, namespace, name)
	}
}

// GetGatewayIsUnlinkedCallback returns a callback that checks if the specific Route (HTTP, TCP, TLS, or UDP)
// is correctly unlinked from a supported gateway
func GetGatewayIsUnlinkedCallback(t *testing.T, c *gatewayclient.Clientset, protocolType gatewayv1alpha2.ProtocolType, namespace, name string) func() bool {
	return func() bool {
		return gatewayLinkStatusMatches(t, c, false, protocolType, namespace, name)
	}
}

// gatewayLinkStatusMatches checks if the specific Route (HTTP, TCP, TLS, or UDP)
// is correctly linked to (or unlinked from) a supported gateway. In order to assert
// that the route must be linked to the gateway, or unlinked from the gateway, the
// verifyLinked boolean arg must be set accordingly.
func gatewayLinkStatusMatches(t *testing.T, c *gatewayclient.Clientset, verifyLinked bool, protocolType gatewayv1alpha2.ProtocolType, namespace, name string) bool {
	var routeParents []gatewayv1alpha2.RouteParentStatus

	// gather a fresh copy of the route, given the specific protocol type
	switch protocolType {
	case gatewayv1alpha2.HTTPProtocolType:
		route, err := c.GatewayV1alpha2().HTTPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		t.Logf("error getting http route: %v", err)
		routeParents = route.Status.Parents
	case gatewayv1alpha2.TCPProtocolType:
		route, err := c.GatewayV1alpha2().TCPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		t.Logf("error getting tcp route: %v", err)
		routeParents = route.Status.Parents
	case gatewayv1alpha2.UDPProtocolType:
		route, err := c.GatewayV1alpha2().UDPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		t.Logf("error getting udp route: %v", err)
		routeParents = route.Status.Parents
	case gatewayv1alpha2.TLSProtocolType:
		route, err := c.GatewayV1alpha2().TLSRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		t.Logf("error getting tls route: %v", err)
		routeParents = route.Status.Parents
	default:
		t.Fatalf("protocol %s not supported", string(protocolType))
	}

	// determine if there is a link to a supported Gateway
	for _, parentStatus := range routeParents {
		if parentStatus.ControllerName == gateway.ControllerName {
			// supported Gateway link was found, hence if we want to ensure
			// the link existence return true
			return verifyLinked
		}
	}

	// supported Gateway link was not found, hence if we want to ensure
	// the link existence return false
	return !verifyLinked
}
