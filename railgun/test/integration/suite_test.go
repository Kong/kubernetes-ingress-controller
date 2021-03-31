//+build integration_tests

package integration

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	ktfkind "github.com/kong/kubernetes-testing-framework/pkg/kind"
)

var (
	// cluster is the object which contains a Kubernetes client for the testing cluster
	cluster ktfkind.Cluster

	// proxyReady is the channel that indicates when the Kong proxyReady is ready to use.
	proxyReady = make(chan *url.URL)
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a new cluster for tests
	config := ktfkind.ClusterConfigurationWithKongProxy{EnableMetalLB: true}
	newCluster, ready, err := config.Deploy(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(10)
	}
	defer newCluster.Cleanup()
	cluster = newCluster

	// deploy the Kong Kubernetes Ingress Controller (KIC) to the cluster
	if err := deployControllers(ctx, ready, cluster.Client(), os.Getenv("KONG_CONTROLLER_TEST_IMAGE"), "kic-under-test"); err != nil {
		newCluster.Cleanup()
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(11)
	}

	code := m.Run()
	newCluster.Cleanup()
	os.Exit(code)
}

// FIXME: this is a total hack for now, in the future we should deploy the controller into the cluster via image or run it as a goroutine.
func deployControllers(ctx context.Context, ready chan ktfkind.ProxyReadinessEvent, kc *kubernetes.Clientset, containerImage, namespace string) error {
	// ensure the controller namespace is created
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	if _, err := kc.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	// run the controller in the background
	go func() {
		event := <-ready
		if event.Err != nil {
			panic(event.Err)
		}
		u := event.URL
		proxyReady <- u

		cmd := exec.CommandContext(ctx, "go", "run", "../../main.go", "--kong-url", fmt.Sprintf("http://%s:8001", u.Hostname()))
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}()

	return nil
}
