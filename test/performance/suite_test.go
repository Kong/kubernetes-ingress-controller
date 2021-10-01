//go:build performance_tests
// +build performance_tests

package performance

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	clusterDeployWait = time.Minute * 5
	waitTick          = time.Second * 1
	ingressWait       = time.Minute * 3
	httpcTimeout      = time.Second * 3
	httpBinImage      = "kennethreitz/httpbin"
	ingressClass      = "kong"
	max_ingress       = 500
)

var (
	httpc         = http.Client{Timeout: httpcTimeout}
	env           environments.Environment
	proxyURL      *url.URL
	proxyAdminURL *url.URL
	proxyUDPURL   *url.URL

	// maxBatchSize indicates the maximum number of objects that should be POSTed per second during testing
	maxBatchSize = determineMaxBatchSize()
)

func TestMain(m *testing.M) {
	var err error
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(clusterDeployWait))
	defer cancel()

	fmt.Println("INFO: configuring testing environment")
	kongAddon := kong.New()
	builder := environments.NewBuilder().WithAddons(metallb.New(), kongAddon)
	if existingClusterName := os.Getenv("KIND_CLUSTER"); existingClusterName != "" {
		fmt.Printf("INFO: using existing cluster %s\n", existingClusterName)
		cluster, err := kind.NewFromExisting(existingClusterName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not use existing cluster for test env: %s", err.Error())
			os.Exit(24)
		}
		builder = builder.WithExistingCluster(cluster)
	}

	fmt.Println("INFO: building test environment (note: can take some time)")
	env, err = builder.Build(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create testing environment: %s", err.Error())
		os.Exit(25)
	}
	fmt.Printf(
		"INFO: environment built CLUSTER_NAME=(%s) CLUSTER_TYPE=(%s) ADDONS=(metallb, kong)\n",
		env.Cluster().Name(), env.Cluster().Type(),
	)
	defer env.Cleanup(ctx)

	fmt.Printf("INFO: waiting for cluster %s and all addons to become ready\n", env.Cluster().Name())
	if err := <-env.WaitForReady(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: testing environment never became ready: %s", err.Error())
		os.Exit(26)
	}

	fmt.Printf("INFO: collecting Kong Proxy URLs from cluster %s for tests to make HTTP calls\n", env.Cluster().Name())
	proxyURL, err = kongAddon.ProxyURL(ctx, env.Cluster())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not get proxy URL from Kong Addon: %s", err.Error())
		os.Exit(27)
	}
	proxyAdminURL, err = kongAddon.ProxyAdminURL(ctx, env.Cluster())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not get proxy URL from Kong Addon: %s", err.Error())
		os.Exit(28)
	}
	proxyUDPURL, err = kongAddon.ProxyUDPURL(ctx, env.Cluster())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not get proxy URL from Kong Addon: %s", err.Error())
		os.Exit(29)
	}

	fmt.Println("INFO: testing environment is ready, running tests")
	code := m.Run()
	os.Exit(code)
}

// CreateNamespace create customized namespace
func CreateNamespace(ctx context.Context, namespace string, t *testing.T) error {
	assert.Eventually(t, func() bool {
		nsList, err := env.Cluster().Client().CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
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
			env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
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
	_, err := env.Cluster().Client().CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed creating namespace %s, err %v", namespace, err)
	}

	nsList, err := env.Cluster().Client().CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
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
	nsList, err := env.Cluster().Client().CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
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
		env.Cluster().Client().CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
		return nil
	}
	return nil
}

// determineMaxBatchSize provides a size limit for the number of resources to POST in a single second during tests, and can be overridden with an ENV var if desired.
func determineMaxBatchSize() int {
	if v := os.Getenv("KONG_BULK_TESTING_BATCH_SIZE"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("Error: invalid batch size %s: %s", v, err))
		}
		return i
	}
	return 50
}
