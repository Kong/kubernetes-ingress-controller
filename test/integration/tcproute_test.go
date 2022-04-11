//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
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
	"k8s.io/apimachinery/pkg/util/intstr"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
)

const tcpEchoPort = 1025

func TestTCPRouteEssentials(t *testing.T) {
	ns, cleanup := namespace(t)
	defer cleanup()
	t.Log("locking TCP port")
	tcpMutex.Lock()
	defer func() {
		// Free up the TCP port
		t.Log("unlocking TCP port")
		tcpMutex.Unlock()
	}()

	// TODO consolidate into suite and use for all GW tests?
	t.Log("deploying a supported gatewayclass to the test cluster")
	c, err := gatewayclient.NewForConfig(env.Cluster().Config())
	require.NoError(t, err)
	gwc := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	gwc, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up gatewayclasses")
		if err := c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("deploying a gateway to the test cluster using unmanaged gateway mode")
	gw := &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kong",
			Annotations: map[string]string{
				unmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(gwc.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "tcp",
				Protocol: gatewayv1alpha2.TCPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTCPServicePort),
			}},
		},
	}
	gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Log("cleaning up gateways")
		if err := c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	container := generators.NewContainer("tcpecho", tcpEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID := uuid.NewString()
	container.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID,
		},
	}
	deployment := generators.NewDeploymentForContainer(container)
	deployment, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	// we have to override the ports so that we can map the default TCP port from
	// the Kong Gateway deployment to the tcpecho port, as this is what will be
	// used to route the traffic at the Gateway (at the time of writing, the
	// Kong Gateway doesn't support an API for dynamically adding these ports. The
	// ports must be added manually to the config or ENV).
	service.Spec.Ports = []corev1.ServicePort{{
		Name:       "tcp",
		Protocol:   corev1.ProtocolTCP,
		Port:       ktfkong.DefaultTCPServicePort,
		TargetPort: intstr.FromInt(tcpEchoPort),
	}}
	service, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating a tcproute to access deployment %s via kong", deployment.Name)
	tcpPortDefault := gatewayv1alpha2.PortNumber(ktfkong.DefaultTCPServicePort)
	tcproute := &gatewayv1alpha2.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:        uuid.NewString(),
			Annotations: map[string]string{},
		},
		Spec: gatewayv1alpha2.TCPRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gw.Name),
				}},
			},
			Rules: []gatewayv1alpha2.TCPRouteRule{{
				BackendRefs: []gatewayv1alpha2.BackendRef{{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service.Name),
						Port: &tcpPortDefault,
					},
				}},
			}},
		},
	}
	tcproute, err = c.GatewayV1alpha2().TCPRoutes(ns.Name).Create(ctx, tcproute, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the tcproute %s", tcproute.Name)
		if err := c.GatewayV1alpha2().TCPRoutes(ns.Name).Delete(ctx, tcproute.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("verifying that the Gateway gets linked to the route via status")
	tcpeventuallyGatewayIsLinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that the tcpecho is responding properly")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("removing the parentrefs from the TCPRoute")
	oldParentRefs := tcproute.Spec.ParentRefs
	require.Eventually(t, func() bool {
		tcproute, err = c.GatewayV1alpha2().TCPRoutes(ns.Name).Get(ctx, tcproute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		tcproute.Spec.ParentRefs = nil
		tcproute, err = c.GatewayV1alpha2().TCPRoutes(ns.Name).Update(ctx, tcproute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	tcpeventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that the tcpecho is no longer responding")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the parentRefs back")
	require.Eventually(t, func() bool {
		tcproute, err = c.GatewayV1alpha2().TCPRoutes(ns.Name).Get(ctx, tcproute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		tcproute.Spec.ParentRefs = oldParentRefs
		tcproute, err = c.GatewayV1alpha2().TCPRoutes(ns.Name).Update(ctx, tcproute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets linked to the route via status")
	tcpeventuallyGatewayIsLinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the GatewayClass")
	oldGWCName := gwc.Name
	require.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	tcpeventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that the data-plane configuration from the TCPRoute gets dropped with the GatewayClass now removed")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the GatewayClass back")
	gwc = &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: oldGWCName,
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	gwc, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	tcpeventuallyGatewayIsLinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of TCPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the Gateway")
	oldGWName := gw.Name
	require.NoError(t, c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	tcpeventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that the data-plane configuration from the TCPRoute gets dropped with the Gateway now removed")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the Gateway back")
	gw = &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: oldGWName,
			Annotations: map[string]string{
				unmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(gwc.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "tcp",
				Protocol: gatewayv1alpha2.TCPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(9999),
			}},
		},
	}
	gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	tcpeventuallyGatewayIsLinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that creating the Gateway again triggers reconciliation of TCPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	tcpeventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that the data-plane configuration from the TCPRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the GatewayClass back")
	gwc = &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: oldGWCName,
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	gwc, err = c.GatewayV1alpha2().GatewayClasses().Create(ctx, gwc, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("putting the Gateway back")
	gw = &gatewayv1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name: oldGWName,
			Annotations: map[string]string{
				unmanagedAnnotation: "true", // trigger the unmanaged gateway mode
			},
		},
		Spec: gatewayv1alpha2.GatewaySpec{
			GatewayClassName: gatewayv1alpha2.ObjectName(gwc.Name),
			Listeners: []gatewayv1alpha2.Listener{{
				Name:     "tcp",
				Protocol: gatewayv1alpha2.TCPProtocolType,
				Port:     gatewayv1alpha2.PortNumber(9999),
			}},
		},
	}
	gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	tcpeventuallyGatewayIsLinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that creating the Gateway again triggers reconciliation of TCPRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	tcpeventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, tcproute.Name)

	t.Log("verifying that the data-plane configuration from the TCPRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		responded, err := tcpEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTCPServicePort), testUUID)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)
}

// TODO consolidate shared util gateway linked funcs
func tcpeventuallyGatewayIsLinkedInStatus(t *testing.T, c *gatewayclient.Clientset, namespace, name string) {
	require.Eventually(t, func() bool {
		// gather a fresh copy of the TCPRoute
		tcproute, err := c.GatewayV1alpha2().TCPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		require.NoError(t, err)

		// determine if there is a link to a supported Gateway
		for _, parentStatus := range tcproute.Status.Parents {
			if parentStatus.ControllerName == gateway.ControllerName {
				// supported Gateway link was found
				return true
			}
		}

		// if no link was found yet retry
		return false
	}, ingressWait, waitTick)
}

func tcpeventuallyGatewayIsUnlinkedInStatus(t *testing.T, c *gatewayclient.Clientset, namespace, name string) {
	require.Eventually(t, func() bool {
		// gather a fresh copy of the TCPRoute
		tcproute, err := c.GatewayV1alpha2().TCPRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		require.NoError(t, err)

		// determine if there is a link to a supported Gateway
		for _, parentStatus := range tcproute.Status.Parents {
			if parentStatus.ControllerName == gateway.ControllerName {
				// a supported Gateway link was found retry
				return false
			}
		}

		// linked gateway is not present, all set
		return true
	}, ingressWait, waitTick)
}

// tcpEchoResponds takes a TCP address URL and a Pod name and checks if a
// go-echo instance is running on that Pod at that address. It compares an
// expected message and its length against an expected message, returning true
// if it is and false and an error explanation if it is not
func tcpEchoResponds(url string, podName string) (bool, error) {
	dialer := net.Dialer{Timeout: time.Second * 10}
	conn, err := dialer.Dial("tcp", url)
	if err != nil {
		return false, err
	}

	header := []byte(fmt.Sprintf("Running on Pod %s.", podName))
	message := []byte("testing tcproute")

	wrote, err := conn.Write(message)
	if err != nil {
		return false, err
	}

	if wrote != len(message) {
		return false, fmt.Errorf("wrote message of size %d, expected %d", wrote, len(message))
	}

	if err := conn.SetDeadline(time.Now().Add(time.Second * 5)); err != nil {
		return false, err
	}

	headerResponse := make([]byte, len(header)+1)
	read, err := conn.Read(headerResponse)
	if err != nil {
		return false, err
	}

	if read != len(header)+1 { // add 1 for newline
		return false, fmt.Errorf("read %d bytes but expected %d", read, len(header)+1)
	}

	if !bytes.Contains(headerResponse, header) {
		return false, fmt.Errorf(`expected header response "%s", received: "%s"`, string(header), string(headerResponse))
	}

	messageResponse := make([]byte, wrote+1)
	read, err = conn.Read(messageResponse)
	if err != nil {
		return false, err
	}

	if read != len(message) {
		return false, fmt.Errorf("read %d bytes but expected %d", read, len(message))
	}

	if !bytes.Contains(messageResponse, message) {
		return false, fmt.Errorf(`expected message response "%s", received: "%s"`, string(message), string(messageResponse))
	}

	return true, nil
}
