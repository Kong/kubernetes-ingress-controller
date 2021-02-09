package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"

	ktfkind "github.com/kong/kubernetes-testing-framework/pkg/kind"
	ktfkong "github.com/kong/kubernetes-testing-framework/pkg/kong"
	ktfmetal "github.com/kong/kubernetes-testing-framework/pkg/metallb"

	"github.com/kong/railgun/controllers"
)

var (
	// ClusterName indicates the name of the Kind test cluster setup for this test suite.
	ClusterName = uuid.New().String()

	// kc is a kubernetes clientset for the default Kind cluster created for this test suite.
	kc *kubernetes.Clientset
)

func TestMain(m *testing.M) {
	// setup a kind cluster for testing
	err := ktfkind.CreateKindCluster(ClusterName)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(10)
	}

	// cleanup the kind cluster when we're done, unless flagged otherwise
	defer func() {
		if v := os.Getenv("KIND_KEEP_CLUSTER"); v == "" { // you can optionally flag the tests to retain the test cluster for inspection.
			if err := ktfkind.DeleteKindCluster(ClusterName); err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(11)
			}
		}
	}()

	// setup Metallb for the cluster for LoadBalancer addresses for Kong
	if err := ktfmetal.DeployMetallbForKindCluster(ClusterName, ktfkind.DefaultKindDockerNetwork); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(12)
	}

	// retrieve the *kubernetes.Clientset for the cluster
	kc, err = ktfkind.ClientForKindCluster(ClusterName)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(13)
	}

	// setup a Kong proxy to test against
	if err := ktfkong.SimpleProxySetup(kc, controllers.DefaultNamespace); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(14)
	}

	// deploy the Kong Kubernetes Ingress Controller (KIC) to the cluster
	cancel, err := ktfkong.DeployControllers(kc, os.Getenv("KONG_CONTROLLER_TEST_IMAGE"), controllers.DefaultNamespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(15)
	}
	defer cancel()

	// run the tests
	code := m.Run()
	os.Exit(code)
}
