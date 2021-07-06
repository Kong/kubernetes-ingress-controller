//+build performance_tests

package performance

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"testing"
	"time"

	ktfkind "github.com/kong/kubernetes-testing-framework/pkg/kind"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	controllerNamespace = "kong-system"
	clusterDeployWait   = time.Minute * 5
)

var cluster ktfkind.Cluster

func TestMain(m *testing.M) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(clusterDeployWait))
	defer cancel()

	var err error
	var existingClusterInUse bool
	ready := make(chan ktfkind.ProxyReadinessEvent)
	if existingClusterName := os.Getenv("KIND_CLUSTER"); existingClusterName != "" {
		existingClusterInUse = true
		cluster, err = ktfkind.GetExistingCluster(existingClusterName)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(10)
		}
		go waitForExistingClusterReadiness(ctx, cluster, existingClusterName, ready)
	}

	code := m.Run()
	if !existingClusterInUse {
		cluster.Cleanup()
	}
	os.Exit(code)
}

func waitForExistingClusterReadiness(ctx context.Context, cluster ktfkind.Cluster, name string, ready chan ktfkind.ProxyReadinessEvent) {
	var proxyAdminURL *url.URL
	var proxyURL *url.URL
	var proxyHTTPSURL *url.URL
	var proxyUDPUrl *url.URL
	var proxyIP *net.IP

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "ERROR: timed out waiting for readiness from existing cluster %s", name)
			os.Exit(11)
		default:
			svcs, err := cluster.Client().CoreV1().Services(controllerNamespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				ready <- ktfkind.ProxyReadinessEvent{Err: err}
				break
			}
			for _, svc := range svcs.Items {
				if svc.Name == "ingress-controller-kong-admin" && len(svc.Status.LoadBalancer.Ingress) == 1 {
					proxyAdminURL, err = url.Parse(fmt.Sprintf("http://%s:%d", svc.Status.LoadBalancer.Ingress[0].IP, 8001))
					if err != nil {
						ready <- ktfkind.ProxyReadinessEvent{Err: err}
						break
					}
				} else if svc.Name == "ingress-controller-kong-proxy" && len(svc.Status.LoadBalancer.Ingress) == 1 {
					proxyURL, err = url.Parse(fmt.Sprintf("http://%s:%d", svc.Status.LoadBalancer.Ingress[0].IP, 80))
					if err != nil {
						ready <- ktfkind.ProxyReadinessEvent{Err: err}
						break
					}
					proxyHTTPSURL, err = url.Parse(fmt.Sprintf("https://%s:%d", svc.Status.LoadBalancer.Ingress[0].IP, 443))
					if err != nil {
						ready <- ktfkind.ProxyReadinessEvent{Err: err}
						break
					}
					addr := net.ParseIP(svc.Status.LoadBalancer.Ingress[0].IP)
					proxyIP = &addr
				} else if svc.Name == "ingress-controller-kong-udp" && len(svc.Status.LoadBalancer.Ingress) == 1 {
					proxyUDPUrl, err = url.Parse(fmt.Sprintf("udp://%s:9999", svc.Status.LoadBalancer.Ingress[0].IP))
					if err != nil {
						ready <- ktfkind.ProxyReadinessEvent{Err: err}
						break
					}
				}
			}
		}
		if proxyAdminURL != nil && proxyURL != nil {
			ready <- ktfkind.ProxyReadinessEvent{
				ProxyAdminURL: proxyAdminURL,
				ProxyURL:      proxyURL,
				ProxyHTTPSURL: proxyHTTPSURL,
				ProxyIP:       proxyIP,
				ProxyUDPUrl:   proxyUDPUrl,
			}
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
}
