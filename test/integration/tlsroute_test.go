//go:build integration_tests
// +build integration_tests

package integration

import (
	"bytes"
	"crypto/tls"
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
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
)

const (
	tlsRouteHostname = "tlsroute.kong.example"
)

func TestTLSRouteEssentials(t *testing.T) {
	backendPort := gatewayv1alpha2.PortNumber(tcpEchoPort)
	t.Log("locking TLS port")
	tlsMutex.Lock()
	defer func() {
		t.Log("unlocking TLS port")
		tlsMutex.Unlock()
	}()
	ns, cleanup := namespace(t)
	defer cleanup()

	// TODO consolidate into suite and use for all GW tests?
	// https://github.com/Kong/kubernetes-ingress-controller/issues/2461
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
				Name:     "tls",
				Protocol: gatewayv1alpha2.TLSProtocolType,
				Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTLSServicePort),
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

	t.Log("creating a tcpecho pod to test TLSRoute traffic routing")
	container1 := generators.NewContainer("tcpecho-1", tcpEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID1 := uuid.NewString()
	container1.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID1,
		},
	}
	deployment1 := generators.NewDeploymentForContainer(container1)
	deployment1, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment1, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("creating an additional tcpecho pod to test TLSRoute multiple backendRef loadbalancing")
	container2 := generators.NewContainer("tcpecho-2", tcpEchoImage, tcpEchoPort)
	// go-echo sends a "Running on Pod <UUID>." immediately on connecting
	testUUID2 := uuid.NewString()
	container2.Env = []corev1.EnvVar{
		{
			Name:  "POD_NAME",
			Value: testUUID2,
		},
	}
	deployment2 := generators.NewDeploymentForContainer(container2)
	deployment2, err = env.Cluster().Client().AppsV1().Deployments(ns.Name).Create(ctx, deployment2, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployments %s/%s and %s/%s", deployment1.Namespace, deployment1.Name, deployment2.Namespace, deployment2.Name)
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment1.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().AppsV1().Deployments(ns.Name).Delete(ctx, deployment2.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("exposing deployment %s/%s via service", deployment1.Namespace, deployment1.Name)
	service1 := generators.NewServiceForDeployment(deployment1, corev1.ServiceTypeLoadBalancer)
	service1, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service1, metav1.CreateOptions{})
	assert.NoError(t, err)

	t.Logf("exposing deployment %s/%s via service", deployment2.Namespace, deployment2.Name)
	service2 := generators.NewServiceForDeployment(deployment2, corev1.ServiceTypeLoadBalancer)
	service2, err = env.Cluster().Client().CoreV1().Services(ns.Name).Create(ctx, service2, metav1.CreateOptions{})
	assert.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service1.Name)
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service1.Name, metav1.DeleteOptions{}))
		assert.NoError(t, env.Cluster().Client().CoreV1().Services(ns.Name).Delete(ctx, service2.Name, metav1.DeleteOptions{}))
	}()

	t.Logf("creating a tlsroute to access deployment %s via kong", deployment1.Name)
	tlsroute := &gatewayv1alpha2.TLSRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:        uuid.NewString(),
			Annotations: map[string]string{},
		},
		Spec: gatewayv1alpha2.TLSRouteSpec{
			CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{
				ParentRefs: []gatewayv1alpha2.ParentReference{{
					Name: gatewayv1alpha2.ObjectName(gw.Name),
				}},
			},
			Hostnames: []gatewayv1alpha2.Hostname{tlsRouteHostname},
			Rules: []gatewayv1alpha2.TLSRouteRule{{
				BackendRefs: []gatewayv1alpha2.BackendRef{{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Name: gatewayv1alpha2.ObjectName(service1.Name),
						Port: &backendPort,
					},
				}},
			}},
		},
	}
	tlsroute, err = c.GatewayV1alpha2().TLSRoutes(ns.Name).Create(ctx, tlsroute, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the tlsroute %s", tlsroute.Name)
		if err := c.GatewayV1alpha2().TLSRoutes(ns.Name).Delete(ctx, tlsroute.Name, metav1.DeleteOptions{}); err != nil {
			if !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
			}
		}
	}()

	t.Log("verifying that the Gateway gets linked to the route via status")
	tlseventuallyGatewayIsLinkedInStatus(t, c, ns.Name, tlsroute.Name)

	t.Log("verifying that the tcpecho is responding properly over TLS")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID1, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("removing the parentrefs from the TLSRoute")
	oldParentRefs := tlsroute.Spec.ParentRefs
	require.Eventually(t, func() bool {
		tlsroute, err = c.GatewayV1alpha2().TLSRoutes(ns.Name).Get(ctx, tlsroute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		tlsroute.Spec.ParentRefs = nil
		tlsroute, err = c.GatewayV1alpha2().TLSRoutes(ns.Name).Update(ctx, tlsroute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	tlseventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, tlsroute.Name)

	t.Log("verifying that the tcpecho is no longer responding")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID1, tlsRouteHostname)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)

	t.Log("putting the parentRefs back")
	require.Eventually(t, func() bool {
		tlsroute, err = c.GatewayV1alpha2().TLSRoutes(ns.Name).Get(ctx, tlsroute.Name, metav1.GetOptions{})
		require.NoError(t, err)
		tlsroute.Spec.ParentRefs = oldParentRefs
		tlsroute, err = c.GatewayV1alpha2().TLSRoutes(ns.Name).Update(ctx, tlsroute, metav1.UpdateOptions{})
		return err == nil
	}, time.Minute, time.Second)

	t.Log("verifying that the Gateway gets linked to the route via status")
	tlseventuallyGatewayIsLinkedInStatus(t, c, ns.Name, tlsroute.Name)

	t.Log("verifying that putting the parentRefs back results in the routes becoming available again")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID1, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the GatewayClass")
	oldGWCName := gwc.Name
	require.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	tlseventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, tlsroute.Name)

	t.Log("verifying that the data-plane configuration from the TLSRoute gets dropped with the GatewayClass now removed")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID1, tlsRouteHostname)
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
	tlseventuallyGatewayIsLinkedInStatus(t, c, ns.Name, tlsroute.Name)

	t.Log("verifying that creating the GatewayClass again triggers reconciliation of TLSRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID1, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting the Gateway")
	oldGWName := gw.Name
	require.NoError(t, c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	tlseventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, tlsroute.Name)

	t.Log("verifying that the data-plane configuration from the TLSRoute gets dropped with the Gateway now removed")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID1, tlsRouteHostname)
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
				Name:     "tls",
				Protocol: gatewayv1alpha2.TLSProtocolType,
				Port:     gatewayv1alpha2.PortNumber(ktfkong.DefaultTLSServicePort),
			}},
		},
	}
	gw, err = c.GatewayV1alpha2().Gateways(ns.Name).Create(ctx, gw, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Log("verifying that the Gateway gets linked to the route via status")
	tlseventuallyGatewayIsLinkedInStatus(t, c, ns.Name, tlsroute.Name)

	t.Log("verifying that creating the Gateway again triggers reconciliation of TLSRoutes and the route becomes available again")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID1, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("adding an additional backendRef to the TLSRoute")
	require.Eventually(t, func() bool {
		tlsroute, err = c.GatewayV1alpha2().TLSRoutes(ns.Name).Get(ctx, tlsroute.Name, metav1.GetOptions{})
		require.NoError(t, err)

		tlsroute.Spec.Rules[0].BackendRefs = []gatewayv1alpha2.BackendRef{
			{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name: gatewayv1alpha2.ObjectName(service1.Name),
					Port: &backendPort,
				},
			},
			{
				BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
					Name: gatewayv1alpha2.ObjectName(service2.Name),
					Port: &backendPort,
				},
			},
		}

		tlsroute, err = c.GatewayV1alpha2().TLSRoutes(ns.Name).Update(ctx, tlsroute, metav1.UpdateOptions{})
		return err == nil
	}, ingressWait, waitTick)

	t.Log("verifying that the TLSRoute is now load-balanced between two services")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID1, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID2, tlsRouteHostname)
		return err == nil && responded == true
	}, ingressWait, waitTick)

	t.Log("deleting both GatewayClass and Gateway rapidly")
	require.NoError(t, c.GatewayV1alpha2().GatewayClasses().Delete(ctx, gwc.Name, metav1.DeleteOptions{}))
	require.NoError(t, c.GatewayV1alpha2().Gateways(ns.Name).Delete(ctx, gw.Name, metav1.DeleteOptions{}))

	t.Log("verifying that the Gateway gets unlinked from the route via status")
	tlseventuallyGatewayIsUnlinkedInStatus(t, c, ns.Name, tlsroute.Name)

	t.Log("verifying that the data-plane configuration from the TLSRoute does not get orphaned with the GatewayClass and Gateway gone")
	require.Eventually(t, func() bool {
		responded, err := tlsEchoResponds(fmt.Sprintf("%s:%d", proxyURL.Hostname(), ktfkong.DefaultTLSServicePort),
			testUUID1, tlsRouteHostname)
		return responded == false && errors.Is(err, io.EOF)
	}, ingressWait, waitTick)
}

