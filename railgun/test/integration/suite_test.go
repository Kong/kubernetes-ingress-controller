//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kong/kubernetes-testing-framework/pkg/kind"
	ktfkind "github.com/kong/kubernetes-testing-framework/pkg/kind"

	"github.com/kong/kubernetes-ingress-controller/railgun/controllers"
)

// -----------------------------------------------------------------------------
// Testing Timeouts
// -----------------------------------------------------------------------------

const (
	// clusterDeployWait is the timeout duration for deploying the kind cluster for testing
	clusterDeployWait = time.Minute * 5

	// waitTick is the default timeout tick interval for checking on ingress resources.
	waitTick = time.Second * 1

	// ingressWait is the default amount of time to wait for any particular ingress resource to be provisioned.
	ingressWait = time.Minute * 7
)

// -----------------------------------------------------------------------------
// Testing Variables
// -----------------------------------------------------------------------------

var (
	// LegacyControllerEnvVar indicates the environment variable which can be used to trigger tests against the legacy KIC controller-manager
	LegacyControllerEnvVar = "KONG_LEGACY_CONTROLLER"

	// cluster is the object which contains a Kubernetes client for the testing cluster
	cluster ktfkind.Cluster

	// proxyReady is the channel that indicates when the Kong proxyReady is ready to use.
	proxyReady = make(chan *url.URL)
)

// -----------------------------------------------------------------------------
// Testing Main
// -----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(clusterDeployWait))
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
	if err := deployControllers(ctx, ready, cluster, os.Getenv("KONG_CONTROLLER_TEST_IMAGE"), controllers.DefaultNamespace); err != nil {
		newCluster.Cleanup()
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(11)
	}

	code := m.Run()
	newCluster.Cleanup()
	os.Exit(code)
}

// -----------------------------------------------------------------------------
// Testing Main - Controller Deployment
// -----------------------------------------------------------------------------

var crds = []string{
	"../../config/crd/bases/configuration.konghq.com_udpingresses.yaml",
	"../../config/crd/bases/configuration.konghq.com_tcpingresses.yaml",
}

// FIXME: this is a total hack for now, in the future we should deploy the controller into the cluster via image or run it as a goroutine.
func deployControllers(ctx context.Context, ready chan ktfkind.ProxyReadinessEvent, cluster kind.Cluster, containerImage, namespace string) error {
	// ensure the controller namespace is created
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	if _, err := cluster.Client().CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
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
		u := event.ProxyAdminURL
		adminHost := u.Hostname()
		proxyReady <- u

		// create a tempfile to hold the cluster kubeconfig that will be used for the controller
		kubeconfig, err := ioutil.TempFile(os.TempDir(), "kubeconfig-")
		if err != nil {
			panic(err)
		}
		defer os.Remove(kubeconfig.Name())

		// dump the kubeconfig from kind into the tempfile
		generateKubeconfig := exec.CommandContext(ctx, "kind", "get", "kubeconfig", "--name", cluster.Name())
		generateKubeconfig.Stdout = kubeconfig
		generateKubeconfig.Stderr = os.Stderr
		if err := generateKubeconfig.Run(); err != nil {
			panic(err)
		}
		kubeconfig.Close()

		// deploy our CRDs to the cluster
		for _, crd := range crds {
			cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig.Name(), "apply", "-f", crd)
			stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
			cmd.Stdout = stdout
			cmd.Stderr = stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stdout, stdout.String())
				panic(fmt.Errorf("%s: %w", stderr.String(), err))
			}
		}

		// if set, allow running the legacy controller for the tests instead of the current controller
		var cmd *exec.Cmd
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
		if useLegacyKIC() {
			cmd = buildLegacyCommand(ctx, kubeconfig.Name(), adminHost, cluster.Client())
		} else {
			cmd = buildControllerCommand(ctx, kubeconfig.Name(), adminHost)
		}
		stdout, stderr = new(bytes.Buffer), new(bytes.Buffer)
		cmd.Stdout = io.MultiWriter(stdout, os.Stdout)
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stdout, stdout.String())
			panic(fmt.Errorf("%s: %w", stderr.String(), err))
		}
	}()

	return nil
}

func useLegacyKIC() bool {
	return os.Getenv(LegacyControllerEnvVar) != ""
}

// TODO: this will be removed as part of KIC 2.0, where the legacy controller will be replaced.
//       for more details see the relevant milestone: https://github.com/Kong/kubernetes-ingress-controller/milestone/12
func buildLegacyCommand(ctx context.Context, kubeconfigPath, adminHost string, kc *kubernetes.Clientset) *exec.Cmd {
	fmt.Fprintln(os.Stdout, "WARNING: deploying legacy Kong Kubernetes Ingress Controller (KIC)")

	// get the proxy pod
	podList, err := kc.CoreV1().Pods("kong-system").List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/component=app,app.kubernetes.io/instance=ingress-controller,app.kubernetes.io/name=kong",
	})
	if err != nil {
		panic(err)
	}
	if len(podList.Items) != 1 {
		panic(fmt.Errorf("expected 1 result, found %d", len(podList.Items)))
	}
	proxyPod := podList.Items[0].Name

	// custom command for the legacy controller as there are several differences in flags.
	cmd := exec.CommandContext(ctx, "go", "run", "../../../cli/ingress-controller/",
		"--publish-service", "kong-system/ingress-controller-kong-proxy",
		"--kubeconfig", kubeconfigPath,
		"--kong-admin-url", fmt.Sprintf("http://%s:8001", adminHost))

	// set the environment according to the legacy controller's needs
	cmd.Env = append(os.Environ(),
		"POD_NAMESPACE=kong-system",
		fmt.Sprintf("POD_NAME=%s", proxyPod),
	)

	return cmd
}

func buildControllerCommand(ctx context.Context, kubeconfigPath, adminHost string) *exec.Cmd {
	return exec.CommandContext(ctx, "go", "run", "../../main.go",
		"--kong-url", fmt.Sprintf("http://%s:8001", adminHost),
		"--kubeconfig", kubeconfigPath)
}
