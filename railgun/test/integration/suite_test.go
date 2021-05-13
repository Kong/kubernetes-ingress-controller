//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kong/kubernetes-testing-framework/pkg/kind"
	ktfkind "github.com/kong/kubernetes-testing-framework/pkg/kind"

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/ctrlutils"
	"github.com/kong/kubernetes-ingress-controller/railgun/manager"
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
	ingressWait = time.Minute * 10

	// httpcTimeout is the default client timeout for HTTP clients used in tests.
	httpcTimeout = time.Second * 3

	// httpBinImage is the container image name we use for deploying the "httpbin" HTTP testing tool.
	// if you need a simple HTTP server for tests you're writing, use this and check the documentation.
	// See: https://github.com/postmanlabs/httpbin
	httpBinImage = "kennethreitz/httpbin"

	// ingressClass indicates the ingress class name which the tests will use for supported object reconcilation
	ingressClass = "kongtests"

	// elsewhere is the name of an alternative namespace
	elsewhere = "elsewhere"
)

// -----------------------------------------------------------------------------
// Testing Variables
// -----------------------------------------------------------------------------

var (
	// LegacyControllerEnvVar indicates the environment variable which can be used to trigger tests against the legacy KIC controller-manager
	LegacyControllerEnvVar = "KONG_LEGACY_CONTROLLER"

	// httpc is the default HTTP client to use for tests
	httpc = http.Client{Timeout: httpcTimeout}

	// cluster is the object which contains a Kubernetes client for the testing cluster
	cluster ktfkind.Cluster

	// watchNamespaces is a list of namespaces the controller watches
	watchNamespaces = strings.Join([]string{elsewhere, corev1.NamespaceDefault}, ",")
)

// -----------------------------------------------------------------------------
// Testing Main
// -----------------------------------------------------------------------------

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
	} else {
		// create a new cluster for tests
		config := ktfkind.ClusterConfigurationWithKongProxy{EnableMetalLB: true}
		if name := os.Getenv("KIND_CLUSTER_NAME"); name != "" {
			cluster, ready, err = config.DeployWithName(ctx, name)
		} else {
			cluster, ready, err = config.Deploy(ctx)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(12)
		}
		defer cluster.Cleanup()
	}

	// deploy the Kong Kubernetes Ingress Controller (KIC) to the cluster
	if err := deployControllers(ctx, ready, cluster, os.Getenv("KONG_CONTROLLER_TEST_IMAGE"), ctrlutils.DefaultNamespace); err != nil {
		cluster.Cleanup()
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(13)
	}

	code := m.Run()
	if !existingClusterInUse {
		cluster.Cleanup()
	}
	os.Exit(code)
}

// -----------------------------------------------------------------------------
// Testing Main - Controller Deployment
// -----------------------------------------------------------------------------

var crds = []string{
	"../../config/crd/bases/configuration.konghq.com_udpingresses.yaml",
	"../../config/crd/bases/configuration.konghq.com_tcpingresses.yaml",
	"../../config/crd/bases/configuration.konghq.com_kongplugins.yaml",
	"../../config/crd/bases/configuration.konghq.com_kongingresses.yaml",
	"../../config/crd/bases/configuration.konghq.com_kongconsumers.yaml",
	"../../config/crd/bases/configuration.konghq.com_kongclusterplugins.yaml",
}

// deployControllers ensures that relevant CRDs and controllers are deployed to the test cluster and supports legacy (KIC 1.x) clusters as well.
func deployControllers(ctx context.Context, ready chan ktfkind.ProxyReadinessEvent, cluster kind.Cluster, containerImage, namespace string) error {
	// ensure the controller namespace is created
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	if _, err := cluster.Client().CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}
	// ensure the alternative namespace is created
	elsewhereNS := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "elsewhere"}}
	if _, err := cluster.Client().CoreV1().Namespaces().Create(ctx, elsewhereNS, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	// run the controller in the background
	go func() {
		// pull the readiness event for the proxy
		event := <-ready

		// if there's an error, all tests fail here
		if event.Err != nil {
			panic(event.Err)
		}

		// grab the admin hostname and pass the readiness event on to the tests
		u := event.ProxyAdminURL
		adminHost := u.Hostname()
		proxyReadyCh <- event

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
		if useLegacyKIC() {
			cmd = buildLegacyCommand(ctx, kubeconfig.Name(), adminHost, cluster.Client())
			stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
			cmd.Stdout = io.MultiWriter(stdout, os.Stdout)
			cmd.Stderr = io.MultiWriter(stderr, os.Stderr)
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stdout, stdout.String())
				panic(fmt.Errorf("%s: %w", stderr.String(), err))
			}
		} else {
			config := manager.Config{}
			flags := config.FlagSet()
			flags.Parse([]string{
				fmt.Sprintf("--kong-admin-url=http://%s:8001", adminHost),
				fmt.Sprintf("--kubeconfig=%s", kubeconfig.Name()),
				"--controller-kongstate=enabled",
				"--controller-ingress-networkingv1=enabled",
				"--controller-ingress-networkingv1beta1=disabled",
				"--controller-ingress-extensionsv1beta1=disabled",
				"--controller-udpingress=enabled",
				"--controller-tcpingress=enabled",
				"--controller-kongingress=enabled",
				"--controller-kongclusterplugin=disabled",
				"--controller-kongplugin=disabled",
				"--controller-kongconsumer=disabled",
				"--election-id=integrationtests.konghq.com",
				fmt.Sprintf("--watch-namespace=%s", watchNamespaces),
				fmt.Sprintf("--ingress-class=%s", ingressClass),
				"--log-level=trace",
				"--log-format=text",
			})
			fmt.Printf("config: %+v\n", config)

			if err := manager.Run(ctx, &config); err != nil {
				panic(fmt.Errorf("controller manager exited with error: %w", err))
			}
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
		"--kong-admin-url", fmt.Sprintf("http://%s:8001", adminHost),
		"--ingress-class", ingressClass)

	// set the environment according to the legacy controller's needs
	cmd.Env = append(os.Environ(),
		"POD_NAMESPACE=kong-system",
		fmt.Sprintf("POD_NAME=%s", proxyPod),
	)

	return cmd
}

func waitForExistingClusterReadiness(ctx context.Context, cluster ktfkind.Cluster, name string, ready chan ktfkind.ProxyReadinessEvent) {
	var proxyAdminURL *url.URL
	var proxyURL *url.URL
	var proxyUDPUrl *url.URL

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "ERROR: timed out waiting for readiness from existing cluster %s", name)
			os.Exit(11)
		default:
			svcs, err := cluster.Client().CoreV1().Services(ctrlutils.DefaultNamespace).List(ctx, metav1.ListOptions{})
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
				ProxyUDPUrl:   proxyUDPUrl,
			}
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
}
