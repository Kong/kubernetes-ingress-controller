//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kong/go-kong/kong"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kong/kubernetes-ingress-controller/pkg/annotations"
	"github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
)

func TestPortCollisions(t *testing.T) {
	// TODO: we need to fix a bug which causes kong.Upstream collisions when UDPIngress and TCPIngress
	//       both share the same port, then we can re-enable this test.
	//       See: https://github.com/Kong/kubernetes-ingress-controller/issues/1446
	t.Skip("this test is temporarily disabled until we resolve " +
		"https://github.com/Kong/kubernetes-ingress-controller/issues/1446")

	// TODO: once KIC 2.0 lands and pre v2 is gone, we can remove this check
	if useLegacyKIC() {
		t.Skip("legacy KIC does not support UDPIngress, skipping")
	}

	testName := "collisionchk"
	namespace := corev1.NamespaceDefault
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	defer cancel()

	// -------------------------------------------------------------------------
	// TCP & UDP Deployment
	// -------------------------------------------------------------------------

	t.Log("configuring coredns corefile")
	cfgmap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: testName}, Data: map[string]string{"Corefile": corefile}}
	cfgmap, err := cluster.Client().CoreV1().ConfigMaps(namespace).Create(ctx, cfgmap, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the coredns corefile %s", cfgmap.Name)
		assert.NoError(t, cluster.Client().CoreV1().ConfigMaps(namespace).Delete(ctx, cfgmap.Name, metav1.DeleteOptions{}))
	}()

	t.Log("configuring TCP service container")
	tcpContainer := k8sgen.NewContainer(fmt.Sprintf("%s-tcp", testName), httpBinImage, 80)

	t.Log("configuring UDP service container")
	udpContainer := k8sgen.NewContainer(fmt.Sprintf("%s-udp", testName), "coredns/coredns", 53)
	udpContainer.Ports[0].Protocol = corev1.ProtocolUDP
	udpContainer.VolumeMounts = []corev1.VolumeMount{{Name: "config-volume", MountPath: "/etc/coredns"}}
	udpContainer.Args = []string{"-conf", "/etc/coredns/Corefile"}

	t.Log("deploying TCP and UDP containers")
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: testName,
			Labels: map[string]string{
				"app": testName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": testName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": testName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						tcpContainer,
						udpContainer,
					},
					Volumes: []corev1.Volume{{
						Name: "config-volume",
						VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: cfgmap.Name},
							Items:                []corev1.KeyToPath{{Key: "Corefile", Path: "Corefile"}},
						}},
					}},
				},
			},
		},
	}
	deployment, err = cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the deployment %s", deployment.Name)
		assert.NoError(t, cluster.Client().AppsV1().Deployments(namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	// -------------------------------------------------------------------------
	// TCP & UDP Service
	// -------------------------------------------------------------------------

	t.Logf("exposing deployment %s via service", deployment.Name)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: testName,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: deployment.Spec.Selector.MatchLabels,
			Ports: []corev1.ServicePort{
				{
					Name:     fmt.Sprintf("%s-tcp", testName),
					Protocol: corev1.ProtocolTCP,
					Port:     80,
				},
				{
					Name:     fmt.Sprintf("%s-udp", testName),
					Protocol: corev1.ProtocolUDP,
					// intentionally using the same port (different proto) from the TCP listener
					// to validate that this doesn't result in kong resource collisions.
					Port:       80,
					TargetPort: intstr.FromInt(53),
				},
			},
		},
	}
	service, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("cleaning up the service %s", service.Name)
		assert.NoError(t, cluster.Client().CoreV1().Services(namespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	// -------------------------------------------------------------------------
	// TCP & UDP Routing
	// -------------------------------------------------------------------------

	t.Logf("routing to service %s via TCPIngress", service.Name)
	c, err := clientset.NewForConfig(cluster.Config())
	require.NoError(t, err)
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: namespace,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Port: 8888,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: service.Name,
						ServicePort: 80,
					},
				},
			},
		},
	}
	tcp, err = c.ConfigurationV1beta1().TCPIngresses(namespace).Create(ctx, tcp, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("ensuring that TCPIngress %s is cleaned up", tcp.Name)
		if err := c.ConfigurationV1beta1().TCPIngresses(namespace).Delete(ctx, tcp.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	t.Log("exposing DNS service via UDPIngress")
	udp := &kongv1beta1.UDPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name: testName,
			Annotations: map[string]string{
				annotations.IngressClassKey: ingressClass,
			},
		},
		Spec: kongv1beta1.UDPIngressSpec{Rules: []kongv1beta1.UDPIngressRule{
			{
				Port: 9999,
				Backend: v1beta1.IngressBackend{
					ServiceName: service.Name,
					// again, the odd port number here is due to us intentionally validating that things
					// work even when the port number between two different pods in a service are the same
					// with a separate protocol between them.
					ServicePort: 80,
				},
			},
		}},
	}
	udp, err = c.ConfigurationV1beta1().UDPIngresses(namespace).Create(ctx, udp, metav1.CreateOptions{})
	require.NoError(t, err)

	defer func() {
		t.Logf("ensuring UDPIngress %s is cleaned up", udp.Name)
		if err := c.ConfigurationV1beta1().UDPIngresses(namespace).Delete(ctx, udp.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				require.NoError(t, err)
			}
		}
	}()

	// -------------------------------------------------------------------------
	// Validation
	// -------------------------------------------------------------------------

	/*
		TODO: add validation that the routes actually resolve

		At the time of writing the tests here are "short circuited" due to a problem
		that was encountered with how Kong Upstreams appear to work when Kong Services
		have two protocols sharing a port. This problem is being investigated further
		but for the purposes of unblocking ourselves we are planning to expand these
		tests once a new feature is added to UDP/TCP Ingresses to support the following
		annotation:

		- ingress.kubernetes.io/service-upstream

		The purpose will be to enable flagging UDPIngress (and TCPIngress) resources
		to use the Kubernetes service DNS directly instead of using a Kong upstream
		which will give a path to work around the problem.

		See: https://github.com/Kong/kubernetes-ingress-controller/issues/1441
	*/

	t.Logf("pulling the kong admin apis /services endpoint data to validate the generated services")
	p := proxyReady()
	var servicesResponse *http.Response
	var kongServices kongServiceList
	assert.Eventually(t, func() bool {
		var err error
		servicesResponse, err = http.Get(fmt.Sprintf("%s/services", p.ProxyAdminURL.String()))
		require.NoError(t, err)
		defer servicesResponse.Body.Close()
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(servicesResponse.Body)
		require.NoError(t, err)

		require.NoError(t, json.Unmarshal(buf.Bytes(), &kongServices))
		return len(kongServices.Data) == 2
	}, ingressWait, waitTick)

	t.Logf("validing that the resulting Kong services for same-named UDP and TCP ingress are not colliding")
	assert.NotEqual(t, kongServices.Data[0].Name, kongServices.Data[1].Name)

	t.Logf("pulling the kong admin apis /routes endpoint data to validate the generated routes")
	routesResponse, err := http.Get(fmt.Sprintf("%s/services", p.ProxyAdminURL.String()))
	require.NoError(t, err)
	defer routesResponse.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(routesResponse.Body)
	require.NoError(t, err)

	t.Logf("validing that the resulting Kong routes for same-named UDP and TCP ingress are not colliding")
	kongRoutes := kongRouteList{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &kongRoutes))
	assert.Len(t, kongRoutes.Data, 2)
	assert.NotEqual(t, kongRoutes.Data[0].Name, kongRoutes.Data[1].Name)

	// -------------------------------------------------------------------------
	// Teardown & Cleanup
	// -------------------------------------------------------------------------

	t.Logf("tearing down UDPIngress %s and ensuring backends are torn down", udp.Name)
	require.NoError(t, c.ConfigurationV1beta1().UDPIngresses(namespace).Delete(ctx, udp.Name, metav1.DeleteOptions{}))

	t.Logf("tearing down TCPIngress %s and ensuring that the relevant backend routes are removed", tcp.Name)
	require.NoError(t, c.ConfigurationV1beta1().TCPIngresses(namespace).Delete(ctx, tcp.Name, metav1.DeleteOptions{}))
}

type kongServiceList struct {
	Data []kong.Service `json:"data,required"`
}

type kongRouteList struct {
	Data []kong.Route `json:"data,required"`
}
