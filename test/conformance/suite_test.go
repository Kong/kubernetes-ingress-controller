//go:build conformance_tests

package conformance

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/kong"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/types/kind"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
	testutils "github.com/kong/kubernetes-ingress-controller/v2/internal/util/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test"
	"github.com/kong/kubernetes-ingress-controller/v2/test/consts"
	"github.com/kong/kubernetes-ingress-controller/v2/test/helpers/certificate"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/helpers"
	"github.com/kong/kubernetes-ingress-controller/v2/test/internal/testenv"
)

var (
	conformanceTestsBaseManifests = fmt.Sprintf("%s/conformance/base/manifests.yaml", consts.GatewayRawRepoURL)
	ingressClass                  = "kong-conformance-tests"

	env                    environments.Environment
	ctx                    context.Context
	globalDeprecatedLogger logrus.FieldLogger
	globalLogger           logr.Logger
)

func TestMain(m *testing.M) {
	var code int
	defer func() {
		os.Exit(code)
	}()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Logger needs to be configured before anything else happens.
	// This is because the controller manager has a timeout for
	// logger initialization, and if the logger isn't configured
	// after 30s from the start of controller manager package init function,
	// the controller manager will set up a no op logger and continue.
	// The logger cannot be configured after that point.
	deprecatedLogger, logger, logOutput, err := testutils.SetupLoggers("trace", "text", false)
	if err != nil {
		exitOnErr(fmt.Errorf("failed to setup loggers: %w", err))
	}
	if logOutput != "" {
		fmt.Printf("INFO: writing manager logs to %s\n", logOutput)
	}
	globalDeprecatedLogger = deprecatedLogger
	globalLogger = logger

	// In order to pass conformance tests, the expression router is required.
	kongBuilder := kong.NewBuilder().WithControllerDisabled().WithProxyAdminServiceTypeLoadBalancer().
		WithNamespace(consts.ControllerNamespace)
	if testenv.ExpressionRoutesEnabled() {
		fmt.Println("INFO: expression routes enabled")
		kongBuilder = kongBuilder.WithProxyEnvVar("router_flavor", "expressions")
	}

	// Pin the Helm chart version.
	kongBuilder.WithHelmChartVersion(consts.KongHelmChartVersion)

	kongAddon := kongBuilder.Build()
	builder := environments.NewBuilder().WithAddons(metallb.New(), kongAddon)
	useExistingClusterIfPresent(builder)

	env, err = builder.Build(ctx)
	exitOnErr(err)

	fmt.Println("INFO: waiting for cluster and addons to be ready")
	envReadyCtx, envReadyCancel := context.WithTimeout(ctx, testenv.EnvironmentReadyTimeout())
	defer envReadyCancel()
	exitOnErr(<-env.WaitForReady(envReadyCtx))

	// To allow running conformance tests in a loop to e.g. detect flaky tests
	// let's ensure that conformance related namespaced are deleted from the cluster.
	exitOnErr(ensureConformanceTestsNamespacesAreNotPresent(ctx, env.Cluster().Client()))

	code = m.Run()
	if testenv.IsCI() {
		fmt.Printf("INFO: running in ephemeral CI environment, skipping cluster %s teardown\n", env.Cluster().Name())
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), test.EnvironmentCleanupTimeout)
		defer cancel()
		exitOnErr(helpers.RemoveCluster(ctx, env.Cluster()))
	}
}

