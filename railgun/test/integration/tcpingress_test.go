//+build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/controllers"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestMinimalTCPIngress(t *testing.T) {
	ctx := context.Background()

	// gather the proxy container as it will need to be specially configured to serve TCP
	proxy, err := cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Get(ctx, "ingress-controller-kong", metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Len(t, proxy.Spec.Template.Spec.Containers, 1)
	container := proxy.Spec.Template.Spec.Containers[0].DeepCopy()

	// override the KONG_STREAM_LISTEN env var in the proxy container
	originalVal, err := overrideEnvVar(container, "KONG_STREAM_LISTEN", "0.0.0.0:32080")
	assert.NoError(t, err)
	proxy.Spec.Template.Spec.Containers[0] = *container

	// add the TCP port to the pod
	container.Ports = append(container.Ports, corev1.ContainerPort{
		Name:          "tcp-test",
		ContainerPort: 32080,
		Protocol:      corev1.ProtocolTCP,
	})

	// update the deployment with the new container configurations
	proxy, err = cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Update(ctx, proxy, metav1.UpdateOptions{})
	assert.NoError(t, err)

	// make sure we clean up after ourselves
	defer func() {
		// retrieve the current proxy
		proxy, err := cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Get(ctx, "ingress-controller-kong", metav1.GetOptions{})
		assert.NoError(t, err)
		container := proxy.Spec.Template.Spec.Containers[0].DeepCopy()
		_, err = overrideEnvVar(container, "KONG_STREAM_LISTEN", originalVal.Value)
		assert.NoError(t, err)

		// remove the added TCP port
		newPorts := make([]corev1.ContainerPort, 0, len(container.Ports)-1)
		for _, port := range container.Ports {
			if port.Name != dnsTestService {
				newPorts = append(newPorts, port)
			}
		}
		container.Ports = newPorts

		// revert to pre-test state
		proxy.Spec.Template.Spec.Containers[0] = *container
		_, err = cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Update(ctx, proxy, metav1.UpdateOptions{})
		assert.NoError(t, err)

		// ensure that the proxy deployment is ready before we proceed
		assert.Eventually(t, func() bool {
			d, err := cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Get(ctx, proxy.Name, metav1.GetOptions{})
			if err != nil {
				t.Logf("WARNING: error while waiting for deployment %s to become ready: %v", proxy, err)
				return false
			}
			if d.Status.ReadyReplicas == d.Status.Replicas && d.Status.AvailableReplicas == d.Status.Replicas && d.Status.UnavailableReplicas < 1 {
				return true
			}
			return false
		}, proxyUpdateWait, waitTick)
	}()

	// ensure that the proxy deployment is ready before we proceed
	assert.Eventually(t, func() bool {
		d, err := cluster.Client().AppsV1().Deployments(controllers.DefaultNamespace).Get(ctx, proxy.Name, metav1.GetOptions{})
		if err != nil {
			t.Logf("WARNING: error while waiting for deployment %s to become ready: %v", proxy, err)
			return false
		}
		if d.Status.ReadyReplicas == d.Status.Replicas && d.Status.AvailableReplicas == d.Status.Replicas && d.Status.UnavailableReplicas < 1 {
			return true
		}
		return false
	}, proxyUpdateWait, waitTick)

	// deploy a minimal deployment to test TCPIngress routes to
	deployment := k8sgen.NewDeploymentForContainer(k8sgen.NewContainer("tcp-test", "nginx", 80))
	_, err = cluster.Client().AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	// expose the deployment via service
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = cluster.Client().CoreV1().Services("default").Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	// make sure we clean up after ourselves
	defer func() {
		assert.NoError(t, cluster.Client().CoreV1().Services("default").Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	// initialize a clientset for the TCPIngress API
	c, err := clientset.NewForConfig(cluster.Config())
	assert.NoError(t, err)

	// deploy the TCPIngress object
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-tcp-test",
			Namespace: "default",
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "kong",
			},
		},
		Spec: kongv1beta1.TCPIngressSpec{
			Rules: []kongv1beta1.IngressRule{
				{
					Port: 32080,
					Backend: kongv1beta1.IngressBackend{
						ServiceName: service.Name,
						ServicePort: 80,
					},
				},
			},
		},
	}
	tcp, err = c.ConfigurationV1beta1().TCPIngresses("default").Create(ctx, tcp, metav1.CreateOptions{})
	assert.NoError(t, err)

	// make sure to cleanup
	defer func() {
		assert.NoError(t, c.ConfigurationV1beta1().TCPIngresses("default").Delete(ctx, tcp.Name, metav1.DeleteOptions{}))
	}()
}
