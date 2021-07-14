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

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/cmd/rootcmd"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/manager"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
)

// -----------------------------------------------------------------------------
// Testing Timeouts
// -----------------------------------------------------------------------------

const (
	// waitTick is the default timeout tick interval for checking on ingress resources.
	waitTick = time.Second * 1

	// ingressWait is the default amount of time to wait for any particular ingress resource to be provisioned.
	ingressWait = time.Minute * 3

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

	// controllerNamespace is the Kubernetes namespace where the controller is deployed
	controllerNamespace = "kong-system"
)

// -----------------------------------------------------------------------------
// Testing Variables
// -----------------------------------------------------------------------------

var (
	// LegacyControllerEnvVar indicates the environment variable which can be used to trigger tests against the legacy KIC controller-manager
	LegacyControllerEnvVar = "KONG_LEGACY_CONTROLLER"

	// httpc is the default HTTP client to use for tests
	httpc = http.Client{Timeout: httpcTimeout}

	// watchNamespaces is a list of namespaces the controller watches
	watchNamespaces = strings.Join([]string{
		elsewhere,
		corev1.NamespaceDefault,
		testIngressNamespace,
		testTCPIngressNamespace,
		testUDPIngressNamespace,
		testPluginsNamespace,
	}, ",")

	// dbmode indicates the database backend of the test cluster ("off" and "postgres" are supported)
	dbmode = os.Getenv("TEST_DATABASE_MODE")

	// env is the primary testing environment object which includes access to the Kubernetes cluster
	// and all the addons deployed in support of the tests.
	env environments.Environment

	// proxyURL provides access to the proxy endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyURL *url.URL

	// proxyAdminURL provides access to the Admin API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyAdminURL *url.URL

	// proxyUDPURL provides access to the UDP API endpoint for the Kong Addon which is deployed to the test environment's cluster.
	proxyUDPURL *url.URL
)

// -----------------------------------------------------------------------------
// Testing Main
// -----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
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

	if v := os.Getenv("KONG_BRING_MY_OWN_KIC"); v == "true" {
		fmt.Println("WARNING: caller indicated that they will manage their own controller")
	} else {
		fmt.Println("INFO: deploying controller manager")
		if err := deployControllers(ctx, os.Getenv("KONG_CONTROLLER_TEST_IMAGE"), controllerNamespace); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(13)
		}
	}

	fmt.Println("INFO: testing environment is ready, running tests")
	code := m.Run()
	os.Exit(code)
}

// -----------------------------------------------------------------------------
// Testing Helper Functions
// -----------------------------------------------------------------------------

func useLegacyKIC() bool {
	return os.Getenv(LegacyControllerEnvVar) != ""
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
	knativeCrds,
}

// deployControllers ensures that relevant CRDs and controllers are deployed to the test cluster and supports legacy (KIC 1.x) clusters as well.
func deployControllers(ctx context.Context, containerImage, namespace string) error {
	// ensure the controller namespace is created
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	if _, err := env.Cluster().Client().CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	// run the controller in the background
	go func() {
		// create a tempfile to hold the cluster kubeconfig that will be used for the controller
		kubeconfig, err := ioutil.TempFile(os.TempDir(), "kubeconfig-")
		if err != nil {
			panic(err)
		}
		defer os.Remove(kubeconfig.Name())

		// dump the kubeconfig from kind into the tempfile
		generateKubeconfig := exec.CommandContext(ctx, "kind", "get", "kubeconfig", "--name", env.Cluster().Name())
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
			cmd = buildLegacyCommand(ctx, kubeconfig.Name())
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
				fmt.Sprintf("--kong-admin-url=http://%s:8001", proxyAdminURL.Hostname()),
				fmt.Sprintf("--kubeconfig=%s", kubeconfig.Name()),
				"--controller-kongstate=enabled",
				"--controller-ingress-networkingv1=enabled",
				"--controller-ingress-networkingv1beta1=enabled",
				"--controller-ingress-extensionsv1beta1=enabled",
				"--controller-tcpingress=enabled",
				"--controller-kongingress=enabled",
				"--controller-knativeingress=enabled",
				"--controller-kongclusterplugin=enabled",
				"--controller-kongplugin=enabled",
				"--controller-kongconsumer=disabled",
				"--election-id=integrationtests.konghq.com",
				"--publish-service=kong-system/ingress-controller-kong-proxy",
				fmt.Sprintf("--watch-namespace=%s", watchNamespaces),
				fmt.Sprintf("--ingress-class=%s", ingressClass),
				"--log-level=trace",
				"--log-format=text",
				"--debug-log-reduce-redundancy",
				"--admission-webhook-listen=172.17.0.1:49023",
				fmt.Sprintf("--admission-webhook-cert=%s", admissionWebhookCert),
				fmt.Sprintf("--admission-webhook-key=%s", admissionWebhookKey),
				"--profiling",
			})
			fmt.Fprintf(os.Stderr, "config: %+v\n", config)

			if err := rootcmd.Run(ctx, &config); err != nil {
				panic(fmt.Errorf("controller manager exited with error: %w", err))
			}
		}
	}()

	return nil
}

func buildLegacyCommand(ctx context.Context, kubeconfigPath string) *exec.Cmd {
	fmt.Fprintln(os.Stdout, "WARNING: deploying legacy Kong Kubernetes Ingress Controller (KIC)")

	// get the proxy pod
	podList, err := env.Cluster().Client().CoreV1().Pods("kong-system").List(ctx, metav1.ListOptions{
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
		"--kong-admin-url", fmt.Sprintf("http://%s:8001", proxyAdminURL.Hostname()),
		"--ingress-class", ingressClass)

	// set the environment according to the legacy controller's needs
	cmd.Env = append(os.Environ(),
		"POD_NAMESPACE=kong-system",
		fmt.Sprintf("POD_NAME=%s", proxyPod),
	)

	return cmd
}