// prepareEnvForGatewayConformanceTests prepares the environment for running the gateway conformance test suite
// from the gateway-api project. Used as a helper function for both the stable and experimental conformance tests.
func prepareEnvForGatewayConformanceTests(t *testing.T) (c client.Client, gatewayClassName string) {
	t.Log("configuring environment for gateway conformance tests")
	client, err := client.New(env.Cluster().Config(), client.Options{})
	require.NoError(t, err)
	require.NoError(t, gatewayv1alpha2.AddToScheme(client.Scheme()))
	require.NoError(t, gatewayv1.AddToScheme(client.Scheme()))

	featureGateFlag := fmt.Sprintf("--feature-gates=%s", consts.DefaultFeatureGates)
	if testenv.ExpressionRoutesEnabled() {
		featureGateFlag = fmt.Sprintf("--feature-gates=%s", consts.ConformanceExpressionRoutesTestsFeatureGates)
	}

	t.Log("starting the controller manager")
	cert, key := certificate.GetKongSystemSelfSignedCerts()
	args := []string{
		fmt.Sprintf("--ingress-class=%s", ingressClass),
		fmt.Sprintf("--admission-webhook-cert=%s", cert),
		fmt.Sprintf("--admission-webhook-key=%s", key),
		fmt.Sprintf("--admission-webhook-listen=%s:%d", testutils.AdmissionWebhookListenHost, testutils.AdmissionWebhookListenPort),
		"--profiling",
		"--dump-config",
		"--log-level=trace",
		"--debug-log-reduce-redundancy",
		featureGateFlag,
		"--anonymous-reports=false",
	}

	require.NoError(t, testutils.DeployControllerManagerForCluster(ctx, globalDeprecatedLogger, globalLogger, env.Cluster(), args...))

	t.Log("creating GatewayClass for gateway conformance tests")
	gatewayClass := &gatewayv1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
			Annotations: map[string]string{
				annotations.GatewayClassUnmanagedAnnotation: annotations.GatewayClassUnmanagedAnnotationValuePlaceholder,
			},
		},
		Spec: gatewayv1.GatewayClassSpec{
			ControllerName: gateway.GetControllerName(),
		},
	}
	require.NoError(t, client.Create(ctx, gatewayClass))
	t.Cleanup(func() { require.NoError(t, client.Delete(ctx, gatewayClass)) })
	t.Logf("created GatewayClass %q", gatewayClass.Name)

	return client, gatewayClass.Name
}

func ensureConformanceTestsNamespacesAreNotPresent(ctx context.Context, client *kubernetes.Clientset) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(3) //nolint:gomnd

	for _, namespace := range []string{"gateway-conformance-infra", "gateway-conformance-web-backend", "gateway-conformance-app-backend"} {
		namespace := namespace
		g.Go(func() error {
			return ensureNamespaceDeleted(ctx, namespace, client)
		})
	}
	return g.Wait()
}

func ensureNamespaceDeleted(ctx context.Context, ns string, client *kubernetes.Clientset) error {
	namespaceClient := client.CoreV1().Namespaces()
	namespace, err := namespaceClient.Get(ctx, ns, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the namespace cannot be found then we're good to go.
			return nil
		}
		return err
	}

	if namespace.Status.Phase == corev1.NamespaceActive {
		fmt.Printf("INFO: removing %s namespace for clean test run\n", ns)
		err := namespaceClient.Delete(ctx, ns, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	w, err := namespaceClient.Watch(ctx, metav1.ListOptions{
		LabelSelector: "kubernetes.io/metadata.name=" + ns,
	})
	if err != nil {
		return err
	}

	defer w.Stop()
	for {
		select {
		case event := <-w.ResultChan():
			if event.Type == watch.Deleted {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func useExistingClusterIfPresent(builder *environments.Builder) {
	if existingCluster := testenv.ExistingClusterName(); existingCluster != "" {
		parts := strings.Split(existingCluster, ":")
		if len(parts) != 2 {
			exitOnErr(fmt.Errorf("%s is not a valid value for KONG_TEST_CLUSTER", existingCluster))
		}
		if parts[0] != "kind" {
			exitOnErr(fmt.Errorf("%s is not a supported cluster type for this test suite yet", parts[0]))
		}
		cluster, err := kind.NewFromExisting(parts[1])
		exitOnErr(err)
		fmt.Printf("INFO: using existing kind cluster for test (name: %s)\n", parts[1])
		builder.WithExistingCluster(cluster)
	} else {
		fmt.Println("INFO: creating new kind cluster for conformance tests")
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func shouldRunExperimentalConformance() bool {
	return os.Getenv("TEST_EXPERIMENTAL_CONFORMANCE") == "true"
}
