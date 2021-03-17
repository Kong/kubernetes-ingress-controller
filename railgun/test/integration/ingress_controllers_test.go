//+build integration_tests

package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sgen "github.com/kong/kubernetes-testing-framework/pkg/generators/k8s"
)

var (
	// ingressTimeout is the maximum amount of time that the tests should wait for an Ingress record to be provisioned and the backend accessible.
	ingressTimeout = time.Minute * 5

	// ingressTimeoutTick is the time to wait between Ingress resource timeout checks
	ingressTimeoutTick = time.Second * 1
)

func TestMinimalIngress(t *testing.T) {
	ctx := context.Background()

	// deploy a minimal deployment to test Ingress routes to
	container := k8sgen.NewContainer("nginx", "nginx", 80)
	deployment := k8sgen.NewDeploymentForContainer(container)
	_, err := cluster.Client().AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	assert.NoError(t, err)

	// expose the deployment via service
	service := k8sgen.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	_, err = cluster.Client().CoreV1().Services("default").Create(ctx, service, metav1.CreateOptions{})
	assert.NoError(t, err)

	// route to the service via Kong Ingress
	ingress := k8sgen.NewIngressForService("/nginx", map[string]string{
		"kubernetes.io/ingress.class": "kong",
		"konghq.com/strip-path":       "true",
	}, service)
	cluster.Client().NetworkingV1().Ingresses("default").Create(ctx, ingress, metav1.CreateOptions{})

	// wait for the ingress backend to be routable
	u := <-proxyURL
	assert.Eventually(t, func() bool {
		resp, err := http.Get(fmt.Sprintf("%s/nginx", u.String()))
		if err != nil {
			t.Logf("WARNING: error while waiting for %s: %v", u.String(), err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, ingressTimeout, ingressTimeoutTick)
}
