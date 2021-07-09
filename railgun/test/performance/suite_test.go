//+build performance_tests

package performance

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	ktfkind "github.com/kong/kubernetes-testing-framework/pkg/kind"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	controllerNamespace = "kong-system"
	clusterDeployWait   = time.Minute * 5
	max_ingress         = 500
	httpBinImage        = "kennethreitz/httpbin"
	ingressClass        = "kong"
	ingressWait         = time.Minute * 3
	waitTick            = time.Second * 1
	httpcTimeout        = time.Second * 3
)

var (
	cluster  ktfkind.Cluster
	KongInfo *ktfkind.ProxyReadinessEvent
	httpc    = http.Client{Timeout: httpcTimeout}
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(clusterDeployWait))
	defer cancel()

	var err error
	var existingClusterInUse bool
	if existingClusterName := os.Getenv("KIND_CLUSTER"); existingClusterName != "" {
		existingClusterInUse = true
		cluster, err = ktfkind.GetExistingCluster(existingClusterName)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(10)
		}
		KongInfo = &ktfkind.ProxyReadinessEvent{}
		svcs, err := cluster.Client().CoreV1().Services(controllerNamespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return
		}
		for _, svc := range svcs.Items {
			if svc.Name == "ingress-controller-kong-proxy" && len(svc.Status.LoadBalancer.Ingress) == 1 {
				proxyURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", svc.Status.LoadBalancer.Ingress[0].IP, 80))
				KongInfo.ProxyURL = proxyURL
			}
			if svc.Name == "ingress-controller-kong-admin" && len(svc.Status.LoadBalancer.Ingress) == 1 {
				proxyAdminURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", svc.Status.LoadBalancer.Ingress[0].IP, 8001))
				KongInfo.ProxyAdminURL = proxyAdminURL
			}
			if svc.Name == "ingress-controller-kong-udp" && len(svc.Status.LoadBalancer.Ingress) == 1 {
				proxyUDPUrl, _ := url.Parse(fmt.Sprintf("udp://%s:9999", svc.Status.LoadBalancer.Ingress[0].IP))
				KongInfo.ProxyUDPUrl = proxyUDPUrl
			}
		}
	}

	code := m.Run()
	if !existingClusterInUse {
		cluster.Cleanup()
	}
	os.Exit(code)
}

// CreateNamespace create customized namespace
func CreateNamespace(ctx context.Context, namespace string, t *testing.T) error {
	assert.Eventually(t, func() bool {
		nsList, err := cluster.Client().CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return false
		}

		needDelete := false
		needWait := false
		for _, item := range nsList.Items {
			if item.Name == namespace {
				if item.Status.Phase == corev1.NamespaceActive {
					t.Logf("namespace %s exists. removing it.", namespace)
					needDelete = true
					break
				}
				if item.Status.Phase == corev1.NamespaceTerminating {
					t.Logf("namespace is being terminating.")
					needWait = true
					break
				}
			}
		}

		if !needDelete && !needWait {
			t.Logf("namespace %s does not exist.", namespace)
			return true
		}

		if needDelete {
			cluster.Client().CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
			return false

		}

		return false
	}, 60*time.Second, 2*time.Second, true)

	nsName := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	t.Logf("creating namespace %s.", namespace)
	_, err := cluster.Client().CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed creating namespace %s, err %v", namespace, err)
	}

	nsList, err := cluster.Client().CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	assert.NoError(t, err)
	for _, item := range nsList.Items {
		if item.Name == namespace && item.Status.Phase == corev1.NamespaceActive {
			t.Logf("created namespace %s successfully.", namespace)
			return nil
		}
	}

	return fmt.Errorf("failed creating namespace %s", namespace)
}

func CleanUpNamespace(ctx context.Context, namespace string, t *testing.T) error {
	nsList, err := cluster.Client().CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	needDelete := false
	for _, item := range nsList.Items {
		if item.Name == namespace {
			if item.Status.Phase == corev1.NamespaceActive {
				t.Logf("namespace %s exists. removing it.", namespace)
				needDelete = true
				break
			}
		}
	}

	if needDelete {
		cluster.Client().CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
		return nil
	}
	return nil
}
