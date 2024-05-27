package helpers

import (
	"context"
	"testing"

	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/samber/lo"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

const (
	// DefaultGatewayName is the default name for the Gateways created during tests.
	DefaultGatewayName = "kong"
)

// DeployGatewayClass creates a default gatewayClass, accepts a variadic set of options,
// and finally deploys it on the Kubernetes cluster by means of the gateway client given as arg.
func DeployGatewayClass(ctx context.Context, client *gatewayclient.Clientset, gatewayClassName string, opts ...func(*gatewayapi.GatewayClass)) (*gatewayapi.GatewayClass, error) {
	gwc := &gatewayapi.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: gatewayClassName,
			Annotations: map[string]string{
				annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
			},
		},
		Spec: gatewayapi.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
	}

	for _, opt := range opts {
		opt(gwc)
	}

	result, err := client.GatewayV1().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
	if apierrors.IsAlreadyExists(err) {
		err = client.GatewayV1().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{})
		if err != nil {
			return result, err
		}
		result, err = client.GatewayV1().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
	}
	return result, err
}

// DeployGateway creates a default gateway, accepts a variadic set of options,
// and finally deploys it on the Kubernetes cluster by means of the gateway client given as arg.
func DeployGateway(ctx context.Context, client *gatewayclient.Clientset, namespace, gatewayClassName string, opts ...func(*gatewayapi.Gateway)) (*gatewayapi.Gateway, error) {
	// Create a default gateway with a listener set to port 80 for HTTP traffic.
	gw := &gatewayapi.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: DefaultGatewayName,
		},
		Spec: gatewayapi.GatewaySpec{
			GatewayClassName: gatewayapi.ObjectName(gatewayClassName),
			Listeners: []gatewayapi.Listener{
				{
					Name:     "http",
					Protocol: gatewayapi.HTTPProtocolType,
					Port:     gatewayapi.PortNumber(ktfkong.DefaultProxyTCPServicePort),
				},
			},
		},
	}

	// Call all the modifiers passed as args.
	for _, opt := range opts {
		opt(gw)
	}

	result, err := client.GatewayV1().Gateways(namespace).Create(ctx, gw, metav1.CreateOptions{})
	if !apierrors.IsAlreadyExists(err) {
		return result, err
	}
	if err := client.GatewayV1().Gateways(namespace).Delete(ctx, gw.Name, metav1.DeleteOptions{}); err != nil {
		return nil, err
	}
	return client.GatewayV1().Gateways(namespace).Create(ctx, gw, metav1.CreateOptions{})
}

// GetGatewayIsLinkedCallback returns a callback that checks if the specific Route (HTTP, TCP, TLS, or UDP)
// is correctly linked to a supported gateway.
func GetGatewayIsLinkedCallback(
	ctx context.Context,
	t *testing.T,
	c *gatewayclient.Clientset,
	protocolType gatewayapi.ProtocolType,
	namespace,
	name string,
) func() bool {
	return func() bool {
		return gatewayLinkStatusMatches(ctx, t, c, true, protocolType, namespace, name)
	}
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
	protocolType gatewayapi.ProtocolType,
	namespace, name string,
) bool {
	// gather a fresh copy of the route, given the specific protocol type
	switch protocolType {
	case gatewayapi.HTTPProtocolType:
		route, err := c.GatewayV1().HTTPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		groute, gerr := c.GatewayV1alpha2().GRPCRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil && gerr != nil {
			t.Logf("error getting http route: %v", err)
			t.Logf("error getting grpc route: %v", gerr)
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
	case gatewayapi.TCPProtocolType:
		route, err := c.GatewayV1alpha2().TCPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting tcp route: %v", err)
		} else {
			return newRouteParentsStatus(route.Status.Parents).
				check(verifyLinked, string(gateway.GetControllerName()))
		}
	case gatewayapi.UDPProtocolType:
		route, err := c.GatewayV1alpha2().UDPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting udp route: %v", err)
		} else {
			return newRouteParentsStatus(route.Status.Parents).
				check(verifyLinked, string(gateway.GetControllerName()))
		}
	case gatewayapi.TLSProtocolType:
		route, err := c.GatewayV1alpha2().TLSRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting tls route: %v", err)
		} else {
			return newRouteParentsStatus(route.Status.Parents).
				check(verifyLinked, string(gateway.GetControllerName()))
		}
	case gatewayapi.HTTPSProtocolType:
		t.Fatalf("protocol %s not supported", protocolType)
	default:
		t.Fatalf("protocol %s not supported", protocolType)
	}

	t.Fatal("this should not happen")
	return false
}

func newRouteParentsStatus(parents []gatewayapi.RouteParentStatus) routeParents {
	return routeParents{
		parents: parents,
	}
}

type routeParents struct {
	parents []gatewayapi.RouteParentStatus
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

func GetVerifyProgrammedConditionCallback(t *testing.T,
	c *gatewayclient.Clientset,
	protocolType gatewayapi.ProtocolType,
	namespace, name string,
	expectedStatus metav1.ConditionStatus,
) func() bool {
	return func() bool {
		return verifyProgrammedConditionStatus(t, c, protocolType, namespace, name, expectedStatus)
	}
}

func verifyProgrammedConditionStatus(t *testing.T,
	c *gatewayclient.Clientset,
	protocolType gatewayapi.ProtocolType,
	namespace, name string,
	expectedStatus metav1.ConditionStatus,
) bool {
	ctx := context.Background()

	// gather a fresh copy of the route, given the specific protocol type
	switch protocolType {
	case gatewayapi.HTTPProtocolType:
		route, err := c.GatewayV1().HTTPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
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
	case gatewayapi.TCPProtocolType:
		route, err := c.GatewayV1alpha2().TCPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting tcp route: %v", err)
		} else {
			return parentStatusContainsProgrammedCondition(route.Status.Parents, gateway.GetControllerName(), expectedStatus)
		}
	case gatewayapi.TLSProtocolType:
		route, err := c.GatewayV1alpha2().TLSRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting tls route: %v", err)
		} else {
			return parentStatusContainsProgrammedCondition(route.Status.Parents, gateway.GetControllerName(), expectedStatus)
		}
	case gatewayapi.UDPProtocolType:
		route, err := c.GatewayV1alpha2().UDPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			t.Logf("error getting udp route: %v", err)
		} else {
			return parentStatusContainsProgrammedCondition(route.Status.Parents, gateway.GetControllerName(), expectedStatus)
		}
	case gatewayapi.HTTPSProtocolType:
		t.Fatalf("protocol %s not supported", string(protocolType))
	default:
		t.Fatalf("protocol %s not supported", string(protocolType))
	}

	return false
}

func parentStatusContainsProgrammedCondition(
	parentStatuses []gatewayapi.RouteParentStatus,
	controllerName gatewayapi.GatewayController,
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

// GetGatewayIsUnlinkedCallback returns a callback that checks if the specific Route (HTTP, TCP, TLS, or UDP)
// is correctly unlinked from a supported gateway.
func GetGatewayIsUnlinkedCallback(
	ctx context.Context,
	t *testing.T,
	c *gatewayclient.Clientset,
	protocolType gatewayapi.ProtocolType,
	namespace,
	name string,
) func() bool {
	return func() bool {
		return gatewayLinkStatusMatches(ctx, t, c, false, protocolType, namespace, name)
	}
}
