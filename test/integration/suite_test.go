//+build integration_tests

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/internal/cmd/rootcmd"
	"github.com/kong/kubernetes-ingress-controller/internal/manager"
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

	// ingressClass indicates the ingress class name which the tests will use for supported object reconciliation
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
	// httpc is the default HTTP client to use for tests
	httpc = http.Client{Timeout: httpcTimeout}

	// watchNamespaces is a list of namespaces the controller watches
	watchNamespaces = strings.Join([]string{
		elsewhere,
		corev1.NamespaceDefault,
		testIngressEssentialsNamespace,
		testIngressClassNameSpecNamespace,
		testIngressHTTPSNamespace,
		testIngressHTTPSRedirectNamespace,
		testBulkIngressNamespace,
		testTCPIngressNamespace,
		testUDPIngressNamespace,
		testPluginsNamespace,
	}, ",")

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
// Testing Variables - Environment Overrides
// -----------------------------------------------------------------------------

var (
	// dbmode indicates the database backend of the test cluster ("off" and "postgres" are supported)
	dbmode = os.Getenv("TEST_DATABASE_MODE")

	// clusterVersion indicates the version of Kubernetes to use for the tests (if the cluster was not provided by the caller)
	clusterVersionStr = os.Getenv("KONG_CLUSTER_VERSION")

	// existingClusterName indicates whether or not the caller is providing their own kind cluster for running the tests,
	// and if so what the name of that cluster is.
	existingClusterName = os.Getenv("KIND_CLUSTER")

	// maxBatchSize indicates the maximum number of objects that should be POSTed per second during testing
	maxBatchSize = determineMaxBatchSize()
)

// -----------------------------------------------------------------------------
// Testing Main
// -----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("INFO: configuring testing environment")
	kongbuilder := kong.NewBuilder()
	if dbmode == "postgres" {
		kongbuilder = kongbuilder.WithPostgreSQL()
	}
	kongbuilder.WithControllerDisabled()
	kongAddon := kongbuilder.Build()
	builder := environments.NewBuilder().WithAddons(metallb.New(), kongAddon)

	fmt.Println("INFO: checking for reusable environment components")
	if existingClusterName != "" {
		if clusterVersionStr != "" {
			fmt.Fprintf(os.Stderr, "Error: can't flag cluster version and provide an existing cluster at the same time")
			os.Exit(ExitCodeIncompatibleOptions)
		}
		fmt.Printf("INFO: using existing cluster %s\n", existingClusterName)
		cluster, err := kind.NewFromExisting(existingClusterName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not use existing cluster for test env: %s", err)
			os.Exit(ExitCodeCantUseExistingCluster)
		}
		builder = builder.WithExistingCluster(cluster)
	}

	fmt.Println("INFO: configuring kubernetes cluster")
	if clusterVersionStr != "" {
		clusterVersion, err := semver.Parse(strings.TrimPrefix(clusterVersionStr, "v"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid cluster version provided (%s): %s", clusterVersionStr, err)
			os.Exit(ExitCodeInvalidOptions)
		}
		cluster, err := kind.NewBuilder().WithClusterVersion(clusterVersion).Build(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to build kind cluster with version %s: %s", clusterVersion, err)
			os.Exit(ExitCodeCantCreateCluster)
		}
		builder.WithExistingCluster(cluster)
	}

	fmt.Println("INFO: building test environment (note: can take some time)")
	env, err = builder.Build(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create testing environment: %s", err)
		os.Exit(ExitCodeCantCreateCluster)
	}
	fmt.Printf(
		"INFO: environment built CLUSTER_NAME=(%s) CLUSTER_TYPE=(%s) ADDONS=(metallb, kong)\n",
		env.Cluster().Name(), env.Cluster().Type(),
	)
	defer func() {
		if err := env.Cleanup(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not cleanup testing environment: %s", err)
			os.Exit(ExitCodeCleanupFailed)
		}
	}()

	fmt.Printf("INFO: waiting for cluster %s and all addons to become ready\n", env.Cluster().Name())
	if err := <-env.WaitForReady(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: testing environment never became ready: %s", err)
		os.Exit(ExitCodeCantCreateCluster)
	}

	fmt.Printf("INFO: collecting Kong Proxy URLs from cluster %s for tests to make HTTP calls\n", env.Cluster().Name())
	proxyURL, err = kongAddon.ProxyURL(ctx, env.Cluster())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not get proxy URL from Kong Addon: %s", err)
		os.Exit(ExitCodeCantCreateCluster)
	}
	proxyAdminURL, err = kongAddon.ProxyAdminURL(ctx, env.Cluster())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not get proxy URL from Kong Addon: %s", err)
		os.Exit(ExitCodeCantCreateCluster)
	}
	proxyUDPURL, err = kongAddon.ProxyUDPURL(ctx, env.Cluster())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not get proxy URL from Kong Addon: %s", err)
		os.Exit(ExitCodeCantCreateCluster)
	}

	if v := os.Getenv("KONG_BRING_MY_OWN_KIC"); v == "true" {
		fmt.Println("WARNING: caller indicated that they will manage their own controller")
	} else {
		fmt.Println("INFO: deploying controller manager")
		if err := deployControllers(ctx, controllerNamespace); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(ExitCodeCantCreateCluster)
		}
	}

	fmt.Printf("INFO: running final testing environment checks")
	serverVersion, err := env.Cluster().Client().ServerVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not retrieve server version for cluster: %s", err)
		os.Exit(ExitCodeCantCreateCluster)
	}

	fmt.Printf("INFO: testing environment is ready KUBERNETES_VERSION=(%v): running tests\n", serverVersion)
	code := m.Run()
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
	knativeCrds,
}

// deployControllers ensures that relevant CRDs and controllers are deployed to the test cluster
func deployControllers(ctx context.Context, namespace string) error {
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
		generateKubeconfig := exec.CommandContext(ctx, "kind", "get", "kubeconfig", "--name", env.Cluster().Name()) //nolint:gosec
		generateKubeconfig.Stdout = kubeconfig
		generateKubeconfig.Stderr = os.Stderr
		if err := generateKubeconfig.Run(); err != nil {
			panic(err)
		}
		kubeconfig.Close()

		// deploy our CRDs to the cluster
		for _, crd := range crds {
			cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig.Name(), "apply", "-f", crd) //nolint:gosec
			stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
			cmd.Stdout = stdout
			cmd.Stderr = stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stdout, stdout.String())
				panic(fmt.Errorf("%s: %w", stderr.String(), err))
			}
		}

		config := manager.Config{}
		flags := config.FlagSet()
		if err := flags.Parse([]string{
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
			"--dump-config",
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
		}); err != nil {
			panic(fmt.Errorf("could not parse controller manager flags: %w", err))
		}
		fmt.Fprintf(os.Stderr, "config: %+v\n", config)

		if err := rootcmd.Run(ctx, &config); err != nil {
			panic(fmt.Errorf("controller manager exited with error: %w", err))
		}
	}()

	return nil
}
