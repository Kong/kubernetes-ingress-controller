//go:build expression_router_tests

package expressionrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func exposeKongAdminService(ctx context.Context, t *testing.T,
	env environments.Environment, namespace string, name string,
) (string, int) {
	t.Helper()

	svcClient := env.Cluster().Client().CoreV1().Services(namespace)
	svc, err := svcClient.Get(ctx, name, metav1.GetOptions{})
	require.NoError(t, err)

	if svc.Spec.Type != corev1.ServiceTypeLoadBalancer {
		svc.Spec.Type = corev1.ServiceTypeLoadBalancer
		_, err = svcClient.Update(ctx, svc, metav1.UpdateOptions{})
		require.NoError(t, err)
	}

	// wait for IP to be present
	var ip string
	var port int
	require.Eventually(t, func() bool {
		svc, err = svcClient.Get(ctx, name, metav1.GetOptions{})
		require.NoError(t, err)
		if len(svc.Status.LoadBalancer.Ingress) == 0 {
			return false
		}

		ip = svc.Status.LoadBalancer.Ingress[0].IP
		for _, svcPort := range svc.Spec.Ports {
			if svcPort.Name == "kong-admin" {
				port = int(svcPort.Port)
			}
		}
		return true
	}, 2*time.Minute, 5*time.Second)

	return ip, port
}

// getKongProxyIP takes a Service with Kong proxy ports and returns and its IP, or fails the test if it cannot.
func getKongProxyIP(ctx context.Context, t *testing.T, env environments.Environment, namespace string) string {
	t.Helper()

	refreshService := func() *corev1.Service {
		svc, err := env.Cluster().Client().CoreV1().Services(namespace).Get(ctx, "ingress-controller-kong-proxy", metav1.GetOptions{})
		require.NoError(t, err)
		return svc
	}

	svc := refreshService()
	require.NotEqual(t, svc.Spec.Type, corev1.ServiceTypeClusterIP, "ClusterIP service is not supported")

	//nolint: exhaustive
	switch svc.Spec.Type {
	case corev1.ServiceTypeLoadBalancer:
		return getKongProxyLoadBalancerIP(t, refreshService)
	default:
		t.Fatalf("unknown service type: %q", svc.Spec.Type)
		return ""
	}
}

func getKongProxyLoadBalancerIP(t *testing.T, refreshSvc func() *corev1.Service) string {
	t.Helper()

	var resIP string
	require.Eventually(t, func() bool {
		svc := refreshSvc()

		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			ip := svc.Status.LoadBalancer.Ingress[0].IP
			t.Logf("found loadbalancer IP for the Kong Proxy: %s", ip)
			resIP = ip
			return true
		}
		return false
	}, 2*time.Minute, time.Second)

	return resIP
}

func marshalKongConfig(t *testing.T, s kong.Service, r kong.Route) io.Reader {
	t.Helper()

	content := &file.Content{
		FormatVersion: "3.0",
		Services: []file.FService{
			{
				Service: s,
				Routes: []*file.FRoute{
					{
						Route: r,
					},
				},
			},
		},
	}
	config, err := json.Marshal(content)
	require.NoError(t, err)

	return bytes.NewReader(config)
}
