//+build integration_tests

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

func TestMinimalTCPIngress(t *testing.T) {
	// test setup
	namespace := "default"
	testName := "mintcp"
	ctx, cancel := context.WithTimeout(context.Background(), ingressWait)
	defer cancel()

	// push a minimal deployment to test TCPIngress routes to
	deployment := k8sgen.NewDeploymentForContainer(k8sgen.NewContainer(testName, "nginx", 80))
	_, err := cluster.Client().AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	// make sure we clean up after ourselves
	defer func() {
		assert.NoError(t, cluster.Client().AppsV1().Deployments(namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{}))
	}()

	// expose the deployment via service
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = cluster.Client().CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	// make sure we clean up after ourselves
	defer func() {
		assert.NoError(t, cluster.Client().CoreV1().Services(namespace).Delete(ctx, service.Name, metav1.DeleteOptions{}))
	}()

	// initialize a clientset for the TCPIngress API
	c, err := clientset.NewForConfig(cluster.Config())
	assert.NoError(t, err)

	// deploy the TCPIngress object
	tcp := &kongv1beta1.TCPIngress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "kong",
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
	assert.NoError(t, err)

	// make sure we clean up after ourselves
	defer func() {
		assert.NoError(t, c.ConfigurationV1beta1().TCPIngresses(namespace).Delete(ctx, tcp.Name, metav1.DeleteOptions{}))
	}()
}