// TODO consolidate shared util gateway linked funcs
// https://github.com/Kong/kubernetes-ingress-controller/issues/2461
func tlseventuallyGatewayIsLinkedInStatus(t *testing.T, c *gatewayclient.Clientset, namespace, name string) {
	require.Eventually(t, func() bool {
		// gather a fresh copy of the TLSRoute
		tlsroute, err := c.GatewayV1alpha2().TLSRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		require.NoError(t, err)

		// determine if there is a link to a supported Gateway
		for _, parentStatus := range tlsroute.Status.Parents {
			if parentStatus.ControllerName == gateway.ControllerName {
				// supported Gateway link was found
				return true
			}
		}

		// if no link was found yet retry
		return false
	}, ingressWait, waitTick)
}

// TODO https://github.com/Kong/kubernetes-ingress-controller/issues/2461
func tlseventuallyGatewayIsUnlinkedInStatus(t *testing.T, c *gatewayclient.Clientset, namespace, name string) {
	require.Eventually(t, func() bool {
		// gather a fresh copy of the TLSRoute
		tlsroute, err := c.GatewayV1alpha2().TLSRoutes(namespace).Get(ctx, name, metav1.GetOptions{})
		require.NoError(t, err)

		// determine if there is a link to a supported Gateway
		for _, parentStatus := range tlsroute.Status.Parents {
			if parentStatus.ControllerName == gateway.ControllerName {
				// a supported Gateway link was found retry
				return false
			}
		}

		// linked gateway is not present, all set
		return true
	}, ingressWait, waitTick)
}

// tlsEchoResponds takes a TLS address URL and a Pod name and checks if a
// go-echo instance is running on that Pod at that address using hostname for SNI.
// It compares an expected message and its length against an expected message, returning true
// if it is and false and an error explanation if it is not
func tlsEchoResponds(url string, podName string, hostname string) (bool, error) {
	dialer := net.Dialer{Timeout: time.Second * 10}
	conn, err := tls.DialWithDialer(&dialer,
		"tcp",
		url,
		&tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
			ServerName:         hostname,
		})
	if err != nil {
		return false, err
	}
	defer conn.Close()

	header := []byte(fmt.Sprintf("Running on Pod %s.", podName))
	message := []byte("testing tlsroute")

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
